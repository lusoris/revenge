package playback_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/playback"
	"github.com/lusoris/revenge/internal/playback/hls"
	"github.com/lusoris/revenge/internal/playback/transcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testVideoPath returns the absolute path to the test video file and skips
// the test if it's not present.
func testVideoPath(t *testing.T) string {
	t.Helper()
	// Resolve relative to repo root
	path, err := filepath.Abs(filepath.Join("..", "..", ".workingdir11", "bbb_sunflower_2160p_30fps_normal.mp4"))
	if err != nil {
		t.Skip("cannot resolve test video path")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("test video not found at " + path + " — copy Big Buck Bunny 4K to .workingdir11/")
	}
	return path
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
}

// ---------------------------------------------------------------------------
// Mock movie.Service — returns a fixed file path for GetMovieFiles
// ---------------------------------------------------------------------------

type fakeMovieService struct {
	movie.Service
	files    []movie.MovieFile
	filesErr error
}

func (m *fakeMovieService) GetMovieFiles(_ context.Context, _ uuid.UUID) ([]movie.MovieFile, error) {
	return m.files, m.filesErr
}

// ---------------------------------------------------------------------------
// Mock tvshow.Service — unused but required for constructor
// ---------------------------------------------------------------------------

type fakeTVService struct {
	tvshow.Service
}

// ---------------------------------------------------------------------------
// Helper: create a fully wired Service + SessionManager for integration tests
// ---------------------------------------------------------------------------

func newIntegrationService(t *testing.T, videoPath string) (*playback.Service, *playback.SessionManager, *transcode.PipelineManager) {
	t.Helper()

	segDir := filepath.Join(t.TempDir(), "segments")
	require.NoError(t, os.MkdirAll(segDir, 0o755))

	cfg := &config.Config{
		Playback: config.PlaybackConfig{
			Enabled:               true,
			SegmentDir:            segDir,
			SegmentDuration:       2,
			MaxConcurrentSessions: 10,
			SessionTimeout:        5 * time.Minute,
			FFmpegPath:            "ffmpeg",
			Transcode: config.TranscodeConfig{
				Enabled:  true,
				Profiles: []string{"original", "720p"},
			},
		},
	}

	fileID := uuid.Must(uuid.NewV7())
	movieSvc := &fakeMovieService{
		files: []movie.MovieFile{
			{ID: fileID, FilePath: videoPath},
		},
	}

	logger := testLogger()

	sm, err := playback.NewSessionManager(10, 5*time.Minute, logger)
	require.NoError(t, err)
	t.Cleanup(sm.Close)

	pm, err := transcode.NewPipelineManager(2, logger)
	require.NoError(t, err)
	t.Cleanup(pm.Close)

	prober := movie.NewMediaInfoProber()

	svc, err := playback.NewService(cfg, sm, pm, prober, movieSvc, &fakeTVService{}, logger)
	require.NoError(t, err)
	t.Cleanup(svc.Close)

	return svc, sm, pm
}

// ===========================================================================
// 1. Probe real video file
// ===========================================================================

func TestIntegration_ProbeRealVideo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	videoPath := testVideoPath(t)

	prober := movie.NewMediaInfoProber()
	info, err := prober.Probe(videoPath)
	require.NoError(t, err)
	require.NotNil(t, info)

	// Video assertions
	assert.Equal(t, "h264", info.VideoCodec)
	assert.Equal(t, 3840, info.Width)
	assert.Equal(t, 2160, info.Height)
	assert.InDelta(t, 10.0, info.DurationSeconds, 1.0)
	assert.Greater(t, info.VideoBitrateKbps, int64(5000))

	// Audio — two tracks: MP3 stereo + AC-3 5.1
	require.Len(t, info.AudioStreams, 2)
	assert.Equal(t, "mp3", info.AudioStreams[0].Codec)
	assert.Equal(t, 2, info.AudioStreams[0].Channels)
	assert.Equal(t, "ac3", info.AudioStreams[1].Codec)
	assert.Equal(t, 6, info.AudioStreams[1].Channels)

	// No subtitles in this file
	assert.Empty(t, info.SubtitleStreams)
}

// ===========================================================================
// 2. Transcode decision for H.264 4K with MP3+AC3
// ===========================================================================

func TestIntegration_TranscodeDecision(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	videoPath := testVideoPath(t)

	prober := movie.NewMediaInfoProber()
	info, err := prober.Probe(videoPath)
	require.NoError(t, err)

	profiles := transcode.GetEnabledProfiles([]string{"original", "1080p", "720p"})
	decision := transcode.AnalyzeMedia(info, profiles)

	// H.264 is HLS-compatible, MP3 is HLS-compatible
	assert.True(t, decision.CanRemux)
	assert.Equal(t, "h264", decision.SourceVideoCodec)
	assert.Equal(t, "mp3", decision.SourceAudioCodec)

	require.Len(t, decision.Profiles, 3)

	// Original: copy/copy (H.264+MP3 are HLS-compatible)
	orig := decision.Profiles[0]
	assert.Equal(t, "original", orig.Name)
	assert.Equal(t, "copy", orig.VideoCodec)
	assert.Equal(t, "copy", orig.AudioCodec)
	assert.False(t, orig.NeedsTranscode)

	// 1080p: must downscale from 4K → transcode
	p1080 := decision.Profiles[1]
	assert.Equal(t, "1080p", p1080.Name)
	assert.Equal(t, "libx264", p1080.VideoCodec)
	assert.True(t, p1080.NeedsTranscode)
	assert.Equal(t, 1080, p1080.Height)

	// 720p: must downscale from 4K → transcode
	p720 := decision.Profiles[2]
	assert.Equal(t, "720p", p720.Name)
	assert.Equal(t, "libx264", p720.VideoCodec)
	assert.True(t, p720.NeedsTranscode)
	assert.Equal(t, 720, p720.Height)
}

// ===========================================================================
// 3. Full session lifecycle: start → generate segments → serve HLS → stop
// ===========================================================================

func TestIntegration_SessionLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	videoPath := testVideoPath(t)
	svc, sm, _ := newIntegrationService(t, videoPath)

	userID := uuid.Must(uuid.NewV7())
	mediaID := uuid.Must(uuid.NewV7())

	// --- Start session ---
	req := &playback.StartPlaybackRequest{
		MediaType:     playback.MediaTypeMovie,
		MediaID:       mediaID,
		AudioTrack:    0,
		StartPosition: 0,
	}
	sess, err := svc.StartSession(context.Background(), userID, req)
	require.NoError(t, err)
	require.NotNil(t, sess)

	assert.Equal(t, userID, sess.UserID)
	assert.Equal(t, playback.MediaTypeMovie, sess.MediaType)
	assert.InDelta(t, 10.0, sess.DurationSeconds, 1.0)
	assert.NotEmpty(t, sess.ActiveProfiles)
	assert.Len(t, sess.AudioTracks, 2)
	assert.Empty(t, sess.SubtitleTracks)

	// Session should be retrievable
	got, ok := svc.GetSession(sess.ID)
	require.True(t, ok)
	assert.Equal(t, sess.ID, got.ID)

	// Active count
	assert.Equal(t, 1, sm.ActiveCount())

	// --- Wait for FFmpeg to produce at least one segment ---
	// "original" profile is copy so it should produce segments very quickly.
	origSegDir := filepath.Join(sess.SegmentDir, "original")
	require.Eventually(t, func() bool {
		entries, err := os.ReadDir(origSegDir)
		if err != nil {
			return false
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".ts") {
				return true
			}
		}
		return false
	}, 15*time.Second, 200*time.Millisecond, "FFmpeg should produce at least one segment")

	// --- Serve master playlist via HTTP handler ---
	handler, err := hls.NewStreamHandler(sm, testLogger())
	require.NoError(t, err)
	t.Cleanup(handler.Close)

	masterReq := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/master.m3u8", nil)
	masterRec := httptest.NewRecorder()
	handler.ServeHTTP(masterRec, masterReq)

	require.Equal(t, http.StatusOK, masterRec.Code)
	masterBody := masterRec.Body.String()
	assert.Contains(t, masterBody, "#EXTM3U")
	assert.Contains(t, masterBody, "original/index.m3u8")
	assert.Contains(t, masterBody, "720p/index.m3u8")
	assert.Contains(t, masterBody, `TYPE=AUDIO`)
	assert.Contains(t, masterBody, `URI="audio/0/index.m3u8"`)
	assert.Contains(t, masterBody, `URI="audio/1/index.m3u8"`)

	t.Logf("Master playlist:\n%s", masterBody)

	// --- Serve media playlist ---
	mediaReq := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/original/index.m3u8", nil)
	mediaRec := httptest.NewRecorder()
	handler.ServeHTTP(mediaRec, mediaReq)

	require.Equal(t, http.StatusOK, mediaRec.Code)
	mediaBody := mediaRec.Body.String()
	assert.Contains(t, mediaBody, "#EXTM3U")
	assert.Contains(t, mediaBody, "#EXTINF:")
	assert.Contains(t, mediaBody, "seg-00000.ts")

	t.Logf("Media playlist:\n%s", mediaBody)

	// --- Serve a real segment ---
	segReq := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/original/seg-00000.ts", nil)
	segRec := httptest.NewRecorder()
	handler.ServeHTTP(segRec, segReq)

	require.Equal(t, http.StatusOK, segRec.Code)
	assert.Equal(t, "video/mp2t", segRec.Header().Get("Content-Type"))
	assert.Greater(t, segRec.Body.Len(), 1000, "segment should contain real video data")
	assert.Contains(t, segRec.Header().Get("Cache-Control"), "immutable")

	// --- Serve audio rendition ---
	// Wait for audio segments to appear
	audioDir := filepath.Join(sess.SegmentDir, "audio", "0")
	require.Eventually(t, func() bool {
		entries, err := os.ReadDir(audioDir)
		if err != nil {
			return false
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".ts") {
				return true
			}
		}
		return false
	}, 15*time.Second, 200*time.Millisecond, "audio rendition should produce segments")

	audioSegReq := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/audio/0/seg-00000.ts", nil)
	audioSegRec := httptest.NewRecorder()
	handler.ServeHTTP(audioSegRec, audioSegReq)

	assert.Equal(t, http.StatusOK, audioSegRec.Code)
	assert.Equal(t, "video/mp2t", audioSegRec.Header().Get("Content-Type"))
	assert.Greater(t, audioSegRec.Body.Len(), 100, "audio segment should contain real data")

	// --- Stop session ---
	err = svc.StopSession(sess.ID)
	require.NoError(t, err)

	assert.Equal(t, 0, sm.ActiveCount())
	_, ok = svc.GetSession(sess.ID)
	assert.False(t, ok)
}

// ===========================================================================
// 4. Seek start — FFmpeg should begin from the specified position
// ===========================================================================

func TestIntegration_SeekStart(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	videoPath := testVideoPath(t)
	svc, sm, _ := newIntegrationService(t, videoPath)

	userID := uuid.Must(uuid.NewV7())
	req := &playback.StartPlaybackRequest{
		MediaType:     playback.MediaTypeMovie,
		MediaID:       uuid.Must(uuid.NewV7()),
		AudioTrack:    0,
		StartPosition: 300, // seek to 5 min in
	}
	sess, err := svc.StartSession(context.Background(), userID, req)
	require.NoError(t, err)
	assert.Equal(t, 300, sess.StartPosition)

	// Segments should still be produced from the seek point
	origSegDir := filepath.Join(sess.SegmentDir, "original")
	require.Eventually(t, func() bool {
		entries, err := os.ReadDir(origSegDir)
		if err != nil {
			return false
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".ts") {
				return true
			}
		}
		return false
	}, 15*time.Second, 200*time.Millisecond, "FFmpeg should produce segments from seek position")

	_ = svc.StopSession(sess.ID)
	assert.Equal(t, 0, sm.ActiveCount())
}

// ===========================================================================
// 5. Max concurrent sessions enforced
// ===========================================================================

func TestIntegration_MaxConcurrentSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	videoPath := testVideoPath(t)

	segDir := filepath.Join(t.TempDir(), "segments")
	require.NoError(t, os.MkdirAll(segDir, 0o755))

	cfg := &config.Config{
		Playback: config.PlaybackConfig{
			Enabled:               true,
			SegmentDir:            segDir,
			SegmentDuration:       2,
			MaxConcurrentSessions: 1, // only allow 1
			SessionTimeout:        5 * time.Minute,
			FFmpegPath:            "ffmpeg",
			Transcode: config.TranscodeConfig{
				Enabled:  true,
				Profiles: []string{"original"},
			},
		},
	}

	fileID := uuid.Must(uuid.NewV7())
	movieSvc := &fakeMovieService{
		files: []movie.MovieFile{
			{ID: fileID, FilePath: videoPath},
		},
	}
	logger := testLogger()

	sm, err := playback.NewSessionManager(1, 5*time.Minute, logger)
	require.NoError(t, err)
	t.Cleanup(sm.Close)

	pm, err := transcode.NewPipelineManager(2, logger)
	require.NoError(t, err)
	t.Cleanup(pm.Close)

	prober := movie.NewMediaInfoProber()
	svc, err := playback.NewService(cfg, sm, pm, prober, movieSvc, nil, logger)
	require.NoError(t, err)
	t.Cleanup(svc.Close)

	userID := uuid.Must(uuid.NewV7())

	// First session succeeds
	sess1, err := svc.StartSession(context.Background(), userID, &playback.StartPlaybackRequest{
		MediaType: playback.MediaTypeMovie,
		MediaID:   uuid.Must(uuid.NewV7()),
	})
	require.NoError(t, err)
	require.NotNil(t, sess1)

	// Second session fails — max reached
	sess2, err := svc.StartSession(context.Background(), userID, &playback.StartPlaybackRequest{
		MediaType: playback.MediaTypeMovie,
		MediaID:   uuid.Must(uuid.NewV7()),
	})
	assert.Error(t, err)
	assert.Nil(t, sess2)
	assert.Contains(t, err.Error(), "maximum concurrent sessions")

	// Stop first → second should succeed
	require.NoError(t, svc.StopSession(sess1.ID))

	sess3, err := svc.StartSession(context.Background(), userID, &playback.StartPlaybackRequest{
		MediaType: playback.MediaTypeMovie,
		MediaID:   uuid.Must(uuid.NewV7()),
	})
	require.NoError(t, err)
	require.NotNil(t, sess3)

	_ = svc.StopSession(sess3.ID)
}

// ===========================================================================
// 6. SessionToResponse round-trip
// ===========================================================================

func TestIntegration_SessionToResponse(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	videoPath := testVideoPath(t)
	svc, _, _ := newIntegrationService(t, videoPath)

	userID := uuid.Must(uuid.NewV7())
	sess, err := svc.StartSession(context.Background(), userID, &playback.StartPlaybackRequest{
		MediaType: playback.MediaTypeMovie,
		MediaID:   uuid.Must(uuid.NewV7()),
	})
	require.NoError(t, err)
	defer svc.StopSession(sess.ID)

	resp := playback.SessionToResponse(sess)
	assert.Equal(t, sess.ID, resp.SessionID)
	assert.Contains(t, resp.MasterPlaylistURL, sess.ID.String())
	assert.Contains(t, resp.MasterPlaylistURL, "master.m3u8")
	assert.InDelta(t, 10.0, resp.DurationSeconds, 1.0)
	assert.NotEmpty(t, resp.Profiles)
	assert.Len(t, resp.AudioTracks, 2)
	assert.Empty(t, resp.SubtitleTracks)
	assert.False(t, resp.CreatedAt.IsZero())
	assert.False(t, resp.ExpiresAt.IsZero())

	// Original profile should be marked IsOriginal
	foundOriginal := false
	for _, p := range resp.Profiles {
		if p.Name == "original" {
			assert.True(t, p.IsOriginal)
			assert.Equal(t, 3840, p.Width)
			assert.Equal(t, 2160, p.Height)
			foundOriginal = true
		}
		if p.Name == "720p" {
			assert.False(t, p.IsOriginal)
			assert.Equal(t, 720, p.Height)
		}
	}
	assert.True(t, foundOriginal, "original profile should be present")
}

// ===========================================================================
// 7. Path traversal security on HLS handler
// ===========================================================================

func TestIntegration_PathTraversalSecurity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}

	logger := testLogger()
	sm, err := playback.NewSessionManager(10, 30*time.Minute, logger)
	require.NoError(t, err)
	t.Cleanup(sm.Close)

	handler, err := hls.NewStreamHandler(sm, logger)
	require.NoError(t, err)
	t.Cleanup(handler.Close)

	// Create a session with a temp segment dir
	segDir := t.TempDir()
	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		FileID:     uuid.Must(uuid.NewV7()),
		FilePath:   "/media/movie.mkv",
		SegmentDir: segDir,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{
				{Name: "original", Width: 1920, Height: 1080, VideoCodec: "copy", AudioCodec: "copy"},
			},
		},
	}
	require.NoError(t, sm.Create(sess))
	base := "/api/v1/playback/stream/" + sess.ID.String()

	// Create a sensitive file outside the segment dir
	sensitiveDir := t.TempDir()
	sensitiveFile := filepath.Join(sensitiveDir, "passwd")
	require.NoError(t, os.WriteFile(sensitiveFile, []byte("root:x:0:0"), 0o644))

	traversalPaths := []struct {
		name string
		path string
	}{
		{"segment traversal up", base + "/../../../etc/passwd/seg-00000.ts"},
		{"segment dot-dot", base + "/..%2f..%2f..%2fetc%2fpasswd/seg-00000.ts"},
		{"profile dot-dot", base + "/../../etc/passwd/index.m3u8"},
		{"media playlist traversal", base + "/../../../etc/passwd/index.m3u8"},
		{"audio traversal", base + "/audio/../../../etc/passwd/index.m3u8"},
		{"subtitle traversal", base + "/subs/../../../etc/passwd.vtt"},
	}

	for _, tt := range traversalPaths {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))
			// Should either 400 or 404 — never 200
			assert.NotEqual(t, http.StatusOK, rec.Code,
				"path traversal attack should not succeed: %s", tt.path)
			assert.NotContains(t, rec.Body.String(), "root:x:0:0",
				"should not leak sensitive file contents")
		})
	}
}

// ===========================================================================
// 8. Multiple audio renditions produced for both tracks
// ===========================================================================

func TestIntegration_MultipleAudioRenditions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	videoPath := testVideoPath(t)
	svc, sm, _ := newIntegrationService(t, videoPath)

	sess, err := svc.StartSession(context.Background(), uuid.Must(uuid.NewV7()), &playback.StartPlaybackRequest{
		MediaType: playback.MediaTypeMovie,
		MediaID:   uuid.Must(uuid.NewV7()),
	})
	require.NoError(t, err)
	defer svc.StopSession(sess.ID)

	// BBB has 2 audio tracks — both should get their own rendition
	assert.Len(t, sess.AudioTracks, 2)

	// Track 0: MP3 stereo
	assert.Equal(t, "mp3", sess.AudioTracks[0].Codec)
	assert.Equal(t, 2, sess.AudioTracks[0].Channels)

	// Track 1: AC-3 5.1
	assert.Equal(t, "ac3", sess.AudioTracks[1].Codec)
	assert.Equal(t, 6, sess.AudioTracks[1].Channels)

	// Wait for both audio renditions to produce segments
	for _, trackIdx := range []int{0, 1} {
		audioDir := filepath.Join(sess.SegmentDir, "audio", fmt.Sprintf("%d", trackIdx))
		require.Eventually(t, func() bool {
			entries, err := os.ReadDir(audioDir)
			if err != nil {
				return false
			}
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".ts") {
					return true
				}
			}
			return false
		}, 15*time.Second, 200*time.Millisecond,
			"audio rendition track %d should produce segments", trackIdx)
	}

	// Serve audio via HTTP handler
	handler, err := hls.NewStreamHandler(sm, testLogger())
	require.NoError(t, err)
	t.Cleanup(handler.Close)

	// Master playlist should list both audio tracks
	masterRec := httptest.NewRecorder()
	handler.ServeHTTP(masterRec, httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/master.m3u8", nil))
	require.Equal(t, http.StatusOK, masterRec.Code)
	assert.Contains(t, masterRec.Body.String(), `URI="audio/0/index.m3u8"`)
	assert.Contains(t, masterRec.Body.String(), `URI="audio/1/index.m3u8"`)

	// Both audio playlists should be servable
	for _, trackIdx := range []int{0, 1} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/v1/playback/stream/%s/audio/%d/index.m3u8", sess.ID, trackIdx), nil))
		assert.Equal(t, http.StatusOK, rec.Code, "audio track %d playlist", trackIdx)
		assert.Contains(t, rec.Body.String(), "#EXTM3U", "audio track %d", trackIdx)
	}
}

// ===========================================================================
// 9. Pipeline processes are cleaned up on session stop
// ===========================================================================

func TestIntegration_PipelineCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	videoPath := testVideoPath(t)
	svc, _, pm := newIntegrationService(t, videoPath)

	sess, err := svc.StartSession(context.Background(), uuid.Must(uuid.NewV7()), &playback.StartPlaybackRequest{
		MediaType: playback.MediaTypeMovie,
		MediaID:   uuid.Must(uuid.NewV7()),
	})
	require.NoError(t, err)

	// Pipeline should have running processes
	time.Sleep(500 * time.Millisecond)
	_, hasOriginal := pm.GetProcess(sess.ID, "original")
	_, hasAudio0 := pm.GetProcess(sess.ID, "audio/0")
	assert.True(t, hasOriginal || hasAudio0, "at least one process should be running")

	// Stop session — should clean up all processes
	err = svc.StopSession(sess.ID)
	require.NoError(t, err)

	// Give cleanup a moment
	time.Sleep(200 * time.Millisecond)

	_, hasOriginal = pm.GetProcess(sess.ID, "original")
	_, has720p := pm.GetProcess(sess.ID, "720p")
	_, hasAudio0 = pm.GetProcess(sess.ID, "audio/0")
	_, hasAudio1 := pm.GetProcess(sess.ID, "audio/1")
	assert.False(t, hasOriginal, "original process should be cleaned up")
	assert.False(t, has720p, "720p process should be cleaned up")
	assert.False(t, hasAudio0, "audio/0 process should be cleaned up")
	assert.False(t, hasAudio1, "audio/1 process should be cleaned up")
}

// ===========================================================================
// 10. Session Touch keeps session alive
// ===========================================================================

func TestIntegration_SessionTouch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}

	logger := testLogger()
	sm, err := playback.NewSessionManager(10, 5*time.Minute, logger)
	require.NoError(t, err)
	t.Cleanup(sm.Close)

	sess := &playback.Session{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    uuid.Must(uuid.NewV7()),
		MediaType: playback.MediaTypeMovie,
		MediaID:   uuid.Must(uuid.NewV7()),
	}
	require.NoError(t, sm.Create(sess))

	original, ok := sm.Get(sess.ID)
	require.True(t, ok)
	origExpiry := original.ExpiresAt
	origAccess := original.LastAccessedAt

	time.Sleep(100 * time.Millisecond)

	ok = sm.Touch(sess.ID)
	assert.True(t, ok)

	// Touch mutates the session pointer in-place before re-setting in cache,
	// so the same pointer we already have reflects the updated timestamps.
	assert.True(t, original.ExpiresAt.After(origExpiry), "touch should extend expiry")
	assert.True(t, original.LastAccessedAt.After(origAccess), "touch should update LastAccessedAt")
}
