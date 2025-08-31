package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"dev.mfr/weather-api/models"
)

// OpenMeteoResponse represents the subset of fields we consume from the
// Open-Meteo API current weather endpoint.
type OpenMeteoResponse struct {
	Latitude             float64      `json:"latitude"`
	Longitude            float64      `json:"longitude"`
	GenerationTimeMs     float64      `json:"generationtime_ms"`
	UtcOffsetSeconds     int          `json:"utc_offset_seconds"`
	Timezone             string       `json:"timezone"`
	TimezoneAbbreviation string       `json:"timezone_abbreviation"`
	Elevation            float64      `json:"elevation"`
	CurrentUnits         Units        `json:"current_units"`
	Current              CurrentBlock `json:"current"`
}

type Units struct {
	Time                string `json:"time"`
	Interval            string `json:"interval"`
	Temperature2m       string `json:"temperature_2m"`
	IsDay               string `json:"is_day"`
	ApparentTemperature string `json:"apparent_temperature"`
	Precipitation       string `json:"precipitation"`
	Rain                string `json:"rain"`
	Showers             string `json:"showers"`
	Snowfall            string `json:"snowfall"`
	CloudCover          string `json:"cloud_cover"`
	WindSpeed10m        string `json:"wind_speed_10m"`
	WindDirection10m    string `json:"wind_direction_10m"`
	WindGusts10m        string `json:"wind_gusts_10m"`
	RelativeHumidity2m  string `json:"relative_humidity_2m"`
	WeatherCode         string `json:"weather_code"`
	PressureMsl         string `json:"pressure_msl"`
	SurfacePressure     string `json:"surface_pressure"`
}

type CurrentBlock struct {
	Time                string  `json:"time"`
	Interval            int     `json:"interval"`
	Temperature2m       float64 `json:"temperature_2m"`
	IsDay               int     `json:"is_day"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	Precipitation       float64 `json:"precipitation"`
	Rain                float64 `json:"rain"`
	Showers             float64 `json:"showers"`
	Snowfall            float64 `json:"snowfall"`
	CloudCover          int     `json:"cloud_cover"`
	WindSpeed10m        float64 `json:"wind_speed_10m"`
	WindDirection10m    int     `json:"wind_direction_10m"`
	WindGusts10m        float64 `json:"wind_gusts_10m"`
	RelativeHumidity2m  int     `json:"relative_humidity_2m"`
	WeatherCode         int     `json:"weather_code"`
	PressureMsl         float64 `json:"pressure_msl"`
	SurfacePressure     float64 `json:"surface_pressure"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/", index)
	r.Get("/weather", weather)
	r.Get("/ip", ipAddress)

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

	// If connecting locally/private and no "loc" from ipinfo, set default coords
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
		if loc, ok := ipInfo["loc"].(string); !ok || strings.TrimSpace(loc) == "" {
			ipInfo["loc"] = "0.5167,101.4417"
		}
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
		json.NewEncoder(w).Encode(map[string]any{
			"request_id": reqID,
			"client_ip":  ip,
			"ip_info":    ipInfo,
			"weather":    "Unknown",
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"request_id": reqID,
		"client_ip":  ip,
		"ip_info":    ipInfo,
		"weather":    cond,
	})
}

func getWeather(lat float64, lon float64) (map[string]any, bool) {

	if lat == 0 && lon == 0 {
		return map[string]any{"description": "Unknown"}, false
	}

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,is_day,apparent_temperature,precipitation,rain,showers,snowfall,cloud_cover,wind_speed_10m,wind_direction_10m,wind_gusts_10m,relative_humidity_2m,weather_code,pressure_msl,surface_pressure&timezone=auto", lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return map[string]any{"description": "Unknown"}, false
	}
	defer resp.Body.Close()

	var weatherData OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return map[string]any{"description": "Unknown"}, false
	}
	isDay := weatherData.Current.IsDay == 1
	weatherCode, success := models.GetCondition(weatherData.Current.WeatherCode, isDay)
	if !success {
		return map[string]any{"description": "Unknown"}, false
	}
	return map[string]any{
		"description":   weatherCode.Description,
		"code":          weatherData.Current.WeatherCode,
		"is_day":        weatherData.Current.IsDay,
		"temp_c":        weatherData.Current.Temperature2m,
		"apparent_c":    weatherData.Current.ApparentTemperature,
		"wind_kmh":      weatherData.Current.WindSpeed10m,
		"gust_kmh":      weatherData.Current.WindGusts10m,
		"humidity":      weatherData.Current.RelativeHumidity2m,
		"pressure_hpa":  weatherData.Current.PressureMsl,
		"cloud_cover":   weatherData.Current.CloudCover,
		"weather_image": weatherCode.Image,
	}, true
}

func getIpGeolocation(ip string) (map[string]any, error) {

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
	// return map[string]any{
	// 	"ip":        ip,
	// 	"latitude":  37.7749,
	// 	"longitude": -122.4194,
	// 	"country":   "Wonderland",
	// 	"region":    "Fictional",
	// 	"city":      "Imaginary",
	// }, nil
}

func ipAddress(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	ip := clientIP(r)

	ipInfo, err := getIpGeolocation(ip)
	if err != nil {
		http.Error(w, "Failed to get IP geolocation", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"request_id": reqID,
		"client_ip":  ip,
		"ip_info":    ipInfo,
	})

}
