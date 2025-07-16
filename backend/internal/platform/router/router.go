package router

import (
	"net/http"

	"tradeoff/backend/internal/handler"
	platformMiddleware "tradeoff/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(h *handler.Handler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	router.Use(chiMiddleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TradeOff Game Server is running."))
	})

	appRouter := chi.NewRouter()

	appRouter.Post("/login", h.Login)
	appRouter.With(platformMiddleware.AuthMiddleware).Get("/player/{id}", h.GetPlayer)

	router.Mount("/api", appRouter)

	router.Get("/ws", h.HandleWebSocket)

	return router
}
