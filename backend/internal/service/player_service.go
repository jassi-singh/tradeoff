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

type CreatePlayerParams struct {
	Username string `json:"username"`
}

func (s *PlayerService) CreatePlayer(payload CreatePlayerParams) (domain.Player, error) {
	player := domain.Player{
		Username: payload.Username,
	}

	return s.repo.CreatePlayer(player)
}

func (s *PlayerService) GetPlayer(id string) (domain.Player, error) {
	return s.repo.GetPlayer(id)
}