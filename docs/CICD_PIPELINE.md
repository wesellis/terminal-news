# CI/CD Pipeline - Zero-Touch Deployment

## Goal: Push to GitHub → Automatically Deploy to Production

---

## GitHub Actions Workflow

### `.github/workflows/deploy.yml`

```yaml
name: Build, Test, and Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

  build-and-push:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v3

      - name: Log in to Container Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=sha
            type=raw,value=latest

      - name: Build and push Docker image
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - name: Deploy to production
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.PRODUCTION_HOST }}
          username: ${{ secrets.PRODUCTION_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /opt/terminal-news
            docker-compose pull
            docker-compose up -d --no-deps --build api worker
            docker system prune -f

      - name: Run smoke tests
        run: |
          sleep 30
          curl -f https://api.terminalnews.com/health || exit 1

      - name: Notify Slack
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: 'Deployment ${{ job.status }}'
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}
        if: always()
```

---

## Automated Testing

### Unit Tests

```bash
# Run on every commit
go test ./...

# With coverage
go test -cover ./...

# Race detection
go test -race ./...
```

### Integration Tests

```bash
# Docker-based integration tests
docker-compose -f docker-compose.test.yml up -d
go test -tags=integration ./...
docker-compose -f docker-compose.test.yml down
```

### E2E Tests (Optional)

```bash
# Playwright or similar for CLI testing
./bin/terminal-news --test-mode
# Run automated interactions
```

---

## Deployment Strategies

### Strategy 1: Blue-Green Deployment

```yaml
# docker-compose.yml
services:
  api-blue:
    image: ghcr.io/terminalnews/api:latest
    # ... config

  api-green:
    image: ghcr.io/terminalnews/api:latest
    # ... config

  nginx:
    # Routes traffic to blue or green
```

**Process:**
1. Deploy to green (inactive)
2. Run smoke tests
3. Switch nginx to green
4. Keep blue as rollback option
5. After 1 hour: Remove blue

### Strategy 2: Rolling Update (Simpler)

```bash
# Update one container at a time
docker-compose up -d --no-deps --scale api=2 api
# Wait for health check
docker-compose up -d --no-deps --scale api=1 api
```

### Strategy 3: Simple Replace (Start)

```bash
# Pull new image
docker-compose pull api

# Restart with new image
docker-compose up -d api

# Auto-rollback if health check fails
```

---

## Automated Rollback

### Health Check Endpoint

```go
// /health endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    // Check database
    err := db.Ping()
    if err != nil {
        http.Error(w, "Database down", 500)
        return
    }

    // Check Redis
    _, err = redis.Ping().Result()
    if err != nil {
        http.Error(w, "Redis down", 500)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
        "version": VERSION,
    })
}
```

### Rollback Script

```bash
#!/bin/bash
# rollback.sh

# Check if deployment healthy
if ! curl -f https://api.terminalnews.com/health; then
    echo "Health check failed! Rolling back..."

    # Get previous image
    PREVIOUS=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep terminalnews/api | sed -n 2p)

    # Rollback
    docker tag $PREVIOUS ghcr.io/terminalnews/api:latest
    docker-compose up -d api

    # Alert
    curl -X POST $SLACK_WEBHOOK -d '{"text":"🚨 Auto-rollback triggered!"}'

    exit 1
fi

echo "Deployment successful ✅"
```

---

## Database Migrations

### Automated with Caution

```yaml
# In GitHub Actions
- name: Run migrations
  run: |
    # Connect to production DB
    ./migrate -path migrations -database $DATABASE_URL up

    # Run with transaction (rollback on error)
```

**Better approach:**
```bash
# Manual migration approval
# 1. Deploy code (backward compatible)
# 2. Run migration manually
# 3. Verify
# 4. Deploy new code using new schema
```

---

## Monitoring Deployment

### Automated Checks

```yaml
- name: Post-deployment checks
  run: |
    # API responding?
    curl -f https://api.terminalnews.com/health

    # Database accessible?
    curl -f https://api.terminalnews.com/db/health

    # No error spike?
    ./scripts/check_error_rate.sh

    # Latency acceptable?
    ./scripts/check_latency.sh
```

### Sentry Release Tracking

```yaml
- name: Create Sentry release
  run: |
    sentry-cli releases new -p terminal-news ${{ github.sha }}
    sentry-cli releases set-commits ${{ github.sha }} --auto
    sentry-cli releases finalize ${{ github.sha }}
```

---

## Secrets Management

### GitHub Secrets

```
Settings → Secrets and variables → Actions

Add secrets:
- PRODUCTION_HOST (IP address)
- PRODUCTION_USER (ssh user)
- SSH_PRIVATE_KEY (for deployment)
- SLACK_WEBHOOK (notifications)
- DATABASE_URL (if needed)
- STRIPE_SECRET_KEY (if deploying payments)
```

### In Production (Environment Variables)

```bash
# /opt/terminal-news/.env
DATABASE_URL=postgres://...
REDIS_URL=redis://...
STRIPE_SECRET_KEY=sk_live_...
JWT_SECRET=...

# Load via docker-compose
env_file:
  - .env
```

---

## Automated Notifications

### Slack Integration

```yaml
- name: Notify deployment start
  run: |
    curl -X POST ${{ secrets.SLACK_WEBHOOK }} \
      -d '{"text":"🚀 Deploying to production..."}'

- name: Notify deployment success
  if: success()
  run: |
    curl -X POST ${{ secrets.SLACK_WEBHOOK }} \
      -d '{"text":"✅ Deployment successful!"}'

- name: Notify deployment failure
  if: failure()
  run: |
    curl -X POST ${{ secrets.SLACK_WEBHOOK }} \
      -d '{"text":"❌ Deployment FAILED! Rolling back..."}'
```

---

## Deployment Frequency

**Goal: Deploy multiple times per day**

- Every push to main = deploy
- Average: 2-5 deploys/day
- No manual intervention
- Rollback automatic if issues

**This is how you move fast.**

---

## Cost

**GitHub Actions:**
- Free for public repos
- $0.008/minute for private repos
- Typical deploy: 5 minutes
- Cost: $0.04 per deploy
- 100 deploys/month = $4

**Worth it:** Saves hours of manual deployment

---

## Summary

**Automation:**
✅ Push code → Auto-test → Auto-deploy
✅ Health checks run automatically
✅ Rollback if anything fails
✅ Slack notifications
✅ Zero-downtime deploys

**Manual work:** 0 hours (maybe review logs)

**Deployment time:** 5 minutes (automated)
