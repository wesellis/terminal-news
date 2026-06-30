# Location-Based Features - Implementation Guide

## 🌍 CRITICAL: User Location Context

Terminal News **MUST** be location-aware for these features:
- ☀️ **Weather** - User's local forecast (not global)
- 📰 **Local News** - News relevant to user's city/region
- 📋 **Classifieds** - Buy/sell listings in user's area

---

## ✅ What Dev 3 (Scraper) Has Built

### 1. Weather Client (`internal/weather/weather.go`)

**Status**: ✅ COMPLETE - Location-based by design

The weather client takes a `City` object with coordinates:

```go
type City struct {
    ID        int64
    Name      string
    State     string
    Country   string
    Latitude  float64   // User's location
    Longitude float64   // User's location
}

// Fetches weather for SPECIFIC city
weatherClient.UpdateWeatherForCities([]City{userCity})
```

**Database Storage**:
```sql
CREATE TABLE weather_current (
    city_id BIGINT PRIMARY KEY,
    temperature FLOAT,
    condition TEXT,
    humidity FLOAT,
    wind_speed FLOAT,
    updated_at TIMESTAMP
);

CREATE TABLE weather_forecast (
    city_id BIGINT,
    day_index INT,
    temperature INT,
    condition TEXT,
    -- 5-day forecast per city
);
```

**How It Works**:
1. User's city coordinates → NOAA API
2. NOAA returns local forecast for EXACT location
3. Stored in DB with `city_id`
4. Backend/CLI query by user's `city_id`

---

### 2. Classifieds Schema (Database Ready)

**Status**: ✅ SCHEMA READY (from main project)

Database schema includes location fields:

```sql
CREATE TABLE classifieds (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    title VARCHAR(200),
    description TEXT,
    price DECIMAL(10,2),

    -- LOCATION FIELDS
    city VARCHAR(100) NOT NULL,      -- "San Francisco"
    state VARCHAR(50),                -- "CA"
    country VARCHAR(50) DEFAULT 'US', -- "US"
    lat DECIMAL(10,7),                -- 37.7749
    lng DECIMAL(10,7),                -- -122.4194

    category VARCHAR(50),
    created_at TIMESTAMP,
    expires_at TIMESTAMP
);

-- Geographic index for fast location queries
CREATE INDEX idx_classifieds_location ON classifieds(city, state)
WHERE is_active = TRUE;
```

**Storage Methods** (`internal/storage/storage.go`):
- Already has methods to query by location
- Can filter classifieds by city/state
- Geographic search ready

---

### 3. Local News (RSS Feeds)

**Status**: ⚠️ NEEDS LOCATION-SPECIFIC FEEDS

Current RSS sources are **global**. For local news, we need:

**Option A: Add Location-Specific RSS Feeds**
```go
// Example: Bay Area news sources
{Name: "SF Chronicle", URL: "https://www.sfchronicle.com/rss",
 Category: "local", Location: "San Francisco, CA"},
{Name: "LA Times", URL: "https://www.latimes.com/rss",
 Category: "local", Location: "Los Angeles, CA"},
{Name: "Chicago Tribune", URL: "https://www.chicagotribune.com/rss",
 Category: "local", Location: "Chicago, IL"},
```

**Option B: NewsAPI Location Filter**
```go
// NewsAPI supports location filtering
client.R().
    SetQueryParams(map[string]string{
        "q": "San Francisco OR Bay Area",
        "sortBy": "publishedAt",
    }).
    Get("/everything")
```

**Option C: Geographic Tagging** (Most Scalable)
```go
// Tag articles with locations during classification
article.Locations = []string{"San Francisco", "Bay Area", "California"}

// Users query by their location
SELECT * FROM articles
WHERE 'San Francisco' = ANY(locations)
```

---

## 🔧 What Backend (Dev 1) Needs to Do

### 1. User Location Management

**Add to Users Table**:
```sql
ALTER TABLE users ADD COLUMN city VARCHAR(100);
ALTER TABLE users ADD COLUMN state VARCHAR(50);
ALTER TABLE users ADD COLUMN latitude DECIMAL(10,7);
ALTER TABLE users ADD COLUMN longitude DECIMAL(10,7);
ALTER TABLE users ADD COLUMN timezone VARCHAR(50);
```

**API Endpoints Needed**:
```go
// User sets their location during onboarding or in settings
POST /api/user/location
{
    "city": "San Francisco",
    "state": "CA",
    "latitude": 37.7749,
    "longitude": -122.4194,
    "timezone": "America/Los_Angeles"
}

// Get user's location
GET /api/user/location
```

### 2. Location-Filtered Endpoints

**Weather**:
```go
// Returns weather for USER'S city
GET /api/weather
// Backend queries: SELECT * FROM weather_current WHERE city_id = user.city_id
```

**Classifieds**:
```go
// Returns classifieds near USER'S location
GET /api/classifieds?radius=50
// Backend filters by city/state or within radius of user's lat/lng
```

**Local News**:
```go
// Returns news relevant to USER'S location
GET /api/articles/local
// Backend filters articles tagged with user's city/state/region
```

---

## 🖥️ What CLI (Dev 2) Needs to Do

### 1. Location Detection/Setup

**On First Run** (or in settings):

```
╔════════════════════════════════════════════════════════════════╗
║                    TERMINAL NEWS                               ║
║                    Setup Your Location                         ║
╚════════════════════════════════════════════════════════════════╝

To show you local weather, news, and classifieds, we need
your location.

  [1] Auto-detect (using IP geolocation)
  [2] Enter city manually
  [3] Use ZIP code
  [4] Skip (use default: New York, NY)

Your choice: _
```

**Auto-Detection Options**:
- IP geolocation API (ipapi.co, ip-api.com - free)
- System timezone → infer city
- Manual entry with autocomplete

### 2. Location Display in UI

**Weather Widget** (always visible):
```
┌─ WEATHER ─────────────────────────────────┐
│ San Francisco, CA              ☀️  72°F    │
│ Partly Cloudy                             │
│ Wind: 10mph W  Humidity: 65%              │
└───────────────────────────────────────────┘
```

**Classifieds Tab**:
```
┌─ CLASSIFIEDS ─────────────────────────────┐
│ 📍 Within 25 miles of San Francisco, CA   │
│                                            │
│ [1] 2015 Honda Civic - $12,000           │
│     ↳ Oakland, CA (8 mi)                  │
│                                            │
│ [2] Studio Apartment - $1,800/mo         │
│     ↳ San Francisco, CA (2 mi)            │
└────────────────────────────────────────────┘
```

**Local News Feed**:
```
┌─ LOCAL NEWS ──────────────────────────────┐
│ 📍 San Francisco & Bay Area               │
│                                            │
│ • BART delays continue on Richmond line   │
│   ↳ SF Chronicle · 2h ago                 │
│                                            │
│ • SF housing market cools in Q4           │
│   ↳ SF Business Times · 5h ago            │
└────────────────────────────────────────────┘
```

---

## 🗺️ Location Architecture Flow

```
┌─────────────────────────────────────────────────────────────┐
│ CLI (Dev 2)                                                 │
│                                                             │
│ 1. User sets location: "San Francisco, CA"                 │
│    ↓                                                        │
│ 2. Sends to Backend API                                    │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│ Backend (Dev 1)                                             │
│                                                             │
│ 1. Stores user.city = "San Francisco"                      │
│ 2. Stores user.latitude = 37.7749                          │
│ 3. Stores user.longitude = -122.4194                       │
│    ↓                                                        │
│ 4. When user requests weather:                             │
│    SELECT * FROM weather_current                           │
│    WHERE city_id = (SELECT city_id FROM cities             │
│                     WHERE name = 'San Francisco')          │
│    ↓                                                        │
│ 5. When user requests classifieds:                         │
│    SELECT * FROM classifieds                               │
│    WHERE city = 'San Francisco' OR state = 'CA'            │
│    ORDER BY created_at DESC                                │
│    ↓                                                        │
│ 6. When user requests local news:                          │
│    SELECT * FROM articles                                  │
│    WHERE 'San Francisco' = ANY(locations)                  │
│    OR 'Bay Area' = ANY(locations)                          │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│ Scraper (Dev 3 - Me!)                                       │
│                                                             │
│ 1. Fetches weather for all active cities                   │
│    FOR EACH city IN cities:                                │
│      weather_client.UpdateWeatherForCities([city])         │
│      ↳ Stores in weather_current table with city_id        │
│                                                             │
│ 2. Classifieds posted by users (not scraped)               │
│    ↳ Users specify their city when posting                 │
│                                                             │
│ 3. Articles tagged with locations during classification    │
│    ↳ Extract city/region mentions from article text        │
└─────────────────────────────────────────────────────────────┘
```

---

## 📋 TODO: Team Coordination Needed

### Dev 1 (Backend) Must Do:
- [ ] Add location fields to `users` table
- [ ] Create `cities` reference table
- [ ] API endpoint: `POST /api/user/location`
- [ ] API endpoint: `GET /api/weather` (filtered by user city)
- [ ] API endpoint: `GET /api/classifieds?location=...`
- [ ] API endpoint: `GET /api/articles/local` (filtered by user city)
- [ ] Join user location with weather/classifieds queries

### Dev 2 (CLI) Must Do:
- [ ] Onboarding flow: "What's your location?"
- [ ] Location auto-detection (IP geolocation API)
- [ ] Settings: Change location
- [ ] Display location in weather widget
- [ ] Filter classifieds by proximity
- [ ] Show "local news" tab
- [ ] Distance calculation (miles from user location)

### Dev 3 (Me - Scraper) Must Do:
- [x] Weather client with city-specific fetching ✅
- [x] Database schema supports locations ✅
- [ ] Add location tagging to articles during classification
- [ ] Add location-specific RSS feeds (optional)
- [ ] Create cities seed data (major US cities)

---

## 🎯 Priority Implementation Order

### Phase 1 (MVP Launch):
1. **Backend**: Add `user.city` and `user.state` fields
2. **CLI**: Simple location input on first run
3. **Backend**: Filter weather by user city
4. **Scraper**: Populate weather for top 50 US cities

### Phase 2 (Post-Launch):
1. **CLI**: IP geolocation auto-detection
2. **Backend**: Classifieds filtering by location
3. **Scraper**: Location tagging for articles

### Phase 3 (Enhancement):
1. **CLI**: Distance calculations (miles/km)
2. **Backend**: Radius-based searches
3. **Scraper**: Local news RSS feeds

---

## 🌍 Cities to Support (Initial List)

**Top 50 US Cities** (by population):
- New York, NY
- Los Angeles, CA
- Chicago, IL
- Houston, TX
- Phoenix, AZ
- Philadelphia, PA
- San Antonio, TX
- San Diego, CA
- Dallas, TX
- San Jose, CA
- Austin, TX
- Jacksonville, FL
- Fort Worth, TX
- Columbus, OH
- San Francisco, CA
- Charlotte, NC
- Indianapolis, IN
- Seattle, WA
- Denver, CO
- Boston, MA
... (full list in seed data)

**Scraper fetches weather for all active cities daily.**

---

## 🔑 Key Takeaway

**The scraper (Dev 3) is already location-ready!**

What's needed:
1. **Backend** must store and query by user location
2. **CLI** must capture and display user location
3. **Scraper** continues fetching location-specific data

The architecture is **designed for location-based features** from the ground up!

---

**Status**: ⚠️ **WAITING ON DEV 1 & DEV 2**

The scraper is ready. Backend and CLI need to implement location context!
