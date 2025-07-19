package handler

import (
	"encoding/json"
	"net/http"
	"tradeoff/backend/internal/helpers"
	"tradeoff/backend/internal/service"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

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

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, helpers.NewCustomError("Invalid request body", http.StatusBadRequest))
		return
	}

	if req.RefreshToken == "" {
		helpers.RespondWithError(w, helpers.NewCustomError("Refresh token is required", http.StatusBadRequest))
		return
	}

	result, err := h.AuthService.RefreshToken(req.RefreshToken)
	if err != nil {
		helpers.RespondWithError(w, helpers.NewCustomError("Invalid or expired refresh token", http.StatusUnauthorized))
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, result)
}
