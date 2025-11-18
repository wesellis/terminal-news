package handlers

import "net/http"

func (h *Handler) HandleGetComments(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandlePostComment(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleUpdateComment(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}
