package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/ogen-go/ogen/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("generates request ID when not present", func(t *testing.T) {
		var capturedID string

		// Create middleware
		mw := RequestIDMiddleware()

		// Create request without X-Request-ID header
		httpReq, err := http.NewRequest("GET", "http://example.com", nil)
		require.NoError(t, err)

		req := middleware.Request{
			Context: context.Background(),
			Raw:     httpReq,
		}

		// Mock next handler
		next := func(req middleware.Request) (middleware.Response, error) {
			capturedID = GetRequestID(req.Context)
			return middleware.Response{}, nil
		}

		// Execute middleware
		_, err = mw(req, next)

		require.NoError(t, err)
		assert.NotEmpty(t, capturedID, "should generate request ID")
		assert.Len(t, capturedID, 36, "should be valid UUID format")
	})

	t.Run("uses existing request ID from header", func(t *testing.T) {
		existingID := "test-request-id-12345"
		var capturedID string

		// Create middleware
		mw := RequestIDMiddleware()

		// Create request with X-Request-ID header
		httpReq, err := http.NewRequest("GET", "http://example.com", nil)
		require.NoError(t, err)
		httpReq.Header.Set("X-Request-ID", existingID)

		req := middleware.Request{
			Context: context.Background(),
			Raw:     httpReq,
		}

		// Mock next handler
		next := func(req middleware.Request) (middleware.Response, error) {
			capturedID = GetRequestID(req.Context)
			return middleware.Response{}, nil
		}

		// Execute middleware
		_, err = mw(req, next)

		require.NoError(t, err)
		assert.Equal(t, existingID, capturedID, "should use existing request ID")
	})

	t.Run("stores request ID in context", func(t *testing.T) {
		// Create middleware
		mw := RequestIDMiddleware()

		// Create request
		httpReq, err := http.NewRequest("GET", "http://example.com", nil)
		require.NoError(t, err)

		req := middleware.Request{
			Context: context.Background(),
			Raw:     httpReq,
		}

		// Mock next handler that checks context
		next := func(req middleware.Request) (middleware.Response, error) {
			id := GetRequestID(req.Context)
			assert.NotEmpty(t, id, "request ID should be in context")
			return middleware.Response{}, nil
		}

		// Execute middleware
		_, err = mw(req, next)
		require.NoError(t, err)
	})
}

func TestGetRequestID(t *testing.T) {
	t.Run("returns request ID from context", func(t *testing.T) {
		testID := "test-id-123"
		ctx := WithRequestID(context.Background(), testID)

		result := GetRequestID(ctx)

		assert.Equal(t, testID, result)
	})

	t.Run("returns empty string when not found", func(t *testing.T) {
		ctx := context.Background()

		result := GetRequestID(ctx)

		assert.Empty(t, result)
	})

	t.Run("returns empty string for wrong type", func(t *testing.T) {
		// Store wrong type in context
		ctx := context.WithValue(context.Background(), requestIDKey{}, 12345)

		result := GetRequestID(ctx)

		assert.Empty(t, result, "should return empty for wrong type")
	})
}

func TestWithRequestID(t *testing.T) {
	t.Run("stores request ID in context", func(t *testing.T) {
		testID := "my-request-id"
		ctx := context.Background()

		ctx = WithRequestID(ctx, testID)

		result := GetRequestID(ctx)
		assert.Equal(t, testID, result)
	})

	t.Run("overwrites existing request ID", func(t *testing.T) {
		ctx := WithRequestID(context.Background(), "old-id")
		ctx = WithRequestID(ctx, "new-id")

		result := GetRequestID(ctx)
		assert.Equal(t, "new-id", result)
	})
}
