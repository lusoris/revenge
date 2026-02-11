// Package content provides shared domain types used across content modules (movie, tvshow).
package content

// ExternalRating represents a rating from an external source (e.g. IMDb, Rotten Tomatoes).
// This type is shared between movie and tvshow modules to avoid duplication.
type ExternalRating struct {
	Source string  `json:"source"` // e.g. "Internet Movie Database", "Rotten Tomatoes", "Metacritic", "TMDb"
	Value  string  `json:"value"`  // e.g. "8.8/10", "96%", "90/100"
	Score  float64 `json:"score"`  // Normalized 0-100 scale
}

// GenreSummary represents a distinct genre with the count of items tagged with it.
// Used by listing endpoints that aggregate genres across content.
type GenreSummary struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	ItemCount int64  `json:"item_count"`
}
