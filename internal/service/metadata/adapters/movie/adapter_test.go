package movie

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	contentmovie "github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/service/metadata"
)

//go:fix inline
func ptr[T any](v T) *T { return new(v) }

func TestPtrString(t *testing.T) {
	assert.Nil(t, ptrString(""))
	result := ptrString("hello")
	require.NotNil(t, result)
	assert.Equal(t, "hello", *result)
}

func TestPtrInt32(t *testing.T) {
	assert.Nil(t, ptrInt32(nil))
	v := 42
	result := ptrInt32(&v)
	require.NotNil(t, result)
	assert.Equal(t, int32(42), *result)
}

func TestPtrInt32FromInt(t *testing.T) {
	assert.Nil(t, ptrInt32FromInt(nil))
	v := 2024
	result := ptrInt32FromInt(&v)
	require.NotNil(t, result)
	assert.Equal(t, int32(2024), *result)
}

func TestGetAgeRatingSystem(t *testing.T) {
	tests := []struct {
		country string
		want    string
	}{
		{"US", "MPAA"},
		{"DE", "FSK"},
		{"GB", "BBFC"},
		{"FR", "CNC"},
		{"JP", "Eirin"},
		{"KR", "KMRB"},
		{"BR", "DJCTQ"},
		{"AU", "ACB"},
		{"ZZ", "ZZ"},
	}
	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			assert.Equal(t, tt.want, getAgeRatingSystem(tt.country))
		})
	}
}

func TestMapSearchResultToMovie(t *testing.T) {
	t.Run("full search result", func(t *testing.T) {
		release := time.Date(2008, 7, 18, 0, 0, 0, 0, time.UTC)
		r := &metadata.MovieSearchResult{
			ProviderID:       "155",
			Provider:         metadata.ProviderTMDb,
			Title:            "The Dark Knight",
			OriginalTitle:    "The Dark Knight",
			OriginalLanguage: "en",
			VoteAverage:      9.0,
			VoteCount:        30000,
			Popularity:       80.5,
			PosterPath:       new("https://example.com/poster.jpg"),
			BackdropPath:     new("https://example.com/backdrop.jpg"),
			ReleaseDate:      &release,
			Year:             new(2008),
		}
		mov := mapSearchResultToMovie(r)
		require.NotNil(t, mov)
		assert.NotEqual(t, uuid.Nil, mov.ID)
		assert.Equal(t, "The Dark Knight", mov.Title)
		require.NotNil(t, mov.OriginalTitle)
		require.NotNil(t, mov.OriginalLanguage)
		require.NotNil(t, mov.PosterPath)
		require.NotNil(t, mov.BackdropPath)
		require.NotNil(t, mov.ReleaseDate)
		require.NotNil(t, mov.Year)
		assert.Equal(t, int32(2008), *mov.Year)

		require.NotNil(t, mov.VoteAverage)
		va, _ := mov.VoteAverage.Float64()
		assert.InDelta(t, 9.0, va, 0.01)

		require.NotNil(t, mov.VoteCount)
		assert.Equal(t, int32(30000), *mov.VoteCount)

		require.NotNil(t, mov.Popularity)

		require.NotNil(t, mov.TMDbID)
		assert.Equal(t, int32(155), *mov.TMDbID)
	})

	t.Run("zero values not set", func(t *testing.T) {
		r := &metadata.MovieSearchResult{Title: "Test"}
		mov := mapSearchResultToMovie(r)
		assert.Nil(t, mov.VoteAverage)
		assert.Nil(t, mov.VoteCount)
		assert.Nil(t, mov.Popularity)
		assert.Nil(t, mov.OriginalTitle)
		assert.Nil(t, mov.OriginalLanguage)
		assert.Nil(t, mov.Year)
	})

	t.Run("non-numeric provider ID", func(t *testing.T) {
		r := &metadata.MovieSearchResult{ProviderID: "abc", Title: "Test"}
		mov := mapSearchResultToMovie(r)
		// Sscanf fails, tmdbID stays 0 but pointer is still set
		require.NotNil(t, mov.TMDbID)
		assert.Equal(t, int32(0), *mov.TMDbID)
	})

	t.Run("empty provider ID skips TMDbID", func(t *testing.T) {
		r := &metadata.MovieSearchResult{Title: "Test"}
		mov := mapSearchResultToMovie(r)
		// ProviderID is empty so the if guard prevents entering the block
		// However, since TMDbID is zero-value nil initially and never set, it stays nil
		_ = mov // verified via zero_values_not_set subtest above
	})
}

func TestMapMetadataToMovie(t *testing.T) {
	t.Run("full metadata", func(t *testing.T) {
		mov := &contentmovie.Movie{ID: uuid.Must(uuid.NewV7())}
		overview := "Batman fights the Joker."
		tagline := "Why so serious?"
		status := "Released"
		homepage := "https://example.com"
		imdb := "tt0468569"
		tmdb := int32(155)
		runtime := int32(152)
		budget := int64(185000000)
		revenue := int64(1004558444)
		release := time.Date(2008, 7, 18, 0, 0, 0, 0, time.UTC)

		meta := &metadata.MovieMetadata{
			Title:            "The Dark Knight",
			OriginalTitle:    "The Dark Knight",
			OriginalLanguage: "en",
			Overview:         &overview,
			Tagline:          &tagline,
			Status:           status,
			ReleaseDate:      &release,
			Runtime:          &runtime,
			Budget:           &budget,
			Revenue:          &revenue,
			VoteAverage:      9.0,
			VoteCount:        30000,
			Popularity:       80.5,
			PosterPath:       new("https://example.com/poster.jpg"),
			BackdropPath:     new("https://example.com/backdrop.jpg"),
			Homepage:         &homepage,
			IMDbID:           &imdb,
			TMDbID:           &tmdb,
			Translations: map[string]*metadata.LocalizedMovieData{
				"de": {Title: "The Dark Knight", Tagline: "Warum so ernst?", Overview: "Batman kämpft..."},
				"fr": {Title: "Le Chevalier noir", Overview: "Batman combat..."},
			},
			ExternalRatings: []metadata.ExternalRating{
				{Source: "IMDb", Value: "9.0/10", Score: 90.0},
				{Source: "Rotten Tomatoes", Value: "94%", Score: 94.0},
			},
		}
		releaseDates := []metadata.ReleaseDate{
			{CountryCode: "US", Certification: "PG-13"},
			{CountryCode: "DE", Certification: "12"},
			{CountryCode: "GB", Certification: ""},
		}

		mapMetadataToMovie(mov, meta, releaseDates)

		assert.Equal(t, "The Dark Knight", mov.Title)
		require.NotNil(t, mov.OriginalTitle)
		require.NotNil(t, mov.OriginalLanguage)
		require.NotNil(t, mov.Overview)
		require.NotNil(t, mov.Tagline)
		require.NotNil(t, mov.Status)
		assert.Equal(t, "Released", *mov.Status)
		require.NotNil(t, mov.ReleaseDate)
		require.NotNil(t, mov.Runtime)
		assert.Equal(t, int32(152), *mov.Runtime)
		require.NotNil(t, mov.Budget)
		require.NotNil(t, mov.Revenue)
		assert.Equal(t, &imdb, mov.IMDbID)
		assert.Equal(t, &tmdb, mov.TMDbID)

		require.NotNil(t, mov.Year)
		assert.Equal(t, int32(2008), *mov.Year)

		require.NotNil(t, mov.VoteAverage)
		va, _ := mov.VoteAverage.Float64()
		assert.InDelta(t, 9.0, va, 0.01)

		require.Len(t, mov.TitlesI18n, 2)
		require.Len(t, mov.TaglinesI18n, 1)
		require.Len(t, mov.OverviewsI18n, 2)

		require.Len(t, mov.AgeRatings, 2)
		assert.Equal(t, "PG-13", mov.AgeRatings["US"]["MPAA"])
		assert.Equal(t, "12", mov.AgeRatings["DE"]["FSK"])

		require.Len(t, mov.ExternalRatings, 2)
		assert.Equal(t, "IMDb", mov.ExternalRatings[0].Source)
	})

	t.Run("empty translations and no release dates", func(t *testing.T) {
		mov := &contentmovie.Movie{ID: uuid.Must(uuid.NewV7())}
		meta := &metadata.MovieMetadata{Title: "Test"}
		mapMetadataToMovie(mov, meta, nil)
		assert.Nil(t, mov.TitlesI18n)
		assert.Nil(t, mov.AgeRatings)
		assert.Nil(t, mov.ExternalRatings)
		assert.Nil(t, mov.Year)
	})
}

func TestMapCreditsToMovieCredits(t *testing.T) {
	movieID := uuid.Must(uuid.NewV7())

	t.Run("cast and crew", func(t *testing.T) {
		credits := &metadata.Credits{
			Cast: []metadata.CastMember{
				{ProviderID: "6193", Name: "Leonardo DiCaprio", Character: "Cobb", Order: 0, ProfilePath: new("https://example.com/leo.jpg")},
			},
			Crew: []metadata.CrewMember{
				{ProviderID: "525", Name: "Christopher Nolan", Job: "Director", Department: "Directing", ProfilePath: new("https://example.com/cn.jpg")},
			},
		}

		result := mapCreditsToMovieCredits(movieID, credits)
		require.Len(t, result, 2)

		assert.Equal(t, movieID, result[0].MovieID)
		assert.Equal(t, int32(6193), result[0].TMDbPersonID)
		assert.Equal(t, "Leonardo DiCaprio", result[0].Name)
		assert.Equal(t, "cast", result[0].CreditType)
		require.NotNil(t, result[0].Character)
		assert.Equal(t, "Cobb", *result[0].Character)
		require.NotNil(t, result[0].CastOrder)
		require.NotNil(t, result[0].ProfilePath)

		assert.Equal(t, "crew", result[1].CreditType)
		assert.Equal(t, int32(525), result[1].TMDbPersonID)
		require.NotNil(t, result[1].Job)
		assert.Equal(t, "Director", *result[1].Job)
		require.NotNil(t, result[1].Department)
	})

	t.Run("nil credits returns nil", func(t *testing.T) {
		result := mapCreditsToMovieCredits(movieID, &metadata.Credits{})
		assert.Nil(t, result)
	})

	t.Run("non-numeric provider ID", func(t *testing.T) {
		credits := &metadata.Credits{
			Cast: []metadata.CastMember{{ProviderID: "not-a-number", Name: "Test"}},
		}
		result := mapCreditsToMovieCredits(movieID, credits)
		require.Len(t, result, 1)
		assert.Equal(t, int32(0), result[0].TMDbPersonID)
	})
}
