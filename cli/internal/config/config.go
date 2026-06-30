package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	API      APIConfig      `mapstructure:"api"`
	Cache    CacheConfig    `mapstructure:"cache"`
	UI       UIConfig       `mapstructure:"ui"`
	Keys     KeyBindings    `mapstructure:"keybindings"`
	User     UserConfig     `mapstructure:"user"`
	Offline  bool           `mapstructure:"offline"`
}

type APIConfig struct {
	BaseURL      string        `mapstructure:"base_url"`
	WebSocketURL string        `mapstructure:"websocket_url"`
	Timeout      time.Duration `mapstructure:"timeout"`
}

type CacheConfig struct {
	DatabasePath  string `mapstructure:"database_path"`
	TTL           int    `mapstructure:"ttl"`
	MaxArticles   int    `mapstructure:"max_articles"`
	MaxComments   int    `mapstructure:"max_comments"`
}

type UIConfig struct {
	Theme            string `mapstructure:"theme"`
	CompactMode      bool   `mapstructure:"compact_mode"`
	ShowEmojis       bool   `mapstructure:"show_emojis"`
	ArticlesPerPage  int    `mapstructure:"articles_per_page"`
	RefreshInterval  int    `mapstructure:"refresh_interval"` // seconds
}

type KeyBindings struct {
	Quit        string `mapstructure:"quit"`
	Refresh     string `mapstructure:"refresh"`
	Help        string `mapstructure:"help"`
	Search      string `mapstructure:"search"`
	NextTab     string `mapstructure:"next_tab"`
	PrevTab     string `mapstructure:"prev_tab"`
	Up          string `mapstructure:"up"`
	Down        string `mapstructure:"down"`
	Select      string `mapstructure:"select"`
	Back        string `mapstructure:"back"`
	Like        string `mapstructure:"like"`
	Dislike     string `mapstructure:"dislike"`
	Comment     string `mapstructure:"comment"`
	OpenBrowser string `mapstructure:"open_browser"`
}

type UserConfig struct {
	Location   string `mapstructure:"location"`
	DefaultTab string `mapstructure:"default_tab"`
	Token      string `mapstructure:"token"`
	Username   string `mapstructure:"username"`
}

func Load() (*Config, error) {
	// Set defaults
	setDefaults()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand home directory in cache path
	if cfg.Cache.DatabasePath != "" {
		home, _ := os.UserHomeDir()
		cfg.Cache.DatabasePath = filepath.Join(home, ".terminal-news", "cache.db")
	}

	return &cfg, nil
}

func setDefaults() {
	// API defaults
	viper.SetDefault("api.base_url", "http://localhost:8080/api")
	viper.SetDefault("api.websocket_url", "ws://localhost:8080/ws")
	viper.SetDefault("api.timeout", "30s")

	// Cache defaults
	viper.SetDefault("cache.database_path", "~/.terminal-news/cache.db")
	viper.SetDefault("cache.ttl", 3600) // 1 hour
	viper.SetDefault("cache.max_articles", 1000)
	viper.SetDefault("cache.max_comments", 5000)

	// UI defaults
	viper.SetDefault("ui.theme", "default")
	viper.SetDefault("ui.compact_mode", false)
	viper.SetDefault("ui.show_emojis", true)
	viper.SetDefault("ui.articles_per_page", 50)
	viper.SetDefault("ui.refresh_interval", 300) // 5 minutes

	// Keybindings defaults
	viper.SetDefault("keybindings.quit", "q")
	viper.SetDefault("keybindings.refresh", "r")
	viper.SetDefault("keybindings.help", "?")
	viper.SetDefault("keybindings.search", "/")
	viper.SetDefault("keybindings.next_tab", "tab")
	viper.SetDefault("keybindings.prev_tab", "shift+tab")
	viper.SetDefault("keybindings.up", "up,k")
	viper.SetDefault("keybindings.down", "down,j")
	viper.SetDefault("keybindings.select", "enter")
	viper.SetDefault("keybindings.back", "esc")
	viper.SetDefault("keybindings.like", "l")
	viper.SetDefault("keybindings.dislike", "d")
	viper.SetDefault("keybindings.comment", "c")
	viper.SetDefault("keybindings.open_browser", "o")

	// User defaults
	viper.SetDefault("user.location", "")
	viper.SetDefault("user.default_tab", "hot")

	// Offline mode
	viper.SetDefault("offline", false)
}

func CreateDefault() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".terminal-news")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")

	// Create default config file
	defaultConfig := `# Terminal News Configuration

api:
  base_url: "http://localhost:8080/api"
  websocket_url: "ws://localhost:8080/ws"
  timeout: 30s

cache:
  database_path: "~/.terminal-news/cache.db"
  ttl: 3600  # 1 hour
  max_articles: 1000
  max_comments: 5000

ui:
  theme: "default"
  compact_mode: false
  show_emojis: true
  articles_per_page: 50
  refresh_interval: 300  # 5 minutes

keybindings:
  quit: "q"
  refresh: "r"
  help: "?"
  search: "/"
  next_tab: "tab"
  prev_tab: "shift+tab"
  up: "up,k"
  down: "down,j"
  select: "enter"
  back: "esc"
  like: "l"
  dislike: "d"
  comment: "c"
  open_browser: "o"

user:
  location: ""
  default_tab: "hot"
  token: ""
  username: ""

offline: false
`

	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}
