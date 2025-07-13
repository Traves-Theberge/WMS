package components

import (
	"log"
	"wms/internal/api"
)

type Moon struct {
	Phase        string
	Illumination float64
	NextPhase    string
	DaysToNext   int
	Icon         string
	MoonName     string
	IsLoading    bool
	Error        error
}

func NewMoon() Moon {
	// Return default/loading state
	return Moon{
		Phase:        "Loading...",
		Illumination: 0.0,
		NextPhase:    "...",
		DaysToNext:   0,
		Icon:         "⏳",
		MoonName:     "",
		IsLoading:    true,
	}
}

// UpdateWithData updates the moon component with API data
func (m *Moon) UpdateWithData(data *api.MoonResponse) {
	if len(*data) > 0 {
		currentMoon := (*data)[0]

		// Update with real data
		m.Phase = currentMoon.Phase
		m.Illumination = currentMoon.Illumination * 100 // Convert to percentage
		m.Icon = GetMoonIcon(currentMoon.Phase)

		// Get moon name if availableu
		if len(currentMoon.Moon) > 0 {
			m.MoonName = currentMoon.Moon[0]
		}

		// Calculate next phase
		m.NextPhase, m.DaysToNext = api.CalculateNextPhase(currentMoon.Age, currentMoon.Phase)
	}

	m.IsLoading = false
	m.Error = nil
}

// UpdateWithError updates the moon component with an error state
func (m *Moon) UpdateWithError(err error) {
	m.Error = err
	m.IsLoading = false
	log.Printf("Failed to fetch moon data: %v", err)

	// Use fallback data
	m.Phase = "Error loading data"
	m.Illumination = 0.0
	m.NextPhase = "Unknown"
	m.DaysToNext = 0
	m.Icon = "🌙"
	m.MoonName = ""
}

// GetMoonIcon returns the appropriate moon emoji based on phase
func GetMoonIcon(phase string) string {
	switch phase {
	case "New Moon":
		return "🌑"
	case "Waxing Crescent":
		return "🌒"
	case "First Quarter":
		return "🌓"
	case "Waxing Gibbous":
		return "🌔"
	case "Full Moon":
		return "🌕"
	case "Waning Gibbous":
		return "🌖"
	case "Last Quarter":
		return "🌗"
	case "Waning Crescent":
		return "🌘"
	default:
		return "🌙"
	}
}
