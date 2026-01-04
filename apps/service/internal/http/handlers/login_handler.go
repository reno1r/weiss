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

type LoginPayload struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string            `json:"message"`
	Data    LoginResponseData `json:"data"`
}

type LoginResponseData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *LoginHandler) Handle(c fiber.Ctx) error {
	var credentials LoginPayload

	if err := c.Bind().Body(&credentials); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := h.loginUsecase.Execute(usecases.LoginData{
		Email:    credentials.Email,
		Phone:    credentials.Phone,
		Password: credentials.Password,
	})
	if err != nil {
		if isValidationError(err) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		if isInvalidCredentialsError(err) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to process login request")
	}

	return c.JSON(LoginResponse{
		Message: "authorized.",
		Data: LoginResponseData{
			AccessToken:  response.AccessToken,
			RefreshToken: response.RefreshToken,
		},
	})
}

func isValidationError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "validation failed")
}

func isInvalidCredentialsError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "invalid credentials")
}
