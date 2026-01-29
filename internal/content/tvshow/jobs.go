package tvshow

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// Job kinds for tvshow module.
const (
	JobKindEnrichSeriesMetadata  = "tvshow.enrich_series_metadata"
	JobKindEnrichSeasonMetadata  = "tvshow.enrich_season_metadata"
	JobKindEnrichEpisodeMetadata = "tvshow.enrich_episode_metadata"
	JobKindScanLibrary           = "tvshow.scan_library"
	JobKindScanFile              = "tvshow.scan_file"
)

// =============================================================================
// Series Metadata Enrichment
// =============================================================================

// EnrichSeriesMetadataArgs contains arguments for series metadata enrichment job.
type EnrichSeriesMetadataArgs struct {
	SeriesID uuid.UUID `json:"series_id"`
	TmdbID   int       `json:"tmdb_id,omitempty"`
	TvdbID   int       `json:"tvdb_id,omitempty"`
	ImdbID   string    `json:"imdb_id,omitempty"`
	Title    string    `json:"title,omitempty"`
	Year     int       `json:"year,omitempty"`
}

// Kind returns the job kind.
func (EnrichSeriesMetadataArgs) Kind() string { return JobKindEnrichSeriesMetadata }

// InsertOpts returns insert options - metadata jobs go to the metadata queue.
func (EnrichSeriesMetadataArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "metadata",
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Hour,
		},
	}
}

// EnrichSeriesMetadataWorker handles series metadata enrichment jobs.
type EnrichSeriesMetadataWorker struct {
	river.WorkerDefaults[EnrichSeriesMetadataArgs]
	service  *Service
	provider MetadataProvider
	client   *river.Client[pgx.Tx]
	logger   *slog.Logger
}

// NewEnrichSeriesMetadataWorker creates a new series metadata enrichment worker.
func NewEnrichSeriesMetadataWorker(
	service *Service,
	provider MetadataProvider,
	client *river.Client[pgx.Tx],
	logger *slog.Logger,
) *EnrichSeriesMetadataWorker {
	return &EnrichSeriesMetadataWorker{
		service:  service,
		provider: provider,
		client:   client,
		logger:   logger.With("worker", JobKindEnrichSeriesMetadata),
	}
}

// Work executes the series metadata enrichment job.
func (w *EnrichSeriesMetadataWorker) Work(ctx context.Context, job *river.Job[EnrichSeriesMetadataArgs]) error {
	if w.provider == nil || !w.provider.IsAvailable() {
		return ErrMetadataUnavailable
	}

	w.logger.Info("enriching series metadata",
		"series_id", job.Args.SeriesID,
		"tmdb_id", job.Args.TmdbID,
		"attempt", job.Attempt,
	)

	var metadata *SeriesMetadata
	var err error

	if job.Args.TmdbID > 0 {
		metadata, err = w.provider.GetSeriesMetadata(ctx, job.Args.TmdbID)
	} else {
		metadata, err = w.provider.MatchSeries(ctx, job.Args.Title, job.Args.Year, job.Args.TvdbID, job.Args.ImdbID)
	}

	if err != nil {
		w.logger.Error("failed to fetch series metadata",
			"series_id", job.Args.SeriesID,
			"error", err,
		)
		return err
	}

	updates := SeriesMetadataUpdate{
		Title:           metadata.Title,
		OriginalTitle:   metadata.OriginalTitle,
		Overview:        metadata.Overview,
		Status:          metadata.Status,
		Type:            metadata.Type,
		CommunityRating: metadata.Rating,
		VoteCount:       metadata.VoteCount,
		PosterPath:      metadata.PosterURL,
		BackdropPath:    metadata.BackdropURL,
		TmdbID:          metadata.TMDbID,
		TvdbID:          metadata.TvdbID,
		ImdbID:          metadata.IMDbID,
	}

	if !metadata.FirstAirDate.IsZero() {
		updates.FirstAirDate = &metadata.FirstAirDate
	}

	if err := w.service.ApplySeriesMetadata(ctx, job.Args.SeriesID, updates); err != nil {
		w.logger.Error("failed to apply series metadata",
			"series_id", job.Args.SeriesID,
			"error", err,
		)
		return err
	}

	w.logger.Info("series metadata enriched",
		"series_id", job.Args.SeriesID,
		"tmdb_id", metadata.TMDbID,
	)

	// Queue season metadata enrichment if TMDb ID is available
	if metadata.TMDbID > 0 {
		series, err := w.service.GetSeries(ctx, job.Args.SeriesID)
		if err == nil {
			seasons, err := w.service.ListSeasons(ctx, job.Args.SeriesID)
			if err == nil {
				for _, season := range seasons {
					_, _ = w.client.Insert(ctx, EnrichSeasonMetadataArgs{
						SeasonID:     season.ID,
						SeriesTmdbID: metadata.TMDbID,
						SeasonNumber: season.SeasonNumber,
						SeriesTitle:  series.Title,
					}, nil)
				}
			}
		}
	}

	return nil
}

// NextRetry implements custom retry delay with exponential backoff.
func (w *EnrichSeriesMetadataWorker) NextRetry(job *river.Job[EnrichSeriesMetadataArgs]) time.Time {
	delays := []time.Duration{
		1 * time.Minute,
		5 * time.Minute,
		15 * time.Minute,
		30 * time.Minute,
		1 * time.Hour,
	}

	idx := job.Attempt - 1
	if idx >= len(delays) {
		idx = len(delays) - 1
	}

	return time.Now().Add(delays[idx])
}

// =============================================================================
// Season Metadata Enrichment
// =============================================================================

// EnrichSeasonMetadataArgs contains arguments for season metadata enrichment job.
type EnrichSeasonMetadataArgs struct {
	SeasonID     uuid.UUID `json:"season_id"`
	SeriesTmdbID int       `json:"series_tmdb_id"`
	SeasonNumber int       `json:"season_number"`
	SeriesTitle  string    `json:"series_title,omitempty"`
}

// Kind returns the job kind.
func (EnrichSeasonMetadataArgs) Kind() string { return JobKindEnrichSeasonMetadata }

// InsertOpts returns insert options.
func (EnrichSeasonMetadataArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "metadata",
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Hour,
		},
	}
}

// EnrichSeasonMetadataWorker handles season metadata enrichment jobs.
type EnrichSeasonMetadataWorker struct {
	river.WorkerDefaults[EnrichSeasonMetadataArgs]
	service  *Service
	provider MetadataProvider
	client   *river.Client[pgx.Tx]
	logger   *slog.Logger
}

// NewEnrichSeasonMetadataWorker creates a new season metadata enrichment worker.
func NewEnrichSeasonMetadataWorker(
	service *Service,
	provider MetadataProvider,
	client *river.Client[pgx.Tx],
	logger *slog.Logger,
) *EnrichSeasonMetadataWorker {
	return &EnrichSeasonMetadataWorker{
		service:  service,
		provider: provider,
		client:   client,
		logger:   logger.With("worker", JobKindEnrichSeasonMetadata),
	}
}

// Work executes the season metadata enrichment job.
func (w *EnrichSeasonMetadataWorker) Work(ctx context.Context, job *river.Job[EnrichSeasonMetadataArgs]) error {
	if w.provider == nil || !w.provider.IsAvailable() {
		return ErrMetadataUnavailable
	}

	w.logger.Info("enriching season metadata",
		"season_id", job.Args.SeasonID,
		"series_tmdb_id", job.Args.SeriesTmdbID,
		"season_number", job.Args.SeasonNumber,
	)

	metadata, err := w.provider.GetSeasonMetadata(ctx, job.Args.SeriesTmdbID, job.Args.SeasonNumber)
	if err != nil {
		w.logger.Error("failed to fetch season metadata",
			"season_id", job.Args.SeasonID,
			"error", err,
		)
		return err
	}

	season, err := w.service.GetSeason(ctx, job.Args.SeasonID)
	if err != nil {
		return err
	}

	if metadata.Name != "" {
		season.Name = metadata.Name
	}
	if metadata.Overview != "" {
		season.Overview = metadata.Overview
	}
	if !metadata.AirDate.IsZero() {
		season.AirDate = &metadata.AirDate
	}
	if metadata.PosterURL != "" {
		season.PosterPath = metadata.PosterURL
	}

	if err := w.service.UpdateSeason(ctx, season); err != nil {
		w.logger.Error("failed to update season",
			"season_id", job.Args.SeasonID,
			"error", err,
		)
		return err
	}

	w.logger.Info("season metadata enriched",
		"season_id", job.Args.SeasonID,
		"tmdb_id", metadata.TMDbID,
	)

	// Queue episode metadata enrichment
	episodes, err := w.service.ListEpisodesBySeason(ctx, job.Args.SeasonID)
	if err == nil {
		for _, ep := range episodes {
			_, _ = w.client.Insert(ctx, EnrichEpisodeMetadataArgs{
				EpisodeID:     ep.ID,
				SeriesTmdbID:  job.Args.SeriesTmdbID,
				SeasonNumber:  job.Args.SeasonNumber,
				EpisodeNumber: ep.EpisodeNumber,
			}, nil)
		}
	}

	return nil
}

// =============================================================================
// Episode Metadata Enrichment
// =============================================================================

// EnrichEpisodeMetadataArgs contains arguments for episode metadata enrichment job.
type EnrichEpisodeMetadataArgs struct {
	EpisodeID     uuid.UUID `json:"episode_id"`
	SeriesTmdbID  int       `json:"series_tmdb_id"`
	SeasonNumber  int       `json:"season_number"`
	EpisodeNumber int       `json:"episode_number"`
}

// Kind returns the job kind.
func (EnrichEpisodeMetadataArgs) Kind() string { return JobKindEnrichEpisodeMetadata }

// InsertOpts returns insert options.
func (EnrichEpisodeMetadataArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "metadata",
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Hour,
		},
	}
}

// EnrichEpisodeMetadataWorker handles episode metadata enrichment jobs.
type EnrichEpisodeMetadataWorker struct {
	river.WorkerDefaults[EnrichEpisodeMetadataArgs]
	service  *Service
	provider MetadataProvider
	logger   *slog.Logger
}

// NewEnrichEpisodeMetadataWorker creates a new episode metadata enrichment worker.
func NewEnrichEpisodeMetadataWorker(
	service *Service,
	provider MetadataProvider,
	logger *slog.Logger,
) *EnrichEpisodeMetadataWorker {
	return &EnrichEpisodeMetadataWorker{
		service:  service,
		provider: provider,
		logger:   logger.With("worker", JobKindEnrichEpisodeMetadata),
	}
}

// Work executes the episode metadata enrichment job.
func (w *EnrichEpisodeMetadataWorker) Work(ctx context.Context, job *river.Job[EnrichEpisodeMetadataArgs]) error {
	if w.provider == nil || !w.provider.IsAvailable() {
		return ErrMetadataUnavailable
	}

	w.logger.Info("enriching episode metadata",
		"episode_id", job.Args.EpisodeID,
		"series_tmdb_id", job.Args.SeriesTmdbID,
		"season", job.Args.SeasonNumber,
		"episode", job.Args.EpisodeNumber,
	)

	metadata, err := w.provider.GetEpisodeMetadata(ctx, job.Args.SeriesTmdbID, job.Args.SeasonNumber, job.Args.EpisodeNumber)
	if err != nil {
		w.logger.Error("failed to fetch episode metadata",
			"episode_id", job.Args.EpisodeID,
			"error", err,
		)
		return err
	}

	episode, err := w.service.GetEpisode(ctx, job.Args.EpisodeID)
	if err != nil {
		return err
	}

	if metadata.Name != "" {
		episode.Title = metadata.Name
	}
	if metadata.Overview != "" {
		episode.Overview = metadata.Overview
	}
	if !metadata.AirDate.IsZero() {
		episode.AirDate = &metadata.AirDate
	}
	if metadata.Runtime > 0 {
		episode.RuntimeTicks = int64(metadata.Runtime) * 60 * 10_000_000
	}
	if metadata.StillURL != "" {
		episode.StillPath = metadata.StillURL
	}
	if metadata.Rating > 0 {
		episode.CommunityRating = metadata.Rating
	}
	if metadata.VoteCount > 0 {
		episode.VoteCount = metadata.VoteCount
	}
	if metadata.TMDbID > 0 {
		episode.TmdbID = metadata.TMDbID
	}
	if metadata.TvdbID > 0 {
		episode.TvdbID = metadata.TvdbID
	}

	if err := w.service.UpdateEpisode(ctx, episode); err != nil {
		w.logger.Error("failed to update episode",
			"episode_id", job.Args.EpisodeID,
			"error", err,
		)
		return err
	}

	w.logger.Info("episode metadata enriched",
		"episode_id", job.Args.EpisodeID,
		"tmdb_id", metadata.TMDbID,
	)

	return nil
}

// =============================================================================
// Library Scanning
// =============================================================================

// ScanLibraryArgs contains arguments for library scan job.
type ScanLibraryArgs struct {
	LibraryID uuid.UUID `json:"library_id"`
	FullScan  bool      `json:"full_scan"`
	FetchMeta bool      `json:"fetch_meta"`
}

// Kind returns the job kind.
func (ScanLibraryArgs) Kind() string { return JobKindScanLibrary }

// InsertOpts returns insert options.
func (ScanLibraryArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "scan",
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 5 * time.Minute,
		},
	}
}

// Scanner defines the interface for file scanning.
type Scanner interface {
	ScanLibrary(ctx context.Context, libraryID uuid.UUID, fullScan bool) ([]ScanResult, error)
}

// ScanResult represents a scanned file.
type ScanResult struct {
	Path     string
	Title    string
	Year     int
	TvdbID   int
	IMDbID   string
	IsNew    bool
	SeriesID uuid.UUID
}

// ScanLibraryWorker handles library scan jobs.
type ScanLibraryWorker struct {
	river.WorkerDefaults[ScanLibraryArgs]
	service *Service
	scanner Scanner
	client  *river.Client[pgx.Tx]
	logger  *slog.Logger
}

// NewScanLibraryWorker creates a new library scan worker.
func NewScanLibraryWorker(
	service *Service,
	scanner Scanner,
	client *river.Client[pgx.Tx],
	logger *slog.Logger,
) *ScanLibraryWorker {
	return &ScanLibraryWorker{
		service: service,
		scanner: scanner,
		client:  client,
		logger:  logger.With("worker", JobKindScanLibrary),
	}
}

// Work executes the library scan job.
func (w *ScanLibraryWorker) Work(ctx context.Context, job *river.Job[ScanLibraryArgs]) error {
	w.logger.Info("scanning library",
		"library_id", job.Args.LibraryID,
		"full_scan", job.Args.FullScan,
	)

	results, err := w.scanner.ScanLibrary(ctx, job.Args.LibraryID, job.Args.FullScan)
	if err != nil {
		w.logger.Error("library scan failed",
			"library_id", job.Args.LibraryID,
			"error", err,
		)
		return err
	}

	w.logger.Info("library scan completed",
		"library_id", job.Args.LibraryID,
		"total_files", len(results),
	)

	// Queue metadata enrichment for new series
	if job.Args.FetchMeta {
		for _, result := range results {
			if result.IsNew {
				_, err := w.client.Insert(ctx, EnrichSeriesMetadataArgs{
					SeriesID: result.SeriesID,
					TvdbID:   result.TvdbID,
					ImdbID:   result.IMDbID,
					Title:    result.Title,
					Year:     result.Year,
				}, nil)
				if err != nil {
					w.logger.Warn("failed to queue metadata enrichment",
						"series_id", result.SeriesID,
						"error", err,
					)
				}
			}
		}
	}

	return nil
}

// =============================================================================
// Single File Scanning
// =============================================================================

// ScanFileArgs contains arguments for single file scan job.
type ScanFileArgs struct {
	LibraryID uuid.UUID `json:"library_id"`
	FilePath  string    `json:"file_path"`
	FetchMeta bool      `json:"fetch_meta"`
}

// Kind returns the job kind.
func (ScanFileArgs) Kind() string { return JobKindScanFile }

// InsertOpts returns insert options.
func (ScanFileArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "scan",
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Minute,
		},
	}
}

// ScanFileWorker handles single file scan jobs.
type ScanFileWorker struct {
	river.WorkerDefaults[ScanFileArgs]
	service *Service
	scanner Scanner
	client  *river.Client[pgx.Tx]
	logger  *slog.Logger
}

// NewScanFileWorker creates a new file scan worker.
func NewScanFileWorker(
	service *Service,
	scanner Scanner,
	client *river.Client[pgx.Tx],
	logger *slog.Logger,
) *ScanFileWorker {
	return &ScanFileWorker{
		service: service,
		scanner: scanner,
		client:  client,
		logger:  logger.With("worker", JobKindScanFile),
	}
}

// Work executes the file scan job.
func (w *ScanFileWorker) Work(ctx context.Context, job *river.Job[ScanFileArgs]) error {
	w.logger.Info("scanning file",
		"library_id", job.Args.LibraryID,
		"path", job.Args.FilePath,
	)

	// TODO: Parse file info and create/update episode entry
	// This is a placeholder - actual implementation needs file parsing logic

	return nil
}

// =============================================================================
// Worker Registration
// =============================================================================

// RegisterWorkers registers all tvshow workers with River.
func RegisterWorkers(
	workers *river.Workers,
	service *Service,
	scanner Scanner,
	provider MetadataProvider,
	client *river.Client[pgx.Tx],
	logger *slog.Logger,
) error {
	// Register enrichment workers if provider is available
	if provider != nil && provider.IsAvailable() {
		if err := river.AddWorkerSafely(workers, NewEnrichSeriesMetadataWorker(service, provider, client, logger)); err != nil {
			return fmt.Errorf("register EnrichSeriesMetadataWorker: %w", err)
		}

		if err := river.AddWorkerSafely(workers, NewEnrichSeasonMetadataWorker(service, provider, client, logger)); err != nil {
			return fmt.Errorf("register EnrichSeasonMetadataWorker: %w", err)
		}

		if err := river.AddWorkerSafely(workers, NewEnrichEpisodeMetadataWorker(service, provider, logger)); err != nil {
			return fmt.Errorf("register EnrichEpisodeMetadataWorker: %w", err)
		}
	}

	// Register scan workers if scanner is available
	if scanner != nil {
		if err := river.AddWorkerSafely(workers, NewScanLibraryWorker(service, scanner, client, logger)); err != nil {
			return fmt.Errorf("register ScanLibraryWorker: %w", err)
		}

		if err := river.AddWorkerSafely(workers, NewScanFileWorker(service, scanner, client, logger)); err != nil {
			return fmt.Errorf("register ScanFileWorker: %w", err)
		}
	}

	return nil
}
