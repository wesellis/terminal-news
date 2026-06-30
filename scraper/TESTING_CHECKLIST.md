# Scraper Testing & Verification Checklist

## ✅ COMPLETED (Code Review)
- [x] All import paths fixed to use `github.com/wesellis/terminal-news`
- [x] All 11 Go files created with proper package structure
- [x] go.mod configured with correct dependencies
- [x] Dockerfile created
- [x] Environment template (.env.example) created
- [x] Documentation written

## 🔨 WEEK 1: BUILD & BASIC TESTING

### Day 1: Environment Setup
- [ ] Install Go 1.21+ on development machine
  ```bash
  # Download from https://go.dev/dl/
  go version  # Should show 1.21 or higher
  ```

- [ ] Clone/navigate to project
  ```bash
  cd C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\scraper
  ```

- [ ] Download dependencies
  ```bash
  go mod download
  go mod tidy
  ```
  **Expected**: All dependencies downloaded without errors

### Day 2: Compilation
- [ ] Build the scraper
  ```bash
  go build -o bin/scraper cmd/scraper/main.go
  ```
  **Expected**: Binary created at `bin/scraper`
  **Potential Issues**: Import errors, missing dependencies

- [ ] Run basic syntax check
  ```bash
  go vet ./...
  ```
  **Expected**: No errors

- [ ] Format code
  ```bash
  go fmt ./...
  ```

### Day 3: Database Setup
- [ ] Install PostgreSQL 15+
- [ ] Create database
  ```sql
  CREATE DATABASE terminalnews;
  ```

- [ ] Run migrations from main project
  ```bash
  cd ../database/migrations
  # Apply migration files
  ```

- [ ] Verify tables exist
  ```sql
  \dt  -- Should show: articles, users, votes, classifieds, etc.
  ```

### Day 4: Configuration
- [ ] Create .env file
  ```bash
  cp .env.example .env
  ```

- [ ] Configure DATABASE_URL
  ```
  DATABASE_URL=postgres://postgres:password@localhost:5432/terminalnews?sslmode=disable
  ```

- [ ] (Optional) Add NewsAPI key
  ```
  NEWSAPI_KEY=get_free_key_from_newsapi.org
  ```

### Day 5: First Run
- [ ] Run the scraper
  ```bash
  ./bin/scraper
  ```

- [ ] Check logs for:
  - [ ] "Starting Terminal News Aggregator..."
  - [ ] Database connection successful
  - [ ] "Running initial article fetch..."
  - [ ] RSS feeds being fetched
  - [ ] Articles being stored

- [ ] Check database
  ```sql
  SELECT COUNT(*) FROM articles;  -- Should be > 0
  SELECT source, COUNT(*) FROM articles GROUP BY source;
  ```

**SUCCESS CRITERIA FOR WEEK 1:**
- ✅ Code compiles without errors
- ✅ Connects to database
- ✅ Fetches at least 100 articles
- ✅ No crashes during 1-hour run

---

## 🔨 WEEK 2: RSS PARSER VERIFICATION

### Day 6-7: RSS Feeds Testing
- [ ] Verify each RSS source
  ```bash
  # Check logs for each source
  grep "Fetching RSS feed" logs.txt
  grep "Parsed .* articles from" logs.txt
  ```

- [ ] Test sources individually:
  - [ ] TechCrunch - Should fetch ~20-30 articles
  - [ ] BBC News - Should fetch ~20-30 articles
  - [ ] Reuters - Should fetch ~20-30 articles
  - [ ] The Verge - Should fetch ~20-30 articles
  - [ ] Wired - Should fetch ~15-25 articles
  - [ ] (All 22 sources...)

- [ ] Check for errors
  ```bash
  grep "ERROR" logs.txt
  grep "Failed to fetch" logs.txt
  ```

### Day 8: Article Quality Check
- [ ] Verify article data quality
  ```sql
  -- Check for missing fields
  SELECT COUNT(*) FROM articles WHERE title IS NULL OR title = '';
  SELECT COUNT(*) FROM articles WHERE url IS NULL OR url = '';
  SELECT COUNT(*) FROM articles WHERE published_at IS NULL;

  -- Check categories
  SELECT category, COUNT(*) FROM articles GROUP BY category;
  -- Should see: tech, business, science, sports, entertainment, general
  ```

### Day 9: Deduplication Testing
- [ ] Check for duplicates by URL
  ```sql
  SELECT url, COUNT(*) FROM articles
  GROUP BY url HAVING COUNT(*) > 1;
  -- Should return 0 rows
  ```

- [ ] Check for duplicate titles (fuzzy)
  ```sql
  SELECT title, COUNT(*) FROM articles
  GROUP BY title HAVING COUNT(*) > 1;
  -- May have some legitimate duplicates from different sources
  ```

### Day 10: Performance Measurement
- [ ] Measure fetch time
  ```bash
  # Check logs for timing
  grep "Fetch complete" logs.txt
  # Should see: "Fetch complete: X articles in Y seconds"
  ```

- [ ] Target: <30 seconds for full fetch cycle
- [ ] Record memory usage
  ```bash
  # While running
  ps aux | grep scraper  # Check RSS memory
  ```

**SUCCESS CRITERIA FOR WEEK 2:**
- ✅ At least 500 articles fetched per day
- ✅ All 22 RSS sources working
- ✅ <1% duplicate rate
- ✅ Categories assigned to >80% of articles
- ✅ Fetch time <30 seconds

---

## 🔨 WEEK 3-4: API INTEGRATION

### Day 11-12: NewsAPI Testing
- [ ] Add NEWSAPI_KEY to .env
- [ ] Restart scraper
- [ ] Check logs for:
  ```
  NewsAPI client initialized successfully
  Cron: Fetching NewsAPI articles...
  Fetched X articles from NewsAPI
  ```

- [ ] Verify NewsAPI articles in database
  ```sql
  SELECT COUNT(*) FROM articles WHERE fetch_source = 'newsapi';
  -- Should be > 0
  ```

### Day 13-14: Article Volume
- [ ] Run for 24 hours
- [ ] Check total articles
  ```sql
  SELECT COUNT(*) FROM articles
  WHERE created_at > NOW() - INTERVAL '24 hours';
  -- Should be 1000-2000
  ```

- [ ] Check articles per source
  ```sql
  SELECT source, COUNT(*) as count
  FROM articles
  WHERE created_at > NOW() - INTERVAL '24 hours'
  GROUP BY source
  ORDER BY count DESC;
  ```

**SUCCESS CRITERIA FOR WEEK 3-4:**
- ✅ 1000+ articles per day
- ✅ NewsAPI contributing 200-400 articles/day
- ✅ All categories covered

---

## 🔨 WEEK 5-6: ACCURACY MEASUREMENT

### Day 15-16: Deduplication Accuracy
- [ ] Manual sampling test
  ```sql
  -- Get 100 random articles
  SELECT id, title, url FROM articles
  ORDER BY RANDOM() LIMIT 100;
  ```

- [ ] Manually check for duplicates
- [ ] Calculate accuracy: (True Negatives + True Positives) / Total
- [ ] Target: >95% accuracy

### Day 17-18: Classification Accuracy
- [ ] Sample 100 random articles
  ```sql
  SELECT id, title, category, source
  FROM articles
  ORDER BY RANDOM() LIMIT 100;
  ```

- [ ] Manually verify categories
- [ ] Record accuracy per category:
  - Tech: ____%
  - Business: ____%
  - Science: ____%
  - Sports: ____%
  - Entertainment: ____%
  - General: ____%

- [ ] Target: >80% overall accuracy

### Day 19-20: Performance Tuning
- [ ] Run with profiling
  ```bash
  go run -cpuprofile=cpu.prof cmd/scraper/main.go
  ```

- [ ] Analyze bottlenecks
  ```bash
  go tool pprof cpu.prof
  ```

- [ ] Optimize slow queries
- [ ] Reduce memory usage if needed

**SUCCESS CRITERIA FOR WEEK 5-6:**
- ✅ Deduplication >95% accurate
- ✅ Classification >80% accurate
- ✅ Fetch cycle <10 seconds for 1000 articles
- ✅ Memory usage <200MB sustained

---

## 🔨 WEEK 7-8: WEATHER & MODERATION

### Day 21-22: Weather Testing
- [ ] Seed cities table
  ```sql
  INSERT INTO cities (name, state, latitude, longitude)
  VALUES
    ('San Francisco', 'CA', 37.7749, -122.4194),
    ('New York', 'NY', 40.7128, -74.0060),
    ('Chicago', 'IL', 41.8781, -87.6298);
  ```

- [ ] Test weather client manually
  ```go
  // Add to main.go temporarily
  weatherClient := weather.NewWeatherClient(db)
  cities := []types.City{{ID: 1, Name: "San Francisco", Latitude: 37.7749, Longitude: -122.4194}}
  weatherClient.UpdateWeatherForCities(cities)
  ```

- [ ] Check weather_current table
  ```sql
  SELECT * FROM weather_current;
  ```

- [ ] Verify NOAA API responses

### Day 23-24: Spam Detection
- [ ] Test with known spam patterns
- [ ] Create test articles with spam keywords
- [ ] Verify moderator flags them
- [ ] Check false positive rate
  ```sql
  SELECT COUNT(*) FROM articles
  WHERE source IN ('BBC', 'Reuters', 'TechCrunch')
  AND is_flagged = true;
  -- Should be 0 or very low
  ```

**SUCCESS CRITERIA FOR WEEK 7-8:**
- ✅ Weather data updates for all seeded cities
- ✅ Spam detection catching >90% of test spam
- ✅ False positive rate <5%
- ✅ Trust scores calculated correctly

---

## 🔨 WEEK 9: PRODUCTION READINESS

### Day 25: Monitoring
- [ ] Check monitoring metrics
  ```sql
  SELECT * FROM fetch_logs
  WHERE created_at > NOW() - INTERVAL '7 days'
  ORDER BY created_at DESC;
  ```

- [ ] Verify source statistics
- [ ] Test alert system
  ```go
  // Temporarily break a source to test alerts
  ```

### Day 26: Testing
- [ ] Run unit tests
  ```bash
  go test ./...
  ```
  **Expected**: All tests pass

- [ ] Run test with coverage
  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```
  **Target**: >70% coverage

### Day 27: Docker
- [ ] Build Docker image
  ```bash
  docker build -t terminal-news-scraper .
  ```

- [ ] Run in Docker
  ```bash
  docker run --env-file .env terminal-news-scraper
  ```

- [ ] Verify it works in container

### Day 28: Load Testing
- [ ] Run for 24 hours continuously
- [ ] Monitor for:
  - Memory leaks
  - Connection pool exhaustion
  - Error rate increases
  - Performance degradation

### Day 29-30: Documentation & Cleanup
- [ ] Update README with actual metrics
- [ ] Document any quirks or issues found
- [ ] Create troubleshooting guide
- [ ] Update deployment docs

**SUCCESS CRITERIA FOR WEEK 9:**
- ✅ All tests passing
- ✅ Docker image builds and runs
- ✅ 24-hour stability test passed
- ✅ Documentation complete and accurate

---

## 📊 FINAL VERIFICATION

### Production Checklist
- [ ] Code compiles without errors
- [ ] All dependencies resolved
- [ ] Environment variables documented
- [ ] Database migrations tested
- [ ] Fetches 1000+ articles/day
- [ ] Deduplication >95% accurate
- [ ] Classification >80% accurate
- [ ] Weather integration working
- [ ] Spam detection functional
- [ ] Monitoring and alerts configured
- [ ] Docker deployment working
- [ ] 48-hour stability test passed
- [ ] Documentation complete

### Performance Benchmarks (Actual Measured)
```
RSS Fetch (22 sources):     ___ seconds
NewsAPI Fetch:               ___ seconds
Deduplication (1000 articles): ___ ms
Classification (1000 articles): ___ ms
Storage (1000 articles):     ___ ms
Total Pipeline:              ___ seconds

Memory Usage:
  Startup:    ___ MB
  Sustained:  ___ MB
  Peak:       ___ MB

Accuracy:
  Deduplication: ___% (target >95%)
  Classification: ___% (target >80%)
  Spam Detection: ___% (target >90%)
```

### Issues Found & Resolved
```
1. Issue: _______________
   Solution: _______________

2. Issue: _______________
   Solution: _______________
```

---

## 🎯 SIGN-OFF

When all items are checked:
- [ ] Dev 3 has completed all deliverables
- [ ] Code is production-ready
- [ ] Handoff to DevOps for deployment
- [ ] Documentation transferred to team

**Signed**: ________________
**Date**: ________________
**Status**: ☐ READY FOR PRODUCTION

---

## 📝 NOTES

Use this space to document anything unusual, gotchas, or improvements needed:

```
[Write notes here during testing]
```
