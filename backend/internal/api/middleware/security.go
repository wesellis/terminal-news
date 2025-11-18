package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// SecurityHeaders adds security-related HTTP headers to responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Enable XSS protection (legacy browsers)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy - allow only same origin
		// Note: This is strict. Adjust if frontend needs external resources
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Prevent browser feature access
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS - Force HTTPS (only set if running on HTTPS)
		// Commented out for development, uncomment in production with HTTPS
		// w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}

// CORS handles Cross-Origin Resource Sharing
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin || allowedOrigin == "*" {
					allowed = true
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			// If no specific origin matched but * is in list
			if !allowed {
				for _, allowedOrigin := range allowedOrigins {
					if allowedOrigin == "*" {
						w.Header().Set("Access-Control-Allow-Origin", "*")
						allowed = true
						break
					}
				}
			}

			// Set other CORS headers
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight OPTIONS request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID adds a unique request ID to each request for tracing
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request already has an ID (from load balancer, etc.)
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate a simple request ID (in production, use UUID)
			requestID = generateRequestID()
		}

		// Add to response header for debugging
		w.Header().Set("X-Request-ID", requestID)

		// Could add to context here for logging throughout the request lifecycle
		// ctx := context.WithValue(r.Context(), "request_id", requestID)
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// generateRequestID creates a simple request ID
// In production, use github.com/google/uuid or similar
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// RecoverPanic recovers from panics and returns a 500 error
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic (in production, use proper logging)
				fmt.Printf("PANIC: %v\n", err)

				// Return 500 error
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
