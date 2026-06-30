package newsapi

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

const (
	baseURL = "https://newsapi.org/v2"
)

type Client struct {
	apiKey string
	client *resty.Client
}

type NewsAPIResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt string    `json:"publishedAt"`
	Content     string    `json:"content"`
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewClient(apiKey string) *Client {
	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetHeader("X-Api-Key", apiKey)
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(5 * time.Second)

	return &Client{
		apiKey: apiKey,
		client: client,
	}
}

// FetchTopHeadlines fetches top headlines by category and country
func (c *Client) FetchTopHeadlines(category, country string, pageSize int) ([]types.ParsedArticle, error) {
	var response NewsAPIResponse

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"category": category,
			"country":  country,
			"pageSize": fmt.Sprintf("%d", pageSize),
		}).
		SetResult(&response).
		Get("/top-headlines")

	if err != nil {
		return nil, fmt.Errorf("failed to fetch top headlines: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("NewsAPI error %d: %s", resp.StatusCode(), resp.String())
	}

	if response.Status != "ok" {
		return nil, fmt.Errorf("NewsAPI returned status: %s", response.Status)
	}

	return c.convertArticles(response.Articles, category), nil
}

// FetchEverything fetches articles matching a query
func (c *Client) FetchEverything(query string, sortBy string, pageSize int) ([]types.ParsedArticle, error) {
	var response NewsAPIResponse

	// Calculate from date (last 7 days)
	from := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"q":        query,
			"sortBy":   sortBy,
			"pageSize": fmt.Sprintf("%d", pageSize),
			"from":     from,
			"language": "en",
		}).
		SetResult(&response).
		Get("/everything")

	if err != nil {
		return nil, fmt.Errorf("failed to fetch everything: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("NewsAPI error %d: %s", resp.StatusCode(), resp.String())
	}

	if response.Status != "ok" {
		return nil, fmt.Errorf("NewsAPI returned status: %s", response.Status)
	}

	return c.convertArticles(response.Articles, "general"), nil
}

// FetchBySource fetches articles from a specific source
func (c *Client) FetchBySource(sourceID string, pageSize int) ([]types.ParsedArticle, error) {
	var response NewsAPIResponse

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"sources":  sourceID,
			"pageSize": fmt.Sprintf("%d", pageSize),
		}).
		SetResult(&response).
		Get("/top-headlines")

	if err != nil {
		return nil, fmt.Errorf("failed to fetch by source: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("NewsAPI error %d: %s", resp.StatusCode(), resp.String())
	}

	if response.Status != "ok" {
		return nil, fmt.Errorf("NewsAPI returned status: %s", response.Status)
	}

	return c.convertArticles(response.Articles, "general"), nil
}

// FetchMultipleCategories fetches articles from multiple categories
func (c *Client) FetchMultipleCategories(country string, pageSize int) ([]types.ParsedArticle, error) {
	categories := []string{"technology", "business", "science", "entertainment", "sports"}
	allArticles := make([]types.ParsedArticle, 0)

	for _, category := range categories {
		articles, err := c.FetchTopHeadlines(category, country, pageSize/len(categories))
		if err != nil {
			log.Printf("Failed to fetch %s headlines: %v", category, err)
			continue
		}

		log.Printf("Fetched %d articles from NewsAPI category: %s", len(articles), category)
		allArticles = append(allArticles, articles...)

		// Rate limiting - NewsAPI has limits
		time.Sleep(500 * time.Millisecond)
	}

	return allArticles, nil
}

func (c *Client) convertArticles(articles []Article, defaultCategory string) []types.ParsedArticle {
	converted := make([]types.ParsedArticle, 0, len(articles))

	for _, article := range articles {
		// Skip articles with removed content
		if article.Title == "[Removed]" || article.URL == "" {
			continue
		}

		parsed := c.convertArticle(article, defaultCategory)
		converted = append(converted, parsed)
	}

	return converted
}

func (c *Client) convertArticle(article Article, defaultCategory string) types.ParsedArticle {
	// Parse published date
	publishedAt, err := time.Parse(time.RFC3339, article.PublishedAt)
	if err != nil {
		publishedAt = time.Now()
	}

	// Determine category
	category := defaultCategory
	if category == "technology" {
		category = "tech"
	}

	// Generate external ID from URL
	externalID := fmt.Sprintf("%x", md5Hash(article.URL))

	// Clean content
	content := article.Description
	if article.Content != "" && len(article.Content) > len(content) {
		content = article.Content
	}

	return types.ParsedArticle{
		Title:       article.Title,
		URL:         article.URL,
		Summary:     content,
		Content:     content,
		Author:      article.Author,
		Source:      article.Source.Name,
		Category:    category,
		Tags:        []string{}, // NewsAPI doesn't provide tags
		ImageURL:    article.URLToImage,
		PublishedAt: publishedAt,
		FetchedAt:   time.Now(),
		ExternalID:  externalID,
		FetchSource: "newsapi",
	}
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// GetSources fetches available news sources
func (c *Client) GetSources() ([]Source, error) {
	var response struct {
		Status  string   `json:"status"`
		Sources []Source `json:"sources"`
	}

	resp, err := c.client.R().
		SetResult(&response).
		Get("/sources")

	if err != nil {
		return nil, fmt.Errorf("failed to fetch sources: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("NewsAPI error %d: %s", resp.StatusCode(), resp.String())
	}

	return response.Sources, nil
}

// GetPopularSources returns a curated list of popular source IDs
func GetPopularSources() []string {
	return []string{
		"techcrunch",
		"the-verge",
		"ars-technica",
		"wired",
		"hacker-news",
		"reuters",
		"bbc-news",
		"cnn",
		"bloomberg",
		"the-wall-street-journal",
		"the-new-york-times",
		"the-washington-post",
		"associated-press",
		"google-news",
	}
}

// ValidateAPIKey checks if the API key is valid
func (c *Client) ValidateAPIKey() error {
	var response NewsAPIResponse

	resp, err := c.client.R().
		SetQueryParam("pageSize", "1").
		SetResult(&response).
		Get("/top-headlines")

	if err != nil {
		return fmt.Errorf("failed to validate API key: %w", err)
	}

	if resp.StatusCode() == 401 {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode() == 429 {
		return fmt.Errorf("rate limit exceeded")
	}

	if response.Status != "ok" {
		return fmt.Errorf("API returned status: %s", response.Status)
	}

	return nil
}
