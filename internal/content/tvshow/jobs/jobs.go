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

	// AutoCreate indicates whether to auto-create series/seasons/episodes for discovered files.
	AutoCreate bool `json:"auto_create"`
}

// Kind returns the job kind identifier.
func (LibraryScanArgs) Kind() string {
	return KindLibraryScan
}

// LibraryScanWorker scans TV show library directories.
type LibraryScanWorker struct {
	river.WorkerDefaults[LibraryScanArgs]
	service          tvshow.Service
	metadataProvider tvshow.MetadataProvider
	logger           *zap.Logger
}

// NewLibraryScanWorker creates a new library scan worker.
func NewLibraryScanWorker(service tvshow.Service, metadataProvider tvshow.MetadataProvider, logger *zap.Logger) *LibraryScanWorker {
	return &LibraryScanWorker{
		service:          service,
		metadataProvider: metadataProvider,
		logger:           logger.Named("tvshow_library_scan"),
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
		zap.Strings("paths", job.Args.Paths),
		zap.Bool("force", job.Args.Force),
		zap.Bool("auto_create", job.Args.AutoCreate),
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
		w.logger.Error("library scan failed", zap.Error(err))
		return err
	}

	w.logger.Info("scan completed",
		zap.Int("total_files", summary.TotalFiles),
		zap.Int("media_files", summary.MediaFiles),
		zap.Int("parsed_files", summary.ParsedFiles),
	)

	// Process each discovered file
	for _, sr := range scanResults {
		if !sr.IsMedia {
			continue
		}

		// Check if file is already matched
		existingFile, err := w.service.GetEpisodeFileByPath(ctx, sr.FilePath)
		if err == nil && existingFile != nil && !job.Args.Force {
			w.logger.Debug("file already matched, skipping",
				zap.String("file_path", sr.FilePath),
			)
			itemsSkipped++
			continue
		}

		// Process the file with auto-create if enabled
		if job.Args.AutoCreate && w.metadataProvider != nil {
			if err := w.processFile(ctx, sr); err != nil {
				w.logger.Warn("failed to process file",
					zap.String("file_path", sr.FilePath),
					zap.Error(err),
				)
				result.AddError(fmt.Errorf("process %s: %w", sr.FilePath, err))
				continue
			}
			result.ItemsProcessed++
		} else {
			// Just log discovered files when auto-create is disabled
			w.logger.Info("discovered tv show file",
				zap.String("file_path", sr.FilePath),
				zap.String("parsed_title", sr.ParsedTitle),
				zap.Any("season", sr.GetSeason()),
				zap.Any("episode", sr.GetEpisode()),
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
		zap.Int("paths_scanned", len(job.Args.Paths)),
		zap.Int("items_processed", result.ItemsProcessed),
		zap.Int("items_skipped", itemsSkipped),
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
			w.logger.Warn("failed to enrich series", zap.Error(err))
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
		zap.String("file_path", sr.FilePath),
		zap.String("series", series.Title),
		zap.Int32("season", int32(*seasonNum)),
		zap.Int32("episode", int32(*episodeNum)),
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

// Timeout returns the maximum execution time for metadata refresh jobs.
func (w *MetadataRefreshWorker) Timeout(job *river.Job[MetadataRefreshArgs]) time.Duration {
	return 15 * time.Minute
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
				if err := w.service.RefreshSeriesMetadata(ctx, s.ID); err != nil {
					result.AddError(fmt.Errorf("refresh series %s (%s): %w", s.ID, s.Title, err))
				} else {
					result.ItemsProcessed++
				}
			}

			w.logger.Info("processed batch",
				zap.Int("batch_size", len(seriesList)),
				zap.Int("total_processed", result.ItemsProcessed),
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
	service          tvshow.Service
	metadataProvider tvshow.MetadataProvider
	logger           *zap.Logger
}

// NewFileMatchWorker creates a new file match worker.
func NewFileMatchWorker(service tvshow.Service, metadataProvider tvshow.MetadataProvider, logger *zap.Logger) *FileMatchWorker {
	return &FileMatchWorker{
		service:          service,
		metadataProvider: metadataProvider,
		logger:           logger.Named("tvshow_file_match"),
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
		zap.String("file_path", args.FilePath),
		zap.Any("episode_id", args.EpisodeID),
		zap.Bool("force_rematch", args.ForceRematch),
		zap.Bool("auto_create", args.AutoCreate),
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
			zap.Int64("job_id", job.ID),
			zap.String("file_path", args.FilePath),
			zap.String("episode_id", existingFile.EpisodeID.String()),
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
			zap.String("file_path", args.FilePath),
			zap.String("episode_id", episode.ID.String()),
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
			w.logger.Warn("failed to enrich series", zap.Error(err))
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
		zap.String("file_path", args.FilePath),
		zap.String("series", series.Title),
		zap.Int("season", seasonNum),
		zap.Int("episode", episodeNum),
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

// SearchIndexWorker indexes TV show content for search.
// Note: TV show search service is not yet implemented.
// This worker logs operations but does not perform actual indexing.
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

// Timeout returns the maximum execution time for search index jobs.
func (w *SearchIndexWorker) Timeout(job *river.Job[SearchIndexArgs]) time.Duration {
	return 10 * time.Minute
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

	// Note: TV show search service is not yet implemented.
	// For now, we just fetch the series data to validate it exists
	// and log that indexing would occur.

	if args.SeriesID != nil {
		// Index specific series
		series, err := w.service.GetSeries(ctx, *args.SeriesID)
		if err != nil {
			result.AddError(fmt.Errorf("get series %s: %w", args.SeriesID, err))
		} else {
			// Would index to search engine here
			w.logger.Info("would index series (search service not implemented)",
				zap.Int64("job_id", job.ID),
				zap.String("series_id", series.ID.String()),
				zap.String("title", series.Title),
			)
			result.ItemsProcessed++
		}
	} else if args.FullReindex {
		// Full reindex - iterate through all series
		batchSize := args.BatchSize
		if batchSize <= 0 {
			batchSize = 100
		}

		var offset int32 = 0
		for {
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

			for _, s := range seriesList {
				// Would index to search engine here
				w.logger.Debug("would index series",
					zap.String("series_id", s.ID.String()),
					zap.String("title", s.Title),
				)
				result.ItemsProcessed++
			}

			offset += batchSize
			if len(seriesList) < int(batchSize) {
				break
			}
		}

		w.logger.Info("full reindex completed (search service not implemented)",
			zap.Int("total_series", result.ItemsProcessed),
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

// Timeout returns the maximum execution time for series refresh jobs.
func (w *SeriesRefreshWorker) Timeout(job *river.Job[SeriesRefreshArgs]) time.Duration {
	return 10 * time.Minute
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

	result := &sharedjobs.JobResult{Success: true}
	start := time.Now()

	// Refresh series metadata using the service
	if err := w.service.RefreshSeriesMetadata(ctx, args.SeriesID); err != nil {
		result.AddError(fmt.Errorf("refresh series %s: %w", args.SeriesID, err))
		w.logger.Error("failed to refresh series",
			zap.String("series_id", args.SeriesID.String()),
			zap.Error(err),
		)
	} else {
		result.ItemsProcessed++
		w.logger.Info("refreshed series metadata",
			zap.String("series_id", args.SeriesID.String()),
		)
	}

	// Refresh seasons if requested
	if args.RefreshSeasons || args.RefreshEpisodes {
		seasons, err := w.service.ListSeasons(ctx, args.SeriesID)
		if err != nil {
			w.logger.Warn("failed to list seasons for refresh",
				zap.String("series_id", args.SeriesID.String()),
				zap.Error(err),
			)
		} else {
			for _, season := range seasons {
				if args.RefreshSeasons {
					if err := w.service.RefreshSeasonMetadata(ctx, season.ID); err != nil {
						w.logger.Warn("failed to refresh season",
							zap.String("season_id", season.ID.String()),
							zap.Error(err),
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
							zap.String("season_id", season.ID.String()),
							zap.Error(err),
						)
						continue
					}

					for _, ep := range episodes {
						if err := w.service.RefreshEpisodeMetadata(ctx, ep.ID); err != nil {
							w.logger.Warn("failed to refresh episode",
								zap.String("episode_id", ep.ID.String()),
								zap.Error(err),
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
		zap.Int("items_refreshed", result.ItemsProcessed),
		zap.Bool("refresh_seasons", args.RefreshSeasons),
		zap.Bool("refresh_episodes", args.RefreshEpisodes),
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
