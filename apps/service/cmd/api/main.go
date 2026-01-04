package main

import (
	"log"

	"github.com/reno1r/weiss/apps/service/internal/config"
	"github.com/reno1r/weiss/apps/service/internal/db"
	"github.com/reno1r/weiss/apps/service/internal/http"
)

func main() {
	config, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.NewDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect database %v", err)
	}

	server := http.NewServer(config, database.DB())

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
