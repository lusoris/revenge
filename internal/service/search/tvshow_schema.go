package search

import (
	"github.com/typesense/typesense-go/v2/typesense/api"
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
	GenreSlugs       []string `json:"genre_slugs"`       // Genre slugs (facet + filter)
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
		EnableNestedFields:  new(false),
		TokenSeparators:     &[]string{"-", "_"},
		SymbolsToIndex:      &[]string{"&"},
		DefaultSortingField: new("popularity"),
		Fields: []api.Field{
			// ID field (primary key)
			{Name: "id", Type: "string"},

			// External IDs
			{Name: "tmdb_id", Type: "int32", Facet: new(false), Index: new(true)},
			{Name: "tvdb_id", Type: "int32", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "imdb_id", Type: "string", Facet: new(false), Index: new(true), Optional: new(true)},

			// Title fields (searchable with infix for partial matching)
			{Name: "title", Type: "string", Facet: new(false), Index: new(true), Infix: new(true)},
			{Name: "original_title", Type: "string", Facet: new(false), Index: new(true), Infix: new(true), Optional: new(true)},

			// Year and dates
			{Name: "year", Type: "int32", Facet: new(true), Index: new(true), Optional: new(true)},
			{Name: "first_air_date", Type: "int64", Facet: new(false), Index: new(true), Sort: new(true), Optional: new(true)},

			// Show details
			{Name: "overview", Type: "string", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "status", Type: "string", Facet: new(true), Index: new(true), Optional: new(true)},
			{Name: "type", Type: "string", Facet: new(true), Index: new(true), Optional: new(true)},
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

			// Credits and networks (searchable)
			{Name: "cast", Type: "string[]", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "networks", Type: "string[]", Facet: new(true), Index: new(true), Optional: new(true)},

			// Counts
			{Name: "total_seasons", Type: "int32", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "total_episodes", Type: "int32", Facet: new(false), Index: new(true), Sort: new(true), Optional: new(true)},

			// Library status
			{Name: "has_file", Type: "bool", Facet: new(true), Index: new(true), Optional: new(true)},

			// Timestamps
			{Name: "created_at", Type: "int64", Facet: new(false), Index: new(false), Optional: new(true)},
			{Name: "updated_at", Type: "int64", Facet: new(false), Index: new(false), Optional: new(true)},
		},
	}
}
