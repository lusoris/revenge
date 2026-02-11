package transcode

import (
	"log/slog"
	"os"
	"os/exec"
	"testing"
	"time"

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
		pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
		require.NoError(t, err)
		require.NotNil(t, pm)
		defer pm.Close()

		assert.Equal(t, "ffmpeg", pm.ffmpegPath)
		assert.Equal(t, 6, pm.segmentDuration)
	})

	t.Run("different ffmpeg path", func(t *testing.T) {
		pm, err := NewPipelineManager("/usr/local/bin/ffmpeg", 4, testLogger())
		require.NoError(t, err)
		require.NotNil(t, pm)
		defer pm.Close()

		assert.Equal(t, "/usr/local/bin/ffmpeg", pm.ffmpegPath)
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
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	_, ok := pm.GetProcess(uuid.New(), "original")
	assert.False(t, ok, "should return false for non-existent process")
}

func TestPipelineManager_StopProcess_NotFound(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	// Stopping a non-existent process should not error
	err = pm.StopProcess(uuid.New(), "original")
	assert.NoError(t, err, "stopping non-existent process should return nil")
}

func TestPipelineManager_StopAllForSession_NoProcesses(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	// StopAllForSession with no running processes should not panic
	pm.StopAllForSession(uuid.New())
}

func TestPipelineManager_Close(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)

	// Close should not panic
	pm.Close()
}

// ---------------------------------------------------------------------------
// analyzeProfile — remaining uncovered branch: sized profile with source
// height 0 (unknown height)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// startProcess / StopProcess with real (harmless) processes
// ---------------------------------------------------------------------------

func TestPipelineManager_StartAndStopProcess(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	// Use "sleep 60" so the process stays running long enough to be stopped
	cmd := exec.Command("sleep", "60")
	key := processKey(sessionID, "test-profile")

	proc, err := pm.startProcess(cmd, key, sessionID, "test-profile", "copy", false)
	require.NoError(t, err)
	require.NotNil(t, proc)

	// Process should be retrievable
	got, ok := pm.GetProcess(sessionID, "test-profile")
	assert.True(t, ok)
	assert.Equal(t, proc, got)

	// Stop the process
	err = pm.StopProcess(sessionID, "test-profile")
	assert.NoError(t, err)

	// Process should be removed
	_, ok = pm.GetProcess(sessionID, "test-profile")
	assert.False(t, ok)

	// Done channel should be closed
	select {
	case <-proc.Done:
		// expected
	case <-time.After(5 * time.Second):
		t.Fatal("process Done channel was not closed after stop")
	}
}

func TestPipelineManager_StartProcess_CommandFails(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	// Use a non-existent command to trigger Start() failure
	cmd := exec.Command("/nonexistent/binary/that/does/not/exist")
	key := processKey(sessionID, "fail")

	proc, err := pm.startProcess(cmd, key, sessionID, "fail", "copy", false)
	assert.Error(t, err, "should fail when command cannot start")
	assert.Nil(t, proc)
	assert.Contains(t, err.Error(), "failed to start FFmpeg")
}

func TestPipelineManager_StartProcess_CompletesSuccessfully(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	// Use "true" which exits immediately with success
	cmd := exec.Command("true")
	key := processKey(sessionID, "quick")

	proc, err := pm.startProcess(cmd, key, sessionID, "quick", "libx264", true)
	require.NoError(t, err)
	require.NotNil(t, proc)

	// Wait for process to complete
	select {
	case <-proc.Done:
		// expected
	case <-time.After(5 * time.Second):
		t.Fatal("process did not complete in time")
	}

	assert.NoError(t, proc.Err, "true command should exit with no error")
}

func TestPipelineManager_StopProcess_AlreadyExited(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	// Use "true" which exits immediately
	cmd := exec.Command("true")
	key := processKey(sessionID, "exited")

	proc, err := pm.startProcess(cmd, key, sessionID, "exited", "copy", false)
	require.NoError(t, err)

	// Wait for process to exit naturally
	<-proc.Done

	// Stopping an already-exited process should still work without error
	err = pm.StopProcess(sessionID, "exited")
	assert.NoError(t, err)
}

func TestPipelineManager_StopAllForSession_WithProcesses(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	// Start two processes for the same session in "original" and "720p" profiles
	cmd1 := exec.Command("sleep", "60")
	key1 := processKey(sessionID, "original")
	_, err = pm.startProcess(cmd1, key1, sessionID, "original", "copy", false)
	require.NoError(t, err)

	cmd2 := exec.Command("sleep", "60")
	key2 := processKey(sessionID, "720p")
	_, err = pm.startProcess(cmd2, key2, sessionID, "720p", "copy", false)
	require.NoError(t, err)

	// Both should exist
	_, ok1 := pm.GetProcess(sessionID, "original")
	_, ok2 := pm.GetProcess(sessionID, "720p")
	assert.True(t, ok1)
	assert.True(t, ok2)

	// StopAllForSession should clean them up
	pm.StopAllForSession(sessionID)

	_, ok1 = pm.GetProcess(sessionID, "original")
	_, ok2 = pm.GetProcess(sessionID, "720p")
	assert.False(t, ok1)
	assert.False(t, ok2)
}

// ---------------------------------------------------------------------------
// analyzeProfile — remaining uncovered branch: sized profile with source
// height 0 (unknown height)
// ---------------------------------------------------------------------------

func TestAnalyzeProfile_SizedProfileWithZeroSourceHeight(t *testing.T) {
	// When source height is 0 (unknown), sized profile cannot determine
	// if source is smaller, so it must transcode.
	info := &movie.MediaInfo{
		VideoCodec: "h264",
		Width:      0,
		Height:     0,
		AudioStreams: []movie.AudioStreamInfo{
			{Index: 0, Codec: "aac", Channels: 2},
		},
	}

	profiles := GetEnabledProfiles([]string{"720p"})
	d := AnalyzeMedia(info, profiles)

	require.Len(t, d.Profiles, 1)
	p := d.Profiles[0]

	// Height is 0, so info.Height (0) > 0 is false, falls to else branch
	assert.True(t, p.NeedsTranscode, "should need transcode when source height is unknown")
	assert.Equal(t, "libx264", p.VideoCodec)
	assert.Equal(t, 720, p.Height)
	assert.Equal(t, 1280, p.Width)
}

func TestAnalyzeProfile_SizedProfileSmallerSourceIncompatibleCodec(t *testing.T) {
	// Source is smaller than profile (480p source, 720p profile), but codec
	// is HEVC (not HLS compatible). Video must be transcoded, but dimensions
	// should use source (no upscale). Audio is compatible so canRemuxAudio=true
	// but since NeedsTranscode=true due to video, audio gets transcoded too.
	info := &movie.MediaInfo{
		VideoCodec: "hevc",
		Width:      854,
		Height:     480,
		AudioStreams: []movie.AudioStreamInfo{
			{Index: 0, Codec: "aac", Channels: 2},
		},
	}

	profiles := GetEnabledProfiles([]string{"720p"})
	d := AnalyzeMedia(info, profiles)

	require.Len(t, d.Profiles, 1)
	p := d.Profiles[0]

	assert.True(t, p.NeedsTranscode)
	assert.Equal(t, "libx264", p.VideoCodec)
	// Source dimensions used since source is smaller
	assert.Equal(t, 854, p.Width)
	assert.Equal(t, 480, p.Height)
	assert.Equal(t, 2800, p.VideoBitrate)
	// Audio transcoded because NeedsTranscode is true
	assert.Equal(t, "aac", p.AudioCodec)
	assert.Equal(t, 128, p.AudioBitrate)
}
