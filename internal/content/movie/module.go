package movie

import (
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
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

// provideService creates movie service with MetadataProvider.
// MetadataProvider is injected from metadatafx module (MovieMetadataAdapter).
func provideService(repo Repository, metadataProvider MetadataProvider) Service {
	return NewService(repo, metadataProvider)
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
