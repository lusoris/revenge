package movie

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService implements Service for testing handlers
type MockService struct {
	mock.Mock
}

func (m *MockService) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockService) GetMovieByTMDbID(ctx context.Context, tmdbID int32) (*Movie, error) {
	args := m.Called(ctx, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockService) GetMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error) {
	args := m.Called(ctx, imdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockService) ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockService) SearchMovies(ctx context.Context, query string, filters SearchFilters) ([]Movie, error) {
	args := m.Called(ctx, query, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockService) ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, error) {
	args := m.Called(ctx, minVotes, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockService) CreateMovie(ctx context.Context, params CreateMovieParams) (*Movie, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockService) UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockService) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) GetMovieFiles(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieFile), args.Error(1)
}

func (m *MockService) CreateMovieFile(ctx context.Context, params CreateMovieFileParams) (*MovieFile, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieFile), args.Error(1)
}

func (m *MockService) DeleteMovieFile(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) GetMovieCast(ctx context.Context, movieID uuid.UUID) ([]MovieCredit, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieCredit), args.Error(1)
}

func (m *MockService) GetMovieCrew(ctx context.Context, movieID uuid.UUID) ([]MovieCredit, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieCredit), args.Error(1)
}

func (m *MockService) GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieCollection), args.Error(1)
}

func (m *MockService) GetMoviesByCollection(ctx context.Context, collectionID uuid.UUID) ([]Movie, error) {
	args := m.Called(ctx, collectionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockService) GetCollectionForMovie(ctx context.Context, movieID uuid.UUID) (*MovieCollection, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieCollection), args.Error(1)
}

func (m *MockService) GetMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieGenre), args.Error(1)
}

func (m *MockService) GetMoviesByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]Movie, error) {
	args := m.Called(ctx, tmdbGenreID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockService) UpdateWatchProgress(ctx context.Context, userID, movieID uuid.UUID, progressSeconds, durationSeconds int32) (*MovieWatched, error) {
	args := m.Called(ctx, userID, movieID, progressSeconds, durationSeconds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieWatched), args.Error(1)
}

func (m *MockService) GetWatchProgress(ctx context.Context, userID, movieID uuid.UUID) (*MovieWatched, error) {
	args := m.Called(ctx, userID, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieWatched), args.Error(1)
}

func (m *MockService) MarkAsWatched(ctx context.Context, userID, movieID uuid.UUID) error {
	args := m.Called(ctx, userID, movieID)
	return args.Error(0)
}

func (m *MockService) RemoveWatchProgress(ctx context.Context, userID, movieID uuid.UUID) error {
	args := m.Called(ctx, userID, movieID)
	return args.Error(0)
}

func (m *MockService) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ContinueWatchingItem), args.Error(1)
}

func (m *MockService) GetWatchHistory(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedMovieItem, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]WatchedMovieItem), args.Error(1)
}

func (m *MockService) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserMovieStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserMovieStats), args.Error(1)
}

func (m *MockService) RefreshMovieMetadata(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Handler Tests

func TestHTTPError(t *testing.T) {
	err := NewHTTPError(404, "movie not found")
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "movie not found", err.Message)
	assert.Equal(t, "movie not found", err.Error())
}

func TestNotFound(t *testing.T) {
	err := NotFound("resource not found")
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "resource not found", err.Message)
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("invalid input")
	assert.Equal(t, 400, err.Code)
	assert.Equal(t, "invalid input", err.Message)
}

func TestInternalError(t *testing.T) {
	err := InternalError("something went wrong")
	assert.Equal(t, 500, err.Code)
	assert.Equal(t, "something went wrong", err.Message)
}

func TestNewHandler(t *testing.T) {
	svc := new(MockService)
	h := NewHandler(svc)
	assert.NotNil(t, h)
}

func TestHandler_GetMovie(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movie := newTestMovie()

		svc.On("GetMovie", ctx, movie.ID).Return(movie, nil)

		result, err := h.GetMovie(ctx, movie.ID.String())
		require.NoError(t, err)
		assert.Equal(t, movie.ID, result.ID)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		result, err := h.GetMovie(ctx, "invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid movie ID")
		assert.Nil(t, result)
	})

	t.Run("Not found", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		id := uuid.New()

		svc.On("GetMovie", ctx, id).Return(nil, errors.New("not found"))

		result, err := h.GetMovie(ctx, id.String())
		assert.Error(t, err)
		assert.Nil(t, result)
		svc.AssertExpectations(t)
	})
}

func TestHandler_ListMovies(t *testing.T) {
	svc := new(MockService)
	h := NewHandler(svc)
	ctx := context.Background()
	movies := []Movie{*newTestMovie(), *newTestMovie()}
	params := ListMoviesParams{
		OrderBy: "title",
		Limit:   10,
		Offset:  0,
	}
	expectedFilters := ListFilters{
		OrderBy: "title",
		Limit:   10,
		Offset:  0,
	}

	svc.On("ListMovies", ctx, expectedFilters).Return(movies, nil)

	result, err := h.ListMovies(ctx, params)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	svc.AssertExpectations(t)
}

func TestHandler_SearchMovies(t *testing.T) {
	svc := new(MockService)
	h := NewHandler(svc)
	ctx := context.Background()
	movies := []Movie{*newTestMovie()}
	params := SearchMoviesParams{
		Query:  "fight club",
		Limit:  10,
		Offset: 0,
	}
	expectedFilters := SearchFilters{
		Limit:  10,
		Offset: 0,
	}

	svc.On("SearchMovies", ctx, "fight club", expectedFilters).Return(movies, nil)

	result, err := h.SearchMovies(ctx, params)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	svc.AssertExpectations(t)
}

func TestHandler_GetRecentlyAdded(t *testing.T) {
	svc := new(MockService)
	h := NewHandler(svc)
	ctx := context.Background()
	movies := []Movie{*newTestMovie()}
	params := PaginationParams{
		Limit:  10,
		Offset: 0,
	}

	svc.On("ListRecentlyAdded", ctx, int32(10), int32(0)).Return(movies, nil)

	result, err := h.GetRecentlyAdded(ctx, params)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	svc.AssertExpectations(t)
}

func TestHandler_GetTopRated(t *testing.T) {
	t.Run("With default min votes", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movies := []Movie{*newTestMovie()}
		params := TopRatedParams{
			Limit:    10,
			Offset:   0,
			MinVotes: nil, // Default to 100
		}

		svc.On("ListTopRated", ctx, int32(100), int32(10), int32(0)).Return(movies, nil)

		result, err := h.GetTopRated(ctx, params)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		svc.AssertExpectations(t)
	})

	t.Run("With custom min votes", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movies := []Movie{*newTestMovie()}
		minVotes := int32(500)
		params := TopRatedParams{
			Limit:    10,
			Offset:   0,
			MinVotes: &minVotes,
		}

		svc.On("ListTopRated", ctx, int32(500), int32(10), int32(0)).Return(movies, nil)

		result, err := h.GetTopRated(ctx, params)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		svc.AssertExpectations(t)
	})
}

func TestHandler_GetMovieFiles(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movieID := uuid.New()
		files := []MovieFile{{ID: uuid.New(), MovieID: movieID, FilePath: "/movies/test.mkv"}}

		svc.On("GetMovieFiles", ctx, movieID).Return(files, nil)

		result, err := h.GetMovieFiles(ctx, movieID.String())
		require.NoError(t, err)
		assert.Len(t, result, 1)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		result, err := h.GetMovieFiles(ctx, "invalid")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestHandler_GetMovieCast(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movieID := uuid.New()
		cast := []MovieCredit{{ID: uuid.New(), MovieID: movieID, Name: "Brad Pitt", CreditType: "cast"}}

		svc.On("GetMovieCast", ctx, movieID).Return(cast, nil)

		result, err := h.GetMovieCast(ctx, movieID.String())
		require.NoError(t, err)
		assert.Len(t, result, 1)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		result, err := h.GetMovieCast(ctx, "not-a-uuid")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestHandler_GetMovieCrew(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movieID := uuid.New()
		crew := []MovieCredit{{ID: uuid.New(), MovieID: movieID, Name: "David Fincher", CreditType: "crew"}}

		svc.On("GetMovieCrew", ctx, movieID).Return(crew, nil)

		result, err := h.GetMovieCrew(ctx, movieID.String())
		require.NoError(t, err)
		assert.Len(t, result, 1)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		result, err := h.GetMovieCrew(ctx, "not-valid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid movie ID")
		assert.Nil(t, result)
	})
}

func TestHandler_GetMovieGenres(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movieID := uuid.New()
		genres := []MovieGenre{
			{ID: uuid.New(), MovieID: movieID, TMDbGenreID: 18, Name: "Drama"},
		}

		svc.On("GetMovieGenres", ctx, movieID).Return(genres, nil)

		result, err := h.GetMovieGenres(ctx, movieID.String())
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Drama", result[0].Name)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		result, err := h.GetMovieGenres(ctx, "bad-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid movie ID")
		assert.Nil(t, result)
	})
}

func TestHandler_GetMoviesByGenre(t *testing.T) {
	svc := new(MockService)
	h := NewHandler(svc)
	ctx := context.Background()
	movies := []Movie{*newTestMovie(), *newTestMovie()}
	params := PaginationParams{
		Limit:  20,
		Offset: 0,
	}

	svc.On("GetMoviesByGenre", ctx, int32(28), int32(20), int32(0)).Return(movies, nil)

	result, err := h.GetMoviesByGenre(ctx, 28, params)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	svc.AssertExpectations(t)
}

func TestHandler_GetMovieCollection(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movieID := uuid.New()
		collection := &MovieCollection{
			ID:   uuid.New(),
			Name: "The Matrix Collection",
		}

		svc.On("GetCollectionForMovie", ctx, movieID).Return(collection, nil)

		result, err := h.GetMovieCollection(ctx, movieID.String())
		require.NoError(t, err)
		assert.Equal(t, "The Matrix Collection", result.Name)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		result, err := h.GetMovieCollection(ctx, "invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid movie ID")
		assert.Nil(t, result)
	})
}

func TestHandler_GetCollection(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		collectionID := uuid.New()
		collection := &MovieCollection{
			ID:   collectionID,
			Name: "Star Wars Collection",
		}

		svc.On("GetMovieCollection", ctx, collectionID).Return(collection, nil)

		result, err := h.GetCollection(ctx, collectionID.String())
		require.NoError(t, err)
		assert.Equal(t, "Star Wars Collection", result.Name)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		result, err := h.GetCollection(ctx, "not-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid collection ID")
		assert.Nil(t, result)
	})
}

func TestHandler_GetCollectionMovies(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		collectionID := uuid.New()
		movies := []Movie{*newTestMovie(), *newTestMovie(), *newTestMovie()}

		svc.On("GetMoviesByCollection", ctx, collectionID).Return(movies, nil)

		result, err := h.GetCollectionMovies(ctx, collectionID.String())
		require.NoError(t, err)
		assert.Len(t, result, 3)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		result, err := h.GetCollectionMovies(ctx, "bad-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid collection ID")
		assert.Nil(t, result)
	})
}

func TestHandler_UpdateWatchProgress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		userID := uuid.New()
		movieID := uuid.New()
		watched := &MovieWatched{
			ID:              uuid.New(),
			UserID:          userID,
			MovieID:         movieID,
			ProgressSeconds: 3600,
			DurationSeconds: 7200,
		}
		params := UpdateWatchProgressParams{
			ProgressSeconds: 3600,
			DurationSeconds: 7200,
		}

		svc.On("UpdateWatchProgress", ctx, userID, movieID, int32(3600), int32(7200)).Return(watched, nil)

		result, err := h.UpdateWatchProgress(ctx, userID, movieID.String(), params)
		require.NoError(t, err)
		assert.Equal(t, int32(3600), result.ProgressSeconds)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid movie UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		userID := uuid.New()
		params := UpdateWatchProgressParams{}

		result, err := h.UpdateWatchProgress(ctx, userID, "invalid", params)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestHandler_GetWatchProgress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		userID := uuid.New()
		movieID := uuid.New()
		watched := &MovieWatched{
			ID:              uuid.New(),
			UserID:          userID,
			MovieID:         movieID,
			ProgressSeconds: 1800,
		}

		svc.On("GetWatchProgress", ctx, userID, movieID).Return(watched, nil)

		result, err := h.GetWatchProgress(ctx, userID, movieID.String())
		require.NoError(t, err)
		assert.Equal(t, int32(1800), result.ProgressSeconds)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid movie UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		userID := uuid.New()

		result, err := h.GetWatchProgress(ctx, userID, "bad-id")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestHandler_MarkAsWatched(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		userID := uuid.New()
		movieID := uuid.New()

		svc.On("MarkAsWatched", ctx, userID, movieID).Return(nil)

		err := h.MarkAsWatched(ctx, userID, movieID.String())
		require.NoError(t, err)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid movie UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		userID := uuid.New()

		err := h.MarkAsWatched(ctx, userID, "nope")
		assert.Error(t, err)
	})
}

func TestHandler_DeleteWatchProgress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		userID := uuid.New()
		movieID := uuid.New()

		svc.On("RemoveWatchProgress", ctx, userID, movieID).Return(nil)

		err := h.DeleteWatchProgress(ctx, userID, movieID.String())
		require.NoError(t, err)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid movie UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		userID := uuid.New()

		err := h.DeleteWatchProgress(ctx, userID, "xyz")
		assert.Error(t, err)
	})
}

func TestHandler_GetContinueWatching(t *testing.T) {
	svc := new(MockService)
	h := NewHandler(svc)
	ctx := context.Background()
	userID := uuid.New()
	items := []ContinueWatchingItem{
		{Movie: *newTestMovie(), ProgressSeconds: 3000},
	}

	svc.On("GetContinueWatching", ctx, userID, int32(10)).Return(items, nil)

	result, err := h.GetContinueWatching(ctx, userID, 10)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	svc.AssertExpectations(t)
}

func TestHandler_GetWatchHistory(t *testing.T) {
	svc := new(MockService)
	h := NewHandler(svc)
	ctx := context.Background()
	userID := uuid.New()
	items := []WatchedMovieItem{
		{Movie: *newTestMovie(), WatchCount: 2},
	}
	params := PaginationParams{Limit: 10, Offset: 0}

	svc.On("GetWatchHistory", ctx, userID, int32(10), int32(0)).Return(items, nil)

	result, err := h.GetWatchHistory(ctx, userID, params)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	svc.AssertExpectations(t)
}

func TestHandler_GetUserStats(t *testing.T) {
	svc := new(MockService)
	h := NewHandler(svc)
	ctx := context.Background()
	userID := uuid.New()
	stats := &UserMovieStats{
		WatchedCount:    50,
		InProgressCount: 5,
	}

	svc.On("GetUserStats", ctx, userID).Return(stats, nil)

	result, err := h.GetUserStats(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(50), result.WatchedCount)
	svc.AssertExpectations(t)
}

func TestHandler_RefreshMetadata(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()
		movieID := uuid.New()

		svc.On("RefreshMovieMetadata", ctx, movieID).Return(nil)

		err := h.RefreshMetadata(ctx, movieID.String())
		require.NoError(t, err)
		svc.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		svc := new(MockService)
		h := NewHandler(svc)
		ctx := context.Background()

		err := h.RefreshMetadata(ctx, "invalid-id")
		assert.Error(t, err)
	})
}
