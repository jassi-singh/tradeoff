package handler

import (
	"net/http"
	"tradeoff/backend/internal/config"
	"tradeoff/backend/internal/helpers"
	"tradeoff/backend/internal/service"
)

type Handler struct {
	Hub          *service.Hub
	RoundManager *service.RoundManager
	AuthService  *service.AuthService
	Config       *config.Config
}

func NewHandler(hub *service.Hub, roundManager *service.RoundManager, authService *service.AuthService, config *config.Config) *Handler {
	return &Handler{
		Hub:          hub,
		RoundManager: roundManager,
		AuthService:  authService,
		Config:       config,
	}
}

func (h *Handler) GetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get player session information
	session := h.RoundManager.GetPlayerSessionOrCreate(userID)

	helpers.RespondWithJSON(w, http.StatusOK, session)
}
