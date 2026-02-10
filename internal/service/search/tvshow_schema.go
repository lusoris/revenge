package search

import (
	"github.com/typesense/typesense-go/v2/typesense/api"

	"github.com/lusoris/revenge/internal/util/ptr"
)

// TVShowCollectionName is the name of the TV shows collection in Typesense.
const TVShowCollectionName = "tvshows"

// TVShowDocument represents a TV show document in the search index.
type TVShowDocument struct {
	ID               string   `json:"id"`                // UUID as string
	TMDbID           int32    `json:"tmdb_id"`           // TMDb ID for external reference
	TVDbID           int32    `json:"tvdb_id"`           // TVDb ID
	IMDbID           string   `json:"imdb_id"`           // IMDb ID (optional)
	Title            string   `json:"title"`             // Main title (searchable)
	OriginalTitle    string   `json:"original_title"`    // Original title (searchable)
	Year             int32    `json:"year"`              // First air year (facet + filter)
	FirstAirDate     int64    `json:"first_air_date"`    // Unix timestamp for sorting
	Overview         string   `json:"overview"`          // Plot overview (searchable)
	Status           string   `json:"status"`            // Series status (facet)
	Type             string   `json:"type"`              // Scripted, Reality, etc. (facet)
	OriginalLanguage string   `json:"original_language"` // Original language (facet)
	PosterPath       string   `json:"poster_path"`       // Poster image path
	BackdropPath     string   `json:"backdrop_path"`     // Backdrop image path
	VoteAverage      float64  `json:"vote_average"`      // Rating (sortable)
	VoteCount        int32    `json:"vote_count"`        // Vote count
	Popularity       float64  `json:"popularity"`        // TMDb popularity (sortable)
	Genres           []string `json:"genres"`            // Genre names (facet + filter)
	GenreIDs         []int32  `json:"genre_ids"`         // Genre IDs
	Cast             []string `json:"cast"`              // Cast member names (searchable)
	Networks         []string `json:"networks"`          // Network names (searchable + facet)
	TotalSeasons     int32    `json:"total_seasons"`     // Number of seasons
	TotalEpisodes    int32    `json:"total_episodes"`    // Number of episodes
	HasFile          bool     `json:"has_file"`          // Whether series has any episode files (facet)
	CreatedAt        int64    `json:"created_at"`        // Document creation time
	UpdatedAt        int64    `json:"updated_at"`        // Document update time
}

// TVShowCollectionSchema returns the Typesense schema for the TV shows collection.
func TVShowCollectionSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name:                TVShowCollectionName,
		EnableNestedFields:  ptr.To(false),
		TokenSeparators:     &[]string{"-", "_"},
		SymbolsToIndex:      &[]string{"&"},
		DefaultSortingField: ptr.To("popularity"),
		Fields: []api.Field{
			// ID field (primary key)
			{Name: "id", Type: "string"},

			// External IDs
			{Name: "tmdb_id", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true)},
			{Name: "tvdb_id", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "imdb_id", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},

			// Title fields (searchable with infix for partial matching)
			{Name: "title", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Infix: ptr.To(true)},
			{Name: "original_title", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Infix: ptr.To(true), Optional: ptr.To(true)},

			// Year and dates
			{Name: "year", Type: "int32", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "first_air_date", Type: "int64", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true), Optional: ptr.To(true)},

			// Show details
			{Name: "overview", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "status", Type: "string", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "type", Type: "string", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},
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

			// Credits and networks (searchable)
			{Name: "cast", Type: "string[]", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "networks", Type: "string[]", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},

			// Counts
			{Name: "total_seasons", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},
			{Name: "total_episodes", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true), Optional: ptr.To(true)},

			// Library status
			{Name: "has_file", Type: "bool", Facet: ptr.To(true), Index: ptr.To(true), Optional: ptr.To(true)},

			// Timestamps
			{Name: "created_at", Type: "int64", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},
			{Name: "updated_at", Type: "int64", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},
		},
	}
}
