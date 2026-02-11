package movie

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

// MetadataProvider defines the interface for metadata operations.
// Implementations include the shared metadata service adapter.
type MetadataProvider interface {
	SearchMovies(ctx context.Context, query string, year *int) ([]*Movie, error)
	EnrichMovie(ctx context.Context, mov *Movie, opts ...MetadataRefreshOptions) error
	GetMovieCredits(ctx context.Context, movieID uuid.UUID, providerID string) ([]MovieCredit, error)
	GetMovieGenres(ctx context.Context, movieID uuid.UUID, providerID string) ([]MovieGenre, error)
	ClearCache()
}
