package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"dev.mfr/go-routine/models"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/", index)
	r.Get("/weather", weather)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	// Safe: helper returns "" if not present
	reqID := middleware.GetReqID(r.Context())

	// Client IP: extract host from "ip:port"
	host := r.RemoteAddr
	var port string
	if ip, prt, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		host = ip
		port = prt
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"request_id":  reqID,
		"client_ip":   host,
		"client_port": port,
	}
	json.NewEncoder(w).Encode(response)
}

// clientIP returns the best-effort real client IP, honoring common proxy headers.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return rip
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}
	return r.RemoteAddr
}

func weather(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	ip := clientIP(r)

	ipInfo, err := getIpGeolocation(ip)
	if err != nil {
		http.Error(w, "Failed to get IP geolocation", http.StatusInternalServerError)
		return
	}

	if loc, ok := ipInfo["loc"].(string); ok {
		var lat, lon float64
		if _, err := fmt.Sscanf(loc, "%f,%f", &lat, &lon); err == nil {
			ipInfo["latitude"] = lat
			ipInfo["longitude"] = lon
		} else {
			fmt.Printf("failed to parse loc %q: %v\n", loc, err)
		}
	}

	// Ensure latitude/longitude are present to avoid panics later
	if _, ok := ipInfo["latitude"].(float64); !ok {
		ipInfo["latitude"] = float64(0)
	}
	if _, ok := ipInfo["longitude"].(float64); !ok {
		ipInfo["longitude"] = float64(0)
	}

	cond, ok := getWeather(ipInfo["latitude"].(float64), ipInfo["longitude"].(float64))
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"request_id": reqID,
			"client_ip":  ip,
			"ip_info":    ipInfo,
			"weather":    "Unknown",
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"request_id": reqID,
		"client_ip":  ip,
		"ip_info":    ipInfo,
		"weather":    cond,
	})
}

func getWeather(lat float64, lon float64) (models.Condition, bool) {
	// Simulate a weather API call
	if lat == 0 && lon == 0 {
		return models.Condition{Description: "Unknown"}, false
	}
	weatherCode, success := models.GetCondition(0, true)
	if !success {
		return models.Condition{Description: "Unknown"}, false
	}
	return weatherCode, true
}

func getIpGeolocation(ip string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get IP geolocation: %s", resp.Status)
	}
	var ipInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		return nil, err
	}
	return ipInfo, nil
	// Simulate an IP geolocation API call
	// return map[string]interface{}{
	// 	"ip":        ip,
	// 	"latitude":  37.7749,
	// 	"longitude": -122.4194,
	// 	"country":   "Wonderland",
	// 	"region":    "Fictional",
	// 	"city":      "Imaginary",
	// }, nil
}
