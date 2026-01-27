// Package domain contains core business entities and repository interfaces.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RatingSystem represents an international content rating system (MPAA, FSK, BBFC, etc.).
type RatingSystem struct {
	ID           uuid.UUID
	Code         string   // Unique code like 'mpaa', 'fsk', 'bbfc'
	Name         string   // Full name like 'Motion Picture Association'
	CountryCodes []string // ISO country codes where this system is used
	IsActive     bool     // Whether this system is active/enabled
	SortOrder    int      // Display order
	CreatedAt    time.Time
}

// Rating represents a specific rating within a rating system.
type Rating struct {
	ID              uuid.UUID
	SystemID        uuid.UUID    // FK to RatingSystem
	System          *RatingSystem // Populated on join
	Code            string       // Rating code like 'PG-13', 'FSK 16'
	Name            string       // Full name like 'Parental Guidance 13'
	Description     *string      // Optional description
	MinAge          *int         // Minimum age for this rating
	NormalizedLevel int          // 0-100 scale for cross-system comparison
	SortOrder       int          // Display order within system
	IsAdult         bool         // Whether this is explicit adult content
	IconURL         *string      // URL to rating icon
	CreatedAt       time.Time
}

// ContentRating represents a rating assigned to a piece of content.
type ContentRating struct {
	ID          uuid.UUID
	ContentID   uuid.UUID // FK to the rated content (media_item, image, etc.)
	ContentType string    // Type of content: 'media_item', 'image', 'person_image'
	RatingID    uuid.UUID // FK to Rating
	Rating      *Rating   // Populated on join
	Source      *string   // Source of rating: 'tmdb', 'manual', 'imdb'
	CreatedAt   time.Time
}

// ContentMinRatingLevel represents the minimum (most restrictive) rating level for content.
// This is typically from the materialized view for efficient filtering.
type ContentMinRatingLevel struct {
	ContentID   uuid.UUID
	ContentType string
	MinLevel    int  // Minimum normalized level (most restrictive)
	IsAdult     bool // Whether any rating is adult
}

// RatingSystemRepository defines the interface for rating system data access.
type RatingSystemRepository interface {
	// GetByID retrieves a rating system by its unique ID.
	GetByID(ctx context.Context, id uuid.UUID) (*RatingSystem, error)

	// GetByCode retrieves a rating system by its code.
	GetByCode(ctx context.Context, code string) (*RatingSystem, error)

	// List retrieves all active rating systems.
	List(ctx context.Context) ([]*RatingSystem, error)

	// ListByCountry retrieves rating systems for a specific country.
	ListByCountry(ctx context.Context, countryCode string) ([]*RatingSystem, error)
}

// RatingRepository defines the interface for rating data access.
type RatingRepository interface {
	// GetByID retrieves a rating by its unique ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Rating, error)

	// GetBySystemAndCode retrieves a rating by system ID and code.
	GetBySystemAndCode(ctx context.Context, systemID uuid.UUID, code string) (*Rating, error)

	// ListBySystem retrieves all ratings for a rating system.
	ListBySystem(ctx context.Context, systemID uuid.UUID) ([]*Rating, error)

	// ListByNormalizedLevel retrieves ratings at or below a normalized level.
	ListByNormalizedLevel(ctx context.Context, maxLevel int) ([]*Rating, error)

	// GetEquivalents retrieves equivalent ratings from other systems.
	GetEquivalents(ctx context.Context, ratingID uuid.UUID) ([]*Rating, error)
}

// ContentRatingRepository defines the interface for content rating data access.
type ContentRatingRepository interface {
	// GetByContent retrieves all ratings for a piece of content.
	GetByContent(ctx context.Context, contentID uuid.UUID, contentType string) ([]*ContentRating, error)

	// GetMinLevel retrieves the minimum (most restrictive) rating level for content.
	GetMinLevel(ctx context.Context, contentID uuid.UUID, contentType string) (*ContentMinRatingLevel, error)

	// GetDisplayRating retrieves the rating to display for content in a preferred system.
	// Falls back to common systems if preferred is not available.
	GetDisplayRating(ctx context.Context, contentID uuid.UUID, contentType string, preferredSystem string) (*ContentRating, error)

	// Create adds a rating to content.
	Create(ctx context.Context, contentID uuid.UUID, contentType string, ratingID uuid.UUID, source *string) (*ContentRating, error)

	// Delete removes a rating from content.
	Delete(ctx context.Context, contentID uuid.UUID, ratingID uuid.UUID) error

	// DeleteAllForContent removes all ratings from content.
	DeleteAllForContent(ctx context.Context, contentID uuid.UUID, contentType string) error

	// IsContentAllowed checks if content is allowed for a user's rating level.
	IsContentAllowed(ctx context.Context, contentID uuid.UUID, contentType string, maxLevel int, includeAdult bool) (bool, error)

	// FilterAllowedContent filters a list of content IDs to only those allowed.
	FilterAllowedContent(ctx context.Context, contentIDs []uuid.UUID, contentType string, maxLevel int, includeAdult bool) ([]uuid.UUID, error)
}
