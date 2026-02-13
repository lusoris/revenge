package tmdb

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
		{"simple date", "2024-06-15", false, time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)},
		{"ISO 8601 with time", "2024-06-15T10:30:00Z", false, time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)},
		{"invalid date", "not-a-date", true, time.Time{}},
		{"partial date", "2024", true, time.Time{}},
		{"date with timezone offset", "2024-06-15T10:30:00+02:00", false, time.Date(2024, 6, 15, 10, 30, 0, 0, time.FixedZone("", 7200))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDate(tt.input)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.True(t, tt.wantDate.Equal(*result), "expected %v, got %v", tt.wantDate, *result)
			}
		})
	}
}

func TestExtractYear(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		want    int
	}{
		{"empty string", "", true, 0},
		{"too short", "202", true, 0},
		{"valid date", "2024-06-15", false, 2024},
		{"just year", "2024", false, 2024},
		{"invalid year prefix", "abcd-01-01", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractYear(tt.input)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.want, *result)
			}
		})
	}
}

func TestMapMovieSearchResult(t *testing.T) {
	input := &MovieSearchResponse{
		ID:               550,
		Title:            "Fight Club",
		OriginalTitle:    "Fight Club",
		OriginalLanguage: "en",
		Overview:         "A ticking-Loss bomb insomniac...",
		ReleaseDate:      "1999-10-15",
		PosterPath:       new("/poster.jpg"),
		BackdropPath:     new("/backdrop.jpg"),
		VoteAverage:      8.4,
		VoteCount:        25000,
		Popularity:       60.5,
		Adult:            false,
		Video:            false,
		GenreIDs:         []int{18, 53},
	}

	result := mapMovieSearchResult(input)

	assert.Equal(t, "550", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	assert.Equal(t, "Fight Club", result.Title)
	assert.Equal(t, "Fight Club", result.OriginalTitle)
	assert.Equal(t, "en", result.OriginalLanguage)
	assert.Equal(t, "A ticking-Loss bomb insomniac...", result.Overview)
	assert.Equal(t, new("/poster.jpg"), result.PosterPath)
	assert.Equal(t, new("/backdrop.jpg"), result.BackdropPath)
	assert.Equal(t, 8.4, result.VoteAverage)
	assert.Equal(t, 25000, result.VoteCount)
	assert.Equal(t, 60.5, result.Popularity)
	assert.False(t, result.Adult)
	assert.Equal(t, []int{18, 53}, result.GenreIDs)
	require.NotNil(t, result.ReleaseDate)
	assert.Equal(t, 1999, result.ReleaseDate.Year())
	require.NotNil(t, result.Year)
	assert.Equal(t, 1999, *result.Year)
}

func TestMapMovieSearchResult_NoReleaseDate(t *testing.T) {
	input := &MovieSearchResponse{
		ID:    123,
		Title: "No Date Movie",
	}

	result := mapMovieSearchResult(input)

	assert.Equal(t, "123", result.ProviderID)
	assert.Nil(t, result.ReleaseDate)
	assert.Nil(t, result.Year)
}

func TestMapMovieSearchResults(t *testing.T) {
	resp := &SearchResultsResponse{
		Page: 1,
		Results: []MovieSearchResponse{
			{ID: 1, Title: "Movie 1"},
			{ID: 2, Title: "Movie 2"},
			{ID: 3, Title: "Movie 3"},
		},
		TotalPages:   1,
		TotalResults: 3,
	}

	results := mapMovieSearchResults(resp)

	assert.Len(t, results, 3)
	assert.Equal(t, "1", results[0].ProviderID)
	assert.Equal(t, "Movie 1", results[0].Title)
	assert.Equal(t, "3", results[2].ProviderID)
}

func TestMapMovieMetadata(t *testing.T) {
	runtime := 139
	budget := int64(63000000)
	revenue := int64(100853753)
	input := &MovieResponse{
		ID:               550,
		IMDbID:           new("tt0137523"),
		Title:            "Fight Club",
		OriginalTitle:    "Fight Club",
		OriginalLanguage: "en",
		Overview:         new("An insomniac office worker..."),
		Tagline:          new("Mischief. Mayhem. Soap."),
		ReleaseDate:      "1999-10-15",
		Runtime:          &runtime,
		Budget:           &budget,
		Revenue:          &revenue,
		Status:           "Released",
		VoteAverage:      8.4,
		VoteCount:        25000,
		Popularity:       60.5,
		Adult:            false,
		PosterPath:       new("/poster.jpg"),
		BackdropPath:     new("/backdrop.jpg"),
		Homepage:         new("https://fightclub.example.com"),
		Genres:           []GenreResponse{{ID: 18, Name: "Drama"}, {ID: 53, Name: "Thriller"}},
		ProductionCompanies: []CompanyResponse{
			{ID: 508, Name: "Regency Enterprises", LogoPath: new("/logo.png"), OriginCountry: "US"},
		},
		ProductionCountries: []CountryResponse{
			{ISO3166_1: "US", Name: "United States"},
		},
		SpokenLanguages: []LanguageResponse{
			{ISO639_1: "en", Name: "English", EnglishName: "English"},
		},
		BelongsToCollection: &CollectionRefResponse{
			ID:           123,
			Name:         "Fight Club Collection",
			PosterPath:   new("/col_poster.jpg"),
			BackdropPath: new("/col_backdrop.jpg"),
		},
		Videos: &VideosResponse{
			Results: []VideoResponse{
				{Key: "SUXWAEX2jlg", Site: "YouTube", Type: "Trailer", Official: true},
			},
		},
	}

	result := mapMovieMetadata(input)

	assert.Equal(t, "550", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(550), *result.TMDbID)
	assert.Equal(t, new("tt0137523"), result.IMDbID)
	assert.Equal(t, "Fight Club", result.Title)
	assert.Equal(t, "Fight Club", result.OriginalTitle)
	assert.Equal(t, "en", result.OriginalLanguage)
	assert.Equal(t, new("Mischief. Mayhem. Soap."), result.Tagline)
	assert.Equal(t, new("An insomniac office worker..."), result.Overview)
	assert.Equal(t, "Released", result.Status)
	assert.Equal(t, 8.4, result.VoteAverage)
	assert.Equal(t, 25000, result.VoteCount)
	assert.Equal(t, 60.5, result.Popularity)
	assert.False(t, result.Adult)
	assert.Equal(t, new("/poster.jpg"), result.PosterPath)
	assert.Equal(t, new("/backdrop.jpg"), result.BackdropPath)

	// Release date
	require.NotNil(t, result.ReleaseDate)
	assert.Equal(t, 1999, result.ReleaseDate.Year())

	// Runtime
	require.NotNil(t, result.Runtime)
	assert.Equal(t, int32(139), *result.Runtime)

	// Budget & Revenue
	require.NotNil(t, result.Budget)
	assert.Equal(t, int64(63000000), *result.Budget)
	require.NotNil(t, result.Revenue)
	assert.Equal(t, int64(100853753), *result.Revenue)

	// Genres
	require.Len(t, result.Genres, 2)
	assert.Equal(t, 18, result.Genres[0].ID)
	assert.Equal(t, "Drama", result.Genres[0].Name)

	// Production companies
	require.Len(t, result.ProductionCompanies, 1)
	assert.Equal(t, "Regency Enterprises", result.ProductionCompanies[0].Name)

	// Production countries
	require.Len(t, result.ProductionCountries, 1)
	assert.Equal(t, "US", result.ProductionCountries[0].ISOCode)

	// Spoken languages
	require.Len(t, result.SpokenLanguages, 1)
	assert.Equal(t, "en", result.SpokenLanguages[0].ISOCode)

	// Collection
	require.NotNil(t, result.Collection)
	assert.Equal(t, "Fight Club Collection", result.Collection.Name)
	assert.Equal(t, 123, result.Collection.ID)

	// Trailer
	require.NotNil(t, result.TrailerURL)
	assert.Equal(t, "https://www.youtube.com/watch?v=SUXWAEX2jlg", *result.TrailerURL)

	// External ratings
	require.Len(t, result.ExternalRatings, 1)
	assert.Equal(t, "TMDb", result.ExternalRatings[0].Source)
	assert.Equal(t, "8.4/10", result.ExternalRatings[0].Value)
	assert.InDelta(t, 84.0, result.ExternalRatings[0].Score, 0.01)
}

func TestMapMovieMetadata_ZeroVoteAverage(t *testing.T) {
	input := &MovieResponse{ID: 1, Title: "Test", VoteAverage: 0}
	result := mapMovieMetadata(input)
	assert.Empty(t, result.ExternalRatings)
}

func TestMapMovieMetadata_ZeroRuntimeBudgetRevenue(t *testing.T) {
	zeroRuntime := 0
	zeroBudget := int64(0)
	zeroRevenue := int64(0)
	input := &MovieResponse{
		ID:      1,
		Title:   "Test",
		Runtime: &zeroRuntime,
		Budget:  &zeroBudget,
		Revenue: &zeroRevenue,
	}
	result := mapMovieMetadata(input)
	assert.Nil(t, result.Runtime)
	assert.Nil(t, result.Budget)
	assert.Nil(t, result.Revenue)
}

func TestMapMovieMetadata_TeaserVideoWhenNoTrailer(t *testing.T) {
	input := &MovieResponse{
		ID:    1,
		Title: "Test",
		Videos: &VideosResponse{
			Results: []VideoResponse{
				{Key: "abc", Site: "Vimeo", Type: "Trailer"},
				{Key: "xyz", Site: "YouTube", Type: "Teaser"},
			},
		},
	}
	result := mapMovieMetadata(input)
	require.NotNil(t, result.TrailerURL)
	assert.Equal(t, "https://www.youtube.com/watch?v=xyz", *result.TrailerURL)
}

func TestMapTVSearchResult(t *testing.T) {
	input := &TVSearchResponse{
		ID:               1399,
		Name:             "Breaking Bad",
		OriginalName:     "Breaking Bad",
		OriginalLanguage: "en",
		Overview:         "A high school chemistry teacher...",
		FirstAirDate:     "2008-01-20",
		PosterPath:       new("/poster.jpg"),
		BackdropPath:     new("/backdrop.jpg"),
		VoteAverage:      8.9,
		VoteCount:        12000,
		Popularity:       120.3,
		Adult:            false,
		GenreIDs:         []int{18, 80},
		OriginCountry:    []string{"US"},
	}

	result := mapTVSearchResult(input)

	assert.Equal(t, "1399", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	assert.Equal(t, "Breaking Bad", result.Name)
	assert.Equal(t, "Breaking Bad", result.OriginalName)
	assert.Equal(t, "en", result.OriginalLanguage)
	assert.Equal(t, 8.9, result.VoteAverage)
	assert.Equal(t, []string{"US"}, result.OriginCountries)
	require.NotNil(t, result.FirstAirDate)
	assert.Equal(t, 2008, result.FirstAirDate.Year())
	require.NotNil(t, result.Year)
	assert.Equal(t, 2008, *result.Year)
}

func TestMapTVShowMetadata(t *testing.T) {
	input := &TVResponse{
		ID:               1399,
		Name:             "Breaking Bad",
		OriginalName:     "Breaking Bad",
		OriginalLanguage: "en",
		Overview:         new("A high school chemistry teacher..."),
		Tagline:          new("All Hail the King"),
		Status:           "Ended",
		Type:             "Scripted",
		FirstAirDate:     "2008-01-20",
		LastAirDate:      "2013-09-29",
		InProduction:     false,
		NumberOfSeasons:  5,
		NumberOfEpisodes: 62,
		EpisodeRunTime:   []int{45, 47},
		VoteAverage:      8.9,
		VoteCount:        12000,
		Popularity:       120.3,
		PosterPath:       new("/poster.jpg"),
		BackdropPath:     new("/backdrop.jpg"),
		Homepage:         new("https://breakingbad.example.com"),
		Genres:           []GenreResponse{{ID: 18, Name: "Drama"}, {ID: 80, Name: "Crime"}},
		Networks:         []NetworkResponse{{ID: 174, Name: "AMC", LogoPath: new("/amc.png"), OriginCountry: "US"}},
		CreatedBy:        []CreatorResponse{{ID: 66633, Name: "Vince Gilligan", Gender: 2, CreditID: "cr1"}},
		OriginCountry:    []string{"US"},
		SpokenLanguages:  []LanguageResponse{{ISO639_1: "en", Name: "English", EnglishName: "English"}},
		Seasons: []SeasonSummaryResponse{
			{ID: 3572, Name: "Season 1", SeasonNumber: 1, AirDate: "2008-01-20", EpisodeCount: 7, VoteAverage: 8.3},
		},
		Videos: &VideosResponse{
			Results: []VideoResponse{
				{Key: "trailer1", Site: "YouTube", Type: "Trailer"},
			},
		},
	}

	result := mapTVShowMetadata(input)

	assert.Equal(t, "1399", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(1399), *result.TMDbID)
	assert.Equal(t, "Breaking Bad", result.Name)
	assert.Equal(t, "Ended", result.Status)
	assert.Equal(t, "Scripted", result.Type)
	assert.False(t, result.InProduction)
	assert.Equal(t, 5, result.NumberOfSeasons)
	assert.Equal(t, 62, result.NumberOfEpisodes)
	assert.Equal(t, []int{45, 47}, result.EpisodeRuntime)

	require.NotNil(t, result.FirstAirDate)
	assert.Equal(t, 2008, result.FirstAirDate.Year())
	require.NotNil(t, result.LastAirDate)
	assert.Equal(t, 2013, result.LastAirDate.Year())

	// Genres
	require.Len(t, result.Genres, 2)
	assert.Equal(t, "Drama", result.Genres[0].Name)

	// Networks
	require.Len(t, result.Networks, 1)
	assert.Equal(t, "AMC", result.Networks[0].Name)

	// CreatedBy
	require.Len(t, result.CreatedBy, 1)
	assert.Equal(t, "Vince Gilligan", result.CreatedBy[0].Name)

	// Seasons
	require.Len(t, result.Seasons, 1)
	assert.Equal(t, "Season 1", result.Seasons[0].Name)
	assert.Equal(t, 1, result.Seasons[0].SeasonNumber)
	assert.Equal(t, 7, result.Seasons[0].EpisodeCount)
	require.NotNil(t, result.Seasons[0].AirDate)

	// Trailer
	require.NotNil(t, result.TrailerURL)
	assert.Equal(t, "https://www.youtube.com/watch?v=trailer1", *result.TrailerURL)

	// External ratings
	require.Len(t, result.ExternalRatings, 1)
	assert.Equal(t, "TMDb", result.ExternalRatings[0].Source)
}

func TestMapSeasonMetadata(t *testing.T) {
	runtime45 := 45
	input := &SeasonResponse{
		ID:           3572,
		Name:         "Season 1",
		Overview:     new("The first season..."),
		SeasonNumber: 1,
		AirDate:      "2008-01-20",
		PosterPath:   new("/season1.jpg"),
		VoteAverage:  8.3,
		Episodes: []EpisodeSummaryResponse{
			{
				ID:            62085,
				Name:          "Pilot",
				Overview:      new("Walter White, a high school chemistry teacher..."),
				SeasonNumber:  1,
				EpisodeNumber: 1,
				AirDate:       "2008-01-20",
				Runtime:       &runtime45,
				StillPath:     new("/still1.jpg"),
				VoteAverage:   8.5,
				VoteCount:     500,
			},
			{
				ID:            62086,
				Name:          "Cat's in the Bag...",
				EpisodeNumber: 2,
				SeasonNumber:  1,
				AirDate:       "",
			},
		},
	}

	result := mapSeasonMetadata(input, "1399")

	assert.Equal(t, "3572", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(3572), *result.TMDbID)
	assert.Equal(t, "1399", result.ShowID)
	assert.Equal(t, 1, result.SeasonNumber)
	assert.Equal(t, "Season 1", result.Name)
	require.NotNil(t, result.AirDate)

	// Episodes
	require.Len(t, result.Episodes, 2)
	ep1 := result.Episodes[0]
	assert.Equal(t, "62085", ep1.ProviderID)
	assert.Equal(t, 1, ep1.EpisodeNumber)
	assert.Equal(t, "Pilot", ep1.Name)
	require.NotNil(t, ep1.Runtime)
	assert.Equal(t, int32(45), *ep1.Runtime)
	require.NotNil(t, ep1.AirDate)

	ep2 := result.Episodes[1]
	assert.Nil(t, ep2.AirDate)
	assert.Nil(t, ep2.Runtime)
}

func TestMapEpisodeMetadata(t *testing.T) {
	runtime47 := 47
	input := &EpisodeResponse{
		ID:            62085,
		Name:          "Pilot",
		Overview:      new("Walter White..."),
		SeasonNumber:  1,
		EpisodeNumber: 1,
		AirDate:       "2008-01-20",
		Runtime:       &runtime47,
		StillPath:     new("/still.jpg"),
		VoteAverage:   8.5,
		VoteCount:     500,
		GuestStars: []CastResponse{
			{ID: 101, Name: "John Doe", Character: "DEA Agent", Order: 0, CreditID: "cr1", Gender: 2, ProfilePath: new("/john.jpg")},
		},
		Crew: []CrewResponse{
			{ID: 202, Name: "Vince Gilligan", Job: "Director", Department: "Directing", CreditID: "cr2", Gender: 2},
		},
	}

	result := mapEpisodeMetadata(input, "1399")

	assert.Equal(t, "62085", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(62085), *result.TMDbID)
	assert.Equal(t, "1399", result.ShowID)
	assert.Equal(t, 1, result.SeasonNumber)
	assert.Equal(t, 1, result.EpisodeNumber)
	assert.Equal(t, "Pilot", result.Name)
	require.NotNil(t, result.Runtime)
	assert.Equal(t, int32(47), *result.Runtime)
	require.NotNil(t, result.AirDate)

	// Guest stars
	require.Len(t, result.GuestStars, 1)
	assert.Equal(t, "101", result.GuestStars[0].ProviderID)
	assert.Equal(t, "John Doe", result.GuestStars[0].Name)
	assert.Equal(t, "DEA Agent", result.GuestStars[0].Character)

	// Crew
	require.Len(t, result.Crew, 1)
	assert.Equal(t, "202", result.Crew[0].ProviderID)
	assert.Equal(t, "Director", result.Crew[0].Job)
	assert.Equal(t, "Directing", result.Crew[0].Department)
}

func TestMapPersonSearchResult(t *testing.T) {
	input := &PersonSearchResponse{
		ID:          17419,
		Name:        "Bryan Cranston",
		ProfilePath: new("/bryan.jpg"),
		Popularity:  40.5,
		Adult:       false,
		KnownFor: []KnownForResponse{
			{MediaType: "tv", ID: 1399, Name: "Breaking Bad", PosterPath: new("/bb.jpg")},
			{MediaType: "movie", ID: 550, Title: "Fight Club", PosterPath: new("/fc.jpg")},
		},
	}

	result := mapPersonSearchResult(input)

	assert.Equal(t, "17419", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	assert.Equal(t, "Bryan Cranston", result.Name)
	assert.Equal(t, new("/bryan.jpg"), result.ProfilePath)
	assert.False(t, result.Adult)

	require.Len(t, result.KnownFor, 2)
	assert.Equal(t, "tv", result.KnownFor[0].MediaType)
	assert.Equal(t, "Breaking Bad", result.KnownFor[0].Title)
	assert.Equal(t, "movie", result.KnownFor[1].MediaType)
	assert.Equal(t, "Fight Club", result.KnownFor[1].Title)
}

func TestMapPersonMetadata(t *testing.T) {
	input := &PersonResponse{
		ID:           17419,
		Name:         "Bryan Cranston",
		AlsoKnownAs:  []string{"Bryan Lee Cranston"},
		Biography:    new("Bryan Lee Cranston is an American actor..."),
		Birthday:     new("1956-03-07"),
		Deathday:     nil,
		Gender:       2,
		PlaceOfBirth: new("Canoga Park, California"),
		ProfilePath:  new("/bryan.jpg"),
		Homepage:     new("https://example.com"),
		Popularity:   40.5,
		Adult:        false,
		IMDbID:       new("nm0186505"),
		KnownForDept: "Acting",
	}

	result := mapPersonMetadata(input)

	assert.Equal(t, "17419", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(17419), *result.TMDbID)
	assert.Equal(t, "Bryan Cranston", result.Name)
	assert.Equal(t, []string{"Bryan Lee Cranston"}, result.AlsoKnownAs)
	assert.Equal(t, 2, result.Gender)
	assert.Equal(t, "Acting", result.KnownForDept)
	require.NotNil(t, result.Birthday)
	assert.Equal(t, 1956, result.Birthday.Year())
	assert.Nil(t, result.Deathday)
}

func TestMapPersonMetadata_WithDeathday(t *testing.T) {
	input := &PersonResponse{
		ID:       12345,
		Name:     "Test Person",
		Birthday: new("1950-01-01"),
		Deathday: new("2020-12-31"),
	}

	result := mapPersonMetadata(input)
	require.NotNil(t, result.Birthday)
	require.NotNil(t, result.Deathday)
	assert.Equal(t, 2020, result.Deathday.Year())
}

func TestMapPersonMetadata_EmptyBirthdayDeathday(t *testing.T) {
	input := &PersonResponse{
		ID:       12345,
		Name:     "Test",
		Birthday: new(""),
		Deathday: new(""),
	}
	result := mapPersonMetadata(input)
	assert.Nil(t, result.Birthday)
	assert.Nil(t, result.Deathday)
}

func TestMapPersonCredits(t *testing.T) {
	episodeCount := 62
	input := &PersonCreditsResponse{
		ID: 17419,
		Cast: []PersonCastCredit{
			{
				MediaType:    "tv",
				ID:           1399,
				Name:         "Breaking Bad",
				Character:    new("Walter White"),
				PosterPath:   new("/bb.jpg"),
				VoteAverage:  8.9,
				FirstAirDate: "2008-01-20",
				EpisodeCount: &episodeCount,
			},
			{
				MediaType:   "movie",
				ID:          550,
				Title:       "Fight Club",
				Character:   new("Tyler Durden"),
				ReleaseDate: "1999-10-15",
			},
		},
		Crew: []PersonCrewCredit{
			{
				MediaType:   "movie",
				ID:          999,
				Title:       "Some Movie",
				Job:         new("Director"),
				Department:  new("Directing"),
				ReleaseDate: "2020-01-01",
			},
			{
				MediaType:    "tv",
				ID:           2000,
				Name:         "Some Show",
				Job:          new("Producer"),
				Department:   new("Production"),
				FirstAirDate: "2018-06-01",
			},
		},
	}

	result := mapPersonCredits(input)

	assert.Equal(t, "17419", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)

	// Cast credits
	require.Len(t, result.CastCredits, 2)
	assert.Equal(t, "tv", result.CastCredits[0].MediaType)
	assert.Equal(t, "Breaking Bad", result.CastCredits[0].Title)
	require.NotNil(t, result.CastCredits[0].ReleaseDate)
	assert.Equal(t, 2008, result.CastCredits[0].ReleaseDate.Year())
	assert.Equal(t, &episodeCount, result.CastCredits[0].EpisodeCount)

	assert.Equal(t, "movie", result.CastCredits[1].MediaType)
	assert.Equal(t, "Fight Club", result.CastCredits[1].Title)
	require.NotNil(t, result.CastCredits[1].ReleaseDate)
	assert.Equal(t, 1999, result.CastCredits[1].ReleaseDate.Year())

	// Crew credits
	require.Len(t, result.CrewCredits, 2)
	assert.Equal(t, "movie", result.CrewCredits[0].MediaType)
	assert.Equal(t, "Some Movie", result.CrewCredits[0].Title)
	assert.Equal(t, "tv", result.CrewCredits[1].MediaType)
	assert.Equal(t, "Some Show", result.CrewCredits[1].Title)
	require.NotNil(t, result.CrewCredits[1].ReleaseDate)
}

func TestMapCredits(t *testing.T) {
	input := &CreditsResponse{
		ID: 550,
		Cast: []CastResponse{
			{ID: 1, Name: "Brad Pitt", Character: "Tyler Durden", Order: 0, CreditID: "cr1", Gender: 2, ProfilePath: new("/brad.jpg")},
			{ID: 2, Name: "Edward Norton", Character: "Narrator", Order: 1, CreditID: "cr2", Gender: 2},
		},
		Crew: []CrewResponse{
			{ID: 3, Name: "David Fincher", Job: "Director", Department: "Directing", CreditID: "cr3", Gender: 2},
		},
	}

	result := mapCredits(input)

	require.Len(t, result.Cast, 2)
	assert.Equal(t, "1", result.Cast[0].ProviderID)
	assert.Equal(t, "Brad Pitt", result.Cast[0].Name)
	assert.Equal(t, "Tyler Durden", result.Cast[0].Character)
	assert.Equal(t, 0, result.Cast[0].Order)
	assert.Equal(t, new("cr1"), result.Cast[0].CreditID)

	require.Len(t, result.Crew, 1)
	assert.Equal(t, "3", result.Crew[0].ProviderID)
	assert.Equal(t, "David Fincher", result.Crew[0].Name)
	assert.Equal(t, "Director", result.Crew[0].Job)
}

func TestMapImages(t *testing.T) {
	input := &ImagesResponse{
		ID: 550,
		Posters: []ImageResponse{
			{FilePath: "/poster1.jpg", AspectRatio: 0.667, Width: 500, Height: 750, VoteAverage: 5.3, VoteCount: 10, ISO639_1: new("en")},
		},
		Backdrops: []ImageResponse{
			{FilePath: "/backdrop1.jpg", AspectRatio: 1.778, Width: 1920, Height: 1080},
		},
		Logos: []ImageResponse{
			{FilePath: "/logo1.png", Width: 500, Height: 200},
		},
		Stills: []ImageResponse{
			{FilePath: "/still1.jpg", Width: 1280, Height: 720},
		},
	}

	result := mapImages(input)

	require.Len(t, result.Posters, 1)
	assert.Equal(t, "/poster1.jpg", result.Posters[0].FilePath)
	assert.Equal(t, 500, result.Posters[0].Width)
	assert.Equal(t, 750, result.Posters[0].Height)
	assert.Equal(t, new("en"), result.Posters[0].Language)

	require.Len(t, result.Backdrops, 1)
	assert.Equal(t, "/backdrop1.jpg", result.Backdrops[0].FilePath)

	require.Len(t, result.Logos, 1)
	require.Len(t, result.Stills, 1)
}

func TestMapPersonImages(t *testing.T) {
	input := &PersonImagesResponse{
		ID: 17419,
		Profiles: []ImageResponse{
			{FilePath: "/profile1.jpg", Width: 300, Height: 450},
			{FilePath: "/profile2.jpg", Width: 600, Height: 900},
		},
	}

	result := mapPersonImages(input)

	require.Len(t, result.Profiles, 2)
	assert.Equal(t, "/profile1.jpg", result.Profiles[0].FilePath)
	assert.Equal(t, "/profile2.jpg", result.Profiles[1].FilePath)
}

func TestMapImage(t *testing.T) {
	input := &ImageResponse{
		FilePath:    "/img.jpg",
		AspectRatio: 1.5,
		Width:       1200,
		Height:      800,
		VoteAverage: 5.5,
		VoteCount:   42,
		ISO639_1:    new("de"),
	}

	result := mapImage(input)

	assert.Equal(t, "/img.jpg", result.FilePath)
	assert.Equal(t, 1.5, result.AspectRatio)
	assert.Equal(t, 1200, result.Width)
	assert.Equal(t, 800, result.Height)
	assert.Equal(t, 5.5, result.VoteAverage)
	assert.Equal(t, 42, result.VoteCount)
	assert.Equal(t, new("de"), result.Language)
}

func TestMapReleaseDates(t *testing.T) {
	input := &ReleaseDatesWrapper{
		ID: 550,
		Results: []CountryReleaseResponse{
			{
				ISO3166_1: "US",
				ReleaseDates: []ReleaseDateResponse{
					{Certification: "R", ISO639_1: "en", ReleaseDate: "1999-10-15T00:00:00.000Z", Type: 3, Note: "wide release"},
					{Certification: "R", ReleaseDate: "", Type: 1},
				},
			},
			{
				ISO3166_1: "DE",
				ReleaseDates: []ReleaseDateResponse{
					{Certification: "16", ReleaseDate: "1999-11-11", Type: 3},
				},
			},
		},
	}

	results := mapReleaseDates(input)

	require.Len(t, results, 3)
	assert.Equal(t, "US", results[0].CountryCode)
	assert.Equal(t, "R", results[0].Certification)
	assert.Equal(t, 3, results[0].ReleaseType)
	assert.Equal(t, "wide release", results[0].Note)
	require.NotNil(t, results[0].ReleaseDate)

	assert.Nil(t, results[1].ReleaseDate)

	assert.Equal(t, "DE", results[2].CountryCode)
	assert.Equal(t, "16", results[2].Certification)
}

func TestMapContentRatings(t *testing.T) {
	input := &ContentRatingsWrapper{
		ID: 1399,
		Results: []ContentRatingResponse{
			{ISO3166_1: "US", Rating: "TV-MA", Descriptors: []string{"Violence", "Language"}},
			{ISO3166_1: "DE", Rating: "16", Descriptors: nil},
		},
	}

	results := mapContentRatings(input)

	require.Len(t, results, 2)
	assert.Equal(t, "US", results[0].CountryCode)
	assert.Equal(t, "TV-MA", results[0].Rating)
	assert.Equal(t, []string{"Violence", "Language"}, results[0].Descriptors)
	assert.Equal(t, "DE", results[1].CountryCode)
}

func TestMapTranslations(t *testing.T) {
	runtime := 139
	input := &TranslationsWrapper{
		ID: 550,
		Translations: []TranslationResponse{
			{
				ISO3166_1:   "DE",
				ISO639_1:    "de",
				Name:        "Deutsch",
				EnglishName: "German",
				Data: TranslationDataResponse{
					Title:    "Kampfclub",
					Overview: "Ein Büroangestellter...",
					Tagline:  "Unfug. Chaos. Seife.",
					Runtime:  &runtime,
				},
			},
			{
				ISO3166_1:   "JP",
				ISO639_1:    "ja",
				Name:        "日本語",
				EnglishName: "Japanese",
				Data: TranslationDataResponse{
					Name:     "ファイト・クラブ",
					Overview: "あるサラリーマンが...",
				},
			},
		},
	}

	results := mapTranslations(input)

	require.Len(t, results, 2)

	de := results[0]
	assert.Equal(t, "DE", de.ISOCode)
	assert.Equal(t, "de", de.Language)
	assert.Equal(t, "Deutsch", de.Name)
	assert.Equal(t, "German", de.EnglishName)
	require.NotNil(t, de.Data)
	assert.Equal(t, "Kampfclub", de.Data.Title)
	assert.Equal(t, "Unfug. Chaos. Seife.", de.Data.Tagline)
	require.NotNil(t, de.Data.Runtime)
	assert.Equal(t, int32(139), *de.Data.Runtime)

	jp := results[1]
	assert.Equal(t, "ファイト・クラブ", jp.Data.Title)
}

func TestMapExternalIDs(t *testing.T) {
	input := &ExternalIDsResponse{
		ID:          550,
		IMDbID:      new("tt0137523"),
		TVDbID:      new(73255),
		WikidataID:  new("Q190050"),
		FacebookID:  new("FightClub"),
		InstagramID: new("fightclub"),
		TwitterID:   new("fightclub"),
		TikTokID:    new("fightclub"),
		YouTubeID:   new("UC123"),
		FreebaseID:  new("/en/fight_club"),
		FreebaseMID: new("/m/0bth"),
		TVRageID:    new(7926),
	}

	result := mapExternalIDs(input, 550)

	require.NotNil(t, result.TMDbID)
	assert.Equal(t, int32(550), *result.TMDbID)
	assert.Equal(t, new("tt0137523"), result.IMDbID)
	require.NotNil(t, result.TVDbID)
	assert.Equal(t, int32(73255), *result.TVDbID)
	require.NotNil(t, result.TVRageID)
	assert.Equal(t, int32(7926), *result.TVRageID)
	assert.Equal(t, new("Q190050"), result.WikidataID)
	assert.Equal(t, new("FightClub"), result.FacebookID)
	assert.Equal(t, new("fightclub"), result.InstagramID)
	assert.Equal(t, new("fightclub"), result.TwitterID)
	assert.Equal(t, new("fightclub"), result.TikTokID)
	assert.Equal(t, new("UC123"), result.YouTubeID)
	assert.Equal(t, new("/en/fight_club"), result.FreebaseID)
	assert.Equal(t, new("/m/0bth"), result.FreebaseMID)
}

func TestMapExternalIDs_NilOptionalFields(t *testing.T) {
	input := &ExternalIDsResponse{
		ID:     550,
		IMDbID: new("tt0137523"),
	}
	result := mapExternalIDs(input, 550)
	assert.Equal(t, new("tt0137523"), result.IMDbID)
	assert.Nil(t, result.TVDbID)
	assert.Nil(t, result.TVRageID)
}

func TestMapExternalIDs_ZeroTVDbID(t *testing.T) {
	input := &ExternalIDsResponse{
		ID:       550,
		TVDbID:   new(0),
		TVRageID: new(0),
	}
	result := mapExternalIDs(input, 550)
	assert.Nil(t, result.TVDbID)
	assert.Nil(t, result.TVRageID)
}

func TestMapCollectionMetadata(t *testing.T) {
	input := &CollectionResponse{
		ID:           86311,
		Name:         "Fight Club Collection",
		Overview:     new("A collection of Fight Club movies."),
		PosterPath:   new("/col_poster.jpg"),
		BackdropPath: new("/col_backdrop.jpg"),
		Parts: []MovieSearchResponse{
			{ID: 550, Title: "Fight Club", ReleaseDate: "1999-10-15"},
			{ID: 551, Title: "Fight Club 2"},
		},
	}

	result := mapCollectionMetadata(input)

	assert.Equal(t, "86311", result.ProviderID)
	assert.Equal(t, metadata.ProviderTMDb, result.Provider)
	assert.Equal(t, "Fight Club Collection", result.Name)
	assert.Equal(t, new("A collection of Fight Club movies."), result.Overview)
	assert.Equal(t, new("/col_poster.jpg"), result.PosterPath)
	assert.Equal(t, new("/col_backdrop.jpg"), result.BackdropPath)

	require.Len(t, result.Parts, 2)
	assert.Equal(t, "550", result.Parts[0].ProviderID)
	assert.Equal(t, "Fight Club", result.Parts[0].Title)
	require.NotNil(t, result.Parts[0].ReleaseDate)
	assert.Equal(t, "551", result.Parts[1].ProviderID)
	assert.Nil(t, result.Parts[1].ReleaseDate)
}
