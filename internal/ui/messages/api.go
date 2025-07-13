package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"wms/internal/api"
)

// WeatherDataMsg is sent when weather data is fetched
type WeatherDataMsg struct {
	Data  *api.WeatherResponse
	Error error
}

// MoonDataMsg is sent when moon phase data is fetched
type MoonDataMsg struct {
	Data  *api.MoonResponse
	Error error
}

// FetchWeatherCmd returns a command that fetches weather data
func FetchWeatherCmd() tea.Cmd {
	return func() tea.Msg {
		client := api.NewWeatherClient()
		data, err := client.GetCurrentWeather()
		return WeatherDataMsg{Data: data, Error: err}
	}
}

// FetchMoonCmd returns a command that fetches moon phase data
func FetchMoonCmd() tea.Cmd {
	return func() tea.Msg {
		client := api.NewMoonClient()
		data, err := client.GetCurrentMoonPhase()
		return MoonDataMsg{Data: data, Error: err}
	}
}
