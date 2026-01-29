// Package voyage provides adult scene domain models (QAR obfuscation: scenes → voyages).
package voyage

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

// Job kinds for voyage module.
// QAR obfuscation: these handle adult scene operations.
const (
	JobKindEnrichMetadata = "qar.voyage.enrich_metadata"
	JobKindScanFleet      = "qar.voyage.scan_fleet"
	JobKindScanFile       = "qar.voyage.scan_file"
	JobKindFingerprint    = "qar.voyage.fingerprint"
)

// Errors for voyage workers.
var (
	ErrMetadataUnavailable = errors.New("metadata provider unavailable")
)

// MetadataProvider defines the interface for adult scene metadata.
// QAR obfuscation: provides metadata from StashDB/TPDB for adult scenes.
type MetadataProvider interface {
	IsAvailable() bool
	GetVoyageMetadata(ctx context.Context, charter string) (*Metadata, error)
	MatchVoyage(ctx context.Context, title string, launchDate *time.Time, coordinates string, oshash string) (*Metadata, error)
}

// Metadata contains adult scene metadata from StashDB/TPDB.
type Metadata struct {
	Title      string
	Overview   string
	LaunchDate time.Time
	Distance   int // runtime_minutes
	PortName   string // studio name
	Charter    string // stashdb_id
	Registry   string // tpdb_id
	StashID    string
	CoverURL   string
	Rating     float64
}

// =============================================================================
// Metadata Enrichment
// =============================================================================

// EnrichMetadataArgs contains arguments for metadata enrichment job.
// This job fetches additional metadata from StashDB/TPDB for a voyage.
type EnrichMetadataArgs struct {
	VoyageID    uuid.UUID  `json:"voyage_id"`
	Charter     string     `json:"charter,omitempty"`     // stashdb_id if known
	Registry    string     `json:"registry,omitempty"`    // tpdb_id if known
	StashID     string     `json:"stash_id,omitempty"`    // stash scene id if known
	Coordinates string     `json:"coordinates,omitempty"` // phash for matching
	Oshash      string     `json:"oshash,omitempty"`      // oshash for matching
	Title       string     `json:"title,omitempty"`       // for search fallback
	LaunchDate  *time.Time `json:"launch_date,omitempty"` // for search fallback
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
			ByPeriod: 1 * time.Hour, // Don't enrich same voyage more than once per hour
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

	w.logger.Info("enriching voyage metadata",
		"voyage_id", job.Args.VoyageID,
		"charter", job.Args.Charter,
		"attempt", job.Attempt,
	)

	var metadata *Metadata
	var err error

	if job.Args.Charter != "" {
		// Direct lookup by StashDB ID
		metadata, err = w.provider.GetVoyageMetadata(ctx, job.Args.Charter)
	} else {
		// Match by title/date/fingerprints
		metadata, err = w.provider.MatchVoyage(ctx, job.Args.Title, job.Args.LaunchDate, job.Args.Coordinates, job.Args.Oshash)
	}

	if err != nil {
		w.logger.Error("failed to fetch voyage metadata",
			"voyage_id", job.Args.VoyageID,
			"error", err,
		)
		return err
	}

	// Get current voyage to update
	voyage, err := w.service.GetByID(ctx, job.Args.VoyageID)
	if err != nil {
		return fmt.Errorf("get voyage: %w", err)
	}

	// Apply metadata updates
	if metadata.Title != "" {
		voyage.Title = metadata.Title
	}
	if metadata.Overview != "" {
		voyage.Overview = metadata.Overview
	}
	if !metadata.LaunchDate.IsZero() {
		voyage.LaunchDate = &metadata.LaunchDate
	}
	if metadata.Distance > 0 {
		voyage.Distance = metadata.Distance
	}
	if metadata.Charter != "" {
		voyage.Charter = metadata.Charter
	}
	if metadata.Registry != "" {
		voyage.Registry = metadata.Registry
	}
	if metadata.StashID != "" {
		voyage.StashID = metadata.StashID
	}
	if metadata.CoverURL != "" {
		voyage.CoverPath = metadata.CoverURL
	}
	// Note: Rating handled via separate rating system

	if err := w.service.Update(ctx, voyage); err != nil {
		w.logger.Error("failed to apply voyage metadata",
			"voyage_id", job.Args.VoyageID,
			"error", err,
		)
		return err
	}

	w.logger.Info("voyage metadata enriched",
		"voyage_id", job.Args.VoyageID,
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
	Path        string
	Title       string
	LaunchDate  *time.Time
	Coordinates string // phash
	Oshash      string
	IsNew       bool
	VoyageID    uuid.UUID
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
	w.logger.Info("scanning fleet for voyages",
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

	// Queue metadata enrichment for new voyages
	if job.Args.FetchMeta {
		for _, result := range results {
			if result.IsNew {
				_, err := w.client.Insert(ctx, EnrichMetadataArgs{
					VoyageID:    result.VoyageID,
					Title:       result.Title,
					LaunchDate:  result.LaunchDate,
					Coordinates: result.Coordinates,
					Oshash:      result.Oshash,
				}, nil)
				if err != nil {
					w.logger.Warn("failed to queue metadata enrichment",
						"voyage_id", result.VoyageID,
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

	// TODO: Parse file info and create voyage entry
	// This is a placeholder - actual implementation needs file parsing logic

	return nil
}

// =============================================================================
// Fingerprinting
// =============================================================================

// FingerprintArgs contains arguments for fingerprinting job.
// This generates perceptual hash (phash) and oshash for matching against StashDB.
type FingerprintArgs struct {
	VoyageID  uuid.UUID `json:"voyage_id"`
	FilePath  string    `json:"file_path"`
	FetchMeta bool      `json:"fetch_meta"` // If true, queue enrichment after fingerprinting
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

// FingerprintResult contains the generated fingerprints.
type FingerprintResult struct {
	Coordinates string // phash
	Oshash      string
	MD5         string
}

// Fingerprinter defines the interface for generating video fingerprints.
type Fingerprinter interface {
	GenerateFingerprints(ctx context.Context, filePath string) (*FingerprintResult, error)
}

// FingerprintWorker handles fingerprinting jobs.
type FingerprintWorker struct {
	river.WorkerDefaults[FingerprintArgs]
	service       *Service
	fingerprinter Fingerprinter
	client        *river.Client[pgx.Tx]
	logger        *slog.Logger
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
	w.logger.Info("fingerprinting voyage",
		"voyage_id", job.Args.VoyageID,
		"path", job.Args.FilePath,
	)

	// Generate fingerprints (phash, oshash, md5)
	result, err := w.fingerprinter.GenerateFingerprints(ctx, job.Args.FilePath)
	if err != nil {
		w.logger.Error("failed to generate fingerprints",
			"voyage_id", job.Args.VoyageID,
			"error", err,
		)
		return err
	}

	// Update voyage with fingerprints
	voyage, err := w.service.GetByID(ctx, job.Args.VoyageID)
	if err != nil {
		return fmt.Errorf("get voyage: %w", err)
	}

	voyage.Coordinates = result.Coordinates
	voyage.Oshash = result.Oshash
	voyage.MD5 = result.MD5

	if err := w.service.Update(ctx, voyage); err != nil {
		w.logger.Error("failed to update voyage fingerprints",
			"voyage_id", job.Args.VoyageID,
			"error", err,
		)
		return err
	}

	w.logger.Info("voyage fingerprinted",
		"voyage_id", job.Args.VoyageID,
		"coordinates", result.Coordinates,
		"oshash", result.Oshash,
	)

	// Queue metadata enrichment if requested
	if job.Args.FetchMeta && (result.Coordinates != "" || result.Oshash != "") {
		_, err := w.client.Insert(ctx, EnrichMetadataArgs{
			VoyageID:    job.Args.VoyageID,
			Coordinates: result.Coordinates,
			Oshash:      result.Oshash,
			Title:       voyage.Title,
		}, nil)
		if err != nil {
			w.logger.Warn("failed to queue metadata enrichment",
				"voyage_id", job.Args.VoyageID,
				"error", err,
			)
		}
	}

	return nil
}

// =============================================================================
// Worker Registration
// =============================================================================

// RegisterWorkers registers all voyage workers with River.
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
