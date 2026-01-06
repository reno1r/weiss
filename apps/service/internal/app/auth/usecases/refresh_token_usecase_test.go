package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/auth/services"
	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	"github.com/reno1r/weiss/apps/service/internal/config"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupRefreshTokenTest(t *testing.T) (*RefreshTokenUsecase, repositories.UserRepository, *services.TokenService) {
	db := testutil.SetupTestDB(t, &entities.User{})
	userRepo := repositories.NewUserRepository(db)

	config := &config.Config{
		JwtSecret:         "test-secret-key-minimum-32-characters-long",
		JwtIssuer:         "test-issuer",
		JwtAccessExpires:  "15m",
		JwtRefreshExpires: "168h",
	}

	tokenService, err := services.NewTokenService(config)
	require.NoError(t, err)

	refreshUsecase := NewRefreshTokenUsecase(userRepo, tokenService)

	return refreshUsecase, userRepo, tokenService
}

func TestRefreshTokenUsecase_Execute(t *testing.T) {
	t.Run("refreshes token successfully with email", func(t *testing.T) {
		usecase, userRepo, tokenService := setupRefreshTokenTest(t)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "hashedpassword",
		}
		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)

		refreshToken, err := tokenService.GenerateRefreshToken(createdUser.ID, createdUser.Email, createdUser.Phone)
		require.NoError(t, err)

		req := RefreshTokenData{
			RefreshToken: refreshToken,
		}

		resp, err := usecase.Execute(req)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, resp.User.ID)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotEqual(t, refreshToken, resp.RefreshToken)
	})

	t.Run("refreshes token successfully with phone", func(t *testing.T) {
		usecase, userRepo, tokenService := setupRefreshTokenTest(t)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "",
			Password: "hashedpassword",
		}
		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)

		refreshToken, err := tokenService.GenerateRefreshToken(createdUser.ID, createdUser.Email, createdUser.Phone)
		require.NoError(t, err)

		req := RefreshTokenData{
			RefreshToken: refreshToken,
		}

		resp, err := usecase.Execute(req)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, resp.User.ID)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	t.Run("returns error when refresh token is empty", func(t *testing.T) {
		usecase, _, _ := setupRefreshTokenTest(t)

		req := RefreshTokenData{
			RefreshToken: "",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "refresh token is required")
		assert.Nil(t, resp)
	})

	t.Run("returns error when refresh token is invalid", func(t *testing.T) {
		usecase, _, _ := setupRefreshTokenTest(t)

		req := RefreshTokenData{
			RefreshToken: "invalid.token.here",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to verify refresh token")
		assert.Nil(t, resp)
	})

	t.Run("returns error when access token is used instead", func(t *testing.T) {
		usecase, userRepo, tokenService := setupRefreshTokenTest(t)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "hashedpassword",
		}
		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)

		accessToken, err := tokenService.GenerateAccessToken(createdUser.ID, createdUser.Email, createdUser.Phone)
		require.NoError(t, err)

		req := RefreshTokenData{
			RefreshToken: accessToken,
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to verify refresh token")
		assert.Nil(t, resp)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		usecase, _, tokenService := setupRefreshTokenTest(t)

		refreshToken, err := tokenService.GenerateRefreshToken(999, "notfound@example.com", "9999999999")
		require.NoError(t, err)

		req := RefreshTokenData{
			RefreshToken: refreshToken,
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		assert.Nil(t, resp)
	})

	t.Run("generates new token pair on refresh", func(t *testing.T) {
		usecase, userRepo, tokenService := setupRefreshTokenTest(t)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "hashedpassword",
		}
		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)

		originalRefreshToken, err := tokenService.GenerateRefreshToken(createdUser.ID, createdUser.Email, createdUser.Phone)
		require.NoError(t, err)

		req := RefreshTokenData{
			RefreshToken: originalRefreshToken,
		}

		resp, err := usecase.Execute(req)
		require.NoError(t, err)

		accessClaims, err := tokenService.VerifyToken(resp.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, "access", accessClaims.Type)

		refreshClaims, err := tokenService.VerifyRefreshToken(resp.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, "refresh", refreshClaims.Type)

		userID, err := tokenService.GetUserID(accessClaims)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, userID)
	})

	t.Run("returns error when token claims have no email or phone", func(t *testing.T) {
		usecase, _, tokenService := setupRefreshTokenTest(t)

		refreshToken, err := tokenService.GenerateRefreshToken(1, "", "")
		require.NoError(t, err)

		req := RefreshTokenData{
			RefreshToken: refreshToken,
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token claims")
		assert.Nil(t, resp)
	})
}
