package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
)

// HelpOverlay displays keyboard shortcuts and help information
type HelpOverlay struct {
	visible bool
	width   int
	height  int
	styles  *styles.Styles
}

// NewHelpOverlay creates a new help overlay
func NewHelpOverlay(styles *styles.Styles) *HelpOverlay {
	return &HelpOverlay{
		visible: false,
		styles:  styles,
	}
}

// SetSize updates dimensions
func (ho *HelpOverlay) SetSize(width, height int) {
	ho.width = width
	ho.height = height
}

// Toggle shows/hides the overlay
func (ho *HelpOverlay) Toggle() {
	ho.visible = !ho.visible
}

// Show displays the overlay
func (ho *HelpOverlay) Show() {
	ho.visible = true
}

// Hide hides the overlay
func (ho *HelpOverlay) Hide() {
	ho.visible = false
}

// IsVisible returns visibility state
func (ho *HelpOverlay) IsVisible() bool {
	return ho.visible
}

// Update handles input
func (ho *HelpOverlay) Update(msg tea.Msg) (*HelpOverlay, tea.Cmd) {
	if !ho.visible {
		return ho, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?", "esc", "q":
			ho.Hide()
		}
	}

	return ho, nil
}

// View renders the help overlay
func (ho *HelpOverlay) View() string {
	if !ho.visible {
		return ""
	}

	content := ho.renderHelp()

	// Center the overlay
	return lipgloss.Place(
		ho.width,
		ho.height,
		lipgloss.Center,
		lipgloss.Center,
		ho.styles.Border.Render(content),
	)
}

// renderHelp renders the help content
func (ho *HelpOverlay) renderHelp() string {
	help := `
┌─ TERMINAL NEWS - KEYBOARD SHORTCUTS ────────────────────────────────┐
│                                                                      │
│  NAVIGATION                                                          │
│  ──────────────────────────────────────────────────────────────────  │
│  Tab             Switch to next tab                                  │
│  Shift+Tab       Switch to previous tab                              │
│  1-6             Jump to specific tab (Hot, Controversial, etc)      │
│  ↑/k             Move up                                             │
│  ↓/j             Move down                                           │
│  g               Go to top                                           │
│  G               Go to bottom                                        │
│  PgUp/PgDn       Page up/down                                        │
│                                                                      │
│  ACTIONS                                                             │
│  ──────────────────────────────────────────────────────────────────  │
│  Enter           Select/Open item                                    │
│  o               Open article in browser                             │
│  l               Like article                                        │
│  d               Dislike article                                     │
│  c               View/Post comments                                  │
│  r               Refresh current view                                │
│  Esc             Go back/Cancel                                      │
│                                                                      │
│  COMMENTS (when viewing)                                             │
│  ──────────────────────────────────────────────────────────────────  │
│  ←/h             Collapse thread / Jump to parent                    │
│  →/l             Expand thread / Jump to first child                 │
│  Space           Toggle collapse                                     │
│  R               Reply to comment                                    │
│                                                                      │
│  CLASSIFIEDS                                                         │
│  ──────────────────────────────────────────────────────────────────  │
│  n               Post new classified                                 │
│  e               Edit your classified                                │
│  Del             Delete your classified                              │
│  /               Search classifieds                                  │
│  f               Filter by category                                  │
│                                                                      │
│  GLOBAL                                                              │
│  ──────────────────────────────────────────────────────────────────  │
│  ?               Show/Hide this help                                 │
│  q / Ctrl+C      Quit Terminal News                                  │
│                                                                      │
│  TIPS                                                                │
│  ──────────────────────────────────────────────────────────────────  │
│  • All actions work offline and sync when connected                  │
│  • Use vim-style keys (h/j/k/l) for navigation                       │
│  • Weather data updates every 15 minutes                             │
│  • Press 'r' frequently to see new content                           │
│                                                                      │
│                      Press ? or Esc to close                         │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
	`

	return strings.TrimSpace(help)
}

// GetShortHelp returns a short help string for status bar
func GetShortHelp() string {
	return "[?] Help  [Tab] Switch  [q] Quit  [r] Refresh  [↑↓] Navigate"
}

// GetTabHelp returns help for specific tab
func GetTabHelp(tab string) string {
	switch tab {
	case "hot", "controversial", "rising":
		return "[o] Open  [l] Like  [d] Dislike  [c] Comments  [Enter] Select"

	case "classifieds":
		return "[n] New Post  [/] Search  [f] Filter  [Enter] View Details"

	case "profile":
		return "[e] Edit Profile  [s] Settings  [l] Logout"

	case "weather":
		return "[e] Toggle View  [r] Refresh  [l] Change Location"

	default:
		return GetShortHelp()
	}
}
