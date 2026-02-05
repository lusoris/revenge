package sonarr

import "errors"

// Common errors returned by the Sonarr client.
var (
	// ErrSeriesNotFound is returned when a series is not found.
	ErrSeriesNotFound = errors.New("sonarr: series not found")

	// ErrEpisodeNotFound is returned when an episode is not found.
	ErrEpisodeNotFound = errors.New("sonarr: episode not found")

	// ErrEpisodeFileNotFound is returned when an episode file is not found.
	ErrEpisodeFileNotFound = errors.New("sonarr: episode file not found")

	// ErrUnauthorized is returned when the API key is invalid.
	ErrUnauthorized = errors.New("sonarr: unauthorized - check API key")

	// ErrConnectionFailed is returned when connection to Sonarr fails.
	ErrConnectionFailed = errors.New("sonarr: connection failed")

	// ErrRateLimited is returned when rate limited by Sonarr.
	ErrRateLimited = errors.New("sonarr: rate limited")
)
