package errors

import (
	"net/http"
)

// APIError represents an error response for the API.
// This should be returned from HTTP handlers.
type APIError struct {
	// Code is the HTTP status code.
	Code int `json:"code"`

	// Message is a human-readable error message.
	Message string `json:"message"`

	// Details provides additional error context (optional).
	Details map[string]interface{} `json:"details,omitempty"`

	// Err is the underlying error (not serialized).
	Err error `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error.
func (e *APIError) Unwrap() error {
	return e.Err
}

// NewAPIError creates a new APIError with the given code and message.
func NewAPIError(code int, message string, err error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithDetails adds details to the APIError.
func (e *APIError) WithDetails(details map[string]interface{}) *APIError {
	e.Details = details
	return e
}

// ToAPIError converts a sentinel error to an APIError.
// If the error is not recognized, it returns a 500 Internal Server Error.
func ToAPIError(err error) *APIError {
	if err == nil {
		return nil
	}

	// Check for sentinel errors
	switch {
	case Is(err, ErrNotFound):
		return NewAPIError(http.StatusNotFound, "Resource not found", err)
	case Is(err, ErrUnauthorized):
		return NewAPIError(http.StatusUnauthorized, "Authentication required", err)
	case Is(err, ErrForbidden):
		return NewAPIError(http.StatusForbidden, "Access forbidden", err)
	case Is(err, ErrConflict):
		return NewAPIError(http.StatusConflict, "Resource conflict", err)
	case Is(err, ErrValidation):
		return NewAPIError(http.StatusBadRequest, "Validation failed", err)
	case Is(err, ErrBadRequest):
		return NewAPIError(http.StatusBadRequest, "Bad request", err)
	case Is(err, ErrUnavailable):
		return NewAPIError(http.StatusServiceUnavailable, "Service unavailable", err)
	case Is(err, ErrTimeout):
		return NewAPIError(http.StatusGatewayTimeout, "Request timeout", err)
	default:
		// Unknown error - don't expose details
		return NewAPIError(http.StatusInternalServerError, "Internal server error", ErrInternal)
	}
}

// Common API error constructors

// NotFoundError creates a 404 Not Found error.
func NotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, message, ErrNotFound)
}

// UnauthorizedError creates a 401 Unauthorized error.
func UnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, message, ErrUnauthorized)
}

// ForbiddenError creates a 403 Forbidden error.
func ForbiddenError(message string) *APIError {
	return NewAPIError(http.StatusForbidden, message, ErrForbidden)
}

// ConflictError creates a 409 Conflict error.
func ConflictError(message string) *APIError {
	return NewAPIError(http.StatusConflict, message, ErrConflict)
}

// ValidationError creates a 400 Bad Request error for validation failures.
func ValidationError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message, ErrValidation)
}

// BadRequestError creates a 400 Bad Request error.
func BadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message, ErrBadRequest)
}

// InternalError creates a 500 Internal Server Error.
func InternalError(message string, err error) *APIError {
	return NewAPIError(http.StatusInternalServerError, message, err)
}

// UnavailableError creates a 503 Service Unavailable error.
func UnavailableError(message string) *APIError {
	return NewAPIError(http.StatusServiceUnavailable, message, ErrUnavailable)
}

// TimeoutError creates a 504 Gateway Timeout error.
func TimeoutError(message string) *APIError {
	return NewAPIError(http.StatusGatewayTimeout, message, ErrTimeout)
}
