package router

import (
	"net/http"

	"tradeoff/backend/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handler.Handler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TradeOff Game Server is running."))
	})

	appRouter := chi.NewRouter()

	appRouter.Post("/player", h.CreatePlayer)
	appRouter.Get("/player/{id}", h.GetPlayer)

	router.Mount("/api", appRouter)

	router.Get("/ws", h.HandleWebSocket)

	return router
}