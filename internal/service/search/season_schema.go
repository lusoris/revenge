package search

import (
	"github.com/typesense/typesense-go/v2/typesense/api"

	"github.com/lusoris/revenge/internal/util/ptr"
)

// SeasonCollectionName is the name of the seasons collection in Typesense.
const SeasonCollectionName = "seasons"

// SeasonDocument represents a season document in the search index.
type SeasonDocument struct {
	ID           string  `json:"id"`            // UUID as string
	SeriesID     string  `json:"series_id"`     // Parent series UUID
	TMDbID       int32   `json:"tmdb_id"`       // TMDb ID for external reference
	SeasonNumber int32   `json:"season_number"` // Season number (facet + filter)
	Name         string  `json:"name"`          // Season name (searchable)
	Overview     string  `json:"overview"`      // Season overview (searchable)
	AirDate      int64   `json:"air_date"`      // Unix timestamp for sorting
	EpisodeCount int32   `json:"episode_count"` // Number of episodes
	VoteAverage  float64 `json:"vote_average"`  // Rating (sortable)
	PosterPath   string  `json:"poster_path"`   // Poster image path

	// Series context for display in search results
	SeriesTitle      string `json:"series_title"`       // Parent series title (searchable)
	SeriesPosterPath string `json:"series_poster_path"` // Parent series poster

	// Timestamps
	CreatedAt int64 `json:"created_at"` // Document creation time
	UpdatedAt int64 `json:"updated_at"` // Document update time
}

// SeasonCollectionSchema returns the Typesense schema for the seasons collection.
func SeasonCollectionSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name:                SeasonCollectionName,
		EnableNestedFields:  ptr.To(false),
		TokenSeparators:     &[]string{"-", "_"},
		SymbolsToIndex:      &[]string{"&"},
		DefaultSortingField: ptr.To("air_date"),
		Fields: []api.Field{
			// ID field (primary key)
			{Name: "id", Type: "string"},

			// Parent reference (filterable)
			{Name: "series_id", Type: "string", Facet: ptr.To(false), Index: ptr.To(true)},

			// External ID
			{Name: "tmdb_id", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},

			// Season numbering (facetable for filtering)
			{Name: "season_number", Type: "int32", Facet: ptr.To(true), Index: ptr.To(true)},

			// Name and overview (searchable with infix)
			{Name: "name", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Infix: ptr.To(true)},
			{Name: "overview", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Optional: ptr.To(true)},

			// Air date and episode count
			{Name: "air_date", Type: "int64", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true)},
			{Name: "episode_count", Type: "int32", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true)},

			// Rating (sortable)
			{Name: "vote_average", Type: "float", Facet: ptr.To(false), Index: ptr.To(true), Sort: ptr.To(true), Optional: ptr.To(true)},

			// Image
			{Name: "poster_path", Type: "string", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},

			// Series context (searchable)
			{Name: "series_title", Type: "string", Facet: ptr.To(false), Index: ptr.To(true), Infix: ptr.To(true)},
			{Name: "series_poster_path", Type: "string", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},

			// Timestamps
			{Name: "created_at", Type: "int64", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},
			{Name: "updated_at", Type: "int64", Facet: ptr.To(false), Index: ptr.To(false), Optional: ptr.To(true)},
		},
	}
}
