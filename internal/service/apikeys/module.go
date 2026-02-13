package apikeys

import (
	"time"

	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
	"go.uber.org/fx"
)

// Module provides the API keys service
var Module = fx.Module("apikeys",
	fx.Provide(
		NewRepositoryPg,
		provideService,
		provideConfig,
	),
)

// provideService creates the API keys service wrapped with caching.
// Cache may be nil if caching is disabled — CachedService handles nil gracefully.
func provideService(repo Repository, c *cache.Cache, logger *slog.Logger, maxKeysPerUser int, defaultExpiry time.Duration) Service {
	base := NewService(repo, logger, maxKeysPerUser, defaultExpiry)
	return NewCachedService(base, repo, c, logger)
}

// provideConfig extracts API keys configuration
func provideConfig(cfg *config.Config) (int, time.Duration) {
	// For now, use defaults until we add config keys
	maxKeysPerUser := 10
	var defaultExpiry time.Duration = 0 // Never expire

	return maxKeysPerUser, defaultExpiry
}
