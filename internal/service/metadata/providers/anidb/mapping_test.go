package anidb

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBestTitle(t *testing.T) {
	titles := Titles{
		Title: []Title{
			{Lang: "ja", Type: "official", Text: "進撃の巨人"},
			{Lang: "en", Type: "official", Text: "Attack on Titan"},
			{Lang: "x-jat", Type: "main", Text: "Shingeki no Kyojin"},
		},
	}

	assert.Equal(t, "Attack on Titan", bestTitle(titles, "en"))
	assert.Equal(t, "進撃の巨人", bestTitle(titles, "ja"))
	// main type fallback when no official title in language
	assert.Equal(t, "Shingeki no Kyojin", bestTitle(titles, "de"))

	// Empty titles
	assert.Equal(t, "", bestTitle(Titles{}, "en"))
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantNil  bool
		wantDate time.Time
	}{
		{"empty", "", true, time.Time{}},
		{"full date", "2013-04-07", false, time.Date(2013, 4, 7, 0, 0, 0, 0, time.UTC)},
		{"year-month", "2013-04", false, time.Date(2013, 4, 1, 0, 0, 0, 0, time.UTC)},
		{"year only", "2013", false, time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"invalid", "bad", true, time.Time{}},
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

func TestCleanDescription(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain text", "A simple description.", "A simple description."},
		{"bold tags", "[b]Bold[/b] text", "Bold text"},
		{"italic tags", "[i]Italic[/i] text", "Italic text"},
		{"url tags", "Visit [url=https://example.com]Example[/url] site", "Visit Example site"},
		{"mixed", "[b]Title[/b]: [i]subtitle[/i] with [url=http://test.com]link[/url]", "Title: subtitle with link"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cleanDescription(tt.input))
		})
	}
}

func TestMapType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"TV Series", "Scripted"},
		{"Movie", "Movie"},
		{"OVA", "OVA"},
		{"Web", "ONA"},
		{"TV Special", "Special"},
		{"Music Video", "Music"},
		{"Other", "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapType(tt.input))
		})
	}
}

func TestMapStatusFromDates(t *testing.T) {
	tests := []struct {
		name  string
		start string
		end   string
		want  string
	}{
		{"no start", "", "", "Planned"},
		{"ongoing", "2020-01-01", "", "Returning Series"},
		{"ended", "2020-01-01", "2021-12-31", "Ended"},
		{"future start", "2099-01-01", "", "Planned"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mapStatusFromDates(tt.start, tt.end))
		})
	}
}

func TestMapGender(t *testing.T) {
	assert.Equal(t, 1, mapGender("female"))
	assert.Equal(t, 1, mapGender("Female"))
	assert.Equal(t, 2, mapGender("male"))
	assert.Equal(t, 2, mapGender("Male"))
	assert.Equal(t, 0, mapGender("unknown"))
	assert.Equal(t, 0, mapGender(""))
}

func TestMapDepartment(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Direction", "Directing"},
		{"Music", "Sound"},
		{"Character Design", "Art"},
		{"Animation Work", "Art"},
		{"Series Composition", "Writing"},
		{"Original Work", "Writing"},
		{"Something New", "Production"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapDepartment(tt.input))
		})
	}
}

func TestMapTitleToTVShowSearchResult(t *testing.T) {
	entry := TitleDumpEntry{AID: 12345, Title: "Test Anime", Lang: "en", Type: "official"}
	result := mapTitleToTVShowSearchResult(entry)
	assert.Equal(t, "12345", result.ProviderID)
	assert.Equal(t, metadata.ProviderAniDB, result.Provider)
	assert.Equal(t, "Test Anime", result.Name)
	assert.Equal(t, "ja", result.OriginalLanguage)
	assert.Equal(t, []string{"JP"}, result.OriginCountries)
}

func TestMapAnimeToTVShowMetadata(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapAnimeToTVShowMetadata(nil))
	})

	t.Run("full anime", func(t *testing.T) {
		a := &AnimeResponse{
			ID:         12345,
			Restricted: false,
			Type:       "TV Series",
			EpCount:    25,
			StartDate:  "2013-04-07",
			EndDate:    "2013-09-28",
			Description: "[b]Synopsis[/b]: Humanity fights for survival.",
			Picture:    "12345.jpg",
			URL:        "https://anidb.net/anime/12345",
			Titles: Titles{
				Title: []Title{
					{Lang: "en", Type: "official", Text: "Attack on Titan"},
					{Lang: "ja", Type: "official", Text: "進撃の巨人"},
				},
			},
			Episodes: Episodes{
				Episode: []Episode{
					{ID: 1, EpNo: EpNo{Type: 1, Text: "1"}, Length: 24, Title: []EpTitle{{Lang: "en", Text: "To You, in 2000 Years"}}},
				},
			},
			Ratings: Ratings{
				Permanent: RatingValue{Count: 5000, Value: 8.5},
				Temporary: RatingValue{Count: 1000, Value: 8.2},
			},
			Tags: Tags{
				Tag: []Tag{
					{ID: 1, Weight: 600, Name: "action"},
					{ID: 2, Weight: 100, Name: "low-weight"},        // filtered
					{ID: 3, Weight: 300, GlobalSpoiler: true, Name: "spoiler"}, // filtered
				},
			},
			Creators: Creators{
				Name: []Creator{
					{ID: 1, Type: "Animation Work", Text: "Wit Studio"},
				},
			},
		}

		result := mapAnimeToTVShowMetadata(a)
		require.NotNil(t, result)
		assert.Equal(t, "12345", result.ProviderID)
		assert.Equal(t, metadata.ProviderAniDB, result.Provider)
		assert.Equal(t, "Attack on Titan", result.Name)
		assert.Equal(t, "進撃の巨人", result.OriginalName)
		assert.Equal(t, "ja", result.OriginalLanguage)
		assert.False(t, result.Adult)
		assert.Equal(t, 25, result.NumberOfEpisodes)
		assert.Equal(t, 1, result.NumberOfSeasons)
		assert.Equal(t, "Scripted", result.Type)
		assert.Equal(t, "Ended", result.Status)
		require.NotNil(t, result.Overview)
		assert.Contains(t, *result.Overview, "Humanity fights for survival")
		assert.NotContains(t, *result.Overview, "[b]")
		require.NotNil(t, result.FirstAirDate)
		require.NotNil(t, result.LastAirDate)
		assert.Equal(t, []int{24}, result.EpisodeRuntime)
		assert.InDelta(t, 8.5, result.VoteAverage, 0.01)
		require.NotNil(t, result.PosterPath)
		assert.Contains(t, *result.PosterPath, "12345.jpg")
		require.NotNil(t, result.Homepage)
		// Tags: only "action" survives filtering
		require.Len(t, result.Genres, 1)
		assert.Equal(t, "action", result.Genres[0].Name)
		// Studios
		require.Len(t, result.Networks, 1)
		assert.Equal(t, "Wit Studio", result.Networks[0].Name)
		// External ratings
		require.Len(t, result.ExternalRatings, 2)
	})
}

func TestMapCredits(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapCredits(nil))
	})

	t.Run("with characters and creators", func(t *testing.T) {
		a := &AnimeResponse{
			Characters: Characters{
				Character: []Character{
					{
						ID:     1,
						Name:   "Eren Yeager",
						Gender: "male",
						Seiyuu: []Seiyuu{{ID: 100, Text: "Yuki Kaji", Picture: "kaji.jpg"}},
					},
					{
						ID:   2,
						Name: "Levi",
						Type: "main character in",
					},
				},
			},
			Creators: Creators{
				Name: []Creator{
					{ID: 200, Type: "Direction", Text: "Tetsuro Araki"},
				},
			},
		}
		result := mapCredits(a)
		require.NotNil(t, result)
		// Cast: VA for Eren + character-only for Levi
		require.Len(t, result.Cast, 2)
		assert.Equal(t, "Yuki Kaji", result.Cast[0].Name)
		assert.Equal(t, "Eren Yeager", result.Cast[0].Character)
		require.NotNil(t, result.Cast[0].ProfilePath)
		assert.Equal(t, "Levi", result.Cast[1].Name)
		// Crew
		require.Len(t, result.Crew, 1)
		assert.Equal(t, "Tetsuro Araki", result.Crew[0].Name)
		assert.Equal(t, "Directing", result.Crew[0].Department)
	})
}

func TestMapImages(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, mapImages(nil))
	})

	t.Run("no picture", func(t *testing.T) {
		assert.Nil(t, mapImages(&AnimeResponse{}))
	})

	t.Run("with picture", func(t *testing.T) {
		result := mapImages(&AnimeResponse{Picture: "test.jpg"})
		require.NotNil(t, result)
		require.Len(t, result.Posters, 1)
		assert.Contains(t, result.Posters[0].FilePath, "test.jpg")
	})
}

func TestMapEpisodes(t *testing.T) {
	a := &AnimeResponse{
		Episodes: Episodes{
			Episode: []Episode{
				{
					ID:      1,
					EpNo:    EpNo{Type: 1, Text: "1"},
					Length:  24,
					Airdate: "2013-04-07",
					Title:   []EpTitle{{Lang: "en", Text: "To You, in 2000 Years"}},
					Rating:  &EpisodeRating{Votes: 100, Value: 8.5},
				},
				{
					ID:    2,
					EpNo:  EpNo{Type: 2, Text: "S1"}, // special - should be skipped
					Title: []EpTitle{{Lang: "en", Text: "Special"}},
				},
				{
					ID:      3,
					EpNo:    EpNo{Type: 1, Text: "2"},
					Length:  24,
					Airdate: "2013-04-14",
					Title:   []EpTitle{{Lang: "x-jat", Text: "Sono Hi"}},
				},
			},
		},
	}

	result := mapEpisodes(a, 1)
	require.Len(t, result, 2) // only regular episodes
	assert.Equal(t, "1", result[0].ProviderID)
	assert.Equal(t, 1, result[0].EpisodeNumber)
	assert.Equal(t, "To You, in 2000 Years", result[0].Name)
	require.NotNil(t, result[0].AirDate)
	require.NotNil(t, result[0].Runtime)
	assert.Equal(t, int32(24), *result[0].Runtime)
	assert.InDelta(t, 8.5, result[0].VoteAverage, 0.01)

	// Second episode falls back to x-jat title
	assert.Equal(t, "Sono Hi", result[1].Name)
}

func TestMapEpisodeToMetadata(t *testing.T) {
	ep := Episode{
		ID:      42,
		EpNo:    EpNo{Type: 1, Text: "5"},
		Length:  24,
		Airdate: "2013-05-05",
		Title:   []EpTitle{{Lang: "en", Text: "First Battle"}},
		Rating:  &EpisodeRating{Value: 9.0},
	}
	result := mapEpisodeToMetadata(ep, "12345")
	require.NotNil(t, result)
	assert.Equal(t, "42", result.ProviderID)
	assert.Equal(t, metadata.ProviderAniDB, result.Provider)
	assert.Equal(t, "12345", result.ShowID)
	assert.Equal(t, 1, result.SeasonNumber)
	assert.Equal(t, 5, result.EpisodeNumber)
	assert.Equal(t, "First Battle", result.Name)
	require.NotNil(t, result.Runtime)
	assert.Equal(t, int32(24), *result.Runtime)
	assert.InDelta(t, 9.0, result.VoteAverage, 0.01)

	// Invalid episode number
	epBad := Episode{EpNo: EpNo{Type: 1, Text: "bad"}}
	assert.Nil(t, mapEpisodeToMetadata(epBad, "1"))
}

func TestFindExternalIDs(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, findExternalIDs(nil))
	})

	t.Run("with resources", func(t *testing.T) {
		a := &AnimeResponse{
			Resources: Resources{
				Resource: []Resource{
					{Type: 6, ExternalID: []ExternalEntity{{Identifier: []string{"Q12345"}}}},
					{Type: 1, ExternalID: []ExternalEntity{{Identifier: []string{"123"}}}},
					{Type: 99, ExternalID: []ExternalEntity{}}, // empty
				},
			},
		}
		result := findExternalIDs(a)
		require.NotNil(t, result)
		require.NotNil(t, result.WikidataID)
		assert.Equal(t, "Q12345", *result.WikidataID)
	})
}
