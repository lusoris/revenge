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
