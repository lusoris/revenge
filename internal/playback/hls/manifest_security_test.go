package hls

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===========================================================================
// ReadMediaPlaylist — path traversal (CWE-22)
// ===========================================================================

func TestReadMediaPlaylist_PathTraversal(t *testing.T) {
	segDir := t.TempDir()

	// Create a file outside the segment dir tree
	parentDir := filepath.Dir(segDir)
	secretFile := filepath.Join(parentDir, "secret", "index.m3u8")
	require.NoError(t, os.MkdirAll(filepath.Dir(secretFile), 0o755))
	require.NoError(t, os.WriteFile(secretFile, []byte("#EXTM3U\n#SECRET"), 0o644))
	t.Cleanup(func() { os.RemoveAll(filepath.Join(parentDir, "secret")) })

	attacks := []struct {
		name    string
		profile string
	}{
		{"dot dot", "../secret"},
		{"dot dot double", "../../secret"},
		{"dot", "."},
		{"dot dot bare", ".."},
		{"empty component", ""},
		{"slash prefix", "/etc"},
	}

	for _, tc := range attacks {
		t.Run(tc.name, func(t *testing.T) {
			content, err := ReadMediaPlaylist(segDir, tc.profile)
			assert.Error(t, err, "profile %q should be rejected", tc.profile)
			assert.Empty(t, content)
			assert.NotContains(t, content, "SECRET")
		})
	}
}

// ===========================================================================
// ReadMediaPlaylist — valid playlist from disk
// ===========================================================================

func TestReadMediaPlaylist_ValidFile(t *testing.T) {
	segDir := t.TempDir()
	profileDir := filepath.Join(segDir, "720p")
	require.NoError(t, os.MkdirAll(profileDir, 0o755))

	playlistContent := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:4\n#EXTINF:4.000,\nseg-00000.ts\n#EXTINF:4.000,\nseg-00001.ts\n"
	require.NoError(t, os.WriteFile(filepath.Join(profileDir, "index.m3u8"), []byte(playlistContent), 0o644))

	content, err := ReadMediaPlaylist(segDir, "720p")
	require.NoError(t, err)
	assert.Contains(t, content, "#EXTM3U")
	assert.Contains(t, content, "seg-00000.ts")
	assert.Contains(t, content, "seg-00001.ts")
}

// ===========================================================================
// ReadMediaPlaylist — file not found (polling exhausted)
// ===========================================================================

func TestReadMediaPlaylist_FileNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow poll test in short mode")
	}

	segDir := t.TempDir()
	// No file created — should poll and time out

	_, err := ReadMediaPlaylist(segDir, "original")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not available")
}

// ===========================================================================
// SegmentPath — defense in depth
// ===========================================================================

func TestSegmentPath_DefenseInDepth(t *testing.T) {
	tests := []struct {
		name        string
		segDir      string
		profile     string
		segFile     string
		expected    string
	}{
		{
			name:     "normal path",
			segDir:   "/tmp/segments",
			profile:  "original",
			segFile:  "seg-00000.ts",
			expected: "/tmp/segments/original/seg-00000.ts",
		},
		{
			name:     "defense: strips directory from profile",
			segDir:   "/tmp/segments",
			profile:  "../../etc",
			segFile:  "seg-00000.ts",
			expected: "/tmp/segments/etc/seg-00000.ts",
		},
		{
			name:     "defense: strips directory from segment",
			segDir:   "/tmp/segments",
			profile:  "original",
			segFile:  "../../passwd",
			expected: "/tmp/segments/original/passwd",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SegmentPath(tc.segDir, tc.profile, tc.segFile)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// ===========================================================================
// SubtitlePath — various track indices
// ===========================================================================

func TestSubtitlePath_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		segDir   string
		track    int
		expected string
	}{
		{"zero index", "/segments", 0, "/segments/subs/0.vtt"},
		{"high index", "/segments", 99, "/segments/subs/99.vtt"},
		{"deep dir", "/data/revenge/sessions/abc", 3, "/data/revenge/sessions/abc/subs/3.vtt"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SubtitlePath(tc.segDir, tc.track)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// ===========================================================================
// AudioRenditionSegmentPath — edge cases
// ===========================================================================

func TestAudioRenditionSegmentPath_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		segDir   string
		track    int
		segFile  string
		expected string
	}{
		{"first track first seg", "/tmp/s", 0, "seg-00000.ts", "/tmp/s/audio/0/seg-00000.ts"},
		{"high track", "/tmp/s", 15, "seg-00100.ts", "/tmp/s/audio/15/seg-00100.ts"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := AudioRenditionSegmentPath(tc.segDir, tc.track, tc.segFile)
			assert.Equal(t, tc.expected, got)
		})
	}
}
