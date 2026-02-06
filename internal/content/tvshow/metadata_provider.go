package tvshow

import (
	"context"

	"github.com/google/uuid"
)

// MetadataRefreshOptions configures how a metadata refresh is performed.
type MetadataRefreshOptions struct {
	// Force bypasses any "recently updated" checks and refreshes unconditionally.
	Force bool
	// Languages overrides the default languages for this refresh.
	// If empty, uses the adapter's configured default languages.
	Languages []string
}

// MetadataProvider defines the interface for TV show metadata operations.
// Implementations include the shared metadata service adapter.
type MetadataProvider interface {
	// SearchSeries searches for TV series by title.
	SearchSeries(ctx context.Context, query string, year *int) ([]*Series, error)

	// EnrichSeries enriches a series with full metadata from external sources.
	EnrichSeries(ctx context.Context, series *Series, opts ...MetadataRefreshOptions) error

	// EnrichSeason enriches a season with full metadata.
	EnrichSeason(ctx context.Context, season *Season, seriesTMDbID int32, opts ...MetadataRefreshOptions) error

	// EnrichEpisode enriches an episode with full metadata.
	EnrichEpisode(ctx context.Context, episode *Episode, seriesTMDbID int32, opts ...MetadataRefreshOptions) error

	// GetSeriesCredits retrieves series credits (cast and crew).
	GetSeriesCredits(ctx context.Context, seriesID uuid.UUID, tmdbID int) ([]SeriesCredit, error)

	// GetSeriesGenres retrieves series genres.
	GetSeriesGenres(ctx context.Context, seriesID uuid.UUID, tmdbID int) ([]SeriesGenre, error)

	// GetSeriesNetworks retrieves series networks.
	GetSeriesNetworks(ctx context.Context, tmdbID int) ([]Network, error)

	// ClearCache clears any cached metadata.
	ClearCache()
}
