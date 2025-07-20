package service

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
	"tradeoff/backend/internal/domain"
	"tradeoff/backend/internal/helpers"
)

type RoundManager struct {
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

func NewRoundManager(hub *Hub, marketService *MarketService, playerService *PlayerService) *RoundManager {
	rm := &RoundManager{
		hub:            hub,
		marketService:  marketService,
		playerService:  playerService,
		chartDataChan:  make(chan []domain.PriceData),
		hourlyDataChan: make(chan []domain.PriceData),
	}

	rm.transitionToLobby()
	return rm
}

func (r *RoundManager) Run() {
	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()

	for range timer.C {
		if time.Now().After(r.phaseEndTime) {
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

func (r *RoundManager) GetGameState(playerId string) GameStatePayload {
	session := r.playerService.GetPlayerSessionOrCreate(playerId)

	totalRealizedPnl, totalUnrealizedPnl := r.playerService.GetPlayerPnlData(playerId)

	longPositions, shortPositions := r.playerService.GetPositionsCount()

	return GameStatePayload{
		RoundID:   r.roundID,
		ChartData: r.chartData,
		PhaseChangePayload: PhaseChangePayload{
			Phase:   r.phase,
			EndTime: r.phaseEndTime,
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
	}
}

func (r *RoundManager) transitionToLive() {
	log.Println("--- Transitioning to Live Phase ---")
	r.phase = domain.Live
	r.phaseEndTime = time.Now().Add(LiveDuration)

	if len(r.chartData) == 0 || len(r.hourlyData) == 0 {
		log.Println("Failed to load chart data for live phase, transitioning to cooldown")
		r.transitionToCooldown()
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

	r.chartData = <-r.chartDataChan
	r.hourlyData = <-r.hourlyDataChan

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

func (r *RoundManager) sendPriceUpdate(priceData domain.PriceData) {
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
		r.sendPriceUpdate(priceData)

		r.sendPnlUpdate()
		i++
	}
	log.Println("--- Live Phase Finished ---")
}

func (r *RoundManager) CreatePosition(playerID string, positionType *domain.PositionType) (*domain.Position, error) {
	if len(r.chartData) == 0 {
		return nil, helpers.NewCustomError("no chart data available", http.StatusBadRequest)
	}

	currentPrice := r.chartData[len(r.chartData)-1].Close
	return r.playerService.CreatePosition(playerID, *positionType, currentPrice)
}

func (r *RoundManager) ClosePosition(playerID string) error {
	if len(r.chartData) == 0 {
		return helpers.NewCustomError("no chart data available", http.StatusBadRequest)
	}

	currentPrice := r.chartData[len(r.chartData)-1].Close
	_, err := r.playerService.ClosePosition(playerID, currentPrice)
	return err
}

func (r *RoundManager) sendPnlUpdate() {
	if len(r.chartData) == 0 {
		return
	}

	currentPrice := r.chartData[len(r.chartData)-1].Close

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


