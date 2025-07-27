package service

import (
	"log"
	"time"
	"tradeoff/backend/internal/domain"
)

type WsMsgType string

const (
	WsMsgTypeGameStateSync WsMsgType = "game_state_sync"
	WsMsgTypeNewRound      WsMsgType = "new_round"
	WsMsgTypePhaseUpdate   WsMsgType = "phase_update"
	WsMsgTypePriceUpdate   WsMsgType = "price_update"
	WsMsgTypeCountUpdate   WsMsgType = "count_update"

	WsMsgTypePnlUpdate         WsMsgType = "pnl_update"
	WsMsgTypeLeaderboardUpdate WsMsgType = "leaderboard_update"
)

type WsMessage struct {
	Type WsMsgType `json:"type"`
	Data any       `json:"data"`
}

// GameStatePayload is the data for the 'game_state_sync' and 'new_round' messages.
// It contains everything a client needs to render the game from scratch.
type GameStatePayload struct {
	RoundID             string             `json:"roundId"`
	ChartData           []domain.PriceData `json:"chartData"`
	TotalPnl            float64            `json:"pnl"`
	ActivePnl           float64            `json:"activePnl"`
	ActivePnlPercentage float64            `json:"activePnlPercentage"`
	PhaseChangePayload
	CountUpdatePayload
	domain.BasePlayerState
}

// PhaseChangePayload is the data for the 'phase_update' message.
type PhaseChangePayload struct {
	Phase   domain.Phase `json:"phase"`
	EndTime time.Time    `json:"endTime"` // Unix milliseconds
}

// CountUpdatePayload is the data for the 'count_update' message.
type CountUpdatePayload struct {
	LongPositions  int `json:"longPositions"`
	ShortPositions int `json:"shortPositions"`
	TotalPlayers   int `json:"totalPlayers"`
}

// PnlUpdatePayload is the data for the 'pnl_update' message.
// This is sent directly to a single player.
type PnlUpdatePayload struct {
	TotalPnl            float64 `json:"pnl"`
	Balance             float64 `json:"balance"`
	ActivePnl           float64 `json:"activePnl"`
	ActivePnlPercentage float64 `json:"activePnlPercentage"`
}

type PriceUpdate struct {
	PriceData  domain.PriceData `json:"priceData"`
	UpdateLast bool             `json:"updateLast"`
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
			for _, client := range h.Clients {
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
