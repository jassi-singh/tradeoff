package handler

import (
	"encoding/json"
	"net/http"
	"tradeoff/backend/internal/domain"
	"tradeoff/backend/internal/helpers"
)

type positionRequest struct {
	Type domain.PositionType `json:"type"`
}

func (h *Handler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("userId").(string)
	if !ok {
		helpers.RespondWithError(w, helpers.NewCustomError("Unauthorized", http.StatusUnauthorized))
		return
	}

	// Parse the request body to get positionReq details
	var positionReq positionRequest
	if err := json.NewDecoder(r.Body).Decode(&positionReq); err != nil {
		helpers.RespondWithError(w, err)
		return
	}

	currentPrice := h.RoundManager.GetCurrentPrice()
	position, err := h.PlayerService.CreatePosition(userID, positionReq.Type, currentPrice)
	if err != nil {
		helpers.RespondWithError(w, err)
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, position)

}

func (h *Handler) ClosePosition(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("userId").(string)
	if !ok {
		helpers.RespondWithError(w, helpers.NewCustomError("Unauthorized", http.StatusUnauthorized))
		return
	}

	currentPrice := h.RoundManager.GetCurrentPrice()
	_, err := h.PlayerService.ClosePosition(userID, currentPrice)
	if err != nil {
		helpers.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
