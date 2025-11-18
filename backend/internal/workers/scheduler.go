package workers

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// Scheduler handles background jobs and scheduled tasks
type Scheduler struct {
	db  *sqlx.DB
	rdb *redis.Client
}

// NewScheduler creates a new scheduler instance
func NewScheduler(db *sqlx.DB, rdb *redis.Client) *Scheduler {
	return &Scheduler{
		db:  db,
		rdb: rdb,
	}
}

// Start starts all scheduled jobs
func (s *Scheduler) Start(ctx context.Context) {
	log.Println("🕐 Starting background scheduler")

	// Refresh article rankings every 5 minutes
	go s.runEvery(ctx, 5*time.Minute, "Refresh article rankings", s.RefreshRankings)

	// Expire old classifieds every hour
	go s.runEvery(ctx, 1*time.Hour, "Expire old classifieds", s.ExpireClassifieds)

	// Clean up old votes (older than 7 days from deleted articles) daily
	go s.runEvery(ctx, 24*time.Hour, "Clean up old votes", s.CleanupOldVotes)

	// Downgrade expired premium classifieds every hour
	go s.runEvery(ctx, 1*time.Hour, "Downgrade expired premiums", s.DowngradeExpiredPremiums)

	log.Println("✅ Background scheduler running")
}

// runEvery runs a function at regular intervals
func (s *Scheduler) runEvery(ctx context.Context, interval time.Duration, name string, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately on start
	if err := fn(); err != nil {
		log.Printf("❌ Error in %s: %v", name, err)
	} else {
		log.Printf("✅ %s completed", name)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("🛑 Stopping scheduler: %s", name)
			return
		case <-ticker.C:
			if err := fn(); err != nil {
				log.Printf("❌ Error in %s: %v", name, err)
			} else {
				log.Printf("✅ %s completed", name)
			}
		}
	}
}

// RefreshRankings refreshes the materialized view for article rankings
func (s *Scheduler) RefreshRankings() error {
	start := time.Now()

	// Refresh materialized view concurrently (non-blocking)
	_, err := s.db.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY article_rankings`)
	if err != nil {
		return err
	}

	// Invalidate Redis caches for article listings
	keys := []string{
		"articles:hot:*",
		"articles:controversial:*",
		"articles:rising:*",
	}

	for _, pattern := range keys {
		iter := s.rdb.Scan(context.Background(), 0, pattern, 0).Iterator()
		for iter.Next(context.Background()) {
			s.rdb.Del(context.Background(), iter.Val())
		}
		if err := iter.Err(); err != nil {
			log.Printf("Warning: Failed to clear cache pattern %s: %v", pattern, err)
		}
	}

	log.Printf("📊 Ranking refresh took %v", time.Since(start))
	return nil
}

// ExpireClassifieds sets expired classifieds to inactive
func (s *Scheduler) ExpireClassifieds() error {
	query := `
		UPDATE classifieds
		SET is_active = FALSE, updated_at = NOW()
		WHERE expires_at < NOW()
		  AND is_active = TRUE
	`

	result, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("🗑️  Expired %d classifieds", rowsAffected)
	}

	return nil
}

// DowngradeExpiredPremiums downgrades premium classifieds that have expired
func (s *Scheduler) DowngradeExpiredPremiums() error {
	query := `
		UPDATE classifieds
		SET is_premium = FALSE, updated_at = NOW()
		WHERE premium_until < NOW()
		  AND is_premium = TRUE
	`

	result, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("⬇️  Downgraded %d premium classifieds", rowsAffected)
	}

	return nil
}

// CleanupOldVotes removes votes from articles older than 7 days
func (s *Scheduler) CleanupOldVotes() error {
	query := `
		DELETE FROM votes
		WHERE created_at < NOW() - INTERVAL '7 days'
		  AND article_id NOT IN (
		    SELECT id FROM articles WHERE published_at > NOW() - INTERVAL '7 days'
		  )
	`

	result, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("🧹 Cleaned up %d old votes", rowsAffected)
	}

	return nil
}
