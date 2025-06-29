package main

import (
	"log"
	"net/http"
	"os"

	"tradeoff/backend/internal/handler"
	"tradeoff/backend/internal/platform/router"
	"tradeoff/backend/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("TradeOff Game Server starting on port %s...", port)

	h := handler.NewHandler(store)
	r := router.NewRouter(h)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
