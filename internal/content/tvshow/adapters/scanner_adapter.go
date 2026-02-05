// Package adapters provides TV show-specific implementations of shared interfaces.
package adapters

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/lusoris/revenge/internal/content/shared/scanner"
)

// TVShowFileParser implements scanner.FileParser for TV show files.
// It extracts series title, season number, and episode number from filenames like:
// - "Breaking Bad S01E01.mkv" -> Series: "Breaking Bad", Season: 1, Episode: 1
// - "Breaking.Bad.S01E01.720p.mkv" -> Series: "Breaking Bad", Season: 1, Episode: 1
// - "Breaking Bad - S01E01 - Pilot.mkv" -> Series: "Breaking Bad", Season: 1, Episode: 1, Title: "Pilot"
// - "Breaking Bad/Season 1/Breaking Bad - S01E01 - Pilot.mkv"
// - "Dark.S01E01.German.1080p.WEB.mkv" -> Series: "Dark", Season: 1, Episode: 1
type TVShowFileParser struct{}

// NewTVShowFileParser creates a new TV show file parser
func NewTVShowFileParser() *TVShowFileParser {
	return &TVShowFileParser{}
}

// Common regex patterns for TV show episode matching
var (
	// SxxExx pattern - most common format
	// Matches: S01E01, S1E1, S01E01E02 (multi-episode)
	sxxexxPattern = regexp.MustCompile(`(?i)[Ss](\d{1,2})[Ee](\d{1,3})(?:[Ee](\d{1,3}))?`)

	// Season x Episode x pattern
	// Matches: Season 1 Episode 1, Season01Episode01
	seasonEpisodePattern = regexp.MustCompile(`(?i)Season\s*(\d{1,2})\s*Episode\s*(\d{1,3})`)

	// x.xx or x-xx pattern (older format)
	// Matches: 1.01, 1-01, 01.01
	dotDashPattern = regexp.MustCompile(`(?:^|[\s\._-])(\d{1,2})[\.\-](\d{2})(?:[\s\._-]|$)`)

	// Episode title pattern - after SxxExx
	// Matches: "S01E01 - Pilot", "S01E01.Pilot"
	episodeTitlePattern = regexp.MustCompile(`(?i)[Ss]\d{1,2}[Ee]\d{1,3}(?:[Ee]\d{1,3})?\s*[\.\-\s]+(.+?)(?:[\.\-](?:720p|1080p|2160p|HDTV|WEB|BluRay)|$)`)

	// Year in series title pattern
	// Matches: "Doctor Who (2005)"
	seriesYearPattern = regexp.MustCompile(`^(.+?)\s*\((\d{4})\)`)

	// Daily show pattern (date-based episodes)
	// Matches: 2024.01.15, 2024-01-15
	dailyShowPattern = regexp.MustCompile(`(\d{4})[\.\-](\d{2})[\.\-](\d{2})`)
)

// Parse extracts series title, season, episode from a TV show filename
func (p *TVShowFileParser) Parse(filename string) (title string, metadata map[string]any) {
	metadata = make(map[string]any)

	// Remove extension
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Try SxxExx pattern first (most common)
	if matches := sxxexxPattern.FindStringSubmatch(nameWithoutExt); len(matches) >= 3 {
		// Everything before the SxxExx is the series title
		idx := sxxexxPattern.FindStringIndex(nameWithoutExt)
		if idx != nil && idx[0] > 0 {
			rawTitle := nameWithoutExt[:idx[0]]

			// Check for series year in raw title before cleaning (e.g., "Doctor Who (2005)")
			if yearMatches := seriesYearPattern.FindStringSubmatch(rawTitle); len(yearMatches) >= 3 {
				title = scanner.CleanTitle(yearMatches[1])
				if year, err := strconv.Atoi(yearMatches[2]); err == nil {
					metadata["series_year"] = year
				}
			} else {
				title = scanner.CleanTitle(rawTitle)
			}
			// Remove trailing separators
			title = strings.TrimRight(title, " .-_")
		}

		// Parse season and episode numbers
		if season, err := strconv.Atoi(matches[1]); err == nil {
			metadata["season"] = season
		}
		if episode, err := strconv.Atoi(matches[2]); err == nil {
			metadata["episode"] = episode
		}
		// Check for multi-episode (S01E01E02)
		if len(matches) >= 4 && matches[3] != "" {
			if endEpisode, err := strconv.Atoi(matches[3]); err == nil {
				metadata["end_episode"] = endEpisode
			}
		}

		// Try to extract episode title
		if epTitleMatches := episodeTitlePattern.FindStringSubmatch(nameWithoutExt); len(epTitleMatches) >= 2 {
			epTitle := scanner.CleanTitle(epTitleMatches[1])
			epTitle = strings.TrimRight(epTitle, " .-_")
			if epTitle != "" {
				metadata["episode_title"] = epTitle
			}
		}

		return title, metadata
	}

	// Try Season x Episode x pattern
	if matches := seasonEpisodePattern.FindStringSubmatch(nameWithoutExt); len(matches) >= 3 {
		idx := seasonEpisodePattern.FindStringIndex(nameWithoutExt)
		if idx != nil && idx[0] > 0 {
			title = scanner.CleanTitle(nameWithoutExt[:idx[0]])
		}

		if season, err := strconv.Atoi(matches[1]); err == nil {
			metadata["season"] = season
		}
		if episode, err := strconv.Atoi(matches[2]); err == nil {
			metadata["episode"] = episode
		}

		return title, metadata
	}

	// Try daily show pattern (date-based)
	if matches := dailyShowPattern.FindStringSubmatch(nameWithoutExt); len(matches) >= 4 {
		idx := dailyShowPattern.FindStringIndex(nameWithoutExt)
		if idx != nil && idx[0] > 0 {
			title = scanner.CleanTitle(nameWithoutExt[:idx[0]])
		}

		if year, err := strconv.Atoi(matches[1]); err == nil {
			metadata["air_year"] = year
		}
		if month, err := strconv.Atoi(matches[2]); err == nil {
			metadata["air_month"] = month
		}
		if day, err := strconv.Atoi(matches[3]); err == nil {
			metadata["air_day"] = day
		}
		metadata["is_daily"] = true

		return title, metadata
	}

	// Try x.xx pattern last (less reliable, more false positives)
	if matches := dotDashPattern.FindStringSubmatch(nameWithoutExt); len(matches) >= 3 {
		idx := dotDashPattern.FindStringIndex(nameWithoutExt)
		if idx != nil && idx[0] > 0 {
			title = scanner.CleanTitle(nameWithoutExt[:idx[0]])
		}

		if season, err := strconv.Atoi(matches[1]); err == nil {
			metadata["season"] = season
		}
		if episode, err := strconv.Atoi(matches[2]); err == nil {
			metadata["episode"] = episode
		}

		return title, metadata
	}

	// Fallback: clean the whole name as title, no episode info
	title = scanner.CleanTitle(nameWithoutExt)
	return title, metadata
}

// GetExtensions returns the video extensions supported for TV shows
func (p *TVShowFileParser) GetExtensions() []string {
	return scanner.ExtensionsToSlice(scanner.VideoExtensions)
}

// ContentType returns the content type identifier
func (p *TVShowFileParser) ContentType() string {
	return "tvshow"
}

// ParseFromPath attempts to extract series info from the full file path.
// This is useful when the series name is in the parent directory.
// Returns updated metadata if path parsing yields additional info.
func (p *TVShowFileParser) ParseFromPath(filePath string) (title string, metadata map[string]any) {
	// First, parse the filename
	filename := filepath.Base(filePath)
	title, metadata = p.Parse(filename)

	// Try to extract additional info from parent directories
	// even if filename provided title/season/episode (may need series_year)
	dir := filepath.Dir(filePath)
	parts := strings.Split(dir, string(filepath.Separator))

	// Look for "Season X" directory pattern
	seasonDirPattern := regexp.MustCompile(`(?i)^Season\s*(\d{1,2})$`)
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]

		// Check if this is a "Season X" directory
		if matches := seasonDirPattern.FindStringSubmatch(part); len(matches) >= 2 {
			if metadata["season"] == nil {
				if season, err := strconv.Atoi(matches[1]); err == nil {
					metadata["season"] = season
				}
			}
			// The parent of "Season X" is likely the series name
			if i > 0 {
				seriesDir := parts[i-1]
				// Check for series year in parent directory
				if yearMatches := seriesYearPattern.FindStringSubmatch(seriesDir); len(yearMatches) >= 3 {
					if title == "" {
						title = strings.TrimSpace(yearMatches[1])
					}
					if metadata["series_year"] == nil {
						if year, err := strconv.Atoi(yearMatches[2]); err == nil {
							metadata["series_year"] = year
						}
					}
				} else if title == "" {
					title = scanner.CleanTitle(seriesDir)
				}
			}
			break
		}

		// Check if this directory looks like a series name (not Season X)
		if !seasonDirPattern.MatchString(part) && part != "" {
			// Check for series year
			if yearMatches := seriesYearPattern.FindStringSubmatch(part); len(yearMatches) >= 3 {
				if title == "" {
					title = strings.TrimSpace(yearMatches[1])
				}
				if metadata["series_year"] == nil {
					if year, err := strconv.Atoi(yearMatches[2]); err == nil {
						metadata["series_year"] = year
					}
				}
			} else if title == "" {
				potentialTitle := scanner.CleanTitle(part)
				if potentialTitle != "" && len(potentialTitle) > 2 {
					title = potentialTitle
				}
			}
			// Don't break - continue looking for Season X
		}
	}

	return title, metadata
}

// Verify interface compliance at compile time
var _ scanner.FileParser = (*TVShowFileParser)(nil)
