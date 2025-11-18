package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter provides rate limiting middleware using Redis
type RateLimiter struct {
	redisClient *redis.Client
	// Requests per window
	limit int
	// Time window duration
	window time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		limit:       limit,
		window:      window,
	}
}

// Limit returns middleware that limits requests per IP address
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		ip := getClientIP(r)

		// Create Redis key for this IP
		key := fmt.Sprintf("ratelimit:%s", ip)

		ctx := context.Background()

		// Increment counter
		count, err := rl.redisClient.Incr(ctx, key).Result()
		if err != nil {
			// If Redis fails, allow the request (fail open)
			next.ServeHTTP(w, r)
			return
		}

		// Set expiration on first request
		if count == 1 {
			rl.redisClient.Expire(ctx, key, rl.window)
		}

		// Get TTL for rate limit reset time
		ttl, err := rl.redisClient.TTL(ctx, key).Result()
		if err != nil {
			ttl = rl.window
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(max(0, rl.limit-int(count))))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

		// Check if over limit
		if count > int64(rl.limit) {
			w.Header().Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LimitByUser returns middleware that limits requests per authenticated user
func (rl *RateLimiter) LimitByUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("user_id").(int64)
		if !ok || userID == 0 {
			// No user ID, fall back to IP-based limiting
			rl.Limit(next).ServeHTTP(w, r)
			return
		}

		// Create Redis key for this user
		key := fmt.Sprintf("ratelimit:user:%d", userID)

		ctx := context.Background()

		// Increment counter
		count, err := rl.redisClient.Incr(ctx, key).Result()
		if err != nil {
			// If Redis fails, allow the request (fail open)
			next.ServeHTTP(w, r)
			return
		}

		// Set expiration on first request
		if count == 1 {
			rl.redisClient.Expire(ctx, key, rl.window)
		}

		// Get TTL for rate limit reset time
		ttl, err := rl.redisClient.TTL(ctx, key).Result()
		if err != nil {
			ttl = rl.window
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(max(0, rl.limit-int(count))))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

		// Check if over limit
		if count > int64(rl.limit) {
			w.Header().Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if multiple are listed
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
