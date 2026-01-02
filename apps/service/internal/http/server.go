package http

import "github.com/gofiber/fiber/v3"

func NewServer() *fiber.App {
	app := fiber.New()

	return app
}
