package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	ws "github.com/wesellis/terminal-news/backend/pkg/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now (restrict in production)
		return true
	},
}

// HandleWebSocket handles WebSocket connection upgrades
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	// Get user ID from context (optional - can be 0 for unauthenticated users)
	userID := h.getUserID(r)

	// Create new client
	client := ws.NewClient(h.wsHub, conn, userID)

	// Register client with hub
	h.wsHub.Register(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()

	log.Printf("WebSocket connection established for user %d", userID)
}

// HandleGetActivity retrieves user activity history
func (h *Handler) HandleGetActivity(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := h.getUserID(r)
	if userID == 0 {
		h.respondError(w, r, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Get user's comments
	comments, err := h.commentService.GetUserComments(r.Context(), userID, 50, 0)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve user activity")
		return
	}

	// Get user's votes
	votes, err := h.voteService.GetUserVotes(r.Context(), userID, 50, 0)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Failed to retrieve user activity")
		return
	}

	activity := map[string]interface{}{
		"comments": comments,
		"votes":    votes,
	}

	h.respondJSON(w, r, http.StatusOK, activity)
}
