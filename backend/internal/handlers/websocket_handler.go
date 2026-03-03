package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	ws "github.com/aung-arata/youtube-clone/backend/internal/websocket"
)

// WebSocketHandler handles WebSocket upgrade requests
type WebSocketHandler struct {
	hub *ws.Hub
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleConnection upgrades an HTTP connection to WebSocket
func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Get user ID from query parameter
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		ClientID: fmt.Sprintf("client-%d-%d", userID, time.Now().UnixNano()),
	}

	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
