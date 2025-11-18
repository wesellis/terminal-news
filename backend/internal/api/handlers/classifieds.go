package handlers

import "net/http"

func (h *Handler) HandleGetClassifieds(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleGetClassified(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandlePostClassified(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleUpdateClassified(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleDeleteClassified(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleBoostClassified(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}
