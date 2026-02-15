package movie

import (
	"fmt"
	"time"

	"github.com/lusoris/revenge/internal/util"
)

// Prober defines the interface for media probing
type Prober interface {
	Probe(filePath string) (*MediaInfo, error)
}

// MediaInfo contains detailed technical information about a media file
type MediaInfo struct {
	// File metadata
	FilePath  string
	Container string
	FileSize  int64

	// Duration and bitrate
	DurationSeconds float64
	BitrateKbps     int64

	// Video stream info (first video stream)
	VideoCodec       string
	VideoCodecLong   string
	VideoProfile     string
	Width            int
	Height           int
	Resolution       string // e.g., "1920x1080", "4K", "1080p"
	ResolutionLabel  string // e.g., "4K UHD", "1080p Full HD"
	Framerate        float64
	PixelFormat      string
	ColorSpace       string
	ColorPrimaries   string
	ColorTransfer    string
	ColorRange       string
	DynamicRange     string // SDR, HDR10, HDR10+, Dolby Vision, HLG
	VideoCodecString string // RFC 6381 CODECS string, e.g. "hvc1.2.4.L150.90"
	VideoBitrateKbps int64

	// Audio streams
	AudioStreams []AudioStreamInfo

	// Subtitle streams
	SubtitleStreams []SubtitleStreamInfo
}

// AudioStreamInfo contains information about an audio stream
type AudioStreamInfo struct {
	Index       int
	Codec       string
	CodecLong   string
	Channels    int
	Layout      string // e.g., "stereo", "5.1", "7.1"
	SampleRate  int
	BitrateKbps int64
	Language    string
	Title       string
	IsDefault   bool
}

// SubtitleStreamInfo contains information about a subtitle stream
type SubtitleStreamInfo struct {
	Index     int
	Codec     string
	Language  string
	Title     string
	IsForced  bool
	IsDefault bool
}

// ToMovieFileInfo converts MediaInfo to MovieFileInfo
func (m *MediaInfo) ToMovieFileInfo() *MovieFileInfo {
	info := &MovieFileInfo{
		Path:            m.FilePath,
		Size:            m.FileSize,
		Container:       m.Container,
		Resolution:      m.Resolution,
		ResolutionLabel: m.ResolutionLabel,
		VideoCodec:      m.VideoCodec,
		VideoProfile:    m.VideoProfile,
		BitrateKbps:     util.SafeInt64ToInt32(m.BitrateKbps),
		DurationSeconds: m.DurationSeconds,
		Framerate:       m.Framerate,
		DynamicRange:    m.DynamicRange,
		ColorSpace:      m.ColorSpace,
	}

	// Get primary audio codec
	if len(m.AudioStreams) > 0 {
		info.AudioCodec = m.AudioStreams[0].Codec
		info.AudioChannels = m.AudioStreams[0].Channels
		info.AudioLayout = m.AudioStreams[0].Layout
	}

	// Collect all audio languages
	for _, audio := range m.AudioStreams {
		if audio.Language != "" {
			info.Languages = append(info.Languages, audio.Language)
		}
	}

	// Collect all subtitle languages
	for _, sub := range m.SubtitleStreams {
		if sub.Language != "" {
			info.SubtitleLangs = append(info.SubtitleLangs, sub.Language)
		}
	}

	return info
}

// GetAudioLanguages returns all unique audio languages
func (m *MediaInfo) GetAudioLanguages() []string {
	seen := make(map[string]bool)
	var languages []string

	for _, audio := range m.AudioStreams {
		if audio.Language != "" && !seen[audio.Language] {
			seen[audio.Language] = true
			languages = append(languages, audio.Language)
		}
	}

	return languages
}

// GetSubtitleLanguages returns all unique subtitle languages
func (m *MediaInfo) GetSubtitleLanguages() []string {
	seen := make(map[string]bool)
	var languages []string

	for _, sub := range m.SubtitleStreams {
		if sub.Language != "" && !seen[sub.Language] {
			seen[sub.Language] = true
			languages = append(languages, sub.Language)
		}
	}

	return languages
}

// GetDurationFormatted returns a human-readable duration string
func (m *MediaInfo) GetDurationFormatted() string {
	d := time.Duration(m.DurationSeconds * float64(time.Second))

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// getResolutionLabel returns a human-readable resolution label
func getResolutionLabel(width, height int) string {
	// Common resolution labels
	switch {
	case height >= 2160 || width >= 3840:
		return "4K UHD"
	case height >= 1440 || width >= 2560:
		return "1440p QHD"
	case height >= 1080 || width >= 1920:
		return "1080p Full HD"
	case height >= 720 || width >= 1280:
		return "720p HD"
	case height >= 576:
		return "576p SD"
	case height >= 480:
		return "480p SD"
	default:
		return fmt.Sprintf("%dp", height)
	}
}

// getChannelLayoutName returns a human-readable channel layout name
func getChannelLayoutName(channels int) string {
	switch channels {
	case 1:
		return "mono"
	case 2:
		return "stereo"
	case 6:
		return "5.1"
	case 8:
		return "7.1"
	default:
		return fmt.Sprintf("%d channels", channels)
	}
}
