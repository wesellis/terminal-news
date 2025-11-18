package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wesellis/terminal-news/backend/internal/services"
)

type Handler struct {
	authService       *services.AuthService
	articleService    *services.ArticleService
	voteService       *services.VoteService
	commentService    *services.CommentService
	classifiedService *services.ClassifiedService
	paymentService    *services.PaymentService
}

func NewHandler(
	authService *services.AuthService,
	articleService *services.ArticleService,
	voteService *services.VoteService,
	commentService *services.CommentService,
	classifiedService *services.ClassifiedService,
	paymentService *services.PaymentService,
) *Handler {
	return &Handler{
		authService:       authService,
		articleService:    articleService,
		voteService:       voteService,
		commentService:    commentService,
		classifiedService: classifiedService,
		paymentService:    paymentService,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// Helper functions

func (h *Handler) respondJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

func (h *Handler) respondSuccess(w http.ResponseWriter, r *http.Request, status int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(SuccessResponse{
		Data:    data,
		Message: message,
	})
}

// getUserID extracts the user ID from the request context (set by auth middleware)
func (h *Handler) getUserID(r *http.Request) int64 {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		return 0
	}
	return userID
}

// getUsername extracts the username from the request context
func (h *Handler) getUsername(r *http.Request) string {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		return ""
	}
	return username
}

// withValue is a helper to add values to request context
func withValue(r *http.Request, key string, value interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}
