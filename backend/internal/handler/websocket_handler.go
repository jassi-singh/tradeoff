package handler

import (
	"log"
	"net/http"
	"tradeoff/backend/internal/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	client := service.NewClient(conn, h.Hub, r.URL.Query().Get("playerId"))

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

