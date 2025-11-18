package handlers

import "net/http"

func (h *Handler) HandleVote(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleDeleteVote(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleVoteComment(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}
