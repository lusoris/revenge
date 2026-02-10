package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx driver
	"github.com/lusoris/revenge/internal/errors"
	"github.com/lusoris/revenge/internal/util"
)

//go:embed migrations/shared/*.sql
var migrationsFS embed.FS

// MigrateUp runs all pending migrations.
func MigrateUp(databaseURL string, logger *slog.Logger) error {
	m, err := newMigrate(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close() //nolint:errcheck // Deferred cleanup, error not actionable

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return errors.Wrap(err, "failed to get migration version")
	}

	logger.Info("running migrations",
		slog.Uint64("current_version", uint64(version)),
		slog.Bool("dirty", dirty),
	)

	if dirty {
		logger.Warn("database is in dirty state, forcing version reset",
			slog.Uint64("version", uint64(version)),
		)
		if err := m.Force(util.SafeUintToInt(version)); err != nil {
			return errors.Wrap(err, "failed to force migration version")
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to run migrations")
	}

	newVersion, _, err := m.Version()
	if err != nil {
		return errors.Wrap(err, "failed to get new migration version")
	}

	logger.Info("migrations completed",
		slog.Uint64("version", uint64(newVersion)),
	)

	return nil
}

// MigrateDown rolls back one migration.
func MigrateDown(databaseURL string, logger *slog.Logger) error {
	m, err := newMigrate(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close() //nolint:errcheck // Deferred cleanup, error not actionable

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return errors.Wrap(err, "failed to get migration version")
	}

	logger.Info("rolling back migration",
		slog.Uint64("current_version", uint64(version)),
		slog.Bool("dirty", dirty),
	)

	if dirty {
		logger.Warn("database is in dirty state, forcing version reset before rollback",
			slog.Uint64("version", uint64(version)),
		)
		if err := m.Force(util.SafeUintToInt(version)); err != nil {
			return errors.Wrap(err, "failed to force migration version")
		}
	}

	if err := m.Steps(-1); err != nil {
		return errors.Wrap(err, "failed to rollback migration")
	}

	newVersion, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return errors.Wrap(err, "failed to get new migration version")
	}

	logger.Info("migration rolled back",
		slog.Uint64("version", uint64(newVersion)),
	)

	return nil
}

// MigrateVersion returns the current migration version.
func MigrateVersion(databaseURL string) (uint, bool, error) {
	m, err := newMigrate(databaseURL)
	if err != nil {
		return 0, false, err
	}
	defer m.Close() //nolint:errcheck // Deferred cleanup, error not actionable

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, errors.Wrap(err, "failed to get migration version")
	}

	return version, dirty, nil
}

// MigrateTo migrates to a specific version.
func MigrateTo(databaseURL string, version uint, logger *slog.Logger) error {
	m, err := newMigrate(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close() //nolint:errcheck // Deferred cleanup, error not actionable

	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return errors.Wrap(err, "failed to get migration version")
	}

	logger.Info("migrating to version",
		slog.Uint64("current_version", uint64(currentVersion)),
		slog.Uint64("target_version", uint64(version)),
		slog.Bool("dirty", dirty),
	)

	if dirty {
		logger.Warn("database is in dirty state, forcing version reset before targeted migration",
			slog.Uint64("version", uint64(currentVersion)),
		)
		if err := m.Force(util.SafeUintToInt(currentVersion)); err != nil {
			return errors.Wrap(err, "failed to force migration version")
		}
	}

	if err := m.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, fmt.Sprintf("failed to migrate to version %d", version))
	}

	logger.Info("migration completed",
		slog.Uint64("version", uint64(version)),
	)

	return nil
}

// newMigrate creates a new migrate instance with embedded migrations.
func newMigrate(databaseURL string) (*migrate.Migrate, error) {
	// Open database connection using pgx driver
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		_ = db.Close() // Best-effort cleanup
		return nil, errors.Wrap(err, "failed to create postgres driver")
	}

	// Create source from embedded filesystem
	sourceDriver, err := iofs.New(migrationsFS, "migrations/shared")
	if err != nil {
		_ = db.Close() // Best-effort cleanup
		return nil, errors.Wrap(err, "failed to create source driver")
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		_ = db.Close() // Best-effort cleanup
		return nil, errors.Wrap(err, "failed to create migrate instance")
	}

	return m, nil
}

func init() {
	// pgx driver is automatically registered by importing github.com/jackc/pgx/v5/stdlib
}
