package http

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"gorm.io/gorm"

	"github.com/reno1r/weiss/apps/service/internal/app/auth/services"
	"github.com/reno1r/weiss/apps/service/internal/app/auth/usecases"
	"github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	"github.com/reno1r/weiss/apps/service/internal/config"
	"github.com/reno1r/weiss/apps/service/internal/http/handlers"
)

type Server struct {
	app    *fiber.App
	config *config.Config
	db     *gorm.DB
}

func NewServer(config *config.Config, db *gorm.DB) *Server {
	server := &Server{
		app: fiber.New(fiber.Config{
			AppName:         config.AppName,
			CaseSensitive:   true,
			StrictRouting:   true,
			BodyLimit:       20 * 1024 * 1024,
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			IdleTimeout:     120 * time.Second,
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			ErrorHandler:    defaultErrorHandler,
			TrustProxy:      true,
			TrustProxyConfig: fiber.TrustProxyConfig{
				Proxies: []string{"0.0.0.0/0"},
			},
			ProxyHeader: fiber.HeaderXForwardedFor,
		}),
		config: config,
		db:     db,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddleware() {
	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: s.config.AppDebug,
	}))

	s.app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	s.app.Use(helmet.New())

	s.app.Use(cors.New(s.getCorsConfig()))

	s.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.app.Use(idempotency.New())
}

func (s *Server) setupRoutes() {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": s.config.AppName,
		})
	})

	s.setupAuthRoutes()

}

func (s *Server) setupAuthRoutes() {
	userRepo := repositories.NewUserRepository(s.db)

	passwordService := services.NewPasswordService(s.config)
	tokenService, err := services.NewTokenService(s.config)

	if err != nil {
		return
	}

	registerUsecase := usecases.NewRegisterUsecase(&userRepo, passwordService)
	loginUsecase := usecases.NewLoginUsecase(&userRepo, tokenService, passwordService)

	registerHandler := handlers.NewRegisterHandler(registerUsecase)
	loginHandler := handlers.NewLoginHandler(loginUsecase)

	s.app.Post("/api/auth/register", registerHandler.Handle)
	s.app.Post("/api/auth/login", loginHandler.Handle)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.AppHost, s.config.AppPort)
	return s.app.Listen(addr)
}

func (s *Server) StartWithContext(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%s", s.config.AppHost, s.config.AppPort)
	return s.app.Listen(addr)
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) getCorsConfig() cors.Config {
	allowedOrigins := s.parseCorsOrigins()
	allowCredentials := true

	if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
		allowCredentials = false
	}

	return cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowCredentials,
		MaxAge:           3600,
	}
}

func (s *Server) parseCorsOrigins() []string {
	if s.config.CorsAllowedOrigins == "" {
		return []string{"*"}
	}

	origins := strings.Split(s.config.CorsAllowedOrigins, ",")
	result := make([]string, 0, len(origins))

	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return []string{"*"}
	}

	return result
}

func defaultErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	detail := "An internal server error occurred"

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		if fiberErr.Message != "" {
			detail = fiberErr.Message
		}
	}

	title := GetTitleForStatus(code)
	if detail == "" {
		detail = title
	}

	problem := NewProblemDetails(
		code,
		title,
		detail,
		GetInstanceFromPath(c.Path()),
	)

	c.Set(fiber.HeaderContentType, ContentTypeProblemJSON)
	return c.Status(code).JSON(problem)
}
