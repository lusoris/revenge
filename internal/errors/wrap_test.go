package errors_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/errors"
)

func TestWrapf(t *testing.T) {
	t.Run("wrap non-nil error with format", func(t *testing.T) {
		baseErr := errors.New("database error")
		wrappedErr := errors.Wrapf(baseErr, "failed to fetch user %s", "john")

		require.NotNil(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to fetch user john")
		assert.Contains(t, wrappedErr.Error(), "database error")
		assert.True(t, errors.Is(wrappedErr, baseErr))
	})

	t.Run("wrap nil error returns nil", func(t *testing.T) {
		wrappedErr := errors.Wrapf(nil, "context %d", 42)
		assert.Nil(t, wrappedErr)
	})

	t.Run("wrap sentinel error with format", func(t *testing.T) {
		wrappedErr := errors.Wrapf(errors.ErrNotFound, "resource %s with id %d", "user", 123)

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, errors.ErrNotFound))
		assert.Contains(t, wrappedErr.Error(), "resource user with id 123")
		assert.Contains(t, wrappedErr.Error(), "not found")
	})

	t.Run("wrap with empty format string", func(t *testing.T) {
		baseErr := errors.New("base")
		wrappedErr := errors.Wrapf(baseErr, "")

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, baseErr))
	})

	t.Run("wrap with multiple format arguments", func(t *testing.T) {
		baseErr := errors.New("connection failed")
		wrappedErr := errors.Wrapf(baseErr, "host=%s port=%d timeout=%v", "localhost", 5432, "30s")

		require.NotNil(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "host=localhost")
		assert.Contains(t, wrappedErr.Error(), "port=5432")
		assert.Contains(t, wrappedErr.Error(), "timeout=30s")
	})

	t.Run("wrap deeply nested error with format", func(t *testing.T) {
		baseErr := errors.ErrTimeout
		err1 := errors.Wrapf(baseErr, "layer %d", 1)
		err2 := errors.Wrapf(err1, "layer %d", 2)
		err3 := errors.Wrapf(err2, "layer %d", 3)

		assert.True(t, errors.Is(err3, errors.ErrTimeout))
		assert.Contains(t, err3.Error(), "layer 1")
		assert.Contains(t, err3.Error(), "layer 2")
		assert.Contains(t, err3.Error(), "layer 3")
	})
}

func TestWithStack(t *testing.T) {
	t.Run("add stack to non-nil error", func(t *testing.T) {
		baseErr := errors.New("base error")
		withStack := errors.WithStack(baseErr)

		require.NotNil(t, withStack)
		assert.True(t, errors.Is(withStack, baseErr))
	})

	t.Run("add stack to nil error returns nil", func(t *testing.T) {
		withStack := errors.WithStack(nil)
		assert.Nil(t, withStack)
	})

	t.Run("add stack to sentinel error", func(t *testing.T) {
		withStack := errors.WithStack(errors.ErrForbidden)

		require.NotNil(t, withStack)
		assert.True(t, errors.Is(withStack, errors.ErrForbidden))
	})

	t.Run("stack trace is added", func(t *testing.T) {
		baseErr := errors.New("base")
		withStack := errors.WithStack(baseErr)

		// Format with %+v to get stack trace
		formatted := errors.FormatError(withStack)
		// Stack trace should contain function name
		assert.True(t, len(formatted) > len(baseErr.Error()),
			"formatted error with stack should be longer than base error")
	})

	t.Run("double stack addition", func(t *testing.T) {
		baseErr := errors.New("base")
		stack1 := errors.WithStack(baseErr)
		stack2 := errors.WithStack(stack1)

		assert.True(t, errors.Is(stack2, baseErr))
	})
}

func TestWrapSentinel(t *testing.T) {
	t.Run("wrap ErrNotFound with string id", func(t *testing.T) {
		wrappedErr := errors.WrapSentinel(errors.ErrNotFound, "user", "abc123")

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, errors.ErrNotFound))
		assert.Contains(t, wrappedErr.Error(), "user abc123")
		assert.Contains(t, wrappedErr.Error(), "not found")
	})

	t.Run("wrap ErrNotFound with int id", func(t *testing.T) {
		wrappedErr := errors.WrapSentinel(errors.ErrNotFound, "post", 42)

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, errors.ErrNotFound))
		assert.Contains(t, wrappedErr.Error(), "post 42")
	})

	t.Run("wrap ErrConflict", func(t *testing.T) {
		wrappedErr := errors.WrapSentinel(errors.ErrConflict, "email", "test@example.com")

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, errors.ErrConflict))
		assert.Contains(t, wrappedErr.Error(), "email test@example.com")
	})

	t.Run("wrap ErrUnauthorized", func(t *testing.T) {
		wrappedErr := errors.WrapSentinel(errors.ErrUnauthorized, "token", "expired_token")

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, errors.ErrUnauthorized))
		assert.Contains(t, wrappedErr.Error(), "token expired_token")
	})

	t.Run("wrap ErrValidation", func(t *testing.T) {
		wrappedErr := errors.WrapSentinel(errors.ErrValidation, "field", "email")

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, errors.ErrValidation))
		assert.Contains(t, wrappedErr.Error(), "field email")
	})

	t.Run("wrap with struct id", func(t *testing.T) {
		type ID struct {
			Namespace string
			Name      string
		}
		id := ID{Namespace: "default", Name: "my-resource"}
		wrappedErr := errors.WrapSentinel(errors.ErrNotFound, "resource", id)

		require.NotNil(t, wrappedErr)
		assert.True(t, errors.Is(wrappedErr, errors.ErrNotFound))
		assert.Contains(t, wrappedErr.Error(), "resource")
		assert.Contains(t, wrappedErr.Error(), "default")
		assert.Contains(t, wrappedErr.Error(), "my-resource")
	})

	t.Run("all sentinel errors work", func(t *testing.T) {
		sentinels := []error{
			errors.ErrNotFound,
			errors.ErrUnauthorized,
			errors.ErrForbidden,
			errors.ErrConflict,
			errors.ErrValidation,
			errors.ErrInternal,
			errors.ErrBadRequest,
			errors.ErrUnavailable,
			errors.ErrTimeout,
		}

		for _, sentinel := range sentinels {
			wrapped := errors.WrapSentinel(sentinel, "resource", "test-id")
			assert.True(t, errors.Is(wrapped, sentinel),
				"WrapSentinel should preserve %v", sentinel)
		}
	})
}

func TestFormatError(t *testing.T) {
	t.Run("format nil error returns empty string", func(t *testing.T) {
		result := errors.FormatError(nil)
		assert.Equal(t, "", result)
	})

	t.Run("format simple error", func(t *testing.T) {
		err := errors.New("simple error")
		result := errors.FormatError(err)

		assert.Contains(t, result, "simple error")
	})

	t.Run("format wrapped error includes context", func(t *testing.T) {
		baseErr := errors.New("base")
		wrappedErr := errors.Wrap(baseErr, "context")

		result := errors.FormatError(wrappedErr)
		assert.Contains(t, result, "context")
		assert.Contains(t, result, "base")
	})

	t.Run("format error includes stack trace", func(t *testing.T) {
		err := errors.New("error with stack")
		result := errors.FormatError(err)

		// go-faster/errors includes stack trace with %+v
		// The result should contain file/line info
		assert.True(t, len(result) > len("error with stack"),
			"formatted error should include stack trace info")
	})

	t.Run("format sentinel error", func(t *testing.T) {
		result := errors.FormatError(errors.ErrNotFound)
		assert.Contains(t, result, "not found")
	})

	t.Run("format deeply nested error", func(t *testing.T) {
		baseErr := errors.ErrInternal
		err1 := errors.Wrap(baseErr, "database layer")
		err2 := errors.Wrap(err1, "service layer")
		err3 := errors.Wrap(err2, "handler layer")

		result := errors.FormatError(err3)
		assert.Contains(t, result, "handler layer")
		assert.Contains(t, result, "service layer")
		assert.Contains(t, result, "database layer")
		assert.Contains(t, result, "internal server error")
	})

	t.Run("format Wrapf error includes format args", func(t *testing.T) {
		baseErr := errors.New("connection failed")
		wrappedErr := errors.Wrapf(baseErr, "host=%s port=%d", "localhost", 5432)

		result := errors.FormatError(wrappedErr)
		assert.Contains(t, result, "host=localhost")
		assert.Contains(t, result, "port=5432")
		assert.Contains(t, result, "connection failed")
	})
}

func TestWrapChaining(t *testing.T) {
	t.Run("chain Wrap and Wrapf", func(t *testing.T) {
		baseErr := errors.ErrNotFound
		err1 := errors.Wrap(baseErr, "user lookup failed")
		err2 := errors.Wrapf(err1, "getting user %s", "john")
		err3 := errors.WithStack(err2)

		assert.True(t, errors.Is(err3, errors.ErrNotFound))
		assert.Contains(t, err3.Error(), "getting user john")
		assert.Contains(t, err3.Error(), "user lookup failed")
	})

	t.Run("WrapSentinel then Wrapf", func(t *testing.T) {
		err1 := errors.WrapSentinel(errors.ErrForbidden, "resource", "secret-doc")
		err2 := errors.Wrapf(err1, "user %s attempted access", "alice")

		assert.True(t, errors.Is(err2, errors.ErrForbidden))
		assert.Contains(t, err2.Error(), "resource secret-doc")
		assert.Contains(t, err2.Error(), "user alice attempted access")
	})
}

func TestConcurrentErrorCreation(t *testing.T) {
	// Test that error creation is safe for concurrent use
	done := make(chan bool, 100)

	for i := 0; i < 100; i++ {
		go func(n int) {
			err := errors.Wrapf(errors.ErrNotFound, "item %d", n)
			assert.True(t, errors.Is(err, errors.ErrNotFound))
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestErrorMessageFormat(t *testing.T) {
	t.Run("Wrap message format", func(t *testing.T) {
		baseErr := errors.New("base")
		wrapped := errors.Wrap(baseErr, "context")

		// go-faster/errors format is "context: base"
		msg := wrapped.Error()
		assert.True(t, strings.HasPrefix(msg, "context"),
			"wrapped error should start with context, got: %s", msg)
	})

	t.Run("Wrapf message format", func(t *testing.T) {
		baseErr := errors.New("base")
		wrapped := errors.Wrapf(baseErr, "context %d", 42)

		msg := wrapped.Error()
		assert.True(t, strings.HasPrefix(msg, "context 42"),
			"wrapped error should start with formatted context, got: %s", msg)
	})

	t.Run("WrapSentinel message format", func(t *testing.T) {
		wrapped := errors.WrapSentinel(errors.ErrNotFound, "user", "123")

		msg := wrapped.Error()
		assert.True(t, strings.HasPrefix(msg, "user 123"),
			"wrapped sentinel should start with resource info, got: %s", msg)
	})
}
