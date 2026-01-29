package metadata

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/service/metadata/radarr"
	"github.com/lusoris/revenge/internal/service/metadata/tmdb"
)

// Module provides the metadata service and its dependencies.
var Module = fx.Module("metadata",
	// Include provider modules
	radarr.Module,
	tmdb.Module,

	// Provide the central metadata service
	fx.Provide(func(
		radarrProvider *radarr.Provider,
		tmdbProvider *tmdb.Provider,
		localCache *cache.LocalCache,
		apiCache *cache.APICache,
		logger *slog.Logger,
	) *Service {
		return NewService(radarrProvider, tmdbProvider, localCache, apiCache, logger)
	}),
)
