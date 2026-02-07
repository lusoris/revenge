package playback

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// StartSession tests — uses mocked prober and service, no real FFmpeg
// ---------------------------------------------------------------------------

func TestStartSession_Movie(t *testing.T) {
	fileID := uuid.New()
	movieSvc := &mockMovieService{
		files: []movie.MovieFile{
			{ID: fileID, FilePath: "/media/movies/test.mkv"},
		},
	}
	prober := &mockProber{
		info: &movie.MediaInfo{
			VideoCodec:       "h264",
			Width:            1920,
			Height:           1080,
			DurationSeconds:  7200,
			VideoBitrateKbps: 5000,
			AudioStreams: []movie.AudioStreamInfo{
				{Index: 0, Codec: "aac", Channels: 2, Language: "eng", IsDefault: true},
			},
		},
	}

	cfg := testConfig()
	cfg.Playback.SegmentDir = t.TempDir()

	svc, _ := newTestService(t, cfg, movieSvc, nil, prober)

	req := &StartPlaybackRequest{
		MediaType:     MediaTypeMovie,
		MediaID:       uuid.New(),
		AudioTrack:    0,
		StartPosition: 0,
	}

	sess, err := svc.StartSession(context.Background(), uuid.New(), req)
	require.NoError(t, err)
	require.NotNil(t, sess)

	assert.Equal(t, MediaTypeMovie, sess.MediaType)
	assert.Equal(t, fileID, sess.FileID)
	assert.Equal(t, "/media/movies/test.mkv", sess.FilePath)
	assert.Equal(t, float64(7200), sess.DurationSeconds)
	assert.NotEmpty(t, sess.ActiveProfiles)
	assert.Len(t, sess.AudioTracks, 1)
	assert.Equal(t, "eng", sess.AudioTracks[0].Language)

	// Session should be retrievable
	got, ok := svc.GetSession(sess.ID)
	assert.True(t, ok)
	assert.Equal(t, sess.ID, got.ID)

	// Clean up
	err = svc.StopSession(sess.ID)
	assert.NoError(t, err)
}

func TestStartSession_WithSubtitles(t *testing.T) {
	fileID := uuid.New()
	movieSvc := &mockMovieService{
		files: []movie.MovieFile{
			{ID: fileID, FilePath: "/media/movies/test.mkv"},
		},
	}
	prober := &mockProber{
		info: &movie.MediaInfo{
			VideoCodec:       "h264",
			Width:            1920,
			Height:           1080,
			DurationSeconds:  3600,
			VideoBitrateKbps: 5000,
			AudioStreams: []movie.AudioStreamInfo{
				{Index: 0, Codec: "aac", Channels: 2, Language: "eng"},
			},
			SubtitleStreams: []movie.SubtitleStreamInfo{
				{Index: 0, Codec: "subrip", Language: "eng", Title: "English"},
				{Index: 1, Codec: "hdmv_pgs_subtitle", Language: "eng", Title: "English PGS"},
			},
		},
	}

	cfg := testConfig()
	cfg.Playback.SegmentDir = t.TempDir()

	svc, _ := newTestService(t, cfg, movieSvc, nil, prober)

	req := &StartPlaybackRequest{
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	}

	sess, err := svc.StartSession(context.Background(), uuid.New(), req)
	require.NoError(t, err)
	require.NotNil(t, sess)

	// Only text subtitles should be included (PGS filtered out)
	assert.Len(t, sess.SubtitleTracks, 1)
	assert.Equal(t, "subrip", sess.SubtitleTracks[0].Codec)
	assert.Contains(t, sess.SubtitleTracks[0].URL, sess.ID.String())

	// Give the goroutine a moment to run (extractSubtitles)
	time.Sleep(10 * time.Millisecond)

	_ = svc.StopSession(sess.ID)
}

func TestStartSession_ResolveFileError(t *testing.T) {
	movieSvc := &mockMovieService{
		files: []movie.MovieFile{},
	}
	cfg := testConfig()
	cfg.Playback.SegmentDir = t.TempDir()

	svc, _ := newTestService(t, cfg, movieSvc, nil, nil)

	req := &StartPlaybackRequest{
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	}

	sess, err := svc.StartSession(context.Background(), uuid.New(), req)
	assert.Error(t, err)
	assert.Nil(t, sess)
	assert.Contains(t, err.Error(), "failed to resolve file")
}

func TestStartSession_ProbeError(t *testing.T) {
	fileID := uuid.New()
	movieSvc := &mockMovieService{
		files: []movie.MovieFile{
			{ID: fileID, FilePath: "/media/movies/test.mkv"},
		},
	}
	prober := &mockProber{
		err: assert.AnError,
	}

	cfg := testConfig()
	cfg.Playback.SegmentDir = t.TempDir()

	svc, _ := newTestService(t, cfg, movieSvc, nil, prober)

	req := &StartPlaybackRequest{
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	}

	sess, err := svc.StartSession(context.Background(), uuid.New(), req)
	assert.Error(t, err)
	assert.Nil(t, sess)
	assert.Contains(t, err.Error(), "failed to probe media")
}

func TestStartSession_MaxSessionsExceeded(t *testing.T) {
	fileID := uuid.New()
	movieSvc := &mockMovieService{
		files: []movie.MovieFile{
			{ID: fileID, FilePath: "/media/movies/test.mkv"},
		},
	}
	prober := &mockProber{
		info: &movie.MediaInfo{
			VideoCodec:      "h264",
			Width:           1920,
			Height:          1080,
			DurationSeconds: 3600,
			AudioStreams: []movie.AudioStreamInfo{
				{Index: 0, Codec: "aac", Channels: 2},
			},
		},
	}

	cfg := testConfig()
	cfg.Playback.SegmentDir = t.TempDir()
	cfg.Playback.MaxConcurrentSessions = 1

	// Create session manager with max 1 session
	sm, err := NewSessionManager(1, 30*time.Minute, testLogger())
	require.NoError(t, err)
	t.Cleanup(sm.Close)

	pm := testPipelineManager(t)
	svc, err := NewService(cfg, sm, pm, prober, movieSvc, nil, testLogger())
	require.NoError(t, err)
	t.Cleanup(svc.Close)

	req := &StartPlaybackRequest{
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	}

	// First session should succeed
	sess1, err := svc.StartSession(context.Background(), uuid.New(), req)
	require.NoError(t, err)
	require.NotNil(t, sess1)

	// Second session should fail (max exceeded)
	sess2, err := svc.StartSession(context.Background(), uuid.New(), req)
	assert.Error(t, err)
	assert.Nil(t, sess2)
	assert.Contains(t, err.Error(), "failed to create session")

	_ = svc.StopSession(sess1.ID)
}

func TestStartSession_Episode(t *testing.T) {
	fileID := uuid.New()
	tvSvc := &mockTVService{
		files: []tvshow.EpisodeFile{
			{ID: fileID, FilePath: "/media/tv/episode.mkv"},
		},
	}
	prober := &mockProber{
		info: &movie.MediaInfo{
			VideoCodec:      "h264",
			Width:           1920,
			Height:          1080,
			DurationSeconds: 2400,
			AudioStreams: []movie.AudioStreamInfo{
				{Index: 0, Codec: "aac", Channels: 2, Language: "eng"},
			},
		},
	}

	cfg := testConfig()
	cfg.Playback.SegmentDir = t.TempDir()

	svc, _ := newTestService(t, cfg, nil, tvSvc, prober)

	req := &StartPlaybackRequest{
		MediaType: MediaTypeEpisode,
		MediaID:   uuid.New(),
	}

	sess, err := svc.StartSession(context.Background(), uuid.New(), req)
	require.NoError(t, err)
	require.NotNil(t, sess)

	assert.Equal(t, MediaTypeEpisode, sess.MediaType)
	assert.Equal(t, fileID, sess.FileID)

	_ = svc.StopSession(sess.ID)
}
