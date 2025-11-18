-- Migration: 002_triggers_and_functions.sql
-- Description: Automated triggers and functions for Terminal News
-- Created: 2024-11-18

BEGIN;

-- ============================================================================
-- AUTO-UPDATE TIMESTAMP FUNCTION
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to all tables with updated_at
CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER articles_updated_at
    BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER comments_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER classifieds_updated_at
    BEFORE UPDATE ON classifieds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER sponsors_updated_at
    BEFORE UPDATE ON sponsors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER sponsorships_updated_at
    BEFORE UPDATE ON sponsorships
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================================================
-- AUTO-UPDATE USER KARMA
-- ============================================================================
CREATE OR REPLACE FUNCTION update_user_karma()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET karma = (
        SELECT COALESCE(SUM(upvotes - downvotes), 0)
        FROM comments
        WHERE user_id = NEW.user_id
          AND is_deleted = FALSE
    )
    WHERE id = NEW.user_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER comment_karma_update
    AFTER INSERT OR UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_user_karma();

-- ============================================================================
-- AUTO-EXPIRE CLASSIFIEDS
-- ============================================================================
CREATE OR REPLACE FUNCTION expire_old_classifieds()
RETURNS INTEGER AS $$
DECLARE
    expired_count INTEGER;
BEGIN
    UPDATE classifieds
    SET is_active = FALSE
    WHERE expires_at < NOW()
      AND is_active = TRUE;

    GET DIAGNOSTICS expired_count = ROW_COUNT;
    RETURN expired_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- AUTO-DOWNGRADE EXPIRED PREMIUM
-- ============================================================================
CREATE OR REPLACE FUNCTION downgrade_expired_premium()
RETURNS INTEGER AS $$
DECLARE
    downgraded_count INTEGER;
BEGIN
    UPDATE classifieds
    SET is_premium = FALSE
    WHERE premium_until < NOW()
      AND is_premium = TRUE;

    GET DIAGNOSTICS downgraded_count = ROW_COUNT;
    RETURN downgraded_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- AUTO-FLAG SUSPICIOUS CONTENT
-- ============================================================================
CREATE OR REPLACE FUNCTION auto_flag_content()
RETURNS TRIGGER AS $$
DECLARE
    spam_score FLOAT;
BEGIN
    -- Simple spam detection (can be enhanced with AI)
    spam_score := 0;

    -- Check for excessive caps
    IF NEW.content ~ '[A-Z]{10,}' THEN
        spam_score := spam_score + 0.3;
    END IF;

    -- Check for excessive links
    IF (SELECT regexp_count(NEW.content, 'https?://')) > 3 THEN
        spam_score := spam_score + 0.4;
    END IF;

    -- Check for repeated characters
    IF NEW.content ~ '(.)\1{5,}' THEN
        spam_score := spam_score + 0.2;
    END IF;

    -- Auto-flag if score is high
    IF spam_score > 0.7 THEN
        NEW.is_flagged := TRUE;
        NEW.flag_count := 1;

        -- Add to moderation queue
        INSERT INTO moderation_queue (
            item_type,
            item_id,
            user_id,
            flag_reason,
            auto_flagged,
            ai_confidence
        ) VALUES (
            TG_TABLE_NAME,
            NEW.id,
            NEW.user_id,
            'auto_spam_detection',
            TRUE,
            spam_score
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER comment_spam_check
    BEFORE INSERT ON comments
    FOR EACH ROW EXECUTE FUNCTION auto_flag_content();

-- ============================================================================
-- REFRESH ARTICLE RANKINGS (Call from cron)
-- ============================================================================
CREATE OR REPLACE FUNCTION refresh_article_rankings()
RETURNS VOID AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY article_rankings;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- CLEANUP OLD ARTICLES (Call from cron)
-- ============================================================================
CREATE OR REPLACE FUNCTION cleanup_old_articles()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Delete articles older than 90 days with no engagement
    DELETE FROM articles
    WHERE published_at < NOW() - INTERVAL '90 days'
      AND id NOT IN (
          SELECT DISTINCT article_id FROM votes
          UNION
          SELECT DISTINCT article_id FROM comments
      );

    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- GET HOT ARTICLES (Helper function)
-- ============================================================================
CREATE OR REPLACE FUNCTION get_hot_articles(
    limit_count INTEGER DEFAULT 50,
    offset_count INTEGER DEFAULT 0
)
RETURNS TABLE (
    article_id BIGINT,
    title TEXT,
    url TEXT,
    source VARCHAR(100),
    published_at TIMESTAMP,
    open_count INTEGER,
    like_count INTEGER,
    dislike_count INTEGER,
    total_score INTEGER,
    hot_rank FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        a.id,
        a.title,
        a.url,
        a.source,
        a.published_at,
        COALESCE(ar.open_count, 0)::INTEGER,
        COALESCE(ar.like_count, 0)::INTEGER,
        COALESCE(ar.dislike_count, 0)::INTEGER,
        COALESCE(ar.total_score, 0)::INTEGER,
        COALESCE(ar.hot_rank, 0)::FLOAT
    FROM articles a
    LEFT JOIN article_rankings ar ON a.id = ar.article_id
    WHERE a.published_at > NOW() - INTERVAL '7 days'
    ORDER BY ar.hot_rank DESC NULLS LAST
    LIMIT limit_count
    OFFSET offset_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- GET CONTROVERSIAL ARTICLES
-- ============================================================================
CREATE OR REPLACE FUNCTION get_controversial_articles(
    limit_count INTEGER DEFAULT 50,
    offset_count INTEGER DEFAULT 0
)
RETURNS TABLE (
    article_id BIGINT,
    title TEXT,
    url TEXT,
    source VARCHAR(100),
    published_at TIMESTAMP,
    like_count INTEGER,
    dislike_count INTEGER,
    controversy_score FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        a.id,
        a.title,
        a.url,
        a.source,
        a.published_at,
        COALESCE(ar.like_count, 0)::INTEGER,
        COALESCE(ar.dislike_count, 0)::INTEGER,
        COALESCE(ar.controversy_score, 0)::FLOAT
    FROM articles a
    INNER JOIN article_rankings ar ON a.id = ar.article_id
    WHERE a.published_at > NOW() - INTERVAL '7 days'
      AND ar.controversy_score > 0.5
    ORDER BY ar.controversy_score DESC, ar.total_engagement DESC
    LIMIT limit_count
    OFFSET offset_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- SEARCH ARTICLES (Full-text search)
-- ============================================================================
CREATE OR REPLACE FUNCTION search_articles(
    search_query TEXT,
    limit_count INTEGER DEFAULT 50
)
RETURNS TABLE (
    article_id BIGINT,
    title TEXT,
    url TEXT,
    source VARCHAR(100),
    published_at TIMESTAMP,
    relevance FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        a.id,
        a.title,
        a.url,
        a.source,
        a.published_at,
        ts_rank(to_tsvector('english', a.title), plainto_tsquery('english', search_query)) AS relevance
    FROM articles a
    WHERE to_tsvector('english', a.title) @@ plainto_tsquery('english', search_query)
    ORDER BY relevance DESC, a.published_at DESC
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;

COMMIT;
