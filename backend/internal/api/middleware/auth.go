package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/wesellis/terminal-news/backend/internal/services"
)

// AuthRequired is middleware that requires authentication
func AuthRequired(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Extract token (format: "Bearer <token>")
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			ctx = context.WithValue(ctx, "email", claims.Email)

			// Call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth is middleware that parses auth if present but doesn't require it
func OptionalAuth(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Extract token
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token := parts[1]

					// Validate token
					claims, err := authService.ValidateToken(token)
					if err == nil {
						// Add user info to context
						ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
						ctx = context.WithValue(ctx, "username", claims.Username)
						ctx = context.WithValue(ctx, "email", claims.Email)
						r = r.WithContext(ctx)
					}
				}
			}

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}
