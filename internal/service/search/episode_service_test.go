package search

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

func TestEpisodeCollectionSchema(t *testing.T) {
	schema := EpisodeCollectionSchema()

	assert.Equal(t, EpisodeCollectionName, schema.Name)
	assert.NotNil(t, schema.DefaultSortingField)
	assert.Equal(t, "air_date", *schema.DefaultSortingField)

	fieldNames := make(map[string]bool)
	for _, f := range schema.Fields {
		fieldNames[f.Name] = true
	}

	requiredFields := []string{
		"id", "series_id", "season_id", "season_number", "episode_number",
		"title", "overview", "air_date", "runtime", "has_file",
		"series_title", "series_poster_path",
		"vote_average", "vote_count", "still_path",
	}

	for _, name := range requiredFields {
		assert.True(t, fieldNames[name], "schema should have field: %s", name)
	}
}

func TestEpisodeCollectionSchemaFacets(t *testing.T) {
	schema := EpisodeCollectionSchema()

	facetFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Facet != nil && *f.Facet {
			facetFields[f.Name] = true
		}
	}

	expectedFacets := []string{"season_number", "has_file"}
	for _, name := range expectedFacets {
		assert.True(t, facetFields[name], "field %s should be facetable", name)
	}
}

func TestEpisodeCollectionSchemaSortable(t *testing.T) {
	schema := EpisodeCollectionSchema()

	sortableFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Sort != nil && *f.Sort {
			sortableFields[f.Name] = true
		}
	}

	expectedSortable := []string{"air_date", "vote_average"}
	for _, name := range expectedSortable {
		assert.True(t, sortableFields[name], "field %s should be sortable", name)
	}
}

func TestEpisodeSchemaFieldTypes(t *testing.T) {
	schema := EpisodeCollectionSchema()

	fieldTypes := make(map[string]string)
	for _, f := range schema.Fields {
		fieldTypes[f.Name] = f.Type
	}

	assert.Equal(t, "string", fieldTypes["id"])
	assert.Equal(t, "string", fieldTypes["series_id"])
	assert.Equal(t, "string", fieldTypes["season_id"])
	assert.Equal(t, "int32", fieldTypes["tmdb_id"])
	assert.Equal(t, "int32", fieldTypes["tvdb_id"])
	assert.Equal(t, "string", fieldTypes["imdb_id"])
	assert.Equal(t, "int32", fieldTypes["season_number"])
	assert.Equal(t, "int32", fieldTypes["episode_number"])
	assert.Equal(t, "string", fieldTypes["title"])
	assert.Equal(t, "string", fieldTypes["overview"])
	assert.Equal(t, "int64", fieldTypes["air_date"])
	assert.Equal(t, "int32", fieldTypes["runtime"])
	assert.Equal(t, "float", fieldTypes["vote_average"])
	assert.Equal(t, "int32", fieldTypes["vote_count"])
	assert.Equal(t, "string", fieldTypes["still_path"])
	assert.Equal(t, "bool", fieldTypes["has_file"])
	assert.Equal(t, "string", fieldTypes["series_title"])
	assert.Equal(t, "string", fieldTypes["series_poster_path"])
	assert.Equal(t, "int64", fieldTypes["created_at"])
	assert.Equal(t, "int64", fieldTypes["updated_at"])
}

func TestEpisodeSchemaInfixSearch(t *testing.T) {
	schema := EpisodeCollectionSchema()

	for _, f := range schema.Fields {
		if f.Name == "title" || f.Name == "series_title" {
			assert.NotNil(t, f.Infix, "field %s should have infix set", f.Name)
			assert.True(t, *f.Infix, "field %s should have infix enabled", f.Name)
		}
	}
}

func TestEpisodeSchemaTokenSeparators(t *testing.T) {
	schema := EpisodeCollectionSchema()

	assert.NotNil(t, schema.TokenSeparators)
	assert.Contains(t, *schema.TokenSeparators, "-")
	assert.Contains(t, *schema.TokenSeparators, "_")
}

func TestEpisodeSchemaSymbolsToIndex(t *testing.T) {
	schema := EpisodeCollectionSchema()

	assert.NotNil(t, schema.SymbolsToIndex)
	assert.Contains(t, *schema.SymbolsToIndex, "&")
}

func TestDefaultEpisodeSearchParams(t *testing.T) {
	params := DefaultEpisodeSearchParams()

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PerPage)
	assert.Equal(t, "air_date:desc", params.SortBy)
	assert.True(t, params.IncludeHighlights)
	assert.Contains(t, params.FacetBy, "season_number")
	assert.Contains(t, params.FacetBy, "has_file")
}

func TestEpisodeToDocument(t *testing.T) {
	s := &EpisodeSearchService{}

	episodeID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	now := time.Now()
	tmdbID := int32(62085)
	tvdbID := int32(349232)
	imdbID := "tt2301451"
	overview := "Walter White and Jesse Pinkman attempt their first cook."
	airDate := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	runtime := int32(58)
	voteAverage, _ := decimal.NewFromFloat64(9.2)
	voteCount := int32(5000)
	stillPath := "/still.jpg"

	ep := &tvshow.Episode{
		ID:            episodeID,
		SeriesID:      seriesID,
		SeasonID:      seasonID,
		TMDbID:        &tmdbID,
		TVDbID:        &tvdbID,
		IMDbID:        &imdbID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
		Overview:      &overview,
		AirDate:       &airDate,
		Runtime:       &runtime,
		VoteAverage:   &voteAverage,
		VoteCount:     &voteCount,
		StillPath:     &stillPath,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	doc := s.episodeToDocument(ep, "Breaking Bad", "/poster.jpg", true)

	assert.Equal(t, episodeID.String(), doc.ID)
	assert.Equal(t, seriesID.String(), doc.SeriesID)
	assert.Equal(t, seasonID.String(), doc.SeasonID)
	assert.Equal(t, tmdbID, doc.TMDbID)
	assert.Equal(t, tvdbID, doc.TVDbID)
	assert.Equal(t, imdbID, doc.IMDbID)
	assert.Equal(t, int32(1), doc.SeasonNumber)
	assert.Equal(t, int32(1), doc.EpisodeNumber)
	assert.Equal(t, "Pilot", doc.Title)
	assert.Equal(t, overview, doc.Overview)
	assert.Equal(t, airDate.Unix(), doc.AirDate)
	assert.Equal(t, runtime, doc.Runtime)
	assert.InDelta(t, 9.2, doc.VoteAverage, 0.01)
	assert.Equal(t, voteCount, doc.VoteCount)
	assert.Equal(t, stillPath, doc.StillPath)
	assert.True(t, doc.HasFile)
	assert.Equal(t, "Breaking Bad", doc.SeriesTitle)
	assert.Equal(t, "/poster.jpg", doc.SeriesPosterPath)
}

func TestEpisodeToDocumentMinimal(t *testing.T) {
	s := &EpisodeSearchService{}

	episodeID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	now := time.Now()

	ep := &tvshow.Episode{
		ID:            episodeID,
		SeriesID:      seriesID,
		SeasonID:      seasonID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Unknown Episode",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	doc := s.episodeToDocument(ep, "Unknown Show", "", false)

	assert.Equal(t, episodeID.String(), doc.ID)
	assert.Equal(t, seriesID.String(), doc.SeriesID)
	assert.Equal(t, "Unknown Episode", doc.Title)
	assert.Equal(t, int32(0), doc.TMDbID)
	assert.Empty(t, doc.IMDbID)
	assert.Empty(t, doc.Overview)
	assert.Empty(t, doc.StillPath)
	assert.Equal(t, int64(0), doc.AirDate)
	assert.Equal(t, int32(0), doc.Runtime)
	assert.False(t, doc.HasFile)
	assert.Equal(t, "Unknown Show", doc.SeriesTitle)
	assert.Empty(t, doc.SeriesPosterPath)
}

func TestEpisodeToDocumentZeroValues(t *testing.T) {
	s := &EpisodeSearchService{}

	episodeID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	now := time.Now()
	zeroDecimal, _ := decimal.NewFromFloat64(0.0)
	zeroInt := int32(0)

	ep := &tvshow.Episode{
		ID:            episodeID,
		SeriesID:      seriesID,
		SeasonID:      seasonID,
		SeasonNumber:  0,
		EpisodeNumber: 1,
		Title:         "Special Episode",
		VoteAverage:   &zeroDecimal,
		VoteCount:     &zeroInt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	doc := s.episodeToDocument(ep, "Show", "", false)

	assert.Equal(t, float64(0), doc.VoteAverage)
	assert.Equal(t, int32(0), doc.VoteCount)
	assert.Equal(t, int32(0), doc.SeasonNumber)
}

func TestParseEpisodeDocument(t *testing.T) {
	data := map[string]interface{}{
		"id":                "550e8400-e29b-41d4-a716-446655440000",
		"series_id":         "660e8400-e29b-41d4-a716-446655440000",
		"season_id":         "770e8400-e29b-41d4-a716-446655440000",
		"tmdb_id":           float64(62085),
		"tvdb_id":           float64(349232),
		"imdb_id":           "tt2301451",
		"season_number":     float64(1),
		"episode_number":    float64(1),
		"title":             "Pilot",
		"overview":          "A chemistry teacher...",
		"air_date":          float64(1200787200),
		"runtime":           float64(58),
		"vote_average":      float64(9.2),
		"vote_count":        float64(5000),
		"still_path":        "/still.jpg",
		"has_file":          true,
		"series_title":      "Breaking Bad",
		"series_poster_path": "/poster.jpg",
		"created_at":        float64(1700000000),
		"updated_at":        float64(1700000000),
	}

	doc := parseEpisodeDocument(data)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", doc.ID)
	assert.Equal(t, "660e8400-e29b-41d4-a716-446655440000", doc.SeriesID)
	assert.Equal(t, "770e8400-e29b-41d4-a716-446655440000", doc.SeasonID)
	assert.Equal(t, int32(62085), doc.TMDbID)
	assert.Equal(t, int32(349232), doc.TVDbID)
	assert.Equal(t, "tt2301451", doc.IMDbID)
	assert.Equal(t, int32(1), doc.SeasonNumber)
	assert.Equal(t, int32(1), doc.EpisodeNumber)
	assert.Equal(t, "Pilot", doc.Title)
	assert.Equal(t, "A chemistry teacher...", doc.Overview)
	assert.Equal(t, int64(1200787200), doc.AirDate)
	assert.Equal(t, int32(58), doc.Runtime)
	assert.InDelta(t, 9.2, doc.VoteAverage, 0.01)
	assert.Equal(t, int32(5000), doc.VoteCount)
	assert.Equal(t, "/still.jpg", doc.StillPath)
	assert.True(t, doc.HasFile)
	assert.Equal(t, "Breaking Bad", doc.SeriesTitle)
	assert.Equal(t, "/poster.jpg", doc.SeriesPosterPath)
	assert.Equal(t, int64(1700000000), doc.CreatedAt)
	assert.Equal(t, int64(1700000000), doc.UpdatedAt)
}

func TestParseEpisodeDocumentEmpty(t *testing.T) {
	data := map[string]interface{}{}
	doc := parseEpisodeDocument(data)

	assert.Empty(t, doc.ID)
	assert.Empty(t, doc.SeriesID)
	assert.Empty(t, doc.Title)
	assert.Equal(t, int32(0), doc.TMDbID)
	assert.False(t, doc.HasFile)
}

func TestParseEpisodeDocumentPartialData(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]interface{}
		verify func(t *testing.T, doc EpisodeDocument)
	}{
		{
			name: "only id and title",
			data: map[string]interface{}{
				"id":    "test-id",
				"title": "Test Episode",
			},
			verify: func(t *testing.T, doc EpisodeDocument) {
				assert.Equal(t, "test-id", doc.ID)
				assert.Equal(t, "Test Episode", doc.Title)
				assert.Equal(t, int32(0), doc.SeasonNumber)
				assert.Empty(t, doc.Overview)
			},
		},
		{
			name: "with nil values",
			data: map[string]interface{}{
				"id":           "test-id",
				"title":        "Test Episode",
				"overview":     nil,
				"series_title": nil,
			},
			verify: func(t *testing.T, doc EpisodeDocument) {
				assert.Equal(t, "test-id", doc.ID)
				assert.Empty(t, doc.Overview)
				assert.Empty(t, doc.SeriesTitle)
			},
		},
		{
			name: "with boolean fields",
			data: map[string]interface{}{
				"id":       "test-id",
				"has_file": true,
			},
			verify: func(t *testing.T, doc EpisodeDocument) {
				assert.True(t, doc.HasFile)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseEpisodeDocument(tt.data)
			tt.verify(t, doc)
		})
	}
}

func TestEpisodeSearchServiceIsEnabled(t *testing.T) {
	s := &EpisodeSearchService{client: nil}
	assert.False(t, s.IsEnabled())
}

func TestEpisodeSearchParamsValidation(t *testing.T) {
	params := DefaultEpisodeSearchParams()

	assert.Greater(t, params.Page, 0)
	assert.Greater(t, params.PerPage, 0)
	assert.LessOrEqual(t, params.PerPage, 100)
	assert.NotEmpty(t, params.SortBy)
	assert.NotEmpty(t, params.FacetBy)
}

func TestDerefHits(t *testing.T) {
	assert.Nil(t, derefHits(nil))

	// Non-nil returns the value
	hits := []api.SearchResultHit{}
	assert.NotNil(t, derefHits(&hits))
	assert.Empty(t, derefHits(&hits))
}

func TestEpisodeWithContext(t *testing.T) {
	episodeID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	now := time.Now()

	ep := &tvshow.Episode{
		ID:            episodeID,
		SeriesID:      seriesID,
		SeasonID:      seasonID,
		SeasonNumber:  2,
		EpisodeNumber: 5,
		Title:         "Dead Freight",
		CreatedAt:     now,
	}

	ctx := EpisodeWithContext{
		Episode:          ep,
		SeriesTitle:      "Breaking Bad",
		SeriesPosterPath: "/poster.jpg",
		HasFile:          true,
	}

	assert.Equal(t, ep, ctx.Episode)
	assert.Equal(t, "Breaking Bad", ctx.SeriesTitle)
	assert.Equal(t, "/poster.jpg", ctx.SeriesPosterPath)
	assert.True(t, ctx.HasFile)
}
