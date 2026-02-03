package movie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMovieFilename(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		expectedTitle string
		expectedYear  *int
	}{
		{
			name:          "Title (YEAR).ext",
			filename:      "The Matrix (1999).mkv",
			expectedTitle: "The Matrix",
			expectedYear:  intPtr(1999),
		},
		{
			name:          "Title.YEAR.ext",
			filename:      "The.Matrix.1999.mkv",
			expectedTitle: "The Matrix",
			expectedYear:  intPtr(1999),
		},
		{
			name:          "Title with quality markers",
			filename:      "The.Matrix.1999.1080p.BluRay.x264-GROUP.mkv",
			expectedTitle: "The Matrix",
			expectedYear:  intPtr(1999),
		},
		{
			name:          "Title with spaces and quality",
			filename:      "The Matrix 1999 2160p UHD BluRay x265.mkv",
			expectedTitle: "The Matrix",
			expectedYear:  intPtr(1999),
		},
		{
			name:          "Title without year",
			filename:      "The Matrix.mkv",
			expectedTitle: "The Matrix",
			expectedYear:  nil,
		},
		{
			name:          "Title with underscore separators",
			filename:      "The_Matrix_(1999).mkv",
			expectedTitle: "The Matrix",
			expectedYear:  intPtr(1999),
		},
		{
			name:          "Complex title with multiple parentheses",
			filename:      "The Matrix Reloaded (2003) (1080p).mkv",
			expectedTitle: "The Matrix Reloaded",
			expectedYear:  intPtr(2003),
		},
		{
			name:          "Title with REMUX",
			filename:      "The.Matrix.1999.UHD.BluRay.2160p.TrueHD.Atmos.7.1.DV.HEVC.REMUX-FraMeSToR.mkv",
			expectedTitle: "The Matrix",
			expectedYear:  intPtr(1999),
		},
		{
			name:          "Title with edition markers",
			filename:      "Blade Runner (1982) (Final Cut) (1080p BluRay x265 HEVC 10bit AAC 5.1).mkv",
			expectedTitle: "Blade Runner",
			expectedYear:  intPtr(1982),
		},
		{
			name:          "Invalid year (too old)",
			filename:      "Some Movie (1800).mkv",
			expectedTitle: "Some Movie",
			expectedYear:  nil, // Year outside valid range (1888-2100)
		},
		{
			name:          "Invalid year (too far future)",
			filename:      "Some Movie (2200).mkv",
			expectedTitle: "Some Movie",
			expectedYear:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, year := parseMovieFilename(tt.filename)
			
			assert.Equal(t, tt.expectedTitle, title, "Title mismatch")
			if tt.expectedYear == nil {
				assert.Nil(t, year, "Expected nil year")
			} else {
				assert.NotNil(t, year, "Expected non-nil year")
				assert.Equal(t, *tt.expectedYear, *year, "Year mismatch")
			}
		})
	}
}

func TestCleanTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Remove quality markers",
			input:    "The Matrix 1080p BluRay",
			expected: "The Matrix",
		},
		{
			name:     "Remove codec markers",
			input:    "The Matrix x264 HEVC",
			expected: "The Matrix",
		},
		{
			name:     "Remove group tags",
			input:    "The Matrix SPARKS",
			expected: "The Matrix",
		},
		{
			name:     "Remove audio markers",
			input:    "The Matrix DTS AC3",
			expected: "The Matrix",
		},
		{
			name:     "Multiple quality markers",
			input:    "The Matrix 2160p UHD BluRay x265 HDR10 DTS-HD MA 7.1",
			expected: "The Matrix",
		},
		{
			name:     "Clean title with dots",
			input:    "The.Matrix.1080p",
			expected: "The Matrix",
		},
		{
			name:     "Clean title with underscores",
			input:    "The_Matrix_BluRay",
			expected: "The Matrix",
		},
		{
			name:     "Remove REMUX",
			input:    "The Matrix REMUX",
			expected: "The Matrix",
		},
		{
			name:     "Remove WEB-DL markers",
			input:    "The Matrix WEB-DL WEBRip",
			expected: "The Matrix",
		},
		{
			name:     "Preserve special characters",
			input:    "The Matrix: Reloaded",
			expected: "The Matrix: Reloaded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanTitle(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsVideoFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "MKV file",
			filename: "movie.mkv",
			expected: true,
		},
		{
			name:     "MP4 file",
			filename: "movie.mp4",
			expected: true,
		},
		{
			name:     "AVI file",
			filename: "movie.avi",
			expected: true,
		},
		{
			name:     "M4V file",
			filename: "movie.m4v",
			expected: true,
		},
		{
			name:     "Text file",
			filename: "readme.txt",
			expected: false,
		},
		{
			name:     "NFO file",
			filename: "movie.nfo",
			expected: false,
		},
		{
			name:     "Subtitle file",
			filename: "movie.srt",
			expected: false,
		},
		{
			name:     "Uppercase extension",
			filename: "MOVIE.MKV",
			expected: true,
		},
		{
			name:     "Mixed case extension",
			filename: "Movie.Mp4",
			expected: true,
		},
		{
			name:     "No extension",
			filename: "movie",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVideoFile(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScanner_ParseMovieFilename_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		expectedTitle string
		expectedYear  *int
	}{
		{
			name:          "Empty filename",
			filename:      "",
			expectedTitle: "",
			expectedYear:  nil,
		},
		{
			name:          "Only extension",
			filename:      ".mkv",
			expectedTitle: "",
			expectedYear:  nil,
		},
		{
			name:          "Year at end without parentheses",
			filename:      "The Matrix 1999.mkv",
			expectedTitle: "The Matrix",
			expectedYear:  intPtr(1999),
		},
		{
			name:          "Multiple years (take first)",
			filename:      "The Matrix (1999) Reloaded (2003).mkv",
			expectedTitle: "The Matrix",
			expectedYear:  intPtr(1999),
		},
		{
			name:          "Year-like numbers that aren't years",
			filename:      "2001 A Space Odyssey (1968).mkv",
			expectedTitle: "2001 A Space Odyssey",
			expectedYear:  intPtr(1968),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, year := parseMovieFilename(tt.filename)
			
			assert.Equal(t, tt.expectedTitle, title, "Title mismatch")
			if tt.expectedYear == nil {
				assert.Nil(t, year, "Expected nil year")
			} else {
				assert.NotNil(t, year, "Expected non-nil year")
				assert.Equal(t, *tt.expectedYear, *year, "Year mismatch")
			}
		})
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
