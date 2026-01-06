package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/reno1r/weiss/apps/service/internal/app/auth/services"
	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	"github.com/reno1r/weiss/apps/service/internal/config"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupLoginTest(t *testing.T) (*LoginUsecase, repositories.UserRepository) {
	db := testutil.SetupTestDB(t, &entities.User{})
	userRepo := repositories.NewUserRepository(db)

	config := &config.Config{
		BcryptCost:        bcrypt.MinCost,
		JwtSecret:         "test-secret-key-minimum-32-characters-long",
		JwtIssuer:         "test-issuer",
		JwtAccessExpires:  "15m",
		JwtRefreshExpires: "168h",
	}

	passwordService := services.NewPasswordService(config)
	tokenService, err := services.NewTokenService(config)
	require.NoError(t, err)

	loginUsecase := NewLoginUsecase(userRepo, tokenService, passwordService)

	return loginUsecase, userRepo
}

func TestLoginUsecase_Execute(t *testing.T) {
	t.Run("logs in successfully with email", func(t *testing.T) {
		usecase, userRepo := setupLoginTest(t)

		password := "password123"
		hashedPassword, err := services.NewPasswordService(&config.Config{BcryptCost: bcrypt.MinCost}).HashPassword(password)
		require.NoError(t, err)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: hashedPassword,
		}
		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)

		credentials := LoginData{
			Email:    "john@example.com",
			Password: password,
		}

		resp, err := usecase.Execute(credentials)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, resp.User.ID)
		assert.Equal(t, "John Doe", resp.User.FullName)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotEqual(t, resp.AccessToken, resp.RefreshToken)
	})

	t.Run("logs in successfully with phone", func(t *testing.T) {
		usecase, userRepo := setupLoginTest(t)

		password := "password123"
		hashedPassword, err := services.NewPasswordService(&config.Config{BcryptCost: bcrypt.MinCost}).HashPassword(password)
		require.NoError(t, err)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: hashedPassword,
		}
		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)

		credentials := LoginData{
			Phone:    "1234567890",
			Password: password,
		}

		resp, err := usecase.Execute(credentials)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, resp.User.ID)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	t.Run("returns error when email not found", func(t *testing.T) {
		usecase, _ := setupLoginTest(t)

		credentials := LoginData{
			Email:    "notfound@example.com",
			Password: "password123",
		}

		resp, err := usecase.Execute(credentials)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		assert.Nil(t, resp)
	})

	t.Run("returns error when phone not found", func(t *testing.T) {
		usecase, _ := setupLoginTest(t)

		credentials := LoginData{
			Phone:    "9999999999",
			Password: "password123",
		}

		resp, err := usecase.Execute(credentials)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		assert.Nil(t, resp)
	})

	t.Run("returns error when password is incorrect", func(t *testing.T) {
		usecase, userRepo := setupLoginTest(t)

		password := "password123"
		hashedPassword, err := services.NewPasswordService(&config.Config{BcryptCost: bcrypt.MinCost}).HashPassword(password)
		require.NoError(t, err)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: hashedPassword,
		}
		_, err = userRepo.Create(user)
		require.NoError(t, err)

		credentials := LoginData{
			Email:    "john@example.com",
			Password: "wrongpassword",
		}

		resp, err := usecase.Execute(credentials)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		assert.Nil(t, resp)
	})

	t.Run("validates password is required", func(t *testing.T) {
		usecase, _ := setupLoginTest(t)

		credentials := LoginData{
			Email:    "john@example.com",
			Password: "",
		}

		resp, err := usecase.Execute(credentials)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates email format when provided", func(t *testing.T) {
		usecase, _ := setupLoginTest(t)

		credentials := LoginData{
			Email:    "invalid-email",
			Password: "password123",
		}

		resp, err := usecase.Execute(credentials)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates phone length when provided", func(t *testing.T) {
		usecase, _ := setupLoginTest(t)

		credentials := LoginData{
			Phone:    "123",
			Password: "password123",
		}

		resp, err := usecase.Execute(credentials)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("returns error when both email and phone are empty", func(t *testing.T) {
		usecase, _ := setupLoginTest(t)

		credentials := LoginData{
			Email:    "",
			Phone:    "",
			Password: "password123",
		}

		resp, err := usecase.Execute(credentials)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email or phone is required")
		assert.Nil(t, resp)
	})

	t.Run("generates valid tokens", func(t *testing.T) {
		usecase, userRepo := setupLoginTest(t)

		password := "password123"
		hashedPassword, err := services.NewPasswordService(&config.Config{BcryptCost: bcrypt.MinCost}).HashPassword(password)
		require.NoError(t, err)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: hashedPassword,
		}
		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)

		credentials := LoginData{
			Email:    "john@example.com",
			Password: password,
		}

		resp, err := usecase.Execute(credentials)
		require.NoError(t, err)

		config := &config.Config{
			JwtSecret:         "test-secret-key-minimum-32-characters-long",
			JwtIssuer:         "test-issuer",
			JwtAccessExpires:  "15m",
			JwtRefreshExpires: "168h",
		}
		tokenService, err := services.NewTokenService(config)
		require.NoError(t, err)

		accessClaims, err := tokenService.VerifyToken(resp.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, "access", accessClaims.Type)
		userID, err := tokenService.GetUserID(accessClaims)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, userID)

		refreshClaims, err := tokenService.VerifyRefreshToken(resp.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, "refresh", refreshClaims.Type)
	})
}
