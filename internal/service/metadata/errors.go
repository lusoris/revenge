package metadata

import (
	"errors"
	"fmt"
)

// Common errors returned by the metadata service.
var (
	// ErrNotFound indicates the requested content was not found.
	ErrNotFound = errors.New("metadata: not found")

	// ErrProviderUnavailable indicates a provider is temporarily unavailable.
	ErrProviderUnavailable = errors.New("metadata: provider unavailable")

	// ErrRateLimited indicates the provider's rate limit was exceeded.
	ErrRateLimited = errors.New("metadata: rate limited")

	// ErrUnauthorized indicates invalid or missing API credentials.
	ErrUnauthorized = errors.New("metadata: unauthorized")

	// ErrInvalidID indicates the provided ID format is invalid.
	ErrInvalidID = errors.New("metadata: invalid id")

	// ErrNoProviders indicates no providers are configured.
	ErrNoProviders = errors.New("metadata: no providers configured")

	// ErrUnsupported indicates the operation is not supported by any provider.
	ErrUnsupported = errors.New("metadata: operation not supported")
)

// ProviderError wraps an error from a specific provider.
type ProviderError struct {
	Provider   ProviderID
	StatusCode int
	Message    string
	Err        error
}

func (e *ProviderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("metadata provider %s (status %d): %s: %v", e.Provider, e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("metadata provider %s (status %d): %s", e.Provider, e.StatusCode, e.Message)
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// IsProviderError checks if an error is a ProviderError.
func IsProviderError(err error) bool {
	_, ok := errors.AsType[*ProviderError](err)
	return ok
}

// NewProviderError creates a new ProviderError.
func NewProviderError(provider ProviderID, statusCode int, message string, err error) *ProviderError {
	return &ProviderError{
		Provider:   provider,
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// NewNotFoundError creates a provider-specific not found error.
func NewNotFoundError(provider ProviderID, resourceType string, id string) *ProviderError {
	return &ProviderError{
		Provider:   provider,
		StatusCode: 404,
		Message:    fmt.Sprintf("%s with id '%s' not found", resourceType, id),
		Err:        ErrNotFound,
	}
}

// NewRateLimitError creates a provider-specific rate limit error.
func NewRateLimitError(provider ProviderID, retryAfter int) *ProviderError {
	msg := "rate limit exceeded"
	if retryAfter > 0 {
		msg = fmt.Sprintf("rate limit exceeded, retry after %d seconds", retryAfter)
	}
	return &ProviderError{
		Provider:   provider,
		StatusCode: 429,
		Message:    msg,
		Err:        ErrRateLimited,
	}
}

// NewUnauthorizedError creates a provider-specific unauthorized error.
func NewUnauthorizedError(provider ProviderID) *ProviderError {
	return &ProviderError{
		Provider:   provider,
		StatusCode: 401,
		Message:    "invalid or missing API key",
		Err:        ErrUnauthorized,
	}
}

// AggregateError contains errors from multiple providers.
type AggregateError struct {
	Errors []error
}

func (e *AggregateError) Error() string {
	if len(e.Errors) == 0 {
		return "metadata: no errors"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("metadata: %d provider errors (first: %v)", len(e.Errors), e.Errors[0])
}

// Add appends an error to the aggregate.
func (e *AggregateError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// IsEmpty returns true if no errors were collected.
func (e *AggregateError) IsEmpty() bool {
	return len(e.Errors) == 0
}

// First returns the first error or nil.
func (e *AggregateError) First() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return e.Errors[0]
}

// HasNotFound returns true if all errors are not-found errors.
func (e *AggregateError) HasNotFound() bool {
	if len(e.Errors) == 0 {
		return false
	}
	for _, err := range e.Errors {
		if !errors.Is(err, ErrNotFound) {
			return false
		}
	}
	return true
}
