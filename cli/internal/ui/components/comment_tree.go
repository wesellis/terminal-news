package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/wesellis/terminal-news/cli/internal/models"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
)

// CommentTree displays threaded comments
type CommentTree struct {
	comments     []models.CommentTree
	cursor       int
	offset       int
	width        int
	height       int
	styles       *styles.Styles
	collapsed    map[int64]bool // Track collapsed comment threads
	maxDepth     int
}

// NewCommentTree creates a new comment tree
func NewCommentTree(comments []models.CommentWithUser, styles *styles.Styles) *CommentTree {
	return &CommentTree{
		comments:  buildCommentTree(comments),
		cursor:    0,
		offset:    0,
		styles:    styles,
		collapsed: make(map[int64]bool),
		maxDepth:  5, // Max nesting level to display
	}
}

// buildCommentTree organizes flat comments into a tree structure
func buildCommentTree(comments []models.CommentWithUser) []models.CommentTree {
	// Create lookup map
	commentMap := make(map[int64]*models.CommentTree)
	for i := range comments {
		commentMap[comments[i].ID] = &models.CommentTree{
			Comment:  comments[i].Comment,
			Username: comments[i].Username,
			Karma:    comments[i].Karma,
			Depth:    0,
			Children: []models.CommentTree{},
		}
	}

	// Build tree
	var roots []models.CommentTree
	for _, comment := range comments {
		node := commentMap[comment.ID]
		if comment.ParentID == nil {
			// Root comment
			roots = append(roots, *node)
		} else {
			// Child comment
			if parent, ok := commentMap[*comment.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	// Flatten tree for display
	return flattenTree(roots, 0)
}

// flattenTree flattens the tree into a display list with depth info
func flattenTree(comments []models.CommentTree, depth int) []models.CommentTree {
	var result []models.CommentTree

	for _, comment := range comments {
		comment.Depth = depth
		result = append(result, comment)

		if len(comment.Children) > 0 {
			children := flattenTree(comment.Children, depth+1)
			result = append(result, children...)
		}
	}

	return result
}

// SetSize updates dimensions
func (ct *CommentTree) SetSize(width, height int) {
	ct.width = width
	ct.height = height
}

// SetComments updates the comment list
func (ct *CommentTree) SetComments(comments []models.CommentWithUser) {
	ct.comments = buildCommentTree(comments)
	ct.cursor = 0
	ct.offset = 0
}

// Update handles input
func (ct *CommentTree) Update(msg tea.Msg) (*CommentTree, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if ct.cursor > 0 {
				ct.cursor--
				ct.adjustOffset()
			}

		case "down", "j":
			if ct.cursor < len(ct.comments)-1 {
				ct.cursor++
				ct.adjustOffset()
			}

		case "left", "h":
			// Collapse current thread
			if ct.cursor >= 0 && ct.cursor < len(ct.comments) {
				comment := ct.comments[ct.cursor]
				if len(comment.Children) > 0 {
					ct.collapsed[comment.ID] = true
				} else if comment.Depth > 0 {
					// Jump to parent
					ct.jumpToParent()
				}
			}

		case "right", "l":
			// Expand current thread
			if ct.cursor >= 0 && ct.cursor < len(ct.comments) {
				comment := ct.comments[ct.cursor]
				if ct.collapsed[comment.ID] {
					ct.collapsed[comment.ID] = false
				} else if len(comment.Children) > 0 {
					// Jump to first child
					ct.cursor++
					ct.adjustOffset()
				}
			}

		case "space":
			// Toggle collapse
			if ct.cursor >= 0 && ct.cursor < len(ct.comments) {
				comment := ct.comments[ct.cursor]
				ct.collapsed[comment.ID] = !ct.collapsed[comment.ID]
			}
		}
	}

	return ct, nil
}

// jumpToParent jumps cursor to parent comment
func (ct *CommentTree) jumpToParent() {
	if ct.cursor <= 0 || ct.cursor >= len(ct.comments) {
		return
	}

	currentComment := ct.comments[ct.cursor]
	if currentComment.ParentID == nil {
		return
	}

	// Search backwards for parent
	for i := ct.cursor - 1; i >= 0; i-- {
		if ct.comments[i].ID == *currentComment.ParentID {
			ct.cursor = i
			ct.adjustOffset()
			return
		}
	}
}

// adjustOffset keeps cursor visible
func (ct *CommentTree) adjustOffset() {
	visibleLines := ct.height / 4 // Approx lines per comment

	if ct.cursor < ct.offset {
		ct.offset = ct.cursor
	}
	if ct.cursor >= ct.offset+visibleLines {
		ct.offset = ct.cursor - visibleLines + 1
	}
}

// GetSelectedComment returns currently selected comment
func (ct *CommentTree) GetSelectedComment() *models.CommentTree {
	if ct.cursor >= 0 && ct.cursor < len(ct.comments) {
		return &ct.comments[ct.cursor]
	}
	return nil
}

// View renders the comment tree
func (ct *CommentTree) View() string {
	if len(ct.comments) == 0 {
		return ct.renderEmpty()
	}

	var items []string

	// Render visible comments
	for i, comment := range ct.comments {
		if i < ct.offset {
			continue
		}

		if len(items) > ct.height/4 {
			break
		}

		// Skip if parent is collapsed
		if ct.isHidden(comment) {
			continue
		}

		isActive := i == ct.cursor
		item := ct.renderComment(comment, isActive)
		items = append(items, item)
	}

	// Add navigation hint
	hint := "\n" + ct.styles.HelpText.Render("[↑/↓] Navigate  [←/→] Collapse/Expand  [Space] Toggle  [R] Reply  [Esc] Back")
	items = append(items, hint)

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// isHidden checks if comment is hidden due to collapsed parent
func (ct *CommentTree) isHidden(comment models.CommentTree) bool {
	if comment.ParentID == nil {
		return false
	}

	// Check all ancestors for collapsed state
	for _, c := range ct.comments {
		if c.ID == *comment.ParentID {
			if ct.collapsed[c.ID] {
				return true
			}
			return ct.isHidden(c)
		}
	}

	return false
}

// renderComment renders a single comment
func (ct *CommentTree) renderComment(comment models.CommentTree, isActive bool) string {
	var parts []string

	// Calculate indentation
	indent := strings.Repeat("  ", comment.Depth)
	if comment.Depth > ct.maxDepth {
		indent = strings.Repeat("  ", ct.maxDepth)
	}

	// Depth indicator
	depthChar := "│"
	if comment.Depth == 0 {
		depthChar = "┌"
	}

	// Collapse indicator
	collapseChar := " "
	if len(comment.Children) > 0 {
		if ct.collapsed[comment.ID] {
			collapseChar = "+" // Collapsed, has children
		} else {
			collapseChar = "−" // Expanded, has children
		}
	}

	// Header: username, time, votes
	timeStr := ct.styles.FormatTime(comment.CreatedAt)
	votesStr := ct.styles.FormatVotes(comment.Upvotes, comment.Downvotes)

	header := fmt.Sprintf("%s%s %s @%s • %s • %s",
		indent,
		depthChar,
		collapseChar,
		comment.Username,
		timeStr,
		votesStr,
	)

	if comment.IsDeleted {
		header += " [deleted]"
	}

	parts = append(parts, ct.styles.Meta.Render(header))

	// Content
	if !comment.IsDeleted {
		content := comment.Content
		maxWidth := ct.width - (comment.Depth * 2) - 6
		if maxWidth < 40 {
			maxWidth = 40
		}

		wrapped := wordwrap.String(content, maxWidth)
		contentLines := strings.Split(wrapped, "\n")

		for _, line := range contentLines {
			parts = append(parts, indent+"  "+line)
		}

		// Show collapsed children count
		if ct.collapsed[comment.ID] && len(comment.Children) > 0 {
			childCount := ct.countDescendants(comment)
			parts = append(parts, ct.styles.Meta.Render(
				fmt.Sprintf("%s  [%d hidden replies]", indent, childCount),
			))
		}
	} else {
		parts = append(parts, ct.styles.Meta.Render(indent+"  [comment deleted]"))
	}

	// Combine
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)

	// Apply active styling
	if isActive {
		content = ct.styles.ActiveItem.Render(content)
	} else {
		content = ct.styles.Item.Render(content)
	}

	return content
}

// countDescendants counts all descendants of a comment
func (ct *CommentTree) countDescendants(comment models.CommentTree) int {
	count := len(comment.Children)
	for _, child := range comment.Children {
		count += ct.countDescendants(child)
	}
	return count
}

// renderEmpty renders empty state
func (ct *CommentTree) renderEmpty() string {
	empty := `
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║                     No Comments Yet                        ║
║                                                            ║
║            Be the first to share your thoughts!            ║
║                     Press 'C' to comment                   ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
	`

	return ct.styles.HelpText.Render(empty)
}

// GetCommentCount returns total comment count
func (ct *CommentTree) GetCommentCount() int {
	return len(ct.comments)
}

// GetVisibleCount returns count of visible (non-collapsed) comments
func (ct *CommentTree) GetVisibleCount() int {
	count := 0
	for _, comment := range ct.comments {
		if !ct.isHidden(comment) {
			count++
		}
	}
	return count
}
