package messages

import (
	"fmt"

	"wms/internal/config"
	"wms/internal/weather"

	tea "github.com/charmbracelet/bubbletea"
)

// WeatherMsg represents a weather update message
type WeatherMsg struct {
	Weather *weather.Weather
	Error   error
}

// FetchWeatherWithConfigCmd creates a command to fetch weather using the new provider system
func FetchWeatherWithConfigCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		// Use configured location or detect from IP
		location := cfg.Location
		if location == "" {
			// Try to detect location from IP
			detectedLocation, err := weather.DetectLocationFromIP()
			if err != nil {
				return WeatherMsg{
					Weather: nil,
					Error:   fmt.Errorf("failed to detect location: %w", err),
				}
			}
			location = detectedLocation
		}

		// Fetch weather using the new provider system
		weatherData, err := weather.FetchWeather(cfg.WeatherProvider, cfg.WeatherAPIKey, location)
		if err != nil {
			return WeatherMsg{
				Weather: nil,
				Error:   fmt.Errorf("failed to fetch weather: %w", err),
			}
		}

		return WeatherMsg{
			Weather: weatherData,
			Error:   nil,
		}
	}
}
