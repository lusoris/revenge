package moviejobs

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	infrasearch "github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/service/search"
)

// =============================================================================
// Test Mock Repository
// =============================================================================

// mockMovieRepo implements movie.Repository for testing.
// Only methods used by search_index.go are implemented; all others panic.
type mockMovieRepo struct {
	movie.Repository // embed to satisfy interface, unused methods panic

	getMovieFunc            func(ctx context.Context, id uuid.UUID) (*movie.Movie, error)
	listMovieGenresFunc     func(ctx context.Context, movieID uuid.UUID) ([]movie.MovieGenre, error)
	listMovieCastFunc       func(ctx context.Context, movieID uuid.UUID) ([]movie.MovieCredit, error)
	listMovieCrewFunc       func(ctx context.Context, movieID uuid.UUID) ([]movie.MovieCredit, error)
	listMovieFilesByMovieID func(ctx context.Context, movieID uuid.UUID) ([]movie.MovieFile, error)
	listMoviesFunc          func(ctx context.Context, filters movie.ListFilters) ([]movie.Movie, error)
	countMoviesFunc         func(ctx context.Context) (int64, error)
}

func (m *mockMovieRepo) GetMovie(ctx context.Context, id uuid.UUID) (*movie.Movie, error) {
	if m.getMovieFunc != nil {
		return m.getMovieFunc(ctx, id)
	}
	return nil, movie.ErrMovieNotFound
}

func (m *mockMovieRepo) ListMovieGenres(ctx context.Context, movieID uuid.UUID) ([]movie.MovieGenre, error) {
	if m.listMovieGenresFunc != nil {
		return m.listMovieGenresFunc(ctx, movieID)
	}
	return nil, nil
}

func (m *mockMovieRepo) ListMovieCast(ctx context.Context, movieID uuid.UUID) ([]movie.MovieCredit, error) {
	if m.listMovieCastFunc != nil {
		return m.listMovieCastFunc(ctx, movieID)
	}
	return nil, nil
}

func (m *mockMovieRepo) ListMovieCrew(ctx context.Context, movieID uuid.UUID) ([]movie.MovieCredit, error) {
	if m.listMovieCrewFunc != nil {
		return m.listMovieCrewFunc(ctx, movieID)
	}
	return nil, nil
}

func (m *mockMovieRepo) ListMovieFilesByMovieID(ctx context.Context, movieID uuid.UUID) ([]movie.MovieFile, error) {
	if m.listMovieFilesByMovieID != nil {
		return m.listMovieFilesByMovieID(ctx, movieID)
	}
	return nil, nil
}

func (m *mockMovieRepo) ListMovies(ctx context.Context, filters movie.ListFilters) ([]movie.Movie, error) {
	if m.listMoviesFunc != nil {
		return m.listMoviesFunc(ctx, filters)
	}
	return nil, nil
}

func (m *mockMovieRepo) CountMovies(ctx context.Context) (int64, error) {
	if m.countMoviesFunc != nil {
		return m.countMoviesFunc(ctx)
	}
	return 0, nil
}

// newEnabledSearchService creates a MovieSearchService with an enabled search client
// for testing. The client points at a non-existent Typesense server, so actual
// network calls will fail, but IsEnabled() returns true.
func newEnabledSearchService(t *testing.T) *search.MovieSearchService {
	t.Helper()

	cfg := &config.Config{
		Search: config.SearchConfig{
			URL:     "http://localhost:59999",
			APIKey:  "test-key",
			Enabled: true,
		},
	}

	client, err := infrasearch.NewClient(cfg, slog.Default())
	require.NoError(t, err)
	require.True(t, client.IsEnabled())

	return search.NewMovieSearchService(client, slog.Default())
}

// =============================================================================
// SearchIndexOperation Constants Tests
// =============================================================================

func TestSearchIndexOperationConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, SearchIndexOperation("index"), SearchIndexOperationIndex)
	assert.Equal(t, SearchIndexOperation("remove"), SearchIndexOperationRemove)
	assert.Equal(t, SearchIndexOperation("reindex"), SearchIndexOperationReindex)
}

func TestSearchIndexOperationConstants_AreUnique(t *testing.T) {
	t.Parallel()

	ops := []SearchIndexOperation{
		SearchIndexOperationIndex,
		SearchIndexOperationRemove,
		SearchIndexOperationReindex,
	}

	seen := make(map[SearchIndexOperation]bool)
	for _, op := range ops {
		assert.False(t, seen[op], "duplicate operation: %s", op)
		seen[op] = true
		assert.NotEmpty(t, string(op))
	}
}

// =============================================================================
// MovieSearchIndexArgs Tests
// =============================================================================

func TestMovieSearchIndexArgs_Kind(t *testing.T) {
	t.Parallel()

	args := MovieSearchIndexArgs{}
	assert.Equal(t, "movie_search_index", args.Kind())
}

func TestMovieSearchIndexArgs_Fields(t *testing.T) {
	t.Parallel()

	movieID := uuid.Must(uuid.NewV7())

	tests := []struct {
		name      string
		operation SearchIndexOperation
		movieID   uuid.UUID
	}{
		{
			name:      "index operation",
			operation: SearchIndexOperationIndex,
			movieID:   movieID,
		},
		{
			name:      "remove operation",
			operation: SearchIndexOperationRemove,
			movieID:   movieID,
		},
		{
			name:      "reindex operation",
			operation: SearchIndexOperationReindex,
			movieID:   uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			args := MovieSearchIndexArgs{
				Operation: tt.operation,
				MovieID:   tt.movieID,
			}

			assert.Equal(t, tt.operation, args.Operation)
			assert.Equal(t, tt.movieID, args.MovieID)
			assert.Equal(t, "movie_search_index", args.Kind())
		})
	}
}

func TestMovieSearchIndexArgs_ZeroValue(t *testing.T) {
	t.Parallel()

	args := MovieSearchIndexArgs{}
	assert.Equal(t, SearchIndexOperation(""), args.Operation)
	assert.Equal(t, uuid.Nil, args.MovieID)
	assert.Equal(t, "movie_search_index", args.Kind())
}

// =============================================================================
// MovieSearchIndexWorker Constructor Tests
// =============================================================================

func TestNewMovieSearchIndexWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieSearchIndexWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.movieRepo)
	assert.Nil(t, worker.searchService)
	assert.NotNil(t, worker.logger)
}

func TestNewMovieSearchIndexWorker_LoggerNamed(t *testing.T) {
	t.Parallel()

	// The constructor calls logger.Named("search_index_worker"), verify it doesn't panic.
	logger := zap.NewNop()
	worker := NewMovieSearchIndexWorker(nil, nil, logger)
	assert.NotNil(t, worker.logger)
}

// =============================================================================
// MovieSearchIndexWorker.Timeout() Tests
// =============================================================================

func TestMovieSearchIndexWorker_Timeout(t *testing.T) {
	t.Parallel()

	worker := NewMovieSearchIndexWorker(nil, nil, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationIndex,
			MovieID:   uuid.Must(uuid.NewV7()),
		},
	}

	assert.Equal(t, 15*time.Minute, worker.Timeout(job))
}

func TestMovieSearchIndexWorker_Timeout_ConsistentForAllOperations(t *testing.T) {
	t.Parallel()

	worker := NewMovieSearchIndexWorker(nil, nil, zap.NewNop())

	operations := []SearchIndexOperation{
		SearchIndexOperationIndex,
		SearchIndexOperationRemove,
		SearchIndexOperationReindex,
	}

	for _, op := range operations {
		job := &river.Job[MovieSearchIndexArgs]{
			JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
			Args: MovieSearchIndexArgs{
				Operation: op,
				MovieID:   uuid.Must(uuid.NewV7()),
			},
		}
		assert.Equal(t, 15*time.Minute, worker.Timeout(job), "timeout should be 15m for operation: %s", op)
	}
}

// =============================================================================
// MovieSearchIndexWorker.Work() Tests - Disabled Search
// =============================================================================

func TestMovieSearchIndexWorker_Work_SearchDisabled(t *testing.T) {
	t.Parallel()

	// A zero-value MovieSearchService with nil client returns IsEnabled() == false.
	searchSvc := &search.MovieSearchService{}
	worker := NewMovieSearchIndexWorker(nil, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationIndex,
			MovieID:   uuid.Must(uuid.NewV7()),
		},
	}

	// When search is disabled, work should return nil immediately.
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestMovieSearchIndexWorker_Work_SearchDisabled_Reindex(t *testing.T) {
	t.Parallel()

	searchSvc := &search.MovieSearchService{}
	worker := NewMovieSearchIndexWorker(nil, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationReindex,
		},
	}

	// Reindex should also return nil when search is disabled.
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestMovieSearchIndexWorker_Work_SearchDisabled_Remove(t *testing.T) {
	t.Parallel()

	searchSvc := &search.MovieSearchService{}
	worker := NewMovieSearchIndexWorker(nil, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationRemove,
			MovieID:   uuid.Must(uuid.NewV7()),
		},
	}

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestMovieSearchIndexWorker_Work_NilSearchService(t *testing.T) {
	t.Parallel()

	// Nil searchService should panic when accessing IsEnabled().
	worker := NewMovieSearchIndexWorker(nil, nil, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationIndex,
			MovieID:   uuid.Must(uuid.NewV7()),
		},
	}

	assert.Panics(t, func() {
		_ = worker.Work(context.Background(), job)
	})
}

// =============================================================================
// MovieSearchIndexWorker.Work() Tests - Enabled Search
// =============================================================================

func TestMovieSearchIndexWorker_Work_UnknownOperation(t *testing.T) {
	t.Parallel()

	searchSvc := newEnabledSearchService(t)
	worker := NewMovieSearchIndexWorker(nil, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperation("bogus"),
			MovieID:   uuid.Must(uuid.NewV7()),
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown operation")
	assert.Contains(t, err.Error(), "bogus")
}

func TestMovieSearchIndexWorker_Work_IndexMovieNotFound(t *testing.T) {
	t.Parallel()

	searchSvc := newEnabledSearchService(t)
	repo := &mockMovieRepo{
		getMovieFunc: func(ctx context.Context, id uuid.UUID) (*movie.Movie, error) {
			return nil, movie.ErrMovieNotFound
		},
	}
	worker := NewMovieSearchIndexWorker(repo, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationIndex,
			MovieID:   uuid.Must(uuid.NewV7()),
		},
	}

	// When the movie is not found, indexMovie returns nil (job succeeds, movie was deleted).
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestMovieSearchIndexWorker_Work_IndexMovieRepoError(t *testing.T) {
	t.Parallel()

	searchSvc := newEnabledSearchService(t)
	repo := &mockMovieRepo{
		getMovieFunc: func(ctx context.Context, id uuid.UUID) (*movie.Movie, error) {
			return nil, assert.AnError
		},
	}
	worker := NewMovieSearchIndexWorker(repo, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationIndex,
			MovieID:   uuid.Must(uuid.NewV7()),
		},
	}

	// When repo returns a non-ErrMovieNotFound error, the job should fail.
	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get movie")
}

func TestMovieSearchIndexWorker_Work_IndexMovieSuccess_SearchFails(t *testing.T) {
	t.Parallel()

	searchSvc := newEnabledSearchService(t)
	movieID := uuid.Must(uuid.NewV7())
	repo := &mockMovieRepo{
		getMovieFunc: func(ctx context.Context, id uuid.UUID) (*movie.Movie, error) {
			return &movie.Movie{ID: movieID, Title: "Test Movie"}, nil
		},
		listMovieGenresFunc: func(ctx context.Context, id uuid.UUID) ([]movie.MovieGenre, error) {
			return []movie.MovieGenre{{Name: "Action"}}, nil
		},
		listMovieCastFunc: func(ctx context.Context, id uuid.UUID) ([]movie.MovieCredit, error) {
			return []movie.MovieCredit{{Name: "Actor 1"}}, nil
		},
		listMovieCrewFunc: func(ctx context.Context, id uuid.UUID) ([]movie.MovieCredit, error) {
			return []movie.MovieCredit{{Name: "Director 1"}}, nil
		},
		listMovieFilesByMovieID: func(ctx context.Context, id uuid.UUID) ([]movie.MovieFile, error) {
			return []movie.MovieFile{{FilePath: "/movies/test.mkv"}}, nil
		},
	}
	worker := NewMovieSearchIndexWorker(repo, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationIndex,
			MovieID:   movieID,
		},
	}

	// The search service will fail trying to connect to Typesense,
	// but all code paths in indexMovie are exercised.
	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to index movie")
}

func TestMovieSearchIndexWorker_Work_IndexMovie_GenreError(t *testing.T) {
	t.Parallel()

	searchSvc := newEnabledSearchService(t)
	movieID := uuid.Must(uuid.NewV7())
	repo := &mockMovieRepo{
		getMovieFunc: func(ctx context.Context, id uuid.UUID) (*movie.Movie, error) {
			return &movie.Movie{ID: movieID, Title: "Test Movie"}, nil
		},
		listMovieGenresFunc: func(ctx context.Context, id uuid.UUID) ([]movie.MovieGenre, error) {
			return nil, assert.AnError
		},
		listMovieCastFunc: func(ctx context.Context, id uuid.UUID) ([]movie.MovieCredit, error) {
			return nil, assert.AnError
		},
		listMovieCrewFunc: func(ctx context.Context, id uuid.UUID) ([]movie.MovieCredit, error) {
			return nil, assert.AnError
		},
		listMovieFilesByMovieID: func(ctx context.Context, id uuid.UUID) ([]movie.MovieFile, error) {
			return nil, assert.AnError
		},
	}
	worker := NewMovieSearchIndexWorker(repo, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationIndex,
			MovieID:   movieID,
		},
	}

	// Errors from genres/cast/crew/files are logged as warnings, not fatal.
	// The indexing still proceeds and will fail at Typesense.
	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to index movie")
}

func TestMovieSearchIndexWorker_Work_IndexMovie_NoFiles(t *testing.T) {
	t.Parallel()

	searchSvc := newEnabledSearchService(t)
	movieID := uuid.Must(uuid.NewV7())
	repo := &mockMovieRepo{
		getMovieFunc: func(ctx context.Context, id uuid.UUID) (*movie.Movie, error) {
			return &movie.Movie{ID: movieID, Title: "Test Movie"}, nil
		},
		// All other funcs return nil/empty by default.
	}
	worker := NewMovieSearchIndexWorker(repo, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationIndex,
			MovieID:   movieID,
		},
	}

	// Exercises the path where no files are found (file == nil).
	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to index movie")
}

func TestMovieSearchIndexWorker_Work_RemoveMovie(t *testing.T) {
	t.Parallel()

	searchSvc := newEnabledSearchService(t)
	worker := NewMovieSearchIndexWorker(nil, searchSvc, zap.NewNop())

	movieID := uuid.Must(uuid.NewV7())
	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationRemove,
			MovieID:   movieID,
		},
	}

	// RemoveMovie will fail trying to connect to Typesense,
	// but the code path through removeMovie is exercised.
	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to remove movie from index")
}

func TestMovieSearchIndexWorker_Work_Reindex(t *testing.T) {
	t.Parallel()

	searchSvc := newEnabledSearchService(t)
	worker := NewMovieSearchIndexWorker(nil, searchSvc, zap.NewNop())

	job := &river.Job[MovieSearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: "movie_search_index"},
		Args: MovieSearchIndexArgs{
			Operation: SearchIndexOperationReindex,
		},
	}

	// ReindexAll will fail trying to connect to Typesense,
	// but the code path through reindexAll is exercised.
	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to reindex all movies")
}
