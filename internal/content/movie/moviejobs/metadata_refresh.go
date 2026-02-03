package moviejobs

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/content/movie"
)

// MovieMetadataRefreshArgs are the arguments for refreshing movie metadata.
type MovieMetadataRefreshArgs struct {
	MovieID uuid.UUID `json:"movie_id"`
	Force   bool      `json:"force"`
}

// Kind returns the unique job kind for River
func (MovieMetadataRefreshArgs) Kind() string {
	return "movie_metadata_refresh"
}

// MovieMetadataRefreshWorker refreshes movie metadata from TMDb
type MovieMetadataRefreshWorker struct {
	river.WorkerDefaults[MovieMetadataRefreshArgs]
	movieRepo       movie.Repository
	metadataService *movie.MetadataService
	logger          *zap.Logger
}

// NewMovieMetadataRefreshWorker creates a new metadata refresh worker
func NewMovieMetadataRefreshWorker(
	movieRepo movie.Repository,
	metadataService *movie.MetadataService,
	logger *zap.Logger,
) *MovieMetadataRefreshWorker {
	return &MovieMetadataRefreshWorker{
		movieRepo:       movieRepo,
		metadataService: metadataService,
		logger:          logger,
	}
}

// Work executes the metadata refresh job
func (w *MovieMetadataRefreshWorker) Work(ctx context.Context, job *river.Job[MovieMetadataRefreshArgs]) error {
	args := job.Args

	w.logger.Info("starting movie metadata refresh",
		zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		zap.String("movie_id", args.MovieID.String()),
		zap.Bool("force", args.Force),
	)

	// Get existing movie
	existingMovie, err := w.movieRepo.GetMovie(ctx, args.MovieID)
	if err != nil {
		w.logger.Error("failed to get movie",
			zap.String("movie_id", args.MovieID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get movie: %w", err)
	}

	// Check if movie has TMDb ID
	if existingMovie.TMDbID == nil {
		w.logger.Warn("movie has no TMDb ID, cannot refresh",
			zap.String("movie_id", args.MovieID.String()),
		)
		return fmt.Errorf("movie has no TMDb ID")
	}

	// Clear cache if force refresh
	if args.Force {
		w.logger.Info("cleared metadata cache for force refresh",
			zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		)
	}

	// Enrich movie with fresh metadata (this also updates it in the database)
	if err := w.metadataService.EnrichMovie(ctx, existingMovie); err != nil {
		w.logger.Error("failed to enrich movie with metadata",
			zap.Int32("tmdb_id", *existingMovie.TMDbID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to enrich movie: %w", err)
	}

	// TODO: Refresh credits if needed (check if EnrichMovie already handles this)

	// TODO: Refresh credits (needs DeleteMovieCredits and CreateMovieCredits methods)
	w.logger.Info("movie metadata refresh completed",
		zap.String("movie_id", args.MovieID.String()),
		zap.String("title", existingMovie.Title),
	)

	return nil
}
