package moviejobs

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMovieMetadataRefreshArgs_Kind(t *testing.T) {
	t.Parallel()

	args := MovieMetadataRefreshArgs{}
	assert.Equal(t, "movie_metadata_refresh", args.Kind())
}

func TestMovieMetadataRefreshArgs_Fields(t *testing.T) {
	t.Parallel()

	movieID := uuid.New()
	args := MovieMetadataRefreshArgs{
		MovieID: movieID,
		Force:   true,
	}

	assert.Equal(t, movieID, args.MovieID)
	assert.True(t, args.Force)
}

func TestNewMovieMetadataRefreshWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieMetadataRefreshWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.movieRepo)
	assert.Nil(t, worker.metadataService)
	assert.NotNil(t, worker.logger)
}

func TestFormatTimePtr(t *testing.T) {
	t.Parallel()

	t.Run("nil input", func(t *testing.T) {
		result := formatTimePtr(nil)
		assert.Nil(t, result)
	})

	t.Run("valid time", func(t *testing.T) {
		tm := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)
		result := formatTimePtr(&tm)

		assert.NotNil(t, result)
		assert.Equal(t, "2024-05-15", *result)
	})
}

func TestFormatDecimalPtr(t *testing.T) {
	t.Parallel()

	t.Run("nil input", func(t *testing.T) {
		result := formatDecimalPtr(nil)
		assert.Nil(t, result)
	})

	t.Run("zero decimal", func(t *testing.T) {
		d := decimal.Zero
		result := formatDecimalPtr(&d)
		assert.Nil(t, result)
	})

	t.Run("valid decimal", func(t *testing.T) {
		d := decimal.NewFromFloat(7.5)
		result := formatDecimalPtr(&d)

		assert.NotNil(t, result)
		assert.Equal(t, "7.5", *result)
	})
}
