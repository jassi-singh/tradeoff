package service

import (
	"tradeoff/backend/internal/domain"
)

type PlayerService struct {
	repo PlayerRepository
}

func NewPlayerService(repo PlayerRepository) *PlayerService {
	return &PlayerService{repo: repo}
}

type LoginParams struct {
	Username string `json:"username"`
}

func (s *PlayerService) CreatePlayer(payload LoginParams) (domain.Player, error) {
	player := domain.Player{
		Username: payload.Username,
	}

	return s.repo.CreatePlayer(player)
}

func (s *PlayerService) GetPlayer(id string) (domain.Player, error) {
	return s.repo.GetPlayer(id)
}

func (s *PlayerService) UpdateRefreshToken(playerId, refreshToken string) error {
	player, err := s.repo.GetPlayer(playerId)
	if err != nil {
		return err
	}

	player.RefreshToken = refreshToken
	_, err = s.repo.CreatePlayer(player) // Assuming CreatePlayer updates the player if it exists
	return err
}
