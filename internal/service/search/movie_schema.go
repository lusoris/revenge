package search

import (
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// MovieCollectionName is the name of the movies collection in Typesense.
const MovieCollectionName = "movies"

// MovieDocument represents a movie document in the search index.
type MovieDocument struct {
	ID               string   `json:"id"`               // UUID as string
	TMDbID           int32    `json:"tmdb_id"`          // TMDb ID for external reference
	IMDbID           string   `json:"imdb_id"`          // IMDb ID (optional)
	Title            string   `json:"title"`            // Main title (searchable)
	OriginalTitle    string   `json:"original_title"`   // Original title (searchable)
	Year             int32    `json:"year"`             // Release year (facet + filter)
	ReleaseDate      int64    `json:"release_date"`     // Unix timestamp for sorting
	Runtime          int32    `json:"runtime"`          // Runtime in minutes
	Overview         string   `json:"overview"`         // Plot overview (searchable)
	Tagline          string   `json:"tagline"`          // Movie tagline
	Status           string   `json:"status"`           // Release status (facet)
	OriginalLanguage string   `json:"original_language"` // Original language (facet)
	PosterPath       string   `json:"poster_path"`       // Poster image path
	BackdropPath     string   `json:"backdrop_path"`     // Backdrop image path
	VoteAverage      float64  `json:"vote_average"`      // Rating (sortable)
	VoteCount        int32    `json:"vote_count"`        // Vote count
	Popularity       float64  `json:"popularity"`        // TMDb popularity (sortable)
	Genres           []string `json:"genres"`            // Genre names (facet + filter)
	GenreIDs         []int32  `json:"genre_ids"`         // Genre IDs
	Cast             []string `json:"cast"`              // Cast member names (searchable)
	Directors        []string `json:"directors"`         // Director names (searchable + facet)
	HasFile          bool     `json:"has_file"`          // Whether movie has a file (facet)
	Resolution       string   `json:"resolution"`        // Video resolution (facet)
	QualityProfile   string   `json:"quality_profile"`   // Quality profile (facet)
	LibraryAddedAt   int64    `json:"library_added_at"`  // When added to library (sortable)
	CreatedAt        int64    `json:"created_at"`        // Document creation time
	UpdatedAt        int64    `json:"updated_at"`        // Document update time
}

// ptr is a helper to create a pointer to a value.
func ptr[T any](v T) *T {
	return &v
}

// MovieCollectionSchema returns the Typesense schema for the movies collection.
func MovieCollectionSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name:                MovieCollectionName,
		EnableNestedFields:  ptr(false),
		TokenSeparators:     &[]string{"-", "_", "'", "'"},
		SymbolsToIndex:      &[]string{"&"},
		DefaultSortingField: ptr("popularity"),
		Fields: []api.Field{
			// ID field (primary key)
			{Name: "id", Type: "string"},

			// External IDs
			{Name: "tmdb_id", Type: "int32", Facet: ptr(false), Index: ptr(true)},
			{Name: "imdb_id", Type: "string", Facet: ptr(false), Index: ptr(true), Optional: ptr(true)},

			// Title fields (searchable with different weights)
			{Name: "title", Type: "string", Facet: ptr(false), Index: ptr(true), Infix: ptr(true)},
			{Name: "original_title", Type: "string", Facet: ptr(false), Index: ptr(true), Infix: ptr(true), Optional: ptr(true)},

			// Year and dates
			{Name: "year", Type: "int32", Facet: ptr(true), Index: ptr(true), Optional: ptr(true)},
			{Name: "release_date", Type: "int64", Facet: ptr(false), Index: ptr(true), Sort: ptr(true), Optional: ptr(true)},

			// Movie details
			{Name: "runtime", Type: "int32", Facet: ptr(false), Index: ptr(true), Optional: ptr(true)},
			{Name: "overview", Type: "string", Facet: ptr(false), Index: ptr(true), Optional: ptr(true)},
			{Name: "tagline", Type: "string", Facet: ptr(false), Index: ptr(false), Optional: ptr(true)},
			{Name: "status", Type: "string", Facet: ptr(true), Index: ptr(true), Optional: ptr(true)},
			{Name: "original_language", Type: "string", Facet: ptr(true), Index: ptr(true), Optional: ptr(true)},

			// Images
			{Name: "poster_path", Type: "string", Facet: ptr(false), Index: ptr(false), Optional: ptr(true)},
			{Name: "backdrop_path", Type: "string", Facet: ptr(false), Index: ptr(false), Optional: ptr(true)},

			// Ratings and popularity (sortable)
			{Name: "vote_average", Type: "float", Facet: ptr(false), Index: ptr(true), Sort: ptr(true), Optional: ptr(true)},
			{Name: "vote_count", Type: "int32", Facet: ptr(false), Index: ptr(true), Optional: ptr(true)},
			{Name: "popularity", Type: "float", Facet: ptr(false), Index: ptr(true), Sort: ptr(true), Optional: ptr(true)},

			// Genres (array, facetable, filterable)
			{Name: "genres", Type: "string[]", Facet: ptr(true), Index: ptr(true), Optional: ptr(true)},
			{Name: "genre_ids", Type: "int32[]", Facet: ptr(false), Index: ptr(true), Optional: ptr(true)},

			// Credits (searchable)
			{Name: "cast", Type: "string[]", Facet: ptr(false), Index: ptr(true), Optional: ptr(true)},
			{Name: "directors", Type: "string[]", Facet: ptr(true), Index: ptr(true), Optional: ptr(true)},

			// Library status
			{Name: "has_file", Type: "bool", Facet: ptr(true), Index: ptr(true), Optional: ptr(true)},
			{Name: "resolution", Type: "string", Facet: ptr(true), Index: ptr(true), Optional: ptr(true)},
			{Name: "quality_profile", Type: "string", Facet: ptr(true), Index: ptr(true), Optional: ptr(true)},

			// Timestamps
			{Name: "library_added_at", Type: "int64", Facet: ptr(false), Index: ptr(true), Sort: ptr(true), Optional: ptr(true)},
			{Name: "created_at", Type: "int64", Facet: ptr(false), Index: ptr(false), Optional: ptr(true)},
			{Name: "updated_at", Type: "int64", Facet: ptr(false), Index: ptr(false), Optional: ptr(true)},
		},
	}
}
