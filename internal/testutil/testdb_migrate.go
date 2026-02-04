package testutil

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// pathToFileURL converts a filesystem path to a proper file:// URL
// Handles both Windows (C:\path) and Unix (/path) paths correctly
func pathToFileURL(path string) string {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	// Convert to URL-safe format
	u := &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(absPath),
	}

	return u.String()
}

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot() (string, error) {
	// Start from the current file's directory
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	dir := filepath.Dir(currentFile)

	// Walk up until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (go.mod)")
		}
		dir = parent
	}
}

// runMigrationsWithMigrate runs migrations from the project's migrations directory
func runMigrationsWithMigrate(databaseURL string) error {
	// Find the actual migrations directory in project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	migrationsPath := filepath.Join(projectRoot, "migrations")
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found: %s", migrationsPath)
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Use file:// source instead of embedded FS
	// Convert to proper file:// URL format (handles Windows paths correctly)
	sourceURL := pathToFileURL(migrationsPath)
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		_ = sourceErr
		_ = dbErr
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
