package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/wesellis/terminal-news/cli/internal/models"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
)

// ArticleList displays a list of articles with navigation
type ArticleList struct {
	articles       []models.ArticleWithRanking
	cursor         int
	selectedID     int64
	offset         int
	visibleItems   int
	width          int
	height         int
	styles         *styles.Styles
	showSummaries  bool
	compactMode    bool
}

// NewArticleList creates a new article list component
func NewArticleList(articles []models.ArticleWithRanking, styles *styles.Styles) *ArticleList {
	return &ArticleList{
		articles:      articles,
		cursor:        0,
		selectedID:    -1,
		offset:        0,
		visibleItems:  10,
		styles:        styles,
		showSummaries: true,
		compactMode:   false,
	}
}

// SetSize updates the dimensions
func (al *ArticleList) SetSize(width, height int) {
	al.width = width
	al.height = height
	// Calculate how many items can fit
	itemHeight := 6 // Default item height with summary
	if al.compactMode {
		itemHeight = 2
	}
	al.visibleItems = height / itemHeight
	if al.visibleItems < 1 {
		al.visibleItems = 1
	}
}

// SetArticles updates the article list
func (al *ArticleList) SetArticles(articles []models.ArticleWithRanking) {
	al.articles = articles
	if al.cursor >= len(articles) {
		al.cursor = len(articles) - 1
	}
	if al.cursor < 0 {
		al.cursor = 0
	}
}

// Update handles keyboard input
func (al *ArticleList) Update(msg tea.Msg) (*ArticleList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			al.MovePrev()
		case "down", "j":
			al.MoveNext()
		case "g": // Go to top
			al.cursor = 0
			al.offset = 0
		case "G": // Go to bottom
			al.cursor = len(al.articles) - 1
			al.adjustOffset()
		case "pageup":
			al.cursor -= al.visibleItems
			if al.cursor < 0 {
				al.cursor = 0
			}
			al.adjustOffset()
		case "pagedown":
			al.cursor += al.visibleItems
			if al.cursor >= len(al.articles) {
				al.cursor = len(al.articles) - 1
			}
			al.adjustOffset()
		case "home":
			al.cursor = 0
			al.offset = 0
		case "end":
			al.cursor = len(al.articles) - 1
			al.adjustOffset()
		case "l": // Like article
			article := al.GetSelectedArticle()
			if article != nil {
				return al, func() tea.Msg {
					return VoteArticleMsg{ArticleID: article.ID, VoteType: "like"}
				}
			}
		case "d": // Dislike article
			article := al.GetSelectedArticle()
			if article != nil {
				return al, func() tea.Msg {
					return VoteArticleMsg{ArticleID: article.ID, VoteType: "dislike"}
				}
			}
		case "c": // View comments
			article := al.GetSelectedArticle()
			if article != nil {
				return al, func() tea.Msg {
					return NavigateToCommentsMsg{
						ArticleID:    article.ID,
						ArticleTitle: article.Title,
					}
				}
			}
		case "o", "enter": // Open article
			article := al.GetSelectedArticle()
			if article != nil {
				return al, func() tea.Msg {
					return NavigateToArticleMsg{
						ArticleID:  article.ID,
						ArticleURL: article.URL,
					}
				}
			}
		}
	}

	return al, nil
}

// MovePrev moves cursor up
func (al *ArticleList) MovePrev() {
	if al.cursor > 0 {
		al.cursor--
		al.adjustOffset()
	}
}

// MoveNext moves cursor down
func (al *ArticleList) MoveNext() {
	if al.cursor < len(al.articles)-1 {
		al.cursor++
		al.adjustOffset()
	}
}

// adjustOffset adjusts the scroll offset to keep cursor visible
func (al *ArticleList) adjustOffset() {
	if al.cursor < al.offset {
		al.offset = al.cursor
	}
	if al.cursor >= al.offset+al.visibleItems {
		al.offset = al.cursor - al.visibleItems + 1
	}
}

// GetSelectedArticle returns the currently selected article
func (al *ArticleList) GetSelectedArticle() *models.ArticleWithRanking {
	if al.cursor >= 0 && al.cursor < len(al.articles) {
		return &al.articles[al.cursor]
	}
	return nil
}

// View renders the article list
func (al *ArticleList) View() string {
	if len(al.articles) == 0 {
		return al.renderEmpty()
	}

	var items []string

	// Calculate visible range
	start := al.offset
	end := al.offset + al.visibleItems
	if end > len(al.articles) {
		end = len(al.articles)
	}

	// Render visible articles
	for i := start; i < end; i++ {
		article := al.articles[i]
		isActive := i == al.cursor
		item := al.renderArticleItem(article, isActive, i+1)
		items = append(items, item)
	}

	// Add scroll indicator if needed
	if len(al.articles) > al.visibleItems {
		scrollInfo := fmt.Sprintf("\n[%d-%d of %d articles]", start+1, end, len(al.articles))
		items = append(items, al.styles.HelpText.Render(scrollInfo))
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderArticleItem renders a single article
func (al *ArticleList) renderArticleItem(article models.ArticleWithRanking, isActive bool, index int) string {
	if al.compactMode {
		return al.renderCompactItem(article, isActive, index)
	}
	return al.renderFullItem(article, isActive, index)
}

// renderFullItem renders article with full details
func (al *ArticleList) renderFullItem(article models.ArticleWithRanking, isActive bool, index int) string {
	var parts []string

	// Line 1: Index, votes, views, time, source
	votesStr := al.styles.FormatVotes(article.LikeCount, article.DislikeCount)
	score := article.LikeCount - article.DislikeCount
	scoreStyle := al.styles.Success
	if score < 0 {
		scoreStyle = al.styles.Error
	}
	scoreStr := scoreStyle.Render(fmt.Sprintf("%+d", score))

	viewsStr := al.styles.Meta.Render(fmt.Sprintf("👁 %d", article.OpenCount))
	timeStr := al.styles.Meta.Render(al.styles.FormatTime(article.PublishedAt))
	sourceStr := al.styles.Meta.Render(fmt.Sprintf("via %s", article.Source))

	metaLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		fmt.Sprintf("#%d ", index),
		votesStr,
		"  ",
		scoreStr,
		"  ",
		viewsStr,
		"  ",
		timeStr,
		"  ",
		sourceStr,
	)
	parts = append(parts, metaLine)

	// Line 2: Title with indicators
	title := article.Title
	if isActive {
		title = al.styles.SelectedTitle.Render("▶ " + title)
	} else {
		title = al.styles.Title.Render("  " + title)
	}

	// Add hot/rising indicators based on ranking
	indicators := ""
	if article.HotRank > 0.7 { // High hot rank = hot article
		indicators += " 🔥"
	}
	if article.HoursSincePublished < 6 && article.TotalEngagement > 10 { // Recent + engaging = rising
		indicators += " ⚡"
	}
	if indicators != "" {
		title = title + indicators
	}

	parts = append(parts, title)

	// Line 3: Summary (if enabled) - use Content field
	if al.showSummaries && article.Content != "" {
		summary := article.Content
		maxWidth := al.width - 6
		if maxWidth < 40 {
			maxWidth = 40
		}
		// Truncate to reasonable length for summary
		if len(summary) > 200 {
			summary = summary[:200] + "..."
		}
		wrapped := wordwrap.String(summary, maxWidth)
		summaryLines := strings.Split(wrapped, "\n")

		// Show first 2 lines max
		displayLines := summaryLines
		if len(summaryLines) > 2 {
			displayLines = summaryLines[:2]
			displayLines[1] = displayLines[1] + "..."
		}

		for _, line := range displayLines {
			parts = append(parts, al.styles.Meta.Render("  "+line))
		}
	}

	// Line 4: Actions
	actions := []string{}

	// Comment count is part of total engagement (estimate as 20% of engagement)
	commentCount := article.TotalEngagement / 5
	if commentCount > 0 {
		actions = append(actions, fmt.Sprintf("💬 %d", commentCount))
	} else {
		actions = append(actions, "💬 0")
	}

	actions = append(actions, "[O]pen", "[L]ike", "[D]islike", "[C]omments")

	actionsLine := al.styles.Actions.Render("  " + strings.Join(actions, "  "))
	parts = append(parts, actionsLine)

	// Combine all parts
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)

	// Apply active item styling
	if isActive {
		return al.styles.ActiveItem.Render(content)
	}

	return al.styles.Item.Render(content)
}

// renderCompactItem renders minimal article view
func (al *ArticleList) renderCompactItem(article models.ArticleWithRanking, isActive bool, index int) string {
	// Format: #123 [+42] Title (5h ago, 12 comments)
	score := article.LikeCount - article.DislikeCount
	scoreStr := fmt.Sprintf("[%+d]", score)

	timeStr := al.styles.FormatTime(article.PublishedAt)
	commentCount := article.TotalEngagement / 5 // Estimate
	commentsStr := fmt.Sprintf("%d💬", commentCount)

	prefix := fmt.Sprintf("#%d ", index)
	suffix := fmt.Sprintf(" (%s, %s)", timeStr, commentsStr)

	// Calculate max title length
	maxTitleLen := al.width - len(prefix) - len(scoreStr) - len(suffix) - 10
	if maxTitleLen < 20 {
		maxTitleLen = 20
	}

	title := article.Title
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}

	line := prefix + scoreStr + " " + title + suffix

	if isActive {
		return al.styles.ActiveItem.Render("▶ " + line)
	}

	return al.styles.Item.Render("  " + line)
}

// renderEmpty renders empty state
func (al *ArticleList) renderEmpty() string {
	emptyMsg := `
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║                    No Articles Available                   ║
║                                                            ║
║  • Check your internet connection                          ║
║  • Press 'r' to refresh                                    ║
║  • Try switching to offline mode (--offline flag)          ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
	`

	return al.styles.HelpText.Render(emptyMsg)
}

// SetCompactMode toggles compact display mode
func (al *ArticleList) SetCompactMode(compact bool) {
	al.compactMode = compact
}

// SetShowSummaries toggles summary display
func (al *ArticleList) SetShowSummaries(show bool) {
	al.showSummaries = show
}

// GetCursor returns current cursor position
func (al *ArticleList) GetCursor() int {
	return al.cursor
}

// GetOffset returns current scroll offset
func (al *ArticleList) GetOffset() int {
	return al.offset
}

// VoteArticleMsg is sent when an article is voted on
type VoteArticleMsg struct {
	ArticleID int64
	VoteType  string // "like" or "dislike"
}

// NavigateToCommentsMsg is sent when user wants to view comments
type NavigateToCommentsMsg struct {
	ArticleID    int64
	ArticleTitle string
}

// NavigateToArticleMsg is sent when user wants to open article
type NavigateToArticleMsg struct {
	ArticleID  int64
	ArticleURL string
}
