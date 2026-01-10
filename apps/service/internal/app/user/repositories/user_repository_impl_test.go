package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func TestUserRepository_All(t *testing.T) {

	t.Run("returns empty slice when no users exist", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		users := repo.All(ctx)
		assert.Empty(t, users)
	})

	t.Run("returns all users", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user1 := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}
		user2 := entities.User{
			FullName: "Jane Smith",
			Phone:    "0987654321",
			Email:    "jane@example.com",
			Password: "password456",
		}

		_, err := repo.Create(ctx, user1)
		require.NoError(t, err)
		_, err = repo.Create(ctx, user2)
		require.NoError(t, err)

		users := repo.All(ctx)
		assert.Len(t, users, 2)
	})

	t.Run("excludes soft deleted users", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "Deleted User",
			Phone:    "1111111111",
			Email:    "deleted@example.com",
			Password: "password",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		err = repo.Delete(ctx, created)
		require.NoError(t, err)

		users := repo.All(ctx)
		assert.Empty(t, users)
	})
}

func TestUserRepository_FindByPhone(t *testing.T) {
	t.Run("returns user when found", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		ctx := context.Background()
		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		found, err := repo.FindByPhone(ctx, "1234567890")
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "John Doe", found.FullName)
		assert.Equal(t, "1234567890", found.Phone)
		assert.Equal(t, "john@example.com", found.Email)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		_, err := repo.FindByPhone(ctx, "9999999999")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("does not find soft deleted users", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "Deleted User",
			Phone:    "1111111111",
			Email:    "deleted@example.com",
			Password: "password",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		err = repo.Delete(ctx, created)
		require.NoError(t, err)

		_, err = repo.FindByPhone(ctx, "1111111111")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	t.Run("returns user when found", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		found, err := repo.FindByEmail(ctx, "john@example.com")
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "John Doe", found.FullName)
		assert.Equal(t, "john@example.com", found.Email)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		_, err := repo.FindByEmail(ctx, "notfound@example.com")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("does not find soft deleted users", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "Deleted User",
			Phone:    "1111111111",
			Email:    "deleted@example.com",
			Password: "password",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		err = repo.Delete(ctx, created)
		require.NoError(t, err)

		_, err = repo.FindByEmail(ctx, "deleted@example.com")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestUserRepository_Create(t *testing.T) {
	t.Run("creates user successfully", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, "John Doe", created.FullName)
		assert.Equal(t, "1234567890", created.Phone)
		assert.Equal(t, "john@example.com", created.Email)
		assert.NotZero(t, created.CreatedAt)
		assert.NotZero(t, created.UpdatedAt)
	})

	t.Run("returns error on duplicate phone", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user1 := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		_, err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := entities.User{
			FullName: "Jane Doe",
			Phone:    "1234567890",
			Email:    "jane@example.com",
			Password: "password456",
		}

		_, err = repo.Create(ctx, user2)
		assert.Error(t, err)
	})

	t.Run("returns error on duplicate email", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user1 := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		_, err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := entities.User{
			FullName: "Jane Doe",
			Phone:    "0987654321",
			Email:    "john@example.com",
			Password: "password456",
		}

		_, err = repo.Create(ctx, user2)
		assert.Error(t, err)
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("updates user successfully", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		originalUpdatedAt := created.UpdatedAt
		time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt changes

		created.FullName = "John Updated"
		created.Email = "john.updated@example.com"

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "John Updated", updated.FullName)
		assert.Equal(t, "john.updated@example.com", updated.Email)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("updates non-zero fields", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		created.Password = "newpassword123"
		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "newpassword123", updated.Password)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	t.Run("soft deletes user successfully", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user := entities.User{
			FullName: "John Doe",
			Phone:    "1234567890",
			Email:    "john@example.com",
			Password: "password123",
		}

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		err = repo.Delete(ctx, created)
		require.NoError(t, err)

		// Verify user is soft deleted
		_, err = repo.FindByPhone(ctx, "1234567890")
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())

		// Verify user still exists in database (soft deleted)
		var deletedUser entities.User
		err = db.Unscoped().Where("id = ?", created.ID).First(&deletedUser).Error
		require.NoError(t, err)
		assert.NotZero(t, deletedUser.DeletedAt)
	})

	t.Run("can delete multiple users", func(t *testing.T) {
		ctx := context.Background()
		db := testutil.SetupTestDB(t, &entities.User{})
		repo := NewUserRepository(db)

		user1 := entities.User{
			FullName: "User One",
			Phone:    "1111111111",
			Email:    "user1@example.com",
			Password: "password1",
		}
		user2 := entities.User{
			FullName: "User Two",
			Phone:    "2222222222",
			Email:    "user2@example.com",
			Password: "password2",
		}

		created1, err := repo.Create(ctx, user1)
		require.NoError(t, err)
		created2, err := repo.Create(ctx, user2)
		require.NoError(t, err)

		err = repo.Delete(ctx, created1)
		require.NoError(t, err)
		err = repo.Delete(ctx, created2)
		require.NoError(t, err)

		users := repo.All(ctx)
		assert.Empty(t, users)
	})
}
