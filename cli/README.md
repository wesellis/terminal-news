# Terminal News - CLI Client

> 🗞 **AM Radio for the Information Age**

A beautiful, keyboard-driven terminal user interface for Terminal News - the terminal-native news aggregator with community curation, local classifieds, and real-time weather.

## Features

- **📰 News Aggregation**: Hot, Controversial, and Rising feeds with community voting
- **💬 Comments**: Threaded discussions on articles
- **📍 Local Classifieds**: Browse and post classified ads for your area
- **🌤 Weather Widget**: Real-time NOAA weather integrated into the interface
- **⚡ Real-time Updates**: WebSocket connections for live content updates
- **💾 Offline Mode**: Local SQLite cache for offline reading
- **⌨️ Keyboard-First**: Fully navigable with vim-style keybindings
- **🎨 Beautiful TUI**: Built with Bubbletea and Lipgloss for a stunning terminal experience

## Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Terminal with 256 colors** - Most modern terminals work great
- **Backend API** - Running instance of the Terminal News backend

## Installation

### From Source

```bash
# Clone the repository
cd C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\cli

# Install dependencies
go mod download

# Build the binary
go build -o bin/terminal-news cmd/terminal-news/main.go

# Run it!
./bin/terminal-news
```

### Quick Start

```bash
# Run directly without building
go run cmd/terminal-news/main.go

# Run in offline mode (uses cached data)
go run cmd/terminal-news/main.go --offline

# Specify custom config file
go run cmd/terminal-news/main.go --config /path/to/config.yaml
```

## Configuration

On first run, Terminal News creates a default configuration file at:
- **Windows**: `C:\Users\YourName\.terminal-news\config.yaml`
- **macOS/Linux**: `~/.terminal-news/config.yaml`

### Configuration Options

```yaml
api:
  base_url: "http://localhost:8080/api"
  websocket_url: "ws://localhost:8080/ws"
  timeout: 30s

cache:
  database_path: "~/.terminal-news/cache.db"
  ttl: 3600  # seconds
  max_articles: 1000
  max_comments: 5000

ui:
  theme: "default"
  compact_mode: false
  show_emojis: true
  articles_per_page: 50
  refresh_interval: 300  # seconds

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
  location: "San Francisco, CA"  # For weather and local classifieds
  default_tab: "hot"
  token: ""  # Filled after login
  username: ""  # Filled after login
```

## Keyboard Shortcuts

### Global Navigation
- `q` or `Ctrl+C` - Quit
- `Tab` - Next tab
- `Shift+Tab` - Previous tab
- `1-6` - Jump to specific tab (Hot, Controversial, Rising, Profile, Weather, Classifieds)
- `r` - Refresh current view
- `?` - Show help

### Article Navigation
- `↑` or `k` - Move up
- `↓` or `j` - Move down
- `Enter` - Open selected article
- `l` - Like article
- `d` - Dislike article
- `c` - View/add comments
- `o` - Open article in browser
- `Esc` - Go back

### Classifieds
- `n` - Post new classified
- `e` - Edit your classified
- `Del` - Delete your classified
- `/` - Search classifieds
- `f` - Filter by category

## Development

### Project Structure

```
cli/
├── cmd/
│   └── terminal-news/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── client.go            # REST API client
│   │   └── websocket.go         # WebSocket client
│   ├── cache/
│   │   └── cache.go             # SQLite local cache
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   └── models.go            # Data models
│   └── ui/
│       ├── app.go               # Main app model
│       ├── components/          # Reusable UI components
│       ├── views/               # Tab views
│       └── styles/
│           └── styles.go        # Lipgloss styles
├── go.mod
└── go.sum
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/cache
```

### Development Mode

```bash
# Watch for changes and auto-reload (requires nodemon)
make dev

# Run linter
make lint

# Format code
make format
```

## Building for Multiple Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o terminal-news-linux cmd/terminal-news/main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o terminal-news-macos-intel cmd/terminal-news/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o terminal-news-macos-arm cmd/terminal-news/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o terminal-news.exe cmd/terminal-news/main.go
```

## Troubleshooting

### "Cannot connect to backend"
- Ensure the backend API is running on the configured URL
- Check `api.base_url` in your config file
- Try running with `--offline` flag to test with cached data

### "Database locked" error
- Another instance might be running
- Delete `~/.terminal-news/cache.db` and restart

### Rendering issues
- Ensure your terminal supports 256 colors
- Try a different terminal emulator (iTerm2, Alacritty, Windows Terminal)
- Set `TERM=xterm-256color` environment variable

### Slow performance
- Reduce `articles_per_page` in config
- Enable `compact_mode` in config
- Clear old cache: `rm ~/.terminal-news/cache.db`

## Contributing

This is part of the larger Terminal News project. See the main project's CONTRIBUTING.md for guidelines.

### Local Development Setup

1. Install Go 1.21+
2. Clone the repository
3. Run `go mod download`
4. Make your changes
5. Run tests: `go test ./...`
6. Format code: `go fmt ./...`
7. Submit a pull request

## Architecture

Terminal News CLI uses the [Elm Architecture](https://guide.elm-lang.org/architecture/) via [Bubbletea](https://github.com/charmbracelet/bubbletea):

1. **Model** - Application state
2. **Update** - Handle messages and update state
3. **View** - Render the current state

### Data Flow

```
User Input → Messages → Update → Model → View → Terminal
                ↑                           │
                └──── WebSocket Events ─────┘
                └──── API Responses ────────┘
```

## Tech Stack

- **[Bubbletea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - TUI components
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling
- **[Glamour](https://github.com/charmbracelet/glamour)** - Markdown rendering
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[Viper](https://github.com/spf13/viper)** - Configuration
- **[Resty](https://github.com/go-resty/resty)** - HTTP client
- **[Gorilla WebSocket](https://github.com/gorilla/websocket)** - WebSocket client
- **[SQLite](https://www.sqlite.org/)** - Local cache database

## License

MIT License - See LICENSE file for details

## Links

- **Documentation**: [Full Docs](../docs/)
- **Backend API**: [Backend Guide](../backend/)
- **UI Mockups**: [Design Mockups](../design/UI_MOCKUPS.md)
- **GitHub**: [Terminal News](https://github.com/wesellis/Terminal-AM)

## Support

- **Issues**: Report bugs on GitHub Issues
- **Discussions**: Join GitHub Discussions
- **Discord**: [Join our Discord](#) (coming soon)

---

**Built with ❤️ by the Terminal News team**

*News at the speed of terminal* ⚡
