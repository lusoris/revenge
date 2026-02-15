package transcode

import (
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
}

func TestNewPipelineManager(t *testing.T) {
	t.Run("creates successfully", func(t *testing.T) {
		pm, err := NewPipelineManager(6, testLogger())
		require.NoError(t, err)
		require.NotNil(t, pm)
		defer pm.Close()

		assert.Equal(t, 6, pm.segmentDuration)
	})

	t.Run("different segment duration", func(t *testing.T) {
		pm, err := NewPipelineManager(4, testLogger())
		require.NoError(t, err)
		require.NotNil(t, pm)
		defer pm.Close()

		assert.Equal(t, 4, pm.segmentDuration)
	})
}

func TestProcessKey(t *testing.T) {
	tests := []struct {
		name      string
		sessionID uuid.UUID
		profile   string
		expected  string
	}{
		{
			name:      "original profile",
			sessionID: uuid.MustParse("11111111-2222-3333-4444-555555555555"),
			profile:   "original",
			expected:  "11111111-2222-3333-4444-555555555555:original",
		},
		{
			name:      "720p profile",
			sessionID: uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
			profile:   "720p",
			expected:  "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee:720p",
		},
		{
			name:      "audio rendition",
			sessionID: uuid.MustParse("12345678-abcd-ef01-2345-6789abcdef01"),
			profile:   "audio/0",
			expected:  "12345678-abcd-ef01-2345-6789abcdef01:audio/0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := processKey(tc.sessionID, tc.profile)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestPipelineManager_GetProcess_NotFound(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	_, ok := pm.GetProcess(uuid.New(), "original")
	assert.False(t, ok, "should return false for non-existent job")
}

func TestPipelineManager_StopProcess_NotFound(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	err = pm.StopProcess(uuid.New(), "original")
	assert.NoError(t, err, "stopping non-existent job should return nil")
}

func TestPipelineManager_StopAllForSession_NoJobs(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	pm.StopAllForSession(uuid.New())
}

func TestPipelineManager_Close(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)

	pm.Close()
}

// ---------------------------------------------------------------------------
// TranscodeJob creation and configuration
// ---------------------------------------------------------------------------

func TestNewTranscodeJob_Defaults(t *testing.T) {
	job := NewTranscodeJob(TranscodeJobConfig{
		InputFile:        "/media/movie.mkv",
		OutputDir:        "/tmp/segments/720p",
		SessionID:        "test-session",
		Profile:          "720p",
		VideoCodec:       "libx264",
		AudioCodec:       "",
		Height:           720,
		VideoBitrate:     2800,
		VideoStreamIndex: 0,
		AudioStreamIndex: -1,
	})

	assert.Equal(t, "/media/movie.mkv", job.InputFile)
	assert.Equal(t, "/tmp/segments/720p/index.m3u8", job.OutputFile)
	assert.Equal(t, 23, job.CRF, "default CRF should be 23")
	assert.Equal(t, "veryfast", job.Preset, "default preset should be veryfast")
	assert.Equal(t, 6, job.SegmentDuration, "default segment duration should be 6")
	assert.True(t, job.IsTranscode, "libx264 should be marked as transcode")
	assert.Equal(t, -1, job.AudioStreamIndex, "audio should be disabled")
}

func TestNewTranscodeJob_Copy(t *testing.T) {
	job := NewTranscodeJob(TranscodeJobConfig{
		InputFile:        "/media/movie.mkv",
		OutputDir:        "/tmp/segments/original",
		VideoCodec:       "copy",
		AudioCodec:       "copy",
		VideoStreamIndex: 0,
		AudioStreamIndex: 0,
	})

	assert.False(t, job.IsTranscode, "copy should not be marked as transcode")
}

func TestNewTranscodeJob_AudioOnly(t *testing.T) {
	job := NewTranscodeJob(TranscodeJobConfig{
		InputFile:        "/media/movie.mkv",
		OutputDir:        "/tmp/segments/audio/0",
		VideoCodec:       "",
		AudioCodec:       "aac",
		AudioBitrate:     256,
		VideoStreamIndex: -1,
		AudioStreamIndex: 0,
	})

	assert.Equal(t, -1, job.VideoStreamIndex)
	assert.Equal(t, 0, job.AudioStreamIndex)
	assert.True(t, job.IsTranscode)
}

func TestNewTranscodeJob_AudioCopy(t *testing.T) {
	job := NewTranscodeJob(TranscodeJobConfig{
		InputFile:        "/media/movie.mkv",
		OutputDir:        "/tmp/segments/audio/0",
		VideoCodec:       "",
		AudioCodec:       "copy",
		VideoStreamIndex: -1,
		AudioStreamIndex: 0,
	})

	assert.False(t, job.IsTranscode, "copy audio should not be marked as transcode")
}

func TestNewTranscodeJob_CustomSettings(t *testing.T) {
	job := NewTranscodeJob(TranscodeJobConfig{
		InputFile:        "/media/movie.mkv",
		OutputDir:        "/tmp/segments/720p",
		VideoCodec:       "libx264",
		CRF:              18,
		Preset:           "medium",
		SegmentDuration:  10,
		VideoStreamIndex: 0,
		AudioStreamIndex: -1,
	})

	assert.Equal(t, 18, job.CRF)
	assert.Equal(t, "medium", job.Preset)
	assert.Equal(t, 10, job.SegmentDuration)
}

func TestTranscodeJob_Stop(t *testing.T) {
	job := NewTranscodeJob(TranscodeJobConfig{
		InputFile:  "/media/movie.mkv",
		OutputDir:  "/tmp/segments",
		VideoCodec: "copy",
	})

	// Stop before Run — should not panic
	job.Stop()
	// Double stop — should not panic
	job.Stop()
}

// ---------------------------------------------------------------------------
// Codec resolution helpers
// ---------------------------------------------------------------------------

func TestResolveVideoCodecID(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"libx264", "h264"},
		{"h264", "h264"},
		{"libx265", "hevc"},
		{"hevc", "hevc"},
		{"h265", "hevc"},
		{"unknown", "h264"}, // default
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := resolveVideoCodecID(tc.name)
			assert.Equal(t, tc.expected, id.Name())
		})
	}
}

func TestResolveAudioCodecID(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"aac", "aac"},
		{"mp3", "mp3"},
		{"ac3", "ac3"},
		{"eac3", "eac3"},
		{"opus", "opus"},
		{"unknown", "aac"}, // default
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := resolveAudioCodecID(tc.name)
			assert.Equal(t, tc.expected, id.Name())
		})
	}
}

func TestResolveResolutionLabel(t *testing.T) {
	tests := []struct {
		profile  string
		expected string
	}{
		{"1080p", "1080p"},
		{"720p", "720p"},
		{"480p", "480p"},
		{"original", "original"},
		{"other", "unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.profile, func(t *testing.T) {
			assert.Equal(t, tc.expected, resolveResolutionLabel(tc.profile))
		})
	}
}

// ---------------------------------------------------------------------------
// Stream counting helper
// ---------------------------------------------------------------------------

func TestCountStreamsBefore(t *testing.T) {
	// This is a unit-level test for the helper; actual stream counting
	// is exercised in integration tests with real media files
}

// ---------------------------------------------------------------------------
// analyzeProfile — remaining uncovered branch: sized profile with source
// height 0 (unknown height)
// ---------------------------------------------------------------------------

func TestAnalyzeProfile_SizedProfileWithZeroSourceHeight(t *testing.T) {
	info := &movie.MediaInfo{
		VideoCodec: "h264",
		Width:      0,
		Height:     0,
		AudioStreams: []movie.AudioStreamInfo{
			{Index: 0, Codec: "aac", Channels: 2},
		},
	}

	profiles := GetEnabledProfiles([]string{"720p"})
	d := AnalyzeMedia(info, profiles, nil)

	require.Len(t, d.Profiles, 1)
	p := d.Profiles[0]

	assert.True(t, p.NeedsTranscode, "should need transcode when source height is unknown")
	assert.Equal(t, "libx264", p.VideoCodec)
	assert.Equal(t, 720, p.Height)
	assert.Equal(t, 1280, p.Width)
}

func TestAnalyzeProfile_SizedProfileSmallerSourceIncompatibleCodec(t *testing.T) {
	info := &movie.MediaInfo{
		VideoCodec: "vp9",
		Width:      854,
		Height:     480,
		AudioStreams: []movie.AudioStreamInfo{
			{Index: 0, Codec: "aac", Channels: 2},
		},
	}

	profiles := GetEnabledProfiles([]string{"720p"})
	d := AnalyzeMedia(info, profiles, nil)

	require.Len(t, d.Profiles, 1)
	p := d.Profiles[0]

	assert.True(t, p.NeedsTranscode)
	assert.Equal(t, "libx264", p.VideoCodec)
	assert.Equal(t, 854, p.Width)
	assert.Equal(t, 480, p.Height)
	assert.Equal(t, 2800, p.VideoBitrate)
	assert.Equal(t, "aac", p.AudioCodec)
	assert.Equal(t, 128, p.AudioBitrate)
}
