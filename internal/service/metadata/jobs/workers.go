package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/image"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/riverqueue/river"
)

// =============================================================================
// TV Show Metadata Refresh Worker
// =============================================================================

// RefreshTVShowWorker handles metadata_refresh_tvshow jobs by delegating
// to the tvshow service's RefreshSeriesMetadata method.
type RefreshTVShowWorker struct {
	river.WorkerDefaults[RefreshTVShowArgs]
	service   tvshow.Service
	jobClient *infrajobs.Client
	logger    *slog.Logger
}

// NewRefreshTVShowWorker creates a new TV show metadata refresh worker.
func NewRefreshTVShowWorker(service tvshow.Service, jobClient *infrajobs.Client, logger *slog.Logger) *RefreshTVShowWorker {
	return &RefreshTVShowWorker{
		service:   service,
		jobClient: jobClient,
		logger:    logger.With("component", "metadata_refresh_tvshow"),
	}
}

// Timeout returns the maximum execution time for TV show metadata refresh jobs.
func (w *RefreshTVShowWorker) Timeout(job *river.Job[RefreshTVShowArgs]) time.Duration {
	return 10 * time.Minute
}

// Work executes the TV show metadata refresh job.
func (w *RefreshTVShowWorker) Work(ctx context.Context, job *river.Job[RefreshTVShowArgs]) error {
	args := job.Args

	w.logger.Info("starting tvshow metadata refresh",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("series_id", args.SeriesID.String()),
		slog.Bool("force", args.Force),
		slog.Bool("include_seasons", args.IncludeSeasons),
		slog.Bool("include_episodes", args.IncludeEpisodes),
	)

	opts := tvshow.MetadataRefreshOptions{
		Force:     args.Force,
		Languages: args.Languages,
	}

	if err := w.service.RefreshSeriesMetadata(ctx, args.SeriesID, opts); err != nil {
		w.logger.Error("tvshow metadata refresh failed",
			slog.String("series_id", args.SeriesID.String()),
			slog.Any("error", err),
		)
		return fmt.Errorf("tvshow metadata refresh failed: %w", err)
	}

	w.logger.Info("tvshow metadata refresh completed",
		slog.String("series_id", args.SeriesID.String()),
	)
	return nil
}

// =============================================================================
// Season Metadata Refresh Worker
// =============================================================================

// RefreshSeasonWorker handles metadata_refresh_season jobs.
type RefreshSeasonWorker struct {
	river.WorkerDefaults[RefreshSeasonArgs]
	service   tvshow.Service
	jobClient *infrajobs.Client
	logger    *slog.Logger
}

// NewRefreshSeasonWorker creates a new season metadata refresh worker.
func NewRefreshSeasonWorker(service tvshow.Service, jobClient *infrajobs.Client, logger *slog.Logger) *RefreshSeasonWorker {
	return &RefreshSeasonWorker{
		service:   service,
		jobClient: jobClient,
		logger:    logger.With("component", "metadata_refresh_season"),
	}
}

// Timeout returns the maximum execution time for season metadata refresh jobs.
func (w *RefreshSeasonWorker) Timeout(job *river.Job[RefreshSeasonArgs]) time.Duration {
	return 5 * time.Minute
}

// Work executes the season metadata refresh job.
func (w *RefreshSeasonWorker) Work(ctx context.Context, job *river.Job[RefreshSeasonArgs]) error {
	args := job.Args

	w.logger.Info("starting season metadata refresh",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("series_id", args.SeriesID.String()),
		slog.String("season_id", args.SeasonID.String()),
		slog.Int("season_number", args.SeasonNumber),
	)

	opts := tvshow.MetadataRefreshOptions{
		Force:     args.Force,
		Languages: args.Languages,
	}

	if err := w.service.RefreshSeasonMetadata(ctx, args.SeasonID, opts); err != nil {
		w.logger.Error("season metadata refresh failed",
			slog.String("season_id", args.SeasonID.String()),
			slog.Any("error", err),
		)
		return fmt.Errorf("season metadata refresh failed: %w", err)
	}

	w.logger.Info("season metadata refresh completed",
		slog.String("season_id", args.SeasonID.String()),
	)
	return nil
}

// =============================================================================
// Episode Metadata Refresh Worker
// =============================================================================

// RefreshEpisodeWorker handles metadata_refresh_episode jobs.
type RefreshEpisodeWorker struct {
	river.WorkerDefaults[RefreshEpisodeArgs]
	service   tvshow.Service
	jobClient *infrajobs.Client
	logger    *slog.Logger
}

// NewRefreshEpisodeWorker creates a new episode metadata refresh worker.
func NewRefreshEpisodeWorker(service tvshow.Service, jobClient *infrajobs.Client, logger *slog.Logger) *RefreshEpisodeWorker {
	return &RefreshEpisodeWorker{
		service:   service,
		jobClient: jobClient,
		logger:    logger.With("component", "metadata_refresh_episode"),
	}
}

// Timeout returns the maximum execution time for episode metadata refresh jobs.
func (w *RefreshEpisodeWorker) Timeout(job *river.Job[RefreshEpisodeArgs]) time.Duration {
	return 5 * time.Minute
}

// Work executes the episode metadata refresh job.
func (w *RefreshEpisodeWorker) Work(ctx context.Context, job *river.Job[RefreshEpisodeArgs]) error {
	args := job.Args

	w.logger.Info("starting episode metadata refresh",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("episode_id", args.EpisodeID.String()),
		slog.Int("season_number", args.SeasonNumber),
		slog.Int("episode_number", args.EpisodeNumber),
	)

	opts := tvshow.MetadataRefreshOptions{
		Force:     args.Force,
		Languages: args.Languages,
	}

	if err := w.service.RefreshEpisodeMetadata(ctx, args.EpisodeID, opts); err != nil {
		w.logger.Error("episode metadata refresh failed",
			slog.String("episode_id", args.EpisodeID.String()),
			slog.Any("error", err),
		)
		return fmt.Errorf("episode metadata refresh failed: %w", err)
	}

	w.logger.Info("episode metadata refresh completed",
		slog.String("episode_id", args.EpisodeID.String()),
	)
	return nil
}

// =============================================================================
// Person Metadata Refresh Worker (stub — no person service yet)
// =============================================================================

// RefreshPersonWorker handles metadata_refresh_person jobs.
// Currently a stub that logs and succeeds — will be implemented when
// the person service is created.
type RefreshPersonWorker struct {
	river.WorkerDefaults[RefreshPersonArgs]
	logger *slog.Logger
}

// NewRefreshPersonWorker creates a new person metadata refresh worker.
func NewRefreshPersonWorker(logger *slog.Logger) *RefreshPersonWorker {
	return &RefreshPersonWorker{
		logger: logger.With("component", "metadata_refresh_person"),
	}
}

// Timeout returns the maximum execution time for person metadata refresh jobs.
func (w *RefreshPersonWorker) Timeout(job *river.Job[RefreshPersonArgs]) time.Duration {
	return 2 * time.Minute
}

// Work executes the person metadata refresh job.
// Person service is not yet implemented — cancel the job so it doesn't retry.
func (w *RefreshPersonWorker) Work(ctx context.Context, job *river.Job[RefreshPersonArgs]) error {
	w.logger.Warn("person metadata refresh not yet implemented — cancelling job",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("provider_id", job.Args.ProviderID),
	)
	return river.JobCancel(fmt.Errorf("person metadata refresh not yet implemented"))
}

// =============================================================================
// Content Enrichment Worker
// =============================================================================

// EnrichContentWorker handles metadata_enrich_content jobs by delegating
// to the metadata service's enrichment methods.
type EnrichContentWorker struct {
	river.WorkerDefaults[EnrichContentArgs]
	metadataService metadata.Service
	movieService    movie.Service
	tvshowService   tvshow.Service
	logger          *slog.Logger
}

// NewEnrichContentWorker creates a new content enrichment worker.
func NewEnrichContentWorker(
	metadataService metadata.Service,
	movieService movie.Service,
	tvshowService tvshow.Service,
	logger *slog.Logger,
) *EnrichContentWorker {
	return &EnrichContentWorker{
		metadataService: metadataService,
		movieService:    movieService,
		tvshowService:   tvshowService,
		logger:          logger.With("component", "metadata_enrich_content"),
	}
}

// Timeout returns the maximum execution time for content enrichment jobs.
func (w *EnrichContentWorker) Timeout(job *river.Job[EnrichContentArgs]) time.Duration {
	return 5 * time.Minute
}

// Work executes the content enrichment job.
func (w *EnrichContentWorker) Work(ctx context.Context, job *river.Job[EnrichContentArgs]) error {
	args := job.Args

	w.logger.Info("starting content enrichment",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("content_type", args.ContentType),
		slog.String("content_id", args.ContentID.String()),
		slog.Any("providers", args.Providers),
	)

	switch args.ContentType {
	case "movie":
		opts := movie.MetadataRefreshOptions{
			Force:     true,
			Languages: args.Languages,
		}
		if err := w.movieService.RefreshMovieMetadata(ctx, args.ContentID, opts); err != nil {
			return fmt.Errorf("movie enrichment failed: %w", err)
		}
	case "tvshow":
		opts := tvshow.MetadataRefreshOptions{
			Force:     true,
			Languages: args.Languages,
		}
		if err := w.tvshowService.RefreshSeriesMetadata(ctx, args.ContentID, opts); err != nil {
			return fmt.Errorf("tvshow enrichment failed: %w", err)
		}
	default:
		w.logger.Warn("unknown content type for enrichment",
			slog.String("content_type", args.ContentType),
		)
	}

	w.logger.Info("content enrichment completed",
		slog.String("content_type", args.ContentType),
		slog.String("content_id", args.ContentID.String()),
	)
	return nil
}

// =============================================================================
// Image Download Worker
// =============================================================================

// DownloadImageWorker handles metadata_download_image jobs.
// It downloads images via the image service and stores them locally.
type DownloadImageWorker struct {
	river.WorkerDefaults[DownloadImageArgs]
	imageService *image.Service
	logger       *slog.Logger
}

// NewDownloadImageWorker creates a new image download worker.
func NewDownloadImageWorker(imageService *image.Service, logger *slog.Logger) *DownloadImageWorker {
	return &DownloadImageWorker{
		imageService: imageService,
		logger:       logger.With("component", "metadata_download_image"),
	}
}

// Timeout returns the maximum execution time for image download jobs.
func (w *DownloadImageWorker) Timeout(job *river.Job[DownloadImageArgs]) time.Duration {
	return 2 * time.Minute
}

// Work executes the image download job by fetching the image into the local cache.
func (w *DownloadImageWorker) Work(ctx context.Context, job *river.Job[DownloadImageArgs]) error {
	args := job.Args

	if args.Path == "" {
		w.logger.Debug("skipping image download — empty path",
			slog.String("content_type", args.ContentType),
			slog.String("content_id", args.ContentID),
		)
		return nil
	}

	w.logger.Debug("downloading image",
		slog.String("content_type", args.ContentType),
		slog.String("image_type", args.ImageType),
		slog.String("path", args.Path),
		slog.String("size", args.Size),
	)

	_, _, err := w.imageService.FetchImage(ctx, args.ImageType, args.Path, args.Size)
	if err != nil {
		return fmt.Errorf("download image %s/%s: %w", args.ImageType, args.Path, err)
	}

	return nil
}
