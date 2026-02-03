package library

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
)

// VideoExtensions are supported video file extensions
var VideoExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
	".m4v":  true,
	".mpg":  true,
	".mpeg": true,
	".3gp":  true,
	".ts":   true,
	".m2ts": true,
}

// Scanner handles file system scanning for movie files
type Scanner struct {
	libraryPaths []string
}

// ScanResult represents a discovered movie file
type ScanResult struct {
	FilePath    string
	FileName    string
	ParsedTitle string
	ParsedYear  *int
	FileSize    int64
	IsVideo     bool
	Error       error
}

// NewScanner creates a new library scanner
func NewScanner(libraryPaths []string) *Scanner {
	return &Scanner{
		libraryPaths: libraryPaths,
	}
}

// Scan scans all library paths for movie files
func (s *Scanner) Scan(ctx context.Context) ([]ScanResult, error) {
	var results []ScanResult

	for _, path := range s.libraryPaths {
		pathResults, err := s.scanPath(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to scan path %s: %w", path, err)
		}
		results = append(results, pathResults...)
	}

	return results, nil
}

// scanPath scans a single library path
func (s *Scanner) scanPath(ctx context.Context, path string) ([]ScanResult, error) {
	var results []ScanResult

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if it's a video file
		ext := strings.ToLower(filepath.Ext(filePath))
		if !VideoExtensions[ext] {
			return nil
		}

		// Parse filename
		fileName := filepath.Base(filePath)
		title, year := parseMovieFilename(fileName)

		results = append(results, ScanResult{
			FilePath:    filePath,
			FileName:    fileName,
			ParsedTitle: title,
			ParsedYear:  year,
			FileSize:    info.Size(),
			IsVideo:     true,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// parseMovieFilename attempts to extract title and year from filename
// Examples:
// - "The Matrix (1999).mkv" -> "The Matrix", 1999
// - "Inception.2010.1080p.BluRay.mkv" -> "Inception", 2010
// - "The.Dark.Knight.2008.mkv" -> "The Dark Knight", 2008
func parseMovieFilename(filename string) (title string, year *int) {
	// Remove extension
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Try pattern: "Title (YEAR)"
	re1 := regexp.MustCompile(`^(.+?)\s*\((\d{4})\)`)
	if matches := re1.FindStringSubmatch(nameWithoutExt); len(matches) == 3 {
		title = cleanTitle(matches[1])
		if y := parseInt(matches[2]); y != nil && *y >= 1900 && *y <= 2100 {
			year = y
		}
		return
	}

	// Try pattern: "Title.YEAR" or "Title YEAR"
	re2 := regexp.MustCompile(`^(.+?)[\s\.](\d{4})`)
	if matches := re2.FindStringSubmatch(nameWithoutExt); len(matches) == 3 {
		title = cleanTitle(matches[1])
		if y := parseInt(matches[2]); y != nil && *y >= 1900 && *y <= 2100 {
			year = y
		}
		return
	}

	// No year found, use whole name as title
	title = cleanTitle(nameWithoutExt)
	return
}

// cleanTitle cleans up the title string
func cleanTitle(title string) string {
	// Replace dots/underscores with spaces
	title = strings.ReplaceAll(title, ".", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// Remove common quality markers
	qualityMarkers := []string{
		"1080p", "720p", "480p", "2160p", "4K",
		"BluRay", "BRRip", "WEBRip", "WEB-DL", "HDRip",
		"x264", "x265", "h264", "h265", "HEVC",
		"AAC", "DTS", "AC3", "DD5.1",
		"EXTENDED", "UNRATED", "REMASTERED",
	}

	titleLower := strings.ToLower(title)
	for _, marker := range qualityMarkers {
		markerLower := strings.ToLower(marker)
		if idx := strings.Index(titleLower, markerLower); idx != -1 {
			title = title[:idx]
			titleLower = titleLower[:idx]
		}
	}

	// Trim whitespace
	title = strings.TrimSpace(title)

	return title
}

// parseInt converts string to int pointer
func parseInt(s string) *int {
	var i int
	if _, err := fmt.Sscanf(s, "%d", &i); err != nil {
		return nil
	}
	return &i
}

// MovieFileInfo contains detailed information about a movie file
type MovieFileInfo struct {
	Path        string
	Size        int64
	Container   string
	Resolution  string
	VideoCodec  string
	AudioCodec  string
	Languages   []string
	BitrateKbps int32
}

// ExtractFileInfo extracts technical details from a video file
func ExtractFileInfo(filePath string) (*MovieFileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	container := strings.TrimPrefix(ext, ".")

	// Basic info without external tools (ffprobe would be needed for full details)
	return &MovieFileInfo{
		Path:      filePath,
		Size:      info.Size(),
		Container: container,
		// Resolution, codecs, etc. would require ffprobe integration
	}, nil
}

// CreateMovieFile creates a domain MovieFile from file info
func CreateMovieFile(movieID uuid.UUID, info *MovieFileInfo) *movie.MovieFile {
	return &movie.MovieFile{
		ID:         uuid.New(),
		MovieID:    movieID,
		FilePath:   info.Path,
		FileSize:   info.Size,
		FileName:   filepath.Base(info.Path),
		Container:  &info.Container,
		Resolution: parseOptionalString(info.Resolution),
		VideoCodec: parseOptionalString(info.VideoCodec),
		AudioCodec: parseOptionalString(info.AudioCodec),
	}
}

func parseOptionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseOptionalInt32(i int) *int32 {
	if i == 0 {
		return nil
	}
	val := int32(i)
	return &val
}
