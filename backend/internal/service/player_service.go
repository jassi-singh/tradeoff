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

func (s *PlayerService) GetPlayer(id string) (domain.Player, error) {
	return s.repo.GetPlayer(id)
}
