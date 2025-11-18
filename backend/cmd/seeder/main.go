package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to database")

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// Seed users
	log.Println("Seeding users...")
	userIDs := seedUsers(db)
	log.Printf("✅ Seeded %d users\n", len(userIDs))

	// Seed articles
	log.Println("Seeding articles...")
	articleIDs := seedArticles(db)
	log.Printf("✅ Seeded %d articles\n", len(articleIDs))

	// Seed votes
	log.Println("Seeding votes...")
	voteCount := seedVotes(db, userIDs, articleIDs)
	log.Printf("✅ Seeded %d votes\n", voteCount)

	// Seed comments
	log.Println("Seeding comments...")
	commentCount := seedComments(db, userIDs, articleIDs)
	log.Printf("✅ Seeded %d comments\n", commentCount)

	// Seed classifieds
	log.Println("Seeding classifieds...")
	classifiedCount := seedClassifieds(db, userIDs)
	log.Printf("✅ Seeded %d classifieds\n", classifiedCount)

	log.Println("\n🎉 Database seeding complete!")
	log.Println("\nTest credentials:")
	log.Println("  Username: testuser / Password: password123")
	log.Println("  Username: alice / Password: password123")
	log.Println("  Username: bob / Password: password123")
}

func seedUsers(db *sqlx.DB) []int64 {
	users := []struct {
		username string
		email    string
	}{
		{"testuser", "test@example.com"},
		{"alice", "alice@example.com"},
		{"bob", "bob@example.com"},
		{"charlie", "charlie@example.com"},
		{"diana", "diana@example.com"},
	}

	// Hash password once
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var userIDs []int64
	for _, u := range users {
		var userID int64
		err := db.Get(&userID, `
			INSERT INTO users (username, email, password_hash, account_tier, created_at)
			VALUES ($1, $2, $3, 'free', NOW())
			ON CONFLICT (username) DO UPDATE SET email = EXCLUDED.email
			RETURNING id
		`, u.username, u.email, string(passwordHash))

		if err != nil {
			log.Printf("Failed to insert user %s: %v", u.username, err)
			continue
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

func seedArticles(db *sqlx.DB) []int64 {
	sources := []string{"TechCrunch", "Reuters", "BBC", "The Verge", "Ars Technica", "Bloomberg", "NPR"}
	categories := []string{"tech", "business", "science", "world", "entertainment", "sports"}

	titles := []string{
		"Major Tech Company Announces New AI Breakthrough",
		"Global Markets React to Economic Policy Changes",
		"Scientists Discover New Species in Deep Ocean",
		"International Summit Addresses Climate Change",
		"New Study Reveals Health Benefits of Exercise",
		"Startup Raises $100M Series B Funding Round",
		"Research Team Makes Progress on Quantum Computing",
		"City Council Approves New Infrastructure Project",
		"Breaking: Major Political Development Unfolds",
		"Technology Giants Form New Industry Alliance",
	}

	var articleIDs []int64

	for i := 0; i < 100; i++ {
		title := titles[rand.Intn(len(titles))]
		if i >= len(titles) {
			title = fmt.Sprintf("%s (Part %d)", title, i/len(titles)+1)
		}

		source := sources[rand.Intn(len(sources))]
		category := categories[rand.Intn(len(categories))]
		publishedAt := time.Now().Add(-time.Duration(rand.Intn(72)) * time.Hour)

		var articleID int64
		err := db.Get(&articleID, `
			INSERT INTO articles (
				title, url, source, summary, category,
				published_at, fetched_at, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
			RETURNING id
		`,
			title,
			fmt.Sprintf("https://example.com/article/%d", i),
			source,
			fmt.Sprintf("This is a summary of article #%d about %s. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer nec odio. Praesent libero. Sed cursus ante dapibus diam.", i, category),
			category,
			publishedAt,
		)

		if err != nil {
			log.Printf("Failed to insert article %d: %v", i, err)
			continue
		}
		articleIDs = append(articleIDs, articleID)
	}

	return articleIDs
}

func seedVotes(db *sqlx.DB, userIDs, articleIDs []int64) int {
	count := 0

	// Each user votes on random articles
	for _, userID := range userIDs {
		// Vote on 20-40 random articles
		numVotes := 20 + rand.Intn(20)

		for i := 0; i < numVotes && i < len(articleIDs); i++ {
			articleID := articleIDs[rand.Intn(len(articleIDs))]

			// 70% like, 20% open, 10% dislike
			voteType := "like"
			randNum := rand.Float64()
			if randNum < 0.1 {
				voteType = "dislike"
			} else if randNum < 0.3 {
				voteType = "open"
			}

			_, err := db.Exec(`
				INSERT INTO votes (user_id, article_id, vote_type, created_at)
				VALUES ($1, $2, $3, NOW())
				ON CONFLICT (user_id, article_id, vote_type) DO NOTHING
			`, userID, articleID, voteType)

			if err != nil {
				log.Printf("Failed to insert vote: %v", err)
				continue
			}
			count++
		}
	}

	return count
}

func seedComments(db *sqlx.DB, userIDs, articleIDs []int64) int {
	commentTexts := []string{
		"Great article! Very informative.",
		"Thanks for sharing this.",
		"I disagree with some of the points made here.",
		"Can you provide more sources?",
		"This is exactly what I was looking for.",
		"Interesting perspective on this topic.",
		"Well written and researched.",
		"I'm not sure I understand the conclusion.",
		"This needs more context.",
		"Excellent analysis!",
	}

	count := 0

	// Add comments to random articles
	for _, articleID := range articleIDs {
		// 50% chance of having comments
		if rand.Float64() < 0.5 {
			continue
		}

		// 1-5 comments per article
		numComments := 1 + rand.Intn(5)

		var topLevelCommentIDs []int64

		for i := 0; i < numComments; i++ {
			userID := userIDs[rand.Intn(len(userIDs))]
			content := commentTexts[rand.Intn(len(commentTexts))]

			var commentID int64
			err := db.Get(&commentID, `
				INSERT INTO comments (user_id, article_id, content, created_at)
				VALUES ($1, $2, $3, NOW())
				RETURNING id
			`, userID, articleID, content)

			if err != nil {
				log.Printf("Failed to insert comment: %v", err)
				continue
			}

			topLevelCommentIDs = append(topLevelCommentIDs, commentID)
			count++
		}

		// Add some replies (nested comments)
		if len(topLevelCommentIDs) > 0 {
			numReplies := rand.Intn(3)
			for i := 0; i < numReplies; i++ {
				userID := userIDs[rand.Intn(len(userIDs))]
				parentID := topLevelCommentIDs[rand.Intn(len(topLevelCommentIDs))]
				content := "Reply: " + commentTexts[rand.Intn(len(commentTexts))]

				_, err := db.Exec(`
					INSERT INTO comments (user_id, article_id, parent_id, content, created_at)
					VALUES ($1, $2, $3, $4, NOW())
				`, userID, articleID, parentID, content)

				if err != nil {
					log.Printf("Failed to insert reply: %v", err)
					continue
				}
				count++
			}
		}
	}

	return count
}

func seedClassifieds(db *sqlx.DB, userIDs []int64) int {
	cities := []struct {
		city  string
		state string
		lat   float64
		lng   float64
	}{
		{"Portland", "OR", 45.5152, -122.6784},
		{"Seattle", "WA", 47.6062, -122.3321},
		{"San Francisco", "CA", 37.7749, -122.4194},
		{"Austin", "TX", 30.2672, -97.7431},
		{"Denver", "CO", 39.7392, -104.9903},
	}

	categories := []string{"for-sale", "housing", "jobs", "services", "community"}

	titles := []string{
		"Vintage bicycle for sale - excellent condition",
		"Looking for roommate in downtown area",
		"Software Developer position available",
		"Professional photography services",
		"Community garden plot available",
		"Couch for sale - must pick up",
		"1BR apartment for rent",
		"Freelance web design services",
		"Guitar lessons - all levels",
		"Desk and office chair - like new",
	}

	count := 0

	for i := 0; i < 50; i++ {
		userID := userIDs[rand.Intn(len(userIDs))]
		title := titles[i%len(titles)]
		category := categories[rand.Intn(len(categories))]
		city := cities[rand.Intn(len(cities))]

		var price *float64
		if category == "for-sale" || category == "housing" {
			p := float64(rand.Intn(1000) + 50)
			price = &p
		}

		_, err := db.Exec(`
			INSERT INTO classifieds (
				user_id, title, description, price, category,
				city, state, country, lat, lng,
				contact_email, contact_method, status,
				expires_at, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, 'US', $8, $9, $10, 'email', 'active', NOW() + INTERVAL '30 days', NOW())
		`,
			userID,
			title,
			fmt.Sprintf("This is a detailed description of %s. Great condition, must see!", title),
			price,
			category,
			city.city,
			city.state,
			city.lat,
			city.lng,
			"seller@example.com",
		)

		if err != nil {
			log.Printf("Failed to insert classified: %v", err)
			continue
		}
		count++
	}

	return count
}
