package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wesellis/terminal-news/cli/internal/api"
	"github.com/wesellis/terminal-news/cli/internal/cache"
	"github.com/wesellis/terminal-news/cli/internal/models"
	"github.com/wesellis/terminal-news/cli/internal/ui/components"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
)

// ArticleDetailView displays an article with its comments
type ArticleDetailView struct {
	BaseView
	article     *models.ArticleWithRanking
	commentTree *components.CommentTree
}

// NewArticleDetailView creates a new article detail view
func NewArticleDetailView(apiClient *api.Client, cache *cache.Cache, styles *styles.Styles) *ArticleDetailView {
	return &ArticleDetailView{
		BaseView: BaseView{
			apiClient: apiClient,
			cache:     cache,
			styles:    styles,
		},
		commentTree: components.NewCommentTree(styles),
	}
}

// LoadArticle loads article and comments
func (v *ArticleDetailView) LoadArticle(articleID int64) tea.Cmd {
	return func() tea.Msg {
		// Load article
		resp, err := v.apiClient.GetArticles("hot", 0, 1) // TODO: Get specific article
		if err != nil {
			return nil
		}

		var article *models.Article
		if len(resp.Articles) > 0 {
			article = &resp.Articles[0]
		}

		// Load comments
		comments, err := v.apiClient.GetComments(articleID)
		if err != nil {
			// Try cache
			comments, _ = v.cache.GetComments(articleID)
		} else {
			// Cache comments
			for _, comment := range comments {
				v.cache.SaveComment(&comment)
			}
		}

		return ArticleDetailLoadedMsg{
			Article:  article,
			Comments: comments,
		}
	}
}

// Update handles messages
func (v *ArticleDetailView) Update(msg tea.Msg) (*ArticleDetailView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ArticleDetailLoadedMsg:
		v.article = msg.Article
		v.commentTree.SetComments(msg.Comments)
		v.loading = false

	case tea.WindowSizeMsg:
		v.SetSize(msg.Width, msg.Height)
		v.commentTree.SetSize(msg.Width, msg.Height-20) // Reserve space for article header

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			// Go back to article list
			return v, func() tea.Msg {
				return BackToListMsg{}
			}
		case "r":
			// Reply to comment
			// TODO: Open reply form
		}
	}

	// Pass message to comment tree
	v.commentTree, cmd = v.commentTree.Update(msg)

	return v, cmd
}

// View renders the article detail view
func (v *ArticleDetailView) View() string {
	if v.loading {
		return v.styles.InfoMessage.Render("\n  ⠋ Loading article...\n")
	}

	if v.article == nil {
		return v.styles.ErrorMessage.Render("\n  Article not found\n")
	}

	var content string

	// Article header
	header := v.styles.Header.Render("📰 ARTICLE")
	content += header + "\n\n"

	// Article details
	commentCount := v.article.TotalEngagement / 5 // Estimate from engagement
	contentPreview := v.article.Content
	if len(contentPreview) > 60 {
		contentPreview = contentPreview[:60] + "..."
	}

	articleBox := fmt.Sprintf(`
┌─ %s ────────────────────────────────────────────────────────────────┐
│                                                                      │
│  %s                                                                  │
│                                                                      │
│  Source: %s  •  %s                                                   │
│  ↑ %d  ↓ %d  •  %d views  •  %d comments                            │
│                                                                      │
│  %s                                                                  │
│                                                                      │
│  URL: %s                                                             │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
`,
		v.article.Source,
		v.styles.Truncate(v.article.Title, 60),
		v.article.Category,
		v.styles.FormatTime(v.article.PublishedAt),
		v.article.LikeCount,
		v.article.DislikeCount,
		v.article.OpenCount,
		commentCount,
		v.styles.Truncate(contentPreview, 60),
		v.styles.Truncate(v.article.URL, 60),
	)

	content += articleBox + "\n\n"

	// Comments section
	commentsHeader := v.styles.Header.Render("💬 COMMENTS")
	content += commentsHeader + "\n\n"
	content += v.commentTree.View()

	// Footer
	footer := v.styles.HelpText.Render("[↑/↓] Navigate [Space] Toggle [←] Collapse [→] Expand [R] Reply [Esc] Back")
	content += "\n\n" + footer

	return content
}

// ArticleDetailLoadedMsg is sent when article and comments are loaded
type ArticleDetailLoadedMsg struct {
	Article  *models.ArticleWithRanking
	Comments []models.CommentWithUser
}

// BackToListMsg is sent when user wants to go back to article list
type BackToListMsg struct{}
