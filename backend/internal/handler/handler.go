package handler

import (
	"tradeoff/backend/internal/config"
	"tradeoff/backend/internal/service"
)

type Handler struct {
	Hub           *service.Hub
	RoundManager  *service.RoundManager
	AuthService   *service.AuthService
	Config        *config.Config
}
func NewHandler(hub *service.Hub, roundManager *service.RoundManager, authService *service.AuthService, config *config.Config) *Handler {
	return &Handler{
		Hub:           hub,
		RoundManager:  roundManager,
		AuthService:   authService,
		Config:        config,
	}
}
