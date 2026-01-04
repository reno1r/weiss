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

type LoginUsecase struct {
	userRepository  repositories.UserRepository
	tokenService    *services.TokenService
	passwordService *services.PasswordService
	validator       *validator.Validate
}

func NewLoginUsecase(userRepository repositories.UserRepository, tokenService *services.TokenService, passwordService *services.PasswordService) *LoginUsecase {
	return &LoginUsecase{
		userRepository:  userRepository,
		tokenService:    tokenService,
		passwordService: passwordService,
		validator:       validator.New(),
	}
}

type LoginData struct {
	Email    string `validate:"omitempty,email"`
	Phone    string `validate:"omitempty,min=10,max=20"`
	Password string `validate:"required,min=1"`
}

type LoginResult struct {
	User         *entities.User
	AccessToken  string
	RefreshToken string
}

func (u *LoginUsecase) Execute(credentials LoginData) (*LoginResult, error) {
	if err := u.validator.Struct(credentials); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validationutil.GetValidationErrorMessage(err))
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	if credentials.Email == "" && credentials.Phone == "" {
		return nil, errors.New("email or phone is required")
	}

	var user entities.User
	var err error

	if credentials.Email != "" {
		user, err = u.userRepository.FindByEmail(credentials.Email)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	} else if credentials.Phone != "" {
		user, err = u.userRepository.FindByPhone(credentials.Phone)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	} else {
		return nil, errors.New("email or phone is required")
	}

	if !u.passwordService.VerifyPassword(user.Password, credentials.Password) {
		return nil, errors.New("invalid credentials")
	}

	tokenPair, err := u.tokenService.GenerateTokenPair(user.ID, user.Email, user.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &LoginResult{
		User:         &user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
