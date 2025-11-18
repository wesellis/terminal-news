-- Migration: 001_initial_schema.sql
-- Description: Initial database schema for Terminal News
-- Created: 2024-11-18

BEGIN;

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For fuzzy text search

-- ============================================================================
-- USERS TABLE
-- ============================================================================
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
    trust_score DECIMAL(3,2) DEFAULT 0.5 CHECK (trust_score >= 0 AND trust_score <= 1),

    -- Status
    email_verified BOOLEAN DEFAULT FALSE,
    is_banned BOOLEAN DEFAULT FALSE,
    is_moderator BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,

    -- Tracking
    last_active_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    -- Constraints
    CONSTRAINT username_length CHECK (LENGTH(username) >= 3),
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_karma ON users(karma DESC);
CREATE INDEX idx_users_last_active ON users(last_active_at DESC);

-- ============================================================================
-- ARTICLES TABLE
-- ============================================================================
CREATE TABLE articles (
    id BIGSERIAL PRIMARY KEY,

    -- Content
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    content TEXT,
    image_url TEXT,

    -- Metadata
    source VARCHAR(100),
    author VARCHAR(255),
    published_at TIMESTAMP,

    -- Categorization
    category VARCHAR(50),
    tags TEXT[],

    -- Aggregation tracking
    external_id VARCHAR(255),
    fetch_source VARCHAR(50),

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
CREATE INDEX idx_articles_title_search ON articles USING GIN(to_tsvector('english', title));

-- ============================================================================
-- VOTES TABLE
-- ============================================================================
CREATE TABLE votes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('open', 'like', 'dislike')),
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT unique_user_article_vote UNIQUE (user_id, article_id, vote_type)
);

CREATE INDEX idx_votes_article ON votes(article_id);
CREATE INDEX idx_votes_user ON votes(user_id);
CREATE INDEX idx_votes_created ON votes(created_at DESC);
CREATE INDEX idx_votes_type ON votes(vote_type);
CREATE INDEX idx_votes_article_type_created ON votes(article_id, vote_type, created_at);

-- ============================================================================
-- COMMENTS TABLE
-- ============================================================================
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,

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

-- ============================================================================
-- CLASSIFIEDS TABLE
-- ============================================================================
CREATE TABLE classifieds (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Content
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10,2),

    -- Categorization
    category VARCHAR(50) NOT NULL,
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
    contact_method VARCHAR(20) DEFAULT 'email',

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
CREATE INDEX idx_classifieds_search ON classifieds USING GIN(to_tsvector('english', title || ' ' || description));

-- ============================================================================
-- ARTICLE RANKINGS (Materialized View)
-- ============================================================================
CREATE MATERIALIZED VIEW article_rankings AS
SELECT
    a.id AS article_id,

    -- Vote counts
    COUNT(*) FILTER (WHERE v.vote_type = 'open') AS open_count,
    COUNT(*) FILTER (WHERE v.vote_type = 'like') AS like_count,
    COUNT(*) FILTER (WHERE v.vote_type = 'dislike') AS dislike_count,

    -- Total score
    (
        COUNT(*) FILTER (WHERE v.vote_type = 'open') * 1 +
        COUNT(*) FILTER (WHERE v.vote_type = 'like') * 2 +
        COUNT(*) FILTER (WHERE v.vote_type = 'dislike') * -1
    ) AS total_score,

    -- Controversy score
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

    -- Total engagement
    COUNT(*) AS total_engagement,

    -- Time metrics
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
WHERE a.published_at > NOW() - INTERVAL '7 days'
GROUP BY a.id, a.published_at;

CREATE UNIQUE INDEX idx_article_rankings_id ON article_rankings(article_id);
CREATE INDEX idx_article_rankings_hot ON article_rankings(hot_rank DESC);
CREATE INDEX idx_article_rankings_controversy ON article_rankings(controversy_score DESC);
CREATE INDEX idx_article_rankings_score ON article_rankings(total_score DESC);

-- ============================================================================
-- SPONSORS TABLE
-- ============================================================================
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

-- ============================================================================
-- SPONSORSHIPS TABLE
-- ============================================================================
CREATE TABLE sponsorships (
    id BIGSERIAL PRIMARY KEY,
    sponsor_id BIGINT NOT NULL REFERENCES sponsors(id) ON DELETE CASCADE,

    -- Placement
    placement_type VARCHAR(50) NOT NULL,
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

-- ============================================================================
-- PAYMENTS TABLE
-- ============================================================================
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    sponsor_id BIGINT REFERENCES sponsors(id),

    -- Payment details
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',

    -- Type
    payment_type VARCHAR(50) NOT NULL,
    reference_id BIGINT,

    -- Stripe
    stripe_payment_intent_id VARCHAR(100),
    stripe_charge_id VARCHAR(100),

    -- Status
    status VARCHAR(20) DEFAULT 'pending',

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

-- ============================================================================
-- AUDIT LOG TABLE
-- ============================================================================
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),

    -- Action
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50),
    entity_id BIGINT,

    -- Metadata
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_user ON audit_log(user_id, created_at DESC);
CREATE INDEX idx_audit_action ON audit_log(action);
CREATE INDEX idx_audit_created ON audit_log(created_at DESC);

-- ============================================================================
-- MODERATION QUEUE TABLE
-- ============================================================================
CREATE TABLE moderation_queue (
    id BIGSERIAL PRIMARY KEY,

    -- Item being moderated
    item_type VARCHAR(20) NOT NULL,
    item_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),

    -- Reason
    flag_reason VARCHAR(50),
    auto_flagged BOOLEAN DEFAULT FALSE,
    ai_confidence DECIMAL(3,2),

    -- Status
    status VARCHAR(20) DEFAULT 'pending',
    reviewed_by BIGINT REFERENCES users(id),
    reviewed_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_moderation_pending ON moderation_queue(status, created_at) WHERE status = 'pending';
CREATE INDEX idx_moderation_item ON moderation_queue(item_type, item_id);

COMMIT;
