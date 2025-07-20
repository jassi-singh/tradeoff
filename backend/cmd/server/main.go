package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	authService := service.NewAuthService(store, config.JWT.Secret, config.JWT.Expiration)

	hub := service.NewHub()
	go hub.Run()

	marketService := service.NewMarketService(hub, config.Polygon.APIKey)
	playerService := service.NewPlayerService()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	roundManager := service.NewRoundManager(ctx, hub, marketService, playerService)
	go roundManager.Run()

	handler := handler.NewHandler(hub, roundManager, authService, config, playerService)
	router := router.NewRouter(handler, config)

	// Create server
	server := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("TradeOff Game Server starting on port %s...", config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown RoundManager
	roundManager.Shutdown()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
