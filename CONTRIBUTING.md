# Contributing to Terminal News

First off, thank you for considering contributing to Terminal News! It's people like you that make Terminal News a great tool for the community.

## Code of Conduct

This project and everyone participating in it is governed by our commitment to maintaining a welcoming, harassment-free environment. By participating, you are expected to uphold this code.

**Our Standards:**
- Be respectful and inclusive
- Welcome newcomers warmly
- Accept constructive criticism gracefully
- Focus on what's best for the community
- Show empathy towards others

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

**Bug Report Template:**
```
**Description:**
A clear description of the bug

**Steps to Reproduce:**
1. Go to '...'
2. Click on '...'
3. Scroll down to '...'
4. See error

**Expected Behavior:**
What you expected to happen

**Actual Behavior:**
What actually happened

**Environment:**
- OS: [e.g. macOS 14.1, Ubuntu 22.04, Windows 11]
- Terminal: [e.g. iTerm2, Alacritty, Windows Terminal]
- Terminal News Version: [e.g. 1.0.0]

**Screenshots/Logs:**
If applicable, add screenshots or error logs
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

**Enhancement Template:**
```
**Problem:**
What problem does this solve?

**Proposed Solution:**
How should this work?

**Alternatives Considered:**
What other approaches did you think about?

**Additional Context:**
Any mockups, examples, or additional details
```

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:
- `good first issue` - Good for newcomers
- `help wanted` - Need community assistance
- `documentation` - Improve docs

### Pull Requests

**Process:**

1. **Fork the repo** and create your branch from `main`
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes**
   - Write clear, concise commit messages
   - Follow the code style (see below)
   - Add tests if applicable
   - Update documentation

3. **Test thoroughly**
   - Run existing tests: `go test ./...`
   - Test manually in terminal
   - Test on your target OS

4. **Commit with good messages**
   ```
   Add feature: Real-time notification system

   - Implements WebSocket connection for live updates
   - Adds notification badge to UI
   - Includes user preference settings

   Closes #123
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/my-new-feature
   ```

6. **Open a Pull Request**
   - Use the PR template
   - Link related issues
   - Describe what changed and why
   - Add screenshots for UI changes

**PR Template:**
```
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## How Has This Been Tested?
Describe your testing process

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Commented complex code
- [ ] Updated documentation
- [ ] No new warnings
- [ ] Added tests
- [ ] All tests pass
- [ ] Works on Mac/Linux/Windows
```

## Development Setup

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+
- Git

### Local Development

1. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/terminal-news.git
   cd terminal-news
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up local database**
   ```bash
   # Install PostgreSQL and Redis
   # Then run:
   docker-compose up -d db redis

   # Run migrations
   go run cmd/migrate/main.go up
   ```

4. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your local settings
   ```

5. **Run the development server**
   ```bash
   # API server
   go run cmd/api/main.go

   # In another terminal, run the client
   go run cmd/terminal-news/main.go
   ```

## Code Style Guidelines

### Go Style

**Follow standard Go conventions:**
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Run `golint` and address warnings
- Use meaningful variable names

**Examples:**

```go
// Good
func FetchArticles(limit int, offset int) ([]Article, error) {
    if limit <= 0 {
        return nil, ErrInvalidLimit
    }

    articles, err := db.Query(
        "SELECT * FROM articles LIMIT $1 OFFSET $2",
        limit, offset,
    )
    if err != nil {
        return nil, fmt.Errorf("fetch articles: %w", err)
    }

    return articles, nil
}

// Bad
func get(l int, o int) ([]Article, error) {
    a, e := db.Query("SELECT * FROM articles LIMIT $1 OFFSET $2", l, o)
    if e != nil {
        return nil, e
    }
    return a, nil
}
```

**Error Handling:**
```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to connect to database: %w", err)
}
```

**Comments:**
```go
// Public functions must have doc comments
// FetchArticles retrieves paginated articles from the database.
// Returns an error if limit is invalid or database query fails.
func FetchArticles(limit int, offset int) ([]Article, error) {
    // Implementation
}
```

### UI Code (Bubbletea)

**Follow Elm Architecture:**
```go
// Model
type Model struct {
    articles []Article
    cursor   int
    loading  bool
}

// Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case ArticlesLoadedMsg:
        m.articles = msg.Articles
        m.loading = false
        return m, nil
    }
    return m, nil
}

// View
func (m Model) View() string {
    if m.loading {
        return "Loading..."
    }
    return m.renderArticles()
}
```

### Database Queries

**Use prepared statements:**
```go
// Good
stmt := `SELECT id, title, url FROM articles WHERE id = $1`
row := db.QueryRow(stmt, articleID)

// Bad
query := fmt.Sprintf("SELECT * FROM articles WHERE id = %d", articleID)
```

**Use transactions for multi-step operations:**
```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// Multiple operations
if err := insertArticle(tx, article); err != nil {
    return err
}
if err := updateVotes(tx, articleID); err != nil {
    return err
}

return tx.Commit()
```

## Testing

### Writing Tests

**Unit tests:**
```go
func TestFetchArticles(t *testing.T) {
    t.Run("returns articles with valid parameters", func(t *testing.T) {
        articles, err := FetchArticles(10, 0)
        assert.NoError(t, err)
        assert.Len(t, articles, 10)
    })

    t.Run("returns error with invalid limit", func(t *testing.T) {
        _, err := FetchArticles(-1, 0)
        assert.Error(t, err)
    })
}
```

**Integration tests:**
```go
func TestArticleVoting(t *testing.T) {
    // Set up test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Create test article
    article := createTestArticle(t, db)

    // Vote
    err := VoteOnArticle(db, article.ID, user.ID, "like")
    assert.NoError(t, err)

    // Verify
    votes := getVoteCount(t, db, article.ID)
    assert.Equal(t, 1, votes)
}
```

**Run tests:**
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/api/...

# Verbose
go test -v ./...
```

## Documentation

### Code Documentation

- All public functions must have doc comments
- Complex logic should have inline comments
- Use examples in doc comments when helpful

```go
// CalculateControversyScore computes how controversial an article is
// based on the ratio of likes to dislikes.
//
// The score ranges from 0-100, where:
//   - 0 = unanimous agreement (all likes or all dislikes)
//   - 100 = perfectly split opinion
//
// Example:
//   score := CalculateControversyScore(100, 95) // Returns ~95
func CalculateControversyScore(likes, dislikes int) float64 {
    // Implementation
}
```

### User Documentation

When adding features, update:
- `README.md` - Overview and quick start
- `docs/` - Detailed feature documentation
- `UI_MOCKUPS.md` - UI changes or new views
- Help screen in app - Keyboard shortcuts

## Git Workflow

### Branches

- `main` - Stable, production-ready code
- `develop` - Integration branch for features
- `feature/feature-name` - New features
- `fix/bug-description` - Bug fixes
- `docs/what-changed` - Documentation only

### Commit Messages

**Format:**
```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style (formatting, no logic change)
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance (dependencies, build, etc.)

**Examples:**
```
feat: Add real-time vote updates via WebSocket

Implements WebSocket connection for live vote count updates.
Users now see vote changes without manual refresh.

- Add WebSocket server endpoint
- Update UI to subscribe to vote events
- Add reconnection logic

Closes #45

---

fix: Prevent duplicate article submissions

Articles with identical URLs within 24 hours are now rejected.

Fixes #67

---

docs: Update installation instructions for Windows

Added Windows-specific setup steps and troubleshooting.
```

## Release Process

**Versioning:**
- Follow [Semantic Versioning](https://semver.org/)
- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes

**Release Checklist:**
1. Update version in code
2. Update CHANGELOG.md
3. Create release branch
4. Run full test suite
5. Build binaries for all platforms
6. Create GitHub release with notes
7. Publish binaries
8. Announce on social media

## Community

### Getting Help

- **Discord/Slack** (when available): Real-time chat
- **GitHub Discussions**: Questions, ideas, show-and-tell
- **GitHub Issues**: Bug reports, feature requests
- **Email**: For security issues only

### Recognition

Contributors are recognized in:
- `CONTRIBUTORS.md` file
- Release notes for their contributions
- Special thanks in major releases

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (TBD).

---

## Questions?

Don't hesitate to ask! We're here to help:
- Open a discussion on GitHub
- Tag issues with `question`
- Reach out to maintainers

**Thank you for contributing to Terminal News!**
