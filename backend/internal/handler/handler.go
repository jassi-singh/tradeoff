package handler

import (
	"tradeoff/backend/internal/service"
)

type Handler struct {
	PlayerService *service.PlayerService
	Hub           *service.Hub
	RoundManager  *service.RoundManager
}

func NewHandler(playerService *service.PlayerService, hub *service.Hub, roundManager *service.RoundManager) *Handler {
	return &Handler{PlayerService: playerService, Hub: hub, RoundManager: roundManager}
}

