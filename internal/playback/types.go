// Package playback provides HLS streaming and playback session management.
package playback

import (
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/playback/transcode"
)

// MediaType distinguishes movie files from episode files.
type MediaType string

const (
	MediaTypeMovie   MediaType = "movie"
	MediaTypeEpisode MediaType = "episode"
)

// Session represents an active playback/streaming session.
type Session struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	MediaType         MediaType
	MediaID           uuid.UUID // movie or episode ID
	FileID            uuid.UUID
	FilePath          string // absolute path to the media file
	SegmentDir        string // directory for this session's HLS segments
	TranscodeDecision transcode.Decision
	ActiveProfiles    []string // profile names being generated
	AudioTrack        int
	SubtitleTrack     *int
	StartPosition     int // seconds
	DurationSeconds   float64
	AudioTracks       []AudioTrackInfo
	SubtitleTracks    []SubtitleTrackInfo
	CreatedAt         time.Time
	LastAccessedAt    time.Time
	ExpiresAt         time.Time
}

// StartPlaybackRequest is the input for creating a playback session.
type StartPlaybackRequest struct {
	MediaType     MediaType  `json:"media_type"`
	MediaID       uuid.UUID  `json:"media_id"`
	FileID        *uuid.UUID `json:"file_id,omitempty"`
	AudioTrack    int        `json:"audio_track"`
	SubtitleTrack *int       `json:"subtitle_track,omitempty"`
	StartPosition int        `json:"start_position"` // seconds
}

// PlaybackSessionResponse is the API response for a playback session.
type PlaybackSessionResponse struct {
	SessionID         uuid.UUID           `json:"session_id"`
	MasterPlaylistURL string              `json:"master_playlist_url"`
	DurationSeconds   float64             `json:"duration_seconds"`
	Profiles          []ProfileInfo       `json:"profiles"`
	AudioTracks       []AudioTrackInfo    `json:"audio_tracks"`
	SubtitleTracks    []SubtitleTrackInfo `json:"subtitle_tracks"`
	CreatedAt         time.Time           `json:"created_at"`
	ExpiresAt         time.Time           `json:"expires_at"`
}

// ProfileInfo describes an available quality profile.
type ProfileInfo struct {
	Name       string `json:"name"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Bitrate    int    `json:"bitrate"` // kbps
	IsOriginal bool   `json:"is_original"`
}

// AudioTrackInfo describes an audio track in the media file.
type AudioTrackInfo struct {
	Index     int    `json:"index"`
	Language  string `json:"language"`
	Title     string `json:"title"`
	Channels  int    `json:"channels"`
	Layout    string `json:"layout"`
	Codec     string `json:"codec"`
	IsDefault bool   `json:"is_default"`
}

// SubtitleTrackInfo describes a subtitle track in the media file.
type SubtitleTrackInfo struct {
	Index    int    `json:"index"`
	Language string `json:"language"`
	Title    string `json:"title"`
	Codec    string `json:"codec"`
	URL      string `json:"url"` // WebVTT URL
	IsForced bool   `json:"is_forced"`
}

// AudioTracksFromMediaInfo converts movie.MediaInfo audio streams to AudioTrackInfo.
func AudioTracksFromMediaInfo(info *movie.MediaInfo) []AudioTrackInfo {
	tracks := make([]AudioTrackInfo, len(info.AudioStreams))
	for i, s := range info.AudioStreams {
		tracks[i] = AudioTrackInfo{
			Index:     s.Index,
			Language:  s.Language,
			Title:     s.Title,
			Channels:  s.Channels,
			Layout:    s.Layout,
			Codec:     s.Codec,
			IsDefault: s.IsDefault,
		}
	}
	return tracks
}

// SubtitleTracksFromMediaInfo converts movie.MediaInfo subtitle streams to SubtitleTrackInfo.
func SubtitleTracksFromMediaInfo(info *movie.MediaInfo, sessionID uuid.UUID) []SubtitleTrackInfo {
	tracks := make([]SubtitleTrackInfo, 0, len(info.SubtitleStreams))
	for _, s := range info.SubtitleStreams {
		// Skip bitmap subtitle formats (PGS, VobSub) for MVP
		if isBitmapSubtitle(s.Codec) {
			continue
		}
		tracks = append(tracks, SubtitleTrackInfo{
			Index:    s.Index,
			Language: s.Language,
			Title:    s.Title,
			Codec:    s.Codec,
			URL:      subtitleURL(sessionID, s.Index),
			IsForced: s.IsForced,
		})
	}
	return tracks
}

// isBitmapSubtitle returns true for subtitle codecs that are image-based.
func isBitmapSubtitle(codec string) bool {
	switch codec {
	case "hdmv_pgs_subtitle", "dvd_subtitle", "dvb_subtitle", "pgssub", "vobsub":
		return true
	}
	return false
}

func subtitleURL(sessionID uuid.UUID, trackIndex int) string {
	return "/api/v1/playback/stream/" + sessionID.String() + "/subs/" + itoa(trackIndex) + ".vtt"
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	// Simple int to string for small positive numbers
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}
