package arrbase

import "errors"

// Shared errors common to all arr integrations.
var (
	// ErrNotConfigured indicates the arr service is not configured.
	ErrNotConfigured = errors.New("arr service not configured")

	// ErrUnauthorized indicates the API key is invalid or missing.
	ErrUnauthorized = errors.New("unauthorized: invalid API key")

	// ErrConnectionFailed indicates the arr service is unreachable.
	ErrConnectionFailed = errors.New("connection to arr service failed")

	// ErrRateLimited indicates the client hit the rate limit.
	ErrRateLimited = errors.New("rate limited by arr service")
)
