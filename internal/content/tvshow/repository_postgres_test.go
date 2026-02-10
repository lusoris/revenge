package tvshow

import (
	"context"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/content"
	"github.com/lusoris/revenge/internal/testutil"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.StopSharedPostgres()
	os.Exit(code)
}

func setupTestRepo(t *testing.T) (Repository, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	repo := NewPostgresRepository(testDB.Pool())
	return repo, testDB
}

func i32Ptr(v int32) *int32     { return &v }
func i64Ptr(v int64) *int64     { return &v }

func createTestSeries(t *testing.T, repo Repository, title string) *Series {
	t.Helper()
	tmdbID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	s, err := repo.CreateSeries(context.Background(), CreateSeriesParams{
		Title:            title,
		TMDbID:           &tmdbID,
		OriginalLanguage: "en",
		Status:           strPtr("Returning Series"),
		TotalSeasons:     3,
		TotalEpisodes:    30,
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	return s
}

func createTestSeason(t *testing.T, repo Repository, seriesID uuid.UUID, number int32) *Season {
	t.Helper()
	tmdbID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	s, err := repo.CreateSeason(context.Background(), CreateSeasonParams{
		SeriesID:     seriesID,
		TMDbID:       &tmdbID,
		SeasonNumber: number,
		Name:         "Season " + string(rune('0'+number)),
		EpisodeCount: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	return s
}

func createTestEpisode(t *testing.T, repo Repository, seriesID, seasonID uuid.UUID, seasonNum, epNum int32) *Episode {
	t.Helper()
	tmdbID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	title := "Test Episode " + string(rune('0'+epNum))
	runtime := int32(45)
	e, err := repo.CreateEpisode(context.Background(), CreateEpisodeParams{
		SeriesID:      seriesID,
		SeasonID:      seasonID,
		TMDbID:        &tmdbID,
		SeasonNumber:  seasonNum,
		EpisodeNumber: epNum,
		Title:         title,
		Overview:      strPtr("Overview for episode"),
		Runtime:       &runtime,
		AirDate:       strPtr("2024-01-15"),
	})
	require.NoError(t, err)
	require.NotNil(t, e)
	return e
}

// ============================================================================
// Series CRUD
// ============================================================================

func TestRepo_CreateAndGetSeries(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	tmdbID := int32(1396)
	tvdbID := int32(81189)
	imdbID := "tt0903747"
	voteAvg := "8.9"
	popularity := "200.5"

	s, err := repo.CreateSeries(ctx, CreateSeriesParams{
		Title:            "Breaking Bad",
		TMDbID:           &tmdbID,
		TVDbID:           &tvdbID,
		IMDbID:           &imdbID,
		OriginalTitle:    strPtr("Breaking Bad"),
		OriginalLanguage: "en",
		Tagline:          strPtr("All Hail the King"),
		Overview:         strPtr("A chemistry teacher turned meth cook"),
		Status:           strPtr("Ended"),
		Type:             strPtr("Scripted"),
		FirstAirDate:     strPtr("2008-01-20"),
		LastAirDate:      strPtr("2013-09-29"),
		VoteAverage:      &voteAvg,
		VoteCount:        i32Ptr(15000),
		Popularity:       &popularity,
		PosterPath:       strPtr("/poster_bb.jpg"),
		BackdropPath:     strPtr("/backdrop_bb.jpg"),
		TotalSeasons:     5,
		TotalEpisodes:    62,
		TrailerURL:       strPtr("https://youtube.com/watch?v=bb"),
		Homepage:         strPtr("https://www.amc.com/breakingbad"),
		TitlesI18n:       map[string]string{"de": "Breaking Bad", "ja": "ブレイキング・バッド"},
		TaglinesI18n:     map[string]string{"de": "Der König ist da"},
		OverviewsI18n:    map[string]string{"de": "Ein Chemielehrer wird zum Drogenkoch"},
		AgeRatings:       map[string]map[string]string{"US": {"MPAA": "TV-MA"}, "DE": {"FSK": "16"}},
		ExternalRatings: []ExternalRating{
			{Source: "imdb", Value: "9.5", Score: 95.0},
			{Source: "rottentomatoes", Value: "96%", Score: 96.0},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	assert.Equal(t, "Breaking Bad", s.Title)
	assert.Equal(t, int32(1396), *s.TMDbID)
	assert.Equal(t, int32(81189), *s.TVDbID)
	assert.Equal(t, "tt0903747", *s.IMDbID)
	assert.Equal(t, "Ended", *s.Status)
	assert.Equal(t, int32(5), s.TotalSeasons)
	assert.Equal(t, int32(62), s.TotalEpisodes)
	assert.Equal(t, "/poster_bb.jpg", *s.PosterPath)
	assert.NotEmpty(t, s.ID)

	// i18n
	assert.Equal(t, "ブレイキング・バッド", s.TitlesI18n["ja"])
	assert.Equal(t, "Der König ist da", s.TaglinesI18n["de"])
	assert.Equal(t, "Ein Chemielehrer wird zum Drogenkoch", s.OverviewsI18n["de"])

	// age ratings
	assert.Equal(t, "TV-MA", s.AgeRatings["US"]["MPAA"])
	assert.Equal(t, "16", s.AgeRatings["DE"]["FSK"])

	// external ratings
	require.Len(t, s.ExternalRatings, 2)
	assert.Equal(t, "imdb", s.ExternalRatings[0].Source)
	assert.Equal(t, "9.5", s.ExternalRatings[0].Value)

	// Get by ID
	got, err := repo.GetSeries(ctx, s.ID)
	require.NoError(t, err)
	assert.Equal(t, s.Title, got.Title)
	assert.Equal(t, *s.TMDbID, *got.TMDbID)
	assert.Len(t, got.ExternalRatings, 2)
}

func TestRepo_GetSeriesByTMDbID(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	s := createTestSeries(t, repo, "TMDb Lookup Series")
	got, err := repo.GetSeriesByTMDbID(ctx, *s.TMDbID)
	require.NoError(t, err)
	assert.Equal(t, s.ID, got.ID)
}

func TestRepo_GetSeriesByTVDbID(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	tvdbID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	tmdbID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	s, err := repo.CreateSeries(ctx, CreateSeriesParams{
		Title:            "TVDb Lookup Series",
		TMDbID:           &tmdbID,
		TVDbID:           &tvdbID,
		OriginalLanguage: "en",
	})
	require.NoError(t, err)

	got, err := repo.GetSeriesByTVDbID(ctx, tvdbID)
	require.NoError(t, err)
	assert.Equal(t, s.ID, got.ID)
}

func TestRepo_GetSeries_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetSeries(ctx, uuid.Must(uuid.NewV7()))
	require.Error(t, err)
}

func TestRepo_UpdateSeries(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	s := createTestSeries(t, repo, "Update Me Series")

	updated, err := repo.UpdateSeries(ctx, UpdateSeriesParams{
		ID:      s.ID,
		Title:   strPtr("Updated Series Title"),
		Status:  strPtr("Ended"),
		Tagline: strPtr("New tagline"),
		TitlesI18n: map[string]string{
			"fr": "Titre Mis à Jour",
		},
		ExternalRatings: []ExternalRating{
			{Source: "trakt", Value: "85"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated Series Title", updated.Title)
	assert.Equal(t, "Ended", *updated.Status)
	assert.Equal(t, "Titre Mis à Jour", updated.TitlesI18n["fr"])
	require.Len(t, updated.ExternalRatings, 1)
	assert.Equal(t, "trakt", updated.ExternalRatings[0].Source)
}

func TestRepo_DeleteSeries(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	s := createTestSeries(t, repo, "Delete Me Series")
	err := repo.DeleteSeries(ctx, s.ID)
	require.NoError(t, err)

	_, err = repo.GetSeries(ctx, s.ID)
	require.Error(t, err)
}

func TestRepo_ListSeries(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestSeries(t, repo, "List Series A")
	createTestSeries(t, repo, "List Series B")
	createTestSeries(t, repo, "List Series C")

	list, err := repo.ListSeries(ctx, SeriesListFilters{Limit: 10, Offset: 0, OrderBy: "title"})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3)
}

func TestRepo_ListSeries_Pagination(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		createTestSeries(t, repo, "Paginated "+string(rune('A'+i)))
	}

	page1, err := repo.ListSeries(ctx, SeriesListFilters{Limit: 2, Offset: 0, OrderBy: "title"})
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := repo.ListSeries(ctx, SeriesListFilters{Limit: 2, Offset: 2, OrderBy: "title"})
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}

func TestRepo_CountSeries(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestSeries(t, repo, "Count Series 1")
	createTestSeries(t, repo, "Count Series 2")

	count, err := repo.CountSeries(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
}

func TestRepo_SearchSeriesByTitle(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestSeries(t, repo, "Searchable Xylophone Show")

	results, err := repo.SearchSeriesByTitle(ctx, "Xylophone", 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)
	assert.Contains(t, results[0].Title, "Xylophone")
}

func TestRepo_ListRecentlyAdded(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestSeries(t, repo, "Recent Addition")

	recent, err := repo.ListRecentlyAddedSeries(ctx, 5, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(recent), 1)
}

func TestRepo_ListSeriesByStatus(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestSeries(t, repo, "Returning Show")

	list, err := repo.ListSeriesByStatus(ctx, "Returning Series", 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
}

func TestRepo_CreateSeries_Minimal(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	s, err := repo.CreateSeries(ctx, CreateSeriesParams{
		Title:            "Minimal Series",
		OriginalLanguage: "en",
	})
	require.NoError(t, err)
	assert.Equal(t, "Minimal Series", s.Title)
	assert.Nil(t, s.TMDbID)
	assert.Nil(t, s.TVDbID)
	assert.Nil(t, s.PosterPath)
}

// ============================================================================
// Seasons
// ============================================================================

func TestRepo_Seasons_CRUD(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Season Test Show")

	// Create
	season, err := repo.CreateSeason(ctx, CreateSeasonParams{
		SeriesID:      series.ID,
		TMDbID:        i32Ptr(88888),
		SeasonNumber:  1,
		Name:          "Season 1",
		Overview:      strPtr("The first season"),
		EpisodeCount:  10,
		AirDate:       strPtr("2024-01-01"),
		VoteAverage:   strPtr("8.0"),
		NamesI18n:     map[string]string{"de": "Staffel 1", "ja": "シーズン1"},
		OverviewsI18n: map[string]string{"de": "Die erste Staffel"},
	})
	require.NoError(t, err)
	assert.Equal(t, "Season 1", season.Name)
	assert.Equal(t, int32(1), season.SeasonNumber)
	assert.Equal(t, "シーズン1", season.NamesI18n["ja"])
	assert.Equal(t, "Die erste Staffel", season.OverviewsI18n["de"])

	// Get by ID
	got, err := repo.GetSeason(ctx, season.ID)
	require.NoError(t, err)
	assert.Equal(t, season.Name, got.Name)

	// Get by number
	byNum, err := repo.GetSeasonByNumber(ctx, series.ID, 1)
	require.NoError(t, err)
	assert.Equal(t, season.ID, byNum.ID)

	// Update
	updated, err := repo.UpdateSeason(ctx, UpdateSeasonParams{
		ID:           season.ID,
		Name:         strPtr("Season One (Updated)"),
		EpisodeCount: i32Ptr(12),
	})
	require.NoError(t, err)
	assert.Equal(t, "Season One (Updated)", updated.Name)

	// List by series
	repo.CreateSeason(ctx, CreateSeasonParams{ //nolint:errcheck
		SeriesID:     series.ID,
		SeasonNumber: 2,
		Name:         "Season 2",
		EpisodeCount: 8,
	})
	seasons, err := repo.ListSeasonsBySeries(ctx, series.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(seasons), 2)

	// Delete
	err = repo.DeleteSeason(ctx, season.ID)
	require.NoError(t, err)
	_, err = repo.GetSeason(ctx, season.ID)
	require.Error(t, err)
}

func TestRepo_UpsertSeason(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Upsert Season Show")

	// First insert
	s1, err := repo.UpsertSeason(ctx, CreateSeasonParams{
		SeriesID:     series.ID,
		SeasonNumber: 1,
		Name:         "Season 1 Original",
		EpisodeCount: 8,
	})
	require.NoError(t, err)
	assert.Equal(t, "Season 1 Original", s1.Name)

	// Upsert same season_number → should update
	s2, err := repo.UpsertSeason(ctx, CreateSeasonParams{
		SeriesID:     series.ID,
		SeasonNumber: 1,
		Name:         "Season 1 Updated",
		EpisodeCount: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, s1.ID, s2.ID, "upsert should return same season")
	assert.Equal(t, "Season 1 Updated", s2.Name)
}

func TestRepo_SeasonWithEpisodeCount(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Episode Count Show")
	season := createTestSeason(t, repo, series.ID, 1)

	// Add episodes
	for i := int32(1); i <= 3; i++ {
		createTestEpisode(t, repo, series.ID, season.ID, 1, i)
	}

	seasons, err := repo.ListSeasonsBySeriesWithEpisodeCount(ctx, series.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(seasons), 1)
	assert.Equal(t, int64(3), seasons[0].ActualEpisodeCount)
}

func TestRepo_DeleteSeasonsBySeries(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Delete Seasons Show")
	createTestSeason(t, repo, series.ID, 1)
	createTestSeason(t, repo, series.ID, 2)

	err := repo.DeleteSeasonsBySeries(ctx, series.ID)
	require.NoError(t, err)

	seasons, err := repo.ListSeasonsBySeries(ctx, series.ID)
	require.NoError(t, err)
	assert.Empty(t, seasons)
}

// ============================================================================
// Episodes
// ============================================================================

func TestRepo_Episodes_CRUD(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Episode CRUD Show")
	season := createTestSeason(t, repo, series.ID, 1)

	// Create with full fields
	tmdbID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	tvdbID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	ep, err := repo.CreateEpisode(ctx, CreateEpisodeParams{
		SeriesID:       series.ID,
		SeasonID:       season.ID,
		TMDbID:         &tmdbID,
		TVDbID:         &tvdbID,
		IMDbID:         strPtr("tt9999999"),
		SeasonNumber:   1,
		EpisodeNumber:  1,
		Title:          "Pilot",
		Overview:       strPtr("The beginning of everything"),
		AirDate:        strPtr("2024-01-15"),
		Runtime:        i32Ptr(58),
		VoteAverage:    strPtr("9.1"),
		VoteCount:      i32Ptr(5000),
		StillPath:      strPtr("/still_pilot.jpg"),
		ProductionCode: strPtr("101"),
		TitlesI18n:     map[string]string{"de": "Pilotfolge", "ja": "パイロット"},
		OverviewsI18n:  map[string]string{"de": "Der Anfang von allem"},
	})
	require.NoError(t, err)
	assert.Equal(t, "Pilot", ep.Title)
	assert.Equal(t, int32(1), ep.SeasonNumber)
	assert.Equal(t, int32(1), ep.EpisodeNumber)
	assert.Equal(t, "S01E01", ep.EpisodeCode())
	assert.Equal(t, "Pilotfolge", ep.TitlesI18n["de"])

	// Get by ID
	got, err := repo.GetEpisode(ctx, ep.ID)
	require.NoError(t, err)
	assert.Equal(t, ep.Title, got.Title)
	assert.Equal(t, "tt9999999", *got.IMDbID)

	// Get by TMDb ID
	byTMDb, err := repo.GetEpisodeByTMDbID(ctx, tmdbID)
	require.NoError(t, err)
	assert.Equal(t, ep.ID, byTMDb.ID)

	// Get by number
	byNum, err := repo.GetEpisodeByNumber(ctx, series.ID, 1, 1)
	require.NoError(t, err)
	assert.Equal(t, ep.ID, byNum.ID)

	// Update
	updated, err := repo.UpdateEpisode(ctx, UpdateEpisodeParams{
		ID:      ep.ID,
		Title:   strPtr("Pilot (Director's Cut)"),
		Runtime: i32Ptr(65),
	})
	require.NoError(t, err)
	assert.Equal(t, "Pilot (Director's Cut)", updated.Title)

	// Delete
	err = repo.DeleteEpisode(ctx, ep.ID)
	require.NoError(t, err)
	_, err = repo.GetEpisode(ctx, ep.ID)
	require.Error(t, err)
}

func TestRepo_UpsertEpisode(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Upsert Episode Show")
	season := createTestSeason(t, repo, series.ID, 1)

	// First insert
	e1, err := repo.UpsertEpisode(ctx, CreateEpisodeParams{
		SeriesID:      series.ID,
		SeasonID:      season.ID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Original Title",
	})
	require.NoError(t, err)
	assert.Equal(t, "Original Title", e1.Title)

	// Upsert same series+season+episode number → should update
	e2, err := repo.UpsertEpisode(ctx, CreateEpisodeParams{
		SeriesID:      series.ID,
		SeasonID:      season.ID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Updated Title",
	})
	require.NoError(t, err)
	assert.Equal(t, e1.ID, e2.ID, "upsert should return same episode")
	assert.Equal(t, "Updated Title", e2.Title)
}

func TestRepo_ListEpisodes(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "List Episodes Show")
	s1 := createTestSeason(t, repo, series.ID, 1)
	s2 := createTestSeason(t, repo, series.ID, 2)

	// 3 episodes in S1, 2 in S2
	for i := int32(1); i <= 3; i++ {
		createTestEpisode(t, repo, series.ID, s1.ID, 1, i)
	}
	for i := int32(1); i <= 2; i++ {
		createTestEpisode(t, repo, series.ID, s2.ID, 2, i)
	}

	// List by series
	allEps, err := repo.ListEpisodesBySeries(ctx, series.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, len(allEps))

	// List by season
	s1Eps, err := repo.ListEpisodesBySeason(ctx, s1.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, len(s1Eps))

	// List by season number
	s2Eps, err := repo.ListEpisodesBySeasonNumber(ctx, series.ID, 2)
	require.NoError(t, err)
	assert.Equal(t, 2, len(s2Eps))

	// Count by series
	totalCount, err := repo.CountEpisodesBySeries(ctx, series.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(5), totalCount)

	// Count by season
	s1Count, err := repo.CountEpisodesBySeason(ctx, s1.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), s1Count)
}

func TestRepo_DeleteEpisodesBySeason(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Delete Season Eps Show")
	season := createTestSeason(t, repo, series.ID, 1)
	createTestEpisode(t, repo, series.ID, season.ID, 1, 1)
	createTestEpisode(t, repo, series.ID, season.ID, 1, 2)

	err := repo.DeleteEpisodesBySeason(ctx, season.ID)
	require.NoError(t, err)

	eps, err := repo.ListEpisodesBySeason(ctx, season.ID)
	require.NoError(t, err)
	assert.Empty(t, eps)
}

func TestRepo_DeleteEpisodesBySeries(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Delete Series Eps Show")
	s1 := createTestSeason(t, repo, series.ID, 1)
	s2 := createTestSeason(t, repo, series.ID, 2)
	createTestEpisode(t, repo, series.ID, s1.ID, 1, 1)
	createTestEpisode(t, repo, series.ID, s2.ID, 2, 1)

	err := repo.DeleteEpisodesBySeries(ctx, series.ID)
	require.NoError(t, err)

	count, err := repo.CountEpisodesBySeries(ctx, series.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// ============================================================================
// Episode Files
// ============================================================================

func TestRepo_EpisodeFiles_CRUD(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Episode Files Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)

	// Create
	f, err := repo.CreateEpisodeFile(ctx, CreateEpisodeFileParams{
		EpisodeID:         ep.ID,
		FilePath:          "/tvshows/show/S01E01.mkv",
		FileName:          "S01E01.mkv",
		FileSize:          3500000000,
		Container:         strPtr("mkv"),
		Resolution:        strPtr("1080p"),
		QualityProfile:    strPtr("Bluray-1080p"),
		VideoCodec:        strPtr("h265"),
		AudioCodec:        strPtr("DTS-HD MA"),
		BitrateKbps:       i32Ptr(12000),
		AudioLanguages:    []string{"eng", "jpn"},
		SubtitleLanguages: []string{"eng", "jpn", "chi"},
		SonarrFileID:      i32Ptr(42),
	})
	require.NoError(t, err)
	assert.Equal(t, "S01E01.mkv", f.FileName)
	assert.Equal(t, "/tvshows/show/S01E01.mkv", f.FilePath)
	assert.Equal(t, int64(3500000000), f.FileSize)
	assert.Equal(t, "1080p", *f.Resolution)
	assert.Equal(t, []string{"eng", "jpn"}, f.AudioLanguages)
	assert.Equal(t, int32(42), *f.SonarrFileID)

	// Get by ID
	got, err := repo.GetEpisodeFile(ctx, f.ID)
	require.NoError(t, err)
	assert.Equal(t, f.FilePath, got.FilePath)

	// Get by path
	byPath, err := repo.GetEpisodeFileByPath(ctx, "/tvshows/show/S01E01.mkv")
	require.NoError(t, err)
	assert.Equal(t, f.ID, byPath.ID)

	// Get by Sonarr ID
	bySonarr, err := repo.GetEpisodeFileBySonarrID(ctx, 42)
	require.NoError(t, err)
	assert.Equal(t, f.ID, bySonarr.ID)

	// Update
	updated, err := repo.UpdateEpisodeFile(ctx, UpdateEpisodeFileParams{
		ID:         f.ID,
		Resolution: strPtr("2160p"),
		FileSize:   i64Ptr(7000000000),
	})
	require.NoError(t, err)
	assert.Equal(t, "2160p", *updated.Resolution)
	assert.Equal(t, int64(7000000000), updated.FileSize)

	// List by episode
	files, err := repo.ListEpisodeFilesByEpisode(ctx, ep.ID)
	require.NoError(t, err)
	assert.Len(t, files, 1)

	// Delete
	err = repo.DeleteEpisodeFile(ctx, f.ID)
	require.NoError(t, err)
	_, err = repo.GetEpisodeFile(ctx, f.ID)
	require.Error(t, err)
}

func TestRepo_DeleteEpisodeFilesByEpisode(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Delete Ep Files Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)

	for i := 0; i < 2; i++ {
		_, err := repo.CreateEpisodeFile(ctx, CreateEpisodeFileParams{
			EpisodeID: ep.ID,
			FilePath:  "/tvshows/deltest/S01E01_v" + string(rune('0'+i)) + ".mkv",
			FileName:  "S01E01_v" + string(rune('0'+i)) + ".mkv",
			FileSize:  1000000000,
		})
		require.NoError(t, err)
	}

	err := repo.DeleteEpisodeFilesByEpisode(ctx, ep.ID)
	require.NoError(t, err)

	files, err := repo.ListEpisodeFilesByEpisode(ctx, ep.ID)
	require.NoError(t, err)
	assert.Empty(t, files)
}

// ============================================================================
// Series Credits
// ============================================================================

func TestRepo_SeriesCredits(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Credits Show")

	// Cast
	cast, err := repo.CreateSeriesCredit(ctx, CreateSeriesCreditParams{
		SeriesID:     series.ID,
		TMDbPersonID: 17419,
		Name:         "Bryan Cranston",
		CreditType:   "cast",
		Character:    strPtr("Walter White"),
		CastOrder:    i32Ptr(0),
		ProfilePath:  strPtr("/cranston.jpg"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Bryan Cranston", cast.Name)
	assert.Equal(t, "Walter White", *cast.Character)

	// Crew
	_, err = repo.CreateSeriesCredit(ctx, CreateSeriesCreditParams{
		SeriesID:     series.ID,
		TMDbPersonID: 66633,
		Name:         "Vince Gilligan",
		CreditType:   "crew",
		Job:          strPtr("Creator"),
		Department:   strPtr("Production"),
	})
	require.NoError(t, err)

	// List cast
	castList, err := repo.ListSeriesCast(ctx, series.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(castList), 1)

	// Count cast
	castCount, err := repo.CountSeriesCast(ctx, series.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), castCount)

	// List crew
	crewList, err := repo.ListSeriesCrew(ctx, series.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(crewList), 1)

	// Count crew
	crewCount, err := repo.CountSeriesCrew(ctx, series.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), crewCount)

	// Delete all credits
	err = repo.DeleteSeriesCredits(ctx, series.ID)
	require.NoError(t, err)

	castList, err = repo.ListSeriesCast(ctx, series.ID, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, castList)
}

// ============================================================================
// Episode Credits
// ============================================================================

func TestRepo_EpisodeCredits(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Episode Credits Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)

	// Guest star
	gs, err := repo.CreateEpisodeCredit(ctx, CreateEpisodeCreditParams{
		EpisodeID:    ep.ID,
		TMDbPersonID: 11111,
		Name:         "Guest Star Actor",
		CreditType:   "guest_star",
		Character:    strPtr("Special Agent"),
		CastOrder:    i32Ptr(0),
	})
	require.NoError(t, err)
	assert.Equal(t, "Guest Star Actor", gs.Name)

	// Crew
	_, err = repo.CreateEpisodeCredit(ctx, CreateEpisodeCreditParams{
		EpisodeID:    ep.ID,
		TMDbPersonID: 22222,
		Name:         "Director Name",
		CreditType:   "crew",
		Job:          strPtr("Director"),
		Department:   strPtr("Directing"),
	})
	require.NoError(t, err)

	// List guest stars
	gsList, err := repo.ListEpisodeGuestStars(ctx, ep.ID)
	require.NoError(t, err)
	assert.Len(t, gsList, 1)

	// List crew
	crewList, err := repo.ListEpisodeCrew(ctx, ep.ID)
	require.NoError(t, err)
	assert.Len(t, crewList, 1)

	// Delete
	err = repo.DeleteEpisodeCredits(ctx, ep.ID)
	require.NoError(t, err)
	gsList, err = repo.ListEpisodeGuestStars(ctx, ep.ID)
	require.NoError(t, err)
	assert.Empty(t, gsList)
}

// ============================================================================
// Genres
// ============================================================================

func TestRepo_Genres(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Genre Show")

	// Add genres
	err := repo.AddSeriesGenre(ctx, series.ID, 18, "Drama")
	require.NoError(t, err)
	err = repo.AddSeriesGenre(ctx, series.ID, 80, "Crime")
	require.NoError(t, err)

	// List genres for series
	genres, err := repo.ListSeriesGenres(ctx, series.ID)
	require.NoError(t, err)
	assert.Len(t, genres, 2)

	// Add genres to another series for distinct test
	series2 := createTestSeries(t, repo, "Another Genre Show")
	err = repo.AddSeriesGenre(ctx, series2.ID, 18, "Drama") // same genre
	require.NoError(t, err)
	err = repo.AddSeriesGenre(ctx, series2.ID, 10765, "Sci-Fi & Fantasy")
	require.NoError(t, err)

	// List distinct genres across all series
	distinct, err := repo.ListDistinctSeriesGenres(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(distinct), 3) // Drama, Crime, Sci-Fi & Fantasy

	// Delete genres for series
	err = repo.DeleteSeriesGenres(ctx, series.ID)
	require.NoError(t, err)
	genres, err = repo.ListSeriesGenres(ctx, series.ID)
	require.NoError(t, err)
	assert.Empty(t, genres)
}

func TestRepo_ListSeriesByGenre(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	s1 := createTestSeries(t, repo, "Drama Show 1")
	s2 := createTestSeries(t, repo, "Drama Show 2")
	_ = repo.AddSeriesGenre(ctx, s1.ID, 18, "Drama")
	_ = repo.AddSeriesGenre(ctx, s2.ID, 18, "Drama")

	list, err := repo.ListSeriesByGenre(ctx, 18, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}

// ============================================================================
// Networks
// ============================================================================

func TestRepo_Networks(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Network Show")

	// Create network
	tmdbNetID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	net, err := repo.CreateNetwork(ctx, CreateNetworkParams{
		TMDbID:        tmdbNetID,
		Name:          "AMC",
		LogoPath:      strPtr("/amc_logo.png"),
		OriginCountry: strPtr("US"),
	})
	require.NoError(t, err)
	assert.Equal(t, "AMC", net.Name)

	// Get by ID
	got, err := repo.GetNetwork(ctx, net.ID)
	require.NoError(t, err)
	assert.Equal(t, "AMC", got.Name)

	// Get by TMDb ID
	byTMDb, err := repo.GetNetworkByTMDbID(ctx, tmdbNetID)
	require.NoError(t, err)
	assert.Equal(t, net.ID, byTMDb.ID)

	// Add network to series
	err = repo.AddSeriesNetwork(ctx, series.ID, net.ID)
	require.NoError(t, err)

	// List networks for series
	networks, err := repo.ListNetworksBySeries(ctx, series.ID)
	require.NoError(t, err)
	assert.Len(t, networks, 1)
	assert.Equal(t, "AMC", networks[0].Name)

	// Delete series networks
	err = repo.DeleteSeriesNetworks(ctx, series.ID)
	require.NoError(t, err)
	networks, err = repo.ListNetworksBySeries(ctx, series.ID)
	require.NoError(t, err)
	assert.Empty(t, networks)
}

func TestRepo_ListSeriesByNetwork(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	tmdbNetID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	net, err := repo.CreateNetwork(ctx, CreateNetworkParams{
		TMDbID: tmdbNetID,
		Name:   "HBO",
	})
	require.NoError(t, err)

	s1 := createTestSeries(t, repo, "HBO Show 1")
	s2 := createTestSeries(t, repo, "HBO Show 2")
	_ = repo.AddSeriesNetwork(ctx, s1.ID, net.ID)
	_ = repo.AddSeriesNetwork(ctx, s2.ID, net.ID)

	list, err := repo.ListSeriesByNetwork(ctx, net.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}

// ============================================================================
// Watch Progress
// ============================================================================

func TestRepo_WatchProgress(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "tvwatcher",
		Email:    "tvwatcher@test.com",
	})

	series := createTestSeries(t, repo, "Watch Progress Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep1 := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)
	ep2 := createTestEpisode(t, repo, series.ID, season.ID, 1, 2)

	// Create watch progress (in-progress)
	wp, err := repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		EpisodeID:       ep1.ID,
		ProgressSeconds: 1200,
		DurationSeconds: 2700,
		IsCompleted:     false,
	})
	require.NoError(t, err)
	assert.Equal(t, int32(1200), wp.ProgressSeconds)
	assert.False(t, wp.IsCompleted)

	// Get progress
	got, err := repo.GetWatchProgress(ctx, user.ID, ep1.ID)
	require.NoError(t, err)
	assert.Equal(t, int32(1200), got.ProgressSeconds)

	// Update progress
	updatedWP, err := repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		EpisodeID:       ep1.ID,
		ProgressSeconds: 2700,
		DurationSeconds: 2700,
		IsCompleted:     true,
	})
	require.NoError(t, err)
	assert.True(t, updatedWP.IsCompleted)

	// Mark another episode watched
	marked, err := repo.MarkEpisodeWatched(ctx, user.ID, ep2.ID, 2400)
	require.NoError(t, err)
	assert.True(t, marked.IsCompleted)

	// List watched episodes by series
	watched, err := repo.ListWatchedEpisodesBySeries(ctx, user.ID, series.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(watched), 2)

	// Get watch stats
	stats, err := repo.GetSeriesWatchStats(ctx, user.ID, series.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), stats.WatchedCount)

	// Delete watch progress for one episode
	err = repo.DeleteWatchProgress(ctx, user.ID, ep1.ID)
	require.NoError(t, err)
	_, err = repo.GetWatchProgress(ctx, user.ID, ep1.ID)
	require.Error(t, err)
}

func TestRepo_MarkEpisodesWatchedBulk(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "bulkwatcher",
		Email:    "bulkwatcher@test.com",
	})

	series := createTestSeries(t, repo, "Bulk Watch Show")
	season := createTestSeason(t, repo, series.ID, 1)
	var epIDs []uuid.UUID
	for i := int32(1); i <= 5; i++ {
		ep := createTestEpisode(t, repo, series.ID, season.ID, 1, i)
		epIDs = append(epIDs, ep.ID)
	}

	// Bulk mark all watched
	count, err := repo.MarkEpisodesWatchedBulk(ctx, user.ID, epIDs)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)

	// Verify all marked
	stats, err := repo.GetSeriesWatchStats(ctx, user.ID, series.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(5), stats.WatchedCount)
}

func TestRepo_ContinueWatching(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "continuewatcher",
		Email:    "continuewatcher@test.com",
	})

	series := createTestSeries(t, repo, "Continue Watching Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep1 := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)

	// Create in-progress watch
	_, err := repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		EpisodeID:       ep1.ID,
		ProgressSeconds: 600,
		DurationSeconds: 2700,
		IsCompleted:     false,
	})
	require.NoError(t, err)

	// Should appear in continue watching
	cwList, err := repo.ListContinueWatchingSeries(ctx, user.ID, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cwList), 1)
}

func TestRepo_NextUnwatchedEpisode(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "nextepwatcher",
		Email:    "nextepwatcher@test.com",
	})

	series := createTestSeries(t, repo, "Next Ep Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep1 := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)
	ep2 := createTestEpisode(t, repo, series.ID, season.ID, 1, 2)
	createTestEpisode(t, repo, series.ID, season.ID, 1, 3)

	// Watch first episode
	_, err := repo.MarkEpisodeWatched(ctx, user.ID, ep1.ID, 2700)
	require.NoError(t, err)

	// Next unwatched should be ep2
	next, err := repo.GetNextUnwatchedEpisode(ctx, user.ID, series.ID)
	require.NoError(t, err)
	require.NotNil(t, next)
	assert.Equal(t, ep2.ID, next.ID)
}

func TestRepo_UserTVStats(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "tvstatuser",
		Email:    "tvstatuser@test.com",
	})

	series := createTestSeries(t, repo, "Stats Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)

	_, err := repo.MarkEpisodeWatched(ctx, user.ID, ep.ID, 2700)
	require.NoError(t, err)

	stats, err := repo.GetUserTVStats(ctx, user.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.EpisodesWatched, int64(1))
	assert.GreaterOrEqual(t, stats.SeriesCount, int64(1))
}

func TestRepo_DeleteSeriesWatchProgress(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "delwatchuser",
		Email:    "delwatchuser@test.com",
	})

	series := createTestSeries(t, repo, "Delete Watch Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep1 := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)
	ep2 := createTestEpisode(t, repo, series.ID, season.ID, 1, 2)

	_, err := repo.MarkEpisodeWatched(ctx, user.ID, ep1.ID, 2700)
	require.NoError(t, err)
	_, err = repo.MarkEpisodeWatched(ctx, user.ID, ep2.ID, 2400)
	require.NoError(t, err)

	// Delete all watch progress for series
	err = repo.DeleteSeriesWatchProgress(ctx, user.ID, series.ID)
	require.NoError(t, err)

	stats, err := repo.GetSeriesWatchStats(ctx, user.ID, series.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.WatchedCount)
}

// ============================================================================
// Cascade deletes — verify referential integrity
// ============================================================================

func TestRepo_CascadeDelete_SeriesDeletesEverything(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "cascadeuser",
		Email:    "cascadeuser@test.com",
	})

	// Build full hierarchy: series → season → episode → file + credits + watch progress + genres + networks
	series := createTestSeries(t, repo, "Cascade Delete Show")
	season := createTestSeason(t, repo, series.ID, 1)
	ep := createTestEpisode(t, repo, series.ID, season.ID, 1, 1)

	_, err := repo.CreateEpisodeFile(ctx, CreateEpisodeFileParams{
		EpisodeID: ep.ID,
		FilePath:  "/cascade/S01E01.mkv",
		FileName:  "S01E01.mkv",
		FileSize:  2000000000,
	})
	require.NoError(t, err)

	_, err = repo.CreateSeriesCredit(ctx, CreateSeriesCreditParams{
		SeriesID: series.ID, TMDbPersonID: 99999, Name: "Actor", CreditType: "cast",
	})
	require.NoError(t, err)

	_, err = repo.CreateEpisodeCredit(ctx, CreateEpisodeCreditParams{
		EpisodeID: ep.ID, TMDbPersonID: 88888, Name: "Guest", CreditType: "guest_star",
	})
	require.NoError(t, err)

	err = repo.AddSeriesGenre(ctx, series.ID, 18, "Drama")
	require.NoError(t, err)

	tmdbNetID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	net, err := repo.CreateNetwork(ctx, CreateNetworkParams{TMDbID: tmdbNetID, Name: "Cascade Net"})
	require.NoError(t, err)
	err = repo.AddSeriesNetwork(ctx, series.ID, net.ID)
	require.NoError(t, err)

	_, err = repo.MarkEpisodeWatched(ctx, user.ID, ep.ID, 2700)
	require.NoError(t, err)

	// Delete the series — everything should cascade
	err = repo.DeleteSeries(ctx, series.ID)
	require.NoError(t, err)

	// Verify cascade: series gone
	_, err = repo.GetSeries(ctx, series.ID)
	require.Error(t, err)

	// Season gone
	_, err = repo.GetSeason(ctx, season.ID)
	require.Error(t, err)

	// Episode gone
	_, err = repo.GetEpisode(ctx, ep.ID)
	require.Error(t, err)

	// Credits should be gone
	castList, err := repo.ListSeriesCast(ctx, series.ID, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, castList)

	gsList, err := repo.ListEpisodeGuestStars(ctx, ep.ID)
	require.NoError(t, err)
	assert.Empty(t, gsList)

	// Genres should be gone
	genres, err := repo.ListSeriesGenres(ctx, series.ID)
	require.NoError(t, err)
	assert.Empty(t, genres)

	// Watch progress should be gone
	stats, err := repo.GetSeriesWatchStats(ctx, user.ID, series.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.WatchedCount)
}

// ============================================================================
// I18n + Complex fields roundtrip
// ============================================================================

func TestRepo_I18n_Roundtrip(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Series with full i18n + age ratings + external ratings
	s, err := repo.CreateSeries(ctx, CreateSeriesParams{
		Title:            "Multilingual Show",
		OriginalLanguage: "ja",
		TitlesI18n:       map[string]string{"en": "Multilingual Show", "ja": "多言語ショー", "ko": "다국어 쇼"},
		TaglinesI18n:     map[string]string{"en": "Words everywhere", "ja": "言葉がどこにでも"},
		OverviewsI18n:    map[string]string{"en": "A show in many languages", "ja": "多くの言語で放映される番組"},
		AgeRatings:       map[string]map[string]string{"JP": {"EIRIN": "G"}, "US": {"MPAA": "PG"}},
		ExternalRatings: []content.ExternalRating{
			{Source: "mal", Value: "8.7", Score: 87.0},
		},
	})
	require.NoError(t, err)

	// Read back
	got, err := repo.GetSeries(ctx, s.ID)
	require.NoError(t, err)
	assert.Equal(t, "多言語ショー", got.TitlesI18n["ja"])
	assert.Equal(t, "다국어 쇼", got.TitlesI18n["ko"])
	assert.Equal(t, "言葉がどこにでも", got.TaglinesI18n["ja"])
	assert.Equal(t, "G", got.AgeRatings["JP"]["EIRIN"])
	assert.Equal(t, "PG", got.AgeRatings["US"]["MPAA"])
	require.Len(t, got.ExternalRatings, 1)
	assert.Equal(t, "mal", got.ExternalRatings[0].Source)

	// Season i18n
	season, err := repo.CreateSeason(ctx, CreateSeasonParams{
		SeriesID:      s.ID,
		SeasonNumber:  1,
		Name:          "Season 1",
		NamesI18n:     map[string]string{"ja": "シーズン1"},
		OverviewsI18n: map[string]string{"ja": "最初のシーズン"},
		EpisodeCount:  12,
	})
	require.NoError(t, err)

	gotSeason, err := repo.GetSeason(ctx, season.ID)
	require.NoError(t, err)
	assert.Equal(t, "シーズン1", gotSeason.NamesI18n["ja"])
	assert.Equal(t, "最初のシーズン", gotSeason.OverviewsI18n["ja"])

	// Episode i18n
	ep, err := repo.CreateEpisode(ctx, CreateEpisodeParams{
		SeriesID:      s.ID,
		SeasonID:      season.ID,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "First Episode",
		TitlesI18n:    map[string]string{"ja": "第1話"},
		OverviewsI18n: map[string]string{"ja": "最初のエピソード"},
	})
	require.NoError(t, err)

	gotEp, err := repo.GetEpisode(ctx, ep.ID)
	require.NoError(t, err)
	assert.Equal(t, "第1話", gotEp.TitlesI18n["ja"])
}

// ============================================================================
// Edge cases
// ============================================================================

func TestRepo_CreateSeries_WithI18nAndAgeRatings(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Empty maps should be handled gracefully
	s, err := repo.CreateSeries(ctx, CreateSeriesParams{
		Title:            "Empty I18n Show",
		OriginalLanguage: "en",
		TitlesI18n:       map[string]string{},
		OverviewsI18n:    map[string]string{},
		AgeRatings:       map[string]map[string]string{},
		ExternalRatings:  []ExternalRating{},
	})
	require.NoError(t, err)

	got, err := repo.GetSeries(ctx, s.ID)
	require.NoError(t, err)
	assert.NotNil(t, got)
}

func TestRepo_SearchSeriesByTitleAnyLanguage(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.CreateSeries(ctx, CreateSeriesParams{
		Title:            "Original English Title",
		OriginalLanguage: "en",
		TitlesI18n:       map[string]string{"ja": "ユニークな日本語タイトル"},
	})
	require.NoError(t, err)

	// Search by Japanese title
	results, err := repo.SearchSeriesByTitleAnyLanguage(ctx, "ユニークな日本語タイトル", 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestRepo_UpdateSeriesStats(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	series := createTestSeries(t, repo, "Stats Update Show")
	season := createTestSeason(t, repo, series.ID, 1)
	createTestEpisode(t, repo, series.ID, season.ID, 1, 1)
	createTestEpisode(t, repo, series.ID, season.ID, 1, 2)

	// Update stats should not error
	err := repo.UpdateSeriesStats(ctx, series.ID)
	require.NoError(t, err)
}
