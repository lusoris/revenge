package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgreSQLContainer represents a PostgreSQL container for integration testing.
type PostgreSQLContainer struct {
	container testcontainers.Container
	Pool      *pgxpool.Pool
	URL       string
	Config    *config.Config
}

// NewPostgreSQLContainer starts a PostgreSQL container for integration testing.
// Uses testcontainers-go to start a real PostgreSQL instance in Docker.
//
// Example:
//
//	func TestIntegration(t *testing.T) {
//	    if testing.Short() {
//	        t.Skip("skipping integration test")
//	    }
//
//	    pg := testutil.NewPostgreSQLContainer(t)
//	    defer pg.Close()
//
//	    // Use pg.Pool for testing
//	}
func NewPostgreSQLContainer(t *testing.T) *PostgreSQLContainer {
	t.Helper()

	ctx := context.Background()

	// Create PostgreSQL container request
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18.1-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "revenge_test",
			"POSTGRES_PASSWORD": "test_pass",
			"POSTGRES_DB":       "revenge_test",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
			wait.ForListeningPort("5432/tcp"),
		),
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start PostgreSQL container: %v", err)
	}

	// Get connection details
	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx) // Best-effort cleanup
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		_ = container.Terminate(ctx) // Best-effort cleanup
		t.Fatalf("failed to get container port: %v", err)
	}

	// Build connection string
	// Note: Include search_path in URL for schema isolation
	// See: https://www.postgresql.org/docs/current/runtime-config-client.html#GUC-SEARCH-PATH
	dbURL := fmt.Sprintf("postgres://revenge_test:test_pass@%s:%s/revenge_test?sslmode=disable&search_path=public,shared,movie,tvshow",
		host, port.Port())

	// Run migrations
	if err := runMigrations(dbURL); err != nil {
		_ = container.Terminate(ctx) // Best-effort cleanup
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:               dbURL,
			MaxConns:          10,
			MinConns:          2,
			MaxConnLifetime:   300000000000, // 5m
			MaxConnIdleTime:   60000000000,  // 1m
			HealthCheckPeriod: 30000000000,  // 30s
		},
		Logging: config.LoggingConfig{
			Level:       "debug",
			Format:      "text",
			Development: true,
		},
	}

	// Create logger
	logger := logging.NewLogger(logging.Config{
		Level:       cfg.Logging.Level,
		Format:      cfg.Logging.Format,
		Development: cfg.Logging.Development,
	})

	// Create database pool (NewPool creates its own context internally)
	pool, err := database.NewPool(cfg, logger)
	if err != nil {
		_ = container.Terminate(ctx) // Best-effort cleanup
		t.Fatalf("failed to create database pool: %v", err)
	}

	return &PostgreSQLContainer{
		container: container,
		Pool:      pool,
		URL:       dbURL,
		Config:    cfg,
	}
}

// Close stops the PostgreSQL container and closes the pool.
func (c *PostgreSQLContainer) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
	if c.container != nil {
		ctx := context.Background()
		_ = c.container.Terminate(ctx) // Best-effort cleanup
	}
}

// Reset truncates all tables in the test database, preserving schema.
func (c *PostgreSQLContainer) Reset(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	// Truncate all tables in shared schema
	_, err := c.Pool.Exec(ctx, `
		TRUNCATE TABLE shared.sessions CASCADE;
		TRUNCATE TABLE shared.users CASCADE;
	`)
	if err != nil {
		t.Fatalf("failed to reset database: %v", err)
	}
}

// DragonflyContainer represents a Dragonfly (Redis-compatible) container for integration testing.
type DragonflyContainer struct {
	container testcontainers.Container
	URL       string
	Host      string
	Port      string
}

// NewDragonflyContainer starts a Dragonfly container for integration testing.
// Uses testcontainers-go to start a real Dragonfly instance in Docker.
//
// Example:
//
//	func TestCacheIntegration(t *testing.T) {
//	    if testing.Short() {
//	        t.Skip("skipping integration test")
//	    }
//
//	    df := testutil.NewDragonflyContainer(t)
//	    defer df.Close()
//
//	    // Use df.URL for cache client connection
//	}
func NewDragonflyContainer(t *testing.T) *DragonflyContainer {
	t.Helper()

	ctx := context.Background()

	// Create Dragonfly container request
	// NOTE: --cluster_mode=emulated is required for rueidis client compatibility
	// rueidis auto-detects cluster mode and needs emulated mode for single-node Dragonfly
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/dragonflydb/dragonfly:latest",
		ExposedPorts: []string{"6379/tcp"},
		Cmd:          []string{"--cluster_mode=emulated", "--maxmemory=256mb"},
		WaitingFor: wait.ForAll(
			wait.ForLog("accepting connections").
				WithStartupTimeout(60*time.Second),
			wait.ForListeningPort("6379/tcp"),
		),
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start Dragonfly container: %v", err)
	}

	// Get connection details
	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("failed to get container port: %v", err)
	}

	// Build connection URL (Redis-compatible)
	url := fmt.Sprintf("redis://%s:%s", host, port.Port())

	return &DragonflyContainer{
		container: container,
		URL:       url,
		Host:      host,
		Port:      port.Port(),
	}
}

// Close stops the Dragonfly container.
func (c *DragonflyContainer) Close() {
	if c.container != nil {
		ctx := context.Background()
		_ = c.container.Terminate(ctx)
	}
}

// TypesenseContainer represents a Typesense container for integration testing.
type TypesenseContainer struct {
	container testcontainers.Container
	URL       string
	Host      string
	Port      string
	APIKey    string
}

// NewTypesenseContainer starts a Typesense container for integration testing.
// Uses testcontainers-go to start a real Typesense instance in Docker.
//
// Example:
//
//	func TestSearchIntegration(t *testing.T) {
//	    if testing.Short() {
//	        t.Skip("skipping integration test")
//	    }
//
//	    ts := testutil.NewTypesenseContainer(t)
//	    defer ts.Close()
//
//	    // Use ts.URL and ts.APIKey for search client connection
//	}
func NewTypesenseContainer(t *testing.T) *TypesenseContainer {
	t.Helper()

	ctx := context.Background()

	// Test API key for integration tests
	apiKey := "test-api-key-for-integration-tests"

	// Create Typesense container request
	req := testcontainers.ContainerRequest{
		Image:        "typesense/typesense:27.1",
		ExposedPorts: []string{"8108/tcp"},
		Env: map[string]string{
			"TYPESENSE_API_KEY":  apiKey,
			"TYPESENSE_DATA_DIR": "/data",
		},
		WaitingFor: wait.ForAll(
			wait.ForHTTP("/health").
				WithPort("8108/tcp").
				WithStartupTimeout(60*time.Second),
			wait.ForListeningPort("8108/tcp"),
		),
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start Typesense container: %v", err)
	}

	// Get connection details
	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "8108")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("failed to get container port: %v", err)
	}

	// Build connection URL
	url := fmt.Sprintf("http://%s:%s", host, port.Port())

	return &TypesenseContainer{
		container: container,
		URL:       url,
		Host:      host,
		Port:      port.Port(),
		APIKey:    apiKey,
	}
}

// Close stops the Typesense container.
func (c *TypesenseContainer) Close() {
	if c.container != nil {
		ctx := context.Background()
		_ = c.container.Terminate(ctx)
	}
}
