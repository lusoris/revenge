package tvdb

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantYr  int
	}{
		{"empty", "", true, 0},
		{"simple date", "2024-06-15", false, 2024},
		{"ISO with time", "2024-06-15T10:30:00Z", false, 2024},
		{"datetime with space", "2024-06-15 10:30:00", false, 2024},
		{"year only", "2008", false, 2008},
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

func TestCountSeasons(t *testing.T) {
	tests := []struct {
		name    string
		seasons []SeasonSummaryResponse
		want    int
	}{
		{"nil", nil, 0},
		{"empty", []SeasonSummaryResponse{}, 0},
		{"all default", []SeasonSummaryResponse{
			{Number: 1, Type: &SeasonTypeResponse{Type: "default"}},
			{Number: 2, Type: &SeasonTypeResponse{Type: "official"}},
			{Number: 0, Type: nil},
		}, 3},
		{"mixed types", []SeasonSummaryResponse{
			{Number: 1, Type: &SeasonTypeResponse{Type: "default"}},
			{Number: 2, Type: &SeasonTypeResponse{Type: "absolute"}},
			{Number: 3, Type: &SeasonTypeResponse{Type: "official"}},
		}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, countSeasons(tt.seasons))
		})
	}
}

func TestCharacterTypeToJob(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{CharacterTypeDirector, "Director"},
		{CharacterTypeWriter, "Writer"},
		{CharacterTypeProducer, "Producer"},
		{99, "Unknown"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, characterTypeToJob(tt.input))
	}
}

func TestCharacterTypeToDepartment(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{CharacterTypeDirector, "Directing"},
		{CharacterTypeWriter, "Writing"},
		{CharacterTypeProducer, "Production"},
		{99, "Unknown"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, characterTypeToDepartment(tt.input))
	}
}

func TestPtrString(t *testing.T) {
	assert.Equal(t, "", ptrString(nil))
	assert.Equal(t, "hello", ptrString(ptr("hello")))
}

func TestMapTVSearchResult(t *testing.T) {
	input := &SearchResult{
		TVDbID:       "73255",
		Name:         "Breaking Bad",
		Overview:     "A high school chemistry teacher...",
		ImageURL:     ptr("https://artworks.thetvdb.com/banners/poster.jpg"),
		Year:         "2008",
		FirstAirTime: "2008-01-20",
	}

	result := mapTVSearchResult(input)
	assert.Equal(t, "73255", result.ProviderID)
	assert.Equal(t, metadata.ProviderTVDb, result.Provider)
	assert.Equal(t, "Breaking Bad", result.Name)
	assert.Equal(t, ptr("https://artworks.thetvdb.com/banners/poster.jpg"), result.PosterPath)
	require.NotNil(t, result.Year)
	assert.Equal(t, 2008, *result.Year)
	require.NotNil(t, result.FirstAirDate)
}

func TestMapTVSearchResult_NoYear(t *testing.T) {
	input := &SearchResult{TVDbID: "1", Name: "Test"}
	result := mapTVSearchResult(input)
	assert.Nil(t, result.Year)
	assert.Nil(t, result.FirstAirDate)
}

func TestMapTVShowMetadata(t *testing.T) {
	input := &SeriesResponse{
		ID:               73255,
		Name:             "Breaking Bad",
		OriginalLanguage: "eng",
		Overview:         ptr("A high school chemistry teacher..."),
		Score:            85,
		Status:           &StatusResponse{Name: "Ended", KeepUpdated: false},
		FirstAired:       "2008-01-20",
		LastAired:        "2013-09-29",
		AverageRuntime:   47,
		OriginalCountry:  "usa",
		Genres:           []GenreResponse{{ID: 18, Name: "Drama"}},
		Networks:         []NetworkResponse{{ID: 174, Name: "AMC", Country: ptr("usa")}},
		Seasons: []SeasonSummaryResponse{
			{ID: 1, Number: 1, Name: "Season 1", Type: &SeasonTypeResponse{Type: "default"}, Year: "2008"},
		},
		RemoteIDs: []RemoteIDResponse{
			{ID: "tt0903747", Type: RemoteIDTypeIMDb},
			{ID: "1396", Type: RemoteIDTypeTMDb},
		},
		Artworks: []ArtworkResponse{
			{Type: ArtworkTypePoster, Image: "https://poster.jpg", Width: 680, Height: 1000},
			{Type: ArtworkTypeBackground, Image: "https://backdrop.jpg", Width: 1920, Height: 1080},
		},
		Trailers: []TrailerResponse{
			{URL: "https://youtube.com/watch?v=trailer"},
		},
	}

	result := mapTVShowMetadata(input)

	assert.Equal(t, "73255", result.ProviderID)
	assert.Equal(t, metadata.ProviderTVDb, result.Provider)
	require.NotNil(t, result.TVDbID)
	assert.Equal(t, int32(73255), *result.TVDbID)
	assert.Equal(t, "Breaking Bad", result.Name)
	assert.Equal(t, "eng", result.OriginalLanguage)
	assert.Equal(t, "Ended", result.Status)
	assert.False(t, result.InProduction)
	assert.InDelta(t, 8.5, result.VoteAverage, 0.01)

	require.NotNil(t, result.FirstAirDate)
	require.NotNil(t, result.LastAirDate)

	assert.Equal(t, []int{47}, result.EpisodeRuntime)
	assert.Equal(t, 1, result.NumberOfSeasons)

	require.NotNil(t, result.PosterPath)
	assert.Equal(t, "https://poster.jpg", *result.PosterPath)
	require.NotNil(t, result.BackdropPath)
	assert.Equal(t, "https://backdrop.jpg", *result.BackdropPath)

	require.Len(t, result.Genres, 1)
	assert.Equal(t, "Drama", result.Genres[0].Name)

	require.Len(t, result.Networks, 1)
	assert.Equal(t, "AMC", result.Networks[0].Name)

	require.Len(t, result.Seasons, 1)
	assert.Equal(t, "Season 1", result.Seasons[0].Name)

	require.NotNil(t, result.IMDbID)
	assert.Equal(t, "tt0903747", *result.IMDbID)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(1396), *result.TMDbID)

	require.NotNil(t, result.TrailerURL)
	assert.Equal(t, "https://youtube.com/watch?v=trailer", *result.TrailerURL)

	assert.Equal(t, []string{"usa"}, result.OriginCountries)

	require.Len(t, result.ExternalRatings, 1)
	assert.Equal(t, "TVDb", result.ExternalRatings[0].Source)
}

func TestMapTVShowMetadata_FallbackImage(t *testing.T) {
	input := &SeriesResponse{
		ID:    1,
		Name:  "Test",
		Image: ptr("https://main-image.jpg"),
	}
	result := mapTVShowMetadata(input)
	require.NotNil(t, result.PosterPath)
	assert.Equal(t, "https://main-image.jpg", *result.PosterPath)
}

func TestMapTVShowMetadata_InProduction(t *testing.T) {
	input := &SeriesResponse{
		ID:     1,
		Name:   "Test",
		Status: &StatusResponse{Name: "Continuing", KeepUpdated: true},
	}
	result := mapTVShowMetadata(input)
	assert.True(t, result.InProduction)
}

func TestMapSeasonMetadata(t *testing.T) {
	input := &SeasonResponse{
		ID:       100,
		SeriesID: 73255,
		Number:   1,
		Name:     "Season 1",
		Overview: ptr("The first season"),
		Image:    ptr("https://season1.jpg"),
		Year:     "2008",
		Overviews: map[string]string{
			"deu": "Die erste Staffel",
			"fra": "La premiere saison",
		},
	}

	result := mapSeasonMetadata(input, "73255")

	assert.Equal(t, "100", result.ProviderID)
	assert.Equal(t, metadata.ProviderTVDb, result.Provider)
	require.NotNil(t, result.TVDbID)
	assert.Equal(t, int32(100), *result.TVDbID)
	assert.Equal(t, "73255", result.ShowID)
	assert.Equal(t, 1, result.SeasonNumber)
	assert.Equal(t, "Season 1", result.Name)
	assert.Equal(t, ptr("The first season"), result.Overview)
	assert.Equal(t, ptr("https://season1.jpg"), result.PosterPath)
	require.NotNil(t, result.AirDate)

	require.Len(t, result.Translations, 2)
	assert.Equal(t, "Die erste Staffel", result.Translations["deu"].Overview)
}

func TestMapEpisodeMetadata(t *testing.T) {
	runtime := 47
	input := &EpisodeResponse{
		ID:           1,
		SeriesID:     73255,
		Name:         "Pilot",
		Aired:        "2008-01-20",
		Runtime:      &runtime,
		SeasonNumber: 1,
		Number:       1,
		Image:        ptr("https://still.jpg"),
		Overview:     ptr("Walter White discovers..."),
		Characters: []CharacterResponse{
			{ID: 10, Name: "Walter White", Type: CharacterTypeActor, PersonName: "Bryan Cranston", Sort: 0, PersonImgURL: ptr("https://bryan.jpg")},
			{ID: 20, Name: "", Type: CharacterTypeDirector, PersonName: "Vince Gilligan", Sort: 0},
		},
		Overviews: map[string]string{"deu": "Walter White entdeckt..."},
	}

	result := mapEpisodeMetadata(input, "73255")

	assert.Equal(t, "1", result.ProviderID)
	assert.Equal(t, metadata.ProviderTVDb, result.Provider)
	assert.Equal(t, "73255", result.ShowID)
	assert.Equal(t, 1, result.SeasonNumber)
	assert.Equal(t, 1, result.EpisodeNumber)
	assert.Equal(t, "Pilot", result.Name)
	require.NotNil(t, result.Runtime)
	assert.Equal(t, int32(47), *result.Runtime)
	require.NotNil(t, result.AirDate)

	require.Len(t, result.GuestStars, 1)
	assert.Equal(t, "Bryan Cranston", result.GuestStars[0].Name)
	assert.Equal(t, "Walter White", result.GuestStars[0].Character)

	require.Len(t, result.Crew, 1)
	assert.Equal(t, "Vince Gilligan", result.Crew[0].Name)
	assert.Equal(t, "Director", result.Crew[0].Job)
	assert.Equal(t, "Directing", result.Crew[0].Department)

	require.Len(t, result.Translations, 1)
	assert.Equal(t, "Walter White entdeckt...", result.Translations["deu"].Overview)
}

func TestMapPersonSearchResult(t *testing.T) {
	input := &SearchResult{
		TVDbID:   "12345",
		Name:     "Bryan Cranston",
		ImageURL: ptr("https://image.jpg"),
	}
	result := mapPersonSearchResult(input)

	assert.Equal(t, "12345", result.ProviderID)
	assert.Equal(t, metadata.ProviderTVDb, result.Provider)
	assert.Equal(t, "Bryan Cranston", result.Name)
	assert.Equal(t, ptr("https://image.jpg"), result.ProfilePath)
}

func TestMapPersonMetadata(t *testing.T) {
	input := &PersonResponse{
		ID:         17419,
		Name:       "Bryan Cranston",
		Image:      ptr("https://bryan.jpg"),
		Birth:      ptr("1956-03-07"),
		Death:      nil,
		BirthPlace: ptr("Canoga Park, California"),
		Gender:     2,
		Score:      50,
		Aliases:    []AliasResponse{{Name: "Bryan Lee Cranston"}},
		RemoteIDs: []RemoteIDResponse{
			{ID: "nm0186505", Type: RemoteIDTypeIMDb},
			{ID: "17419", Type: RemoteIDTypeTMDb},
		},
		Biographies: []BiographyResponse{
			{Biography: "Bryan is an actor", Language: "eng"},
			{Biography: "Bryan est un acteur", Language: "fra"},
		},
	}

	result := mapPersonMetadata(input, "eng")

	assert.Equal(t, "17419", result.ProviderID)
	assert.Equal(t, metadata.ProviderTVDb, result.Provider)
	assert.Equal(t, "Bryan Cranston", result.Name)
	assert.Equal(t, 2, result.Gender)
	assert.Equal(t, float64(50), result.Popularity)
	require.NotNil(t, result.Birthday)
	assert.Nil(t, result.Deathday)
	assert.Equal(t, ptr("Canoga Park, California"), result.PlaceOfBirth)

	require.Len(t, result.AlsoKnownAs, 1)
	assert.Equal(t, "Bryan Lee Cranston", result.AlsoKnownAs[0])

	require.NotNil(t, result.Biography)
	assert.Equal(t, "Bryan is an actor", *result.Biography)

	require.Len(t, result.Translations, 2)

	require.NotNil(t, result.IMDbID)
	assert.Equal(t, "nm0186505", *result.IMDbID)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(17419), *result.TMDbID)
}

func TestMapPersonMetadata_FallbackBiography(t *testing.T) {
	input := &PersonResponse{
		ID:   1,
		Name: "Test",
		Biographies: []BiographyResponse{
			{Biography: "English bio", Language: "eng"},
		},
	}
	result := mapPersonMetadata(input, "deu")
	require.NotNil(t, result.Biography)
	assert.Equal(t, "English bio", *result.Biography)
}

func TestMapPersonCredits(t *testing.T) {
	seriesID := 73255
	movieID := 550
	input := &PersonResponse{
		ID:   17419,
		Name: "Bryan Cranston",
		Characters: []CharacterResponse{
			{ID: 1, Name: "Walter White", Type: CharacterTypeActor, SeriesID: &seriesID, Image: ptr("https://img.jpg"), PersonName: "Bryan Cranston"},
			{ID: 2, Name: "", Type: CharacterTypeDirector, MovieID: &movieID, PersonName: "Bryan Cranston"},
		},
	}

	result := mapPersonCredits(input)

	assert.Equal(t, "17419", result.ProviderID)
	assert.Equal(t, metadata.ProviderTVDb, result.Provider)

	require.Len(t, result.CastCredits, 1)
	assert.Equal(t, "tv", result.CastCredits[0].MediaType)
	assert.Equal(t, "73255", result.CastCredits[0].MediaID)

	require.Len(t, result.CrewCredits, 1)
	assert.Equal(t, "movie", result.CrewCredits[0].MediaType)
	assert.Equal(t, "550", result.CrewCredits[0].MediaID)
}

func TestMapCharactersToCredits(t *testing.T) {
	characters := []CharacterResponse{
		{ID: 1, Name: "Walter White", Type: CharacterTypeActor, PersonName: "Bryan Cranston", Sort: 0},
		{ID: 2, Name: "Jesse Pinkman", Type: CharacterTypeActor, PersonName: "Aaron Paul", Sort: 1},
		{ID: 3, Name: "", Type: CharacterTypeWriter, PersonName: "Vince Gilligan", Sort: 0},
	}

	result := mapCharactersToCredits(characters)

	require.Len(t, result.Cast, 2)
	assert.Equal(t, "Bryan Cranston", result.Cast[0].Name)
	assert.Equal(t, "Walter White", result.Cast[0].Character)
	assert.Equal(t, "Aaron Paul", result.Cast[1].Name)

	require.Len(t, result.Crew, 1)
	assert.Equal(t, "Writer", result.Crew[0].Job)
	assert.Equal(t, "Writing", result.Crew[0].Department)
}

func TestMapArtworksToImages(t *testing.T) {
	artworks := []ArtworkResponse{
		{Type: ArtworkTypePoster, Image: "poster.jpg", Width: 680, Height: 1000, Score: 5, Language: ptr("en")},
		{Type: ArtworkTypeBackground, Image: "bg.jpg", Width: 1920, Height: 1080, Score: 8},
		{Type: ArtworkTypeBanner, Image: "banner.jpg", Width: 758, Height: 140},
		{Type: ArtworkTypeClearLogo, Image: "logo.jpg", Width: 500, Height: 185},
		{Type: ArtworkTypeClearArt, Image: "art.png", Width: 1000, Height: 562},
		{Type: ArtworkTypeIcon, Image: "icon.jpg", Width: 64, Height: 64},
	}

	result := mapArtworksToImages(artworks)

	require.Len(t, result.Posters, 1)
	assert.Equal(t, "poster.jpg", result.Posters[0].FilePath)
	assert.InDelta(t, 0.68, result.Posters[0].AspectRatio, 0.01)
	assert.Equal(t, ptr("en"), result.Posters[0].Language)

	require.Len(t, result.Backdrops, 1)
	assert.Equal(t, "bg.jpg", result.Backdrops[0].FilePath)

	require.Len(t, result.Logos, 3)
}

func TestMapArtworksToImages_ZeroWidthHeight(t *testing.T) {
	artworks := []ArtworkResponse{
		{Type: ArtworkTypePoster, Image: "poster.jpg", Width: 0, Height: 0},
	}
	result := mapArtworksToImages(artworks)
	assert.Equal(t, float64(0), result.Posters[0].AspectRatio)
}

func TestMapContentRatings(t *testing.T) {
	ratings := []ContentRatingResponse{
		{ID: 1, Name: "TV-MA", Country: "usa"},
		{ID: 2, Name: "16+", Country: "deu"},
	}
	result := mapContentRatings(ratings)
	require.Len(t, result, 2)
	assert.Equal(t, "usa", result[0].CountryCode)
	assert.Equal(t, "TV-MA", result[0].Rating)
}

func TestMapOverviewsToTranslations(t *testing.T) {
	overviews := map[string]string{
		"eng": "English overview",
		"deu": "German overview",
	}
	result := mapOverviewsToTranslations(overviews, nil)
	require.Len(t, result, 2)

	found := map[string]bool{}
	for _, tr := range result {
		found[tr.Language] = true
		require.NotNil(t, tr.Data)
	}
	assert.True(t, found["eng"])
	assert.True(t, found["deu"])
}

func TestMapRemoteIDsToExternalIDs(t *testing.T) {
	remoteIDs := []RemoteIDResponse{
		{ID: "tt0903747", Type: RemoteIDTypeIMDb},
		{ID: "1396", Type: RemoteIDTypeTMDb},
		{ID: "unknown", Type: 99},
	}
	result := mapRemoteIDsToExternalIDs(remoteIDs, 73255)

	require.NotNil(t, result.TVDbID)
	assert.Equal(t, int32(73255), *result.TVDbID)
	require.NotNil(t, result.IMDbID)
	assert.Equal(t, "tt0903747", *result.IMDbID)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(1396), *result.TMDbID)
}

func TestMapRemoteIDsToExternalIDs_InvalidTMDb(t *testing.T) {
	remoteIDs := []RemoteIDResponse{
		{ID: "not-a-number", Type: RemoteIDTypeTMDb},
	}
	result := mapRemoteIDsToExternalIDs(remoteIDs, 1)
	assert.Nil(t, result.TMDbID)
}

func TestCacheEntry_IsExpired(t *testing.T) {
	expired := &CacheEntry{ExpiresAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	assert.True(t, expired.IsExpired())

	notExpired := &CacheEntry{ExpiresAt: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)}
	assert.False(t, notExpired.IsExpired())
}
