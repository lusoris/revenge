package playback

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/playback/transcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testPipelineManager creates a PipelineManager for tests.
func testPipelineManager(t *testing.T) *transcode.PipelineManager {
	t.Helper()
	pm, err := transcode.NewPipelineManager(4, testLogger())
	require.NoError(t, err)
	t.Cleanup(pm.Close)
	return pm
}

// ---------------------------------------------------------------------------
// Minimal mock: movie.Service — only implements methods used by playback
// ---------------------------------------------------------------------------

type mockMovieService struct {
	movie.Service // embed interface; panics on unimplemented methods
	files         []movie.MovieFile
	filesErr      error
}

func (m *mockMovieService) GetMovieFiles(_ context.Context, _ uuid.UUID) ([]movie.MovieFile, error) {
	return m.files, m.filesErr
}

// ---------------------------------------------------------------------------
// Minimal mock: tvshow.Service — only implements methods used by playback
// ---------------------------------------------------------------------------

type mockTVService struct {
	tvshow.Service // embed interface
	files          []tvshow.EpisodeFile
	filesErr       error
	file           *tvshow.EpisodeFile
	fileErr        error
}

func (m *mockTVService) ListEpisodeFiles(_ context.Context, _ uuid.UUID) ([]tvshow.EpisodeFile, error) {
	return m.files, m.filesErr
}

func (m *mockTVService) GetEpisodeFile(_ context.Context, _ uuid.UUID) (*tvshow.EpisodeFile, error) {
	return m.file, m.fileErr
}

// ---------------------------------------------------------------------------
// Minimal mock: movie.Prober
// ---------------------------------------------------------------------------

type mockProber struct {
	info *movie.MediaInfo
	err  error
}

func (m *mockProber) Probe(_ string) (*movie.MediaInfo, error) {
	return m.info, m.err
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// testConfig returns a minimal Config suitable for unit tests.
func testConfig() *config.Config {
	return &config.Config{
		Playback: config.PlaybackConfig{
			Enabled:               true,
			SegmentDir:            "/tmp/revenge-test-segments",
			SegmentDuration:       4,
			MaxConcurrentSessions: 10,
			SessionTimeout:        30 * time.Minute,
			FFmpegPath:            "ffmpeg",
			Transcode: config.TranscodeConfig{
				Enabled:  true,
				Profiles: []string{"original", "1080p", "720p"},
			},
		},
	}
}

func newTestService(t *testing.T, cfg *config.Config, movieSvc movie.Service, tvSvc tvshow.Service, prober movie.Prober) (*Service, *SessionManager) {
	t.Helper()
	sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
	require.NoError(t, err)
	t.Cleanup(sm.Close)

	pm := testPipelineManager(t)
	svc, err := NewService(cfg, sm, pm, prober, movieSvc, tvSvc, testLogger())
	require.NoError(t, err)
	t.Cleanup(svc.Close)

	return svc, sm
}

// ---------------------------------------------------------------------------
// NewService tests
// ---------------------------------------------------------------------------

func TestNewService(t *testing.T) {
	t.Run("creates service with valid config", func(t *testing.T) {
		cfg := testConfig()
		sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
		require.NoError(t, err)
		defer sm.Close()

		pm := testPipelineManager(t)
		svc, err := NewService(cfg, sm, pm, nil, nil, nil, testLogger())
		require.NoError(t, err)
		require.NotNil(t, svc)
		defer svc.Close()

		assert.Equal(t, cfg, svc.cfg)
		assert.Equal(t, sm, svc.sessions)
		assert.NotNil(t, svc.probeCache)
		assert.NotNil(t, svc.logger)
	})

	t.Run("loads enabled profiles from config", func(t *testing.T) {
		cfg := testConfig()
		svc, _ := newTestService(t, cfg, nil, nil, nil)
		// Config has "original", "1080p", "720p"
		assert.Len(t, svc.profiles, 3)
	})

	t.Run("empty profiles list produces no profiles", func(t *testing.T) {
		cfg := testConfig()
		cfg.Playback.Transcode.Profiles = []string{}
		svc, _ := newTestService(t, cfg, nil, nil, nil)
		assert.Empty(t, svc.profiles)
	})

	t.Run("unknown profile names are silently skipped", func(t *testing.T) {
		cfg := testConfig()
		cfg.Playback.Transcode.Profiles = []string{"original", "nonexistent", "720p"}
		svc, _ := newTestService(t, cfg, nil, nil, nil)
		assert.Len(t, svc.profiles, 2, "only valid profiles should be loaded")
	})

	t.Run("nil optional services are accepted", func(t *testing.T) {
		cfg := testConfig()
		svc, _ := newTestService(t, cfg, nil, nil, nil)
		assert.NotNil(t, svc)
	})
}

// ---------------------------------------------------------------------------
// GetSession tests
// ---------------------------------------------------------------------------

func TestGetSession(t *testing.T) {
	cfg := testConfig()
	svc, sm := newTestService(t, cfg, nil, nil, nil)

	t.Run("returns false for nonexistent session", func(t *testing.T) {
		_, ok := svc.GetSession(uuid.New())
		assert.False(t, ok)
	})

	t.Run("returns session after creation via session manager", func(t *testing.T) {
		id := uuid.New()
		sess := &Session{
			ID:        id,
			UserID:    uuid.New(),
			MediaType: MediaTypeMovie,
			MediaID:   uuid.New(),
			FilePath:  "/media/test.mkv",
		}
		err := sm.Create(sess)
		require.NoError(t, err)

		got, ok := svc.GetSession(id)
		require.True(t, ok)
		assert.Equal(t, id, got.ID)
	})
}

// ---------------------------------------------------------------------------
// StopSession tests
// ---------------------------------------------------------------------------

func TestStopSession(t *testing.T) {
	t.Run("not found returns error", func(t *testing.T) {
		cfg := testConfig()
		svc, _ := newTestService(t, cfg, nil, nil, nil)

		err := svc.StopSession(uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("stops existing session", func(t *testing.T) {
		cfg := testConfig()
		svc, sm := newTestService(t, cfg, nil, nil, nil)

		id := uuid.New()
		err := sm.Create(&Session{
			ID:         id,
			UserID:     uuid.New(),
			MediaType:  MediaTypeMovie,
			MediaID:    uuid.New(),
			SegmentDir: t.TempDir(), // real temp dir so cleanup works
		})
		require.NoError(t, err)

		err = svc.StopSession(id)
		assert.NoError(t, err)

		// Session should be gone
		_, ok := svc.GetSession(id)
		assert.False(t, ok)
	})

	t.Run("active count decremented after stop", func(t *testing.T) {
		cfg := testConfig()
		svc, sm := newTestService(t, cfg, nil, nil, nil)

		id := uuid.New()
		err := sm.Create(&Session{
			ID:         id,
			UserID:     uuid.New(),
			MediaType:  MediaTypeMovie,
			MediaID:    uuid.New(),
			SegmentDir: t.TempDir(),
		})
		require.NoError(t, err)
		assert.Equal(t, 1, sm.ActiveCount())

		err = svc.StopSession(id)
		require.NoError(t, err)
		assert.Equal(t, 0, sm.ActiveCount())
	})
}

// ---------------------------------------------------------------------------
// resolveFilePath tests
// ---------------------------------------------------------------------------

func TestResolveFilePath(t *testing.T) {
	ctx := context.Background()

	t.Run("movie: returns first file when no FileID specified", func(t *testing.T) {
		fileID := uuid.New()
		movieSvc := &mockMovieService{
			files: []movie.MovieFile{
				{ID: fileID, FilePath: "/media/movies/test.mkv"},
			},
		}
		svc, _ := newTestService(t, testConfig(), movieSvc, nil, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeMovie,
			MediaID:   uuid.New(),
		}
		path, id, err := svc.resolveFilePath(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, "/media/movies/test.mkv", path)
		assert.Equal(t, fileID, id)
	})

	t.Run("movie: returns specific file when FileID specified", func(t *testing.T) {
		file1ID := uuid.New()
		file2ID := uuid.New()
		movieSvc := &mockMovieService{
			files: []movie.MovieFile{
				{ID: file1ID, FilePath: "/media/movies/test-720p.mkv"},
				{ID: file2ID, FilePath: "/media/movies/test-1080p.mkv"},
			},
		}
		svc, _ := newTestService(t, testConfig(), movieSvc, nil, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeMovie,
			MediaID:   uuid.New(),
			FileID:    &file2ID,
		}
		path, id, err := svc.resolveFilePath(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, "/media/movies/test-1080p.mkv", path)
		assert.Equal(t, file2ID, id)
	})

	t.Run("movie: error when file ID not found in files list", func(t *testing.T) {
		movieSvc := &mockMovieService{
			files: []movie.MovieFile{
				{ID: uuid.New(), FilePath: "/media/movies/test.mkv"},
			},
		}
		svc, _ := newTestService(t, testConfig(), movieSvc, nil, nil)

		wrongID := uuid.New()
		req := &StartPlaybackRequest{
			MediaType: MediaTypeMovie,
			MediaID:   uuid.New(),
			FileID:    &wrongID,
		}
		_, _, err := svc.resolveFilePath(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("movie: error when no files available", func(t *testing.T) {
		movieSvc := &mockMovieService{
			files: []movie.MovieFile{},
		}
		svc, _ := newTestService(t, testConfig(), movieSvc, nil, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeMovie,
			MediaID:   uuid.New(),
		}
		_, _, err := svc.resolveFilePath(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no files available")
	})

	t.Run("movie: error from movie service is propagated", func(t *testing.T) {
		movieSvc := &mockMovieService{
			filesErr: fmt.Errorf("db connection failed"),
		}
		svc, _ := newTestService(t, testConfig(), movieSvc, nil, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeMovie,
			MediaID:   uuid.New(),
		}
		_, _, err := svc.resolveFilePath(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "movie files not found")
	})

	t.Run("episode: returns first file when no FileID specified", func(t *testing.T) {
		fileID := uuid.New()
		tvSvc := &mockTVService{
			files: []tvshow.EpisodeFile{
				{ID: fileID, FilePath: "/media/tv/episode.mkv"},
			},
		}
		svc, _ := newTestService(t, testConfig(), nil, tvSvc, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeEpisode,
			MediaID:   uuid.New(),
		}
		path, id, err := svc.resolveFilePath(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, "/media/tv/episode.mkv", path)
		assert.Equal(t, fileID, id)
	})

	t.Run("episode: returns specific file when FileID specified", func(t *testing.T) {
		fileID := uuid.New()
		tvSvc := &mockTVService{
			file: &tvshow.EpisodeFile{ID: fileID, FilePath: "/media/tv/ep-1080p.mkv"},
		}
		svc, _ := newTestService(t, testConfig(), nil, tvSvc, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeEpisode,
			MediaID:   uuid.New(),
			FileID:    &fileID,
		}
		path, id, err := svc.resolveFilePath(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, "/media/tv/ep-1080p.mkv", path)
		assert.Equal(t, fileID, id)
	})

	t.Run("episode: error when GetEpisodeFile fails", func(t *testing.T) {
		fileID := uuid.New()
		tvSvc := &mockTVService{
			fileErr: fmt.Errorf("not found"),
		}
		svc, _ := newTestService(t, testConfig(), nil, tvSvc, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeEpisode,
			MediaID:   uuid.New(),
			FileID:    &fileID,
		}
		_, _, err := svc.resolveFilePath(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "episode file not found")
	})

	t.Run("episode: error when no files available", func(t *testing.T) {
		tvSvc := &mockTVService{
			files: []tvshow.EpisodeFile{},
		}
		svc, _ := newTestService(t, testConfig(), nil, tvSvc, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeEpisode,
			MediaID:   uuid.New(),
		}
		_, _, err := svc.resolveFilePath(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no files available")
	})

	t.Run("episode: error when ListEpisodeFiles fails", func(t *testing.T) {
		tvSvc := &mockTVService{
			filesErr: fmt.Errorf("db timeout"),
		}
		svc, _ := newTestService(t, testConfig(), nil, tvSvc, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeEpisode,
			MediaID:   uuid.New(),
		}
		_, _, err := svc.resolveFilePath(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "episode files not found")
	})

	t.Run("episode: error when tv service is nil", func(t *testing.T) {
		svc, _ := newTestService(t, testConfig(), nil, nil, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaTypeEpisode,
			MediaID:   uuid.New(),
		}
		_, _, err := svc.resolveFilePath(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TV show service not available")
	})

	t.Run("unsupported media type returns error", func(t *testing.T) {
		svc, _ := newTestService(t, testConfig(), nil, nil, nil)

		req := &StartPlaybackRequest{
			MediaType: MediaType("audiobook"),
			MediaID:   uuid.New(),
		}
		_, _, err := svc.resolveFilePath(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported media type")
	})
}

// ---------------------------------------------------------------------------
// probeFile tests
// ---------------------------------------------------------------------------

func TestProbeFile(t *testing.T) {
	t.Run("returns probed info on success", func(t *testing.T) {
		expected := &movie.MediaInfo{
			FilePath:        "/media/test.mkv",
			VideoCodec:      "h264",
			Width:           1920,
			Height:          1080,
			DurationSeconds: 3600,
		}
		prober := &mockProber{info: expected}
		svc, _ := newTestService(t, testConfig(), nil, nil, prober)

		fileID := uuid.New()
		info, err := svc.probeFile(fileID, "/media/test.mkv")
		require.NoError(t, err)
		assert.Equal(t, expected, info)
	})

	t.Run("returns error on probe failure", func(t *testing.T) {
		prober := &mockProber{err: fmt.Errorf("ffprobe not found")}
		svc, _ := newTestService(t, testConfig(), nil, nil, prober)

		_, err := svc.probeFile(uuid.New(), "/media/test.mkv")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ffprobe not found")
	})

	t.Run("caches probe result for same file ID", func(t *testing.T) {
		callCount := 0
		expected := &movie.MediaInfo{VideoCodec: "hevc", Width: 3840, Height: 2160}
		prober := &mockProber{info: expected}

		// Wrap prober to count calls
		countingProber := &countingProber{inner: prober}
		svc, _ := newTestService(t, testConfig(), nil, nil, countingProber)

		fileID := uuid.New()

		// First call should invoke prober
		info1, err := svc.probeFile(fileID, "/media/test.mkv")
		require.NoError(t, err)
		assert.Equal(t, expected, info1)
		callCount = countingProber.count
		assert.Equal(t, 1, callCount)

		// Second call with same fileID should return cached
		info2, err := svc.probeFile(fileID, "/media/test.mkv")
		require.NoError(t, err)
		assert.Equal(t, expected, info2)
		assert.Equal(t, 1, countingProber.count, "prober should not be called again for cached file")
	})

	t.Run("different file IDs are probed separately", func(t *testing.T) {
		expected := &movie.MediaInfo{VideoCodec: "h264"}
		countingProber := &countingProber{inner: &mockProber{info: expected}}
		svc, _ := newTestService(t, testConfig(), nil, nil, countingProber)

		_, err := svc.probeFile(uuid.New(), "/media/a.mkv")
		require.NoError(t, err)
		_, err = svc.probeFile(uuid.New(), "/media/b.mkv")
		require.NoError(t, err)

		assert.Equal(t, 2, countingProber.count, "each unique file ID should be probed")
	})
}

// countingProber wraps a Prober and counts calls.
type countingProber struct {
	inner movie.Prober
	count int
}

func (c *countingProber) Probe(filePath string) (*movie.MediaInfo, error) {
	c.count++
	return c.inner.Probe(filePath)
}

// ---------------------------------------------------------------------------
// audioRenditionCodec tests
// ---------------------------------------------------------------------------

func TestAudioRenditionCodec(t *testing.T) {
	tests := []struct {
		name        string
		sourceCodec string
		wantCodec   string
		wantBitrate int
	}{
		// HLS-compatible codecs should be copied
		{name: "aac is copied", sourceCodec: "aac", wantCodec: "copy", wantBitrate: 0},
		{name: "mp3 is copied", sourceCodec: "mp3", wantCodec: "copy", wantBitrate: 0},
		{name: "ac3 is copied", sourceCodec: "ac3", wantCodec: "copy", wantBitrate: 0},
		{name: "eac3 is copied", sourceCodec: "eac3", wantCodec: "copy", wantBitrate: 0},

		// Non-HLS codecs should be transcoded to AAC at 256 kbps
		{name: "dts is transcoded", sourceCodec: "dts", wantCodec: "aac", wantBitrate: 256},
		{name: "truehd is transcoded", sourceCodec: "truehd", wantCodec: "aac", wantBitrate: 256},
		{name: "flac is transcoded", sourceCodec: "flac", wantCodec: "aac", wantBitrate: 256},
		{name: "pcm_s16le is transcoded", sourceCodec: "pcm_s16le", wantCodec: "aac", wantBitrate: 256},
		{name: "opus is transcoded", sourceCodec: "opus", wantCodec: "aac", wantBitrate: 256},
		{name: "vorbis is transcoded", sourceCodec: "vorbis", wantCodec: "aac", wantBitrate: 256},
		{name: "empty codec is transcoded", sourceCodec: "", wantCodec: "aac", wantBitrate: 256},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			codec, bitrate := audioRenditionCodec(tc.sourceCodec)
			assert.Equal(t, tc.wantCodec, codec)
			assert.Equal(t, tc.wantBitrate, bitrate)
		})
	}
}

// ---------------------------------------------------------------------------
// profileNames tests
// ---------------------------------------------------------------------------

func TestProfileNames(t *testing.T) {
	t.Run("empty profiles", func(t *testing.T) {
		names := profileNames(nil)
		assert.Empty(t, names)
	})

	t.Run("single profile", func(t *testing.T) {
		profiles := []transcode.ProfileDecision{
			{Name: "original"},
		}
		names := profileNames(profiles)
		assert.Equal(t, []string{"original"}, names)
	})

	t.Run("multiple profiles preserve order", func(t *testing.T) {
		profiles := []transcode.ProfileDecision{
			{Name: "original"},
			{Name: "1080p"},
			{Name: "720p"},
			{Name: "480p"},
		}
		names := profileNames(profiles)
		assert.Equal(t, []string{"original", "1080p", "720p", "480p"}, names)
	})
}

// ---------------------------------------------------------------------------
// SessionToResponse tests
// ---------------------------------------------------------------------------

func TestSessionToResponse(t *testing.T) {
	sessionID := uuid.MustParse("11111111-2222-3333-4444-555555555555")
	now := time.Now()
	expires := now.Add(30 * time.Minute)

	t.Run("basic response fields", func(t *testing.T) {
		sess := &Session{
			ID:              sessionID,
			DurationSeconds: 7200.5,
			TranscodeDecision: transcode.Decision{
				Profiles: []transcode.ProfileDecision{},
			},
			AudioTracks:    []AudioTrackInfo{},
			SubtitleTracks: []SubtitleTrackInfo{},
			CreatedAt:      now,
			ExpiresAt:      expires,
		}

		resp := SessionToResponse(sess)
		assert.Equal(t, sessionID, resp.SessionID)
		assert.Equal(t, "/api/v1/playback/stream/11111111-2222-3333-4444-555555555555/master.m3u8", resp.MasterPlaylistURL)
		assert.Equal(t, 7200.5, resp.DurationSeconds)
		assert.Empty(t, resp.Profiles)
		assert.Empty(t, resp.AudioTracks)
		assert.Empty(t, resp.SubtitleTracks)
		assert.Equal(t, now, resp.CreatedAt)
		assert.Equal(t, expires, resp.ExpiresAt)
	})

	t.Run("profiles with transcode", func(t *testing.T) {
		sess := &Session{
			ID:              sessionID,
			DurationSeconds: 3600,
			TranscodeDecision: transcode.Decision{
				Profiles: []transcode.ProfileDecision{
					{
						Name:         "original",
						Width:        1920,
						Height:       1080,
						VideoBitrate: 0,
						VideoCodec:   "copy",
						AudioCodec:   "copy",
					},
					{
						Name:         "720p",
						Width:        1280,
						Height:       720,
						VideoBitrate: 2800,
						VideoCodec:   "libx264",
						AudioCodec:   "aac",
					},
				},
			},
			AudioTracks:    []AudioTrackInfo{},
			SubtitleTracks: []SubtitleTrackInfo{},
			CreatedAt:      now,
			ExpiresAt:      expires,
		}

		resp := SessionToResponse(sess)
		require.Len(t, resp.Profiles, 2)

		// Original profile: copy/copy should be marked IsOriginal
		assert.Equal(t, "original", resp.Profiles[0].Name)
		assert.Equal(t, 1920, resp.Profiles[0].Width)
		assert.Equal(t, 1080, resp.Profiles[0].Height)
		assert.Equal(t, 0, resp.Profiles[0].Bitrate)
		assert.True(t, resp.Profiles[0].IsOriginal)

		// 720p profile: libx264/aac should NOT be marked IsOriginal
		assert.Equal(t, "720p", resp.Profiles[1].Name)
		assert.Equal(t, 1280, resp.Profiles[1].Width)
		assert.Equal(t, 720, resp.Profiles[1].Height)
		assert.Equal(t, 2800, resp.Profiles[1].Bitrate)
		assert.False(t, resp.Profiles[1].IsOriginal)
	})

	t.Run("IsOriginal requires both video and audio copy", func(t *testing.T) {
		sess := &Session{
			ID: sessionID,
			TranscodeDecision: transcode.Decision{
				Profiles: []transcode.ProfileDecision{
					{
						Name:       "video-copy-audio-transcode",
						VideoCodec: "copy",
						AudioCodec: "aac",
					},
					{
						Name:       "video-transcode-audio-copy",
						VideoCodec: "libx264",
						AudioCodec: "copy",
					},
					{
						Name:       "both-transcode",
						VideoCodec: "libx264",
						AudioCodec: "aac",
					},
				},
			},
			AudioTracks:    []AudioTrackInfo{},
			SubtitleTracks: []SubtitleTrackInfo{},
			CreatedAt:      now,
			ExpiresAt:      expires,
		}

		resp := SessionToResponse(sess)
		require.Len(t, resp.Profiles, 3)

		assert.False(t, resp.Profiles[0].IsOriginal, "copy video + aac audio should not be original")
		assert.False(t, resp.Profiles[1].IsOriginal, "libx264 video + copy audio should not be original")
		assert.False(t, resp.Profiles[2].IsOriginal, "both transcoded should not be original")
	})

	t.Run("audio and subtitle tracks are passed through", func(t *testing.T) {
		audioTracks := []AudioTrackInfo{
			{Index: 1, Language: "eng", Codec: "aac", Channels: 2, IsDefault: true},
			{Index: 2, Language: "ger", Codec: "dts", Channels: 6, IsDefault: false},
		}
		subtitleTracks := []SubtitleTrackInfo{
			{Index: 4, Language: "eng", Codec: "subrip", URL: "/api/v1/playback/stream/" + sessionID.String() + "/subs/4.vtt"},
		}

		sess := &Session{
			ID: sessionID,
			TranscodeDecision: transcode.Decision{
				Profiles: []transcode.ProfileDecision{},
			},
			AudioTracks:    audioTracks,
			SubtitleTracks: subtitleTracks,
			CreatedAt:      now,
			ExpiresAt:      expires,
		}

		resp := SessionToResponse(sess)
		assert.Equal(t, audioTracks, resp.AudioTracks)
		assert.Equal(t, subtitleTracks, resp.SubtitleTracks)
	})

	t.Run("nil profiles in decision creates empty response profiles", func(t *testing.T) {
		sess := &Session{
			ID: sessionID,
			TranscodeDecision: transcode.Decision{
				Profiles: nil,
			},
			AudioTracks:    nil,
			SubtitleTracks: nil,
			CreatedAt:      now,
			ExpiresAt:      expires,
		}

		resp := SessionToResponse(sess)
		assert.Empty(t, resp.Profiles)
	})
}

// ---------------------------------------------------------------------------
// Close tests
// ---------------------------------------------------------------------------

func TestServiceClose(t *testing.T) {
	cfg := testConfig()
	sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	pm := testPipelineManager(t)
	svc, err := NewService(cfg, sm, pm, nil, nil, nil, testLogger())
	require.NoError(t, err)

	// Should not panic
	svc.Close()
}
