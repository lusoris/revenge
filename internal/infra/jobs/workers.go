// Package jobs provides River job queue setup and workers.
package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

// Queue names for different job types.
const (
	QueueScanning = "scanning"
	QueueMetadata = "metadata"
	QueueIndexing = "indexing"
	QueueImages   = "images"
	QueueCleanup  = "cleanup"
)

// =============================================================================
// Job Args - Define the arguments for each job type
// =============================================================================

// ScanLibraryArgs are arguments for scanning a library.
type ScanLibraryArgs struct {
	LibraryID uuid.UUID `json:"library_id"`
	FullScan  bool      `json:"full_scan"` // If false, only scan for new files
}

func (ScanLibraryArgs) Kind() string { return "scan_library" }

func (args ScanLibraryArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:    QueueScanning,
		Priority: 2, // Lower priority = runs first
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 5 * time.Minute, // Prevent duplicate scans within 5 min
		},
	}
}

// FetchMetadataArgs are arguments for fetching metadata.
type FetchMetadataArgs struct {
	ContentType string    `json:"content_type"` // "movie", "tvshow", "music", etc.
	ContentID   uuid.UUID `json:"content_id"`
	Provider    string    `json:"provider,omitempty"` // Specific provider or empty for auto
	Force       bool      `json:"force"`              // Force refresh even if metadata exists
}

func (FetchMetadataArgs) Kind() string { return "fetch_metadata" }

func (args FetchMetadataArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:    QueueMetadata,
		Priority: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Minute,
		},
	}
}

// DownloadImageArgs are arguments for downloading an image.
type DownloadImageArgs struct {
	ContentType string    `json:"content_type"` // "movie", "person", "studio", etc.
	ContentID   uuid.UUID `json:"content_id"`
	ImageType   string    `json:"image_type"` // "poster", "backdrop", "profile", etc.
	URL         string    `json:"url"`
	Priority    int       `json:"priority"` // Higher = more important (posters > backdrops)
}

func (DownloadImageArgs) Kind() string { return "download_image" }

func (args DownloadImageArgs) InsertOpts() river.InsertOpts {
	prio := 4
	if args.Priority > 0 {
		prio = args.Priority
	}
	return river.InsertOpts{
		Queue:    QueueImages,
		Priority: prio,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 10 * time.Minute,
		},
	}
}

// IndexSearchArgs are arguments for indexing content in search.
type IndexSearchArgs struct {
	Operation   string    `json:"operation"`    // "upsert" or "delete"
	ContentType string    `json:"content_type"` // "movie", "tvshow", "music", etc.
	ContentID   uuid.UUID `json:"content_id"`
}

func (IndexSearchArgs) Kind() string { return "index_search" }

func (args IndexSearchArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:    QueueIndexing,
		Priority: 3,
	}
}

// CleanupArgs are arguments for cleanup operations.
type CleanupArgs struct {
	CleanupType string `json:"cleanup_type"` // "orphaned_files", "expired_sessions", "old_activity"
}

func (CleanupArgs) Kind() string { return "cleanup" }

func (args CleanupArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:    QueueCleanup,
		Priority: 4, // Low priority
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Hour, // Only one cleanup per type per hour
		},
	}
}

// RefreshLibraryArgs are arguments for refreshing library metadata.
type RefreshLibraryArgs struct {
	LibraryID uuid.UUID `json:"library_id"`
}

func (RefreshLibraryArgs) Kind() string { return "refresh_library" }

func (args RefreshLibraryArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:    QueueMetadata,
		Priority: 4,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 15 * time.Minute,
		},
	}
}

// GenerateTrickplayArgs are arguments for generating trickplay images.
type GenerateTrickplayArgs struct {
	ContentType string    `json:"content_type"` // "movie" or "episode"
	ContentID   uuid.UUID `json:"content_id"`
	StreamID    uuid.UUID `json:"stream_id"`
}

func (GenerateTrickplayArgs) Kind() string { return "generate_trickplay" }

func (args GenerateTrickplayArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:    QueueScanning, // Uses scanning queue due to heavy IO
		Priority: 4,             // Lower priority than library scans
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

// =============================================================================
// Worker Dependencies - Interfaces for worker dependencies
// =============================================================================

// LibraryScanner handles library scanning operations.
type LibraryScanner interface {
	ScanLibrary(ctx context.Context, libraryID uuid.UUID, fullScan bool) error
}

// MetadataFetcher handles metadata fetching operations.
type MetadataFetcher interface {
	FetchMetadata(ctx context.Context, contentType string, contentID uuid.UUID, provider string, force bool) error
}

// ImageDownloader handles image download operations.
type ImageDownloader interface {
	DownloadImage(ctx context.Context, contentType string, contentID uuid.UUID, imageType, url string) error
}

// SearchIndexer handles search indexing operations.
type SearchIndexer interface {
	Index(ctx context.Context, contentType string, contentID uuid.UUID) error
	Delete(ctx context.Context, contentType string, contentID uuid.UUID) error
}

// CleanupService handles cleanup operations.
type CleanupService interface {
	CleanupOrphanedFiles(ctx context.Context) error
	CleanupExpiredSessions(ctx context.Context) error
	CleanupOldActivity(ctx context.Context) error
}

// TrickplayGenerator handles trickplay generation.
type TrickplayGenerator interface {
	Generate(ctx context.Context, contentType string, contentID, streamID uuid.UUID) error
}

// =============================================================================
// Worker Implementations
// =============================================================================

// ScanLibraryWorker handles library scanning jobs.
type ScanLibraryWorker struct {
	river.WorkerDefaults[ScanLibraryArgs]
	scanner LibraryScanner
	logger  *slog.Logger
}

func (w *ScanLibraryWorker) Work(ctx context.Context, job *river.Job[ScanLibraryArgs]) error {
	w.logger.Info("Starting library scan",
		slog.String("library_id", job.Args.LibraryID.String()),
		slog.Bool("full_scan", job.Args.FullScan),
	)

	if w.scanner == nil {
		w.logger.Warn("Library scanner not configured, skipping")
		return nil
	}

	if err := w.scanner.ScanLibrary(ctx, job.Args.LibraryID, job.Args.FullScan); err != nil {
		return fmt.Errorf("scan library: %w", err)
	}

	w.logger.Info("Library scan completed",
		slog.String("library_id", job.Args.LibraryID.String()),
	)
	return nil
}

// FetchMetadataWorker handles metadata fetching jobs.
type FetchMetadataWorker struct {
	river.WorkerDefaults[FetchMetadataArgs]
	fetcher MetadataFetcher
	logger  *slog.Logger
}

func (w *FetchMetadataWorker) Work(ctx context.Context, job *river.Job[FetchMetadataArgs]) error {
	w.logger.Info("Fetching metadata",
		slog.String("content_type", job.Args.ContentType),
		slog.String("content_id", job.Args.ContentID.String()),
		slog.String("provider", job.Args.Provider),
	)

	if w.fetcher == nil {
		w.logger.Warn("Metadata fetcher not configured, skipping")
		return nil
	}

	if err := w.fetcher.FetchMetadata(ctx, job.Args.ContentType, job.Args.ContentID, job.Args.Provider, job.Args.Force); err != nil {
		return fmt.Errorf("fetch metadata: %w", err)
	}

	w.logger.Info("Metadata fetched",
		slog.String("content_type", job.Args.ContentType),
		slog.String("content_id", job.Args.ContentID.String()),
	)
	return nil
}

// DownloadImageWorker handles image download jobs.
type DownloadImageWorker struct {
	river.WorkerDefaults[DownloadImageArgs]
	downloader ImageDownloader
	logger     *slog.Logger
}

func (w *DownloadImageWorker) Work(ctx context.Context, job *river.Job[DownloadImageArgs]) error {
	w.logger.Debug("Downloading image",
		slog.String("content_type", job.Args.ContentType),
		slog.String("content_id", job.Args.ContentID.String()),
		slog.String("image_type", job.Args.ImageType),
	)

	if w.downloader == nil {
		w.logger.Warn("Image downloader not configured, skipping")
		return nil
	}

	if err := w.downloader.DownloadImage(ctx, job.Args.ContentType, job.Args.ContentID, job.Args.ImageType, job.Args.URL); err != nil {
		return fmt.Errorf("download image: %w", err)
	}

	return nil
}

// IndexSearchWorker handles search indexing jobs.
type IndexSearchWorker struct {
	river.WorkerDefaults[IndexSearchArgs]
	indexer SearchIndexer
	logger  *slog.Logger
}

func (w *IndexSearchWorker) Work(ctx context.Context, job *river.Job[IndexSearchArgs]) error {
	w.logger.Debug("Indexing content",
		slog.String("operation", job.Args.Operation),
		slog.String("content_type", job.Args.ContentType),
		slog.String("content_id", job.Args.ContentID.String()),
	)

	if w.indexer == nil {
		w.logger.Warn("Search indexer not configured, skipping")
		return nil
	}

	switch job.Args.Operation {
	case "upsert":
		if err := w.indexer.Index(ctx, job.Args.ContentType, job.Args.ContentID); err != nil {
			return fmt.Errorf("index content: %w", err)
		}
	case "delete":
		if err := w.indexer.Delete(ctx, job.Args.ContentType, job.Args.ContentID); err != nil {
			return fmt.Errorf("delete from index: %w", err)
		}
	default:
		return fmt.Errorf("unknown operation: %s", job.Args.Operation)
	}

	return nil
}

// CleanupWorker handles cleanup jobs.
type CleanupWorker struct {
	river.WorkerDefaults[CleanupArgs]
	cleaner CleanupService
	logger  *slog.Logger
}

func (w *CleanupWorker) Work(ctx context.Context, job *river.Job[CleanupArgs]) error {
	w.logger.Info("Running cleanup",
		slog.String("cleanup_type", job.Args.CleanupType),
	)

	if w.cleaner == nil {
		w.logger.Warn("Cleanup service not configured, skipping")
		return nil
	}

	var err error
	switch job.Args.CleanupType {
	case "orphaned_files":
		err = w.cleaner.CleanupOrphanedFiles(ctx)
	case "expired_sessions":
		err = w.cleaner.CleanupExpiredSessions(ctx)
	case "old_activity":
		err = w.cleaner.CleanupOldActivity(ctx)
	default:
		return fmt.Errorf("unknown cleanup type: %s", job.Args.CleanupType)
	}

	if err != nil {
		return fmt.Errorf("cleanup %s: %w", job.Args.CleanupType, err)
	}

	w.logger.Info("Cleanup completed",
		slog.String("cleanup_type", job.Args.CleanupType),
	)
	return nil
}

// RefreshLibraryWorker handles library metadata refresh jobs.
type RefreshLibraryWorker struct {
	river.WorkerDefaults[RefreshLibraryArgs]
	fetcher MetadataFetcher
	logger  *slog.Logger
}

func (w *RefreshLibraryWorker) Work(ctx context.Context, job *river.Job[RefreshLibraryArgs]) error {
	w.logger.Info("Refreshing library metadata",
		slog.String("library_id", job.Args.LibraryID.String()),
	)

	// TODO: Implement library metadata refresh
	// This should fetch updated metadata for all items in the library
	w.logger.Warn("Library refresh not yet implemented")
	return nil
}

// GenerateTrickplayWorker handles trickplay generation jobs.
type GenerateTrickplayWorker struct {
	river.WorkerDefaults[GenerateTrickplayArgs]
	generator TrickplayGenerator
	logger    *slog.Logger
}

func (w *GenerateTrickplayWorker) Work(ctx context.Context, job *river.Job[GenerateTrickplayArgs]) error {
	w.logger.Info("Generating trickplay images",
		slog.String("content_type", job.Args.ContentType),
		slog.String("content_id", job.Args.ContentID.String()),
		slog.String("stream_id", job.Args.StreamID.String()),
	)

	if w.generator == nil {
		w.logger.Warn("Trickplay generator not configured, skipping")
		return nil
	}

	if err := w.generator.Generate(ctx, job.Args.ContentType, job.Args.ContentID, job.Args.StreamID); err != nil {
		return fmt.Errorf("generate trickplay: %w", err)
	}

	return nil
}

// =============================================================================
// Worker Registration
// =============================================================================

// WorkerDeps contains optional dependencies for workers.
type WorkerDeps struct {
	Scanner     LibraryScanner
	Fetcher     MetadataFetcher
	Downloader  ImageDownloader
	Indexer     SearchIndexer
	Cleaner     CleanupService
	Trickplay   TrickplayGenerator
	Logger      *slog.Logger
}

// RegisterWorkers registers all workers with the River workers registry.
func RegisterWorkers(workers *river.Workers, deps WorkerDeps) error {
	logger := deps.Logger
	if logger == nil {
		logger = slog.Default()
	}

	// Register scan library worker
	river.AddWorker(workers, &ScanLibraryWorker{
		scanner: deps.Scanner,
		logger:  logger.With(slog.String("worker", "scan_library")),
	})

	// Register fetch metadata worker
	river.AddWorker(workers, &FetchMetadataWorker{
		fetcher: deps.Fetcher,
		logger:  logger.With(slog.String("worker", "fetch_metadata")),
	})

	// Register download image worker
	river.AddWorker(workers, &DownloadImageWorker{
		downloader: deps.Downloader,
		logger:     logger.With(slog.String("worker", "download_image")),
	})

	// Register index search worker
	river.AddWorker(workers, &IndexSearchWorker{
		indexer: deps.Indexer,
		logger:  logger.With(slog.String("worker", "index_search")),
	})

	// Register cleanup worker
	river.AddWorker(workers, &CleanupWorker{
		cleaner: deps.Cleaner,
		logger:  logger.With(slog.String("worker", "cleanup")),
	})

	// Register refresh library worker
	river.AddWorker(workers, &RefreshLibraryWorker{
		fetcher: deps.Fetcher,
		logger:  logger.With(slog.String("worker", "refresh_library")),
	})

	// Register generate trickplay worker
	river.AddWorker(workers, &GenerateTrickplayWorker{
		generator: deps.Trickplay,
		logger:    logger.With(slog.String("worker", "generate_trickplay")),
	})

	logger.Info("Registered job workers",
		slog.Int("count", 7),
	)

	return nil
}
