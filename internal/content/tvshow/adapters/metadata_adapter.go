// Package adapters provides adapter implementations that bridge TV show-specific
// functionality with shared frameworks.
package adapters

import (
	"github.com/lusoris/revenge/internal/content/shared/metadata"
	"golang.org/x/time/rate"
)

// TMDbBaseURL is the base URL for the TMDb API.
const TMDbBaseURL = "https://api.themoviedb.org/3"

// TVShowTMDbClientConfig wraps metadata.ClientConfig with TV show-specific defaults.
type TVShowTMDbClientConfig struct {
	APIKey   string
	ProxyURL string
}

// NewTVShowTMDbClient creates a metadata.BaseClient configured for TMDb TV API access.
// This provides rate limiting, caching, and retry functionality from the shared package.
func NewTVShowTMDbClient(config TVShowTMDbClientConfig) *metadata.BaseClient {
	clientConfig := metadata.ClientConfig{
		BaseURL:   TMDbBaseURL,
		APIKey:    config.APIKey,
		RateLimit: rate.Limit(4.0), // TMDb rate limit: 40 requests per 10 seconds
		RateBurst: 10,
		ProxyURL:  config.ProxyURL,
	}

	return metadata.NewBaseClient(clientConfig)
}

// NewTVShowImageDownloader creates an ImageDownloader for TV show poster/backdrop downloads.
func NewTVShowImageDownloader(client *metadata.BaseClient) *metadata.ImageDownloader {
	return metadata.NewImageDownloader(client)
}

// NewTVShowImageURLBuilder creates an ImageURLBuilder for constructing TV show image URLs.
func NewTVShowImageURLBuilder() *metadata.ImageURLBuilder {
	return metadata.NewImageURLBuilder()
}

// TVShowTMDbEndpoints contains the TMDb API endpoints for TV shows.
var TVShowTMDbEndpoints = struct {
	// Series endpoints
	SearchTV        string
	TVDetails       string
	TVCredits       string
	TVImages        string
	TVAlternative   string
	TVContentRatings string
	TVTranslations  string
	TVExternalIDs   string

	// Season endpoints
	SeasonDetails      string
	SeasonCredits      string
	SeasonImages       string
	SeasonTranslations string

	// Episode endpoints
	EpisodeDetails      string
	EpisodeCredits      string
	EpisodeImages       string
	EpisodeTranslations string
}{
	// Series endpoints
	SearchTV:         "/search/tv",
	TVDetails:        "/tv/%d",                   // tv_id
	TVCredits:        "/tv/%d/credits",           // tv_id
	TVImages:         "/tv/%d/images",            // tv_id
	TVAlternative:    "/tv/%d/alternative_titles", // tv_id
	TVContentRatings: "/tv/%d/content_ratings",   // tv_id
	TVTranslations:   "/tv/%d/translations",      // tv_id
	TVExternalIDs:    "/tv/%d/external_ids",      // tv_id

	// Season endpoints
	SeasonDetails:      "/tv/%d/season/%d",              // tv_id, season_number
	SeasonCredits:      "/tv/%d/season/%d/credits",      // tv_id, season_number
	SeasonImages:       "/tv/%d/season/%d/images",       // tv_id, season_number
	SeasonTranslations: "/tv/%d/season/%d/translations", // tv_id, season_number

	// Episode endpoints
	EpisodeDetails:      "/tv/%d/season/%d/episode/%d",              // tv_id, season_number, episode_number
	EpisodeCredits:      "/tv/%d/season/%d/episode/%d/credits",      // tv_id, season_number, episode_number
	EpisodeImages:       "/tv/%d/season/%d/episode/%d/images",       // tv_id, season_number, episode_number
	EpisodeTranslations: "/tv/%d/season/%d/episode/%d/translations", // tv_id, season_number, episode_number
}

// TVShowTMDbAppendToResponse contains common append_to_response values for TV shows.
var TVShowTMDbAppendToResponse = struct {
	// Series append options
	SeriesFull     string
	SeriesBasic    string
	SeriesImages   string
	SeriesCredits  string

	// Season append options
	SeasonFull   string
	SeasonImages string

	// Episode append options
	EpisodeFull   string
	EpisodeImages string
}{
	// Series: Get everything in one request
	SeriesFull:   "credits,images,content_ratings,external_ids,translations,alternative_titles",
	SeriesBasic:  "credits,content_ratings,external_ids",
	SeriesImages: "images",
	SeriesCredits: "credits",

	// Season
	SeasonFull:   "credits,images,translations",
	SeasonImages: "images",

	// Episode
	EpisodeFull:   "credits,images,translations",
	EpisodeImages: "images",
}

// TVShowGenreMap maps TMDb genre IDs to genre names for TV shows.
// Note: TV show genres differ from movie genres on TMDb.
var TVShowGenreMap = map[int]string{
	10759: "Action & Adventure",
	16:    "Animation",
	35:    "Comedy",
	80:    "Crime",
	99:    "Documentary",
	18:    "Drama",
	10751: "Family",
	10762: "Kids",
	9648:  "Mystery",
	10763: "News",
	10764: "Reality",
	10765: "Sci-Fi & Fantasy",
	10766: "Soap",
	10767: "Talk",
	10768: "War & Politics",
	37:    "Western",
}

// GetTVShowGenreName returns the genre name for a given TMDb genre ID.
// Returns empty string if genre ID is not found.
func GetTVShowGenreName(genreID int) string {
	if name, ok := TVShowGenreMap[genreID]; ok {
		return name
	}
	return ""
}

// TVShowStatusMap maps TMDb status values to standardized status strings.
var TVShowStatusMap = map[string]string{
	"Returning Series": "Returning Series",
	"Ended":            "Ended",
	"Canceled":         "Canceled",
	"In Production":    "In Production",
	"Planned":          "Planned",
	"Pilot":            "Pilot",
}

// GetTVShowStatus returns the standardized status for a TMDb status value.
func GetTVShowStatus(tmdbStatus string) string {
	if status, ok := TVShowStatusMap[tmdbStatus]; ok {
		return status
	}
	return tmdbStatus
}

// TVShowTypeMap maps TMDb type values to standardized type strings.
var TVShowTypeMap = map[string]string{
	"Scripted":    "Scripted",
	"Reality":     "Reality",
	"Documentary": "Documentary",
	"Miniseries":  "Miniseries",
	"Talk Show":   "Talk Show",
	"News":        "News",
	"Video":       "Video",
}

// GetTVShowType returns the standardized type for a TMDb type value.
func GetTVShowType(tmdbType string) string {
	if t, ok := TVShowTypeMap[tmdbType]; ok {
		return t
	}
	return tmdbType
}
