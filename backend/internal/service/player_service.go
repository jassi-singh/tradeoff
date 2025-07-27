package service

import (
	"errors"
	"sort"
	"sync"
	"time"
	"tradeoff/backend/internal/domain"
)

// PlayerService is the sole, concurrent-safe owner of all live player state for a round.
type PlayerService struct {
	playerSessions map[string]*domain.PlayerState
	mu             sync.RWMutex
}

func NewPlayerService() *PlayerService {
	return &PlayerService{
		playerSessions: make(map[string]*domain.PlayerState),
	}
}

// GetPlayerSessionOrCreate is the safe way to get or create a session.
// It handles the "check-then-act" concurrency problem correctly.
func (s *PlayerService) GetPlayerSessionOrCreate(playerID string, username *string) *domain.PlayerState {
	// First, try with just a read lock for performance.
	s.mu.RLock()
	session, exists := s.playerSessions[playerID]
	s.mu.RUnlock()

	if exists {
		return session
	}

	// If the session doesn't exist, we need a full write lock to create it.
	s.mu.Lock()
	defer s.mu.Unlock()

	// CRITICAL: We must check again after acquiring the write lock.
	// Another goroutine might have created the session in the tiny gap
	// between our RUnlock and Lock calls.
	session, exists = s.playerSessions[playerID]
	if exists {
		return session
	}

	// If it still doesn't exist, we are safe to create it.
	newSession := &domain.PlayerState{
		PlayerId: playerID,
		Username: *username,
		BasePlayerState: domain.BasePlayerState{
			Balance:         StartingBalance,
			ActivePosition:  nil,
			ClosedPositions: []domain.ClosedPosition{},
		},
	}
	s.playerSessions[playerID] = newSession
	return newSession
}

// CreatePosition now returns an error if an action is invalid.
func (s *PlayerService) CreatePosition(playerID string, positionType domain.PositionType, entryPrice float64) (*domain.Position, error) {
	s.mu.Lock() // We need a full write lock since we are modifying the session.
	defer s.mu.Unlock()

	session, exists := s.playerSessions[playerID]
	if !exists {
		// This case should ideally not happen if GetPlayerSessionOrCreate is called on connect.
		return nil, errors.New("player session not found")
	}

	if session.ActivePosition != nil {
		return nil, errors.New("player already has an active position")
	}

	if session.Balance == 0 {
		return nil, errors.New("player has no balance")
	}

	quantity := session.Balance / entryPrice

	position := &domain.Position{
		Type:       positionType,
		EntryPrice: entryPrice,
		EntryTime:  time.Now(),
		Quantity:   quantity,
	}
	session.ActivePosition = position
	session.Balance = 0

	return position, nil
}

// ClosePosition now returns an error for invalid states.
func (s *PlayerService) ClosePosition(playerID string, closePrice float64) (*domain.ClosedPosition, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.playerSessions[playerID]
	if !exists {
		return nil, errors.New("player session not found")
	}

	if session.ActivePosition == nil {
		return nil, errors.New("no active position to close")
	}

	activePosition := session.ActivePosition
	pnl, pnlPercentage := s.calculatePnl(activePosition, closePrice)
	initialInvestment := activePosition.Quantity * activePosition.EntryPrice

	closedPosition := domain.ClosedPosition{
		Position: domain.Position{
			Quantity:      activePosition.Quantity,
			Type:          activePosition.Type,
			EntryPrice:    activePosition.EntryPrice,
			EntryTime:     activePosition.EntryTime,
			Pnl:           pnl,
			PnlPercentage: pnlPercentage,
		},
		ExitPrice: closePrice,
		ExitTime:  time.Now(),
	}

	session.ClosedPositions = append(session.ClosedPositions, closedPosition)
	session.ActivePosition = nil
	session.Balance = initialInvestment + pnl

	return &closedPosition, nil
}

func (s *PlayerService) calculatePnl(position *domain.Position, currentPrice float64) (float64, float64) {
	pnl := (currentPrice - position.EntryPrice) * position.Quantity
	if position.Type == domain.PositionTypeShort {
		pnl *= -1
	}
	pnlPercentage := (pnl / (position.Quantity * position.EntryPrice)) * 100
	return pnl, pnlPercentage
}

func (s *PlayerService) GetPlayerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.playerSessions)
}

func (s *PlayerService) GetPositionsCount() (int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	longPositions := 0
	shortPositions := 0
	for _, session := range s.playerSessions {
		if session.ActivePosition != nil {
			if session.ActivePosition.Type == domain.PositionTypeLong {
				longPositions++
			} else {
				shortPositions++
			}
		}
	}
	return longPositions, shortPositions
}

func (s *PlayerService) ResetAllPlayers() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, session := range s.playerSessions {
		session.Balance = StartingBalance
		session.ActivePosition = nil
		session.ClosedPositions = []domain.ClosedPosition{}
	}
}

// GetAllSessions returns a copy of all player sessions for iteration
func (s *PlayerService) GetAllSessions() map[string]*domain.PlayerState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make(map[string]*domain.PlayerState)
	for playerID, session := range s.playerSessions {
		sessions[playerID] = session
	}
	return sessions
}

// UpdateAllPlayerPnl updates PnL for all players with active positions
func (s *PlayerService) UpdateAllPlayerPnl(currentPrice float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	pnlUpdated := false

	for _, session := range s.playerSessions {
		if session.ActivePosition != nil {
			pnl, pnlPercentage := s.calculatePnl(session.ActivePosition, currentPrice)
			session.ActivePosition.Pnl = pnl
			session.ActivePosition.PnlPercentage = pnlPercentage
			pnlUpdated = true
		}
	}
	return pnlUpdated
}

// GetPlayerStat returns PnL data for a specific player
func (s *PlayerService) GetPlayerStat(playerID string) (float64, float64, float64, float64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.playerSessions[playerID]
	if !exists {
		return 0, 0, 0, 0
	}

	totalRealizedPnl := 0.0
	for _, closedPosition := range session.ClosedPositions {
		totalRealizedPnl += closedPosition.Pnl
	}

	balance := session.Balance
	activePnl := 0.0
	activePnlPercentage := 0.0
	if session.ActivePosition != nil {
		activePnl = session.ActivePosition.Pnl
		activePnlPercentage = session.ActivePosition.PnlPercentage
	}

	return totalRealizedPnl, activePnl, balance, activePnlPercentage
}

func (s *PlayerService) GetLeaderboard() []domain.LeaderboardPlayer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	leaderboard := []domain.LeaderboardPlayer{}

	for _, session := range s.playerSessions {
		activeBalance := session.Balance
		if session.ActivePosition != nil {
			activeBalance += session.ActivePosition.EntryPrice*session.ActivePosition.Quantity + session.ActivePosition.Pnl
		}
		leaderboard = append(leaderboard, domain.LeaderboardPlayer{
			PlayerId:      session.PlayerId,
			Username:      session.Username,
			ActiveBalance: activeBalance,
		})
	}

	// Sort by balance in descending order
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].ActiveBalance > leaderboard[j].ActiveBalance
	})

	// Return only top 20 entries
	if len(leaderboard) > 20 {
		leaderboard = leaderboard[:20]
	}

	return leaderboard
}
