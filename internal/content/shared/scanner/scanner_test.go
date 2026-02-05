package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockParser is a simple test parser
type mockParser struct {
	extensions []string
}

func (p *mockParser) Parse(filename string) (string, map[string]any) {
	// Simple mock: just return the filename without extension
	ext := filepath.Ext(filename)
	title := filename[:len(filename)-len(ext)]
	return title, map[string]any{"mock": true}
}

func (p *mockParser) GetExtensions() []string {
	return p.extensions
}

func (p *mockParser) ContentType() string {
	return "mock"
}

func TestNewFilesystemScanner(t *testing.T) {
	parser := &mockParser{extensions: []string{".mp4", ".mkv"}}
	scanner := NewFilesystemScanner([]string{"/test/path"}, parser)

	assert.NotNil(t, scanner)
	assert.Equal(t, []string{"/test/path"}, scanner.GetPaths())
	assert.Equal(t, parser, scanner.GetParser())
	assert.False(t, scanner.GetOptions().FollowSymlinks)
}

func TestFilesystemScanner_Scan(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()

	// Create test files
	testFiles := []struct {
		path    string
		content string
	}{
		{"movie1.mp4", "test content"},
		{"movie2.mkv", "test content"},
		{"document.txt", "text document"},
		{"subdir/movie3.mp4", "test content"},
		{".hidden/movie4.mp4", "hidden content"},
	}

	for _, tf := range testFiles {
		fullPath := filepath.Join(tempDir, tf.path)
		dir := filepath.Dir(fullPath)
		require.NoError(t, os.MkdirAll(dir, 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(tf.content), 0644))
	}

	parser := &mockParser{extensions: []string{".mp4", ".mkv"}}
	scanner := NewFilesystemScanner([]string{tempDir}, parser)

	ctx := context.Background()
	results, err := scanner.Scan(ctx)
	require.NoError(t, err)

	// Should find 3 files (movie1.mp4, movie2.mkv, subdir/movie3.mp4)
	// Hidden directory should be skipped by default
	assert.Len(t, results, 3)

	// All results should be marked as media
	for _, r := range results {
		assert.True(t, r.IsMedia)
		assert.Greater(t, r.FileSize, int64(0))
	}
}

func TestFilesystemScanner_ScanWithHiddenFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files
	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, ".hidden"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "movie1.mp4"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, ".hidden", "movie2.mp4"), []byte("test"), 0644))

	parser := &mockParser{extensions: []string{".mp4"}}

	// Without hidden files
	scanner := NewFilesystemScanner([]string{tempDir}, parser)
	results, err := scanner.Scan(context.Background())
	require.NoError(t, err)
	assert.Len(t, results, 1)

	// With hidden files
	opts := DefaultScanOptions()
	opts.IncludeHidden = true
	scanner = NewFilesystemScanner([]string{tempDir}, parser, opts)
	results, err = scanner.Scan(context.Background())
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestFilesystemScanner_ScanWithMaxDepth(t *testing.T) {
	tempDir := t.TempDir()

	// Create nested structure
	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "a", "b", "c"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "movie1.mp4"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "a", "movie2.mp4"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "a", "b", "movie3.mp4"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "a", "b", "c", "movie4.mp4"), []byte("test"), 0644))

	parser := &mockParser{extensions: []string{".mp4"}}

	// MaxDepth 1 should only find files in root and one level down
	opts := DefaultScanOptions()
	opts.MaxDepth = 1
	scanner := NewFilesystemScanner([]string{tempDir}, parser, opts)
	results, err := scanner.Scan(context.Background())
	require.NoError(t, err)
	assert.Len(t, results, 2) // movie1.mp4, a/movie2.mp4
}

func TestFilesystemScanner_ScanWithExcludePatterns(t *testing.T) {
	tempDir := t.TempDir()

	// Create directories
	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "movies"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "@eaDir"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, ".Trash-1000"), 0755))

	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "movies", "movie1.mp4"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "@eaDir", "movie2.mp4"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, ".Trash-1000", "movie3.mp4"), []byte("test"), 0644))

	parser := &mockParser{extensions: []string{".mp4"}}
	scanner := NewFilesystemScanner([]string{tempDir}, parser)

	results, err := scanner.Scan(context.Background())
	require.NoError(t, err)

	// Should only find movies/movie1.mp4
	// @eaDir is in default exclude patterns, .Trash-* is excluded
	assert.Len(t, results, 1)
	assert.Contains(t, results[0].FilePath, "movie1.mp4")
}

func TestFilesystemScanner_ScanWithSummary(t *testing.T) {
	tempDir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "movie1.mp4"), []byte("test content 1"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "movie2.mkv"), []byte("test content 2"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("text"), 0644))

	parser := &mockParser{extensions: []string{".mp4", ".mkv"}}
	scanner := NewFilesystemScanner([]string{tempDir}, parser)

	results, summary, err := scanner.ScanWithSummary(context.Background())
	require.NoError(t, err)

	assert.Len(t, results, 2)
	assert.Equal(t, 2, summary.TotalFiles)
	assert.Equal(t, 2, summary.MediaFiles)
	assert.Equal(t, 2, summary.ParsedFiles)
	assert.Equal(t, 0, summary.FailedParses)
}

func TestFilesystemScanner_ScanContextCancellation(t *testing.T) {
	tempDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "movie1.mp4"), []byte("test"), 0644))

	parser := &mockParser{extensions: []string{".mp4"}}
	scanner := NewFilesystemScanner([]string{tempDir}, parser)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := scanner.Scan(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestFilesystemScanner_ScanNonExistentPath(t *testing.T) {
	parser := &mockParser{extensions: []string{".mp4"}}
	scanner := NewFilesystemScanner([]string{"/nonexistent/path"}, parser)

	_, summary, err := scanner.ScanWithSummary(context.Background())
	require.NoError(t, err) // Should not error, just collect errors in summary
	assert.Len(t, summary.Errors, 1)
}

func TestFilesystemScanner_MultiplePaths(t *testing.T) {
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(tempDir1, "movie1.mp4"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir2, "movie2.mp4"), []byte("test"), 0644))

	parser := &mockParser{extensions: []string{".mp4"}}
	scanner := NewFilesystemScanner([]string{tempDir1, tempDir2}, parser)

	results, err := scanner.Scan(context.Background())
	require.NoError(t, err)
	assert.Len(t, results, 2)
}
