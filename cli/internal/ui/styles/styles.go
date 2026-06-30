package styles

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Styles contains all visual styles for the application
type Styles struct {
	// Colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Muted     lipgloss.Color
	Hot       lipgloss.Color
	Rising    lipgloss.Color

	// Layout
	Header      lipgloss.Style
	Tab         lipgloss.Style
	ActiveTab   lipgloss.Style
	StatusBar   lipgloss.Style
	Container   lipgloss.Style
	Border      lipgloss.Style

	// Articles
	Item           lipgloss.Style
	ActiveItem     lipgloss.Style
	Title          lipgloss.Style
	SelectedTitle  lipgloss.Style
	Meta           lipgloss.Style
	Actions        lipgloss.Style
	VoteUp         lipgloss.Style
	VoteDown       lipgloss.Style

	// Components
	WeatherWidget  lipgloss.Style
	HelpText       lipgloss.Style

	// Forms
	Input       lipgloss.Style
	ActiveInput lipgloss.Style
	Label       lipgloss.Style
	Button      lipgloss.Style
	ActiveButton lipgloss.Style

	// Messages
	InfoMessage    lipgloss.Style
	ErrorMessage   lipgloss.Style
	SuccessMessage lipgloss.Style
}

// DefaultStyles returns the default styling configuration
func DefaultStyles() *Styles {
	return &Styles{
		// Brand colors
		Primary:   lipgloss.Color("#7D56F4"),
		Secondary: lipgloss.Color("#04B575"),
		Success:   lipgloss.Color("#10B981"),
		Warning:   lipgloss.Color("#F59E0B"),
		Error:     lipgloss.Color("#EF4444"),
		Muted:     lipgloss.Color("#6B7280"),
		Hot:       lipgloss.Color("#FF6B6B"),
		Rising:    lipgloss.Color("#FFD93D"),

		// Header and navigation
		Header: lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 1).
			Foreground(lipgloss.Color("#FFFFFF")),

		Tab: lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#9CA3AF")),

		ActiveTab: lipgloss.NewStyle().
			Padding(0, 2).
			Bold(true).
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF")),

		StatusBar: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(0, 1).
			Foreground(lipgloss.Color("#9CA3AF")),

		Container: lipgloss.NewStyle().
			Padding(1, 2),

		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")),

		// Article list items
		Item: lipgloss.NewStyle().
			Padding(1, 2).
			MarginBottom(1),

		ActiveItem: lipgloss.NewStyle().
			Padding(1, 2).
			MarginBottom(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#1F2937")),

		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true),

		SelectedTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Underline(true),

		Meta: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true),

		Actions: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")),

		VoteUp: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true),

		VoteDown: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true),

		// Weather widget
		WeatherWidget: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(1, 2).
			MarginBottom(1),

		HelpText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true),

		// Form elements
		Input: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(0, 1).
			Foreground(lipgloss.Color("#FFFFFF")),

		ActiveInput: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Foreground(lipgloss.Color("#FFFFFF")),

		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Bold(true).
			MarginBottom(1),

		Button: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(0, 2).
			Foreground(lipgloss.Color("#9CA3AF")),

		ActiveButton: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true),

		// Status messages
		InfoMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6")).
			Bold(true),

		ErrorMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true),

		SuccessMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true),
	}
}

// Helper functions for consistent spacing and formatting

func (s *Styles) FormatVotes(upvotes, downvotes int) string {
	up := s.VoteUp.Render(fmt.Sprintf("▲ %d", upvotes))
	down := s.VoteDown.Render(fmt.Sprintf("▼ %d", downvotes))
	return lipgloss.JoinHorizontal(lipgloss.Left, up, "  ", down)
}

func (s *Styles) FormatTime(t time.Time) string {
	diff := time.Since(t)

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

func (s *Styles) Truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}
