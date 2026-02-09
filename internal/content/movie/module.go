package movie

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
)

// Module provides the movie content module
var Module = fx.Module("movie",
	fx.Provide(
		// Core movie services
		NewPostgresRepository,
		provideService,
		NewHandler,

		// Library service
		provideLibraryService,
	),
)

// provideService creates movie service wrapped with caching.
// MetadataProvider is injected from metadatafx module (MovieMetadataAdapter).
// Cache may be nil if caching is disabled — CachedService handles nil gracefully.
func provideService(repo Repository, metadataProvider MetadataProvider, c *cache.Cache, logger *slog.Logger) Service {
	base := NewService(repo, metadataProvider)
	return NewCachedService(base, c, logger)
}

// provideLibraryService creates library service from config.
// MetadataProvider is injected from metadatafx module (MovieMetadataAdapter).
func provideLibraryService(
	repo Repository,
	metadataProvider MetadataProvider,
	cfg *config.Config,
) *LibraryService {
	return NewLibraryService(repo, metadataProvider, cfg.Movie.Library, NewMediaInfoProber())
}
