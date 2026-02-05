package api

import (
	"testing"

	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{
			name:     "German with quality values",
			header:   "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7",
			expected: "de",
		},
		{
			name:     "French with region",
			header:   "fr-FR,fr;q=0.9,en;q=0.8",
			expected: "fr",
		},
		{
			name:     "Simple English",
			header:   "en-US",
			expected: "en",
		},
		{
			name:     "No quality values",
			header:   "de-DE",
			expected: "de",
		},
		{
			name:     "Simple language code",
			header:   "ja",
			expected: "ja",
		},
		{
			name:     "Empty header",
			header:   "",
			expected: "",
		},
		{
			name:     "Complex with spaces",
			header:   " es-ES , es ; q=0.9 , en ; q=0.8 ",
			expected: "es",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAcceptLanguage(tt.header)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLocalizeMovie(t *testing.T) {
	overview := "English overview"
	tagline := "English tagline"
	originalTitle := "Original Title"
	originalLang := "en"

	t.Run("localizes to requested language", func(t *testing.T) {
		m := &movie.Movie{
			Title:            "English Title",
			Overview:         &overview,
			Tagline:          &tagline,
			OriginalTitle:    &originalTitle,
			OriginalLanguage: &originalLang,
			TitlesI18n: map[string]string{
				"en": "English Title",
				"de": "Deutscher Titel",
				"fr": "Titre Français",
			},
			OverviewsI18n: map[string]string{
				"en": "English overview",
				"de": "Deutsche Übersicht",
				"fr": "Aperçu français",
			},
			TaglinesI18n: map[string]string{
				"en": "English tagline",
				"de": "Deutscher Slogan",
				"fr": "Slogan français",
			},
		}

		localized := LocalizeMovie(m, "de")

		require.NotNil(t, localized)
		assert.Equal(t, "Deutscher Titel", localized.Title)
		assert.NotNil(t, localized.Overview)
		assert.Equal(t, "Deutsche Übersicht", *localized.Overview)
		assert.NotNil(t, localized.Tagline)
		assert.Equal(t, "Deutscher Slogan", *localized.Tagline)
	})

	t.Run("falls back to English when language not available", func(t *testing.T) {
		m := &movie.Movie{
			Title:            "English Title",
			Overview:         &overview,
			Tagline:          &tagline,
			OriginalTitle:    &originalTitle,
			OriginalLanguage: &originalLang,
			TitlesI18n: map[string]string{
				"en": "English Title",
				"de": "Deutscher Titel",
			},
			OverviewsI18n: map[string]string{
				"en": "English overview",
			},
			TaglinesI18n: map[string]string{
				"en": "English tagline",
			},
		}

		localized := LocalizeMovie(m, "ja") // Japanese not available

		require.NotNil(t, localized)
		assert.Equal(t, "English Title", localized.Title) // Falls back to English
		assert.Equal(t, "English overview", *localized.Overview)
		assert.Equal(t, "English tagline", *localized.Tagline)
	})

	t.Run("falls back to original title when English not available", func(t *testing.T) {
		origTitle := "Originaltitel"
		m := &movie.Movie{
			Title:            "Title",
			OriginalTitle:    &origTitle,
			OriginalLanguage: &originalLang,
			TitlesI18n: map[string]string{
				"de": "Deutscher Titel",
			},
		}

		localized := LocalizeMovie(m, "ja")

		require.NotNil(t, localized)
		assert.Equal(t, "Originaltitel", localized.Title)
	})

	t.Run("returns original movie when nil", func(t *testing.T) {
		result := LocalizeMovie(nil, "en")
		assert.Nil(t, result)
	})

	t.Run("handles empty i18n maps", func(t *testing.T) {
		m := &movie.Movie{
			Title:         "Default Title",
			Overview:      &overview,
			Tagline:       &tagline,
			OriginalTitle: &originalTitle,
		}

		localized := LocalizeMovie(m, "de")

		require.NotNil(t, localized)
		// Falls back through chain: TitlesI18n[de] → TitlesI18n[en] → OriginalTitle → Title
		// Since TitlesI18n is nil/empty, falls back to OriginalTitle
		assert.Equal(t, "Original Title", localized.Title)
		// For Overview and Tagline, falls back through OverviewsI18n → default field
		assert.Equal(t, "English overview", *localized.Overview)
		assert.Equal(t, "English tagline", *localized.Tagline)
	})
}

func TestLocalizeMovies(t *testing.T) {
	overview1 := "Overview 1"
	overview2 := "Overview 2"
	originalTitle := "Original"
	originalLang := "en"

	movies := []movie.Movie{
		{
			Title:            "Movie 1",
			Overview:         &overview1,
			OriginalTitle:    &originalTitle,
			OriginalLanguage: &originalLang,
			TitlesI18n: map[string]string{
				"en": "Movie 1",
				"de": "Film 1",
			},
			OverviewsI18n: map[string]string{
				"en": "Overview 1",
				"de": "Übersicht 1",
			},
		},
		{
			Title:            "Movie 2",
			Overview:         &overview2,
			OriginalTitle:    &originalTitle,
			OriginalLanguage: &originalLang,
			TitlesI18n: map[string]string{
				"en": "Movie 2",
				"de": "Film 2",
			},
			OverviewsI18n: map[string]string{
				"en": "Overview 2",
				"de": "Übersicht 2",
			},
		},
	}

	localized := LocalizeMovies(movies, "de")

	require.Len(t, localized, 2)
	assert.Equal(t, "Film 1", localized[0].Title)
	assert.Equal(t, "Übersicht 1", *localized[0].Overview)
	assert.Equal(t, "Film 2", localized[1].Title)
	assert.Equal(t, "Übersicht 2", *localized[1].Overview)
}

func TestLocalizeContinueWatchingItem(t *testing.T) {
	overview := "English overview"
	tagline := "English tagline"
	originalTitle := "Original"
	originalLang := "en"

	t.Run("localizes continue watching item", func(t *testing.T) {
		item := &movie.ContinueWatchingItem{
			Movie: movie.Movie{
				Title:            "English Title",
				Overview:         &overview,
				Tagline:          &tagline,
				OriginalTitle:    &originalTitle,
				OriginalLanguage: &originalLang,
				TitlesI18n: map[string]string{
					"en": "English Title",
					"de": "Deutscher Titel",
				},
				OverviewsI18n: map[string]string{
					"en": "English overview",
					"de": "Deutsche Übersicht",
				},
				TaglinesI18n: map[string]string{
					"en": "English tagline",
					"de": "Deutscher Slogan",
				},
			},
			ProgressSeconds: 1800,
			DurationSeconds: 7200,
		}

		localized := LocalizeContinueWatchingItem(item, "de")

		require.NotNil(t, localized)
		assert.Equal(t, "Deutscher Titel", localized.Title)
		assert.Equal(t, "Deutsche Übersicht", *localized.Overview)
		assert.Equal(t, "Deutscher Slogan", *localized.Tagline)
		assert.Equal(t, int32(1800), localized.ProgressSeconds)
		assert.Equal(t, int32(7200), localized.DurationSeconds)
	})

	t.Run("returns nil when item is nil", func(t *testing.T) {
		result := LocalizeContinueWatchingItem(nil, "en")
		assert.Nil(t, result)
	})
}

func TestLocalizeWatchedMovieItem(t *testing.T) {
	overview := "English overview"
	tagline := "English tagline"
	originalTitle := "Original"
	originalLang := "en"

	t.Run("localizes watched movie item", func(t *testing.T) {
		item := &movie.WatchedMovieItem{
			Movie: movie.Movie{
				Title:            "English Title",
				Overview:         &overview,
				Tagline:          &tagline,
				OriginalTitle:    &originalTitle,
				OriginalLanguage: &originalLang,
				TitlesI18n: map[string]string{
					"en": "English Title",
					"fr": "Titre Français",
				},
				OverviewsI18n: map[string]string{
					"en": "English overview",
					"fr": "Aperçu français",
				},
				TaglinesI18n: map[string]string{
					"en": "English tagline",
					"fr": "Slogan français",
				},
			},
			WatchCount: 3,
		}

		localized := LocalizeWatchedMovieItem(item, "fr")

		require.NotNil(t, localized)
		assert.Equal(t, "Titre Français", localized.Title)
		assert.Equal(t, "Aperçu français", *localized.Overview)
		assert.Equal(t, "Slogan français", *localized.Tagline)
		assert.Equal(t, int32(3), localized.WatchCount)
	})

	t.Run("returns nil when item is nil", func(t *testing.T) {
		result := LocalizeWatchedMovieItem(nil, "en")
		assert.Nil(t, result)
	})
}
