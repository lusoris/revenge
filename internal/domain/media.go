// Package domain contains core business entities and repository interfaces.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// MediaType represents the type of a media item.
type MediaType string

const (
	// Video types
	MediaTypeMovie      MediaType = "movie"
	MediaTypeEpisode    MediaType = "episode"
	MediaTypeMusicVideo MediaType = "musicvideo"
	MediaTypeTrailer    MediaType = "trailer"
	MediaTypeHomeVideo  MediaType = "homevideo"

	// Audio types
	MediaTypeAudio           MediaType = "audio"
	MediaTypeAudiobookChapter MediaType = "audiobook_chapter"
	MediaTypePodcastEpisode  MediaType = "podcast_episode"

	// Image types
	MediaTypePhoto MediaType = "photo"

	// Document types
	MediaTypeBook MediaType = "book"

	// Collection/container types
	MediaTypeSeries     MediaType = "series"
	MediaTypeSeason     MediaType = "season"
	MediaTypeAlbum      MediaType = "album"
	MediaTypeAudiobook  MediaType = "audiobook"
	MediaTypePodcast    MediaType = "podcast"
	MediaTypeArtist     MediaType = "artist"
	MediaTypeBoxSet     MediaType = "boxset"
	MediaTypePlaylist   MediaType = "playlist"
	MediaTypeFolder     MediaType = "folder"
	MediaTypePhotoAlbum MediaType = "photo_album"

	// Live TV types
	MediaTypeChannel   MediaType = "channel"
	MediaTypeProgram   MediaType = "program"
	MediaTypeRecording MediaType = "recording"
)

// IsValid returns true if the media type is a known valid type.
func (t MediaType) IsValid() bool {
	switch t {
	case MediaTypeMovie, MediaTypeEpisode, MediaTypeMusicVideo, MediaTypeTrailer, MediaTypeHomeVideo,
		MediaTypeAudio, MediaTypeAudiobookChapter, MediaTypePodcastEpisode,
		MediaTypePhoto, MediaTypeBook,
		MediaTypeSeries, MediaTypeSeason, MediaTypeAlbum, MediaTypeAudiobook, MediaTypePodcast,
		MediaTypeArtist, MediaTypeBoxSet, MediaTypePlaylist, MediaTypeFolder, MediaTypePhotoAlbum,
		MediaTypeChannel, MediaTypeProgram, MediaTypeRecording:
		return true
	default:
		return false
	}
}

// IsContainer returns true if this media type is a container for other items.
func (t MediaType) IsContainer() bool {
	switch t {
	case MediaTypeSeries, MediaTypeSeason, MediaTypeAlbum, MediaTypeAudiobook, MediaTypePodcast,
		MediaTypeArtist, MediaTypeBoxSet, MediaTypePlaylist, MediaTypeFolder, MediaTypePhotoAlbum:
		return true
	default:
		return false
	}
}

// IsPlayable returns true if this media type is directly playable.
func (t MediaType) IsPlayable() bool {
	switch t {
	case MediaTypeMovie, MediaTypeEpisode, MediaTypeMusicVideo, MediaTypeTrailer, MediaTypeHomeVideo,
		MediaTypeAudio, MediaTypeAudiobookChapter, MediaTypePodcastEpisode, MediaTypeRecording:
		return true
	default:
		return false
	}
}

// MediaItem represents a media item in a library.
type MediaItem struct {
	ID        uuid.UUID
	LibraryID uuid.UUID
	ParentID  *uuid.UUID // Parent item (e.g., series for episode)
	Type      MediaType
	Name      string
	SortName  *string
	Path      *string // Filesystem path (nil for virtual items)

	// Common metadata
	Overview     *string
	Tagline      *string
	Year         *int
	PremiereDate *time.Time
	EndDate      *time.Time
	RuntimeTicks *int64 // Duration in ticks (100ns units)

	// Series/Episode specific
	SeasonNumber          *int
	EpisodeNumber         *int
	AbsoluteEpisodeNumber *int

	// Music specific
	AlbumArtist *string
	TrackNumber *int
	DiscNumber  *int

	// Ratings
	CommunityRating *float64 // e.g., 8.5
	CriticRating    *float64

	// External IDs
	ProviderIDs map[string]string // e.g., {"imdb": "tt123", "tmdb": "456"}

	// Metadata arrays
	Genres  []string
	Tags    []string
	Studios []string

	// File info (for actual media files)
	Container  *string // mkv, mp4, etc.
	VideoCodec *string
	AudioCodec *string
	Width      *int
	Height     *int
	Bitrate    *int

	// Timestamps
	DateCreated  *time.Time // File creation date
	DateModified *time.Time // File modification date
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RuntimeDuration returns the runtime as a time.Duration.
func (m *MediaItem) RuntimeDuration() time.Duration {
	if m.RuntimeTicks == nil {
		return 0
	}
	// Ticks are in 100ns units
	return time.Duration(*m.RuntimeTicks * 100)
}

// CreateMediaItemParams contains parameters for creating a new media item.
type CreateMediaItemParams struct {
	LibraryID   uuid.UUID
	ParentID    *uuid.UUID
	Type        MediaType
	Name        string
	SortName    *string
	Path        *string
	Overview    *string
	Year        *int
	ProviderIDs map[string]string
	Genres      []string
	Tags        []string
	Studios     []string
}

// MediaItemRepository defines the interface for media item data access.
type MediaItemRepository interface {
	// GetByID retrieves a media item by its unique ID.
	GetByID(ctx context.Context, id uuid.UUID) (*MediaItem, error)

	// GetByPath retrieves a media item by its filesystem path.
	GetByPath(ctx context.Context, path string) (*MediaItem, error)

	// List retrieves media items with pagination and filtering.
	List(ctx context.Context, params ListMediaItemsParams) ([]*MediaItem, error)

	// ListByLibrary retrieves media items in a library.
	ListByLibrary(ctx context.Context, libraryID uuid.UUID, limit, offset int32) ([]*MediaItem, error)

	// ListByParent retrieves children of a media item.
	ListByParent(ctx context.Context, parentID uuid.UUID) ([]*MediaItem, error)

	// ListByType retrieves media items of a specific type.
	ListByType(ctx context.Context, mediaType MediaType, limit, offset int32) ([]*MediaItem, error)

	// Create creates a new media item and returns the created entity.
	Create(ctx context.Context, params CreateMediaItemParams) (*MediaItem, error)

	// Update updates an existing media item.
	Update(ctx context.Context, item *MediaItem) error

	// Delete removes a media item by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByLibrary removes all media items in a library.
	DeleteByLibrary(ctx context.Context, libraryID uuid.UUID) error

	// Count returns the total number of media items.
	Count(ctx context.Context) (int64, error)

	// CountByLibrary returns the number of media items in a library.
	CountByLibrary(ctx context.Context, libraryID uuid.UUID) (int64, error)

	// Search searches for media items by name.
	Search(ctx context.Context, query string, limit int32) ([]*MediaItem, error)
}

// ListMediaItemsParams contains parameters for listing media items.
type ListMediaItemsParams struct {
	LibraryID    *uuid.UUID
	ParentID     *uuid.UUID
	MediaType    *MediaType
	Genres       []string
	Years        []int
	Tags         []string
	SortBy       string
	SortOrder    string // "asc" or "desc"
	Limit        int32
	Offset       int32

	// Content filtering
	MaxRatingLevel int  // Maximum normalized rating level
	IncludeAdult   bool // Whether to include adult content
	UserID         *uuid.UUID // User for personalized filtering
}
