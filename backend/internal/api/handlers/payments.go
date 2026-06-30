package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/wesellis/terminal-news/backend/internal/services"
)

// CreatePaymentIntentRequest represents a payment intent creation request
type CreatePaymentIntentRequest struct {
	Type         string `json:"type"`          // "classified_boost" or "sponsor_subscription"
	ClassifiedID int64  `json:"classified_id,omitempty"`
	DurationDays int    `json:"duration_days,omitempty"` // for classified boost
	Tier         string `json:"tier,omitempty"`          // for sponsor subscription
}

// HandleCreatePaymentIntent creates a payment intent for classified boost or sponsor subscription
func (h *Handler) HandleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Parse request
	var req CreatePaymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	var payment *services.Payment
	var err error

	switch req.Type {
	case "classified_boost":
		if req.ClassifiedID == 0 {
			h.respondError(w, r, http.StatusBadRequest, "classified_id is required")
			return
		}
		if req.DurationDays == 0 {
			req.DurationDays = 7 // default to 7 days
		}

		// Note: Ownership verification done in service layer
		payment, err = h.paymentService.CreateClassifiedBoostPayment(r.Context(), userID, req.ClassifiedID, req.DurationDays)

	case "sponsor_subscription":
		if req.Tier == "" {
			h.respondError(w, r, http.StatusBadRequest, "tier is required (bronze, silver, gold)")
			return
		}

		payment, err = h.paymentService.CreateSponsorSubscription(r.Context(), userID, req.Tier)

	default:
		h.respondError(w, r, http.StatusBadRequest, "Invalid payment type. Use 'classified_boost' or 'sponsor_subscription'")
		return
	}

	if err != nil {
		if err == services.ErrInvalidAmount {
			h.respondError(w, r, http.StatusBadRequest, "Invalid payment amount or duration")
			return
		}
		h.respondError(w, r, http.StatusInternalServerError, "Failed to create payment: "+err.Error())
		return
	}

	h.respondJSON(w, r, http.StatusCreated, payment)
}

// HandleGetPaymentHistory retrieves payment history for the authenticated user
func (h *Handler) HandleGetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	payments, err := h.paymentService.GetPaymentHistory(r.Context(), userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve payment history")
		return
	}

	h.respondJSON(w, r, http.StatusOK, map[string]interface{}{
		"payments": payments,
		"count":    len(payments),
	})
}

// HandleStripeWebhook handles Stripe webhook events
func (h *Handler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, r, http.StatusServiceUnavailable, "Error reading request body")
		return
	}

	// Verify webhook signature
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Webhook signature verification failed")
		return
	}

	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent map[string]interface{}
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "Error parsing webhook JSON")
			return
		}

		paymentIntentID, ok := paymentIntent["id"].(string)
		if !ok {
			h.respondError(w, r, http.StatusBadRequest, "Invalid payment intent ID")
			return
		}

		// Handle successful payment
		err = h.paymentService.HandlePaymentSuccess(r.Context(), paymentIntentID)
		if err != nil {
			// Log error but return 200 to Stripe
			h.respondError(w, r, http.StatusInternalServerError, "Error processing payment success")
			return
		}

	case "customer.subscription.created":
		// Subscription created - already handled in CreateSponsorSubscription
		// Just acknowledge the webhook

	case "customer.subscription.deleted":
		var subscription map[string]interface{}
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "Error parsing webhook JSON")
			return
		}

		subscriptionID, ok := subscription["id"].(string)
		if !ok {
			h.respondError(w, r, http.StatusBadRequest, "Invalid subscription ID")
			return
		}

		// Handle subscription cancellation
		err = h.paymentService.HandleSubscriptionCanceled(r.Context(), subscriptionID)
		if err != nil {
			// Log error but return 200 to Stripe
			h.respondError(w, r, http.StatusInternalServerError, "Error processing subscription cancellation")
			return
		}

	case "invoice.payment_succeeded":
		// Recurring payment succeeded - update subscription status if needed
		// For now, just acknowledge

	case "invoice.payment_failed":
		// Payment failed - potentially notify user or downgrade subscription
		// For now, just acknowledge

	default:
		// Unexpected event type
		h.respondError(w, r, http.StatusBadRequest, "Unhandled event type: "+string(event.Type))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleSetupStripeProducts is an admin endpoint to setup Stripe products (one-time use)
func (h *Handler) HandleSetupStripeProducts(w http.ResponseWriter, r *http.Request) {
	// This should be protected by admin auth in production
	err := h.paymentService.SetupStripeProducts()
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to setup Stripe products: "+err.Error())
		return
	}

	h.respondSuccess(w, r, http.StatusOK, nil, "Stripe products and prices created successfully")
}
