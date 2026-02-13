package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple dots",
			input:    "The.Matrix",
			expected: "The Matrix",
		},
		{
			name:     "dots and underscores",
			input:    "The_Dark_Knight",
			expected: "The Dark Knight",
		},
		{
			name:     "quality markers",
			input:    "Inception.1080p.BluRay.x264",
			expected: "Inception",
		},
		{
			name:     "release group",
			input:    "Interstellar.SPARKS",
			expected: "Interstellar",
		},
		{
			name:     "multiple markers",
			input:    "Avatar.2160p.UHD.BluRay.REMUX.HEVC.DTS-HD.MA.5.1-FGT",
			expected: "Avatar",
		},
		{
			name:     "web release",
			input:    "The.Mandalorian.S01E01.WEBRip.x264-ION10",
			expected: "The Mandalorian S01E01",
		},
		{
			name:     "service markers",
			input:    "Stranger.Things.S04E01.NF.WEB-DL",
			expected: "Stranger Things S04E01",
		},
		{
			name:     "already clean",
			input:    "The Shawshank Redemption",
			expected: "The Shawshank Redemption",
		},
		{
			name:     "HDR content",
			input:    "Dune.2160p.UHD.BluRay.HDR10.HEVC.Atmos",
			expected: "Dune",
		},
		{
			name:     "year preserved - not a quality marker",
			input:    "Inception.2010",
			expected: "Inception 2010",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanTitle(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase",
			input:    "The Matrix",
			expected: "matrix",
		},
		{
			name:     "remove article 'the'",
			input:    "The Dark Knight",
			expected: "dark knight",
		},
		{
			name:     "remove article 'a'",
			input:    "A Beautiful Mind",
			expected: "beautiful mind",
		},
		{
			name:     "remove punctuation",
			input:    "Spider-Man: No Way Home",
			expected: "spiderman no way home",
		},
		{
			name:     "numbers preserved",
			input:    "2001: A Space Odyssey",
			expected: "2001 a space odyssey",
		},
		{
			name:     "special characters removed",
			input:    "Se7en",
			expected: "se7en",
		},
		{
			name:     "multiple spaces cleaned",
			input:    "The   Matrix   Reloaded",
			expected: "matrix reloaded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeTitle(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractYear(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *int
	}{
		{
			name:     "year in title",
			input:    "The Matrix 1999",
			expected: new(1999),
		},
		{
			name:     "year with dots",
			input:    "Inception.2010.1080p",
			expected: new(2010),
		},
		{
			name:     "recent year",
			input:    "Dune Part Two 2024",
			expected: new(2024),
		},
		{
			name:     "old year",
			input:    "Metropolis 1927",
			expected: new(1927),
		},
		{
			name:     "no year",
			input:    "Some Movie Title",
			expected: nil,
		},
		{
			name:     "invalid year range",
			input:    "Movie 2500",
			expected: nil,
		},
		{
			name:     "year in brackets",
			input:    "Movie (1999)",
			expected: new(1999),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractYear(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestParseYearFromBrackets(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *int
	}{
		{
			name:     "year in brackets",
			input:    "The Matrix (1999)",
			expected: new(1999),
		},
		{
			name:     "multiple brackets",
			input:    "Movie (2020) (Extended)",
			expected: new(2020),
		},
		{
			name:     "no brackets",
			input:    "Movie 2020",
			expected: nil,
		},
		{
			name:     "non-year in brackets",
			input:    "Movie (Extended)",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseYearFromBrackets(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestExtractResolution(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "1080p",
			input:    "Movie.2020.1080p.BluRay",
			expected: "1080p",
		},
		{
			name:     "720p",
			input:    "Movie.720p.WEBRip",
			expected: "720p",
		},
		{
			name:     "4K",
			input:    "Movie.4K.UHD",
			expected: "4K",
		},
		{
			name:     "2160p",
			input:    "Movie.2160p.HDR",
			expected: "2160p",
		},
		{
			name:     "no resolution",
			input:    "Movie.BluRay",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractResolution(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractSource(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "BluRay",
			input:    "Movie.1080p.BluRay.x264",
			expected: "BluRay",
		},
		{
			name:     "WEB-DL",
			input:    "Movie.WEB-DL.1080p",
			expected: "WEB-DL",
		},
		{
			name:     "WEBRip",
			input:    "Movie.WEBRip.720p",
			expected: "WEBRip",
		},
		{
			name:     "HDTV",
			input:    "Show.S01E01.HDTV.x264",
			expected: "HDTV",
		},
		{
			name:     "REMUX",
			input:    "Movie.REMUX.2160p",
			expected: "REMUX",
		},
		{
			name:     "BDRip",
			input:    "Movie.BDRip.1080p",
			expected: "BluRay",
		},
		{
			name:     "no source",
			input:    "Movie.1080p",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSource(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
