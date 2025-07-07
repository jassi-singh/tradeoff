package service

import (
	"log"
)

type WsMessageType string

const (
	WsMessageTypePriceUpdate WsMessageType = "price_update"
	WsMessageTypeChartData   WsMessageType = "chart_data"
	WsMessageTypeRoundStatus WsMessageType = "round_status"
	WsMessageTypeGameState   WsMessageType = "game_state"
)

type WsMessage struct {
	Type WsMessageType `json:"type"`
	Data any           `json:"data"`
}

type DirectMessage struct {
	Client  *Client   `json:"client"`
	Message WsMessage `json:"message"`
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan WsMessage
	Register   chan *Client
	Unregister chan *Client
	SendDirect chan DirectMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan WsMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		SendDirect: make(chan DirectMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("Client registered", client.PlayerId)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
				log.Println("Client unregistered", client.PlayerId)
			}

		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		case directMessage := <-h.SendDirect:
			select {
			case directMessage.Client.send <- directMessage.Message:
			}
		}

	}
}
