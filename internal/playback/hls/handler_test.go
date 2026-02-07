package hls

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/playback"
	"github.com/lusoris/revenge/internal/playback/transcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHandler(t *testing.T) (*StreamHandler, *playback.SessionManager) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	sm, err := playback.NewSessionManager(10, 30*time.Minute, logger)
	require.NoError(t, err)
	t.Cleanup(sm.Close)

	handler, err := NewStreamHandler(sm, logger)
	require.NoError(t, err)
	t.Cleanup(handler.Close)

	return handler, sm
}

func createTestSession(t *testing.T, sm *playback.SessionManager) *playback.Session {
	t.Helper()
	sess := &playback.Session{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    uuid.Must(uuid.NewV7()),
		MediaType: playback.MediaTypeMovie,
		MediaID:   uuid.Must(uuid.NewV7()),
		FileID:    uuid.Must(uuid.NewV7()),
		FilePath:  "/media/movie.mkv",
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{
				{Name: "original", Width: 1920, Height: 1080, VideoCodec: "copy", AudioCodec: "copy"},
				{Name: "720p", Width: 1280, Height: 720, VideoBitrate: 2800, AudioBitrate: 128, VideoCodec: "libx264", AudioCodec: "aac"},
			},
		},
		AudioTracks: []playback.AudioTrackInfo{
			{Index: 0, Language: "en", Title: "English", Channels: 6, IsDefault: true},
			{Index: 1, Language: "de", Title: "German", Channels: 2, IsDefault: false},
		},
		SubtitleTracks: []playback.SubtitleTrackInfo{
			{Index: 0, Language: "en", Title: "English"},
		},
		DurationSeconds: 7200,
	}
	require.NoError(t, sm.Create(sess))
	return sess
}

func TestStreamHandler_MasterPlaylist(t *testing.T) {
	handler, sm := newTestHandler(t)
	sess := createTestSession(t, sm)

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/master.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/vnd.apple.mpegurl", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "#EXTM3U")
	assert.Contains(t, rec.Body.String(), "original/index.m3u8")
	assert.Contains(t, rec.Body.String(), "720p/index.m3u8")
	assert.Contains(t, rec.Body.String(), `TYPE=AUDIO,GROUP-ID="audio",NAME="English"`)
	assert.Contains(t, rec.Body.String(), `URI="audio/0/index.m3u8"`)
	assert.Contains(t, rec.Body.String(), `TYPE=AUDIO,GROUP-ID="audio",NAME="German"`)
	assert.Contains(t, rec.Body.String(), `URI="audio/1/index.m3u8"`)
	assert.Contains(t, rec.Body.String(), `AUDIO="audio"`)
	assert.Contains(t, rec.Body.String(), `SUBTITLES="subs"`)
}

func TestStreamHandler_MasterPlaylistCached(t *testing.T) {
	handler, sm := newTestHandler(t)
	sess := createTestSession(t, sm)

	url := "/api/v1/playback/stream/" + sess.ID.String() + "/master.m3u8"

	// First request generates
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, httptest.NewRequest(http.MethodGet, url, nil))
	assert.Equal(t, http.StatusOK, rec1.Code)
	body1 := rec1.Body.String()

	// Second request served from cache
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, url, nil))
	assert.Equal(t, http.StatusOK, rec2.Code)
	assert.Equal(t, body1, rec2.Body.String())
}

func TestStreamHandler_InvalidSession(t *testing.T) {
	handler, _ := newTestHandler(t)

	fakeID := uuid.Must(uuid.NewV7())
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+fakeID.String()+"/master.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestStreamHandler_InvalidSessionID(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/not-a-uuid/master.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestStreamHandler_MethodNotAllowed(t *testing.T) {
	handler, sm := newTestHandler(t)
	sess := createTestSession(t, sm)

	req := httptest.NewRequest(http.MethodPost,
		"/api/v1/playback/stream/"+sess.ID.String()+"/master.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestStreamHandler_CORSHeaders(t *testing.T) {
	handler, sm := newTestHandler(t)
	sess := createTestSession(t, sm)

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/master.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET", rec.Header().Get("Access-Control-Allow-Methods"))
}

func TestStreamHandler_NotFoundPaths(t *testing.T) {
	handler, sm := newTestHandler(t)
	sess := createTestSession(t, sm)
	base := "/api/v1/playback/stream/" + sess.ID.String()

	tests := []struct {
		name string
		path string
	}{
		{"unknown file", base + "/unknown.txt"},
		{"too short", "/api/v1/playback/stream/" + sess.ID.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))
			assert.Equal(t, http.StatusNotFound, rec.Code)
		})
	}
}
