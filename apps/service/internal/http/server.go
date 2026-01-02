package http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/reno1r/weiss/apps/service/internal/config"
)

type Server struct {
	app    *fiber.App
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		app: fiber.New(fiber.Config{
			AppName: config.AppName,
		}),
		config: config,
	}
}

func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf("%s:%s", s.config.AppHost, s.config.AppPort))
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}
