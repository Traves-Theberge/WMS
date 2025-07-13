package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const moonAPIURL = "https://api.farmsense.net/v1/moonphases/"

type MoonData struct {
	Error            int      `json:"Error"`
	ErrorMsg         string   `json:"ErrorMsg"`
	TargetDate       string   `json:"TargetDate"`
	Moon             []string `json:"Moon"`
	Index            int      `json:"Index"`
	Age              float64  `json:"Age"`
	Phase            string   `json:"Phase"`
	Distance         float64  `json:"Distance"`
	Illumination     float64  `json:"Illumination"`
	AngularDiameter  float64  `json:"AngularDiameter"`
	DistanceToSun    float64  `json:"DistanceToSun"`
	SunAngularDiameter float64 `json:"SunAngularDiameter"`
}

type MoonResponse []MoonData

type MoonClient struct {
	client *http.Client
}

func NewMoonClient() *MoonClient {
	return &MoonClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (m *MoonClient) GetCurrentMoonPhase() (*MoonResponse, error) {
	// Get current Unix timestamp
	timestamp := time.Now().Unix()
	url := fmt.Sprintf("%s?d=%d", moonAPIURL, timestamp)
	
	resp, err := m.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch moon data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moon API returned status %d", resp.StatusCode)
	}

	var moonData MoonResponse
	if err := json.NewDecoder(resp.Body).Decode(&moonData); err != nil {
		return nil, fmt.Errorf("failed to decode moon response: %w", err)
	}

	if len(moonData) == 0 {
		return nil, fmt.Errorf("no moon data returned")
	}

	return &moonData, nil
}

// CalculateNextPhase estimates the next moon phase based on current age
func CalculateNextPhase(currentAge float64, currentPhase string) (string, int) {
	// Moon cycle is approximately 29.53 days
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
	
	// Find the next phase
	for _, phase := range phases {
		if currentAge < phase.age {
			daysToNext := int(phase.age - currentAge)
			if daysToNext == 0 {
				daysToNext = 1
			}
			return phase.name, daysToNext
		}
	}
	
	// If we're past the last phase, next is New Moon
	daysToNext := int(moonCycle - currentAge)
	if daysToNext == 0 {
		daysToNext = 1
	}
	return "New Moon", daysToNext
}
