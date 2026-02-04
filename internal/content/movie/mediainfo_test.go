package movie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResolutionLabel(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected string
	}{
		{"4K UHD from height", 3840, 2160, "4K UHD"},
		{"4K UHD from width", 4096, 1716, "4K UHD"},
		{"1440p QHD", 2560, 1440, "1440p QHD"},
		{"1080p Full HD", 1920, 1080, "1080p Full HD"},
		{"720p HD", 1280, 720, "720p HD"},
		{"576p SD", 720, 576, "576p SD"},
		{"480p SD", 720, 480, "480p SD"},
		{"Low resolution", 320, 240, "240p"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getResolutionLabel(tt.width, tt.height)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetChannelLayoutName(t *testing.T) {
	tests := []struct {
		channels int
		expected string
	}{
		{1, "mono"},
		{2, "stereo"},
		{6, "5.1"},
		{8, "7.1"},
		{4, "4 channels"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getChannelLayoutName(tt.channels)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMediaInfo_GetDurationFormatted(t *testing.T) {
	tests := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{"Short clip", 65.5, "1:05"},
		{"Movie length", 7260.0, "2:01:00"},
		{"Under a minute", 45.0, "0:45"},
		{"Exactly 1 hour", 3600.0, "1:00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MediaInfo{DurationSeconds: tt.seconds}
			result := m.GetDurationFormatted()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMediaInfo_GetAudioLanguages(t *testing.T) {
	m := &MediaInfo{
		AudioStreams: []AudioStreamInfo{
			{Language: "eng"},
			{Language: "deu"},
			{Language: "eng"}, // Duplicate
			{Language: ""},    // Empty
			{Language: "fra"},
		},
	}

	languages := m.GetAudioLanguages()
	assert.Equal(t, []string{"eng", "deu", "fra"}, languages)
}

func TestMediaInfo_GetSubtitleLanguages(t *testing.T) {
	m := &MediaInfo{
		SubtitleStreams: []SubtitleStreamInfo{
			{Language: "eng"},
			{Language: "deu"},
			{Language: "eng"}, // Duplicate
			{Language: ""},    // Empty
		},
	}

	languages := m.GetSubtitleLanguages()
	assert.Equal(t, []string{"eng", "deu"}, languages)
}

func TestMediaInfo_ToMovieFileInfo(t *testing.T) {
	m := &MediaInfo{
		FilePath:        "/path/to/movie.mkv",
		FileSize:        1024 * 1024 * 500,
		Container:       "matroska",
		Resolution:      "1920x1080",
		ResolutionLabel: "1080p Full HD",
		VideoCodec:      "h264",
		VideoProfile:    "High",
		BitrateKbps:     5000,
		DurationSeconds: 7200,
		Framerate:       23.976,
		DynamicRange:    "SDR",
		ColorSpace:      "BT.709",
		AudioStreams: []AudioStreamInfo{
			{Codec: "aac", Channels: 2, Layout: "stereo", Language: "eng"},
			{Codec: "ac3", Channels: 6, Layout: "5.1", Language: "deu"},
		},
		SubtitleStreams: []SubtitleStreamInfo{
			{Language: "eng"},
			{Language: "deu"},
		},
	}

	info := m.ToMovieFileInfo()

	assert.Equal(t, m.FilePath, info.Path)
	assert.Equal(t, m.Container, info.Container)
	assert.Equal(t, m.Resolution, info.Resolution)
	assert.Equal(t, m.ResolutionLabel, info.ResolutionLabel)
	assert.Equal(t, m.VideoCodec, info.VideoCodec)
	assert.Equal(t, m.VideoProfile, info.VideoProfile)
	assert.Equal(t, int32(m.BitrateKbps), info.BitrateKbps)
	assert.Equal(t, m.DurationSeconds, info.DurationSeconds)
	assert.Equal(t, m.Framerate, info.Framerate)
	assert.Equal(t, m.DynamicRange, info.DynamicRange)
	assert.Equal(t, m.ColorSpace, info.ColorSpace)

	// First audio stream values
	assert.Equal(t, "aac", info.AudioCodec)
	assert.Equal(t, 2, info.AudioChannels)
	assert.Equal(t, "stereo", info.AudioLayout)

	// Languages
	assert.Equal(t, []string{"eng", "deu"}, info.Languages)
	assert.Equal(t, []string{"eng", "deu"}, info.SubtitleLangs)
}

func TestNewMediaInfoProber(t *testing.T) {
	prober := NewMediaInfoProber()
	assert.NotNil(t, prober)
}

// Integration test - skipped unless test video files are available
// To run: go test -run TestMediaInfoProber_Probe_Integration -tags=integration
func TestMediaInfoProber_Probe_Integration(t *testing.T) {
	t.Skip("Skipping integration test - requires test video file and FFmpeg libraries")

	prober := NewMediaInfoProber()
	info, err := prober.Probe("/path/to/test/video.mkv")

	if err != nil {
		t.Fatalf("Failed to probe video: %v", err)
	}

	// Basic assertions for a real video file
	assert.NotEmpty(t, info.VideoCodec)
	assert.Greater(t, info.Width, 0)
	assert.Greater(t, info.Height, 0)
	assert.Greater(t, info.DurationSeconds, 0.0)
}
