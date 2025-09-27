package ipgeolocation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GeoLocation represents a normalized geolocation response.
// Fields are adapted from ipinfo.io response plus parsed latitude/longitude
// extracted from the "loc" composite field (format: "lat,lon").
// Sample ipinfo JSON:
//
//	{
//	  "ip": "125.165.110.51",
//	  "city": "Pekanbaru",
//	  "region": "Riau",
//	  "country": "ID",
//	  "loc": "0.5167,101.4417",
//	  "org": "AS7713 PT Telekomunikasi Indonesia",
//	  "postal": "28115",
//	  "timezone": "Asia/Jakarta",
//	  "readme": "https://ipinfo.io/missingauth"
//	}
type GeoLocation struct {
	IP        string  `json:"ip"`
	City      string  `json:"city"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Org       string  `json:"org"`
	Postal    string  `json:"postal"`
	Timezone  string  `json:"timezone"`
}

// rawIPInfo mirrors the raw response from ipinfo.io before normalization.
type rawIPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"` // "lat,lon"
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

func GetGeoLocation(ip string) (*GeoLocation, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get geolocation: %s", resp.Status)
	}

	var raw rawIPInfo
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	locParts := strings.Split(raw.Loc, ",")
	var lat, lon float64
	if len(locParts) == 2 {
		if v, err := strconv.ParseFloat(strings.TrimSpace(locParts[0]), 64); err == nil {
			lat = v
		}
		if v, err := strconv.ParseFloat(strings.TrimSpace(locParts[1]), 64); err == nil {
			lon = v
		}
	}

	return &GeoLocation{
		IP:        raw.IP,
		City:      raw.City,
		Region:    raw.Region,
		Country:   raw.Country,
		Latitude:  lat,
		Longitude: lon,
		Org:       raw.Org,
		Postal:    raw.Postal,
		Timezone:  raw.Timezone,
	}, nil
}
