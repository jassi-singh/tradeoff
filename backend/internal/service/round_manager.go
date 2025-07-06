package service

import (
	"context"
	"log"
	"math/rand/v2"
	"time"
	"tradeoff/backend/internal/domain"
)

type RoundManager struct {
	hub            *Hub
	marketService  *MarketService
	phase          domain.Phase
	phaseCountDown uint
	chartData      []domain.PriceData
	hourlyData     []domain.PriceData
}

const (
	LobbyDuration    = 15 * time.Second
	LiveDuration     = 45 * time.Second
	CooldownDuration = 10 * time.Second
	RoundDuration    = LobbyDuration + LiveDuration + CooldownDuration
	BroadcastTick    = 500 * time.Millisecond
	Ticker           = "X:BTCUSD"
)

func NewRoundManager(hub *Hub, marketService *MarketService) *RoundManager {
	return &RoundManager{hub: hub, marketService: marketService}
}

func (r *RoundManager) Run() {

	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()

	for range timer.C {
		r.phaseCountDown--

		if r.phaseCountDown <= 0 {
			switch r.phase {
			case domain.Lobby:
				r.transitionToLive()
			case domain.Live:
				r.transitionToCooldown()
			case domain.Closed:
				r.transitionToLobby()
			}
		}

	}
}

func (r *RoundManager) transitionToLive() {
	log.Println("--- Transitioning to Live Phase ---")
	r.phase = domain.Live
	r.phaseCountDown = uint(LiveDuration.Seconds())

	data := map[string]any{
		"phase":        r.phase,
		"phaseEndTime": time.Now().Add(LiveDuration).Unix(),
	}
	msg := WsMessage{
		Type: WsMessageTypeRoundStatus,
		Data: data,
	}
	r.hub.Broadcast <- msg
}

func (r *RoundManager) transitionToCooldown() {
	log.Println("--- Transitioning to Cooldown Phase ---")
	r.phase = domain.Closed
	r.phaseCountDown = uint(CooldownDuration.Seconds())

	data := map[string]any{
		"phase":         r.phase,
		"nextPhaseTime": time.Now().Add(CooldownDuration),
	}
	msg := WsMessage{
		Type: WsMessageTypeRoundStatus,
		Data: data,
	}
	r.hub.Broadcast <- msg
}

func (r *RoundManager) transitionToLobby() {
	log.Println("--- Transitioning to Lobby Phase ---")

	chartDataChan := make(chan []domain.PriceData)
	hourlyDataChan := make(chan []domain.PriceData)
	randomDecrease := -3 - int(rand.Float64()*10)
	go func() {
		from := time.Now().AddDate(-2, 0, 0) // 2 years ago
		to := time.Now().AddDate(0, randomDecrease, 0)

		chartData, err := r.marketService.LoadPriceData(context.Background(), Ticker, from, to, "day")
		if err != nil {
			log.Println("Error loading price data:", err)
			chartDataChan <- nil
			return
		}

		chartDataChan <- chartData

	}()

	go func() {
		from := time.Now().AddDate(0, randomDecrease, 1)
		to := time.Now().AddDate(0, randomDecrease, 10)

		hourlyData, err := r.marketService.LoadPriceData(context.Background(), Ticker, from, to, "hour")
		if err != nil {
			log.Println("Error loading hourly price data:", err)
			hourlyDataChan <- nil
			return
		}
		hourlyDataChan <- hourlyData
	}()

	r.phase = domain.Lobby
	r.phaseCountDown = uint(LobbyDuration.Seconds())

	data := map[string]any{
		"phase":         r.phase,
		"nextPhaseTime": time.Now().Add(LobbyDuration),
	}
	msg := WsMessage{
		Type: WsMessageTypeRoundStatus,
		Data: data,
	}
	r.hub.Broadcast <- msg
	r.chartData = <-chartDataChan
	r.hourlyData = <-hourlyDataChan
}
