package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"tradeoff/backend/internal/domain"
	"tradeoff/backend/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Store *storage.PostgresStore
}

func NewHandler(store *storage.PostgresStore) *Handler {
	return &Handler{Store: store}
}

func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var player domain.Player
	err = json.Unmarshal(body, &player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err = h.Store.CreatePlayer(player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(player)
}

func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	player, err := h.Store.GetPlayer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(player)
}