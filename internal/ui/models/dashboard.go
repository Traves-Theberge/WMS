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
	ViewMoon
	ViewSolar
	ViewSettings      // A new view for the settings menu
	ViewLocationInput // For text input, accessed from settings
	ViewAPIKeyInput   // For API key input, accessed from settings
)

// Model represents the state of the entire application. It contains all the
// data and settings needed to render the TUI.
type Model struct {
	// Core components for weather, moon, and sun data
	weather components.Weather
	moon    components.Moon
	sun     components.Sun

	// TUI dimensions
	width  int
	height int

	// Time and refresh data
	time        time.Time
	lastRefresh time.Time

	// View management
	viewMode    ViewMode
	refreshing  bool
	statusMsg   string
	statusTimer time.Time

	// Configuration
	config config.Config

	// New weather system state
	stormyWeather *weather.Weather
	weatherError  error

	// Location input state
	isEditingLocation bool
	locationInput     string
	settingsCursor    int // For navigating the settings menu

	// API key input state
	isEditingAPIKey bool
	apiKeyInput     string
}

// InitialModel creates the initial model with default settings.
func InitialModel() Model {
	return InitialModelWithConfig(config.DefaultConfig())
}

// InitialModelWithConfig creates the initial model with a given configuration.
// This is the main entry point for initializing the application's state.
func InitialModelWithConfig(cfg config.Config) Model {
	now := time.Now()
	return Model{
		weather:           components.NewWeather(),
		moon:              components.NewMoon(),
		sun:               components.NewSun(),
		time:              now,
		lastRefresh:       now,
		viewMode:          ViewWeather,
		refreshing:        false,
		statusMsg:         "",
		statusTimer:       now,
		config:            cfg,
		stormyWeather:     nil,
		weatherError:      nil,
		isEditingLocation: false,
		locationInput:     cfg.Location,
		settingsCursor:    0,
		isEditingAPIKey:   false,
		apiKeyInput:       cfg.WeatherAPIKey,
	}
}

// Init is the first command that is executed when the application starts. It
// initializes the timers and fetches the initial data.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		refreshCmd(),
		tea.WindowSize(),
		messages.FetchWeatherWithConfigCmd(m.config),
		m.fetchMoonDataCmd(), // Fetch moon data on init
	)
}

// tickCmd creates a command that sends a tick message every second. This is
// used to update the live clock.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// refreshCmd creates a command that sends a refresh message at the configured
// interval. This is used to automatically refresh the weather data.
func refreshCmd() tea.Cmd {
	return tea.Tick(5*time.Minute, func(t time.Time) tea.Msg {
		return refreshMsg(t)
	})
}

// fetchMoonDataCmd creates a command to fetch moon data.
func (m *Model) fetchMoonDataCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := components.FetchMoonData()
		if err != nil {
			return messages.MoonDataMsg{Error: err}
		}
		return messages.MoonDataMsg{Data: data}
	}
}

// Update is the main message loop for the application. It handles all incoming
// messages, from key presses to data updates, and returns an updated model
// and an optional command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// When editing location, only handle input-specific keys
		if m.isEditingLocation {
			return m.updateLocationInputView(msg)
		}
		// When editing API key, only handle input-specific keys
		if m.isEditingAPIKey {
			return m.updateAPIKeyInputView(msg)
		}

		// Global keybindings that work in any view
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.viewMode = ViewWeather
			return m, nil
		case "2":
			m.viewMode = ViewMoon
			return m, nil
		case "3":
			m.viewMode = ViewSolar
			return m, nil
		case "r":
			m.refreshing = true
			m.statusMsg = "Refreshing..."
			m.stormyWeather = nil
			m.weatherError = nil
			return m, messages.FetchWeatherWithConfigCmd(m.config)
		case "u":
			// Cycle through all combinations of units and time formats
			switch {
			case m.config.Units == "metric" && m.config.TimeFormat == "24":
				m.config.TimeFormat = "12"
				m.statusMsg = "Units: Metric, Time: 12h"
			case m.config.Units == "metric" && m.config.TimeFormat == "12":
				m.config.Units = "imperial"
				m.config.TimeFormat = "24"
				m.statusMsg = "Units: Imperial, Time: 24h"
			case m.config.Units == "imperial" && m.config.TimeFormat == "24":
				m.config.TimeFormat = "12"
				m.statusMsg = "Units: Imperial, Time: 12h"
			case m.config.Units == "imperial" && m.config.TimeFormat == "12":
				m.config.Units = "metric"
				m.config.TimeFormat = "24"
				m.statusMsg = "Units: Metric, Time: 24h"
			}
			m.statusTimer = time.Now()
			return m, messages.FetchWeatherWithConfigCmd(m.config)
		case "t":
			// Toggle time format
			if m.config.TimeFormat == "24" {
				m.config.TimeFormat = "12"
				m.statusMsg = "Time: 12-hour"
			} else {
				m.config.TimeFormat = "24"
				m.statusMsg = "Time: 24-hour"
			}
			m.statusTimer = time.Now()
			return m, nil
		case "s":
			// Open the settings menu
			m.viewMode = ViewSettings
			m.statusMsg = "Settings"
			m.statusTimer = time.Now()
			return m, nil
		}

		// Mode-specific keybindings
		switch m.viewMode {
		case ViewWeather, ViewMoon, ViewSolar:
			return m.updateMainView(msg)
		case ViewSettings:
			return m.updateSettingsView(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.time = time.Now()
		m.sun = components.NewSun()
		if time.Since(m.statusTimer) > 3*time.Second {
			m.statusMsg = ""
		}
		return m, tickCmd()

	case refreshMsg:
		return m, messages.FetchWeatherWithConfigCmd(m.config)

	case messages.WeatherMsg:
		m.refreshing = false
		if msg.Error != nil {
			m.weatherError = msg.Error
			m.stormyWeather = nil
		} else {
			m.stormyWeather = msg.Weather
			m.weatherError = nil
		}
		m.statusTimer = time.Now()
		return m, nil

	case messages.MoonDataMsg:
		if msg.Error != nil {
			m.moon.UpdateWithError(msg.Error)
		} else if msg.Data != nil {
			m.moon.UpdateWithData(msg.Data)
		}
		return m, nil
	}
	return m, nil
}

// updateMainView handles keybindings for the main tabbed view.
func (m Model) updateMainView(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.viewMode = (m.viewMode + 1) % 3 // Simple cycle through main views
	case "shift+tab":
		m.viewMode = (m.viewMode - 1 + 3) % 3 // Reverse cycle through main views
	}
	return m, nil
}

// updateSettingsView handles keybindings for the settings menu.
func (m Model) updateSettingsView(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = ViewWeather
		m.statusMsg = ""
		return m, nil
	case "enter":
		switch m.settingsCursor {
		case 0: // Toggle Location Mode
			if m.config.LocationMode == "ip" {
				m.config.LocationMode = "manual"
				m.statusMsg = "Location: Manual"
			} else {
				m.config.LocationMode = "ip"
				m.statusMsg = "Location: IP Detection"
			}
			return m, messages.FetchWeatherWithConfigCmd(m.config)
		case 1: // Set Manual Location
			// Only allow setting location in manual mode
			if m.config.LocationMode == "manual" {
				m.viewMode = ViewLocationInput
				m.isEditingLocation = true
				m.statusMsg = "Enter new location"
			}
		case 2: // Set API Key
			m.viewMode = ViewAPIKeyInput
			m.isEditingAPIKey = true
			m.statusMsg = "Enter WeatherAPI key"
		case 3: // Save and Exit
			err := config.WriteConfig(m.config)
			if err != nil {
				m.statusMsg = "Error saving config"
			} else {
				m.statusMsg = "Config saved!"
			}
			m.viewMode = ViewWeather
			return m, nil
		}
	}

	// Handle cursor navigation
	if msg.String() == "up" {
		m.settingsCursor = (m.settingsCursor - 1 + 4) % 4 // Cycle through 4 options
	} else if msg.String() == "down" {
		m.settingsCursor = (m.settingsCursor + 1) % 4 // Cycle through 4 options
	}

	return m, nil
}

// updateLocationInputView handles keybindings for the location input screen.
func (m Model) updateLocationInputView(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel editing and return to settings
		m.isEditingLocation = false
		m.viewMode = ViewSettings
		m.statusMsg = "Cancelled"
		m.statusTimer = time.Now()
		return m, nil
	case "enter":
		// Save the new location and refresh the weather
		m.config.Location = m.locationInput
		m.isEditingLocation = false
		m.viewMode = ViewSettings // Return to settings after saving
		m.statusMsg = "Location saved!"
		m.statusTimer = time.Now()
		return m, messages.FetchWeatherWithConfigCmd(m.config)
	case "backspace":
		if len(m.locationInput) > 0 {
			m.locationInput = m.locationInput[:len(m.locationInput)-1]
		}
		return m, nil
	default:
		// Only handle printable characters for text input
		if len(msg.String()) == 1 {
			m.locationInput += msg.String()
		}
		return m, nil
	}
}

// updateAPIKeyInputView handles keybindings for the API key input screen.
func (m Model) updateAPIKeyInputView(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel editing and return to settings
		m.isEditingAPIKey = false
		m.viewMode = ViewSettings
		m.statusMsg = "Cancelled"
		m.statusTimer = time.Now()
		return m, nil
	case "enter":
		// Save the new API key to .env file
		m.config.WeatherAPIKey = m.apiKeyInput
		err := config.SaveAPIKey(m.apiKeyInput)
		m.isEditingAPIKey = false
		m.viewMode = ViewSettings
		if err != nil {
			m.statusMsg = "Error saving API key"
		} else {
			m.statusMsg = "API key saved!"
		}
		m.statusTimer = time.Now()
		return m, messages.FetchWeatherWithConfigCmd(m.config)
	case "backspace":
		if len(m.apiKeyInput) > 0 {
			m.apiKeyInput = m.apiKeyInput[:len(m.apiKeyInput)-1]
		}
		return m, nil
	default:
		// Only handle printable characters for text input
		if len(msg.String()) == 1 {
			m.apiKeyInput += msg.String()
		}
		return m, nil
	}
}

func (m Model) formatTime(t time.Time) string {
	if m.config.TimeFormat == "12" {
		return t.Format("3:04:05 PM")
	}
	return t.Format("15:04:05")
}

// View is the main rendering function for the application. It determines which
// view to render based on the current mode and returns it as a string.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	header := m.createTabHeader()
	footer := m.createTabFooter()
	contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)

	var finalContent string

	var activeContent string
	var activeColor lipgloss.Color

	// Pre-generate all panel content to calculate max dimensions
	weatherContent := m.createWeatherPanelContent()
	moonContent := m.createMoonPanelContent()
	solarContent := m.createSolarPanelContent()

	// Calculate max dimensions of the content itself
	maxContentWidth := max(lipgloss.Width(weatherContent), lipgloss.Width(moonContent), lipgloss.Width(solarContent))
	maxContentHeight := max(lipgloss.Height(weatherContent), lipgloss.Height(moonContent), lipgloss.Height(solarContent))

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
	case ViewSettings:
		activeContent = m.renderSettings()
		activeColor = styles.Primary
	case ViewLocationInput:
		activeContent = m.renderLocationInput()
		activeColor = styles.Primary
	case ViewAPIKeyInput:
		activeContent = m.renderAPIKeyInput()
		activeColor = styles.Primary
	}

	// Calculate available space for the card
	availableWidth := m.width - 4  // Leave some margin
	availableHeight := contentHeight - 4

	// Determine card dimensions - use max content size but constrain to available space
	cardWidth := maxContentWidth + 4
	if cardWidth > availableWidth {
		cardWidth = availableWidth
	}

	cardHeight := maxContentHeight + 2
	if cardHeight > availableHeight {
		cardHeight = availableHeight
	}

	// Define a responsive card style
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(activeColor).
		Padding(1, 2).
		Width(cardWidth).
		Height(cardHeight).
		Align(lipgloss.Left, lipgloss.Top)

	// Render the card with the active content
	renderedCard := cardStyle.Render(activeContent)

	// Center the card on the screen
	finalContent = lipgloss.NewStyle().
		Width(m.width).
		Height(contentHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(renderedCard)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		finalContent,
		footer,
	)
}

// createTabHeader creates the header component, which includes the tab navigation
// and status information.
func (m Model) createTabHeader() string {
	// --- Left Block: Time and Location ---
	timeLocationDisplay := styles.ClockStyle.Render(fmt.Sprintf("%s • 📍 %s",
		m.formatTime(m.time),
		getLocationDisplay(m)))

	// --- Center Block: Tabs ---
	weatherTab := "[1] Weather"
	moonTab := "[2] Moon"
	solarTab := "[3] Solar"

	switch m.viewMode {
	case ViewWeather:
		weatherTab = styles.H2Style.Copy().Foreground(styles.WeatherColor).Render("● WEATHER")
	case ViewMoon:
		moonTab = styles.H2Style.Copy().Foreground(styles.MoonColor).Render("● MOON")
	case ViewSolar:
		solarTab = styles.H2Style.Copy().Foreground(styles.SunColor).Render("● SOLAR")
	}
	tabsLine := fmt.Sprintf("%s    %s    %s", weatherTab, moonTab, solarTab)

	// --- Layout with a flexible spring ---
	headerWidth := m.width
	leftWidth := lipgloss.Width(timeLocationDisplay)
	centerWidth := lipgloss.Width(tabsLine)

	sideWidth := (headerWidth - centerWidth) / 2

	if leftWidth > sideWidth {
		return lipgloss.JoinHorizontal(lipgloss.Top, timeLocationDisplay, "   ", tabsLine)
	}

	springWidth := sideWidth - leftWidth
	spring := strings.Repeat(" ", springWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top, timeLocationDisplay, spring, tabsLine)
}

// createTabFooter creates the footer component, which displays the keybindings.
func (m Model) createTabFooter() string {
	// A cleaner footer with a unified units toggle and settings key
	controls := fmt.Sprintf("[R] Refresh    [U] Units (%s, %s)    [S] Settings    [Tab] Switch Tabs    [Q] Quit",
		m.config.Units,
		m.config.TimeFormat+"h")

	return styles.CaptionStyle.Copy().
		Align(lipgloss.Center).
		Render(controls)
}

func getLocationDisplay(m Model) string {
	if m.stormyWeather != nil && m.stormyWeather.Location.Name != "" {
		return m.stormyWeather.Location.Name
	}
	if m.config.Location != "" {
		return m.config.Location
	}
	return "IP Lookup"
}

// createWeatherPanelContent generates the content for the weather tab.
func (m Model) createWeatherPanelContent() string {
	if m.stormyWeather != nil {
		return weather.RenderWeatherCompact(m.stormyWeather, m.config)
	}
	if m.weatherError != nil {
		return lipgloss.JoinVertical(lipgloss.Center, "⚠️ Weather data unavailable")
	}
	return lipgloss.JoinVertical(lipgloss.Center, "⏳ Loading weather...")
}

// createMoonPanelContent generates the content for the moon tab.
func (m Model) createMoonPanelContent() string {
	if m.moon.Error != nil {
		return lipgloss.JoinVertical(lipgloss.Center, "⚠️ Moon data unavailable")
	}
	if m.moon.IsLoading {
		return lipgloss.JoinVertical(lipgloss.Center, "⏳ Loading moon data...")
	}

	moonIcon := getMoonPhaseIcon(m.moon.Phase)
	labelStyle := lipgloss.NewStyle().Foreground(styles.MoonColor)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F4F6"))

	// Create text lines to match weather format exactly, with enough space for all lines
	textLines := []string{
		"", // Empty line to match icon spacing
		labelStyle.Render("Phase") + "    " + valueStyle.Render(m.moon.Phase),
		labelStyle.Render("Illuminated") + "  " + valueStyle.Render(fmt.Sprintf("%.0f%%", m.moon.Illumination)),
		labelStyle.Render("Next") + "     " + valueStyle.Render(m.moon.NextPhase),
		"", // Placeholder for moon name
		"",
		"",
	}

	// Add moon name if available
	if m.moon.MoonName != "" {
		textLines[4] = labelStyle.Render("Name") + "     " + valueStyle.Render(m.moon.MoonName)
	}

	return m.formatTwoColumnContent(moonIcon, textLines)
}

// createSolarPanelContent generates the content for the solar tab.
func (m Model) createSolarPanelContent() string {
	var solarIcon []string
	if strings.ToLower(m.sun.CurrentPos) == "night" {
		solarIcon = getMoonPhaseIcon("Waning Crescent")
	} else {
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

	sunriseStr := m.formatTime(m.sun.Sunrise)
	sunsetStr := m.formatTime(m.sun.Sunset)
	hours := int(m.sun.DayLength.Hours())
	minutes := int(m.sun.DayLength.Minutes()) % 60
	daylightStr := fmt.Sprintf("%dh %dm", hours, minutes)

	// Create styles
	labelStyle := lipgloss.NewStyle().Foreground(styles.SunColor)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F4F6"))

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

	return m.formatTwoColumnContent(solarIcon, textLines)
}

func (m Model) formatTwoColumnContent(iconLines, textLines []string) string {
	maxLines := max(len(iconLines), len(textLines))
	for len(iconLines) < maxLines {
		iconLines = append(iconLines, "")
	}
	for len(textLines) < maxLines {
		textLines = append(textLines, "")
	}

	iconBlock := lipgloss.JoinVertical(lipgloss.Left, iconLines...)
	textBlock := lipgloss.JoinVertical(lipgloss.Left, textLines...)

	return lipgloss.JoinHorizontal(lipgloss.Top, iconBlock, "    ", textBlock)
}

// renderSettings creates the view for the settings menu.
func (m Model) renderSettings() string {
	var b strings.Builder

	b.WriteString(styles.H2Style.Render("Settings"))
	b.WriteString("\n\n")

	// --- Location Mode Setting ---
	cursor := " "
	if m.settingsCursor == 0 {
		cursor = ">"
	}
	modeStatus := fmt.Sprintf("Location Mode: %s", strings.Title(m.config.LocationMode))
	b.WriteString(fmt.Sprintf("%s %s\n", cursor, modeStatus))

	// --- Manual Location Setting ---
	cursor = " "
	locationStyle := lipgloss.NewStyle()
	if m.config.LocationMode == "ip" {
		locationStyle = locationStyle.Foreground(styles.TextMuted)
	}
	if m.settingsCursor == 1 {
		cursor = ">"
	}
	locationStatus := fmt.Sprintf("Set Location:  %s", m.config.Location)
	b.WriteString(fmt.Sprintf("%s %s\n", cursor, locationStyle.Render(locationStatus)))

	// --- API Key Setting ---
	cursor = " "
	if m.settingsCursor == 2 {
		cursor = ">"
	}
	apiKeyDisplay := "Not Set"
	if m.config.WeatherAPIKey != "" {
		// Show only last 4 characters for security
		if len(m.config.WeatherAPIKey) > 4 {
			apiKeyDisplay = "****" + m.config.WeatherAPIKey[len(m.config.WeatherAPIKey)-4:]
		} else {
			apiKeyDisplay = "****"
		}
	}
	apiKeyStatus := fmt.Sprintf("Set API Key:   %s", apiKeyDisplay)
	b.WriteString(fmt.Sprintf("%s %s\n", cursor, apiKeyStatus))

	// --- Save and Exit Setting ---
	cursor = " "
	if m.settingsCursor == 3 {
		cursor = ">"
	}
	saveStatus := "Save and Exit"
	b.WriteString(fmt.Sprintf("%s %s\n", cursor, saveStatus))

	b.WriteString("\n\n")
	b.WriteString(styles.CaptionStyle.Render("Get your free API key at: https://www.weatherapi.com/signup.aspx"))
	b.WriteString("\n")
	b.WriteString(styles.CaptionStyle.Render("(Use ↑/↓ to navigate, Enter to select, Esc to cancel)"))

	return b.String()
}

// renderLocationInput creates the view for the location input screen.
func (m Model) renderLocationInput() string {
	// Create a simple input field for location, styled as a card
	prompt := "Enter new location:"
	inputField := fmt.Sprintf("%s\n\n> %s█", prompt, m.locationInput)

	// Return the content, which will be wrapped in a card by the View function
	return inputField
}

// renderAPIKeyInput creates the view for the API key input screen.
func (m Model) renderAPIKeyInput() string {
	prompt := "Enter your WeatherAPI key:"
	help := "\nGet a free API key at:\nhttps://www.weatherapi.com/signup.aspx"

	// Mask the API key for security (show asterisks)
	maskedInput := strings.Repeat("*", len(m.apiKeyInput))

	inputField := fmt.Sprintf("%s%s\n\n> %s█\n\n%s",
		prompt,
		help,
		maskedInput,
		styles.CaptionStyle.Render("(Press Enter to save, Esc to cancel)"))

	return inputField
}

// getMoonPhaseIcon returns the appropriate ASCII art for the moon phase.
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

// max returns the largest of a list of integers.
func max(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	maxNum := nums[0]
	for _, num := range nums {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}
