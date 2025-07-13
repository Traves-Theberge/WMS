package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Modern Color Palette - Dark Theme
	Primary   = lipgloss.Color("#60A5FA") // Blue-400
	Secondary = lipgloss.Color("#A78BFA") // Violet-400
	Success   = lipgloss.Color("#34D399") // Emerald-400
	Warning   = lipgloss.Color("#FBBF24") // Amber-400
	Error     = lipgloss.Color("#F87171") // Red-400
	Info      = lipgloss.Color("#38BDF8") // Sky-400

	// Grayscale
	White   = lipgloss.Color("#FFFFFF")
	Gray50  = lipgloss.Color("#F9FAFB")
	Gray100 = lipgloss.Color("#F3F4F6")
	Gray200 = lipgloss.Color("#E5E7EB")
	Gray300 = lipgloss.Color("#D1D5DB")
	Gray400 = lipgloss.Color("#9CA3AF")
	Gray500 = lipgloss.Color("#6B7280")
	Gray600 = lipgloss.Color("#4B5563")
	Gray700 = lipgloss.Color("#374151")
	Gray800 = lipgloss.Color("#1F2937")
	Gray900 = lipgloss.Color("#111827")

	// Component-specific Colors
	WeatherColor = lipgloss.Color("#06B6D4") // Cyan-500
	MoonColor    = lipgloss.Color("#8B5CF6") // Violet-500
	SunColor     = lipgloss.Color("#F59E0B") // Amber-500
	TimeColor    = lipgloss.Color("#10B981") // Emerald-500

	// Typography Scale
	TextPrimary   = Gray50
	TextSecondary = Gray300
	TextMuted     = Gray500
	TextInverse   = Gray900
)

// Base Styles
var (
	BaseStyle = lipgloss.NewStyle().
			Foreground(TextPrimary)

	// Typography
	H1Style = BaseStyle.Copy().
		Bold(true).
		Foreground(Primary).
		MarginBottom(1)

	H2Style = BaseStyle.Copy().
		Bold(true).
		Foreground(TextPrimary)

	H3Style = BaseStyle.Copy().
		Bold(true).
		Foreground(TextSecondary)

	BodyStyle = BaseStyle.Copy().
			Foreground(TextPrimary)

	CaptionStyle = BaseStyle.Copy().
			Foreground(TextMuted)

	// Layout Components - Ultra minimal padding
	ContainerStyle = BaseStyle.Copy().
			Padding(0, 0)

	CardStyle = BaseStyle.Copy().
			Padding(0, 0).
			Margin(0, 0)

	CardHeaderStyle = BaseStyle.Copy().
			Bold(true).
			Foreground(TextPrimary).
			MarginBottom(1)

	// Navigation & Header - Ultra minimal padding
	HeaderStyle = BaseStyle.Copy().
			Bold(true).
			Foreground(Primary).
			Padding(0, 0).
			Align(lipgloss.Center)

	StatusBarStyle = BaseStyle.Copy().
			Foreground(TextMuted).
			Padding(0, 0)

	// Data Display
	MetricLabelStyle = BaseStyle.Copy().
				Foreground(TextMuted).
				Bold(false)

	MetricValueStyle = BaseStyle.Copy().
				Foreground(TextPrimary).
				Bold(true)

	MetricLargeStyle = BaseStyle.Copy().
				Foreground(TextPrimary).
				Bold(true).
				MarginRight(1)

	// Icons and Indicators
	IconStyle = BaseStyle.Copy().
			Bold(true).
			MarginRight(1)

	IconLargeStyle = BaseStyle.Copy().
			Bold(true).
			MarginRight(1)

	// States
	LoadingStyle = BaseStyle.Copy().
			Foreground(Info).
			Italic(true).
			Align(lipgloss.Center)

	ErrorStyle = BaseStyle.Copy().
			Foreground(Error).
			Bold(true).
			Align(lipgloss.Center)

	SuccessStyle = BaseStyle.Copy().
			Foreground(Success).
			Bold(true)

	WarningStyle = BaseStyle.Copy().
			Foreground(Warning).
			Bold(true)

	// Interactive Elements
	ButtonStyle = BaseStyle.Copy().
			Foreground(Primary).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Bold(true)

	ButtonSecondaryStyle = BaseStyle.Copy().
				Foreground(Primary).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Bold(true)

	KeybindStyle = BaseStyle.Copy().
			Foreground(Primary).
			Bold(true)

	// Dividers and Separators
	DividerStyle = BaseStyle.Copy().
			Foreground(Gray600).
			MarginTop(1).
			MarginBottom(1)

	SeparatorStyle = BaseStyle.Copy().
			Foreground(Gray700)

	// Progress and Charts
	ProgressBarStyle = BaseStyle.Copy().
				Foreground(Primary).
				Bold(true)

	ProgressTrackStyle = BaseStyle.Copy().
				Foreground(Gray600)

		// Specialized Component Styles - Clean borderless design
	WeatherCardStyle = CardStyle.Copy()

	MoonCardStyle = CardStyle.Copy()

	SunCardStyle = CardStyle.Copy()

	TimeCardStyle = CardStyle.Copy()

	// Weather-specific
	TemperatureStyle = BaseStyle.Copy().
				Foreground(WeatherColor).
				Bold(true)

	ConditionStyle = BaseStyle.Copy().
			Foreground(TextSecondary).
			Italic(true)

	// Moon-specific
	MoonPhaseStyle = BaseStyle.Copy().
			Foreground(MoonColor).
			Bold(true)

	IlluminationStyle = BaseStyle.Copy().
				Foreground(MoonColor)

	// Sun-specific
	SunTimeStyle = BaseStyle.Copy().
			Foreground(SunColor).
			Bold(true)

	DayLengthStyle = BaseStyle.Copy().
			Foreground(SunColor)

	// Time-specific
	ClockStyle = BaseStyle.Copy().
			Foreground(TimeColor).
			Bold(true)

	DateStyle = BaseStyle.Copy().
			Foreground(TextSecondary)

	// Utility Styles
	CenterStyle = BaseStyle.Copy().
			Align(lipgloss.Center)

	RightStyle = BaseStyle.Copy().
			Align(lipgloss.Right)

	// Responsive helpers
	CompactStyle = BaseStyle.Copy().
			Padding(0, 1)

	SpacingXS = BaseStyle.Copy().Margin(0, 1)
	SpacingSM = BaseStyle.Copy().Margin(0, 2)
	SpacingMD = BaseStyle.Copy().Margin(1, 2)
	SpacingLG = BaseStyle.Copy().Margin(1, 3)
)

// Layout Constants
const (
	MinTerminalWidth  = 80
	MinTerminalHeight = 24
	CardMinWidth      = 25
	CardMinHeight     = 6
)

// Helper Functions
func GetAdaptiveWidth(terminalWidth int, columns int) int {
	if terminalWidth < MinTerminalWidth {
		return CardMinWidth
	}

	padding := 6 // Minimal padding for borders
	availableWidth := terminalWidth - padding
	return (availableWidth / columns) - 2
}

func GetAdaptiveHeight(terminalHeight int, rows int) int {
	if terminalHeight < MinTerminalHeight {
		return CardMinHeight
	}

	padding := 6 // Minimal padding for header and footer
	availableHeight := terminalHeight - padding
	return (availableHeight / rows) - 1
}

func GetResponsiveLayout(width, height int) (columns, rows int) {
	// Determine optimal layout based on terminal size
	if width >= 120 && height >= 30 {
		return 3, 1 // 3 columns, 1 row for large screens
	} else if width >= 90 && height >= 24 {
		return 3, 1 // 3 columns, 1 row for medium screens
	} else {
		return 1, 3 // 1 column, 3 rows for small screens
	}
}
