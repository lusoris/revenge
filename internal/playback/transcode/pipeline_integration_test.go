package transcode

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===========================================================================
// StartVideoSegmenting — directory creation and command construction
// ===========================================================================

func TestPipelineManager_StartVideoSegmenting_CreatesDir(t *testing.T) {
	pm, err := NewPipelineManager("sleep", 6, testLogger()) // use sleep so process runs
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

	// StartVideoSegmenting with "sleep" as ffmpeg binary — it'll fail immediately
	// but should still create the profile directory
	proc, err := pm.StartVideoSegmenting(context.Background(), sessionID, "/dev/null", segDir, pd, 0)
	if err == nil && proc != nil {
		// Process started — wait for it to finish (sleep with bad args exits quickly)
		<-proc.Done
	}

	// Profile directory should have been created
	profileDir := filepath.Join(segDir, "720p")
	info, statErr := os.Stat(profileDir)
	require.NoError(t, statErr, "profile directory should be created")
	assert.True(t, info.IsDir())
}

func TestPipelineManager_StartVideoSegmenting_WithSeek(t *testing.T) {
	pm, err := NewPipelineManager("sleep", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	pd := ProfileDecision{
		Name:       "original",
		VideoCodec: "copy",
	}

	proc, err := pm.StartVideoSegmenting(context.Background(), sessionID, "/dev/null", segDir, pd, 120)
	if err == nil && proc != nil {
		<-proc.Done
	}

	// Verify directory creation even though the command fails
	_, statErr := os.Stat(filepath.Join(segDir, "original"))
	assert.NoError(t, statErr)
}

// ===========================================================================
// StartAudioRendition — directory creation
// ===========================================================================

func TestPipelineManager_StartAudioRendition_CreatesDir(t *testing.T) {
	pm, err := NewPipelineManager("sleep", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	proc, err := pm.StartAudioRendition(context.Background(), sessionID, "/dev/null", segDir, 0, "aac", 256, 0)
	if err == nil && proc != nil {
		<-proc.Done
	}

	audioDir := filepath.Join(segDir, "audio", "0")
	info, statErr := os.Stat(audioDir)
	require.NoError(t, statErr, "audio rendition directory should be created")
	assert.True(t, info.IsDir())
}

func TestPipelineManager_StartAudioRendition_MultipleTrack(t *testing.T) {
	pm, err := NewPipelineManager("sleep", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	for i := range 3 {
		proc, err := pm.StartAudioRendition(context.Background(), sessionID, "/dev/null", segDir, i, "copy", 0, 0)
		if err == nil && proc != nil {
			<-proc.Done
		}
	}

	// All three directories should exist
	for i := range 3 {
		dir := filepath.Join(segDir, "audio", itoa(i))
		_, statErr := os.Stat(dir)
		assert.NoError(t, statErr, "audio dir for track %d should exist", i)
	}
}

func TestPipelineManager_StartAudioRendition_WithSeek(t *testing.T) {
	pm, err := NewPipelineManager("sleep", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()
	segDir := t.TempDir()

	proc, err := pm.StartAudioRendition(context.Background(), sessionID, "/dev/null", segDir, 1, "aac", 128, 300)
	if err == nil && proc != nil {
		<-proc.Done
	}

	audioDir := filepath.Join(segDir, "audio", "1")
	_, statErr := os.Stat(audioDir)
	assert.NoError(t, statErr)
}

// ===========================================================================
// StopAllForSession — covers audio rendition cleanup
// ===========================================================================

func TestPipelineManager_StopAllForSession_AudioRenditions(t *testing.T) {
	pm, err := NewPipelineManager("sleep", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	// Start multiple sleep processes under audio rendition names
	for i := range 4 {
		key := processKey(sessionID, "audio/"+itoa(i))
		cmd := newSleepCmd()
		_, startErr := pm.startProcess(cmd, key, sessionID, "audio/"+itoa(i), "copy", false)
		require.NoError(t, startErr)
	}

	// Verify they exist
	for i := range 4 {
		_, ok := pm.GetProcess(sessionID, "audio/"+itoa(i))
		assert.True(t, ok, "audio/%d should exist before StopAll", i)
	}

	pm.StopAllForSession(sessionID)

	// All should be gone
	for i := range 4 {
		_, ok := pm.GetProcess(sessionID, "audio/"+itoa(i))
		assert.False(t, ok, "audio/%d should be gone after StopAll", i)
	}
}

// ===========================================================================
// Transcode metrics — process marked as transcode
// ===========================================================================

func TestPipelineManager_TranscodeMetric(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	// Start a "true" process marked as transcode to verify metric codepath
	cmd := newTrueCmd()
	key := processKey(sessionID, "720p")
	proc, err := pm.startProcess(cmd, key, sessionID, "720p", "libx264", true)
	require.NoError(t, err)
	require.NotNil(t, proc)
	assert.True(t, proc.IsTranscode)

	// Wait for it to complete
	<-proc.Done
	assert.NoError(t, proc.Err)
}

func TestPipelineManager_TranscodeMetric_Original(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	cmd := newTrueCmd()
	key := processKey(sessionID, "original")
	proc, err := pm.startProcess(cmd, key, sessionID, "original", "libx264", true)
	require.NoError(t, err)

	<-proc.Done
	assert.NoError(t, proc.Err)
}

func TestPipelineManager_TranscodeMetric_1080p(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	cmd := newTrueCmd()
	key := processKey(sessionID, "1080p")
	proc, err := pm.startProcess(cmd, key, sessionID, "1080p", "libx264", true)
	require.NoError(t, err)

	<-proc.Done
	assert.NoError(t, proc.Err)
}

func TestPipelineManager_TranscodeMetric_480p(t *testing.T) {
	pm, err := NewPipelineManager("ffmpeg", 6, testLogger())
	require.NoError(t, err)
	defer pm.Close()

	sessionID := uuid.New()

	cmd := newTrueCmd()
	key := processKey(sessionID, "480p")
	proc, err := pm.startProcess(cmd, key, sessionID, "480p", "libx264", true)
	require.NoError(t, err)

	<-proc.Done
	assert.NoError(t, proc.Err)
}

// ===========================================================================
// BuildVideoOnlyCommand — additional edge cases
// ===========================================================================

func TestBuildVideoOnlyCommand_TranscodeNoBitrateNoHeight(t *testing.T) {
	pd := ProfileDecision{
		Name:         "custom",
		VideoCodec:   "libx264",
		VideoBitrate: 0,
		Height:       0,
	}

	cmd := BuildVideoOnlyCommand("ffmpeg", "/media/movie.mkv", "/tmp/segments/custom", pd, 4, 0)
	args := CommandArgs(cmd)

	assert.Contains(t, args, "-c:v")
	assert.Contains(t, args, "libx264")
	// No -vf scale since Height is 0
	assert.NotContains(t, args, "-vf")
	// No -maxrate since VideoBitrate is 0
	assert.NotContains(t, args, "-maxrate")
}

func TestBuildVideoOnlyCommand_TranscodeWithBitrate(t *testing.T) {
	pd := ProfileDecision{
		Name:         "480p",
		Width:        854,
		Height:       480,
		VideoCodec:   "libx264",
		VideoBitrate: 1400,
	}

	cmd := BuildVideoOnlyCommand("ffmpeg", "/media/movie.mkv", "/tmp/segments/480p", pd, 4, 0)
	args := CommandArgs(cmd)

	assert.Contains(t, args, "1400k")
	assert.Contains(t, args, "2800k") // bufsize = bitrate * 2
	assert.Contains(t, args, "scale=-2:480")
}

func TestBuildAudioRenditionCommand_NoBitrate(t *testing.T) {
	cmd := BuildAudioRenditionCommand("ffmpeg", "/media/movie.mkv", "/tmp/audio/0", 0, "aac", 0, 6, 0)
	args := CommandArgs(cmd)

	assert.Contains(t, args, "-c:a")
	assert.Contains(t, args, "aac")
	// With zero bitrate, should NOT have -b:a
	assert.NotContains(t, args, "-b:a")
}

// ===========================================================================
// Helpers
// ===========================================================================

func newSleepCmd() *exec.Cmd {
	return exec.Command("sleep", "60")
}

func newTrueCmd() *exec.Cmd {
	return exec.Command("true")
}

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
