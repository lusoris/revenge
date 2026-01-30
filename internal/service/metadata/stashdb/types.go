// Package stashdb provides a client for the StashDB GraphQL API.
// QAR obfuscation: This provider fetches metadata for adult content.
package stashdb

import "time"

// Scene represents a scene from StashDB.
type Scene struct {
	ID             string              `json:"id"`
	Title          string              `json:"title"`
	Details        string              `json:"details"`
	Date           string              `json:"date"`
	ReleaseDate    string              `json:"release_date"`
	ProductionDate string              `json:"production_date"`
	Duration       int                 `json:"duration"` // in seconds
	Director       string              `json:"director"`
	Code           string              `json:"code"`
	Studio         *Studio             `json:"studio"`
	Performers     []PerformerAppearance `json:"performers"`
	Tags           []Tag               `json:"tags"`
	Images         []Image             `json:"images"`
	Fingerprints   []Fingerprint       `json:"fingerprints"`
	URLs           []URL               `json:"urls"`
	Created        time.Time           `json:"created"`
	Updated        time.Time           `json:"updated"`
}

// Performer represents a performer from StashDB.
type Performer struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Disambiguation  string   `json:"disambiguation"`
	Aliases         []string `json:"aliases"`
	Gender          string   `json:"gender"`
	Birthdate       string   `json:"birth_date"`
	DeathDate       string   `json:"death_date"`
	Age             int      `json:"age"`
	Ethnicity       string   `json:"ethnicity"`
	Country         string   `json:"country"`
	EyeColor        string   `json:"eye_color"`
	HairColor       string   `json:"hair_color"`
	Height          int      `json:"height"` // in cm
	CupSize         string   `json:"cup_size"`
	BandSize        int      `json:"band_size"`
	WaistSize       int      `json:"waist_size"`
	HipSize         int      `json:"hip_size"`
	BreastType      string   `json:"breast_type"`
	CareerStartYear int      `json:"career_start_year"`
	CareerEndYear   int      `json:"career_end_year"`
	Images          []Image  `json:"images"`
	SceneCount      int      `json:"scene_count"`
}

// PerformerAppearance represents a performer's appearance in a scene.
type PerformerAppearance struct {
	Performer Performer `json:"performer"`
	As        string    `json:"as"` // Alias used in this scene
}

// Studio represents a studio from StashDB.
type Studio struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	URLs     []URL   `json:"urls"`
	Parent   *Studio `json:"parent"`
	Images   []Image `json:"images"`
}

// Tag represents a tag from StashDB.
type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// Image represents an image from StashDB.
type Image struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Fingerprint represents a scene fingerprint for matching.
type Fingerprint struct {
	Algorithm string `json:"algorithm"` // PHASH, OSHASH, MD5
	Hash      string `json:"hash"`
	Duration  int    `json:"duration"`
}

// URL represents a URL with type information.
type URL struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// QueryResult wraps paginated query results.
type QueryResult[T any] struct {
	Count int `json:"count"`
	Data  []T `json:"data"`
}

// SceneQueryResult is the result of a scene query.
type SceneQueryResult = QueryResult[Scene]

// PerformerQueryResult is the result of a performer query.
type PerformerQueryResult = QueryResult[Performer]

// FingerprintMatch represents a fingerprint match result.
type FingerprintMatch struct {
	Scene Scene `json:"scene"`
}

// GraphQL request/response types.

// GraphQLRequest represents a GraphQL request.
type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response.
type GraphQLResponse[T any] struct {
	Data   T              `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error.
type GraphQLError struct {
	Message    string         `json:"message"`
	Locations  []ErrorLocation `json:"locations,omitempty"`
	Path       []any          `json:"path,omitempty"`
	Extensions map[string]any `json:"extensions,omitempty"`
}

// ErrorLocation represents the location of a GraphQL error.
type ErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Error returns the error message.
func (e GraphQLError) Error() string {
	return e.Message
}
