package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wesellis/terminal-news/cli/internal/models"
)

// WebSocketClient manages the WebSocket connection for real-time updates
type WebSocketClient struct {
	conn   *websocket.Conn
	url    string
	events chan models.WSMessage
	done   chan struct{}
}

// NewWebSocketClient creates and connects a new WebSocket client
func NewWebSocketClient(url string) (*WebSocketClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	ws := &WebSocketClient{
		conn:   conn,
		url:    url,
		events: make(chan models.WSMessage, 100),
		done:   make(chan struct{}),
	}

	// Start listening for messages
	go ws.listen()

	// Start ping/pong to keep connection alive
	go ws.keepAlive()

	return ws, nil
}

// listen continuously reads messages from the WebSocket
func (ws *WebSocketClient) listen() {
	defer func() {
		ws.conn.Close()
		close(ws.events)
		close(ws.done)
	}()

	for {
		select {
		case <-ws.done:
			return
		default:
			var msg models.WSMessage
			err := ws.conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				// Attempt to reconnect
				if err := ws.reconnect(); err != nil {
					log.Printf("Failed to reconnect: %v", err)
					return
				}
				continue
			}

			ws.events <- msg
		}
	}
}

// keepAlive sends periodic ping messages to keep the connection alive
func (ws *WebSocketClient) keepAlive() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ws.done:
			return
		case <-ticker.C:
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("WebSocket ping error: %v", err)
				return
			}
		}
	}
}

// reconnect attempts to reconnect the WebSocket
func (ws *WebSocketClient) reconnect() error {
	log.Println("Attempting to reconnect WebSocket...")

	// Close existing connection
	if ws.conn != nil {
		ws.conn.Close()
	}

	// Try to reconnect with exponential backoff
	maxRetries := 5
	baseDelay := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(ws.url, nil)
		if err == nil {
			ws.conn = conn
			log.Println("WebSocket reconnected successfully")
			return nil
		}

		delay := baseDelay * time.Duration(1<<uint(i))
		log.Printf("Reconnect attempt %d/%d failed, retrying in %v", i+1, maxRetries, delay)
		time.Sleep(delay)
	}

	return fmt.Errorf("failed to reconnect after %d attempts", maxRetries)
}

// Events returns the channel of WebSocket events
func (ws *WebSocketClient) Events() <-chan models.WSMessage {
	return ws.events
}

// Send sends a message through the WebSocket
func (ws *WebSocketClient) Send(msgType string, data interface{}) error {
	msg := models.WSMessage{
		Type: msgType,
		Data: data,
	}

	return ws.conn.WriteJSON(msg)
}

// Close closes the WebSocket connection
func (ws *WebSocketClient) Close() error {
	close(ws.done)
	return ws.conn.Close()
}

// Subscribe subscribes to specific event types
func (ws *WebSocketClient) Subscribe(eventTypes ...string) error {
	return ws.Send("subscribe", map[string]interface{}{
		"events": eventTypes,
	})
}

// Unsubscribe unsubscribes from specific event types
func (ws *WebSocketClient) Unsubscribe(eventTypes ...string) error {
	return ws.Send("unsubscribe", map[string]interface{}{
		"events": eventTypes,
	})
}
