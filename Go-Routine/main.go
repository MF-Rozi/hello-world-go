package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

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
func weather(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	host := r.RemoteAddr
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		host = ip
	}

	ipInfo, err := getIpGeolocation(host)
	if err != nil {
		http.Error(w, "Failed to get IP geolocation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"request_id": reqID,
		"client_ip":  host,
		"weather":    getWeather(ipInfo["latitude"].(float64), ipInfo["longitude"].(float64)),
	}
	json.NewEncoder(w).Encode(response)
}

func getWeather(lat float64, lon float64) string {
	// Simulate a weather API call
	if lat == 0 && lon == 0 {
		return "Unknown"
	}
	return "Sunny, 25°C"
}

func getIpGeolocation(ip string) (map[string]interface{}, error) {
	// Simulate an IP geolocation API call
	return map[string]interface{}{
		"ip":        ip,
		"latitude":  37.7749,
		"longitude": -122.4194,
		"country":   "Wonderland",
		"region":    "Fictional",
		"city":      "Imaginary",
	}, nil
}
