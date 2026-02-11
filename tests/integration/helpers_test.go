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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/api"
	"github.com/lusoris/revenge/internal/api/sse"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/movie/moviejobs"
	"github.com/lusoris/revenge/internal/content/tvshow"
	tvshowjobs "github.com/lusoris/revenge/internal/content/tvshow/jobs"
	appcrypto "github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/image"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/infra/raft"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/integration/radarr"
	"github.com/lusoris/revenge/internal/integration/sonarr"
	"github.com/lusoris/revenge/internal/playback/playbackfx"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/analytics"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/email"
	"github.com/lusoris/revenge/internal/service/library"
	metadatajobs "github.com/lusoris/revenge/internal/service/metadata/jobs"
	"github.com/lusoris/revenge/internal/service/metadata/metadatafx"
	"github.com/lusoris/revenge/internal/service/mfa"
	"github.com/lusoris/revenge/internal/service/notification"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/rbac"
	searchsvc "github.com/lusoris/revenge/internal/service/search"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/storage"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/lusoris/revenge/internal/testutil"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// TestServer represents a running test server instance.
type TestServer struct {
	App         *fxtest.App
	BaseURL     string
	HTTPClient  *http.Client
	DB          *testutil.PostgreSQLContainer
	AppPool     *pgxpool.Pool  // DI-managed database pool (for health tests)
	RBACService *rbac.Service // RBAC service for policy reload after direct DB changes
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

	// Use a temp directory for local storage in tests
	storageDir := t.TempDir()
	cfg.Storage.Backend = "local"
	cfg.Storage.Local.Path = storageDir

	// Auth config (empty JWT secret / zero expiry causes instant token expiry)
	cfg.Auth.JWTSecret = "integration-test-secret-key-must-be-32chars!!"
	cfg.Auth.JWTExpiry = 15 * time.Minute
	cfg.Auth.RefreshExpiry = 24 * time.Hour
	cfg.Auth.LockoutThreshold = 5
	cfg.Auth.LockoutWindow = 15 * time.Minute
	cfg.Auth.LockoutEnabled = false

	// Capture the DI-managed pool for health tests
	var appPool *pgxpool.Pool
	var rbacService *rbac.Service

	// Create fx app for testing
	// Mirrors app.Module but replaces config.Module with fx.Supply(cfg)
	// and omits observability.Module (starts extra HTTP server for metrics).
	app := fxtest.New(t,
		// Provide configuration directly (from PostgreSQL container)
		fx.Supply(cfg),

		// Infrastructure
		logging.Module,
		database.Module,
		cache.Module,
		search.Module,
		jobs.Module,
		raft.Module,
		health.Module,
		image.Module,
		appcrypto.Module,

		// Periodic jobs stub (no real periodic jobs in tests)
		fx.Provide(func() []*river.PeriodicJob { return nil }),
		fx.Invoke(
			func(workers *river.Workers, w *activity.ActivityCleanupWorker) { river.AddWorker(workers, w) },
			func(workers *river.Workers, w *library.LibraryScanCleanupWorker) { river.AddWorker(workers, w) },
		),

		// Bridge: metadata jobs Queue → movie.MetadataQueue interface
		fx.Provide(func(q *metadatajobs.Queue) movie.MetadataQueue { return q }),

		// Services
		settings.Module,
		user.Module,
		auth.Module,
		email.Module,
		session.Module,
		rbac.Module,
		apikeys.Module,
		mfa.Module,
		oidc.Module,
		activity.Module,
		analytics.Module,
		notification.Module,
		storage.Module,
		library.Module,
		searchsvc.Module,

		// Content modules
		movie.Module,
		tvshow.Module,

		// Playback / HLS Streaming
		playbackfx.Module,

		// Job Workers
		moviejobs.Module,
		tvshowjobs.Module,
		metadatajobs.Module,

		// Integrations
		radarr.Module,
		sonarr.Module,

		// Metadata Service
		metadatafx.Module,

		// SSE Real-Time Events
		sse.Module,

		// HTTP API Server (ogen-generated)
		api.Module,

		// Extract DI-managed services for test use
		fx.Populate(&appPool),
		fx.Populate(&rbacService),
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
		DB:          pgContainer,
		AppPool:     appPool,
		RBACService: rbacService,
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
	// go test sets cwd to the package directory (tests/integration/).
	// Many modules expect to find files relative to the repo root
	// (e.g. config/casbin_model.conf), so chdir to repo root.
	if err := os.Chdir("../.."); err != nil {
		fmt.Fprintf(os.Stderr, "failed to chdir to repo root: %v\n", err)
		os.Exit(1)
	}

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
