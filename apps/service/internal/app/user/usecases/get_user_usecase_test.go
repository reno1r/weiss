package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupGetUserTest(t *testing.T) (*GetUserUsecase, repositories.UserRepository) {
	db := testutil.SetupTestDB(t, &entities.User{})
	userRepo := repositories.NewUserRepository(db)
	usecase := NewGetUserUsecase(userRepo)
	return usecase, userRepo
}

func TestGetUserUsecase_Execute(t *testing.T) {
	t.Run("returns user when found", func(t *testing.T) {
		usecase, userRepo := setupGetUserTest(t)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "hashedpassword",
		}

		created, err := userRepo.Create(user)
		require.NoError(t, err)

		result, err := usecase.Execute(created.ID)
		require.NoError(t, err)
		assert.NotNil(t, result.User)
		assert.Equal(t, created.ID, result.User.ID)
		assert.Equal(t, "John Doe", result.User.FullName)
		assert.Equal(t, "1234567890", result.User.Phone)
		assert.Equal(t, "john@example.com", result.User.Email)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		usecase, _ := setupGetUserTest(t)

		result, err := usecase.Execute(999)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, result)
	})

	t.Run("does not find soft deleted users", func(t *testing.T) {
		usecase, userRepo := setupGetUserTest(t)

		user := entities.User{
			FullName: "Deleted User",
			Phone:    "1111111111",
			Email:    "deleted@example.com",
			Password: "hashedpassword",
		}

		created, err := userRepo.Create(user)
		require.NoError(t, err)

		err = userRepo.Delete(created)
		require.NoError(t, err)

		result, err := usecase.Execute(created.ID)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, result)
	})
}

