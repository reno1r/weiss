package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	accessentities "github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	shopentities "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/config"
	weisshttp "github.com/reno1r/weiss/apps/service/internal/http"
)

// TestEnv holds the test environment configuration
type TestEnv struct {
	App       *fiber.App
	DB        *gorm.DB
	Container testcontainers.Container
	Ctx       context.Context
	isLive    bool
	liveURL   string
}

// SetupTestEnv creates a new test environment with testcontainers
func SetupTestEnv(t *testing.T) *TestEnv {
	// Check if testing against live API
	if liveURL := os.Getenv("E2E_LIVE_URL"); liveURL != "" {
		return &TestEnv{
			isLive:  true,
			liveURL: liveURL,
			Ctx:     context.Background(),
		}
	}

	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "test",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=test sslmode=disable",
		host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate models
	err = db.AutoMigrate(
		&entities.User{},
		&shopentities.Shop{},
		&accessentities.Role{},
		&accessentities.Staff{},
	)
	require.NoError(t, err)

	// Create test config
	cfg := &config.Config{
		AppName:            "weiss-test",
		AppHost:            "127.0.0.1",
		AppPort:            "0",
		AppDebug:           true,
		BcryptCost:         4, // Lower cost for faster tests
		JwtSecret:          "test-secret-key-for-e2e-testing-purposes",
		JwtIssuer:          "weiss-test",
		JwtAccessExpires:   "15m",
		JwtRefreshExpires:  "168h", // 7 days
		CorsAllowedOrigins: "*",
	}

	// Create test middleware to set user_id from X-Test-User-ID header
	testMiddleware := func(c fiber.Ctx) error {
		if userIDHeader := c.Get("X-Test-User-ID"); userIDHeader != "" {
			var userID uint64
			if _, err := fmt.Sscanf(userIDHeader, "%d", &userID); err == nil && userID > 0 {
				c.Locals("user_id", userID)
			}
		}
		return c.Next()
	}

	server := weisshttp.NewServerWithTestMiddleware(cfg, db, testMiddleware)
	app := server.App()

	return &TestEnv{
		App:       app,
		DB:        db,
		Container: pgContainer,
		Ctx:       ctx,
	}
}

// Cleanup tears down the test environment
func (e *TestEnv) Cleanup(t *testing.T) {
	if e.isLive {
		return
	}

	if e.Container != nil {
		if err := e.Container.Terminate(e.Ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}
}

// CleanupDB truncates all tables
func (e *TestEnv) CleanupDB(t *testing.T) {
	if e.isLive || e.DB == nil {
		return
	}
	// Truncate in order to respect foreign key constraints
	err := e.DB.WithContext(e.Ctx).Exec("TRUNCATE TABLE staffs, roles, shops, users RESTART IDENTITY CASCADE").Error
	require.NoError(t, err)
}

// Request makes an HTTP request to the API
func (e *TestEnv) Request(t *testing.T, method, path string, body any) *Response {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(jsonBytes)
	}

	if e.isLive {
		return e.doLiveRequest(t, method, e.liveURL+path, reqBody)
	}

	return e.doTestRequest(t, method, path, reqBody)
}

func (e *TestEnv) doTestRequest(t *testing.T, method, path string, body io.Reader) *Response {
	req, err := http.NewRequest(method, path, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.App.Test(req, fiber.TestConfig{
		Timeout: 10 * time.Second,
	})
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}
}

// RequestWithAuth makes an authenticated HTTP request to the API
func (e *TestEnv) RequestWithAuth(t *testing.T, method, path string, body any, userID uint64) *Response {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(jsonBytes)
	}

	if e.isLive {
		return e.doLiveRequest(t, method, e.liveURL+path, reqBody)
	}

	return e.doTestRequestWithAuth(t, method, path, reqBody, userID)
}

func (e *TestEnv) doTestRequestWithAuth(t *testing.T, method, path string, body io.Reader, userID uint64) *Response {
	req, err := http.NewRequest(method, path, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if userID > 0 {
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", userID))
	}

	resp, err := e.App.Test(req, fiber.TestConfig{
		Timeout: 10 * time.Second,
	})
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}
}

func (e *TestEnv) doLiveRequest(t *testing.T, method, url string, body io.Reader) *Response {
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}
}

// Response wraps an HTTP response
type Response struct {
	StatusCode int
	Body       []byte
}

// JSON unmarshals the response body into v
func (r *Response) JSON(t *testing.T, v any) {
	err := json.Unmarshal(r.Body, v)
	require.NoError(t, err)
}
