package moviejobs

import (
	"context"
	"errors"
	"time"

	"github.com/riverqueue/river"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/content/movie"
)

const MovieFileMatchJobKind = "movie_file_match"

// MovieFileMatchArgs are the arguments for the movie file match job.
type MovieFileMatchArgs struct {
	FilePath     string `json:"file_path"`
	ForceRematch bool   `json:"force_rematch"`
}

// Kind returns the job kind for the movie file match job.
func (MovieFileMatchArgs) Kind() string {
	return MovieFileMatchJobKind
}

// MovieFileMatchWorker is a worker that matches movie files to movies.
type MovieFileMatchWorker struct {
	river.WorkerDefaults[MovieFileMatchArgs]
	libraryService *movie.LibraryService
	logger         *zap.Logger
}

// NewMovieFileMatchWorker creates a new movie file match worker.
func NewMovieFileMatchWorker(
	libraryService *movie.LibraryService,
	logger *zap.Logger,
) *MovieFileMatchWorker {
	return &MovieFileMatchWorker{
		libraryService: libraryService,
		logger:         logger,
	}
}

// Kind returns the job kind.
func (w *MovieFileMatchWorker) Kind() string {
	return MovieFileMatchJobKind
}

// Timeout returns the maximum execution time for movie file match jobs.
func (w *MovieFileMatchWorker) Timeout(job *river.Job[MovieFileMatchArgs]) time.Duration {
	return 5 * time.Minute
}

// Work performs the movie file match job.
func (w *MovieFileMatchWorker) Work(ctx context.Context, job *river.Job[MovieFileMatchArgs]) error {
	args := job.Args

	w.logger.Info("starting movie file match",
		zap.String("file_path", args.FilePath),
		zap.Bool("force_rematch", args.ForceRematch),
	)

	// Match the file using the library service
	result, err := w.libraryService.MatchFile(ctx, args.FilePath, args.ForceRematch)
	if err != nil {
		w.logger.Error("file match failed",
			zap.String("file_path", args.FilePath),
			zap.Error(err),
		)
		return err
	}

	// Check for match errors
	if result.Error != nil {
		w.logger.Warn("file matched with warnings",
			zap.String("file_path", args.FilePath),
			zap.Error(result.Error),
		)
	}

	// Log match result
	if result.Movie != nil {
		w.logger.Info("movie file matched successfully",
			zap.String("file_path", args.FilePath),
			zap.String("movie_id", result.Movie.ID.String()),
			zap.String("movie_title", result.Movie.Title),
			zap.String("match_type", string(result.MatchType)),
			zap.Float64("confidence", result.Confidence),
			zap.Bool("created_new_movie", result.CreatedNewMovie),
		)
	} else {
		w.logger.Warn("file could not be matched",
			zap.String("file_path", args.FilePath),
			zap.String("match_type", string(result.MatchType)),
		)
		return errors.New("file could not be matched to any movie")
	}

	return nil
}
