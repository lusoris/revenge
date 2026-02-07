package search

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/stretchr/testify/assert"
)

func TestTVShowCollectionSchema(t *testing.T) {
	schema := TVShowCollectionSchema()

	assert.Equal(t, TVShowCollectionName, schema.Name)
	assert.NotNil(t, schema.DefaultSortingField)
	assert.Equal(t, "popularity", *schema.DefaultSortingField)

	fieldNames := make(map[string]bool)
	for _, f := range schema.Fields {
		fieldNames[f.Name] = true
	}

	requiredFields := []string{
		"id", "tmdb_id", "title", "year", "first_air_date",
		"overview", "genres", "cast", "networks", "has_file",
		"vote_average", "popularity", "status", "type",
		"total_seasons", "total_episodes",
	}

	for _, name := range requiredFields {
		assert.True(t, fieldNames[name], "schema should have field: %s", name)
	}
}

func TestTVShowCollectionSchemaFacets(t *testing.T) {
	schema := TVShowCollectionSchema()

	facetFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Facet != nil && *f.Facet {
			facetFields[f.Name] = true
		}
	}

	expectedFacets := []string{"year", "status", "type", "original_language", "genres", "networks", "has_file"}
	for _, name := range expectedFacets {
		assert.True(t, facetFields[name], "field %s should be facetable", name)
	}
}

func TestTVShowCollectionSchemaSortable(t *testing.T) {
	schema := TVShowCollectionSchema()

	sortableFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Sort != nil && *f.Sort {
			sortableFields[f.Name] = true
		}
	}

	expectedSortable := []string{"first_air_date", "vote_average", "popularity", "total_episodes"}
	for _, name := range expectedSortable {
		assert.True(t, sortableFields[name], "field %s should be sortable", name)
	}
}

func TestTVShowSchemaFieldTypes(t *testing.T) {
	schema := TVShowCollectionSchema()

	fieldTypes := make(map[string]string)
	for _, f := range schema.Fields {
		fieldTypes[f.Name] = f.Type
	}

	assert.Equal(t, "string", fieldTypes["id"])
	assert.Equal(t, "int32", fieldTypes["tmdb_id"])
	assert.Equal(t, "int32", fieldTypes["tvdb_id"])
	assert.Equal(t, "string", fieldTypes["imdb_id"])
	assert.Equal(t, "string", fieldTypes["title"])
	assert.Equal(t, "int32", fieldTypes["year"])
	assert.Equal(t, "int64", fieldTypes["first_air_date"])
	assert.Equal(t, "float", fieldTypes["vote_average"])
	assert.Equal(t, "float", fieldTypes["popularity"])
	assert.Equal(t, "string[]", fieldTypes["genres"])
	assert.Equal(t, "string[]", fieldTypes["cast"])
	assert.Equal(t, "string[]", fieldTypes["networks"])
	assert.Equal(t, "bool", fieldTypes["has_file"])
	assert.Equal(t, "int32", fieldTypes["total_seasons"])
	assert.Equal(t, "int32", fieldTypes["total_episodes"])
}

func TestTVShowSchemaInfixSearch(t *testing.T) {
	schema := TVShowCollectionSchema()

	for _, f := range schema.Fields {
		if f.Name == "title" || f.Name == "original_title" {
			assert.NotNil(t, f.Infix, "field %s should have infix set", f.Name)
			assert.True(t, *f.Infix, "field %s should have infix enabled", f.Name)
		}
	}
}

func TestTVShowSchemaTokenSeparators(t *testing.T) {
	schema := TVShowCollectionSchema()

	assert.NotNil(t, schema.TokenSeparators)
	assert.Contains(t, *schema.TokenSeparators, "-")
	assert.Contains(t, *schema.TokenSeparators, "_")
}

func TestTVShowSchemaSymbolsToIndex(t *testing.T) {
	schema := TVShowCollectionSchema()

	assert.NotNil(t, schema.SymbolsToIndex)
	assert.Contains(t, *schema.SymbolsToIndex, "&")
}

func TestDefaultTVShowSearchParams(t *testing.T) {
	params := DefaultTVShowSearchParams()

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PerPage)
	assert.Equal(t, "popularity:desc", params.SortBy)
	assert.True(t, params.IncludeHighlights)
	assert.Contains(t, params.FacetBy, "genres")
	assert.Contains(t, params.FacetBy, "year")
	assert.Contains(t, params.FacetBy, "status")
	assert.Contains(t, params.FacetBy, "type")
	assert.Contains(t, params.FacetBy, "networks")
	assert.Contains(t, params.FacetBy, "has_file")
}

func TestSeriesToDocument(t *testing.T) {
	s := &TVShowSearchService{}

	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()
	tmdbID := int32(1396)
	tvdbID := int32(81189)
	imdbID := "tt0903747"
	originalTitle := "Breaking Bad"
	firstAirDate := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	overview := "A chemistry teacher diagnosed with inoperable lung cancer."
	status := "Ended"
	showType := "Scripted"
	posterPath := "/poster.jpg"
	backdropPath := "/backdrop.jpg"
	voteAverage, _ := decimal.NewFromFloat64(9.5)
	voteCount := int32(15000)
	popularity, _ := decimal.NewFromFloat64(120.3)

	series := &tvshow.Series{
		ID:               seriesID,
		TMDbID:           &tmdbID,
		TVDbID:           &tvdbID,
		IMDbID:           &imdbID,
		Title:            "Breaking Bad",
		OriginalTitle:    &originalTitle,
		OriginalLanguage: "en",
		FirstAirDate:     &firstAirDate,
		Overview:         &overview,
		Status:           &status,
		Type:             &showType,
		PosterPath:       &posterPath,
		BackdropPath:     &backdropPath,
		VoteAverage:      &voteAverage,
		VoteCount:        &voteCount,
		Popularity:       &popularity,
		TotalSeasons:     5,
		TotalEpisodes:    62,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	genres := []tvshow.SeriesGenre{
		{TMDbGenreID: 18, Name: "Drama"},
		{TMDbGenreID: 80, Name: "Crime"},
	}

	credits := []tvshow.SeriesCredit{
		{CreditType: "cast", Name: "Bryan Cranston"},
		{CreditType: "cast", Name: "Aaron Paul"},
		{CreditType: "crew", Name: "Vince Gilligan", Job: ptr("Creator")},
	}

	networks := []tvshow.Network{
		{Name: "AMC"},
	}

	doc := s.seriesToDocument(series, genres, credits, networks, true)

	assert.Equal(t, seriesID.String(), doc.ID)
	assert.Equal(t, tmdbID, doc.TMDbID)
	assert.Equal(t, tvdbID, doc.TVDbID)
	assert.Equal(t, imdbID, doc.IMDbID)
	assert.Equal(t, "Breaking Bad", doc.Title)
	assert.Equal(t, "Breaking Bad", doc.OriginalTitle)
	assert.Equal(t, int32(2008), doc.Year)
	assert.Equal(t, firstAirDate.Unix(), doc.FirstAirDate)
	assert.Equal(t, overview, doc.Overview)
	assert.Equal(t, status, doc.Status)
	assert.Equal(t, showType, doc.Type)
	assert.Equal(t, "en", doc.OriginalLanguage)
	assert.Equal(t, posterPath, doc.PosterPath)
	assert.Equal(t, backdropPath, doc.BackdropPath)
	assert.InDelta(t, 9.5, doc.VoteAverage, 0.01)
	assert.Equal(t, voteCount, doc.VoteCount)
	assert.InDelta(t, 120.3, doc.Popularity, 0.01)
	assert.Equal(t, int32(5), doc.TotalSeasons)
	assert.Equal(t, int32(62), doc.TotalEpisodes)

	// Verify genres
	assert.Equal(t, []string{"Drama", "Crime"}, doc.Genres)
	assert.Equal(t, []int32{18, 80}, doc.GenreIDs)

	// Verify credits (only cast)
	assert.Equal(t, []string{"Bryan Cranston", "Aaron Paul"}, doc.Cast)

	// Verify networks
	assert.Equal(t, []string{"AMC"}, doc.Networks)

	// Verify file info
	assert.True(t, doc.HasFile)
}

func TestSeriesToDocumentMinimal(t *testing.T) {
	s := &TVShowSearchService{}

	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()

	series := &tvshow.Series{
		ID:        seriesID,
		Title:     "Unknown Show",
		CreatedAt: now,
		UpdatedAt: now,
	}

	doc := s.seriesToDocument(series, nil, nil, nil, false)

	assert.Equal(t, seriesID.String(), doc.ID)
	assert.Equal(t, "Unknown Show", doc.Title)
	assert.Equal(t, int32(0), doc.TMDbID)
	assert.Empty(t, doc.IMDbID)
	assert.Empty(t, doc.Overview)
	assert.Empty(t, doc.Genres)
	assert.Empty(t, doc.Cast)
	assert.Empty(t, doc.Networks)
	assert.False(t, doc.HasFile)
	assert.Equal(t, int32(0), doc.TotalSeasons)
	assert.Equal(t, int32(0), doc.TotalEpisodes)
}

func TestSeriesToDocumentCastLimit(t *testing.T) {
	s := &TVShowSearchService{}

	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()

	series := &tvshow.Series{
		ID:        seriesID,
		Title:     "Show With Many Cast",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create 30 cast members
	credits := make([]tvshow.SeriesCredit, 30)
	for i := range 30 {
		credits[i] = tvshow.SeriesCredit{
			CreditType: "cast",
			Name:       "Actor " + string(rune('A'+i)),
		}
	}

	doc := s.seriesToDocument(series, nil, credits, nil, false)

	// Should be limited to 20
	assert.Len(t, doc.Cast, 20)
}

func TestSeriesToDocumentWithEmptyCredits(t *testing.T) {
	s := &TVShowSearchService{}

	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()

	series := &tvshow.Series{
		ID:        seriesID,
		Title:     "Empty Credits Show",
		CreatedAt: now,
		UpdatedAt: now,
	}

	genres := []tvshow.SeriesGenre{}
	credits := []tvshow.SeriesCredit{}
	networks := []tvshow.Network{}

	doc := s.seriesToDocument(series, genres, credits, networks, false)

	assert.Equal(t, seriesID.String(), doc.ID)
	assert.Empty(t, doc.Genres)
	assert.Empty(t, doc.GenreIDs)
	assert.Empty(t, doc.Cast)
	assert.Empty(t, doc.Networks)
}

func TestSeriesToDocumentWithCrewOnly(t *testing.T) {
	s := &TVShowSearchService{}

	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()

	series := &tvshow.Series{
		ID:        seriesID,
		Title:     "Crew Only Show",
		CreatedAt: now,
		UpdatedAt: now,
	}

	credits := []tvshow.SeriesCredit{
		{CreditType: "crew", Name: "Vince Gilligan", Job: ptr("Creator")},
		{CreditType: "crew", Name: "Peter Gould", Job: ptr("Writer")},
	}

	doc := s.seriesToDocument(series, nil, credits, nil, false)

	// No cast members extracted
	assert.Empty(t, doc.Cast)
}

func TestSeriesToDocumentWithZeroValues(t *testing.T) {
	s := &TVShowSearchService{}

	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()
	zeroDecimal, _ := decimal.NewFromFloat64(0.0)
	zeroInt := int32(0)

	series := &tvshow.Series{
		ID:          seriesID,
		Title:       "Zero Values Show",
		VoteAverage: &zeroDecimal,
		VoteCount:   &zeroInt,
		Popularity:  &zeroDecimal,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	doc := s.seriesToDocument(series, nil, nil, nil, false)

	assert.Equal(t, float64(0), doc.VoteAverage)
	assert.Equal(t, int32(0), doc.VoteCount)
	assert.Equal(t, float64(0), doc.Popularity)
}

func TestParseTVShowDocument(t *testing.T) {
	data := map[string]interface{}{
		"id":                "550e8400-e29b-41d4-a716-446655440000",
		"tmdb_id":           float64(1396),
		"tvdb_id":           float64(81189),
		"imdb_id":           "tt0903747",
		"title":             "Breaking Bad",
		"original_title":    "Breaking Bad",
		"year":              float64(2008),
		"first_air_date":    float64(1200787200),
		"overview":          "A chemistry teacher...",
		"status":            "Ended",
		"type":              "Scripted",
		"original_language": "en",
		"poster_path":       "/poster.jpg",
		"backdrop_path":     "/backdrop.jpg",
		"vote_average":      float64(9.5),
		"vote_count":        float64(15000),
		"popularity":        float64(120.3),
		"has_file":          true,
		"total_seasons":     float64(5),
		"total_episodes":    float64(62),
		"created_at":        float64(1700000000),
		"updated_at":        float64(1700000000),
		"genres":            []interface{}{"Drama", "Crime"},
		"cast":              []interface{}{"Bryan Cranston", "Aaron Paul"},
		"networks":          []interface{}{"AMC"},
		"genre_ids":         []interface{}{float64(18), float64(80)},
	}

	doc := parseTVShowDocument(data)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", doc.ID)
	assert.Equal(t, int32(1396), doc.TMDbID)
	assert.Equal(t, int32(81189), doc.TVDbID)
	assert.Equal(t, "tt0903747", doc.IMDbID)
	assert.Equal(t, "Breaking Bad", doc.Title)
	assert.Equal(t, "Breaking Bad", doc.OriginalTitle)
	assert.Equal(t, int32(2008), doc.Year)
	assert.Equal(t, int64(1200787200), doc.FirstAirDate)
	assert.Equal(t, "A chemistry teacher...", doc.Overview)
	assert.Equal(t, "Ended", doc.Status)
	assert.Equal(t, "Scripted", doc.Type)
	assert.Equal(t, "en", doc.OriginalLanguage)
	assert.Equal(t, "/poster.jpg", doc.PosterPath)
	assert.Equal(t, "/backdrop.jpg", doc.BackdropPath)
	assert.InDelta(t, 9.5, doc.VoteAverage, 0.01)
	assert.Equal(t, int32(15000), doc.VoteCount)
	assert.InDelta(t, 120.3, doc.Popularity, 0.01)
	assert.True(t, doc.HasFile)
	assert.Equal(t, int32(5), doc.TotalSeasons)
	assert.Equal(t, int32(62), doc.TotalEpisodes)
	assert.Equal(t, int64(1700000000), doc.CreatedAt)
	assert.Equal(t, []string{"Drama", "Crime"}, doc.Genres)
	assert.Equal(t, []string{"Bryan Cranston", "Aaron Paul"}, doc.Cast)
	assert.Equal(t, []string{"AMC"}, doc.Networks)
	assert.Equal(t, []int32{18, 80}, doc.GenreIDs)
}

func TestParseTVShowDocumentEmpty(t *testing.T) {
	data := map[string]interface{}{}
	doc := parseTVShowDocument(data)

	assert.Empty(t, doc.ID)
	assert.Equal(t, int32(0), doc.TMDbID)
	assert.Empty(t, doc.Title)
	assert.Nil(t, doc.Genres)
}

func TestParseTVShowDocumentPartialData(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]interface{}
		verify func(t *testing.T, doc TVShowDocument)
	}{
		{
			name: "only id and title",
			data: map[string]interface{}{
				"id":    "test-id",
				"title": "Test Show",
			},
			verify: func(t *testing.T, doc TVShowDocument) {
				assert.Equal(t, "test-id", doc.ID)
				assert.Equal(t, "Test Show", doc.Title)
				assert.Equal(t, int32(0), doc.Year)
				assert.Empty(t, doc.Overview)
			},
		},
		{
			name: "with nil values",
			data: map[string]interface{}{
				"id":       "test-id",
				"title":    "Test Show",
				"year":     nil,
				"genres":   nil,
				"networks": nil,
			},
			verify: func(t *testing.T, doc TVShowDocument) {
				assert.Equal(t, "test-id", doc.ID)
				assert.Equal(t, "Test Show", doc.Title)
				assert.Equal(t, int32(0), doc.Year)
				assert.Nil(t, doc.Genres)
				assert.Nil(t, doc.Networks)
			},
		},
		{
			name: "with boolean fields",
			data: map[string]interface{}{
				"id":       "test-id",
				"has_file": true,
			},
			verify: func(t *testing.T, doc TVShowDocument) {
				assert.True(t, doc.HasFile)
			},
		},
		{
			name: "with empty arrays",
			data: map[string]interface{}{
				"id":       "test-id",
				"genres":   []interface{}{},
				"networks": []interface{}{},
			},
			verify: func(t *testing.T, doc TVShowDocument) {
				assert.Empty(t, doc.Genres)
				assert.Empty(t, doc.Networks)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseTVShowDocument(tt.data)
			tt.verify(t, doc)
		})
	}
}

func TestTVShowSearchServiceIsEnabled(t *testing.T) {
	s := &TVShowSearchService{client: nil}
	assert.False(t, s.IsEnabled())
}

func TestTVShowSearchParamsValidation(t *testing.T) {
	params := DefaultTVShowSearchParams()

	assert.Greater(t, params.Page, 0)
	assert.Greater(t, params.PerPage, 0)
	assert.LessOrEqual(t, params.PerPage, 100)
	assert.NotEmpty(t, params.SortBy)
	assert.NotEmpty(t, params.FacetBy)
}
