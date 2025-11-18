package handlers

import "net/http"

func (h *Handler) HandleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleGetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, r, http.StatusNotImplemented, "Not yet implemented")
}

func (h *Handler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Stripe webhook handling
	w.WriteHeader(http.StatusOK)
}
