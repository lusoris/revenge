package movie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovie_GetTitle(t *testing.T) {
	originalLang := "ja"
	originalTitle := "千と千尋の神隠し"

	t.Run("returns title in requested language", func(t *testing.T) {
		movie := &Movie{
			Title: "Spirited Away",
			TitlesI18n: map[string]string{
				"en": "Spirited Away",
				"de": "Chihiros Reise ins Zauberland",
				"fr": "Le Voyage de Chihiro",
			},
		}

		assert.Equal(t, "Chihiros Reise ins Zauberland", movie.GetTitle("de"))
		assert.Equal(t, "Le Voyage de Chihiro", movie.GetTitle("fr"))
	})

	t.Run("falls back to English when requested language not available", func(t *testing.T) {
		movie := &Movie{
			Title: "Spirited Away",
			TitlesI18n: map[string]string{
				"en": "Spirited Away",
				"de": "Chihiros Reise ins Zauberland",
			},
		}

		assert.Equal(t, "Spirited Away", movie.GetTitle("fr")) // French not available, returns English
	})

	t.Run("falls back to original language when English not available", func(t *testing.T) {
		movie := &Movie{
			Title:            "Spirited Away",
			OriginalLanguage: &originalLang,
			TitlesI18n: map[string]string{
				"ja": originalTitle,
				"de": "Chihiros Reise ins Zauberland",
			},
		}

		assert.Equal(t, originalTitle, movie.GetTitle("fr")) // French and English not available, returns original
	})

	t.Run("falls back to OriginalTitle field", func(t *testing.T) {
		movie := &Movie{
			Title:         "Spirited Away",
			OriginalTitle: &originalTitle,
			TitlesI18n:    map[string]string{},
		}

		assert.Equal(t, originalTitle, movie.GetTitle("de"))
	})

	t.Run("falls back to default Title field", func(t *testing.T) {
		movie := &Movie{
			Title: "Spirited Away",
		}

		assert.Equal(t, "Spirited Away", movie.GetTitle("de"))
	})

	t.Run("handles nil TitlesI18n", func(t *testing.T) {
		movie := &Movie{
			Title:      "Spirited Away",
			TitlesI18n: nil,
		}

		assert.Equal(t, "Spirited Away", movie.GetTitle("en"))
	})

	t.Run("handles empty strings in TitlesI18n", func(t *testing.T) {
		movie := &Movie{
			Title: "Spirited Away",
			TitlesI18n: map[string]string{
				"de": "", // Empty string should be skipped
				"en": "Spirited Away",
			},
		}

		assert.Equal(t, "Spirited Away", movie.GetTitle("de")) // Skip empty de, return en
	})
}

func TestMovie_GetTagline(t *testing.T) {
	tagline := "Fear can hold you prisoner. Hope can set you free."

	t.Run("returns tagline in requested language", func(t *testing.T) {
		movie := &Movie{
			Tagline: &tagline,
			TaglinesI18n: map[string]string{
				"en": "Fear can hold you prisoner. Hope can set you free.",
				"de": "Angst kann dich gefangen halten. Hoffnung kann dich befreien.",
			},
		}

		assert.Equal(t, "Angst kann dich gefangen halten. Hoffnung kann dich befreien.", movie.GetTagline("de"))
	})

	t.Run("falls back to English", func(t *testing.T) {
		movie := &Movie{
			Tagline: &tagline,
			TaglinesI18n: map[string]string{
				"en": "Fear can hold you prisoner. Hope can set you free.",
			},
		}

		assert.Equal(t, "Fear can hold you prisoner. Hope can set you free.", movie.GetTagline("fr"))
	})

	t.Run("falls back to default Tagline field", func(t *testing.T) {
		movie := &Movie{
			Tagline: &tagline,
		}

		assert.Equal(t, tagline, movie.GetTagline("de"))
	})

	t.Run("returns empty string when no tagline available", func(t *testing.T) {
		movie := &Movie{}

		assert.Empty(t, movie.GetTagline("en"))
	})
}

func TestMovie_GetOverview(t *testing.T) {
	overview := "Imprisoned in the 1940s for the double murder..."

	t.Run("returns overview in requested language", func(t *testing.T) {
		movie := &Movie{
			Overview: &overview,
			OverviewsI18n: map[string]string{
				"en": "Imprisoned in the 1940s for the double murder...",
				"de": "In den 1940er Jahren eingesperrt...",
			},
		}

		assert.Equal(t, "In den 1940er Jahren eingesperrt...", movie.GetOverview("de"))
	})

	t.Run("falls back to English", func(t *testing.T) {
		movie := &Movie{
			Overview: &overview,
			OverviewsI18n: map[string]string{
				"en": "Imprisoned in the 1940s for the double murder...",
			},
		}

		assert.Equal(t, "Imprisoned in the 1940s for the double murder...", movie.GetOverview("fr"))
	})

	t.Run("falls back to default Overview field", func(t *testing.T) {
		movie := &Movie{
			Overview: &overview,
		}

		assert.Equal(t, overview, movie.GetOverview("de"))
	})

	t.Run("returns empty string when no overview available", func(t *testing.T) {
		movie := &Movie{}

		assert.Empty(t, movie.GetOverview("en"))
	})
}

func TestMovie_GetAgeRating(t *testing.T) {
	t.Run("returns age rating for country and system", func(t *testing.T) {
		movie := &Movie{
			AgeRatings: map[string]map[string]string{
				"US": {"MPAA": "R"},
				"DE": {"FSK": "12"},
				"GB": {"BBFC": "15"},
			},
		}

		assert.Equal(t, "R", movie.GetAgeRating("US", "MPAA"))
		assert.Equal(t, "12", movie.GetAgeRating("DE", "FSK"))
		assert.Equal(t, "15", movie.GetAgeRating("GB", "BBFC"))
	})

	t.Run("returns empty string when country not found", func(t *testing.T) {
		movie := &Movie{
			AgeRatings: map[string]map[string]string{
				"US": {"MPAA": "R"},
			},
		}

		assert.Empty(t, movie.GetAgeRating("FR", "CNC"))
	})

	t.Run("returns empty string when system not found", func(t *testing.T) {
		movie := &Movie{
			AgeRatings: map[string]map[string]string{
				"US": {"MPAA": "R"},
			},
		}

		assert.Empty(t, movie.GetAgeRating("US", "FSK"))
	})

	t.Run("handles nil AgeRatings", func(t *testing.T) {
		movie := &Movie{
			AgeRatings: nil,
		}

		assert.Empty(t, movie.GetAgeRating("US", "MPAA"))
	})
}

func TestMovie_GetAvailableLanguages(t *testing.T) {
	t.Run("returns list of available languages", func(t *testing.T) {
		movie := &Movie{
			TitlesI18n: map[string]string{
				"en": "The Shawshank Redemption",
				"de": "Die Verurteilten",
				"fr": "Les Évadés",
			},
		}

		langs := movie.GetAvailableLanguages()

		assert.Len(t, langs, 3)
		assert.Contains(t, langs, "en")
		assert.Contains(t, langs, "de")
		assert.Contains(t, langs, "fr")
	})

	t.Run("returns empty slice when no translations", func(t *testing.T) {
		movie := &Movie{
			TitlesI18n: nil,
		}

		langs := movie.GetAvailableLanguages()

		assert.Empty(t, langs)
		assert.NotNil(t, langs) // Should return empty slice, not nil
	})

	t.Run("returns empty slice for empty map", func(t *testing.T) {
		movie := &Movie{
			TitlesI18n: map[string]string{},
		}

		langs := movie.GetAvailableLanguages()

		assert.Empty(t, langs)
	})
}

func TestMovie_GetAvailableAgeRatingCountries(t *testing.T) {
	t.Run("returns list of countries with age ratings", func(t *testing.T) {
		movie := &Movie{
			AgeRatings: map[string]map[string]string{
				"US": {"MPAA": "R"},
				"DE": {"FSK": "12"},
				"GB": {"BBFC": "15"},
			},
		}

		countries := movie.GetAvailableAgeRatingCountries()

		assert.Len(t, countries, 3)
		assert.Contains(t, countries, "US")
		assert.Contains(t, countries, "DE")
		assert.Contains(t, countries, "GB")
	})

	t.Run("returns empty slice when no ratings", func(t *testing.T) {
		movie := &Movie{
			AgeRatings: nil,
		}

		countries := movie.GetAvailableAgeRatingCountries()

		assert.Empty(t, countries)
		assert.NotNil(t, countries) // Should return empty slice, not nil
	})

	t.Run("returns empty slice for empty map", func(t *testing.T) {
		movie := &Movie{
			AgeRatings: map[string]map[string]string{},
		}

		countries := movie.GetAvailableAgeRatingCountries()

		assert.Empty(t, countries)
	})
}
