package movie

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTMDbMapper_MapMultiLanguageMovie(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{APIKey: "test-key"})
	mapper := NewTMDbMapper(client)

	t.Run("maps movie with multiple languages", func(t *testing.T) {
		enOverview := "A depressed man suffering from insomnia meets a strange soap salesman."
		enTagline := "Mischief. Mayhem. Soap."
		deOverview := "Ein depressiver Mann, der unter Schlaflosigkeit leidet..."
		deTagline := "Unfug. Chaos. Seife."
		frOverview := "Un homme déprimé souffrant d'insomnie rencontre..."
		frTagline := "Espièglerie. Chaos. Savon."

		multiLang := &TMDbMultiLanguageResult{
			Movies: map[string]*TMDbMovie{
				"en": {
					ID:               550,
					Title:            "Fight Club",
					OriginalTitle:    "Fight Club",
					OriginalLanguage: "en",
					Overview:         &enOverview,
					Tagline:          &enTagline,
					ReleaseDate:      "1999-10-15",
				},
				"de": {
					ID:               550,
					Title:            "Fight Club",
					OriginalTitle:    "Fight Club",
					OriginalLanguage: "en",
					Overview:         &deOverview,
					Tagline:          &deTagline,
					ReleaseDate:      "1999-10-15",
				},
				"fr": {
					ID:               550,
					Title:            "Fight Club",
					OriginalTitle:    "Fight Club",
					OriginalLanguage: "en",
					Overview:         &frOverview,
					Tagline:          &frTagline,
					ReleaseDate:      "1999-10-15",
				},
			},
		}

		releaseDates := &TMDbReleaseDatesResponse{
			ID: 550,
			Results: []TMDbCountryRelease{
				{
					ISO3166_1: "US",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "R", Type: 3},
					},
				},
				{
					ISO3166_1: "DE",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "16", Type: 3},
					},
				},
			},
		}

		movie := mapper.MapMultiLanguageMovie(multiLang, releaseDates)

		require.NotNil(t, movie)

		// Check base fields from English
		assert.Equal(t, "Fight Club", movie.Title)
		assert.NotNil(t, movie.TMDbID)
		assert.Equal(t, int32(550), *movie.TMDbID)

		// Check i18n fields
		assert.Len(t, movie.TitlesI18n, 3)
		assert.Equal(t, "Fight Club", movie.TitlesI18n["en"])
		assert.Equal(t, "Fight Club", movie.TitlesI18n["de"])
		assert.Equal(t, "Fight Club", movie.TitlesI18n["fr"])

		assert.Len(t, movie.TaglinesI18n, 3)
		assert.Equal(t, "Mischief. Mayhem. Soap.", movie.TaglinesI18n["en"])
		assert.Equal(t, "Unfug. Chaos. Seife.", movie.TaglinesI18n["de"])
		assert.Equal(t, "Espièglerie. Chaos. Savon.", movie.TaglinesI18n["fr"])

		assert.Len(t, movie.OverviewsI18n, 3)
		assert.Contains(t, movie.OverviewsI18n["en"], "insomnia")
		assert.Contains(t, movie.OverviewsI18n["de"], "Schlaflosigkeit")
		assert.Contains(t, movie.OverviewsI18n["fr"], "insomnie")

		// Check age ratings
		assert.Len(t, movie.AgeRatings, 2)
		assert.Equal(t, "R", movie.AgeRatings["US"]["MPAA"])
		assert.Equal(t, "16", movie.AgeRatings["DE"]["FSK"])
	})

	t.Run("returns nil when English not present", func(t *testing.T) {
		multiLang := &TMDbMultiLanguageResult{
			Movies: map[string]*TMDbMovie{
				"de": {
					ID:               550,
					Title:            "Fight Club",
					OriginalTitle:    "Fight Club",
					OriginalLanguage: "en",
				},
			},
		}

		movie := mapper.MapMultiLanguageMovie(multiLang, nil)

		assert.Nil(t, movie, "Should return nil when English is not present")
	})

	t.Run("handles empty taglines and overviews", func(t *testing.T) {
		emptyStr := ""
		overview := "Some overview"

		multiLang := &TMDbMultiLanguageResult{
			Movies: map[string]*TMDbMovie{
				"en": {
					ID:               550,
					Title:            "Fight Club",
					OriginalTitle:    "Fight Club",
					OriginalLanguage: "en",
					Overview:         &overview,
					Tagline:          &emptyStr, // Empty tagline
				},
				"de": {
					ID:               550,
					Title:            "Fight Club",
					OriginalTitle:    "Fight Club",
					OriginalLanguage: "en",
					Overview:         &emptyStr, // Empty overview
					Tagline:          nil,       // Nil tagline
				},
			},
		}

		movie := mapper.MapMultiLanguageMovie(multiLang, nil)

		require.NotNil(t, movie)

		// Titles should always be mapped
		assert.Len(t, movie.TitlesI18n, 2)

		// Taglines: only "en" has empty string, de has nil - neither should be mapped
		assert.Empty(t, movie.TaglinesI18n, "Empty and nil taglines should not be mapped")

		// Overviews: en has content, de is empty
		assert.Len(t, movie.OverviewsI18n, 1)
		assert.Contains(t, movie.OverviewsI18n, "en")
	})

	t.Run("handles nil release dates", func(t *testing.T) {
		overview := "Overview"
		multiLang := &TMDbMultiLanguageResult{
			Movies: map[string]*TMDbMovie{
				"en": {
					ID:               550,
					Title:            "Fight Club",
					OriginalTitle:    "Fight Club",
					OriginalLanguage: "en",
					Overview:         &overview,
				},
			},
		}

		movie := mapper.MapMultiLanguageMovie(multiLang, nil)

		require.NotNil(t, movie)
		assert.Nil(t, movie.AgeRatings, "Age ratings should be nil when release dates not provided")
	})

	t.Run("handles single language", func(t *testing.T) {
		overview := "Overview"
		tagline := "Tagline"

		multiLang := &TMDbMultiLanguageResult{
			Movies: map[string]*TMDbMovie{
				"en": {
					ID:               550,
					Title:            "Fight Club",
					OriginalTitle:    "Fight Club",
					OriginalLanguage: "en",
					Overview:         &overview,
					Tagline:          &tagline,
				},
			},
		}

		movie := mapper.MapMultiLanguageMovie(multiLang, nil)

		require.NotNil(t, movie)
		assert.Len(t, movie.TitlesI18n, 1)
		assert.Len(t, movie.TaglinesI18n, 1)
		assert.Len(t, movie.OverviewsI18n, 1)
	})
}

func TestTMDbMapper_MapAgeRatings(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{APIKey: "test-key"})
	mapper := NewTMDbMapper(client)

	t.Run("maps age ratings from release dates", func(t *testing.T) {
		releaseDates := &TMDbReleaseDatesResponse{
			ID: 550,
			Results: []TMDbCountryRelease{
				{
					ISO3166_1: "US",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "R", Type: 3}, // Theatrical
					},
				},
				{
					ISO3166_1: "DE",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "16", Type: 3},
					},
				},
				{
					ISO3166_1: "GB",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "18", Type: 3},
					},
				},
				{
					ISO3166_1: "FR",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "12", Type: 3},
					},
				},
			},
		}

		ratings := mapper.MapAgeRatings(releaseDates)

		assert.Len(t, ratings, 4)
		assert.Equal(t, "R", ratings["US"]["MPAA"])
		assert.Equal(t, "16", ratings["DE"]["FSK"])
		assert.Equal(t, "18", ratings["GB"]["BBFC"])
		assert.Equal(t, "12", ratings["FR"]["CNC"])
	})

	t.Run("ignores non-theatrical releases", func(t *testing.T) {
		releaseDates := &TMDbReleaseDatesResponse{
			ID: 550,
			Results: []TMDbCountryRelease{
				{
					ISO3166_1: "US",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "PG", Type: 1},  // Premiere - should be ignored
						{Certification: "R", Type: 3},   // Theatrical - should be used
						{Certification: "PG", Type: 4},  // Digital - should be ignored
					},
				},
			},
		}

		ratings := mapper.MapAgeRatings(releaseDates)

		assert.Len(t, ratings, 1)
		assert.Equal(t, "R", ratings["US"]["MPAA"], "Should use theatrical release (type 3)")
	})

	t.Run("ignores empty certifications", func(t *testing.T) {
		releaseDates := &TMDbReleaseDatesResponse{
			ID: 550,
			Results: []TMDbCountryRelease{
				{
					ISO3166_1: "US",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "", Type: 3}, // Empty certification
					},
				},
				{
					ISO3166_1: "DE",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "16", Type: 3}, // Valid
					},
				},
			},
		}

		ratings := mapper.MapAgeRatings(releaseDates)

		assert.Len(t, ratings, 1, "Should only include DE, not US")
		assert.NotContains(t, ratings, "US")
		assert.Equal(t, "16", ratings["DE"]["FSK"])
	})

	t.Run("uses first theatrical release when multiple present", func(t *testing.T) {
		releaseDates := &TMDbReleaseDatesResponse{
			ID: 550,
			Results: []TMDbCountryRelease{
				{
					ISO3166_1: "US",
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "R", Type: 3, ReleaseDate: "1999-10-15"},   // First
						{Certification: "PG-13", Type: 3, ReleaseDate: "1999-11-01"}, // Second
					},
				},
			},
		}

		ratings := mapper.MapAgeRatings(releaseDates)

		assert.Equal(t, "R", ratings["US"]["MPAA"], "Should use first theatrical release")
	})

	t.Run("handles unknown countries with fallback", func(t *testing.T) {
		releaseDates := &TMDbReleaseDatesResponse{
			ID: 550,
			Results: []TMDbCountryRelease{
				{
					ISO3166_1: "XX", // Unknown country
					ReleaseDates: []TMDbReleaseDate{
						{Certification: "18", Type: 3},
					},
				},
			},
		}

		ratings := mapper.MapAgeRatings(releaseDates)

		assert.Len(t, ratings, 1)
		assert.Equal(t, "18", ratings["XX"]["XX"], "Should use country code as system name for unknown countries")
	})

	t.Run("handles empty release dates", func(t *testing.T) {
		releaseDates := &TMDbReleaseDatesResponse{
			ID:      550,
			Results: []TMDbCountryRelease{},
		}

		ratings := mapper.MapAgeRatings(releaseDates)

		assert.Empty(t, ratings)
	})

	t.Run("handles all supported rating systems", func(t *testing.T) {
		releaseDates := &TMDbReleaseDatesResponse{
			ID: 550,
			Results: []TMDbCountryRelease{
				{ISO3166_1: "US", ReleaseDates: []TMDbReleaseDate{{Certification: "R", Type: 3}}},
				{ISO3166_1: "DE", ReleaseDates: []TMDbReleaseDate{{Certification: "16", Type: 3}}},
				{ISO3166_1: "GB", ReleaseDates: []TMDbReleaseDate{{Certification: "15", Type: 3}}},
				{ISO3166_1: "FR", ReleaseDates: []TMDbReleaseDate{{Certification: "12", Type: 3}}},
				{ISO3166_1: "JP", ReleaseDates: []TMDbReleaseDate{{Certification: "PG12", Type: 3}}},
				{ISO3166_1: "KR", ReleaseDates: []TMDbReleaseDate{{Certification: "15", Type: 3}}},
				{ISO3166_1: "BR", ReleaseDates: []TMDbReleaseDate{{Certification: "16", Type: 3}}},
				{ISO3166_1: "AU", ReleaseDates: []TMDbReleaseDate{{Certification: "MA15+", Type: 3}}},
			},
		}

		ratings := mapper.MapAgeRatings(releaseDates)

		assert.Equal(t, "R", ratings["US"]["MPAA"])
		assert.Equal(t, "16", ratings["DE"]["FSK"])
		assert.Equal(t, "15", ratings["GB"]["BBFC"])
		assert.Equal(t, "12", ratings["FR"]["CNC"])
		assert.Equal(t, "PG12", ratings["JP"]["Eirin"])
		assert.Equal(t, "15", ratings["KR"]["KMRB"])
		assert.Equal(t, "16", ratings["BR"]["DJCTQ"])
		assert.Equal(t, "MA15+", ratings["AU"]["ACB"])
	})
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
		{"XX", "XX"}, // Unknown country uses country code
	}

	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			got := getAgeRatingSystem(tt.country)
			assert.Equal(t, tt.want, got)
		})
	}
}
