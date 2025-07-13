package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Config holds the application configuration
type Config struct {
	// Weather settings
	WeatherProvider string `toml:"weather_provider"`
	WeatherAPIKey   string `toml:"weather_api_key"`
	Location        string `toml:"location"`

	// Display settings
	Units        string `toml:"units"`       // "metric" or "imperial"
	TimeFormat   string `toml:"time_format"` // "12" or "24"
	UseColors    bool   `toml:"use_colors"`
	Compact      bool   `toml:"compact"`
	ShowCityName bool   `toml:"show_city_name"`

	// Update settings
	RefreshInterval int `toml:"refresh_interval"` // minutes

	// Moon settings
	MoonProvider string `toml:"moon_provider"`
	MoonAPIKey   string `toml:"moon_api_key"`
}

// Flags holds command line flags
type Flags struct {
	Location        string
	Units           string
	TimeFormat      string
	Compact         bool
	Help            bool
	RefreshInterval int
}

const (
	ProviderWeatherAPI = "WeatherAPI"
	ProviderOpenMeteo  = "OpenMeteo"
	ProviderIPGeo      = "IPGeolocation"
)

// DefaultConfig returns a new Config with default values
func DefaultConfig() Config {
	return Config{
		WeatherProvider: ProviderWeatherAPI,
		WeatherAPIKey:   "33253c8d785646d18fd184607251207",
		Location:        "",
		Units:           "metric",
		TimeFormat:      "24",
		UseColors:       true,
		Compact:         false,
		ShowCityName:    true,
		RefreshInterval: 5,
		MoonProvider:    ProviderIPGeo,
		MoonAPIKey:      "7b5a5c79f6e04d6e8b8c9d0e1f2a3b4c",
	}
}

// GetConfigPath returns the path to the config file
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

// ValidateConfig checks if the config is valid
func ValidateConfig(config *Config) {
	// Validate weather provider
	if config.WeatherProvider != ProviderWeatherAPI && config.WeatherProvider != ProviderOpenMeteo {
		fmt.Fprintln(os.Stderr, "Warning: Invalid weather provider in config. Using 'WeatherAPI' as default.")
		config.WeatherProvider = ProviderWeatherAPI
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

// ReadConfig reads/creates the config file and returns the configuration
func ReadConfig() Config {
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
		fmt.Fprintln(os.Stderr, "Failed to parse config file, using defaults with available values:", err)

		// Try to load partial config
		defaultConfig := DefaultConfig()
		var partialConfig map[string]interface{}

		if err := toml.Unmarshal(data, &partialConfig); err == nil {
			// Apply any valid values from partial config
			if provider, ok := partialConfig["weather_provider"].(string); ok {
				defaultConfig.WeatherProvider = provider
			}
			if apiKey, ok := partialConfig["weather_api_key"].(string); ok {
				defaultConfig.WeatherAPIKey = apiKey
			}
			if location, ok := partialConfig["location"].(string); ok {
				defaultConfig.Location = location
			}
			if units, ok := partialConfig["units"].(string); ok {
				defaultConfig.Units = units
			}
			if timeFormat, ok := partialConfig["time_format"].(string); ok {
				defaultConfig.TimeFormat = timeFormat
			}
			if useColors, ok := partialConfig["use_colors"].(bool); ok {
				defaultConfig.UseColors = useColors
			}
			if compact, ok := partialConfig["compact"].(bool); ok {
				defaultConfig.Compact = compact
			}
			if showCityName, ok := partialConfig["show_city_name"].(bool); ok {
				defaultConfig.ShowCityName = showCityName
			}
			if refreshInterval, ok := partialConfig["refresh_interval"].(int); ok {
				defaultConfig.RefreshInterval = refreshInterval
			}
		}

		// Write corrected config back
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to update config file:", err)
			return defaultConfig
		}
		defer file.Close()

		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(defaultConfig); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write merged config:", err)
		}

		config = defaultConfig
	}

	// Validate configuration
	ValidateConfig(&config)

	return config
}

// ParseFlags parses command line flags
func ParseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&flags.Location, "location", "", "Location to get weather for")
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

// ApplyFlags applies command line flags to the config
func ApplyFlags(config *Config, flags Flags) {
	if flags.Location != "" {
		config.Location = flags.Location
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
