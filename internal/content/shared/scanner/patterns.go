package scanner

import (
	"fmt"
	"regexp"
	"strings"
)

// QualityMarkers contains common quality markers found in release filenames.
// These are shared across movies and TV shows.
var QualityMarkers = []string{
	// Resolutions
	"2160p", "1080p", "720p", "480p", "360p",
	"4K", "UHD", "FHD", "HD",
	// Sources
	"BluRay", "Blu-Ray", "BRRip", "BDRip", "BDRIP",
	"WEBRip", "WEB-Rip", "WEB-DL", "WEBDL", "WEB",
	"HDRip", "HDRIP", "DVDRip", "DVDRIP",
	"HDTV", "PDTV", "SDTV", "DSR", "DSRip",
	"CAM", "HDCAM", "TS", "TELESYNC", "TC", "TELECINE",
	"SCR", "SCREENER", "DVDSCR", "BDSCR",
	// Video codecs
	"x264", "x265", "H264", "H.264", "H265", "H.265",
	"HEVC", "AVC", "VP9", "AV1",
	"XviD", "DivX",
	// Audio codecs
	"AAC", "AC3", "DTS", "DTS-HD", "DTSHD",
	"TrueHD", "Atmos", "DD5.1", "DD7.1",
	"FLAC", "EAC3", "E-AC3",
	// HDR
	"HDR", "HDR10", "HDR10+", "DV", "DoVi", "Dolby Vision",
	// Release types
	"EXTENDED", "UNRATED", "UNCUT", "DIRECTORS.CUT", "THEATRICAL",
	"REMASTERED", "REPACK", "PROPER", "REAL", "INTERNAL",
	"REMUX", "HYBRID",
	// 3D
	"3D", "SBS", "HSBS", "OU", "HOU",
}

// ReleaseGroups contains common release group tags found in filenames.
var ReleaseGroups = []string{
	// Scene groups
	"SPARKS", "GECKOS", "RARBG", "FGT", "EVO", "FLAME",
	"DRONES", "AMIABLE", "ROVERS", "VETO", "WARHD",
	"SiMPLE", "DEFLATE", "STRIFE", "BRMP", "CADAVER",
	// P2P groups
	"YTS", "YIFY", "ETRG", "MkvCage", "Tigole", "QxR",
	"PSA", "RARBG", "AMZN", "NTb", "NTG", "SiGMA",
	"STUTTERSHIT", "PSYCHD", "CMRG", "ION10",
	"D-Z0N3", "AZHD", "CtrlHD", "DON", "EbP", "EPSiLON",
	"FraMeSToR", "HiFi", "MainFrame", "NCmt", "SbR",
	// Web-DL groups
	"FLUX", "PECULATE", "MZABI", "TOMMY", "HONE",
	"CAKES", "SMURF", "DSNP", "HMAX", "AMZN", "NF", "ATVP",
	"DSNP+", "ROKU", "PCOK", "STAN", "PMTP", "iT",
}

// ServiceMarkers contains common words to strip from titles (typically service/platform markers).
var ServiceMarkers = []string{
	"AMZN", "NF", "NETFLIX", "ATVP", "DSNP", "HMAX", "HULU",
	"ROKU", "PCOK", "STAN", "iT", "PMTP",
}

// removeFromTitle removes all occurrences of markers from a title
// It performs case-insensitive matching with word boundaries to avoid false positives
// (e.g., "TS" should not match inside "DTS")
func removeFromTitle(title string, markers []string) string {
	result := title
	resultLower := strings.ToLower(title)

	for _, marker := range markers {
		markerLower := strings.ToLower(marker)
		idx := findWordBoundary(resultLower, markerLower)
		if idx != -1 {
			// Truncate at this point (markers typically appear at end of title)
			result = result[:idx]
			resultLower = resultLower[:idx]
		}
	}

	return strings.TrimSpace(result)
}

// findWordBoundary finds the index of a marker only if it appears as a whole word
// (i.e., not as a substring of another word)
func findWordBoundary(text, marker string) int {
	idx := strings.Index(text, marker)
	if idx == -1 {
		return -1
	}

	// Check if this is a word boundary match
	// Before: must be start of string or non-alphanumeric
	if idx > 0 {
		prevChar := text[idx-1]
		if isAlphanumeric(prevChar) {
			// Not a word boundary, search further
			remaining := text[idx+1:]
			nextIdx := findWordBoundary(remaining, marker)
			if nextIdx == -1 {
				return -1
			}
			return idx + 1 + nextIdx
		}
	}

	// After: must be end of string or non-alphanumeric
	endIdx := idx + len(marker)
	if endIdx < len(text) {
		nextChar := text[endIdx]
		if isAlphanumeric(nextChar) {
			// Not a word boundary, search further
			remaining := text[idx+1:]
			nextIdx := findWordBoundary(remaining, marker)
			if nextIdx == -1 {
				return -1
			}
			return idx + 1 + nextIdx
		}
	}

	return idx
}

// isAlphanumeric checks if a byte is a letter or digit
func isAlphanumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// CleanTitle removes quality markers, release groups, and normalizes a title string
func CleanTitle(title string) string {
	// Replace dots/underscores with spaces
	title = strings.ReplaceAll(title, ".", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// Remove quality markers (case-insensitive, truncate at match)
	title = removeFromTitle(title, QualityMarkers)

	// Remove release groups (case-insensitive, truncate at match)
	title = removeFromTitle(title, ReleaseGroups)

	// Remove service markers
	title = removeFromTitle(title, ServiceMarkers)

	// Clean up multiple spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	title = spaceRegex.ReplaceAllString(title, " ")

	return strings.TrimSpace(title)
}

// NormalizeTitle prepares a title for comparison by lowercasing,
// removing articles, and normalizing punctuation
func NormalizeTitle(title string) string {
	// Lowercase
	title = strings.ToLower(title)

	// Remove leading articles for matching purposes
	articles := []string{"the ", "a ", "an "}
	for _, article := range articles {
		if after, ok := strings.CutPrefix(title, article); ok {
			title = after
			break
		}
	}

	// Remove punctuation except alphanumeric and spaces
	var normalized strings.Builder
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			normalized.WriteRune(r)
		}
	}

	// Clean up multiple spaces
	result := normalized.String()
	spaceRegex := regexp.MustCompile(`\s+`)
	result = spaceRegex.ReplaceAllString(result, " ")

	return strings.TrimSpace(result)
}

// ExtractYear attempts to extract a year from text
// Returns nil if no valid year (1900-2100) is found
func ExtractYear(text string) *int {
	// Pattern: 4-digit year
	yearRegex := regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`)
	if matches := yearRegex.FindStringSubmatch(text); len(matches) == 2 {
		var year int
		if _, err := fmt.Sscanf(matches[1], "%d", &year); err == nil {
			if year >= 1900 && year <= 2100 {
				return &year
			}
		}
	}
	return nil
}

// ParseYearFromBrackets extracts year from "(YYYY)" pattern
func ParseYearFromBrackets(text string) *int {
	yearRegex := regexp.MustCompile(`\((\d{4})\)`)
	if matches := yearRegex.FindStringSubmatch(text); len(matches) == 2 {
		var year int
		if _, err := fmt.Sscanf(matches[1], "%d", &year); err == nil {
			if year >= 1900 && year <= 2100 {
				return &year
			}
		}
	}
	return nil
}

// ExtractResolution extracts resolution info from text
// Returns resolution string like "1080p" or empty string
func ExtractResolution(text string) string {
	resolutions := []string{"2160p", "1080p", "720p", "480p", "4K", "UHD"}
	textLower := strings.ToLower(text)
	for _, res := range resolutions {
		if strings.Contains(textLower, strings.ToLower(res)) {
			return res
		}
	}
	return ""
}

// ExtractSource extracts source info from text (BluRay, WEB-DL, etc)
func ExtractSource(text string) string {
	sources := map[string]string{
		"bluray":  "BluRay",
		"blu-ray": "BluRay",
		"bdrip":   "BluRay",
		"brrip":   "BluRay",
		"web-dl":  "WEB-DL",
		"webdl":   "WEB-DL",
		"webrip":  "WEBRip",
		"web-rip": "WEBRip",
		"hdtv":    "HDTV",
		"dvdrip":  "DVDRip",
		"remux":   "REMUX",
	}
	textLower := strings.ToLower(text)
	for pattern, source := range sources {
		if strings.Contains(textLower, pattern) {
			return source
		}
	}
	return ""
}
