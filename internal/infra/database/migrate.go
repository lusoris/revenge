// Package database provides database infrastructure for Jellyfin Go.
package database

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrator handles database migrations.
type Migrator struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewMigrator creates a new database migrator.
func NewMigrator(pool *pgxpool.Pool, logger *slog.Logger) *Migrator {
	return &Migrator{
		pool:   pool,
		logger: logger.With(slog.String("component", "migrator")),
	}
}

// Up runs all pending migrations.
func (m *Migrator) Up(ctx context.Context) error {
	migrator, err := m.getMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	version, dirty, verErr := migrator.Version()
	if verErr != nil && !errors.Is(verErr, migrate.ErrNilVersion) {
		m.logger.Warn("failed to get current version", slog.Any("error", verErr))
	}
	m.logger.Info("starting migrations",
		slog.Uint64("current_version", uint64(version)),
		slog.Bool("dirty", dirty))

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("no migrations to run")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	newVersion, _, verErr := migrator.Version()
	if verErr != nil {
		m.logger.Warn("failed to get new version", slog.Any("error", verErr))
	}
	m.logger.Info("migrations completed",
		slog.Uint64("new_version", uint64(newVersion)))

	return nil
}

// Down rolls back all migrations.
func (m *Migrator) Down(ctx context.Context) error {
	migrator, err := m.getMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("no migrations to roll back")
			return nil
		}
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	m.logger.Info("all migrations rolled back")
	return nil
}

// Steps runs n migrations (positive = up, negative = down).
func (m *Migrator) Steps(ctx context.Context, n int) error {
	migrator, err := m.getMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Steps(n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("no migrations to run")
			return nil
		}
		return fmt.Errorf("failed to run migration steps: %w", err)
	}

	version, _, verErr := migrator.Version()
	if verErr != nil {
		m.logger.Warn("failed to get version after steps", slog.Any("error", verErr))
	}
	m.logger.Info("migration steps completed",
		slog.Int("steps", n),
		slog.Uint64("current_version", uint64(version)))

	return nil
}

// Version returns the current migration version.
func (m *Migrator) Version(ctx context.Context) (uint, bool, error) {
	migrator, err := m.getMigrate()
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Version()
}

// Force sets the migration version without running migrations.
// Use with caution - only for fixing dirty state.
func (m *Migrator) Force(ctx context.Context, version int) error {
	migrator, err := m.getMigrate()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Force(version); err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	m.logger.Warn("migration version forced",
		slog.Int("version", version))

	return nil
}

func (m *Migrator) getMigrate() (*migrate.Migrate, error) {
	// Create source from embedded filesystem
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create source driver: %w", err)
	}

	// Get a *sql.DB from the pool for the migrate driver
	db := stdlib.OpenDBFromPool(m.pool)

	// Create database driver
	dbDriver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create migrate instance
	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return migrator, nil
}
