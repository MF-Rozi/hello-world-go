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

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"request_id": reqID,
		"client_ip":  host,
		"weather":    getWeather(),
	}
	json.NewEncoder(w).Encode(response)
}

func getWeather() string {
	// Simulate a weather API call
	return "Sunny, 25°C"
}
