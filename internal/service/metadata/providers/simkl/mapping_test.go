package simkl

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantNil  bool
		wantDate time.Time
	}{
		{"empty string", "", true, time.Time{}},
		{"valid date", "2024-06-15", false, time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)},
		{"invalid format", "15/06/2024", true, time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDate(tt.input)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.True(t, tt.wantDate.Equal(*result))
			}
		})
	}
}

func TestMapMovieStatus(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"released", "Released"},
		{"upcoming", "Upcoming"},
		{"rumored", "Rumored"},
		{"planned", "Planned"},
		{"in production", "In Production"},
		{"post production", "Post Production"},
		{"cancelled", "Canceled"},
		{"Something Else", "Something Else"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapMovieStatus(tt.input))
		})
	}
}

func TestMapShowStatus(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"airing", "Returning Series"},
		{"returning series", "Returning Series"},
		{"ongoing", "Returning Series"},
		{"ended", "Ended"},
		{"cancelled", "Canceled"},
		{"canceled", "Canceled"},
		{"tba", "In Production"},
		{"upcoming", "In Production"},
		{"planned", "In Production"},
		{"in production", "In Production"},
		{"Something Else", "Something Else"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapShowStatus(tt.input))
		})
	}
}

func TestGenreNameToID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"Action", 28},
		{"action", 28},
		{"Comedy", 35},
		{"comedy", 35},
		{"Drama", 18},
		{"Unknown Genre", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, genreNameToID(tt.name))
		})
	}
}

func TestMapSearchResultToMovieSearchResult(t *testing.T) {
	r := SearchResult{
		Title:  "Inception",
		Year:   2010,
		Poster: "abc123",
		IDs:    IDs{Simkl: 42},
	}
	result := mapSearchResultToMovieSearchResult(r)
	assert.Equal(t, "42", result.ProviderID)
	assert.Equal(t, metadata.ProviderSimkl, result.Provider)
	assert.Equal(t, "Inception", result.Title)
	require.NotNil(t, result.Year)
	assert.Equal(t, 2010, *result.Year)
	require.NotNil(t, result.PosterPath)
	assert.Contains(t, *result.PosterPath, "abc123")
	assert.Contains(t, *result.PosterPath, "_ca")

	// No year
	r2 := SearchResult{Title: "Test", IDs: IDs{Simkl: 1}}
	result2 := mapSearchResultToMovieSearchResult(r2)
	assert.Nil(t, result2.Year)
	assert.Nil(t, result2.PosterPath)
}

func TestMapSearchResultToTVShowSearchResult(t *testing.T) {
	r := SearchResult{
		Title:  "Breaking Bad",
		Year:   2008,
		Poster: "poster123",
		IDs:    IDs{Simkl: 99},
	}
	result := mapSearchResultToTVShowSearchResult(r)
	assert.Equal(t, "99", result.ProviderID)
	assert.Equal(t, metadata.ProviderSimkl, result.Provider)
	assert.Equal(t, "Breaking Bad", result.Name)
	require.NotNil(t, result.Year)
	assert.Equal(t, 2008, *result.Year)
	require.NotNil(t, result.PosterPath)
}

func TestMapMovieToMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapMovieToMetadata(nil))
	})

	t.Run("full movie", func(t *testing.T) {
		m := &Movie{
			Title:       "Inception",
			IDs:         IDs{Simkl: 42, IMDb: "tt1375666", TMDb: 27205, TVDb: 12345},
			Overview:    "A thief who steals corporate secrets.",
			Trailer:     "https://youtube.com/watch?v=abc",
			Runtime:     148,
			Status:      "released",
			ReleaseDate: "2010-07-16",
			Poster:      "poster123",
			Fanart:      "fanart123",
			Country:     "us",
			Genres:      []string{"Action", "Drama"},
			Ratings: &Ratings{
				Simkl: &RatingInfo{Rating: 8.5, Votes: 1000},
				IMDb:  &RatingInfo{Rating: 8.8},
			},
		}
		result := mapMovieToMetadata(m)
		require.NotNil(t, result)
		assert.Equal(t, "42", result.ProviderID)
		assert.Equal(t, metadata.ProviderSimkl, result.Provider)
		assert.Equal(t, "Inception", result.Title)
		assert.Equal(t, "Released", result.Status)
		require.NotNil(t, result.Overview)
		assert.Equal(t, "A thief who steals corporate secrets.", *result.Overview)
		require.NotNil(t, result.TrailerURL)
		require.NotNil(t, result.Runtime)
		assert.Equal(t, int32(148), *result.Runtime)
		require.NotNil(t, result.IMDbID)
		assert.Equal(t, "tt1375666", *result.IMDbID)
		require.NotNil(t, result.TMDbID)
		assert.Equal(t, int32(27205), *result.TMDbID)
		require.NotNil(t, result.TVDbID)
		require.NotNil(t, result.ReleaseDate)
		assert.InDelta(t, 8.5, result.VoteAverage, 0.01)
		assert.Equal(t, 1000, result.VoteCount)
		require.NotNil(t, result.PosterPath)
		assert.Contains(t, *result.PosterPath, "poster123")
		require.NotNil(t, result.BackdropPath)
		assert.Contains(t, *result.BackdropPath, "fanart123")
		require.Len(t, result.Genres, 2)
		require.Len(t, result.ProductionCountries, 1)
		assert.Equal(t, "US", result.ProductionCountries[0].ISOCode)
		// External ratings: Simkl + IMDb
		require.Len(t, result.ExternalRatings, 2)
	})

	t.Run("minimal movie", func(t *testing.T) {
		m := &Movie{Title: "Test", IDs: IDs{Simkl: 1}}
		result := mapMovieToMetadata(m)
		require.NotNil(t, result)
		assert.Nil(t, result.Overview)
		assert.Nil(t, result.Runtime)
		assert.Nil(t, result.PosterPath)
	})
}

func TestMapShowToMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapShowToMetadata(nil))
	})

	t.Run("full show", func(t *testing.T) {
		s := &Show{
			Title:         "Breaking Bad",
			ENTitle:       "Breaking Bad EN",
			IDs:           IDs{Simkl: 99, IMDb: "tt0903747", TMDb: 1396, TVDb: 81189},
			Overview:      "A chemistry teacher.",
			Trailer:       "https://youtube.com/watch?v=xyz",
			Runtime:       47,
			Status:        "ended",
			TotalEpisodes: 62,
			Network:       "AMC",
			Country:       "us",
			Genres:        []string{"Drama", "Crime"},
			Ratings: &Ratings{
				Simkl: &RatingInfo{Rating: 9.0, Votes: 5000},
			},
			Poster: "bbposter",
			Fanart: "bbfanart",
		}
		result := mapShowToMetadata(s)
		require.NotNil(t, result)
		assert.Equal(t, "99", result.ProviderID)
		assert.Equal(t, "Breaking Bad EN", result.Name)
		assert.Equal(t, "Breaking Bad", result.OriginalName)
		assert.Equal(t, "Ended", result.Status)
		assert.False(t, result.InProduction)
		assert.Equal(t, 62, result.NumberOfEpisodes)
		assert.Equal(t, []int{47}, result.EpisodeRuntime)
		require.Len(t, result.Networks, 1)
		assert.Equal(t, "AMC", result.Networks[0].Name)
		assert.InDelta(t, 9.0, result.VoteAverage, 0.01)
		assert.Equal(t, 5000, result.VoteCount)
		assert.Equal(t, []string{"US"}, result.OriginCountries)
		require.NotNil(t, result.PosterPath)
		require.NotNil(t, result.BackdropPath)
	})

	t.Run("airing show is in production", func(t *testing.T) {
		s := &Show{Title: "Test", IDs: IDs{Simkl: 1}, Status: "airing"}
		result := mapShowToMetadata(s)
		require.NotNil(t, result)
		assert.True(t, result.InProduction)
	})
}

func TestMapEpisodesToSummaries(t *testing.T) {
	now := time.Now()
	episodes := []Episode{
		{Title: "Pilot", Season: 1, Episode: 1, Img: "ep1img", Date: &now, IDs: EpisodeIDs{Simkl: 100}},
		{Title: "Cat's in the Bag", Season: 1, Episode: 2, IDs: EpisodeIDs{Simkl: 101}},
		{Title: "S2E1", Season: 2, Episode: 1, IDs: EpisodeIDs{Simkl: 200}},
	}

	result := mapEpisodesToSummaries(episodes, 1)
	require.Len(t, result, 2)
	assert.Equal(t, "100", result[0].ProviderID)
	assert.Equal(t, "Pilot", result[0].Name)
	assert.Equal(t, 1, result[0].EpisodeNumber)
	require.NotNil(t, result[0].StillPath)
	assert.Contains(t, *result[0].StillPath, "ep1img")

	// Season 2 filter
	result2 := mapEpisodesToSummaries(episodes, 2)
	require.Len(t, result2, 1)

	// Empty
	result3 := mapEpisodesToSummaries(episodes, 99)
	assert.Empty(t, result3)
}

func TestMapExternalIDs(t *testing.T) {
	ids := IDs{IMDb: "tt1375666", TMDb: 27205, TVDb: 12345}
	result := mapExternalIDs(ids)
	require.NotNil(t, result)
	require.NotNil(t, result.IMDbID)
	assert.Equal(t, "tt1375666", *result.IMDbID)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(27205), *result.TMDbID)
	require.NotNil(t, result.TVDbID)

	// Empty IDs
	result2 := mapExternalIDs(IDs{})
	require.NotNil(t, result2)
	assert.Nil(t, result2.IMDbID)
	assert.Nil(t, result2.TMDbID)
	assert.Nil(t, result2.TVDbID)
}
