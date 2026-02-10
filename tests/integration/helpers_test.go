//go:build integration

package integration

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/api"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/lusoris/revenge/internal/testutil"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// TestServer represents a running test server instance.
type TestServer struct {
	App        *fxtest.App
	BaseURL    string
	HTTPClient *http.Client
	DB         *testutil.PostgreSQLContainer
}

// setupServer starts a test server with all dependencies.
func setupServer(t *testing.T) *TestServer {
	t.Helper()

	// Start PostgreSQL container
	pgContainer := testutil.NewPostgreSQLContainer(t)

	// Create test configuration
	cfg := pgContainer.Config

	// Find available port for HTTP server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		pgContainer.Close()
		t.Fatalf("failed to get available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()

	cfg.Server.Host = "localhost"
	cfg.Server.Port = port
	cfg.Server.ReadTimeout = 5000000000  // 5s
	cfg.Server.WriteTimeout = 5000000000 // 5s

	// Create fx app for testing
	// Note: We don't use app.Module because it includes config.Module
	// Instead, we provide cfg directly and include only the needed modules
	app := fxtest.New(t,
		// Provide configuration directly (from PostgreSQL container)
		fx.Supply(cfg),

		// Infrastructure modules (logging provides zap.Logger)
		logging.Module,
		database.Module,
		cache.Module,
		search.Module,
		jobs.Module,
		health.Module,
		image.Module,
		appcrypto.Module,

		// Service modules (required by API)
		user.Module,
		auth.Module,
		session.Module,
		settings.Module,
		rbac.Module,
		apikeys.Module,
		oidc.Module,
		activity.Module,
		library.Module,
		mfa.Module,

		// Content modules
		movie.Module,

		// Stubs for dependencies not needed in integration tests
		fx.Provide(func() []*river.PeriodicJob { return nil }),
		fx.Provide(func() movie.MetadataProvider { return nil }),
		fx.Provide(func() movie.MetadataQueue { return nil }),

		// API module
		api.Module,
	)

	// Start app
	app.RequireStart()

	// Wait for server to be ready
	baseURL := fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	waitForServer(t, baseURL, 10*time.Second)

	return &TestServer{
		App:     app,
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		DB: pgContainer,
	}
}

// teardownServer stops the test server and cleans up resources.
func teardownServer(t *testing.T, ts *TestServer) {
	t.Helper()

	if ts.App != nil {
		ts.App.RequireStop()
	}

	if ts.DB != nil {
		ts.DB.Close()
	}
}

// waitForServer waits for the HTTP server to be ready.
func waitForServer(t *testing.T, baseURL string, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	client := &http.Client{Timeout: 1 * time.Second}
	healthURL := baseURL + "/healthz"

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("server did not become ready within %v", timeout)
		case <-ticker.C:
			resp, err := client.Get(healthURL)
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				return
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

// resetDatabase truncates all tables in the test database.
func resetDatabase(t *testing.T, ts *TestServer) {
	t.Helper()
	ts.DB.Reset(t)
}

// TestMain runs before all tests.
func TestMain(m *testing.M) {
	// Check if Docker is available
	if !isDockerAvailable() {
		fmt.Println("Docker is not available - skipping integration tests")
		os.Exit(0)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// isDockerAvailable checks if Docker is available on the system.
func isDockerAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// On Windows, check via Docker CLI which uses named pipes internally
	if _, err := os.Stat(`\\.\pipe\docker_engine`); err == nil {
		return true
	}

	// Try Unix socket (Linux/macOS)
	conn, err := net.DialTimeout("unix", "/var/run/docker.sock", 2*time.Second)
	if err == nil {
		conn.Close()
		return true
	}

	// Try Docker Desktop via TCP
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:2375/_ping")
	if err == nil {
		resp.Body.Close()
		return true
	}

	_ = ctx // Satisfy linter
	return false
}
