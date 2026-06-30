package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wesellis/terminal-news/cli/internal/models"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
)

// WeatherWidget displays weather information
type WeatherWidget struct {
	location    string
	weather     *models.Weather
	sponsor     string
	expanded    bool
	width       int
	height      int
	styles      *styles.Styles
	lastUpdated time.Time
}

// NewWeatherWidget creates a new weather widget
func NewWeatherWidget(location string, styles *styles.Styles) *WeatherWidget {
	return &WeatherWidget{
		location: location,
		sponsor:  "NOAA National Weather Service",
		expanded: false,
		styles:   styles,
	}
}

// SetWeather updates the weather data
func (ww *WeatherWidget) SetWeather(weather *models.Weather) {
	ww.weather = weather
	ww.lastUpdated = time.Now()
}

// SetSponsor sets the sponsor name
func (ww *WeatherWidget) SetSponsor(sponsor string) {
	ww.sponsor = sponsor
}

// SetExpanded toggles between compact and full view
func (ww *WeatherWidget) SetExpanded(expanded bool) {
	ww.expanded = expanded
}

// SetSize updates dimensions
func (ww *WeatherWidget) SetSize(width, height int) {
	ww.width = width
	ww.height = height
}

// Update handles input (for full view)
func (ww *WeatherWidget) Update(msg tea.Msg) (*WeatherWidget, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			ww.expanded = !ww.expanded
		}
	}
	return ww, nil
}

// View renders the weather widget
func (ww *WeatherWidget) View() string {
	if ww.expanded {
		return ww.renderExpanded()
	}
	return ww.renderCompact()
}

// CompactView renders compact header widget
func (ww *WeatherWidget) CompactView() string {
	if ww.weather == nil {
		return ww.styles.WeatherWidget.Render(fmt.Sprintf("Weather: %s (no data)", ww.location))
	}

	temp := fmt.Sprintf("%d°F", ww.weather.Current.Temperature)
	condition := ww.weather.Current.Condition
	icon := ww.getConditionIcon(condition)

	// Format: 📍 Location | 🌤 72°F Clear | Updated 2m ago | ☕ Sponsored by Coffee Shop
	compact := fmt.Sprintf("📍 %s  │  %s %s %s  │  Updated %s  │  ☕ %s",
		ww.location,
		icon,
		temp,
		condition,
		ww.formatLastUpdated(),
		ww.sponsor,
	)

	return ww.styles.WeatherWidget.Render(compact)
}

// renderCompact renders compact widget
func (ww *WeatherWidget) renderCompact() string {
	return ww.CompactView()
}

// renderExpanded renders full weather view
func (ww *WeatherWidget) renderExpanded() string {
	if ww.weather == nil {
		return ww.renderNoData()
	}

	var sections []string

	// Header
	header := ww.renderHeader()
	sections = append(sections, header)

	// Current conditions
	current := ww.renderCurrentConditions()
	sections = append(sections, current)

	// 5-day forecast
	forecast := ww.renderForecast()
	sections = append(sections, forecast)

	// Footer with sponsor
	footer := ww.renderFooter()
	sections = append(sections, footer)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return ww.styles.WeatherWidget.Render(content)
}

// renderHeader renders the header section
func (ww *WeatherWidget) renderHeader() string {
	title := fmt.Sprintf("WEATHER - %s", strings.ToUpper(ww.location))
	subtitle := fmt.Sprintf("Updated %s  •  %s", ww.formatLastUpdated(), ww.sponsor)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		ww.styles.Header.Render(title),
		ww.styles.Meta.Render(subtitle),
		"",
	)
}

// renderCurrentConditions renders current weather
func (ww *WeatherWidget) renderCurrentConditions() string {
	current := ww.weather.Current

	// ASCII art weather icon
	ascii := ww.getWeatherASCII(current.Condition)

	// Temperature and condition
	temp := fmt.Sprintf("%d°F", current.Temperature)
	feelsLike := fmt.Sprintf("Feels like %d°F", current.FeelsLike)
	condition := current.Condition

	// Details
	details := []string{
		fmt.Sprintf("💧 Humidity: %d%%", current.Humidity),
		fmt.Sprintf("💨 Wind: %d mph %s", current.WindSpeed, current.WindDir),
		fmt.Sprintf("🌡 Pressure: %.2f in", current.Pressure),
	}

	// Left side: ASCII art
	left := ascii

	// Right side: Details
	right := lipgloss.JoinVertical(
		lipgloss.Left,
		ww.styles.Title.Render(temp+" "+condition),
		ww.styles.Meta.Render(feelsLike),
		"",
		strings.Join(details, "\n"),
	)

	// Combine side by side
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		"    ",
		right,
	)

	return combined
}

// renderForecast renders 5-day forecast
func (ww *WeatherWidget) renderForecast() string {
	if len(ww.weather.Forecast) == 0 {
		return ""
	}

	var days []string
	var icons []string
	var highs []string
	var lows []string

	// Take up to 5 days
	count := len(ww.weather.Forecast)
	if count > 5 {
		count = 5
	}

	for i := 0; i < count; i++ {
		forecast := ww.weather.Forecast[i]
		days = append(days, fmt.Sprintf("%-10s", forecast.Day))
		icons = append(icons, fmt.Sprintf("%-10s", ww.getConditionIcon(forecast.Condition)))
		highs = append(highs, fmt.Sprintf("%-10s", fmt.Sprintf("H:%d°", forecast.High)))
		lows = append(lows, fmt.Sprintf("%-10s", fmt.Sprintf("L:%d°", forecast.Low)))
	}

	forecastBox := fmt.Sprintf(`
┌─ 5-Day Forecast ────────────────────────────────────────────────────┐
│                                                                      │
│ %s │
│ %s │
│ %s │
│ %s │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
	`,
		strings.Join(days, " "),
		strings.Join(icons, " "),
		strings.Join(highs, " "),
		strings.Join(lows, " "),
	)

	return forecastBox
}

// renderFooter renders footer with data source
func (ww *WeatherWidget) renderFooter() string {
	return ww.styles.HelpText.Render(
		"\nData from NOAA National Weather Service  •  Press 'E' to toggle view",
	)
}

// renderNoData renders no data state
func (ww *WeatherWidget) renderNoData() string {
	noData := fmt.Sprintf(`
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║                  Weather Data Unavailable                  ║
║                                                            ║
║  Location: %s                                 ║
║                                                            ║
║  • Check your internet connection                          ║
║  • Verify location in config.yaml                          ║
║  • Press 'R' to refresh                                    ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
	`, ww.location)

	return ww.styles.ErrorMessage.Render(noData)
}

// getWeatherASCII returns ASCII art for weather condition
func (ww *WeatherWidget) getWeatherASCII(condition string) string {
	condition = strings.ToLower(condition)

	if strings.Contains(condition, "clear") || strings.Contains(condition, "sunny") {
		return "      \\   /\n       .-.\n    ― (   ) ―\n       `-'\n      /   \\"
	}

	if strings.Contains(condition, "partly cloudy") || strings.Contains(condition, "partly") {
		return "      \\  /\n    _ /\"\".-.  \n      \\_(   ). \n      /(___(__)"
	}

	if strings.Contains(condition, "cloudy") || strings.Contains(condition, "overcast") {
		return "\n     .--.\n  .-(    ).\n (___.__)__)"
	}

	if strings.Contains(condition, "rain") || strings.Contains(condition, "shower") {
		return "     .-.\n    (   ).\n   (___(__)\n    ' ' ' '\n   ' ' ' '"
	}

	if strings.Contains(condition, "storm") || strings.Contains(condition, "thunder") {
		return "     .-.\n    (   ).\n   (___(__)\n   ⚡' ⚡'\n    ' ' ' '"
	}

	if strings.Contains(condition, "snow") {
		return "     .-.\n    (   ).\n   (___(__)\n    *  *  *\n   *  *  *"
	}

	if strings.Contains(condition, "fog") || strings.Contains(condition, "mist") {
		return "\n  _ - _ - _ -\n   _ - _ - _\n  _ - _ - _ -\n"
	}

	// Default
	return "     .-.\n    (   ).\n   (___(__)\n"
}

// getConditionIcon returns emoji icon for condition
func (ww *WeatherWidget) getConditionIcon(condition string) string {
	condition = strings.ToLower(condition)

	switch {
	case strings.Contains(condition, "clear"), strings.Contains(condition, "sunny"):
		return "☀️"
	case strings.Contains(condition, "partly cloudy"):
		return "⛅"
	case strings.Contains(condition, "cloudy"), strings.Contains(condition, "overcast"):
		return "☁️"
	case strings.Contains(condition, "rain"), strings.Contains(condition, "shower"):
		return "🌧"
	case strings.Contains(condition, "storm"), strings.Contains(condition, "thunder"):
		return "⛈"
	case strings.Contains(condition, "snow"):
		return "❄️"
	case strings.Contains(condition, "fog"), strings.Contains(condition, "mist"):
		return "🌫"
	case strings.Contains(condition, "wind"):
		return "💨"
	default:
		return "🌤"
	}
}

// formatLastUpdated formats the last update time
func (ww *WeatherWidget) formatLastUpdated() string {
	if ww.lastUpdated.IsZero() {
		return "never"
	}

	diff := time.Since(ww.lastUpdated)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}

// IsExpanded returns whether widget is in expanded view
func (ww *WeatherWidget) IsExpanded() bool {
	return ww.expanded
}

// GetLocation returns the current location
func (ww *WeatherWidget) GetLocation() string {
	return ww.location
}

// SetLocation updates the location
func (ww *WeatherWidget) SetLocation(location string) {
	ww.location = location
	ww.weather = nil // Clear old data
}
