package search

import (
	"github.com/typesense/typesense-go/v2/typesense/api"

	"github.com/lusoris/revenge/internal/util/ptr"
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

// MovieCollectionSchema returns the Typesense schema for the movies collection.
func MovieCollectionSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name:                MovieCollectionName,
		EnableNestedFields:  ptr.To(false),
		TokenSeparators:     &[]string{"-", "_"},
		SymbolsToIndex:      &[]string{"&"},
		DefaultSortingField: ptr.To("popularity"),
		Fields: []api.Field{
			// ID field (primary key)
			{Name: "id", Type: "string"},

			// External IDs
			{Name: "tmdb_id", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true)},
			{Name: "imdb_id", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},

			// Title fields (searchable with different weights)
			{Name: "title", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Infix: ptr.To(true)},
			{Name: "original_title", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Infix: ptr.To(true), Optional: ptr.To(true)},

			// Year and dates
			{Name: "year", Type: "int32", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "release_date", Type: "int64", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true), Optional: ptr.To(true)},

			// Movie details
			{Name: "runtime", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "overview", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "tagline", Type: "string", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},
			{Name: "status", Type: "string", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "original_language", Type: "string", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},

			// Images
			{Name: "poster_path", Type: "string", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},
			{Name: "backdrop_path", Type: "string", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},

			// Ratings and popularity (sortable)
			{Name: "vote_average", Type: "float", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true), Optional: ptr.To(true)},
			{Name: "vote_count", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "popularity", Type: "float", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true)},

			// Genres (array, facetable, filterable)
			{Name: "genres", Type: "string[]", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "genre_ids", Type: "int32[]", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},

			// Credits (searchable)
			{Name: "cast", Type: "string[]", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "directors", Type: "string[]", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},

			// Library status
			{Name: "has_file", Type: "bool", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "resolution", Type: "string", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "quality_profile", Type: "string", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},

			// Timestamps
			{Name: "library_added_at", Type: "int64", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true), Optional: ptr.To(true)},
			{Name: "created_at", Type: "int64", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},
			{Name: "updated_at", Type: "int64", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},
		},
	}
}
