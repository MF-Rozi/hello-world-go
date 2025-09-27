package server

//TODO: - make Chi Server and Routes
import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"dev.mfr/next-js-backend/internal/ipgeolocation"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Next.js Backend!")
	})
	r.Get("/ip", clientIP)
	return r
}

func StartServer() {
	r := routes()

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func clientIP(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)

	location, err := ipgeolocation.GetGeoLocation(ip)
	if err != nil {
		http.Error(w, "Failed to get geolocation", http.StatusInternalServerError)
		return
	}

	// If connecting locally/private and no location from ipinfo, set default coords
	if func(ip string) bool {
		p := net.ParseIP(ip)
		if p == nil {
			return false
		}
		if p.IsLoopback() {
			return true
		}
		if v4 := p.To4(); v4 != nil {
			if v4[0] == 10 {
				return true
			}
			if v4[0] == 172 && v4[1] >= 16 && v4[1] <= 31 {
				return true
			}
			if v4[0] == 192 && v4[1] == 168 {
				return true
			}
		}
		return false
	}(ip) {
		// Only override if upstream service returned no coordinates
		if location.Latitude == 0 && location.Longitude == 0 {
			location.City = "Localhost/Pekanbaru"
			location.Region = "Local"
			location.Country = "Local"
			location.Latitude = 0.5167
			location.Longitude = 101.4417
			location.Timezone = "Asia/Jakarta"
			location.Postal = "00000"
			location.Org = "Local Dev"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message":  "success",
		"ip":       ip,
		"location": location,
	})

}

func getClientIP(r *http.Request) string {

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	// Try X-Real-IP
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return rip
	}
	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
