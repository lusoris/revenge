package simkl

import "time"

// SearchResult is a single result from the Simkl search API.
type SearchResult struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Year      int    `json:"year,omitempty"`
	Poster    string `json:"poster,omitempty"`
	Fanart    string `json:"fanart,omitempty"`
	AnimeType string `json:"anime_type,omitempty"`
	IDs       IDs    `json:"ids"`
}

// IDs contains cross-referenced IDs for a media item.
type IDs struct {
	Simkl   int    `json:"simkl"`
	Slug    string `json:"slug,omitempty"`
	IMDb    string `json:"imdb,omitempty"`
	TMDb    int    `json:"tmdb,omitempty"`
	TVDb    int    `json:"tvdb,omitempty"`
	MAL     int    `json:"mal,omitempty"`
	AniDB   int    `json:"anidb,omitempty"`
	AniList int    `json:"anilist,omitempty"`
	Kitsu   int    `json:"kitsu,omitempty"`
}

// Movie represents a Simkl movie with full details.
type Movie struct {
	Title         string     `json:"title"`
	Year          int        `json:"year,omitempty"`
	Type          string     `json:"type,omitempty"`
	IDs           IDs        `json:"ids"`
	Overview      string     `json:"overview,omitempty"`
	Poster        string     `json:"poster,omitempty"`
	Fanart        string     `json:"fanart,omitempty"`
	Genres        []string   `json:"genres,omitempty"`
	Country       string     `json:"country,omitempty"`
	Runtime       int        `json:"runtime,omitempty"`
	Status        string     `json:"status,omitempty"`
	Certification string     `json:"certification,omitempty"`
	ReleaseDate   string     `json:"release_date,omitempty"`
	Trailer       string     `json:"trailer,omitempty"`
	Ratings       *Ratings   `json:"ratings,omitempty"`
	Network       string     `json:"network,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// Show represents a Simkl TV show with full details.
type Show struct {
	Title         string     `json:"title"`
	Year          int        `json:"year,omitempty"`
	Type          string     `json:"type,omitempty"`
	IDs           IDs        `json:"ids"`
	Overview      string     `json:"overview,omitempty"`
	Poster        string     `json:"poster,omitempty"`
	Fanart        string     `json:"fanart,omitempty"`
	Genres        []string   `json:"genres,omitempty"`
	Country       string     `json:"country,omitempty"`
	Runtime       int        `json:"runtime,omitempty"`
	Status        string     `json:"status,omitempty"`
	Certification string     `json:"certification,omitempty"`
	Network       string     `json:"network,omitempty"`
	Trailer       string     `json:"trailer,omitempty"`
	Ratings       *Ratings   `json:"ratings,omitempty"`
	TotalEpisodes int        `json:"total_episodes,omitempty"`
	AnimeType     string     `json:"anime_type,omitempty"`
	ENTitle       string     `json:"en_title,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// Episode represents a Simkl episode.
type Episode struct {
	Title   string     `json:"title,omitempty"`
	Season  int        `json:"season"`
	Episode int        `json:"episode"`
	Type    string     `json:"type,omitempty"`
	Img     string     `json:"img,omitempty"`
	Date    *time.Time `json:"date,omitempty"`
	IDs     EpisodeIDs `json:"ids,omitempty"`
}

// EpisodeIDs contains IDs for an episode.
type EpisodeIDs struct {
	Simkl int `json:"simkl_id,omitempty"`
}

// Ratings contains rating information from Simkl and external sources.
type Ratings struct {
	Simkl *RatingInfo `json:"simkl,omitempty"`
	IMDb  *RatingInfo `json:"imdb,omitempty"`
	MAL   *RatingInfo `json:"mal,omitempty"`
}

// RatingInfo is a single rating source.
type RatingInfo struct {
	Rating float64 `json:"rating"`
	Votes  int     `json:"votes,omitempty"`
	Rank   int     `json:"rank,omitempty"`
}

// IDLookupResult represents a result from the ID lookup endpoint.
type IDLookupResult struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Year      int    `json:"year,omitempty"`
	Poster    string `json:"poster,omitempty"`
	Fanart    string `json:"fanart,omitempty"`
	AnimeType string `json:"anime_type,omitempty"`
	IDs       IDs    `json:"ids"`
}
