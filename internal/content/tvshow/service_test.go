package tvshow

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/lusoris/revenge/internal/content"
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

func (m *MockRepository) ListSeriesCast(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, error) {
	args := m.Called(ctx, seriesID, limit, offset)
	return args.Get(0).([]SeriesCredit), args.Error(1)
}

func (m *MockRepository) ListSeriesCrew(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, error) {
	args := m.Called(ctx, seriesID, limit, offset)
	return args.Get(0).([]SeriesCredit), args.Error(1)
}

func (m *MockRepository) CountSeriesCast(ctx context.Context, seriesID uuid.UUID) (int64, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CountSeriesCrew(ctx context.Context, seriesID uuid.UUID) (int64, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).(int64), args.Error(1)
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

func (m *MockRepository) ListDistinctSeriesGenres(ctx context.Context) ([]content.GenreSummary, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]content.GenreSummary), args.Error(1)
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

func (m *MockRepository) MarkEpisodesWatchedBulk(ctx context.Context, userID uuid.UUID, episodeIDs []uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID, episodeIDs)
	return args.Get(0).(int64), args.Error(1)
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
	svc := NewService(repo, nil)
	assert.NotNil(t, svc)
}

func TestGetSeries(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

	params := CreateSeriesParams{Title: ""}

	result, err := svc.CreateSeries(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
	assert.Nil(t, result)
}

func TestCreateSeries_DuplicateTMDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

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
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())

	repo.On("DeleteEpisodesBySeason", ctx, seasonID).Return(nil)
	repo.On("DeleteSeason", ctx, seasonID).Return(nil)

	err := svc.DeleteSeason(ctx, seasonID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRefreshSeriesMetadata_NoProvider(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	err := svc.RefreshSeriesMetadata(ctx, uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "metadata provider not configured")
}

// =============================================================================
// MockMetadataProvider
// =============================================================================

type MockMetadataProvider struct {
	mock.Mock
}

func (m *MockMetadataProvider) SearchSeries(ctx context.Context, query string, year *int) ([]*Series, error) {
	args := m.Called(ctx, query, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Series), args.Error(1)
}

func (m *MockMetadataProvider) EnrichSeries(ctx context.Context, series *Series, opts ...MetadataRefreshOptions) error {
	args := m.Called(ctx, series, opts)
	return args.Error(0)
}

func (m *MockMetadataProvider) EnrichSeason(ctx context.Context, season *Season, seriesProviderID string, opts ...MetadataRefreshOptions) error {
	args := m.Called(ctx, season, seriesProviderID, opts)
	return args.Error(0)
}

func (m *MockMetadataProvider) EnrichEpisode(ctx context.Context, episode *Episode, seriesProviderID string, opts ...MetadataRefreshOptions) error {
	args := m.Called(ctx, episode, seriesProviderID, opts)
	return args.Error(0)
}

func (m *MockMetadataProvider) GetSeriesCredits(ctx context.Context, seriesID uuid.UUID, providerID string) ([]SeriesCredit, error) {
	args := m.Called(ctx, seriesID, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SeriesCredit), args.Error(1)
}

func (m *MockMetadataProvider) GetSeriesGenres(ctx context.Context, seriesID uuid.UUID, providerID string) ([]SeriesGenre, error) {
	args := m.Called(ctx, seriesID, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SeriesGenre), args.Error(1)
}

func (m *MockMetadataProvider) GetSeriesNetworks(ctx context.Context, providerID string) ([]Network, error) {
	args := m.Called(ctx, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Network), args.Error(1)
}

func (m *MockMetadataProvider) ClearCache() {
	m.Called()
}

var _ MetadataProvider = (*MockMetadataProvider)(nil)

// =============================================================================
// Series Delegation Tests
// =============================================================================

func TestListSeries(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	filters := SeriesListFilters{OrderBy: "title", Limit: 10, Offset: 0}
	expected := []Series{
		{ID: uuid.Must(uuid.NewV7()), Title: "Breaking Bad"},
		{ID: uuid.Must(uuid.NewV7()), Title: "Better Call Saul"},
	}

	repo.On("ListSeries", ctx, filters).Return(expected, nil)

	result, err := svc.ListSeries(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestListSeries_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	filters := SeriesListFilters{OrderBy: "title", Limit: 10, Offset: 0}
	repo.On("ListSeries", ctx, filters).Return([]Series{}, errors.New("db error"))

	result, err := svc.ListSeries(ctx, filters)
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestCountSeries(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	repo.On("CountSeries", ctx).Return(int64(42), nil)

	result, err := svc.CountSeries(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), result)
	repo.AssertExpectations(t)
}

func TestCountSeries_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	repo.On("CountSeries", ctx).Return(int64(0), errors.New("db error"))

	result, err := svc.CountSeries(ctx)
	assert.Error(t, err)
	assert.Equal(t, int64(0), result)
	repo.AssertExpectations(t)
}

func TestGetSeriesByTMDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	tmdbID := int32(1396)
	expected := &Series{ID: uuid.Must(uuid.NewV7()), Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetSeriesByTMDbID", ctx, tmdbID).Return(expected, nil)

	result, err := svc.GetSeriesByTMDbID(ctx, tmdbID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSeriesByTMDbID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	repo.On("GetSeriesByTMDbID", ctx, int32(99999)).Return(nil, errors.New("not found"))

	result, err := svc.GetSeriesByTMDbID(ctx, 99999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetSeriesByTVDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	tvdbID := int32(81189)
	expected := &Series{ID: uuid.Must(uuid.NewV7()), Title: "Breaking Bad", TVDbID: &tvdbID}

	repo.On("GetSeriesByTVDbID", ctx, tvdbID).Return(expected, nil)

	result, err := svc.GetSeriesByTVDbID(ctx, tvdbID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSeriesByTVDbID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	repo.On("GetSeriesByTVDbID", ctx, int32(99999)).Return(nil, errors.New("not found"))

	result, err := svc.GetSeriesByTVDbID(ctx, 99999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetSeriesBySonarrID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	sonarrID := int32(5)
	expected := &Series{ID: uuid.Must(uuid.NewV7()), Title: "Breaking Bad", SonarrID: &sonarrID}

	repo.On("GetSeriesBySonarrID", ctx, sonarrID).Return(expected, nil)

	result, err := svc.GetSeriesBySonarrID(ctx, sonarrID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSeriesBySonarrID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	repo.On("GetSeriesBySonarrID", ctx, int32(99999)).Return(nil, errors.New("not found"))

	result, err := svc.GetSeriesBySonarrID(ctx, 99999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestSearchSeries(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	expected := []Series{
		{ID: uuid.Must(uuid.NewV7()), Title: "Breaking Bad"},
	}

	repo.On("SearchSeriesByTitle", ctx, "breaking", int32(10), int32(0)).Return(expected, nil)

	result, err := svc.SearchSeries(ctx, "breaking", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestSearchSeries_Empty(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	repo.On("SearchSeriesByTitle", ctx, "nonexistent", int32(10), int32(0)).Return([]Series{}, nil)

	result, err := svc.SearchSeries(ctx, "nonexistent", 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestListRecentlyAdded(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	expected := []Series{
		{ID: uuid.Must(uuid.NewV7()), Title: "New Show"},
	}

	repo.On("ListRecentlyAddedSeries", ctx, int32(10), int32(0)).Return(expected, nil)
	repo.On("CountSeries", ctx).Return(int64(1), nil)

	result, total, err := svc.ListRecentlyAdded(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestListByGenre(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	expected := []Series{
		{ID: uuid.Must(uuid.NewV7()), Title: "Drama Show"},
	}

	repo.On("ListSeriesByGenre", ctx, int32(18), int32(10), int32(0)).Return(expected, nil)

	result, err := svc.ListByGenre(ctx, 18, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestListByNetwork(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	networkID := uuid.Must(uuid.NewV7())
	expected := []Series{
		{ID: uuid.Must(uuid.NewV7()), Title: "AMC Show"},
	}

	repo.On("ListSeriesByNetwork", ctx, networkID, int32(10), int32(0)).Return(expected, nil)

	result, err := svc.ListByNetwork(ctx, networkID, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestListByStatus(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	expected := []Series{
		{ID: uuid.Must(uuid.NewV7()), Title: "Ended Show"},
	}

	repo.On("ListSeriesByStatus", ctx, "Ended", int32(10), int32(0)).Return(expected, nil)

	result, err := svc.ListByStatus(ctx, "Ended", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestDeleteSeries_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())

	repo.On("DeleteSeries", ctx, seriesID).Return(nil)

	err := svc.DeleteSeries(ctx, seriesID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteSeries_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())

	repo.On("DeleteSeries", ctx, seriesID).Return(errors.New("db error"))

	err := svc.DeleteSeries(ctx, seriesID)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateSeries_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	newTitle := "New Title"
	params := UpdateSeriesParams{
		ID:    seriesID,
		Title: &newTitle,
	}

	repo.On("GetSeries", ctx, seriesID).Return(nil, errors.New("not found"))

	result, err := svc.UpdateSeries(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "series not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// =============================================================================
// Season Delegation Tests
// =============================================================================

func TestGetSeason(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())
	expected := &Season{ID: seasonID, SeasonNumber: 1, Name: "Season 1"}

	repo.On("GetSeason", ctx, seasonID).Return(expected, nil)

	result, err := svc.GetSeason(ctx, seasonID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSeason_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())
	repo.On("GetSeason", ctx, seasonID).Return(nil, errors.New("not found"))

	result, err := svc.GetSeason(ctx, seasonID)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetSeasonByNumber(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := &Season{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 2, Name: "Season 2"}

	repo.On("GetSeasonByNumber", ctx, seriesID, int32(2)).Return(expected, nil)

	result, err := svc.GetSeasonByNumber(ctx, seriesID, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSeasonByNumber_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	repo.On("GetSeasonByNumber", ctx, seriesID, int32(99)).Return(nil, errors.New("not found"))

	result, err := svc.GetSeasonByNumber(ctx, seriesID, 99)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestListSeasons(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := []Season{
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 1, Name: "Season 1"},
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 2, Name: "Season 2"},
	}

	repo.On("ListSeasonsBySeries", ctx, seriesID).Return(expected, nil)

	result, err := svc.ListSeasons(ctx, seriesID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestListSeasons_Empty(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	repo.On("ListSeasonsBySeries", ctx, seriesID).Return([]Season{}, nil)

	result, err := svc.ListSeasons(ctx, seriesID)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestListSeasonsWithEpisodeCount(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := []SeasonWithEpisodeCount{
		{
			Season:             Season{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 1},
			ActualEpisodeCount: 7,
		},
		{
			Season:             Season{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 2},
			ActualEpisodeCount: 13,
		},
	}

	repo.On("ListSeasonsBySeriesWithEpisodeCount", ctx, seriesID).Return(expected, nil)

	result, err := svc.ListSeasonsWithEpisodeCount(ctx, seriesID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(7), result[0].ActualEpisodeCount)
	assert.Equal(t, int64(13), result[1].ActualEpisodeCount)
	repo.AssertExpectations(t)
}

func TestCreateSeason_DuplicateSeasonNumber(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	params := CreateSeasonParams{
		SeriesID:     seriesID,
		SeasonNumber: 1,
		Name:         "Season 1",
	}
	series := &Series{ID: seriesID, Title: "Breaking Bad"}
	existing := &Season{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 1}

	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	repo.On("GetSeasonByNumber", ctx, seriesID, int32(1)).Return(existing, nil)

	result, err := svc.CreateSeason(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpsertSeason_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	params := CreateSeasonParams{
		SeriesID:     seriesID,
		SeasonNumber: 1,
		Name:         "Season 1",
	}
	series := &Series{ID: seriesID, Title: "Breaking Bad"}
	expected := &Season{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 1}

	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	repo.On("UpsertSeason", ctx, params).Return(expected, nil)

	result, err := svc.UpsertSeason(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestUpsertSeason_SeriesNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	params := CreateSeasonParams{
		SeriesID:     seriesID,
		SeasonNumber: 1,
	}

	repo.On("GetSeries", ctx, seriesID).Return(nil, errors.New("not found"))

	result, err := svc.UpsertSeason(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "series not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateSeason_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())
	newName := "Season One"
	params := UpdateSeasonParams{ID: seasonID, Name: &newName}
	existing := &Season{ID: seasonID, SeasonNumber: 1, Name: "Season 1"}
	expected := &Season{ID: seasonID, SeasonNumber: 1, Name: "Season One"}

	repo.On("GetSeason", ctx, seasonID).Return(existing, nil)
	repo.On("UpdateSeason", ctx, params).Return(expected, nil)

	result, err := svc.UpdateSeason(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestUpdateSeason_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())
	newName := "Season One"
	params := UpdateSeasonParams{ID: seasonID, Name: &newName}

	repo.On("GetSeason", ctx, seasonID).Return(nil, errors.New("not found"))

	result, err := svc.UpdateSeason(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "season not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestDeleteSeason_ErrorDeletingEpisodes(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())

	repo.On("DeleteEpisodesBySeason", ctx, seasonID).Return(errors.New("failed to delete episodes"))

	err := svc.DeleteSeason(ctx, seasonID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete episodes")
	repo.AssertExpectations(t)
}

// =============================================================================
// Episode Delegation Tests
// =============================================================================

func TestGetEpisode(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	expected := &Episode{ID: episodeID, Title: "Pilot", SeasonNumber: 1, EpisodeNumber: 1}

	repo.On("GetEpisode", ctx, episodeID).Return(expected, nil)

	result, err := svc.GetEpisode(ctx, episodeID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetEpisode_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	repo.On("GetEpisode", ctx, episodeID).Return(nil, errors.New("not found"))

	result, err := svc.GetEpisode(ctx, episodeID)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeByTMDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	tmdbID := int32(62085)
	expected := &Episode{ID: uuid.Must(uuid.NewV7()), TMDbID: &tmdbID, Title: "Pilot"}

	repo.On("GetEpisodeByTMDbID", ctx, tmdbID).Return(expected, nil)

	result, err := svc.GetEpisodeByTMDbID(ctx, tmdbID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeByTMDbID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	repo.On("GetEpisodeByTMDbID", ctx, int32(99999)).Return(nil, errors.New("not found"))

	result, err := svc.GetEpisodeByTMDbID(ctx, 99999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeByNumber(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := &Episode{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 1, EpisodeNumber: 3, Title: "...And the Bag's in the River"}

	repo.On("GetEpisodeByNumber", ctx, seriesID, int32(1), int32(3)).Return(expected, nil)

	result, err := svc.GetEpisodeByNumber(ctx, seriesID, 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeByNumber_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	repo.On("GetEpisodeByNumber", ctx, seriesID, int32(99), int32(99)).Return(nil, errors.New("not found"))

	result, err := svc.GetEpisodeByNumber(ctx, seriesID, 99, 99)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeByFile(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	filePath := "/tv/Breaking Bad/Season 1/S01E01.mkv"
	file := &EpisodeFile{ID: uuid.Must(uuid.NewV7()), EpisodeID: episodeID, FilePath: filePath}
	expected := &Episode{ID: episodeID, Title: "Pilot"}

	repo.On("GetEpisodeFileByPath", ctx, filePath).Return(file, nil)
	repo.On("GetEpisode", ctx, episodeID).Return(expected, nil)

	result, err := svc.GetEpisodeByFile(ctx, filePath)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeByFile_FileNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	filePath := "/tv/nonexistent.mkv"
	repo.On("GetEpisodeFileByPath", ctx, filePath).Return(nil, errors.New("not found"))

	result, err := svc.GetEpisodeByFile(ctx, filePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestListEpisodesBySeries(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := []Episode{
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 1, EpisodeNumber: 1, Title: "Pilot"},
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 1, EpisodeNumber: 2, Title: "Cat's in the Bag..."},
	}

	repo.On("ListEpisodesBySeries", ctx, seriesID).Return(expected, nil)

	result, err := svc.ListEpisodesBySeries(ctx, seriesID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestListEpisodesBySeason(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())
	expected := []Episode{
		{ID: uuid.Must(uuid.NewV7()), SeasonID: seasonID, SeasonNumber: 1, EpisodeNumber: 1, Title: "Pilot"},
		{ID: uuid.Must(uuid.NewV7()), SeasonID: seasonID, SeasonNumber: 1, EpisodeNumber: 2, Title: "Cat's in the Bag..."},
	}

	repo.On("ListEpisodesBySeason", ctx, seasonID).Return(expected, nil)

	result, err := svc.ListEpisodesBySeason(ctx, seasonID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestListEpisodesBySeasonNumber(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := []Episode{
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, SeasonNumber: 2, EpisodeNumber: 1, Title: "Seven Thirty-Seven"},
	}

	repo.On("ListEpisodesBySeasonNumber", ctx, seriesID, int32(2)).Return(expected, nil)

	result, err := svc.ListEpisodesBySeasonNumber(ctx, seriesID, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestListRecentEpisodes(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	expected := []EpisodeWithSeriesInfo{
		{
			Episode:     Episode{ID: uuid.Must(uuid.NewV7()), Title: "Recent Episode"},
			SeriesTitle: "Some Series",
		},
	}

	repo.On("ListRecentEpisodes", ctx, int32(10), int32(0)).Return(expected, nil)

	result, err := svc.ListRecentEpisodes(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestListUpcomingEpisodes(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	expected := []EpisodeWithSeriesInfo{
		{
			Episode:     Episode{ID: uuid.Must(uuid.NewV7()), Title: "Upcoming Episode"},
			SeriesTitle: "Some Series",
		},
	}

	repo.On("ListUpcomingEpisodes", ctx, int32(10), int32(0)).Return(expected, nil)

	result, err := svc.ListUpcomingEpisodes(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCreateEpisode_SeasonNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())
	params := CreateEpisodeParams{
		SeriesID:      uuid.Must(uuid.NewV7()),
		SeasonID:      seasonID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
	}

	repo.On("GetSeason", ctx, seasonID).Return(nil, errors.New("not found"))

	result, err := svc.CreateEpisode(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "season not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestCreateEpisode_DuplicateEpisodeNumber(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

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
	existing := &Episode{ID: uuid.Must(uuid.NewV7()), SeasonNumber: 1, EpisodeNumber: 1}

	repo.On("GetSeason", ctx, seasonID).Return(season, nil)
	repo.On("GetEpisodeByNumber", ctx, seriesID, int32(1), int32(1)).Return(existing, nil)

	result, err := svc.CreateEpisode(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpsertEpisode_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

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
	repo.On("UpsertEpisode", ctx, params).Return(expected, nil)

	result, err := svc.UpsertEpisode(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestUpsertEpisode_SeasonNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seasonID := uuid.Must(uuid.NewV7())
	params := CreateEpisodeParams{
		SeriesID:      uuid.Must(uuid.NewV7()),
		SeasonID:      seasonID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
	}

	repo.On("GetSeason", ctx, seasonID).Return(nil, errors.New("not found"))

	result, err := svc.UpsertEpisode(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "season not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateEpisode_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	newTitle := "Updated Pilot"
	params := UpdateEpisodeParams{ID: episodeID, Title: &newTitle}
	existing := &Episode{ID: episodeID, Title: "Pilot"}
	expected := &Episode{ID: episodeID, Title: "Updated Pilot"}

	repo.On("GetEpisode", ctx, episodeID).Return(existing, nil)
	repo.On("UpdateEpisode", ctx, params).Return(expected, nil)

	result, err := svc.UpdateEpisode(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestUpdateEpisode_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	newTitle := "Updated"
	params := UpdateEpisodeParams{ID: episodeID, Title: &newTitle}

	repo.On("GetEpisode", ctx, episodeID).Return(nil, errors.New("not found"))

	result, err := svc.UpdateEpisode(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "episode not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestDeleteEpisode_ErrorDeletingFiles(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())

	repo.On("DeleteEpisodeFilesByEpisode", ctx, episodeID).Return(errors.New("failed to delete files"))

	err := svc.DeleteEpisode(ctx, episodeID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete episode files")
	repo.AssertExpectations(t)
}

// =============================================================================
// Episode File Delegation Tests
// =============================================================================

func TestGetEpisodeFile(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	fileID := uuid.Must(uuid.NewV7())
	expected := &EpisodeFile{ID: fileID, FilePath: "/tv/show/s01e01.mkv"}

	repo.On("GetEpisodeFile", ctx, fileID).Return(expected, nil)

	result, err := svc.GetEpisodeFile(ctx, fileID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeFile_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	fileID := uuid.Must(uuid.NewV7())
	repo.On("GetEpisodeFile", ctx, fileID).Return(nil, errors.New("not found"))

	result, err := svc.GetEpisodeFile(ctx, fileID)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeFileByPath(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	filePath := "/tv/show/s01e01.mkv"
	expected := &EpisodeFile{ID: uuid.Must(uuid.NewV7()), FilePath: filePath}

	repo.On("GetEpisodeFileByPath", ctx, filePath).Return(expected, nil)

	result, err := svc.GetEpisodeFileByPath(ctx, filePath)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeFileBySonarrID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	sonarrFileID := int32(42)
	expected := &EpisodeFile{ID: uuid.Must(uuid.NewV7()), SonarrFileID: &sonarrFileID}

	repo.On("GetEpisodeFileBySonarrID", ctx, sonarrFileID).Return(expected, nil)

	result, err := svc.GetEpisodeFileBySonarrID(ctx, sonarrFileID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestListEpisodeFiles(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	expected := []EpisodeFile{
		{ID: uuid.Must(uuid.NewV7()), EpisodeID: episodeID, FilePath: "/tv/show/s01e01.mkv"},
		{ID: uuid.Must(uuid.NewV7()), EpisodeID: episodeID, FilePath: "/tv/show/s01e01.mp4"},
	}

	repo.On("ListEpisodeFilesByEpisode", ctx, episodeID).Return(expected, nil)

	result, err := svc.ListEpisodeFiles(ctx, episodeID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestCreateEpisodeFile_EpisodeNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	params := CreateEpisodeFileParams{
		EpisodeID: episodeID,
		FilePath:  "/tv/show/s01e01.mkv",
		FileName:  "s01e01.mkv",
	}

	repo.On("GetEpisode", ctx, episodeID).Return(nil, errors.New("not found"))

	result, err := svc.CreateEpisodeFile(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "episode not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateEpisodeFile_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	fileID := uuid.Must(uuid.NewV7())
	newPath := "/tv/show/s01e01-updated.mkv"
	params := UpdateEpisodeFileParams{ID: fileID, FilePath: &newPath}
	existing := &EpisodeFile{ID: fileID, FilePath: "/tv/show/s01e01.mkv"}
	expected := &EpisodeFile{ID: fileID, FilePath: newPath}

	repo.On("GetEpisodeFile", ctx, fileID).Return(existing, nil)
	repo.On("UpdateEpisodeFile", ctx, params).Return(expected, nil)

	result, err := svc.UpdateEpisodeFile(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestUpdateEpisodeFile_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	fileID := uuid.Must(uuid.NewV7())
	newPath := "/tv/show/s01e01-updated.mkv"
	params := UpdateEpisodeFileParams{ID: fileID, FilePath: &newPath}

	repo.On("GetEpisodeFile", ctx, fileID).Return(nil, errors.New("not found"))

	result, err := svc.UpdateEpisodeFile(ctx, params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "episode file not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestDeleteEpisodeFile(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	fileID := uuid.Must(uuid.NewV7())

	repo.On("DeleteEpisodeFile", ctx, fileID).Return(nil)

	err := svc.DeleteEpisodeFile(ctx, fileID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteEpisodeFile_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	fileID := uuid.Must(uuid.NewV7())

	repo.On("DeleteEpisodeFile", ctx, fileID).Return(errors.New("db error"))

	err := svc.DeleteEpisodeFile(ctx, fileID)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// =============================================================================
// Credits Delegation Tests
// =============================================================================

func TestGetSeriesCast(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := []SeriesCredit{
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, Name: "Bryan Cranston", CreditType: "cast"},
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, Name: "Aaron Paul", CreditType: "cast"},
	}

	repo.On("ListSeriesCast", ctx, seriesID, int32(50), int32(0)).Return(expected, nil)
	repo.On("CountSeriesCast", ctx, seriesID).Return(int64(2), nil)

	result, total, err := svc.GetSeriesCast(ctx, seriesID, 50, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestGetSeriesCrew(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := []SeriesCredit{
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, Name: "Vince Gilligan", CreditType: "crew"},
	}

	repo.On("ListSeriesCrew", ctx, seriesID, int32(50), int32(0)).Return(expected, nil)
	repo.On("CountSeriesCrew", ctx, seriesID).Return(int64(1), nil)

	result, total, err := svc.GetSeriesCrew(ctx, seriesID, 50, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestGetEpisodeGuestStars(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	expected := []EpisodeCredit{
		{ID: uuid.Must(uuid.NewV7()), EpisodeID: episodeID, Name: "Guest Star", CreditType: "cast"},
	}

	repo.On("ListEpisodeGuestStars", ctx, episodeID).Return(expected, nil)

	result, err := svc.GetEpisodeGuestStars(ctx, episodeID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeCrew(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	episodeID := uuid.Must(uuid.NewV7())
	expected := []EpisodeCredit{
		{ID: uuid.Must(uuid.NewV7()), EpisodeID: episodeID, Name: "Director", CreditType: "crew"},
	}

	repo.On("ListEpisodeCrew", ctx, episodeID).Return(expected, nil)

	result, err := svc.GetEpisodeCrew(ctx, episodeID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

// =============================================================================
// Genres & Networks Delegation Tests
// =============================================================================

func TestGetSeriesGenres(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := []SeriesGenre{
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, TMDbGenreID: 18, Name: "Drama"},
		{ID: uuid.Must(uuid.NewV7()), SeriesID: seriesID, TMDbGenreID: 80, Name: "Crime"},
	}

	repo.On("ListSeriesGenres", ctx, seriesID).Return(expected, nil)

	result, err := svc.GetSeriesGenres(ctx, seriesID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestGetSeriesGenres_Empty(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	repo.On("ListSeriesGenres", ctx, seriesID).Return([]SeriesGenre{}, nil)

	result, err := svc.GetSeriesGenres(ctx, seriesID)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestGetSeriesNetworks(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	expected := []Network{
		{ID: uuid.Must(uuid.NewV7()), TMDbID: 174, Name: "AMC"},
	}

	repo.On("ListNetworksBySeries", ctx, seriesID).Return(expected, nil)

	result, err := svc.GetSeriesNetworks(ctx, seriesID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSeriesNetworks_Empty(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	seriesID := uuid.Must(uuid.NewV7())
	repo.On("ListNetworksBySeries", ctx, seriesID).Return([]Network{}, nil)

	result, err := svc.GetSeriesNetworks(ctx, seriesID)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

// =============================================================================
// Watch Progress Delegation Tests
// =============================================================================

func TestUpdateEpisodeProgress_EpisodeNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	repo.On("GetEpisode", ctx, episodeID).Return(nil, errors.New("not found"))

	result, err := svc.UpdateEpisodeProgress(ctx, userID, episodeID, 1200, 2700)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "episode not found")
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateEpisodeProgress_ZeroDuration(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	episode := &Episode{ID: episodeID, EpisodeNumber: 1}
	expected := &EpisodeWatched{UserID: userID, EpisodeID: episodeID}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("CreateOrUpdateWatchProgress", ctx, mock.MatchedBy(func(params CreateWatchProgressParams) bool {
		// With zero duration, should NOT be marked as completed
		return params.IsCompleted == false && params.DurationSeconds == 0
	})).Return(expected, nil)

	result, err := svc.UpdateEpisodeProgress(ctx, userID, episodeID, 100, 0)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeProgress(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	expected := &EpisodeWatched{UserID: userID, EpisodeID: episodeID, ProgressSeconds: 500}

	repo.On("GetWatchProgress", ctx, userID, episodeID).Return(expected, nil)

	result, err := svc.GetEpisodeProgress(ctx, userID, episodeID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetEpisodeProgress_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	repo.On("GetWatchProgress", ctx, userID, episodeID).Return(nil, errors.New("not found"))

	result, err := svc.GetEpisodeProgress(ctx, userID, episodeID)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestMarkEpisodeWatched_NoRuntime(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	// Episode with nil runtime - should use default 2700s (45 min)
	episode := &Episode{ID: episodeID, EpisodeNumber: 1, Runtime: nil}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("MarkEpisodeWatched", ctx, userID, episodeID, int32(2700)).Return(&EpisodeWatched{}, nil)

	err := svc.MarkEpisodeWatched(ctx, userID, episodeID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMarkEpisodeWatched_ZeroRuntime(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	runtime := int32(0)
	episode := &Episode{ID: episodeID, EpisodeNumber: 1, Runtime: &runtime}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	// Zero runtime should use default 2700s
	repo.On("MarkEpisodeWatched", ctx, userID, episodeID, int32(2700)).Return(&EpisodeWatched{}, nil)

	err := svc.MarkEpisodeWatched(ctx, userID, episodeID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMarkEpisodeWatched_EpisodeNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	repo.On("GetEpisode", ctx, episodeID).Return(nil, errors.New("not found"))

	err := svc.MarkEpisodeWatched(ctx, userID, episodeID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "episode not found")
	repo.AssertExpectations(t)
}

func TestMarkSeasonWatched_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	repo.On("ListEpisodesBySeason", ctx, seasonID).Return([]Episode{}, errors.New("db error"))

	err := svc.MarkSeasonWatched(ctx, userID, seasonID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list episodes")
	repo.AssertExpectations(t)
}

func TestMarkSeriesWatched_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	runtime := int32(45)

	season1ID := uuid.Must(uuid.NewV7())
	season2ID := uuid.Must(uuid.NewV7())
	seasons := []Season{
		{ID: season1ID, SeriesID: seriesID, SeasonNumber: 1},
		{ID: season2ID, SeriesID: seriesID, SeasonNumber: 2},
	}

	ep1 := Episode{ID: uuid.Must(uuid.NewV7()), SeasonID: season1ID, EpisodeNumber: 1, Runtime: &runtime}
	ep2 := Episode{ID: uuid.Must(uuid.NewV7()), SeasonID: season2ID, EpisodeNumber: 1, Runtime: &runtime}

	repo.On("ListSeasonsBySeries", ctx, seriesID).Return(seasons, nil)
	repo.On("ListEpisodesBySeason", ctx, season1ID).Return([]Episode{ep1}, nil)
	repo.On("ListEpisodesBySeason", ctx, season2ID).Return([]Episode{ep2}, nil)
	repo.On("GetEpisode", ctx, ep1.ID).Return(&ep1, nil)
	repo.On("GetEpisode", ctx, ep2.ID).Return(&ep2, nil)
	repo.On("MarkEpisodeWatched", ctx, userID, ep1.ID, int32(2700)).Return(&EpisodeWatched{}, nil)
	repo.On("MarkEpisodeWatched", ctx, userID, ep2.ID, int32(2700)).Return(&EpisodeWatched{}, nil)

	err := svc.MarkSeriesWatched(ctx, userID, seriesID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMarkSeriesWatched_ErrorListingSeasons(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())

	repo.On("ListSeasonsBySeries", ctx, seriesID).Return([]Season{}, errors.New("db error"))

	err := svc.MarkSeriesWatched(ctx, userID, seriesID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list seasons")
	repo.AssertExpectations(t)
}

func TestMarkEpisodesWatchedBulk_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeIDs := []uuid.UUID{uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7())}

	repo.On("MarkEpisodesWatchedBulk", ctx, userID, episodeIDs).Return(int64(3), nil)

	affected, err := svc.MarkEpisodesWatchedBulk(ctx, userID, episodeIDs)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), affected)
	repo.AssertExpectations(t)
}

func TestMarkEpisodesWatchedBulk_Empty(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())

	affected, err := svc.MarkEpisodesWatchedBulk(ctx, userID, []uuid.UUID{})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), affected)
	// Should not call repository for empty slice
	repo.AssertNotCalled(t, "MarkEpisodesWatchedBulk")
}

func TestMarkEpisodesWatchedBulk_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeIDs := []uuid.UUID{uuid.Must(uuid.NewV7())}

	repo.On("MarkEpisodesWatchedBulk", ctx, userID, episodeIDs).Return(int64(0), errors.New("db error"))

	affected, err := svc.MarkEpisodesWatchedBulk(ctx, userID, episodeIDs)
	assert.Error(t, err)
	assert.Equal(t, int64(0), affected)
	repo.AssertExpectations(t)
}

func TestRemoveEpisodeProgress(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	repo.On("DeleteWatchProgress", ctx, userID, episodeID).Return(nil)

	err := svc.RemoveEpisodeProgress(ctx, userID, episodeID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRemoveSeriesProgress(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())

	repo.On("DeleteSeriesWatchProgress", ctx, userID, seriesID).Return(nil)

	err := svc.RemoveSeriesProgress(ctx, userID, seriesID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetSeriesWatchStats(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	expected := &SeriesWatchStats{WatchedCount: 30, TotalEpisodes: 62}

	repo.On("GetSeriesWatchStats", ctx, userID, seriesID).Return(expected, nil)

	result, err := svc.GetSeriesWatchStats(ctx, userID, seriesID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetUserStats(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	expected := &UserTVStats{SeriesCount: 15, EpisodesWatched: 250}

	repo.On("GetUserTVStats", ctx, userID).Return(expected, nil)

	result, err := svc.GetUserStats(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

// =============================================================================
// Metadata Refresh Tests (with MockMetadataProvider)
// =============================================================================

func TestRefreshSeriesMetadata_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID, OriginalLanguage: "en"}

	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichSeries", ctx, series, []MetadataRefreshOptions(nil)).Return(nil)
	repo.On("UpdateSeries", ctx, mock.AnythingOfType("UpdateSeriesParams")).Return(series, nil)

	// Credits
	credits := []SeriesCredit{
		{SeriesID: seriesID, TMDbPersonID: 17419, Name: "Bryan Cranston", CreditType: "cast"},
	}
	provider.On("GetSeriesCredits", ctx, seriesID, int(1396)).Return(credits, nil)
	repo.On("DeleteSeriesCredits", ctx, seriesID).Return(nil)
	repo.On("CreateSeriesCredit", ctx, mock.AnythingOfType("CreateSeriesCreditParams")).Return(&SeriesCredit{}, nil)

	// Genres
	genres := []SeriesGenre{
		{SeriesID: seriesID, TMDbGenreID: 18, Name: "Drama"},
	}
	provider.On("GetSeriesGenres", ctx, seriesID, int(1396)).Return(genres, nil)
	repo.On("DeleteSeriesGenres", ctx, seriesID).Return(nil)
	repo.On("AddSeriesGenre", ctx, seriesID, int32(18), "Drama").Return(nil)

	err := svc.RefreshSeriesMetadata(ctx, seriesID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestRefreshSeriesMetadata_SeriesNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	repo.On("GetSeries", ctx, seriesID).Return(nil, errors.New("not found"))

	err := svc.RefreshSeriesMetadata(ctx, seriesID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get series")
	repo.AssertExpectations(t)
}

func TestRefreshSeriesMetadata_EnrichFails(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID, OriginalLanguage: "en"}

	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichSeries", ctx, series, []MetadataRefreshOptions(nil)).Return(errors.New("API error"))

	err := svc.RefreshSeriesMetadata(ctx, seriesID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enrich series")
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestRefreshSeriesMetadata_NoTMDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	// Series without TMDbID - should still succeed but skip credits/genres
	series := &Series{ID: seriesID, Title: "No TMDb", TMDbID: nil, OriginalLanguage: "en"}

	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichSeries", ctx, series, []MetadataRefreshOptions(nil)).Return(nil)
	repo.On("UpdateSeries", ctx, mock.AnythingOfType("UpdateSeriesParams")).Return(series, nil)

	err := svc.RefreshSeriesMetadata(ctx, seriesID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestRefreshSeriesMetadata_UpdateFails(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID, OriginalLanguage: "en"}

	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichSeries", ctx, series, []MetadataRefreshOptions(nil)).Return(nil)
	repo.On("UpdateSeries", ctx, mock.AnythingOfType("UpdateSeriesParams")).Return(nil, errors.New("update failed"))

	err := svc.RefreshSeriesMetadata(ctx, seriesID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update series")
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestRefreshSeasonMetadata_NoProvider(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	err := svc.RefreshSeasonMetadata(ctx, uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "metadata provider not configured")
}

func TestRefreshSeasonMetadata_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	season := &Season{ID: seasonID, SeriesID: seriesID, SeasonNumber: 1, Name: "Season 1"}
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetSeason", ctx, seasonID).Return(season, nil)
	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichSeason", ctx, season, tmdbID, []MetadataRefreshOptions(nil)).Return(nil)
	repo.On("UpdateSeason", ctx, mock.AnythingOfType("UpdateSeasonParams")).Return(season, nil)

	err := svc.RefreshSeasonMetadata(ctx, seasonID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestRefreshSeasonMetadata_SeasonNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seasonID := uuid.Must(uuid.NewV7())
	repo.On("GetSeason", ctx, seasonID).Return(nil, errors.New("not found"))

	err := svc.RefreshSeasonMetadata(ctx, seasonID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get season")
	repo.AssertExpectations(t)
}

func TestRefreshSeasonMetadata_SeriesNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	season := &Season{ID: seasonID, SeriesID: seriesID, SeasonNumber: 1}

	repo.On("GetSeason", ctx, seasonID).Return(season, nil)
	repo.On("GetSeries", ctx, seriesID).Return(nil, errors.New("not found"))

	err := svc.RefreshSeasonMetadata(ctx, seasonID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get series")
	repo.AssertExpectations(t)
}

func TestRefreshSeasonMetadata_NoTMDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	season := &Season{ID: seasonID, SeriesID: seriesID, SeasonNumber: 1}
	series := &Series{ID: seriesID, Title: "No TMDb", TMDbID: nil}

	repo.On("GetSeason", ctx, seasonID).Return(season, nil)
	repo.On("GetSeries", ctx, seriesID).Return(series, nil)

	err := svc.RefreshSeasonMetadata(ctx, seasonID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "series has no TMDb ID")
	repo.AssertExpectations(t)
}

func TestRefreshSeasonMetadata_EnrichFails(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	season := &Season{ID: seasonID, SeriesID: seriesID, SeasonNumber: 1}
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetSeason", ctx, seasonID).Return(season, nil)
	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichSeason", ctx, season, tmdbID, []MetadataRefreshOptions(nil)).Return(errors.New("API error"))

	err := svc.RefreshSeasonMetadata(ctx, seasonID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enrich season")
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestRefreshEpisodeMetadata_NoProvider(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	err := svc.RefreshEpisodeMetadata(ctx, uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "metadata provider not configured")
}

func TestRefreshEpisodeMetadata_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	episode := &Episode{ID: episodeID, SeriesID: seriesID, SeasonNumber: 1, EpisodeNumber: 1, Title: "Pilot"}
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichEpisode", ctx, episode, tmdbID, []MetadataRefreshOptions(nil)).Return(nil)
	repo.On("UpdateEpisode", ctx, mock.AnythingOfType("UpdateEpisodeParams")).Return(episode, nil)

	err := svc.RefreshEpisodeMetadata(ctx, episodeID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestRefreshEpisodeMetadata_EpisodeNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	episodeID := uuid.Must(uuid.NewV7())
	repo.On("GetEpisode", ctx, episodeID).Return(nil, errors.New("not found"))

	err := svc.RefreshEpisodeMetadata(ctx, episodeID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get episode")
	repo.AssertExpectations(t)
}

func TestRefreshEpisodeMetadata_SeriesNotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	episode := &Episode{ID: episodeID, SeriesID: seriesID, SeasonNumber: 1, EpisodeNumber: 1}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("GetSeries", ctx, seriesID).Return(nil, errors.New("not found"))

	err := svc.RefreshEpisodeMetadata(ctx, episodeID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get series")
	repo.AssertExpectations(t)
}

func TestRefreshEpisodeMetadata_NoTMDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	episode := &Episode{ID: episodeID, SeriesID: seriesID, SeasonNumber: 1, EpisodeNumber: 1}
	series := &Series{ID: seriesID, Title: "No TMDb", TMDbID: nil}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("GetSeries", ctx, seriesID).Return(series, nil)

	err := svc.RefreshEpisodeMetadata(ctx, episodeID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "series has no TMDb ID")
	repo.AssertExpectations(t)
}

func TestRefreshEpisodeMetadata_EnrichFails(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	episode := &Episode{ID: episodeID, SeriesID: seriesID, SeasonNumber: 1, EpisodeNumber: 1}
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichEpisode", ctx, episode, tmdbID, []MetadataRefreshOptions(nil)).Return(errors.New("API error"))

	err := svc.RefreshEpisodeMetadata(ctx, episodeID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enrich episode")
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestRefreshEpisodeMetadata_UpdateFails(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	episode := &Episode{ID: episodeID, SeriesID: seriesID, SeasonNumber: 1, EpisodeNumber: 1}
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetEpisode", ctx, episodeID).Return(episode, nil)
	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichEpisode", ctx, episode, tmdbID, []MetadataRefreshOptions(nil)).Return(nil)
	repo.On("UpdateEpisode", ctx, mock.AnythingOfType("UpdateEpisodeParams")).Return(nil, errors.New("update failed"))

	err := svc.RefreshEpisodeMetadata(ctx, episodeID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update episode")
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

// =============================================================================
// Conversion Helper Tests
// =============================================================================

func TestSeriesToUpdateParams(t *testing.T) {
	seriesID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	series := &Series{
		ID:               seriesID,
		TMDbID:           &tmdbID,
		Title:            "Breaking Bad",
		OriginalLanguage: "en",
		TotalSeasons:     5,
		TotalEpisodes:    62,
	}

	params := seriesToUpdateParams(series)
	assert.Equal(t, seriesID, params.ID)
	assert.Equal(t, &tmdbID, params.TMDbID)
	assert.Equal(t, "Breaking Bad", *params.Title)
	assert.Equal(t, "en", *params.OriginalLanguage)
	assert.Equal(t, int32(5), *params.TotalSeasons)
	assert.Equal(t, int32(62), *params.TotalEpisodes)
}

func TestSeasonToUpdateParams(t *testing.T) {
	seasonID := uuid.Must(uuid.NewV7())
	tmdbID := int32(3572)
	season := &Season{
		ID:           seasonID,
		TMDbID:       &tmdbID,
		SeasonNumber: 1,
		Name:         "Season 1",
		EpisodeCount: 7,
	}

	params := seasonToUpdateParams(season)
	assert.Equal(t, seasonID, params.ID)
	assert.Equal(t, &tmdbID, params.TMDbID)
	assert.Equal(t, int32(1), *params.SeasonNumber)
	assert.Equal(t, "Season 1", *params.Name)
	assert.Equal(t, int32(7), *params.EpisodeCount)
}

func TestEpisodeToUpdateParams(t *testing.T) {
	episodeID := uuid.Must(uuid.NewV7())
	tmdbID := int32(62085)
	episode := &Episode{
		ID:            episodeID,
		TMDbID:        &tmdbID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
	}

	params := episodeToUpdateParams(episode)
	assert.Equal(t, episodeID, params.ID)
	assert.Equal(t, &tmdbID, params.TMDbID)
	assert.Equal(t, int32(1), *params.SeasonNumber)
	assert.Equal(t, int32(1), *params.EpisodeNumber)
	assert.Equal(t, "Pilot", *params.Title)
}

func TestCreateSeries_NoTMDbID(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	params := CreateSeriesParams{
		Title: "Custom Series",
		// No TMDbID - should skip duplicate check
	}
	expected := &Series{ID: uuid.Must(uuid.NewV7()), Title: "Custom Series"}

	repo.On("CreateSeries", ctx, params).Return(expected, nil)

	result, err := svc.CreateSeries(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestMarkSeasonWatched_EpisodeMarkFails(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	runtime := int32(45)

	ep := Episode{ID: uuid.Must(uuid.NewV7()), SeasonID: seasonID, EpisodeNumber: 1, Runtime: &runtime}

	repo.On("ListEpisodesBySeason", ctx, seasonID).Return([]Episode{ep}, nil)
	repo.On("GetEpisode", ctx, ep.ID).Return(&ep, nil)
	repo.On("MarkEpisodeWatched", ctx, userID, ep.ID, int32(2700)).Return(nil, errors.New("mark failed"))

	err := svc.MarkSeasonWatched(ctx, userID, seasonID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to mark episode")
	repo.AssertExpectations(t)
}

func TestMarkSeriesWatched_SeasonMarkFails(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	svc := NewService(repo, nil)

	userID := uuid.Must(uuid.NewV7())
	seriesID := uuid.Must(uuid.NewV7())
	runtime := int32(45)

	season1ID := uuid.Must(uuid.NewV7())
	seasons := []Season{
		{ID: season1ID, SeriesID: seriesID, SeasonNumber: 1},
	}

	ep := Episode{ID: uuid.Must(uuid.NewV7()), SeasonID: season1ID, EpisodeNumber: 1, Runtime: &runtime}

	repo.On("ListSeasonsBySeries", ctx, seriesID).Return(seasons, nil)
	repo.On("ListEpisodesBySeason", ctx, season1ID).Return([]Episode{ep}, nil)
	repo.On("GetEpisode", ctx, ep.ID).Return(&ep, nil)
	repo.On("MarkEpisodeWatched", ctx, userID, ep.ID, int32(2700)).Return(nil, errors.New("mark failed"))

	err := svc.MarkSeriesWatched(ctx, userID, seriesID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to mark season")
	repo.AssertExpectations(t)
}

func TestRefreshSeasonMetadata_UpdateFails(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	provider := new(MockMetadataProvider)
	svc := NewService(repo, provider)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	season := &Season{ID: seasonID, SeriesID: seriesID, SeasonNumber: 1, Name: "Season 1"}
	series := &Series{ID: seriesID, Title: "Breaking Bad", TMDbID: &tmdbID}

	repo.On("GetSeason", ctx, seasonID).Return(season, nil)
	repo.On("GetSeries", ctx, seriesID).Return(series, nil)
	provider.On("EnrichSeason", ctx, season, tmdbID, []MetadataRefreshOptions(nil)).Return(nil)
	repo.On("UpdateSeason", ctx, mock.AnythingOfType("UpdateSeasonParams")).Return(nil, errors.New("update failed"))

	err := svc.RefreshSeasonMetadata(ctx, seasonID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update season")
	repo.AssertExpectations(t)
	provider.AssertExpectations(t)
}

func TestSeriesToUpdateParams_WithDateAndDecimalFields(t *testing.T) {
	seriesID := uuid.Must(uuid.NewV7())
	tmdbID := int32(1396)
	firstAirDate := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	lastAirDate := time.Date(2013, 9, 29, 0, 0, 0, 0, time.UTC)
	voteAvg, _ := decimal.New(86, 1) // 8.6
	popularity, _ := decimal.New(1234, 1)
	series := &Series{
		ID:               seriesID,
		TMDbID:           &tmdbID,
		Title:            "Breaking Bad",
		OriginalLanguage: "en",
		TotalSeasons:     5,
		TotalEpisodes:    62,
		FirstAirDate:     &firstAirDate,
		LastAirDate:      &lastAirDate,
		VoteAverage:      &voteAvg,
		Popularity:       &popularity,
	}

	params := seriesToUpdateParams(series)
	assert.Equal(t, seriesID, params.ID)
	assert.NotNil(t, params.FirstAirDate)
	assert.Equal(t, "2008-01-20", *params.FirstAirDate)
	assert.NotNil(t, params.LastAirDate)
	assert.Equal(t, "2013-09-29", *params.LastAirDate)
	assert.NotNil(t, params.VoteAverage)
	assert.NotNil(t, params.Popularity)
}

func TestSeasonToUpdateParams_WithDateAndDecimalFields(t *testing.T) {
	seasonID := uuid.Must(uuid.NewV7())
	tmdbID := int32(3572)
	airDate := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	voteAvg, _ := decimal.New(85, 1) // 8.5
	season := &Season{
		ID:           seasonID,
		TMDbID:       &tmdbID,
		SeasonNumber: 1,
		Name:         "Season 1",
		EpisodeCount: 7,
		AirDate:      &airDate,
		VoteAverage:  &voteAvg,
	}

	params := seasonToUpdateParams(season)
	assert.Equal(t, seasonID, params.ID)
	assert.NotNil(t, params.AirDate)
	assert.Equal(t, "2008-01-20", *params.AirDate)
	assert.NotNil(t, params.VoteAverage)
}

func TestEpisodeToUpdateParams_WithDateAndDecimalFields(t *testing.T) {
	episodeID := uuid.Must(uuid.NewV7())
	tmdbID := int32(62085)
	airDate := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	voteAvg, _ := decimal.New(88, 1)
	episode := &Episode{
		ID:            episodeID,
		TMDbID:        &tmdbID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
		AirDate:       &airDate,
		VoteAverage:   &voteAvg,
	}

	params := episodeToUpdateParams(episode)
	assert.Equal(t, episodeID, params.ID)
	assert.NotNil(t, params.AirDate)
	assert.Equal(t, "2008-01-20", *params.AirDate)
	assert.NotNil(t, params.VoteAverage)
}

// =============================================================================
// Module Tests
// =============================================================================

func TestProvideService(t *testing.T) {
	t.Parallel()

	repo := new(MockRepository)
	mdp := new(MockMetadataProvider)

	svc := provideService(repo, mdp, nil, slog.Default())
	assert.NotNil(t, svc)
}
