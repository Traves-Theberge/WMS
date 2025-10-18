package components

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Moon holds the state of the moon component, including phase, illumination,
// and other data.
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

// MoonData represents the structure of the JSON response from the Farmsense API.
type MoonData struct {
	Phase        string   `json:"Phase"`
	Illumination float64  `json:"Illumination"`
	Age          float64  `json:"Age"`
	Moon         []string `json:"Moon"`
}

// MoonResponse is a wrapper for a slice of MoonData.
type MoonResponse []MoonData

// NewMoon creates a new Moon component with a default loading state.
func NewMoon() Moon {
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

// FetchMoonData fetches the current moon phase data from the Farmsense API.
// If the API is unavailable, it falls back to calculating moon phase locally.
func FetchMoonData() (*MoonResponse, error) {
	client := &http.Client{Timeout: 5 * time.Second} // Reduced timeout
	timestamp := time.Now().Unix()
	url := fmt.Sprintf("https://api.farmsense.net/v1/moonphases/?d=%d", timestamp)

	resp, err := client.Get(url)
	if err != nil {
		// Fallback to local calculation if API is unavailable
		return calculateMoonPhaseLocally()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Fallback to local calculation if API returns error
		return calculateMoonPhaseLocally()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return calculateMoonPhaseLocally()
	}

	var moonData MoonResponse
	if err := json.Unmarshal(body, &moonData); err != nil {
		return calculateMoonPhaseLocally()
	}

	return &moonData, nil
}

// UpdateWithData updates the moon component's state with new data from the API.
func (m *Moon) UpdateWithData(data *MoonResponse) {
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
		m.NextPhase, m.DaysToNext = calculateNextPhase(currentMoon.Age)
	}

	m.IsLoading = false
	m.Error = nil
}

// UpdateWithError updates the moon component with an error state
func (m *Moon) UpdateWithError(err error) {
	m.Error = err
	m.IsLoading = false
	// Use fallback data
	m.Phase = "Error loading data"
	m.Illumination = 0.0
	m.NextPhase = "Unknown"
	m.DaysToNext = 0
	m.Icon = "🌙"
	m.MoonName = ""
}

// calculateNextPhase is a helper function to determine the next moon phase.
func calculateNextPhase(currentAge float64) (string, int) {
	const moonCycle = 29.53
	phases := []struct {
		name string
		age  float64
	}{
		{"New Moon", 0},
		{"Waxing Crescent", 3.69},
		{"First Quarter", 7.38},
		{"Waxing Gibbous", 11.07},
		{"Full Moon", 14.77},
		{"Waning Gibbous", 18.46},
		{"Last Quarter", 22.15},
		{"Waning Crescent", 25.84},
	}

	for _, phase := range phases {
		if currentAge < phase.age {
			daysToNext := int(phase.age - currentAge)
			if daysToNext == 0 {
				daysToNext = 1
			}
			return phase.name, daysToNext
		}
	}

	daysToNext := int(moonCycle - currentAge)
	if daysToNext == 0 {
		daysToNext = 1
	}
	return "New Moon", daysToNext
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

// calculateMoonPhaseLocally calculates moon phase using astronomical formulas
// This is a fallback when the API is unavailable
func calculateMoonPhaseLocally() (*MoonResponse, error) {
	now := time.Now()

	// Calculate days since J2000 (January 1, 2000, 12:00 TT)
	j2000 := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	daysSinceJ2000 := now.Sub(j2000).Hours() / 24.0

	// Calculate lunar cycle position (synodic month ≈ 29.53059 days)
	lunarCycle := 29.53059
	cyclePosition := daysSinceJ2000 / lunarCycle
	cyclePosition = cyclePosition - float64(int(cyclePosition)) // Get fractional part

	// Calculate moon age in days
	moonAge := cyclePosition * lunarCycle

	// Calculate illumination based on cycle position
	var illumination float64
	if cyclePosition <= 0.5 {
		illumination = cyclePosition * 2 // 0 to 1 (new to full)
	} else {
		illumination = 2 - (cyclePosition * 2) // 1 to 0 (full to new)
	}

	// Determine phase name based on moon age
	var phaseName string
	switch {
	case moonAge < 1.84566:
		phaseName = "New Moon"
	case moonAge < 5.53699:
		phaseName = "Waxing Crescent"
	case moonAge < 9.22831:
		phaseName = "First Quarter"
	case moonAge < 12.91963:
		phaseName = "Waxing Gibbous"
	case moonAge < 16.61096:
		phaseName = "Full Moon"
	case moonAge < 20.30228:
		phaseName = "Waning Gibbous"
	case moonAge < 23.99361:
		phaseName = "Last Quarter"
	case moonAge < 27.68493:
		phaseName = "Waning Crescent"
	default:
		phaseName = "New Moon"
	}

	// Create moon data response
	moonData := MoonResponse{
		{
			Phase:        phaseName,
			Illumination: illumination,
			Age:          moonAge,
			Moon:         []string{"Calculated"},
		},
	}

	return &moonData, nil
}
