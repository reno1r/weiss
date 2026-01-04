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
	testutil "github.com/reno1r/weiss/apps/service/internal/test_util"
)

func setupRegisterTest(t *testing.T) (*RegisterUsecase, repositories.UserRepository) {
	db := testutil.SetupTestDB(t, &entities.User{})
	userRepo := repositories.NewUserRepository(db)

	config := &config.Config{
		BcryptCost: bcrypt.MinCost,
	}
	passwordService := services.NewPasswordService(config)

	registerUsecase := NewRegisterUsecase(&userRepo, passwordService)

	return registerUsecase, userRepo
}

func TestRegisterUsecase_Execute(t *testing.T) {
	t.Run("registers user successfully", func(t *testing.T) {
		usecase, userRepo := setupRegisterTest(t)

		req := RegisterData{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		require.NoError(t, err)
		assert.NotZero(t, resp.ID)
		assert.Equal(t, "John Doe", resp.FullName)
		assert.Equal(t, "1234567890", resp.Phone)
		assert.Equal(t, "john@example.com", resp.Email)
		assert.NotEqual(t, "password123", resp.Password)
		assert.NotEmpty(t, resp.Password)

		verifyUser, err := userRepo.FindByEmail("john@example.com")
		require.NoError(t, err)
		assert.Equal(t, resp.ID, verifyUser.ID)
	})

	t.Run("returns error when email already exists", func(t *testing.T) {
		usecase, userRepo := setupRegisterTest(t)

		existingUser := entities.User{
			FullName: "Existing User",
			Phone:    "1111111111",
			Email:    "existing@example.com",
			Password: "hashedpassword",
		}
		_, err := userRepo.Create(existingUser)
		require.NoError(t, err)

		req := RegisterData{
			FullName: "New User",
			Phone:    "2222222222",
			Email:    "existing@example.com",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user with this email already exists")
		assert.Nil(t, resp)
	})

	t.Run("returns error when phone already exists", func(t *testing.T) {
		usecase, userRepo := setupRegisterTest(t)

		existingUser := entities.User{
			FullName: "Existing User",
			Phone:    "1234567890",
			Email:    "existing@example.com",
			Password: "hashedpassword",
		}
		_, err := userRepo.Create(existingUser)
		require.NoError(t, err)

		req := RegisterData{
			FullName: "New User",
			Phone:    "1234567890",
			Email:    "new@example.com",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user with this phone already exists")
		assert.Nil(t, resp)
	})

	t.Run("validates full name is required", func(t *testing.T) {
		usecase, _ := setupRegisterTest(t)

		req := RegisterData{
			FullName: "",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates full name minimum length", func(t *testing.T) {
		usecase, _ := setupRegisterTest(t)

		req := RegisterData{
			FullName: "A",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates phone is required", func(t *testing.T) {
		usecase, _ := setupRegisterTest(t)

		req := RegisterData{
			FullName: "John Doe",
			Phone:    "",
			Email:    "john@example.com",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates phone minimum length", func(t *testing.T) {
		usecase, _ := setupRegisterTest(t)

		req := RegisterData{
			FullName: "John Doe",
			Phone:    "123",
			Email:    "john@example.com",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates email is required", func(t *testing.T) {
		usecase, _ := setupRegisterTest(t)

		req := RegisterData{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates email format", func(t *testing.T) {
		usecase, _ := setupRegisterTest(t)

		req := RegisterData{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "invalid-email",
			Password: "password123",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates password is required", func(t *testing.T) {
		usecase, _ := setupRegisterTest(t)

		req := RegisterData{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("validates password minimum length", func(t *testing.T) {
		usecase, _ := setupRegisterTest(t)

		req := RegisterData{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "12345",
		}

		resp, err := usecase.Execute(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, resp)
	})

	t.Run("hashes password before storing", func(t *testing.T) {
		usecase, userRepo := setupRegisterTest(t)

		req := RegisterData{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		_, err := usecase.Execute(req)
		require.NoError(t, err)

		user, err := userRepo.FindByEmail("john@example.com")
		require.NoError(t, err)
		assert.NotEqual(t, "password123", user.Password)
		assert.True(t, len(user.Password) > 50)
	})
}
