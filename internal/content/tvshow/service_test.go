package tvshow

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	mock.Mock
}

// Series operations
func (m *MockRepository) GetSeries(ctx context.Context, id uuid.UUID) (*Series, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Series), args.Error(1)
}

func (m *MockRepository) GetSeriesByTMDbID(ctx context.Context, tmdbID int32) (*Series, error) {
	args := m.Called(ctx, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Series), args.Error(1)
}

func (m *MockRepository) GetSeriesByTVDbID(ctx context.Context, tvdbID int32) (*Series, error) {
	args := m.Called(ctx, tvdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Series), args.Error(1)
}

func (m *MockRepository) GetSeriesBySonarrID(ctx context.Context, sonarrID int32) (*Series, error) {
	args := m.Called(ctx, sonarrID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Series), args.Error(1)
}

func (m *MockRepository) ListSeries(ctx context.Context, filters SeriesListFilters) ([]Series, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]Series), args.Error(1)
}

func (m *MockRepository) CountSeries(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) SearchSeriesByTitle(ctx context.Context, query string, limit, offset int32) ([]Series, error) {
	args := m.Called(ctx, query, limit, offset)
	return args.Get(0).([]Series), args.Error(1)
}

func (m *MockRepository) SearchSeriesByTitleAnyLanguage(ctx context.Context, query string, limit, offset int32) ([]Series, error) {
	args := m.Called(ctx, query, limit, offset)
	return args.Get(0).([]Series), args.Error(1)
}

func (m *MockRepository) ListRecentlyAddedSeries(ctx context.Context, limit, offset int32) ([]Series, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]Series), args.Error(1)
}

func (m *MockRepository) ListSeriesByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]Series, error) {
	args := m.Called(ctx, tmdbGenreID, limit, offset)
	return args.Get(0).([]Series), args.Error(1)
}

func (m *MockRepository) ListSeriesByNetwork(ctx context.Context, networkID uuid.UUID, limit, offset int32) ([]Series, error) {
	args := m.Called(ctx, networkID, limit, offset)
	return args.Get(0).([]Series), args.Error(1)
}

func (m *MockRepository) ListSeriesByStatus(ctx context.Context, status string, limit, offset int32) ([]Series, error) {
	args := m.Called(ctx, status, limit, offset)
	return args.Get(0).([]Series), args.Error(1)
}

func (m *MockRepository) CreateSeries(ctx context.Context, params CreateSeriesParams) (*Series, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Series), args.Error(1)
}

func (m *MockRepository) UpdateSeries(ctx context.Context, params UpdateSeriesParams) (*Series, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Series), args.Error(1)
}

func (m *MockRepository) UpdateSeriesStats(ctx context.Context, seriesID uuid.UUID) error {
	args := m.Called(ctx, seriesID)
	return args.Error(0)
}

func (m *MockRepository) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Season operations
func (m *MockRepository) GetSeason(ctx context.Context, id uuid.UUID) (*Season, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Season), args.Error(1)
}

func (m *MockRepository) GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) (*Season, error) {
	args := m.Called(ctx, seriesID, seasonNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Season), args.Error(1)
}

func (m *MockRepository) ListSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) ([]Season, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]Season), args.Error(1)
}

func (m *MockRepository) ListSeasonsBySeriesWithEpisodeCount(ctx context.Context, seriesID uuid.UUID) ([]SeasonWithEpisodeCount, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]SeasonWithEpisodeCount), args.Error(1)
}

func (m *MockRepository) CreateSeason(ctx context.Context, params CreateSeasonParams) (*Season, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Season), args.Error(1)
}

func (m *MockRepository) UpsertSeason(ctx context.Context, params CreateSeasonParams) (*Season, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Season), args.Error(1)
}

func (m *MockRepository) UpdateSeason(ctx context.Context, params UpdateSeasonParams) (*Season, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Season), args.Error(1)
}

func (m *MockRepository) DeleteSeason(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) DeleteSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) error {
	args := m.Called(ctx, seriesID)
	return args.Error(0)
}

// Episode operations
func (m *MockRepository) GetEpisode(ctx context.Context, id uuid.UUID) (*Episode, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Episode), args.Error(1)
}

func (m *MockRepository) GetEpisodeByTMDbID(ctx context.Context, tmdbID int32) (*Episode, error) {
	args := m.Called(ctx, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Episode), args.Error(1)
}

func (m *MockRepository) GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int32) (*Episode, error) {
	args := m.Called(ctx, seriesID, seasonNumber, episodeNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Episode), args.Error(1)
}

func (m *MockRepository) ListEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) ([]Episode, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]Episode), args.Error(1)
}

func (m *MockRepository) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]Episode, error) {
	args := m.Called(ctx, seasonID)
	return args.Get(0).([]Episode), args.Error(1)
}

func (m *MockRepository) ListEpisodesBySeasonNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) ([]Episode, error) {
	args := m.Called(ctx, seriesID, seasonNumber)
	return args.Get(0).([]Episode), args.Error(1)
}

func (m *MockRepository) ListRecentEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]EpisodeWithSeriesInfo), args.Error(1)
}

func (m *MockRepository) ListUpcomingEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]EpisodeWithSeriesInfo), args.Error(1)
}

func (m *MockRepository) CountEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) (int64, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CountEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) (int64, error) {
	args := m.Called(ctx, seasonID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CreateEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Episode), args.Error(1)
}

func (m *MockRepository) UpsertEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Episode), args.Error(1)
}

func (m *MockRepository) UpdateEpisode(ctx context.Context, params UpdateEpisodeParams) (*Episode, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Episode), args.Error(1)
}

func (m *MockRepository) DeleteEpisode(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) DeleteEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) error {
	args := m.Called(ctx, seasonID)
	return args.Error(0)
}

func (m *MockRepository) DeleteEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) error {
	args := m.Called(ctx, seriesID)
	return args.Error(0)
}

// Episode file operations
func (m *MockRepository) GetEpisodeFile(ctx context.Context, id uuid.UUID) (*EpisodeFile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeFile), args.Error(1)
}

func (m *MockRepository) GetEpisodeFileByPath(ctx context.Context, filePath string) (*EpisodeFile, error) {
	args := m.Called(ctx, filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeFile), args.Error(1)
}

func (m *MockRepository) GetEpisodeFileBySonarrID(ctx context.Context, sonarrFileID int32) (*EpisodeFile, error) {
	args := m.Called(ctx, sonarrFileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeFile), args.Error(1)
}

func (m *MockRepository) ListEpisodeFilesByEpisode(ctx context.Context, episodeID uuid.UUID) ([]EpisodeFile, error) {
	args := m.Called(ctx, episodeID)
	return args.Get(0).([]EpisodeFile), args.Error(1)
}

func (m *MockRepository) CreateEpisodeFile(ctx context.Context, params CreateEpisodeFileParams) (*EpisodeFile, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeFile), args.Error(1)
}

func (m *MockRepository) UpdateEpisodeFile(ctx context.Context, params UpdateEpisodeFileParams) (*EpisodeFile, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeFile), args.Error(1)
}

func (m *MockRepository) DeleteEpisodeFile(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) DeleteEpisodeFilesByEpisode(ctx context.Context, episodeID uuid.UUID) error {
	args := m.Called(ctx, episodeID)
	return args.Error(0)
}

// Credits operations
func (m *MockRepository) CreateSeriesCredit(ctx context.Context, params CreateSeriesCreditParams) (*SeriesCredit, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SeriesCredit), args.Error(1)
}

func (m *MockRepository) ListSeriesCast(ctx context.Context, seriesID uuid.UUID) ([]SeriesCredit, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]SeriesCredit), args.Error(1)
}

func (m *MockRepository) ListSeriesCrew(ctx context.Context, seriesID uuid.UUID) ([]SeriesCredit, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]SeriesCredit), args.Error(1)
}

func (m *MockRepository) DeleteSeriesCredits(ctx context.Context, seriesID uuid.UUID) error {
	args := m.Called(ctx, seriesID)
	return args.Error(0)
}

func (m *MockRepository) CreateEpisodeCredit(ctx context.Context, params CreateEpisodeCreditParams) (*EpisodeCredit, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeCredit), args.Error(1)
}

func (m *MockRepository) ListEpisodeGuestStars(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error) {
	args := m.Called(ctx, episodeID)
	return args.Get(0).([]EpisodeCredit), args.Error(1)
}

func (m *MockRepository) ListEpisodeCrew(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error) {
	args := m.Called(ctx, episodeID)
	return args.Get(0).([]EpisodeCredit), args.Error(1)
}

func (m *MockRepository) DeleteEpisodeCredits(ctx context.Context, episodeID uuid.UUID) error {
	args := m.Called(ctx, episodeID)
	return args.Error(0)
}

// Genre operations
func (m *MockRepository) AddSeriesGenre(ctx context.Context, seriesID uuid.UUID, tmdbGenreID int32, name string) error {
	args := m.Called(ctx, seriesID, tmdbGenreID, name)
	return args.Error(0)
}

func (m *MockRepository) ListSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]SeriesGenre, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]SeriesGenre), args.Error(1)
}

func (m *MockRepository) DeleteSeriesGenres(ctx context.Context, seriesID uuid.UUID) error {
	args := m.Called(ctx, seriesID)
	return args.Error(0)
}

// Network operations
func (m *MockRepository) CreateNetwork(ctx context.Context, params CreateNetworkParams) (*Network, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Network), args.Error(1)
}

func (m *MockRepository) GetNetwork(ctx context.Context, id uuid.UUID) (*Network, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Network), args.Error(1)
}

func (m *MockRepository) GetNetworkByTMDbID(ctx context.Context, tmdbNetworkID int32) (*Network, error) {
	args := m.Called(ctx, tmdbNetworkID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Network), args.Error(1)
}

func (m *MockRepository) ListNetworksBySeries(ctx context.Context, seriesID uuid.UUID) ([]Network, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]Network), args.Error(1)
}

func (m *MockRepository) AddSeriesNetwork(ctx context.Context, seriesID, networkID uuid.UUID) error {
	args := m.Called(ctx, seriesID, networkID)
	return args.Error(0)
}

func (m *MockRepository) DeleteSeriesNetworks(ctx context.Context, seriesID uuid.UUID) error {
	args := m.Called(ctx, seriesID)
	return args.Error(0)
}

// Watch progress operations
func (m *MockRepository) CreateOrUpdateWatchProgress(ctx context.Context, params CreateWatchProgressParams) (*EpisodeWatched, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeWatched), args.Error(1)
}

func (m *MockRepository) MarkEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID, durationSeconds int32) (*EpisodeWatched, error) {
	args := m.Called(ctx, userID, episodeID, durationSeconds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeWatched), args.Error(1)
}

func (m *MockRepository) GetWatchProgress(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeWatched, error) {
	args := m.Called(ctx, userID, episodeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*EpisodeWatched), args.Error(1)
}

func (m *MockRepository) DeleteWatchProgress(ctx context.Context, userID, episodeID uuid.UUID) error {
	args := m.Called(ctx, userID, episodeID)
	return args.Error(0)
}

func (m *MockRepository) DeleteSeriesWatchProgress(ctx context.Context, userID, seriesID uuid.UUID) error {
	args := m.Called(ctx, userID, seriesID)
	return args.Error(0)
}

func (m *MockRepository) ListContinueWatchingSeries(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]ContinueWatchingItem), args.Error(1)
}

func (m *MockRepository) ListWatchedEpisodesBySeries(ctx context.Context, userID, seriesID uuid.UUID) ([]WatchedEpisodeItem, error) {
	args := m.Called(ctx, userID, seriesID)
	return args.Get(0).([]WatchedEpisodeItem), args.Error(1)
}

func (m *MockRepository) ListWatchedEpisodesByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedEpisodeItem, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]WatchedEpisodeItem), args.Error(1)
}

func (m *MockRepository) GetSeriesWatchStats(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesWatchStats, error) {
	args := m.Called(ctx, userID, seriesID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SeriesWatchStats), args.Error(1)
}

func (m *MockRepository) GetUserTVStats(ctx context.Context, userID uuid.UUID) (*UserTVStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserTVStats), args.Error(1)
}

func (m *MockRepository) GetNextUnwatchedEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*Episode, error) {
	args := m.Called(ctx, userID, seriesID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Episode), args.Error(1)
}

// Ensure MockRepository implements Repository
var _ Repository = (*MockRepository)(nil)

// =============================================================================
// Service Tests
// =============================================================================

func TestNewService(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)
	assert.NotNil(t, svc)
}

func TestGetSeries(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	seriesID := uuid.Must(uuid.NewV7())
	expected := &Series{ID: seriesID, Title: "Breaking Bad"}

	repo.On("GetSeries", ctx, seriesID).Return(expected, nil)

	result, err := svc.GetSeries(ctx, seriesID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSeries_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	seriesID := uuid.Must(uuid.NewV7())
	repo.On("GetSeries", ctx, seriesID).Return(nil, errors.New("not found"))

	result, err := svc.GetSeries(ctx, seriesID)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestCreateSeries_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	tmdbID := int32(1396)
	params := CreateSeriesParams{
		Title:  "Breaking Bad",
		TMDbID: &tmdbID,
	}
	expected := &Series{ID: uuid.Must(uuid.NewV7()), Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetSeriesByTMDbID", ctx, tmdbID).Return(nil, errors.New("not found"))
	repo.On("CreateSeries", ctx, params).Return(expected, nil)

	result, err := svc.CreateSeries(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCreateSeries_EmptyTitle(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	params := CreateSeriesParams{Title: ""}

	result, err := svc.CreateSeries(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
	assert.Nil(t, result)
}

func TestCreateSeries_DuplicateTMDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	tmdbID := int32(1396)
	params := CreateSeriesParams{
		Title:  "Breaking Bad",
		TMDbID: &tmdbID,
	}
	existing := &Series{ID: uuid.Must(uuid.NewV7()), Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetSeriesByTMDbID", ctx, tmdbID).Return(existing, nil)

	result, err := svc.CreateSeries(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateSeries_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	seriesID := uuid.Must(uuid.NewV7())
	newTitle := "Better Call Saul"
	params := UpdateSeriesParams{
		ID:    seriesID,
		Title: &newTitle,
	}
	existing := &Series{ID: seriesID, Title: "Breaking Bad"}
	expected := &Series{ID: seriesID, Title: newTitle}

	repo.On("GetSeries", ctx, seriesID).Return(existing, nil)
	repo.On("UpdateSeries", ctx, params).Return(expected, nil)

	result, err := svc.UpdateSeries(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestUpdateSeries_EmptyTitle(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	seriesID := uuid.Must(uuid.NewV7())
	emptyTitle := ""
	params := UpdateSeriesParams{
		ID:    seriesID,
		Title: &emptyTitle,
	}
	existing := &Series{ID: seriesID, Title: "Breaking Bad"}

	repo.On("GetSeries", ctx, seriesID).Return(existing, nil)

	result, err := svc.UpdateSeries(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title cannot be empty")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestCreateSeason_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	seriesID := uuid.Must(uuid.NewV7())
	params := CreateSeasonParams{
		SeriesID:     seriesID,
		SeasonNumber: 1,
		Name:         "Season 1",
	}
	series := &Series{ID: seriesID, Title: "Breaking Bad"}
	expected := &Season{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 1}

	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	repo.On("GetSeasonByNumber", ctx, seriesID, int32(1)).Return(nil, errors.New("not found"))
	repo.On("CreateSeason", ctx, params).Return(expected, nil)

	result, err := svc.CreateSeason(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCreateSeason_SeriesNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	seriesID := uuid.Must(uuid.NewV7())
	params := CreateSeasonParams{
		SeriesID:     seriesID,
		SeasonNumber: 1,
	}

	repo.On("GetSeries", ctx, seriesID).Return(nil, errors.New("not found"))

	result, err := svc.CreateSeason(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "series not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestCreateEpisode_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	params := CreateEpisodeParams{
		SeriesID:      seriesID,
		SeasonID:      seasonID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
	}
	season := &Season{ID: seasonID, SeasonNumber: 1}
	expected := &Episode{ID: uuid.Must(uuid.NewV7()), SeasonID: seasonID, EpisodeNumber: 1}

	repo.On("GetSeason", ctx, seasonID).Return(season, nil)
	repo.On("GetEpisodeByNumber", ctx, seriesID, int32(1), int32(1)).Return(nil, errors.New("not found"))
	repo.On("CreateEpisode", ctx, params).Return(expected, nil)

	result, err := svc.CreateEpisode(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCreateEpisodeFile_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	episodeID := uuid.Must(uuid.NewV7())
	params := CreateEpisodeFileParams{
		EpisodeID: episodeID,
		FilePath:  "/tv/Breaking Bad/Season 1/S01E01.mkv",
		FileName:  "S01E01.mkv",
	}
	episode := &Episode{ID: episodeID, EpisodeNumber: 1}
	expected := &EpisodeFile{ID: uuid.Must(uuid.NewV7()), EpisodeID: episodeID}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("GetEpisodeFileByPath", ctx, params.FilePath).Return(nil, errors.New("not found"))
	repo.On("CreateEpisodeFile", ctx, params).Return(expected, nil)

	result, err := svc.CreateEpisodeFile(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCreateEpisodeFile_DuplicatePath(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	episodeID := uuid.Must(uuid.NewV7())
	params := CreateEpisodeFileParams{
		EpisodeID: episodeID,
		FilePath:  "/tv/Breaking Bad/Season 1/S01E01.mkv",
		FileName:  "S01E01.mkv",
	}
	episode := &Episode{ID: episodeID, EpisodeNumber: 1}
	existing := &EpisodeFile{ID: uuid.Must(uuid.NewV7()), EpisodeID: episodeID}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("GetEpisodeFileByPath", ctx, params.FilePath).Return(existing, nil)

	result, err := svc.CreateEpisodeFile(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file already exists")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateEpisodeProgress_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	episode := &Episode{ID: episodeID, EpisodeNumber: 1}
	expected := &EpisodeWatched{UserID: userID, EpisodeID: episodeID, ProgressSeconds: 1200}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("CreateOrUpdateWatchProgress", ctx, mock.AnythingOfType("CreateWatchProgressParams")).Return(expected, nil)

	result, err := svc.UpdateEpisodeProgress(ctx, userID, episodeID, 1200, 2700)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestUpdateEpisodeProgress_Completed(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	episode := &Episode{ID: episodeID, EpisodeNumber: 1}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("CreateOrUpdateWatchProgress", ctx, mock.MatchedBy(func(params CreateWatchProgressParams) bool {
		// >90% watched should be marked as completed
		return params.IsCompleted == true
	})).Return(&EpisodeWatched{}, nil)

	_, err := svc.UpdateEpisodeProgress(ctx, userID, episodeID, 2500, 2700)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMarkEpisodeWatched(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	runtime := int32(45)
	episode := &Episode{ID: episodeID, EpisodeNumber: 1, Runtime: &runtime}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("MarkEpisodeWatched", ctx, userID, episodeID, int32(2700)).Return(&EpisodeWatched{}, nil)

	err := svc.MarkEpisodeWatched(ctx, userID, episodeID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMarkSeasonWatched(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	userID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	runtime := int32(45)

	episodes := []Episode{
		{ID: uuid.Must(uuid.NewV7()), SeasonID: seasonID, EpisodeNumber: 1, Runtime: &runtime},
		{ID: uuid.Must(uuid.NewV7()), SeasonID: seasonID, EpisodeNumber: 2, Runtime: &runtime},
	}

	repo.On("ListEpisodesBySeason", ctx, seasonID).Return(episodes, nil)
	repo.On("GetEpisode", ctx, episodes[0].ID).Return(&episodes[0], nil)
	repo.On("GetEpisode", ctx, episodes[1].ID).Return(&episodes[1], nil)
	repo.On("MarkEpisodeWatched", ctx, userID, episodes[0].ID, int32(2700)).Return(&EpisodeWatched{}, nil)
	repo.On("MarkEpisodeWatched", ctx, userID, episodes[1].ID, int32(2700)).Return(&EpisodeWatched{}, nil)

	err := svc.MarkSeasonWatched(ctx, userID, seasonID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetContinueWatching(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	userID := uuid.Must(uuid.NewV7())
	expected := []ContinueWatchingItem{
		{
			Series:           &Series{ID: uuid.Must(uuid.NewV7()), Title: "Breaking Bad"},
			LastEpisodeID:    uuid.Must(uuid.NewV7()),
			LastEpisodeTitle: "Pilot",
			ProgressSeconds:  1200,
			DurationSeconds:  2700,
		},
	}

	repo.On("ListContinueWatchingSeries", ctx, userID, int32(10)).Return(expected, nil)

	result, err := svc.GetContinueWatching(ctx, userID, 10)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetNextEpisode(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	userID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	expected := &Episode{ID: uuid.Must(uuid.NewV7()), EpisodeNumber: 2, Title: "Next Episode"}

	repo.On("GetNextUnwatchedEpisode", ctx, userID, seriesID).Return(expected, nil)

	result, err := svc.GetNextEpisode(ctx, userID, seriesID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestDeleteEpisode_DeletesFiles(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	episodeID := uuid.Must(uuid.NewV7())

	repo.On("DeleteEpisodeFilesByEpisode", ctx, episodeID).Return(nil)
	repo.On("DeleteEpisode", ctx, episodeID).Return(nil)

	err := svc.DeleteEpisode(ctx, episodeID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteSeason_DeletesEpisodes(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	seasonID := uuid.Must(uuid.NewV7())

	repo.On("DeleteEpisodesBySeason", ctx, seasonID).Return(nil)
	repo.On("DeleteSeason", ctx, seasonID).Return(nil)

	err := svc.DeleteSeason(ctx, seasonID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRefreshSeriesMetadata_NotImplemented(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo)

	err := svc.RefreshSeriesMetadata(ctx, uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}
