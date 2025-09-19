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
