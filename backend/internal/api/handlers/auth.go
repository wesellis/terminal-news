package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wesellis/terminal-news/backend/internal/services"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	User   *services.User      `json:"user"`
	Tokens *services.TokenPair `json:"tokens"`
}

// HandleRegister handles user registration
func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		h.respondError(w, r, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// Create user
	user, err := h.authService.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrUsernameTaken:
			h.respondError(w, r, http.StatusConflict, "Username already taken")
		case services.ErrEmailTaken:
			h.respondError(w, r, http.StatusConflict, "Email already taken")
		default:
			h.respondError(w, r, http.StatusInternalServerError, "Failed to create user")
		}
		return
	}

	// Generate tokens
	tokens, err := h.authService.GenerateTokens(user)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	h.respondJSON(w, r, http.StatusCreated, AuthResponse{
		User:   user,
		Tokens: tokens,
	})
}

// HandleLogin handles user login
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		h.respondError(w, r, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Authenticate user
	user, tokens, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			h.respondError(w, r, http.StatusUnauthorized, "Invalid username or password")
		case services.ErrInvalidPassword:
			h.respondError(w, r, http.StatusUnauthorized, "Invalid username or password")
		case services.ErrUserBanned:
			h.respondError(w, r, http.StatusForbidden, "Your account has been banned")
		default:
			h.respondError(w, r, http.StatusInternalServerError, "Login failed")
		}
		return
	}

	h.respondJSON(w, r, http.StatusOK, AuthResponse{
		User:   user,
		Tokens: tokens,
	})
}

// HandleRefreshToken handles token refresh
func (h *Handler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		h.respondError(w, r, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Refresh token
	tokens, err := h.authService.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case services.ErrInvalidToken:
			h.respondError(w, r, http.StatusUnauthorized, "Invalid refresh token")
		case services.ErrUserNotFound:
			h.respondError(w, r, http.StatusUnauthorized, "User not found")
		case services.ErrUserBanned:
			h.respondError(w, r, http.StatusForbidden, "Your account has been banned")
		default:
			h.respondError(w, r, http.StatusInternalServerError, "Failed to refresh token")
		}
		return
	}

	h.respondJSON(w, r, http.StatusOK, map[string]interface{}{
		"tokens": tokens,
	})
}

// HandleGetProfile returns the current user's profile
func (h *Handler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	h.respondJSON(w, r, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

// HandleUpdateProfile updates the current user's profile
func (h *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement profile update logic
	h.respondJSON(w, r, http.StatusOK, map[string]interface{}{
		"message": "Profile update not yet implemented",
	})
}
