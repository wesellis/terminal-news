package handlers

import (
	"net/http"
	"strconv"

	"github.com/wesellis/terminal-news/backend/internal/external"
)

// HandleGetWeather fetches weather data from NOAA
func (h *Handler) HandleGetWeather(w http.ResponseWriter, r *http.Request) {
	// Get coordinates from query parameters
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		h.respondError(w, r, http.StatusBadRequest, "Both lat and lng parameters are required")
		return
	}

	// Parse coordinates
	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)

	if err1 != nil || err2 != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid lat/lng coordinates")
		return
	}

	// Validate coordinate ranges
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		h.respondError(w, r, http.StatusBadRequest, "Coordinates out of valid range")
		return
	}

	// Create weather service
	weatherService := external.NewWeatherService()

	// Fetch weather
	weather, err := weatherService.GetWeather(r.Context(), lat, lng)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to fetch weather data: "+err.Error())
		return
	}

	h.respondJSON(w, r, http.StatusOK, weather)
}
