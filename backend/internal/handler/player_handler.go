package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"tradeoff/backend/internal/helpers"
)

func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	player, err := h.PlayerService.GetPlayer(id)
	if err != nil {
		helpers.RespondWithError(w, err)
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, player)
}
