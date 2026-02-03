package movie

import "errors"

// Error definitions for movie operations
var (
	ErrMovieNotFound    = errors.New("movie not found")
	ErrProgressNotFound = errors.New("watch progress not found")
	ErrNotInCollection  = errors.New("movie is not in a collection")
)
