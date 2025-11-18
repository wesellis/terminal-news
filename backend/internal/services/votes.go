package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidVoteType = errors.New("invalid vote type")
	ErrVoteNotFound    = errors.New("vote not found")
)

type VoteService struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewVoteService(db *sqlx.DB, rdb *redis.Client) *VoteService {
	return &VoteService{db: db, rdb: rdb}
}

// Vote represents a user's vote on an article
type Vote struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	ArticleID int64     `json:"article_id" db:"article_id"`
	VoteType  string    `json:"vote_type" db:"vote_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// VoteRequest is the request body for voting on an article
type VoteRequest struct {
	VoteType string `json:"vote_type"`
}

// VoteResponse is the response after voting
type VoteResponse struct {
	Vote    *Vote `json:"vote"`
	Message string `json:"message"`
}

// CreateVote creates a vote for an article
// Vote types: "open", "like", "dislike"
func (s *VoteService) CreateVote(ctx context.Context, userID, articleID int64, voteType string) (*VoteResponse, error) {
	// Validate vote type
	if voteType != "open" && voteType != "like" && voteType != "dislike" {
		return nil, ErrInvalidVoteType
	}

	// Insert vote (ON CONFLICT DO NOTHING to handle duplicate votes)
	query := `
		INSERT INTO votes (user_id, article_id, vote_type)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, article_id, vote_type) DO NOTHING
		RETURNING id, user_id, article_id, vote_type, created_at
	`

	var vote Vote
	err := s.db.GetContext(ctx, &vote, query, userID, articleID, voteType)
	if err != nil {
		// Check if no rows returned (already voted)
		if err.Error() == "sql: no rows in result set" {
			return &VoteResponse{
				Vote:    nil,
				Message: "Already voted with this type",
			}, nil
		}
		return nil, fmt.Errorf("failed to create vote: %w", err)
	}

	// Invalidate cache for this article
	cachePattern := fmt.Sprintf("article*:%d", articleID)
	s.rdb.Del(ctx, cachePattern)
	// Invalidate article list caches
	s.rdb.Del(ctx, "articles:hot:*", "articles:controversial:*", "articles:rising:*")

	return &VoteResponse{
		Vote:    &vote,
		Message: "Vote recorded successfully",
	}, nil
}

// RemoveVote removes a user's vote from an article
func (s *VoteService) RemoveVote(ctx context.Context, userID, articleID int64, voteType string) error {
	// Validate vote type
	if voteType != "open" && voteType != "like" && voteType != "dislike" {
		return ErrInvalidVoteType
	}

	// Delete vote
	query := `
		DELETE FROM votes
		WHERE user_id = $1 AND article_id = $2 AND vote_type = $3
	`

	result, err := s.db.ExecContext(ctx, query, userID, articleID, voteType)
	if err != nil {
		return fmt.Errorf("failed to remove vote: %w", err)
	}

	// Check if vote was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVoteNotFound
	}

	// Invalidate cache for this article
	cachePattern := fmt.Sprintf("article*:%d", articleID)
	s.rdb.Del(ctx, cachePattern)
	// Invalidate article list caches
	s.rdb.Del(ctx, "articles:hot:*", "articles:controversial:*", "articles:rising:*")

	return nil
}

// GetUserVotes retrieves all votes for a specific user
func (s *VoteService) GetUserVotes(ctx context.Context, userID int64, limit, offset int) ([]Vote, error) {
	query := `
		SELECT id, user_id, article_id, vote_type, created_at
		FROM votes
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var votes []Vote
	if err := s.db.SelectContext(ctx, &votes, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get user votes: %w", err)
	}

	return votes, nil
}

// GetArticleVotes retrieves all votes for a specific article
func (s *VoteService) GetArticleVotes(ctx context.Context, articleID int64) (map[string]int, error) {
	query := `
		SELECT
			vote_type,
			COUNT(*) as count
		FROM votes
		WHERE article_id = $1
		GROUP BY vote_type
	`

	rows, err := s.db.QueryContext(ctx, query, articleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get article votes: %w", err)
	}
	defer rows.Close()

	voteCounts := make(map[string]int)
	for rows.Next() {
		var voteType string
		var count int
		if err := rows.Scan(&voteType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan vote count: %w", err)
		}
		voteCounts[voteType] = count
	}

	return voteCounts, nil
}

// GetUserVoteForArticle checks if user has voted on a specific article
func (s *VoteService) GetUserVoteForArticle(ctx context.Context, userID, articleID int64) ([]string, error) {
	query := `
		SELECT vote_type
		FROM votes
		WHERE user_id = $1 AND article_id = $2
	`

	var voteTypes []string
	if err := s.db.SelectContext(ctx, &voteTypes, query, userID, articleID); err != nil {
		return nil, fmt.Errorf("failed to get user vote: %w", err)
	}

	return voteTypes, nil
}
