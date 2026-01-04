package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/config"
)

func setupTestConfig() *config.Config {
	return &config.Config{
		JwtSecret:         "test-secret-key-minimum-32-characters-long",
		JwtIssuer:         "test-issuer",
		JwtAccessExpires:  "15m",
		JwtRefreshExpires: "168h",
	}
}

func TestNewTokenService(t *testing.T) {
	t.Run("creates service with valid config", func(t *testing.T) {
		config := setupTestConfig()
		service, err := NewTokenService(config)
		require.NoError(t, err)
		assert.NotNil(t, service)
		assert.Equal(t, "test-issuer", service.issuer)
		assert.Equal(t, 15*time.Minute, service.accessExpiresIn)
		assert.Equal(t, 7*24*time.Hour, service.refreshExpiresIn)
	})

	t.Run("returns error when secret is empty", func(t *testing.T) {
		config := setupTestConfig()
		config.JwtSecret = ""
		service, err := NewTokenService(config)
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "JWT secret is required")
	})

	t.Run("returns error when issuer is empty", func(t *testing.T) {
		config := setupTestConfig()
		config.JwtIssuer = ""
		service, err := NewTokenService(config)
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "JWT issuer is required")
	})

	t.Run("returns error when access expiration is invalid", func(t *testing.T) {
		config := setupTestConfig()
		config.JwtAccessExpires = "invalid-duration"
		service, err := NewTokenService(config)
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "invalid JWT access expiration")
	})

	t.Run("returns error when refresh expiration is invalid", func(t *testing.T) {
		config := setupTestConfig()
		config.JwtRefreshExpires = "invalid-duration"
		service, err := NewTokenService(config)
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "invalid JWT refresh expiration")
	})
}

func TestTokenService_GenerateAccessToken(t *testing.T) {
	config := setupTestConfig()
	service, err := NewTokenService(config)
	require.NoError(t, err)

	t.Run("generates access token successfully", func(t *testing.T) {
		token, err := service.GenerateAccessToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generates different tokens for same user", func(t *testing.T) {
		token1, err := service.GenerateAccessToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		token2, err := service.GenerateAccessToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		assert.NotEqual(t, token1, token2)
	})

	t.Run("token contains correct claims", func(t *testing.T) {
		userID := uint64(123)
		email := "test@example.com"
		phone := "1234567890"

		token, err := service.GenerateAccessToken(userID, email, phone)
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)

		assert.Equal(t, "access", claims.Type)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, phone, claims.Phone)
		assert.Equal(t, "test-issuer", claims.Issuer)
		assert.Equal(t, "123", claims.Subject)
		assert.NotEmpty(t, claims.ID)
	})

	t.Run("token has correct expiration", func(t *testing.T) {
		token, err := service.GenerateAccessToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)

		now := time.Now()
		expectedExp := now.Add(15 * time.Minute)
		actualExp := claims.ExpiresAt.Time

		assert.WithinDuration(t, expectedExp, actualExp, 1*time.Second)
	})
}

func TestTokenService_GenerateRefreshToken(t *testing.T) {
	config := setupTestConfig()
	service, err := NewTokenService(config)
	require.NoError(t, err)

	t.Run("generates refresh token successfully", func(t *testing.T) {
		token, err := service.GenerateRefreshToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("token contains correct claims", func(t *testing.T) {
		userID := uint64(123)
		email := "test@example.com"
		phone := "1234567890"

		token, err := service.GenerateRefreshToken(userID, email, phone)
		require.NoError(t, err)

		claims, err := service.VerifyRefreshToken(token)
		require.NoError(t, err)

		assert.Equal(t, "refresh", claims.Type)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, phone, claims.Phone)
		assert.Equal(t, "test-issuer", claims.Issuer)
		assert.Equal(t, "123", claims.Subject)
		assert.NotEmpty(t, claims.ID)
	})

	t.Run("token has correct expiration", func(t *testing.T) {
		token, err := service.GenerateRefreshToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyRefreshToken(token)
		require.NoError(t, err)

		now := time.Now()
		expectedExp := now.Add(168 * time.Hour)
		actualExp := claims.ExpiresAt.Time

		assert.WithinDuration(t, expectedExp, actualExp, 1*time.Second)
	})
}

func TestTokenService_GenerateTokenPair(t *testing.T) {
	config := setupTestConfig()
	service, err := NewTokenService(config)
	require.NoError(t, err)

	t.Run("generates both tokens successfully", func(t *testing.T) {
		tokenPair, err := service.GenerateTokenPair(1, "test@example.com", "1234567890")
		require.NoError(t, err)
		assert.NotEmpty(t, tokenPair.AccessToken)
		assert.NotEmpty(t, tokenPair.RefreshToken)
		assert.NotEqual(t, tokenPair.AccessToken, tokenPair.RefreshToken)
	})

	t.Run("access token is valid", func(t *testing.T) {
		tokenPair, err := service.GenerateTokenPair(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyToken(tokenPair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, "access", claims.Type)
	})

	t.Run("refresh token is valid", func(t *testing.T) {
		tokenPair, err := service.GenerateTokenPair(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyRefreshToken(tokenPair.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, "refresh", claims.Type)
	})
}

func TestTokenService_VerifyToken(t *testing.T) {
	config := setupTestConfig()
	service, err := NewTokenService(config)
	require.NoError(t, err)

	t.Run("verifies valid token", func(t *testing.T) {
		token, err := service.GenerateAccessToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "1", claims.Subject)
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		claims, err := service.VerifyToken(invalidToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("rejects empty token", func(t *testing.T) {
		claims, err := service.VerifyToken("")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("rejects token with wrong secret", func(t *testing.T) {
		otherconfig := setupTestConfig()
		otherconfig.JwtSecret = "different-secret-key-minimum-32-characters-long"
		otherService, err := NewTokenService(otherconfig)
		require.NoError(t, err)

		token, err := service.GenerateAccessToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := otherService.VerifyToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestTokenService_VerifyRefreshToken(t *testing.T) {
	config := setupTestConfig()
	service, err := NewTokenService(config)
	require.NoError(t, err)

	t.Run("verifies valid refresh token", func(t *testing.T) {
		token, err := service.GenerateRefreshToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyRefreshToken(token)
		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "refresh", claims.Type)
	})

	t.Run("rejects access token", func(t *testing.T) {
		token, err := service.GenerateAccessToken(1, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyRefreshToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token type")
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		claims, err := service.VerifyRefreshToken(invalidToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestTokenService_GetUserID(t *testing.T) {
	config := setupTestConfig()
	service, err := NewTokenService(config)
	require.NoError(t, err)

	t.Run("extracts user ID from claims", func(t *testing.T) {
		token, err := service.GenerateAccessToken(123, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)

		userID, err := service.GetUserID(claims)
		require.NoError(t, err)
		assert.Equal(t, uint64(123), userID)
	})

	t.Run("returns error when subject is missing", func(t *testing.T) {
		claims := &Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "",
			},
		}

		userID, err := service.GetUserID(claims)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), userID)
		assert.Contains(t, err.Error(), "subject claim is missing")
	})

	t.Run("returns error when subject is invalid", func(t *testing.T) {
		claims := &Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "invalid",
			},
		}

		userID, err := service.GetUserID(claims)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), userID)
		assert.Contains(t, err.Error(), "invalid subject claim")
	})

	t.Run("handles large user IDs", func(t *testing.T) {
		largeID := uint64(999999999999)
		token, err := service.GenerateAccessToken(largeID, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)

		userID, err := service.GetUserID(claims)
		require.NoError(t, err)
		assert.Equal(t, largeID, userID)
	})
}

func TestTokenService_RefreshAccessToken(t *testing.T) {
	config := setupTestConfig()
	service, err := NewTokenService(config)
	require.NoError(t, err)

	t.Run("refreshes token successfully", func(t *testing.T) {
		refreshToken, err := service.GenerateRefreshToken(123, "test@example.com", "1234567890")
		require.NoError(t, err)

		tokenPair, err := service.RefreshAccessToken(refreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenPair.AccessToken)
		assert.NotEmpty(t, tokenPair.RefreshToken)
	})

	t.Run("new tokens are valid", func(t *testing.T) {
		refreshToken, err := service.GenerateRefreshToken(123, "test@example.com", "1234567890")
		require.NoError(t, err)

		tokenPair, err := service.RefreshAccessToken(refreshToken)
		require.NoError(t, err)

		accessClaims, err := service.VerifyToken(tokenPair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, "access", accessClaims.Type)

		refreshClaims, err := service.VerifyRefreshToken(tokenPair.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, "refresh", refreshClaims.Type)
	})

	t.Run("returns error with access token", func(t *testing.T) {
		accessToken, err := service.GenerateAccessToken(123, "test@example.com", "1234567890")
		require.NoError(t, err)

		tokenPair, err := service.RefreshAccessToken(accessToken)
		assert.Error(t, err)
		assert.Nil(t, tokenPair)
	})

	t.Run("returns error with invalid token", func(t *testing.T) {
		tokenPair, err := service.RefreshAccessToken("invalid.token")
		assert.Error(t, err)
		assert.Nil(t, tokenPair)
	})
}

func TestTokenService_RFC7519_Compliance(t *testing.T) {
	config := setupTestConfig()
	service, err := NewTokenService(config)
	require.NoError(t, err)

	t.Run("includes all RFC 7519 standard claims", func(t *testing.T) {
		token, err := service.GenerateAccessToken(123, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)

		assert.NotEmpty(t, claims.Issuer, "iss claim should be present")
		assert.NotEmpty(t, claims.Subject, "sub claim should be present")
		assert.NotNil(t, claims.ExpiresAt, "exp claim should be present")
		assert.NotNil(t, claims.IssuedAt, "iat claim should be present")
		assert.NotNil(t, claims.NotBefore, "nbf claim should be present")
		assert.NotEmpty(t, claims.ID, "jti claim should be present")
	})

	t.Run("nbf claim prevents early use", func(t *testing.T) {
		token, err := service.GenerateAccessToken(123, "test@example.com", "1234567890")
		require.NoError(t, err)

		claims, err := service.VerifyToken(token)
		require.NoError(t, err)

		now := time.Now()
		assert.WithinDuration(t, now, claims.NotBefore.Time, 1*time.Second)
	})
}
