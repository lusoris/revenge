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
	dbURL := fmt.Sprintf("postgres://revenge_test:test_pass@%s:%s/revenge_test?sslmode=disable&search_path=public,shared",
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
	URL string
}

// NewDragonflyContainer starts a Dragonfly container for integration testing.
// Currently returns a stub - will be implemented when cache is needed.
func NewDragonflyContainer(t *testing.T) *DragonflyContainer {
	t.Helper()
	t.Skip("Dragonfly container not yet implemented - implement when cache module is needed")
	return nil
}

// TypesenseContainer represents a Typesense container for integration testing.
type TypesenseContainer struct {
	URL    string
	APIKey string
}

// NewTypesenseContainer starts a Typesense container for integration testing.
// Currently returns a stub - will be implemented when search is needed.
func NewTypesenseContainer(t *testing.T) *TypesenseContainer {
	t.Helper()
	t.Skip("Typesense container not yet implemented - implement when search module is needed")
	return nil
}
