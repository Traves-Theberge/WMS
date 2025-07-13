package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ProviderWeatherAPI = "WeatherAPI"
	ProviderOpenMeteo  = "OpenMeteo"
)

// Weather holds the standardized weather data
type Weather struct {
	Location struct {
		Name      string  `json:"name"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		LocalTime string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC      float64 `json:"temp_c"`
		TempF      float64 `json:"temp_f"`
		IsDay      int     `json:"is_day"`
		Condition  string  `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDir    string  `json:"wind_dir"`
		Humidity   int     `json:"humidity"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		UV         float64 `json:"uv"`
		PrecipMm   float64 `json:"precip_mm"`
		PressureMb float64 `json:"pressure_mb"`
		Cloud      int     `json:"cloud"`
		Visibility float64 `json:"vis_km"`
	} `json:"current"`
}

// WeatherAPIResponse represents the response from WeatherAPI
type WeatherAPIResponse struct {
	Location struct {
		Name      string  `json:"name"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		LocalTime string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		TempF     float64 `json:"temp_f"`
		IsDay     int     `json:"is_day"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDir    string  `json:"wind_dir"`
		Humidity   int     `json:"humidity"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		UV         float64 `json:"uv"`
		PrecipMm   float64 `json:"precip_mm"`
		PressureMb float64 `json:"pressure_mb"`
		Cloud      int     `json:"cloud"`
		Visibility float64 `json:"vis_km"`
	} `json:"current"`
}

// OpenMeteoResponse represents the response from Open-Meteo
type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Time               string  `json:"time"`
		Interval           int     `json:"interval"`
		Temperature2m      float64 `json:"temperature_2m"`
		WeatherCode        int     `json:"weather_code"`
		Precipitation      float64 `json:"precipitation"`
		RelativeHumidity2m int     `json:"relative_humidity_2m"`
		WindSpeed10m       float64 `json:"wind_speed_10m"`
		WindDirection10m   int     `json:"wind_direction_10m"`
		IsDay              int     `json:"is_day"`
	} `json:"current"`
}

// GeoResult represents a geocoding result
type GeoResult struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Admin1    string  `json:"admin1"`
}

// GeoResponse represents the geocoding API response
type GeoResponse struct {
	Results []GeoResult `json:"results"`
}

// WeatherProvider interface defines the contract for weather providers
type WeatherProvider interface {
	FetchWeather(location string) (*Weather, error)
	GetProviderName() string
}

// WeatherAPIProvider implements WeatherProvider for WeatherAPI
type WeatherAPIProvider struct {
	APIKey string
	Client *http.Client
}

// OpenMeteoProvider implements WeatherProvider for Open-Meteo
type OpenMeteoProvider struct {
	Client *http.Client
}

// NewWeatherAPIProvider creates a new WeatherAPI provider
func NewWeatherAPIProvider(apiKey string) *WeatherAPIProvider {
	return &WeatherAPIProvider{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// NewOpenMeteoProvider creates a new Open-Meteo provider
func NewOpenMeteoProvider() *OpenMeteoProvider {
	return &OpenMeteoProvider{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchWeather implements WeatherProvider for WeatherAPI
func (w *WeatherAPIProvider) FetchWeather(location string) (*Weather, error) {
	encodedLocation := url.QueryEscape(location)
	apiURL := fmt.Sprintf(
		"http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		w.APIKey,
		encodedLocation,
	)

	resp, err := w.Client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid API key - please check your configuration")
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("location '%s' not found - please check the spelling", location)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var weatherAPIResp WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherAPIResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to standardized format
	weather := &Weather{
		Location: weatherAPIResp.Location,
		Current: struct {
			TempC      float64 `json:"temp_c"`
			TempF      float64 `json:"temp_f"`
			IsDay      int     `json:"is_day"`
			Condition  string  `json:"condition"`
			WindMph    float64 `json:"wind_mph"`
			WindKph    float64 `json:"wind_kph"`
			WindDir    string  `json:"wind_dir"`
			Humidity   int     `json:"humidity"`
			FeelslikeC float64 `json:"feelslike_c"`
			FeelslikeF float64 `json:"feelslike_f"`
			UV         float64 `json:"uv"`
			PrecipMm   float64 `json:"precip_mm"`
			PressureMb float64 `json:"pressure_mb"`
			Cloud      int     `json:"cloud"`
			Visibility float64 `json:"vis_km"`
		}{
			TempC:      weatherAPIResp.Current.TempC,
			TempF:      weatherAPIResp.Current.TempF,
			IsDay:      weatherAPIResp.Current.IsDay,
			Condition:  weatherAPIResp.Current.Condition.Text,
			WindMph:    weatherAPIResp.Current.WindMph,
			WindKph:    weatherAPIResp.Current.WindKph,
			WindDir:    weatherAPIResp.Current.WindDir,
			Humidity:   weatherAPIResp.Current.Humidity,
			FeelslikeC: weatherAPIResp.Current.FeelslikeC,
			FeelslikeF: weatherAPIResp.Current.FeelslikeF,
			UV:         weatherAPIResp.Current.UV,
			PrecipMm:   weatherAPIResp.Current.PrecipMm,
			PressureMb: weatherAPIResp.Current.PressureMb,
			Cloud:      weatherAPIResp.Current.Cloud,
			Visibility: weatherAPIResp.Current.Visibility,
		},
	}

	return weather, nil
}

// GetProviderName returns the provider name
func (w *WeatherAPIProvider) GetProviderName() string {
	return ProviderWeatherAPI
}

// FetchWeather implements WeatherProvider for Open-Meteo
func (o *OpenMeteoProvider) FetchWeather(location string) (*Weather, error) {
	// First, get coordinates for the location
	geoResult, err := o.getFirstGeoResult(location)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Then fetch weather data
	apiURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,weather_code,precipitation,relative_humidity_2m,wind_speed_10m,wind_direction_10m,is_day&wind_speed_unit=kmh&temperature_unit=celsius",
		geoResult.Latitude,
		geoResult.Longitude,
	)

	resp, err := o.Client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openMeteoResp OpenMeteoResponse
	if err := json.Unmarshal(body, &openMeteoResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to standardized format
	weather := &Weather{
		Location: struct {
			Name      string  `json:"name"`
			Region    string  `json:"region"`
			Country   string  `json:"country"`
			Lat       float64 `json:"lat"`
			Lon       float64 `json:"lon"`
			LocalTime string  `json:"localtime"`
		}{
			Name:      geoResult.Name,
			Region:    geoResult.Admin1,
			Country:   geoResult.Country,
			Lat:       geoResult.Latitude,
			Lon:       geoResult.Longitude,
			LocalTime: openMeteoResp.Current.Time,
		},
		Current: struct {
			TempC      float64 `json:"temp_c"`
			TempF      float64 `json:"temp_f"`
			IsDay      int     `json:"is_day"`
			Condition  string  `json:"condition"`
			WindMph    float64 `json:"wind_mph"`
			WindKph    float64 `json:"wind_kph"`
			WindDir    string  `json:"wind_dir"`
			Humidity   int     `json:"humidity"`
			FeelslikeC float64 `json:"feelslike_c"`
			FeelslikeF float64 `json:"feelslike_f"`
			UV         float64 `json:"uv"`
			PrecipMm   float64 `json:"precip_mm"`
			PressureMb float64 `json:"pressure_mb"`
			Cloud      int     `json:"cloud"`
			Visibility float64 `json:"vis_km"`
		}{
			TempC:      openMeteoResp.Current.Temperature2m,
			TempF:      celsiusToFahrenheit(openMeteoResp.Current.Temperature2m),
			IsDay:      openMeteoResp.Current.IsDay,
			Condition:  weatherCodeToCondition(openMeteoResp.Current.WeatherCode),
			WindMph:    kmhToMph(openMeteoResp.Current.WindSpeed10m),
			WindKph:    openMeteoResp.Current.WindSpeed10m,
			WindDir:    degreeToDirection(openMeteoResp.Current.WindDirection10m),
			Humidity:   openMeteoResp.Current.RelativeHumidity2m,
			FeelslikeC: openMeteoResp.Current.Temperature2m, // Open-Meteo doesn't provide feels-like
			FeelslikeF: celsiusToFahrenheit(openMeteoResp.Current.Temperature2m),
			UV:         0, // Open-Meteo doesn't provide UV in basic plan
			PrecipMm:   openMeteoResp.Current.Precipitation,
			PressureMb: 0, // Open-Meteo doesn't provide pressure in basic plan
			Cloud:      0, // Open-Meteo doesn't provide cloud cover in basic plan
			Visibility: 0, // Open-Meteo doesn't provide visibility in basic plan
		},
	}

	return weather, nil
}

// GetProviderName returns the provider name
func (o *OpenMeteoProvider) GetProviderName() string {
	return ProviderOpenMeteo
}

// getFirstGeoResult gets the first geocoding result for a location
func (o *OpenMeteoProvider) getFirstGeoResult(location string) (*GeoResult, error) {
	encodedLocation := url.QueryEscape(location)
	geoURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", encodedLocation)

	resp, err := o.Client.Get(geoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geo GeoResponse
	if err := json.Unmarshal(body, &geo); err != nil {
		return nil, err
	}

	if len(geo.Results) == 0 {
		return nil, fmt.Errorf("no results found for location: %s", location)
	}

	return &geo.Results[0], nil
}

// Helper functions
func celsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

func kmhToMph(kmh float64) float64 {
	return kmh * 0.621371
}

func degreeToDirection(degree int) string {
	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	index := int((float64(degree)+11.25)/22.5) % 16
	return directions[index]
}

func weatherCodeToCondition(code int) string {
	switch code {
	case 0:
		return "Clear"
	case 1, 2, 3:
		return "Partly cloudy"
	case 45, 48:
		return "Fog"
	case 51, 53, 55:
		return "Light rain"
	case 56, 57:
		return "Light rain"
	case 61, 63, 65:
		return "Moderate rain"
	case 66, 67:
		return "Heavy rain"
	case 71, 73, 75:
		return "Light snow"
	case 77:
		return "Heavy snow"
	case 80, 81, 82:
		return "Heavy rain"
	case 85, 86:
		return "Heavy snow"
	case 95:
		return "Thunderstorm"
	case 96, 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

// CreateWeatherProvider creates a weather provider based on the provider name
func CreateWeatherProvider(providerName, apiKey string) (WeatherProvider, error) {
	switch strings.ToLower(providerName) {
	case strings.ToLower(ProviderWeatherAPI):
		if apiKey == "" {
			return nil, fmt.Errorf("API key is required for WeatherAPI provider")
		}
		return NewWeatherAPIProvider(apiKey), nil
	case strings.ToLower(ProviderOpenMeteo):
		return NewOpenMeteoProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported weather provider: %s", providerName)
	}
}

// FetchWeather is a convenience function that creates a provider and fetches weather
func FetchWeather(providerName, apiKey, location string) (*Weather, error) {
	provider, err := CreateWeatherProvider(providerName, apiKey)
	if err != nil {
		return nil, err
	}

	return provider.FetchWeather(location)
}
