package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wesellis/terminal-news/backend/internal/services"
)

// HandleGetArticles retrieves articles with optional feed type
func (h *Handler) HandleGetArticles(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	feed := r.URL.Query().Get("feed")

	// Default values
	limit := 25
	offset := 0

	// Parse limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse offset
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get articles
	response, err := h.articleService.GetArticles(r.Context(), feed, limit, offset)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve articles")
		return
	}

	h.respondJSON(w, r, http.StatusOK, response)
}

// HandleGetHotArticles retrieves hot articles
func (h *Handler) HandleGetHotArticles(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default values
	limit := 25
	offset := 0

	// Parse limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse offset
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get hot articles
	response, err := h.articleService.GetHotArticles(r.Context(), limit, offset)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve hot articles")
		return
	}

	h.respondJSON(w, r, http.StatusOK, response)
}

// HandleGetControversialArticles retrieves controversial articles
func (h *Handler) HandleGetControversialArticles(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default values
	limit := 25
	offset := 0

	// Parse limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse offset
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get controversial articles
	response, err := h.articleService.GetControversialArticles(r.Context(), limit, offset)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve controversial articles")
		return
	}

	h.respondJSON(w, r, http.StatusOK, response)
}

// HandleGetRisingArticles retrieves rising articles
func (h *Handler) HandleGetRisingArticles(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default values
	limit := 25
	offset := 0

	// Parse limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse offset
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get rising articles
	response, err := h.articleService.GetRisingArticles(r.Context(), limit, offset)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve rising articles")
		return
	}

	h.respondJSON(w, r, http.StatusOK, response)
}

// HandleGetArticle retrieves a single article by ID
func (h *Handler) HandleGetArticle(w http.ResponseWriter, r *http.Request) {
	// Get article ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Article ID is required")
		return
	}

	// Parse article ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid article ID")
		return
	}

	// Get article
	article, err := h.articleService.GetArticle(r.Context(), id)
	if err != nil {
		if err == services.ErrArticleNotFound {
			h.respondError(w, r, http.StatusNotFound, "Article not found")
			return
		}
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve article")
		return
	}

	h.respondJSON(w, r, http.StatusOK, article)
}
