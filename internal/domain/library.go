// Package domain contains core business entities and repository interfaces.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LibraryType represents the type of content a library contains.
type LibraryType string

const (
	LibraryTypeMovies      LibraryType = "movies"
	LibraryTypeTVShows     LibraryType = "tvshows"
	LibraryTypeMusic       LibraryType = "music"
	LibraryTypeMusicVideos LibraryType = "musicvideos"
	LibraryTypePhotos      LibraryType = "photos"
	LibraryTypeHomeVideos  LibraryType = "homevideos"
	LibraryTypeBoxSets     LibraryType = "boxsets"
	LibraryTypeLiveTV      LibraryType = "livetv"
	LibraryTypePlaylists   LibraryType = "playlists"
	LibraryTypeMixed       LibraryType = "mixed"
	LibraryTypeBooks       LibraryType = "books"
	LibraryTypeAudiobooks  LibraryType = "audiobooks"
	LibraryTypePodcasts    LibraryType = "podcasts"
	LibraryTypeAdultMovies LibraryType = "adult_movies"
	LibraryTypeAdultShows  LibraryType = "adult_shows"
)

// IsAdultType returns true if this library type is for adult content.
func (t LibraryType) IsAdultType() bool {
	return t == LibraryTypeAdultMovies || t == LibraryTypeAdultShows
}

// IsValid returns true if the library type is a known valid type.
func (t LibraryType) IsValid() bool {
	switch t {
	case LibraryTypeMovies, LibraryTypeTVShows, LibraryTypeMusic, LibraryTypeMusicVideos,
		LibraryTypePhotos, LibraryTypeHomeVideos, LibraryTypeBoxSets, LibraryTypeLiveTV,
		LibraryTypePlaylists, LibraryTypeMixed, LibraryTypeBooks, LibraryTypeAudiobooks,
		LibraryTypePodcasts, LibraryTypeAdultMovies, LibraryTypeAdultShows:
		return true
	default:
		return false
	}
}

// Library represents a media library containing files of a specific type.
type Library struct {
	ID                uuid.UUID
	Name              string
	Type              LibraryType
	Paths             []string           // Filesystem paths to scan
	Settings          map[string]any     // Library-specific settings (JSON)
	IsVisible         bool               // Whether library is visible to users
	IsAdult           bool               // Whether library contains adult content
	ScanIntervalHours *int               // Automatic scan interval, nil = no auto-scan
	LastScanAt        *time.Time         // Last scan timestamp
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// IsAdultLibrary returns true if this library contains adult content.
// This checks both the explicit flag and the library type.
func (l *Library) IsAdultLibrary() bool {
	return l.IsAdult || l.Type.IsAdultType()
}

// CreateLibraryParams contains parameters for creating a new library.
type CreateLibraryParams struct {
	Name              string
	Type              LibraryType
	Paths             []string
	Settings          map[string]any
	IsVisible         bool
	ScanIntervalHours *int
}

// UpdateLibraryParams contains parameters for updating an existing library.
type UpdateLibraryParams struct {
	ID                uuid.UUID
	Name              *string
	Paths             []string // nil means don't update
	Settings          map[string]any
	IsVisible         *bool
	ScanIntervalHours *int
}

// LibraryRepository defines the interface for library data access.
type LibraryRepository interface {
	// GetByID retrieves a library by its unique ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Library, error)

	// GetByName retrieves a library by its name.
	GetByName(ctx context.Context, name string) (*Library, error)

	// List retrieves all libraries.
	List(ctx context.Context) ([]*Library, error)

	// ListByType retrieves libraries of a specific type.
	ListByType(ctx context.Context, libType LibraryType) ([]*Library, error)

	// ListVisible retrieves only visible libraries.
	ListVisible(ctx context.Context) ([]*Library, error)

	// ListNonAdult retrieves libraries that are not adult content.
	ListNonAdult(ctx context.Context) ([]*Library, error)

	// ListForUser retrieves libraries accessible to a specific user.
	// Takes into account adult content settings.
	ListForUser(ctx context.Context, userID uuid.UUID) ([]*Library, error)

	// Create creates a new library and returns the created entity.
	Create(ctx context.Context, params CreateLibraryParams) (*Library, error)

	// Update updates an existing library.
	Update(ctx context.Context, params UpdateLibraryParams) error

	// Delete removes a library by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateLastScan updates the library's last scan timestamp.
	UpdateLastScan(ctx context.Context, id uuid.UUID) error

	// Count returns the total number of libraries.
	Count(ctx context.Context) (int64, error)

	// NameExists checks if a library name is already taken.
	NameExists(ctx context.Context, name string) (bool, error)
}
