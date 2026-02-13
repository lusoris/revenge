package movie

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockMetadataProvider implements MetadataProvider for testing
type MockMetadataProvider struct {
	mock.Mock
}

func (m *MockMetadataProvider) SearchMovies(ctx context.Context, query string, year *int) ([]*Movie, error) {
	args := m.Called(ctx, query, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Movie), args.Error(1)
}

func (m *MockMetadataProvider) EnrichMovie(ctx context.Context, mov *Movie, opts ...MetadataRefreshOptions) error {
	args := m.Called(ctx, mov)
	return args.Error(0)
}

func (m *MockMetadataProvider) GetMovieCredits(ctx context.Context, movieID uuid.UUID, providerID string) ([]MovieCredit, error) {
	args := m.Called(ctx, movieID, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieCredit), args.Error(1)
}

func (m *MockMetadataProvider) GetMovieGenres(ctx context.Context, movieID uuid.UUID, providerID string) ([]MovieGenre, error) {
	args := m.Called(ctx, movieID, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MovieGenre), args.Error(1)
}

func (m *MockMetadataProvider) ClearCache() {
	m.Called()
}

func TestLibraryService_ScanLibrary(t *testing.T) {
	// Setup temporary directory with a dummy movie file
	tempDir := t.TempDir()

	movieFile := filepath.Join(tempDir, "The Matrix (1999).mkv")
	err := os.WriteFile(movieFile, []byte("dummy content"), 0644)
	require.NoError(t, err)

	t.Run("New movie found", func(t *testing.T) {
		repo := new(MockMovieRepository)
		metadata := new(MockMetadataProvider)
		prober := new(MockProber)

		libConfig := config.LibraryConfig{
			Paths: []string{tempDir},
		}

		svc := NewLibraryService(repo, metadata, libConfig, prober)
		ctx := context.Background()

		// Mock Prober
		mediaInfo := &MediaInfo{
			FilePath:        movieFile,
			Container:       "mkv",
			DurationSeconds: 7200,
			VideoCodec:      "h264",
			AudioStreams:    []AudioStreamInfo{{Codec: "aac"}},
		}
		prober.On("Probe", movieFile).Return(mediaInfo, nil)

		// Mock Repository - first search for existing movies (returns empty)
		repo.On("SearchMoviesByTitle", ctx, "The Matrix", int32(10), int32(0)).Return([]Movie{}, nil)

		// Mock Metadata Search (TMDb)
		tmdbMovie := &Movie{
			Title:      "The Matrix",
			Year:       new(int32(1999)),
			TMDbID:     new(int32(603)),
			Popularity: decimalPtr("100.0"),
		}
		metadata.On("SearchMovies", ctx, "The Matrix", new(1999)).Return([]*Movie{tmdbMovie}, nil)

		// Mock Metadata Enrich
		metadata.On("EnrichMovie", ctx, mock.AnythingOfType("*movie.Movie")).Return(nil)
		metadata.On("GetMovieCredits", ctx, mock.Anything, "603").Return([]MovieCredit{}, nil)
		metadata.On("GetMovieGenres", ctx, mock.Anything, "603").Return([]MovieGenre{}, nil)

		// Mock Repository
		// 1. Create Movie
		repo.On("CreateMovie", ctx, mock.MatchedBy(func(p CreateMovieParams) bool {
			return p.Title == "The Matrix" && *p.TMDbID == 603
		})).Return(&Movie{ID: uuid.Must(uuid.NewV7()), Title: "The Matrix", TMDbID: new(int32(603))}, nil)

		// 2. Create Movie File
		repo.On("CreateMovieFile", ctx, mock.MatchedBy(func(p CreateMovieFileParams) bool {
			return p.FilePath == movieFile && p.Container != nil && *p.Container == "mkv"
		})).Return(&MovieFile{}, nil)

		summary, err := svc.ScanLibrary(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, summary.TotalFiles)
		assert.Equal(t, 1, summary.MatchedFiles)
		assert.Equal(t, 1, summary.NewMovies)
		assert.Equal(t, 0, summary.ExistingMovies)
		assert.Empty(t, summary.Errors)

		repo.AssertExpectations(t)
		metadata.AssertExpectations(t)
		prober.AssertExpectations(t)
	})

	t.Run("Unmatched file", func(t *testing.T) {
		repo := new(MockMovieRepository)
		metadata := new(MockMetadataProvider)
		prober := new(MockProber)

		libConfig := config.LibraryConfig{
			Paths: []string{tempDir},
		}

		svc := NewLibraryService(repo, metadata, libConfig, prober)
		ctx := context.Background()

		// Mock Repository - first search for existing movies (returns empty)
		repo.On("SearchMoviesByTitle", ctx, "The Matrix", int32(10), int32(0)).Return([]Movie{}, nil)

		// Mock Metadata Search (return empty - no TMDb results)
		metadata.On("SearchMovies", ctx, "The Matrix", new(1999)).Return([]*Movie{}, nil)

		summary, err := svc.ScanLibrary(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, summary.TotalFiles)
		assert.Equal(t, 0, summary.MatchedFiles)
		assert.Equal(t, 1, summary.UnmatchedFiles)

		repo.AssertExpectations(t)
		metadata.AssertExpectations(t)
	})
}

func TestLibraryService_GetLibraryStats(t *testing.T) {
	t.Run("returns movie count from repo", func(t *testing.T) {
		repo := new(MockMovieRepository)
		svc := NewLibraryService(repo, nil, config.LibraryConfig{}, nil)

		repo.On("CountMovies", mock.Anything).Return(int64(42), nil)

		stats, err := svc.GetLibraryStats(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 42, stats["total_movies"])
		repo.AssertExpectations(t)
	})

	t.Run("returns error on repo failure", func(t *testing.T) {
		repo := new(MockMovieRepository)
		svc := NewLibraryService(repo, nil, config.LibraryConfig{}, nil)

		repo.On("CountMovies", mock.Anything).Return(int64(0), assert.AnError)

		stats, err := svc.GetLibraryStats(context.Background())
		assert.Error(t, err)
		assert.Nil(t, stats)
		repo.AssertExpectations(t)
	})
}

func TestLibraryService_RefreshMovie(t *testing.T) {
	repo := new(MockMovieRepository)
	metadata := new(MockMetadataProvider)
	prober := new(MockProber)
	svc := NewLibraryService(repo, metadata, config.LibraryConfig{}, prober)
	ctx := context.Background()

	movieID := uuid.Must(uuid.NewV7())
	movie := &Movie{ID: movieID, TMDbID: new(int32(603)), Title: "The Matrix"}

	repo.On("GetMovie", ctx, movieID).Return(movie, nil)
	metadata.On("EnrichMovie", ctx, movie).Return(nil)
	metadata.On("GetMovieCredits", ctx, movieID, "603").Return([]MovieCredit{}, nil)
	metadata.On("GetMovieGenres", ctx, movieID, "603").Return([]MovieGenre{}, nil)
	repo.On("UpdateMovie", ctx, mock.AnythingOfType("UpdateMovieParams")).Return(movie, nil)

	err := svc.RefreshMovie(ctx, movieID)
	require.NoError(t, err)

	repo.AssertExpectations(t)
	metadata.AssertExpectations(t)
}
