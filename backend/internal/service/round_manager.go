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
	phase          domain.Phase
	phaseCountDown int
	roundID        string
	chartDataChan  chan []domain.PriceData
	chartData      []domain.PriceData
	hourlyDataChan chan []domain.PriceData
	hourlyData     []domain.PriceData
	playerSessions map[string]domain.PlayerState
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

func NewRoundManager(hub *Hub, marketService *MarketService) *RoundManager {
	return &RoundManager{
		hub:            hub,
		marketService:  marketService,
		chartDataChan:  make(chan []domain.PriceData),
		hourlyDataChan: make(chan []domain.PriceData),
		playerSessions: make(map[string]domain.PlayerState),
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

func (r *RoundManager) GetGameState(playerId string) GameStatePayload {
	session := r.GetPlayerSessionOrCreate(playerId)

	totalRealizedPnl := 0.0
	for _, closedPosition := range session.ClosedPositions {
		totalRealizedPnl += closedPosition.Pnl
	}
	totalUnrealizedPnl := 0.0
	if session.ActivePosition != nil {
		totalUnrealizedPnl = session.ActivePosition.Pnl
	}

	longPositions := 0
	shortPositions := 0
	for _, playerSession := range r.playerSessions {
		if playerSession.ActivePosition != nil {
			if playerSession.ActivePosition.Type == domain.PositionTypeLong {
				longPositions++
			} else {
				shortPositions++
			}
		}
	}

	return GameStatePayload{
		RoundID: r.roundID,
		ChartData: r.chartData,
		PhaseChangePayload: PhaseChangePayload{
			Phase:   r.phase,
			EndTime: time.Now().Add(time.Duration(r.phaseCountDown) * time.Second),
		},
		CountUpdatePayload: CountUpdatePayload{
			TotalPlayers:   r.GetActivePlayerCount(),
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

	data := PhaseChangePayload{
		Phase:   r.phase,
		EndTime: time.Now().Add(LiveDuration),
	}
	r.broadcastPhaseUpdate(data)

	go r.runLivePhase()
}

func (r *RoundManager) transitionToCooldown() {
	log.Println("--- Transitioning to Cooldown Phase ---")
	r.phase = domain.Closed
	r.phaseCountDown = int(CooldownDuration.Seconds())

	data := PhaseChangePayload{
		Phase:   r.phase,
		EndTime: time.Now().Add(CooldownDuration),
	}
	r.broadcastPhaseUpdate(data)
}

func (r *RoundManager) CreatePlayerSession(playerID string) {
	session := domain.PlayerState{
		PlayerId: playerID,
		BasePlayerState: domain.BasePlayerState{
			Balance:         StartingBalance,
			ActivePosition:  nil,
			ClosedPositions: []domain.ClosedPosition{},
		},
	}
	r.playerSessions[playerID] = session
	log.Printf("Created player session for %s with balance %.2f in round %s", playerID, StartingBalance, r.roundID)
}

func (r *RoundManager) ResetAllPlayers() {
	log.Println("--- Resetting all player balances and positions ---")
	for playerID := range r.playerSessions {
		session := r.playerSessions[playerID]
		session.Balance = StartingBalance
		session.ActivePosition = nil
		session.ClosedPositions = []domain.ClosedPosition{}
		r.playerSessions[playerID] = session
		log.Printf("Reset player %s: balance=%.2f, positions=0", playerID, StartingBalance)
	}
}

func (r *RoundManager) GetPlayerSessionOrCreate(playerID string) domain.PlayerState {
	session, exists := r.playerSessions[playerID]
	if !exists {
		r.CreatePlayerSession(playerID)
		session = r.playerSessions[playerID]
	}
	return session
}

func (r *RoundManager) GetActivePlayerCount() int {
	return len(r.playerSessions)
}

func (r *RoundManager) transitionToLobby() {
	log.Println("--- Transitioning to Lobby Phase ---")
	r.phase = domain.Lobby
	r.phaseCountDown = int(LobbyDuration.Seconds())
	r.roundID = generateUUID()

	// Reset all existing players for the new round
	playerCount := r.GetActivePlayerCount()
	if playerCount > 0 {
		r.ResetAllPlayers()
		log.Printf("Reset %d players for new round %s", playerCount, r.roundID)
	}

	r.loadMarketData()

	data := PhaseChangePayload{
		Phase:   r.phase,
		EndTime: time.Now().Add(LobbyDuration),
	}
	r.broadcastPhaseUpdate(data)

	r.chartData = <-r.chartDataChan
	r.hourlyData = <-r.hourlyDataChan

	log.Printf("Loaded %d daily chart data and %d hourly data", len(r.chartData), len(r.hourlyData))

	gameState := GameStatePayload{
		RoundID:            r.roundID,
		ChartData:          r.chartData,
		PhaseChangePayload: data,
		CountUpdatePayload: CountUpdatePayload{
			TotalPlayers:   r.GetActivePlayerCount(),
			LongPositions:  0,
			ShortPositions: 0,
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
	session, exists := r.playerSessions[playerID]
	if !exists {
		return nil, helpers.NewCustomError("player session not found", http.StatusNotFound)
	}

	if session.ActivePosition != nil {
		return nil, helpers.NewCustomError("position already active", http.StatusBadRequest)
	}

	session.ActivePosition = &domain.Position{
		Type:          *positionType,
		EntryPrice:    r.chartData[len(r.chartData)-1].Close,
		EntryTime:     time.Now(),
		Pnl:           0,
		PnlPercentage: 0,
	}

	r.playerSessions[playerID] = session
	return session.ActivePosition, nil
}

func (r *RoundManager) ClosePosition(playerID string) error {
	session, exists := r.playerSessions[playerID]
	if !exists {
		return helpers.NewCustomError("player session not found", http.StatusNotFound)
	}

	if session.ActivePosition == nil {
		return helpers.NewCustomError("no active position to close", http.StatusBadRequest)
	}

	exitPrice := r.chartData[len(r.chartData)-1].Close
	profitLoss, profitLossPercentage := r.calculatePnl(exitPrice, session.ActivePosition.EntryPrice, session.ActivePosition.Type)

	closedPosition := domain.ClosedPosition{
		Position: domain.Position{
			Type:          session.ActivePosition.Type,
			EntryPrice:    session.ActivePosition.EntryPrice,
			EntryTime:     session.ActivePosition.EntryTime,
			Pnl:           profitLoss,
			PnlPercentage: profitLossPercentage,
		},
		ExitPrice: exitPrice,
		ExitTime:  time.Now(),
	}

	session.ClosedPositions = append(session.ClosedPositions, closedPosition)
	session.ActivePosition = nil

	r.playerSessions[playerID] = session
	return nil
}

func (r *RoundManager) sendPnlUpdate() {
	for _, playerSession := range r.playerSessions {
		if playerSession.ActivePosition != nil {

			currentPrice := r.chartData[len(r.chartData)-1].Close
			pnl, pnlPercentage := r.calculatePnl(currentPrice, playerSession.ActivePosition.EntryPrice, playerSession.ActivePosition.Type)

			playerSession.ActivePosition.Pnl = pnl
			playerSession.ActivePosition.PnlPercentage = pnlPercentage

			totalRealizedPnl := 0.0
			for _, closedPosition := range playerSession.ClosedPositions {
				totalRealizedPnl += closedPosition.Pnl
			}

			msg := WsMessage{
				Type: WsMsgTypePnlUpdate,
				Data: PnlUpdatePayload{
					TotalUnrealizedPnl: playerSession.ActivePosition.Pnl,
					TotalRealizedPnl:   totalRealizedPnl,
				},
			}

			client, exists := r.hub.Clients[playerSession.PlayerId]
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

func (r *RoundManager) calculatePnl(currentPrice float64, entryPrice float64, positionType domain.PositionType) (float64, float64) {
	pnl := currentPrice - entryPrice
	if positionType == domain.PositionTypeShort {
		pnl = pnl * -1
	}
	pnlPercentage := (pnl / entryPrice) * 100
	return pnl, pnlPercentage
}
