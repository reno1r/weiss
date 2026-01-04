package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/reno1r/weiss/apps/service/internal/config"
)

func TestNewPasswordService(t *testing.T) {
	t.Run("uses configured cost", func(t *testing.T) {
		config := &config.Config{
			BcryptCost: 12,
		}

		service := NewPasswordService(config)
		assert.Equal(t, 12, service.cost)
	})

	t.Run("uses default cost when config is 0", func(t *testing.T) {
		config := &config.Config{
			BcryptCost: 0,
		}

		service := NewPasswordService(config)
		assert.Equal(t, bcrypt.DefaultCost, service.cost)
	})

	t.Run("uses provided cost even when negative", func(t *testing.T) {
		config := &config.Config{
			BcryptCost: -1,
		}

		service := NewPasswordService(config)
		assert.Equal(t, -1, service.cost)
	})
}

func TestPasswordService_HashPassword(t *testing.T) {
	config := &config.Config{
		BcryptCost: bcrypt.MinCost,
	}
	service := NewPasswordService(config)

	t.Run("hashes password successfully", func(t *testing.T) {
		password := "testpassword123"
		hashed, err := service.HashPassword(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)
	})

	t.Run("generates different hashes for same password", func(t *testing.T) {
		password := "testpassword123"
		hashed1, err := service.HashPassword(password)
		require.NoError(t, err)

		hashed2, err := service.HashPassword(password)
		require.NoError(t, err)

		assert.NotEqual(t, hashed1, hashed2)
	})

	t.Run("hashed password can be verified", func(t *testing.T) {
		password := "testpassword123"
		hashed, err := service.HashPassword(password)
		require.NoError(t, err)

		isValid := service.VerifyPassword(hashed, password)
		assert.True(t, isValid)
	})
}

func TestPasswordService_VerifyPassword(t *testing.T) {
	config := &config.Config{
		BcryptCost: bcrypt.MinCost,
	}
	service := NewPasswordService(config)

	t.Run("verifies correct password", func(t *testing.T) {
		password := "testpassword123"
		hashed, err := service.HashPassword(password)
		require.NoError(t, err)

		isValid := service.VerifyPassword(hashed, password)
		assert.True(t, isValid)
	})

	t.Run("rejects incorrect password", func(t *testing.T) {
		password := "testpassword123"
		wrongPassword := "wrongpassword"
		hashed, err := service.HashPassword(password)
		require.NoError(t, err)

		isValid := service.VerifyPassword(hashed, wrongPassword)
		assert.False(t, isValid)
	})

	t.Run("rejects empty password", func(t *testing.T) {
		password := "testpassword123"
		hashed, err := service.HashPassword(password)
		require.NoError(t, err)

		isValid := service.VerifyPassword(hashed, "")
		assert.False(t, isValid)
	})

	t.Run("rejects invalid hash", func(t *testing.T) {
		invalidHash := "invalidhash"
		isValid := service.VerifyPassword(invalidHash, "password")
		assert.False(t, isValid)
	})

	t.Run("handles different password lengths", func(t *testing.T) {
		testCases := []string{
			"short",
			"mediumlengthpassword",
			"verylongpasswordthatexceedsnormallimitsandshouldstillworkcorrectly",
		}

		for _, password := range testCases {
			hashed, err := service.HashPassword(password)
			require.NoError(t, err)

			isValid := service.VerifyPassword(hashed, password)
			assert.True(t, isValid, "password: %s", password)
		}
	})
}
