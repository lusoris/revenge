// Package expedition provides adult movie domain models (QAR obfuscation: movies → expeditions).
package expedition

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// Job kinds for expedition module.
// QAR obfuscation: these handle adult movie operations.
const (
	JobKindEnrichMetadata = "qar.expedition.enrich_metadata"
	JobKindScanFleet      = "qar.expedition.scan_fleet"
	JobKindScanFile       = "qar.expedition.scan_file"
	JobKindFingerprint    = "qar.expedition.fingerprint"
)

// Errors for expedition workers.
var (
	ErrMetadataUnavailable = errors.New("metadata provider unavailable")
)

// MetadataProvider defines the interface for adult content metadata.
// QAR obfuscation: provides metadata from StashDB/TPDB for adult movies.
type MetadataProvider interface {
	IsAvailable() bool
	GetExpeditionMetadata(ctx context.Context, charter string) (*Metadata, error)
	MatchExpedition(ctx context.Context, title string, launchYear int, coordinates string) (*Metadata, error)
}

// Metadata contains adult movie metadata from StashDB/TPDB.
type Metadata struct {
	Title        string
	Overview     string
	LaunchDate   time.Time
	RuntimeTicks int64
	Director     string
	Series       string
	PortName     string // studio name
	Charter      string // stashdb_id
	Registry     string // tpdb_id
	PosterURL    string
	BackdropURL  string
	Rating       float64
}

// =============================================================================
// Metadata Enrichment
// =============================================================================

// EnrichMetadataArgs contains arguments for metadata enrichment job.
// This job fetches additional metadata from StashDB/TPDB for an expedition.
type EnrichMetadataArgs struct {
	ExpeditionID uuid.UUID `json:"expedition_id"`
	Charter      string    `json:"charter,omitempty"`      // stashdb_id if known
	Registry     string    `json:"registry,omitempty"`     // tpdb_id if known
	Coordinates  string    `json:"coordinates,omitempty"`  // phash for matching
	Title        string    `json:"title,omitempty"`        // for search fallback
	LaunchYear   int       `json:"launch_year,omitempty"`  // for search fallback
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
			ByPeriod: 1 * time.Hour, // Don't enrich same expedition more than once per hour
		},
	}
}

// EnrichMetadataWorker handles metadata enrichment jobs.
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

	w.logger.Info("enriching expedition metadata",
		"expedition_id", job.Args.ExpeditionID,
		"charter", job.Args.Charter,
		"attempt", job.Attempt,
	)

	var metadata *Metadata
	var err error

	if job.Args.Charter != "" {
		// Direct lookup by StashDB ID
		metadata, err = w.provider.GetExpeditionMetadata(ctx, job.Args.Charter)
	} else {
		// Match by title/year/fingerprint
		metadata, err = w.provider.MatchExpedition(ctx, job.Args.Title, job.Args.LaunchYear, job.Args.Coordinates)
	}

	if err != nil {
		w.logger.Error("failed to fetch expedition metadata",
			"expedition_id", job.Args.ExpeditionID,
			"error", err,
		)
		return err
	}

	// Get current expedition to update
	expedition, err := w.service.GetByID(ctx, job.Args.ExpeditionID)
	if err != nil {
		return fmt.Errorf("get expedition: %w", err)
	}

	// Apply metadata updates
	if metadata.Title != "" {
		expedition.Title = metadata.Title
	}
	if metadata.Overview != "" {
		expedition.Overview = metadata.Overview
	}
	if !metadata.LaunchDate.IsZero() {
		expedition.LaunchDate = &metadata.LaunchDate
	}
	if metadata.RuntimeTicks > 0 {
		expedition.RuntimeTicks = metadata.RuntimeTicks
	}
	if metadata.Director != "" {
		expedition.Director = metadata.Director
	}
	if metadata.Series != "" {
		expedition.Series = metadata.Series
	}
	if metadata.Charter != "" {
		expedition.Charter = metadata.Charter
	}
	if metadata.Registry != "" {
		expedition.Registry = metadata.Registry
	}
	// Note: PosterURL/BackdropURL/Rating handled via separate image/rating systems

	if err := w.service.Update(ctx, expedition); err != nil {
		w.logger.Error("failed to apply expedition metadata",
			"expedition_id", job.Args.ExpeditionID,
			"error", err,
		)
		return err
	}

	w.logger.Info("expedition metadata enriched",
		"expedition_id", job.Args.ExpeditionID,
		"charter", metadata.Charter,
	)
	return nil
}

// NextRetry implements custom retry delay with exponential backoff.
func (w *EnrichMetadataWorker) NextRetry(job *river.Job[EnrichMetadataArgs]) time.Time {
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
// Fleet Scanning
// =============================================================================

// ScanFleetArgs contains arguments for fleet scan job.
// QAR obfuscation: fleet = library for adult content.
type ScanFleetArgs struct {
	FleetID   uuid.UUID `json:"fleet_id"`
	FullScan  bool      `json:"full_scan"`  // If true, re-scan all files
	FetchMeta bool      `json:"fetch_meta"` // If true, queue enrichment for new items
}

// Kind returns the job kind.
func (ScanFleetArgs) Kind() string { return JobKindScanFleet }

// InsertOpts returns insert options - scan jobs go to the scan queue.
func (ScanFleetArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "scan",
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 5 * time.Minute, // Don't scan same fleet more than once per 5 min
		},
	}
}

// Scanner defines the interface for file scanning.
type Scanner interface {
	ScanFleet(ctx context.Context, fleetID uuid.UUID, fullScan bool) ([]ScanResult, error)
}

// ScanResult represents a scanned file.
type ScanResult struct {
	Path         string
	Title        string
	LaunchYear   int
	Coordinates  string // phash
	IsNew        bool
	ExpeditionID uuid.UUID
}

// ScanFleetWorker handles fleet scan jobs.
type ScanFleetWorker struct {
	river.WorkerDefaults[ScanFleetArgs]
	service *Service
	scanner Scanner
	client  *river.Client[pgx.Tx]
	logger  *slog.Logger
}

// NewScanFleetWorker creates a new fleet scan worker.
func NewScanFleetWorker(
	service *Service,
	scanner Scanner,
	client *river.Client[pgx.Tx],
	logger *slog.Logger,
) *ScanFleetWorker {
	return &ScanFleetWorker{
		service: service,
		scanner: scanner,
		client:  client,
		logger:  logger.With("worker", JobKindScanFleet),
	}
}

// Work executes the fleet scan job.
func (w *ScanFleetWorker) Work(ctx context.Context, job *river.Job[ScanFleetArgs]) error {
	w.logger.Info("scanning fleet",
		"fleet_id", job.Args.FleetID,
		"full_scan", job.Args.FullScan,
	)

	results, err := w.scanner.ScanFleet(ctx, job.Args.FleetID, job.Args.FullScan)
	if err != nil {
		w.logger.Error("fleet scan failed",
			"fleet_id", job.Args.FleetID,
			"error", err,
		)
		return err
	}

	w.logger.Info("fleet scan completed",
		"fleet_id", job.Args.FleetID,
		"total_files", len(results),
	)

	// Queue metadata enrichment for new expeditions
	if job.Args.FetchMeta {
		for _, result := range results {
			if result.IsNew {
				_, err := w.client.Insert(ctx, EnrichMetadataArgs{
					ExpeditionID: result.ExpeditionID,
					Title:        result.Title,
					LaunchYear:   result.LaunchYear,
					Coordinates:  result.Coordinates,
				}, nil)
				if err != nil {
					w.logger.Warn("failed to queue metadata enrichment",
						"expedition_id", result.ExpeditionID,
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
	FleetID   uuid.UUID `json:"fleet_id"`
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
		"fleet_id", job.Args.FleetID,
		"path", job.Args.FilePath,
	)

	// TODO: Parse file info and create expedition entry
	// This is a placeholder - actual implementation needs file parsing logic

	return nil
}

// =============================================================================
// Fingerprinting
// =============================================================================

// FingerprintArgs contains arguments for fingerprinting job.
// This generates perceptual hash (phash) for matching against StashDB.
type FingerprintArgs struct {
	ExpeditionID uuid.UUID `json:"expedition_id"`
	FilePath     string    `json:"file_path"`
	FetchMeta    bool      `json:"fetch_meta"` // If true, queue enrichment after fingerprinting
}

// Kind returns the job kind.
func (FingerprintArgs) Kind() string { return JobKindFingerprint }

// InsertOpts returns insert options.
func (FingerprintArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       "fingerprint",
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 24 * time.Hour, // Don't re-fingerprint same file within 24 hours
		},
	}
}

// Fingerprinter defines the interface for generating video fingerprints.
type Fingerprinter interface {
	GeneratePhash(ctx context.Context, filePath string) (string, error)
}

// FingerprintWorker handles fingerprinting jobs.
type FingerprintWorker struct {
	river.WorkerDefaults[FingerprintArgs]
	service      *Service
	fingerprinter Fingerprinter
	client       *river.Client[pgx.Tx]
	logger       *slog.Logger
}

// NewFingerprintWorker creates a new fingerprinting worker.
func NewFingerprintWorker(
	service *Service,
	fingerprinter Fingerprinter,
	client *river.Client[pgx.Tx],
	logger *slog.Logger,
) *FingerprintWorker {
	return &FingerprintWorker{
		service:       service,
		fingerprinter: fingerprinter,
		client:        client,
		logger:        logger.With("worker", JobKindFingerprint),
	}
}

// Work executes the fingerprinting job.
func (w *FingerprintWorker) Work(ctx context.Context, job *river.Job[FingerprintArgs]) error {
	w.logger.Info("fingerprinting expedition",
		"expedition_id", job.Args.ExpeditionID,
		"path", job.Args.FilePath,
	)

	// Generate perceptual hash
	coordinates, err := w.fingerprinter.GeneratePhash(ctx, job.Args.FilePath)
	if err != nil {
		w.logger.Error("failed to generate phash",
			"expedition_id", job.Args.ExpeditionID,
			"error", err,
		)
		return err
	}

	// Update expedition with fingerprint
	expedition, err := w.service.GetByID(ctx, job.Args.ExpeditionID)
	if err != nil {
		return fmt.Errorf("get expedition: %w", err)
	}

	expedition.Coordinates = coordinates
	if err := w.service.Update(ctx, expedition); err != nil {
		w.logger.Error("failed to update expedition coordinates",
			"expedition_id", job.Args.ExpeditionID,
			"error", err,
		)
		return err
	}

	w.logger.Info("expedition fingerprinted",
		"expedition_id", job.Args.ExpeditionID,
		"coordinates", coordinates,
	)

	// Queue metadata enrichment if requested
	if job.Args.FetchMeta && coordinates != "" {
		_, err := w.client.Insert(ctx, EnrichMetadataArgs{
			ExpeditionID: job.Args.ExpeditionID,
			Coordinates:  coordinates,
			Title:        expedition.Title,
		}, nil)
		if err != nil {
			w.logger.Warn("failed to queue metadata enrichment",
				"expedition_id", job.Args.ExpeditionID,
				"error", err,
			)
		}
	}

	return nil
}

// =============================================================================
// Worker Registration
// =============================================================================

// RegisterWorkers registers all expedition workers with River.
func RegisterWorkers(
	workers *river.Workers,
	service *Service,
	scanner Scanner,
	fingerprinter Fingerprinter,
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
		if err := river.AddWorkerSafely(workers, NewScanFleetWorker(service, scanner, client, logger)); err != nil {
			return fmt.Errorf("register ScanFleetWorker: %w", err)
		}

		if err := river.AddWorkerSafely(workers, NewScanFileWorker(service, scanner, client, logger)); err != nil {
			return fmt.Errorf("register ScanFileWorker: %w", err)
		}
	}

	// Register fingerprint worker if fingerprinter is available
	if fingerprinter != nil {
		if err := river.AddWorkerSafely(workers, NewFingerprintWorker(service, fingerprinter, client, logger)); err != nil {
			return fmt.Errorf("register FingerprintWorker: %w", err)
		}
	}

	return nil
}
