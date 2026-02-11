package movie

import (
	"context"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func createTestMovie(t *testing.T, repo Repository, title string) *Movie {
	t.Helper()
	tmdbID := rand.Int32N(9000000) + 1000000 //nolint:gosec
	year := int32(2024)
	overview := "Test overview for " + title
	m, err := repo.CreateMovie(context.Background(), CreateMovieParams{
		Title:    title,
		TMDbID:   &tmdbID,
		Year:     &year,
		Overview: &overview,
		Status:   strPtr("Released"),
	})
	require.NoError(t, err)
	require.NotNil(t, m)
	return m
}

func strPtr(s string) *string   { return &s }
func i32Ptr(v int32) *int32     { return &v }

// ============================================================================
// Movie CRUD
// ============================================================================

func TestRepo_CreateAndGetMovie(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	tmdbID := int32(27205)
	imdbID := "tt1375666"
	year := int32(2010)
	runtime := int32(148)
	overview := "A thief enters dreams"
	tagline := "Your mind is the scene of the crime"
	voteAvg := "8.5"
	popularity := "100.5"
	budget := int64(200000000)
	revenue := int64(1000000000)

	m, err := repo.CreateMovie(ctx, CreateMovieParams{
		Title:            "Inception",
		TMDbID:           &tmdbID,
		IMDbID:           &imdbID,
		Year:             &year,
		Runtime:          &runtime,
		Overview:         &overview,
		Tagline:          &tagline,
		Status:           strPtr("Released"),
		OriginalTitle:    strPtr("Inception"),
		OriginalLanguage: strPtr("en"),
		PosterPath:       strPtr("/poster.jpg"),
		BackdropPath:     strPtr("/backdrop.jpg"),
		TrailerURL:       strPtr("https://youtube.com/watch?v=abc"),
		VoteAverage:      &voteAvg,
		VoteCount:        &runtime,
		Popularity:       &popularity,
		Budget:           &budget,
		Revenue:          &revenue,
		TitlesI18n:       map[string]string{"de": "Inception", "fr": "Inception"},
		OverviewsI18n:    map[string]string{"de": "Ein Dieb betritt Träume"},
		ExternalRatings: []ExternalRating{
			{Source: "IMDb", Value: "8.8/10", Score: 88.0},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "Inception", m.Title)
	assert.NotEqual(t, uuid.Nil, m.ID)

	// Get by ID
	got, err := repo.GetMovie(ctx, m.ID)
	require.NoError(t, err)
	assert.Equal(t, "Inception", got.Title)
	assert.Equal(t, int32(27205), *got.TMDbID)
	assert.Equal(t, "tt1375666", *got.IMDbID)
	assert.Equal(t, int32(2010), *got.Year)
	assert.Equal(t, int32(148), *got.Runtime)
	assert.Equal(t, "A thief enters dreams", *got.Overview)
	assert.NotNil(t, got.Budget)
	assert.Equal(t, int64(200000000), *got.Budget)
	assert.NotNil(t, got.Revenue)
	assert.Equal(t, int64(1000000000), *got.Revenue)
	assert.NotNil(t, got.VoteAverage)
	assert.Len(t, got.TitlesI18n, 2)
	assert.Equal(t, "Inception", got.TitlesI18n["de"])
	require.Len(t, got.ExternalRatings, 1)
	assert.Equal(t, "IMDb", got.ExternalRatings[0].Source)
}

func TestRepo_GetMovieByTMDbID(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	m := createTestMovie(t, repo, "TMDb Lookup Test")

	got, err := repo.GetMovieByTMDbID(ctx, *m.TMDbID)
	require.NoError(t, err)
	assert.Equal(t, m.ID, got.ID)
}

func TestRepo_GetMovieByIMDbID(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	imdbID := "tt9999999"
	m, err := repo.CreateMovie(ctx, CreateMovieParams{
		Title:  "IMDb Test",
		IMDbID: &imdbID,
	})
	require.NoError(t, err)

	got, err := repo.GetMovieByIMDbID(ctx, imdbID)
	require.NoError(t, err)
	assert.Equal(t, m.ID, got.ID)
}

func TestRepo_GetMovie_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)

	_, err := repo.GetMovie(context.Background(), uuid.New())
	assert.ErrorIs(t, err, ErrMovieNotFound)
}

func TestRepo_UpdateMovie(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	m := createTestMovie(t, repo, "Update Me")

	newTitle := "Updated Title"
	newOverview := "Updated overview"
	newYear := int32(2025)
	updated, err := repo.UpdateMovie(ctx, UpdateMovieParams{
		ID:       m.ID,
		Title:    &newTitle,
		Overview: &newOverview,
		Year:     &newYear,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "Updated overview", *updated.Overview)
	assert.Equal(t, int32(2025), *updated.Year)

	// Verify via fresh get
	got, err := repo.GetMovie(ctx, m.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", got.Title)
}

func TestRepo_DeleteMovie(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	m := createTestMovie(t, repo, "Delete Me")

	err := repo.DeleteMovie(ctx, m.ID)
	require.NoError(t, err)

	_, err = repo.GetMovie(ctx, m.ID)
	assert.ErrorIs(t, err, ErrMovieNotFound)
}

// ============================================================================
// Listing & Counting
// ============================================================================

func TestRepo_ListMovies(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestMovie(t, repo, "Alpha Movie")
	createTestMovie(t, repo, "Beta Movie")
	createTestMovie(t, repo, "Gamma Movie")

	movies, err := repo.ListMovies(ctx, ListFilters{
		OrderBy: "created_at",
		Limit:   10,
		Offset:  0,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(movies), 3)
}

func TestRepo_ListMovies_Pagination(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		createTestMovie(t, repo, "Page"+string(rune('A'+i)))
	}

	page1, err := repo.ListMovies(ctx, ListFilters{OrderBy: "created_at", Limit: 2, Offset: 0})
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := repo.ListMovies(ctx, ListFilters{OrderBy: "created_at", Limit: 2, Offset: 2})
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	// Different movies on each page
	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}

func TestRepo_CountMovies(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestMovie(t, repo, "Count1")
	createTestMovie(t, repo, "Count2")

	count, err := repo.CountMovies(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
}

func TestRepo_SearchMoviesByTitle(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestMovie(t, repo, "Unique Inception Search")

	results, err := repo.SearchMoviesByTitle(ctx, "Unique Inception", 10, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 1)
	assert.Contains(t, results[0].Title, "Inception")
}

func TestRepo_ListRecentlyAdded(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	createTestMovie(t, repo, "Recent Film One")
	createTestMovie(t, repo, "Recent Film Two")

	movies, err := repo.ListRecentlyAdded(ctx, 5, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(movies), 2)
}

func TestRepo_ListTopRated(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	voteAvg := "9.0"
	voteCount := int32(500)
	_, err := repo.CreateMovie(ctx, CreateMovieParams{
		Title:       "Top Rated Film",
		VoteAverage: &voteAvg,
		VoteCount:   &voteCount,
	})
	require.NoError(t, err)

	movies, err := repo.ListTopRated(ctx, 100, 10, 0)
	require.NoError(t, err)
	// At least our movie should be in top rated with 500 votes > 100 min
	found := false
	for _, m := range movies {
		if m.Title == "Top Rated Film" {
			found = true
		}
	}
	assert.True(t, found, "expected to find 'Top Rated Film' in top rated list")
}

// ============================================================================
// Movie Files
// ============================================================================

func TestRepo_MovieFiles_CRUD(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	m := createTestMovie(t, repo, "File Test Movie")

	// Create file
	f, err := repo.CreateMovieFile(ctx, CreateMovieFileParams{
		MovieID:           m.ID,
		FilePath:          "/movies/file_test.mkv",
		FileSize:          42000000000,
		FileName:          "file_test.mkv",
		Resolution:        strPtr("2160p"),
		QualityProfile:    strPtr("UHD Remux"),
		VideoCodec:        strPtr("HEVC"),
		AudioCodec:        strPtr("TrueHD Atmos"),
		Container:         strPtr("mkv"),
		BitrateKbps:       i32Ptr(35000),
		AudioLanguages:    []string{"eng", "fra"},
		SubtitleLanguages: []string{"eng", "deu"},
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, f.ID)
	assert.Equal(t, "/movies/file_test.mkv", f.FilePath)
	assert.Equal(t, "2160p", *f.Resolution)

	// Get file
	got, err := repo.GetMovieFile(ctx, f.ID)
	require.NoError(t, err)
	assert.Equal(t, f.ID, got.ID)
	assert.Equal(t, int64(42000000000), got.FileSize)

	// Get file by path
	gotByPath, err := repo.GetMovieFileByPath(ctx, "/movies/file_test.mkv")
	require.NoError(t, err)
	assert.Equal(t, f.ID, gotByPath.ID)

	// List files for movie
	files, err := repo.ListMovieFilesByMovieID(ctx, m.ID)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, f.ID, files[0].ID)

	// Delete file
	err = repo.DeleteMovieFile(ctx, f.ID)
	require.NoError(t, err)

	files, err = repo.ListMovieFilesByMovieID(ctx, m.ID)
	require.NoError(t, err)
	assert.Len(t, files, 0)
}

// ============================================================================
// Credits
// ============================================================================

func TestRepo_Credits(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	m := createTestMovie(t, repo, "Credit Test Movie")

	// Create cast credit
	cast, err := repo.CreateMovieCredit(ctx, CreateMovieCreditParams{
		MovieID:      m.ID,
		TMDbPersonID: 6193,
		Name:         "Leonardo DiCaprio",
		CreditType:   "cast",
		Character:    strPtr("Cobb"),
		CastOrder:    i32Ptr(0),
		ProfilePath:  strPtr("/leo.jpg"),
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, cast.ID)

	// Create crew credit
	_, err = repo.CreateMovieCredit(ctx, CreateMovieCreditParams{
		MovieID:      m.ID,
		TMDbPersonID: 525,
		Name:         "Christopher Nolan",
		CreditType:   "crew",
		Job:          strPtr("Director"),
		Department:   strPtr("Directing"),
	})
	require.NoError(t, err)

	// List cast
	castList, err := repo.ListMovieCast(ctx, m.ID, 10, 0)
	require.NoError(t, err)
	require.Len(t, castList, 1)
	assert.Equal(t, "Leonardo DiCaprio", castList[0].Name)
	assert.Equal(t, "Cobb", *castList[0].Character)

	// Count cast
	castCount, err := repo.CountMovieCast(ctx, m.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), castCount)

	// List crew
	crewList, err := repo.ListMovieCrew(ctx, m.ID, 10, 0)
	require.NoError(t, err)
	require.Len(t, crewList, 1)
	assert.Equal(t, "Christopher Nolan", crewList[0].Name)

	// Count crew
	crewCount, err := repo.CountMovieCrew(ctx, m.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), crewCount)

	// Delete all credits
	err = repo.DeleteMovieCredits(ctx, m.ID)
	require.NoError(t, err)

	castList, err = repo.ListMovieCast(ctx, m.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, castList, 0)
}

// ============================================================================
// Collections
// ============================================================================

func TestRepo_Collections(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create collection
	tmdbCollID := int32(263)
	coll, err := repo.CreateMovieCollection(ctx, CreateMovieCollectionParams{
		TMDbCollectionID: &tmdbCollID,
		Name:             "The Dark Knight Collection",
		Overview:         strPtr("Batman trilogy"),
		PosterPath:       strPtr("/collection.jpg"),
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, coll.ID)

	// Get collection
	got, err := repo.GetMovieCollection(ctx, coll.ID)
	require.NoError(t, err)
	assert.Equal(t, "The Dark Knight Collection", got.Name)

	// Get by TMDb ID
	gotByTMDb, err := repo.GetMovieCollectionByTMDbID(ctx, 263)
	require.NoError(t, err)
	assert.Equal(t, coll.ID, gotByTMDb.ID)

	// Create movies and add to collection
	m1 := createTestMovie(t, repo, "Batman Begins")
	m2 := createTestMovie(t, repo, "The Dark Knight")

	err = repo.AddMovieToCollection(ctx, coll.ID, m1.ID, i32Ptr(1))
	require.NoError(t, err)
	err = repo.AddMovieToCollection(ctx, coll.ID, m2.ID, i32Ptr(2))
	require.NoError(t, err)

	// List movies in collection
	collMovies, err := repo.ListMoviesByCollection(ctx, coll.ID)
	require.NoError(t, err)
	assert.Len(t, collMovies, 2)

	// Get collection for movie
	movieColl, err := repo.GetCollectionForMovie(ctx, m1.ID)
	require.NoError(t, err)
	assert.Equal(t, coll.ID, movieColl.ID)

	// Remove movie from collection
	err = repo.RemoveMovieFromCollection(ctx, coll.ID, m1.ID)
	require.NoError(t, err)

	collMovies, err = repo.ListMoviesByCollection(ctx, coll.ID)
	require.NoError(t, err)
	assert.Len(t, collMovies, 1)
}

// ============================================================================
// Genres
// ============================================================================

func TestRepo_Genres(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	m := createTestMovie(t, repo, "Genre Test Movie")

	// Add genres
	err := repo.AddMovieGenre(ctx, m.ID, "action", "Action")
	require.NoError(t, err)
	err = repo.AddMovieGenre(ctx, m.ID, "science-fiction", "Science Fiction")
	require.NoError(t, err)

	// List genres for movie
	genres, err := repo.ListMovieGenres(ctx, m.ID)
	require.NoError(t, err)
	assert.Len(t, genres, 2)

	genreNames := make([]string, len(genres))
	for i, g := range genres {
		genreNames[i] = g.Name
	}
	assert.Contains(t, genreNames, "Action")
	assert.Contains(t, genreNames, "Science Fiction")

	// List distinct genres (all movies)
	distinct, err := repo.ListDistinctMovieGenres(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(distinct), 2)

	// List movies by genre
	genreMovies, err := repo.ListMoviesByGenre(ctx, "action", 10, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(genreMovies), 1)
	assert.Equal(t, m.ID, genreMovies[0].ID)

	// Delete genres
	err = repo.DeleteMovieGenres(ctx, m.ID)
	require.NoError(t, err)

	genres, err = repo.ListMovieGenres(ctx, m.ID)
	require.NoError(t, err)
	assert.Len(t, genres, 0)
}

// ============================================================================
// Watch Progress
// ============================================================================

func TestRepo_WatchProgress(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	m := createTestMovie(t, repo, "Watch Progress Film")
	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "watcher",
		Email:    "watcher@example.com",
	})

	// Create watch progress
	watched, err := repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		MovieID:         m.ID,
		ProgressSeconds: 3600,
		DurationSeconds: 7200,
		IsCompleted:     false,
	})
	require.NoError(t, err)
	assert.Equal(t, int32(3600), watched.ProgressSeconds)

	// Get progress
	got, err := repo.GetWatchProgress(ctx, user.ID, m.ID)
	require.NoError(t, err)
	assert.Equal(t, int32(3600), got.ProgressSeconds)

	// Update progress
	updated, err := repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		MovieID:         m.ID,
		ProgressSeconds: 5400,
		DurationSeconds: 7200,
		IsCompleted:     false,
	})
	require.NoError(t, err)
	assert.Equal(t, int32(5400), updated.ProgressSeconds)

	// Delete progress
	err = repo.DeleteWatchProgress(ctx, user.ID, m.ID)
	require.NoError(t, err)

	_, err = repo.GetWatchProgress(ctx, user.ID, m.ID)
	assert.Error(t, err)
}

func TestRepo_ContinueWatching(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "cwatcher",
		Email:    "cwatcher@example.com",
	})

	m1 := createTestMovie(t, repo, "Continue Film One")
	m2 := createTestMovie(t, repo, "Continue Film Two")

	// Create in-progress watches
	_, err := repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		MovieID:         m1.ID,
		ProgressSeconds: 1800,
		DurationSeconds: 7200,
		IsCompleted:     false,
	})
	require.NoError(t, err)

	_, err = repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		MovieID:         m2.ID,
		ProgressSeconds: 900,
		DurationSeconds: 5400,
		IsCompleted:     false,
	})
	require.NoError(t, err)

	items, err := repo.ListContinueWatching(ctx, user.ID, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 2)

	// Verify items have movie data + progress
	for _, item := range items {
		assert.NotEmpty(t, item.Title)
		assert.Greater(t, item.ProgressSeconds, int32(0))
		assert.Greater(t, item.DurationSeconds, int32(0))
	}
}

func TestRepo_WatchHistory(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "historian",
		Email:    "historian@example.com",
	})

	m := createTestMovie(t, repo, "History Film")

	// Mark as completed
	_, err := repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		MovieID:         m.ID,
		ProgressSeconds: 7200,
		DurationSeconds: 7200,
		IsCompleted:     true,
	})
	require.NoError(t, err)

	items, err := repo.ListWatchedMovies(ctx, user.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 1)

	found := false
	for _, item := range items {
		if item.Title == "History Film" {
			found = true
			assert.Greater(t, item.WatchCount, int32(0))
		}
	}
	assert.True(t, found, "expected History Film in watch history")
}

func TestRepo_UserMovieStats(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepo(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "statuser",
		Email:    "statuser@example.com",
	})

	m1 := createTestMovie(t, repo, "Stat Film One")
	m2 := createTestMovie(t, repo, "Stat Film Two")

	// One completed, one in-progress
	_, err := repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		MovieID:         m1.ID,
		ProgressSeconds: 7200,
		DurationSeconds: 7200,
		IsCompleted:     true,
	})
	require.NoError(t, err)

	_, err = repo.CreateOrUpdateWatchProgress(ctx, CreateWatchProgressParams{
		UserID:          user.ID,
		MovieID:         m2.ID,
		ProgressSeconds: 1000,
		DurationSeconds: 5400,
		IsCompleted:     false,
	})
	require.NoError(t, err)

	stats, err := repo.GetUserMovieStats(ctx, user.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.WatchedCount, int64(1))
	assert.GreaterOrEqual(t, stats.InProgressCount, int64(1))
}

// ============================================================================
// Edge cases
// ============================================================================

func TestRepo_CreateMovie_Minimal(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)

	m, err := repo.CreateMovie(context.Background(), CreateMovieParams{
		Title: "Minimal Movie",
	})
	require.NoError(t, err)
	assert.Equal(t, "Minimal Movie", m.Title)
	assert.Nil(t, m.TMDbID)
	assert.Nil(t, m.IMDbID)
	assert.Nil(t, m.Year)
	assert.Nil(t, m.Runtime)
}

func TestRepo_UpdateMovie_PartialFields(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	m := createTestMovie(t, repo, "Partial Update")

	// Only update title, everything else stays
	newTitle := "New Partial Title"
	updated, err := repo.UpdateMovie(ctx, UpdateMovieParams{
		ID:    m.ID,
		Title: &newTitle,
	})
	require.NoError(t, err)
	assert.Equal(t, "New Partial Title", updated.Title)
	assert.Equal(t, *m.Year, *updated.Year)       // unchanged
	assert.Equal(t, *m.Overview, *updated.Overview) // unchanged
}

func TestRepo_CreateMovie_WithI18nAndAgeRatings(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	m, err := repo.CreateMovie(ctx, CreateMovieParams{
		Title:      "I18n Movie",
		TitlesI18n: map[string]string{"de": "I18n Film", "fr": "Film I18n", "ja": "I18n映画"},
		AgeRatings: map[string]map[string]string{
			"US":  {"rating": "PG-13", "source": "MPAA"},
			"DE":  {"rating": "FSK 12", "source": "FSK"},
		},
	})
	require.NoError(t, err)

	got, err := repo.GetMovie(ctx, m.ID)
	require.NoError(t, err)
	assert.Len(t, got.TitlesI18n, 3)
	assert.Equal(t, "I18n Film", got.TitlesI18n["de"])
	assert.Len(t, got.AgeRatings, 2)
	assert.Equal(t, "PG-13", got.AgeRatings["US"]["rating"])
}
