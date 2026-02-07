package hls

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/playback"
	"github.com/lusoris/revenge/internal/playback/transcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// audioDisplayName / subtitleDisplayName coverage
// ---------------------------------------------------------------------------

func TestAudioDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		track    playback.AudioTrackInfo
		expected string
	}{
		{
			name:     "uses title when present",
			track:    playback.AudioTrackInfo{Index: 0, Title: "English 5.1", Language: "en"},
			expected: "English 5.1",
		},
		{
			name:     "falls back to language when no title",
			track:    playback.AudioTrackInfo{Index: 0, Title: "", Language: "de"},
			expected: "de",
		},
		{
			name:     "falls back to track index when no title or language",
			track:    playback.AudioTrackInfo{Index: 3, Title: "", Language: ""},
			expected: "Track 3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := audioDisplayName(tc.track)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestSubtitleDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		track    playback.SubtitleTrackInfo
		expected string
	}{
		{
			name:     "uses title when present",
			track:    playback.SubtitleTrackInfo{Index: 0, Title: "English (SDH)", Language: "en"},
			expected: "English (SDH)",
		},
		{
			name:     "falls back to language when no title",
			track:    playback.SubtitleTrackInfo{Index: 1, Title: "", Language: "fr"},
			expected: "fr",
		},
		{
			name:     "falls back to track index when no title or language",
			track:    playback.SubtitleTrackInfo{Index: 5, Title: "", Language: ""},
			expected: "Track 5",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := subtitleDisplayName(tc.track)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// ---------------------------------------------------------------------------
// serveSegment handler tests
// ---------------------------------------------------------------------------

func TestStreamHandler_ServeSegment(t *testing.T) {
	handler, sm := newTestHandler(t)

	// Create a session with a real temp directory containing a segment
	segDir := t.TempDir()
	profileDir := filepath.Join(segDir, "original")
	require.NoError(t, os.MkdirAll(profileDir, 0o755))
	segFile := filepath.Join(profileDir, "seg-00000.ts")
	require.NoError(t, os.WriteFile(segFile, []byte("fake-ts-data"), 0o644))

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
		DurationSeconds: 7200,
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/original/seg-00000.ts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "video/mp2t", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Cache-Control"), "immutable")
	assert.Equal(t, "fake-ts-data", rec.Body.String())
}

func TestStreamHandler_ServeSegment_NotFoundFile(t *testing.T) {
	handler, sm := newTestHandler(t)

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

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/original/seg-99999.ts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// serveSubtitle handler tests
// ---------------------------------------------------------------------------

func TestStreamHandler_ServeSubtitle(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	subsDir := filepath.Join(segDir, "subs")
	require.NoError(t, os.MkdirAll(subsDir, 0o755))
	vttFile := filepath.Join(subsDir, "0.vtt")
	require.NoError(t, os.WriteFile(vttFile, []byte("WEBVTT\n\n00:00:01.000 --> 00:00:04.000\nHello!"), 0o644))

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		FileID:     uuid.Must(uuid.NewV7()),
		FilePath:   "/media/movie.mkv",
		SegmentDir: segDir,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
		SubtitleTracks: []playback.SubtitleTrackInfo{
			{Index: 0, Language: "en", Title: "English"},
		},
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/subs/0.vtt", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/vtt", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "WEBVTT")
}

func TestStreamHandler_ServeSubtitle_InvalidTrackIndex(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: segDir,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/subs/abc.vtt", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// serveAudioRendition handler tests
// ---------------------------------------------------------------------------

func TestStreamHandler_ServeAudioRenditionSegment(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	audioDir := filepath.Join(segDir, "audio", "0")
	require.NoError(t, os.MkdirAll(audioDir, 0o755))
	segFile := filepath.Join(audioDir, "seg-00001.ts")
	require.NoError(t, os.WriteFile(segFile, []byte("fake-audio-data"), 0o644))

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: segDir,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
		AudioTracks: []playback.AudioTrackInfo{
			{Index: 0, Language: "en", Title: "English", Channels: 2, IsDefault: true},
		},
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/audio/0/seg-00001.ts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "video/mp2t", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Cache-Control"), "immutable")
	assert.Equal(t, "fake-audio-data", rec.Body.String())
}

func TestStreamHandler_ServeAudioRendition_InvalidTrackIndex(t *testing.T) {
	handler, sm := newTestHandler(t)

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: t.TempDir(),
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/audio/xyz/index.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestStreamHandler_ServeAudioRendition_NoSlash(t *testing.T) {
	handler, sm := newTestHandler(t)

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: t.TempDir(),
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
	}
	require.NoError(t, sm.Create(sess))

	// "audio/" without track subpath should 404
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/audio/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Empty track string after "audio/" => slashPos < 0, so 404
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestStreamHandler_ServeAudioRendition_UnknownFile(t *testing.T) {
	handler, sm := newTestHandler(t)

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: t.TempDir(),
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
	}
	require.NoError(t, sm.Create(sess))

	// Request for unknown file type under audio rendition
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/audio/0/unknown.txt", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// serveMediaPlaylist handler tests
// ---------------------------------------------------------------------------

func TestStreamHandler_ServeMediaPlaylist(t *testing.T) {
	handler, sm := newTestHandler(t)

	// Create segment dir with a valid media playlist
	segDir := t.TempDir()
	profileDir := filepath.Join(segDir, "original")
	require.NoError(t, os.MkdirAll(profileDir, 0o755))
	playlistContent := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:4\n#EXTINF:4.0,\nseg-00000.ts\n"
	require.NoError(t, os.WriteFile(filepath.Join(profileDir, "index.m3u8"), []byte(playlistContent), 0o644))

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: segDir,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{
				{Name: "original", Width: 1920, Height: 1080, VideoCodec: "copy", AudioCodec: "copy"},
			},
		},
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/original/index.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/vnd.apple.mpegurl", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "#EXTM3U")
	assert.Contains(t, rec.Body.String(), "seg-00000.ts")
}

func TestStreamHandler_ServeMediaPlaylist_Cached(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	profileDir := filepath.Join(segDir, "720p")
	require.NoError(t, os.MkdirAll(profileDir, 0o755))
	playlistContent := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:4\n#EXTINF:4.0,\nseg-00000.ts\n"
	require.NoError(t, os.WriteFile(filepath.Join(profileDir, "index.m3u8"), []byte(playlistContent), 0o644))

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: segDir,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{
				{Name: "720p", Width: 1280, Height: 720, VideoBitrate: 2800, VideoCodec: "libx264", AudioCodec: "aac"},
			},
		},
	}
	require.NoError(t, sm.Create(sess))

	url := "/api/v1/playback/stream/" + sess.ID.String() + "/720p/index.m3u8"

	// First request (populates cache)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, httptest.NewRequest(http.MethodGet, url, nil))
	assert.Equal(t, http.StatusOK, rec1.Code)
	body1 := rec1.Body.String()

	// Delete the file on disk to verify second request serves from cache
	require.NoError(t, os.Remove(filepath.Join(profileDir, "index.m3u8")))

	// Second request (from cache)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, url, nil))
	assert.Equal(t, http.StatusOK, rec2.Code)
	assert.Equal(t, body1, rec2.Body.String())
}

func TestStreamHandler_ServeAudioRenditionPlaylist(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	audioDir := filepath.Join(segDir, "audio", "0")
	require.NoError(t, os.MkdirAll(audioDir, 0o755))
	playlistContent := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:4\n#EXTINF:4.0,\nseg-00000.ts\n"
	require.NoError(t, os.WriteFile(filepath.Join(audioDir, "index.m3u8"), []byte(playlistContent), 0o644))

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: segDir,
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
		AudioTracks: []playback.AudioTrackInfo{
			{Index: 0, Language: "en", Title: "English", Channels: 2, IsDefault: true},
		},
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/audio/0/index.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/vnd.apple.mpegurl", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "#EXTM3U")
}

// ---------------------------------------------------------------------------
// Segment path edge cases without slash
// ---------------------------------------------------------------------------

func TestStreamHandler_SegmentNoSlash(t *testing.T) {
	handler, sm := newTestHandler(t)

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: t.TempDir(),
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{},
		},
	}
	require.NoError(t, sm.Create(sess))

	// seg- at root without profile directory
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/seg-00000.ts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// "seg-00000.ts" has no "/" before "seg-" so slashPos would be -1 => 404
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Master playlist with audio tracks that have no title/language
// ---------------------------------------------------------------------------

func TestStreamHandler_MasterPlaylistAudioFallbackNames(t *testing.T) {
	handler, sm := newTestHandler(t)

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: t.TempDir(),
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{
				{Name: "original", Width: 1920, Height: 1080, VideoCodec: "copy", AudioCodec: "copy"},
			},
		},
		AudioTracks: []playback.AudioTrackInfo{
			{Index: 0, Language: "", Title: "", Channels: 2, IsDefault: true},
		},
		DurationSeconds: 3600,
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/master.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	// Should fall back to "Track 0" since no title or language
	assert.Contains(t, rec.Body.String(), `NAME="Track 0"`)
}

func TestStreamHandler_MasterPlaylistSubtitleFallbackNames(t *testing.T) {
	handler, sm := newTestHandler(t)

	sess := &playback.Session{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     uuid.Must(uuid.NewV7()),
		MediaType:  playback.MediaTypeMovie,
		MediaID:    uuid.Must(uuid.NewV7()),
		SegmentDir: t.TempDir(),
		TranscodeDecision: transcode.Decision{
			Profiles: []transcode.ProfileDecision{
				{Name: "original", Width: 1920, Height: 1080, VideoCodec: "copy", AudioCodec: "copy"},
			},
		},
		SubtitleTracks: []playback.SubtitleTrackInfo{
			{Index: 0, Language: "fr", Title: "", Codec: "subrip"},
			{Index: 1, Language: "", Title: "", Codec: "ass"},
		},
		DurationSeconds: 3600,
	}
	require.NoError(t, sm.Create(sess))

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/master.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	// First subtitle falls back to language
	assert.Contains(t, body, `NAME="fr"`)
	// Second subtitle falls back to "Track 1"
	assert.Contains(t, body, `NAME="Track 1"`)
}
