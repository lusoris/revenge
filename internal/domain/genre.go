// Package domain contains the core business entities and interfaces.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GenreDomain represents the content type a genre belongs to.
// This ensures genres are scoped to their appropriate media type.
type GenreDomain string

const (
	GenreDomainMovie   GenreDomain = "movie"
	GenreDomainTV      GenreDomain = "tv"
	GenreDomainMusic   GenreDomain = "music"
	GenreDomainBook    GenreDomain = "book"
	GenreDomainPodcast GenreDomain = "podcast"
	GenreDomainGame    GenreDomain = "game"
)

// ValidGenreDomains returns all valid genre domains.
func ValidGenreDomains() []GenreDomain {
	return []GenreDomain{
		GenreDomainMovie,
		GenreDomainTV,
		GenreDomainMusic,
		GenreDomainBook,
		GenreDomainPodcast,
		GenreDomainGame,
	}
}

// IsValid checks if the genre domain is valid.
func (d GenreDomain) IsValid() bool {
	switch d {
	case GenreDomainMovie, GenreDomainTV, GenreDomainMusic,
		GenreDomainBook, GenreDomainPodcast, GenreDomainGame:
		return true
	}
	return false
}

// Genre represents a content genre with domain scoping.
type Genre struct {
	ID          uuid.UUID
	Domain      GenreDomain
	Name        string
	Slug        string  // URL-safe identifier
	Description *string // Optional description
	ParentID    *uuid.UUID
	ExternalIDs map[string]string // Provider IDs (tmdb, musicbrainz, etc.)
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationships (populated on demand)
	Parent   *Genre   // Parent genre (for hierarchical genres)
	Children []*Genre // Child genres
}

// MediaItemGenre represents the association between a media item and a genre.
type MediaItemGenre struct {
	MediaItemID uuid.UUID
	GenreID     uuid.UUID
	Source      string  // "tmdb", "musicbrainz", "manual", "nfo"
	Confidence  float64 // 0.00-1.00 for auto-tagged content
	CreatedAt   time.Time

	// Relationship (populated on demand)
	Genre *Genre
}

// GenreSource constants for tracking where genre assignments came from.
const (
	GenreSourceManual      = "manual"
	GenreSourceTMDB        = "tmdb"
	GenreSourceTVDB        = "tvdb"
	GenreSourceMusicBrainz = "musicbrainz"
	GenreSourceNFO         = "nfo"
	GenreSourceFilename    = "filename"
)

// CreateGenreParams contains parameters for creating a new genre.
type CreateGenreParams struct {
	Domain      GenreDomain
	Name        string
	Slug        string
	Description *string
	ParentID    *uuid.UUID
	ExternalIDs map[string]string
}

// UpdateGenreParams contains parameters for updating a genre.
type UpdateGenreParams struct {
	ID          uuid.UUID
	Name        *string
	Slug        *string
	Description *string
	ParentID    *uuid.UUID
	ExternalIDs map[string]string
}

// AssignGenreParams contains parameters for assigning a genre to a media item.
type AssignGenreParams struct {
	MediaItemID uuid.UUID
	GenreID     uuid.UUID
	Source      string
	Confidence  float64
}

// ListGenresParams contains parameters for listing genres.
type ListGenresParams struct {
	Domain     *GenreDomain // Filter by domain
	ParentID   *uuid.UUID   // Filter by parent (NULL for top-level)
	Query      string       // Search by name
	Limit      int
	Offset     int
	IncludeAll bool // Include all levels, not just top-level
}

// GenreRepository defines the interface for genre data access.
type GenreRepository interface {
	// Single genre operations
	GetByID(ctx context.Context, id uuid.UUID) (*Genre, error)
	GetBySlug(ctx context.Context, domain GenreDomain, slug string) (*Genre, error)

	// List operations
	List(ctx context.Context, params ListGenresParams) ([]*Genre, error)
	ListByDomain(ctx context.Context, domain GenreDomain) ([]*Genre, error)
	ListChildren(ctx context.Context, parentID uuid.UUID) ([]*Genre, error)
	ListForMediaItem(ctx context.Context, mediaItemID uuid.UUID) ([]*Genre, error)

	// Search
	Search(ctx context.Context, domain GenreDomain, query string, limit int) ([]*Genre, error)

	// CRUD
	Create(ctx context.Context, params CreateGenreParams) (*Genre, error)
	Update(ctx context.Context, params UpdateGenreParams) (*Genre, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Media item associations
	AssignToMediaItem(ctx context.Context, params AssignGenreParams) error
	RemoveFromMediaItem(ctx context.Context, mediaItemID, genreID uuid.UUID) error
	RemoveAllFromMediaItem(ctx context.Context, mediaItemID uuid.UUID) error

	// Bulk operations
	BulkAssignToMediaItem(ctx context.Context, mediaItemID uuid.UUID, genreIDs []uuid.UUID, source string) error
}

// GenreService defines the interface for genre business logic.
type GenreService interface {
	// Single genre operations
	GetGenre(ctx context.Context, id uuid.UUID) (*Genre, error)
	GetGenreBySlug(ctx context.Context, domain GenreDomain, slug string) (*Genre, error)

	// List operations
	ListGenres(ctx context.Context, params ListGenresParams) ([]*Genre, error)
	ListGenresByDomain(ctx context.Context, domain GenreDomain) ([]*Genre, error)
	ListGenresForMediaItem(ctx context.Context, mediaItemID uuid.UUID) ([]*Genre, error)
	GetGenreHierarchy(ctx context.Context, domain GenreDomain) ([]*Genre, error)

	// Search
	SearchGenres(ctx context.Context, domain GenreDomain, query string, limit int) ([]*Genre, error)

	// CRUD (admin only)
	CreateGenre(ctx context.Context, params CreateGenreParams) (*Genre, error)
	UpdateGenre(ctx context.Context, params UpdateGenreParams) (*Genre, error)
	DeleteGenre(ctx context.Context, id uuid.UUID) error

	// Media item associations
	AssignGenreToMediaItem(ctx context.Context, mediaItemID, genreID uuid.UUID, source string) error
	RemoveGenreFromMediaItem(ctx context.Context, mediaItemID, genreID uuid.UUID) error
	SetMediaItemGenres(ctx context.Context, mediaItemID uuid.UUID, genreIDs []uuid.UUID, source string) error

	// Domain mapping helpers
	GetDomainForMediaType(mediaType string) GenreDomain
}
