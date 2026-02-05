// Package adapters provides adapter implementations that bridge movie-specific
// functionality with shared frameworks.
package adapters

import (
	"github.com/lusoris/revenge/internal/content/shared/metadata"
	"golang.org/x/time/rate"
)

// TMDbBaseURL is the base URL for the TMDb API.
const TMDbBaseURL = "https://api.themoviedb.org/3"

// MovieTMDbClientConfig wraps metadata.ClientConfig with movie-specific defaults.
type MovieTMDbClientConfig struct {
	APIKey   string
	ProxyURL string
}

// NewMovieTMDbClient creates a metadata.BaseClient configured for TMDb movie API access.
// This provides rate limiting, caching, and retry functionality from the shared package.
func NewMovieTMDbClient(config MovieTMDbClientConfig) *metadata.BaseClient {
	clientConfig := metadata.ClientConfig{
		BaseURL:   TMDbBaseURL,
		APIKey:    config.APIKey,
		RateLimit: rate.Limit(4.0), // TMDb rate limit: 40 requests per 10 seconds
		RateBurst: 10,
		ProxyURL:  config.ProxyURL,
	}

	return metadata.NewBaseClient(clientConfig)
}

// NewMovieImageDownloader creates an ImageDownloader for movie poster/backdrop downloads.
func NewMovieImageDownloader(client *metadata.BaseClient) *metadata.ImageDownloader {
	return metadata.NewImageDownloader(client)
}

// NewMovieImageURLBuilder creates an ImageURLBuilder for constructing movie image URLs.
func NewMovieImageURLBuilder() *metadata.ImageURLBuilder {
	return metadata.NewImageURLBuilder()
}
