package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/reno1r/weiss/apps/service/internal/app/auth/usecases"
)

type RegisterHandler struct {
	registerUsecase *usecases.RegisterUsecase
}

func NewRegisterHandler(registerUsecase *usecases.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		registerUsecase: registerUsecase,
	}
}

func (h *RegisterHandler) Handle(c fiber.Ctx) error {
	var request usecases.RegisterRequest

	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := h.registerUsecase.Execute(request)
	if err != nil {
		if isValidationError(err) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		if isConflictError(err) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to process registration request")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": fiber.Map{
			"id":        response.User.ID,
			"full_name": response.User.FullName,
			"phone":     response.User.Phone,
			"email":     response.User.Email,
		},
	})
}

func isConflictError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already exists")
}
