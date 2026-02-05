package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTVShowTMDbClient(t *testing.T) {
	config := TVShowTMDbClientConfig{
		APIKey:   "test-api-key",
		ProxyURL: "",
	}

	client := NewTVShowTMDbClient(config)
	assert.NotNil(t, client)
}

func TestNewTVShowImageDownloader(t *testing.T) {
	config := TVShowTMDbClientConfig{
		APIKey: "test-api-key",
	}
	client := NewTVShowTMDbClient(config)
	downloader := NewTVShowImageDownloader(client)
	assert.NotNil(t, downloader)
}

func TestNewTVShowImageURLBuilder(t *testing.T) {
	builder := NewTVShowImageURLBuilder()
	assert.NotNil(t, builder)
}

func TestTVShowTMDbEndpoints(t *testing.T) {
	// Verify all endpoints are defined
	assert.NotEmpty(t, TVShowTMDbEndpoints.SearchTV)
	assert.NotEmpty(t, TVShowTMDbEndpoints.TVDetails)
	assert.NotEmpty(t, TVShowTMDbEndpoints.TVCredits)
	assert.NotEmpty(t, TVShowTMDbEndpoints.TVImages)
	assert.NotEmpty(t, TVShowTMDbEndpoints.TVContentRatings)
	assert.NotEmpty(t, TVShowTMDbEndpoints.TVTranslations)
	assert.NotEmpty(t, TVShowTMDbEndpoints.TVExternalIDs)

	assert.NotEmpty(t, TVShowTMDbEndpoints.SeasonDetails)
	assert.NotEmpty(t, TVShowTMDbEndpoints.SeasonCredits)
	assert.NotEmpty(t, TVShowTMDbEndpoints.SeasonImages)

	assert.NotEmpty(t, TVShowTMDbEndpoints.EpisodeDetails)
	assert.NotEmpty(t, TVShowTMDbEndpoints.EpisodeCredits)
	assert.NotEmpty(t, TVShowTMDbEndpoints.EpisodeImages)
}

func TestTVShowTMDbAppendToResponse(t *testing.T) {
	// Verify all append options are defined
	assert.NotEmpty(t, TVShowTMDbAppendToResponse.SeriesFull)
	assert.NotEmpty(t, TVShowTMDbAppendToResponse.SeriesBasic)
	assert.NotEmpty(t, TVShowTMDbAppendToResponse.SeasonFull)
	assert.NotEmpty(t, TVShowTMDbAppendToResponse.EpisodeFull)

	// SeriesFull should include key components
	assert.Contains(t, TVShowTMDbAppendToResponse.SeriesFull, "credits")
	assert.Contains(t, TVShowTMDbAppendToResponse.SeriesFull, "images")
	assert.Contains(t, TVShowTMDbAppendToResponse.SeriesFull, "content_ratings")
}

func TestGetTVShowGenreName(t *testing.T) {
	tests := []struct {
		genreID  int
		expected string
	}{
		{10759, "Action & Adventure"},
		{16, "Animation"},
		{35, "Comedy"},
		{80, "Crime"},
		{99, "Documentary"},
		{18, "Drama"},
		{10751, "Family"},
		{10762, "Kids"},
		{9648, "Mystery"},
		{10763, "News"},
		{10764, "Reality"},
		{10765, "Sci-Fi & Fantasy"},
		{10766, "Soap"},
		{10767, "Talk"},
		{10768, "War & Politics"},
		{37, "Western"},
		{99999, ""}, // Unknown ID
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetTVShowGenreName(tt.genreID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTVShowStatus(t *testing.T) {
	tests := []struct {
		tmdbStatus string
		expected   string
	}{
		{"Returning Series", "Returning Series"},
		{"Ended", "Ended"},
		{"Canceled", "Canceled"},
		{"In Production", "In Production"},
		{"Planned", "Planned"},
		{"Pilot", "Pilot"},
		{"Unknown Status", "Unknown Status"}, // Passthrough for unknown
	}

	for _, tt := range tests {
		t.Run(tt.tmdbStatus, func(t *testing.T) {
			result := GetTVShowStatus(tt.tmdbStatus)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTVShowType(t *testing.T) {
	tests := []struct {
		tmdbType string
		expected string
	}{
		{"Scripted", "Scripted"},
		{"Reality", "Reality"},
		{"Documentary", "Documentary"},
		{"Miniseries", "Miniseries"},
		{"Talk Show", "Talk Show"},
		{"News", "News"},
		{"Video", "Video"},
		{"Unknown Type", "Unknown Type"}, // Passthrough for unknown
	}

	for _, tt := range tests {
		t.Run(tt.tmdbType, func(t *testing.T) {
			result := GetTVShowType(tt.tmdbType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTVShowGenreMapCompleteness(t *testing.T) {
	// Ensure we have the common TV genres from TMDb
	// Note: TMDb has different genre IDs for TV vs movies
	expectedGenres := map[int]string{
		10759: "Action & Adventure",
		16:    "Animation",
		35:    "Comedy",
		18:    "Drama",
		10765: "Sci-Fi & Fantasy",
	}

	for id, name := range expectedGenres {
		result := GetTVShowGenreName(id)
		assert.Equal(t, name, result, "Genre ID %d should map to %s", id, name)
	}
}
