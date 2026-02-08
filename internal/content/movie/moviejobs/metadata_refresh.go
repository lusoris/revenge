package moviejobs

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
	"log/slog"

	"github.com/lusoris/revenge/internal/content/movie"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	metadatajobs "github.com/lusoris/revenge/internal/service/metadata/jobs"
)

// MovieMetadataRefreshWorker refreshes movie metadata using the movie service.
// It handles jobs of type metadatajobs.RefreshMovieArgs (kind: "metadata_refresh_movie").
type MovieMetadataRefreshWorker struct {
	river.WorkerDefaults[metadatajobs.RefreshMovieArgs]
	service   movie.Service
	jobClient *infrajobs.Client
	logger    *slog.Logger
}

// NewMovieMetadataRefreshWorker creates a new metadata refresh worker.
func NewMovieMetadataRefreshWorker(
	service movie.Service,
	jobClient *infrajobs.Client,
	logger *slog.Logger,
) *MovieMetadataRefreshWorker {
	return &MovieMetadataRefreshWorker{
		service:   service,
		jobClient: jobClient,
		logger:    logger,
	}
}

// Timeout returns the maximum execution time for movie metadata refresh jobs.
func (w *MovieMetadataRefreshWorker) Timeout(job *river.Job[metadatajobs.RefreshMovieArgs]) time.Duration {
	return 5 * time.Minute
}

// Work executes the metadata refresh job by delegating to the movie service.
func (w *MovieMetadataRefreshWorker) Work(ctx context.Context, job *river.Job[metadatajobs.RefreshMovieArgs]) error {
	args := job.Args

	w.logger.Info("starting movie metadata refresh",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("movie_id", args.MovieID.String()),
		slog.Bool("force", args.Force),
		slog.Any("languages", args.Languages),
	)

	_ = w.jobClient.ReportProgress(ctx, job.ID, &infrajobs.JobProgress{
		Phase:   "refreshing",
		Current: 0,
		Total:   1,
		Message: "refreshing movie metadata",
	})

	opts := movie.MetadataRefreshOptions{
		Force:     args.Force,
		Languages: args.Languages,
	}

	if err := w.service.RefreshMovieMetadata(ctx, args.MovieID, opts); err != nil {
		w.logger.Error("movie metadata refresh failed",
			slog.String("movie_id", args.MovieID.String()),
			slog.Any("error",err),
		)
		return fmt.Errorf("movie metadata refresh failed: %w", err)
	}

	_ = w.jobClient.ReportProgress(ctx, job.ID, &infrajobs.JobProgress{
		Phase:   "completed",
		Current: 1,
		Total:   1,
		Percent: 100,
		Message: "metadata refresh completed",
	})

	w.logger.Info("movie metadata refresh completed",
		slog.String("movie_id", args.MovieID.String()),
	)

	return nil
}
