package moviejobs

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	metadatajobs "github.com/lusoris/revenge/internal/service/metadata/jobs"
)

func TestRefreshMovieArgs_Kind(t *testing.T) {
	t.Parallel()

	args := metadatajobs.RefreshMovieArgs{}
	assert.Equal(t, "metadata_refresh_movie", args.Kind())
}

func TestRefreshMovieArgs_Fields(t *testing.T) {
	t.Parallel()

	movieID := uuid.Must(uuid.NewV7())
	args := metadatajobs.RefreshMovieArgs{
		MovieID:   movieID,
		Force:     true,
		Languages: []string{"en", "de"},
	}

	assert.Equal(t, movieID, args.MovieID)
	assert.True(t, args.Force)
	assert.Equal(t, []string{"en", "de"}, args.Languages)
}

func TestNewMovieMetadataRefreshWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieMetadataRefreshWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.service)
	assert.Nil(t, worker.jobClient)
	assert.NotNil(t, worker.logger)
}
