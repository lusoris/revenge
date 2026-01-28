// Package shared provides common types and interfaces for all content modules.
package shared

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity contains fields common to all entities.
type BaseEntity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ContentEntity contains fields common to all content items.
type ContentEntity struct {
	BaseEntity
	LibraryID uuid.UUID `json:"library_id"`
	Path      string    `json:"path"`
	Title     string    `json:"title"`
	SortTitle string    `json:"sort_title"`
}

// UserDataEntity contains fields for user-specific data.
type UserDataEntity struct {
	UserID    uuid.UUID `json:"user_id"`
	ItemID    uuid.UUID `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Rating represents a user's rating for an item.
type Rating struct {
	UserDataEntity
	Score float64 `json:"score"` // 1.0 - 10.0
}

// Favorite represents a user's favorite item.
type Favorite struct {
	UserDataEntity
	AddedAt time.Time `json:"added_at"`
}

// WatchHistory represents watch/play history for video content.
type WatchHistory struct {
	UserDataEntity
	PositionTicks int64     `json:"position_ticks"` // Current position in ticks (100ns units)
	Completed     bool      `json:"completed"`
	WatchCount    int       `json:"watch_count"`
	LastWatched   time.Time `json:"last_watched"`
}

// PlayHistory represents play history for audio content.
type PlayHistory struct {
	UserDataEntity
	PlayCount  int       `json:"play_count"`
	LastPlayed time.Time `json:"last_played"`
}

// Image represents an image associated with content.
type Image struct {
	ID         uuid.UUID `json:"id"`
	ItemID     uuid.UUID `json:"item_id"`
	Type       ImageType `json:"type"`
	Language   *string   `json:"language,omitempty"`
	URL        string    `json:"url"`
	LocalPath  *string   `json:"local_path,omitempty"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Blurhash   *string   `json:"blurhash,omitempty"`
	IsPrimary  bool      `json:"is_primary"`
	Source     string    `json:"source"` // "arr", "tmdb", "fanart", "manual"
	CreatedAt  time.Time `json:"created_at"`
}

// ImageType represents the type of image.
type ImageType string

const (
	ImageTypePoster    ImageType = "poster"
	ImageTypeBackdrop  ImageType = "backdrop"
	ImageTypeLogo      ImageType = "logo"
	ImageTypeThumb     ImageType = "thumb"
	ImageTypeBanner    ImageType = "banner"
	ImageTypeDisc      ImageType = "disc"
	ImageTypeClearart  ImageType = "clearart"
	ImageTypeProfile   ImageType = "profile"   // For people/performers
	ImageTypeCover     ImageType = "cover"     // For albums
	ImageTypeArtist    ImageType = "artist"    // For artists
	ImageTypeScreenshot ImageType = "screenshot" // For scenes
)

// Stream represents a media stream (video, audio, subtitle).
type Stream struct {
	ID         uuid.UUID  `json:"id"`
	ItemID     uuid.UUID  `json:"item_id"`
	Type       StreamType `json:"type"`
	Index      int        `json:"index"`
	Codec      string     `json:"codec"`
	Language   *string    `json:"language,omitempty"`
	Title      *string    `json:"title,omitempty"`
	IsDefault  bool       `json:"is_default"`
	IsForced   bool       `json:"is_forced"`
	IsExternal bool       `json:"is_external"`

	// Video-specific
	Width       *int     `json:"width,omitempty"`
	Height      *int     `json:"height,omitempty"`
	Bitrate     *int     `json:"bitrate,omitempty"`
	Framerate   *float64 `json:"framerate,omitempty"`
	AspectRatio *string  `json:"aspect_ratio,omitempty"`

	// Audio-specific
	Channels      *int `json:"channels,omitempty"`
	SampleRate    *int `json:"sample_rate,omitempty"`
	BitDepth      *int `json:"bit_depth,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// StreamType represents the type of media stream.
type StreamType string

const (
	StreamTypeVideo    StreamType = "video"
	StreamTypeAudio    StreamType = "audio"
	StreamTypeSubtitle StreamType = "subtitle"
)

// ExternalID represents an external identifier (e.g., TMDb, IMDb).
type ExternalID struct {
	Provider string `json:"provider"` // "tmdb", "imdb", "tvdb", "musicbrainz", etc.
	ID       string `json:"id"`
}

// Genre represents a genre.
type Genre struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

// Tag represents a user-defined tag.
type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Person represents a person (actor, director, artist, etc.).
type Person struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	SortName   string     `json:"sort_name"`
	Birthdate  *time.Time `json:"birthdate,omitempty"`
	Deathdate  *time.Time `json:"deathdate,omitempty"`
	Birthplace *string    `json:"birthplace,omitempty"`
	Bio        *string    `json:"bio,omitempty"`
	ImagePath  *string    `json:"image_path,omitempty"`
	TmdbID     *int       `json:"tmdb_id,omitempty"`
	ImdbID     *string    `json:"imdb_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// CastMember represents a person's role in content.
type CastMember struct {
	Person    Person  `json:"person"`
	Character *string `json:"character,omitempty"`
	Order     int     `json:"order"`
}

// CrewMember represents a person's crew role in content.
type CrewMember struct {
	Person     Person `json:"person"`
	Department string `json:"department"`
	Job        string `json:"job"`
}

// TicksToDuration converts ticks (100ns units) to time.Duration.
func TicksToDuration(ticks int64) time.Duration {
	return time.Duration(ticks * 100)
}

// DurationToTicks converts time.Duration to ticks (100ns units).
func DurationToTicks(d time.Duration) int64 {
	return int64(d / 100)
}

// GenerateSortTitle generates a sort title from a regular title.
// Removes common prefixes like "The", "A", "An" for better sorting.
func GenerateSortTitle(title string) string {
	prefixes := []string{"the ", "a ", "an ", "der ", "die ", "das ", "ein ", "eine "}
	lower := stringToLower(title)

	for _, prefix := range prefixes {
		if len(lower) > len(prefix) && lower[:len(prefix)] == prefix {
			return title[len(prefix):]
		}
	}

	return title
}

// stringToLower is a simple lowercase function without importing strings.
func stringToLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
