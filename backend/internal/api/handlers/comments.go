package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wesellis/terminal-news/backend/internal/services"
)

// HandleGetComments retrieves all comments for an article in a tree structure
func (h *Handler) HandleGetComments(w http.ResponseWriter, r *http.Request) {
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

	// Get comments
	response, err := h.commentService.GetArticleComments(r.Context(), articleID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve comments")
		return
	}

	h.respondJSON(w, r, http.StatusOK, response)
}

// HandlePostComment creates a new comment on an article
func (h *Handler) HandlePostComment(w http.ResponseWriter, r *http.Request) {
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
	var req services.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create comment
	comment, err := h.commentService.CreateComment(r.Context(), userID, articleID, req.Content, req.ParentID)
	if err != nil {
		switch err {
		case services.ErrCommentTooLong:
			h.respondError(w, r, http.StatusBadRequest, "Comment must be between 1 and 10,000 characters")
		case services.ErrInvalidParentID:
			h.respondError(w, r, http.StatusBadRequest, "Invalid parent comment ID")
		default:
			h.respondError(w, r, http.StatusInternalServerError, "Failed to create comment")
		}
		return
	}

	h.respondJSON(w, r, http.StatusCreated, comment)
}

// HandleUpdateComment updates an existing comment
func (h *Handler) HandleUpdateComment(w http.ResponseWriter, r *http.Request) {
	// Get comment ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Comment ID is required")
		return
	}

	// Parse comment ID
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Parse request body
	var req services.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update comment
	comment, err := h.commentService.UpdateComment(r.Context(), commentID, userID, req.Content)
	if err != nil {
		switch err {
		case services.ErrCommentNotFound:
			h.respondError(w, r, http.StatusNotFound, "Comment not found")
		case services.ErrCommentDeleted:
			h.respondError(w, r, http.StatusGone, "Comment has been deleted")
		case services.ErrUnauthorized:
			h.respondError(w, r, http.StatusForbidden, "You can only edit your own comments")
		case services.ErrCommentTooLong:
			h.respondError(w, r, http.StatusBadRequest, "Comment must be between 1 and 10,000 characters")
		default:
			h.respondError(w, r, http.StatusInternalServerError, "Failed to update comment")
		}
		return
	}

	h.respondJSON(w, r, http.StatusOK, comment)
}

// HandleDeleteComment soft-deletes a comment
func (h *Handler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	// Get comment ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Comment ID is required")
		return
	}

	// Parse comment ID
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Delete comment
	if err := h.commentService.DeleteComment(r.Context(), commentID, userID); err != nil {
		switch err {
		case services.ErrCommentNotFound:
			h.respondError(w, r, http.StatusNotFound, "Comment not found")
		case services.ErrUnauthorized:
			h.respondError(w, r, http.StatusForbidden, "You can only delete your own comments")
		default:
			h.respondError(w, r, http.StatusInternalServerError, "Failed to delete comment")
		}
		return
	}

	h.respondSuccess(w, r, http.StatusOK, nil, "Comment deleted successfully")
}
