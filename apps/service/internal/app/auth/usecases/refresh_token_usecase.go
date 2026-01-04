package usecases

import (
	"errors"
	"fmt"

	"github.com/reno1r/weiss/apps/service/internal/app/auth/services"
	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
)

type RefreshTokenUsecase struct {
	userRepository repositories.UserRepository
	tokenService   *services.TokenService
}

func NewRefreshTokenUsecase(userRepository repositories.UserRepository, tokenService *services.TokenService) *RefreshTokenUsecase {
	return &RefreshTokenUsecase{
		userRepository: userRepository,
		tokenService:   tokenService,
	}
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
	User         entities.User
}

func (u *RefreshTokenUsecase) Execute(req RefreshTokenRequest) (RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return RefreshTokenResponse{}, errors.New("refresh token is required")
	}

	claims, err := u.tokenService.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return RefreshTokenResponse{}, fmt.Errorf("failed to verify refresh token: %w", err)
	}

	userID, err := u.tokenService.GetUserID(claims)
	if err != nil {
		return RefreshTokenResponse{}, fmt.Errorf("failed to get user ID from token: %w", err)
	}

	tokenPair, err := u.tokenService.GenerateTokenPair(userID, claims.Email, claims.Phone)
	if err != nil {
		return RefreshTokenResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}

	var user entities.User
	if claims.Email != "" {
		user, err = u.userRepository.FindByEmail(claims.Email)
	} else if claims.Phone != "" {
		user, err = u.userRepository.FindByPhone(claims.Phone)
	} else {
		return RefreshTokenResponse{}, errors.New("invalid token claims")
	}

	if err != nil {
		return RefreshTokenResponse{}, errors.New("user not found")
	}

	return RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         user,
	}, nil
}
