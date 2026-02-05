// Package jobs provides background job workers for TV show content processing.
// It uses River for job queueing and processing with support for library scanning,
// metadata refresh, file matching, and search indexing.
package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	sharedjobs "github.com/lusoris/revenge/internal/content/shared/jobs"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// Job kind constants for TV show jobs
const (
	KindLibraryScan     = "tvshow_library_scan"
	KindMetadataRefresh = "tvshow_metadata_refresh"
	KindFileMatch       = "tvshow_file_match"
	KindSearchIndex     = "tvshow_search_index"
	KindSeriesRefresh   = "tvshow_series_refresh"
	KindSeasonRefresh   = "tvshow_season_refresh"
	KindEpisodeRefresh  = "tvshow_episode_refresh"
)

// =============================================================================
// Library Scan Job
// =============================================================================

// LibraryScanArgs defines arguments for TV show library scan jobs.
type LibraryScanArgs struct {
	// Paths are the library paths to scan for TV shows.
	Paths []string `json:"paths"`

	// Force indicates whether to force a full rescan, ignoring last scan time.
	Force bool `json:"force"`

	// LibraryID is the optional library ID to scan.
	LibraryID *uuid.UUID `json:"library_id,omitempty"`
}

// Kind returns the job kind identifier.
func (LibraryScanArgs) Kind() string {
	return KindLibraryScan
}

// LibraryScanWorker scans TV show library directories.
type LibraryScanWorker struct {
	river.WorkerDefaults[LibraryScanArgs]
	service tvshow.Service
	logger  *zap.Logger
}

// NewLibraryScanWorker creates a new library scan worker.
func NewLibraryScanWorker(service tvshow.Service, logger *zap.Logger) *LibraryScanWorker {
	return &LibraryScanWorker{
		service: service,
		logger:  logger.Named("tvshow_library_scan"),
	}
}

// Work executes the library scan job.
func (w *LibraryScanWorker) Work(ctx context.Context, job *river.Job[LibraryScanArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindLibraryScan)
	jctx.LogStart(
		zap.Strings("paths", job.Args.Paths),
		zap.Bool("force", job.Args.Force),
	)

	result := &sharedjobs.JobResult{Success: true}
	start := time.Now()

	// TODO: Implement actual library scanning logic
	// 1. Walk directories to find video files
	// 2. Parse filenames using TVShowFileParser
	// 3. Queue FileMatch jobs for new files
	// 4. Update library scan timestamps

	for _, path := range job.Args.Paths {
		w.logger.Info("scanning path",
			zap.Int64("job_id", job.ID),
			zap.String("path", path),
		)
		// Placeholder: actual scanning logic would go here
		result.ItemsProcessed++
	}

	result.Duration = time.Since(start)
	result.LogSummary(w.logger, KindLibraryScan)

	jctx.LogComplete(
		zap.Int("paths_scanned", len(job.Args.Paths)),
		zap.Int("items_processed", result.ItemsProcessed),
	)

	return nil
}

// =============================================================================
// Metadata Refresh Job
// =============================================================================

// MetadataRefreshArgs defines arguments for TV show metadata refresh jobs.
type MetadataRefreshArgs struct {
	// SeriesID is the series to refresh. If nil, refreshes all series.
	SeriesID *uuid.UUID `json:"series_id,omitempty"`

	// SeasonID is the specific season to refresh.
	SeasonID *uuid.UUID `json:"season_id,omitempty"`

	// EpisodeID is the specific episode to refresh.
	EpisodeID *uuid.UUID `json:"episode_id,omitempty"`

	// Force indicates whether to force a refresh even if recently updated.
	Force bool `json:"force"`

	// RefreshImages indicates whether to also refresh images.
	RefreshImages bool `json:"refresh_images"`
}

// Kind returns the job kind identifier.
func (MetadataRefreshArgs) Kind() string {
	return KindMetadataRefresh
}

// MetadataRefreshWorker refreshes TV show metadata from external sources.
type MetadataRefreshWorker struct {
	river.WorkerDefaults[MetadataRefreshArgs]
	service tvshow.Service
	logger  *zap.Logger
}

// NewMetadataRefreshWorker creates a new metadata refresh worker.
func NewMetadataRefreshWorker(service tvshow.Service, logger *zap.Logger) *MetadataRefreshWorker {
	return &MetadataRefreshWorker{
		service: service,
		logger:  logger.Named("tvshow_metadata_refresh"),
	}
}

// Work executes the metadata refresh job.
func (w *MetadataRefreshWorker) Work(ctx context.Context, job *river.Job[MetadataRefreshArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindMetadataRefresh)
	args := job.Args

	jctx.LogStart(
		zap.Any("series_id", args.SeriesID),
		zap.Any("season_id", args.SeasonID),
		zap.Any("episode_id", args.EpisodeID),
		zap.Bool("force", args.Force),
		zap.Bool("refresh_images", args.RefreshImages),
	)

	result := &sharedjobs.JobResult{Success: true}
	start := time.Now()

	// Determine what to refresh
	switch {
	case args.EpisodeID != nil:
		if err := w.service.RefreshEpisodeMetadata(ctx, *args.EpisodeID); err != nil {
			result.AddError(fmt.Errorf("refresh episode %s: %w", args.EpisodeID, err))
		} else {
			result.ItemsProcessed++
		}

	case args.SeasonID != nil:
		if err := w.service.RefreshSeasonMetadata(ctx, *args.SeasonID); err != nil {
			result.AddError(fmt.Errorf("refresh season %s: %w", args.SeasonID, err))
		} else {
			result.ItemsProcessed++
		}

	case args.SeriesID != nil:
		if err := w.service.RefreshSeriesMetadata(ctx, *args.SeriesID); err != nil {
			result.AddError(fmt.Errorf("refresh series %s: %w", args.SeriesID, err))
		} else {
			result.ItemsProcessed++
		}

	default:
		// Refresh all series (batch operation)
		// TODO: Implement batch refresh with pagination
		w.logger.Info("batch metadata refresh not implemented yet")
	}

	result.Duration = time.Since(start)
	result.Success = !result.HasErrors()
	result.LogSummary(w.logger, KindMetadataRefresh)

	if result.HasErrors() {
		result.LogErrors(w.logger, 10)
		return fmt.Errorf("metadata refresh completed with errors: %d failed", result.ItemsFailed)
	}

	jctx.LogComplete(zap.Int("items_refreshed", result.ItemsProcessed))
	return nil
}

// =============================================================================
// File Match Job
// =============================================================================

// FileMatchArgs defines arguments for TV show file matching jobs.
type FileMatchArgs struct {
	// FilePath is the path to the file to match.
	FilePath string `json:"file_path"`

	// EpisodeID is set if the file should be matched to a specific episode.
	EpisodeID *uuid.UUID `json:"episode_id,omitempty"`

	// ForceRematch indicates whether to rematch even if already matched.
	ForceRematch bool `json:"force_rematch"`

	// AutoCreate indicates whether to create series/season/episode if not found.
	AutoCreate bool `json:"auto_create"`
}

// Kind returns the job kind identifier.
func (FileMatchArgs) Kind() string {
	return KindFileMatch
}

// FileMatchWorker matches scanned files to TV show episodes.
type FileMatchWorker struct {
	river.WorkerDefaults[FileMatchArgs]
	service tvshow.Service
	logger  *zap.Logger
}

// NewFileMatchWorker creates a new file match worker.
func NewFileMatchWorker(service tvshow.Service, logger *zap.Logger) *FileMatchWorker {
	return &FileMatchWorker{
		service: service,
		logger:  logger.Named("tvshow_file_match"),
	}
}

// Work executes the file match job.
func (w *FileMatchWorker) Work(ctx context.Context, job *river.Job[FileMatchArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindFileMatch)
	args := job.Args

	jctx.LogStart(
		zap.String("file_path", args.FilePath),
		zap.Any("episode_id", args.EpisodeID),
		zap.Bool("force_rematch", args.ForceRematch),
		zap.Bool("auto_create", args.AutoCreate),
	)

	// Check if file is already matched
	existingFile, err := w.service.GetEpisodeFileByPath(ctx, args.FilePath)
	if err == nil && existingFile != nil && !args.ForceRematch {
		w.logger.Info("file already matched, skipping",
			zap.Int64("job_id", job.ID),
			zap.String("file_path", args.FilePath),
			zap.String("episode_id", existingFile.EpisodeID.String()),
		)
		return nil
	}

	// TODO: Implement file matching logic
	// 1. Parse filename to extract series title, season, episode
	// 2. Search for matching series in database
	// 3. If auto_create and not found, search TMDb and create series
	// 4. Find or create season
	// 5. Find or create episode
	// 6. Create episode file record

	w.logger.Info("file matching not fully implemented yet",
		zap.Int64("job_id", job.ID),
		zap.String("file_path", args.FilePath),
	)

	jctx.LogComplete()
	return nil
}

// =============================================================================
// Search Index Job
// =============================================================================

// SearchIndexArgs defines arguments for TV show search indexing jobs.
type SearchIndexArgs struct {
	// SeriesID is the specific series to index. If nil, indexes all.
	SeriesID *uuid.UUID `json:"series_id,omitempty"`

	// FullReindex indicates whether to do a full reindex.
	FullReindex bool `json:"full_reindex"`
}

// Kind returns the job kind identifier.
func (SearchIndexArgs) Kind() string {
	return KindSearchIndex
}

// SearchIndexWorker indexes TV show content for search.
type SearchIndexWorker struct {
	river.WorkerDefaults[SearchIndexArgs]
	service tvshow.Service
	logger  *zap.Logger
}

// NewSearchIndexWorker creates a new search index worker.
func NewSearchIndexWorker(service tvshow.Service, logger *zap.Logger) *SearchIndexWorker {
	return &SearchIndexWorker{
		service: service,
		logger:  logger.Named("tvshow_search_index"),
	}
}

// Work executes the search index job.
func (w *SearchIndexWorker) Work(ctx context.Context, job *river.Job[SearchIndexArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindSearchIndex)
	args := job.Args

	jctx.LogStart(
		zap.Any("series_id", args.SeriesID),
		zap.Bool("full_reindex", args.FullReindex),
	)

	result := &sharedjobs.JobResult{Success: true}
	start := time.Now()

	if args.SeriesID != nil {
		// Index specific series
		series, err := w.service.GetSeries(ctx, *args.SeriesID)
		if err != nil {
			result.AddError(fmt.Errorf("get series %s: %w", args.SeriesID, err))
		} else {
			// TODO: Index to Typesense
			w.logger.Info("indexing series",
				zap.Int64("job_id", job.ID),
				zap.String("series_id", series.ID.String()),
				zap.String("title", series.Title),
			)
			result.ItemsProcessed++
		}
	} else {
		// Full reindex
		// TODO: Implement batch indexing with pagination
		w.logger.Info("full search reindex not implemented yet",
			zap.Int64("job_id", job.ID),
		)
	}

	result.Duration = time.Since(start)
	result.Success = !result.HasErrors()
	result.LogSummary(w.logger, KindSearchIndex)

	if result.HasErrors() {
		result.LogErrors(w.logger, 10)
		return fmt.Errorf("search indexing completed with errors: %d failed", result.ItemsFailed)
	}

	jctx.LogComplete(zap.Int("items_indexed", result.ItemsProcessed))
	return nil
}

// =============================================================================
// Series Refresh Job (Individual Series)
// =============================================================================

// SeriesRefreshArgs defines arguments for refreshing a single series.
type SeriesRefreshArgs struct {
	// SeriesID is the series to refresh.
	SeriesID uuid.UUID `json:"series_id"`

	// TMDbID is the TMDb ID to use for lookup.
	TMDbID int32 `json:"tmdb_id,omitempty"`

	// RefreshSeasons indicates whether to also refresh all seasons.
	RefreshSeasons bool `json:"refresh_seasons"`

	// RefreshEpisodes indicates whether to also refresh all episodes.
	RefreshEpisodes bool `json:"refresh_episodes"`

	// Languages specifies which languages to fetch metadata for.
	Languages []string `json:"languages,omitempty"`
}

// Kind returns the job kind identifier.
func (SeriesRefreshArgs) Kind() string {
	return KindSeriesRefresh
}

// SeriesRefreshWorker refreshes a single series from TMDb.
type SeriesRefreshWorker struct {
	river.WorkerDefaults[SeriesRefreshArgs]
	service tvshow.Service
	logger  *zap.Logger
}

// NewSeriesRefreshWorker creates a new series refresh worker.
func NewSeriesRefreshWorker(service tvshow.Service, logger *zap.Logger) *SeriesRefreshWorker {
	return &SeriesRefreshWorker{
		service: service,
		logger:  logger.Named("tvshow_series_refresh"),
	}
}

// Work executes the series refresh job.
func (w *SeriesRefreshWorker) Work(ctx context.Context, job *river.Job[SeriesRefreshArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindSeriesRefresh)
	args := job.Args

	jctx.LogStart(
		zap.String("series_id", args.SeriesID.String()),
		zap.Int32("tmdb_id", args.TMDbID),
		zap.Bool("refresh_seasons", args.RefreshSeasons),
		zap.Bool("refresh_episodes", args.RefreshEpisodes),
	)

	// TODO: Implement TMDb API call and update logic
	// 1. Fetch series details from TMDb
	// 2. Update series record in database
	// 3. If RefreshSeasons, queue SeasonRefresh jobs for each season
	// 4. Update metadata_updated_at timestamp

	w.logger.Info("series refresh not fully implemented yet",
		zap.Int64("job_id", job.ID),
		zap.String("series_id", args.SeriesID.String()),
	)

	jctx.LogComplete()
	return nil
}

// =============================================================================
// Job Helpers
// =============================================================================

// JobInsertOpts returns common job insert options for TV show jobs.
func JobInsertOpts(priority int, scheduledAt *time.Time) *river.InsertOpts {
	opts := &river.InsertOpts{
		Priority: priority,
	}
	if scheduledAt != nil {
		opts.ScheduledAt = *scheduledAt
	}
	return opts
}

// DefaultPriority is the default priority for TV show jobs.
const DefaultPriority = 2

// HighPriority is for urgent jobs like user-triggered refreshes.
const HighPriority = 1

// LowPriority is for batch/background jobs.
const LowPriority = 3
