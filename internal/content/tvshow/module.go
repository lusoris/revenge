package tvshow

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/infra/cache"
	"go.uber.org/fx"
)

// Module provides the TV show content module.
var Module = fx.Module("tvshow",
	fx.Provide(
		// Core TV show services
		NewPostgresRepository,
		provideService,
	),
)

// provideService creates TV show service wrapped with caching.
// MetadataProvider is injected from metadatafx module (TVShowMetadataAdapter).
// Cache may be nil if caching is disabled — CachedService handles nil gracefully.
func provideService(repo Repository, metadataProvider MetadataProvider, c *cache.Cache, logger *slog.Logger) Service {
	base := NewService(repo, metadataProvider)
	return NewCachedService(base, c, logger)
}
