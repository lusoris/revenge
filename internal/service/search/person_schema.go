package search

import (
	"github.com/lusoris/revenge/internal/util/ptr"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// PersonCollectionName is the Typesense collection name for people.
const PersonCollectionName = "people"

// PersonDocument represents a person in the search index, aggregated from all credit sources.
type PersonDocument struct {
	ID           string   `json:"id"`            // tmdb_person_id as string
	TMDbID       int32    `json:"tmdb_id"`       // TMDb person ID
	Name         string   `json:"name"`          // Person name
	ProfilePath  string   `json:"profile_path"`  // Profile image path
	KnownFor     []string `json:"known_for"`     // Movie/show titles they appeared in
	Characters   []string `json:"characters"`    // Character names played
	Departments  []string `json:"departments"`   // Unique departments: Acting, Directing, etc.
	MovieCount   int32    `json:"movie_count"`   // Number of movie credits
	TVShowCount  int32    `json:"tvshow_count"`  // Number of TV show credits
	TotalCredits int32    `json:"total_credits"` // Total credit count
}

// PersonCollectionSchema returns the Typesense collection schema for people.
func PersonCollectionSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name:                PersonCollectionName,
		DefaultSortingField: ptr.To("total_credits"),
		TokenSeparators:     &[]string{"-", "_"},
		SymbolsToIndex:      &[]string{"&"},
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "tmdb_id", Type: "int32"},
			{Name: "name", Type: "string", Infix: ptr.To(true)},
			{Name: "profile_path", Type: "string", Index: ptr.To(false), Optional: ptr.To(true)},
			{Name: "known_for", Type: "string[]", Optional: ptr.To(true)},
			{Name: "characters", Type: "string[]", Optional: ptr.To(true)},
			{Name: "departments", Type: "string[]", Facet: ptr.To(true), Optional: ptr.To(true)},
			{Name: "movie_count", Type: "int32", Sort: ptr.To(true)},
			{Name: "tvshow_count", Type: "int32", Sort: ptr.To(true)},
			{Name: "total_credits", Type: "int32", Sort: ptr.To(true)},
		},
	}
}
