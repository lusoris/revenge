package main

import (
	"fmt"
	"os"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/logging"
)

// runMigrate handles database migration subcommands.
func runMigrate() {
	if len(os.Args) < 3 {
		printMigrateUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load(config.DefaultConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Execute migration command
	subcommand := os.Args[2]
	switch subcommand {
	case "up":
		if err := database.MigrateUp(cfg.Database.URL, logger); err != nil {
			fmt.Fprintf(os.Stderr, "Migration up failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		if err := database.MigrateDown(cfg.Database.URL, logger); err != nil {
			fmt.Fprintf(os.Stderr, "Migration down failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migration rolled back successfully")

	case "version":
		version, dirty, err := database.MigrateVersion(cfg.Database.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get migration version: %v\n", err)
			os.Exit(1)
		}
		dirtyStr := ""
		if dirty {
			dirtyStr = " (dirty)"
		}
		fmt.Printf("Current migration version: %d%s\n", version, dirtyStr)

	case "to":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "Usage: revenge migrate to <version>\n")
			os.Exit(1)
		}
		var targetVersion uint
		if _, err := fmt.Sscanf(os.Args[3], "%d", &targetVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid version number: %v\n", err)
			os.Exit(1)
		}
		if err := database.MigrateTo(cfg.Database.URL, targetVersion, logger); err != nil {
			fmt.Fprintf(os.Stderr, "Migration to version %d failed: %v\n", targetVersion, err)
			os.Exit(1)
		}
		fmt.Printf("Migrated to version %d successfully\n", targetVersion)

	default:
		fmt.Fprintf(os.Stderr, "Unknown migration command: %s\n", subcommand)
		printMigrateUsage()
		os.Exit(1)
	}
}

func printMigrateUsage() {
	fmt.Println("Database migration commands")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  revenge migrate up              Apply all pending migrations")
	fmt.Println("  revenge migrate down            Rollback the last migration")
	fmt.Println("  revenge migrate version         Show current migration version")
	fmt.Println("  revenge migrate to <version>    Migrate to a specific version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  revenge migrate up")
	fmt.Println("  revenge migrate down")
	fmt.Println("  revenge migrate version")
	fmt.Println("  revenge migrate to 3")
}
