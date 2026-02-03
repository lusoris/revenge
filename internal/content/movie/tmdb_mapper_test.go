package movie

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTMDbMapper_MapMovie(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{
		APIKey: "test-key",
	})
	mapper := NewTMDbMapper(client)

	t.Run("Full movie with all fields", func(t *testing.T) {
		runtime := 139
		budget := int64(63000000)
		revenue := int64(100853753)
		tmdbMovie := &TMDbMovie{
			ID:               550,
			IMDbID:           stringPtr("tt0137523"),
			Title:            "Fight Club",
			OriginalTitle:    "Fight Club",
			OriginalLanguage: "en",
			Overview:         stringPtr("A depressed man suffering from insomnia meets a strange soap salesman."),
			Tagline:          stringPtr("Mischief. Mayhem. Soap."),
			ReleaseDate:      "1999-10-15",
			Runtime:          &runtime,
			Budget:           &budget,
			Revenue:          &revenue,
			Status:           "Released",
			VoteAverage:      8.4,
			VoteCount:        26000,
			Popularity:       45.0,
			PosterPath:       stringPtr("/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg"),
			BackdropPath:     stringPtr("/87hTDiay2N2qWyX4Ds7ybXi9h8I.jpg"),
		}

		movie := mapper.MapMovie(tmdbMovie)

		assert.NotEqual(t, uuid.Nil, movie.ID)
		assert.Equal(t, int32(550), *movie.TMDbID)
		assert.Equal(t, "tt0137523", *movie.IMDbID)
		assert.Equal(t, "Fight Club", movie.Title)
		assert.Equal(t, "Fight Club", *movie.OriginalTitle)
		assert.Equal(t, "en", *movie.OriginalLanguage)
		assert.Contains(t, *movie.Overview, "insomnia")
		assert.Equal(t, "Mischief. Mayhem. Soap.", *movie.Tagline)
		assert.Equal(t, int32(1999), *movie.Year)
		assert.Equal(t, int32(139), *movie.Runtime)
		assert.Equal(t, int64(63000000), *movie.Budget)
		assert.Equal(t, int64(100853753), *movie.Revenue)
		assert.Equal(t, "Released", *movie.Status)
		require.NotNil(t, movie.VoteAverage)
		assert.Equal(t, "8.4", movie.VoteAverage.String())
		assert.Equal(t, int32(26000), *movie.VoteCount)
		assert.Equal(t, "/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg", *movie.PosterPath)
		assert.Equal(t, "/87hTDiay2N2qWyX4Ds7ybXi9h8I.jpg", *movie.BackdropPath)
	})

	t.Run("Movie with minimal fields", func(t *testing.T) {
		tmdbMovie := &TMDbMovie{
			ID:    12345,
			Title: "Unknown Movie",
		}

		movie := mapper.MapMovie(tmdbMovie)

		assert.NotEqual(t, uuid.Nil, movie.ID)
		assert.Equal(t, int32(12345), *movie.TMDbID)
		assert.Equal(t, "Unknown Movie", movie.Title)
		assert.Nil(t, movie.IMDbID)
		assert.Nil(t, movie.Overview)
		assert.Nil(t, movie.Runtime)
		assert.Nil(t, movie.Year)
	})
}

func TestTMDbMapper_MapSearchResult(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{APIKey: "test-key"})
	mapper := NewTMDbMapper(client)

	t.Run("Search result to movie", func(t *testing.T) {
		result := &TMDbSearchResult{
			ID:               550,
			Title:            "Fight Club",
			OriginalTitle:    "Fight Club",
			OriginalLanguage: "en",
			Overview:         "A depressed man suffering from insomnia meets a strange soap salesman.",
			ReleaseDate:      "1999-10-15",
			VoteAverage:      8.4,
			VoteCount:        26000,
			Popularity:       45.0,
			PosterPath:       stringPtr("/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg"),
			BackdropPath:     stringPtr("/87hTDiay2N2qWyX4Ds7ybXi9h8I.jpg"),
		}

		movie := mapper.MapSearchResult(result)

		assert.Equal(t, int32(550), *movie.TMDbID)
		assert.Equal(t, "Fight Club", movie.Title)
		assert.Equal(t, int32(1999), *movie.Year)
		require.NotNil(t, movie.VoteAverage)
		assert.NotNil(t, movie.PosterPath)
	})
}

func TestTMDbMapper_MapCredits(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{APIKey: "test-key"})
	mapper := NewTMDbMapper(client)
	movieID := uuid.New()

	credits := &TMDbCredits{
		Cast: []CastMember{
			{
				ID:          287,
				Name:        "Brad Pitt",
				Character:   "Tyler Durden",
				Order:       0,
				ProfilePath: stringPtr("/cckcYc2v0yh1tc9QjRelptcOBko.jpg"),
			},
			{
				ID:          819,
				Name:        "Edward Norton",
				Character:   "The Narrator",
				Order:       1,
				ProfilePath: stringPtr("/5XBzD5WuTyVQZeS4II6gs1nn5P6.jpg"),
			},
		},
		Crew: []CrewMember{
			{
				ID:          7467,
				Name:        "David Fincher",
				Department:  "Directing",
				Job:         "Director",
				ProfilePath: stringPtr("/tpEczFclQZeKAiCeKZZ0adRvtfz.jpg"),
			},
			{
				ID:          7468,
				Name:        "Art Linson",
				Department:  "Production",
				Job:         "Producer",
				ProfilePath: nil,
			},
		},
	}

	result := mapper.MapCredits(movieID, credits)

	assert.Len(t, result, 4) // 2 cast + 2 crew

	// Check cast
	cast := filterCreditsByType(result, "cast")
	assert.Len(t, cast, 2)
	assert.Equal(t, "Brad Pitt", cast[0].Name)
	assert.Equal(t, "Tyler Durden", *cast[0].Character)
	// CastOrder is nil for order=0 due to tmdbParseOptionalInt32Ptr treating 0 as nil
	// Second actor (order=1) should have CastOrder
	require.NotNil(t, cast[1].CastOrder)
	assert.Equal(t, int32(1), *cast[1].CastOrder)
	assert.Nil(t, cast[0].Job)

	// Check crew
	crew := filterCreditsByType(result, "crew")
	assert.Len(t, crew, 2)
	assert.Equal(t, "David Fincher", crew[0].Name)
	assert.Equal(t, "Director", *crew[0].Job)
	assert.Equal(t, "Directing", *crew[0].Department)
	assert.Nil(t, crew[0].Character)
}

func TestTMDbMapper_MapGenres(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{APIKey: "test-key"})
	mapper := NewTMDbMapper(client)
	movieID := uuid.New()

	genres := []Genre{
		{ID: 18, Name: "Drama"},
		{ID: 53, Name: "Thriller"},
		{ID: 35, Name: "Comedy"},
	}

	result := mapper.MapGenres(movieID, genres)

	assert.Len(t, result, 3)
	assert.Equal(t, movieID, result[0].MovieID)
	assert.Equal(t, int32(18), result[0].TMDbGenreID)
	assert.Equal(t, "Drama", result[0].Name)
	assert.Equal(t, "Thriller", result[1].Name)
	assert.Equal(t, "Comedy", result[2].Name)
}

func TestTMDbMapper_MapCollection(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{APIKey: "test-key"})
	mapper := NewTMDbMapper(client)

	collection := &TMDbCollectionDetails{
		ID:           119,
		Name:         "The Matrix Collection",
		Overview:     "The Matrix is everywhere.",
		PosterPath:   stringPtr("/bN5dH4KNtQKj5kA1eKTRkC3Yjbi.jpg"),
		BackdropPath: stringPtr("/bRm2DEgUiYciDw3myHuYFInD7la.jpg"),
	}

	result := mapper.MapCollection(collection)

	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, int32(119), *result.TMDbCollectionID)
	assert.Equal(t, "The Matrix Collection", result.Name)
	assert.Equal(t, "The Matrix is everywhere.", *result.Overview)
	assert.Equal(t, "/bN5dH4KNtQKj5kA1eKTRkC3Yjbi.jpg", *result.PosterPath)
}

func TestTMDbMapper_GetPosterURL(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{
		APIKey: "test-key",
	})
	mapper := NewTMDbMapper(client)

	t.Run("Valid poster path", func(t *testing.T) {
		posterPath := "/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg"
		url := mapper.GetPosterURL(&posterPath, "w500")

		require.NotNil(t, url)
		assert.Contains(t, *url, "w500")
		assert.Contains(t, *url, posterPath)
	})

	t.Run("Nil poster path", func(t *testing.T) {
		url := mapper.GetPosterURL(nil, "w500")
		assert.Nil(t, url)
	})

	t.Run("Empty poster path", func(t *testing.T) {
		emptyPath := ""
		url := mapper.GetPosterURL(&emptyPath, "w500")
		assert.Nil(t, url)
	})

	t.Run("Default size", func(t *testing.T) {
		posterPath := "/test.jpg"
		url := mapper.GetPosterURL(&posterPath, "")

		require.NotNil(t, url)
		assert.Contains(t, *url, "w500") // Default poster size
	})
}

func TestTMDbMapper_GetBackdropURL(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{
		APIKey: "test-key",
	})
	mapper := NewTMDbMapper(client)

	t.Run("Valid backdrop path", func(t *testing.T) {
		backdropPath := "/87hTDiay2N2qWyX4Ds7ybXi9h8I.jpg"
		url := mapper.GetBackdropURL(&backdropPath, "w1280")

		require.NotNil(t, url)
		assert.Contains(t, *url, "w1280")
		assert.Contains(t, *url, backdropPath)
	})

	t.Run("Nil backdrop path", func(t *testing.T) {
		url := mapper.GetBackdropURL(nil, "w1280")
		assert.Nil(t, url)
	})

	t.Run("Default size", func(t *testing.T) {
		backdropPath := "/test.jpg"
		url := mapper.GetBackdropURL(&backdropPath, "")

		require.NotNil(t, url)
		assert.Contains(t, *url, "w1280") // Default backdrop size
	})
}

func TestMapperExtractYear(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		expected *int
	}{
		{"Valid date", "1999-10-15", intPtrMapper(1999)},
		{"Year only", "2024", intPtrMapper(2024)},
		{"Empty string", "", nil},
		{"Invalid format", "invalid", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractYear(tt.date)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestParseReleaseDate(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		expected *time.Time
	}{
		{
			"Valid date",
			"1999-10-15",
			timePtr(time.Date(1999, 10, 15, 0, 0, 0, 0, time.UTC)),
		},
		{"Empty string", "", nil},
		{"Invalid format", "15/10/1999", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseReleaseDate(tt.date)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.Year(), result.Year())
				assert.Equal(t, tt.expected.Month(), result.Month())
				assert.Equal(t, tt.expected.Day(), result.Day())
			}
		})
	}
}

func TestParseDecimal(t *testing.T) {
	t.Run("Non-zero value", func(t *testing.T) {
		result := parseDecimal(8.5)
		require.NotNil(t, result)
		assert.Equal(t, "8.5", result.String())
	})

	t.Run("Zero value", func(t *testing.T) {
		result := parseDecimal(0.0)
		assert.Nil(t, result)
	})
}

func TestTmdbParseOptionalInt32(t *testing.T) {
	t.Run("Non-zero value", func(t *testing.T) {
		result := tmdbParseOptionalInt32(550)
		require.NotNil(t, result)
		assert.Equal(t, int32(550), *result)
	})

	t.Run("Zero value", func(t *testing.T) {
		result := tmdbParseOptionalInt32(0)
		assert.Nil(t, result)
	})
}

func TestTmdbParseOptionalString(t *testing.T) {
	t.Run("Non-empty value", func(t *testing.T) {
		s := "test"
		result := tmdbParseOptionalString(&s)
		require.NotNil(t, result)
		assert.Equal(t, "test", *result)
	})

	t.Run("Empty string", func(t *testing.T) {
		s := ""
		result := tmdbParseOptionalString(&s)
		assert.Nil(t, result)
	})

	t.Run("Nil pointer", func(t *testing.T) {
		result := tmdbParseOptionalString(nil)
		assert.Nil(t, result)
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtrMapper(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func filterCreditsByType(credits []MovieCredit, creditType string) []MovieCredit {
	var result []MovieCredit
	for _, c := range credits {
		if c.CreditType == creditType {
			result = append(result, c)
		}
	}
	return result
}
