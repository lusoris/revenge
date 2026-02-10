package omdb

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantY   int
		wantM   time.Month
		wantD   int
	}{
		{"empty", "", true, 0, 0, 0},
		{"N/A", "N/A", true, 0, 0, 0},
		{"valid", "15 Jun 2024", false, 2024, 6, 15},
		{"another", "01 Jan 2006", false, 2006, 1, 1},
		{"invalid format", "2024-06-15", true, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDate(tt.input)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.wantY, result.Year())
				assert.Equal(t, tt.wantM, result.Month())
				assert.Equal(t, tt.wantD, result.Day())
			}
		})
	}
}

func TestParseRuntime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		want    int32
	}{
		{"empty", "", true, 0},
		{"N/A", "N/A", true, 0},
		{"valid", "142 min", false, 142},
		{"short", "5 min", false, 5},
		{"invalid", "abc min", true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseRuntime(tt.input)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.want, *result)
			}
		})
	}
}

func TestParseGenres(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []metadata.Genre
	}{
		{"empty", "", nil},
		{"N/A", "N/A", nil},
		{"single", "Action", []metadata.Genre{{Name: "Action"}}},
		{"multiple", "Action, Drama, Thriller", []metadata.Genre{
			{Name: "Action"}, {Name: "Drama"}, {Name: "Thriller"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseGenres(tt.input))
		})
	}
}

func TestNormalizeScore(t *testing.T) {
	tests := []struct {
		name   string
		source string
		value  string
		want   float64
	}{
		{"IMDb", "Internet Movie Database", "8.8/10", 88.0},
		{"RT", "Rotten Tomatoes", "96%", 96.0},
		{"Metacritic", "Metacritic", "90/100", 90.0},
		{"unknown", "Unknown", "5/5", 0},
		{"IMDb invalid", "Internet Movie Database", "bad", 0},
		{"IMDb no slash", "Internet Movie Database", "8.8", 0},
		{"RT invalid", "Rotten Tomatoes", "bad%", 0},
		{"Metacritic invalid", "Metacritic", "bad/100", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, normalizeScore(tt.source, tt.value), 0.01)
		})
	}
}

func TestMapExternalRatings(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapExternalRatings(nil))
	})

	t.Run("empty ratings", func(t *testing.T) {
		result := mapExternalRatings(&Response{})
		assert.Nil(t, result)
	})

	t.Run("with ratings", func(t *testing.T) {
		resp := &Response{
			Ratings: []Rating{
				{Source: "Internet Movie Database", Value: "8.8/10"},
				{Source: "Rotten Tomatoes", Value: "96%"},
			},
		}
		result := mapExternalRatings(resp)
		require.Len(t, result, 2)
		assert.Equal(t, "Internet Movie Database", result[0].Source)
		assert.Equal(t, "8.8/10", result[0].Value)
		assert.InDelta(t, 88.0, result[0].Score, 0.01)
		assert.Equal(t, "Rotten Tomatoes", result[1].Source)
		assert.InDelta(t, 96.0, result[1].Score, 0.01)
	})
}

func TestMapMovieMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapMovieMetadata(nil))
	})

	t.Run("wrong type", func(t *testing.T) {
		assert.Nil(t, mapMovieMetadata(&Response{Type: "series"}))
	})

	t.Run("full movie", func(t *testing.T) {
		resp := &Response{
			Title:      "Inception",
			IMDbID:     "tt1375666",
			Type:       "movie",
			Plot:       "A thief who steals corporate secrets.",
			Runtime:    "148 min",
			Released:   "16 Jul 2010",
			Genre:      "Action, Sci-Fi, Thriller",
			Poster:     "https://example.com/poster.jpg",
			IMDbRating: "8.8",
			IMDbVotes:  "2,500,000",
			Ratings:    []Rating{{Source: "Internet Movie Database", Value: "8.8/10"}},
		}
		result := mapMovieMetadata(resp)
		require.NotNil(t, result)
		assert.Equal(t, "tt1375666", result.ProviderID)
		assert.Equal(t, metadata.ProviderOMDb, result.Provider)
		assert.Equal(t, "Inception", result.Title)
		require.NotNil(t, result.IMDbID)
		assert.Equal(t, "tt1375666", *result.IMDbID)
		require.NotNil(t, result.Overview)
		assert.Equal(t, "A thief who steals corporate secrets.", *result.Overview)
		require.NotNil(t, result.Runtime)
		assert.Equal(t, int32(148), *result.Runtime)
		require.NotNil(t, result.ReleaseDate)
		assert.Equal(t, 2010, result.ReleaseDate.Year())
		require.Len(t, result.Genres, 3)
		assert.Equal(t, "Action", result.Genres[0].Name)
		require.NotNil(t, result.PosterPath)
		assert.Equal(t, "https://example.com/poster.jpg", *result.PosterPath)
		assert.InDelta(t, 8.8, result.VoteAverage, 0.01)
		assert.Equal(t, 2500000, result.VoteCount)
		require.Len(t, result.ExternalRatings, 1)
	})

	t.Run("N/A fields", func(t *testing.T) {
		resp := &Response{
			Title:      "Test",
			IMDbID:     "tt0000001",
			Type:       "movie",
			Plot:       "N/A",
			Runtime:    "N/A",
			Released:   "N/A",
			Genre:      "N/A",
			Poster:     "N/A",
			IMDbRating: "N/A",
			IMDbVotes:  "N/A",
		}
		result := mapMovieMetadata(resp)
		require.NotNil(t, result)
		assert.Nil(t, result.Overview)
		assert.Nil(t, result.Runtime)
		assert.Nil(t, result.ReleaseDate)
		assert.Nil(t, result.Genres)
		assert.Nil(t, result.PosterPath)
		assert.InDelta(t, 0.0, result.VoteAverage, 0.01)
	})
}

func TestMapTVShowMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapTVShowMetadata(nil))
	})

	t.Run("wrong type", func(t *testing.T) {
		assert.Nil(t, mapTVShowMetadata(&Response{Type: "movie"}))
	})

	t.Run("full series", func(t *testing.T) {
		resp := &Response{
			Title:        "Breaking Bad",
			IMDbID:       "tt0903747",
			Type:         "series",
			Plot:         "A chemistry teacher turns to cooking meth.",
			Genre:        "Crime, Drama, Thriller",
			Poster:       "https://example.com/bb.jpg",
			Released:     "20 Jan 2008",
			TotalSeasons: "5",
			IMDbRating:   "9.5",
			IMDbVotes:    "2,000,000",
		}
		result := mapTVShowMetadata(resp)
		require.NotNil(t, result)
		assert.Equal(t, "tt0903747", result.ProviderID)
		assert.Equal(t, metadata.ProviderOMDb, result.Provider)
		assert.Equal(t, "Breaking Bad", result.Name)
		require.NotNil(t, result.IMDbID)
		assert.Equal(t, 5, result.NumberOfSeasons)
		assert.InDelta(t, 9.5, result.VoteAverage, 0.01)
		assert.Equal(t, 2000000, result.VoteCount)
	})
}

func TestMapMovieSearchResults(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapMovieSearchResults(nil))
	})

	t.Run("mixed results", func(t *testing.T) {
		resp := &SearchResponse{
			Search: []SearchResult{
				{Title: "Inception", Year: "2010", IMDbID: "tt1375666", Type: "movie", Poster: "https://example.com/p1.jpg"},
				{Title: "Breaking Bad", Year: "2008", IMDbID: "tt0903747", Type: "series"},
				{Title: "Interstellar", Year: "2014", IMDbID: "tt0816692", Type: "movie", Poster: "N/A"},
			},
		}
		result := mapMovieSearchResults(resp)
		require.Len(t, result, 2) // only movies
		assert.Equal(t, "Inception", result[0].Title)
		assert.Equal(t, metadata.ProviderOMDb, result[0].Provider)
		require.NotNil(t, result[0].Year)
		assert.Equal(t, 2010, *result[0].Year)
		require.NotNil(t, result[0].PosterPath)
		assert.Equal(t, "https://example.com/p1.jpg", *result[0].PosterPath)
		// N/A poster should be nil
		assert.Nil(t, result[1].PosterPath)
	})

	t.Run("invalid year", func(t *testing.T) {
		resp := &SearchResponse{
			Search: []SearchResult{
				{Title: "Test", Year: "bad", IMDbID: "tt0000001", Type: "movie"},
			},
		}
		result := mapMovieSearchResults(resp)
		require.Len(t, result, 1)
		assert.Nil(t, result[0].Year)
	})
}

func TestMapTVShowSearchResults(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapTVShowSearchResults(nil))
	})

	t.Run("filters to series only", func(t *testing.T) {
		resp := &SearchResponse{
			Search: []SearchResult{
				{Title: "Breaking Bad", Year: "2008", IMDbID: "tt0903747", Type: "series", Poster: "https://example.com/bb.jpg"},
				{Title: "Inception", Year: "2010", IMDbID: "tt1375666", Type: "movie"},
			},
		}
		result := mapTVShowSearchResults(resp)
		require.Len(t, result, 1)
		assert.Equal(t, "Breaking Bad", result[0].Name)
		assert.Equal(t, metadata.ProviderOMDb, result[0].Provider)
		require.NotNil(t, result[0].Year)
		assert.Equal(t, 2008, *result[0].Year)
		require.NotNil(t, result[0].PosterPath)
	})
}
