package main

import (
	"log"
	"net/http"
	"os"

	"tradeoff/backend/internal/platform/router"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("TradeOff Game Server starting on port %s...", port)

	r := router.NewRouter()

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
