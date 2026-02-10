package tvshow

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	contenttvshow "github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/service/metadata"
)

func ptr[T any](v T) *T { return &v }

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

func TestGetTVAgeRatingSystem(t *testing.T) {
	tests := []struct {
		country string
		want    string
	}{
		{"US", "TV Parental Guidelines"},
		{"DE", "FSK"},
		{"GB", "BBFC"},
		{"FR", "CSA"},
		{"JP", "EIRIN"},
		{"KR", "KMRB"},
		{"BR", "DJCTQ"},
		{"AU", "ACB"},
		{"CA", "CHVRS"},
		{"ZZ", "ZZ"},
	}
	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			assert.Equal(t, tt.want, getTVAgeRatingSystem(tt.country))
		})
	}
}

func TestMapSearchResultToSeries(t *testing.T) {
	t.Run("full search result", func(t *testing.T) {
		r := &metadata.TVShowSearchResult{
			ProviderID:       "12345",
			Provider:         metadata.ProviderTMDb,
			Name:             "Breaking Bad",
			OriginalName:     "Breaking Bad",
			OriginalLanguage: "en",
			VoteAverage:      9.5,
			VoteCount:        50000,
			Popularity:       100.5,
			PosterPath:       ptr("https://example.com/poster.jpg"),
			BackdropPath:     ptr("https://example.com/backdrop.jpg"),
			FirstAirDate:     ptr(time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)),
		}
		series := mapSearchResultToSeries(r)
		require.NotNil(t, series)
		assert.NotEqual(t, uuid.Nil, series.ID)
		assert.Equal(t, "Breaking Bad", series.Title)
		assert.Equal(t, "en", series.OriginalLanguage)
		require.NotNil(t, series.OriginalTitle)
		require.NotNil(t, series.PosterPath)
		require.NotNil(t, series.BackdropPath)
		require.NotNil(t, series.FirstAirDate)

		require.NotNil(t, series.VoteAverage)
		va, _ := series.VoteAverage.Float64()
		assert.InDelta(t, 9.5, va, 0.01)

		require.NotNil(t, series.VoteCount)
		assert.Equal(t, int32(50000), *series.VoteCount)

		require.NotNil(t, series.Popularity)

		require.NotNil(t, series.TMDbID)
		assert.Equal(t, int32(12345), *series.TMDbID)
	})

	t.Run("zero values not set", func(t *testing.T) {
		r := &metadata.TVShowSearchResult{Name: "Test"}
		series := mapSearchResultToSeries(r)
		assert.Nil(t, series.VoteAverage)
		assert.Nil(t, series.VoteCount)
		assert.Nil(t, series.Popularity)
		assert.Nil(t, series.TMDbID)
		assert.Nil(t, series.OriginalTitle)
	})

	t.Run("non-numeric provider ID", func(t *testing.T) {
		r := &metadata.TVShowSearchResult{ProviderID: "abc", Name: "Test"}
		series := mapSearchResultToSeries(r)
		assert.Nil(t, series.TMDbID)
	})
}

func TestMapMetadataToSeries(t *testing.T) {
	t.Run("full metadata", func(t *testing.T) {
		series := &contenttvshow.Series{ID: uuid.Must(uuid.NewV7())}
		overview := "A chemistry teacher turns to cooking meth."
		tagline := "All bad things must come to an end."
		homepage := "https://www.amc.com/shows/breaking-bad"
		trailer := "https://youtube.com/watch?v=abc"
		imdb := "tt0903747"
		tmdb := int32(1396)
		tvdb := int32(81189)

		meta := &metadata.TVShowMetadata{
			Name:             "Breaking Bad",
			OriginalName:     "Breaking Bad",
			OriginalLanguage: "en",
			Overview:         &overview,
			Tagline:          &tagline,
			Status:           "Ended",
			Type:             "Scripted",
			NumberOfSeasons:  5,
			NumberOfEpisodes: 62,
			FirstAirDate:     ptr(time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)),
			LastAirDate:      ptr(time.Date(2013, 9, 29, 0, 0, 0, 0, time.UTC)),
			VoteAverage:      9.5,
			VoteCount:        50000,
			Popularity:       100.5,
			PosterPath:       ptr("https://example.com/poster.jpg"),
			BackdropPath:     ptr("https://example.com/backdrop.jpg"),
			Homepage:         &homepage,
			TrailerURL:       &trailer,
			IMDbID:           &imdb,
			TMDbID:           &tmdb,
			TVDbID:           &tvdb,
			Translations: map[string]*metadata.LocalizedTVShowData{
				"de": {Name: "Breaking Bad", Tagline: "Alle bösen Dinge", Overview: "Ein Chemielehrer..."},
				"fr": {Name: "Breaking Bad", Overview: "Un prof de chimie..."},
			},
			ExternalRatings: []metadata.ExternalRating{
				{Source: "IMDb", Value: "9.5/10", Score: 95.0},
			},
		}
		contentRatings := []metadata.ContentRating{
			{CountryCode: "US", Rating: "TV-MA"},
			{CountryCode: "DE", Rating: "16"},
			{CountryCode: "US", Rating: ""},
		}

		mapMetadataToSeries(series, meta, contentRatings)

		assert.Equal(t, "Breaking Bad", series.Title)
		require.NotNil(t, series.OriginalTitle)
		assert.Equal(t, "en", series.OriginalLanguage)
		require.NotNil(t, series.Overview)
		require.NotNil(t, series.Status)
		assert.Equal(t, "Ended", *series.Status)
		require.NotNil(t, series.Type)
		assert.Equal(t, int32(5), series.TotalSeasons)
		assert.Equal(t, int32(62), series.TotalEpisodes)
		require.NotNil(t, series.FirstAirDate)
		require.NotNil(t, series.LastAirDate)
		assert.Equal(t, &imdb, series.IMDbID)
		assert.Equal(t, &tmdb, series.TMDbID)
		assert.Equal(t, &tvdb, series.TVDbID)

		require.NotNil(t, series.VoteAverage)
		va, _ := series.VoteAverage.Float64()
		assert.InDelta(t, 9.5, va, 0.01)

		require.Len(t, series.TitlesI18n, 2)
		require.Len(t, series.TaglinesI18n, 1)
		require.Len(t, series.OverviewsI18n, 2)

		require.Len(t, series.AgeRatings, 2)
		assert.Equal(t, "TV-MA", series.AgeRatings["US"]["TV Parental Guidelines"])
		assert.Equal(t, "16", series.AgeRatings["DE"]["FSK"])

		require.Len(t, series.ExternalRatings, 1)
		assert.Equal(t, "IMDb", series.ExternalRatings[0].Source)
	})

	t.Run("empty translations and ratings", func(t *testing.T) {
		series := &contenttvshow.Series{ID: uuid.Must(uuid.NewV7())}
		meta := &metadata.TVShowMetadata{Name: "Test"}
		mapMetadataToSeries(series, meta, nil)
		assert.Nil(t, series.TitlesI18n)
		assert.Nil(t, series.AgeRatings)
		assert.Nil(t, series.ExternalRatings)
	})
}

func TestMapSeasonMetadataToSeason(t *testing.T) {
	t.Run("full season", func(t *testing.T) {
		season := &contenttvshow.Season{ID: uuid.Must(uuid.NewV7())}
		tmdbID := int32(3572)
		meta := &metadata.SeasonMetadata{
			Name:        "Season 1",
			Overview:    ptr("The first season."),
			PosterPath:  ptr("https://example.com/s1.jpg"),
			AirDate:     ptr(time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)),
			VoteAverage: 9.2,
			TMDbID:      &tmdbID,
			Episodes: []metadata.EpisodeSummary{
				{Name: "Pilot"},
				{Name: "Cat's in the Bag..."},
			},
			Translations: map[string]*metadata.LocalizedSeasonData{
				"de": {Name: "Staffel 1", Overview: "Die erste Staffel."},
			},
		}

		mapSeasonMetadataToSeason(season, meta)

		assert.Equal(t, "Season 1", season.Name)
		require.NotNil(t, season.Overview)
		require.NotNil(t, season.PosterPath)
		require.NotNil(t, season.AirDate)
		assert.Equal(t, int32(2), season.EpisodeCount)
		assert.Equal(t, &tmdbID, season.TMDbID)

		require.NotNil(t, season.VoteAverage)
		va, _ := season.VoteAverage.Float64()
		assert.InDelta(t, 9.2, va, 0.01)

		require.Len(t, season.NamesI18n, 1)
		assert.Equal(t, "Staffel 1", season.NamesI18n["de"])
	})

	t.Run("no translations", func(t *testing.T) {
		season := &contenttvshow.Season{}
		meta := &metadata.SeasonMetadata{Name: "Season 1"}
		mapSeasonMetadataToSeason(season, meta)
		assert.Nil(t, season.NamesI18n)
	})
}

func TestMapEpisodeMetadataToEpisode(t *testing.T) {
	t.Run("full episode", func(t *testing.T) {
		episode := &contenttvshow.Episode{ID: uuid.Must(uuid.NewV7())}
		tmdbID := int32(62085)
		imdbID := "tt0959621"
		prodCode := "1ABE79"
		meta := &metadata.EpisodeMetadata{
			Name:           "Pilot",
			Overview:       ptr("Walter White turns to cooking meth."),
			AirDate:        ptr(time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)),
			Runtime:        ptr(int32(58)),
			StillPath:      ptr("https://example.com/ep1.jpg"),
			ProductionCode: &prodCode,
			TMDbID:         &tmdbID,
			IMDbID:         &imdbID,
			VoteAverage:    9.0,
			VoteCount:      1500,
			Translations: map[string]*metadata.LocalizedEpisodeData{
				"de": {Name: "Pilot", Overview: "Walter White beginnt..."},
			},
		}

		mapEpisodeMetadataToEpisode(episode, meta)

		assert.Equal(t, "Pilot", episode.Title)
		require.NotNil(t, episode.Overview)
		require.NotNil(t, episode.AirDate)
		require.NotNil(t, episode.Runtime)
		assert.Equal(t, int32(58), *episode.Runtime)
		require.NotNil(t, episode.StillPath)
		require.NotNil(t, episode.ProductionCode)
		assert.Equal(t, &tmdbID, episode.TMDbID)
		assert.Equal(t, &imdbID, episode.IMDbID)

		require.NotNil(t, episode.VoteAverage)
		va, _ := episode.VoteAverage.Float64()
		assert.InDelta(t, 9.0, va, 0.01)
		require.NotNil(t, episode.VoteCount)
		assert.Equal(t, int32(1500), *episode.VoteCount)

		require.Len(t, episode.TitlesI18n, 1)
		require.Len(t, episode.OverviewsI18n, 1)
	})

	t.Run("zero ratings not set", func(t *testing.T) {
		episode := &contenttvshow.Episode{}
		meta := &metadata.EpisodeMetadata{Name: "Test"}
		mapEpisodeMetadataToEpisode(episode, meta)
		assert.Nil(t, episode.VoteAverage)
		assert.Nil(t, episode.VoteCount)
	})
}

func TestMapCreditsToSeriesCredits(t *testing.T) {
	seriesID := uuid.Must(uuid.NewV7())

	t.Run("cast and crew", func(t *testing.T) {
		credits := &metadata.Credits{
			Cast: []metadata.CastMember{
				{ProviderID: "17419", Name: "Bryan Cranston", Character: "Walter White", Order: 0, ProfilePath: ptr("https://example.com/bc.jpg")},
			},
			Crew: []metadata.CrewMember{
				{ProviderID: "66633", Name: "Vince Gilligan", Job: "Creator", Department: "Writing", ProfilePath: ptr("https://example.com/vg.jpg")},
			},
		}

		result := mapCreditsToSeriesCredits(seriesID, credits)
		require.Len(t, result, 2)

		assert.Equal(t, seriesID, result[0].SeriesID)
		assert.Equal(t, int32(17419), result[0].TMDbPersonID)
		assert.Equal(t, "Bryan Cranston", result[0].Name)
		assert.Equal(t, "cast", result[0].CreditType)
		require.NotNil(t, result[0].Character)
		assert.Equal(t, "Walter White", *result[0].Character)
		require.NotNil(t, result[0].CastOrder)
		require.NotNil(t, result[0].ProfilePath)

		assert.Equal(t, "crew", result[1].CreditType)
		assert.Equal(t, int32(66633), result[1].TMDbPersonID)
		require.NotNil(t, result[1].Job)
		assert.Equal(t, "Creator", *result[1].Job)
		require.NotNil(t, result[1].Department)
	})

	t.Run("empty credits", func(t *testing.T) {
		result := mapCreditsToSeriesCredits(seriesID, &metadata.Credits{})
		assert.Nil(t, result)
	})

	t.Run("non-numeric provider ID", func(t *testing.T) {
		credits := &metadata.Credits{
			Cast: []metadata.CastMember{{ProviderID: "not-a-number", Name: "Test"}},
		}
		result := mapCreditsToSeriesCredits(seriesID, credits)
		require.Len(t, result, 1)
		assert.Equal(t, int32(0), result[0].TMDbPersonID)
	})
}
