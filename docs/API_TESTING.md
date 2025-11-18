# API Testing Guide

Quick start guide for testing the Terminal News backend API.

## Prerequisites

- PostgreSQL 15+ running
- Redis 7+ running
- Go 1.21+ installed
- Backend server running

## Quick Start

### 1. Start Backend Services

```bash
# Start PostgreSQL and Redis (using Docker)
cd terminal-news
docker-compose -f docker-compose.dev.yml up -d postgres redis

# Navigate to backend
cd backend

# Run migrations
psql $DATABASE_URL -f ../database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f ../database/migrations/002_triggers_and_functions.sql

# Start server
go run cmd/server/main.go
```

Server will start on http://localhost:8080

### 2. Seed Test Data

```bash
# In another terminal, from backend/ directory
go run cmd/seeder/main.go
```

This creates:
- 5 test users (testuser, alice, bob, charlie, diana)
- 100 test articles across categories
- Random votes on articles
- Comments and replies
- 50 test classifieds

**Test credentials:**
- Username: `testuser` / Password: `password123`
- Username: `alice` / Password: `password123`
- Username: `bob` / Password: `password123`

### 3. Test API Endpoints

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Register New User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "new@example.com",
    "password": "password123"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

Save the `access_token` from the response!

#### Get Hot Articles
```bash
curl http://localhost:8080/api/v1/articles/hot?limit=10
```

#### Get Article by ID
```bash
curl http://localhost:8080/api/v1/articles/1
```

#### Vote on Article (requires auth)
```bash
curl -X POST http://localhost:8080/api/v1/articles/1/vote \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"vote_type": "like"}'
```

#### Get Comments
```bash
curl http://localhost:8080/api/v1/articles/1/comments
```

#### Post Comment (requires auth)
```bash
curl -X POST http://localhost:8080/api/v1/articles/1/comments \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "Great article!"}'
```

#### Get Classifieds
```bash
curl "http://localhost:8080/api/v1/classifieds?limit=10"
```

#### Search Classifieds by Location
```bash
curl "http://localhost:8080/api/v1/classifieds?lat=45.5152&lng=-122.6784&radius=25"
```

#### Get Weather
```bash
curl "http://localhost:8080/api/v1/weather?lat=45.5152&lng=-122.6784"
```

### 4. Test WebSocket

#### Using wscat (Node.js)
```bash
# Install wscat
npm install -g wscat

# Connect without auth
wscat -c ws://localhost:8080/ws

# Connect with auth
wscat -c "ws://localhost:8080/ws?token=YOUR_ACCESS_TOKEN"
```

#### Using Browser Console
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => console.log('Connected!');
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};
ws.onerror = (error) => console.error('Error:', error);
```

### 5. Test Payments (Stripe)

**Note:** Requires Stripe API keys in .env

#### Create Classified Boost Payment
```bash
curl -X POST http://localhost:8080/api/v1/payments/create-intent \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "classified_boost",
    "classified_id": 1,
    "duration_days": 7
  }'
```

#### Create Sponsor Subscription
```bash
curl -X POST http://localhost:8080/api/v1/payments/create-intent \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "sponsor_subscription",
    "tier": "bronze"
  }'
```

#### Get Payment History
```bash
curl http://localhost:8080/api/v1/payments/history \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Testing Checklist

### Authentication ✓
- [ ] Register new user
- [ ] Login with existing user
- [ ] Refresh access token
- [ ] Access protected endpoint with token
- [ ] Access protected endpoint without token (should fail)

### Articles ✓
- [ ] Get hot articles
- [ ] Get controversial articles
- [ ] Get rising articles
- [ ] Get single article by ID
- [ ] Pagination works (offset/limit)

### Voting ✓
- [ ] Vote "like" on article
- [ ] Vote "dislike" on article
- [ ] Vote "open" on article
- [ ] Remove vote
- [ ] Vote count updates correctly

### Comments ✓
- [ ] Get comments for article (tree structure)
- [ ] Post top-level comment
- [ ] Post reply to comment
- [ ] Update own comment
- [ ] Delete own comment
- [ ] Cannot update/delete other's comments

### Classifieds ✓
- [ ] List all classifieds
- [ ] Filter by category
- [ ] Filter by city/state
- [ ] Geographic search (lat/lng/radius)
- [ ] Create classified
- [ ] Update own classified
- [ ] Delete own classified
- [ ] Cannot update/delete other's classifieds

### Weather ✓
- [ ] Get weather by coordinates
- [ ] Returns current conditions
- [ ] Returns forecast

### WebSocket ✓
- [ ] Connect without auth
- [ ] Connect with auth token
- [ ] Receive real-time updates
- [ ] Auto-reconnect on disconnect

### Payments ✓
- [ ] Create classified boost payment
- [ ] Create sponsor subscription
- [ ] View payment history
- [ ] Webhook receives events (need ngrok for testing)

## Rate Limiting

The API has rate limiting enabled (100 requests/minute per IP). Check headers:
```bash
curl -i http://localhost:8080/api/v1/articles/hot
```

Look for:
- `X-RateLimit-Limit: 100`
- `X-RateLimit-Remaining: 99`
- `X-RateLimit-Reset: 1234567890`

## Common Issues

### "connect: connection refused"
- PostgreSQL or Redis not running
- Check `DATABASE_URL` and `REDIS_URL` in .env

### "database does not exist"
```bash
createdb terminalnews_dev
# Then run migrations
```

### "JWT_SECRET not set"
```bash
cp .env.example .env
# Edit .env and set JWT_SECRET
```

### "Stripe error"
- Check `STRIPE_SECRET_KEY` in .env
- Use test mode keys from Stripe dashboard

## Integration Testing for Dev 2 (CLI)

The CLI should be able to:

1. **Authenticate**
   - Register new users
   - Login and store tokens
   - Refresh expired tokens

2. **Browse Articles**
   - Fetch hot/controversial/rising feeds
   - Display article metadata
   - Handle pagination

3. **Interact**
   - Vote on articles
   - Post and read comments
   - Navigate comment trees

4. **Classifieds**
   - Browse listings
   - Post new classified
   - Edit/delete own listings

5. **Real-time**
   - Connect to WebSocket
   - Receive live updates
   - Auto-reconnect on disconnect

6. **Offline Mode**
   - Cache articles locally
   - Queue actions when offline
   - Sync when reconnected

## Integration Testing for Dev 3 (Scraper)

The scraper should be able to:

1. **Insert Articles**
   - Use the articles table schema
   - Set correct categories
   - Include all metadata fields

2. **Avoid Duplicates**
   - Check existing articles by URL
   - Use deduplication logic

3. **Performance**
   - Insert 1000+ articles efficiently
   - Use connection pooling
   - Batch inserts where possible

## Monitoring

### Check Logs
```bash
# Backend logs show all requests
tail -f backend/logs/server.log
```

### Check Redis Cache
```bash
redis-cli
> KEYS *
> GET "cache:articles:hot"
> TTL "cache:articles:hot"
```

### Check Database
```bash
psql $DATABASE_URL

# Count articles
SELECT COUNT(*) FROM articles;

# Check recent articles
SELECT title, source, published_at
FROM articles
ORDER BY published_at DESC
LIMIT 10;

# Check vote counts
SELECT a.title, COUNT(v.id) as vote_count
FROM articles a
LEFT JOIN votes v ON a.id = v.article_id
GROUP BY a.id
ORDER BY vote_count DESC
LIMIT 10;
```

## Performance Testing

### Load Test with Apache Bench
```bash
# Test articles endpoint
ab -n 1000 -c 10 http://localhost:8080/api/v1/articles/hot

# Test with auth
ab -n 1000 -c 10 -H "Authorization: Bearer TOKEN" \
  http://localhost:8080/api/v1/user/profile
```

### Expected Performance
- Articles endpoint: <50ms avg response time
- With Redis cache: <10ms avg response time
- WebSocket connections: 1000+ concurrent
- Rate limiting: 100 req/min sustained

## Next Steps

1. Run through all testing checklist items
2. Report any errors or unexpected behavior
3. Test edge cases (invalid input, missing fields, etc.)
4. Test with CLI client
5. Test with scraper integration
6. Load test with realistic traffic
7. Security test (SQL injection, XSS, etc.)

## Support

- Backend README: `backend/README.md`
- API Endpoints: See backend/README.md for full list
- Database Schema: `database/migrations/001_initial_schema.sql`
- Issues: Report in project issue tracker
