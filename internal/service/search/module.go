// Package search provides search functionality using Typesense.
package search

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

// Module provides search service dependencies.
var Module = fx.Module("search_service",
	fx.Provide(NewMovieSearchService),
	fx.Provide(NewTVShowSearchService),
	fx.Invoke(initializeCollections),
)

// initializeCollections creates Typesense collections on startup if they don't exist.
func initializeCollections(lc fx.Lifecycle, movieSearch *MovieSearchService, tvshowSearch *TVShowSearchService, logger *slog.Logger) {
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
			return nil
		},
	})
}
