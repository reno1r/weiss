package usecases

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/reno1r/weiss/apps/service/internal/app/auth/services"
	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	validationutil "github.com/reno1r/weiss/apps/service/internal/validation_util"
)

type RegisterUsecase struct {
	userRepository  repositories.UserRepository
	passwordService *services.PasswordService
	validator       *validator.Validate
}

func NewRegisterUsecase(userRepository repositories.UserRepository, passwordService *services.PasswordService) *RegisterUsecase {
	return &RegisterUsecase{
		userRepository:  userRepository,
		passwordService: passwordService,
		validator:       validator.New(),
	}
}

type RegisterRequest struct {
	FullName string `validate:"required,min=2,max=255"`
	Phone    string `validate:"required,min=10,max=20"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=100"`
}

type RegisterResponse struct {
	User entities.User
}

func (u *RegisterUsecase) Execute(req RegisterRequest) (RegisterResponse, error) {
	if err := u.validator.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validationutil.GetValidationErrorMessage(err))
		}
		return RegisterResponse{}, fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	_, err := u.userRepository.FindByEmail(req.Email)
	if err == nil {
		return RegisterResponse{}, errors.New("user with this email already exists")
	}

	_, err = u.userRepository.FindByPhone(req.Phone)
	if err == nil {
		return RegisterResponse{}, errors.New("user with this phone already exists")
	}

	hashedPassword, err := u.passwordService.HashPassword(req.Password)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user := entities.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: hashedPassword,
	}

	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	return RegisterResponse{
		User: createdUser,
	}, nil
}
