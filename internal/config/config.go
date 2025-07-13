// Package config handles the configuration management for the WMS application.
// It supports loading settings from a TOML file and overriding them with
// command-line flags. It also manages environment variables for API keys.
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

// Config holds all the user-configurable settings for the application.
// Tags are used to map fields to the TOML configuration file.
type Config struct {
	// Weather settings
	WeatherProvider string `toml:"weather_provider"` // The weather API provider to use (e.g., "WeatherAPI")
	Location        string `toml:"location"`         // The default location for weather data
	LocationMode    string `toml:"location_mode"`    // How the location is determined ("ip" or "manual")

	// Display settings
	Units        string `toml:"units"`          // The unit system for temperature and speed ("metric" or "imperial")
	TimeFormat   string `toml:"time_format"`    // The time format ("12" or "24")
	UseColors    bool   `toml:"use_colors"`     // Whether to use colors in the TUI
	Compact      bool   `toml:"compact"`        // Whether to use a compact display mode
	ShowCityName bool   `toml:"show_city_name"` // Whether to show the city name in the display

	// Update settings
	RefreshInterval int `toml:"refresh_interval"` // The refresh interval in minutes

	// API Keys are loaded from a .env file and are not stored in the TOML config.
	WeatherAPIKey string `toml:"-"`
}

// Flags represents the command-line flags that can be used to override the configuration.
type Flags struct {
	Location        string
	LocationMode    string
	Units           string
	TimeFormat      string
	Compact         bool
	Help            bool
	RefreshInterval int
}

// Constants for the supported weather providers.
const (
	ProviderWeatherAPI = "WeatherAPI"
	ProviderOpenMeteo  = "OpenMeteo"
	ProviderIPGeo      = "IPGeolocation"
)

// DefaultConfig returns a new Config with sensible default values.
func DefaultConfig() Config {
	return Config{
		WeatherProvider: ProviderWeatherAPI,
		Location:        "New York",
		LocationMode:    "ip",
		Units:           "metric",
		TimeFormat:      "24",
		UseColors:       true,
		Compact:         false,
		ShowCityName:    true,
		RefreshInterval: 5,
	}
}

// GetConfigPath determines the appropriate path for the configuration file based
// on the user's operating system.
func GetConfigPath() string {
	var configDir string

	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get config directory:", err)
			dir, err = os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to get home directory:", err)
				return ""
			}
			return filepath.Join(dir, "wms", "wms.toml")
		}
		configDir = filepath.Join(dir, "wms")
	} else {
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get home directory:", err)
			return ""
		}
		configDir = filepath.Join(dir, ".config", "wms")
	}

	return filepath.Join(configDir, "wms.toml")
}

// ValidateConfig checks the configuration for valid values and sets defaults
// if any are invalid.
func ValidateConfig(config *Config) {
	// Validate weather provider
	if config.WeatherProvider != ProviderWeatherAPI && config.WeatherProvider != ProviderOpenMeteo {
		fmt.Fprintln(os.Stderr, "Warning: Invalid weather provider in config. Using 'WeatherAPI' as default.")
		config.WeatherProvider = ProviderWeatherAPI
	}

	// Validate location mode
	if config.LocationMode != "ip" && config.LocationMode != "manual" {
		fmt.Fprintln(os.Stderr, "Warning: Invalid location mode in config. Using 'ip' as default.")
		config.LocationMode = "ip"
	}

	// Validate units
	validUnits := map[string]bool{
		"metric":   true,
		"imperial": true,
	}

	if !validUnits[config.Units] {
		fmt.Fprintln(os.Stderr, "Warning: Invalid units in config. Using 'metric' as default.")
		config.Units = "metric"
	}

	// Validate time format
	validTimeFormats := map[string]bool{
		"12": true,
		"24": true,
	}

	if !validTimeFormats[config.TimeFormat] {
		fmt.Fprintln(os.Stderr, "Warning: Invalid time format in config. Using '24' as default.")
		config.TimeFormat = "24"
	}

	// Validate refresh interval
	if config.RefreshInterval < 1 || config.RefreshInterval > 60 {
		fmt.Fprintln(os.Stderr, "Warning: Invalid refresh interval in config. Using 5 minutes as default.")
		config.RefreshInterval = 5
	}

	// Validate API key requirement
	if config.WeatherProvider == ProviderWeatherAPI && config.WeatherAPIKey == "" {
		fmt.Fprintln(os.Stderr, "Warning: 'weather_api_key' is required for WeatherAPI provider.")
	}
}

// LoadEnv loads environment variables from a .env file in the project root.
// It is safe to call even if the file does not exist.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		// It's okay if the .env file doesn't exist, we'll just use env vars
	}
}

// ReadConfig reads the configuration from the TOML file. If the file does not
// exist, it creates a default one. It also loads API keys from the environment.
func ReadConfig() Config {
	LoadEnv() // Load .env file first

	configPath := GetConfigPath()
	if configPath == "" {
		return DefaultConfig()
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create config directory:", err)
			return DefaultConfig()
		}
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := DefaultConfig()
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create config file:", err)
			return defaultConfig
		}
		defer file.Close()

		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(defaultConfig); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write default config:", err)
			return defaultConfig
		}

		fmt.Printf("Config created at %s\n", configPath)
		return defaultConfig
	}

	// Read existing config
	var config Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read config file:", err)
		return DefaultConfig()
	}

	if err := toml.Unmarshal(data, &config); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config file, using defaults.", err)
		return DefaultConfig()
	}

	// Load API keys from environment variables
	config.WeatherAPIKey = os.Getenv("WEATHER_API_KEY")

	// Validate configuration
	ValidateConfig(&config)

	return config
}

// ParseFlags parses the command-line flags and returns them in a Flags struct.
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&flags.Location, "location", "", "Location to get weather for")
	flag.StringVar(&flags.LocationMode, "location-mode", "", "Location mode (ip, manual)")
	flag.StringVar(&flags.Units, "units", "", "Units (metric, imperial)")
	flag.StringVar(&flags.TimeFormat, "time", "", "Time format (12, 24)")
	flag.BoolVar(&flags.Compact, "compact", false, "Compact display mode")
	flag.BoolVar(&flags.Help, "help", false, "Show help")
	flag.IntVar(&flags.RefreshInterval, "refresh", 0, "Refresh interval in minutes")

	// Add usage information
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Weather Management System (WMS) - A comprehensive weather dashboard")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nConfig file is located at:", GetConfigPath())
		fmt.Fprintln(os.Stderr, "\nKeyboard shortcuts:")
		fmt.Fprintln(os.Stderr, "  [U] - Toggle temperature units")
		fmt.Fprintln(os.Stderr, "  [T] - Toggle time format")
		fmt.Fprintln(os.Stderr, "  [S] - Toggle speed units")
		fmt.Fprintln(os.Stderr, "  [R] - Refresh data")
		fmt.Fprintln(os.Stderr, "  [Q] - Quit")
	}

	flag.Parse()

	if flags.Help {
		flag.Usage()
		os.Exit(0)
	}

	return flags
}

// ApplyFlags applies the command-line flags to the Config struct, overriding
// any values that were set in the configuration file.
func ApplyFlags(config *Config, flags Flags) {
	if flags.Location != "" {
		config.Location = flags.Location
	}
	if flags.LocationMode != "" {
		config.LocationMode = flags.LocationMode
		ValidateConfig(config)
	}
	if flags.Units != "" {
		config.Units = flags.Units
		ValidateConfig(config)
	}
	if flags.TimeFormat != "" {
		config.TimeFormat = flags.TimeFormat
		ValidateConfig(config)
	}
	if flags.Compact {
		config.Compact = true
	}
	if flags.RefreshInterval > 0 {
		config.RefreshInterval = flags.RefreshInterval
		ValidateConfig(config)
	}
}

// WriteConfig saves the provided Config struct to the TOML configuration file.
func WriteConfig(config Config) error {
	configPath := GetConfigPath()
	if configPath == "" {
		return fmt.Errorf("could not determine config path")
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
