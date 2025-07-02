package main

import (
	"log"
	"net/http"
	"os"

	"tradeoff/backend/internal/handler"
	"tradeoff/backend/internal/platform/router"
	"tradeoff/backend/internal/service"
	"tradeoff/backend/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	playerService := service.NewPlayerService(store)

	hub := service.NewHub()
	go hub.Run()

	marketService := service.NewMarketService(hub)
	marketService.LoadPriceData()
	go marketService.StartPriceFeed()

	handler := handler.NewHandler(playerService, hub)
	router := router.NewRouter(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("TradeOff Game Server starting on port %s...", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
