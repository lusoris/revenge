package storage

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
)

// Module provides the storage service for fx dependency injection.
var Module = fx.Options(
	fx.Provide(provideStorage),
)

// provideStorage creates the storage service based on configured backend.
// Supports "local" (filesystem) and "s3" (S3-compatible) backends.
func provideStorage(cfg *config.Config, logger *slog.Logger) (Storage, error) {
	switch cfg.Storage.Backend {
	case "s3":
		logger.Info("Using S3 storage backend",
			slog.String("endpoint", cfg.Storage.S3.Endpoint),
			slog.String("bucket", cfg.Storage.S3.Bucket),
			slog.String("region", cfg.Storage.S3.Region))
		return NewS3Storage(cfg.Storage.S3, logger)

	case "local":
		logger.Info("Using local filesystem storage backend",
			slog.String("path", cfg.Storage.Local.Path))
		// For backwards compatibility, use Avatar.StoragePath if Storage.Local.Path is empty
		path := cfg.Storage.Local.Path
		if path == "" {
			path = cfg.Avatar.StoragePath
		}
		localCfg := config.AvatarConfig{
			StoragePath:  path,
			MaxSizeBytes: cfg.Avatar.MaxSizeBytes,
			AllowedTypes: cfg.Avatar.AllowedTypes,
		}
		return NewLocalStorage(localCfg, logger)

	default:
		logger.Error("Unknown storage backend, falling back to local",
			slog.String("backend", cfg.Storage.Backend))
		return NewLocalStorage(cfg.Avatar, logger)
	}
}
