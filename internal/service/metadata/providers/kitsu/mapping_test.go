package kitsu

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBestImage(t *testing.T) {
	tests := []struct {
		name string
		img  *ImageSet
		want string
	}{
		{"nil", nil, ""},
		{"original first", &ImageSet{Original: new("orig"), Large: new("lg")}, "orig"},
		{"large fallback", &ImageSet{Large: new("lg")}, "lg"},
		{"medium fallback", &ImageSet{Medium: new("md")}, "md"},
		{"small fallback", &ImageSet{Small: new("sm")}, "sm"},
		{"empty", &ImageSet{}, ""},
		{"empty original", &ImageSet{Original: new(""), Large: new("lg")}, "lg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, bestImage(tt.img))
		})
	}
}

func TestImageURLs(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, imageURLs(nil))
	})

	t.Run("with URLs", func(t *testing.T) {
		img := &ImageSet{
			Original: new("orig"),
			Large:    new("lg"),
			Medium:   new(""),
			Small:    new("sm"),
			Tiny:     new("tiny"),
		}
		result := imageURLs(img)
		assert.Equal(t, []string{"orig", "lg", "sm", "tiny"}, result)
	})
}

func TestParseDate(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, parseDate(nil))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, parseDate(new("")))
	})

	t.Run("valid", func(t *testing.T) {
		result := parseDate(new("1998-04-03"))
		require.NotNil(t, result)
		assert.Equal(t, time.Date(1998, 4, 3, 0, 0, 0, 0, time.UTC), *result)
	})

	t.Run("invalid", func(t *testing.T) {
		assert.Nil(t, parseDate(new("bad")))
	})
}

func TestMapStatus(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"finished", "Ended"},
		{"current", "Returning Series"},
		{"upcoming", "Planned"},
		{"tba", "Planned"},
		{"unreleased", "Planned"},
		{"other", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapStatus(tt.input))
		})
	}
}

func TestMapSubtype(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"TV", "Scripted"},
		{"movie", "Movie"},
		{"OVA", "OVA"},
		{"ONA", "ONA"},
		{"special", "Special"},
		{"music", "Music"},
		{"other", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapSubtype(tt.input))
		})
	}
}

func TestCategoryToGenreID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"Action", 10759},
		{"action", 10759},
		{"Comedy", 35},
		{"Drama", 18},
		{"Fantasy", 10765},
		{"Horror", 27},
		{"Mystery", 9648},
		{"Romance", 10749},
		{"Thriller", 53},
		{"Mecha", 90003},
		{"Slice of Life", 90005},
		{"Unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, categoryToGenreID(tt.name))
		})
	}
}

func TestMapAnimeToTVShowSearchResult(t *testing.T) {
	t.Run("full anime", func(t *testing.T) {
		res := ResourceObject[AnimeAttributes]{
			ID:   "12345",
			Type: "anime",
			Attributes: AnimeAttributes{
				CanonicalTitle: "Cowboy Bebop",
				Titles:         map[string]string{"en": "Cowboy Bebop", "ja_jp": "カウボーイビバップ"},
				Synopsis:       "A bounty hunter story in space.",
				AverageRating:  new("87.50"),
				UserCount:      50000,
				StartDate:      new("1998-04-03"),
				NSFW:           false,
				PosterImage:    &ImageSet{Original: new("https://example.com/poster.jpg")},
				CoverImage:     &ImageSet{Original: new("https://example.com/cover.jpg")},
			},
		}
		result := mapAnimeToTVShowSearchResult(res)
		assert.Equal(t, "12345", result.ProviderID)
		assert.Equal(t, metadata.ProviderKitsu, result.Provider)
		assert.Equal(t, "Cowboy Bebop", result.Name)
		assert.Equal(t, "カウボーイビバップ", result.OriginalName)
		assert.Equal(t, "ja", result.OriginalLanguage)
		assert.Equal(t, "A bounty hunter story in space.", result.Overview)
		assert.InDelta(t, 8.75, result.VoteAverage, 0.01) // 87.5/10
		assert.Equal(t, 50000, result.VoteCount)
		assert.InDelta(t, 50000.0, result.Popularity, 0.01)
		require.NotNil(t, result.FirstAirDate)
		require.NotNil(t, result.Year)
		assert.Equal(t, 1998, *result.Year)
		require.NotNil(t, result.PosterPath)
		require.NotNil(t, result.BackdropPath)
		assert.False(t, result.Adult)
	})

	t.Run("en_jp as original name", func(t *testing.T) {
		res := ResourceObject[AnimeAttributes]{
			ID: "1",
			Attributes: AnimeAttributes{
				CanonicalTitle: "Test",
				Titles:         map[string]string{"en_jp": "Romaji Title"},
			},
		}
		result := mapAnimeToTVShowSearchResult(res)
		assert.Equal(t, "Romaji Title", result.OriginalName)
	})
}

func TestMapAnimeToTVShowMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapAnimeToTVShowMetadata(nil))
	})

	t.Run("full response", func(t *testing.T) {
		resp := &SingleResponse[AnimeAttributes]{
			Data: ResourceObject[AnimeAttributes]{
				ID:   "12345",
				Type: "anime",
				Attributes: AnimeAttributes{
					CanonicalTitle: "Cowboy Bebop",
					Slug:           "cowboy-bebop",
					Titles:         map[string]string{"en": "Cowboy Bebop", "ja_jp": "カウボーイビバップ"},
					Synopsis:       "A sci-fi western.",
					AverageRating:  new("87.50"),
					UserCount:      50000,
					Status:         "finished",
					Subtype:        "TV",
					StartDate:      new("1998-04-03"),
					EndDate:        new("1999-04-24"),
					EpisodeCount:   new(26),
					EpisodeLength:  new(24),
					AgeRating:      new("R"),
					NSFW:           false,
					PosterImage:    &ImageSet{Original: new("https://example.com/poster.jpg")},
					CoverImage:     &ImageSet{Original: new("https://example.com/cover.jpg")},
					YoutubeVideoID: new("qig4KOK2R2g"),
				},
			},
			Included: []IncludedResource{
				{Type: "categories", Attributes: map[string]any{"title": "Action"}},
				{Type: "categories", Attributes: map[string]any{"title": "Drama"}},
				{Type: "mappings", Attributes: map[string]any{"externalSite": "myanimelist/anime", "externalId": "1"}},
				{Type: "mappings", Attributes: map[string]any{"externalSite": "thetvdb/series", "externalId": "76885"}},
			},
		}

		result := mapAnimeToTVShowMetadata(resp)
		require.NotNil(t, result)
		assert.Equal(t, "12345", result.ProviderID)
		assert.Equal(t, metadata.ProviderKitsu, result.Provider)
		assert.Equal(t, "Cowboy Bebop", result.Name)
		assert.Equal(t, "カウボーイビバップ", result.OriginalName)
		require.NotNil(t, result.Overview)
		assert.Equal(t, "Ended", result.Status)
		assert.Equal(t, "Scripted", result.Type)
		assert.False(t, result.InProduction)
		assert.Equal(t, 26, result.NumberOfEpisodes)
		assert.Equal(t, 1, result.NumberOfSeasons)
		assert.Equal(t, []int{24}, result.EpisodeRuntime)
		assert.InDelta(t, 8.75, result.VoteAverage, 0.01)
		require.NotNil(t, result.PosterPath)
		require.NotNil(t, result.BackdropPath)
		require.NotNil(t, result.TrailerURL)
		assert.Equal(t, "https://www.youtube.com/watch?v=qig4KOK2R2g", *result.TrailerURL)
		require.NotNil(t, result.Homepage)
		assert.Contains(t, *result.Homepage, "kitsu.io/anime/cowboy-bebop")
		require.Len(t, result.Genres, 2)
		require.NotNil(t, result.TVDbID)
		assert.Equal(t, int32(76885), *result.TVDbID)
		// External rating
		require.Len(t, result.ExternalRatings, 1)
		assert.Equal(t, "Kitsu", result.ExternalRatings[0].Source)
	})

	t.Run("current status is in production", func(t *testing.T) {
		resp := &SingleResponse[AnimeAttributes]{
			Data: ResourceObject[AnimeAttributes]{
				ID: "1",
				Attributes: AnimeAttributes{
					CanonicalTitle: "Test",
					Status:         "current",
				},
			},
		}
		result := mapAnimeToTVShowMetadata(resp)
		require.NotNil(t, result)
		assert.True(t, result.InProduction)
		assert.Equal(t, "Returning Series", result.Status)
	})
}

func TestMapEpisodesToSummary(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapEpisodesToSummary(nil, 1))
	})

	t.Run("with episodes", func(t *testing.T) {
		episodes := &ListResponse[EpisodeAttributes]{
			Data: []ResourceObject[EpisodeAttributes]{
				{
					ID: "ep1",
					Attributes: EpisodeAttributes{
						CanonicalTitle: "Asteroid Blues",
						Number:         new(1),
						SeasonNumber:   new(1),
						Length:         new(24),
						Airdate:        new("1998-10-24"),
						Synopsis:       "First episode.",
						Thumbnail:      &ImageSet{Original: new("https://example.com/ep1.jpg")},
					},
				},
				{
					ID: "ep2",
					Attributes: EpisodeAttributes{
						CanonicalTitle: "Stray Dog Strut",
						Number:         new(2),
						SeasonNumber:   new(1),
					},
				},
				{
					ID: "ep-s2",
					Attributes: EpisodeAttributes{
						CanonicalTitle: "S2 Ep",
						Number:         new(1),
						SeasonNumber:   new(2),
					},
				},
			},
		}

		result := mapEpisodesToSummary(episodes, 1)
		require.Len(t, result, 2)
		assert.Equal(t, "ep1", result[0].ProviderID)
		assert.Equal(t, "Asteroid Blues", result[0].Name)
		assert.Equal(t, 1, result[0].EpisodeNumber)
		require.NotNil(t, result[0].Runtime)
		assert.Equal(t, int32(24), *result[0].Runtime)
		require.NotNil(t, result[0].AirDate)
		require.NotNil(t, result[0].Overview)
		assert.Equal(t, "First episode.", *result[0].Overview)
		require.NotNil(t, result[0].StillPath)
	})
}

func TestMapEpisodeToMetadata(t *testing.T) {
	ep := ResourceObject[EpisodeAttributes]{
		ID: "ep42",
		Attributes: EpisodeAttributes{
			CanonicalTitle: "Asteroid Blues",
			Number:         new(1),
			SeasonNumber:   new(1),
			Length:         new(24),
			Airdate:        new("1998-10-24"),
			Synopsis:       "Test synopsis.",
			Thumbnail:      &ImageSet{Original: new("https://example.com/thumb.jpg")},
		},
	}
	result := mapEpisodeToMetadata(ep, "anime-123")
	require.NotNil(t, result)
	assert.Equal(t, "ep42", result.ProviderID)
	assert.Equal(t, metadata.ProviderKitsu, result.Provider)
	assert.Equal(t, "anime-123", result.ShowID)
	assert.Equal(t, 1, result.SeasonNumber)
	assert.Equal(t, 1, result.EpisodeNumber)
	assert.Equal(t, "Asteroid Blues", result.Name)
	require.NotNil(t, result.Runtime)
	assert.Equal(t, int32(24), *result.Runtime)
	require.NotNil(t, result.Overview)
	require.NotNil(t, result.StillPath)

	// Default season 1 when no season number
	ep2 := ResourceObject[EpisodeAttributes]{
		ID:         "ep2",
		Attributes: EpisodeAttributes{CanonicalTitle: "Test", Number: new(5)},
	}
	result2 := mapEpisodeToMetadata(ep2, "anime-1")
	assert.Equal(t, 1, result2.SeasonNumber)
}

func TestMapImages(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapImages(nil))
	})

	t.Run("with images", func(t *testing.T) {
		resp := &SingleResponse[AnimeAttributes]{
			Data: ResourceObject[AnimeAttributes]{
				Attributes: AnimeAttributes{
					PosterImage: &ImageSet{Original: new("poster"), Large: new("poster-lg")},
					CoverImage:  &ImageSet{Original: new("cover")},
				},
			},
		}
		result := mapImages(resp)
		require.NotNil(t, result)
		assert.Len(t, result.Posters, 2)
		assert.Len(t, result.Backdrops, 1)
	})

	t.Run("empty images returns nil", func(t *testing.T) {
		resp := &SingleResponse[AnimeAttributes]{
			Data: ResourceObject[AnimeAttributes]{Attributes: AnimeAttributes{}},
		}
		assert.Nil(t, mapImages(resp))
	})
}

func TestMapMappingsToExternalIDs(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapMappingsToExternalIDs(nil))
	})

	t.Run("empty data", func(t *testing.T) {
		assert.Nil(t, mapMappingsToExternalIDs(&ListResponse[MappingAttributes]{}))
	})

	t.Run("with mappings", func(t *testing.T) {
		mappings := &ListResponse[MappingAttributes]{
			Data: []ResourceObject[MappingAttributes]{
				{ID: "1", Attributes: MappingAttributes{ExternalSite: "thetvdb/series", ExternalID: "76885"}},
				{ID: "2", Attributes: MappingAttributes{ExternalSite: "myanimelist/anime", ExternalID: "1"}},
				{ID: "3", Attributes: MappingAttributes{ExternalSite: "anilist/anime", ExternalID: "1"}},
			},
		}
		result := mapMappingsToExternalIDs(mappings)
		require.NotNil(t, result)
		require.NotNil(t, result.TVDbID)
		assert.Equal(t, int32(76885), *result.TVDbID)
	})
}
