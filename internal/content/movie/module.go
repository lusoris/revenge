package movie

import (
	"go.uber.org/fx"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/config"
)

// Module provides the movie content module
var Module = fx.Module("movie",
	fx.Provide(
		// Core movie services
		NewPostgresRepository,
		NewService,
		NewHandler,

		// Metadata service
		provideMetadataService,

		// Library service
		provideLibraryService,
	),
)

// provideMetadataService creates TMDb metadata service from config.
func provideMetadataService(cfg *config.Config) *MetadataService {
	return NewMetadataService(TMDbConfig{
		APIKey:    cfg.Movie.TMDb.APIKey,
		RateLimit: rate.Limit(cfg.Movie.TMDb.RateLimit),
		CacheTTL:  cfg.Movie.TMDb.CacheTTL,
		ProxyURL:  cfg.Movie.TMDb.ProxyURL,
	})
}

// provideLibraryService creates library service from config.
func provideLibraryService(
	repo Repository,
	metadataService *MetadataService,
	cfg *config.Config,
) *LibraryService {
	return NewLibraryService(repo, metadataService, cfg.Movie.Library)
}
