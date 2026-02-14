package transcode

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===========================================================================
// StartVideoSegmenting — directory creation
// ===========================================================================

func TestPipelineManager_StartVideoSegmenting_CreatesDir(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	pd := ProfileDecision{
		Name:       "720p",
		Width:      1280,
		Height:     720,
		VideoCodec: "libx264",
	}

	// Using /dev/null as input will cause the job to fail, but the
	// directory creation should still happen before the error.
	_, err = pm.StartVideoSegmenting(context.Background(), sessionID, "/dev/null", segDir, pd, 0)
	// Error may or may not occur depending on whether the job goroutine
	// started before we check; the directory is created synchronously.

	// Profile directory should have been created
	profileDir := filepath.Join(segDir, "720p")
	info, statErr := os.Stat(profileDir)
	_ = err // job start may fail but directory should exist
	require.NoError(t, statErr, "profile directory should be created")
	assert.True(t, info.IsDir())
}

func TestPipelineManager_StartVideoSegmenting_WithSeek(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	pd := ProfileDecision{
		Name:       "original",
		VideoCodec: "copy",
	}

	_, _ = pm.StartVideoSegmenting(context.Background(), sessionID, "/dev/null", segDir, pd, 120)

	// Verify directory creation even though the command fails
	_, statErr := os.Stat(filepath.Join(segDir, "original"))
	assert.NoError(t, statErr)
}

// ===========================================================================
// StartAudioRendition — directory creation
// ===========================================================================

func TestPipelineManager_StartAudioRendition_CreatesDir(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	_, _ = pm.StartAudioRendition(context.Background(), sessionID, "/dev/null", segDir, 0, "aac", 256, 0)

	audioDir := filepath.Join(segDir, "audio", "0")
	info, statErr := os.Stat(audioDir)
	require.NoError(t, statErr, "audio rendition directory should be created")
	assert.True(t, info.IsDir())
}

func TestPipelineManager_StartAudioRendition_MultipleTracks(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	for i := range 3 {
		_, _ = pm.StartAudioRendition(context.Background(), sessionID, "/dev/null", segDir, i, "copy", 0, 0)
	}

	// All three directories should exist
	for i := range 3 {
		dir := filepath.Join(segDir, "audio", itoa(i))
		_, statErr := os.Stat(dir)
		assert.NoError(t, statErr, "audio dir for track %d should exist", i)
	}
}

func TestPipelineManager_StartAudioRendition_WithSeek(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	_, _ = pm.StartAudioRendition(context.Background(), sessionID, "/dev/null", segDir, 1, "aac", 128, 300)

	audioDir := filepath.Join(segDir, "audio", "1")
	_, statErr := os.Stat(audioDir)
	assert.NoError(t, statErr)
}

// ===========================================================================
// StopAllForSession — stops all jobs
// ===========================================================================

func TestPipelineManager_StopAllForSession_WithJobs(t *testing.T) {
	pm, err := NewPipelineManager(6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	// Start video jobs for two profiles
	pd1 := ProfileDecision{Name: "original", VideoCodec: "copy"}
	pd2 := ProfileDecision{Name: "720p", Width: 1280, Height: 720, VideoCodec: "libx264"}
	_, _ = pm.StartVideoSegmenting(context.Background(), sessionID, "/dev/null", segDir, pd1, 0)
	_, _ = pm.StartVideoSegmenting(context.Background(), sessionID, "/dev/null", segDir, pd2, 0)

	// Allow jobs to start (they'll fail on /dev/null but the cache entry exists briefly)
	time.Sleep(50 * time.Millisecond)

	// StopAllForSession should not panic even if jobs already finished
	pm.StopAllForSession(sessionID)
}

// ===========================================================================
// Helpers
// ===========================================================================

func itoa(i int) string {
	switch {
	case i == 0:
		return "0"
	case i < 10:
		return string(rune('0' + i))
	default:
		s := ""
		for i > 0 {
			s = string(rune('0'+i%10)) + s
			i /= 10
		}
		return s
	}
}
