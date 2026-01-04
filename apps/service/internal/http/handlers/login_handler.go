package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/reno1r/weiss/apps/service/internal/app/auth/usecases"
)

type LoginHandler struct {
	loginUsecase *usecases.LoginUsecase
}

func NewLoginHandler(loginUsecase *usecases.LoginUsecase) *LoginHandler {
	return &LoginHandler{
		loginUsecase: loginUsecase,
	}
}

func (h *LoginHandler) Handle(c fiber.Ctx) error {
	var credentials usecases.LoginCredential

	if err := c.Bind().Body(&credentials); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := h.loginUsecase.Execute(credentials)
	if err != nil {
		if isValidationError(err) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		if isInvalidCredentialsError(err) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to process login request")
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":        response.User.ID,
			"full_name": response.User.FullName,
			"phone":     response.User.Phone,
			"email":     response.User.Email,
		},
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
	})
}

func isValidationError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "validation failed")
}

func isInvalidCredentialsError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "invalid credentials")
}
