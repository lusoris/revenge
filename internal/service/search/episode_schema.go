package search

import (
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// EpisodeCollectionName is the name of the episodes collection in Typesense.
const EpisodeCollectionName = "episodes"

// EpisodeDocument represents an episode document in the search index.
type EpisodeDocument struct {
	ID            string  `json:"id"`             // UUID as string
	SeriesID      string  `json:"series_id"`      // Parent series UUID
	SeasonID      string  `json:"season_id"`      // Parent season UUID
	TMDbID        int32   `json:"tmdb_id"`        // TMDb ID for external reference
	TVDbID        int32   `json:"tvdb_id"`        // TVDb ID
	IMDbID        string  `json:"imdb_id"`        // IMDb ID (optional)
	SeasonNumber  int32   `json:"season_number"`  // Season number (facet + filter)
	EpisodeNumber int32   `json:"episode_number"` // Episode number within season
	Title         string  `json:"title"`          // Episode title (searchable)
	Overview      string  `json:"overview"`       // Episode overview (searchable)
	AirDate       int64   `json:"air_date"`       // Unix timestamp for sorting
	Runtime       int32   `json:"runtime"`        // Duration in minutes
	VoteAverage   float64 `json:"vote_average"`   // Rating (sortable)
	VoteCount     int32   `json:"vote_count"`     // Vote count
	StillPath     string  `json:"still_path"`     // Still image path
	HasFile       bool    `json:"has_file"`       // Whether episode has a media file (facet)

	// Series context for display in search results
	SeriesTitle      string `json:"series_title"`       // Parent series title (searchable)
	SeriesPosterPath string `json:"series_poster_path"` // Parent series poster

	// Timestamps
	CreatedAt int64 `json:"created_at"` // Document creation time
	UpdatedAt int64 `json:"updated_at"` // Document update time
}

// EpisodeCollectionSchema returns the Typesense schema for the episodes collection.
func EpisodeCollectionSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name:                EpisodeCollectionName,
		EnableNestedFields:  new(false),
		TokenSeparators:     &[]string{"-", "_"},
		SymbolsToIndex:      &[]string{"&"},
		DefaultSortingField: new("air_date"),
		Fields: []api.Field{
			// ID field (primary key)
			{Name: "id", Type: "string"},

			// Parent references (filterable)
			{Name: "series_id", Type: "string", Facet: new(false), Index: new(true)},
			{Name: "season_id", Type: "string", Facet: new(false), Index: new(true)},

			// External IDs
			{Name: "tmdb_id", Type: "int32", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "tvdb_id", Type: "int32", Facet: new(false), Index: new(true), Optional: new(true)},
			{Name: "imdb_id", Type: "string", Facet: new(false), Index: new(true), Optional: new(true)},

			// Episode numbering (facet season_number for filtering by season)
			{Name: "season_number", Type: "int32", Facet: new(true), Index: new(true)},
			{Name: "episode_number", Type: "int32", Facet: new(false), Index: new(true)},

			// Title fields (searchable with infix for partial matching)
			{Name: "title", Type: "string", Facet: new(false), Index: new(true), Infix: new(true)},
			{Name: "overview", Type: "string", Facet: new(false), Index: new(true), Optional: new(true)},

			// Air date and runtime
			{Name: "air_date", Type: "int64", Facet: new(false), Index: new(true), Sort: new(true)},
			{Name: "runtime", Type: "int32", Facet: new(false), Index: new(true), Optional: new(true)},

			// Ratings (sortable)
			{Name: "vote_average", Type: "float", Facet: new(false), Index: new(true), Sort: new(true), Optional: new(true)},
			{Name: "vote_count", Type: "int32", Facet: new(false), Index: new(true), Optional: new(true)},

			// Image
			{Name: "still_path", Type: "string", Facet: new(false), Index: new(false), Optional: new(true)},

			// Library status
			{Name: "has_file", Type: "bool", Facet: new(true), Index: new(true)},

			// Series context (searchable so users can find episodes by series name)
			{Name: "series_title", Type: "string", Facet: new(false), Index: new(true), Infix: new(true)},
			{Name: "series_poster_path", Type: "string", Facet: new(false), Index: new(false), Optional: new(true)},

			// Timestamps
			{Name: "created_at", Type: "int64", Facet: new(false), Index: new(false), Optional: new(true)},
			{Name: "updated_at", Type: "int64", Facet: new(false), Index: new(false), Optional: new(true)},
		},
	}
}
