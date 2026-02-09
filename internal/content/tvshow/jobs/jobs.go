// Package jobs provides background job workers for TV show content processing.
// It uses River for job queueing and processing with support for library scanning,
// metadata refresh, file matching, and search indexing.
package jobs

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	sharedjobs "github.com/lusoris/revenge/internal/content/shared/jobs"
	"github.com/lusoris/revenge/internal/content/shared/scanner"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/content/tvshow/adapters"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/service/search"
	"github.com/riverqueue/river"
	"log/slog"
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

	// AutoCreate indicates whether to auto-create series/seasons/episodes for discovered files.
	AutoCreate bool `json:"auto_create"`
}

// Kind returns the job kind identifier.
func (LibraryScanArgs) Kind() string {
	return KindLibraryScan
}

// InsertOpts returns the default insert options for TV show library scan jobs.
// Library scans run on the bulk queue since they're resource-intensive batch operations.
func (LibraryScanArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: infrajobs.QueueBulk,
	}
}

// LibraryScanWorker scans TV show library directories.
type LibraryScanWorker struct {
	river.WorkerDefaults[LibraryScanArgs]
	service          tvshow.Service
	metadataProvider tvshow.MetadataProvider
	jobClient        *infrajobs.Client
	logger           *slog.Logger
}

// NewLibraryScanWorker creates a new library scan worker.
func NewLibraryScanWorker(service tvshow.Service, metadataProvider tvshow.MetadataProvider, jobClient *infrajobs.Client, logger *slog.Logger) *LibraryScanWorker {
	return &LibraryScanWorker{
		service:          service,
		metadataProvider: metadataProvider,
		jobClient:        jobClient,
		logger:           logger.With("component", "tvshow_library_scan"),
	}
}

// Timeout returns the maximum execution time for library scan jobs.
func (w *LibraryScanWorker) Timeout(job *river.Job[LibraryScanArgs]) time.Duration {
	return 30 * time.Minute
}

// Work executes the library scan job.
func (w *LibraryScanWorker) Work(ctx context.Context, job *river.Job[LibraryScanArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindLibraryScan)
	jctx.LogStart(
		slog.Any("paths", job.Args.Paths),
		slog.Bool("force", job.Args.Force),
		slog.Bool("auto_create", job.Args.AutoCreate),
	)

	result := &sharedjobs.JobResult{Success: true}
	start := time.Now()
	itemsSkipped := 0

	if len(job.Args.Paths) == 0 {
		w.logger.Warn("no paths provided for library scan")
		return nil
	}

	// Create the TV show file parser and scanner
	parser := adapters.NewTVShowFileParser()
	fsScanner := scanner.NewFilesystemScanner(job.Args.Paths, parser)

	// Scan all paths
	scanResults, summary, err := fsScanner.ScanWithSummary(ctx)
	if err != nil {
		result.AddError(fmt.Errorf("scan failed: %w", err))
		w.logger.Error("library scan failed", slog.Any("error", err))
		return err
	}

	w.logger.Info("scan completed",
		slog.Int("total_files", summary.TotalFiles),
		slog.Int("media_files", summary.MediaFiles),
		slog.Int("parsed_files", summary.ParsedFiles),
	)

	// Count media files for progress tracking
	mediaFiles := 0
	for _, sr := range scanResults {
		if sr.IsMedia {
			mediaFiles++
		}
	}
	processed := 0

	// Process each discovered file
	for _, sr := range scanResults {
		if !sr.IsMedia {
			continue
		}
		processed++

		// Report progress
		_ = w.jobClient.ReportProgress(ctx, job.ID, &infrajobs.JobProgress{
			Phase:   "processing",
			Current: processed,
			Total:   mediaFiles,
			Message: sr.ParsedTitle,
		})

		// Check if file is already matched
		existingFile, err := w.service.GetEpisodeFileByPath(ctx, sr.FilePath)
		if err == nil && existingFile != nil && !job.Args.Force {
			w.logger.Debug("file already matched, skipping",
				slog.String("file_path", sr.FilePath),
			)
			itemsSkipped++
			continue
		}

		// Process the file with auto-create if enabled
		if job.Args.AutoCreate && w.metadataProvider != nil {
			if err := w.processFile(ctx, sr); err != nil {
				w.logger.Warn("failed to process file",
					slog.String("file_path", sr.FilePath),
					slog.Any("error", err),
				)
				result.AddError(fmt.Errorf("process %s: %w", sr.FilePath, err))
				continue
			}
			result.ItemsProcessed++
		} else {
			// Just log discovered files when auto-create is disabled
			w.logger.Info("discovered tv show file",
				slog.String("file_path", sr.FilePath),
				slog.String("parsed_title", sr.ParsedTitle),
				slog.Any("season", sr.GetSeason()),
				slog.Any("episode", sr.GetEpisode()),
			)
			result.ItemsProcessed++
		}
	}

	result.Duration = time.Since(start)
	result.Success = !result.HasErrors()
	result.LogSummary(w.logger, KindLibraryScan)

	if result.HasErrors() {
		result.LogErrors(w.logger, 10)
	}

	jctx.LogComplete(
		slog.Int("paths_scanned", len(job.Args.Paths)),
		slog.Int("items_processed", result.ItemsProcessed),
		slog.Int("items_skipped", itemsSkipped),
	)

	return nil
}

// processFile processes a single scanned file, creating series/season/episode as needed.
func (w *LibraryScanWorker) processFile(ctx context.Context, sr scanner.ScanResult) error {
	// Extract metadata from scan result
	seriesTitle := sr.ParsedTitle
	if seriesTitle == "" {
		return fmt.Errorf("could not parse series title from filename")
	}

	seasonNum := sr.GetSeason()
	episodeNum := sr.GetEpisode()
	if seasonNum == nil || episodeNum == nil {
		return fmt.Errorf("could not parse season/episode from filename")
	}

	// Search for existing series by title
	seriesList, err := w.service.SearchSeries(ctx, seriesTitle, 5, 0)
	if err != nil {
		return fmt.Errorf("search series: %w", err)
	}

	var series *tvshow.Series

	// Check for exact match
	for _, s := range seriesList {
		if normalizeTitle(s.Title) == normalizeTitle(seriesTitle) {
			series = &s
			break
		}
	}

	// If no match found, search TMDb and create series
	if series == nil {
		searchResults, err := w.metadataProvider.SearchSeries(ctx, seriesTitle, nil)
		if err != nil || len(searchResults) == 0 {
			return fmt.Errorf("series not found: %s", seriesTitle)
		}

		// Use first result and enrich it
		newSeries := searchResults[0]
		if err := w.metadataProvider.EnrichSeries(ctx, newSeries); err != nil {
			w.logger.Warn("failed to enrich series", slog.Any("error", err))
		}

		// Create series
		params := seriesToCreateParams(newSeries)
		created, err := w.service.CreateSeries(ctx, params)
		if err != nil {
			return fmt.Errorf("create series: %w", err)
		}
		series = created
	}

	// Find or create season
	season, err := w.service.GetSeasonByNumber(ctx, series.ID, int32(*seasonNum))
	if err != nil {
		// Create season
		seasonParams := tvshow.CreateSeasonParams{
			SeriesID:     series.ID,
			SeasonNumber: int32(*seasonNum),
			Name:         fmt.Sprintf("Season %d", *seasonNum),
			EpisodeCount: 0,
		}

		// Try to enrich season from TMDb if series has TMDbID
		if series.TMDbID != nil {
			tmpSeason := &tvshow.Season{
				SeriesID:     series.ID,
				SeasonNumber: int32(*seasonNum),
			}
			if err := w.metadataProvider.EnrichSeason(ctx, tmpSeason, *series.TMDbID); err == nil {
				seasonParams.TMDbID = tmpSeason.TMDbID
				seasonParams.Name = tmpSeason.Name
				seasonParams.Overview = tmpSeason.Overview
				seasonParams.PosterPath = tmpSeason.PosterPath
				seasonParams.EpisodeCount = tmpSeason.EpisodeCount
				if tmpSeason.AirDate != nil {
					d := tmpSeason.AirDate.Format("2006-01-02")
					seasonParams.AirDate = &d
				}
			}
		}

		season, err = w.service.CreateSeason(ctx, seasonParams)
		if err != nil {
			return fmt.Errorf("create season: %w", err)
		}
	}

	// Find or create episode
	episode, err := w.service.GetEpisodeByNumber(ctx, series.ID, int32(*seasonNum), int32(*episodeNum))
	if err != nil {
		// Create episode
		episodeParams := tvshow.CreateEpisodeParams{
			SeriesID:      series.ID,
			SeasonID:      season.ID,
			SeasonNumber:  int32(*seasonNum),
			EpisodeNumber: int32(*episodeNum),
			Title:         fmt.Sprintf("Episode %d", *episodeNum),
		}

		// Check for parsed episode title
		if epTitle := sr.GetString("episode_title"); epTitle != "" {
			episodeParams.Title = epTitle
		}

		// Try to enrich episode from TMDb if series has TMDbID
		if series.TMDbID != nil {
			tmpEpisode := &tvshow.Episode{
				SeriesID:      series.ID,
				SeasonID:      season.ID,
				SeasonNumber:  int32(*seasonNum),
				EpisodeNumber: int32(*episodeNum),
			}
			if err := w.metadataProvider.EnrichEpisode(ctx, tmpEpisode, *series.TMDbID); err == nil {
				episodeParams.TMDbID = tmpEpisode.TMDbID
				episodeParams.Title = tmpEpisode.Title
				episodeParams.Overview = tmpEpisode.Overview
				episodeParams.Runtime = tmpEpisode.Runtime
				episodeParams.StillPath = tmpEpisode.StillPath
				if tmpEpisode.AirDate != nil {
					d := tmpEpisode.AirDate.Format("2006-01-02")
					episodeParams.AirDate = &d
				}
			}
		}

		episode, err = w.service.CreateEpisode(ctx, episodeParams)
		if err != nil {
			return fmt.Errorf("create episode: %w", err)
		}
	}

	// Create episode file record
	fileParams := tvshow.CreateEpisodeFileParams{
		EpisodeID: episode.ID,
		FilePath:  sr.FilePath,
		FileSize:  sr.FileSize,
	}

	// Try to get file stats
	if fileInfo, err := os.Stat(sr.FilePath); err == nil {
		fileParams.FileSize = fileInfo.Size()
	}

	_, err = w.service.CreateEpisodeFile(ctx, fileParams)
	if err != nil {
		return fmt.Errorf("create episode file: %w", err)
	}

	w.logger.Info("processed tv show file",
		slog.String("file_path", sr.FilePath),
		slog.String("series", series.Title),
		slog.Any("season", int32(*seasonNum)),
		slog.Any("episode", int32(*episodeNum)),
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

	// BatchSize is the number of series to process in each batch (for full refresh).
	BatchSize int32 `json:"batch_size,omitempty"`
}

// Kind returns the job kind identifier.
func (MetadataRefreshArgs) Kind() string {
	return KindMetadataRefresh
}

// MetadataRefreshWorker refreshes TV show metadata from external sources.
type MetadataRefreshWorker struct {
	river.WorkerDefaults[MetadataRefreshArgs]
	service   tvshow.Service
	jobClient *infrajobs.Client
	logger    *slog.Logger
}

// NewMetadataRefreshWorker creates a new metadata refresh worker.
func NewMetadataRefreshWorker(service tvshow.Service, jobClient *infrajobs.Client, logger *slog.Logger) *MetadataRefreshWorker {
	return &MetadataRefreshWorker{
		service:   service,
		jobClient: jobClient,
		logger:    logger.With("component", "tvshow_metadata_refresh"),
	}
}

// Timeout returns the maximum execution time for metadata refresh jobs.
func (w *MetadataRefreshWorker) Timeout(job *river.Job[MetadataRefreshArgs]) time.Duration {
	return 15 * time.Minute
}

// Work executes the metadata refresh job.
func (w *MetadataRefreshWorker) Work(ctx context.Context, job *river.Job[MetadataRefreshArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindMetadataRefresh)
	args := job.Args

	jctx.LogStart(
		slog.Any("series_id", args.SeriesID),
		slog.Any("season_id", args.SeasonID),
		slog.Any("episode_id", args.EpisodeID),
		slog.Bool("force", args.Force),
		slog.Bool("refresh_images", args.RefreshImages),
	)

	result := &sharedjobs.JobResult{Success: true}
	start := time.Now()

	// Build refresh options from job args
	opts := tvshow.MetadataRefreshOptions{
		Force: args.Force,
	}

	// Determine what to refresh
	switch {
	case args.EpisodeID != nil:
		if err := w.service.RefreshEpisodeMetadata(ctx, *args.EpisodeID, opts); err != nil {
			result.AddError(fmt.Errorf("refresh episode %s: %w", args.EpisodeID, err))
		} else {
			result.ItemsProcessed++
		}

	case args.SeasonID != nil:
		if err := w.service.RefreshSeasonMetadata(ctx, *args.SeasonID, opts); err != nil {
			result.AddError(fmt.Errorf("refresh season %s: %w", args.SeasonID, err))
		} else {
			result.ItemsProcessed++
		}

	case args.SeriesID != nil:
		if err := w.service.RefreshSeriesMetadata(ctx, *args.SeriesID, opts); err != nil {
			result.AddError(fmt.Errorf("refresh series %s: %w", args.SeriesID, err))
		} else {
			result.ItemsProcessed++
		}

	default:
		// Refresh all series (batch operation)
		batchSize := args.BatchSize
		if batchSize <= 0 {
			batchSize = 50
		}

		var offset int32 = 0
		for {
			// List series in batches
			seriesList, err := w.service.ListSeries(ctx, tvshow.SeriesListFilters{
				Limit:  batchSize,
				Offset: offset,
			})
			if err != nil {
				result.AddError(fmt.Errorf("list series at offset %d: %w", offset, err))
				break
			}

			if len(seriesList) == 0 {
				break
			}

			// Refresh each series in the batch
			for _, s := range seriesList {
				if err := w.service.RefreshSeriesMetadata(ctx, s.ID, opts); err != nil {
					result.AddError(fmt.Errorf("refresh series %s (%s): %w", s.ID, s.Title, err))
				} else {
					result.ItemsProcessed++
				}

				_ = w.jobClient.ReportProgress(ctx, job.ID, &infrajobs.JobProgress{
					Phase:   "refreshing",
					Current: result.ItemsProcessed,
					Message: s.Title,
				})
			}

			w.logger.Info("processed batch",
				slog.Int("batch_size", len(seriesList)),
				slog.Int("total_processed", result.ItemsProcessed),
			)

			offset += batchSize
			if len(seriesList) < int(batchSize) {
				break
			}
		}
	}

	result.Duration = time.Since(start)
	result.Success = !result.HasErrors()
	result.LogSummary(w.logger, KindMetadataRefresh)

	if result.HasErrors() {
		result.LogErrors(w.logger, 10)
		return fmt.Errorf("metadata refresh completed with errors: %d failed", result.ItemsFailed)
	}

	jctx.LogComplete(slog.Int("items_refreshed", result.ItemsProcessed))
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
	service          tvshow.Service
	metadataProvider tvshow.MetadataProvider
	logger           *slog.Logger
}

// NewFileMatchWorker creates a new file match worker.
func NewFileMatchWorker(service tvshow.Service, metadataProvider tvshow.MetadataProvider, logger *slog.Logger) *FileMatchWorker {
	return &FileMatchWorker{
		service:          service,
		metadataProvider: metadataProvider,
		logger:           logger.With("component", "tvshow_file_match"),
	}
}

// Timeout returns the maximum execution time for file match jobs.
func (w *FileMatchWorker) Timeout(job *river.Job[FileMatchArgs]) time.Duration {
	return 5 * time.Minute
}

// Work executes the file match job.
func (w *FileMatchWorker) Work(ctx context.Context, job *river.Job[FileMatchArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindFileMatch)
	args := job.Args

	jctx.LogStart(
		slog.String("file_path", args.FilePath),
		slog.Any("episode_id", args.EpisodeID),
		slog.Bool("force_rematch", args.ForceRematch),
		slog.Bool("auto_create", args.AutoCreate),
	)

	// Check if file exists
	fileInfo, err := os.Stat(args.FilePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if file is already matched
	existingFile, err := w.service.GetEpisodeFileByPath(ctx, args.FilePath)
	if err == nil && existingFile != nil && !args.ForceRematch {
		w.logger.Info("file already matched, skipping",
			slog.Int64("job_id", job.ID),
			slog.String("file_path", args.FilePath),
			slog.String("episode_id", existingFile.EpisodeID.String()),
		)
		return nil
	}

	// If EpisodeID is provided, link file directly to that episode
	if args.EpisodeID != nil {
		episode, err := w.service.GetEpisode(ctx, *args.EpisodeID)
		if err != nil {
			return fmt.Errorf("episode not found: %w", err)
		}

		fileParams := tvshow.CreateEpisodeFileParams{
			EpisodeID: episode.ID,
			FilePath:  args.FilePath,
			FileSize:  fileInfo.Size(),
		}

		_, err = w.service.CreateEpisodeFile(ctx, fileParams)
		if err != nil {
			return fmt.Errorf("create episode file: %w", err)
		}

		w.logger.Info("file matched to episode",
			slog.String("file_path", args.FilePath),
			slog.String("episode_id", episode.ID.String()),
		)
		jctx.LogComplete()
		return nil
	}

	// Parse filename to extract series title, season, episode
	parser := adapters.NewTVShowFileParser()
	seriesTitle, metadata := parser.ParseFromPath(args.FilePath)

	if seriesTitle == "" {
		return fmt.Errorf("could not parse series title from filename: %s", args.FilePath)
	}

	var seasonNum, episodeNum int
	if v, ok := metadata["season"].(int); ok {
		seasonNum = v
	} else {
		return fmt.Errorf("could not parse season number from filename")
	}
	if v, ok := metadata["episode"].(int); ok {
		episodeNum = v
	} else {
		return fmt.Errorf("could not parse episode number from filename")
	}

	// Search for existing series by title
	seriesList, err := w.service.SearchSeries(ctx, seriesTitle, 5, 0)
	if err != nil {
		return fmt.Errorf("search series: %w", err)
	}

	var series *tvshow.Series
	for _, s := range seriesList {
		if normalizeTitle(s.Title) == normalizeTitle(seriesTitle) {
			series = &s
			break
		}
	}

	// If no match and auto_create is enabled, search TMDb and create
	if series == nil {
		if !args.AutoCreate || w.metadataProvider == nil {
			return fmt.Errorf("series not found: %s (auto_create disabled)", seriesTitle)
		}

		searchResults, err := w.metadataProvider.SearchSeries(ctx, seriesTitle, nil)
		if err != nil || len(searchResults) == 0 {
			return fmt.Errorf("series not found in TMDb: %s", seriesTitle)
		}

		newSeries := searchResults[0]
		if err := w.metadataProvider.EnrichSeries(ctx, newSeries); err != nil {
			w.logger.Warn("failed to enrich series", slog.Any("error", err))
		}

		params := seriesToCreateParams(newSeries)
		created, err := w.service.CreateSeries(ctx, params)
		if err != nil {
			return fmt.Errorf("create series: %w", err)
		}
		series = created
	}

	// Find or create season
	season, err := w.service.GetSeasonByNumber(ctx, series.ID, int32(seasonNum))
	if err != nil {
		if !args.AutoCreate {
			return fmt.Errorf("season %d not found for series %s", seasonNum, series.Title)
		}

		seasonParams := tvshow.CreateSeasonParams{
			SeriesID:     series.ID,
			SeasonNumber: int32(seasonNum),
			Name:         fmt.Sprintf("Season %d", seasonNum),
			EpisodeCount: 0,
		}

		season, err = w.service.CreateSeason(ctx, seasonParams)
		if err != nil {
			return fmt.Errorf("create season: %w", err)
		}
	}

	// Find or create episode
	episode, err := w.service.GetEpisodeByNumber(ctx, series.ID, int32(seasonNum), int32(episodeNum))
	if err != nil {
		if !args.AutoCreate {
			return fmt.Errorf("episode S%02dE%02d not found for series %s", seasonNum, episodeNum, series.Title)
		}

		episodeParams := tvshow.CreateEpisodeParams{
			SeriesID:      series.ID,
			SeasonID:      season.ID,
			SeasonNumber:  int32(seasonNum),
			EpisodeNumber: int32(episodeNum),
			Title:         fmt.Sprintf("Episode %d", episodeNum),
		}

		// Check for parsed episode title
		if epTitle, ok := metadata["episode_title"].(string); ok && epTitle != "" {
			episodeParams.Title = epTitle
		}

		episode, err = w.service.CreateEpisode(ctx, episodeParams)
		if err != nil {
			return fmt.Errorf("create episode: %w", err)
		}
	}

	// Create episode file record
	fileParams := tvshow.CreateEpisodeFileParams{
		EpisodeID: episode.ID,
		FilePath:  args.FilePath,
		FileSize:  fileInfo.Size(),
	}

	_, err = w.service.CreateEpisodeFile(ctx, fileParams)
	if err != nil {
		return fmt.Errorf("create episode file: %w", err)
	}

	w.logger.Info("file matched successfully",
		slog.String("file_path", args.FilePath),
		slog.String("series", series.Title),
		slog.Int("season", seasonNum),
		slog.Int("episode", episodeNum),
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

	// BatchSize is the number of series to process in each batch.
	BatchSize int32 `json:"batch_size,omitempty"`
}

// Kind returns the job kind identifier.
func (SearchIndexArgs) Kind() string {
	return KindSearchIndex
}

// InsertOpts returns the default insert options for TV show search index jobs.
// Search indexing runs on the bulk queue since it's batch-heavy.
func (SearchIndexArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: infrajobs.QueueBulk,
	}
}

// SearchIndexWorker indexes TV show content for search.
type SearchIndexWorker struct {
	river.WorkerDefaults[SearchIndexArgs]
	service              tvshow.Service
	searchService        *search.TVShowSearchService
	episodeSearchService *search.EpisodeSearchService
	logger               *slog.Logger
}

// NewSearchIndexWorker creates a new search index worker.
func NewSearchIndexWorker(service tvshow.Service, searchService *search.TVShowSearchService, episodeSearchService *search.EpisodeSearchService, logger *slog.Logger) *SearchIndexWorker {
	return &SearchIndexWorker{
		service:              service,
		searchService:        searchService,
		episodeSearchService: episodeSearchService,
		logger:               logger.With("component", "tvshow_search_index"),
	}
}

// Timeout returns the maximum execution time for search index jobs.
func (w *SearchIndexWorker) Timeout(job *river.Job[SearchIndexArgs]) time.Duration {
	return 10 * time.Minute
}

// Work executes the search index job.
func (w *SearchIndexWorker) Work(ctx context.Context, job *river.Job[SearchIndexArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindSearchIndex)
	args := job.Args

	jctx.LogStart(
		slog.Any("series_id", args.SeriesID),
		slog.Bool("full_reindex", args.FullReindex),
	)

	// Check if search is enabled
	if !w.searchService.IsEnabled() {
		w.logger.Debug("search is disabled, skipping index operation")
		return nil
	}

	result := &sharedjobs.JobResult{Success: true}
	start := time.Now()

	if args.SeriesID != nil {
		// Index specific series
		if err := w.indexSeries(ctx, *args.SeriesID); err != nil {
			result.AddError(fmt.Errorf("index series %s: %w", args.SeriesID, err))
		} else {
			result.ItemsProcessed++
		}
	} else if args.FullReindex {
		// Full reindex of series
		if err := w.searchService.ReindexAll(ctx, w.service); err != nil {
			result.AddError(fmt.Errorf("full series reindex: %w", err))
		} else {
			w.logger.Info("full series reindex completed")
		}
		// Full reindex of episodes
		if err := w.episodeSearchService.ReindexAll(ctx, w.service); err != nil {
			result.AddError(fmt.Errorf("full episode reindex: %w", err))
		} else {
			w.logger.Info("full episode reindex completed")
		}
	}

	result.Duration = time.Since(start)
	result.Success = !result.HasErrors()
	result.LogSummary(w.logger, KindSearchIndex)

	if result.HasErrors() {
		result.LogErrors(w.logger, 10)
		return fmt.Errorf("search indexing completed with errors: %d failed", result.ItemsFailed)
	}

	jctx.LogComplete(slog.Int("items_indexed", result.ItemsProcessed))
	return nil
}

// indexSeries indexes a single series with all its related data.
func (w *SearchIndexWorker) indexSeries(ctx context.Context, seriesID uuid.UUID) error {
	series, err := w.service.GetSeries(ctx, seriesID)
	if err != nil {
		w.logger.Warn("series not found, skipping index",
			slog.String("series_id", seriesID.String()),
		)
		return nil
	}

	genres, err := w.service.GetSeriesGenres(ctx, seriesID)
	if err != nil {
		w.logger.Warn("failed to get genres", slog.Any("error", err))
		genres = nil
	}

	// Fetch all credits (use high limit for indexing)
	cast, _, err := w.service.GetSeriesCast(ctx, seriesID, 1000, 0)
	if err != nil {
		w.logger.Warn("failed to get cast", slog.Any("error", err))
		cast = nil
	}

	crew, _, err := w.service.GetSeriesCrew(ctx, seriesID, 1000, 0)
	if err != nil {
		w.logger.Warn("failed to get crew", slog.Any("error", err))
		crew = nil
	}

	credits := append(cast, crew...)

	networks, err := w.service.GetSeriesNetworks(ctx, seriesID)
	if err != nil {
		w.logger.Warn("failed to get networks", slog.Any("error", err))
		networks = nil
	}

	// Check if series has any episode files
	hasFile := false
	episodes, err := w.service.ListEpisodesBySeries(ctx, seriesID)
	if err == nil {
		for _, ep := range episodes {
			files, err := w.service.ListEpisodeFiles(ctx, ep.ID)
			if err == nil && len(files) > 0 {
				hasFile = true
				break
			}
		}
	}

	if err := w.searchService.UpdateSeries(ctx, series, genres, credits, networks, hasFile); err != nil {
		return fmt.Errorf("failed to index series: %w", err)
	}

	// Also index all episodes for this series
	if w.episodeSearchService.IsEnabled() {
		if err := w.indexSeriesEpisodes(ctx, series, episodes); err != nil {
			w.logger.Warn("failed to index episodes for series",
				slog.String("series_id", seriesID.String()),
				slog.Any("error", err),
			)
		}
	}

	w.logger.Info("series indexed successfully",
		slog.String("series_id", seriesID.String()),
		slog.String("title", series.Title),
	)

	return nil
}

// indexSeriesEpisodes indexes all episodes for a series into the episode search index.
func (w *SearchIndexWorker) indexSeriesEpisodes(ctx context.Context, series *tvshow.Series, episodes []tvshow.Episode) error {
	if len(episodes) == 0 {
		return nil
	}

	posterPath := ""
	if series.PosterPath != nil {
		posterPath = *series.PosterPath
	}

	batch := make([]search.EpisodeWithContext, 0, len(episodes))
	for i := range episodes {
		hasFile := false
		files, err := w.service.ListEpisodeFiles(ctx, episodes[i].ID)
		if err == nil && len(files) > 0 {
			hasFile = true
		}

		batch = append(batch, search.EpisodeWithContext{
			Episode:          &episodes[i],
			SeriesTitle:      series.Title,
			SeriesPosterPath: posterPath,
			HasFile:          hasFile,
		})
	}

	if err := w.episodeSearchService.BulkIndexEpisodes(ctx, batch); err != nil {
		return fmt.Errorf("failed to bulk index episodes: %w", err)
	}

	w.logger.Debug("indexed episodes for series",
		slog.String("series_id", series.ID.String()),
		slog.Int("count", len(batch)),
	)

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
	service   tvshow.Service
	jobClient *infrajobs.Client
	logger    *slog.Logger
}

// NewSeriesRefreshWorker creates a new series refresh worker.
func NewSeriesRefreshWorker(service tvshow.Service, jobClient *infrajobs.Client, logger *slog.Logger) *SeriesRefreshWorker {
	return &SeriesRefreshWorker{
		service:   service,
		jobClient: jobClient,
		logger:    logger.With("component", "tvshow_series_refresh"),
	}
}

// Timeout returns the maximum execution time for series refresh jobs.
func (w *SeriesRefreshWorker) Timeout(job *river.Job[SeriesRefreshArgs]) time.Duration {
	return 10 * time.Minute
}

// Work executes the series refresh job.
func (w *SeriesRefreshWorker) Work(ctx context.Context, job *river.Job[SeriesRefreshArgs]) error {
	jctx := sharedjobs.NewJobContext(ctx, w.logger, job.ID, KindSeriesRefresh)
	args := job.Args

	jctx.LogStart(
		slog.String("series_id", args.SeriesID.String()),
		slog.Any("tmdb_id", args.TMDbID),
		slog.Bool("refresh_seasons", args.RefreshSeasons),
		slog.Bool("refresh_episodes", args.RefreshEpisodes),
	)

	result := &sharedjobs.JobResult{Success: true}
	start := time.Now()

	// Build refresh options from job args
	opts := tvshow.MetadataRefreshOptions{
		Languages: args.Languages,
	}

	// Refresh series metadata using the service
	if err := w.service.RefreshSeriesMetadata(ctx, args.SeriesID, opts); err != nil {
		result.AddError(fmt.Errorf("refresh series %s: %w", args.SeriesID, err))
		w.logger.Error("failed to refresh series",
			slog.String("series_id", args.SeriesID.String()),
			slog.Any("error", err),
		)
	} else {
		result.ItemsProcessed++
		w.logger.Info("refreshed series metadata",
			slog.String("series_id", args.SeriesID.String()),
		)
	}

	// Refresh seasons if requested
	if args.RefreshSeasons || args.RefreshEpisodes {
		seasons, err := w.service.ListSeasons(ctx, args.SeriesID)
		if err != nil {
			w.logger.Warn("failed to list seasons for refresh",
				slog.String("series_id", args.SeriesID.String()),
				slog.Any("error", err),
			)
		} else {
			for i, season := range seasons {
				_ = w.jobClient.ReportProgress(ctx, job.ID, &infrajobs.JobProgress{
					Phase:   "refreshing seasons",
					Current: i + 1,
					Total:   len(seasons),
					Message: season.Name,
				})

				if args.RefreshSeasons {
					if err := w.service.RefreshSeasonMetadata(ctx, season.ID, opts); err != nil {
						w.logger.Warn("failed to refresh season",
							slog.String("season_id", season.ID.String()),
							slog.Any("error", err),
						)
					} else {
						result.ItemsProcessed++
					}
				}

				// Refresh episodes if requested
				if args.RefreshEpisodes {
					episodes, err := w.service.ListEpisodesBySeason(ctx, season.ID)
					if err != nil {
						w.logger.Warn("failed to list episodes for refresh",
							slog.String("season_id", season.ID.String()),
							slog.Any("error", err),
						)
						continue
					}

					for _, ep := range episodes {
						if err := w.service.RefreshEpisodeMetadata(ctx, ep.ID, opts); err != nil {
							w.logger.Warn("failed to refresh episode",
								slog.String("episode_id", ep.ID.String()),
								slog.Any("error", err),
							)
						} else {
							result.ItemsProcessed++
						}
					}
				}
			}
		}
	}

	result.Duration = time.Since(start)
	result.Success = !result.HasErrors()
	result.LogSummary(w.logger, KindSeriesRefresh)

	if result.HasErrors() {
		result.LogErrors(w.logger, 10)
	}

	jctx.LogComplete(
		slog.Int("items_refreshed", result.ItemsProcessed),
		slog.Bool("refresh_seasons", args.RefreshSeasons),
		slog.Bool("refresh_episodes", args.RefreshEpisodes),
	)
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

// normalizeTitle normalizes a title for comparison by lowercasing
// and removing common punctuation.
func normalizeTitle(title string) string {
	// Simple normalization - lowercase
	return strings.ToLower(title)
}

// seriesToCreateParams converts a Series to CreateSeriesParams.
func seriesToCreateParams(s *tvshow.Series) tvshow.CreateSeriesParams {
	params := tvshow.CreateSeriesParams{
		TMDbID:           s.TMDbID,
		TVDbID:           s.TVDbID,
		IMDbID:           s.IMDbID,
		Title:            s.Title,
		OriginalTitle:    s.OriginalTitle,
		OriginalLanguage: s.OriginalLanguage,
		Tagline:          s.Tagline,
		Overview:         s.Overview,
		Status:           s.Status,
		Type:             s.Type,
		VoteCount:        s.VoteCount,
		PosterPath:       s.PosterPath,
		BackdropPath:     s.BackdropPath,
		TotalSeasons:     s.TotalSeasons,
		TotalEpisodes:    s.TotalEpisodes,
		Homepage:         s.Homepage,
		TrailerURL:       s.TrailerURL,
	}
	if s.FirstAirDate != nil {
		d := s.FirstAirDate.Format("2006-01-02")
		params.FirstAirDate = &d
	}
	if s.LastAirDate != nil {
		d := s.LastAirDate.Format("2006-01-02")
		params.LastAirDate = &d
	}
	if s.VoteAverage != nil {
		v := s.VoteAverage.String()
		params.VoteAverage = &v
	}
	if s.Popularity != nil {
		p := s.Popularity.String()
		params.Popularity = &p
	}
	return params
}
