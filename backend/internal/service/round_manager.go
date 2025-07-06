package service

import (
	"log"
	"time"
	"tradeoff/backend/internal/domain"
)

type RoundManager struct {
	hub           *Hub
	marketService *MarketService
	roundStatus   domain.RoundStatus
	roundId       string
	roundStart    time.Time
	roundEnd      time.Time
}

const (
	LobbyDuration    = 15 * time.Second
	LiveDuration     = 45 * time.Second
	CooldownDuration = 10 * time.Second
	RoundDuration    = LobbyDuration + LiveDuration + CooldownDuration
	BroadcastTick    = 500 * time.Millisecond
)

func NewRoundManager(hub *Hub, marketService *MarketService) *RoundManager {
	return &RoundManager{hub: hub, marketService: marketService}
}

func (r *RoundManager) Run() {
	for {
		log.Println("--- New Round : Lobby ---")
		r.roundStatus = domain.RoundStatusLobby
		r.hub.Broadcast <- r.hub.NewRoundStatusMessage(r.roundStatus)
		<-time.After(LobbyDuration)
		log.Println("---Round : Live ---")
		r.roundStatus = domain.RoundStatusLive
		r.hub.Broadcast <- r.hub.NewRoundStatusMessage(r.roundStatus)
		<-time.After(LiveDuration)

		log.Println("---Round : Cooldown ---")
	}
}
