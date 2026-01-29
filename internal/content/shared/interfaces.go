// Package shared provides base types and interfaces for all content modules.
package shared

import (
	"context"

	"github.com/google/uuid"
)

// Repository is the base interface for all content repositories.
type Repository[T any] interface {
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ListParams contains common pagination parameters.
type ListParams struct {
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string // "asc" or "desc"
}

// FilterParams contains common filter parameters.
type FilterParams struct {
	LibraryID      *uuid.UUID
	GenreIDs       []uuid.UUID
	Years          []int
	MinRating      *float64
	MaxRatingLevel int  // Content rating level (0-100)
	IncludeAdult   bool // Only relevant if user has adult access
}

// SearchParams contains parameters for search operations.
type SearchParams struct {
	Query  string
	Limit  int
	Offset int
	Types  []string // Content types to search
}

// ContentService is the base interface for all content services.
type ContentService interface {
	// ScanLibrary triggers a library scan.
	ScanLibrary(ctx context.Context, libraryID uuid.UUID, fullScan bool) error
}

// MetadataProvider is the base interface for metadata providers.
type MetadataProvider interface {
	// Name returns the provider name (e.g., "tmdb", "musicbrainz").
	Name() string

	// Priority returns the provider priority (lower = higher priority).
	Priority() int

	// IsAvailable checks if the provider is configured and available.
	IsAvailable() bool
}

// Scanner is the interface for library scanners.
type Scanner interface {
	// Scan scans a library for new/changed content.
	Scan(ctx context.Context, libraryID uuid.UUID, paths []string, fullScan bool) error

	// ScanFile scans a single file.
	ScanFile(ctx context.Context, libraryID uuid.UUID, path string) error
}

// Indexer is the interface for search indexers.
type Indexer interface {
	// Index adds or updates an item in the search index.
	Index(ctx context.Context, item any) error

	// Delete removes an item from the search index.
	Delete(ctx context.Context, itemID uuid.UUID) error

	// Reindex rebuilds the entire search index for a module.
	Reindex(ctx context.Context) error
}

// LibraryProvider is the interface for module-specific library management.
// Each content module (movie, tvshow, music, qar) implements this interface
// to provide unified library operations.
type LibraryProvider interface {
	// ModuleName returns the module identifier (e.g., "movie", "tvshow", "qar").
	ModuleName() string

	// ListLibraries returns all libraries for this module accessible by the user.
	ListLibraries(ctx context.Context, userID uuid.UUID) ([]LibraryInfo, error)

	// GetLibrary returns a specific library by ID.
	GetLibrary(ctx context.Context, libraryID uuid.UUID) (*LibraryInfo, error)

	// CreateLibrary creates a new library for this module.
	CreateLibrary(ctx context.Context, req CreateLibraryRequest) (*LibraryInfo, error)

	// UpdateLibrary updates library settings.
	UpdateLibrary(ctx context.Context, libraryID uuid.UUID, req UpdateLibraryRequest) (*LibraryInfo, error)

	// DeleteLibrary removes a library and optionally its content.
	DeleteLibrary(ctx context.Context, libraryID uuid.UUID, deleteContent bool) error

	// ScanLibrary triggers a library scan.
	ScanLibrary(ctx context.Context, libraryID uuid.UUID, fullScan bool) error
}

// LibraryInfo is the common representation of a library across all modules.
type LibraryInfo struct {
	ID        uuid.UUID `json:"id"`
	Module    string    `json:"module"`    // "movie", "tvshow", "music", "qar"
	Name      string    `json:"name"`
	Paths     []string  `json:"paths"`
	IsAdult   bool      `json:"is_adult"`
	ItemCount int64     `json:"item_count"`
	Settings  any       `json:"settings,omitempty"` // Module-specific settings
}

// CreateLibraryRequest contains parameters for creating a library.
type CreateLibraryRequest struct {
	Name     string   `json:"name"`
	Paths    []string `json:"paths"`
	Settings any      `json:"settings,omitempty"` // Module-specific settings
}

// UpdateLibraryRequest contains parameters for updating a library.
type UpdateLibraryRequest struct {
	Name     *string  `json:"name,omitempty"`
	Paths    []string `json:"paths,omitempty"`
	Settings any      `json:"settings,omitempty"`
}
