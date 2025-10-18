// Package messages defines the messages used for communication between different
// parts of the WMS application.
package messages

import (
	"fmt"

	"wms/internal/config"
	"wms/internal/weather"

	tea "github.com/charmbracelet/bubbletea"
)

// WeatherMsg is a message that is sent when weather data has been fetched. It
// contains either the weather data or an error if the fetch failed.
type WeatherMsg struct {
	Weather *weather.Weather
	Error   error
}

// FetchWeatherWithConfigCmd creates a Bubble Tea command that fetches weather
// data using the new provider system. It takes a Config struct and returns a
// command function that can be executed by the Bubble Tea runtime.
func FetchWeatherWithConfigCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		// Determine location based on LocationMode setting
		var location string
		if cfg.LocationMode == "ip" || cfg.Location == "" {
			// Attempt to automatically detect the user's location via their IP address.
			detectedLocation, err := weather.DetectLocationFromIP()
			if err != nil {
				return WeatherMsg{
					Weather: nil,
					Error:   fmt.Errorf("failed to detect location: %w", err),
				}
			}
			location = detectedLocation
		} else {
			// Use the manually specified location
			location = cfg.Location
		}

		// Create a weather provider based on the configuration.
		provider, err := weather.CreateWeatherProvider(cfg.WeatherProvider, cfg.WeatherAPIKey)
		if err != nil {
			return WeatherMsg{Error: fmt.Errorf("failed to create weather provider: %w", err)}
		}

		// Fetch the weather data using the provider.
		weatherData, err := provider.FetchWeather(location)
		if err != nil {
			return WeatherMsg{
				Weather: nil,
				Error:   fmt.Errorf("failed to fetch weather: %w", err),
			}
		}

		// Return the weather data in a WeatherMsg.
		return WeatherMsg{
			Weather: weatherData,
			Error:   nil,
		}
	}
}
