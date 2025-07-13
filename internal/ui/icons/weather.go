package icons

import (
	"github.com/charmbracelet/lipgloss"
)

// WeatherIcon represents a weather icon with color support
type WeatherIcon struct {
	Lines     []string
	UseColors bool
}

// GetWeatherIcon returns the appropriate weather icon based on condition and day/night
func GetWeatherIcon(condition string, isDay bool, useColors bool) *WeatherIcon {
	iconName := mapConditionToIcon(condition, isDay)
	return &WeatherIcon{
		Lines:     getIcon(iconName, useColors),
		UseColors: useColors,
	}
}

// mapConditionToIcon maps weather conditions to icon names
func mapConditionToIcon(condition string, isDay bool) string {
	switch condition {
	case "Sunny", "Clear":
		if isDay {
			return "Sunny"
		}
		return "Clear"
	case "Partly cloudy", "Partly Cloudy":
		if isDay {
			return "PartlyCloudy"
		}
		return "PartlyCloudyNight"
	case "Cloudy", "Overcast":
		return "Cloudy"
	case "Mist", "Fog":
		return "Fog"
	case "Patchy rain possible", "Light rain", "Moderate rain at times", "Moderate rain", "Light drizzle", "Patchy light drizzle":
		return "LightRain"
	case "Heavy rain at times", "Heavy rain", "Moderate or heavy rain shower", "Torrential rain shower":
		return "HeavyRain"
	case "Patchy snow possible", "Light snow", "Patchy light snow", "Light snow showers":
		return "LightSnow"
	case "Moderate snow", "Heavy snow", "Patchy heavy snow", "Moderate or heavy snow showers", "Blizzard":
		return "HeavySnow"
	case "Thundery outbreaks possible", "Patchy light rain with thunder", "Moderate or heavy rain with thunder":
		return "Thunderstorm"
	case "Patchy sleet possible", "Light sleet", "Moderate or heavy sleet":
		return "Sleet"
	case "Ice pellets", "Light showers of ice pellets", "Moderate or heavy showers of ice pellets":
		return "IcePellets"
	default:
		return "Unknown"
	}
}

// getIcon returns the ASCII art for a given weather condition
func getIcon(name string, useColors bool) []string {
	if useColors {
		return getColoredIcon(name)
	}
	return getMonochromeIcon(name)
}

// getMonochromeIcon returns monochrome ASCII art icons
func getMonochromeIcon(name string) []string {
	icons := map[string][]string{
		"Unknown": {
			"             ",
			"    .-.      ",
			"     __)     ",
			"    (        ",
			"     `-'     ",
			"      •      ",
			"             ",
		},
		"Sunny": {
			"             ",
			"    \\   /    ",
			"     .-.     ",
			"  ― (   ) ―  ",
			"     `-'     ",
			"    /   \\    ",
			"             ",
		},
		"Clear": {
			"             ",
			"      .-.    ",
			"   .-(   )   ",
			"  (       )  ",
			"   `-(   )-' ",
			"      `-'    ",
			"             ",
		},
		"PartlyCloudy": {
			"             ",
			"   \\  /      ",
			" _ /\"\".-.    ",
			"   \\_(   ).  ",
			"   /(___(__) ",
			"             ",
			"             ",
		},
		"PartlyCloudyNight": {
			"             ",
			"   )  (      ",
			" .-(  ).-.   ",
			"(  (    ).  )",
			" `-(____)-'  ",
			"             ",
			"             ",
		},
		"Cloudy": {
			"             ",
			"             ",
			"     .--.    ",
			"  .-(    ).  ",
			" (___.__)__) ",
			"             ",
			"             ",
		},
		"LightRain": {
			"             ",
			" _`/\"\".-.    ",
			"  ,\\_(   ).  ",
			"   /(___(__) ",
			"     ' ' ' ' ",
			"    ' ' ' '  ",
			"             ",
		},
		"HeavyRain": {
			"             ",
			" _`/\"\".-.    ",
			"  ,\\_(   ).  ",
			"   /(___(__) ",
			"   ‚'‚'‚'‚'  ",
			"   ‚'‚'‚'‚'  ",
			"             ",
		},
		"LightSnow": {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    *  *  *  ",
			"   *  *  *   ",
			"             ",
		},
		"HeavySnow": {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"   * * * *   ",
			"  * * * *    ",
			"             ",
		},
		"Thunderstorm": {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    ⚡\"\"⚡\"\" ",
			"  ‚'‚'‚'‚'   ",
			"             ",
		},
		"Fog": {
			"             ",
			"             ",
			" _ - _ - _ - ",
			"  _ - _ - _  ",
			" _ - _ - _ - ",
			"             ",
			"             ",
		},
		"Sleet": {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    ‚ * ‚ *  ",
			"   * ‚ * ‚   ",
			"             ",
		},
		"IcePellets": {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    ° ° ° °  ",
			"   ° ° ° °   ",
			"             ",
		},
	}

	if icon, ok := icons[name]; ok {
		return icon
	}
	return icons["Unknown"]
}

// getColoredIcon returns colored ASCII art icons using lipgloss
func getColoredIcon(name string) []string {
	// Define colors
	sunColor := lipgloss.Color("#FCD34D")       // amber-300
	cloudColor := lipgloss.Color("#9CA3AF")     // gray-400
	darkCloudColor := lipgloss.Color("#6B7280") // gray-500
	rainColor := lipgloss.Color("#60A5FA")      // blue-400
	snowColor := lipgloss.Color("#F3F4F6")      // gray-50
	thunderColor := lipgloss.Color("#FBBF24")   // amber-400
	moonColor := lipgloss.Color("#E5E7EB")      // gray-200
	fogColor := lipgloss.Color("#D1D5DB")       // gray-300

	coloredIcons := map[string][]string{
		"Sunny": {
			"             ",
			lipgloss.NewStyle().Foreground(sunColor).Render("    \\   /    "),
			lipgloss.NewStyle().Foreground(sunColor).Render("     .-.     "),
			lipgloss.NewStyle().Foreground(sunColor).Render("  ― (   ) ―  "),
			lipgloss.NewStyle().Foreground(sunColor).Render("     `-'     "),
			lipgloss.NewStyle().Foreground(sunColor).Render("    /   \\    "),
			"             ",
		},
		"Clear": {
			"             ",
			lipgloss.NewStyle().Foreground(moonColor).Render("      .-.    "),
			lipgloss.NewStyle().Foreground(moonColor).Render("   .-(   )   "),
			lipgloss.NewStyle().Foreground(moonColor).Render("  (       )  "),
			lipgloss.NewStyle().Foreground(moonColor).Render("   `-(   )-' "),
			lipgloss.NewStyle().Foreground(moonColor).Render("      `-'    "),
			"             ",
		},
		"PartlyCloudy": {
			"             ",
			lipgloss.NewStyle().Foreground(sunColor).Render("   \\  /") + "      ",
			lipgloss.NewStyle().Foreground(sunColor).Render(" _ /\"\"") + lipgloss.NewStyle().Foreground(cloudColor).Render(".-.    "),
			lipgloss.NewStyle().Foreground(sunColor).Render("   \\_") + lipgloss.NewStyle().Foreground(cloudColor).Render("(   ).  "),
			lipgloss.NewStyle().Foreground(sunColor).Render("   /") + lipgloss.NewStyle().Foreground(cloudColor).Render("(___(__) "),
			"             ",
			"             ",
		},
		"PartlyCloudyNight": {
			"             ",
			lipgloss.NewStyle().Foreground(moonColor).Render("   )  (      "),
			lipgloss.NewStyle().Foreground(moonColor).Render(" .-(  )") + lipgloss.NewStyle().Foreground(cloudColor).Render(".-.   "),
			lipgloss.NewStyle().Foreground(moonColor).Render("(  (    )") + lipgloss.NewStyle().Foreground(cloudColor).Render(".  )"),
			lipgloss.NewStyle().Foreground(cloudColor).Render(" `-(____)-'  "),
			"             ",
			"             ",
		},
		"Cloudy": {
			"             ",
			"             ",
			lipgloss.NewStyle().Foreground(cloudColor).Render("     .--.    "),
			lipgloss.NewStyle().Foreground(cloudColor).Render("  .-(    ).  "),
			lipgloss.NewStyle().Foreground(cloudColor).Render(" (___.__)__) "),
			"             ",
			"             ",
		},
		"LightRain": {
			"             ",
			lipgloss.NewStyle().Foreground(sunColor).Render(" _`/\"\"") + lipgloss.NewStyle().Foreground(cloudColor).Render(".-.    "),
			lipgloss.NewStyle().Foreground(sunColor).Render("  ,\\_") + lipgloss.NewStyle().Foreground(cloudColor).Render("(   ).  "),
			lipgloss.NewStyle().Foreground(sunColor).Render("   /") + lipgloss.NewStyle().Foreground(cloudColor).Render("(___(__) "),
			lipgloss.NewStyle().Foreground(rainColor).Render("     ' ' ' ' "),
			lipgloss.NewStyle().Foreground(rainColor).Render("    ' ' ' '  "),
			"             ",
		},
		"HeavyRain": {
			"             ",
			lipgloss.NewStyle().Foreground(sunColor).Render(" _`/\"\"") + lipgloss.NewStyle().Foreground(darkCloudColor).Render(".-.    "),
			lipgloss.NewStyle().Foreground(sunColor).Render("  ,\\_") + lipgloss.NewStyle().Foreground(darkCloudColor).Render("(   ).  "),
			lipgloss.NewStyle().Foreground(sunColor).Render("   /") + lipgloss.NewStyle().Foreground(darkCloudColor).Render("(___(__) "),
			lipgloss.NewStyle().Foreground(rainColor).Render("   ‚'‚'‚'‚'  "),
			lipgloss.NewStyle().Foreground(rainColor).Render("   ‚'‚'‚'‚'  "),
			"             ",
		},
		"LightSnow": {
			"             ",
			lipgloss.NewStyle().Foreground(cloudColor).Render("     .-.     "),
			lipgloss.NewStyle().Foreground(cloudColor).Render("    (   ).   "),
			lipgloss.NewStyle().Foreground(cloudColor).Render("   (___(__)  "),
			lipgloss.NewStyle().Foreground(snowColor).Render("    *  *  *  "),
			lipgloss.NewStyle().Foreground(snowColor).Render("   *  *  *   "),
			"             ",
		},
		"HeavySnow": {
			"             ",
			lipgloss.NewStyle().Foreground(darkCloudColor).Render("     .-.     "),
			lipgloss.NewStyle().Foreground(darkCloudColor).Render("    (   ).   "),
			lipgloss.NewStyle().Foreground(darkCloudColor).Render("   (___(__)  "),
			lipgloss.NewStyle().Foreground(snowColor).Render("   * * * *   "),
			lipgloss.NewStyle().Foreground(snowColor).Render("  * * * *    "),
			"             ",
		},
		"Thunderstorm": {
			"             ",
			lipgloss.NewStyle().Foreground(darkCloudColor).Render("     .-.     "),
			lipgloss.NewStyle().Foreground(darkCloudColor).Render("    (   ).   "),
			lipgloss.NewStyle().Foreground(darkCloudColor).Render("   (___(__)  "),
			lipgloss.NewStyle().Foreground(thunderColor).Render("    ⚡") + lipgloss.NewStyle().Foreground(rainColor).Render("\"\"") + lipgloss.NewStyle().Foreground(thunderColor).Render("⚡") + lipgloss.NewStyle().Foreground(rainColor).Render("\"\" "),
			lipgloss.NewStyle().Foreground(rainColor).Render("  ‚'‚'‚'‚'   "),
			"             ",
		},
		"Fog": {
			"             ",
			"             ",
			lipgloss.NewStyle().Foreground(fogColor).Render(" _ - _ - _ - "),
			lipgloss.NewStyle().Foreground(fogColor).Render("  _ - _ - _  "),
			lipgloss.NewStyle().Foreground(fogColor).Render(" _ - _ - _ - "),
			"             ",
			"             ",
		},
		"Sleet": {
			"             ",
			lipgloss.NewStyle().Foreground(cloudColor).Render("     .-.     "),
			lipgloss.NewStyle().Foreground(cloudColor).Render("    (   ).   "),
			lipgloss.NewStyle().Foreground(cloudColor).Render("   (___(__)  "),
			lipgloss.NewStyle().Foreground(rainColor).Render("    ‚ ") + lipgloss.NewStyle().Foreground(snowColor).Render("* ") + lipgloss.NewStyle().Foreground(rainColor).Render("‚ ") + lipgloss.NewStyle().Foreground(snowColor).Render("*  "),
			lipgloss.NewStyle().Foreground(snowColor).Render("   * ") + lipgloss.NewStyle().Foreground(rainColor).Render("‚ ") + lipgloss.NewStyle().Foreground(snowColor).Render("* ") + lipgloss.NewStyle().Foreground(rainColor).Render("‚   "),
			"             ",
		},
		"IcePellets": {
			"             ",
			lipgloss.NewStyle().Foreground(cloudColor).Render("     .-.     "),
			lipgloss.NewStyle().Foreground(cloudColor).Render("    (   ).   "),
			lipgloss.NewStyle().Foreground(cloudColor).Render("   (___(__)  "),
			lipgloss.NewStyle().Foreground(snowColor).Render("    ° ° ° °  "),
			lipgloss.NewStyle().Foreground(snowColor).Render("   ° ° ° °   "),
			"             ",
		},
	}

	if icon, ok := coloredIcons[name]; ok {
		return icon
	}
	return getMonochromeIcon("Unknown")
}

// RenderIcon renders the weather icon as a string
func (w *WeatherIcon) RenderIcon() string {
	result := ""
	for _, line := range w.Lines {
		result += line + "\n"
	}
	return result
}

// GetHeight returns the height of the icon
func (w *WeatherIcon) GetHeight() int {
	return len(w.Lines)
}

// GetWidth returns the width of the icon (assumes all lines are same width)
func (w *WeatherIcon) GetWidth() int {
	if len(w.Lines) == 0 {
		return 0
	}
	// Return the width of the first line (assuming all lines are the same width)
	return len(w.Lines[0])
}
