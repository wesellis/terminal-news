package views

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wesellis/terminal-news/cli/internal/api"
	"github.com/wesellis/terminal-news/cli/internal/cache"
	"github.com/wesellis/terminal-news/cli/internal/models"
	"github.com/wesellis/terminal-news/cli/internal/ui/components"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
)

// BaseView contains common functionality for all views
type BaseView struct {
	apiClient *api.Client
	cache     *cache.Cache
	styles    *styles.Styles
	width     int
	height    int
	loading   bool
	error     string
}

// SetSize updates the view dimensions
func (v *BaseView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// ArticlesLoadedMsg is sent when articles are loaded
type ArticlesLoadedMsg struct {
	Articles []models.ArticleWithRanking
	Feed     string
}

// CommentsLoadedMsg is sent when comments are loaded
type CommentsLoadedMsg struct {
	Comments []models.CommentWithUser
}

// ClassifiedsLoadedMsg is sent when classifieds are loaded
type ClassifiedsLoadedMsg struct {
	Classifieds []models.Classified
}

// WeatherLoadedMsg is sent when weather is loaded
type WeatherLoadedMsg struct {
	Weather *models.Weather
}

// ProfileLoadedMsg is sent when profile is loaded
type ProfileLoadedMsg struct {
	User       *models.User
	Activity   []models.Activity
	Classifieds []models.Classified
}

// HotView displays the hot news feed
type HotView struct {
	BaseView
	articleList *components.ArticleList
}

// NewHotView creates a new hot view
func NewHotView(apiClient *api.Client, cache *cache.Cache, styles *styles.Styles) *HotView {
	return &HotView{
		BaseView: BaseView{
			apiClient: apiClient,
			cache:     cache,
			styles:    styles,
		},
		articleList: components.NewArticleList(styles),
	}
}

// LoadArticles loads articles from API or cache
func (v *HotView) LoadArticles() tea.Cmd {
	return func() tea.Msg {
		// Try API first
		resp, err := v.apiClient.GetArticles("hot", 0, 50)
		if err == nil {
			// Cache the articles
			v.cache.SaveArticles(resp.Articles)
			return ArticlesLoadedMsg{Articles: resp.Articles, Feed: "hot"}
		}

		// Fallback to cache
		articles, _ := v.cache.GetArticles("hot", 50)
		return ArticlesLoadedMsg{Articles: articles, Feed: "hot"}
	}
}

// Update handles messages
func (v *HotView) Update(msg tea.Msg) (*HotView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ArticlesLoadedMsg:
		if msg.Feed == "hot" {
			v.articleList.SetArticles(msg.Articles)
			v.loading = false
		}
	case tea.WindowSizeMsg:
		v.SetSize(msg.Width, msg.Height)
		v.articleList.SetSize(msg.Width, msg.Height-10) // Reserve space for header/footer
	case components.VoteArticleMsg:
		// Handle voting
		return v, func() tea.Msg {
			err := v.apiClient.VoteArticle(msg.ArticleID, msg.VoteType)
			if err != nil {
				// Queue for offline sync
				v.cache.QueueAction("vote", map[string]interface{}{
					"article_id": msg.ArticleID,
					"type":       msg.VoteType,
				})
			}
			// Reload articles to get updated vote counts
			return v.LoadArticles()()
		}
	}

	// Pass message to article list
	v.articleList, cmd = v.articleList.Update(msg)

	return v, cmd
}

// View renders the hot view
func (v *HotView) View() string {
	if v.loading {
		return v.styles.InfoMessage.Render("\n  ⠋ Loading hot articles...\n")
	}

	header := v.styles.Header.Render("🔥 HOT ARTICLES")
	footer := v.styles.HelpText.Render("[↑/↓] Navigate [Enter] Open [L] Like [D] Dislike [C] Comments [R] Refresh")

	return header + "\n\n" + v.articleList.View() + "\n\n" + footer
}

// ControversialView displays controversial articles
type ControversialView struct {
	BaseView
	articleList *components.ArticleList
}

// NewControversialView creates a new controversial view
func NewControversialView(apiClient *api.Client, cache *cache.Cache, styles *styles.Styles) *ControversialView {
	return &ControversialView{
		BaseView: BaseView{
			apiClient: apiClient,
			cache:     cache,
			styles:    styles,
		},
		articleList: components.NewArticleList(styles),
	}
}

func (v *ControversialView) LoadArticles() tea.Cmd {
	return func() tea.Msg {
		resp, err := v.apiClient.GetArticles("controversial", 0, 50)
		if err == nil {
			return ArticlesLoadedMsg{Articles: resp.Articles, Feed: "controversial"}
		}
		articles, _ := v.cache.GetArticles("controversial", 50)
		return ArticlesLoadedMsg{Articles: articles, Feed: "controversial"}
	}
}

func (v *ControversialView) Update(msg tea.Msg) (*ControversialView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ArticlesLoadedMsg:
		if msg.Feed == "controversial" {
			v.articleList.SetArticles(msg.Articles)
			v.loading = false
		}
	case tea.WindowSizeMsg:
		v.SetSize(msg.Width, msg.Height)
		v.articleList.SetSize(msg.Width, msg.Height-10)
	case components.VoteArticleMsg:
		return v, func() tea.Msg {
			err := v.apiClient.VoteArticle(msg.ArticleID, msg.VoteType)
			if err != nil {
				v.cache.QueueAction("vote", map[string]interface{}{
					"article_id": msg.ArticleID,
					"type":       msg.VoteType,
				})
			}
			return v.LoadArticles()()
		}
	}

	v.articleList, cmd = v.articleList.Update(msg)
	return v, cmd
}

func (v *ControversialView) View() string {
	if v.loading {
		return v.styles.InfoMessage.Render("\n  ⠋ Loading controversial articles...\n")
	}

	header := v.styles.Header.Render("⚡ CONTROVERSIAL ARTICLES")
	footer := v.styles.HelpText.Render("[↑/↓] Navigate [Enter] Open [L] Like [D] Dislike [C] Comments [R] Refresh")

	return header + "\n\n" + v.articleList.View() + "\n\n" + footer
}

// RisingView displays rising articles
type RisingView struct {
	BaseView
	articleList *components.ArticleList
}

func NewRisingView(apiClient *api.Client, cache *cache.Cache, styles *styles.Styles) *RisingView {
	return &RisingView{
		BaseView: BaseView{
			apiClient: apiClient,
			cache:     cache,
			styles:    styles,
		},
		articleList: components.NewArticleList(styles),
	}
}

func (v *RisingView) LoadArticles() tea.Cmd {
	return func() tea.Msg {
		resp, err := v.apiClient.GetArticles("rising", 0, 50)
		if err == nil {
			return ArticlesLoadedMsg{Articles: resp.Articles, Feed: "rising"}
		}
		articles, _ := v.cache.GetArticles("rising", 50)
		return ArticlesLoadedMsg{Articles: articles, Feed: "rising"}
	}
}

func (v *RisingView) Update(msg tea.Msg) (*RisingView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ArticlesLoadedMsg:
		if msg.Feed == "rising" {
			v.articleList.SetArticles(msg.Articles)
			v.loading = false
		}
	case tea.WindowSizeMsg:
		v.SetSize(msg.Width, msg.Height)
		v.articleList.SetSize(msg.Width, msg.Height-10)
	case components.VoteArticleMsg:
		return v, func() tea.Msg {
			err := v.apiClient.VoteArticle(msg.ArticleID, msg.VoteType)
			if err != nil {
				v.cache.QueueAction("vote", map[string]interface{}{
					"article_id": msg.ArticleID,
					"type":       msg.VoteType,
				})
			}
			return v.LoadArticles()()
		}
	}

	v.articleList, cmd = v.articleList.Update(msg)
	return v, cmd
}

func (v *RisingView) View() string {
	if v.loading {
		return v.styles.InfoMessage.Render("\n  ⠋ Loading rising articles...\n")
	}

	header := v.styles.Header.Render("📈 RISING ARTICLES")
	footer := v.styles.HelpText.Render("[↑/↓] Navigate [Enter] Open [L] Like [D] Dislike [C] Comments [R] Refresh")

	return header + "\n\n" + v.articleList.View() + "\n\n" + footer
}

// ProfileView displays user profile
type ProfileView struct {
	BaseView
	user        *models.User
	activity    []models.Activity
	classifieds []models.Classified
	activeTab   int
	tabs        []string
}

func NewProfileView(apiClient *api.Client, cache *cache.Cache, styles *styles.Styles) *ProfileView {
	return &ProfileView{
		BaseView: BaseView{
			apiClient: apiClient,
			cache:     cache,
			styles:    styles,
		},
		tabs: []string{"Stats", "Activity", "Classifieds", "Settings"},
		activeTab: 0,
	}
}

func (v *ProfileView) LoadProfile() tea.Cmd {
	return func() tea.Msg {
		user, err := v.apiClient.GetProfile()
		if err != nil {
			return nil
		}

		// Load activity
		activity, _ := v.apiClient.GetUserActivity(user.ID)

		// Load user's classifieds
		classifieds, _ := v.apiClient.GetUserClassifieds(user.ID)

		return ProfileLoadedMsg{
			User:        user,
			Activity:    activity,
			Classifieds: classifieds,
		}
	}
}

func (v *ProfileView) Update(msg tea.Msg) (*ProfileView, tea.Cmd) {
	switch msg := msg.(type) {
	case ProfileLoadedMsg:
		v.user = msg.User
		v.activity = msg.Activity
		v.classifieds = msg.Classifieds
		v.loading = false
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if v.activeTab > 0 {
				v.activeTab--
			}
		case "right", "l":
			if v.activeTab < len(v.tabs)-1 {
				v.activeTab++
			}
		}
	}
	return v, nil
}

func (v *ProfileView) View() string {
	if v.loading {
		return v.styles.InfoMessage.Render("\n  ⠋ Loading profile...\n")
	}

	if v.user == nil {
		return v.styles.ErrorMessage.Render("\n  Please login first (press 'l')\n")
	}

	var content string

	// Header
	header := v.styles.Header.Render("👤 PROFILE - " + v.user.Username)
	content += header + "\n\n"

	// Tab bar
	var tabItems []string
	for i, tab := range v.tabs {
		if i == v.activeTab {
			tabItems = append(tabItems, v.styles.ActiveTab.Render(tab))
		} else {
			tabItems = append(tabItems, v.styles.Tab.Render(tab))
		}
	}
	content += v.styles.Tabs.Render(tabItems...) + "\n\n"

	// Tab content
	switch v.activeTab {
	case 0: // Stats
		content += v.renderStats()
	case 1: // Activity
		content += v.renderActivity()
	case 2: // Classifieds
		content += v.renderClassifieds()
	case 3: // Settings
		content += v.renderSettings()
	}

	// Footer
	footer := v.styles.HelpText.Render("[←/→] Switch Tabs [L] Logout [E] Edit Profile")
	content += "\n\n" + footer

	return content
}

func (v *ProfileView) renderStats() string {
	if v.user == nil {
		return ""
	}

	stats := fmt.Sprintf(`
╔═══════════════════════════════════════════════════════════════════════╗
║  USER STATISTICS                                                      ║
╠═══════════════════════════════════════════════════════════════════════╣
║                                                                       ║
║  Karma Score:           %6d                                          ║
║  Articles Posted:       %6d                                          ║
║  Comments Posted:       %6d                                          ║
║  Classifieds Posted:    %6d                                          ║
║  Votes Cast:            %6d                                          ║
║                                                                       ║
║  Member Since:          %s                                           ║
║  Last Active:           %s                                           ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝
	`,
		v.user.Karma,
		v.user.ArticleCount,
		v.user.CommentCount,
		len(v.classifieds),
		v.user.VoteCount,
		v.styles.FormatTime(v.user.CreatedAt),
		v.styles.FormatTime(v.user.LastActive),
	)

	return stats
}

func (v *ProfileView) renderActivity() string {
	if len(v.activity) == 0 {
		return "No recent activity"
	}

	var content string
	content += "RECENT ACTIVITY\n\n"

	for i, act := range v.activity {
		if i >= 10 {
			break
		}
		icon := "•"
		switch act.Type {
		case "comment":
			icon = "💬"
		case "vote":
			icon = "👍"
		case "article":
			icon = "📰"
		case "classified":
			icon = "📋"
		}

		content += fmt.Sprintf("  %s %s - %s\n",
			icon,
			v.styles.FormatTime(act.CreatedAt),
			v.styles.Truncate(act.Description, 60),
		)
	}

	return content
}

func (v *ProfileView) renderClassifieds() string {
	if len(v.classifieds) == 0 {
		return "No classifieds posted yet"
	}

	var content string
	content += "YOUR CLASSIFIEDS\n\n"

	for i, classified := range v.classifieds {
		if i >= 10 {
			break
		}

		premiumBadge := ""
		if classified.IsPremium {
			premiumBadge = " ⭐"
		}

		status := "Active"
		if classified.ExpiresAt.Before(time.Now()) {
			status = "Expired"
		}

		content += fmt.Sprintf("  %s%s\n    $%s • %s • %s\n\n",
			classified.Title,
			premiumBadge,
			v.styles.FormatPrice(classified.Price),
			classified.Category,
			status,
		)
	}

	return content
}

func (v *ProfileView) renderSettings() string {
	settings := `
SETTINGS

[ ] Email Notifications
[ ] Desktop Notifications
[✓] Show Compact View
[✓] Offline Mode
[ ] Hide Controversial Content

Location: San Francisco, CA
Theme: Default
Language: English

[Press 'e' to edit settings]
	`

	return settings
}

// WeatherView displays weather information
type WeatherView struct {
	BaseView
	widget *components.WeatherWidget
}

func NewWeatherView(apiClient *api.Client, cache *cache.Cache, styles *styles.Styles, location string) *WeatherView {
	return &WeatherView{
		BaseView: BaseView{
			apiClient: apiClient,
			cache:     cache,
			styles:    styles,
		},
		widget: components.NewWeatherWidget(styles, location),
	}
}

func (v *WeatherView) LoadWeather() tea.Cmd {
	return func() tea.Msg {
		weather, err := v.apiClient.GetWeather(v.widget.GetLocation())
		if err == nil {
			v.cache.SaveWeather(v.widget.GetLocation(), weather)
			return WeatherLoadedMsg{Weather: weather}
		}
		// Try cache
		cached, _ := v.cache.GetWeather(v.widget.GetLocation(), 15*60) // 15 min TTL
		if cached != nil {
			return WeatherLoadedMsg{Weather: cached}
		}
		return nil
	}
}

func (v *WeatherView) Update(msg tea.Msg) (*WeatherView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case WeatherLoadedMsg:
		v.widget.SetWeather(msg.Weather)
		v.loading = false
	case tea.WindowSizeMsg:
		v.SetSize(msg.Width, msg.Height)
		v.widget.SetSize(msg.Width, msg.Height)
	}

	v.widget, cmd = v.widget.Update(msg)
	return v, cmd
}

func (v *WeatherView) View() string {
	if v.loading {
		return v.styles.InfoMessage.Render("\n  ⠋ Loading weather data...\n")
	}

	return v.widget.View()
}

func (v *WeatherView) CompactView() string {
	return v.widget.CompactView()
}

// ClassifiedsView displays classifieds
type ClassifiedsView struct {
	BaseView
	classifieds    []models.Classified
	cursor         int
	showingForm    bool
	classifiedForm *components.ClassifiedForm
	selectedCategory string
	categories     []string
}

func NewClassifiedsView(apiClient *api.Client, cache *cache.Cache, styles *styles.Styles) *ClassifiedsView {
	return &ClassifiedsView{
		BaseView: BaseView{
			apiClient: apiClient,
			cache:     cache,
			styles:    styles,
		},
		classifieds:    []models.Classified{},
		cursor:         0,
		showingForm:    false,
		classifiedForm: components.NewClassifiedForm(styles),
		selectedCategory: "All",
		categories:     []string{"All", "For Sale", "Jobs", "Housing", "Services", "Events", "Gigs"},
	}
}

func (v *ClassifiedsView) LoadClassifieds() tea.Cmd {
	return func() tea.Msg {
		category := ""
		if v.selectedCategory != "All" {
			category = v.selectedCategory
		}
		resp, err := v.apiClient.GetClassifieds(category, "", 0, 50)
		if err == nil {
			v.cache.SaveClassifieds(resp.Classifieds)
			return ClassifiedsLoadedMsg{Classifieds: resp.Classifieds}
		}
		// Try cache
		cached, _ := v.cache.GetClassifieds(category, "", 50)
		if cached != nil {
			return ClassifiedsLoadedMsg{Classifieds: cached}
		}
		return nil
	}
}

func (v *ClassifiedsView) Update(msg tea.Msg) (*ClassifiedsView, tea.Cmd) {
	var cmd tea.Cmd

	// If showing form, delegate to form
	if v.showingForm {
		switch msg := msg.(type) {
		case components.ClassifiedSubmitMsg:
			// Submit classified via API
			v.showingForm = false
			return v, func() tea.Msg {
				_, err := v.apiClient.PostClassified(msg.Classified)
				if err == nil {
					return v.LoadClassifieds()()
				}
				return nil
			}
		case components.ClassifiedCancelMsg:
			v.showingForm = false
			return v, nil
		default:
			v.classifiedForm, cmd = v.classifiedForm.Update(msg)
			return v, cmd
		}
	}

	// Normal browsing mode
	switch msg := msg.(type) {
	case ClassifiedsLoadedMsg:
		v.classifieds = msg.Classifieds
		v.loading = false
		v.cursor = 0
	case tea.WindowSizeMsg:
		v.SetSize(msg.Width, msg.Height)
		v.classifiedForm.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			if v.cursor < len(v.classifieds)-1 {
				v.cursor++
			}
		case "n":
			// New classified
			v.showingForm = true
			return v, nil
		case "f":
			// Filter by category - cycle through
			for i, cat := range v.categories {
				if cat == v.selectedCategory {
					v.selectedCategory = v.categories[(i+1)%len(v.categories)]
					return v, v.LoadClassifieds()
				}
			}
		}
	}

	return v, cmd
}

func (v *ClassifiedsView) View() string {
	if v.showingForm {
		return v.classifiedForm.View()
	}

	if v.loading {
		return v.styles.InfoMessage.Render("\n  ⠋ Loading classifieds...\n")
	}

	var content string

	// Header
	header := v.styles.Header.Render("📋 CLASSIFIEDS - " + v.selectedCategory)
	content += header + "\n\n"

	// Classifieds list
	if len(v.classifieds) == 0 {
		content += v.styles.ErrorMessage.Render("No classifieds found. Press 'n' to post one!")
	} else {
		for i, classified := range v.classifieds {
			prefix := "  "
			if i == v.cursor {
				prefix = "▶ "
			}

			itemStyle := v.styles.Article
			if i == v.cursor {
				itemStyle = v.styles.ActiveArticle
			}

			premiumBadge := ""
			if classified.IsPremium {
				premiumBadge = " ⭐"
			}

			line := itemStyle.Render(
				prefix + classified.Title + premiumBadge + "\n" +
				"    " + classified.Category + " • $" + v.styles.FormatPrice(classified.Price) +
				" • " + classified.City + ", " + classified.State + "\n" +
				"    " + v.styles.Truncate(classified.Description, 80) + "\n",
			)
			content += line + "\n"
		}
	}

	// Footer
	footer := v.styles.HelpText.Render("[↑/↓] Navigate [N] New Post [F] Filter [Enter] View Details [R] Refresh")
	content += "\n" + footer

	return content
}
