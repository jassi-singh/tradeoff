package service

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"
	"tradeoff/backend/internal/domain"
)

type RoundManager struct {
	mu             sync.RWMutex
	hub            *Hub
	marketService  *MarketService
	playerService  *PlayerService
	phase          domain.Phase
	phaseEndTime   time.Time
	roundID        string
	chartDataChan  chan []domain.PriceData
	chartData      []domain.PriceData
	hourlyDataChan chan []domain.PriceData
	hourlyData     []domain.PriceData
	ctx            context.Context
	cancel         context.CancelFunc
}

const (
	LobbyDuration     = 15 * time.Second
	LiveDuration      = 1 * time.Minute
	CooldownDuration  = 10 * time.Second
	HourlyDataForDays = 10
	RoundDuration     = LobbyDuration + LiveDuration + CooldownDuration
	Ticker            = "X:BTCUSD"
	StartingBalance   = 100.0
)

func NewRoundManager(ctx context.Context, hub *Hub, marketService *MarketService, playerService *PlayerService) *RoundManager {
	rmCtx, cancel := context.WithCancel(ctx)
	rm := &RoundManager{
		hub:            hub,
		marketService:  marketService,
		playerService:  playerService,
		chartDataChan:  make(chan []domain.PriceData),
		hourlyDataChan: make(chan []domain.PriceData),
		ctx:            rmCtx,
		cancel:         cancel,
	}

	rm.transitionToLobby()
	return rm
}

// Shutdown gracefully stops the round manager
func (r *RoundManager) Shutdown() {
	log.Println("Shutting down RoundManager...")
	r.cancel()
}

func (r *RoundManager) Run() {
	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-r.ctx.Done():
			log.Println("RoundManager context cancelled, stopping...")
			return
		case <-timer.C:
			r.mu.RLock()
			phaseEndTime := r.phaseEndTime
			currentPhase := r.phase
			r.mu.RUnlock()

			if time.Now().After(phaseEndTime) {
				switch currentPhase {
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
}

func (r *RoundManager) GetGameState(playerId string) (GameStatePayload, error) {
	if playerId == "" {
		return GameStatePayload{}, fmt.Errorf("player ID cannot be empty")
	}

	r.mu.RLock()
	roundID := r.roundID
	chartData := r.chartData
	phase := r.phase
	phaseEndTime := r.phaseEndTime
	r.mu.RUnlock()

	session := r.playerService.GetPlayerSessionOrCreate(playerId)
	totalRealizedPnl, totalUnrealizedPnl := r.playerService.GetPlayerPnlData(playerId)
	longPositions, shortPositions := r.playerService.GetPositionsCount()

	return GameStatePayload{
		RoundID:   roundID,
		ChartData: chartData,
		PhaseChangePayload: PhaseChangePayload{
			Phase:   phase,
			EndTime: phaseEndTime,
		},
		CountUpdatePayload: CountUpdatePayload{
			TotalPlayers:   r.playerService.GetPlayerCount(),
			LongPositions:  longPositions,
			ShortPositions: shortPositions,
		},
		BasePlayerState: domain.BasePlayerState{
			Balance:         session.Balance,
			ActivePosition:  session.ActivePosition,
			ClosedPositions: session.ClosedPositions,
		},
		PnlUpdatePayload: PnlUpdatePayload{
			TotalRealizedPnl:   totalRealizedPnl,
			TotalUnrealizedPnl: totalUnrealizedPnl,
		},
	}, nil
}

func (r *RoundManager) transitionToLive() {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Println("--- Transitioning to Live Phase ---")
	r.phase = domain.Live
	r.phaseEndTime = time.Now().Add(LiveDuration)

	if len(r.chartData) == 0 || len(r.hourlyData) == 0 {
		log.Println("Failed to load chart data for live phase, transitioning to cooldown")
		r.transitionToCooldownUnsafe()
		return
	}

	data := PhaseChangePayload{
		Phase:   r.phase,
		EndTime: r.phaseEndTime,
	}
	r.broadcastPhaseUpdate(data)

	go r.runLivePhase()
}

func (r *RoundManager) transitionToCooldown() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transitionToCooldownUnsafe()
}

func (r *RoundManager) transitionToCooldownUnsafe() {
	log.Println("--- Transitioning to Cooldown Phase ---")
	r.phase = domain.Closed
	r.phaseEndTime = time.Now().Add(CooldownDuration)

	data := PhaseChangePayload{
		Phase:   r.phase,
		EndTime: r.phaseEndTime,
	}
	r.broadcastPhaseUpdate(data)
}

func (r *RoundManager) transitionToLobby() {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Println("--- Transitioning to Lobby Phase ---")
	r.phase = domain.Lobby
	r.phaseEndTime = time.Now().Add(LobbyDuration)
	r.roundID = generateUUID()

	// Reset all existing players for the new round
	playerCount := r.playerService.GetPlayerCount()
	if playerCount > 0 {
		r.playerService.ResetAllPlayers()
		log.Printf("Reset %d players for new round %s", playerCount, r.roundID)
	}

	r.loadMarketData()

	data := PhaseChangePayload{
		Phase:   r.phase,
		EndTime: r.phaseEndTime,
	}
	r.broadcastPhaseUpdate(data)

	// Wait for market data with timeout
	select {
	case r.chartData = <-r.chartDataChan:
		log.Printf("Loaded %d daily chart data points", len(r.chartData))
	case <-r.ctx.Done():
		log.Println("Context cancelled while loading chart data")
		return
	case <-time.After(30 * time.Second):
		log.Println("Timeout loading chart data")
		r.chartData = []domain.PriceData{}
	}

	select {
	case r.hourlyData = <-r.hourlyDataChan:
		log.Printf("Loaded %d hourly data points", len(r.hourlyData))
	case <-r.ctx.Done():
		log.Println("Context cancelled while loading hourly data")
		return
	case <-time.After(30 * time.Second):
		log.Println("Timeout loading hourly data")
		r.hourlyData = []domain.PriceData{}
	}

	log.Printf("Loaded %d daily chart data and %d hourly data", len(r.chartData), len(r.hourlyData))

	longPositions, shortPositions := r.playerService.GetPositionsCount()

	gameState := GameStatePayload{
		RoundID:            r.roundID,
		ChartData:          r.chartData,
		PhaseChangePayload: data,
		CountUpdatePayload: CountUpdatePayload{
			TotalPlayers:   r.playerService.GetPlayerCount(),
			LongPositions:  longPositions,
			ShortPositions: shortPositions,
		},
		BasePlayerState: domain.BasePlayerState{
			Balance:         StartingBalance,
			ActivePosition:  nil,
			ClosedPositions: []domain.ClosedPosition{},
		},
		PnlUpdatePayload: PnlUpdatePayload{
			TotalRealizedPnl:   0,
			TotalUnrealizedPnl: 0,
		},
	}

	r.hub.Broadcast <- WsMessage{
		Type: WsMsgTypeNewRound,
		Data: gameState,
	}
}

func (r *RoundManager) broadcastPhaseUpdate(data PhaseChangePayload) {
	msg := WsMessage{
		Type: WsMsgTypePhaseUpdate,
		Data: data,
	}
	r.hub.Broadcast <- msg
}

func generateUUID() string {
	b := make([]byte, 16)
	_, err := cryptorand.Read(b)
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
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

	chartData, err := r.marketService.LoadPriceData(r.ctx, Ticker, from, to, "day", &limit)
	if err != nil {
		log.Printf("Error loading daily price data: %v", err)
		r.chartDataChan <- nil
		return
	}
	r.chartDataChan <- chartData
}

func (r *RoundManager) loadHourlyChartData(randomDecrease int) {
	from := truncateToDate(time.Now().UTC().AddDate(0, randomDecrease, 0))
	to := truncateToDate(from.AddDate(0, 0, HourlyDataForDays))
	log.Printf("Loading hourly chart data from %s to %s", from, to)

	limit := HourlyDataForDays * 24 * 60
	hourlyData, err := r.marketService.LoadPriceData(r.ctx, Ticker, from, to, "hour", &limit)
	if err != nil {
		log.Printf("Error loading hourly price data: %v", err)
		r.hourlyDataChan <- nil
		return
	}
	r.hourlyDataChan <- hourlyData
}

func (r *RoundManager) sendPriceUpdate(priceData domain.PriceData) {
	r.mu.Lock()
	defer r.mu.Unlock()

	updateLast := false
	if len(r.chartData) == 0 {
		r.chartData = append(r.chartData, priceData)
	} else {
		lastChartData := &r.chartData[len(r.chartData)-1]
		lastDataTime := time.Unix(lastChartData.Time, 0)
		priceDataTime := time.Unix(priceData.Time, 0)

		if lastDataTime.Day() != priceDataTime.Day() {
			r.chartData = append(r.chartData, priceData)
		} else {
			updateLast = true
			lastChartData.High = max(lastChartData.High, priceData.High)
			lastChartData.Low = min(lastChartData.Low, priceData.Low)
			lastChartData.Close = priceData.Close
			lastChartData.Volume += priceData.Volume
		}
	}

	lastChartData := r.chartData[len(r.chartData)-1]
	msg := WsMessage{
		Type: WsMsgTypePriceUpdate,
		Data: PriceUpdate{
			PriceData:  lastChartData,
			UpdateLast: updateLast,
		},
	}
	r.hub.Broadcast <- msg
}

func (r *RoundManager) runLivePhase() {
	log.Println("--- Running Live Phase ---")

	r.mu.RLock()
	hourlyData := r.hourlyData
	currentPhase := r.phase
	r.mu.RUnlock()

	if len(hourlyData) == 0 {
		log.Println("No hourly data to broadcast")
		return
	}

	livePhaseTick := LiveDuration / time.Duration(len(hourlyData))
	log.Printf("Live phase tick duration: %s", livePhaseTick)
	ticker := time.NewTicker(livePhaseTick)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-r.ctx.Done():
			log.Println("Live phase cancelled")
			return
		case <-ticker.C:
			r.mu.RLock()
			currentPhase = r.phase
			r.mu.RUnlock()

			if i >= len(hourlyData) || currentPhase != domain.Live {
				log.Println("--- Live Phase Finished ---")
				return
			}

			priceData := hourlyData[i]
			r.sendPriceUpdate(priceData)
			r.sendPnlUpdate()
			i++
		}
	}
}

func (r *RoundManager) sendPnlUpdate() {
	if len(r.chartData) == 0 {
		return
	}

	currentPrice := r.GetCurrentPrice()
	if currentPrice == 0 {
		log.Println("Warning: Current price is 0, skipping PnL update")
		return
	}

	// Update PnL for all players with active positions
	r.playerService.UpdateAllPlayerPnl(currentPrice)

	// Get all sessions and send individual PnL updates
	sessions := r.playerService.GetAllSessions()
	for playerID, session := range sessions {
		if session.ActivePosition != nil {
			totalRealizedPnl, totalUnrealizedPnl := r.playerService.GetPlayerPnlData(playerID)

			msg := WsMessage{
				Type: WsMsgTypePnlUpdate,
				Data: PnlUpdatePayload{
					TotalUnrealizedPnl: totalUnrealizedPnl,
					TotalRealizedPnl:   totalRealizedPnl,
				},
			}

			client, exists := r.hub.Clients[playerID]
			if !exists {
				log.Printf("Warning: Client not found for player %s", playerID)
				continue
			}

			directMsg := DirectMessage{
				Client:  client,
				Message: msg,
			}
			r.hub.SendDirect <- directMsg
		}
	}
}

func (r *RoundManager) GetCurrentPrice() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.chartData) == 0 {
		return 0
	}
	return r.chartData[len(r.chartData)-1].Close
}