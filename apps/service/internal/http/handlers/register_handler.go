package handlers

import (
	"strings"
	"time"

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

type RegisterPayload struct {
	FulllName string `json:"full_name" example:"John Doe" binding:"required"`           // User's full name
	Email     string `json:"email" example:"john@example.com" binding:"required,email"` // User's email address
	Phone     string `json:"phone" example:"1234567890" binding:"required"`             // User's phone number
	Password  string `json:"password" example:"password123" binding:"required,min=6"`   // User's password (minimum 6 characters)
}

type RegisterResponse struct {
	Message string               `json:"message"`
	Data    RegisterResponseData `json:"data"`
}

type RegisterResponseData struct {
	User *UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uint64    `json:"id" example:"1"`
	FullName  string    `json:"full_name" example:"John Doe"`
	Phone     string    `json:"phone" example:"1234567890"`
	Email     string    `json:"email" example:"john@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with email, phone, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterPayload  true  "Registration data"
// @Success      201      {object}  RegisterResponse
// @Failure      400      {object}  map[string]string  "Invalid request body"
// @Failure      422      {object}  map[string]string  "Validation failed"
// @Failure      409      {object}  map[string]string  "User already exists"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /auth/register [post]
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

	userResponse := &UserResponse{
		ID:        result.User.ID,
		FullName:  result.User.FullName,
		Phone:     result.User.Phone,
		Email:     result.User.Email,
		CreatedAt: result.User.CreatedAt,
		UpdatedAt: result.User.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(RegisterResponse{
		Message: "user created successfully.",
		Data: RegisterResponseData{
			User: userResponse,
		},
	})
}

func isConflictError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already exists")
}
