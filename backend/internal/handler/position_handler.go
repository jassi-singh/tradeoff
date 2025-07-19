package handler

import (
	"encoding/json"
	"net/http"
	"tradeoff/backend/internal/domain"
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

	// Parse the request body to get position details
	var position positionRequest
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create the position using the RoundManager
	if err := h.RoundManager.CreatePosition(userID, &position.Type); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(position)
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
