package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/errors"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrNotFound", errors.ErrNotFound, "not found"},
		{"ErrUnauthorized", errors.ErrUnauthorized, "unauthorized"},
		{"ErrForbidden", errors.ErrForbidden, "forbidden"},
		{"ErrConflict", errors.ErrConflict, "conflict"},
		{"ErrValidation", errors.ErrValidation, "validation failed"},
		{"ErrInternal", errors.ErrInternal, "internal server error"},
		{"ErrBadRequest", errors.ErrBadRequest, "bad request"},
		{"ErrUnavailable", errors.ErrUnavailable, "service unavailable"},
		{"ErrTimeout", errors.ErrTimeout, "timeout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err, "sentinel error should not be nil")
			assert.Equal(t, tt.msg, tt.err.Error(), "error message should match")
		})
	}
}

func TestNew(t *testing.T) {
	msg := "test error"
	err := errors.New(msg)

	require.NotNil(t, err)
	assert.Equal(t, msg, err.Error())
}

func TestWrap(t *testing.T) {
	t.Run("wrap non-nil error", func(t *testing.T) {
		baseErr := errors.New("base error")
		wrappedErr := errors.Wrap(baseErr, "additional context")

		require.NotNil(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "additional context")
		assert.Contains(t, wrappedErr.Error(), "base error")
		assert.True(t, errors.Is(wrappedErr, baseErr))
	})

	t.Run("wrap nil error returns nil", func(t *testing.T) {
		wrappedErr := errors.Wrap(nil, "context")
		assert.Nil(t, wrappedErr)
	})

	t.Run("wrap sentinel error", func(t *testing.T) {
		wrappedErr := errors.Wrap(errors.ErrNotFound, "user not found")

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, errors.ErrNotFound))
		assert.Contains(t, wrappedErr.Error(), "user not found")
	})
}

func TestErrorf(t *testing.T) {
	t.Run("format with no wrapping", func(t *testing.T) {
		err := errors.Errorf("failed to process item %d", 42)

		require.NotNil(t, err)
		assert.Equal(t, "failed to process item 42", err.Error())
	})

	t.Run("format with wrapping using Wrap", func(t *testing.T) {
		baseErr := errors.New("base error")
		err := errors.Wrap(baseErr, "wrapped")

		require.NotNil(t, err)
		assert.True(t, errors.Is(err, baseErr))
		assert.Contains(t, err.Error(), "wrapped")
		assert.Contains(t, err.Error(), "base error")
	})

	t.Run("format with sentinel error wrapping", func(t *testing.T) {
		err := errors.Wrap(errors.ErrNotFound, "resource not found")

		require.NotNil(t, err)
		assert.True(t, errors.Is(err, errors.ErrNotFound))
	})
}

func TestIs(t *testing.T) {
	t.Run("direct match", func(t *testing.T) {
		assert.True(t, errors.Is(errors.ErrNotFound, errors.ErrNotFound))
	})

	t.Run("wrapped match", func(t *testing.T) {
		wrappedErr := errors.Wrap(errors.ErrUnauthorized, "context")
		assert.True(t, errors.Is(wrappedErr, errors.ErrUnauthorized))
	})

	t.Run("no match", func(t *testing.T) {
		assert.False(t, errors.Is(errors.ErrNotFound, errors.ErrForbidden))
	})

	t.Run("nil error", func(t *testing.T) {
		assert.False(t, errors.Is(nil, errors.ErrNotFound))
	})

	t.Run("deeply nested wrapped error", func(t *testing.T) {
		baseErr := errors.ErrValidation
		err1 := errors.Wrap(baseErr, "level 1")
		err2 := errors.Wrap(err1, "level 2")
		err3 := errors.Wrap(err2, "level 3")

		assert.True(t, errors.Is(err3, errors.ErrValidation))
	})
}

// customError is a test error type that implements error interface
type customError struct {
	code int
	msg  string
}

func (e *customError) Error() string {
	return e.msg
}

func TestAs(t *testing.T) {
	t.Run("custom error type found", func(t *testing.T) {
		customErr := &customError{code: 404, msg: "not found"}
		wrappedErr := errors.Wrap(customErr, "wrapped")

		var target *customError
		found := errors.As(wrappedErr, &target)

		assert.True(t, found)
		require.NotNil(t, target)
		assert.Equal(t, 404, target.code)
		assert.Equal(t, "not found", target.msg)
	})

	t.Run("custom error type not found", func(t *testing.T) {
		baseErr := errors.New("simple error")
		var target *customError
		found := errors.As(baseErr, &target)

		assert.False(t, found)
		assert.Nil(t, target)
	})
}

func TestAsType(t *testing.T) {
	t.Run("custom error type found", func(t *testing.T) {
		customErr := &customError{code: 404, msg: "not found"}
		wrappedErr := errors.Wrap(customErr, "wrapped")

		target, ok := errors.AsType[*customError](wrappedErr)

		assert.True(t, ok)
		require.NotNil(t, target)
		assert.Equal(t, 404, target.code)
		assert.Equal(t, "not found", target.msg)
	})

	t.Run("custom error type not found", func(t *testing.T) {
		baseErr := errors.New("simple error")

		_, ok := errors.AsType[*customError](baseErr)
		assert.False(t, ok)
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("unwrap wrapped error", func(t *testing.T) {
		baseErr := errors.New("base")
		wrappedErr := errors.Wrap(baseErr, "wrapped")

		unwrapped := errors.Unwrap(wrappedErr)
		assert.NotNil(t, unwrapped)
		assert.True(t, errors.Is(unwrapped, baseErr))
	})

	t.Run("unwrap sentinel error returns nil", func(t *testing.T) {
		unwrapped := errors.Unwrap(errors.ErrNotFound)
		assert.Nil(t, unwrapped)
	})

	t.Run("unwrap nil returns nil", func(t *testing.T) {
		unwrapped := errors.Unwrap(nil)
		assert.Nil(t, unwrapped)
	})

	t.Run("multiple levels of wrapping", func(t *testing.T) {
		baseErr := errors.New("base")
		err1 := errors.Wrap(baseErr, "level 1")
		err2 := errors.Wrap(err1, "level 2")

		unwrapped1 := errors.Unwrap(err2)
		assert.True(t, errors.Is(unwrapped1, err1))

		unwrapped2 := errors.Unwrap(unwrapped1)
		assert.True(t, errors.Is(unwrapped2, baseErr))

		unwrapped3 := errors.Unwrap(unwrapped2)
		assert.Nil(t, unwrapped3)
	})
}

// TestErrorChaining verifies complete error chain behavior
func TestErrorChaining(t *testing.T) {
	baseErr := errors.ErrNotFound
	err1 := errors.Wrap(baseErr, "database query failed")
	err2 := errors.Wrap(err1, "user service error")
	//nolint:govet // go-faster/errors supports %w but Go vet doesn't know that
	err3 := errors.Wrap(err2, "API handler failed")

	// Should be able to detect sentinel error through all layers
	assert.True(t, errors.Is(err3, errors.ErrNotFound))

	// Error message should contain all context
	errMsg := err3.Error()
	assert.Contains(t, errMsg, "API handler failed")
	assert.Contains(t, errMsg, "user service error")
	assert.Contains(t, errMsg, "database query failed")
}
