package moviejobs

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
	"go.uber.org/zap"

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
	logger    *zap.Logger
}

// NewMovieMetadataRefreshWorker creates a new metadata refresh worker.
func NewMovieMetadataRefreshWorker(
	service movie.Service,
	jobClient *infrajobs.Client,
	logger *zap.Logger,
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
		zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		zap.String("movie_id", args.MovieID.String()),
		zap.Bool("force", args.Force),
		zap.Strings("languages", args.Languages),
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
			zap.String("movie_id", args.MovieID.String()),
			zap.Error(err),
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
		zap.String("movie_id", args.MovieID.String()),
	)

	return nil
}
