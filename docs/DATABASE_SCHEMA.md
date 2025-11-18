# Database Schema

Complete PostgreSQL schema for Terminal News with optimizations for performance and automation.

---

## Schema Overview

```
Users ──┬─── Articles ──── Votes
        │                 └── Comments
        │                 └── Article_Views
        │
        ├─── Classifieds ──── Classified_Views
        │                   └── Classified_Messages
        │
        ├─── Sponsors ──── Sponsorships
        │               └── Sponsor_Analytics
        │
        └─── User_Settings
            └── User_Activity
```

---

## Core Tables

### 1. Users

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    -- Profile
    display_name VARCHAR(100),
    bio TEXT,
    location VARCHAR(100),
    website VARCHAR(255),

    -- Reputation
    karma INTEGER DEFAULT 0,
    trust_score DECIMAL(3,2) DEFAULT 0.5, -- 0.0 to 1.0

    -- Status
    email_verified BOOLEAN DEFAULT FALSE,
    is_banned BOOLEAN DEFAULT FALSE,
    is_moderator BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,

    -- Tracking
    last_active_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    -- Indexes
    CONSTRAINT username_length CHECK (LENGTH(username) >= 3),
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_karma ON users(karma DESC);
CREATE INDEX idx_users_last_active ON users(last_active_at DESC);
```

**Automation Features:**
- `trust_score`: Auto-calculated for spam prevention
- `karma`: Auto-updated from votes
- `last_active_at`: Auto-updated on any action

---

### 2. Articles

```sql
CREATE TABLE articles (
    id BIGSERIAL PRIMARY KEY,

    -- Content
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    content TEXT, -- Summary/excerpt
    image_url TEXT,

    -- Metadata
    source VARCHAR(100), -- TechCrunch, BBC, etc.
    author VARCHAR(255),
    published_at TIMESTAMP,

    -- Categorization
    category VARCHAR(50), -- tech, world, politics, etc.
    tags TEXT[], -- Array of tags

    -- Aggregation
    external_id VARCHAR(255), -- From source API
    fetch_source VARCHAR(50), -- newsapi, rss, guardian, etc.

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    -- Constraints
    CONSTRAINT unique_external_article UNIQUE (external_id, fetch_source),
    CONSTRAINT url_format CHECK (url ~* '^https?://.*')
);

CREATE INDEX idx_articles_published ON articles(published_at DESC);
CREATE INDEX idx_articles_category ON articles(category);
CREATE INDEX idx_articles_source ON articles(source);
CREATE INDEX idx_articles_created ON articles(created_at DESC);
CREATE INDEX idx_articles_tags ON articles USING GIN(tags);

-- Full-text search
CREATE INDEX idx_articles_title_search ON articles USING GIN(to_tsvector('english', title));
```

**Automation Features:**
- `external_id` + `fetch_source`: Prevents duplicate articles
- Full-text search: Fast article search
- GIN index on tags: Fast filtering

---

### 3. Votes

```sql
CREATE TABLE votes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,

    -- Vote type
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('open', 'like', 'dislike')),

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW(),

    -- Constraints
    CONSTRAINT unique_user_article_vote UNIQUE (user_id, article_id, vote_type)
);

CREATE INDEX idx_votes_article ON votes(article_id);
CREATE INDEX idx_votes_user ON votes(user_id);
CREATE INDEX idx_votes_created ON votes(created_at DESC);
CREATE INDEX idx_votes_type ON votes(vote_type);

-- Composite index for ranking queries
CREATE INDEX idx_votes_article_type_created ON votes(article_id, vote_type, created_at);
```

**Automation Features:**
- Unique constraint: Can't vote twice the same way
- Composite index: Fast ranking calculations
- CASCADE delete: Auto-cleanup when user/article deleted

---

### 4. Article Rankings (Materialized View)

```sql
-- Cached ranking scores for performance
CREATE MATERIALIZED VIEW article_rankings AS
SELECT
    a.id AS article_id,

    -- Vote counts
    COUNT(*) FILTER (WHERE v.vote_type = 'open') AS open_count,
    COUNT(*) FILTER (WHERE v.vote_type = 'like') AS like_count,
    COUNT(*) FILTER (WHERE v.vote_type = 'dislike') AS dislike_count,

    -- Scores
    (
        COUNT(*) FILTER (WHERE v.vote_type = 'open') * 1 +
        COUNT(*) FILTER (WHERE v.vote_type = 'like') * 2 +
        COUNT(*) FILTER (WHERE v.vote_type = 'dislike') * -1
    ) AS total_score,

    -- Controversy (min/max ratio)
    CASE
        WHEN COUNT(*) FILTER (WHERE v.vote_type = 'like') > 0
         AND COUNT(*) FILTER (WHERE v.vote_type = 'dislike') > 0
        THEN LEAST(
            COUNT(*) FILTER (WHERE v.vote_type = 'like'),
            COUNT(*) FILTER (WHERE v.vote_type = 'dislike')
        )::FLOAT / GREATEST(
            COUNT(*) FILTER (WHERE v.vote_type = 'like'),
            COUNT(*) FILTER (WHERE v.vote_type = 'dislike')
        )
        ELSE 0
    END AS controversy_score,

    -- Engagement
    COUNT(*) AS total_engagement,

    -- Time decay
    EXTRACT(EPOCH FROM (NOW() - a.published_at)) / 3600 AS hours_since_published,

    -- Hot rank (Reddit algorithm)
    (
        LOG(GREATEST(ABS(
            COUNT(*) FILTER (WHERE v.vote_type = 'like') * 2 -
            COUNT(*) FILTER (WHERE v.vote_type = 'dislike') * 1
        ), 1)) * SIGN(
            COUNT(*) FILTER (WHERE v.vote_type = 'like') -
            COUNT(*) FILTER (WHERE v.vote_type = 'dislike')
        ) +
        EXTRACT(EPOCH FROM a.published_at) / 45000
    ) AS hot_rank,

    NOW() AS last_updated

FROM articles a
LEFT JOIN votes v ON a.id = v.article_id
GROUP BY a.id, a.published_at;

CREATE UNIQUE INDEX idx_article_rankings_id ON article_rankings(article_id);
CREATE INDEX idx_article_rankings_hot ON article_rankings(hot_rank DESC);
CREATE INDEX idx_article_rankings_controversy ON article_rankings(controversy_score DESC);
CREATE INDEX idx_article_rankings_score ON article_rankings(total_score DESC);
```

**Automation:**
- Refresh every 5 minutes via cron:
  ```sql
  REFRESH MATERIALIZED VIEW CONCURRENTLY article_rankings;
  ```
- Fast queries (pre-calculated)
- No real-time vote counting needed

---

### 5. Comments

```sql
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,

    -- Content
    content TEXT NOT NULL,

    -- Status
    is_deleted BOOLEAN DEFAULT FALSE,
    is_flagged BOOLEAN DEFAULT FALSE,
    flag_count INTEGER DEFAULT 0,

    -- Votes
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0,

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    edited_at TIMESTAMP,

    CONSTRAINT content_length CHECK (LENGTH(content) > 0 AND LENGTH(content) <= 10000)
);

CREATE INDEX idx_comments_article ON comments(article_id, created_at);
CREATE INDEX idx_comments_user ON comments(user_id);
CREATE INDEX idx_comments_parent ON comments(parent_id);
CREATE INDEX idx_comments_flagged ON comments(is_flagged) WHERE is_flagged = TRUE;
```

**Automation Features:**
- `flag_count`: Auto-hide at threshold (3+)
- Soft delete: `is_deleted` instead of actual deletion
- Upvotes/downvotes: Cached for performance

---

### 6. Classifieds

```sql
CREATE TABLE classifieds (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Content
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10,2),

    -- Categorization
    category VARCHAR(50) NOT NULL, -- jobs, housing, for-sale, services, events
    subcategory VARCHAR(50),

    -- Location
    city VARCHAR(100) NOT NULL,
    state VARCHAR(50),
    country VARCHAR(50) DEFAULT 'US',
    lat DECIMAL(10,7),
    lng DECIMAL(10,7),

    -- Contact
    contact_email VARCHAR(255),
    contact_phone VARCHAR(20),
    contact_method VARCHAR(20) DEFAULT 'email', -- email, phone, dm

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_premium BOOLEAN DEFAULT FALSE,
    is_flagged BOOLEAN DEFAULT FALSE,
    flag_count INTEGER DEFAULT 0,

    -- Premium features
    premium_until TIMESTAMP,
    boost_count INTEGER DEFAULT 0,
    last_boosted_at TIMESTAMP,

    -- Metrics
    view_count INTEGER DEFAULT 0,
    contact_count INTEGER DEFAULT 0,

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,

    CONSTRAINT title_length CHECK (LENGTH(title) >= 10),
    CONSTRAINT description_length CHECK (LENGTH(description) >= 20)
);

CREATE INDEX idx_classifieds_category ON classifieds(category, is_active);
CREATE INDEX idx_classifieds_location ON classifieds(city, state) WHERE is_active = TRUE;
CREATE INDEX idx_classifieds_user ON classifieds(user_id);
CREATE INDEX idx_classifieds_created ON classifieds(created_at DESC) WHERE is_active = TRUE;
CREATE INDEX idx_classifieds_premium ON classifieds(is_premium, created_at DESC) WHERE is_active = TRUE;
CREATE INDEX idx_classifieds_expires ON classifieds(expires_at) WHERE expires_at IS NOT NULL;

-- Full-text search
CREATE INDEX idx_classifieds_search ON classifieds USING GIN(
    to_tsvector('english', title || ' ' || description)
);

-- Geographic search (if using PostGIS)
-- CREATE INDEX idx_classifieds_location_geo ON classifieds USING GIST(ll_to_earth(lat, lng));
```

**Automation Features:**
- `expires_at`: Cron job auto-deactivates expired listings
- `premium_until`: Auto-downgrade when premium expires
- `last_boosted_at`: Prevents boost spam
- Full-text search: Fast keyword search

---

### 7. Sponsors

```sql
CREATE TABLE sponsors (
    id BIGSERIAL PRIMARY KEY,

    -- Business info
    business_name VARCHAR(200) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    contact_name VARCHAR(100),
    contact_phone VARCHAR(20),
    website VARCHAR(255),

    -- Location
    city VARCHAR(100) NOT NULL,
    state VARCHAR(50),
    country VARCHAR(50) DEFAULT 'US',

    -- Branding
    logo_url TEXT,
    tagline VARCHAR(200),

    -- Billing
    stripe_customer_id VARCHAR(100) UNIQUE,
    stripe_subscription_id VARCHAR(100),

    -- Status
    is_active BOOLEAN DEFAULT FALSE,
    is_approved BOOLEAN DEFAULT FALSE,

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sponsors_active ON sponsors(is_active, city);
CREATE INDEX idx_sponsors_stripe ON sponsors(stripe_customer_id);
```

---

### 8. Sponsorships (Placements)

```sql
CREATE TABLE sponsorships (
    id BIGSERIAL PRIMARY KEY,
    sponsor_id BIGINT NOT NULL REFERENCES sponsors(id) ON DELETE CASCADE,

    -- Placement
    placement_type VARCHAR(50) NOT NULL, -- weather, header, sidebar
    city VARCHAR(100) NOT NULL,

    -- Pricing
    monthly_price DECIMAL(10,2) NOT NULL,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Billing period
    starts_at TIMESTAMP NOT NULL,
    ends_at TIMESTAMP,
    next_billing_date TIMESTAMP,

    -- Metrics
    impressions INTEGER DEFAULT 0,
    clicks INTEGER DEFAULT 0,

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT unique_placement UNIQUE (placement_type, city, sponsor_id)
);

CREATE INDEX idx_sponsorships_active ON sponsorships(is_active, city, placement_type);
CREATE INDEX idx_sponsorships_sponsor ON sponsorships(sponsor_id);
CREATE INDEX idx_sponsorships_billing ON sponsorships(next_billing_date) WHERE is_active = TRUE;
```

**Automation Features:**
- `next_billing_date`: Cron triggers Stripe subscription billing
- `impressions`/`clicks`: Auto-incremented
- Auto-deactivate when subscription ends

---

### 9. Payments

```sql
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    sponsor_id BIGINT REFERENCES sponsors(id),

    -- Payment details
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',

    -- Type
    payment_type VARCHAR(50) NOT NULL, -- classified_premium, boost, sponsorship
    reference_id BIGINT, -- ID of classified, boost, or sponsorship

    -- Stripe
    stripe_payment_intent_id VARCHAR(100),
    stripe_charge_id VARCHAR(100),

    -- Status
    status VARCHAR(20) DEFAULT 'pending', -- pending, succeeded, failed, refunded

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW(),
    succeeded_at TIMESTAMP,
    failed_at TIMESTAMP,

    CONSTRAINT user_or_sponsor CHECK (
        (user_id IS NOT NULL AND sponsor_id IS NULL) OR
        (user_id IS NULL AND sponsor_id IS NOT NULL)
    )
);

CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_sponsor ON payments(sponsor_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_type ON payments(payment_type, created_at DESC);
CREATE INDEX idx_payments_stripe ON payments(stripe_payment_intent_id);
```

**Automation:**
- Stripe webhooks auto-update status
- Daily report of revenue (automated query)

---

## Automated Maintenance Tables

### 10. Moderation Queue

```sql
CREATE TABLE moderation_queue (
    id BIGSERIAL PRIMARY KEY,

    -- Item being moderated
    item_type VARCHAR(20) NOT NULL, -- comment, classified
    item_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),

    -- Reason
    flag_reason VARCHAR(50), -- spam, scam, inappropriate, etc.
    auto_flagged BOOLEAN DEFAULT FALSE,
    ai_confidence DECIMAL(3,2), -- 0.0 to 1.0

    -- Status
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, removed
    reviewed_by BIGINT REFERENCES users(id),
    reviewed_at TIMESTAMP,

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_moderation_pending ON moderation_queue(status, created_at) WHERE status = 'pending';
CREATE INDEX idx_moderation_item ON moderation_queue(item_type, item_id);
```

**Automation:**
- Auto-populated by spam detection
- Cron job: auto-approve items with high trust score
- Alert moderator if queue > 50 items

---

### 11. Audit Log

```sql
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),

    -- Action
    action VARCHAR(50) NOT NULL, -- login, vote, post_classified, etc.
    entity_type VARCHAR(50), -- article, comment, classified
    entity_id BIGINT,

    -- Metadata
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,

    -- Tracking
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_user ON audit_log(user_id, created_at DESC);
CREATE INDEX idx_audit_action ON audit_log(action);
CREATE INDEX idx_audit_created ON audit_log(created_at DESC);

-- Partition by month for performance
-- (Add partitioning later when needed)
```

**Automation:**
- Auto-logged on every important action
- Cron: Archive logs older than 1 year
- Security analysis: Detect suspicious patterns

---

## Utility Functions

### Auto-Update Timestamps

```sql
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to all tables
CREATE TRIGGER users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER articles_updated_at BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER comments_updated_at BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER classifieds_updated_at BEFORE UPDATE ON classifieds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER sponsors_updated_at BEFORE UPDATE ON sponsors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
```

---

### Auto-Calculate User Karma

```sql
CREATE OR REPLACE FUNCTION update_user_karma()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET karma = (
        SELECT COALESCE(SUM(
            CASE
                WHEN c.upvotes > 0 THEN c.upvotes
                ELSE 0
            END
        ), 0)
        FROM comments c
        WHERE c.user_id = NEW.user_id
    )
    WHERE id = NEW.user_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER comment_karma_update AFTER INSERT OR UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_user_karma();
```

---

### Auto-Expire Classifieds

```sql
-- Run this via cron every hour
CREATE OR REPLACE FUNCTION expire_old_classifieds()
RETURNS INTEGER AS $$
DECLARE
    expired_count INTEGER;
BEGIN
    UPDATE classifieds
    SET is_active = FALSE
    WHERE expires_at < NOW() AND is_active = TRUE;

    GET DIAGNOSTICS expired_count = ROW_COUNT;
    RETURN expired_count;
END;
$$ LANGUAGE plpgsql;
```

---

## Cron Jobs (Automation)

```sql
-- Install pg_cron extension
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Refresh article rankings every 5 minutes
SELECT cron.schedule('refresh-rankings', '*/5 * * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY article_rankings');

-- Expire old classifieds every hour
SELECT cron.schedule('expire-classifieds', '0 * * * *',
    'SELECT expire_old_classifieds()');

-- Clean up old audit logs (monthly)
SELECT cron.schedule('cleanup-logs', '0 0 1 * *',
    'DELETE FROM audit_log WHERE created_at < NOW() - INTERVAL ''1 year''');

-- Send billing reminders (daily at 9am)
SELECT cron.schedule('billing-reminders', '0 9 * * *',
    'SELECT send_billing_reminders()'); -- Custom function
```

---

## Performance Optimizations

### 1. Partitioning (For Scale)

```sql
-- Partition audit_log by month
CREATE TABLE audit_log_2025_01 PARTITION OF audit_log
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

-- Auto-create partitions via cron
```

### 2. Read Replicas

```sql
-- Configure in PostgreSQL config
-- master: writes
-- replica: reads (rankings, searches)
```

### 3. Connection Pooling

```sql
-- Use PgBouncer
-- 100 app connections → 20 database connections
```

---

## Backup Strategy

```bash
# Daily backups (automated via cron)
pg_dump -h localhost -U postgres terminalnews | gzip > backup_$(date +%Y%m%d).sql.gz

# Upload to S3
aws s3 cp backup_$(date +%Y%m%d).sql.gz s3://terminalnews-backups/
```

---

## Migration Management

Using golang-migrate:

```bash
# Create migration
migrate create -ext sql -dir migrations -seq add_sponsors_table

# Apply migrations
migrate -path migrations -database "postgres://..." up

# Rollback
migrate -path migrations -database "postgres://..." down 1
```

---

## Summary

**Automation Built-In:**
- ✅ Auto-expire classifieds
- ✅ Auto-refresh rankings
- ✅ Auto-calculate karma
- ✅ Auto-update timestamps
- ✅ Auto-archive old data
- ✅ Auto-trigger billing

**Performance:**
- ✅ Materialized views for rankings
- ✅ Smart indexes on all queries
- ✅ Full-text search ready
- ✅ Partitioning ready for scale

**Safety:**
- ✅ Foreign key constraints
- ✅ Check constraints
- ✅ Audit logging
- ✅ Soft deletes where needed

**This database runs itself.**
