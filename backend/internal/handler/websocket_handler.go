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

	h.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}