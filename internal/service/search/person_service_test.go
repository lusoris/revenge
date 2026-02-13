package search

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersonCollectionSchema(t *testing.T) {
	schema := PersonCollectionSchema()

	assert.Equal(t, PersonCollectionName, schema.Name)
	assert.NotNil(t, schema.DefaultSortingField)
	assert.Equal(t, "total_credits", *schema.DefaultSortingField)

	fieldNames := make(map[string]bool)
	for _, f := range schema.Fields {
		fieldNames[f.Name] = true
	}

	requiredFields := []string{
		"id", "tmdb_id", "name", "profile_path", "known_for",
		"characters", "departments", "movie_count", "tvshow_count", "total_credits",
	}

	for _, name := range requiredFields {
		assert.True(t, fieldNames[name], "schema should have field: %s", name)
	}
}

func TestPersonCollectionSchemaFacets(t *testing.T) {
	schema := PersonCollectionSchema()

	facetFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Facet != nil && *f.Facet {
			facetFields[f.Name] = true
		}
	}

	assert.True(t, facetFields["departments"], "departments should be facetable")
}

func TestPersonCollectionSchemaSortable(t *testing.T) {
	schema := PersonCollectionSchema()

	sortableFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Sort != nil && *f.Sort {
			sortableFields[f.Name] = true
		}
	}

	expectedSortable := []string{"movie_count", "tvshow_count", "total_credits"}
	for _, name := range expectedSortable {
		assert.True(t, sortableFields[name], "field %s should be sortable", name)
	}
}

func TestPersonSchemaFieldTypes(t *testing.T) {
	schema := PersonCollectionSchema()

	fieldTypes := make(map[string]string)
	for _, f := range schema.Fields {
		fieldTypes[f.Name] = f.Type
	}

	assert.Equal(t, "string", fieldTypes["id"])
	assert.Equal(t, "int32", fieldTypes["tmdb_id"])
	assert.Equal(t, "string", fieldTypes["name"])
	assert.Equal(t, "string", fieldTypes["profile_path"])
	assert.Equal(t, "string[]", fieldTypes["known_for"])
	assert.Equal(t, "string[]", fieldTypes["characters"])
	assert.Equal(t, "string[]", fieldTypes["departments"])
	assert.Equal(t, "int32", fieldTypes["movie_count"])
	assert.Equal(t, "int32", fieldTypes["tvshow_count"])
	assert.Equal(t, "int32", fieldTypes["total_credits"])
}

func TestPersonSchemaInfixSearch(t *testing.T) {
	schema := PersonCollectionSchema()

	for _, f := range schema.Fields {
		if f.Name == "name" {
			assert.NotNil(t, f.Infix, "name should have infix set")
			assert.True(t, *f.Infix, "name should have infix enabled")
		}
	}
}

func TestPersonSchemaTokenSeparators(t *testing.T) {
	schema := PersonCollectionSchema()

	assert.NotNil(t, schema.TokenSeparators)
	assert.Contains(t, *schema.TokenSeparators, "-")
	assert.Contains(t, *schema.TokenSeparators, "_")
}

func TestPersonSchemaSymbolsToIndex(t *testing.T) {
	schema := PersonCollectionSchema()

	assert.NotNil(t, schema.SymbolsToIndex)
	assert.Contains(t, *schema.SymbolsToIndex, "&")
}

func TestDefaultPersonSearchParams(t *testing.T) {
	params := DefaultPersonSearchParams()

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PerPage)
	assert.Equal(t, "total_credits:desc", params.SortBy)
	assert.True(t, params.IncludeHighlights)
	assert.Contains(t, params.FacetBy, "departments")
}

func TestPersonToDocument(t *testing.T) {
	agg := PersonAggregate{
		TMDbPersonID: 287,
		Name:         "Brad Pitt",
		ProfilePath:  "/profile.jpg",
		KnownFor:     []string{"Fight Club", "Se7en", "Troy", "Inglourious Basterds"},
		Characters:   []string{"Tyler Durden", "David Mills", "Achilles"},
		Departments:  []string{"Acting", "Production"},
		MovieCount:   45,
		TVShowCount:  3,
	}

	doc := personToDocument(agg)

	assert.Equal(t, "287", doc.ID)
	assert.Equal(t, int32(287), doc.TMDbID)
	assert.Equal(t, "Brad Pitt", doc.Name)
	assert.Equal(t, "/profile.jpg", doc.ProfilePath)
	assert.Equal(t, 4, len(doc.KnownFor))
	assert.Equal(t, "Fight Club", doc.KnownFor[0])
	assert.Equal(t, 3, len(doc.Characters))
	assert.Equal(t, "Tyler Durden", doc.Characters[0])
	assert.Equal(t, 2, len(doc.Departments))
	assert.Equal(t, int32(45), doc.MovieCount)
	assert.Equal(t, int32(3), doc.TVShowCount)
	assert.Equal(t, int32(48), doc.TotalCredits) // 45 + 3
}

func TestPersonToDocumentMinimal(t *testing.T) {
	agg := PersonAggregate{
		TMDbPersonID: 100,
		Name:         "Unknown Actor",
	}

	doc := personToDocument(agg)

	assert.Equal(t, "100", doc.ID)
	assert.Equal(t, int32(100), doc.TMDbID)
	assert.Equal(t, "Unknown Actor", doc.Name)
	assert.Empty(t, doc.ProfilePath)
	assert.Nil(t, doc.KnownFor)
	assert.Nil(t, doc.Characters)
	assert.Nil(t, doc.Departments)
	assert.Equal(t, int32(0), doc.MovieCount)
	assert.Equal(t, int32(0), doc.TVShowCount)
	assert.Equal(t, int32(0), doc.TotalCredits)
}

func TestPersonToDocumentDeduplication(t *testing.T) {
	agg := PersonAggregate{
		TMDbPersonID: 500,
		Name:         "Actor",
		KnownFor:     []string{"Film A", "Film B", "Film A", "Film C", "Film B"},
		Characters:   []string{"Hero", "Villain", "Hero"},
		Departments:  []string{"Acting", "Acting", "Directing"},
	}

	doc := personToDocument(agg)

	assert.Equal(t, 3, len(doc.KnownFor), "known_for should be deduplicated")
	assert.Equal(t, 2, len(doc.Characters), "characters should be deduplicated")
	assert.Equal(t, 2, len(doc.Departments), "departments should be deduplicated")
}

func TestPersonToDocumentCapsKnownForAt20(t *testing.T) {
	titles := make([]string, 30)
	for i := range titles {
		titles[i] = fmt.Sprintf("Movie %d", i+1)
	}

	agg := PersonAggregate{
		TMDbPersonID: 1,
		Name:         "Prolific Actor",
		KnownFor:     titles,
	}

	doc := personToDocument(agg)

	assert.Equal(t, 20, len(doc.KnownFor), "known_for should be capped at 20")
}

func TestParsePersonDocument(t *testing.T) {
	data := map[string]any{
		"id":            "287",
		"tmdb_id":       float64(287),
		"name":          "Brad Pitt",
		"profile_path":  "/profile.jpg",
		"known_for":     []any{"Fight Club", "Se7en"},
		"characters":    []any{"Tyler Durden", "David Mills"},
		"departments":   []any{"Acting", "Production"},
		"movie_count":   float64(45),
		"tvshow_count":  float64(3),
		"total_credits": float64(48),
	}

	doc := parsePersonDocument(data)

	assert.Equal(t, "287", doc.ID)
	assert.Equal(t, int32(287), doc.TMDbID)
	assert.Equal(t, "Brad Pitt", doc.Name)
	assert.Equal(t, "/profile.jpg", doc.ProfilePath)
	assert.Equal(t, []string{"Fight Club", "Se7en"}, doc.KnownFor)
	assert.Equal(t, []string{"Tyler Durden", "David Mills"}, doc.Characters)
	assert.Equal(t, []string{"Acting", "Production"}, doc.Departments)
	assert.Equal(t, int32(45), doc.MovieCount)
	assert.Equal(t, int32(3), doc.TVShowCount)
	assert.Equal(t, int32(48), doc.TotalCredits)
}

func TestParsePersonDocumentEmpty(t *testing.T) {
	data := map[string]any{}
	doc := parsePersonDocument(data)

	assert.Empty(t, doc.ID)
	assert.Empty(t, doc.Name)
	assert.Equal(t, int32(0), doc.TMDbID)
	assert.Nil(t, doc.KnownFor)
}

func TestParsePersonDocumentPartialData(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]any
		verify func(t *testing.T, doc PersonDocument)
	}{
		{
			name: "only id and name",
			data: map[string]any{
				"id":   "123",
				"name": "Jane Doe",
			},
			verify: func(t *testing.T, doc PersonDocument) {
				assert.Equal(t, "123", doc.ID)
				assert.Equal(t, "Jane Doe", doc.Name)
				assert.Equal(t, int32(0), doc.MovieCount)
				assert.Nil(t, doc.KnownFor)
			},
		},
		{
			name: "with nil values",
			data: map[string]any{
				"id":           "456",
				"name":         "John Doe",
				"known_for":    nil,
				"characters":   nil,
				"profile_path": nil,
			},
			verify: func(t *testing.T, doc PersonDocument) {
				assert.Equal(t, "456", doc.ID)
				assert.Nil(t, doc.KnownFor)
				assert.Nil(t, doc.Characters)
				assert.Empty(t, doc.ProfilePath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parsePersonDocument(tt.data)
			tt.verify(t, doc)
		})
	}
}

func TestPersonSearchServiceIsEnabled(t *testing.T) {
	s := &PersonSearchService{client: nil}
	assert.False(t, s.IsEnabled())
}

func TestPersonSearchParamsValidation(t *testing.T) {
	params := DefaultPersonSearchParams()

	assert.Greater(t, params.Page, 0)
	assert.Greater(t, params.PerPage, 0)
	assert.LessOrEqual(t, params.PerPage, 100)
	assert.NotEmpty(t, params.SortBy)
	assert.NotEmpty(t, params.FacetBy)
}

func TestPersonAggregate(t *testing.T) {
	agg := PersonAggregate{
		TMDbPersonID: 1000,
		Name:         "Test Person",
		ProfilePath:  "/test.jpg",
		KnownFor:     []string{"Movie A"},
		Characters:   []string{"Character X"},
		Departments:  []string{"Acting"},
		MovieCount:   5,
		TVShowCount:  2,
	}

	assert.Equal(t, int32(1000), agg.TMDbPersonID)
	assert.Equal(t, "Test Person", agg.Name)
	assert.Equal(t, "/test.jpg", agg.ProfilePath)
	assert.Equal(t, []string{"Movie A"}, agg.KnownFor)
	assert.Equal(t, []string{"Character X"}, agg.Characters)
	assert.Equal(t, []string{"Acting"}, agg.Departments)
	assert.Equal(t, int32(5), agg.MovieCount)
	assert.Equal(t, int32(2), agg.TVShowCount)
}
