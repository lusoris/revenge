package tvshow

import (
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

// provideService creates TV show service with MetadataProvider.
// MetadataProvider is injected from metadatafx module (TVShowMetadataAdapter).
func provideService(repo Repository, metadataProvider MetadataProvider) Service {
	return NewService(repo, metadataProvider)
}
