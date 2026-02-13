package mal

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
		{"full date", "2017-10-23", false, time.Date(2017, 10, 23, 0, 0, 0, 0, time.UTC)},
		{"year-month", "2017-10", false, time.Date(2017, 10, 1, 0, 0, 0, 0, time.UTC)},
		{"year only", "2017", false, time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"invalid", "not-a-date", true, time.Time{}},
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

func TestParseYear(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"full date", "2017-10-23", 2017},
		{"year only", "2017", 2017},
		{"too short", "20", 0},
		{"empty", "", 0},
		{"invalid", "abcd", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseYear(tt.input))
		})
	}
}

func TestMapStatus(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"finished_airing", "Ended"},
		{"currently_airing", "Returning Series"},
		{"not_yet_aired", "Planned"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapStatus(tt.input))
		})
	}
}

func TestMapMediaType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"tv", "Scripted"},
		{"movie", "Movie"},
		{"ova", "OVA"},
		{"ona", "ONA"},
		{"special", "Special"},
		{"music", "Music"},
		{"other", "OTHER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapMediaType(tt.input))
		})
	}
}

func TestMapRating(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"g", "G"},
		{"pg", "PG"},
		{"pg_13", "PG-13"},
		{"r", "R"},
		{"r+", "R+"},
		{"rx", "Rx"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapRating(tt.input))
		})
	}
}

func TestMalGenreToStandardID(t *testing.T) {
	tests := []struct {
		name  string
		malID int
		want  int
	}{
		{"Action", 1, 10759},
		{"Adventure", 2, 10759},
		{"Comedy", 4, 35},
		{"Drama", 8, 18},
		{"Fantasy", 10, 10765},
		{"Horror", 14, 27},
		{"Mystery", 7, 9648},
		{"Romance", 22, 10749},
		{"Sci-Fi", 24, 10765},
		{"unmapped - pass through", 99, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, malGenreToStandardID(tt.malID))
		})
	}
}

func TestMapAnimeToTVShowSearchResult(t *testing.T) {
	t.Run("full anime", func(t *testing.T) {
		anime := Anime{
			ID:    1,
			Title: "Cowboy Bebop",
			AlternativeTitles: AlternativeTitles{
				En: "Cowboy Bebop",
				Ja: "カウボーイビバップ",
			},
			Synopsis:        "A bounty hunter story in space.",
			Mean:            new(8.78),
			NumScoringUsers: 100000,
			Popularity:      new(42),
			StartDate:       "1998-04-03",
			StartSeason:     &Season{Year: 1998, Season: "spring"},
			NSFW:            "white",
			MainPicture:     &Picture{Large: "https://example.com/large.jpg", Medium: "https://example.com/medium.jpg"},
			Genres:          []Genre{{ID: 1, Name: "Action"}, {ID: 8, Name: "Drama"}},
		}
		result := mapAnimeToTVShowSearchResult(anime)
		assert.Equal(t, "1", result.ProviderID)
		assert.Equal(t, metadata.ProviderMAL, result.Provider)
		assert.Equal(t, "Cowboy Bebop", result.Name)
		assert.Equal(t, "カウボーイビバップ", result.OriginalName)
		assert.Equal(t, "ja", result.OriginalLanguage)
		assert.Equal(t, []string{"JP"}, result.OriginCountries)
		assert.Equal(t, "A bounty hunter story in space.", result.Overview)
		assert.InDelta(t, 8.78, result.VoteAverage, 0.01)
		assert.Equal(t, 100000, result.VoteCount)
		assert.InDelta(t, 42.0, result.Popularity, 0.01)
		require.NotNil(t, result.FirstAirDate)
		require.NotNil(t, result.Year)
		assert.Equal(t, 1998, *result.Year)
		require.NotNil(t, result.PosterPath)
		assert.Equal(t, "https://example.com/large.jpg", *result.PosterPath)
		assert.False(t, result.Adult)
		assert.Equal(t, []int{10759, 18}, result.GenreIDs)
	})

	t.Run("NSFW black is adult", func(t *testing.T) {
		anime := Anime{ID: 2, Title: "Test", NSFW: "black"}
		result := mapAnimeToTVShowSearchResult(anime)
		assert.True(t, result.Adult)
	})

	t.Run("no mean score", func(t *testing.T) {
		anime := Anime{ID: 3, Title: "Test", NumListUsers: 5000}
		result := mapAnimeToTVShowSearchResult(anime)
		assert.InDelta(t, 0.0, result.VoteAverage, 0.01)
		assert.InDelta(t, 5000.0, result.Popularity, 0.01)
	})

	t.Run("year from startDate when no startSeason", func(t *testing.T) {
		anime := Anime{ID: 4, Title: "Test", StartDate: "2020-05-10"}
		result := mapAnimeToTVShowSearchResult(anime)
		require.NotNil(t, result.Year)
		assert.Equal(t, 2020, *result.Year)
	})

	t.Run("poster fallback to medium", func(t *testing.T) {
		anime := Anime{ID: 5, Title: "Test", MainPicture: &Picture{Medium: "https://example.com/med.jpg"}}
		result := mapAnimeToTVShowSearchResult(anime)
		require.NotNil(t, result.PosterPath)
		assert.Equal(t, "https://example.com/med.jpg", *result.PosterPath)
	})
}

func TestMapAnimeToTVShowMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapAnimeToTVShowMetadata(nil))
	})

	t.Run("full anime", func(t *testing.T) {
		anime := &Anime{
			ID:    1,
			Title: "Attack on Titan",
			AlternativeTitles: AlternativeTitles{
				En: "Attack on Titan",
				Ja: "進撃の巨人",
			},
			Synopsis:               "Humanity fights for survival.",
			Status:                 "finished_airing",
			MediaType:              "tv",
			Mean:                   new(8.53),
			NumEpisodes:            25,
			AverageEpisodeDuration: 1440, // 24 minutes in seconds
			Popularity:             new(5),
			NSFW:                   "gray",
			StartDate:              "2013-04-07",
			EndDate:                "2013-09-29",
			MainPicture:            &Picture{Large: "https://example.com/aot.jpg"},
			Genres:                 []Genre{{ID: 1, Name: "Action"}, {ID: 10, Name: "Fantasy"}},
			Studios:                []Studio{{ID: 858, Name: "Wit Studio"}},
		}
		result := mapAnimeToTVShowMetadata(anime)
		require.NotNil(t, result)
		assert.Equal(t, "1", result.ProviderID)
		assert.Equal(t, metadata.ProviderMAL, result.Provider)
		assert.Equal(t, "Attack on Titan", result.Name)
		assert.Equal(t, "進撃の巨人", result.OriginalName)
		require.NotNil(t, result.Overview)
		assert.Equal(t, "Ended", result.Status)
		assert.Equal(t, "Scripted", result.Type)
		assert.False(t, result.InProduction)
		assert.Equal(t, 25, result.NumberOfEpisodes)
		assert.Equal(t, 1, result.NumberOfSeasons)
		assert.Equal(t, []int{24}, result.EpisodeRuntime)
		assert.InDelta(t, 8.53, result.VoteAverage, 0.01)
		require.NotNil(t, result.PosterPath)
		require.NotNil(t, result.Homepage)
		assert.Contains(t, *result.Homepage, "myanimelist.net/anime/1")
		require.Len(t, result.Genres, 2)
		require.Len(t, result.Networks, 1)
		assert.Equal(t, "Wit Studio", result.Networks[0].Name)
		assert.Equal(t, "JP", result.Networks[0].OriginCountry)
		require.Len(t, result.ExternalRatings, 1)
		assert.Equal(t, "MyAnimeList", result.ExternalRatings[0].Source)
		assert.False(t, result.Adult)
	})

	t.Run("currently airing", func(t *testing.T) {
		anime := &Anime{ID: 2, Title: "Ongoing", Status: "currently_airing", NumEpisodes: 12}
		result := mapAnimeToTVShowMetadata(anime)
		require.NotNil(t, result)
		assert.True(t, result.InProduction)
		assert.Equal(t, "Returning Series", result.Status)
	})
}

func TestMapImages(t *testing.T) {
	t.Run("nil anime", func(t *testing.T) {
		assert.Nil(t, mapImages(nil))
	})

	t.Run("no pictures", func(t *testing.T) {
		anime := &Anime{ID: 1}
		assert.Nil(t, mapImages(anime))
	})

	t.Run("with main picture and additional", func(t *testing.T) {
		anime := &Anime{
			ID:          1,
			MainPicture: &Picture{Large: "https://example.com/large.jpg", Medium: "https://example.com/med.jpg"},
			Pictures:    []Picture{{Large: "https://example.com/pic1.jpg"}, {Medium: "https://example.com/pic2.jpg"}},
		}
		result := mapImages(anime)
		require.NotNil(t, result)
		// Main (large + medium) + additional (2)
		assert.Len(t, result.Posters, 4)
	})
}
