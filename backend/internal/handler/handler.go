package handler

import (
	"tradeoff/backend/internal/config"
	"tradeoff/backend/internal/service"
)

type Handler struct {
	PlayerService *service.PlayerService
	Hub           *service.Hub
	RoundManager  *service.RoundManager
	AuthService   *service.AuthService
	Config        *config.Config
}

func NewHandler(playerService *service.PlayerService, hub *service.Hub, roundManager *service.RoundManager, authService *service.AuthService, config *config.Config) *Handler {
	return &Handler{
		PlayerService: playerService,
		Hub:           hub,
		RoundManager:  roundManager,
		AuthService:   authService,
		Config:        config,
	}
}
