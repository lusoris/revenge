package search

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

func TestMovieCollectionSchema(t *testing.T) {
	schema := MovieCollectionSchema()

	assert.Equal(t, MovieCollectionName, schema.Name)
	assert.NotNil(t, schema.DefaultSortingField)
	assert.Equal(t, "popularity", *schema.DefaultSortingField)

	// Check that required fields exist
	fieldNames := make(map[string]bool)
	for _, f := range schema.Fields {
		fieldNames[f.Name] = true
	}

	requiredFields := []string{
		"id", "tmdb_id", "title", "year", "release_date",
		"overview", "genres", "cast", "directors", "has_file",
		"vote_average", "popularity",
	}

	for _, name := range requiredFields {
		assert.True(t, fieldNames[name], "schema should have field: %s", name)
	}
}

func TestMovieCollectionSchemaFacets(t *testing.T) {
	schema := MovieCollectionSchema()

	facetFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Facet != nil && *f.Facet {
			facetFields[f.Name] = true
		}
	}

	// These fields should be facetable
	expectedFacets := []string{"year", "status", "original_language", "genres", "directors", "has_file", "resolution", "quality_profile"}
	for _, name := range expectedFacets {
		assert.True(t, facetFields[name], "field %s should be facetable", name)
	}
}

func TestMovieCollectionSchemaSortable(t *testing.T) {
	schema := MovieCollectionSchema()

	sortableFields := make(map[string]bool)
	for _, f := range schema.Fields {
		if f.Sort != nil && *f.Sort {
			sortableFields[f.Name] = true
		}
	}

	// These fields should be sortable
	expectedSortable := []string{"release_date", "vote_average", "popularity", "library_added_at"}
	for _, name := range expectedSortable {
		assert.True(t, sortableFields[name], "field %s should be sortable", name)
	}
}

func TestDefaultSearchParams(t *testing.T) {
	params := DefaultSearchParams()

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PerPage)
	assert.Equal(t, "popularity:desc", params.SortBy)
	assert.True(t, params.IncludeHighlights)
	assert.Contains(t, params.FacetBy, "genres")
	assert.Contains(t, params.FacetBy, "year")
}

func TestMovieToDocument(t *testing.T) {
	s := &MovieSearchService{}

	movieID := uuid.New()
	now := time.Now()
	tmdbID := int32(603)
	imdbID := "tt0133093"
	year := int32(1999)
	runtime := int32(136)
	overview := "A computer hacker learns about the true nature of reality."
	tagline := "Welcome to the Real World"
	status := "Released"
	originalLanguage := "en"
	posterPath := "/poster.jpg"
	backdropPath := "/backdrop.jpg"
	voteAverage := decimal.NewFromFloat(8.7)
	voteCount := int32(20000)
	popularity := decimal.NewFromFloat(80.5)
	releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
	originalTitle := "The Matrix"

	m := &movie.Movie{
		ID:               movieID,
		TMDbID:           &tmdbID,
		IMDbID:           &imdbID,
		Title:            "The Matrix",
		OriginalTitle:    &originalTitle,
		Year:             &year,
		ReleaseDate:      &releaseDate,
		Runtime:          &runtime,
		Overview:         &overview,
		Tagline:          &tagline,
		Status:           &status,
		OriginalLanguage: &originalLanguage,
		PosterPath:       &posterPath,
		BackdropPath:     &backdropPath,
		VoteAverage:      &voteAverage,
		VoteCount:        &voteCount,
		Popularity:       &popularity,
		LibraryAddedAt:   now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	genres := []movie.MovieGenre{
		{TMDbGenreID: 28, Name: "Action"},
		{TMDbGenreID: 878, Name: "Science Fiction"},
	}

	directorJob := "Director"
	credits := []movie.MovieCredit{
		{CreditType: "cast", Name: "Keanu Reeves"},
		{CreditType: "cast", Name: "Laurence Fishburne"},
		{CreditType: "crew", Name: "Lana Wachowski", Job: &directorJob},
		{CreditType: "crew", Name: "Lilly Wachowski", Job: &directorJob},
	}

	resolution := "2160p"
	qualityProfile := "Ultra-HD"
	file := &movie.MovieFile{
		ID:             uuid.New(),
		MovieID:        movieID,
		Resolution:     &resolution,
		QualityProfile: &qualityProfile,
	}

	doc := s.movieToDocument(m, genres, credits, file)

	// Verify basic fields
	assert.Equal(t, movieID.String(), doc.ID)
	assert.Equal(t, tmdbID, doc.TMDbID)
	assert.Equal(t, imdbID, doc.IMDbID)
	assert.Equal(t, "The Matrix", doc.Title)
	assert.Equal(t, "The Matrix", doc.OriginalTitle)
	assert.Equal(t, year, doc.Year)
	assert.Equal(t, releaseDate.Unix(), doc.ReleaseDate)
	assert.Equal(t, runtime, doc.Runtime)
	assert.Equal(t, overview, doc.Overview)
	assert.Equal(t, tagline, doc.Tagline)
	assert.Equal(t, status, doc.Status)
	assert.Equal(t, originalLanguage, doc.OriginalLanguage)
	assert.Equal(t, posterPath, doc.PosterPath)
	assert.Equal(t, backdropPath, doc.BackdropPath)
	assert.InDelta(t, 8.7, doc.VoteAverage, 0.01)
	assert.Equal(t, voteCount, doc.VoteCount)
	assert.InDelta(t, 80.5, doc.Popularity, 0.01)

	// Verify genres
	assert.Equal(t, []string{"Action", "Science Fiction"}, doc.Genres)
	assert.Equal(t, []int32{28, 878}, doc.GenreIDs)

	// Verify credits
	assert.Equal(t, []string{"Keanu Reeves", "Laurence Fishburne"}, doc.Cast)
	assert.Equal(t, []string{"Lana Wachowski", "Lilly Wachowski"}, doc.Directors)

	// Verify file info
	assert.True(t, doc.HasFile)
	assert.Equal(t, "2160p", doc.Resolution)
	assert.Equal(t, "Ultra-HD", doc.QualityProfile)
}

func TestMovieToDocumentMinimal(t *testing.T) {
	s := &MovieSearchService{}

	movieID := uuid.New()
	now := time.Now()

	m := &movie.Movie{
		ID:             movieID,
		Title:          "Unknown Movie",
		LibraryAddedAt: now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	doc := s.movieToDocument(m, nil, nil, nil)

	assert.Equal(t, movieID.String(), doc.ID)
	assert.Equal(t, "Unknown Movie", doc.Title)
	assert.Equal(t, int32(0), doc.TMDbID)
	assert.Empty(t, doc.IMDbID)
	assert.Empty(t, doc.Overview)
	assert.Empty(t, doc.Genres)
	assert.Empty(t, doc.Cast)
	assert.Empty(t, doc.Directors)
	assert.False(t, doc.HasFile)
}

func TestMovieToDocumentCastLimit(t *testing.T) {
	s := &MovieSearchService{}

	movieID := uuid.New()
	now := time.Now()

	m := &movie.Movie{
		ID:             movieID,
		Title:          "Movie With Many Cast",
		LibraryAddedAt: now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Create 30 cast members
	credits := make([]movie.MovieCredit, 30)
	for i := range 30 {
		credits[i] = movie.MovieCredit{
			CreditType: "cast",
			Name:       "Actor " + string(rune('A'+i)),
		}
	}

	doc := s.movieToDocument(m, nil, credits, nil)

	// Should be limited to 20
	assert.Len(t, doc.Cast, 20)
}

func TestParseMovieDocument(t *testing.T) {
	data := map[string]interface{}{
		"id":                "550e8400-e29b-41d4-a716-446655440000",
		"tmdb_id":           float64(603),
		"imdb_id":           "tt0133093",
		"title":             "The Matrix",
		"original_title":    "The Matrix",
		"year":              float64(1999),
		"release_date":      float64(922838400),
		"runtime":           float64(136),
		"overview":          "A computer hacker learns...",
		"tagline":           "Welcome to the Real World",
		"status":            "Released",
		"original_language": "en",
		"poster_path":       "/poster.jpg",
		"backdrop_path":     "/backdrop.jpg",
		"vote_average":      float64(8.7),
		"vote_count":        float64(20000),
		"popularity":        float64(80.5),
		"has_file":          true,
		"resolution":        "2160p",
		"quality_profile":   "Ultra-HD",
		"library_added_at":  float64(1700000000),
		"created_at":        float64(1700000000),
		"updated_at":        float64(1700000000),
		"genres":            []interface{}{"Action", "Science Fiction"},
		"cast":              []interface{}{"Keanu Reeves", "Laurence Fishburne"},
		"directors":         []interface{}{"Lana Wachowski", "Lilly Wachowski"},
		"genre_ids":         []interface{}{float64(28), float64(878)},
	}

	doc := parseMovieDocument(data)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", doc.ID)
	assert.Equal(t, int32(603), doc.TMDbID)
	assert.Equal(t, "tt0133093", doc.IMDbID)
	assert.Equal(t, "The Matrix", doc.Title)
	assert.Equal(t, "The Matrix", doc.OriginalTitle)
	assert.Equal(t, int32(1999), doc.Year)
	assert.Equal(t, int64(922838400), doc.ReleaseDate)
	assert.Equal(t, int32(136), doc.Runtime)
	assert.Equal(t, "A computer hacker learns...", doc.Overview)
	assert.Equal(t, "Welcome to the Real World", doc.Tagline)
	assert.Equal(t, "Released", doc.Status)
	assert.Equal(t, "en", doc.OriginalLanguage)
	assert.Equal(t, "/poster.jpg", doc.PosterPath)
	assert.Equal(t, "/backdrop.jpg", doc.BackdropPath)
	assert.InDelta(t, 8.7, doc.VoteAverage, 0.01)
	assert.Equal(t, int32(20000), doc.VoteCount)
	assert.InDelta(t, 80.5, doc.Popularity, 0.01)
	assert.True(t, doc.HasFile)
	assert.Equal(t, "2160p", doc.Resolution)
	assert.Equal(t, "Ultra-HD", doc.QualityProfile)
	assert.Equal(t, int64(1700000000), doc.LibraryAddedAt)
	assert.Equal(t, []string{"Action", "Science Fiction"}, doc.Genres)
	assert.Equal(t, []string{"Keanu Reeves", "Laurence Fishburne"}, doc.Cast)
	assert.Equal(t, []string{"Lana Wachowski", "Lilly Wachowski"}, doc.Directors)
	assert.Equal(t, []int32{28, 878}, doc.GenreIDs)
}

func TestParseMovieDocumentEmpty(t *testing.T) {
	data := map[string]interface{}{}
	doc := parseMovieDocument(data)

	assert.Empty(t, doc.ID)
	assert.Equal(t, int32(0), doc.TMDbID)
	assert.Empty(t, doc.Title)
	assert.Nil(t, doc.Genres)
}

func TestToStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []string
	}{
		{
			name:     "valid strings",
			input:    []interface{}{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty slice",
			input:    []interface{}{},
			expected: []string{},
		},
		{
			name:     "mixed types",
			input:    []interface{}{"a", 123, "b", nil},
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toStringSlice(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToInt32Slice(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []int32
	}{
		{
			name:     "valid floats",
			input:    []interface{}{float64(1), float64(2), float64(3)},
			expected: []int32{1, 2, 3},
		},
		{
			name:     "empty slice",
			input:    []interface{}{},
			expected: []int32{},
		},
		{
			name:     "mixed types",
			input:    []interface{}{float64(1), "two", float64(3)},
			expected: []int32{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toInt32Slice(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeref(t *testing.T) {
	v := 42
	assert.Equal(t, 42, deref(&v))
	assert.Equal(t, 0, deref(nil))
}

func TestPtr(t *testing.T) {
	s := "test"
	p := ptr(s)
	assert.NotNil(t, p)
	assert.Equal(t, "test", *p)

	i := 42
	pi := ptr(i)
	assert.NotNil(t, pi)
	assert.Equal(t, 42, *pi)
}

func TestMovieSearchServiceIsEnabled(t *testing.T) {
	// nil client
	s := &MovieSearchService{client: nil}
	assert.False(t, s.IsEnabled())
}

func TestSchemaFieldTypes(t *testing.T) {
	schema := MovieCollectionSchema()

	fieldTypes := make(map[string]string)
	for _, f := range schema.Fields {
		fieldTypes[f.Name] = f.Type
	}

	// Verify field types
	assert.Equal(t, "string", fieldTypes["id"])
	assert.Equal(t, "int32", fieldTypes["tmdb_id"])
	assert.Equal(t, "string", fieldTypes["imdb_id"])
	assert.Equal(t, "string", fieldTypes["title"])
	assert.Equal(t, "int32", fieldTypes["year"])
	assert.Equal(t, "int64", fieldTypes["release_date"])
	assert.Equal(t, "float", fieldTypes["vote_average"])
	assert.Equal(t, "float", fieldTypes["popularity"])
	assert.Equal(t, "string[]", fieldTypes["genres"])
	assert.Equal(t, "string[]", fieldTypes["cast"])
	assert.Equal(t, "string[]", fieldTypes["directors"])
	assert.Equal(t, "bool", fieldTypes["has_file"])
}

func TestSchemaTokenSeparators(t *testing.T) {
	schema := MovieCollectionSchema()

	assert.NotNil(t, schema.TokenSeparators)
	assert.Contains(t, *schema.TokenSeparators, "-")
	assert.Contains(t, *schema.TokenSeparators, "_")
}

func TestSchemaInfixSearch(t *testing.T) {
	schema := MovieCollectionSchema()

	// title and original_title should have infix enabled for better search
	for _, f := range schema.Fields {
		if f.Name == "title" || f.Name == "original_title" {
			assert.NotNil(t, f.Infix, "field %s should have infix set", f.Name)
			assert.True(t, *f.Infix, "field %s should have infix enabled", f.Name)
		}
	}
}

// MockTypesenseClient for testing without actual Typesense server
type MockTypesenseClient struct {
	Collections map[string]*api.CollectionResponse
	Documents   map[string]map[string]interface{}
}

func NewMockTypesenseClient() *MockTypesenseClient {
	return &MockTypesenseClient{
		Collections: make(map[string]*api.CollectionResponse),
		Documents:   make(map[string]map[string]interface{}),
	}
}
