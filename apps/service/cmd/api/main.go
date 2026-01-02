package main

import (
	"log"

	"github.com/reno1r/weiss/apps/service/internal/config"
)

func main() {
	config, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println(config)
}
