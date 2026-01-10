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
	Email    string `json:"email" example:"john@example.com"`                  // User's email address (optional if phone is provided)
	Phone    string `json:"phone" example:"1234567890"`                        // User's phone number (optional if email is provided)
	Password string `json:"password" example:"password123" binding:"required"` // User's password
}

type LoginResponse struct {
	Message string            `json:"message"`
	Data    LoginResponseData `json:"data"`
}

type LoginResponseData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user with email/phone and password, returns JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginPayload  true  "Login credentials"
// @Success      200         {object}  LoginResponse
// @Failure      400         {object}  map[string]string  "Invalid request body"
// @Failure      401         {object}  map[string]string  "Invalid credentials"
// @Failure      422         {object}  map[string]string  "Validation failed"
// @Failure      500         {object}  map[string]string  "Internal server error"
// @Router       /auth/login [post]
func (h *LoginHandler) Handle(c fiber.Ctx) error {
	var credentials LoginPayload

	if err := c.Bind().Body(&credentials); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	ctx := c.Context()
	response, err := h.loginUsecase.Execute(ctx, usecases.LoginParam{
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
