package handler

import (
	"tradeoff/backend/internal/service"
)

type Handler struct {
	PlayerService *service.PlayerService
	Hub           *service.Hub
	RoundManager  *service.RoundManager
	AuthService   *service.AuthService
}

func NewHandler(playerService *service.PlayerService, hub *service.Hub, roundManager *service.RoundManager, authService *service.AuthService) *Handler {
	return &Handler{PlayerService: playerService, Hub: hub, RoundManager: roundManager, AuthService: authService}
}
