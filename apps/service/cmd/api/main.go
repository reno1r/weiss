package main

import (
	"log"

	"github.com/reno1r/weiss/apps/service/internal/config"
	"github.com/reno1r/weiss/apps/service/internal/db"
	"github.com/reno1r/weiss/apps/service/internal/http"
)

// @title           Weiss API
// @version         1.0
// @description     This is the Weiss service API documentation for authentication endpoints.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@weiss.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
