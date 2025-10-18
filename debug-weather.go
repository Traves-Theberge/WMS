package main

import (
	"fmt"
	"os"

	"wms/internal/config"
	"wms/internal/weather"
)

func main() {
	// Load configuration
	cfg := config.ReadConfig()

	fmt.Printf("Config loaded:\n")
	fmt.Printf("  Weather Provider: %s\n", cfg.WeatherProvider)
	fmt.Printf("  Location: %s\n", cfg.Location)
	fmt.Printf("  API Key present: %t\n", cfg.WeatherAPIKey != "")
	fmt.Printf("  API Key length: %d\n", len(cfg.WeatherAPIKey))

	if cfg.WeatherAPIKey != "" {
		fmt.Printf("  API Key first 10 chars: %s...\n", cfg.WeatherAPIKey[:10])
	}

	// Test weather provider creation
	provider, err := weather.CreateWeatherProvider(cfg.WeatherProvider, cfg.WeatherAPIKey)
	if err != nil {
		fmt.Printf("Error creating provider: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nProvider created successfully")

	// Test weather fetch
	fmt.Println("Fetching weather for New York...")
	weatherData, err := provider.FetchWeather("New York")
	if err != nil {
		fmt.Printf("Error fetching weather: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Weather fetched successfully!\n")
	fmt.Printf("Location: %s, %s\n", weatherData.Location.Name, weatherData.Location.Country)
	fmt.Printf("Temperature: %.1f°C\n", weatherData.Current.TempC)
	fmt.Printf("Condition: %s\n", weatherData.Current.Condition)
}
