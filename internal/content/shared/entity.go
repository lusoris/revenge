// Package shared provides base types and interfaces for all content modules.
package shared

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity contains fields common to all database entities.
type BaseEntity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ContentEntity contains fields common to all content items (movies, episodes, tracks, etc.).
type ContentEntity struct {
	BaseEntity
	LibraryID uuid.UUID `json:"libraryId"`
	Path      string    `json:"path"`      // Filesystem path
	Title     string    `json:"title"`     // Display title
	SortTitle string    `json:"sortTitle"` // For alphabetical sorting
}

// MediaInfo contains technical information about a media file.
type MediaInfo struct {
	Container    string `json:"container,omitempty"`    // mkv, mp4, avi
	Size         int64  `json:"size,omitempty"`         // File size in bytes
	Bitrate      int    `json:"bitrate,omitempty"`      // Total bitrate in kbps
	DurationTicks int64 `json:"durationTicks,omitempty"` // Duration in ticks (100ns units)
}

// DurationSeconds returns the duration as seconds.
func (m *MediaInfo) DurationSeconds() float64 {
	if m.DurationTicks == 0 {
		return 0
	}
	return float64(m.DurationTicks) / 10_000_000
}

// DurationDuration returns the duration as time.Duration.
func (m *MediaInfo) DurationDuration() time.Duration {
	return time.Duration(m.DurationTicks * 100)
}

// VideoStream represents a video stream in a media file.
type VideoStream struct {
	Index         int     `json:"index"`
	Codec         string  `json:"codec"`         // h264, hevc, av1, vp9
	Profile       string  `json:"profile"`       // main, high, etc.
	Level         string  `json:"level"`         // 4.1, 5.1, etc.
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	AspectRatio   string  `json:"aspectRatio"`   // 16:9, 2.35:1
	Framerate     float64 `json:"framerate"`
	Bitrate       int     `json:"bitrate"`       // kbps
	BitDepth      int     `json:"bitDepth"`      // 8, 10, 12
	ColorSpace    string  `json:"colorSpace"`    // bt709, bt2020
	HDRFormat     string  `json:"hdrFormat"`     // hdr10, dolby_vision, hlg
	IsInterlaced  bool    `json:"isInterlaced"`
	IsDefault     bool    `json:"isDefault"`
}

// AudioStream represents an audio stream in a media file.
type AudioStream struct {
	Index        int    `json:"index"`
	Codec        string `json:"codec"`        // aac, ac3, eac3, dts, flac, opus
	Profile      string `json:"profile"`      // lc, he-aac, etc.
	Channels     int    `json:"channels"`     // 2, 6, 8
	ChannelLayout string `json:"channelLayout"` // stereo, 5.1, 7.1
	SampleRate   int    `json:"sampleRate"`   // 44100, 48000, etc.
	Bitrate      int    `json:"bitrate"`      // kbps
	BitDepth     int    `json:"bitDepth"`     // 16, 24
	Language     string `json:"language"`     // ISO 639-1
	Title        string `json:"title"`        // Track title
	IsDefault    bool   `json:"isDefault"`
	IsForced     bool   `json:"isForced"`
}

// SubtitleStream represents a subtitle stream in a media file.
type SubtitleStream struct {
	Index     int    `json:"index"`
	Codec     string `json:"codec"`    // srt, ass, pgs, vobsub
	Language  string `json:"language"` // ISO 639-1
	Title     string `json:"title"`
	IsDefault bool   `json:"isDefault"`
	IsForced  bool   `json:"isForced"`
	IsExternal bool  `json:"isExternal"` // External subtitle file
	Path      string `json:"path"`       // Path if external
}

// ChapterInfo represents a chapter in a media file.
type ChapterInfo struct {
	Index     int    `json:"index"`
	Title     string `json:"title"`
	StartTicks int64 `json:"startTicks"` // Start position in ticks
	EndTicks   int64 `json:"endTicks"`   // End position in ticks
}

// ImageInfo represents an image associated with content.
type ImageInfo struct {
	ID        uuid.UUID `json:"id"`
	Type      ImageType `json:"type"`
	URL       string    `json:"url"`       // Remote URL or local path
	LocalPath string    `json:"localPath"` // Cached local path
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Blurhash  string    `json:"blurhash"` // Loading placeholder
	Language  string    `json:"language"` // For text-containing images
	IsPrimary bool      `json:"isPrimary"`
	Source    string    `json:"source"` // tmdb, tvdb, manual
}

// ImageType represents the type of image.
type ImageType string

const (
	ImageTypePoster    ImageType = "poster"
	ImageTypeBackdrop  ImageType = "backdrop"
	ImageTypeBanner    ImageType = "banner"
	ImageTypeLogo      ImageType = "logo"
	ImageTypeThumb     ImageType = "thumb"
	ImageTypeDisc      ImageType = "disc"
	ImageTypeArt       ImageType = "art"
	ImageTypeProfile   ImageType = "profile"   // For people
	ImageTypeClearArt  ImageType = "clearart"
	ImageTypeClearLogo ImageType = "clearlogo"
)

// ExternalIDs holds external provider IDs.
type ExternalIDs struct {
	TmdbID       *int       `json:"tmdbId,omitempty"`
	ImdbID       *string    `json:"imdbId,omitempty"`
	TvdbID       *int       `json:"tvdbId,omitempty"`
	MusicBrainzID *uuid.UUID `json:"musicbrainzId,omitempty"`
}

// CreditRole represents the role of a person in a production.
type CreditRole string

const (
	CreditRoleActor       CreditRole = "actor"
	CreditRoleDirector    CreditRole = "director"
	CreditRoleWriter      CreditRole = "writer"
	CreditRoleProducer    CreditRole = "producer"
	CreditRoleComposer    CreditRole = "composer"
	CreditRoleCinematographer CreditRole = "cinematographer"
	CreditRoleEditor      CreditRole = "editor"
	CreditRoleGuest       CreditRole = "guest" // TV show guest
)

// Credit represents a person's involvement in content.
type Credit struct {
	PersonID  uuid.UUID  `json:"personId"`
	Role      CreditRole `json:"role"`
	Character string     `json:"character,omitempty"` // For actors
	Order     int        `json:"order"`               // Billing order
}

// UserItemData represents user-specific data for any content item.
type UserItemData struct {
	UserID      uuid.UUID  `json:"userId"`
	ProfileID   uuid.UUID  `json:"profileId"`
	ItemID      uuid.UUID  `json:"itemId"`

	// Playback state
	PlaybackPositionTicks int64      `json:"playbackPositionTicks"`
	PlayCount             int        `json:"playCount"`
	LastPlayedAt          *time.Time `json:"lastPlayedAt,omitempty"`
	Played                bool       `json:"played"` // Marked as watched/listened

	// User actions
	IsFavorite bool       `json:"isFavorite"`
	Rating     *int       `json:"rating,omitempty"` // 1-10
	RatedAt    *time.Time `json:"ratedAt,omitempty"`
}

// PlayedPercentage calculates the percentage played based on duration.
func (u *UserItemData) PlayedPercentage(durationTicks int64) float64 {
	if durationTicks == 0 {
		return 0
	}
	return float64(u.PlaybackPositionTicks) / float64(durationTicks) * 100
}
