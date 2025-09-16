package ipgeolocation

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GeoLocation struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TimeZone    string  `json:"time_zone"`
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

	var location GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil, err
	}

	return &location, nil
}
