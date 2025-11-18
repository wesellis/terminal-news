package handlers

import "net/http"

func (h *Handler) HandleGetWeather(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Weather API not yet implemented")
}
