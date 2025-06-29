package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"tradeoff/backend/internal/handler"
)

func NewRouter(h *handler.Handler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TradeOff Game Server is running."))
	})

	appRouter := chi.NewRouter()

	appRouter.Post("/player", h.CreatePlayer)
	appRouter.Get("/player/{id}", h.GetPlayer)

	router.Mount("/api", appRouter)

	return router
}