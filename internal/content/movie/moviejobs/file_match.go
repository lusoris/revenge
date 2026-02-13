package moviejobs

import (
	"context"
	"time"

	"log/slog"

	"github.com/riverqueue/river"

	"github.com/lusoris/revenge/internal/content/movie"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
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

// InsertOpts returns the default insert options for movie file match jobs.
// File matching is deterministic — limit retries to avoid wasting resources.
func (MovieFileMatchArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 3,
	}
}

// MovieFileMatchWorker is a worker that matches movie files to movies.
type MovieFileMatchWorker struct {
	river.WorkerDefaults[MovieFileMatchArgs]
	libraryService *movie.LibraryService
	logger         *slog.Logger
}

// NewMovieFileMatchWorker creates a new movie file match worker.
func NewMovieFileMatchWorker(
	libraryService *movie.LibraryService,
	logger *slog.Logger,
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
		slog.String("file_path", args.FilePath),
		slog.Bool("force_rematch", args.ForceRematch),
	)

	// Match the file using the library service
	result, err := w.libraryService.MatchFile(ctx, args.FilePath, args.ForceRematch)
	if err != nil {
		w.logger.Error("file match failed",
			slog.String("file_path", args.FilePath),
			slog.Any("error", err),
		)
		return err
	}

	// Check for match errors
	if result.Error != nil {
		w.logger.Warn("file matched with warnings",
			slog.String("file_path", args.FilePath),
			slog.Any("error", result.Error),
		)
	}

	// Log match result
	if result.Movie != nil {
		w.logger.Info("movie file matched successfully",
			slog.String("file_path", args.FilePath),
			slog.String("movie_id", result.Movie.ID.String()),
			slog.String("movie_title", result.Movie.Title),
			slog.String("match_type", string(result.MatchType)),
			slog.Float64("confidence", result.Confidence),
			slog.Bool("created_new_movie", result.CreatedNewMovie),
		)
	} else {
		w.logger.Warn("file could not be matched — skipping (not a retryable error)",
			slog.String("file_path", args.FilePath),
			slog.String("match_type", string(result.MatchType)),
		)
		// Return nil: unmatched files are a valid outcome, not a job failure.
		// Retrying won't change the result since the file simply doesn't match.
		return nil
	}

	return nil
}
