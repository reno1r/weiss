package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup(t)

	resp := env.Request(t, http.MethodGet, "/health", nil)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	resp.JSON(t, &body)
	assert.Equal(t, "ok", body["status"])
}

func TestRegisterEndpoint(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup(t)

	t.Run("successful registration", func(t *testing.T) {
		env.CleanupDB(t)

		payload := map[string]string{
			"full_name": "John Doe",
			"email":     "john@example.com",
			"phone":     "+1234567890",
			"password":  "SecurePass123!",
		}

		resp := env.Request(t, http.MethodPost, "/api/auth/register", payload)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var body map[string]any
		resp.JSON(t, &body)
		assert.Equal(t, "user created successfully.", body["message"])

		data := body["data"].(map[string]any)
		user := data["user"].(map[string]any)
		assert.Equal(t, "John Doe", user["full_name"])
		assert.Equal(t, "john@example.com", user["email"])
		assert.Equal(t, "+1234567890", user["phone"])
		assert.Empty(t, user["password"]) // Password should not be returned
	})

	t.Run("duplicate email", func(t *testing.T) {
		env.CleanupDB(t)

		payload := map[string]string{
			"full_name": "John Doe",
			"email":     "duplicate@example.com",
			"phone":     "+1234567890",
			"password":  "SecurePass123!",
		}

		// First registration
		resp := env.Request(t, http.MethodPost, "/api/auth/register", payload)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Second registration with same email
		payload["phone"] = "+0987654321" // Different phone
		resp = env.Request(t, http.MethodPost, "/api/auth/register", payload)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("invalid request body", func(t *testing.T) {
		resp := env.Request(t, http.MethodPost, "/api/auth/register", "invalid json")

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing required fields", func(t *testing.T) {
		payload := map[string]string{
			"email": "test@example.com",
		}

		resp := env.Request(t, http.MethodPost, "/api/auth/register", payload)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestLoginEndpoint(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup(t)

	t.Run("successful login with email", func(t *testing.T) {
		env.CleanupDB(t)

		// First register a user
		registerPayload := map[string]string{
			"full_name": "Jane Doe",
			"email":     "jane@example.com",
			"phone":     "+1234567890",
			"password":  "SecurePass123!",
		}
		resp := env.Request(t, http.MethodPost, "/api/auth/register", registerPayload)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Now login
		loginPayload := map[string]string{
			"email":    "jane@example.com",
			"password": "SecurePass123!",
		}
		resp = env.Request(t, http.MethodPost, "/api/auth/login", loginPayload)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		resp.JSON(t, &body)
		assert.Equal(t, "authorized.", body["message"])

		data := body["data"].(map[string]any)
		assert.NotEmpty(t, data["access_token"])
		assert.NotEmpty(t, data["refresh_token"])
	})

	t.Run("successful login with phone", func(t *testing.T) {
		env.CleanupDB(t)

		// First register a user
		registerPayload := map[string]string{
			"full_name": "Jane Doe",
			"email":     "jane2@example.com",
			"phone":     "+9876543210",
			"password":  "SecurePass123!",
		}
		resp := env.Request(t, http.MethodPost, "/api/auth/register", registerPayload)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Now login with phone
		loginPayload := map[string]string{
			"phone":    "+9876543210",
			"password": "SecurePass123!",
		}
		resp = env.Request(t, http.MethodPost, "/api/auth/login", loginPayload)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		resp.JSON(t, &body)
		assert.Equal(t, "authorized.", body["message"])
	})

	t.Run("invalid credentials", func(t *testing.T) {
		env.CleanupDB(t)

		// Register a user
		registerPayload := map[string]string{
			"full_name": "Jane Doe",
			"email":     "jane3@example.com",
			"phone":     "+1111111111",
			"password":  "SecurePass123!",
		}
		resp := env.Request(t, http.MethodPost, "/api/auth/register", registerPayload)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Login with wrong password
		loginPayload := map[string]string{
			"email":    "jane3@example.com",
			"password": "WrongPassword!",
		}
		resp = env.Request(t, http.MethodPost, "/api/auth/login", loginPayload)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("non-existent user", func(t *testing.T) {
		loginPayload := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "SomePassword123!",
		}
		resp := env.Request(t, http.MethodPost, "/api/auth/login", loginPayload)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestAuthFlow(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup(t)

	t.Run("complete auth flow", func(t *testing.T) {
		env.CleanupDB(t)

		// Generate unique identifiers for this test run
		timestamp := time.Now().UnixNano()
		email := fmt.Sprintf("flow_%d@example.com", timestamp)
		phone := fmt.Sprintf("+1%010d", timestamp%10000000000)

		// 1. Register
		registerPayload := map[string]string{
			"full_name": "Flow Test User",
			"email":     email,
			"phone":     phone,
			"password":  "FlowTestPass123!",
		}
		resp := env.Request(t, http.MethodPost, "/api/auth/register", registerPayload)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// 2. Login
		loginPayload := map[string]string{
			"email":    email,
			"password": "FlowTestPass123!",
		}
		resp = env.Request(t, http.MethodPost, "/api/auth/login", loginPayload)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp map[string]any
		resp.JSON(t, &loginResp)

		data := loginResp["data"].(map[string]any)
		accessToken := data["access_token"].(string)
		refreshToken := data["refresh_token"].(string)

		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.NotEqual(t, accessToken, refreshToken)
	})
}
