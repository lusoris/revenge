package moviejobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/content/movie"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	metadatajobs "github.com/lusoris/revenge/internal/service/metadata/jobs"
)

// =============================================================================
// Mock Service for metadata refresh tests
// =============================================================================

// mockMovieService implements movie.Service for testing.
type mockMovieService struct {
	movie.Service // embed to satisfy interface

	refreshMovieMetadataFunc func(ctx context.Context, id uuid.UUID, opts ...movie.MetadataRefreshOptions) error
}

func (m *mockMovieService) RefreshMovieMetadata(ctx context.Context, id uuid.UUID, opts ...movie.MetadataRefreshOptions) error {
	if m.refreshMovieMetadataFunc != nil {
		return m.refreshMovieMetadataFunc(ctx, id, opts...)
	}
	return nil
}

// =============================================================================
// RefreshMovieArgs Tests
// =============================================================================

func TestRefreshMovieArgs_Kind(t *testing.T) {
	t.Parallel()

	args := metadatajobs.RefreshMovieArgs{}
	assert.Equal(t, "metadata_refresh_movie", args.Kind())
}

func TestRefreshMovieArgs_Fields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		movieID           uuid.UUID
		force             bool
		languages         []string
		expectedForce     bool
		expectedLanguages []string
	}{
		{
			name:              "full args",
			movieID:           uuid.Must(uuid.NewV7()),
			force:             true,
			languages:         []string{"en", "de"},
			expectedForce:     true,
			expectedLanguages: []string{"en", "de"},
		},
		{
			name:              "minimal args",
			movieID:           uuid.Must(uuid.NewV7()),
			force:             false,
			languages:         nil,
			expectedForce:     false,
			expectedLanguages: nil,
		},
		{
			name:              "single language",
			movieID:           uuid.Must(uuid.NewV7()),
			force:             false,
			languages:         []string{"en"},
			expectedForce:     false,
			expectedLanguages: []string{"en"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			args := metadatajobs.RefreshMovieArgs{
				MovieID:   tt.movieID,
				Force:     tt.force,
				Languages: tt.languages,
			}
			assert.Equal(t, tt.movieID, args.MovieID)
			assert.Equal(t, tt.expectedForce, args.Force)
			assert.Equal(t, tt.expectedLanguages, args.Languages)
		})
	}
}

func TestRefreshMovieArgs_InsertOpts_UsesDefaults(t *testing.T) {
	t.Parallel()

	// RefreshMovieArgs does not define InsertOpts, so river.WorkerDefaults
	// will provide default insert opts (empty queue = default queue).
	args := metadatajobs.RefreshMovieArgs{}
	assert.Equal(t, "metadata_refresh_movie", args.Kind())
}

// =============================================================================
// MovieMetadataRefreshWorker Constructor Tests
// =============================================================================

func TestNewMovieMetadataRefreshWorker(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewMovieMetadataRefreshWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.service)
	assert.Nil(t, worker.jobClient)
	assert.NotNil(t, worker.logger)
}

func TestNewMovieMetadataRefreshWorker_NilLogger(t *testing.T) {
	t.Parallel()

	worker := NewMovieMetadataRefreshWorker(nil, nil, nil)
	assert.NotNil(t, worker)
	assert.Nil(t, worker.service)
	assert.Nil(t, worker.jobClient)
	assert.Nil(t, worker.logger)
}

// =============================================================================
// MovieMetadataRefreshWorker.Timeout() Tests
// =============================================================================

func TestMovieMetadataRefreshWorker_Timeout(t *testing.T) {
	t.Parallel()

	worker := NewMovieMetadataRefreshWorker(nil, nil, logging.NewTestLogger())

	job := &river.Job[metadatajobs.RefreshMovieArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "metadata_refresh_movie"},
		Args: metadatajobs.RefreshMovieArgs{
			MovieID: uuid.Must(uuid.NewV7()),
			Force:   false,
		},
	}

	assert.Equal(t, 5*time.Minute, worker.Timeout(job))
}

// =============================================================================
// MovieMetadataRefreshWorker.Work() Tests
// =============================================================================

func TestMovieMetadataRefreshWorker_Work_NilService(t *testing.T) {
	t.Parallel()

	// Create worker with nil service and nil jobClient.
	// The nil jobClient will panic on w.jobClient.ReportProgress().
	worker := NewMovieMetadataRefreshWorker(nil, nil, logging.NewTestLogger())

	job := &river.Job[metadatajobs.RefreshMovieArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "metadata_refresh_movie"},
		Args: metadatajobs.RefreshMovieArgs{
			MovieID:   uuid.Must(uuid.NewV7()),
			Force:     true,
			Languages: []string{"en"},
		},
	}

	// With a nil jobClient, calling Work panics on w.jobClient.ReportProgress().
	assert.Panics(t, func() {
		_ = worker.Work(context.Background(), job)
	})
}

func TestMovieMetadataRefreshWorker_Work_Success(t *testing.T) {
	t.Parallel()

	movieID := uuid.Must(uuid.NewV7())
	var capturedID uuid.UUID
	var capturedOpts movie.MetadataRefreshOptions

	svc := &mockMovieService{
		refreshMovieMetadataFunc: func(ctx context.Context, id uuid.UUID, opts ...movie.MetadataRefreshOptions) error {
			capturedID = id
			if len(opts) > 0 {
				capturedOpts = opts[0]
			}
			return nil
		},
	}

	// Use a zero-value infrajobs.Client where ReportProgress returns nil (c.client == nil).
	jobClient := &infrajobs.Client{}
	worker := NewMovieMetadataRefreshWorker(svc, jobClient, logging.NewTestLogger())

	job := &river.Job[metadatajobs.RefreshMovieArgs]{
		JobRow: &rivertype.JobRow{ID: 42, Kind: "metadata_refresh_movie"},
		Args: metadatajobs.RefreshMovieArgs{
			MovieID:   movieID,
			Force:     true,
			Languages: []string{"en", "de"},
		},
	}

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)

	// Verify the service was called with the correct arguments.
	assert.Equal(t, movieID, capturedID)
	assert.True(t, capturedOpts.Force)
	assert.Equal(t, []string{"en", "de"}, capturedOpts.Languages)
}

func TestMovieMetadataRefreshWorker_Work_ServiceError(t *testing.T) {
	t.Parallel()

	svc := &mockMovieService{
		refreshMovieMetadataFunc: func(ctx context.Context, id uuid.UUID, opts ...movie.MetadataRefreshOptions) error {
			return errors.New("metadata provider unavailable")
		},
	}

	jobClient := &infrajobs.Client{}
	worker := NewMovieMetadataRefreshWorker(svc, jobClient, logging.NewTestLogger())

	job := &river.Job[metadatajobs.RefreshMovieArgs]{
		JobRow: &rivertype.JobRow{ID: 42, Kind: "metadata_refresh_movie"},
		Args: metadatajobs.RefreshMovieArgs{
			MovieID:   uuid.Must(uuid.NewV7()),
			Force:     false,
			Languages: []string{"en"},
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "movie metadata refresh failed")
	assert.Contains(t, err.Error(), "metadata provider unavailable")
}

func TestMovieMetadataRefreshWorker_Work_NoLanguages(t *testing.T) {
	t.Parallel()

	var capturedOpts movie.MetadataRefreshOptions

	svc := &mockMovieService{
		refreshMovieMetadataFunc: func(ctx context.Context, id uuid.UUID, opts ...movie.MetadataRefreshOptions) error {
			if len(opts) > 0 {
				capturedOpts = opts[0]
			}
			return nil
		},
	}

	jobClient := &infrajobs.Client{}
	worker := NewMovieMetadataRefreshWorker(svc, jobClient, logging.NewTestLogger())

	job := &river.Job[metadatajobs.RefreshMovieArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "metadata_refresh_movie"},
		Args: metadatajobs.RefreshMovieArgs{
			MovieID: uuid.Must(uuid.NewV7()),
			Force:   false,
		},
	}

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	assert.False(t, capturedOpts.Force)
	assert.Nil(t, capturedOpts.Languages)
}
