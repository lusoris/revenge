// Package errors provides sentinel errors and error wrapping utilities for the Revenge server.
package errors

import (
	stderrors "errors"

	"github.com/go-faster/errors"
)

// Sentinel errors for common use cases.
// These can be used with errors.Is() for error detection.
var (
	// ErrNotFound indicates a resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrUnauthorized indicates authentication is required.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates the user does not have permission.
	ErrForbidden = errors.New("forbidden")

	// ErrConflict indicates a resource conflict (e.g., duplicate).
	ErrConflict = errors.New("conflict")

	// ErrValidation indicates validation failed.
	ErrValidation = errors.New("validation failed")

	// ErrInternal indicates an internal server error.
	ErrInternal = errors.New("internal server error")

	// ErrBadRequest indicates a malformed request.
	ErrBadRequest = errors.New("bad request")

	// ErrUnavailable indicates a service is temporarily unavailable.
	ErrUnavailable = errors.New("service unavailable")

	// ErrTimeout indicates an operation timed out.
	ErrTimeout = errors.New("timeout")
)

// New creates a new error with a message.
// The error includes a stack trace.
func New(msg string) error {
	return errors.New(msg)
}

// Wrap wraps an error with additional context.
// The original error can be retrieved with errors.Unwrap().
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, msg)
}

// Errorf creates a new error with a formatted message.
// Use %w to wrap an error and preserve the error chain.
func Errorf(format string, args ...any) error {
	return errors.Errorf(format, args...)
}

// Is reports whether any error in err's tree matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's tree that matches target,
// and if one is found, sets target to that error value and returns true.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// AsType is a type-safe alternative to [As] that uses generics.
// It finds the first error in err's tree that matches type E.
//
// Example:
//
//	if pe, ok := errors.AsType[*ProviderError](err); ok {
//	    log.Printf("provider %s failed: %v", pe.Provider, pe.Err)
//	}
func AsType[E error](err error) (E, bool) {
	return stderrors.AsType[E](err)
}

// Unwrap returns the result of calling the Unwrap method on err,
// if err's type contains an Unwrap method returning error.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
