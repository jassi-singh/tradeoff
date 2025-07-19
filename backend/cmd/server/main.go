package main

import (
	"log"
	"net/http"

	"tradeoff/backend/internal/config"
	"tradeoff/backend/internal/handler"
	"tradeoff/backend/internal/platform/router"
	"tradeoff/backend/internal/service"
	"tradeoff/backend/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	store, err := storage.NewPostgresStore(*config)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	playerService := service.NewPlayerService(store)
	authService := service.NewAuthService(store, config.JWT.Secret, config.JWT.Expiration)

	hub := service.NewHub()
	go hub.Run()

	marketService := service.NewMarketService(hub, config.Polygon.APIKey)
	roundManager := service.NewRoundManager(hub, marketService)
	go roundManager.Run()

	handler := handler.NewHandler(playerService, hub, roundManager, authService, config)
	router := router.NewRouter(handler, config)

	log.Printf("TradeOff Game Server starting on port %s...", config.Server.Port)

	if err := http.ListenAndServe(":"+config.Server.Port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
