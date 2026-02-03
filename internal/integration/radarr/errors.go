package radarr

import "errors"

// Errors for Radarr client operations.
var (
	ErrMovieNotFound = errors.New("movie not found in radarr")
	ErrNotConfigured = errors.New("radarr integration not configured")
	ErrUnauthorized  = errors.New("radarr api key is invalid")
)
