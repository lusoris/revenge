// Package metadata provides shared utilities and types for external metadata providers.
// It includes common HTTP client functionality, caching, rate limiting, and shared
// types used across movies, TV shows, and other content types.
package metadata

import (
	"context"
	"time"
)

// SearchResult represents a generic search result from a metadata provider.
// Content-type-specific packages wrap this with additional fields.
type SearchResult struct {
	ExternalID       int      // Provider-specific ID (TMDb ID, TVDb ID, etc.)
	Title            string   // Primary title
	OriginalTitle    string   // Title in original language
	OriginalLanguage string   // ISO 639-1 language code
	Overview         string   // Description/summary
	ReleaseDate      string   // Release/air date (ISO format)
	PosterPath       *string  // Poster image path
	BackdropPath     *string  // Backdrop image path
	VoteAverage      float64  // Rating (0-10 scale)
	VoteCount        int      // Number of votes
	Popularity       float64  // Provider-specific popularity score
	Adult            bool     // Adult content flag
	GenreIDs         []int    // Provider-specific genre IDs
	MediaType        string   // "movie", "tv", "person", etc.
}

// Genre represents a content genre from a metadata provider.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CastMember represents an actor/performer credit.
type CastMember struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Character   string  `json:"character"`
	Order       int     `json:"order"`
	CreditID    string  `json:"credit_id"`
	Gender      *int    `json:"gender"`
	ProfilePath *string `json:"profile_path"`
}

// CrewMember represents a crew (non-cast) credit.
type CrewMember struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Job         string  `json:"job"`
	Department  string  `json:"department"`
	CreditID    string  `json:"credit_id"`
	Gender      *int    `json:"gender"`
	ProfilePath *string `json:"profile_path"`
}

// Credits contains cast and crew information.
type Credits struct {
	ID   int          `json:"id"`
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

// Image represents an image from a metadata provider.
type Image struct {
	AspectRatio float64 `json:"aspect_ratio"`
	FilePath    string  `json:"file_path"`
	Height      int     `json:"height"`
	Width       int     `json:"width"`
	LanguageISO *string `json:"iso_639_1"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

// ImageCollection contains images categorized by type.
type ImageCollection struct {
	ID        int     `json:"id"`
	Backdrops []Image `json:"backdrops"`
	Posters   []Image `json:"posters"`
	Logos     []Image `json:"logos"`
}

// ProductionCompany represents a company involved in production.
type ProductionCompany struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	LogoPath      *string `json:"logo_path"`
	OriginCountry string  `json:"origin_country"`
}

// ProductionCountry represents a country where content was produced.
type ProductionCountry struct {
	ISOCode string `json:"iso_3166_1"`
	Name    string `json:"name"`
}

// SpokenLanguage represents a language spoken in the content.
type SpokenLanguage struct {
	ISOCode     string `json:"iso_639_1"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
}

// CacheEntry wraps cached data with an expiration time.
type CacheEntry struct {
	Data      any
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired.
func (c *CacheEntry) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// ReleaseInfo contains release date and certification information.
type ReleaseInfo struct {
	CountryISO    string // Country code (US, DE, GB, etc.)
	Certification string // Age rating (R, PG-13, FSK 12, etc.)
	ReleaseDate   string // ISO 8601 date
	ReleaseType   int    // Release type (1=Premiere, 2=Limited, 3=Theatrical, etc.)
	LanguageISO   string // Language code
	Note          string // Optional note
}

// MultiLanguageResult contains content data fetched in multiple languages.
type MultiLanguageResult[T any] struct {
	Items map[string]T // Key is language code (en, de, fr, etc.)
}

// Provider is the interface for metadata providers.
// Content-specific implementations (movie, TV) extend this with type-specific methods.
type Provider interface {
	// GetImageURL constructs a full image URL from a path and size.
	GetImageURL(path string, size string) string

	// DownloadImage downloads an image by path and size.
	DownloadImage(ctx context.Context, path string, size string) ([]byte, error)

	// ClearCache clears all cached data.
	ClearCache()
}

// SearchProvider extends Provider with search capabilities.
type SearchProvider interface {
	Provider

	// Search performs a generic search query.
	Search(ctx context.Context, query string, options SearchOptions) ([]SearchResult, error)
}

// SearchOptions configures a search query.
type SearchOptions struct {
	Year       *int   // Filter by release year
	Language   string // Response language (default: en-US)
	Page       int    // Page number for pagination
	Adult      bool   // Include adult content
	Region     string // ISO 3166-1 region code
	MediaTypes []string // Filter by media types ("movie", "tv", etc.)
}

// DefaultSearchOptions returns sensible defaults for search.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		Language: "en-US",
		Page:     1,
		Adult:    false,
	}
}
