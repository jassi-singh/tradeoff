package router

import (
	"net/http"

	"tradeoff/backend/internal/config"
	"tradeoff/backend/internal/handler"
	"tradeoff/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(h *handler.Handler, config *config.Config) *chi.Mux {
	router := chi.NewRouter()

	// CORS configuration
	corsMiddleware := cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})

	router.Use(corsMiddleware)
	router.Use(chiMiddleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TradeOff Game Server is running."))
	})

	appRouter := chi.NewRouter()

	appRouter.Post("/login", h.Login)
	appRouter.Post("/refresh", h.RefreshToken)

	// appRouter.With(middleware.AuthMiddleware(h.Config)).Get("/player", h.GetPlayerInfo)
	appRouter.With(middleware.AuthMiddleware(h.Config)).Post("/position", h.CreatePosition)
	appRouter.With(middleware.AuthMiddleware(h.Config)).Post("/close-position", h.ClosePosition)

	router.Mount("/api", appRouter)

	router.Get("/ws", h.HandleWebSocket)

	return router
}
