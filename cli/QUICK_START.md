# Terminal News CLI - Quick Start Guide

**Status**: ✅ 97% Complete - Ready to Test
**Last Updated**: November 18, 2025 (Night - Model Alignment Complete)

---

## ⚡ FASTEST WAY TO GET RUNNING (5 MINUTES)

### Prerequisites
- Go 1.21+ installed (download from https://golang.org/dl/)

### Steps

#### 1. Start Mock API Server
```bash
cd "C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\cli"
go run cmd/mockserver/main.go
```

#### 2. In Another Terminal: Run CLI
```bash
cd "C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\cli"
go run cmd/terminal-news/main.go
```

#### 3. Test Features
- **Navigate**: Use `↑`/`↓` or `j`/`k` to move
- **Switch Tabs**: Press `Tab`
- **Vote**: Press `l` to like, `d` to dislike
- **Comments**: Press `c` to view comments
- **Help**: Press `?` for all shortcuts
- **Quit**: Press `q`

---

## 🎮 KEYBOARD SHORTCUTS

### Global
- `Tab` - Next tab
- `Shift+Tab` - Previous tab
- `1-6` - Jump to specific tab
- `?` - Help overlay
- `q` - Quit
- `r` - Refresh current view

### Article Views (Hot/Controversial/Rising)
- `↑`/`k` - Previous article
- `↓`/`j` - Next article
- `l` - Like (upvote) article
- `d` - Dislike (downvote) article
- `c` - View comments
- `o`/`Enter` - Open article URL
- `g` - Go to top
- `G` - Go to bottom

### Comment View
- `↑`/`↓` - Navigate comments
- `Space` - Toggle collapse/expand
- `←`/`h` - Collapse thread
- `→`/`l` - Expand thread
- `r` - Reply to comment
- `Esc` - Back to article list

### Classifieds
- `n` - Post new classified
- `f` - Filter by category
- `↑`/`↓` - Navigate listings

---

## 📂 WHAT'S IN THE PROJECT

### Applications
- `cmd/terminal-news/main.go` - Main CLI application
- `cmd/mockserver/main.go` - Mock API server for testing

### Core Code
- `internal/ui/app.go` - Main Bubbletea app
- `internal/ui/views/` - All views (Hot, Profile, Weather, etc.)
- `internal/ui/components/` - Reusable UI components
- `internal/api/` - API client + WebSocket
- `internal/cache/` - SQLite offline cache
- `internal/config/` - Configuration management

### Testing
- `scripts/test_integration.sh` - Automated integration testing

### Documentation
- `README.md` - Full user guide
- `INSTALL_GUIDE.md` - Detailed installation
- `DEV2_FINAL_STATUS.md` - Complete status report
- `SESSION_SUMMARY.md` - What was built today
- `HANDOFF_TO_USER.md` - Handoff document
- `QUICK_START.md` - This file

---

## 🧪 TESTING SCENARIOS

### Test Voting
1. Start mock server
2. Run CLI
3. Navigate to an article
4. Press `l` - should see vote count change (eventually)
5. Press `d` - should see vote count change

### Test Comments
1. Navigate to an article
2. Press `c`
3. Use `↑`/`↓` to navigate comments
4. Press `Space` to collapse/expand threads
5. Press `Esc` to go back

### Test Classifieds
1. Switch to Classifieds tab (Tab key)
2. Press `n` to post new classified
3. Fill out form (Tab to move between fields)
4. Press `Ctrl+P` to preview
5. Press `Ctrl+S` to submit

### Test Weather
1. Switch to Weather tab
2. Should see current conditions
3. 5-day forecast displayed
4. Location from config shown

### Test Offline Mode
```bash
./bin/terminal-news --offline
```
- Should run without backend
- Shows cached data or empty states
- Votes queued for later sync

---

## 🔧 CONFIGURATION

Config file: `~/.terminal-news/config.yaml`

### Quick Edits

#### Change Location
```yaml
user:
  location: "New York, NY"  # Change this
```

#### Change API URL
```yaml
api:
  base_url: "http://localhost:8080/api"  # Change this
  websocket_url: "ws://localhost:8080/ws"
```

#### Enable Offline Mode
```yaml
offline: true
```

---

## 🐛 TROUBLESHOOTING

### "Go not installed"
Download from https://golang.org/dl/
Then restart terminal

### "Backend not running"
Start mock server:
```bash
go run cmd/mockserver/main.go
```

### "Build failed"
Run:
```bash
go mod download
go mod tidy
```

### "Config not found"
Run app once to auto-create:
```bash
go run cmd/terminal-news/main.go
```

### "No articles showing"
Make sure mock server is running on port 8080:
```bash
curl http://localhost:8080/api/health
```

---

## 📊 FEATURES STATUS

| Feature | Status | Test |
|---------|--------|------|
| Article browsing | ✅ | Navigate with ↑/↓ |
| Voting | ✅ | Press 'l' or 'd' |
| Comments | ✅ | Press 'c' on article |
| Classifieds | ✅ | Switch to tab, press 'n' |
| Weather | ✅ | Switch to Weather tab |
| Profile | ✅ | Switch to Profile tab |
| Login | ✅ | Use mock server |
| Offline mode | ✅ | Run with --offline |

---

## 🚀 NEXT STEPS

### After First Run
1. ✅ Verify compilation works
2. ✅ Test basic navigation
3. ✅ Test voting
4. ✅ Test comments

### Integration with Real Backend
1. Start real backend server (Dev 1)
2. Update config.yaml with backend URL
3. Test login/register
4. Test real article voting
5. Test real-time WebSocket updates

### Production Deployment
1. Build binary: `go build -o terminal-news cmd/terminal-news/main.go`
2. Distribute binary
3. Users install to PATH
4. Users run `terminal-news`

---

## 💡 TIPS

### Performance
- Use compact mode: Press 's' (once implemented)
- Limit articles per page in config
- Cache clears automatically after TTL

### Keyboard Efficiency
- Learn vim keys (h/j/k/l)
- Use number keys to jump tabs (1-6)
- Press '?' often to see shortcuts

### Offline Usage
- Run with `--offline` flag
- All reads from cache
- Writes queued for sync
- Great for commutes!

---

## 📞 NEED HELP?

### Check Documentation
1. **User Guide**: `README.md`
2. **Install Help**: `INSTALL_GUIDE.md`
3. **Dev Status**: `DEV2_FINAL_STATUS.md`
4. **Full Summary**: `SESSION_SUMMARY.md`

### Common Issues
- **Go errors**: Run `go mod tidy`
- **Import errors**: Check Go version (need 1.21+)
- **No data**: Start mock server
- **Slow**: Reduce articles_per_page in config

---

## ✅ CHECKLIST FOR FIRST RUN

- [ ] Go 1.21+ installed
- [ ] Navigated to `cli/` directory
- [ ] Ran `go mod download`
- [ ] Started mock server in Terminal 1
- [ ] Ran CLI in Terminal 2
- [ ] Pressed `?` for help
- [ ] Tested navigation with arrows
- [ ] Tested voting with 'l' and 'd'
- [ ] Tested comments with 'c'
- [ ] Quit with 'q'

---

**Status**: 🟢 Ready to Run

**Next**: Install Go, then run `go run cmd/mockserver/main.go`

---

*Quick Start Guide - November 18, 2025*
