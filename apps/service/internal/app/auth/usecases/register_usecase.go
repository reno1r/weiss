package usecases

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/reno1r/weiss/apps/service/internal/app/auth/services"
	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	"github.com/reno1r/weiss/apps/service/internal/validationutil"
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

type RegisterParam struct {
	FullName string `validate:"required,min=2,max=255"`
	Phone    string `validate:"required,min=10,max=20"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=100"`
}

type RegisterResult struct {
	User *entities.User
}

func (u *RegisterUsecase) Execute(param RegisterParam) (*RegisterResult, error) {
	if err := u.validator.Struct(param); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validationutil.GetValidationErrorMessage(err))
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	_, err := u.userRepository.FindByEmail(param.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	_, err = u.userRepository.FindByPhone(param.Phone)
	if err == nil {
		return nil, errors.New("user with this phone already exists")
	}

	hashedPassword, err := u.passwordService.HashPassword(param.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := entities.User{
		FullName: param.FullName,
		Phone:    param.Phone,
		Email:    param.Email,
		Password: hashedPassword,
	}

	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &RegisterResult{
		User: &createdUser,
	}, nil
}
