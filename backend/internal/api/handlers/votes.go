package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wesellis/terminal-news/backend/internal/services"
)

// HandleVote creates a vote on an article
func (h *Handler) HandleVote(w http.ResponseWriter, r *http.Request) {
	// Get article ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Article ID is required")
		return
	}

	// Parse article ID
	articleID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid article ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Parse request body
	var req services.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create vote
	response, err := h.voteService.CreateVote(r.Context(), userID, articleID, req.VoteType)
	if err != nil {
		if err == services.ErrInvalidVoteType {
			h.respondError(w, r, http.StatusBadRequest, "Invalid vote type. Must be 'open', 'like', or 'dislike'")
			return
		}
		h.respondError(w, r, http.StatusInternalServerError, "Failed to create vote")
		return
	}

	h.respondJSON(w, r, http.StatusCreated, response)
}

// HandleDeleteVote removes a vote from an article
func (h *Handler) HandleDeleteVote(w http.ResponseWriter, r *http.Request) {
	// Get article ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Article ID is required")
		return
	}

	// Parse article ID
	articleID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid article ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Get vote type from query parameter
	voteType := r.URL.Query().Get("vote_type")
	if voteType == "" {
		h.respondError(w, r, http.StatusBadRequest, "Vote type is required")
		return
	}

	// Remove vote
	if err := h.voteService.RemoveVote(r.Context(), userID, articleID, voteType); err != nil {
		if err == services.ErrInvalidVoteType {
			h.respondError(w, r, http.StatusBadRequest, "Invalid vote type. Must be 'open', 'like', or 'dislike'")
			return
		}
		if err == services.ErrVoteNotFound {
			h.respondError(w, r, http.StatusNotFound, "Vote not found")
			return
		}
		h.respondError(w, r, http.StatusInternalServerError, "Failed to remove vote")
		return
	}

	h.respondSuccess(w, r, http.StatusOK, nil, "Vote removed successfully")
}

// HandleVoteComment creates a vote on a comment (stub for now)
func (h *Handler) HandleVoteComment(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Comment voting not yet implemented")
}
