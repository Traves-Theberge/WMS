package components

import (
	"fmt"
	"log"
	"wms/internal/api"
)

type Weather struct {
	Temperature float64
	Condition   string
	Humidity    int
	WindSpeed   float64
	FeelsLike   float64
	Location    string
	Icon        string
	UV          float64
	IsLoading   bool
	Error       error
}

func NewWeather() Weather {
	// Return default/loading state
	return Weather{
		Temperature: 0.0,
		Condition:   "Loading...",
		Humidity:    0,
		WindSpeed:   0.0,
		FeelsLike:   0.0,
		Location:    "Detecting location...",
		Icon:        "⏳",
		IsLoading:   true,
	}
}

// UpdateWithData updates the weather component with API data
func (w *Weather) UpdateWithData(data *api.WeatherResponse) {
	w.Temperature = data.Current.TempF
	w.Condition = data.Current.Condition.Text
	w.Humidity = data.Current.Humidity
	w.WindSpeed = data.Current.WindMph
	w.FeelsLike = data.Current.FeelslikeF
	w.UV = data.Current.UV
	w.Location = fmt.Sprintf("%s, %s", data.Location.Name, data.Location.Region)

	// Check if it's day based on the is_day field from API
	isDay := data.Current.IsDay == 1
	w.Icon = api.GetWeatherIcon(data.Current.Condition.Text, isDay)

	w.IsLoading = false
	w.Error = nil
}

// UpdateWithError updates the weather component with an error state
func (w *Weather) UpdateWithError(err error) {
	w.Error = err
	w.IsLoading = false
	log.Printf("Failed to fetch weather data: %v", err)

	// Use fallback data
	w.Temperature = 0.0
	w.Condition = "Error loading data"
	w.Humidity = 0
	w.WindSpeed = 0.0
	w.FeelsLike = 0.0
	w.Location = "Unknown Location"
	w.Icon = "❓"
}
