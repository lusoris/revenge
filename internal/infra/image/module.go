package image

import (
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

// Module provides image service dependencies.
var Module = fx.Module("image",
	fx.Provide(
		NewImageServiceFromConfig,
	),
)

// NewImageServiceFromConfig creates a new image service using app config defaults.
// The image proxy doesn't require external API keys — it proxies TMDb's public image CDN.
func NewImageServiceFromConfig(cfg *config.Config, logger *slog.Logger) (*Service, error) {
	cacheDir := ""
	if dataDir := os.Getenv("REVENGE_DATA_DIR"); dataDir != "" {
		cacheDir = filepath.Join(dataDir, "image-cache")
	} else if home, err := os.UserCacheDir(); err == nil {
		cacheDir = filepath.Join(home, "revenge", "image-cache")
	}

	return NewService(Config{
		CacheDir: cacheDir,
		// All other fields use sensible defaults:
		// BaseURL: https://image.tmdb.org/t/p
		// CacheTTL: 7 days
		// MaxSize: 10MB
	}, logger)
}
