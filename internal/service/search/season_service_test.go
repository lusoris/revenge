package search

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/stretchr/testify/assert"
)

func TestSeasonCollectionSchema(t *testing.T) {
	schema := SeasonCollectionSchema()

	assert.Equal(t, SeasonCollectionName, schema.Name)
	assert.NotNil(t, schema.DefaultSortingField)
	assert.Equal(t, "air_date", *schema.DefaultSortingField)

	fieldNames := make(map[string]bool)
	for _, f := range schema.Fields {
		fieldNames[f.Name] = true
	}

	requiredFields := []string{
		"id", "series_id", "season_number", "name", "overview",
		"air_date", "episode_count", "vote_average", "poster_path",
		"series_title", "series_poster_path",
	}

	for _, name := range requiredFields {
		assert.True(t, fieldNames[name], "schema should have field: %s", name)
	}
}

func TestSeasonCollectionSchemaFacets(t *testing.T) {
	schema := SeasonCollectionSchema()

	facetFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Facet != nil && *f.Facet {
			facetFields[f.Name] = true
		}
	}

	assert.True(t, facetFields["season_number"], "season_number should be facetable")
}

func TestSeasonCollectionSchemaSortable(t *testing.T) {
	schema := SeasonCollectionSchema()

	sortableFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Sort != nil && *f.Sort {
			sortableFields[f.Name] = true
		}
	}

	expectedSortable := []string{"air_date", "vote_average", "episode_count"}
	for _, name := range expectedSortable {
		assert.True(t, sortableFields[name], "field %s should be sortable", name)
	}
}

func TestSeasonSchemaFieldTypes(t *testing.T) {
	schema := SeasonCollectionSchema()

	fieldTypes := make(map[string]string)
	for _, f := range schema.Fields {
		fieldTypes[f.Name] = f.Type
	}

	assert.Equal(t, "string", fieldTypes["id"])
	assert.Equal(t, "string", fieldTypes["series_id"])
	assert.Equal(t, "int32", fieldTypes["tmdb_id"])
	assert.Equal(t, "int32", fieldTypes["season_number"])
	assert.Equal(t, "string", fieldTypes["name"])
	assert.Equal(t, "string", fieldTypes["overview"])
	assert.Equal(t, "int64", fieldTypes["air_date"])
	assert.Equal(t, "int32", fieldTypes["episode_count"])
	assert.Equal(t, "float", fieldTypes["vote_average"])
	assert.Equal(t, "string", fieldTypes["poster_path"])
	assert.Equal(t, "string", fieldTypes["series_title"])
	assert.Equal(t, "string", fieldTypes["series_poster_path"])
	assert.Equal(t, "int64", fieldTypes["created_at"])
	assert.Equal(t, "int64", fieldTypes["updated_at"])
}

func TestSeasonSchemaInfixSearch(t *testing.T) {
	schema := SeasonCollectionSchema()

	for _, f := range schema.Fields {
		if f.Name == "name" || f.Name == "series_title" {
			assert.NotNil(t, f.Infix, "field %s should have infix set", f.Name)
			assert.True(t, *f.Infix, "field %s should have infix enabled", f.Name)
		}
	}
}

func TestSeasonSchemaTokenSeparators(t *testing.T) {
	schema := SeasonCollectionSchema()

	assert.NotNil(t, schema.TokenSeparators)
	assert.Contains(t, *schema.TokenSeparators, "-")
	assert.Contains(t, *schema.TokenSeparators, "_")
}

func TestSeasonSchemaSymbolsToIndex(t *testing.T) {
	schema := SeasonCollectionSchema()

	assert.NotNil(t, schema.SymbolsToIndex)
	assert.Contains(t, *schema.SymbolsToIndex, "&")
}

func TestDefaultSeasonSearchParams(t *testing.T) {
	params := DefaultSeasonSearchParams()

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PerPage)
	assert.Equal(t, "air_date:desc", params.SortBy)
	assert.True(t, params.IncludeHighlights)
	assert.Contains(t, params.FacetBy, "season_number")
}

func TestSeasonToDocument(t *testing.T) {
	s := &SeasonSearchService{}

	seasonID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()
	tmdbID := int32(3572)
	overview := "The first season of Breaking Bad."
	airDate := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	voteAverage, _ := decimal.NewFromFloat64(8.9)
	posterPath := "/season1_poster.jpg"

	season := &tvshow.Season{
		ID:           seasonID,
		SeriesID:     seriesID,
		TMDbID:       &tmdbID,
		SeasonNumber: 1,
		Name:         "Season 1",
		Overview:     &overview,
		PosterPath:   &posterPath,
		EpisodeCount: 7,
		AirDate:      &airDate,
		VoteAverage:  &voteAverage,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	doc := s.seasonToDocument(season, "Breaking Bad", "/series_poster.jpg")

	assert.Equal(t, seasonID.String(), doc.ID)
	assert.Equal(t, seriesID.String(), doc.SeriesID)
	assert.Equal(t, tmdbID, doc.TMDbID)
	assert.Equal(t, int32(1), doc.SeasonNumber)
	assert.Equal(t, "Season 1", doc.Name)
	assert.Equal(t, overview, doc.Overview)
	assert.Equal(t, airDate.Unix(), doc.AirDate)
	assert.Equal(t, int32(7), doc.EpisodeCount)
	assert.InDelta(t, 8.9, doc.VoteAverage, 0.01)
	assert.Equal(t, posterPath, doc.PosterPath)
	assert.Equal(t, "Breaking Bad", doc.SeriesTitle)
	assert.Equal(t, "/series_poster.jpg", doc.SeriesPosterPath)
}

func TestSeasonToDocumentMinimal(t *testing.T) {
	s := &SeasonSearchService{}

	seasonID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()

	season := &tvshow.Season{
		ID:           seasonID,
		SeriesID:     seriesID,
		SeasonNumber: 0,
		Name:         "Specials",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	doc := s.seasonToDocument(season, "Unknown Show", "")

	assert.Equal(t, seasonID.String(), doc.ID)
	assert.Equal(t, seriesID.String(), doc.SeriesID)
	assert.Equal(t, "Specials", doc.Name)
	assert.Equal(t, int32(0), doc.TMDbID)
	assert.Empty(t, doc.Overview)
	assert.Empty(t, doc.PosterPath)
	assert.Equal(t, int64(0), doc.AirDate)
	assert.Equal(t, int32(0), doc.EpisodeCount)
	assert.Equal(t, "Unknown Show", doc.SeriesTitle)
	assert.Empty(t, doc.SeriesPosterPath)
}

func TestSeasonToDocumentZeroValues(t *testing.T) {
	s := &SeasonSearchService{}

	seasonID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()
	zeroDecimal, _ := decimal.NewFromFloat64(0.0)

	season := &tvshow.Season{
		ID:           seasonID,
		SeriesID:     seriesID,
		SeasonNumber: 0,
		Name:         "Specials",
		VoteAverage:  &zeroDecimal,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	doc := s.seasonToDocument(season, "Show", "")

	assert.Equal(t, float64(0), doc.VoteAverage)
	assert.Equal(t, int32(0), doc.SeasonNumber)
}

func TestParseSeasonDocument(t *testing.T) {
	data := map[string]interface{}{
		"id":                 "550e8400-e29b-41d4-a716-446655440000",
		"series_id":          "660e8400-e29b-41d4-a716-446655440000",
		"tmdb_id":            float64(3572),
		"season_number":      float64(1),
		"name":               "Season 1",
		"overview":           "The first season...",
		"air_date":           float64(1200787200),
		"episode_count":      float64(7),
		"vote_average":       float64(8.9),
		"poster_path":        "/season1.jpg",
		"series_title":       "Breaking Bad",
		"series_poster_path": "/poster.jpg",
		"created_at":         float64(1700000000),
		"updated_at":         float64(1700000000),
	}

	doc := parseSeasonDocument(data)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", doc.ID)
	assert.Equal(t, "660e8400-e29b-41d4-a716-446655440000", doc.SeriesID)
	assert.Equal(t, int32(3572), doc.TMDbID)
	assert.Equal(t, int32(1), doc.SeasonNumber)
	assert.Equal(t, "Season 1", doc.Name)
	assert.Equal(t, "The first season...", doc.Overview)
	assert.Equal(t, int64(1200787200), doc.AirDate)
	assert.Equal(t, int32(7), doc.EpisodeCount)
	assert.InDelta(t, 8.9, doc.VoteAverage, 0.01)
	assert.Equal(t, "/season1.jpg", doc.PosterPath)
	assert.Equal(t, "Breaking Bad", doc.SeriesTitle)
	assert.Equal(t, "/poster.jpg", doc.SeriesPosterPath)
	assert.Equal(t, int64(1700000000), doc.CreatedAt)
	assert.Equal(t, int64(1700000000), doc.UpdatedAt)
}

func TestParseSeasonDocumentEmpty(t *testing.T) {
	data := map[string]interface{}{}
	doc := parseSeasonDocument(data)

	assert.Empty(t, doc.ID)
	assert.Empty(t, doc.SeriesID)
	assert.Empty(t, doc.Name)
	assert.Equal(t, int32(0), doc.TMDbID)
}

func TestParseSeasonDocumentPartialData(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]interface{}
		verify func(t *testing.T, doc SeasonDocument)
	}{
		{
			name: "only id and name",
			data: map[string]interface{}{
				"id":   "test-id",
				"name": "Season 1",
			},
			verify: func(t *testing.T, doc SeasonDocument) {
				assert.Equal(t, "test-id", doc.ID)
				assert.Equal(t, "Season 1", doc.Name)
				assert.Equal(t, int32(0), doc.SeasonNumber)
				assert.Empty(t, doc.Overview)
			},
		},
		{
			name: "with nil values",
			data: map[string]interface{}{
				"id":           "test-id",
				"name":         "Season 1",
				"overview":     nil,
				"series_title": nil,
			},
			verify: func(t *testing.T, doc SeasonDocument) {
				assert.Equal(t, "test-id", doc.ID)
				assert.Empty(t, doc.Overview)
				assert.Empty(t, doc.SeriesTitle)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseSeasonDocument(tt.data)
			tt.verify(t, doc)
		})
	}
}

func TestSeasonSearchServiceIsEnabled(t *testing.T) {
	s := &SeasonSearchService{client: nil}
	assert.False(t, s.IsEnabled())
}

func TestSeasonSearchParamsValidation(t *testing.T) {
	params := DefaultSeasonSearchParams()

	assert.Greater(t, params.Page, 0)
	assert.Greater(t, params.PerPage, 0)
	assert.LessOrEqual(t, params.PerPage, 100)
	assert.NotEmpty(t, params.SortBy)
	assert.NotEmpty(t, params.FacetBy)
}

func TestSeasonWithContext(t *testing.T) {
	seasonID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	now := time.Now()

	season := &tvshow.Season{
		ID:           seasonID,
		SeriesID:     seriesID,
		SeasonNumber: 3,
		Name:         "Season 3",
		CreatedAt:    now,
	}

	ctx := SeasonWithContext{
		Season:           season,
		SeriesTitle:      "Breaking Bad",
		SeriesPosterPath: "/poster.jpg",
	}

	assert.Equal(t, season, ctx.Season)
	assert.Equal(t, "Breaking Bad", ctx.SeriesTitle)
	assert.Equal(t, "/poster.jpg", ctx.SeriesPosterPath)
}
