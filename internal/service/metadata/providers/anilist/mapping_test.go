package anilist

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBestTitle(t *testing.T) {
	tests := []struct {
		name  string
		title MediaTitle
		want  string
	}{
		{"english first", MediaTitle{English: new("Attack on Titan"), Romaji: new("Shingeki no Kyojin")}, "Attack on Titan"},
		{"romaji fallback", MediaTitle{Romaji: new("Shingeki no Kyojin")}, "Shingeki no Kyojin"},
		{"userPreferred fallback", MediaTitle{UserPreferred: new("AoT")}, "AoT"},
		{"native fallback", MediaTitle{Native: new("進撃の巨人")}, "進撃の巨人"},
		{"empty", MediaTitle{}, ""},
		{"empty english", MediaTitle{English: new(""), Romaji: new("Test")}, "Test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, bestTitle(tt.title))
		})
	}
}

func TestBestCoverImage(t *testing.T) {
	tests := []struct {
		name string
		ci   CoverImage
		want string
	}{
		{"extraLarge first", CoverImage{ExtraLarge: new("xl"), Large: new("l")}, "xl"},
		{"large fallback", CoverImage{Large: new("l")}, "l"},
		{"medium fallback", CoverImage{Medium: new("m")}, "m"},
		{"empty", CoverImage{}, ""},
		{"empty extraLarge", CoverImage{ExtraLarge: new(""), Large: new("l")}, "l"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, bestCoverImage(tt.ci))
		})
	}
}

func TestFuzzyDateToTime(t *testing.T) {
	t.Run("nil year", func(t *testing.T) {
		assert.Nil(t, fuzzyDateToTime(FuzzyDate{}))
	})

	t.Run("year only", func(t *testing.T) {
		result := fuzzyDateToTime(FuzzyDate{Year: new(2020)})
		require.NotNil(t, result)
		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), *result)
	})

	t.Run("year and month", func(t *testing.T) {
		result := fuzzyDateToTime(FuzzyDate{Year: new(2020), Month: new(6)})
		require.NotNil(t, result)
		assert.Equal(t, time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC), *result)
	})

	t.Run("full date", func(t *testing.T) {
		result := fuzzyDateToTime(FuzzyDate{Year: new(2020), Month: new(6), Day: new(15)})
		require.NotNil(t, result)
		assert.Equal(t, time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC), *result)
	})
}

func TestMapStatus(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"FINISHED", "Ended"},
		{"RELEASING", "Returning Series"},
		{"NOT_YET_RELEASED", "Planned"},
		{"CANCELLED", "Canceled"},
		{"HIATUS", "Returning Series"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapStatus(tt.input))
		})
	}
}

func TestMapFormat(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"TV", "Scripted"},
		{"TV_SHORT", "Scripted"},
		{"MOVIE", "Movie"},
		{"SPECIAL", "Special"},
		{"OVA", "OVA"},
		{"ONA", "ONA"},
		{"MUSIC", "Music"},
		{"OTHER", "OTHER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapFormat(tt.input))
		})
	}
}

func TestMapGender(t *testing.T) {
	assert.Equal(t, 0, mapGender(nil))
	assert.Equal(t, 1, mapGender(new("Female")))
	assert.Equal(t, 2, mapGender(new("Male")))
	assert.Equal(t, 3, mapGender(new("Non-binary")))
	assert.Equal(t, 0, mapGender(new("Unknown")))
}

func TestMapDepartment(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Director", "Directing"},
		{"Chief Director", "Directing"},
		{"Producer", "Production"},
		{"Writer", "Writing"},
		{"Script", "Writing"},
		{"Series Composition", "Writing"},
		{"Story Creator", "Writing"},
		{"Original Creator", "Writing"},
		{"Music", "Sound"},
		{"Sound Director", "Directing"},
		{"Character Design", "Art"},
		{"Art Director", "Directing"},
		{"Animation Director", "Directing"},
		{"Something Else", "Production"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapDepartment(tt.input))
		})
	}
}

func TestSafeStr(t *testing.T) {
	assert.Equal(t, "", safeStr(nil))
	assert.Equal(t, "hello", safeStr(new("hello")))
}

func TestExtractIMDbID(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"standard URL", "https://www.imdb.com/title/tt1234567/", "tt1234567"},
		{"no path", "https://www.imdb.com/", ""},
		{"different prefix", "https://www.imdb.com/name/nm0001234/", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, extractIMDbID(tt.url))
		})
	}
}

func TestAnimeGenreToID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"Action", 10759},
		{"Comedy", 35},
		{"Drama", 18},
		{"Fantasy", 10765},
		{"Horror", 27},
		{"Mystery", 9648},
		{"Romance", 10749},
		{"Sci-Fi", 10765},
		{"Thriller", 53},
		{"Ecchi", 90001},
		{"Mecha", 90003},
		{"Slice of Life", 90005},
		{"Unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, animeGenreToID(tt.name))
		})
	}
}

func TestMapMediaToTVShowSearchResult(t *testing.T) {
	t.Run("full media", func(t *testing.T) {
		m := Media{
			ID:              12345,
			IsAdult:         false,
			Title:           MediaTitle{English: new("Attack on Titan"), Romaji: new("Shingeki no Kyojin"), Native: new("進撃の巨人")},
			CountryOfOrigin: new("JP"),
			Description:     new("Humanity fights giants."),
			AverageScore:    new(85),
			Popularity:      1000,
			StartDate:       FuzzyDate{Year: new(2013), Month: new(4), Day: new(7)},
			CoverImage:      CoverImage{ExtraLarge: new("https://example.com/cover.jpg")},
			BannerImage:     new("https://example.com/banner.jpg"),
			Genres:          []string{"Action", "Drama"},
		}
		result := mapMediaToTVShowSearchResult(m)
		assert.Equal(t, "12345", result.ProviderID)
		assert.Equal(t, metadata.ProviderAniList, result.Provider)
		assert.Equal(t, "Attack on Titan", result.Name)
		assert.Equal(t, "Shingeki no Kyojin", result.OriginalName)
		assert.Equal(t, "jp", result.OriginalLanguage)
		assert.Equal(t, []string{"JP"}, result.OriginCountries)
		assert.Equal(t, "Humanity fights giants.", result.Overview)
		assert.InDelta(t, 8.5, result.VoteAverage, 0.01) // 85/10
		assert.InDelta(t, 1000.0, result.Popularity, 0.01)
		require.NotNil(t, result.FirstAirDate)
		require.NotNil(t, result.Year)
		assert.Equal(t, 2013, *result.Year)
		require.NotNil(t, result.PosterPath)
		require.NotNil(t, result.BackdropPath)
		assert.Equal(t, []int{10759, 18}, result.GenreIDs)
		assert.False(t, result.Adult)
	})

	t.Run("no country - defaults to JP", func(t *testing.T) {
		m := Media{ID: 1, Title: MediaTitle{Romaji: new("Test")}}
		result := mapMediaToTVShowSearchResult(m)
		assert.Equal(t, "ja", result.OriginalLanguage)
		assert.Equal(t, []string{"JP"}, result.OriginCountries)
	})

	t.Run("originalName from native when titles match", func(t *testing.T) {
		m := Media{
			ID:    1,
			Title: MediaTitle{English: new("Same"), Romaji: new("Same"), Native: new("ネイティブ")},
		}
		result := mapMediaToTVShowSearchResult(m)
		assert.Equal(t, "ネイティブ", result.OriginalName)
	})
}

func TestMapMediaToTVShowMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapMediaToTVShowMetadata(nil))
	})

	t.Run("full media", func(t *testing.T) {
		m := &Media{
			ID:              12345,
			IsAdult:         true,
			Title:           MediaTitle{English: new("AoT"), Native: new("進撃の巨人")},
			CountryOfOrigin: new("JP"),
			Description:     new("A description."),
			Status:          "RELEASING",
			Format:          "TV",
			Episodes:        new(25),
			Duration:        new(24),
			AverageScore:    new(85),
			MeanScore:       new(84),
			Popularity:      5000,
			StartDate:       FuzzyDate{Year: new(2013)},
			EndDate:         FuzzyDate{Year: new(2023)},
			CoverImage:      CoverImage{Large: new("https://example.com/poster.jpg")},
			BannerImage:     new("https://example.com/banner.jpg"),
			SiteURL:         "https://anilist.co/anime/12345",
			Trailer:         &Trailer{ID: new("abc123"), Site: new("youtube")},
			Genres:          []string{"Action"},
			Studios:         StudioConnection{Edges: []StudioEdge{{Node: Studio{ID: 1, Name: "WIT"}}}},
		}
		result := mapMediaToTVShowMetadata(m)
		require.NotNil(t, result)
		assert.Equal(t, "12345", result.ProviderID)
		assert.Equal(t, metadata.ProviderAniList, result.Provider)
		assert.True(t, result.Adult)
		assert.Equal(t, "AoT", result.Name)
		assert.Equal(t, "進撃の巨人", result.OriginalName)
		require.NotNil(t, result.Overview)
		assert.Equal(t, "Returning Series", result.Status)
		assert.Equal(t, "Scripted", result.Type)
		assert.True(t, result.InProduction)
		assert.Equal(t, 25, result.NumberOfEpisodes)
		assert.Equal(t, 1, result.NumberOfSeasons)
		assert.Equal(t, []int{24}, result.EpisodeRuntime)
		assert.InDelta(t, 8.5, result.VoteAverage, 0.01)
		require.NotNil(t, result.PosterPath)
		require.NotNil(t, result.BackdropPath)
		require.NotNil(t, result.Homepage)
		assert.Equal(t, "https://anilist.co/anime/12345", *result.Homepage)
		require.NotNil(t, result.TrailerURL)
		assert.Equal(t, "https://www.youtube.com/watch?v=abc123", *result.TrailerURL)
		require.Len(t, result.Genres, 1)
		require.Len(t, result.Networks, 1)
		assert.Equal(t, "WIT", result.Networks[0].Name)
		// External rating from MeanScore
		require.Len(t, result.ExternalRatings, 1)
		assert.Equal(t, "AniList", result.ExternalRatings[0].Source)
		assert.Equal(t, "84%", result.ExternalRatings[0].Value)
	})

	t.Run("no trailer if not youtube", func(t *testing.T) {
		m := &Media{
			ID:      1,
			Title:   MediaTitle{Romaji: new("Test")},
			Trailer: &Trailer{ID: new("abc"), Site: new("dailymotion")},
		}
		result := mapMediaToTVShowMetadata(m)
		require.NotNil(t, result)
		assert.Nil(t, result.TrailerURL)
	})
}

func TestMapCredits(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapCredits(nil))
	})

	t.Run("with characters and staff", func(t *testing.T) {
		m := &Media{
			Characters: CharacterConnection{
				Edges: []CharacterEdge{
					{
						Node: Character{
							ID:     1,
							Name:   CharacterName{Full: new("Eren Yeager")},
							Gender: new("Male"),
						},
						Role: "MAIN",
						VoiceActors: []Staff{
							{
								ID:     100,
								Name:   StaffName{Full: new("Yuki Kaji")},
								Gender: new("Male"),
								Image:  StaffImage{Large: new("https://example.com/va.jpg")},
							},
						},
					},
					{
						Node: Character{
							ID:    2,
							Name:  CharacterName{Full: new("Levi")},
							Image: CharacterImage{Large: new("https://example.com/char.jpg")},
						},
						Role: "SUPPORTING",
					},
				},
			},
			Staff: StaffConnection{
				Edges: []StaffEdge{
					{
						Node: Staff{
							ID:     200,
							Name:   StaffName{Full: new("Tetsuro Araki")},
							Gender: new("Male"),
							Image:  StaffImage{Large: new("https://example.com/staff.jpg")},
						},
						Role: "Director",
					},
				},
			},
		}
		result := mapCredits(m)
		require.NotNil(t, result)
		// Cast: 1 VA for Eren + 1 character-only for Levi
		require.Len(t, result.Cast, 2)
		assert.Equal(t, "Yuki Kaji", result.Cast[0].Name)
		assert.Equal(t, "Eren Yeager", result.Cast[0].Character)
		assert.Equal(t, 2, result.Cast[0].Gender) // Male
		require.NotNil(t, result.Cast[0].ProfilePath)
		// Levi has no VA, so character itself is added
		assert.Equal(t, "Levi", result.Cast[1].Name)
		assert.Equal(t, "SUPPORTING", result.Cast[1].Character)

		// Crew
		require.Len(t, result.Crew, 1)
		assert.Equal(t, "Tetsuro Araki", result.Crew[0].Name)
		assert.Equal(t, "Director", result.Crew[0].Job)
		assert.Equal(t, "Directing", result.Crew[0].Department)
	})
}

func TestMapImages(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapImages(nil))
	})

	t.Run("with images", func(t *testing.T) {
		m := &Media{
			CoverImage:  CoverImage{ExtraLarge: new("xl"), Large: new("l"), Medium: new("m")},
			BannerImage: new("banner"),
		}
		result := mapImages(m)
		require.NotNil(t, result)
		assert.Len(t, result.Posters, 3)
		assert.Len(t, result.Backdrops, 1)
	})

	t.Run("empty images returns nil", func(t *testing.T) {
		m := &Media{CoverImage: CoverImage{}}
		assert.Nil(t, mapImages(m))
	})
}

func TestFindExternalIDs(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, findExternalIDs(nil))
	})

	t.Run("with IMDb link", func(t *testing.T) {
		m := &Media{
			ExternalLinks: []ExternalLink{
				{Site: "IMDb", URL: new("https://www.imdb.com/title/tt1234567/")},
				{Site: "Twitter", URL: new("https://twitter.com/test")},
				{Site: "YouTube", URL: new("https://youtube.com/channel/abc")},
			},
		}
		result := findExternalIDs(m)
		require.NotNil(t, result)
		require.NotNil(t, result.IMDbID)
		assert.Equal(t, "tt1234567", *result.IMDbID)
		require.NotNil(t, result.TwitterID)
		require.NotNil(t, result.YouTubeID)
	})
}
