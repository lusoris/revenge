package tvmaze

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain text", "Hello world", "Hello world"},
		{"with tags", "<p>Hello <b>world</b></p>", "Hello world"},
		{"empty", "", ""},
		{"self-closing", "Text <br/> more", "Text  more"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, stripHTML(tt.input))
		})
	}
}

func TestParseDate(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, parseDate(nil))
	})
	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, parseDate(new("")))
	})
	t.Run("valid", func(t *testing.T) {
		result := parseDate(new("2020-01-15"))
		require.NotNil(t, result)
		assert.Equal(t, time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC), *result)
	})
	t.Run("invalid", func(t *testing.T) {
		assert.Nil(t, parseDate(new("bad")))
	})
}

func TestParseAirdate(t *testing.T) {
	assert.Nil(t, parseAirdate(""))
	assert.Nil(t, parseAirdate("bad"))

	result := parseAirdate("2020-01-15")
	require.NotNil(t, result)
	assert.Equal(t, time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC), *result)
}

func TestMapGender(t *testing.T) {
	assert.Equal(t, 0, mapGender(nil))
	assert.Equal(t, 1, mapGender(new("Female")))
	assert.Equal(t, 2, mapGender(new("Male")))
	assert.Equal(t, 0, mapGender(new("Unknown")))
}

func TestMapDepartment(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Creator", "Writing"},
		{"Developer", "Writing"},
		{"Executive Producer", "Production"},
		{"Producer", "Production"},
		{"Co-Executive Producer", "Production"},
		{"Director", "Directing"},
		{"Other", "Production"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapDepartment(tt.input))
		})
	}
}

func TestMapOriginCountries(t *testing.T) {
	t.Run("from network", func(t *testing.T) {
		show := Show{Network: &Network{Country: &Country{Code: "US"}}}
		assert.Equal(t, []string{"US"}, mapOriginCountries(show))
	})
	t.Run("from web channel", func(t *testing.T) {
		show := Show{WebChannel: &Network{Country: &Country{Code: "GB"}}}
		assert.Equal(t, []string{"GB"}, mapOriginCountries(show))
	})
	t.Run("no country", func(t *testing.T) {
		assert.Nil(t, mapOriginCountries(Show{}))
	})
}

func TestGenreNameToID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"Action", 10759},
		{"Comedy", 35},
		{"Drama", 18},
		{"Science-Fiction", 10765},
		{"Horror", 27},
		{"Unknown", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, genreNameToID(tt.name))
		})
	}
}

func TestMapShowToTVShowSearchResult(t *testing.T) {
	show := Show{
		ID:        1,
		Name:      "Breaking Bad",
		Language:  "English",
		Status:    "Ended",
		Genres:    []string{"Drama", "Crime"},
		Premiered: new("2008-01-20"),
		Summary:   new("<p>A chemistry teacher turns to meth.</p>"),
		Rating:    Rating{Average: new(9.2)},
		Image:     &ImageSet{Original: "https://example.com/bb.jpg"},
		Network:   &Network{ID: 1, Name: "AMC", Country: &Country{Code: "US"}},
	}
	result := mapShowToTVShowSearchResult(show)
	assert.Equal(t, "1", result.ProviderID)
	assert.Equal(t, metadata.ProviderTVmaze, result.Provider)
	assert.Equal(t, "Breaking Bad", result.Name)
	assert.Equal(t, "English", result.OriginalLanguage)
	assert.Equal(t, []string{"US"}, result.OriginCountries)
	assert.Equal(t, "A chemistry teacher turns to meth.", result.Overview)
	assert.InDelta(t, 9.2, result.VoteAverage, 0.01)
	require.NotNil(t, result.FirstAirDate)
	require.NotNil(t, result.Year)
	assert.Equal(t, 2008, *result.Year)
	require.NotNil(t, result.PosterPath)
	assert.Equal(t, []int{18, 80}, result.GenreIDs)
}

func TestMapShowToTVShowMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapShowToTVShowMetadata(nil))
	})

	t.Run("full show", func(t *testing.T) {
		show := &Show{
			ID:           1,
			Name:         "Breaking Bad",
			Language:     "English",
			Status:       "Running",
			Type:         "Scripted",
			Summary:      new("<p>Summary here.</p>"),
			Premiered:    new("2008-01-20"),
			Ended:        new("2013-09-29"),
			Runtime:      new(60),
			Rating:       Rating{Average: new(9.2)},
			Image:        &ImageSet{Original: "https://example.com/bb.jpg"},
			OfficialSite: new("https://www.amc.com/shows/breaking-bad"),
			Genres:       []string{"Drama"},
			Network:      &Network{ID: 1, Name: "AMC", Country: &Country{Code: "US"}},
			Externals: Externals{
				IMDb:   new("tt0903747"),
				TVDb:   new(81189),
				TVRage: new(18164),
			},
		}
		result := mapShowToTVShowMetadata(show)
		require.NotNil(t, result)
		assert.Equal(t, "1", result.ProviderID)
		assert.Equal(t, metadata.ProviderTVmaze, result.Provider)
		assert.Equal(t, "Breaking Bad", result.Name)
		assert.Equal(t, "Scripted", result.Type)
		assert.Equal(t, "Running", result.Status)
		assert.True(t, result.InProduction)
		require.NotNil(t, result.IMDbID)
		assert.Equal(t, "tt0903747", *result.IMDbID)
		require.NotNil(t, result.TVDbID)
		assert.Equal(t, int32(81189), *result.TVDbID)
		require.NotNil(t, result.Overview)
		assert.Equal(t, "Summary here.", *result.Overview)
		assert.Equal(t, []int{60}, result.EpisodeRuntime)
		assert.InDelta(t, 9.2, result.VoteAverage, 0.01)
		require.NotNil(t, result.PosterPath)
		require.NotNil(t, result.Homepage)
		require.Len(t, result.Genres, 1)
		assert.Equal(t, "AMC", result.Networks[0].Name)
		assert.Equal(t, "US", result.Networks[0].OriginCountry)
		require.Len(t, result.ExternalRatings, 1)
		assert.Equal(t, "TVmaze", result.ExternalRatings[0].Source)
	})

	t.Run("web channel fallback", func(t *testing.T) {
		show := &Show{
			ID:         2,
			Name:       "Test",
			WebChannel: &Network{ID: 2, Name: "Netflix", Country: &Country{Code: "US"}},
		}
		result := mapShowToTVShowMetadata(show)
		require.NotNil(t, result)
		require.Len(t, result.Networks, 1)
		assert.Equal(t, "Netflix", result.Networks[0].Name)
	})

	t.Run("averageRuntime fallback", func(t *testing.T) {
		show := &Show{ID: 3, Name: "Test", AverageRuntime: new(45)}
		result := mapShowToTVShowMetadata(show)
		require.NotNil(t, result)
		assert.Equal(t, []int{45}, result.EpisodeRuntime)
	})
}

func TestMapSeasons(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, mapSeasons(nil))
	})

	t.Run("with seasons", func(t *testing.T) {
		seasons := []Season{
			{
				ID:           1,
				Number:       1,
				Name:         "Season 1",
				EpisodeOrder: new(7),
				PremiereDate: new("2008-01-20"),
				Summary:      new("<p>First season.</p>"),
				Image:        &ImageSet{Original: "https://example.com/s1.jpg"},
			},
			{
				ID:     2,
				Number: 2,
				Name:   "Season 2",
			},
		}
		result := mapSeasons(seasons)
		require.Len(t, result, 2)
		assert.Equal(t, "1", result[0].ProviderID)
		assert.Equal(t, 1, result[0].SeasonNumber)
		assert.Equal(t, "Season 1", result[0].Name)
		assert.Equal(t, 7, result[0].EpisodeCount)
		require.NotNil(t, result[0].AirDate)
		require.NotNil(t, result[0].Overview)
		assert.Equal(t, "First season.", *result[0].Overview)
		require.NotNil(t, result[0].PosterPath)
		assert.Equal(t, 0, result[1].EpisodeCount)
	})
}

func TestMapEpisodes(t *testing.T) {
	episodes := []Episode{
		{
			ID:      1,
			Name:    "Pilot",
			Season:  1,
			Number:  new(1),
			Airdate: "2008-01-20",
			Runtime: new(58),
			Rating:  Rating{Average: new(9.0)},
			Summary: new("<b>The beginning.</b>"),
			Image:   &ImageSet{Original: "https://example.com/ep1.jpg"},
		},
		{
			ID:     2,
			Name:   "S2E1",
			Season: 2,
			Number: new(1),
		},
	}
	result := mapEpisodes(episodes, 1)
	require.Len(t, result, 1)
	assert.Equal(t, "1", result[0].ProviderID)
	assert.Equal(t, "Pilot", result[0].Name)
	assert.Equal(t, 1, result[0].EpisodeNumber)
	require.NotNil(t, result[0].Runtime)
	assert.Equal(t, int32(58), *result[0].Runtime)
	assert.InDelta(t, 9.0, result[0].VoteAverage, 0.01)
	require.NotNil(t, result[0].AirDate)
	require.NotNil(t, result[0].Overview)
	assert.Equal(t, "The beginning.", *result[0].Overview)
	require.NotNil(t, result[0].StillPath)
}

func TestMapCast(t *testing.T) {
	cast := []CastMember{
		{
			Person:    Person{ID: 1, Name: "Bryan Cranston", Gender: new("Male"), Image: &ImageSet{Original: "https://example.com/bc.jpg"}},
			Character: Character{ID: 10, Name: "Walter White"},
		},
	}
	crew := []CrewMember{
		{
			Type:   "Creator",
			Person: Person{ID: 2, Name: "Vince Gilligan", Gender: new("Male")},
		},
	}
	result := mapCast(cast, crew)
	require.NotNil(t, result)
	require.Len(t, result.Cast, 1)
	assert.Equal(t, "1", result.Cast[0].ProviderID)
	assert.Equal(t, "Bryan Cranston", result.Cast[0].Name)
	assert.Equal(t, "Walter White", result.Cast[0].Character)
	assert.Equal(t, 0, result.Cast[0].Order)
	assert.Equal(t, 2, result.Cast[0].Gender)
	require.NotNil(t, result.Cast[0].ProfilePath)
	require.Len(t, result.Crew, 1)
	assert.Equal(t, "Vince Gilligan", result.Crew[0].Name)
	assert.Equal(t, "Creator", result.Crew[0].Job)
	assert.Equal(t, "Writing", result.Crew[0].Department)
}

func TestMapImages(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, mapImages(nil))
	})

	t.Run("with images", func(t *testing.T) {
		imgs := []ShowImage{
			{
				ID:   1,
				Type: "poster",
				Resolutions: struct {
					Original *ImageResolution `json:"original"`
					Medium   *ImageResolution `json:"medium"`
				}{
					Original: &ImageResolution{URL: "https://example.com/poster.jpg", Width: 680, Height: 1000},
				},
			},
			{
				ID:   2,
				Type: "background",
				Resolutions: struct {
					Original *ImageResolution `json:"original"`
					Medium   *ImageResolution `json:"medium"`
				}{
					Original: &ImageResolution{URL: "https://example.com/bg.jpg", Width: 1920, Height: 1080},
				},
			},
			{
				ID:   3,
				Type: "typography",
				Resolutions: struct {
					Original *ImageResolution `json:"original"`
					Medium   *ImageResolution `json:"medium"`
				}{
					Original: &ImageResolution{URL: "https://example.com/logo.jpg", Width: 400, Height: 100},
				},
			},
			{
				ID:   4,
				Type: "poster",
				Resolutions: struct {
					Original *ImageResolution `json:"original"`
					Medium   *ImageResolution `json:"medium"`
				}{
					Original: nil, // nil original should be filtered
				},
			},
		}
		result := mapImages(imgs)
		require.NotNil(t, result)
		assert.Len(t, result.Posters, 1)
		assert.Len(t, result.Backdrops, 1)
		assert.Len(t, result.Logos, 1)
		assert.InDelta(t, 0.68, result.Posters[0].AspectRatio, 0.01)
	})
}

func TestMapExternalIDs(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapExternalIDs(nil))
	})

	t.Run("with all IDs", func(t *testing.T) {
		show := &Show{
			Externals: Externals{
				IMDb:   new("tt0903747"),
				TVDb:   new(81189),
				TVRage: new(18164),
			},
		}
		result := mapExternalIDs(show)
		require.NotNil(t, result)
		require.NotNil(t, result.IMDbID)
		assert.Equal(t, "tt0903747", *result.IMDbID)
		require.NotNil(t, result.TVDbID)
		assert.Equal(t, int32(81189), *result.TVDbID)
		require.NotNil(t, result.TVRageID)
	})

	t.Run("empty IDs", func(t *testing.T) {
		show := &Show{}
		result := mapExternalIDs(show)
		require.NotNil(t, result)
		assert.Nil(t, result.IMDbID)
		assert.Nil(t, result.TVDbID)
		assert.Nil(t, result.TVRageID)
	})
}
