package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wesellis/terminal-news/cli/internal/api"
	"github.com/wesellis/terminal-news/cli/internal/cache"
	"github.com/wesellis/terminal-news/cli/internal/config"
	"github.com/wesellis/terminal-news/cli/internal/models"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
	"github.com/wesellis/terminal-news/cli/internal/ui/views"
)

// Tab represents different application views
type Tab int

const (
	HotTab Tab = iota
	ControversialTab
	RisingTab
	ProfileTab
	WeatherTab
	ClassifiedsTab
)

// App is the main application model
type App struct {
	// Configuration
	config *config.Config

	// Navigation
	activeTab Tab
	tabs      []string

	// Views
	hotView           *views.HotView
	controversialView *views.ControversialView
	risingView        *views.RisingView
	profileView       *views.ProfileView
	weatherView       *views.WeatherView
	classifiedsView   *views.ClassifiedsView

	// State
	width         int
	height        int
	authenticated bool
	username      string
	ready         bool
	quitting      bool

	// Services
	apiClient *api.Client
	wsClient  *api.WebSocketClient
	cache     *cache.Cache

	// UI
	styles *styles.Styles

	// Messages
	statusMessage string
	errorMessage  string
}

// NewApp creates and initializes a new application
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize cache
	cacheDB, err := cache.New(cfg.Cache.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Initialize API client
	apiClient := api.NewClient(cfg.API.BaseURL, cfg.API.Timeout)

	// Load auth token if exists
	if cfg.User.Token != "" {
		apiClient.SetAuthToken(cfg.User.Token)
	}

	app := &App{
		config:        cfg,
		tabs:          []string{"Hot", "Controversial", "Rising", "Profile", "Weather", "Classifieds"},
		activeTab:     getDefaultTab(cfg.User.DefaultTab),
		apiClient:     apiClient,
		cache:         cacheDB,
		styles:        styles.DefaultStyles(),
		authenticated: cfg.User.Token != "",
		username:      cfg.User.Username,
	}

	return app, nil
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	// Initialize all views
	a.hotView = views.NewHotView(a.apiClient, a.cache, a.styles)
	a.controversialView = views.NewControversialView(a.apiClient, a.cache, a.styles)
	a.risingView = views.NewRisingView(a.apiClient, a.cache, a.styles)
	a.profileView = views.NewProfileView(a.apiClient, a.cache, a.styles)
	a.weatherView = views.NewWeatherView(a.apiClient, a.cache, a.styles, a.config.User.Location)
	a.classifiedsView = views.NewClassifiedsView(a.apiClient, a.cache, a.styles)

	// Connect WebSocket if not in offline mode
	if !a.config.Offline {
		wsClient, err := api.NewWebSocketClient(a.config.API.WebSocketURL)
		if err == nil {
			a.wsClient = wsClient
		}
	}

	// Batch initial commands
	cmds := []tea.Cmd{
		tea.EnterAltScreen,
	}

	// Load initial data based on active tab
	switch a.activeTab {
	case HotTab:
		cmds = append(cmds, a.hotView.LoadArticles())
	case ControversialTab:
		cmds = append(cmds, a.controversialView.LoadArticles())
	case RisingTab:
		cmds = append(cmds, a.risingView.LoadArticles())
	case WeatherTab:
		cmds = append(cmds, a.weatherView.LoadWeather())
	case ProfileTab:
		if a.authenticated {
			cmds = append(cmds, a.profileView.LoadProfile())
		}
	case ClassifiedsTab:
		cmds = append(cmds, a.classifiedsView.LoadClassifieds())
	}

	// Always load weather widget data
	cmds = append(cmds, a.weatherView.LoadWeather())

	// Start listening for WebSocket events
	if a.wsClient != nil {
		cmds = append(cmds, a.listenForWebSocketEvents())
	}

	a.ready = true
	return tea.Batch(cmds...)
}

// Update handles all messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !a.ready {
			return a, nil
		}

		// Global keyboard shortcuts
		switch msg.String() {
		case "ctrl+c", "q":
			a.quitting = true
			return a, tea.Quit

		case "tab":
			a.nextTab()
			cmds = append(cmds, a.loadTabData())
			return a, tea.Batch(cmds...)

		case "shift+tab":
			a.prevTab()
			cmds = append(cmds, a.loadTabData())
			return a, tea.Batch(cmds...)

		case "r":
			a.statusMessage = "Refreshing..."
			cmds = append(cmds, a.refresh())
			return a, tea.Batch(cmds...)

		case "?":
			return a, a.showHelp()

		case "1":
			a.activeTab = HotTab
			cmds = append(cmds, a.loadTabData())
			return a, tea.Batch(cmds...)

		case "2":
			a.activeTab = ControversialTab
			cmds = append(cmds, a.loadTabData())
			return a, tea.Batch(cmds...)

		case "3":
			a.activeTab = RisingTab
			cmds = append(cmds, a.loadTabData())
			return a, tea.Batch(cmds...)

		case "4":
			a.activeTab = ProfileTab
			cmds = append(cmds, a.loadTabData())
			return a, tea.Batch(cmds...)

		case "5":
			a.activeTab = WeatherTab
			cmds = append(cmds, a.loadTabData())
			return a, tea.Batch(cmds...)

		case "6":
			a.activeTab = ClassifiedsTab
			cmds = append(cmds, a.loadTabData())
			return a, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		// Update all views with new size
		contentHeight := a.height - 10 // Account for header, weather, status bar

		if a.hotView != nil {
			a.hotView.SetSize(msg.Width, contentHeight)
		}
		if a.controversialView != nil {
			a.controversialView.SetSize(msg.Width, contentHeight)
		}
		if a.risingView != nil {
			a.risingView.SetSize(msg.Width, contentHeight)
		}
		if a.profileView != nil {
			a.profileView.SetSize(msg.Width, contentHeight)
		}
		if a.weatherView != nil {
			a.weatherView.SetSize(msg.Width, contentHeight)
		}
		if a.classifiedsView != nil {
			a.classifiedsView.SetSize(msg.Width, contentHeight)
		}

		return a, nil

	case StatusMsg:
		a.statusMessage = string(msg)
		return a, nil

	case ErrorMsg:
		a.errorMessage = string(msg)
		return a, nil

	case models.WSMessage:
		return a, a.handleWebSocketMessage(msg)
	}

	// Forward messages to active view
	var cmd tea.Cmd
	switch a.activeTab {
	case HotTab:
		if a.hotView != nil {
			a.hotView, cmd = a.hotView.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ControversialTab:
		if a.controversialView != nil {
			a.controversialView, cmd = a.controversialView.Update(msg)
			cmds = append(cmds, cmd)
		}
	case RisingTab:
		if a.risingView != nil {
			a.risingView, cmd = a.risingView.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ProfileTab:
		if a.profileView != nil {
			a.profileView, cmd = a.profileView.Update(msg)
			cmds = append(cmds, cmd)
		}
	case WeatherTab:
		if a.weatherView != nil {
			a.weatherView, cmd = a.weatherView.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ClassifiedsTab:
		if a.classifiedsView != nil {
			a.classifiedsView, cmd = a.classifiedsView.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return a, tea.Batch(cmds...)
}

// View renders the application
func (a *App) View() string {
	if !a.ready {
		return "\n  Loading Terminal News..."
	}

	if a.quitting {
		return "\n  Thanks for reading! 📰\n"
	}

	var content string

	// Render header with tabs
	header := a.renderHeader()

	// Render compact weather widget (always visible)
	weather := ""
	if a.weatherView != nil && a.activeTab != WeatherTab {
		weather = a.weatherView.CompactView()
	}

	// Render active view
	switch a.activeTab {
	case HotTab:
		if a.hotView != nil {
			content = a.hotView.View()
		}
	case ControversialTab:
		if a.controversialView != nil {
			content = a.controversialView.View()
		}
	case RisingTab:
		if a.risingView != nil {
			content = a.risingView.View()
		}
	case ProfileTab:
		if a.profileView != nil {
			content = a.profileView.View()
		}
	case WeatherTab:
		if a.weatherView != nil {
			content = a.weatherView.View()
		}
	case ClassifiedsTab:
		if a.classifiedsView != nil {
			content = a.classifiedsView.View()
		}
	}

	// Render status bar
	statusBar := a.renderStatusBar()

	// Combine all sections
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		weather,
		content,
		statusBar,
	)
}

// renderHeader creates the header with tabs and branding
func (a *App) renderHeader() string {
	title := a.styles.Header.Render("🗞  TERMINAL NEWS  •  AM Radio for the Information Age")

	var tabViews []string
	for i, tabName := range a.tabs {
		if Tab(i) == a.activeTab {
			tabViews = append(tabViews, a.styles.ActiveTab.Render(tabName))
		} else {
			tabViews = append(tabViews, a.styles.Tab.Render(tabName))
		}
	}
	tabs := lipgloss.JoinHorizontal(lipgloss.Left, tabViews...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		tabs,
	)
}

// renderStatusBar creates the bottom status bar
func (a *App) renderStatusBar() string {
	// Left side: status message or help hint
	leftStatus := ""
	if a.statusMessage != "" {
		leftStatus = a.styles.InfoMessage.Render(a.statusMessage)
	} else if a.errorMessage != "" {
		leftStatus = a.styles.ErrorMessage.Render(a.errorMessage)
	} else {
		leftStatus = a.styles.HelpText.Render("Press ? for help • Tab to switch • q to quit")
	}

	// Right side: connection status and user
	rightStatus := ""
	if a.config.Offline {
		rightStatus = a.styles.Warning.Render("OFFLINE")
	} else if a.authenticated {
		rightStatus = a.styles.Success.Render(fmt.Sprintf("@%s", a.username))
	} else {
		rightStatus = a.styles.Muted.Render("Not logged in")
	}

	// Calculate spacing
	gap := a.width - lipgloss.Width(leftStatus) - lipgloss.Width(rightStatus) - 4
	if gap < 0 {
		gap = 0
	}

	statusLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftStatus,
		strings.Repeat(" ", gap),
		rightStatus,
	)

	return a.styles.StatusBar.Render(statusLine)
}

// Navigation helpers

func (a *App) nextTab() {
	a.activeTab = (a.activeTab + 1) % Tab(len(a.tabs))
	a.statusMessage = ""
	a.errorMessage = ""
}

func (a *App) prevTab() {
	a.activeTab = (a.activeTab - 1 + Tab(len(a.tabs))) % Tab(len(a.tabs))
	a.statusMessage = ""
	a.errorMessage = ""
}

func (a *App) loadTabData() tea.Cmd {
	switch a.activeTab {
	case HotTab:
		return a.hotView.LoadArticles()
	case ControversialTab:
		return a.controversialView.LoadArticles()
	case RisingTab:
		return a.risingView.LoadArticles()
	case ProfileTab:
		if a.authenticated {
			return a.profileView.LoadProfile()
		}
	case WeatherTab:
		return a.weatherView.LoadWeather()
	case ClassifiedsTab:
		return a.classifiedsView.LoadClassifieds()
	}
	return nil
}

func (a *App) refresh() tea.Cmd {
	return a.loadTabData()
}

func (a *App) showHelp() tea.Cmd {
	return func() tea.Msg {
		// Show help overlay (implement later)
		return StatusMsg("Help: q=quit, Tab=switch tabs, r=refresh, 1-6=jump to tab")
	}
}

func (a *App) listenForWebSocketEvents() tea.Cmd {
	if a.wsClient == nil {
		return nil
	}

	return func() tea.Msg {
		for event := range a.wsClient.Events() {
			return event
		}
		return nil
	}
}

func (a *App) handleWebSocketMessage(msg models.WSMessage) tea.Cmd {
	switch msg.Type {
	case "new_article":
		// Update appropriate feeds
		a.statusMessage = "New article available"
	case "vote_update":
		// Update vote counts in views
	case "new_comment":
		// Update comment counts
	}
	return nil
}

// Helper function to get default tab from config
func getDefaultTab(tabName string) Tab {
	switch tabName {
	case "hot":
		return HotTab
	case "controversial":
		return ControversialTab
	case "rising":
		return RisingTab
	case "profile":
		return ProfileTab
	case "weather":
		return WeatherTab
	case "classifieds":
		return ClassifiedsTab
	default:
		return HotTab
	}
}

// Custom message types
type StatusMsg string
type ErrorMsg string
