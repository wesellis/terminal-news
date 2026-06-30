package monitor

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Monitor struct {
	db *sqlx.DB
}

type SourceStats struct {
	Source       string
	TotalFetched int
	SuccessRate  float64
	AvgFetchTime time.Duration
	LastSuccess  time.Time
	LastError    *string
	QualityScore float64
}

func NewMonitor(db *sqlx.DB) *Monitor {
	return &Monitor{db: db}
}

func (m *Monitor) LogFetch(source string, success bool, fetchTime time.Duration, articlesCount int, err error) {
	query := `
		INSERT INTO fetch_logs (
			source, success, fetch_time_ms, articles_count, error_message, created_at
		) VALUES ($1, $2, $3, $4, $5, NOW())
	`

	var errorMsg *string
	if err != nil {
		msg := err.Error()
		errorMsg = &msg
	}

	_, dbErr := m.db.Exec(query,
		source,
		success,
		fetchTime.Milliseconds(),
		articlesCount,
		errorMsg,
	)

	if dbErr != nil {
		log.Printf("Failed to log fetch: %v", dbErr)
	}
}

func (m *Monitor) GetSourceStats() ([]SourceStats, error) {
	query := `
		SELECT
			source,
			COUNT(*) as total_fetched,
			AVG(CASE WHEN success THEN 1 ELSE 0 END) as success_rate,
			AVG(fetch_time_ms) as avg_fetch_time_ms,
			MAX(CASE WHEN success THEN created_at END) as last_success,
			MAX(CASE WHEN NOT success THEN error_message END) as last_error
		FROM fetch_logs
		WHERE created_at > NOW() - INTERVAL '7 days'
		GROUP BY source
		ORDER BY success_rate DESC
	`

	var stats []SourceStats
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s SourceStats
		var avgFetchMs float64

		err := rows.Scan(
			&s.Source,
			&s.TotalFetched,
			&s.SuccessRate,
			&avgFetchMs,
			&s.LastSuccess,
			&s.LastError,
		)

		if err != nil {
			continue
		}

		s.AvgFetchTime = time.Duration(avgFetchMs) * time.Millisecond

		// Calculate quality score
		s.QualityScore = m.calculateQualityScore(s)

		stats = append(stats, s)
	}

	return stats, nil
}

func (m *Monitor) calculateQualityScore(stats SourceStats) float64 {
	score := 0.0

	// Success rate (40% weight)
	score += stats.SuccessRate * 0.4

	// Speed (20% weight) - faster is better
	if stats.AvgFetchTime < 5*time.Second {
		score += 0.2
	} else if stats.AvgFetchTime < 10*time.Second {
		score += 0.1
	}

	// Consistency (20% weight) - recent success
	if time.Since(stats.LastSuccess) < 1*time.Hour {
		score += 0.2
	} else if time.Since(stats.LastSuccess) < 6*time.Hour {
		score += 0.1
	}

	// Volume (20% weight) - more articles is better
	avgArticles := float64(stats.TotalFetched) / 7.0 // Per day
	if avgArticles > 100 {
		score += 0.2
	} else if avgArticles > 50 {
		score += 0.1
	}

	return score
}

func (m *Monitor) AlertOnFailures() error {
	query := `
		SELECT source, COUNT(*) as fail_count
		FROM fetch_logs
		WHERE success = FALSE
		  AND created_at > NOW() - INTERVAL '1 hour'
		GROUP BY source
		HAVING COUNT(*) > 3
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var source string
		var failCount int

		if err := rows.Scan(&source, &failCount); err != nil {
			continue
		}

		log.Printf("ALERT: Source %s has failed %d times in the last hour", source, failCount)

		// Send alert (email, Slack, etc.)
		// m.sendAlert(source, failCount)
	}

	return nil
}
