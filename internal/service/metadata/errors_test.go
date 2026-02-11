package metadata

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderError_Error(t *testing.T) {
	t.Run("with wrapped error", func(t *testing.T) {
		pe := &ProviderError{
			Provider:   ProviderTMDb,
			StatusCode: 404,
			Message:    "not found",
			Err:        ErrNotFound,
		}
		result := pe.Error()
		assert.Contains(t, result, string(ProviderTMDb))
		assert.Contains(t, result, "404")
		assert.Contains(t, result, "not found")
		assert.Contains(t, result, ErrNotFound.Error())
	})

	t.Run("without wrapped error", func(t *testing.T) {
		pe := &ProviderError{
			Provider:   ProviderTMDb,
			StatusCode: 500,
			Message:    "server error",
		}
		result := pe.Error()
		assert.Contains(t, result, "500")
		assert.Contains(t, result, "server error")
		assert.NotContains(t, result, "<nil>")
	})
}

func TestProviderError_Unwrap(t *testing.T) {
	inner := errors.New("inner error")
	pe := &ProviderError{Err: inner}
	assert.Equal(t, inner, pe.Unwrap())
	assert.True(t, errors.Is(pe, inner))

	pe2 := &ProviderError{}
	assert.Nil(t, pe2.Unwrap())
}

func TestIsProviderError(t *testing.T) {
	pe := NewProviderError(ProviderTMDb, 500, "fail", nil)
	assert.True(t, IsProviderError(pe))

	wrapped := fmt.Errorf("wrap: %w", pe)
	assert.True(t, IsProviderError(wrapped))

	assert.False(t, IsProviderError(errors.New("plain")))
	assert.False(t, IsProviderError(nil))
}

func TestNewProviderError(t *testing.T) {
	inner := errors.New("inner")
	pe := NewProviderError(ProviderOMDb, 503, "unavailable", inner)
	assert.Equal(t, ProviderOMDb, pe.Provider)
	assert.Equal(t, 503, pe.StatusCode)
	assert.Equal(t, "unavailable", pe.Message)
	assert.Equal(t, inner, pe.Err)
}

func TestNewNotFoundError(t *testing.T) {
	pe := NewNotFoundError(ProviderTVDb, "movie", "12345")
	assert.Equal(t, ProviderTVDb, pe.Provider)
	assert.Equal(t, 404, pe.StatusCode)
	assert.Contains(t, pe.Message, "movie")
	assert.Contains(t, pe.Message, "12345")
	assert.True(t, errors.Is(pe, ErrNotFound))
}

func TestNewRateLimitError(t *testing.T) {
	t.Run("with retry", func(t *testing.T) {
		pe := NewRateLimitError(ProviderTMDb, 30)
		assert.Equal(t, 429, pe.StatusCode)
		assert.Contains(t, pe.Message, "30")
		assert.True(t, errors.Is(pe, ErrRateLimited))
	})

	t.Run("without retry", func(t *testing.T) {
		pe := NewRateLimitError(ProviderTMDb, 0)
		assert.Equal(t, 429, pe.StatusCode)
		assert.Equal(t, "rate limit exceeded", pe.Message)
	})
}

func TestNewUnauthorizedError(t *testing.T) {
	pe := NewUnauthorizedError(ProviderOMDb)
	assert.Equal(t, ProviderOMDb, pe.Provider)
	assert.Equal(t, 401, pe.StatusCode)
	assert.Contains(t, pe.Message, "API key")
	assert.True(t, errors.Is(pe, ErrUnauthorized))
}

func TestAggregateError_Error(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		agg := &AggregateError{}
		assert.Equal(t, "metadata: no errors", agg.Error())
	})

	t.Run("single error", func(t *testing.T) {
		agg := &AggregateError{Errors: []error{errors.New("one")}}
		assert.Equal(t, "one", agg.Error())
	})

	t.Run("multiple errors", func(t *testing.T) {
		agg := &AggregateError{Errors: []error{errors.New("first"), errors.New("second")}}
		result := agg.Error()
		assert.Contains(t, result, "2 provider errors")
		assert.Contains(t, result, "first")
	})
}

func TestAggregateError_Add(t *testing.T) {
	agg := &AggregateError{}
	agg.Add(nil)
	assert.Empty(t, agg.Errors)

	agg.Add(errors.New("one"))
	assert.Len(t, agg.Errors, 1)

	agg.Add(errors.New("two"))
	assert.Len(t, agg.Errors, 2)
}

func TestAggregateError_IsEmpty(t *testing.T) {
	agg := &AggregateError{}
	assert.True(t, agg.IsEmpty())

	agg.Add(errors.New("err"))
	assert.False(t, agg.IsEmpty())
}

func TestAggregateError_First(t *testing.T) {
	agg := &AggregateError{}
	assert.Nil(t, agg.First())

	first := errors.New("first")
	agg.Add(first)
	agg.Add(errors.New("second"))
	assert.Equal(t, first, agg.First())
}

func TestAggregateError_HasNotFound(t *testing.T) {
	t.Run("empty returns false", func(t *testing.T) {
		agg := &AggregateError{}
		assert.False(t, agg.HasNotFound())
	})

	t.Run("all not found", func(t *testing.T) {
		agg := &AggregateError{}
		agg.Add(NewNotFoundError(ProviderTMDb, "movie", "1"))
		agg.Add(NewNotFoundError(ProviderOMDb, "movie", "2"))
		assert.True(t, agg.HasNotFound())
	})

	t.Run("mixed errors", func(t *testing.T) {
		agg := &AggregateError{}
		agg.Add(NewNotFoundError(ProviderTMDb, "movie", "1"))
		agg.Add(errors.New("other error"))
		assert.False(t, agg.HasNotFound())
	})
}
