package api

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/wesellis/terminal-news/cli/internal/models"
)

// Client handles all API communications
type Client struct {
	http  *resty.Client
	token string
}

// NewClient creates a new API client
func NewClient(baseURL string, timeout time.Duration) *Client {
	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetTimeout(timeout)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("User-Agent", "Terminal-News/0.1.0")

	return &Client{
		http: client,
	}
}

// SetAuthToken sets the authentication token
func (c *Client) SetAuthToken(token string) {
	c.token = token
	c.http.SetAuthToken(token)
}

// Auth endpoints

func (c *Client) Login(username, password string) (*models.LoginResponse, error) {
	var result models.LoginResponse

	resp, err := c.http.R().
		SetBody(map[string]string{
			"username": username,
			"password": password,
		}).
		SetResult(&result).
		Post("/auth/login")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("login failed: %s", resp.Status())
	}

	c.SetAuthToken(result.Token)
	return &result, nil
}

func (c *Client) Register(username, email, password string) error {
	resp, err := c.http.R().
		SetBody(map[string]string{
			"username": username,
			"email":    email,
			"password": password,
		}).
		Post("/auth/register")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("registration failed: %s", resp.Status())
	}

	return nil
}

func (c *Client) Logout() error {
	resp, err := c.http.R().
		Post("/auth/logout")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("logout failed: %s", resp.Status())
	}

	c.token = ""
	c.http.SetAuthToken("")
	return nil
}

// Article endpoints

func (c *Client) GetArticles(feed string, offset, limit int) (*models.ArticlesResponse, error) {
	var result models.ArticlesResponse

	resp, err := c.http.R().
		SetQueryParams(map[string]string{
			"feed":   feed,
			"offset": fmt.Sprintf("%d", offset),
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&result).
		Get("/articles")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get articles: %s", resp.Status())
	}

	return &result, nil
}

func (c *Client) GetArticle(id int64) (*models.ArticleWithRanking, error) {
	var result models.ArticleWithRanking

	resp, err := c.http.R().
		SetResult(&result).
		Get(fmt.Sprintf("/articles/%d", id))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get article: %s", resp.Status())
	}

	return &result, nil
}

func (c *Client) VoteArticle(articleID int64, voteType string) error {
	resp, err := c.http.R().
		SetBody(map[string]string{
			"type": voteType,
		}).
		Post(fmt.Sprintf("/articles/%d/vote", articleID))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to vote: %s", resp.Status())
	}

	return nil
}

// Comment endpoints

func (c *Client) GetComments(articleID int64) (*models.CommentsResponse, error) {
	var result models.CommentsResponse

	resp, err := c.http.R().
		SetResult(&result).
		Get(fmt.Sprintf("/articles/%d/comments", articleID))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get comments: %s", resp.Status())
	}

	return &result, nil
}

func (c *Client) PostComment(articleID int64, content string, parentID *int64) error {
	body := map[string]interface{}{
		"content": content,
	}

	if parentID != nil {
		body["parent_id"] = *parentID
	}

	resp, err := c.http.R().
		SetBody(body).
		Post(fmt.Sprintf("/articles/%d/comments", articleID))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to post comment: %s", resp.Status())
	}

	return nil
}

// Classifieds endpoints

func (c *Client) GetClassifieds(category, location string, offset, limit int) (*models.ClassifiedsResponse, error) {
	var result models.ClassifiedsResponse

	params := map[string]string{
		"offset": fmt.Sprintf("%d", offset),
		"limit":  fmt.Sprintf("%d", limit),
	}

	if category != "" {
		params["category"] = category
	}
	if location != "" {
		params["location"] = location
	}

	resp, err := c.http.R().
		SetQueryParams(params).
		SetResult(&result).
		Get("/classifieds")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get classifieds: %s", resp.Status())
	}

	return &result, nil
}

func (c *Client) GetClassified(id int64) (*models.Classified, error) {
	var result models.Classified

	resp, err := c.http.R().
		SetResult(&result).
		Get(fmt.Sprintf("/classifieds/%d", id))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get classified: %s", resp.Status())
	}

	return &result, nil
}

func (c *Client) PostClassified(classified *models.Classified) error {
	resp, err := c.http.R().
		SetBody(classified).
		Post("/classifieds")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to post classified: %s", resp.Status())
	}

	return nil
}

func (c *Client) UpdateClassified(id int64, classified *models.Classified) error {
	resp, err := c.http.R().
		SetBody(classified).
		Put(fmt.Sprintf("/classifieds/%d", id))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to update classified: %s", resp.Status())
	}

	return nil
}

func (c *Client) DeleteClassified(id int64) error {
	resp, err := c.http.R().
		Delete(fmt.Sprintf("/classifieds/%d", id))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to delete classified: %s", resp.Status())
	}

	return nil
}

// User endpoints

func (c *Client) GetProfile() (*models.User, error) {
	var result models.User

	resp, err := c.http.R().
		SetResult(&result).
		Get("/user/profile")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get profile: %s", resp.Status())
	}

	return &result, nil
}

func (c *Client) GetActivity() ([]models.Activity, error) {
	var result []models.Activity

	resp, err := c.http.R().
		SetResult(&result).
		Get("/user/activity")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get activity: %s", resp.Status())
	}

	return result, nil
}

func (c *Client) GetUserClassifieds() ([]models.Classified, error) {
	var result []models.Classified

	resp, err := c.http.R().
		SetResult(&result).
		Get("/user/classifieds")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get user classifieds: %s", resp.Status())
	}

	return result, nil
}

// Weather endpoint

func (c *Client) GetWeather(location string) (*models.Weather, error) {
	var result models.Weather

	resp, err := c.http.R().
		SetQueryParam("location", location).
		SetResult(&result).
		Get("/weather")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get weather: %s", resp.Status())
	}

	return &result, nil
}
