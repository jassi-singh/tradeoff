package handler

import (
	"encoding/json"
	"net/http"
	"tradeoff/backend/internal/helpers"
	"tradeoff/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req service.CreatePlayerParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	player, err := h.PlayerService.CreatePlayer(req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to create player")
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, player)
}

func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	player, err := h.PlayerService.GetPlayer(id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to get player")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, player)
}