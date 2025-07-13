package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// IPLocationResponse represents the response from IP geolocation service
type IPLocationResponse struct {
	City    string  `json:"city"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Query   string  `json:"query"`
}

// DetectLocationFromIP detects location based on IP address using ip-api.com (free service)
func DetectLocationFromIP() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Use ip-api.com free service (no API key required)
	resp, err := client.Get("http://ip-api.com/json/")
	if err != nil {
		return "", fmt.Errorf("failed to get location from IP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IP geolocation service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var location IPLocationResponse
	if err := json.Unmarshal(body, &location); err != nil {
		return "", fmt.Errorf("failed to parse location response: %w", err)
	}

	// Return city name, or city + region if available
	if location.City != "" {
		if location.Region != "" && location.Region != location.City {
			return fmt.Sprintf("%s, %s", location.City, location.Region), nil
		}
		return location.City, nil
	}

	// Fallback to country if city not available
	if location.Country != "" {
		return location.Country, nil
	}

	return "", fmt.Errorf("no location information available")
}
