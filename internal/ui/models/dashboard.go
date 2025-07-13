package models

import (
	"fmt"
	"strings"
	"time"

	"wms/internal/config"
	"wms/internal/ui/components"
	"wms/internal/ui/messages"
	"wms/internal/ui/styles"
	"wms/internal/weather"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time
type refreshMsg time.Time

type ViewMode int

const (
	ViewWeather ViewMode = iota // Stormy-style weather tab
	ViewMoon                    // Moon phase tab
	ViewSolar                   // Solar information tab
)

type Model struct {
	weather     components.Weather
	moon        components.Moon
	sun         components.Sun
	width       int
	height      int
	time        time.Time
	lastRefresh time.Time
	viewMode    ViewMode
	refreshing  bool
	statusMsg   string
	statusTimer time.Time
	config      config.Config // Add configuration

	// New weather system
	stormyWeather *weather.Weather
	weatherError  error
}

func InitialModel() Model {
	now := time.Now()
	return Model{
		weather:       components.NewWeather(),
		moon:          components.NewMoon(),
		sun:           components.NewSun(),
		time:          now,
		lastRefresh:   now,
		viewMode:      ViewWeather, // Default to weather tab
		refreshing:    false,
		statusMsg:     "",
		statusTimer:   now,
		config:        config.DefaultConfig(),
		stormyWeather: nil,
		weatherError:  nil,
	}
}

func InitialModelWithConfig(cfg config.Config) Model {
	now := time.Now()

	// Map config units to internal format
	timeFormat := cfg.TimeFormat
	if timeFormat == "" {
		timeFormat = "24"
	}

	return Model{
		weather:       components.NewWeather(),
		moon:          components.NewMoon(),
		sun:           components.NewSun(),
		time:          now,
		lastRefresh:   now,
		viewMode:      ViewWeather, // Default to weather tab
		refreshing:    false,
		statusMsg:     "",
		statusTimer:   now,
		config:        cfg,
		stormyWeather: nil,
		weatherError:  nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		refreshCmd(),
		tea.WindowSize(),
		messages.FetchWeatherWithConfigCmd(m.config), // Use new weather system
		messages.FetchMoonCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func refreshCmd() tea.Cmd {
	return tea.Tick(5*time.Minute, func(t time.Time) tea.Msg {
		return refreshMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.refreshing = true
			m.statusMsg = "Refreshing..."
			m.statusTimer = time.Now()
			m.weather = components.NewWeather()
			m.moon = components.NewMoon()
			m.sun = components.NewSun()
			// Clear new weather data
			m.stormyWeather = nil
			m.weatherError = nil
			return m, tea.Batch(
				messages.FetchWeatherWithConfigCmd(m.config), // Use new weather system
				messages.FetchMoonCmd(),
			)
		case "u":
			// Toggle temperature unit
			if m.config.Units == "imperial" {
				m.config.Units = "metric"
				m.statusMsg = "Units: Metric"
			} else {
				m.config.Units = "imperial"
				m.statusMsg = "Units: Imperial"
			}
			m.statusTimer = time.Now()
			return m, messages.FetchWeatherWithConfigCmd(m.config)
		case "t":
			// Toggle time format
			if m.config.TimeFormat == "24" {
				m.config.TimeFormat = "12"
				m.statusMsg = "Time: 12-hour format"
			} else {
				m.config.TimeFormat = "24"
				m.statusMsg = "Time: 24-hour format"
			}
			m.statusTimer = time.Now()
			return m, nil
		case "s":
			// This key is now part of the 'u' (units) toggle
			m.statusMsg = "Unit toggling is now handled by 'u'"
			m.statusTimer = time.Now()
			return m, nil
		case "1", "w":
			// Switch to weather tab
			m.viewMode = ViewWeather
			m.statusMsg = "Weather Tab"
			m.statusTimer = time.Now()
			return m, nil
		case "2", "m":
			// Switch to moon tab
			m.viewMode = ViewMoon
			m.statusMsg = "Moon Tab"
			m.statusTimer = time.Now()
			return m, nil
		case "3", "o":
			// Switch to solar tab
			m.viewMode = ViewSolar
			m.statusMsg = "Solar Tab"
			m.statusTimer = time.Now()
			return m, nil
		case "tab":
			// Cycle through tabs
			switch m.viewMode {
			case ViewWeather:
				m.viewMode = ViewMoon
				m.statusMsg = "Moon Tab"
			case ViewMoon:
				m.viewMode = ViewSolar
				m.statusMsg = "Solar Tab"
			case ViewSolar:
				m.viewMode = ViewWeather
				m.statusMsg = "Weather Tab"
			}
			m.statusTimer = time.Now()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.time = time.Time(msg)
		m.sun = components.NewSun()

		// Clear status message after 3 seconds
		if time.Since(m.statusTimer) > 3*time.Second {
			m.statusMsg = ""
		}

		return m, tickCmd()

	case refreshMsg:
		m.lastRefresh = time.Time(msg)
		return m, tea.Batch(
			messages.FetchWeatherWithConfigCmd(m.config), // Use new weather system
			messages.FetchMoonCmd(),
			refreshCmd(),
		)

	case messages.WeatherMsg: // Handle new weather message
		m.refreshing = false
		if msg.Error != nil {
			m.weatherError = msg.Error
			m.stormyWeather = nil
			m.statusMsg = "Weather failed"
		} else if msg.Weather != nil {
			m.stormyWeather = msg.Weather
			m.weatherError = nil
			m.lastRefresh = time.Now()
			m.statusMsg = "Weather updated"
		}
		m.statusTimer = time.Now()
		return m, nil

	case messages.WeatherDataMsg: // Keep legacy support
		m.refreshing = false
		if msg.Error != nil {
			m.weather.UpdateWithError(msg.Error)
			m.statusMsg = "Weather data failed"
		} else if msg.Data != nil {
			m.weather.UpdateWithData(msg.Data)
			m.lastRefresh = time.Now()
			m.statusMsg = "Updated"
		}
		m.statusTimer = time.Now()
		return m, nil

	case messages.MoonDataMsg:
		if msg.Error != nil {
			m.moon.UpdateWithError(msg.Error)
			m.statusMsg = "Moon data failed"
		} else if msg.Data != nil {
			m.moon.UpdateWithData(msg.Data)
			m.statusMsg = "Updated"
		}
		m.statusTimer = time.Now()
		return m, nil
	}

	return m, nil
}

// Temperature conversion helpers
func (m Model) convertTemp(tempF float64) float64 {
	if m.config.Units == "metric" {
		return (tempF - 32) * 5 / 9
	}
	return tempF
}

func (m Model) getTempUnit() string {
	if m.config.Units == "metric" {
		return "°C"
	}
	return "°F"
}

// Time formatting helper
func (m Model) formatTime(t time.Time) string {
	if m.config.TimeFormat == "12" {
		return t.Format("3:04:05 PM")
	}
	return t.Format("15:04:05")
}

// Speed conversion helpers
func (m Model) convertSpeed(mph float64) float64 {
	if m.config.Units == "metric" {
		return mph * 1.60934
	}
	return mph
}

func (m Model) getSpeedUnit() string {
	if m.config.Units == "metric" {
		return "km/h"
	}
	return "mph"
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Create header with tab navigation
	header := m.createTabHeader()

	// Create footer with controls
	footer := m.createTabFooter()

	// Calculate content area height
	contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)

	// Pre-generate all panel content to calculate max dimensions
	weatherContent := m.createWeatherPanelContent()
	moonContent := m.createMoonPanelContent()
	solarContent := m.createSolarPanelContent()

	// Calculate max dimensions of the content itself
	maxContentWidth := lipgloss.Width(weatherContent)
	if w := lipgloss.Width(moonContent); w > maxContentWidth {
		maxContentWidth = w
	}
	if w := lipgloss.Width(solarContent); w > maxContentWidth {
		maxContentWidth = w
	}

	maxContentHeight := lipgloss.Height(weatherContent)
	if h := lipgloss.Height(moonContent); h > maxContentHeight {
		maxContentHeight = h
	}
	if h := lipgloss.Height(solarContent); h > maxContentHeight {
		maxContentHeight = h
	}

	// Select the content and color for the active tab
	var activeContent string
	var activeColor lipgloss.Color

	switch m.viewMode {
	case ViewWeather:
		activeContent = weatherContent
		activeColor = styles.WeatherColor
	case ViewMoon:
		activeContent = moonContent
		activeColor = styles.MoonColor
	case ViewSolar:
		activeContent = solarContent
		activeColor = styles.SunColor
	default:
		activeContent = weatherContent
		activeColor = styles.WeatherColor
	}

	// Create a fixed-size content block. This ensures all content areas are the same
	// size before the card styling (padding, border) is applied.
	contentBlock := lipgloss.NewStyle().
		Width(maxContentWidth).
		Height(maxContentHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(activeContent)

	// Define a unified card style that will wrap the content block.
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(activeColor).
		Padding(1, 2).
		Background(lipgloss.Color("#2E2E2E"))

	// Render the fixed-size content block inside the card.
	// The background will now correctly fill the padded area.
	renderedCard := cardStyle.Render(contentBlock)

	// Center the card on the screen
	finalContent := lipgloss.NewStyle().
		Width(m.width).
		Height(contentHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(renderedCard)

	// Combine all parts for the final view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		finalContent,
		footer,
	)
}

func (m Model) createTabHeader() string {
	// Tab indicators
	weatherTab := "[1] Weather"
	moonTab := "[2] Moon"
	solarTab := "[3] Solar"

	// Highlight active tab
	switch m.viewMode {
	case ViewWeather:
		weatherTab = styles.H2Style.Copy().Foreground(styles.WeatherColor).Render("● [1] WEATHER")
		moonTab = styles.CaptionStyle.Render("[2] Moon")
		solarTab = styles.CaptionStyle.Render("[3] Solar")
	case ViewMoon:
		weatherTab = styles.CaptionStyle.Render("[1] Weather")
		moonTab = styles.H2Style.Copy().Foreground(styles.MoonColor).Render("● [2] MOON")
		solarTab = styles.CaptionStyle.Render("[3] Solar")
	case ViewSolar:
		weatherTab = styles.CaptionStyle.Render("[1] Weather")
		moonTab = styles.CaptionStyle.Render("[2] Moon")
		solarTab = styles.H2Style.Copy().Foreground(styles.SunColor).Render("● [3] SOLAR")
	}

	// Time with location and date beside it
	timeLocationDisplay := styles.ClockStyle.Render(fmt.Sprintf("%s • 📍 %s • %s",
		m.formatTime(m.time),
		getLocationDisplay(m),
		m.time.Format("Monday, January 2, 2006")))

	// Status indicator
	var statusIndicator string
	if m.refreshing {
		statusIndicator = styles.LoadingStyle.Render("⟳ Refreshing...")
	} else if m.statusMsg != "" {
		statusIndicator = styles.CaptionStyle.Render(m.statusMsg)
	}

	// Layout tabs and info
	tabsLine := fmt.Sprintf("%s  %s  %s", weatherTab, moonTab, solarTab)

	// Calculate spacing for header layout using full width
	headerWidth := m.width
	spacingNeeded := headerWidth - lipgloss.Width(timeLocationDisplay) - lipgloss.Width(tabsLine) - lipgloss.Width(statusIndicator)
	if spacingNeeded < 0 {
		spacingNeeded = 0
	}

	headerLine := fmt.Sprintf("%s%s%s%s%s",
		timeLocationDisplay,
		strings.Repeat(" ", spacingNeeded/3),
		tabsLine,
		strings.Repeat(" ", spacingNeeded-spacingNeeded/3),
		statusIndicator)

	return headerLine
}

func (m Model) createTabFooter() string {
	controls := "[R] Refresh  [U] Units  [T] Time  [S] Speed  [Tab] Switch  [Q] Quit"
	return styles.CaptionStyle.Copy().
		Align(lipgloss.Center).
		Render(controls)
}

func (m Model) createWeatherPanelContent() string {
	var content string

	// Use Stormy-style weather display
	if m.stormyWeather != nil {
		content = weather.RenderWeatherCompact(m.stormyWeather, m.config)
	} else if m.weatherError != nil {
		// Handle error states
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			styles.IconLargeStyle.Render("⚠️"),
			styles.ErrorStyle.Render("Weather Unavailable"),
			styles.CaptionStyle.Render(m.weatherError.Error()),
			"",
			styles.CaptionStyle.Render("Press [R] to refresh"),
		)
	} else if m.weather.IsLoading || m.refreshing {
		// Loading state
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			styles.IconLargeStyle.Render("⏳"),
			styles.LoadingStyle.Render("Loading weather data..."),
			"",
			styles.CaptionStyle.Render("Fetching from "+m.config.WeatherProvider),
		)
	} else if m.weather.Location != "" && m.weather.Location != "Detecting location..." {
		// Legacy weather data fallback
		legacyWeather := &weather.Weather{
			Location: struct {
				Name      string  `json:"name"`
				Region    string  `json:"region"`
				Country   string  `json:"country"`
				Lat       float64 `json:"lat"`
				Lon       float64 `json:"lon"`
				LocalTime string  `json:"localtime"`
			}{
				Name:    m.weather.Location,
				Region:  "",
				Country: "",
			},
			Current: struct {
				TempC      float64 `json:"temp_c"`
				TempF      float64 `json:"temp_f"`
				IsDay      int     `json:"is_day"`
				Condition  string  `json:"condition"`
				WindMph    float64 `json:"wind_mph"`
				WindKph    float64 `json:"wind_kph"`
				WindDir    string  `json:"wind_dir"`
				Humidity   int     `json:"humidity"`
				FeelslikeC float64 `json:"feelslike_c"`
				FeelslikeF float64 `json:"feelslike_f"`
				UV         float64 `json:"uv"`
				PrecipMm   float64 `json:"precip_mm"`
				PressureMb float64 `json:"pressure_mb"`
				Cloud      int     `json:"cloud"`
				Visibility float64 `json:"vis_km"`
			}{
				TempC:      (m.weather.Temperature - 32) * 5 / 9,
				TempF:      m.weather.Temperature,
				IsDay:      1,
				Condition:  m.weather.Condition,
				WindMph:    m.weather.WindSpeed,
				WindKph:    m.weather.WindSpeed * 1.60934,
				WindDir:    "N",
				Humidity:   m.weather.Humidity,
				FeelslikeC: (m.weather.FeelsLike - 32) * 5 / 9,
				FeelslikeF: m.weather.FeelsLike,
				UV:         m.weather.UV,
				PrecipMm:   0,
				PressureMb: 1013.25,
				Cloud:      50,
				Visibility: 10,
			},
		}
		content = weather.RenderWeatherCompact(legacyWeather, m.config)
	} else {
		// Default state for weather tab
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			styles.IconLargeStyle.Render("🌤️"),
			styles.H1Style.Render("Weather Dashboard"),
			"",
			styles.CaptionStyle.Render("Configure location: wms --location \"Your City\""),
			styles.CaptionStyle.Render("Or edit config file: "+config.GetConfigPath()),
		)
	}

	return content
}

func getLocationDisplay(m Model) string {
	if m.stormyWeather != nil && m.stormyWeather.Location.Name != "" {
		return m.stormyWeather.Location.Name
	}
	if m.weather.Location != "" && m.weather.Location != "Detecting location..." {
		return m.weather.Location
	}
	if m.config.Location != "" {
		return m.config.Location
	}
	return "Unknown"
}

// All tab rendering is now handled in the main View() method.
// The render...Tab functions are no longer needed.

// Extract content creation for moon panel
func (m Model) createMoonPanelContent() string {
	if m.moon.IsLoading {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			styles.IconLargeStyle.Render("🌙"),
			styles.LoadingStyle.Render("Loading moon data..."),
		)
	}

	if m.moon.Error != nil {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			styles.IconLargeStyle.Render("⚠️"),
			styles.ErrorStyle.Render("Moon data unavailable"),
		)
	}

	// Get moon phase ASCII art based on current phase
	moonIcon := getMoonPhaseIcon(m.moon.Phase)

	// Use themed colors for the moon card
	labelColor := styles.MoonColor
	valueColor := lipgloss.Color("#F3F4F6") // Keep values bright white

	// Create styles
	labelStyle := lipgloss.NewStyle().Foreground(labelColor)
	valueStyle := lipgloss.NewStyle().Foreground(valueColor)

	// Create text lines to match weather format exactly
	textLines := []string{
		"", // Empty line to match icon spacing
		labelStyle.Render("Phase") + "    " + valueStyle.Render(m.moon.Phase),
		labelStyle.Render("Illumin") + "  " + valueStyle.Render(fmt.Sprintf("%.0f%%", m.moon.Illumination)),
		labelStyle.Render("Next") + "     " + valueStyle.Render(m.moon.NextPhase),
		"",
		"",
		"",
	}

	// Add moon name if available
	if m.moon.MoonName != "" {
		textLines[4] = labelStyle.Render("Name") + "     " + valueStyle.Render(m.moon.MoonName)
	}

	// Combine icon and text with a robust two-column layout
	maxLines := max(len(moonIcon), len(textLines))
	for len(moonIcon) < maxLines {
		moonIcon = append(moonIcon, "")
	}
	for len(textLines) < maxLines {
		textLines = append(textLines, "")
	}

	iconBlock := lipgloss.JoinVertical(lipgloss.Left, moonIcon...)
	textBlock := lipgloss.JoinVertical(lipgloss.Left, textLines...)

	return lipgloss.JoinHorizontal(lipgloss.Top, iconBlock, "    ", textBlock)
}

// Extract content creation for solar panel
func (m Model) createSolarPanelContent() string {
	// Select ASCII art based on day/night status
	var solarIcon []string
	if strings.ToLower(m.sun.CurrentPos) == "night" {
		// Use a simple moon icon for nighttime in the solar view
		solarIcon = getMoonPhaseIcon("Waning Crescent") // Example phase
	} else {
		// Use sun icon for daytime
		solarIcon = []string{
			"              ",
			"    \\   /    ",
			"     .-.      ",
			"  ― (   ) ―   ",
			"     `-'      ",
			"    /   \\    ",
			"              ",
		}
	}

	// Use themed colors for the solar card
	labelColor := styles.SunColor
	valueColor := lipgloss.Color("#F3F4F6") // Keep values bright white

	// Create styles
	labelStyle := lipgloss.NewStyle().Foreground(labelColor)
	valueStyle := lipgloss.NewStyle().Foreground(valueColor)

	// Format times according to user preference
	sunriseStr := m.formatTime(m.sun.Sunrise)
	sunsetStr := m.formatTime(m.sun.Sunset)

	// Calculate daylight duration
	hours := int(m.sun.DayLength.Hours())
	minutes := int(m.sun.DayLength.Minutes()) % 60
	daylightStr := fmt.Sprintf("%dh %dm", hours, minutes)

	// Create text lines to match weather format exactly
	textLines := []string{
		"", // Empty line to match icon spacing
		labelStyle.Render("Status") + "   " + valueStyle.Render(strings.Title(m.sun.CurrentPos)),
		labelStyle.Render("Sunrise") + "  " + valueStyle.Render(sunriseStr),
		labelStyle.Render("Sunset") + "   " + valueStyle.Render(sunsetStr),
		labelStyle.Render("Daylight") + " " + valueStyle.Render(daylightStr),
		"",
		"",
	}

	// Combine icon and text with a robust two-column layout
	maxLines := max(len(solarIcon), len(textLines))
	for len(solarIcon) < maxLines {
		solarIcon = append(solarIcon, "")
	}
	for len(textLines) < maxLines {
		textLines = append(textLines, "")
	}

	iconBlock := lipgloss.JoinVertical(lipgloss.Left, solarIcon...)
	textBlock := lipgloss.JoinVertical(lipgloss.Left, textLines...)

	return lipgloss.JoinHorizontal(lipgloss.Top, iconBlock, "    ", textBlock)
}

// getMoonPhaseIcon returns the appropriate ASCII art for the moon phase
func getMoonPhaseIcon(phase string) []string {
	switch phase {
	case "New Moon":
		return []string{
			"       _..._  ",
			"     .:::::::.  ",
			"    :::::::::::",
			"    :::::::::::",
			"    `:::::::::'",
			"      `':::'  ",
			"              ",
		}
	case "Waxing Crescent":
		return []string{
			"       _..._  ",
			"     .::::. `.",
			"    :::::::.  :",
			"    ::::::::  :",
			"    `::::::' .'",
			"      `'::'-' ",
			"              ",
		}
	case "First Quarter":
		return []string{
			"       _..._  ",
			"     .::::  `.",
			"    ::::::    :",
			"    ::::::    :",
			"    `:::::   .'",
			"      `'::.-' ",
			"              ",
		}
	case "Waxing Gibbous":
		return []string{
			"       _..._  ",
			"     .::'   `.",
			"    :::       :",
			"    :::       :",
			"    `::.     .'",
			"      `':..-' ",
			"              ",
		}
	case "Full Moon":
		return []string{
			"       _..._  ",
			"     .'     `.",
			"    :         :",
			"    :         :",
			"    `.       .'",
			"      `-...-' ",
			"              ",
		}
	case "Waning Gibbous":
		return []string{
			"       _..._  ",
			"     .'   `::.",
			"    :       :::",
			"    :       :::",
			"    `.     .::'",
			"      `-..:'' ",
			"              ",
		}
	case "Last Quarter":
		return []string{
			"       _..._  ",
			"     .'  ::::.",
			"    :    ::::::",
			"    :    ::::::",
			"    `.   :::::'",
			"      `-.::'' ",
			"              ",
		}
	case "Waning Crescent":
		return []string{
			"       _..._  ",
			"     .' .::::.",
			"    :  ::::::::",
			"    :  ::::::::",
			"    `. '::::::'",
			"      `-.::'' ",
			"              ",
		}
	default:
		// Default to full moon if phase is unknown
		return []string{
			"       _..._  ",
			"     .'     `.",
			"    :         :",
			"    :         :",
			"    `.       .'",
			"      `-...-' ",
			"              ",
		}
	}
}
