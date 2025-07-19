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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body to get positionReq details
	var positionReq positionRequest
	if err := json.NewDecoder(r.Body).Decode(&positionReq); err != nil {
		helpers.RespondWithError(w, err)
		return
	}

	// Create the position using the RoundManager
	position, err := h.RoundManager.CreatePosition(userID, &positionReq.Type)
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Close the position using the RoundManager
	if err := h.RoundManager.ClosePosition(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
