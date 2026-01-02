package main

import (
	"log"

	"github.com/reno1r/weiss/apps/service/internal/config"
	"github.com/reno1r/weiss/apps/service/internal/http"
)

func main() {
	config, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	server := http.NewServer(config)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
