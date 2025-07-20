package handler

import (
	"tradeoff/backend/internal/config"
	"tradeoff/backend/internal/service"
)

type Handler struct {
	Hub          *service.Hub
	RoundManager *service.RoundManager
	PlayerService *service.PlayerService
	AuthService  *service.AuthService
	Config       *config.Config
}

func NewHandler(hub *service.Hub, roundManager *service.RoundManager, authService *service.AuthService, config *config.Config, playerService *service.PlayerService) *Handler {
	return &Handler{
		Hub:          hub,
		RoundManager: roundManager,
		AuthService:  authService,
		Config:       config,
		PlayerService: playerService,
	}
}