# Terminal News CLI - Installation & Quick Start Guide

## 🎯 Prerequisites

### 1. Install Go

**You need Go 1.21 or higher installed on your system.**

#### Windows:
```powershell
# Using Chocolatey
choco install golang

# Or download installer from:
# https://golang.org/dl/
```

#### macOS:
```bash
# Using Homebrew
brew install go

# Or download from:
# https://golang.org/dl/
```

#### Linux:
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Or download from:
# https://golang.org/dl/
```

### 2. Verify Installation

```bash
go version
# Should output: go version go1.21.x or higher
```

---

## 🚀 Quick Start

### Step 1: Navigate to Project

```bash
cd "C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\cli"
```

### Step 2: Download Dependencies

```bash
go mod download
```

This will install all required packages:
- Bubbletea (TUI framework)
- Bubbles (TUI components)
- Lipgloss (Styling)
- Resty (HTTP client)
- Gorilla WebSocket
- SQLite driver
- Cobra (CLI framework)
- Viper (Configuration)

### Step 3: Run the Application

```bash
# Run directly (recommended for development)
go run cmd/terminal-news/main.go

# Or build first, then run
go build -o bin/terminal-news cmd/terminal-news/main.go
./bin/terminal-news
```

### Step 4: First Launch

On first run, Terminal News will:
1. Create `~/.terminal-news/` directory
2. Generate `config.yaml` with default settings
3. Create `cache.db` SQLite database
4. Start the TUI application

---

## ⚙️ Configuration

### Config File Location

- **Windows**: `C:\Users\YourName\.terminal-news\config.yaml`
- **macOS/Linux**: `~/.terminal-news/config.yaml`

### Edit Configuration

```yaml
api:
  base_url: "http://localhost:8080/api"  # Change if backend runs elsewhere
  websocket_url: "ws://localhost:8080/ws"
  timeout: 30s

user:
  location: "Your City, State"  # For weather and local classifieds
  default_tab: "hot"           # Tab to show on startup
```

### Custom Config File

```bash
# Use a different config file
go run cmd/terminal-news/main.go --config /path/to/config.yaml
```

---

## 🎮 Keyboard Shortcuts

### Navigation
- `Tab` - Next tab
- `Shift+Tab` - Previous tab
- `1-6` - Jump to specific tab
- `↑/↓` or `k/j` - Navigate lists
- `q` or `Ctrl+C` - Quit

### Actions
- `r` - Refresh current view
- `?` - Show help
- `Enter` - Select/Open
- `Esc` - Go back
- `l` - Like article
- `d` - Dislike article
- `c` - View comments
- `o` - Open in browser

---

## 🛠️ Development Commands

### Using Make

```bash
# Show all available commands
make help

# Build the binary
make build

# Run the application
make run

# Run in offline mode
make offline

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make format

# Lint code
make lint

# Clean build artifacts
make clean
```

### Manual Commands

```bash
# Download dependencies
go mod download

# Tidy modules
go mod tidy

# Build for current platform
go build -o bin/terminal-news cmd/terminal-news/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o terminal-news-linux cmd/terminal-news/main.go

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o terminal-news-macos cmd/terminal-news/main.go

# Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o terminal-news-macos-arm cmd/terminal-news/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o terminal-news.exe cmd/terminal-news/main.go

# Run tests
go test ./...

# Run specific package tests
go test -v ./internal/cache

# Generate test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🐛 Troubleshooting

### Issue: "go: command not found"

**Solution**: Go is not installed or not in PATH
```bash
# Verify Go installation
which go  # macOS/Linux
where go  # Windows

# Add Go to PATH (if needed)
export PATH=$PATH:/usr/local/go/bin  # macOS/Linux
```

### Issue: "Cannot connect to backend"

**Solution**: Backend API is not running
```bash
# The CLI needs the backend API server running
# Check backend status at: http://localhost:8080/api

# Or run in offline mode to test with cached data
go run cmd/terminal-news/main.go --offline
```

### Issue: "Database locked"

**Solution**: Another instance is running
```bash
# Close all Terminal News instances
# Or delete the database
rm ~/.terminal-news/cache.db
```

### Issue: "Module not found" or "Import errors"

**Solution**: Dependencies not downloaded
```bash
# Download all dependencies
go mod download

# Verify modules
go mod verify

# Clean and re-download
go clean -modcache
go mod download
```

### Issue: "Permission denied" on Linux/macOS

**Solution**: Make binary executable
```bash
chmod +x bin/terminal-news
./bin/terminal-news
```

### Issue: Weird rendering or colors

**Solution**: Terminal color support
```bash
# Set terminal type
export TERM=xterm-256color

# Or try a modern terminal emulator:
# - Windows: Windows Terminal
# - macOS: iTerm2, Alacritty
# - Linux: Alacritty, Kitty, GNOME Terminal
```

---

## 📦 Installation Methods

### Method 1: Run from Source (Recommended for Development)

```bash
cd C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\cli
go run cmd/terminal-news/main.go
```

**Pros**: Easy to modify code, see changes immediately
**Cons**: Slower startup, requires source code

### Method 2: Build and Install

```bash
# Build
make build

# Install to $GOPATH/bin
make install

# Now run from anywhere
terminal-news
```

**Pros**: Fast startup, can run from anywhere
**Cons**: Need to rebuild after changes

### Method 3: Go Install

```bash
# Install directly from source
go install github.com/wesellis/terminal-news/cli/cmd/terminal-news@latest

# Run
terminal-news
```

**Pros**: Simple, automatic
**Cons**: Requires code to be on GitHub

---

## 🧪 Testing the Installation

### 1. Check Version
```bash
go run cmd/terminal-news/main.go version
# Should output: Terminal News v0.1.0
```

### 2. Test Offline Mode
```bash
go run cmd/terminal-news/main.go --offline
# Should start with cached data (or empty if first run)
```

### 3. Test Configuration
```bash
# Check if config was created
cat ~/.terminal-news/config.yaml  # macOS/Linux
type %USERPROFILE%\.terminal-news\config.yaml  # Windows
```

### 4. Run Tests
```bash
make test
# All tests should pass
```

---

## 📝 Next Steps

### After Installation:

1. **Configure Your Location**
   - Edit `~/.terminal-news/config.yaml`
   - Set `user.location` for weather and local classifieds

2. **Start Backend API**
   - The CLI needs the backend server running
   - See `/backend` directory for setup instructions

3. **Create an Account**
   - Run Terminal News
   - Navigate to Profile tab
   - Register a new account

4. **Explore Features**
   - Browse Hot articles
   - Vote on content
   - Post a classified
   - Check the weather

---

## 🔧 Advanced Configuration

### Environment Variables

```bash
# Override config file location
export TERMINAL_NEWS_CONFIG=/path/to/config.yaml

# Override API URL
export TERMINAL_NEWS_API_URL=http://api.terminal-news.com

# Enable debug logging
export TERMINAL_NEWS_DEBUG=true
```

### Cache Management

```bash
# Clear cache
rm ~/.terminal-news/cache.db

# View cache size
du -h ~/.terminal-news/cache.db

# Backup cache
cp ~/.terminal-news/cache.db ~/cache-backup.db
```

### Performance Tuning

Edit `config.yaml`:
```yaml
ui:
  articles_per_page: 25  # Reduce for better performance
  compact_mode: true     # Less visual elements
  show_emojis: false     # Disable emojis

cache:
  ttl: 7200             # Cache longer (2 hours)
  max_articles: 500      # Reduce cache size
```

---

## 🆘 Getting Help

### Documentation
- **README.md** - Full documentation
- **DEV_STATUS.md** - Development progress
- **Makefile** - Available commands
- **UI_MOCKUPS.md** - Design specifications

### Support
- **GitHub Issues**: Report bugs
- **GitHub Discussions**: Ask questions
- **Email**: support@terminal-news.com (coming soon)

### Debugging

Enable verbose output:
```bash
# Run with debug flag
go run cmd/terminal-news/main.go --debug

# View logs
tail -f ~/.terminal-news/debug.log
```

---

## ✅ Verification Checklist

- [ ] Go 1.21+ installed (`go version`)
- [ ] Dependencies downloaded (`go mod download`)
- [ ] Config file created (`~/.terminal-news/config.yaml` exists)
- [ ] Application runs (`go run cmd/terminal-news/main.go`)
- [ ] Tests pass (`make test`)
- [ ] Can navigate tabs (Tab key works)
- [ ] Can quit (q key works)

---

## 🎉 You're Ready!

Terminal News CLI is now installed and ready to use!

```bash
# Start reading news
go run cmd/terminal-news/main.go

# Or if you installed it
terminal-news
```

**Welcome to Terminal News - AM Radio for the Information Age** 📰⚡

---

*For more information, see README.md or visit the full documentation.*
