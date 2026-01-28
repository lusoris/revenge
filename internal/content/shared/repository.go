package shared

import (
	"context"

	"github.com/google/uuid"
)

// ListParams contains common parameters for listing items.
type ListParams struct {
	LibraryID *uuid.UUID
	Limit     int32
	Offset    int32
	SortBy    string
	SortOrder string // "asc" or "desc"
}

// FilterParams contains common filter parameters.
type FilterParams struct {
	Genres    []string
	Years     []int
	Tags      []string
	Query     *string // Search query
}

// UserDataParams contains parameters for user-specific queries.
type UserDataParams struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

// Repository is the base interface for all content repositories.
type Repository[T any] interface {
	// GetByID retrieves an item by its unique ID.
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)

	// List retrieves items with pagination.
	List(ctx context.Context, params ListParams) ([]*T, error)

	// ListByLibrary retrieves items in a specific library.
	ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*T, error)

	// Create creates a new item.
	Create(ctx context.Context, item *T) error

	// Update updates an existing item.
	Update(ctx context.Context, item *T) error

	// Delete removes an item by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// Count returns the total number of items.
	Count(ctx context.Context) (int64, error)

	// CountByLibrary returns the number of items in a library.
	CountByLibrary(ctx context.Context, libraryID uuid.UUID) (int64, error)
}

// UserDataRepository is the interface for user-specific data operations.
type UserDataRepository[T any] interface {
	// Get retrieves user data for an item.
	Get(ctx context.Context, userID, itemID uuid.UUID) (*T, error)

	// Set creates or updates user data for an item.
	Set(ctx context.Context, data *T) error

	// Delete removes user data for an item.
	Delete(ctx context.Context, userID, itemID uuid.UUID) error

	// ListByUser retrieves all user data for a user.
	ListByUser(ctx context.Context, userID uuid.UUID, params ListParams) ([]*T, error)
}

// RatingRepository is the interface for rating operations.
type RatingRepository interface {
	GetRating(ctx context.Context, userID, itemID uuid.UUID) (*Rating, error)
	SetRating(ctx context.Context, userID, itemID uuid.UUID, score float64) error
	DeleteRating(ctx context.Context, userID, itemID uuid.UUID) error
	GetAverageRating(ctx context.Context, itemID uuid.UUID) (float64, int, error) // avg, count
}

// FavoriteRepository is the interface for favorite operations.
type FavoriteRepository interface {
	IsFavorite(ctx context.Context, userID, itemID uuid.UUID) (bool, error)
	AddFavorite(ctx context.Context, userID, itemID uuid.UUID) error
	RemoveFavorite(ctx context.Context, userID, itemID uuid.UUID) error
	ListFavorites(ctx context.Context, userID uuid.UUID, params ListParams) ([]uuid.UUID, error)
}

// WatchHistoryRepository is the interface for watch history (video content).
type WatchHistoryRepository interface {
	GetHistory(ctx context.Context, userID, itemID uuid.UUID) (*WatchHistory, error)
	UpdateHistory(ctx context.Context, userID, itemID uuid.UUID, positionTicks int64, completed bool) error
	DeleteHistory(ctx context.Context, userID, itemID uuid.UUID) error
	ListHistory(ctx context.Context, userID uuid.UUID, params ListParams) ([]*WatchHistory, error)
	GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]*WatchHistory, error)
}

// PlayHistoryRepository is the interface for play history (audio content).
type PlayHistoryRepository interface {
	GetHistory(ctx context.Context, userID, itemID uuid.UUID) (*PlayHistory, error)
	IncrementPlayCount(ctx context.Context, userID, itemID uuid.UUID) error
	ListHistory(ctx context.Context, userID uuid.UUID, params ListParams) ([]*PlayHistory, error)
	GetRecentlyPlayed(ctx context.Context, userID uuid.UUID, limit int32) ([]*PlayHistory, error)
}

// ImageRepository is the interface for image operations.
type ImageRepository interface {
	GetImages(ctx context.Context, itemID uuid.UUID) ([]*Image, error)
	GetImagesByType(ctx context.Context, itemID uuid.UUID, imageType ImageType) ([]*Image, error)
	GetPrimaryImage(ctx context.Context, itemID uuid.UUID, imageType ImageType) (*Image, error)
	AddImage(ctx context.Context, image *Image) error
	SetPrimaryImage(ctx context.Context, itemID uuid.UUID, imageID uuid.UUID) error
	DeleteImage(ctx context.Context, imageID uuid.UUID) error
	DeleteImagesByItem(ctx context.Context, itemID uuid.UUID) error
}

// StreamRepository is the interface for media stream operations.
type StreamRepository interface {
	GetStreams(ctx context.Context, itemID uuid.UUID) ([]*Stream, error)
	GetStreamsByType(ctx context.Context, itemID uuid.UUID, streamType StreamType) ([]*Stream, error)
	AddStream(ctx context.Context, stream *Stream) error
	DeleteStreams(ctx context.Context, itemID uuid.UUID) error
}

// SearchableRepository adds search capability.
type SearchableRepository[T any] interface {
	Repository[T]
	Search(ctx context.Context, query string, params ListParams) ([]*T, error)
}
