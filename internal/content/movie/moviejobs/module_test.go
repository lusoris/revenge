package moviejobs

import (
	"testing"

	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterWorkers(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	workers := river.NewWorkers()

	metadataRefreshWorker := NewMovieMetadataRefreshWorker(nil, nil, logger)
	libraryScanWorker := NewMovieLibraryScanWorker(nil, logger)
	fileMatchWorker := NewMovieFileMatchWorker(nil, logger)
	searchIndexWorker := NewMovieSearchIndexWorker(nil, nil, logger)

	err := RegisterWorkers(workers, metadataRefreshWorker, libraryScanWorker, fileMatchWorker, searchIndexWorker)
	require.NoError(t, err)
}

func TestRegisterWorkers_ReturnsNil(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	workers := river.NewWorkers()

	metadataRefreshWorker := NewMovieMetadataRefreshWorker(nil, nil, logger)
	libraryScanWorker := NewMovieLibraryScanWorker(nil, logger)
	fileMatchWorker := NewMovieFileMatchWorker(nil, logger)
	searchIndexWorker := NewMovieSearchIndexWorker(nil, nil, logger)

	// RegisterWorkers always returns nil.
	err := RegisterWorkers(workers, metadataRefreshWorker, libraryScanWorker, fileMatchWorker, searchIndexWorker)
	assert.NoError(t, err)
}

func TestModule_IsNotNil(t *testing.T) {
	t.Parallel()

	// Verify the fx.Module is defined and not nil.
	assert.NotNil(t, Module)
}

func TestRegisterWorkersParams_Fields(t *testing.T) {
	t.Parallel()

	// Verify RegisterWorkersParams struct can be constructed with nil fields.
	params := RegisterWorkersParams{}
	assert.Nil(t, params.Workers)
	assert.Nil(t, params.MetadataRefreshWorker)
	assert.Nil(t, params.LibraryScanWorker)
	assert.Nil(t, params.FileMatchWorker)
	assert.Nil(t, params.SearchIndexWorker)
	assert.Nil(t, params.MovieService)
	assert.Nil(t, params.LibraryService)
	assert.Nil(t, params.SearchService)
	assert.Nil(t, params.Logger)
}

// =============================================================================
// Job Kind Uniqueness Test
// =============================================================================

func TestAllJobKinds_AreUnique(t *testing.T) {
	t.Parallel()

	kinds := []string{
		MovieLibraryScanJobKind,
		MovieFileMatchJobKind,
		"movie_search_index", // From MovieSearchIndexArgs.Kind()
		"metadata_refresh_movie", // From metadatajobs.RefreshMovieArgs.Kind()
	}

	seen := make(map[string]bool)
	for _, kind := range kinds {
		assert.False(t, seen[kind], "duplicate job kind: %s", kind)
		seen[kind] = true
		assert.NotEmpty(t, kind)
	}
}

func TestAllJobKinds_HaveExpectedPrefix(t *testing.T) {
	t.Parallel()

	// Movie-specific job kinds should contain "movie" in their name.
	assert.Contains(t, MovieLibraryScanJobKind, "movie")
	assert.Contains(t, MovieFileMatchJobKind, "movie")
	assert.Contains(t, MovieSearchIndexArgs{}.Kind(), "movie")
}
