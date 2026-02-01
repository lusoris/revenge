package errors

import (
	"fmt"

	"github.com/go-faster/errors"
)

// Wrapf wraps an error with a formatted message and preserves the stack trace.
// If err is nil, Wrapf returns nil.
//
// Example:
//
//	if err := fetchUser(id); err != nil {
//	    return errors.Wrapf(err, "failed to fetch user %s", id)
//	}
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, format, args...)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
//
// Example:
//
//	if err != nil {
//	    return errors.WithStack(err)
//	}
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "")
}

// WrapSentinel wraps a sentinel error with additional context.
// This is useful for adding context to predefined sentinel errors like ErrNotFound.
//
// Example:
//
//	if user == nil {
//	    return errors.WrapSentinel(ErrNotFound, "user", id)
//	}
func WrapSentinel(sentinel error, resource string, id interface{}) error {
	return errors.Wrapf(sentinel, "%s %v", resource, id)
}

// FormatError returns a formatted error message with stack trace if available.
// Useful for logging errors with full context.
//
// Example:
//
//	logger.Error("operation failed", "error", errors.FormatError(err))
func FormatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%+v", err)
}
