package tvshow

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSeries_GetTitle(t *testing.T) {
	tests := []struct {
		name     string
		series   *Series
		lang     string
		expected string
	}{
		{
			name: "Returns requested language",
			series: &Series{
				Title:            "Breaking Bad",
				OriginalLanguage: "en",
				TitlesI18n: map[string]string{
					"en": "Breaking Bad",
					"de": "Breaking Bad (DE)",
					"es": "Breaking Bad (ES)",
				},
			},
			lang:     "de",
			expected: "Breaking Bad (DE)",
		},
		{
			name: "Falls back to English",
			series: &Series{
				Title:            "Breaking Bad",
				OriginalLanguage: "en",
				TitlesI18n: map[string]string{
					"en": "Breaking Bad",
				},
			},
			lang:     "fr",
			expected: "Breaking Bad",
		},
		{
			name: "Falls back to original language",
			series: &Series{
				Title:            "Dark",
				OriginalLanguage: "de",
				TitlesI18n: map[string]string{
					"de": "Dark (Original)",
				},
			},
			lang:     "fr",
			expected: "Dark (Original)",
		},
		{
			name: "Falls back to OriginalTitle",
			series: &Series{
				Title:            "Dark",
				OriginalTitle:    strPtr("Dark Original"),
				OriginalLanguage: "de",
				TitlesI18n:       nil,
			},
			lang:     "en",
			expected: "Dark Original",
		},
		{
			name: "Falls back to default Title",
			series: &Series{
				Title:            "Breaking Bad",
				OriginalLanguage: "en",
				TitlesI18n:       nil,
			},
			lang:     "de",
			expected: "Breaking Bad",
		},
		{
			name: "Handles nil TitlesI18n",
			series: &Series{
				Title:            "Test Series",
				OriginalLanguage: "en",
			},
			lang:     "de",
			expected: "Test Series",
		},
		{
			name: "Skips empty strings in i18n",
			series: &Series{
				Title:            "Breaking Bad",
				OriginalLanguage: "en",
				TitlesI18n: map[string]string{
					"de": "",
					"en": "Breaking Bad EN",
				},
			},
			lang:     "de",
			expected: "Breaking Bad EN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.series.GetTitle(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSeries_GetTagline(t *testing.T) {
	tests := []struct {
		name     string
		series   *Series
		lang     string
		expected string
	}{
		{
			name: "Returns requested language",
			series: &Series{
				Tagline:          strPtr("All Hail the King"),
				OriginalLanguage: "en",
				TaglinesI18n: map[string]string{
					"en": "All Hail the King",
					"de": "Heil dem König",
				},
			},
			lang:     "de",
			expected: "Heil dem König",
		},
		{
			name: "Falls back to English",
			series: &Series{
				Tagline:          strPtr("All Hail the King"),
				OriginalLanguage: "en",
				TaglinesI18n: map[string]string{
					"en": "All Hail the King",
				},
			},
			lang:     "fr",
			expected: "All Hail the King",
		},
		{
			name: "Falls back to default Tagline",
			series: &Series{
				Tagline:          strPtr("Default Tagline"),
				OriginalLanguage: "en",
			},
			lang:     "de",
			expected: "Default Tagline",
		},
		{
			name: "Returns empty for nil tagline",
			series: &Series{
				OriginalLanguage: "en",
			},
			lang:     "de",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.series.GetTagline(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSeries_GetOverview(t *testing.T) {
	tests := []struct {
		name     string
		series   *Series
		lang     string
		expected string
	}{
		{
			name: "Returns requested language",
			series: &Series{
				Overview:         strPtr("A chemistry teacher..."),
				OriginalLanguage: "en",
				OverviewsI18n: map[string]string{
					"en": "A chemistry teacher...",
					"de": "Ein Chemielehrer...",
				},
			},
			lang:     "de",
			expected: "Ein Chemielehrer...",
		},
		{
			name: "Falls back to English",
			series: &Series{
				Overview:         strPtr("A chemistry teacher..."),
				OriginalLanguage: "en",
				OverviewsI18n: map[string]string{
					"en": "A chemistry teacher...",
				},
			},
			lang:     "fr",
			expected: "A chemistry teacher...",
		},
		{
			name: "Falls back to default Overview",
			series: &Series{
				Overview:         strPtr("Default Overview"),
				OriginalLanguage: "en",
			},
			lang:     "de",
			expected: "Default Overview",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.series.GetOverview(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSeries_GetAgeRating(t *testing.T) {
	series := &Series{
		AgeRatings: map[string]map[string]string{
			"US": {"TV": "TV-MA"},
			"DE": {"FSK": "16"},
			"GB": {"BBFC": "15"},
		},
	}

	tests := []struct {
		country  string
		system   string
		expected string
	}{
		{"US", "TV", "TV-MA"},
		{"DE", "FSK", "16"},
		{"GB", "BBFC", "15"},
		{"FR", "CSA", ""},  // Country not found
		{"US", "FSK", ""},  // System not found
	}

	for _, tt := range tests {
		t.Run(tt.country+"_"+tt.system, func(t *testing.T) {
			result := series.GetAgeRating(tt.country, tt.system)
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("Handles nil AgeRatings", func(t *testing.T) {
		s := &Series{}
		assert.Equal(t, "", s.GetAgeRating("US", "TV"))
	})
}

func TestSeries_GetAvailableLanguages(t *testing.T) {
	t.Run("Returns all languages", func(t *testing.T) {
		series := &Series{
			TitlesI18n: map[string]string{
				"en": "Breaking Bad",
				"de": "Breaking Bad",
				"es": "Breaking Bad",
			},
		}
		langs := series.GetAvailableLanguages()
		assert.Len(t, langs, 3)
		assert.Contains(t, langs, "en")
		assert.Contains(t, langs, "de")
		assert.Contains(t, langs, "es")
	})

	t.Run("Handles nil TitlesI18n", func(t *testing.T) {
		series := &Series{}
		langs := series.GetAvailableLanguages()
		assert.Empty(t, langs)
	})
}

func TestSeries_GetAvailableAgeRatingCountries(t *testing.T) {
	t.Run("Returns all countries", func(t *testing.T) {
		series := &Series{
			AgeRatings: map[string]map[string]string{
				"US": {"TV": "TV-MA"},
				"DE": {"FSK": "16"},
			},
		}
		countries := series.GetAvailableAgeRatingCountries()
		assert.Len(t, countries, 2)
		assert.Contains(t, countries, "US")
		assert.Contains(t, countries, "DE")
	})

	t.Run("Handles nil AgeRatings", func(t *testing.T) {
		series := &Series{}
		countries := series.GetAvailableAgeRatingCountries()
		assert.Empty(t, countries)
	})
}

func TestSeries_IsEnded(t *testing.T) {
	tests := []struct {
		name     string
		status   *string
		expected bool
	}{
		{"Ended status", strPtr("Ended"), true},
		{"Canceled status", strPtr("Canceled"), true},
		{"Returning Series", strPtr("Returning Series"), false},
		{"In Production", strPtr("In Production"), false},
		{"Nil status", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			series := &Series{Status: tt.status}
			assert.Equal(t, tt.expected, series.IsEnded())
		})
	}
}

func TestSeason_GetName(t *testing.T) {
	tests := []struct {
		name     string
		season   *Season
		lang     string
		expected string
	}{
		{
			name: "Returns requested language",
			season: &Season{
				Name: "Season 1",
				NamesI18n: map[string]string{
					"en": "Season 1",
					"de": "Staffel 1",
				},
			},
			lang:     "de",
			expected: "Staffel 1",
		},
		{
			name: "Falls back to English",
			season: &Season{
				Name: "Season 1",
				NamesI18n: map[string]string{
					"en": "Season 1",
				},
			},
			lang:     "fr",
			expected: "Season 1",
		},
		{
			name: "Falls back to default Name",
			season: &Season{
				Name: "Season 1",
			},
			lang:     "de",
			expected: "Season 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.season.GetName(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSeason_GetOverview(t *testing.T) {
	tests := []struct {
		name     string
		season   *Season
		lang     string
		expected string
	}{
		{
			name: "Returns requested language",
			season: &Season{
				Overview: strPtr("Season overview"),
				OverviewsI18n: map[string]string{
					"en": "Season overview",
					"de": "Staffelübersicht",
				},
			},
			lang:     "de",
			expected: "Staffelübersicht",
		},
		{
			name: "Falls back to default Overview",
			season: &Season{
				Overview: strPtr("Default overview"),
			},
			lang:     "de",
			expected: "Default overview",
		},
		{
			name: "Returns empty for nil overview",
			season: &Season{},
			lang:   "de",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.season.GetOverview(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSeason_IsSpecials(t *testing.T) {
	tests := []struct {
		name         string
		seasonNumber int32
		expected     bool
	}{
		{"Season 0 is specials", 0, true},
		{"Season 1 is not specials", 1, false},
		{"Season 5 is not specials", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			season := &Season{SeasonNumber: tt.seasonNumber}
			assert.Equal(t, tt.expected, season.IsSpecials())
		})
	}
}

func TestEpisode_GetTitle(t *testing.T) {
	tests := []struct {
		name     string
		episode  *Episode
		lang     string
		expected string
	}{
		{
			name: "Returns requested language",
			episode: &Episode{
				Title: "Pilot",
				TitlesI18n: map[string]string{
					"en": "Pilot",
					"de": "Der Anfang",
				},
			},
			lang:     "de",
			expected: "Der Anfang",
		},
		{
			name: "Falls back to English",
			episode: &Episode{
				Title: "Pilot",
				TitlesI18n: map[string]string{
					"en": "Pilot",
				},
			},
			lang:     "fr",
			expected: "Pilot",
		},
		{
			name: "Falls back to default Title",
			episode: &Episode{
				Title: "Pilot",
			},
			lang:     "de",
			expected: "Pilot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.episode.GetTitle(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEpisode_EpisodeCode(t *testing.T) {
	tests := []struct {
		season   int32
		episode  int32
		expected string
	}{
		{1, 1, "S01E01"},
		{1, 10, "S01E10"},
		{10, 1, "S10E01"},
		{10, 10, "S10E10"},
		{0, 1, "S00E01"}, // Specials
		{99, 99, "S99E99"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			ep := &Episode{
				SeasonNumber:  tt.season,
				EpisodeNumber: tt.episode,
			}
			assert.Equal(t, tt.expected, ep.EpisodeCode())
		})
	}
}

func TestEpisode_HasAired(t *testing.T) {
	now := time.Now()
	pastDate := now.AddDate(0, 0, -7)
	futureDate := now.AddDate(0, 0, 7)

	tests := []struct {
		name     string
		airDate  *time.Time
		expected bool
	}{
		{"Past date has aired", &pastDate, true},
		{"Future date has not aired", &futureDate, false},
		{"Nil date has not aired", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &Episode{AirDate: tt.airDate}
			assert.Equal(t, tt.expected, ep.HasAired())
		})
	}
}

func TestCredit_IsCastAndIsCrew(t *testing.T) {
	t.Run("SeriesCredit", func(t *testing.T) {
		cast := &SeriesCredit{CreditType: "cast"}
		crew := &SeriesCredit{CreditType: "crew"}

		assert.True(t, cast.IsCast())
		assert.False(t, cast.IsCrew())
		assert.False(t, crew.IsCast())
		assert.True(t, crew.IsCrew())
	})

	t.Run("EpisodeCredit", func(t *testing.T) {
		cast := &EpisodeCredit{CreditType: "cast"}
		crew := &EpisodeCredit{CreditType: "crew"}

		assert.True(t, cast.IsCast())
		assert.False(t, cast.IsCrew())
		assert.False(t, crew.IsCast())
		assert.True(t, crew.IsCrew())
	})
}

func TestEpisodeWatched_ProgressPercent(t *testing.T) {
	tests := []struct {
		name            string
		progressSeconds int32
		durationSeconds int32
		expected        float64
	}{
		{"0% progress", 0, 3600, 0},
		{"50% progress", 1800, 3600, 50},
		{"100% progress", 3600, 3600, 100},
		{"Zero duration", 100, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watched := &EpisodeWatched{
				ProgressSeconds: tt.progressSeconds,
				DurationSeconds: tt.durationSeconds,
			}
			assert.Equal(t, tt.expected, watched.ProgressPercent())
		})
	}
}

func TestSeriesWatchStats_CompletionPercent(t *testing.T) {
	tests := []struct {
		name          string
		watchedCount  int64
		totalEpisodes int64
		expected      float64
	}{
		{"0% completion", 0, 62, 0},
		{"50% completion", 31, 62, 50},
		{"100% completion", 62, 62, 100},
		{"Zero episodes", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &SeriesWatchStats{
				WatchedCount:  tt.watchedCount,
				TotalEpisodes: tt.totalEpisodes,
			}
			assert.Equal(t, tt.expected, stats.CompletionPercent())
		})
	}
}

func TestSeriesListFilters(t *testing.T) {
	filters := SeriesListFilters{
		OrderBy: "title",
		Limit:   25,
		Offset:  50,
	}

	assert.Equal(t, "title", filters.OrderBy)
	assert.Equal(t, int32(25), filters.Limit)
	assert.Equal(t, int32(50), filters.Offset)
}

func TestEpisodeListFilters(t *testing.T) {
	seriesID := uuid.New()
	seasonID := uuid.New()

	filters := EpisodeListFilters{
		SeriesID: &seriesID,
		SeasonID: &seasonID,
		Limit:    10,
		Offset:   0,
	}

	assert.Equal(t, &seriesID, filters.SeriesID)
	assert.Equal(t, &seasonID, filters.SeasonID)
	assert.Equal(t, int32(10), filters.Limit)
	assert.Equal(t, int32(0), filters.Offset)
}

func TestContinueWatchingItem(t *testing.T) {
	series := &Series{
		ID:    uuid.New(),
		Title: "Breaking Bad",
	}

	item := ContinueWatchingItem{
		Series:            series,
		LastEpisodeID:     uuid.New(),
		LastSeasonNumber:  2,
		LastEpisodeNumber: 5,
		LastEpisodeTitle:  "Breakage",
		ProgressSeconds:   1200,
		DurationSeconds:   2700,
		LastWatchedAt:     time.Now(),
	}

	assert.Equal(t, series, item.Series)
	assert.Equal(t, int32(2), item.LastSeasonNumber)
	assert.Equal(t, int32(5), item.LastEpisodeNumber)
	assert.Equal(t, "Breakage", item.LastEpisodeTitle)
	assert.Equal(t, int32(1200), item.ProgressSeconds)
}

func TestNextEpisode(t *testing.T) {
	episode := &Episode{
		ID:            uuid.New(),
		SeasonNumber:  3,
		EpisodeNumber: 1,
		Title:         "No Más",
	}

	next := NextEpisode{
		Episode:        episode,
		IsNewSeason:   true,
		IsSeriesFinale: false,
	}

	assert.Equal(t, episode, next.Episode)
	assert.True(t, next.IsNewSeason)
	assert.False(t, next.IsSeriesFinale)
}

func TestEpisodeFile(t *testing.T) {
	file := EpisodeFile{
		ID:             uuid.New(),
		EpisodeID:      uuid.New(),
		FilePath:       "/media/tv/Breaking Bad/Season 1/Breaking.Bad.S01E01.mkv",
		FileName:       "Breaking.Bad.S01E01.mkv",
		FileSize:       1500000000,
		Container:      strPtr("mkv"),
		Resolution:     strPtr("1920x1080"),
		QualityProfile: strPtr("HDTV-1080p"),
		VideoCodec:     strPtr("h264"),
		AudioCodec:     strPtr("aac"),
		BitrateKbps:    int32Ptr(5000),
		DurationSeconds: decimalPtr(decimal.NewFromInt(2700)),
		AudioLanguages:    []string{"en"},
		SubtitleLanguages: []string{"en", "de", "es"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	assert.NotEqual(t, uuid.Nil, file.ID)
	assert.Equal(t, "Breaking.Bad.S01E01.mkv", file.FileName)
	assert.Equal(t, int64(1500000000), file.FileSize)
	assert.Equal(t, "mkv", *file.Container)
	assert.Len(t, file.SubtitleLanguages, 3)
}

func TestNetwork(t *testing.T) {
	network := Network{
		ID:            uuid.New(),
		TMDbID:        174,
		Name:          "AMC",
		LogoPath:      strPtr("/path/to/logo.png"),
		OriginCountry: strPtr("US"),
		CreatedAt:     time.Now(),
	}

	assert.Equal(t, int32(174), network.TMDbID)
	assert.Equal(t, "AMC", network.Name)
	assert.Equal(t, "US", *network.OriginCountry)
}

func TestSeriesGenre(t *testing.T) {
	genre := SeriesGenre{
		ID:          uuid.New(),
		SeriesID:    uuid.New(),
		TMDbGenreID: 18,
		Name:        "Drama",
		CreatedAt:   time.Now(),
	}

	assert.Equal(t, int32(18), genre.TMDbGenreID)
	assert.Equal(t, "Drama", genre.Name)
}

func TestUserTVStats(t *testing.T) {
	stats := UserTVStats{
		SeriesCount:        15,
		EpisodesWatched:    250,
		EpisodesInProgress: 5,
		TotalWatches:       280,
	}

	assert.Equal(t, int64(15), stats.SeriesCount)
	assert.Equal(t, int64(250), stats.EpisodesWatched)
	assert.Equal(t, int64(5), stats.EpisodesInProgress)
	assert.Equal(t, int64(280), stats.TotalWatches)
}

// Helper functions for creating pointers in tests
func strPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func decimalPtr(d decimal.Decimal) *decimal.Decimal {
	return &d
}
