package moviejobs

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/service/search"
)

// SearchIndexOperation defines the type of search index operation.
type SearchIndexOperation string

const (
	// SearchIndexOperationIndex indexes or updates a movie in the search index.
	SearchIndexOperationIndex SearchIndexOperation = "index"
	// SearchIndexOperationRemove removes a movie from the search index.
	SearchIndexOperationRemove SearchIndexOperation = "remove"
	// SearchIndexOperationReindex triggers a full reindex of all movies.
	SearchIndexOperationReindex SearchIndexOperation = "reindex"
)

// MovieSearchIndexArgs are the arguments for search index operations.
type MovieSearchIndexArgs struct {
	// Operation is the type of index operation to perform.
	Operation SearchIndexOperation `json:"operation"`
	// MovieID is the ID of the movie to index/remove (not used for reindex).
	MovieID uuid.UUID `json:"movie_id,omitempty"`
}

// Kind returns the unique job kind for River.
func (MovieSearchIndexArgs) Kind() string {
	return "movie_search_index"
}

// MovieSearchIndexWorker handles search index operations for movies.
type MovieSearchIndexWorker struct {
	river.WorkerDefaults[MovieSearchIndexArgs]
	movieRepo     movie.Repository
	searchService *search.MovieSearchService
	logger        *zap.Logger
}

// NewMovieSearchIndexWorker creates a new search index worker.
func NewMovieSearchIndexWorker(
	movieRepo movie.Repository,
	searchService *search.MovieSearchService,
	logger *zap.Logger,
) *MovieSearchIndexWorker {
	return &MovieSearchIndexWorker{
		movieRepo:     movieRepo,
		searchService: searchService,
		logger:        logger.Named("search_index_worker"),
	}
}

// Work executes the search index job.
func (w *MovieSearchIndexWorker) Work(ctx context.Context, job *river.Job[MovieSearchIndexArgs]) error {
	args := job.Args

	w.logger.Info("starting search index operation",
		zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		zap.String("operation", string(args.Operation)),
		zap.String("movie_id", args.MovieID.String()),
	)

	// Check if search is enabled
	if !w.searchService.IsEnabled() {
		w.logger.Debug("search is disabled, skipping index operation")
		return nil
	}

	switch args.Operation {
	case SearchIndexOperationIndex:
		return w.indexMovie(ctx, args.MovieID)
	case SearchIndexOperationRemove:
		return w.removeMovie(ctx, args.MovieID)
	case SearchIndexOperationReindex:
		return w.reindexAll(ctx)
	default:
		return fmt.Errorf("unknown operation: %s", args.Operation)
	}
}

// indexMovie indexes a single movie in the search engine.
func (w *MovieSearchIndexWorker) indexMovie(ctx context.Context, movieID uuid.UUID) error {
	// Get the movie
	m, err := w.movieRepo.GetMovie(ctx, movieID)
	if err != nil {
		if err == movie.ErrMovieNotFound {
			w.logger.Warn("movie not found, skipping index",
				zap.String("movie_id", movieID.String()),
			)
			return nil // Don't fail the job, movie might have been deleted
		}
		return fmt.Errorf("failed to get movie: %w", err)
	}

	// Get related data
	genres, err := w.movieRepo.ListMovieGenres(ctx, movieID)
	if err != nil {
		w.logger.Warn("failed to get genres", zap.Error(err))
		genres = nil
	}

	cast, err := w.movieRepo.ListMovieCast(ctx, movieID)
	if err != nil {
		w.logger.Warn("failed to get cast", zap.Error(err))
		cast = nil
	}

	crew, err := w.movieRepo.ListMovieCrew(ctx, movieID)
	if err != nil {
		w.logger.Warn("failed to get crew", zap.Error(err))
		crew = nil
	}

	// Combine cast and crew
	credits := append(cast, crew...)

	// Get primary file
	var file *movie.MovieFile
	files, err := w.movieRepo.ListMovieFilesByMovieID(ctx, movieID)
	if err != nil {
		w.logger.Warn("failed to get files", zap.Error(err))
	} else if len(files) > 0 {
		file = &files[0]
	}

	// Index the movie
	if err := w.searchService.UpdateMovie(ctx, m, genres, credits, file); err != nil {
		return fmt.Errorf("failed to index movie: %w", err)
	}

	w.logger.Info("movie indexed successfully",
		zap.String("movie_id", movieID.String()),
		zap.String("title", m.Title),
	)

	return nil
}

// removeMovie removes a movie from the search index.
func (w *MovieSearchIndexWorker) removeMovie(ctx context.Context, movieID uuid.UUID) error {
	if err := w.searchService.RemoveMovie(ctx, movieID); err != nil {
		return fmt.Errorf("failed to remove movie from index: %w", err)
	}

	w.logger.Info("movie removed from index",
		zap.String("movie_id", movieID.String()),
	)

	return nil
}

// reindexAll triggers a full reindex of all movies.
func (w *MovieSearchIndexWorker) reindexAll(ctx context.Context) error {
	w.logger.Info("starting full search reindex")

	if err := w.searchService.ReindexAll(ctx, w.movieRepo); err != nil {
		return fmt.Errorf("failed to reindex all movies: %w", err)
	}

	w.logger.Info("full search reindex completed")

	return nil
}
