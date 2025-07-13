// Package weather provides the core logic for fetching, formatting, and displaying weather data.
package weather

import (
	"fmt"
	"strings"

	"wms/internal/config"
	"wms/internal/ui/icons"
	"wms/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

// DisplayOptions holds all the configuration options that affect how weather data is displayed.
type DisplayOptions struct {
	UseColors    bool
	ShowCityName bool
	Units        string // "metric" or "imperial"
}

// WeatherDisplay is a struct that holds pre-formatted weather data ready for display.
type WeatherDisplay struct {
	Icon          *icons.WeatherIcon
	Location      string
	Condition     string
	Temperature   string
	FeelsLike     string
	Wind          string
	Humidity      string
	UV            string
	Pressure      string
	Visibility    string
	Precipitation string
}

// FormatWeatherDisplay takes raw weather data and a configuration, and returns a
// WeatherDisplay struct with all fields formatted for presentation.
func FormatWeatherDisplay(weather *Weather, cfg config.Config) *WeatherDisplay {
	opts := DisplayOptions{
		UseColors:    cfg.UseColors,
		ShowCityName: cfg.ShowCityName,
		Units:        cfg.Units,
	}

	// Get weather icon
	weatherIcon := icons.GetWeatherIcon(weather.Current.Condition, weather.Current.IsDay == 1, opts.UseColors)

	// Format location
	location := ""
	if opts.ShowCityName {
		location = weather.Location.Name
		if weather.Location.Region != "" && weather.Location.Region != weather.Location.Name {
			location += ", " + weather.Location.Region
		}
		if weather.Location.Country != "" {
			location += ", " + weather.Location.Country
		}
	}

	// Format temperature and units
	var temp, feelsLike string
	var tempUnit string
	if opts.Units == "imperial" {
		temp = fmt.Sprintf("%.1f", weather.Current.TempF)
		feelsLike = fmt.Sprintf("%.1f", weather.Current.FeelslikeF)
		tempUnit = "°F"
	} else {
		temp = fmt.Sprintf("%.1f", weather.Current.TempC)
		feelsLike = fmt.Sprintf("%.1f", weather.Current.FeelslikeC)
		tempUnit = "°C"
	}

	// Format wind
	var windSpeed string
	var windUnit string
	if opts.Units == "imperial" {
		windSpeed = fmt.Sprintf("%.1f", weather.Current.WindMph)
		windUnit = "mph"
	} else {
		windSpeed = fmt.Sprintf("%.1f", weather.Current.WindKph)
		windUnit = "km/h"
	}

	wind := fmt.Sprintf("%s %s %s", windSpeed, windUnit, getWindDirectionSymbol(weather.Current.WindDir))

	// Format other metrics
	humidity := fmt.Sprintf("%d%%", weather.Current.Humidity)
	uv := fmt.Sprintf("%.1f", weather.Current.UV)
	pressure := fmt.Sprintf("%.1f mb", weather.Current.PressureMb)
	visibility := fmt.Sprintf("%.1f km", weather.Current.Visibility)
	precipitation := fmt.Sprintf("%.1f mm", weather.Current.PrecipMm)

	return &WeatherDisplay{
		Icon:          weatherIcon,
		Location:      location,
		Condition:     weather.Current.Condition,
		Temperature:   temp + tempUnit,
		FeelsLike:     feelsLike + tempUnit,
		Wind:          wind,
		Humidity:      humidity,
		UV:            uv,
		Pressure:      pressure,
		Visibility:    visibility,
		Precipitation: precipitation,
	}
}

// RenderWeatherCompact creates a compact, two-column string representation of the
// weather, inspired by the Stormy TUI. It features an ASCII art icon on the
// left and formatted weather data on the right.
func RenderWeatherCompact(weather *Weather, cfg config.Config) string {
	display := FormatWeatherDisplay(weather, cfg)

	// Use themed colors for the weather card
	labelColor := styles.WeatherColor
	valueColor := styles.TextPrimary // Use primary text color for values

	// Create styles
	labelStyle := lipgloss.NewStyle().Foreground(labelColor)
	valueStyle := lipgloss.NewStyle().Foreground(valueColor)

	// Get icon lines
	iconLines := display.Icon.Lines

	// Format precipitation with percentage
	precipPercent := int(weather.Current.Cloud) // Use cloud coverage as precipitation chance
	precipText := fmt.Sprintf("%.1f mm | %d%%", weather.Current.PrecipMm, precipPercent)

	// Prepare text lines to match the exact format from the image - left aligned
	var textLines []string
	textLines = append(textLines, "") // Empty line to match icon spacing
	textLines = append(textLines, labelStyle.Render("Weather")+"  "+valueStyle.Render(display.Condition))
	textLines = append(textLines, labelStyle.Render("Temp")+"     "+valueStyle.Render(display.Temperature))
	textLines = append(textLines, labelStyle.Render("Wind")+"     "+valueStyle.Render(display.Wind))
	textLines = append(textLines, labelStyle.Render("Humidity")+" "+valueStyle.Render(display.Humidity))
	textLines = append(textLines, labelStyle.Render("Precip")+"   "+valueStyle.Render(precipText))
	textLines = append(textLines, "") // Empty line to match icon spacing

	// Combine icon and text with a robust two-column layout
	maxLines := max(len(iconLines), len(textLines))
	for len(iconLines) < maxLines {
		iconLines = append(iconLines, strings.Repeat(" ", 13))
	}
	for len(textLines) < maxLines {
		textLines = append(textLines, "")
	}

	iconBlock := lipgloss.JoinVertical(lipgloss.Left, iconLines...)
	textBlock := lipgloss.JoinVertical(lipgloss.Left, textLines...)

	// Apply background color to the entire block
	return lipgloss.NewStyle().
		Render(lipgloss.JoinHorizontal(lipgloss.Top, iconBlock, "    ", textBlock))
}

// getWindDirectionSymbol converts a wind direction string (e.g., "N", "SSW")
// into a corresponding arrow symbol for a more visual representation.
func getWindDirectionSymbol(dir string) string {
	switch dir {
	case "N":
		return "↑"
	case "NNE", "NE", "ENE":
		return "↗"
	case "E":
		return "→"
	case "ESE", "SE", "SSE":
		return "↘"
	case "S":
		return "↓"
	case "SSW", "SW", "WSW":
		return "↙"
	case "W":
		return "←"
	case "WNW", "NW", "NNW":
		return "↖"
	default:
		return "•"
	}
}

// min is a utility function that returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max is a utility function that returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
