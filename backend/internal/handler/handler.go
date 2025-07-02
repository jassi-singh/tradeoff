package handler

import (
	"tradeoff/backend/internal/service"
)

type Handler struct {
	PlayerService *service.PlayerService
	Hub           *service.Hub
}

func NewHandler(playerService *service.PlayerService, hub *service.Hub) *Handler {
	return &Handler{PlayerService: playerService, Hub: hub}
}