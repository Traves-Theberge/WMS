package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	weatherAPIURL = "http://api.weatherapi.com/v1/current.json"
	weatherAPIKey = "33253c8d785646d18fd184607251207"
)

type WeatherResponse struct {
	Location struct {
		Name      string  `json:"name"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		LocalTime string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC      float64 `json:"temp_c"`
		TempF      float64 `json:"temp_f"`
		IsDay      int     `json:"is_day"`
		Condition  struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph     float64 `json:"wind_mph"`
		WindKph     float64 `json:"wind_kph"`
		WindDir     string  `json:"wind_dir"`
		Humidity    int     `json:"humidity"`
		FeelslikeC  float64 `json:"feelslike_c"`
		FeelslikeF  float64 `json:"feelslike_f"`
		UV          float64 `json:"uv"`
		PrecipMm    float64 `json:"precip_mm"`
		PressureMb  float64 `json:"pressure_mb"`
		Cloud       int     `json:"cloud"`
		Visibility  float64 `json:"vis_km"`
	} `json:"current"`
}

type WeatherClient struct {
	client *http.Client
}

func NewWeatherClient() *WeatherClient {
	return &WeatherClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *WeatherClient) GetCurrentWeather() (*WeatherResponse, error) {
	url := fmt.Sprintf("%s?key=%s&q=auto:ip", weatherAPIURL, weatherAPIKey)
	
	resp, err := w.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var weatherData WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return &weatherData, nil
}

// GetWeatherIcon returns an appropriate weather emoji based on condition
func GetWeatherIcon(condition string, isDay bool) string {
	switch condition {
	case "Sunny", "Clear":
		if isDay {
			return "☀️"
		}
		return "🌙"
	case "Partly cloudy":
		if isDay {
			return "⛅"
		}
		return "☁️"
	case "Cloudy", "Overcast":
		return "☁️"
	case "Mist", "Fog":
		return "🌫️"
	case "Light rain", "Patchy rain possible", "Light drizzle":
		return "🌦️"
	case "Moderate rain", "Heavy rain", "Rain":
		return "🌧️"
	case "Thunderstorm", "Thundery outbreaks possible":
		return "⛈️"
	case "Snow", "Light snow", "Heavy snow":
		return "❄️"
	case "Sleet", "Light sleet":
		return "🌨️"
	default:
		if isDay {
			return "🌤️"
		}
		return "🌙"
	}
}
