package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrCommentNotFound   = errors.New("comment not found")
	ErrCommentDeleted    = errors.New("comment has been deleted")
	ErrUnauthorized      = errors.New("unauthorized to perform this action")
	ErrCommentTooLong    = errors.New("comment exceeds maximum length")
	ErrInvalidParentID   = errors.New("invalid parent comment ID")
)

type CommentService struct {
	db *sqlx.DB
}

func NewCommentService(db *sqlx.DB) *CommentService {
	return &CommentService{db: db}
}

// Comment represents a user comment on an article
type Comment struct {
	ID         int64          `json:"id" db:"id"`
	UserID     int64          `json:"user_id" db:"user_id"`
	ArticleID  int64          `json:"article_id" db:"article_id"`
	ParentID   sql.NullInt64  `json:"parent_id,omitempty" db:"parent_id"`
	Content    string         `json:"content" db:"content"`
	IsDeleted  bool           `json:"is_deleted" db:"is_deleted"`
	IsFlagged  bool           `json:"is_flagged" db:"is_flagged"`
	FlagCount  int            `json:"flag_count" db:"flag_count"`
	Upvotes    int            `json:"upvotes" db:"upvotes"`
	Downvotes  int            `json:"downvotes" db:"downvotes"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	EditedAt   sql.NullTime   `json:"edited_at,omitempty" db:"edited_at"`

	// Joined data
	Username   string         `json:"username,omitempty" db:"username"`
	Karma      int            `json:"user_karma,omitempty" db:"karma"`

	// Nested replies (populated for tree structure)
	Replies    []Comment      `json:"replies,omitempty"`
}

// CreateCommentRequest is the request body for creating a comment
type CreateCommentRequest struct {
	Content  string  `json:"content"`
	ParentID *int64  `json:"parent_id,omitempty"`
}

// UpdateCommentRequest is the request body for updating a comment
type UpdateCommentRequest struct {
	Content string `json:"content"`
}

// CommentListResponse is the response for comment list endpoints
type CommentListResponse struct {
	Comments   []Comment `json:"comments"`
	TotalCount int       `json:"total_count"`
}

// CreateComment creates a new comment on an article
func (s *CommentService) CreateComment(ctx context.Context, userID, articleID int64, content string, parentID *int64) (*Comment, error) {
	// Validate content length
	if len(content) == 0 || len(content) > 10000 {
		return nil, ErrCommentTooLong
	}

	// If parent_id is provided, validate it exists and belongs to same article
	if parentID != nil {
		var parentArticleID int64
		err := s.db.GetContext(ctx, &parentArticleID,
			`SELECT article_id FROM comments WHERE id = $1 AND is_deleted = FALSE`,
			*parentID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrInvalidParentID
			}
			return nil, fmt.Errorf("failed to validate parent comment: %w", err)
		}

		if parentArticleID != articleID {
			return nil, ErrInvalidParentID
		}
	}

	// Insert comment
	query := `
		INSERT INTO comments (user_id, article_id, parent_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, article_id, parent_id, content, is_deleted, is_flagged,
		          flag_count, upvotes, downvotes, created_at, updated_at, edited_at
	`

	var comment Comment
	err := s.db.GetContext(ctx, &comment, query, userID, articleID, parentID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}

// GetArticleComments retrieves all comments for an article in a tree structure
func (s *CommentService) GetArticleComments(ctx context.Context, articleID int64) (*CommentListResponse, error) {
	// Get all comments for the article
	query := `
		SELECT
			c.id, c.user_id, c.article_id, c.parent_id, c.content,
			c.is_deleted, c.is_flagged, c.flag_count, c.upvotes, c.downvotes,
			c.created_at, c.updated_at, c.edited_at,
			u.username, u.karma
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.article_id = $1
		ORDER BY c.created_at ASC
	`

	var allComments []Comment
	if err := s.db.SelectContext(ctx, &allComments, query, articleID); err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	// Build tree structure
	commentMap := make(map[int64]*Comment)
	var rootComments []Comment

	// First pass: create map of all comments
	for i := range allComments {
		comment := allComments[i]
		commentMap[comment.ID] = &comment
	}

	// Second pass: build tree
	for i := range allComments {
		comment := &allComments[i]

		if comment.ParentID.Valid {
			// This is a reply
			if parent, ok := commentMap[comment.ParentID.Int64]; ok {
				parent.Replies = append(parent.Replies, *comment)
			}
		} else {
			// This is a root comment
			if replies, ok := commentMap[comment.ID]; ok {
				comment.Replies = replies.Replies
			}
			rootComments = append(rootComments, *comment)
		}
	}

	return &CommentListResponse{
		Comments:   rootComments,
		TotalCount: len(allComments),
	}, nil
}

// GetComment retrieves a single comment by ID
func (s *CommentService) GetComment(ctx context.Context, commentID int64) (*Comment, error) {
	query := `
		SELECT
			c.id, c.user_id, c.article_id, c.parent_id, c.content,
			c.is_deleted, c.is_flagged, c.flag_count, c.upvotes, c.downvotes,
			c.created_at, c.updated_at, c.edited_at,
			u.username, u.karma
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`

	var comment Comment
	if err := s.db.GetContext(ctx, &comment, query, commentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return &comment, nil
}

// UpdateComment updates a comment's content
func (s *CommentService) UpdateComment(ctx context.Context, commentID, userID int64, content string) (*Comment, error) {
	// Validate content length
	if len(content) == 0 || len(content) > 10000 {
		return nil, ErrCommentTooLong
	}

	// Check if user owns the comment
	var ownerID int64
	err := s.db.GetContext(ctx, &ownerID, `SELECT user_id FROM comments WHERE id = $1`, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to check comment ownership: %w", err)
	}

	if ownerID != userID {
		return nil, ErrUnauthorized
	}

	// Update comment
	query := `
		UPDATE comments
		SET content = $1, edited_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND is_deleted = FALSE
		RETURNING id, user_id, article_id, parent_id, content, is_deleted, is_flagged,
		          flag_count, upvotes, downvotes, created_at, updated_at, edited_at
	`

	var comment Comment
	if err := s.db.GetContext(ctx, &comment, query, content, commentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCommentDeleted
		}
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &comment, nil
}

// DeleteComment soft-deletes a comment
func (s *CommentService) DeleteComment(ctx context.Context, commentID, userID int64) error {
	// Check if user owns the comment
	var ownerID int64
	err := s.db.GetContext(ctx, &ownerID, `SELECT user_id FROM comments WHERE id = $1`, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCommentNotFound
		}
		return fmt.Errorf("failed to check comment ownership: %w", err)
	}

	if ownerID != userID {
		return ErrUnauthorized
	}

	// Soft delete (set is_deleted = true, clear content)
	query := `
		UPDATE comments
		SET is_deleted = TRUE, content = '[deleted]', updated_at = NOW()
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCommentNotFound
	}

	return nil
}

// VoteComment updates upvote/downvote counts on a comment
func (s *CommentService) VoteComment(ctx context.Context, commentID int64, isUpvote bool) error {
	var query string
	if isUpvote {
		query = `UPDATE comments SET upvotes = upvotes + 1, updated_at = NOW() WHERE id = $1`
	} else {
		query = `UPDATE comments SET downvotes = downvotes + 1, updated_at = NOW() WHERE id = $1`
	}

	result, err := s.db.ExecContext(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("failed to vote on comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCommentNotFound
	}

	return nil
}

// GetUserComments retrieves all comments by a specific user
func (s *CommentService) GetUserComments(ctx context.Context, userID int64, limit, offset int) ([]Comment, error) {
	query := `
		SELECT
			c.id, c.user_id, c.article_id, c.parent_id, c.content,
			c.is_deleted, c.is_flagged, c.flag_count, c.upvotes, c.downvotes,
			c.created_at, c.updated_at, c.edited_at,
			u.username, u.karma
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.user_id = $1 AND c.is_deleted = FALSE
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var comments []Comment
	if err := s.db.SelectContext(ctx, &comments, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get user comments: %w", err)
	}

	return comments, nil
}
