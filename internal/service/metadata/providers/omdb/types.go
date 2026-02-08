package omdb

// Response is the OMDb API response for a single title lookup.
type Response struct {
	Title      string   `json:"Title"`
	Year       string   `json:"Year"`
	Rated      string   `json:"Rated"`
	Released   string   `json:"Released"`
	Runtime    string   `json:"Runtime"`
	Genre      string   `json:"Genre"`
	Director   string   `json:"Director"`
	Writer     string   `json:"Writer"`
	Actors     string   `json:"Actors"`
	Plot       string   `json:"Plot"`
	Language   string   `json:"Language"`
	Country    string   `json:"Country"`
	Awards     string   `json:"Awards"`
	Poster     string   `json:"Poster"`
	Ratings    []Rating `json:"Ratings"`
	Metascore  string   `json:"Metascore"`
	IMDbRating string   `json:"imdbRating"`
	IMDbVotes  string   `json:"imdbVotes"`
	IMDbID     string   `json:"imdbID"`
	Type       string   `json:"Type"` // "movie", "series", "episode"
	DVD        string   `json:"DVD"`
	BoxOffice  string   `json:"BoxOffice"`
	Production string   `json:"Production"`
	Website    string   `json:"Website"`
	Response   string   `json:"Response"` // "True" or "False"
	Error      string   `json:"Error"`

	// Series-specific fields
	TotalSeasons string `json:"totalSeasons"`
}

// Rating represents an individual rating source.
type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

// SearchResponse is the OMDb API response for a search query.
type SearchResponse struct {
	Search       []SearchResult `json:"Search"`
	TotalResults string         `json:"totalResults"`
	Response     string         `json:"Response"`
	Error        string         `json:"Error"`
}

// SearchResult is a single search result entry.
type SearchResult struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	IMDbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}
