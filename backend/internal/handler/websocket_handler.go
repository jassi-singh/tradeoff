package handler

import (
	"log"
	"net/http"
	"net/url"
	"tradeoff/backend/internal/config"
	"tradeoff/backend/internal/helpers"
	"tradeoff/backend/internal/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get token from query parameter
	encodedToken := r.URL.Query().Get("token")
	if encodedToken == "" {
		log.Println("WebSocket connection rejected: missing token")
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// URL decode the token
	token, err := url.QueryUnescape(encodedToken)
	if err != nil {
		log.Printf("WebSocket connection rejected: failed to decode token - %v", err)
		http.Error(w, "Invalid token encoding", http.StatusBadRequest)
		return
	}

	// Load config to get JWT secret
	conf, err := config.LoadConfig()
	if err != nil {
		log.Printf("WebSocket connection rejected: config error - %v", err)
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	// Validate token and extract player ID
	playerId, err := helpers.ValidateJWTAndGetPlayerID(token, conf.JWT.Secret)
	if err != nil {
		log.Printf("WebSocket connection rejected: invalid token - %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	client := service.NewClient(conn, h.Hub, playerId)

	go client.ReadPump()
	go client.WritePump()

	rm := h.RoundManager
	wsMessage := service.WsMessage{
		Type: service.WsMessageTypeGameState,
		Data: rm.GetGameState(),
	}

	directMessage := service.DirectMessage{
		Client:  client,
		Message: wsMessage,
	}

	h.Hub.SendDirect <- directMessage
	h.Hub.Register <- client
}
