package main

import (
	"fmt"
	"log"

	"github.com/wesellis/terminal-news/scraper/internal/classifier"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

func main() {
	log.Println("=== Classifier Test (No Database) ===")

	c := classifier.NewClassifier()

	// Test articles from different categories
	testArticles := []types.ParsedArticle{
		{
			Title:   "Apple Announces New iPhone with AI Features",
			Summary: "Apple unveiled the latest iPhone featuring advanced artificial intelligence capabilities and machine learning",
		},
		{
			Title:   "Stock Market Hits Record High",
			Summary: "The stock market reached an all-time high today as investors reacted positively to strong earnings reports from major corporations",
		},
		{
			Title:   "Scientists Discover New Species in Amazon",
			Summary: "Researchers have discovered a new species of frog in the Amazon rainforest during an ecology expedition",
		},
		{
			Title:   "Lakers Win Championship Game",
			Summary: "The Lakers defeated their opponents in an exciting basketball game to win the NBA championship",
		},
		{
			Title:   "New Movie Breaks Box Office Records",
			Summary: "The latest Hollywood film starring famous actors has broken box office records in its opening weekend",
		},
		{
			Title:   "Biden Announces New Policy on Immigration",
			Summary: "President Biden announced new legislation regarding immigration policy in congress today",
		},
	}

	log.Printf("Testing classification on %d articles\n", len(testArticles))

	correctCount := 0
	expected := []string{"tech", "business", "science", "sports", "entertainment", "politics"}

	for i, article := range testArticles {
		c.Classify(&article)

		match := "❌"
		if article.Category == expected[i] {
			match = "✅"
			correctCount++
		}

		fmt.Printf("%s Article: %s\n", match, article.Title)
		fmt.Printf("   Classified as: %s (expected: %s)\n", article.Category, expected[i])
		if len(article.Tags) > 0 {
			fmt.Printf("   Tags: %v\n", article.Tags)
		}
		fmt.Println()
	}

	accuracy := float64(correctCount) / float64(len(testArticles)) * 100
	log.Printf("=== Results ===")
	log.Printf("Correct: %d/%d", correctCount, len(testArticles))
	log.Printf("Accuracy: %.1f%%", accuracy)

	if accuracy >= 80.0 {
		log.Println("✅ Classifier meets 80% accuracy target!")
	} else {
		log.Printf("⚠️ Classifier below 80%% target (current: %.1f%%)", accuracy)
	}
}
