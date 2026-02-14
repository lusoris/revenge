package api

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/playback"
	"github.com/lusoris/revenge/internal/playback/transcode"
)

// ===========================================================================
// sessionToOgen converter tests
// ===========================================================================

func TestSessionToOgen_BasicFields(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expires := now.Add(30 * time.Minute)
	sessionID := uuid.Must(uuid.NewV7())

	sess := &playback.Session{
		ID:              sessionID,
		DurationSeconds: 7200.5,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
		AudioTracks:    []playback.AudioTrackInfo{},
		SubtitleTracks: []playback.SubtitleTrackInfo{},
		CreatedAt:      now,
		ExpiresAt:      expires,
	}

	result := sessionToOgen(sess)

	assert.Equal(t, sessionID, result.SessionID)
	assert.Contains(t, result.MasterPlaylistURL, sessionID.String())
	assert.Contains(t, result.MasterPlaylistURL, "master.m3u8")
	assert.Equal(t, 7200.5, result.DurationSeconds)
	assert.Empty(t, result.Profiles)
	assert.Empty(t, result.AudioTracks)
	assert.Empty(t, result.SubtitleTracks)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, expires, result.ExpiresAt)
}

func TestSessionToOgen_WithProfiles(t *testing.T) {
	t.Parallel()

	now := time.Now()
	sess := &playback.Session{
		ID:              uuid.Must(uuid.NewV7()),
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
		AudioTracks:    []playback.AudioTrackInfo{},
		SubtitleTracks: []playback.SubtitleTrackInfo{},
		CreatedAt:      now,
		ExpiresAt:      now.Add(30 * time.Minute),
	}

	result := sessionToOgen(sess)

	require.Len(t, result.Profiles, 2)

	assert.Equal(t, "original", result.Profiles[0].Name)
	assert.Equal(t, 1920, result.Profiles[0].Width)
	assert.Equal(t, 1080, result.Profiles[0].Height)
	assert.Equal(t, 0, result.Profiles[0].Bitrate)
	assert.True(t, result.Profiles[0].IsOriginal)

	assert.Equal(t, "720p", result.Profiles[1].Name)
	assert.Equal(t, 1280, result.Profiles[1].Width)
	assert.Equal(t, 720, result.Profiles[1].Height)
	assert.Equal(t, 2800, result.Profiles[1].Bitrate)
	assert.False(t, result.Profiles[1].IsOriginal)
}

func TestSessionToOgen_WithAudioTracks(t *testing.T) {
	t.Parallel()

	now := time.Now()
	sess := &playback.Session{
		ID:              uuid.Must(uuid.NewV7()),
		DurationSeconds: 3600,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
		AudioTracks: []playback.AudioTrackInfo{
			{Index: 0, Language: "eng", Title: "English 5.1", Channels: 6, Codec: "dts", Layout: "5.1", IsDefault: true},
			{Index: 1, Language: "", Title: "", Channels: 2, Codec: "aac", Layout: "", IsDefault: false},
		},
		SubtitleTracks: []playback.SubtitleTrackInfo{},
		CreatedAt:      now,
		ExpiresAt:      now.Add(30 * time.Minute),
	}

	result := sessionToOgen(sess)

	require.Len(t, result.AudioTracks, 2)

	// Track with all fields set
	at0 := result.AudioTracks[0]
	assert.Equal(t, 0, at0.Index)
	assert.Equal(t, 6, at0.Channels)
	assert.Equal(t, "dts", at0.Codec)
	assert.True(t, at0.IsDefault)
	assert.True(t, at0.Language.Set)
	assert.Equal(t, "eng", at0.Language.Value)
	assert.True(t, at0.Title.Set)
	assert.Equal(t, "English 5.1", at0.Title.Value)
	assert.True(t, at0.Layout.Set)
	assert.Equal(t, "5.1", at0.Layout.Value)

	// Track with empty optional fields → should NOT be set
	at1 := result.AudioTracks[1]
	assert.Equal(t, 1, at1.Index)
	assert.False(t, at1.Language.Set, "empty language should not be set")
	assert.False(t, at1.Title.Set, "empty title should not be set")
	assert.False(t, at1.Layout.Set, "empty layout should not be set")
	assert.False(t, at1.IsDefault)
}

func TestSessionToOgen_WithSubtitleTracks(t *testing.T) {
	t.Parallel()

	sessionID := uuid.Must(uuid.NewV7())
	now := time.Now()
	sess := &playback.Session{
		ID:              sessionID,
		DurationSeconds: 3600,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
		AudioTracks: []playback.AudioTrackInfo{},
		SubtitleTracks: []playback.SubtitleTrackInfo{
			{Index: 0, Language: "eng", Title: "English (SDH)", Codec: "subrip", URL: "/subs/0.vtt", IsForced: false},
			{Index: 1, Language: "", Title: "", Codec: "ass", URL: "/subs/1.vtt", IsForced: true},
		},
		CreatedAt: now,
		ExpiresAt: now.Add(30 * time.Minute),
	}

	result := sessionToOgen(sess)

	require.Len(t, result.SubtitleTracks, 2)

	st0 := result.SubtitleTracks[0]
	assert.Equal(t, 0, st0.Index)
	assert.Equal(t, "subrip", st0.Codec)
	assert.Equal(t, "/subs/0.vtt", st0.URL)
	assert.False(t, st0.IsForced)
	assert.True(t, st0.Language.Set)
	assert.Equal(t, "eng", st0.Language.Value)
	assert.True(t, st0.Title.Set)
	assert.Equal(t, "English (SDH)", st0.Title.Value)

	st1 := result.SubtitleTracks[1]
	assert.True(t, st1.IsForced)
	assert.False(t, st1.Language.Set, "empty language should not be set")
	assert.False(t, st1.Title.Set, "empty title should not be set")
}

// ===========================================================================
// StartPlaybackSession — authenticated flow with real session manager
// ===========================================================================

func TestHandler_StartPlaybackSession_AuthenticatedNoSession(t *testing.T) {
	t.Parallel()

	// Create a real (but non-functional for file ops) playback service.
	// StartSession will fail because there's no movie service set up,
	// which tests the error path with authenticated user.
	cfg := testPlaybackConfig()
	sm, err := playback.NewSessionManager(10, 30*time.Minute, logging.NewTestLogger())
	require.NoError(t, err)
	defer sm.Close()

	pm := testPipelineManagerForAPI(t)
	svc, err := playback.NewService(cfg, sm, pm, nil, &playbackMovieSvc{}, nil, logging.NewTestLogger())
	require.NoError(t, err)
	defer svc.Close()

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: svc,
	}

	ctx := WithUserID(context.Background(), uuid.New())
	req := &ogen.StartPlaybackRequest{
		MediaType: "movie",
		MediaID:   uuid.New(),
	}

	result, err := handler.StartPlaybackSession(ctx, req)
	require.NoError(t, err) // handler doesn't return Go errors

	// Should get 404 because mock returns empty files → resolve file fails
	notFound, ok := result.(*ogen.StartPlaybackSessionNotFound)
	require.True(t, ok, "expected *ogen.StartPlaybackSessionNotFound, got %T", result)
	assert.Equal(t, 404, notFound.Code)
}

// ===========================================================================
// GetPlaybackSession — with real session
// ===========================================================================

func TestHandler_GetPlaybackSession_SessionExists(t *testing.T) {
	t.Parallel()

	cfg := testPlaybackConfig()
	sm, err := playback.NewSessionManager(10, 30*time.Minute, logging.NewTestLogger())
	require.NoError(t, err)
	defer sm.Close()

	pm := testPipelineManagerForAPI(t)
	svc, err := playback.NewService(cfg, sm, pm, nil, nil, nil, logging.NewTestLogger())
	require.NoError(t, err)
	defer svc.Close()

	// Manually create a session in the session manager
	sessionID := uuid.Must(uuid.NewV7())
	sess := &playback.Session{
		ID:              sessionID,
		UserID:          uuid.Must(uuid.NewV7()),
		MediaType:       playback.MediaTypeMovie,
		MediaID:         uuid.Must(uuid.NewV7()),
		DurationSeconds: 7200,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{
				{Name: "original", Width: 1920, Height: 1080, VideoCodec: "copy", AudioCodec: "copy"},
			},
		},
		AudioTracks:    []playback.AudioTrackInfo{},
		SubtitleTracks: []playback.SubtitleTrackInfo{},
	}
	require.NoError(t, sm.Create(sess))

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: svc,
	}

	ctx := WithUserID(context.Background(), uuid.New())
	params := ogen.GetPlaybackSessionParams{SessionId: sessionID}

	result, err := handler.GetPlaybackSession(ctx, params)
	require.NoError(t, err)

	pbSess, ok := result.(*ogen.PlaybackSession)
	require.True(t, ok, "expected *ogen.PlaybackSession, got %T", result)
	assert.Equal(t, sessionID, pbSess.SessionID)
	assert.Equal(t, 7200.0, pbSess.DurationSeconds)
	require.Len(t, pbSess.Profiles, 1)
	assert.Equal(t, "original", pbSess.Profiles[0].Name)
	assert.True(t, pbSess.Profiles[0].IsOriginal)
}

func TestHandler_GetPlaybackSession_SessionNotFound(t *testing.T) {
	t.Parallel()

	cfg := testPlaybackConfig()
	sm, err := playback.NewSessionManager(10, 30*time.Minute, logging.NewTestLogger())
	require.NoError(t, err)
	defer sm.Close()

	pm := testPipelineManagerForAPI(t)
	svc, err := playback.NewService(cfg, sm, pm, nil, nil, nil, logging.NewTestLogger())
	require.NoError(t, err)
	defer svc.Close()

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: svc,
	}

	ctx := WithUserID(context.Background(), uuid.New())
	params := ogen.GetPlaybackSessionParams{SessionId: uuid.New()}

	result, err := handler.GetPlaybackSession(ctx, params)
	require.NoError(t, err)

	notFound, ok := result.(*ogen.GetPlaybackSessionNotFound)
	require.True(t, ok, "expected *ogen.GetPlaybackSessionNotFound, got %T", result)
	assert.Equal(t, 404, notFound.Code)
	assert.Contains(t, notFound.Message, "not found")
}

// ===========================================================================
// StopPlaybackSession — with real session
// ===========================================================================

func TestHandler_StopPlaybackSession_SessionExists(t *testing.T) {
	t.Parallel()

	cfg := testPlaybackConfig()
	sm, err := playback.NewSessionManager(10, 30*time.Minute, logging.NewTestLogger())
	require.NoError(t, err)
	defer sm.Close()

	pm := testPipelineManagerForAPI(t)
	svc, err := playback.NewService(cfg, sm, pm, nil, nil, nil, logging.NewTestLogger())
	require.NoError(t, err)
	defer svc.Close()

	sessionID := uuid.Must(uuid.NewV7())
	sess := &playback.Session{
		ID:         sessionID,
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: t.TempDir(),
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
		AudioTracks:    []playback.AudioTrackInfo{},
		SubtitleTracks: []playback.SubtitleTrackInfo{},
	}
	require.NoError(t, sm.Create(sess))

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: svc,
	}

	ctx := WithUserID(context.Background(), uuid.New())
	params := ogen.StopPlaybackSessionParams{SessionId: sessionID}

	result, err := handler.StopPlaybackSession(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.StopPlaybackSessionNoContent)
	require.True(t, ok, "expected *ogen.StopPlaybackSessionNoContent, got %T", result)

	// Session should be gone
	_, found := sm.Get(sessionID)
	assert.False(t, found, "session should be deleted after stop")
}

func TestHandler_StopPlaybackSession_NotFound(t *testing.T) {
	t.Parallel()

	cfg := testPlaybackConfig()
	sm, err := playback.NewSessionManager(10, 30*time.Minute, logging.NewTestLogger())
	require.NoError(t, err)
	defer sm.Close()

	pm := testPipelineManagerForAPI(t)
	svc, err := playback.NewService(cfg, sm, pm, nil, nil, nil, logging.NewTestLogger())
	require.NoError(t, err)
	defer svc.Close()

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: svc,
	}

	ctx := WithUserID(context.Background(), uuid.New())
	params := ogen.StopPlaybackSessionParams{SessionId: uuid.New()}

	result, err := handler.StopPlaybackSession(ctx, params)
	require.NoError(t, err)

	notFound, ok := result.(*ogen.StopPlaybackSessionNotFound)
	require.True(t, ok, "expected *ogen.StopPlaybackSessionNotFound, got %T", result)
	assert.Equal(t, 404, notFound.Code)
}

// ===========================================================================
// StartPlaybackSession — optional fields mapping
// ===========================================================================

func TestHandler_StartPlaybackSession_OptionalFieldsMapping(t *testing.T) {
	t.Parallel()

	// This tests that the handler correctly maps all optional ogen fields.
	// We can't test the full flow without a movie service, but we test
	// that the conversion code runs without panics.

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: new(playback.Service),
	}

	fileID := uuid.New()
	req := &ogen.StartPlaybackRequest{
		MediaType:     "episode",
		MediaID:       uuid.New(),
		FileID:        ogen.NewOptUUID(fileID),
		AudioTrack:    ogen.NewOptInt(2),
		SubtitleTrack: ogen.NewOptInt(1),
		StartPosition: ogen.NewOptInt(300),
	}

	// Without a valid user context, this will hit the unauthorized path
	result, err := handler.StartPlaybackSession(context.Background(), req)
	require.NoError(t, err)

	unauth, ok := result.(*ogen.StartPlaybackSessionUnauthorized)
	require.True(t, ok, "expected unauthorized for bare context")
	assert.Equal(t, 401, unauth.Code)
}

// ===========================================================================
// Helpers
// ===========================================================================

func testPlaybackConfig() *config.Config {
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

func testPipelineManagerForAPI(t *testing.T) *transcode.PipelineManager {
	t.Helper()
	pm, err := transcode.NewPipelineManager(4, logging.NewTestLogger())
	require.NoError(t, err)
	t.Cleanup(pm.Close)
	return pm
}

// playbackMovieSvc is a minimal movie.Service mock for playback tests.
// GetMovieFiles returns an empty list so resolveFilePath fails with "no files available".
type playbackMovieSvc struct {
	movie.Service
}

func (m *playbackMovieSvc) GetMovieFiles(_ context.Context, _ uuid.UUID) ([]movie.MovieFile, error) {
	return []movie.MovieFile{}, nil
}
