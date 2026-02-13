package search

import (
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// MovieCollectionName is the name of the movies collection in Typesense.
const MovieCollectionName = "movies"

// MovieDocument represents a movie document in the search index.
type MovieDocument struct {
	ID               string   `json:"id"`                // UUID as string
	TMDbID           int32    `json:"tmdb_id"`           // TMDb ID for external reference
	IMDbID           string   `json:"imdb_id"`           // IMDb ID (optional)
	Title            string   `json:"title"`             // Main title (searchable)
	OriginalTitle    string   `json:"original_title"`    // Original title (searchable)
	Year             int32    `json:"year"`              // Release year (facet + filter)
	ReleaseDate      int64    `json:"release_date"`      // Unix timestamp for sorting
	Runtime          int32    `json:"runtime"`           // Runtime in minutes
	Overview         string   `json:"overview"`          // Plot overview (searchable)
	Tagline          string   `json:"tagline"`           // Movie tagline
	Status           string   `json:"status"`            // Release status (facet)
	OriginalLanguage string   `json:"original_language"` // Original language (facet)
	PosterPath       string   `json:"poster_path"`       // Poster image path
	BackdropPath     string   `json:"backdrop_path"`     // Backdrop image path
	VoteAverage      float64  `json:"vote_average"`      // Rating (sortable)
	VoteCount        int32    `json:"vote_count"`        // Vote count
	Popularity       float64  `json:"popularity"`        // TMDb popularity (sortable)
	Genres           []string `json:"genres"`            // Genre names (facet + filter)
	GenreSlugs       []string `json:"genre_slugs"`       // Genre slugs (facet + filter)
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
		EnableNestedFields:  new(false),
		TokenSeparators:     &[]string{"-", "_"},
		SymbolsToIndex:      &[]string{"&"},
		DefaultSortingField: new("popularity"),
		Fields: []api.Field{
			// ID field (primary key)
			{Name: "id", Type: "string"},

			// External IDs
			{Name: "tmdb_id", Type: "int32", Facet: new(false), Index: new(true)},
			{Name: "imdb_id", Type: "string", Facet: new(false), Index: new(true), Optional: new(true)},

			// Title fields (searchable with different weights)
			{Name: "title", Type: "string", Facet: new(false), Index: new(true), Infix: new(true)},
			{Name: "original_title", Type: "string", Facet: new(false), Index: new(true), Infix: new(true), Optional: new(true)},

			// Year and dates
			{Name: "year", Type: "int32", Facet: new(true), Index: new(true), Optional: new(true)},
			{Name: "release_date", Type: "int64", Facet: new(false), Index: new(true), Sort: new(true), Optional: new(true)},

			// Movie details
			{Name: "runtime", Type: "int32", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "overview", Type: "string", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "tagline", Type: "string", Facet: new(false), Index: new(false), Optional: new(true)},
			{Name: "status", Type: "string", Facet: new(true), Index: new(true), Optional: new(true)},
			{Name: "original_language", Type: "string", Facet: new(true), Index: new(true), Optional: new(true)},

			// Images
			{Name: "poster_path", Type: "string", Facet: new(false), Index: new(false), Optional: new(true)},
			{Name: "backdrop_path", Type: "string", Facet: new(false), Index: new(false), Optional: new(true)},

			// Ratings and popularity (sortable)
			{Name: "vote_average", Type: "float", Facet: new(false), Index: new(true), Sort: new(true), Optional: new(true)},
			{Name: "vote_count", Type: "int32", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "popularity", Type: "float", Facet: new(false), Index: new(true), Sort: new(true)},

			// Genres (array, facetable, filterable)
			{Name: "genres", Type: "string[]", Facet: new(true), Index: new(true), Optional: new(true)},
			{Name: "genre_slugs", Type: "string[]", Facet: new(true), Index: new(true), Optional: new(true)},

			// Credits (searchable)
			{Name: "cast", Type: "string[]", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "directors", Type: "string[]", Facet: new(true), Index: new(true), Optional: new(true)},

			// Library status
			{Name: "has_file", Type: "bool", Facet: new(true), Index: new(true), Optional: new(true)},
			{Name: "resolution", Type: "string", Facet: new(true), Index: new(true), Optional: new(true)},
			{Name: "quality_profile", Type: "string", Facet: new(true), Index: new(true), Optional: new(true)},

			// Timestamps
			{Name: "library_added_at", Type: "int64", Facet: new(false), Index: new(true), Sort: new(true), Optional: new(true)},
			{Name: "created_at", Type: "int64", Facet: new(false), Index: new(false), Optional: new(true)},
			{Name: "updated_at", Type: "int64", Facet: new(false), Index: new(false), Optional: new(true)},
		},
	}
}
