// Package scanner provides a generic file scanning framework for content libraries.
// It supports different content types (movies, TV shows, music) through pluggable
// file parsers that implement the FileParser interface.
package scanner

import (
	"context"
)

// ScanResult represents a discovered media file with parsed metadata.
// The Metadata field is flexible to accommodate different content types.
type ScanResult struct {
	// FilePath is the absolute path to the file
	FilePath string

	// FileName is the base name of the file
	FileName string

	// FileSize in bytes
	FileSize int64

	// ParsedTitle is the extracted title from the filename
	ParsedTitle string

	// Metadata contains content-type-specific parsed data
	// For movies: may contain "year" (int)
	// For TV: may contain "season" (int), "episode" (int), "series" (string)
	// For music: may contain "artist" (string), "album" (string), "track" (int)
	Metadata map[string]any

	// IsMedia indicates if the file was recognized as a media file
	IsMedia bool

	// Error contains any error that occurred during parsing
	Error error
}

// GetYear extracts year from metadata (commonly used for movies)
func (r ScanResult) GetYear() *int {
	if r.Metadata == nil {
		return nil
	}
	if v, ok := r.Metadata["year"]; ok {
		if year, ok := v.(int); ok {
			return &year
		}
	}
	return nil
}

// GetSeason extracts season number from metadata (used for TV shows)
func (r ScanResult) GetSeason() *int {
	if r.Metadata == nil {
		return nil
	}
	if v, ok := r.Metadata["season"]; ok {
		if season, ok := v.(int); ok {
			return &season
		}
	}
	return nil
}

// GetEpisode extracts episode number from metadata (used for TV shows)
func (r ScanResult) GetEpisode() *int {
	if r.Metadata == nil {
		return nil
	}
	if v, ok := r.Metadata["episode"]; ok {
		if episode, ok := v.(int); ok {
			return &episode
		}
	}
	return nil
}

// GetString extracts a string value from metadata
func (r ScanResult) GetString(key string) string {
	if r.Metadata == nil {
		return ""
	}
	if v, ok := r.Metadata[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// ScanOptions configures the scanning behavior
type ScanOptions struct {
	// FollowSymlinks determines whether to follow symbolic links
	FollowSymlinks bool

	// MaxDepth limits the directory traversal depth (0 = unlimited)
	MaxDepth int

	// ExcludePatterns are glob patterns for paths to skip
	ExcludePatterns []string

	// IncludeHidden determines whether to scan hidden files/directories
	IncludeHidden bool
}

// DefaultScanOptions returns sensible defaults for library scanning
func DefaultScanOptions() ScanOptions {
	return ScanOptions{
		FollowSymlinks:  false,
		MaxDepth:        0, // unlimited
		ExcludePatterns: []string{".Trash*", ".recycle*", "@eaDir", ".DS_Store"},
		IncludeHidden:   false,
	}
}

// ScanSummary provides statistics about a completed scan
type ScanSummary struct {
	TotalFiles      int
	MediaFiles      int
	SkippedFiles    int
	ParsedFiles     int
	FailedParses    int
	DirectoriesRead int
	Errors          []error
}

// FileParser is the interface that content-specific parsers must implement.
// Each content type (movie, TV, music) provides its own parser implementation.
type FileParser interface {
	// Parse extracts metadata from a filename.
	// Returns the cleaned title and any additional metadata.
	Parse(filename string) (title string, metadata map[string]any)

	// GetExtensions returns the file extensions this parser handles.
	// Extensions should include the leading dot, e.g., ".mp4", ".mkv"
	GetExtensions() []string

	// ContentType returns a string identifier for this parser's content type.
	// Used for logging and debugging.
	ContentType() string
}

// Scanner handles file system scanning for media files using pluggable parsers
type Scanner interface {
	// Scan scans all configured paths and returns discovered media files
	Scan(ctx context.Context) ([]ScanResult, error)

	// ScanWithSummary scans and returns both results and statistics
	ScanWithSummary(ctx context.Context) ([]ScanResult, *ScanSummary, error)

	// ScanPath scans a single path
	ScanPath(ctx context.Context, path string) ([]ScanResult, error)
}
