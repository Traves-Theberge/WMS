// Package weather provides core logic for fetching weather and location data.
package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// IPLocationResponse represents the structure of the JSON response from the
// ip-api.com geolocation service.
type IPLocationResponse struct {
	City    string  `json:"city"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Query   string  `json:"query"`
}

// DetectLocationFromIP attempts to determine the user's location based on their
// public IP address. It uses the free ip-api.com service, which requires no
// API key.
func DetectLocationFromIP() (string, error) {
	// Initialize an HTTP client with a 10-second timeout to prevent the
	// application from hanging on slow network requests.
	client := &http.Client{Timeout: 10 * time.Second}

	// Make a GET request to the ip-api.com JSON endpoint.
	resp, err := client.Get("http://ip-api.com/json/")
	if err != nil {
		return "", fmt.Errorf("failed to get location from IP: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful HTTP status code.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IP geolocation service returned status %d", resp.StatusCode)
	}

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal the JSON response into the IPLocationResponse struct.
	var location IPLocationResponse
	if err := json.Unmarshal(body, &location); err != nil {
		return "", fmt.Errorf("failed to parse location response: %w", err)
	}

	// Return the most specific location information available.
	if location.City != "" {
		if location.Region != "" && location.Region != location.City {
			return fmt.Sprintf("%s, %s", location.City, location.Region), nil
		}
		return location.City, nil
	}

	// Fallback to the country name if the city is not available.
	if location.Country != "" {
		return location.Country, nil
	}

	// If no location information can be determined, return an error.
	return "", fmt.Errorf("no location information available")
}
