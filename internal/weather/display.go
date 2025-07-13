package weather

import (
	"fmt"
	"strings"

	"wms/internal/config"
	"wms/internal/ui/icons"

	"github.com/charmbracelet/lipgloss"
)

// DisplayOptions holds display configuration
type DisplayOptions struct {
	UseColors    bool
	Compact      bool
	ShowCityName bool
	Units        string // "metric" or "imperial"
}

// WeatherDisplay holds formatted weather display data
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

// FormatWeatherDisplay formats weather data for display
func FormatWeatherDisplay(weather *Weather, cfg config.Config) *WeatherDisplay {
	opts := DisplayOptions{
		UseColors:    cfg.UseColors,
		Compact:      cfg.Compact,
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

// RenderWeatherPanel renders the weather panel in WMS style
func RenderWeatherPanel(weather *Weather, cfg config.Config, width, height int) string {
	display := FormatWeatherDisplay(weather, cfg)

	// Define colors for different elements
	var (
		titleColor = lipgloss.Color("#60A5FA") // blue-400
		labelColor = lipgloss.Color("#9CA3AF") // gray-400
		valueColor = lipgloss.Color("#F3F4F6") // gray-50
		tempColor  = lipgloss.Color("#06B6D4") // cyan-500
	)

	// Create styles
	titleStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true).
		Align(lipgloss.Center)

	labelStyle := lipgloss.NewStyle().
		Foreground(labelColor).
		Width(10)

	valueStyle := lipgloss.NewStyle().
		Foreground(valueColor).
		Bold(true)

	tempStyle := lipgloss.NewStyle().
		Foreground(tempColor).
		Bold(true)

	conditionStyle := lipgloss.NewStyle().
		Foreground(valueColor).
		Italic(true)

	// Build the panel content
	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render("🌤️  Weather"))
	content.WriteString("\n\n")

	// Location (if enabled)
	if cfg.ShowCityName && display.Location != "" {
		content.WriteString(titleStyle.Render(display.Location))
		content.WriteString("\n\n")
	}

	// Weather icon and condition (side by side)
	if cfg.Compact {
		// Compact mode - icon and basic info
		iconLines := display.Icon.Lines
		if len(iconLines) > 0 {
			// Show first few lines of icon with key info
			for i, line := range iconLines[:min(4, len(iconLines))] {
				content.WriteString(line)
				if i == 0 {
					content.WriteString("  " + conditionStyle.Render(display.Condition))
				} else if i == 1 {
					content.WriteString("  " + tempStyle.Render(display.Temperature))
				} else if i == 2 {
					content.WriteString("  " + valueStyle.Render(display.Wind))
				} else if i == 3 {
					content.WriteString("  " + valueStyle.Render(display.Humidity))
				}
				content.WriteString("\n")
			}
		}
	} else {
		// Full mode - show icon and detailed info
		iconLines := display.Icon.Lines

		// Create info lines
		infoLines := []string{
			"",
			conditionStyle.Render(display.Condition),
			tempStyle.Render(display.Temperature),
			labelStyle.Render("Feels like:") + " " + valueStyle.Render(display.FeelsLike),
			labelStyle.Render("Wind:") + " " + valueStyle.Render(display.Wind),
			labelStyle.Render("Humidity:") + " " + valueStyle.Render(display.Humidity),
		}

		// Add optional fields if they have meaningful values
		if weather.Current.UV > 0 {
			infoLines = append(infoLines, labelStyle.Render("UV Index:")+" "+valueStyle.Render(display.UV))
		}
		if weather.Current.PressureMb > 0 {
			infoLines = append(infoLines, labelStyle.Render("Pressure:")+" "+valueStyle.Render(display.Pressure))
		}
		if weather.Current.Visibility > 0 {
			infoLines = append(infoLines, labelStyle.Render("Visibility:")+" "+valueStyle.Render(display.Visibility))
		}
		if weather.Current.PrecipMm > 0 {
			infoLines = append(infoLines, labelStyle.Render("Precip:")+" "+valueStyle.Render(display.Precipitation))
		}

		// Ensure we have enough info lines to match icon height
		for len(infoLines) < len(iconLines) {
			infoLines = append(infoLines, "")
		}

		// Combine icon and info side by side
		for i := 0; i < len(iconLines); i++ {
			content.WriteString(iconLines[i])
			if i < len(infoLines) {
				content.WriteString("  " + infoLines[i])
			}
			content.WriteString("\n")
		}
	}

	// Apply panel styling
	panelStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Align(lipgloss.Left, lipgloss.Top)

	return panelStyle.Render(content.String())
}

// RenderWeatherCompact renders a compact weather display similar to Stormy
func RenderWeatherCompact(weather *Weather, cfg config.Config) string {
	display := FormatWeatherDisplay(weather, cfg)

	// Define colors to match the reference image
	var (
		labelColor = lipgloss.Color("#9CA3AF") // gray-400 for labels
		valueColor = lipgloss.Color("#F3F4F6") // gray-100 for values
	)

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

	return lipgloss.JoinHorizontal(lipgloss.Top, iconBlock, "    ", textBlock)
}

// Helper functions
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

func getConditionColor(condition string, useColors bool) lipgloss.Color {
	if !useColors {
		return lipgloss.Color("#F3F4F6") // gray-50
	}

	switch {
	case strings.Contains(strings.ToLower(condition), "clear") || strings.Contains(strings.ToLower(condition), "sunny"):
		return lipgloss.Color("#FCD34D") // amber-300
	case strings.Contains(strings.ToLower(condition), "cloud"):
		return lipgloss.Color("#A78BFA") // violet-400
	case strings.Contains(strings.ToLower(condition), "rain"):
		return lipgloss.Color("#60A5FA") // blue-400
	case strings.Contains(strings.ToLower(condition), "snow"):
		return lipgloss.Color("#F3F4F6") // gray-50
	case strings.Contains(strings.ToLower(condition), "thunder"):
		return lipgloss.Color("#FBBF24") // amber-400
	case strings.Contains(strings.ToLower(condition), "fog") || strings.Contains(strings.ToLower(condition), "mist"):
		return lipgloss.Color("#D1D5DB") // gray-300
	default:
		return lipgloss.Color("#F87171") // red-400
	}
}

// Utility functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
