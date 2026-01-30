// Package playback provides playback session management and progress tracking.
package playback

import (
	"time"

	"github.com/google/uuid"
)

// MediaType represents the type of media being played.
type MediaType string

const (
	MediaTypeMovie   MediaType = "movie"
	MediaTypeEpisode MediaType = "episode"
	MediaTypeAdult   MediaType = "adult"
)

// PlaybackSession represents an active playback session.
type PlaybackSession struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"userId"`
	MediaID        uuid.UUID  `json:"mediaId"`
	MediaType      MediaType  `json:"mediaType"`
	PositionTicks  int64      `json:"positionTicks"`
	RuntimeTicks   int64      `json:"runtimeTicks"`
	PlayedPercent  float64    `json:"playedPercent"`
	IsPaused       bool       `json:"isPaused"`
	StartedAt      time.Time  `json:"startedAt"`
	LastActivityAt time.Time  `json:"lastActivityAt"`
	DeviceID       *string    `json:"deviceId,omitempty"`
	DeviceName     *string    `json:"deviceName,omitempty"`
	ClientName     *string    `json:"clientName,omitempty"`
}

// StartPlaybackParams contains parameters for starting playback.
type StartPlaybackParams struct {
	UserID        uuid.UUID
	MediaID       uuid.UUID
	MediaType     MediaType
	RuntimeTicks  int64
	PositionTicks int64 // Optional: resume position
	DeviceID      *string
	DeviceName    *string
	ClientName    *string
}

// UpdateProgressParams contains parameters for updating playback progress.
type UpdateProgressParams struct {
	SessionID     uuid.UUID
	PositionTicks int64
	IsPaused      bool
}

// UpNextItem represents an item in the up-next queue.
type UpNextItem struct {
	MediaID     uuid.UUID `json:"mediaId"`
	MediaType   MediaType `json:"mediaType"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle,omitempty"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	RuntimeTicks int64    `json:"runtimeTicks"`
}

// UpNextQueue represents the queue of items to play next.
type UpNextQueue struct {
	Items       []UpNextItem `json:"items"`
	CurrentItem *UpNextItem  `json:"currentItem,omitempty"`
}
