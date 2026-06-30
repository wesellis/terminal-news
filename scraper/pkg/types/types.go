package types

import "time"

// ParsedArticle represents an article parsed from any source
type ParsedArticle struct {
	Title       string
	URL         string
	Summary     string
	Content     string
	Author      string
	Source      string
	Category    string
	Tags        []string
	ImageURL    string
	PublishedAt time.Time
	FetchedAt   time.Time
	ExternalID  string
	FetchSource string
	TrustScore  float64
}

// City represents a city for weather data
type City struct {
	ID        int64
	Name      string
	State     string
	Country   string
	Latitude  float64
	Longitude float64
}

// FeedSource represents an RSS feed source
type FeedSource struct {
	Name     string
	URL      string
	Category string
	Enabled  bool
}

// WeatherData represents weather information
type WeatherData struct {
	CityID        int64
	Temperature   float64
	FeelsLike     float64
	Condition     string
	Humidity      float64
	WindSpeed     float64
	WindDirection string
	Pressure      float64
	UpdatedAt     time.Time
}
