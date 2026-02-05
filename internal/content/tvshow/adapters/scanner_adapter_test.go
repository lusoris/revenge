package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTVShowFileParser_Parse(t *testing.T) {
	parser := NewTVShowFileParser()

	tests := []struct {
		name            string
		filename        string
		expectedTitle   string
		expectedSeason  int
		expectedEpisode int
		expectedMeta    map[string]any
	}{
		// Standard SxxExx patterns
		{
			name:            "Standard SxxExx format",
			filename:        "Breaking Bad S01E01.mkv",
			expectedTitle:   "Breaking Bad",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "SxxExx with dots",
			filename:        "Breaking.Bad.S01E01.720p.BluRay.mkv",
			expectedTitle:   "Breaking Bad",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "SxxExx with dashes",
			filename:        "Breaking Bad - S01E01 - Pilot.mkv",
			expectedTitle:   "Breaking Bad",
			expectedSeason:  1,
			expectedEpisode: 1,
			expectedMeta:    map[string]any{"episode_title": "Pilot"},
		},
		{
			name:            "SxxExx lowercase",
			filename:        "breaking.bad.s01e01.mkv",
			expectedTitle:   "breaking bad",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "SxxExx double digit season and episode",
			filename:        "Game.of.Thrones.S08E06.mkv",
			expectedTitle:   "Game of Thrones",
			expectedSeason:  8,
			expectedEpisode: 6,
		},
		{
			name:            "SxxExx single digit no padding",
			filename:        "The Office S1E1.mkv",
			expectedTitle:   "The Office",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "Multi-episode SxxExxExx",
			filename:        "Breaking.Bad.S01E01E02.720p.mkv",
			expectedTitle:   "Breaking Bad",
			expectedSeason:  1,
			expectedEpisode: 1,
			expectedMeta:    map[string]any{"end_episode": 2},
		},

		// Series with year in title
		{
			name:            "Series with year in parentheses",
			filename:        "Doctor Who (2005) S01E01.mkv",
			expectedTitle:   "Doctor Who",
			expectedSeason:  1,
			expectedEpisode: 1,
			expectedMeta:    map[string]any{"series_year": 2005},
		},
		{
			name:            "Series with year in parentheses with dots",
			filename:        "Doctor.Who.(2005).S01E01.mkv",
			expectedTitle:   "Doctor Who",
			expectedSeason:  1,
			expectedEpisode: 1,
			expectedMeta:    map[string]any{"series_year": 2005},
		},

		// Quality and source markers
		{
			name:            "With quality markers",
			filename:        "Dark.S01E01.German.1080p.WEB.x264-GROUP.mkv",
			expectedTitle:   "Dark",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "With BluRay source",
			filename:        "Stranger.Things.S04E09.2160p.UHD.BluRay.REMUX.mkv",
			expectedTitle:   "Stranger Things",
			expectedSeason:  4,
			expectedEpisode: 9,
		},

		// Season Episode word format
		{
			name:            "Season Episode words",
			filename:        "Friends Season 1 Episode 1.mkv",
			expectedTitle:   "Friends",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "Season Episode words no space",
			filename:        "Friends Season01Episode01.mkv",
			expectedTitle:   "Friends",
			expectedSeason:  1,
			expectedEpisode: 1,
		},

		// Daily shows
		{
			name:           "Daily show with dots",
			filename:       "The Daily Show 2024.01.15.mkv",
			expectedTitle:  "The Daily Show",
			expectedSeason: 0,
			expectedMeta: map[string]any{
				"air_year":  2024,
				"air_month": 1,
				"air_day":   15,
				"is_daily":  true,
			},
		},
		{
			name:           "Daily show with dashes",
			filename:       "Last.Week.Tonight.2024-01-15.mkv",
			expectedTitle:  "Last Week Tonight",
			expectedSeason: 0,
			expectedMeta: map[string]any{
				"air_year":  2024,
				"air_month": 1,
				"air_day":   15,
				"is_daily":  true,
			},
		},

		// Dot-dash format (older style)
		{
			name:            "Dot format x.xx",
			filename:        "Friends 1.01 The Pilot.mkv",
			expectedTitle:   "Friends",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "Dash format x-xx",
			filename:        "Seinfeld 1-01.mkv",
			expectedTitle:   "Seinfeld",
			expectedSeason:  1,
			expectedEpisode: 1,
		},

		// Edge cases
		{
			name:            "Episode in triple digits",
			filename:        "One.Piece.S01E100.mkv",
			expectedTitle:   "One Piece",
			expectedSeason:  1,
			expectedEpisode: 100,
		},
		{
			name:            "Title with numbers",
			filename:        "24 S01E01.mkv",
			expectedTitle:   "24",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "Title with apostrophe",
			filename:        "Grey's Anatomy S01E01.mkv",
			expectedTitle:   "Grey's Anatomy",
			expectedSeason:  1,
			expectedEpisode: 1,
		},

		// No episode info (should extract title only)
		{
			name:            "No episode info",
			filename:        "Some.Random.Video.mkv",
			expectedTitle:   "Some Random Video",
			expectedSeason:  0,
			expectedEpisode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, metadata := parser.Parse(tt.filename)

			assert.Equal(t, tt.expectedTitle, title, "title mismatch")

			// Check season
			if tt.expectedSeason > 0 {
				assert.NotNil(t, metadata["season"], "season should be set")
				assert.Equal(t, tt.expectedSeason, metadata["season"], "season mismatch")
			}

			// Check episode
			if tt.expectedEpisode > 0 {
				assert.NotNil(t, metadata["episode"], "episode should be set")
				assert.Equal(t, tt.expectedEpisode, metadata["episode"], "episode mismatch")
			}

			// Check additional metadata
			if tt.expectedMeta != nil {
				for key, expectedValue := range tt.expectedMeta {
					actualValue, ok := metadata[key]
					if !ok {
						t.Errorf("expected metadata key %q not found", key)
						continue
					}
					assert.Equal(t, expectedValue, actualValue, "metadata[%s] mismatch", key)
				}
			}
		})
	}
}

func TestTVShowFileParser_ParseFromPath(t *testing.T) {
	parser := NewTVShowFileParser()

	tests := []struct {
		name            string
		filePath        string
		expectedTitle   string
		expectedSeason  int
		expectedEpisode int
		expectedMeta    map[string]any
	}{
		{
			name:            "Standard path with Season folder",
			filePath:        "/media/tv/Breaking Bad/Season 1/Breaking.Bad.S01E01.mkv",
			expectedTitle:   "Breaking Bad",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "Path with Season folder, minimal filename",
			filePath:        "/media/tv/Breaking Bad/Season 1/S01E01.mkv",
			expectedTitle:   "Breaking Bad",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "Path with series year",
			filePath:        "/media/tv/Doctor Who (2005)/Season 1/Doctor.Who.S01E01.mkv",
			expectedTitle:   "Doctor Who",
			expectedSeason:  1,
			expectedEpisode: 1,
			expectedMeta:    map[string]any{"series_year": 2005},
		},
		{
			name:            "Season folder with no space",
			filePath:        "/media/tv/The Office/Season01/S01E01.mkv",
			expectedTitle:   "The Office",
			expectedSeason:  1,
			expectedEpisode: 1,
		},
		{
			name:            "Filename has all info, path ignored",
			filePath:        "/random/path/Breaking.Bad.S02E05.mkv",
			expectedTitle:   "Breaking Bad",
			expectedSeason:  2,
			expectedEpisode: 5,
		},
		{
			name:            "Deep path structure",
			filePath:        "/data/media/TV Shows/Game of Thrones/Season 8/Game.of.Thrones.S08E06.mkv",
			expectedTitle:   "Game of Thrones",
			expectedSeason:  8,
			expectedEpisode: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, metadata := parser.ParseFromPath(tt.filePath)

			assert.Equal(t, tt.expectedTitle, title, "title mismatch")

			// Check season
			if tt.expectedSeason > 0 {
				assert.NotNil(t, metadata["season"], "season should be set")
				assert.Equal(t, tt.expectedSeason, metadata["season"], "season mismatch")
			}

			// Check episode
			if tt.expectedEpisode > 0 {
				assert.NotNil(t, metadata["episode"], "episode should be set")
				assert.Equal(t, tt.expectedEpisode, metadata["episode"], "episode mismatch")
			}

			// Check additional metadata
			if tt.expectedMeta != nil {
				for key, expectedValue := range tt.expectedMeta {
					actualValue, ok := metadata[key]
					if !ok {
						t.Errorf("expected metadata key %q not found", key)
						continue
					}
					assert.Equal(t, expectedValue, actualValue, "metadata[%s] mismatch", key)
				}
			}
		})
	}
}

func TestTVShowFileParser_GetExtensions(t *testing.T) {
	parser := NewTVShowFileParser()
	extensions := parser.GetExtensions()

	// Should have common video extensions
	assert.NotEmpty(t, extensions)
	assert.Contains(t, extensions, ".mkv")
	assert.Contains(t, extensions, ".mp4")
	assert.Contains(t, extensions, ".avi")
}

func TestTVShowFileParser_ContentType(t *testing.T) {
	parser := NewTVShowFileParser()
	assert.Equal(t, "tvshow", parser.ContentType())
}

func TestTVShowFileParser_InterfaceCompliance(t *testing.T) {
	// This test ensures compile-time interface compliance
	parser := NewTVShowFileParser()
	assert.NotNil(t, parser)
}
