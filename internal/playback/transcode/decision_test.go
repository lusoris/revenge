package transcode

import (
	"testing"

	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeMedia_H264AAC_CanRemux(t *testing.T) {
	info := &movie.MediaInfo{
		VideoCodec:       "h264",
		Width:            1920,
		Height:           1080,
		VideoBitrateKbps: 5000,
		AudioStreams: []movie.AudioStreamInfo{
			{Index: 0, Codec: "aac", Channels: 2, Layout: "stereo"},
		},
	}

	profiles := GetEnabledProfiles([]string{"original", "720p", "480p"})
	d := AnalyzeMedia(info, profiles, nil)

	assert.True(t, d.CanRemux, "H.264+AAC should be remuxable")
	assert.Equal(t, "h264", d.SourceVideoCodec)
	assert.Equal(t, "aac", d.SourceAudioCodec)

	// Original profile should be copy/copy
	require.Len(t, d.Profiles, 3)
	orig := d.Profiles[0]
	assert.Equal(t, "original", orig.Name)
	assert.Equal(t, "copy", orig.VideoCodec)
	assert.Equal(t, "copy", orig.AudioCodec)
	assert.False(t, orig.NeedsTranscode)

	// 720p should need transcode (downscale from 1080p)
	p720 := d.Profiles[1]
	assert.Equal(t, "720p", p720.Name)
	assert.True(t, p720.NeedsTranscode)
	assert.Equal(t, "libx264", p720.VideoCodec)
}

func TestAnalyzeMedia_HEVC_CanRemux(t *testing.T) {
	info := &movie.MediaInfo{
		VideoCodec:       "hevc",
		Width:            3840,
		Height:           2160,
		VideoBitrateKbps: 30000,
		AudioStreams: []movie.AudioStreamInfo{
			{Index: 0, Codec: "eac3", Channels: 6, Layout: "5.1"},
		},
	}

	profiles := GetEnabledProfiles([]string{"original", "1080p", "720p"})
	d := AnalyzeMedia(info, profiles, nil)

	// Video is HEVC (fMP4-compatible) but audio is E-AC-3 (not browser-decodable)
	assert.False(t, d.CanRemux, "HEVC+EAC3 can't fully remux because browsers can't decode EAC3")
	assert.Equal(t, "hevc", d.SourceVideoCodec)

	// Original profile: video copied (HEVC fMP4), audio transcoded to AAC
	orig := d.Profiles[0]
	assert.Equal(t, "original", orig.Name)
	assert.True(t, orig.NeedsTranscode, "must transcode because audio needs AAC conversion")
	assert.Equal(t, "copy", orig.VideoCodec, "HEVC should be copied in fMP4")
	assert.Equal(t, "aac", orig.AudioCodec, "eac3 must be transcoded to AAC for browsers")

	// Sized profiles that are smaller than source: must transcode video (scale down)
	for _, p := range d.Profiles[1:] {
		assert.True(t, p.NeedsTranscode, "sized profile %s should need transcode for scaling", p.Name)
		assert.Equal(t, "libx264", p.VideoCodec, "sized profile uses libx264 for scaling")
		assert.Equal(t, "aac", p.AudioCodec, "sized profile %s should transcode audio", p.Name)
	}
}

func TestAnalyzeMedia_H264DTS_MixedRemux(t *testing.T) {
	info := &movie.MediaInfo{
		VideoCodec:       "h264",
		Width:            1920,
		Height:           1080,
		VideoBitrateKbps: 8000,
		AudioStreams: []movie.AudioStreamInfo{
			{Index: 0, Codec: "dts", Channels: 6, Layout: "5.1"},
		},
	}

	profiles := GetEnabledProfiles([]string{"original"})
	d := AnalyzeMedia(info, profiles, nil)

	assert.False(t, d.CanRemux, "H.264+DTS cannot fully remux")

	// Original profile: video copy, audio transcode
	require.Len(t, d.Profiles, 1)
	orig := d.Profiles[0]
	assert.Equal(t, "copy", orig.VideoCodec, "H.264 video should be copied")
	assert.Equal(t, "aac", orig.AudioCodec, "DTS audio must be transcoded to AAC")
	assert.True(t, orig.NeedsTranscode)
}

func TestAnalyzeMedia_SmallSource_NoUpscale(t *testing.T) {
	info := &movie.MediaInfo{
		VideoCodec: "h264",
		Width:      1280,
		Height:     720,
		AudioStreams: []movie.AudioStreamInfo{
			{Index: 0, Codec: "aac", Channels: 2},
		},
	}

	profiles := GetEnabledProfiles([]string{"original", "1080p", "720p", "480p"})
	d := AnalyzeMedia(info, profiles, nil)

	// 1080p profile should use source dimensions (720p source)
	for _, p := range d.Profiles {
		if p.Name == "1080p" {
			assert.Equal(t, 1280, p.Width, "1080p should use source width for 720p content")
			assert.Equal(t, 720, p.Height, "1080p should use source height for 720p content")
		}
	}
}

func TestAnalyzeMedia_NoAudioStreams(t *testing.T) {
	info := &movie.MediaInfo{
		VideoCodec: "h264",
		Width:      1920,
		Height:     1080,
	}

	profiles := GetEnabledProfiles([]string{"original"})
	d := AnalyzeMedia(info, profiles, nil)

	assert.False(t, d.CanRemux, "no audio means can't fully remux")
	assert.Equal(t, "", d.SourceAudioCodec)
}

func TestEstimateBitrate(t *testing.T) {
	tests := []struct {
		name     string
		info     *movie.MediaInfo
		expected int
	}{
		{"4K", &movie.MediaInfo{Height: 2160}, 15000},
		{"1440p", &movie.MediaInfo{Height: 1440}, 10000},
		{"1080p", &movie.MediaInfo{Height: 1080}, 5000},
		{"720p", &movie.MediaInfo{Height: 720}, 2800},
		{"480p", &movie.MediaInfo{Height: 480}, 1400},
		{"known bitrate", &movie.MediaInfo{Height: 1080, VideoBitrateKbps: 8000}, 8000},
		{"capped bitrate", &movie.MediaInfo{Height: 2160, VideoBitrateKbps: 50000}, 20000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := estimateOriginalBitrate(tt.info)
			assert.Equal(t, tt.expected, got)
		})
	}
}
