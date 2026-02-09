// Package search provides search functionality using Typesense.
package search

import (
	"context"
	"log/slog"

	"github.com/lusoris/revenge/internal/infra/cache"
	"go.uber.org/fx"
)

// Module provides search service dependencies.
var Module = fx.Module("search_service",
	fx.Provide(NewMovieSearchService),
	fx.Provide(NewTVShowSearchService),
	fx.Provide(NewEpisodeSearchService),
	fx.Provide(NewSeasonSearchService),
	fx.Provide(provideCachedMovieSearchService),
	fx.Provide(provideCachedTVShowSearchService),
	fx.Provide(provideCachedEpisodeSearchService),
	fx.Provide(provideCachedSeasonSearchService),
	fx.Invoke(initializeCollections),
)

// provideCachedMovieSearchService wraps the movie search service with caching.
func provideCachedMovieSearchService(svc *MovieSearchService, c *cache.Cache, logger *slog.Logger) *CachedMovieSearchService {
	return NewCachedMovieSearchService(svc, c, logger)
}

// provideCachedTVShowSearchService wraps the TV show search service with caching.
func provideCachedTVShowSearchService(svc *TVShowSearchService, c *cache.Cache, logger *slog.Logger) *CachedTVShowSearchService {
	return NewCachedTVShowSearchService(svc, c, logger)
}

// provideCachedEpisodeSearchService wraps the episode search service with caching.
func provideCachedEpisodeSearchService(svc *EpisodeSearchService, c *cache.Cache, logger *slog.Logger) *CachedEpisodeSearchService {
	return NewCachedEpisodeSearchService(svc, c, logger)
}

// provideCachedSeasonSearchService wraps the season search service with caching.
func provideCachedSeasonSearchService(svc *SeasonSearchService, c *cache.Cache, logger *slog.Logger) *CachedSeasonSearchService {
	return NewCachedSeasonSearchService(svc, c, logger)
}

// initializeCollections creates Typesense collections on startup if they don't exist.
func initializeCollections(lc fx.Lifecycle, movieSearch *MovieSearchService, tvshowSearch *TVShowSearchService, episodeSearch *EpisodeSearchService, seasonSearch *SeasonSearchService, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if movieSearch.IsEnabled() {
				if err := movieSearch.InitializeCollection(ctx); err != nil {
					logger.Warn("failed to initialize movies search collection", slog.Any("error", err))
				}
			}
			if tvshowSearch.IsEnabled() {
				if err := tvshowSearch.InitializeCollection(ctx); err != nil {
					logger.Warn("failed to initialize tvshows search collection", slog.Any("error", err))
				}
			}
			if episodeSearch.IsEnabled() {
				if err := episodeSearch.InitializeCollection(ctx); err != nil {
					logger.Warn("failed to initialize episodes search collection", slog.Any("error", err))
				}
			}
			if seasonSearch.IsEnabled() {
				if err := seasonSearch.InitializeCollection(ctx); err != nil {
					logger.Warn("failed to initialize seasons search collection", slog.Any("error", err))
				}
			}
			return nil
		},
	})
}
