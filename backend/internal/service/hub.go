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

	// user msgs
	WsMessageTypePositionUpdate WsMessageType = "position_update"
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
	Clients    map[string]*Client
	Broadcast  chan WsMessage
	Register   chan *Client
	Unregister chan *Client
	SendDirect chan DirectMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
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
			h.Clients[client.PlayerId] = client
			log.Println("Client registered", client.PlayerId)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.PlayerId]; ok {
				delete(h.Clients, client.PlayerId)
				close(client.send)
				log.Println("Client unregistered", client.PlayerId)
			}

		case message := <-h.Broadcast:
			for _,client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send) // TODO: check if this is correct
					delete(h.Clients, client.PlayerId)
				}
			}
		case directMessage := <-h.SendDirect:
			directMessage.Client.send <- directMessage.Message
		}

	}
}
