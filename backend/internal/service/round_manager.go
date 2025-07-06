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
	phaseCountDown int
	chartDataChan  chan []domain.PriceData
	chartData      []domain.PriceData
	hourlyDataChan chan []domain.PriceData
	hourlyData     []domain.PriceData
}

const (
	LobbyDuration     = 15 * time.Second
	LiveDuration      = 45 * time.Second
	CooldownDuration  = 10 * time.Second
	HourlyDataForDays = 10
	RoundDuration     = LobbyDuration + LiveDuration + CooldownDuration
	LivePhaseTick     = LiveDuration / (HourlyDataForDays * 24)
	Ticker            = "X:BTCUSD"
)

func NewRoundManager(hub *Hub, marketService *MarketService) *RoundManager {
	return &RoundManager{hub: hub, marketService: marketService}
}

func (r *RoundManager) Run() {

	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()

	for range timer.C {
		if r.phase == "" {
			log.Println("--- Starting Round Manager ---")
			r.transitionToLobby()
			continue
		}

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
	r.phaseCountDown = int(LiveDuration.Seconds())
	log.Println(r.chartDataChan, r.hourlyDataChan)
	r.chartData = <-r.chartDataChan
	r.hourlyData = <-r.hourlyDataChan

	log.Printf("Loaded %d daily chart data and %d hourly data", len(r.chartData), len(r.hourlyData))

	data := map[string]any{
		"phase":        r.phase,
		"phaseEndTime": time.Now().Add(LiveDuration).Unix(),
		"chartData":    r.chartData,
	}
	msg := WsMessage{
		Type: WsMessageTypeRoundStatus,
		Data: data,
	}
	r.hub.Broadcast <- msg

	go r.runLivePhase()
}

func (r *RoundManager) transitionToCooldown() {
	log.Println("--- Transitioning to Cooldown Phase ---")
	r.phase = domain.Closed
	r.phaseCountDown = int(CooldownDuration.Seconds())

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

	randomDecrease := -3 - int(rand.Float64()*10)
	go func() {
		from := time.Now().AddDate(-2, 0, 0) // 2 years ago
		to := time.Now().AddDate(0, randomDecrease, 0)
		log.Printf("Loading daily chart data from %s to %s", from, to)

		chartData, err := r.marketService.LoadPriceData(context.Background(), Ticker, from, to, "day")
		if err != nil {
			log.Println("Error loading price data:", err)
			r.chartDataChan <- nil
			return
		}

		log.Printf("Loaded %d daily chart data points", len(chartData))
		r.chartDataChan <- chartData
	}()

	go func() {
		from := time.Now().AddDate(0, randomDecrease, 1)
		to := time.Now().AddDate(0, randomDecrease, HourlyDataForDays)
		log.Printf("Loading hourly chart data from %s to %s", from, to)

		hourlyData, err := r.marketService.LoadPriceData(context.Background(), Ticker, from, to, "hour")
		if err != nil {
			log.Println("Error loading hourly price data:", err)
			r.hourlyDataChan <- nil
			return
		}
		log.Printf("Loaded %d hourly chart data points", len(hourlyData))
		r.hourlyDataChan <- hourlyData
	}()

	r.phase = domain.Lobby
	r.phaseCountDown = int(LobbyDuration.Seconds())

	data := map[string]any{
		"phase":         r.phase,
		"nextPhaseTime": time.Now().Add(LobbyDuration),
	}
	msg := WsMessage{
		Type: WsMessageTypeRoundStatus,
		Data: data,
	}
	r.hub.Broadcast <- msg
}

func (r *RoundManager) runLivePhase() {
	log.Println("--- Running Live Phase ---")
	if len(r.hourlyData) == 0 {
		log.Println("no hourly data to broadcast")
		return
	}
	ticker := time.NewTicker(LivePhaseTick)
	defer ticker.Stop()

	i := 0
	for range ticker.C {
		if i >= len(r.hourlyData) {
			break
		}
		lastChartData := r.chartData[len(r.chartData)-1]
		priceData := r.hourlyData[i]

		if lastChartData.Time.Day() != priceData.Time.Day() {
			r.chartData = append(r.chartData, priceData)
		} else {
			lastChartData.High = max(lastChartData.High, priceData.High)
			lastChartData.Low = min(lastChartData.Low, priceData.Low)
			lastChartData.Close = priceData.Close
			lastChartData.Volume += priceData.Volume
			r.chartData[len(r.chartData)-1] = priceData
		}

		msg := WsMessage{
			Type: WsMessageTypePriceData,
			Data: priceData,
		}
		r.hub.Broadcast <- msg
		i++
	}
	log.Println("--- Live Phase Finished ---")
}
