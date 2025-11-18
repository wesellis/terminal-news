package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WeatherService handles NOAA weather API integration
type WeatherService struct {
	client  *http.Client
	baseURL string
}

// NewWeatherService creates a new weather service instance
func NewWeatherService() *WeatherService {
	return &WeatherService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.weather.gov",
	}
}

// Weather represents weather data in terminal-friendly format
type Weather struct {
	Location     string  `json:"location"`
	Temperature  int     `json:"temperature"` // Fahrenheit
	Condition    string  `json:"condition"`
	WindSpeed    string  `json:"wind_speed"`
	Humidity     int     `json:"humidity"`
	ShortForecast string `json:"short_forecast"`
	DetailedForecast string `json:"detailed_forecast"`
	LastUpdated  string  `json:"last_updated"`
}

// NOAAPointResponse is the response from NOAA points API
type NOAAPointResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
		ForecastHourly string `json:"forecastHourly"`
		RelativeLocation struct {
			Properties struct {
				City  string `json:"city"`
				State string `json:"state"`
			} `json:"properties"`
		} `json:"relativeLocation"`
	} `json:"properties"`
}

// NOAAForecastResponse is the response from NOAA forecast API
type NOAAForecastResponse struct {
	Properties struct {
		Periods []struct {
			Number           int    `json:"number"`
			Name             string `json:"name"`
			Temperature      int    `json:"temperature"`
			TemperatureUnit  string `json:"temperatureUnit"`
			WindSpeed        string `json:"windSpeed"`
			WindDirection    string `json:"windDirection"`
			ShortForecast    string `json:"shortForecast"`
			DetailedForecast string `json:"detailedForecast"`
			RelativeHumidity struct {
				Value int `json:"value"`
			} `json:"relativeHumidity"`
		} `json:"periods"`
	} `json:"properties"`
}

// GetWeather fetches weather for a lat/lng coordinate
func (w *WeatherService) GetWeather(ctx context.Context, lat, lng float64) (*Weather, error) {
	// Step 1: Get the forecast URL for this location
	pointURL := fmt.Sprintf("%s/points/%.4f,%.4f", w.baseURL, lat, lng)

	req, err := http.NewRequestWithContext(ctx, "GET", pointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Terminal-News/1.0")
	req.Header.Set("Accept", "application/geo+json")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch point data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("NOAA API error: %d - %s", resp.StatusCode, string(body))
	}

	var pointResp NOAAPointResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointResp); err != nil {
		return nil, fmt.Errorf("failed to decode point response: %w", err)
	}

	// Step 2: Get the forecast from the forecast URL
	forecastURL := pointResp.Properties.Forecast

	req, err = http.NewRequestWithContext(ctx, "GET", forecastURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecast request: %w", err)
	}

	req.Header.Set("User-Agent", "Terminal-News/1.0")
	req.Header.Set("Accept", "application/geo+json")

	resp, err = w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("NOAA forecast API error: %d - %s", resp.StatusCode, string(body))
	}

	var forecastResp NOAAForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	// Get current period (first one)
	if len(forecastResp.Properties.Periods) == 0 {
		return nil, fmt.Errorf("no forecast periods available")
	}

	current := forecastResp.Properties.Periods[0]

	// Build location string
	location := "Unknown"
	if pointResp.Properties.RelativeLocation.Properties.City != "" {
		location = fmt.Sprintf("%s, %s",
			pointResp.Properties.RelativeLocation.Properties.City,
			pointResp.Properties.RelativeLocation.Properties.State)
	}

	return &Weather{
		Location:         location,
		Temperature:      current.Temperature,
		Condition:        current.ShortForecast,
		WindSpeed:        current.WindSpeed,
		Humidity:         current.RelativeHumidity.Value,
		ShortForecast:    current.ShortForecast,
		DetailedForecast: current.DetailedForecast,
		LastUpdated:      time.Now().Format("3:04 PM MST"),
	}, nil
}

// GetWeatherByZip fetches weather for a ZIP code (simplified - uses approximate coordinates)
func (w *WeatherService) GetWeatherByZip(ctx context.Context, zip string) (*Weather, error) {
	// This is a simplified version. In production, you'd use a geocoding service
	// to convert ZIP to lat/lng. For now, return an error message.
	return nil, fmt.Errorf("ZIP code lookup not implemented - please use lat/lng coordinates")
}