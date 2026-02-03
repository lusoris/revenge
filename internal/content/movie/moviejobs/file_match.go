package moviejobs

import (
	"context"
	"errors"

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

// Work performs the movie file match job.
func (w *MovieFileMatchWorker) Work(ctx context.Context, job *river.Job[MovieFileMatchArgs]) error {
	args := job.Args

	w.logger.Info("starting movie file match",
		zap.String("file_path", args.FilePath),
		zap.Bool("force_rematch", args.ForceRematch),
	)

	// TODO: Implement once library.Service.MatchFile method is available
	w.logger.Error("movie file match not implemented",
		zap.String("file_path", args.FilePath),
		zap.String("reason", "library.Service.MatchFile method not available"),
	)

	return errors.New("movie file match not implemented: library.Service.MatchFile method not available")
}
