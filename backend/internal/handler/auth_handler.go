package handler

import (
	"encoding/json"
	"net/http"
	"tradeoff/backend/internal/helpers"
	"tradeoff/backend/internal/service"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, err)
		return
	}

	player, err := h.AuthService.Login(req)
	if err != nil {
		helpers.RespondWithError(w, err)
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, player)
}
