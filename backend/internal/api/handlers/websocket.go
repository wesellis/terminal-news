package handlers

import "net/http"

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "WebSocket not yet implemented")
}

func (h *Handler) HandleGetActivity(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}
