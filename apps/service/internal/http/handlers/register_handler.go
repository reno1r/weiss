package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/reno1r/weiss/apps/service/internal/app/auth/usecases"
	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
)

type RegisterHandler struct {
	registerUsecase *usecases.RegisterUsecase
}

func NewRegisterHandler(registerUsecase *usecases.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{
		registerUsecase: registerUsecase,
	}
}

type RegisterPayload struct {
	FulllName string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type RegisterResponse struct {
	Data struct {
		User *entities.User `json:"user"`
	} `json:"data"`
}

func (h *RegisterHandler) Handle(c fiber.Ctx) error {
	var request RegisterPayload

	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result, err := h.registerUsecase.Execute(usecases.RegisterData{
		FullName: request.FulllName,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: request.Password,
	})

	if err != nil {
		if isValidationError(err) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		if isConflictError(err) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to process registration request")
	}

	result.User.Password = ""

	return c.Status(fiber.StatusCreated).JSON(RegisterResponse{
		Data: struct {
			User *entities.User `json:"user"`
		}{
			User: result.User,
		},
	})
}

func isConflictError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already exists")
}
