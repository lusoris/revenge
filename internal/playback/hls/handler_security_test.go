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

// ===========================================================================
// isSafePathComponent — exhaustive coverage
// ===========================================================================

func TestIsSafePathComponent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		safe  bool
	}{
		// Unsafe
		{name: "empty string", input: "", safe: false},
		{name: "single dot", input: ".", safe: false},
		{name: "double dot", input: "..", safe: false},
		{name: "forward slash", input: "a/b", safe: false},
		{name: "backslash", input: `a\b`, safe: false},
		{name: "dot dot in middle", input: "foo..bar", safe: false},
		{name: "dot dot prefix", input: "../etc", safe: false},
		{name: "dot dot suffix", input: "etc/..", safe: false},
		{name: "slash only", input: "/", safe: false},

		// Safe
		{name: "simple name", input: "original", safe: true},
		{name: "720p", input: "720p", safe: true},
		{name: "segment file", input: "seg-00001.m4s", safe: true},
		{name: "single char", input: "a", safe: true},
		{name: "with hyphen", input: "my-profile", safe: true},
		{name: "with underscore", input: "my_profile", safe: true},
		{name: "with number", input: "42", safe: true},
		{name: "index.m3u8", input: "index.m3u8", safe: true},
		{name: "single dot in name", input: "file.m4s", safe: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isSafePathComponent(tc.input)
			assert.Equal(t, tc.safe, got, "isSafePathComponent(%q)", tc.input)
		})
	}
}

// ===========================================================================
// Path traversal attacks on segment serving (CWE-22)
// ===========================================================================

func TestStreamHandler_SegmentPathTraversal(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	// Create a secret file outside the profile dirs to verify it's unreachable
	secretFile := filepath.Join(segDir, "secret.txt")
	require.NoError(t, os.WriteFile(secretFile, []byte("TOP SECRET"), 0o644))

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

	base := "/api/v1/playback/stream/" + sess.ID.String()

	attacks := []struct {
		name string
		path string
	}{
		{"dot dot profile", base + "/../../../etc/passwd/seg-00000.m4s"},
		{"dot dot segment file", base + "/original/..%2F..%2Fsecret.txt"},
		{"backslash in profile", base + `/original\..\..\seg-00000.m4s`},
		{"dot dot in profile name", base + "/../seg-00000.m4s"},
	}

	for _, tc := range attacks {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			// Must not serve the file — either 400 or 404
			assert.NotEqual(t, http.StatusOK, rec.Code,
				"path traversal attempt should be rejected")
			assert.NotContains(t, rec.Body.String(), "TOP SECRET",
				"secret content must not leak")
		})
	}
}

// ===========================================================================
// Path traversal attacks on audio rendition segments
// ===========================================================================

func TestStreamHandler_AudioRenditionPathTraversal(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(segDir, "secret.txt"), []byte("LEAKED"), 0o644))

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

	base := "/api/v1/playback/stream/" + sess.ID.String()

	attacks := []struct {
		name string
		path string
	}{
		{"dot dot audio segment", base + "/audio/0/../../secret.txt"},
		{"backslash audio segment", base + `/audio/0/..\..\secret.txt`},
	}

	for _, tc := range attacks {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			assert.NotContains(t, rec.Body.String(), "LEAKED",
				"path traversal should not leak content")
		})
	}
}

// ===========================================================================
// Media playlist path traversal via malicious profile names
// ===========================================================================

func TestStreamHandler_MediaPlaylistPathTraversal(t *testing.T) {
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

	base := "/api/v1/playback/stream/" + sess.ID.String()

	attacks := []struct {
		name string
		path string
	}{
		{"dot dot profile", base + "/../../../etc/passwd/index.m3u8"},
		{"empty profile", base + "//index.m3u8"},
		{"dot profile", base + "/./index.m3u8"},
		{"dot dot only", base + "/../index.m3u8"},
	}

	for _, tc := range attacks {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			// Should be 400 (invalid profile) or 404
			code := rec.Code
			assert.True(t, code == http.StatusBadRequest || code == http.StatusNotFound || code == http.StatusServiceUnavailable,
				"traversal profile should be rejected, got %d", code)
		})
	}
}

// ===========================================================================
// Subtitle track injection
// ===========================================================================

func TestStreamHandler_SubtitlePathTraversal(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(segDir, "secret.vtt"), []byte("CONFIDENTIAL"), 0o644))

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

	base := "/api/v1/playback/stream/" + sess.ID.String()

	// Non-numeric track index should 404 (atoi fails)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, base+"/subs/../secret.vtt", nil))
	assert.NotContains(t, rec.Body.String(), "CONFIDENTIAL")
}

// ===========================================================================
// HTTP method enforcement on all paths
// ===========================================================================

func TestStreamHandler_MethodNotAllowed_AllPaths(t *testing.T) {
	handler, sm := newTestHandler(t)

	sess := createTestSession(t, sm)
	base := "/api/v1/playback/stream/" + sess.ID.String()

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	paths := []string{
		base + "/master.m3u8",
		base + "/original/index.m3u8",
		base + "/original/seg-00000.m4s",
		base + "/audio/0/index.m3u8",
		base + "/subs/0.vtt",
	}

	for _, method := range methods {
		for _, path := range paths {
			t.Run(method+"_"+path, func(t *testing.T) {
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
				assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
			})
		}
	}
}

// ===========================================================================
// Master playlist no-cache header
// ===========================================================================

func TestStreamHandler_MasterPlaylistNoCacheHeader(t *testing.T) {
	handler, sm := newTestHandler(t)
	sess := createTestSession(t, sm)

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/playback/stream/"+sess.ID.String()+"/master.m3u8", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
}

// ===========================================================================
// Segment immutable cache header
// ===========================================================================

func TestStreamHandler_SegmentImmutableCacheHeader(t *testing.T) {
	handler, sm := newTestHandler(t)

	segDir := t.TempDir()
	profileDir := filepath.Join(segDir, "original")
	require.NoError(t, os.MkdirAll(profileDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(profileDir, "seg-00000.m4s"), []byte("fmp4-data"), 0o644))

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
		"/api/v1/playback/stream/"+sess.ID.String()+"/original/seg-00000.m4s", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	cacheControl := rec.Header().Get("Cache-Control")
	assert.Contains(t, cacheControl, "immutable")
	assert.Contains(t, cacheControl, "max-age=31536000")
}

// ===========================================================================
// Handler Close is safe to call twice
// ===========================================================================

func TestStreamHandler_Close_Idempotent(t *testing.T) {
	handler, _ := newTestHandler(t)

	// Should not panic
	handler.Close()
	handler.Close()
}
