// Package library provides shared types and interfaces for content library management.
// It defines common patterns for scanning, matching, and managing media libraries
// across different content types (movies, TV shows, music).
package library

import (
	"context"

	"github.com/google/uuid"
)

// ScanSummary contains statistics from a library scan.
// This is used by all content types to report scan results.
type ScanSummary struct {
	// TotalFiles is the number of files scanned.
	TotalFiles int

	// MatchedFiles is the number of files successfully matched to content.
	MatchedFiles int

	// UnmatchedFiles is the number of files that could not be matched.
	UnmatchedFiles int

	// NewContent is the number of new content items created.
	NewContent int

	// ExistingContent is the number of files matched to existing content.
	ExistingContent int

	// Errors contains any errors encountered during scanning.
	Errors []error
}

// AddError adds an error to the summary.
func (s *ScanSummary) AddError(err error) {
	s.Errors = append(s.Errors, err)
}

// HasErrors returns true if any errors occurred.
func (s *ScanSummary) HasErrors() bool {
	return len(s.Errors) > 0
}

// MatchType indicates how content was matched.
type MatchType string

const (
	// MatchTypeExact indicates an exact ID match (e.g., TMDb ID).
	MatchTypeExact MatchType = "exact"

	// MatchTypeTitle indicates a title and year match.
	MatchTypeTitle MatchType = "title"

	// MatchTypeFuzzy indicates a fuzzy title match.
	MatchTypeFuzzy MatchType = "fuzzy"

	// MatchTypeManual indicates a manually matched item.
	MatchTypeManual MatchType = "manual"

	// MatchTypeUnmatched indicates the item could not be matched.
	MatchTypeUnmatched MatchType = "unmatched"
)

// MatchResult represents the result of matching a file to content.
// This is a generic result that can be extended by content-specific packages.
type MatchResult[T any] struct {
	// FilePath is the path to the matched file.
	FilePath string

	// Content is the matched content item (nil if unmatched).
	Content *T

	// MatchType indicates how the match was made.
	MatchType MatchType

	// Confidence is the match confidence score (0.0 to 1.0).
	Confidence float64

	// Error contains any error that occurred during matching.
	Error error

	// CreatedNew indicates whether new content was created.
	CreatedNew bool
}

// IsMatched returns true if the file was matched to content.
func (r MatchResult[T]) IsMatched() bool {
	return r.Content != nil && r.Error == nil
}

// LibraryScanner defines the interface for scanning media libraries.
// Content-specific implementations handle different media types.
type LibraryScanner interface {
	// Scan scans all configured library paths and returns discovered files.
	Scan(ctx context.Context) ([]ScanItem, error)
}

// ScanItem represents a discovered file during scanning.
type ScanItem struct {
	// FilePath is the absolute path to the file.
	FilePath string

	// FileName is the base name of the file.
	FileName string

	// FileSize in bytes.
	FileSize int64

	// ParsedTitle is the extracted title from the filename.
	ParsedTitle string

	// Metadata contains additional parsed data (year, season, episode, etc.).
	Metadata map[string]any

	// IsMedia indicates if the file was recognized as a media file.
	IsMedia bool

	// Error contains any error that occurred during parsing.
	Error error
}

// GetYear extracts the year from metadata (convenience method).
func (s ScanItem) GetYear() *int {
	if s.Metadata == nil {
		return nil
	}
	if v, ok := s.Metadata["year"]; ok {
		if year, ok := v.(int); ok {
			return &year
		}
	}
	return nil
}

// ContentMatcher defines the interface for matching files to content.
// Content-specific implementations handle different matching strategies.
type ContentMatcher[T any] interface {
	// MatchFile attempts to match a single scan item to content.
	MatchFile(ctx context.Context, item ScanItem) MatchResult[T]

	// MatchFiles attempts to match multiple scan items to content.
	MatchFiles(ctx context.Context, items []ScanItem) []MatchResult[T]
}

// MediaFileInfo contains technical information about a media file.
// This is populated by media probers (FFmpeg, MediaInfo, etc.).
type MediaFileInfo struct {
	// Path is the file path.
	Path string

	// Size in bytes.
	Size int64

	// Container format (mkv, mp4, avi, etc.).
	Container string

	// Resolution (e.g., "1920x1080").
	Resolution string

	// ResolutionLabel (e.g., "1080p", "4K").
	ResolutionLabel string

	// VideoCodec (e.g., "h264", "hevc").
	VideoCodec string

	// VideoProfile (e.g., "High", "Main 10").
	VideoProfile string

	// AudioCodec (e.g., "aac", "dts").
	AudioCodec string

	// BitrateKbps is the overall bitrate in kilobits per second.
	BitrateKbps int64

	// DurationSeconds is the duration in seconds.
	DurationSeconds float64

	// Framerate (e.g., 23.976, 24, 30).
	Framerate float64

	// DynamicRange (SDR, HDR10, Dolby Vision, etc.).
	DynamicRange string

	// ColorSpace (e.g., "bt709", "bt2020").
	ColorSpace string

	// AudioChannels is the number of audio channels.
	AudioChannels int

	// AudioLayout (e.g., "stereo", "5.1", "7.1").
	AudioLayout string

	// Languages is the list of audio languages (ISO 639-1).
	Languages []string

	// SubtitleLangs is the list of subtitle languages (ISO 639-1).
	SubtitleLangs []string
}

// MediaProber defines the interface for extracting technical information from media files.
type MediaProber interface {
	// Probe extracts technical information from a media file.
	Probe(filePath string) (*MediaFileInfo, error)
}

// ContentFileRepository defines the interface for managing content file records.
// This is a generic interface that content-specific repositories implement.
type ContentFileRepository[T any] interface {
	// GetFileByPath retrieves a file record by its path.
	GetFileByPath(ctx context.Context, path string) (*T, error)

	// CreateFile creates a new file record.
	CreateFile(ctx context.Context, file *T) (*T, error)

	// UpdateFile updates an existing file record.
	UpdateFile(ctx context.Context, file *T) (*T, error)

	// DeleteFile deletes a file record.
	DeleteFile(ctx context.Context, id uuid.UUID) error
}

// RefreshOptions configures content refresh behavior.
type RefreshOptions struct {
	// RefreshCredits indicates whether to refresh credits (cast/crew).
	RefreshCredits bool

	// RefreshGenres indicates whether to refresh genres.
	RefreshGenres bool

	// RefreshImages indicates whether to refresh images.
	RefreshImages bool

	// Languages specifies which languages to fetch metadata for.
	Languages []string
}

// DefaultRefreshOptions returns refresh options with all features enabled.
func DefaultRefreshOptions() RefreshOptions {
	return RefreshOptions{
		RefreshCredits: true,
		RefreshGenres:  true,
		RefreshImages:  true,
		Languages:      []string{"en-US"},
	}
}
