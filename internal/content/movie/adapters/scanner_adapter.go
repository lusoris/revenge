// Package adapters provides movie-specific implementations of shared interfaces.
package adapters

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lusoris/revenge/internal/content/shared/scanner"
)

// MovieFileParser implements scanner.FileParser for movie files.
// It extracts movie title and year from filenames like:
// - "The Matrix (1999).mkv" -> Title: "The Matrix", Year: 1999
// - "Inception.2010.1080p.BluRay.mkv" -> Title: "Inception", Year: 2010
type MovieFileParser struct{}

// NewMovieFileParser creates a new movie file parser
func NewMovieFileParser() *MovieFileParser {
	return &MovieFileParser{}
}

// Parse extracts title and year from a movie filename
func (p *MovieFileParser) Parse(filename string) (title string, metadata map[string]any) {
	metadata = make(map[string]any)

	// Remove extension
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Try pattern: "Title (YEAR)" or "Title ( YEAR )" (with optional spaces inside brackets)
	re1 := regexp.MustCompile(`^(.+?)\s*\(\s*(\d{4})\s*\)`)
	if matches := re1.FindStringSubmatch(nameWithoutExt); len(matches) == 3 {
		title = scanner.CleanTitle(matches[1])
		if year := scanner.ExtractYear(matches[2]); year != nil {
			metadata["year"] = *year
		}
		return title, metadata
	}

	// Try pattern: "Title.YEAR", "Title YEAR", or "Title_YEAR" (before quality markers)
	// Match separators: space, dot, or underscore
	re2 := regexp.MustCompile(`^(.+?)[\s\._](\d{4})`)
	if matches := re2.FindStringSubmatch(nameWithoutExt); len(matches) == 3 {
		title = scanner.CleanTitle(matches[1])
		if year := scanner.ExtractYear(matches[2]); year != nil {
			metadata["year"] = *year
		}
		return title, metadata
	}

	// No year found, clean and use whole name as title
	title = scanner.CleanTitle(nameWithoutExt)
	return title, metadata
}

// GetExtensions returns the video extensions supported for movies
func (p *MovieFileParser) GetExtensions() []string {
	return scanner.ExtensionsToSlice(scanner.VideoExtensions)
}

// ContentType returns the content type identifier
func (p *MovieFileParser) ContentType() string {
	return "movie"
}

// Verify interface compliance at compile time
var _ scanner.FileParser = (*MovieFileParser)(nil)
