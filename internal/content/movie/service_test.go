package movie

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository implements Repository for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockRepository) GetMovieByTMDbID(ctx context.Context, tmdbID int32) (*Movie, error) {
	args := m.Called(ctx, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockRepository) GetMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error) {
	args := m.Called(ctx, imdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockRepository) GetMovieByRadarrID(ctx context.Context, radarrID int32) (*Movie, error) {
	args := m.Called(ctx, radarrID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockRepository) ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockRepository) CountMovies(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) SearchMoviesByTitle(ctx context.Context, query string, limit, offset int32) ([]Movie, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockRepository) ListMoviesByYear(ctx context.Context, year int32, limit, offset int32) ([]Movie, error) {
	args := m.Called(ctx, year, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockRepository) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockRepository) ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, error) {
	args := m.Called(ctx, minVotes, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockRepository) CreateMovie(ctx context.Context, params CreateMovieParams) (*Movie, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockRepository) UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Movie), args.Error(1)
}

func (m *MockRepository) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateMovieFile(ctx context.Context, params CreateMovieFileParams) (*MovieFile, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieFile), args.Error(1)
}

func (m *MockRepository) GetMovieFile(ctx context.Context, id uuid.UUID) (*MovieFile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieFile), args.Error(1)
}

func (m *MockRepository) GetMovieFileByPath(ctx context.Context, path string) (*MovieFile, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieFile), args.Error(1)
}

func (m *MockRepository) ListMovieFilesByMovieID(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieFile), args.Error(1)
}

func (m *MockRepository) UpdateMovieFile(ctx context.Context, params UpdateMovieFileParams) (*MovieFile, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieFile), args.Error(1)
}

func (m *MockRepository) DeleteMovieFile(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateMovieCredit(ctx context.Context, params CreateMovieCreditParams) (*MovieCredit, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieCredit), args.Error(1)
}

func (m *MockRepository) ListMovieCast(ctx context.Context, movieID uuid.UUID) ([]MovieCredit, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieCredit), args.Error(1)
}

func (m *MockRepository) ListMovieCrew(ctx context.Context, movieID uuid.UUID) ([]MovieCredit, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieCredit), args.Error(1)
}

func (m *MockRepository) DeleteMovieCredits(ctx context.Context, movieID uuid.UUID) error {
	args := m.Called(ctx, movieID)
	return args.Error(0)
}

func (m *MockRepository) CreateMovieCollection(ctx context.Context, params CreateMovieCollectionParams) (*MovieCollection, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieCollection), args.Error(1)
}

func (m *MockRepository) GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieCollection), args.Error(1)
}

func (m *MockRepository) GetMovieCollectionByTMDbID(ctx context.Context, tmdbCollectionID int32) (*MovieCollection, error) {
	args := m.Called(ctx, tmdbCollectionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieCollection), args.Error(1)
}

func (m *MockRepository) AddMovieToCollection(ctx context.Context, collectionID, movieID uuid.UUID, collectionOrder *int32) error {
	args := m.Called(ctx, collectionID, movieID, collectionOrder)
	return args.Error(0)
}

func (m *MockRepository) RemoveMovieFromCollection(ctx context.Context, collectionID, movieID uuid.UUID) error {
	args := m.Called(ctx, collectionID, movieID)
	return args.Error(0)
}

func (m *MockRepository) ListMoviesByCollection(ctx context.Context, collectionID uuid.UUID) ([]Movie, error) {
	args := m.Called(ctx, collectionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockRepository) GetCollectionForMovie(ctx context.Context, movieID uuid.UUID) (*MovieCollection, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieCollection), args.Error(1)
}

func (m *MockRepository) AddMovieGenre(ctx context.Context, movieID uuid.UUID, tmdbGenreID int32, name string) error {
	args := m.Called(ctx, movieID, tmdbGenreID, name)
	return args.Error(0)
}

func (m *MockRepository) ListMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error) {
	args := m.Called(ctx, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieGenre), args.Error(1)
}

func (m *MockRepository) DeleteMovieGenres(ctx context.Context, movieID uuid.UUID) error {
	args := m.Called(ctx, movieID)
	return args.Error(0)
}

func (m *MockRepository) ListMoviesByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]Movie, error) {
	args := m.Called(ctx, tmdbGenreID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Movie), args.Error(1)
}

func (m *MockRepository) CreateOrUpdateWatchProgress(ctx context.Context, params CreateWatchProgressParams) (*MovieWatched, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieWatched), args.Error(1)
}

func (m *MockRepository) GetWatchProgress(ctx context.Context, userID, movieID uuid.UUID) (*MovieWatched, error) {
	args := m.Called(ctx, userID, movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MovieWatched), args.Error(1)
}

func (m *MockRepository) DeleteWatchProgress(ctx context.Context, userID, movieID uuid.UUID) error {
	args := m.Called(ctx, userID, movieID)
	return args.Error(0)
}

func (m *MockRepository) ListContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ContinueWatchingItem), args.Error(1)
}

func (m *MockRepository) ListWatchedMovies(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedMovieItem, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]WatchedMovieItem), args.Error(1)
}

func (m *MockRepository) GetUserMovieStats(ctx context.Context, userID uuid.UUID) (*UserMovieStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserMovieStats), args.Error(1)
}

// Helper functions for creating test data
func ptr[T any](v T) *T {
	return &v
}

func newTestMovie() *Movie {
	return &Movie{
		ID:        uuid.New(),
		TMDbID:    ptr(int32(550)),
		IMDbID:    ptr("tt0137523"),
		Title:     "Fight Club",
		Year:      ptr(int32(1999)),
		Runtime:   ptr(int32(139)),
		Overview:  ptr("A depressed man suffering from insomnia meets a strange soap salesman."),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Tests

func TestNewService(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	assert.NotNil(t, svc)
}

func TestService_GetMovie(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)

		result, err := svc.GetMovie(ctx, movie.ID)
		require.NoError(t, err)
		assert.Equal(t, movie.ID, result.ID)
		assert.Equal(t, movie.Title, result.Title)
		repo.AssertExpectations(t)
	})

	t.Run("Not found", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		id := uuid.New()

		repo.On("GetMovie", ctx, id).Return(nil, errors.New("not found"))

		result, err := svc.GetMovie(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestService_GetMovieByTMDbID(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movie := newTestMovie()

	repo.On("GetMovieByTMDbID", ctx, int32(550)).Return(movie, nil)

	result, err := svc.GetMovieByTMDbID(ctx, 550)
	require.NoError(t, err)
	assert.Equal(t, movie.ID, result.ID)
	repo.AssertExpectations(t)
}

func TestService_GetMovieByIMDbID(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movie := newTestMovie()

	repo.On("GetMovieByIMDbID", ctx, "tt0137523").Return(movie, nil)

	result, err := svc.GetMovieByIMDbID(ctx, "tt0137523")
	require.NoError(t, err)
	assert.Equal(t, movie.ID, result.ID)
	repo.AssertExpectations(t)
}

func TestService_ListMovies(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movies := []Movie{*newTestMovie(), *newTestMovie()}
	filters := ListFilters{Limit: 10, Offset: 0}

	repo.On("ListMovies", ctx, filters).Return(movies, nil)

	result, err := svc.ListMovies(ctx, filters)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestService_SearchMovies(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movies := []Movie{*newTestMovie()}
	filters := SearchFilters{Limit: 10, Offset: 0}

	repo.On("SearchMoviesByTitle", ctx, "fight", int32(10), int32(0)).Return(movies, nil)

	result, err := svc.SearchMovies(ctx, "fight", filters)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestService_ListRecentlyAdded(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movies := []Movie{*newTestMovie()}

	repo.On("ListRecentlyAdded", ctx, int32(10), int32(0)).Return(movies, nil)

	result, err := svc.ListRecentlyAdded(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestService_ListTopRated(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movies := []Movie{*newTestMovie()}

	repo.On("ListTopRated", ctx, int32(100), int32(10), int32(0)).Return(movies, nil)

	result, err := svc.ListTopRated(ctx, 100, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestService_CreateMovie(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()
		params := CreateMovieParams{
			Title:  "Fight Club",
			TMDbID: ptr(int32(550)),
			Year:   ptr(int32(1999)),
		}

		repo.On("CreateMovie", ctx, params).Return(movie, nil)

		result, err := svc.CreateMovie(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, movie.ID, result.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Empty title", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		params := CreateMovieParams{
			Title: "",
		}

		result, err := svc.CreateMovie(ctx, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
		assert.Nil(t, result)
	})
}

func TestService_UpdateMovie(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()
		updatedTitle := "Fight Club (1999)"
		params := UpdateMovieParams{
			ID:    movie.ID,
			Title: &updatedTitle,
		}

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)
		repo.On("UpdateMovie", ctx, params).Return(movie, nil)

		result, err := svc.UpdateMovie(ctx, params)
		require.NoError(t, err)
		assert.NotNil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Movie not found", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		id := uuid.New()
		params := UpdateMovieParams{ID: id}

		repo.On("GetMovie", ctx, id).Return(nil, errors.New("not found"))

		result, err := svc.UpdateMovie(ctx, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "movie not found")
		assert.Nil(t, result)
	})

	t.Run("Empty title rejected", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()
		emptyTitle := ""
		params := UpdateMovieParams{
			ID:    movie.ID,
			Title: &emptyTitle,
		}

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)

		result, err := svc.UpdateMovie(ctx, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title cannot be empty")
		assert.Nil(t, result)
	})
}

func TestService_DeleteMovie(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	id := uuid.New()

	repo.On("DeleteMovie", ctx, id).Return(nil)

	err := svc.DeleteMovie(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_GetMovieFiles(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movieID := uuid.New()
	files := []MovieFile{
		{ID: uuid.New(), MovieID: movieID, FilePath: "/movies/fight-club.mkv"},
	}

	repo.On("ListMovieFilesByMovieID", ctx, movieID).Return(files, nil)

	result, err := svc.GetMovieFiles(ctx, movieID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestService_CreateMovieFile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()
		file := &MovieFile{
			ID:       uuid.New(),
			MovieID:  movie.ID,
			FilePath: "/movies/fight-club.mkv",
		}
		params := CreateMovieFileParams{
			MovieID:  movie.ID,
			FilePath: "/movies/fight-club.mkv",
		}

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)
		repo.On("GetMovieFileByPath", ctx, params.FilePath).Return(nil, errors.New("not found"))
		repo.On("CreateMovieFile", ctx, params).Return(file, nil)

		result, err := svc.CreateMovieFile(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, file.ID, result.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Movie not found", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movieID := uuid.New()
		params := CreateMovieFileParams{
			MovieID:  movieID,
			FilePath: "/movies/movie.mkv",
		}

		repo.On("GetMovie", ctx, movieID).Return(nil, errors.New("not found"))

		result, err := svc.CreateMovieFile(ctx, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "movie not found")
		assert.Nil(t, result)
	})

	t.Run("File already exists", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()
		existingFile := &MovieFile{
			ID:       uuid.New(),
			MovieID:  movie.ID,
			FilePath: "/movies/fight-club.mkv",
		}
		params := CreateMovieFileParams{
			MovieID:  movie.ID,
			FilePath: "/movies/fight-club.mkv",
		}

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)
		repo.On("GetMovieFileByPath", ctx, params.FilePath).Return(existingFile, nil)

		result, err := svc.CreateMovieFile(ctx, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file already exists")
		assert.Nil(t, result)
	})
}

func TestService_DeleteMovieFile(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	fileID := uuid.New()

	repo.On("DeleteMovieFile", ctx, fileID).Return(nil)

	err := svc.DeleteMovieFile(ctx, fileID)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_GetMovieCast(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movieID := uuid.New()
	cast := []MovieCredit{
		{ID: uuid.New(), MovieID: movieID, Name: "Brad Pitt", CreditType: "cast"},
	}

	repo.On("ListMovieCast", ctx, movieID).Return(cast, nil)

	result, err := svc.GetMovieCast(ctx, movieID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Brad Pitt", result[0].Name)
	repo.AssertExpectations(t)
}

func TestService_GetMovieCrew(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movieID := uuid.New()
	crew := []MovieCredit{
		{ID: uuid.New(), MovieID: movieID, Name: "David Fincher", CreditType: "crew", Job: ptr("Director")},
	}

	repo.On("ListMovieCrew", ctx, movieID).Return(crew, nil)

	result, err := svc.GetMovieCrew(ctx, movieID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "David Fincher", result[0].Name)
	repo.AssertExpectations(t)
}

func TestService_GetMovieCollection(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	collectionID := uuid.New()
	collection := &MovieCollection{
		ID:   collectionID,
		Name: "James Bond Collection",
	}

	repo.On("GetMovieCollection", ctx, collectionID).Return(collection, nil)

	result, err := svc.GetMovieCollection(ctx, collectionID)
	require.NoError(t, err)
	assert.Equal(t, "James Bond Collection", result.Name)
	repo.AssertExpectations(t)
}

func TestService_GetMoviesByCollection(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	collectionID := uuid.New()
	movies := []Movie{*newTestMovie(), *newTestMovie()}

	repo.On("ListMoviesByCollection", ctx, collectionID).Return(movies, nil)

	result, err := svc.GetMoviesByCollection(ctx, collectionID)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestService_GetCollectionForMovie(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movieID := uuid.New()
	collection := &MovieCollection{
		ID:   uuid.New(),
		Name: "The Matrix Collection",
	}

	repo.On("GetCollectionForMovie", ctx, movieID).Return(collection, nil)

	result, err := svc.GetCollectionForMovie(ctx, movieID)
	require.NoError(t, err)
	assert.Equal(t, "The Matrix Collection", result.Name)
	repo.AssertExpectations(t)
}

func TestService_GetMovieGenres(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movieID := uuid.New()
	genres := []MovieGenre{
		{ID: uuid.New(), MovieID: movieID, TMDbGenreID: 18, Name: "Drama"},
		{ID: uuid.New(), MovieID: movieID, TMDbGenreID: 53, Name: "Thriller"},
	}

	repo.On("ListMovieGenres", ctx, movieID).Return(genres, nil)

	result, err := svc.GetMovieGenres(ctx, movieID)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestService_GetMoviesByGenre(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movies := []Movie{*newTestMovie()}

	repo.On("ListMoviesByGenre", ctx, int32(18), int32(10), int32(0)).Return(movies, nil)

	result, err := svc.GetMoviesByGenre(ctx, 18, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestService_UpdateWatchProgress(t *testing.T) {
	t.Run("Partial progress", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()
		userID := uuid.New()
		watched := &MovieWatched{
			ID:              uuid.New(),
			UserID:          userID,
			MovieID:         movie.ID,
			ProgressSeconds: 3000,
			DurationSeconds: 8340,
			IsCompleted:     false,
		}

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)
		repo.On("CreateOrUpdateWatchProgress", ctx, mock.MatchedBy(func(p CreateWatchProgressParams) bool {
			return p.UserID == userID && p.MovieID == movie.ID && !p.IsCompleted
		})).Return(watched, nil)

		result, err := svc.UpdateWatchProgress(ctx, userID, movie.ID, 3000, 8340)
		require.NoError(t, err)
		assert.False(t, result.IsCompleted)
		repo.AssertExpectations(t)
	})

	t.Run("Auto-complete at 90%", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()
		userID := uuid.New()
		watched := &MovieWatched{
			ID:              uuid.New(),
			UserID:          userID,
			MovieID:         movie.ID,
			ProgressSeconds: 7600,
			DurationSeconds: 8340,
			IsCompleted:     true,
		}

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)
		repo.On("CreateOrUpdateWatchProgress", ctx, mock.MatchedBy(func(p CreateWatchProgressParams) bool {
			return p.IsCompleted // Should be true at 91%
		})).Return(watched, nil)

		result, err := svc.UpdateWatchProgress(ctx, userID, movie.ID, 7600, 8340) // ~91%
		require.NoError(t, err)
		assert.True(t, result.IsCompleted)
		repo.AssertExpectations(t)
	})

	t.Run("Movie not found", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		userID := uuid.New()
		movieID := uuid.New()

		repo.On("GetMovie", ctx, movieID).Return(nil, errors.New("not found"))

		result, err := svc.UpdateWatchProgress(ctx, userID, movieID, 100, 8340)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "movie not found")
		assert.Nil(t, result)
	})
}

func TestService_GetWatchProgress(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	userID := uuid.New()
	movieID := uuid.New()
	watched := &MovieWatched{
		ID:              uuid.New(),
		UserID:          userID,
		MovieID:         movieID,
		ProgressSeconds: 1800,
	}

	repo.On("GetWatchProgress", ctx, userID, movieID).Return(watched, nil)

	result, err := svc.GetWatchProgress(ctx, userID, movieID)
	require.NoError(t, err)
	assert.Equal(t, int32(1800), result.ProgressSeconds)
	repo.AssertExpectations(t)
}

func TestService_MarkAsWatched(t *testing.T) {
	t.Run("With runtime", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie() // Runtime is 139 minutes
		userID := uuid.New()
		watched := &MovieWatched{
			ID:          uuid.New(),
			UserID:      userID,
			MovieID:     movie.ID,
			IsCompleted: true,
		}

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)
		repo.On("CreateOrUpdateWatchProgress", ctx, mock.MatchedBy(func(p CreateWatchProgressParams) bool {
			expectedDuration := int32(139 * 60) // 139 minutes in seconds
			return p.IsCompleted && p.DurationSeconds == expectedDuration && p.ProgressSeconds == expectedDuration
		})).Return(watched, nil)

		err := svc.MarkAsWatched(ctx, userID, movie.ID)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Without runtime defaults to 2 hours", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		movie := newTestMovie()
		movie.Runtime = nil // No runtime
		userID := uuid.New()
		watched := &MovieWatched{
			ID:          uuid.New(),
			UserID:      userID,
			MovieID:     movie.ID,
			IsCompleted: true,
		}

		repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)
		repo.On("CreateOrUpdateWatchProgress", ctx, mock.MatchedBy(func(p CreateWatchProgressParams) bool {
			return p.IsCompleted && p.DurationSeconds == 7200 && p.ProgressSeconds == 7200
		})).Return(watched, nil)

		err := svc.MarkAsWatched(ctx, userID, movie.ID)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Movie not found", func(t *testing.T) {
		repo := new(MockRepository)
		svc := NewService(repo)
		ctx := context.Background()
		userID := uuid.New()
		movieID := uuid.New()

		repo.On("GetMovie", ctx, movieID).Return(nil, errors.New("not found"))

		err := svc.MarkAsWatched(ctx, userID, movieID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "movie not found")
	})
}

func TestService_RemoveWatchProgress(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	userID := uuid.New()
	movieID := uuid.New()

	repo.On("DeleteWatchProgress", ctx, userID, movieID).Return(nil)

	err := svc.RemoveWatchProgress(ctx, userID, movieID)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_GetContinueWatching(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	userID := uuid.New()
	items := []ContinueWatchingItem{
		{
			Movie:           Movie{ID: uuid.New(), Title: "Fight Club"},
			ProgressSeconds: 3000,
			DurationSeconds: 8340,
			ProgressPercent: ptr(int32(36)),
			LastWatchedAt:   time.Now(),
		},
	}

	repo.On("ListContinueWatching", ctx, userID, int32(10)).Return(items, nil)

	result, err := svc.GetContinueWatching(ctx, userID, 10)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestService_GetWatchHistory(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	userID := uuid.New()
	items := []WatchedMovieItem{
		{
			Movie:         Movie{ID: uuid.New(), Title: "Fight Club"},
			WatchCount:    2,
			LastWatchedAt: time.Now(),
		},
	}

	repo.On("ListWatchedMovies", ctx, userID, int32(10), int32(0)).Return(items, nil)

	result, err := svc.GetWatchHistory(ctx, userID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestService_GetUserStats(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	userID := uuid.New()
	stats := &UserMovieStats{
		WatchedCount:    42,
		InProgressCount: 5,
		TotalWatches:    ptr(int64(50)),
	}

	repo.On("GetUserMovieStats", ctx, userID).Return(stats, nil)

	result, err := svc.GetUserStats(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(42), result.WatchedCount)
	repo.AssertExpectations(t)
}

func TestService_RefreshMovieMetadata(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	movieID := uuid.New()

	// Currently returns "not implemented" error
	err := svc.RefreshMovieMetadata(ctx, movieID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}
