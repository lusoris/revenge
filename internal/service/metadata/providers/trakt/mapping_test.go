package trakt

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantYr  int
	}{
		{"empty", "", true, 0},
		{"valid", "2024-06-15", false, 2024},
		{"invalid", "not-a-date", true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDate(tt.input)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.wantYr, result.Year())
			}
		})
	}
}

func TestMapStatus(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"released", "Released"},
		{"in production", "In Production"},
		{"post production", "Post Production"},
		{"planned", "Planned"},
		{"rumored", "Rumored"},
		{"canceled", "Canceled"},
		{"Unknown", "Unknown"},
		{"Released", "Released"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapStatus(tt.input))
		})
	}
}

func TestMapShowStatus(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"returning series", "Returning Series"},
		{"ended", "Ended"},
		{"canceled", "Canceled"},
		{"in production", "In Production"},
		{"planned", "Planned"},
		{"Unknown", "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapShowStatus(tt.input))
		})
	}
}

func TestGenreNameToID(t *testing.T) {
	assert.Equal(t, 28, genreNameToID("action"))
	assert.Equal(t, 18, genreNameToID("drama"))
	assert.Equal(t, 878, genreNameToID("science-fiction"))
	assert.Equal(t, 53, genreNameToID("thriller"))
	assert.Equal(t, 53, genreNameToID("suspense"))
	assert.Equal(t, 0, genreNameToID("nonexistent"))
}

func TestCapitalizeGenre(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"action", "Action"},
		{"science-fiction", "Science Fiction"},
		{"war", "War"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, capitalizeGenre(tt.input))
		})
	}
}

func TestMapMovieToSearchResult(t *testing.T) {
	m := &Movie{
		Title:    "Fight Club",
		Year:     1999,
		IDs:      IDs{Trakt: 432, IMDb: "tt0137523", TMDb: 550},
		Overview: "An insomniac office worker...",
		Released: "1999-10-15",
		Rating:   8.4,
		Votes:    25000,
		Language: "en",
		Genres:   []string{"drama", "thriller"},
	}

	result := mapMovieToSearchResult(m)

	assert.Equal(t, "432", result.ProviderID)
	assert.Equal(t, metadata.ProviderTrakt, result.Provider)
	assert.Equal(t, "Fight Club", result.Title)
	assert.Equal(t, "An insomniac office worker...", result.Overview)
	assert.Equal(t, "en", result.OriginalLanguage)
	assert.Equal(t, 8.4, result.VoteAverage)
	assert.Equal(t, 25000, result.VoteCount)
	require.NotNil(t, result.Year)
	assert.Equal(t, 1999, *result.Year)
	require.NotNil(t, result.ReleaseDate)
	assert.Len(t, result.GenreIDs, 2)
}

func TestMapMovieToSearchResult_Nil(t *testing.T) {
	result := mapMovieToSearchResult(nil)
	assert.Equal(t, "", result.ProviderID)
}

func TestMapMovieToMetadata(t *testing.T) {
	m := &Movie{
		Title:    "Fight Club",
		Year:     1999,
		IDs:      IDs{Trakt: 432, IMDb: "tt0137523", TMDb: 550, TVDb: 100},
		Tagline:  "Mischief. Mayhem. Soap.",
		Overview: "An insomniac office worker...",
		Released: "1999-10-15",
		Runtime:  139,
		Country:  "us",
		Trailer:  "https://trailer.url",
		Homepage: "https://homepage.url",
		Status:   "released",
		Rating:   8.4,
		Votes:    25000,
		Language: "en",
		Genres:   []string{"drama", "action"},
	}

	result := mapMovieToMetadata(m)
	require.NotNil(t, result)

	assert.Equal(t, "432", result.ProviderID)
	assert.Equal(t, metadata.ProviderTrakt, result.Provider)
	assert.Equal(t, "Fight Club", result.Title)
	assert.Equal(t, "en", result.OriginalLanguage)
	assert.Equal(t, "Released", result.Status)
	assert.Equal(t, 8.4, result.VoteAverage)

	assert.Equal(t, new("Mischief. Mayhem. Soap."), result.Tagline)
	assert.Equal(t, new("An insomniac office worker..."), result.Overview)
	assert.Equal(t, new("https://homepage.url"), result.Homepage)
	assert.Equal(t, new("https://trailer.url"), result.TrailerURL)

	require.NotNil(t, result.Runtime)
	assert.Equal(t, int32(139), *result.Runtime)

	require.NotNil(t, result.IMDbID)
	assert.Equal(t, "tt0137523", *result.IMDbID)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(550), *result.TMDbID)
	require.NotNil(t, result.TVDbID)
	assert.Equal(t, int32(100), *result.TVDbID)

	require.NotNil(t, result.ReleaseDate)

	require.Len(t, result.ExternalRatings, 1)
	assert.Equal(t, "Trakt", result.ExternalRatings[0].Source)
	assert.InDelta(t, 84.0, result.ExternalRatings[0].Score, 0.01)

	require.Len(t, result.Genres, 2)
	assert.Equal(t, "Drama", result.Genres[0].Name)

	require.Len(t, result.ProductionCountries, 1)
	assert.Equal(t, "US", result.ProductionCountries[0].ISOCode)
}

func TestMapMovieToMetadata_Nil(t *testing.T) {
	assert.Nil(t, mapMovieToMetadata(nil))
}

func TestMapMovieToMetadata_EmptyOptionals(t *testing.T) {
	m := &Movie{IDs: IDs{Trakt: 1}, Rating: 0}
	result := mapMovieToMetadata(m)
	assert.Nil(t, result.Tagline)
	assert.Nil(t, result.Overview)
	assert.Nil(t, result.Homepage)
	assert.Nil(t, result.TrailerURL)
	assert.Nil(t, result.Runtime)
	assert.Empty(t, result.ExternalRatings)
	assert.Nil(t, result.IMDbID)
}

func TestMapShowToSearchResult(t *testing.T) {
	aired := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	s := &Show{
		Title:      "Breaking Bad",
		Year:       2008,
		IDs:        IDs{Trakt: 1388},
		Overview:   "A chemistry teacher...",
		Language:   "en",
		Country:    "us",
		Rating:     8.9,
		Votes:      12000,
		Genres:     []string{"drama", "crime"},
		FirstAired: &aired,
	}

	result := mapShowToSearchResult(s)

	assert.Equal(t, "1388", result.ProviderID)
	assert.Equal(t, metadata.ProviderTrakt, result.Provider)
	assert.Equal(t, "Breaking Bad", result.Name)
	require.NotNil(t, result.Year)
	assert.Equal(t, 2008, *result.Year)
	require.NotNil(t, result.FirstAirDate)
	assert.Equal(t, []string{"US"}, result.OriginCountries)
	assert.Len(t, result.GenreIDs, 2)
}

func TestMapShowToSearchResult_Nil(t *testing.T) {
	result := mapShowToSearchResult(nil)
	assert.Equal(t, "", result.ProviderID)
}

func TestMapShowToMetadata(t *testing.T) {
	aired := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	s := &Show{
		Title:      "Breaking Bad",
		Year:       2008,
		IDs:        IDs{Trakt: 1388, IMDb: "tt0903747", TMDb: 1396, TVDb: 73255},
		Overview:   "A chemistry teacher...",
		Homepage:   "https://homepage.url",
		Trailer:    "https://trailer.url",
		Runtime:    47,
		Network:    "AMC",
		Country:    "us",
		Status:     "ended",
		Rating:     8.9,
		Votes:      12000,
		Language:   "en",
		Genres:     []string{"drama"},
		FirstAired: &aired,
	}

	result := mapShowToMetadata(s)
	require.NotNil(t, result)

	assert.Equal(t, "1388", result.ProviderID)
	assert.Equal(t, metadata.ProviderTrakt, result.Provider)
	assert.Equal(t, "Breaking Bad", result.Name)
	assert.Equal(t, "Ended", result.Status)
	assert.False(t, result.InProduction)
	assert.Equal(t, []int{47}, result.EpisodeRuntime)

	assert.Equal(t, new("A chemistry teacher..."), result.Overview)
	assert.Equal(t, new("https://homepage.url"), result.Homepage)
	assert.Equal(t, new("https://trailer.url"), result.TrailerURL)

	require.NotNil(t, result.FirstAirDate)
	assert.Equal(t, []string{"US"}, result.OriginCountries)

	require.Len(t, result.Networks, 1)
	assert.Equal(t, "AMC", result.Networks[0].Name)

	require.Len(t, result.ExternalRatings, 1)
	assert.Equal(t, "Trakt", result.ExternalRatings[0].Source)
}

func TestMapShowToMetadata_ReturningSeriesInProduction(t *testing.T) {
	s := &Show{IDs: IDs{Trakt: 1}, Status: "returning series"}
	result := mapShowToMetadata(s)
	assert.True(t, result.InProduction)
	assert.Equal(t, "Returning Series", result.Status)
}

func TestMapShowToMetadata_Nil(t *testing.T) {
	assert.Nil(t, mapShowToMetadata(nil))
}

func TestMapSeasons(t *testing.T) {
	aired := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	seasons := []Season{
		{Number: 1, IDs: IDs{Trakt: 100}, Title: "Season 1", EpisodeCount: 7, Rating: 8.3, FirstAired: &aired, Overview: "The first season"},
		{Number: 2, IDs: IDs{Trakt: 101}, Title: "Season 2", EpisodeCount: 13, Rating: 8.7},
	}

	result := mapSeasons(seasons)

	require.Len(t, result, 2)
	assert.Equal(t, "100", result[0].ProviderID)
	assert.Equal(t, 1, result[0].SeasonNumber)
	assert.Equal(t, "Season 1", result[0].Name)
	assert.Equal(t, 7, result[0].EpisodeCount)
	require.NotNil(t, result[0].Overview)
	assert.Equal(t, "The first season", *result[0].Overview)
	require.NotNil(t, result[0].AirDate)
}

func TestMapSeasons_Empty(t *testing.T) {
	assert.Nil(t, mapSeasons(nil))
	assert.Nil(t, mapSeasons([]Season{}))
}

func TestMapEpisodesToSummaries(t *testing.T) {
	aired := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	episodes := []Episode{
		{Number: 1, IDs: IDs{Trakt: 200}, Title: "Pilot", Overview: "Walter White...", Rating: 8.5, Votes: 500, FirstAired: &aired, Runtime: 58},
		{Number: 2, IDs: IDs{Trakt: 201}, Title: "Cat's in the Bag...", Runtime: 0},
	}

	result := mapEpisodesToSummaries(episodes)

	require.Len(t, result, 2)
	assert.Equal(t, "200", result[0].ProviderID)
	assert.Equal(t, 1, result[0].EpisodeNumber)
	assert.Equal(t, "Pilot", result[0].Name)
	require.NotNil(t, result[0].Overview)
	require.NotNil(t, result[0].Runtime)
	assert.Equal(t, int32(58), *result[0].Runtime)

	assert.Nil(t, result[1].Overview)
	assert.Nil(t, result[1].Runtime)
}

func TestMapEpisodesToSummaries_Empty(t *testing.T) {
	assert.Nil(t, mapEpisodesToSummaries(nil))
}

func TestMapCredits(t *testing.T) {
	c := &Credits{
		Cast: []CastMember{
			{Characters: []string{"Walter White", "Heisenberg"}, Person: Person{Name: "Bryan Cranston", IDs: IDs{Trakt: 17419}}},
			{Person: Person{Name: "Aaron Paul", IDs: IDs{Trakt: 17420}}},
		},
		Crew: map[string][]CrewMember{
			"directing": {
				{Jobs: []string{"Director", "Executive Producer"}, Person: Person{Name: "Vince Gilligan", IDs: IDs{Trakt: 100}}},
			},
			"writing": {
				{Person: Person{Name: "Peter Gould", IDs: IDs{Trakt: 101}}},
			},
		},
	}

	result := mapCredits(c)

	require.Len(t, result.Cast, 2)
	assert.Equal(t, "17419", result.Cast[0].ProviderID)
	assert.Equal(t, "Bryan Cranston", result.Cast[0].Name)
	assert.Equal(t, "Walter White, Heisenberg", result.Cast[0].Character)
	assert.Equal(t, 0, result.Cast[0].Order)
	assert.Equal(t, "", result.Cast[1].Character) // No characters

	require.Len(t, result.Crew, 2)
}

func TestMapCredits_Nil(t *testing.T) {
	assert.Nil(t, mapCredits(nil))
}

func TestMapTranslations(t *testing.T) {
	translations := []Translation{
		{Title: "Titre", Overview: "Vue d'ensemble", Language: "fr", Country: "FR", Tagline: "Slogan"},
		{Title: "Titel", Overview: "Ubersicht", Language: "de", Country: "DE"},
	}

	result := mapTranslations(translations)

	require.Len(t, result, 2)
	assert.Equal(t, "fr", result[0].Language)
	assert.Equal(t, "FR", result[0].ISOCode)
	require.NotNil(t, result[0].Data)
	assert.Equal(t, "Titre", result[0].Data.Title)
	assert.Equal(t, "Slogan", result[0].Data.Tagline)
}

func TestMapTranslations_Empty(t *testing.T) {
	assert.Nil(t, mapTranslations(nil))
	assert.Nil(t, mapTranslations([]Translation{}))
}

func TestMapExternalIDs(t *testing.T) {
	ids := IDs{Trakt: 432, IMDb: "tt0137523", TMDb: 550, TVDb: 100}

	result := mapExternalIDs(ids)

	require.NotNil(t, result.IMDbID)
	assert.Equal(t, "tt0137523", *result.IMDbID)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(550), *result.TMDbID)
	require.NotNil(t, result.TVDbID)
	assert.Equal(t, int32(100), *result.TVDbID)
}

func TestMapExternalIDs_Empty(t *testing.T) {
	result := mapExternalIDs(IDs{})
	assert.Nil(t, result.IMDbID)
	assert.Nil(t, result.TMDbID)
	assert.Nil(t, result.TVDbID)
}
