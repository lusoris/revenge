package movie

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// Job kinds for movie module.
const (
	JobKindEnrichMetadata = "movie.enrich_metadata"
	JobKindScanLibrary    = "movie.scan_library"
	JobKindScanFile       = "movie.scan_file"
)

// EnrichMetadataArgs contains arguments for metadata enrichment job.
// This job fetches additional metadata from TMDb for a movie.
type EnrichMetadataArgs struct {
	MovieID uuid.UUID `json:"movie_id"`
	TmdbID  int       `json:"tmdb_id,omitempty"`  // Optional: if known
	ImdbID  string    `json:"imdb_id,omitempty"`  // Optional: for lookup
	Title   string    `json:"title,omitempty"`    // For search fallback
	Year    int       `json:"year,omitempty"`     // For search fallback
}

// Kind returns the job kind.
func (EnrichMetadataArgs) Kind() string { return JobKindEnrichMetadata }

// InsertOpts returns insert options - metadata jobs go to the metadata queue.
func (EnrichMetadataArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "metadata",
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Hour, // Don't enrich same movie more than once per hour
		},
	}
}

// EnrichMetadataWorker handles metadata enrichment jobs.
// It fetches metadata from TMDb and applies it to the movie.
type EnrichMetadataWorker struct {
	river.WorkerDefaults[EnrichMetadataArgs]
	service  *Service
	provider MetadataProvider
	logger   *slog.Logger
}

// NewEnrichMetadataWorker creates a new metadata enrichment worker.
func NewEnrichMetadataWorker(service *Service, provider MetadataProvider, logger *slog.Logger) *EnrichMetadataWorker {
	return &EnrichMetadataWorker{
		service:  service,
		provider: provider,
		logger:   logger.With("worker", JobKindEnrichMetadata),
	}
}

// Work executes the metadata enrichment job.
func (w *EnrichMetadataWorker) Work(ctx context.Context, job *river.Job[EnrichMetadataArgs]) error {
	if w.provider == nil || !w.provider.IsAvailable() {
		return ErrMetadataUnavailable
	}

	w.logger.Info("enriching movie metadata",
		"movie_id", job.Args.MovieID,
		"tmdb_id", job.Args.TmdbID,
		"attempt", job.Attempt,
	)

	// Get metadata from TMDb
	var metadata *Metadata
	var err error

	if job.Args.TmdbID > 0 {
		// Direct lookup by TMDb ID
		metadata, err = w.provider.GetMovieMetadata(ctx, job.Args.TmdbID)
	} else {
		// Match by title/year/imdb
		metadata, err = w.provider.MatchMovie(ctx, job.Args.Title, job.Args.Year, job.Args.ImdbID)
	}

	if err != nil {
		w.logger.Error("failed to fetch TMDb metadata",
			"movie_id", job.Args.MovieID,
			"error", err,
		)
		return err
	}

	// Convert TMDb metadata to movie update
	updates := MovieMetadataUpdate{
		Title:           metadata.Title,
		OriginalTitle:   metadata.OriginalTitle,
		Overview:        metadata.Overview,
		Tagline:         metadata.Tagline,
		RuntimeTicks:    int64(metadata.RuntimeMinutes) * 60 * 10_000_000, // Convert to ticks
		Budget:          metadata.Budget,
		Revenue:         metadata.Revenue,
		CommunityRating: metadata.Rating,
		VoteCount:       metadata.VoteCount,
		PosterPath:      metadata.PosterURL,
		BackdropPath:    metadata.BackdropURL,
		TmdbID:          metadata.TMDbID,
		ImdbID:          metadata.IMDbID,
	}

	if !metadata.ReleaseDate.IsZero() {
		updates.ReleaseDate = &metadata.ReleaseDate
	}

	// Apply metadata to movie
	if err := w.service.ApplyMetadata(ctx, job.Args.MovieID, updates); err != nil {
		w.logger.Error("failed to apply metadata",
			"movie_id", job.Args.MovieID,
			"error", err,
		)
		return err
	}

	w.logger.Info("movie metadata enriched",
		"movie_id", job.Args.MovieID,
		"tmdb_id", metadata.TMDbID,
	)
	return nil
}

// NextRetry implements custom retry delay with exponential backoff.
func (w *EnrichMetadataWorker) NextRetry(job *river.Job[EnrichMetadataArgs]) time.Time {
	// Exponential backoff: 1min, 5min, 15min, 30min, 1hr
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

// ScanLibraryArgs contains arguments for library scan job.
type ScanLibraryArgs struct {
	LibraryID uuid.UUID `json:"library_id"`
	FullScan  bool      `json:"full_scan"`  // If true, re-scan all files
	FetchMeta bool      `json:"fetch_meta"` // If true, queue enrichment for new items
}

// Kind returns the job kind.
func (ScanLibraryArgs) Kind() string { return JobKindScanLibrary }

// InsertOpts returns insert options - scan jobs go to the scan queue.
func (ScanLibraryArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "scan",
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 5 * time.Minute, // Don't scan same library more than once per 5 min
		},
	}
}

// ScanLibraryWorker handles library scan jobs.
type ScanLibraryWorker struct {
	river.WorkerDefaults[ScanLibraryArgs]
	service *Service
	scanner Scanner
	client  *river.Client[pgx.Tx]
	logger  *slog.Logger
}

// Scanner defines the interface for file scanning.
type Scanner interface {
	ScanLibrary(ctx context.Context, libraryID uuid.UUID, fullScan bool) ([]ScanResult, error)
}

// ScanResult represents a scanned file.
type ScanResult struct {
	Path    string
	Title   string
	Year    int
	IMDbID  string
	IsNew   bool
	MovieID uuid.UUID
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

	// Queue metadata enrichment for new movies
	if job.Args.FetchMeta {
		for _, result := range results {
			if result.IsNew {
				_, err := w.client.Insert(ctx, EnrichMetadataArgs{
					MovieID: result.MovieID,
					Title:   result.Title,
					Year:    result.Year,
					ImdbID:  result.IMDbID,
				}, nil)
				if err != nil {
					w.logger.Warn("failed to queue metadata enrichment",
						"movie_id", result.MovieID,
						"error", err,
					)
				}
			}
		}
	}

	return nil
}

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

	// Check if file already exists
	exists, err := w.service.repo.ExistsByPath(ctx, job.Args.FilePath)
	if err != nil {
		return fmt.Errorf("check file exists: %w", err)
	}

	if exists {
		w.logger.Debug("file already exists", "path", job.Args.FilePath)
		return nil
	}

	// TODO: Parse file info and create movie entry
	// This is a placeholder - actual implementation needs file parsing logic

	return nil
}

// RegisterWorkers registers all movie workers with River.
func RegisterWorkers(
	workers *river.Workers,
	service *Service,
	scanner Scanner,
	provider MetadataProvider,
	client *river.Client[pgx.Tx],
	logger *slog.Logger,
) error {
	// Register enrichment worker if provider is available
	if provider != nil && provider.IsAvailable() {
		if err := river.AddWorkerSafely(workers, NewEnrichMetadataWorker(service, provider, logger)); err != nil {
			return fmt.Errorf("register EnrichMetadataWorker: %w", err)
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
