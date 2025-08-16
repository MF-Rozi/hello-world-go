package models

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"sync"
)

//go:embed weather_codes.json
var weatherCodes []byte

var (
	once       sync.Once
	codesByInt map[int]WeatherInfo
	loadErr    error
)

type WeatherInfo struct {
	Day   Condition `json:"day"`
	Night Condition `json:"night"`
}

type Condition struct {
	Description string `json:"description"`
	Image       string `json:"image"`
}

func LoadWeatherCodes() (map[int]WeatherInfo, error) {
	once.Do(func() {
		tmp := make(map[string]WeatherInfo)
		if err := json.Unmarshal(weatherCodes, &tmp); err != nil {
			loadErr = err
			return
		}
		out := make(map[int]WeatherInfo, len(tmp))
		for k, v := range tmp {
			if n, err := strconv.Atoi(k); err == nil {
				out[n] = v // skip non-numeric keys
			}
		}
		codesByInt = out
	})
	return codesByInt, loadErr
}

func GetCondition(code int, isDay bool) (Condition, bool) {
	m, err := LoadWeatherCodes()
	if err != nil {
		return Condition{}, false
	}
	info, exists := m[code]
	if !exists {
		return Condition{}, false
	}
	if isDay {
		return info.Day, true
	}
	return info.Night, true
}
