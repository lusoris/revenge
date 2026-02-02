// Package library provides library management and access control functionality.
package library

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for library persistence.
type Repository interface {
	// Library CRUD
	Create(ctx context.Context, lib *Library) error
	Get(ctx context.Context, id uuid.UUID) (*Library, error)
	GetByName(ctx context.Context, name string) (*Library, error)
	List(ctx context.Context) ([]Library, error)
	ListEnabled(ctx context.Context) ([]Library, error)
	ListByType(ctx context.Context, libType string) ([]Library, error)
	Update(ctx context.Context, id uuid.UUID, update *LibraryUpdate) (*Library, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	CountByType(ctx context.Context, libType string) (int64, error)

	// Library Scans
	CreateScan(ctx context.Context, scan *LibraryScan) error
	GetScan(ctx context.Context, id uuid.UUID) (*LibraryScan, error)
	ListScans(ctx context.Context, libraryID uuid.UUID, limit, offset int32) ([]LibraryScan, error)
	CountScans(ctx context.Context, libraryID uuid.UUID) (int64, error)
	GetLatestScan(ctx context.Context, libraryID uuid.UUID) (*LibraryScan, error)
	GetRunningScans(ctx context.Context) ([]LibraryScan, error)
	UpdateScanStatus(ctx context.Context, id uuid.UUID, status *ScanStatusUpdate) (*LibraryScan, error)
	UpdateScanProgress(ctx context.Context, id uuid.UUID, progress *ScanProgress) (*LibraryScan, error)
	DeleteOldScans(ctx context.Context, olderThan time.Time) (int64, error)

	// Library Permissions
	GrantPermission(ctx context.Context, perm *Permission) error
	GetPermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) (*Permission, error)
	ListPermissions(ctx context.Context, libraryID uuid.UUID) ([]Permission, error)
	ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)
	CheckPermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) (bool, error)
	GetUserAccessibleLibraries(ctx context.Context, userID uuid.UUID) ([]Library, error)
	RevokePermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) error
	RevokeAllPermissions(ctx context.Context, libraryID uuid.UUID) error
	RevokeUserPermissions(ctx context.Context, userID uuid.UUID) error
	CountPermissions(ctx context.Context, libraryID uuid.UUID) (int64, error)
}

// Library represents a media library.
type Library struct {
	ID                 uuid.UUID              `json:"id"`
	Name               string                 `json:"name"`
	Type               string                 `json:"type"`
	Paths              []string               `json:"paths"`
	Enabled            bool                   `json:"enabled"`
	ScanOnStartup      bool                   `json:"scan_on_startup"`
	RealtimeMonitoring bool                   `json:"realtime_monitoring"`
	MetadataProvider   *string                `json:"metadata_provider,omitempty"`
	PreferredLanguage  string                 `json:"preferred_language"`
	ScannerConfig      map[string]interface{} `json:"scanner_config,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

// LibraryUpdate represents fields that can be updated on a library.
type LibraryUpdate struct {
	Name               *string                `json:"name,omitempty"`
	Type               *string                `json:"type,omitempty"`
	Paths              []string               `json:"paths,omitempty"`
	Enabled            *bool                  `json:"enabled,omitempty"`
	ScanOnStartup      *bool                  `json:"scan_on_startup,omitempty"`
	RealtimeMonitoring *bool                  `json:"realtime_monitoring,omitempty"`
	MetadataProvider   *string                `json:"metadata_provider,omitempty"`
	PreferredLanguage  *string                `json:"preferred_language,omitempty"`
	ScannerConfig      map[string]interface{} `json:"scanner_config,omitempty"`
}

// LibraryScan represents a library scan job.
type LibraryScan struct {
	ID              uuid.UUID  `json:"id"`
	LibraryID       uuid.UUID  `json:"library_id"`
	ScanType        string     `json:"scan_type"`
	Status          string     `json:"status"`
	ItemsScanned    int32      `json:"items_scanned"`
	ItemsAdded      int32      `json:"items_added"`
	ItemsUpdated    int32      `json:"items_updated"`
	ItemsRemoved    int32      `json:"items_removed"`
	ErrorsCount     int32      `json:"errors_count"`
	ErrorMessage    *string    `json:"error_message,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	DurationSeconds *int32     `json:"duration_seconds,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ScanStatusUpdate represents a scan status update.
type ScanStatusUpdate struct {
	Status          string
	StartedAt       *time.Time
	CompletedAt     *time.Time
	DurationSeconds *int32
	ErrorMessage    *string
}

// ScanProgress represents scan progress update.
type ScanProgress struct {
	ItemsScanned int32
	ItemsAdded   int32
	ItemsUpdated int32
	ItemsRemoved int32
	ErrorsCount  int32
}

// Permission represents a library permission for a user.
type Permission struct {
	ID         uuid.UUID `json:"id"`
	LibraryID  uuid.UUID `json:"library_id"`
	UserID     uuid.UUID `json:"user_id"`
	Permission string    `json:"permission"`
	CreatedAt  time.Time `json:"created_at"`
}

// Library types
const (
	LibraryTypeMovie     = "movie"
	LibraryTypeTVShow    = "tvshow"
	LibraryTypeMusic     = "music"
	LibraryTypePhoto     = "photo"
	LibraryTypeBook      = "book"
	LibraryTypeAudiobook = "audiobook"
	LibraryTypeComic     = "comic"
	LibraryTypePodcast   = "podcast"
	LibraryTypeAdult     = "adult"
)

// Scan types
const (
	ScanTypeFull        = "full"
	ScanTypeIncremental = "incremental"
	ScanTypeMetadata    = "metadata"
)

// Scan statuses
const (
	ScanStatusPending   = "pending"
	ScanStatusRunning   = "running"
	ScanStatusCompleted = "completed"
	ScanStatusFailed    = "failed"
	ScanStatusCancelled = "cancelled"
)

// Permission types
const (
	PermissionView     = "view"
	PermissionDownload = "download"
	PermissionManage   = "manage"
)

// ValidLibraryTypes returns all valid library types.
func ValidLibraryTypes() []string {
	return []string{
		LibraryTypeMovie,
		LibraryTypeTVShow,
		LibraryTypeMusic,
		LibraryTypePhoto,
		LibraryTypeBook,
		LibraryTypeAudiobook,
		LibraryTypeComic,
		LibraryTypePodcast,
		LibraryTypeAdult,
	}
}

// ValidScanTypes returns all valid scan types.
func ValidScanTypes() []string {
	return []string{
		ScanTypeFull,
		ScanTypeIncremental,
		ScanTypeMetadata,
	}
}

// ValidPermissions returns all valid permission types.
func ValidPermissions() []string {
	return []string{
		PermissionView,
		PermissionDownload,
		PermissionManage,
	}
}

// IsValidLibraryType checks if the given type is valid.
func IsValidLibraryType(t string) bool {
	for _, valid := range ValidLibraryTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

// IsValidScanType checks if the given scan type is valid.
func IsValidScanType(t string) bool {
	for _, valid := range ValidScanTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

// IsValidPermission checks if the given permission is valid.
func IsValidPermission(p string) bool {
	for _, valid := range ValidPermissions() {
		if p == valid {
			return true
		}
	}
	return false
}
