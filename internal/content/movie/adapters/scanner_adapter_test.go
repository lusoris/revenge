package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMovieFileParser_Parse(t *testing.T) {
	parser := NewMovieFileParser()

	tests := []struct {
		name         string
		filename     string
		wantTitle    string
		wantYear     *int
		wantMetadata bool
	}{
		{
			name:         "title with year in brackets",
			filename:     "The Matrix (1999).mkv",
			wantTitle:    "The Matrix",
			wantYear:     new(1999),
			wantMetadata: true,
		},
		{
			name:         "title with year and quality markers",
			filename:     "Inception.2010.1080p.BluRay.x264.mkv",
			wantTitle:    "Inception",
			wantYear:     new(2010),
			wantMetadata: true,
		},
		{
			name:         "title with dots",
			filename:     "The.Dark.Knight.2008.BluRay.mkv",
			wantTitle:    "The Dark Knight",
			wantYear:     new(2008),
			wantMetadata: true,
		},
		{
			name:         "title with underscores",
			filename:     "Avatar_2009_1080p.mp4",
			wantTitle:    "Avatar",
			wantYear:     new(2009),
			wantMetadata: true,
		},
		{
			name:         "title without year",
			filename:     "Some Movie.mkv",
			wantTitle:    "Some Movie",
			wantYear:     nil,
			wantMetadata: false,
		},
		{
			name:         "title with quality markers no year",
			filename:     "Movie.1080p.BluRay.mkv",
			wantTitle:    "Movie",
			wantYear:     nil,
			wantMetadata: false,
		},
		{
			name:         "complex release name",
			filename:     "Interstellar.2014.2160p.UHD.BluRay.REMUX.HDR.HEVC.Atmos-FGT.mkv",
			wantTitle:    "Interstellar",
			wantYear:     new(2014),
			wantMetadata: true,
		},
		{
			name:         "year in brackets with spaces",
			filename:     "The Godfather ( 1972 ).mkv",
			wantTitle:    "The Godfather",
			wantYear:     new(1972),
			wantMetadata: true,
		},
		{
			name:         "movie with multi-part title",
			filename:     "Lord.of.the.Rings.The.Return.of.the.King.2003.Extended.mkv",
			wantTitle:    "Lord of the Rings The Return of the King",
			wantYear:     new(2003),
			wantMetadata: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, metadata := parser.Parse(tt.filename)

			assert.Equal(t, tt.wantTitle, title)

			if tt.wantMetadata {
				require.NotNil(t, metadata)
				if tt.wantYear != nil {
					year, ok := metadata["year"]
					require.True(t, ok, "expected year in metadata")
					assert.Equal(t, *tt.wantYear, year)
				}
			} else {
				// Year should not be in metadata
				_, hasYear := metadata["year"]
				assert.False(t, hasYear)
			}
		})
	}
}

func TestMovieFileParser_GetExtensions(t *testing.T) {
	parser := NewMovieFileParser()
	extensions := parser.GetExtensions()

	// Should have video extensions
	assert.NotEmpty(t, extensions)

	// Check common extensions are present
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}

	assert.True(t, extMap[".mp4"], "should have .mp4")
	assert.True(t, extMap[".mkv"], "should have .mkv")
	assert.True(t, extMap[".avi"], "should have .avi")
}

func TestMovieFileParser_ContentType(t *testing.T) {
	parser := NewMovieFileParser()
	assert.Equal(t, "movie", parser.ContentType())
}
