package movie

import (
	"context"

	"github.com/google/uuid"
)

// MetadataProvider defines the interface for metadata operations.
// Implementations include the shared metadata service adapter.
type MetadataProvider interface {
	SearchMovies(ctx context.Context, query string, year *int) ([]*Movie, error)
	EnrichMovie(ctx context.Context, mov *Movie) error
	GetMovieCredits(ctx context.Context, movieID uuid.UUID, tmdbID int) ([]MovieCredit, error)
	GetMovieGenres(ctx context.Context, movieID uuid.UUID, tmdbID int) ([]MovieGenre, error)
	ClearCache()
}
