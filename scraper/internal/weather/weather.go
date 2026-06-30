package weather

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

type WeatherClient struct {
	client *resty.Client
	db     *sqlx.DB
}

type NOAAPoint struct {
	Properties struct {
		Forecast            string `json:"forecast"`
		ForecastHourly      string `json:"forecastHourly"`
		ObservationStations string `json:"observationStations"`
	} `json:"properties"`
}

type NOAAForecast struct {
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
		} `json:"periods"`
	} `json:"properties"`
}

type CurrentConditions struct {
	Properties struct {
		Temperature struct {
			Value    float64 `json:"value"`
			UnitCode string  `json:"unitCode"`
		} `json:"temperature"`
		Dewpoint struct {
			Value float64 `json:"value"`
		} `json:"dewpoint"`
		RelativeHumidity struct {
			Value float64 `json:"value"`
		} `json:"relativeHumidity"`
		WindSpeed struct {
			Value float64 `json:"value"`
		} `json:"windSpeed"`
		WindDirection struct {
			Value float64 `json:"value"`
		} `json:"windDirection"`
		BarometricPressure struct {
			Value float64 `json:"value"`
		} `json:"barometricPressure"`
		TextDescription string `json:"textDescription"`
	} `json:"properties"`
}

func NewWeatherClient(db *sqlx.DB) *WeatherClient {
	client := resty.New()
	client.SetBaseURL("https://api.weather.gov")
	client.SetHeader("User-Agent", "TerminalNews Weather Client")
	client.SetTimeout(30 * time.Second)

	return &WeatherClient{
		client: client,
		db:     db,
	}
}

func (w *WeatherClient) UpdateWeatherForCities(cities []types.City) error {
	for _, city := range cities {
		if err := w.updateCityWeather(city); err != nil {
			log.Printf("Failed to update weather for %s: %v", city.Name, err)
			continue
		}

		// Rate limit to avoid overwhelming NOAA
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (w *WeatherClient) updateCityWeather(city types.City) error {
	// Get grid point for coordinates
	point, err := w.getGridPoint(city.Latitude, city.Longitude)
	if err != nil {
		return err
	}

	// Get current conditions
	current, err := w.getCurrentConditions(point.Properties.ObservationStations)
	if err != nil {
		return err
	}

	// Get forecast
	forecast, err := w.getForecast(point.Properties.Forecast)
	if err != nil {
		return err
	}

	// Store in database
	return w.storeWeatherData(city, current, forecast)
}

func (w *WeatherClient) getGridPoint(lat, lon float64) (*NOAAPoint, error) {
	var point NOAAPoint

	resp, err := w.client.R().
		SetResult(&point).
		Get(fmt.Sprintf("/points/%f,%f", lat, lon))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("NOAA API error: %s", resp.Status())
	}

	return &point, nil
}

func (w *WeatherClient) getCurrentConditions(stationURL string) (*CurrentConditions, error) {
	// Get list of stations
	var stations struct {
		Features []struct {
			Properties struct {
				StationIdentifier string `json:"stationIdentifier"`
			} `json:"properties"`
		} `json:"features"`
	}

	_, err := w.client.R().
		SetResult(&stations).
		Get(stationURL)

	if err != nil {
		return nil, err
	}

	if len(stations.Features) == 0 {
		return nil, fmt.Errorf("no weather stations found")
	}

	// Get observation from first station
	stationID := stations.Features[0].Properties.StationIdentifier

	var conditions CurrentConditions
	_, err = w.client.R().
		SetResult(&conditions).
		Get(fmt.Sprintf("/stations/%s/observations/latest", stationID))

	if err != nil {
		return nil, err
	}

	return &conditions, nil
}

func (w *WeatherClient) getForecast(forecastURL string) (*NOAAForecast, error) {
	var forecast NOAAForecast

	// Use full URL as it's provided by the API
	_, err := w.client.SetBaseURL("").R().
		SetResult(&forecast).
		Get(forecastURL)

	if err != nil {
		return nil, err
	}

	// Reset base URL
	w.client.SetBaseURL("https://api.weather.gov")

	return &forecast, nil
}

func (w *WeatherClient) storeWeatherData(city types.City, current *CurrentConditions, forecast *NOAAForecast) error {
	// Convert Celsius to Fahrenheit if needed
	temp := current.Properties.Temperature.Value
	if current.Properties.Temperature.UnitCode == "wmoUnit:degC" {
		temp = temp*9/5 + 32
	}

	// Store current conditions
	query := `
		INSERT INTO weather_current (
			city_id, temperature, feels_like, condition,
			humidity, wind_speed, wind_direction, pressure,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (city_id) DO UPDATE SET
			temperature = EXCLUDED.temperature,
			feels_like = EXCLUDED.feels_like,
			condition = EXCLUDED.condition,
			humidity = EXCLUDED.humidity,
			wind_speed = EXCLUDED.wind_speed,
			wind_direction = EXCLUDED.wind_direction,
			pressure = EXCLUDED.pressure,
			updated_at = NOW()
	`

	_, err := w.db.Exec(query,
		city.ID,
		temp,
		temp, // Calculate feels like based on wind/humidity
		current.Properties.TextDescription,
		current.Properties.RelativeHumidity.Value,
		current.Properties.WindSpeed.Value,
		w.degreeToCompass(current.Properties.WindDirection.Value),
		current.Properties.BarometricPressure.Value,
	)

	if err != nil {
		return err
	}

	// Store forecast
	for i, period := range forecast.Properties.Periods {
		if i >= 10 { // Store up to 5 days (2 periods per day)
			break
		}

		forecastQuery := `
			INSERT INTO weather_forecast (
				city_id, day_index, period_name, temperature,
				wind_speed, condition, detailed_forecast
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (city_id, day_index) DO UPDATE SET
				period_name = EXCLUDED.period_name,
				temperature = EXCLUDED.temperature,
				wind_speed = EXCLUDED.wind_speed,
				condition = EXCLUDED.condition,
				detailed_forecast = EXCLUDED.detailed_forecast
		`

		_, err = w.db.Exec(forecastQuery,
			city.ID,
			i/2, // Day index (2 periods per day)
			period.Name,
			period.Temperature,
			period.WindSpeed,
			period.ShortForecast,
			period.DetailedForecast,
		)

		if err != nil {
			log.Printf("Failed to store forecast: %v", err)
		}
	}

	return nil
}

func (w *WeatherClient) degreeToCompass(degree float64) string {
	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}

	index := int((degree + 11.25) / 22.5) % 16
	return directions[index]
}
