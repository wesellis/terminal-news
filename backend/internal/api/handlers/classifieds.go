package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wesellis/terminal-news/backend/internal/services"
)

// HandleGetClassifieds retrieves classifieds with optional filtering
func (h *Handler) HandleGetClassifieds(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	category := r.URL.Query().Get("category")
	city := r.URL.Query().Get("city")
	state := r.URL.Query().Get("state")
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")
	radiusStr := r.URL.Query().Get("radius")

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

	// Check if geographic search
	if latStr != "" && lngStr != "" && radiusStr != "" {
		lat, err1 := strconv.ParseFloat(latStr, 64)
		lng, err2 := strconv.ParseFloat(lngStr, 64)
		radius, err3 := strconv.ParseFloat(radiusStr, 64)

		if err1 == nil && err2 == nil && err3 == nil {
			// Geographic search
			response, err := h.classifiedService.SearchClassifiedsByLocation(r.Context(), lat, lng, radius, limit, offset)
			if err != nil {
				h.respondError(w, r, http.StatusInternalServerError, "Failed to search classifieds by location")
				return
			}
			h.respondJSON(w, r, http.StatusOK, response)
			return
		}
	}

	// Standard search with filters
	response, err := h.classifiedService.GetClassifieds(r.Context(), category, city, state, limit, offset)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve classifieds")
		return
	}

	h.respondJSON(w, r, http.StatusOK, response)
}

// HandleGetClassified retrieves a single classified by ID
func (h *Handler) HandleGetClassified(w http.ResponseWriter, r *http.Request) {
	// Get classified ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Classified ID is required")
		return
	}

	// Parse classified ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid classified ID")
		return
	}

	// Get classified
	classified, err := h.classifiedService.GetClassified(r.Context(), id)
	if err != nil {
		if err == services.ErrClassifiedNotFound {
			h.respondError(w, r, http.StatusNotFound, "Classified not found")
			return
		}
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve classified")
		return
	}

	h.respondJSON(w, r, http.StatusOK, classified)
}

// HandlePostClassified creates a new classified ad
func (h *Handler) HandlePostClassified(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Parse request body
	var req services.CreateClassifiedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create classified
	classified, err := h.classifiedService.CreateClassified(r.Context(), userID, &req)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, r, http.StatusCreated, classified)
}

// HandleUpdateClassified updates an existing classified
func (h *Handler) HandleUpdateClassified(w http.ResponseWriter, r *http.Request) {
	// Get classified ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Classified ID is required")
		return
	}

	// Parse classified ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid classified ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Parse request body
	var req services.UpdateClassifiedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update classified
	classified, err := h.classifiedService.UpdateClassified(r.Context(), id, userID, &req)
	if err != nil {
		switch err {
		case services.ErrClassifiedNotFound:
			h.respondError(w, r, http.StatusNotFound, "Classified not found")
		case services.ErrClassifiedInactive:
			h.respondError(w, r, http.StatusGone, "Classified is inactive")
		case services.ErrUnauthorized:
			h.respondError(w, r, http.StatusForbidden, "You can only edit your own classifieds")
		default:
			h.respondError(w, r, http.StatusInternalServerError, "Failed to update classified")
		}
		return
	}

	h.respondJSON(w, r, http.StatusOK, classified)
}

// HandleDeleteClassified soft-deletes a classified
func (h *Handler) HandleDeleteClassified(w http.ResponseWriter, r *http.Request) {
	// Get classified ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Classified ID is required")
		return
	}

	// Parse classified ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid classified ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Delete classified
	if err := h.classifiedService.DeleteClassified(r.Context(), id, userID); err != nil {
		switch err {
		case services.ErrClassifiedNotFound:
			h.respondError(w, r, http.StatusNotFound, "Classified not found")
		case services.ErrUnauthorized:
			h.respondError(w, r, http.StatusForbidden, "You can only delete your own classifieds")
		default:
			h.respondError(w, r, http.StatusInternalServerError, "Failed to delete classified")
		}
		return
	}

	h.respondSuccess(w, r, http.StatusOK, nil, "Classified deleted successfully")
}

// HandleBoostClassified increments contact count (stub for future payment integration)
func (h *Handler) HandleBoostClassified(w http.ResponseWriter, r *http.Request) {
	// Get classified ID from URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Classified ID is required")
		return
	}

	// Parse classified ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid classified ID")
		return
	}

	// Increment contact count
	if err := h.classifiedService.IncrementContactCount(r.Context(), id); err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to track contact")
		return
	}

	h.respondSuccess(w, r, http.StatusOK, nil, "Contact tracked successfully")
}
