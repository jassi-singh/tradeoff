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
	LiveDuration      = 1 * time.Minute
	CooldownDuration  = 10 * time.Second
	HourlyDataForDays = 10
	RoundDuration     = LobbyDuration + LiveDuration + CooldownDuration
	Ticker            = "X:BTCUSD"
)

func NewRoundManager(hub *Hub, marketService *MarketService) *RoundManager {
	return &RoundManager{
		hub:            hub,
		marketService:  marketService,
		chartDataChan:  make(chan []domain.PriceData),
		hourlyDataChan: make(chan []domain.PriceData),
	}
}

func (r *RoundManager) Run() {
	time.Sleep(3 * time.Second) // Give time for the hub to start

	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()

	for {
		if r.phase == "" {
			log.Println("--- Starting Round Manager ---")
			r.transitionToLobby()
		}

		<-timer.C
		r.updatePhase()
	}
}

func (r *RoundManager) updatePhase() {
	if r.phaseCountDown <= 0 {
		return
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

func (r *RoundManager) transitionToLive() {
	log.Println("--- Transitioning to Live Phase ---")
	r.phase = domain.Live
	r.phaseCountDown = int(LiveDuration.Seconds())

	if len(r.chartData) == 0 || len(r.hourlyData) == 0 {
		log.Println("Failed to load chart data for live phase, transitioning to cooldown")
		r.transitionToCooldown()
		return
	}

	data := map[string]any{
		"phase":         r.phase,
		"nextPhaseTime": time.Now().Add(LiveDuration),
	}
	r.broadcastRoundStatus(data)

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
	r.broadcastRoundStatus(data)
}

func (r *RoundManager) transitionToLobby() {
	log.Println("--- Transitioning to Lobby Phase ---")
	r.phase = domain.Lobby
	r.phaseCountDown = int(LobbyDuration.Seconds())

	r.loadMarketData()

	data := map[string]any{
		"phase":         r.phase,
		"nextPhaseTime": time.Now().Add(LobbyDuration),
	}
	r.broadcastRoundStatus(data)

	r.chartData = <-r.chartDataChan
	r.hourlyData = <-r.hourlyDataChan

	log.Printf("Loaded %d daily chart data and %d hourly data", len(r.chartData), len(r.hourlyData))

	data = map[string]any{
		"chartData": r.chartData,
	}

	r.hub.Broadcast <- WsMessage{
		Type: WsMessageTypeChartData,
		Data: data,
	}

}

func (r *RoundManager) broadcastRoundStatus(data map[string]any) {
	msg := WsMessage{
		Type: WsMessageTypeRoundStatus,
		Data: data,
	}
	r.hub.Broadcast <- msg
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (r *RoundManager) loadMarketData() {
	randomDecrease := -3 - int(rand.Float64()*10)
	go r.loadDailyChartData(randomDecrease)
	go r.loadHourlyChartData(randomDecrease)
}

func (r *RoundManager) loadDailyChartData(randomDecrease int) {
	from := truncateToDate(time.Now().UTC().AddDate(-2, 0, 0))
	to := truncateToDate(time.Now().UTC().AddDate(0, randomDecrease, 0))
	log.Printf("Loading daily chart data from %s to %s", from, to)
	limit := int(to.Sub(from).Hours() / 24)

	chartData, err := r.marketService.LoadPriceData(context.Background(), Ticker, from, to, "day", &limit)
	if err != nil {
		log.Println("Error loading price data:", err)
		r.chartDataChan <- nil
		return
	}
	log.Printf("Loaded %d daily chart data points", len(chartData))
	r.chartDataChan <- chartData
}

func (r *RoundManager) loadHourlyChartData(randomDecrease int) {
	from := truncateToDate(time.Now().UTC().AddDate(0, randomDecrease, 0))
	to := truncateToDate(from.AddDate(0, 0, HourlyDataForDays))
	log.Printf("Loading hourly chart data from %s to %s", from, to)

	limit := HourlyDataForDays * 24 * 60
	hourlyData, err := r.marketService.LoadPriceData(context.Background(), Ticker, from, to, "hour", &limit)
	if err != nil {
		log.Println("Error loading hourly price data:", err)
		r.hourlyDataChan <- nil
		return
	}
	log.Printf("Loaded %d hourly chart data points", len(hourlyData))
	r.hourlyDataChan <- hourlyData
}

func (r *RoundManager) runLivePhase() {
	log.Println("--- Running Live Phase ---")
	if len(r.hourlyData) == 0 {
		log.Println("no hourly data to broadcast")
		return
	}
	livePhaseTick := LiveDuration / time.Duration(len(r.hourlyData))
	log.Printf("Live phase tick duration: %s", livePhaseTick)
	ticker := time.NewTicker(livePhaseTick)
	defer ticker.Stop()

	i := 0
	for range ticker.C {
		if i >= len(r.hourlyData) || r.phase != domain.Live {
			break
		}

		priceData := r.hourlyData[i]

		if len(r.chartData) == 0 {
			r.chartData = append(r.chartData, priceData)
		} else {
			lastChartData := &r.chartData[len(r.chartData)-1]
			lastDataTime := time.Unix(lastChartData.Time, 0)
			priceDataTime := time.Unix(priceData.Time, 0)

			if lastDataTime.Day() != priceDataTime.Day() {
				r.chartData = append(r.chartData, priceData)
			} else {
				lastChartData.High = max(lastChartData.High, priceData.High)
				lastChartData.Low = min(lastChartData.Low, priceData.Low)
				lastChartData.Close = priceData.Close
				lastChartData.Volume += priceData.Volume
			}
		}

		msg := WsMessage{
			Type: WsMessageTypePriceUpdate,
			Data: priceData,
		}
		r.hub.Broadcast <- msg
		i++
	}
	log.Println("--- Live Phase Finished ---")
}
