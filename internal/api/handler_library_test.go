package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/testutil"
)

// TestHandler_ListLibraries_NoAuth verifies that ListLibraries returns 401
// when no user is present in the context.
func TestHandler_ListLibraries_NoAuth(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()

	result, err := handler.ListLibraries(ctx)
	require.NoError(t, err)

	errResp, ok := result.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", result)
	assert.Equal(t, 401, errResp.Code)
	assert.Equal(t, "Authentication required", errResp.Message)
}

// TestHandler_CreateLibrary_NotAdmin verifies that CreateLibrary returns 401
// when no user is in the context (unauthenticated).
func TestHandler_CreateLibrary_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.CreateLibraryRequest{}

	result, err := handler.CreateLibrary(ctx, req)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.CreateLibraryUnauthorized)
	require.True(t, ok, "expected *ogen.CreateLibraryUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_GetLibrary_NoAuth verifies that GetLibrary returns 401
// when no user is present in the context.
func TestHandler_GetLibrary_NoAuth(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.GetLibraryParams{LibraryID: uuid.New()}

	result, err := handler.GetLibrary(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.GetLibraryUnauthorized)
	require.True(t, ok, "expected *ogen.GetLibraryUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_UpdateLibrary_NotAdmin verifies that UpdateLibrary returns 401
// when no user is in the context (unauthenticated).
func TestHandler_UpdateLibrary_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.UpdateLibraryRequest{}
	params := ogen.UpdateLibraryParams{LibraryID: uuid.New()}

	result, err := handler.UpdateLibrary(ctx, req, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.UpdateLibraryUnauthorized)
	require.True(t, ok, "expected *ogen.UpdateLibraryUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_DeleteLibrary_NotAdmin verifies that DeleteLibrary returns 401
// when no user is in the context (unauthenticated).
func TestHandler_DeleteLibrary_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.DeleteLibraryParams{LibraryID: uuid.New()}

	result, err := handler.DeleteLibrary(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.DeleteLibraryUnauthorized)
	require.True(t, ok, "expected *ogen.DeleteLibraryUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_TriggerLibraryScan_NotAdmin verifies that TriggerLibraryScan returns 401
// when no user is in the context (unauthenticated).
func TestHandler_TriggerLibraryScan_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.TriggerLibraryScanReq{ScanType: ogen.TriggerLibraryScanReqScanType("full")}
	params := ogen.TriggerLibraryScanParams{LibraryID: uuid.New()}

	result, err := handler.TriggerLibraryScan(ctx, req, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.TriggerLibraryScanUnauthorized)
	require.True(t, ok, "expected *ogen.TriggerLibraryScanUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_ListLibraryScans_NoAuth verifies that ListLibraryScans returns 401
// when no user is present in the context.
func TestHandler_ListLibraryScans_NoAuth(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.ListLibraryScansParams{LibraryID: uuid.New()}

	result, err := handler.ListLibraryScans(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.ListLibraryScansUnauthorized)
	require.True(t, ok, "expected *ogen.ListLibraryScansUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_ListLibraryPermissions_NotAdmin verifies that ListLibraryPermissions returns 401
// when no user is in the context (unauthenticated).
func TestHandler_ListLibraryPermissions_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.ListLibraryPermissionsParams{LibraryID: uuid.New()}

	result, err := handler.ListLibraryPermissions(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.ListLibraryPermissionsUnauthorized)
	require.True(t, ok, "expected *ogen.ListLibraryPermissionsUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_GrantLibraryPermission_NotAdmin verifies that GrantLibraryPermission returns 401
// when no user is in the context (unauthenticated).
func TestHandler_GrantLibraryPermission_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.GrantLibraryPermissionReq{
		UserID:     uuid.New(),
		Permission: ogen.GrantLibraryPermissionReqPermission("read"),
	}
	params := ogen.GrantLibraryPermissionParams{LibraryID: uuid.New()}

	result, err := handler.GrantLibraryPermission(ctx, req, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.GrantLibraryPermissionUnauthorized)
	require.True(t, ok, "expected *ogen.GrantLibraryPermissionUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_RevokeLibraryPermission_NotAdmin verifies that RevokeLibraryPermission returns 401
// when no user is in the context (unauthenticated).
func TestHandler_RevokeLibraryPermission_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.RevokeLibraryPermissionParams{
		LibraryID:  uuid.New(),
		UserID:     uuid.New(),
		Permission: ogen.RevokeLibraryPermissionPermission("read"),
	}

	result, err := handler.RevokeLibraryPermission(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.RevokeLibraryPermissionUnauthorized)
	require.True(t, ok, "expected *ogen.RevokeLibraryPermissionUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// ============================================================================
// Success Path Tests (with real database + library service)
// ============================================================================

// setupLibraryTestHandler creates a Handler with a real library service backed by
// an embedded PostgreSQL test database. It returns the handler, the test DB, and
// the admin user's UUID. The admin user has the "admin" RBAC role assigned.
func setupLibraryTestHandler(t *testing.T) (*Handler, testutil.DB, uuid.UUID) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)

	// Clear any existing policies from the table to ensure test isolation
	_, err := testDB.Pool().Exec(context.Background(), "DELETE FROM shared.casbin_rule")
	require.NoError(t, err)

	// Set up RBAC service with Casbin
	adapter := rbac.NewAdapter(testDB.Pool())
	modelPath := "../../config/casbin_model.conf"
	enforcer, err := casbin.NewSyncedEnforcer(modelPath, adapter)
	require.NoError(t, err)
	rbacService := rbac.NewService(enforcer, logging.NewTestLogger(), activity.NewNoopLogger())

	// Create admin user
	adminUser := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "libadmin",
		Email:    "libadmin@example.com",
	})

	// Grant admin role
	err = rbacService.AssignRole(context.Background(), adminUser.ID, "admin")
	require.NoError(t, err)

	// Create library service with real repository
	queries := db.New(testDB.Pool())
	repo := library.NewRepositoryPg(queries)
	libService := library.NewService(repo, logging.NewTestLogger(), activity.NewNoopLogger())

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTExpiry: 15 * time.Minute,
		},
	}

	handler := &Handler{
		logger:         logging.NewTestLogger(),
		rbacService:    rbacService,
		cfg:            cfg,
		libraryService: libService,
	}

	return handler, testDB, adminUser.ID
}

// TestHandler_ListLibraries_AdminSuccess verifies that an admin can list libraries.
// It first checks that an empty database returns an empty list, then creates a
// library and verifies it appears in the list.
func TestHandler_ListLibraries_AdminSuccess(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupLibraryTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	// List with no libraries should return empty list
	result, err := handler.ListLibraries(ctx)
	require.NoError(t, err)

	listResp, ok := result.(*ogen.LibraryListResponse)
	require.True(t, ok, "expected *ogen.LibraryListResponse, got %T", result)
	assert.Equal(t, int64(0), listResp.Total)
	assert.Empty(t, listResp.Libraries)

	// Create a library so we have something to list
	createReq := &ogen.CreateLibraryRequest{
		Name:  "Movies",
		Type:  ogen.CreateLibraryRequestTypeMovie,
		Paths: []string{"/media/movies"},
	}
	_, err = handler.CreateLibrary(ctx, createReq)
	require.NoError(t, err)

	// List again and verify the library is present
	result, err = handler.ListLibraries(ctx)
	require.NoError(t, err)

	listResp, ok = result.(*ogen.LibraryListResponse)
	require.True(t, ok, "expected *ogen.LibraryListResponse, got %T", result)
	assert.Equal(t, int64(1), listResp.Total)
	require.Len(t, listResp.Libraries, 1)
	assert.Equal(t, "Movies", listResp.Libraries[0].Name)
	assert.Equal(t, ogen.LibraryType("movie"), listResp.Libraries[0].Type)
	assert.Equal(t, []string{"/media/movies"}, listResp.Libraries[0].Paths)
}

// TestHandler_CreateLibrary_AdminSuccess verifies that an admin can create a
// library and receives the created library details in the response.
func TestHandler_CreateLibrary_AdminSuccess(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupLibraryTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	req := &ogen.CreateLibraryRequest{
		Name:  "TV Shows",
		Type:  ogen.CreateLibraryRequestTypeTvshow,
		Paths: []string{"/media/tv", "/media/tv2"},
	}

	result, err := handler.CreateLibrary(ctx, req)
	require.NoError(t, err)

	lib, ok := result.(*ogen.Library)
	require.True(t, ok, "expected *ogen.Library, got %T", result)
	assert.NotEqual(t, uuid.Nil, lib.ID)
	assert.Equal(t, "TV Shows", lib.Name)
	assert.Equal(t, ogen.LibraryType("tvshow"), lib.Type)
	assert.Equal(t, []string{"/media/tv", "/media/tv2"}, lib.Paths)
	assert.True(t, lib.Enabled, "library should be enabled by default")
	assert.False(t, lib.CreatedAt.IsZero())
	assert.False(t, lib.UpdatedAt.IsZero())
}

// TestHandler_CreateLibrary_InvalidType verifies that creating a library with an
// invalid type returns a 400 Bad Request response.
func TestHandler_CreateLibrary_InvalidType(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupLibraryTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	req := &ogen.CreateLibraryRequest{
		Name:  "Bad Library",
		Type:  ogen.CreateLibraryRequestType("invalid_type"),
		Paths: []string{"/media/bad"},
	}

	result, err := handler.CreateLibrary(ctx, req)
	require.NoError(t, err)

	badReq, ok := result.(*ogen.CreateLibraryBadRequest)
	require.True(t, ok, "expected *ogen.CreateLibraryBadRequest, got %T", result)
	assert.Equal(t, 400, badReq.Code)
	assert.Equal(t, "Invalid library type", badReq.Message)
}

// TestHandler_GetLibrary_AdminSuccess verifies that an admin can retrieve a
// specific library by its ID.
func TestHandler_GetLibrary_AdminSuccess(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupLibraryTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	// First create a library
	createReq := &ogen.CreateLibraryRequest{
		Name:  "Music Collection",
		Type:  ogen.CreateLibraryRequestTypeMusic,
		Paths: []string{"/media/music"},
	}

	createResult, err := handler.CreateLibrary(ctx, createReq)
	require.NoError(t, err)
	created, ok := createResult.(*ogen.Library)
	require.True(t, ok, "expected *ogen.Library from create, got %T", createResult)

	// Now get it by ID
	getParams := ogen.GetLibraryParams{LibraryID: created.ID}
	getResult, err := handler.GetLibrary(ctx, getParams)
	require.NoError(t, err)

	lib, ok := getResult.(*ogen.Library)
	require.True(t, ok, "expected *ogen.Library, got %T", getResult)
	assert.Equal(t, created.ID, lib.ID)
	assert.Equal(t, "Music Collection", lib.Name)
	assert.Equal(t, ogen.LibraryType("music"), lib.Type)
	assert.Equal(t, []string{"/media/music"}, lib.Paths)
}

// TestHandler_GetLibrary_NotFound verifies that requesting a non-existent library
// returns a 404 Not Found response.
func TestHandler_GetLibrary_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupLibraryTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	params := ogen.GetLibraryParams{LibraryID: uuid.New()}
	result, err := handler.GetLibrary(ctx, params)
	require.NoError(t, err)

	notFound, ok := result.(*ogen.GetLibraryNotFound)
	require.True(t, ok, "expected *ogen.GetLibraryNotFound, got %T", result)
	assert.Equal(t, 404, notFound.Code)
	assert.Equal(t, "Library not found", notFound.Message)
}

// TestHandler_DeleteLibrary_AdminSuccess verifies that an admin can delete a
// library and that it is no longer retrievable afterwards.
func TestHandler_DeleteLibrary_AdminSuccess(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupLibraryTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	// Create a library to delete
	createReq := &ogen.CreateLibraryRequest{
		Name:  "Temp Library",
		Type:  ogen.CreateLibraryRequestTypeMovie,
		Paths: []string{"/media/temp"},
	}
	createResult, err := handler.CreateLibrary(ctx, createReq)
	require.NoError(t, err)
	created, ok := createResult.(*ogen.Library)
	require.True(t, ok, "expected *ogen.Library from create, got %T", createResult)

	// Delete it
	deleteParams := ogen.DeleteLibraryParams{LibraryID: created.ID}
	deleteResult, err := handler.DeleteLibrary(ctx, deleteParams)
	require.NoError(t, err)

	_, ok = deleteResult.(*ogen.DeleteLibraryNoContent)
	require.True(t, ok, "expected *ogen.DeleteLibraryNoContent, got %T", deleteResult)

	// Verify it no longer exists
	getParams := ogen.GetLibraryParams{LibraryID: created.ID}
	getResult, err := handler.GetLibrary(ctx, getParams)
	require.NoError(t, err)

	notFound, ok := getResult.(*ogen.GetLibraryNotFound)
	require.True(t, ok, "expected *ogen.GetLibraryNotFound after deletion, got %T", getResult)
	assert.Equal(t, 404, notFound.Code)
}

// ============================================================================
// ListGenres tests
// ============================================================================

// mockMovieService is a minimal mock for movie.Service used by ListGenres tests.
type mockMovieService struct {
	movie.Service // embed to satisfy interface; only override what's needed
	genres        []content.GenreSummary
	err           error
}

func (m *mockMovieService) ListDistinctGenres(_ context.Context) ([]content.GenreSummary, error) {
	return m.genres, m.err
}

// mockTVService is a minimal mock for tvshow.Service used by ListGenres tests.
type mockTVService struct {
	tvshow.Service // embed to satisfy interface; only override what's needed
	genres         []content.GenreSummary
	err            error
}

func (m *mockTVService) ListDistinctGenres(_ context.Context) ([]content.GenreSummary, error) {
	return m.genres, m.err
}

func TestHandler_ListGenres_Success(t *testing.T) {
	t.Parallel()

	movieSvc := &mockMovieService{
		genres: []content.GenreSummary{
			{TMDbGenreID: 28, Name: "Action", ItemCount: 10},
			{TMDbGenreID: 18, Name: "Drama", ItemCount: 25},
			{TMDbGenreID: 35, Name: "Comedy", ItemCount: 5},
		},
	}
	tvSvc := &mockTVService{
		genres: []content.GenreSummary{
			{TMDbGenreID: 18, Name: "Drama", ItemCount: 15},
			{TMDbGenreID: 10765, Name: "Sci-Fi & Fantasy", ItemCount: 8},
		},
	}

	handler := &Handler{
		logger:        logging.NewTestLogger(),
		movieHandler:  movie.NewHandler(movieSvc, nil),
		tvshowService: tvSvc,
	}

	result, err := handler.ListGenres(context.Background())
	require.NoError(t, err)

	genres, ok := result.(*ogen.ListGenresOKApplicationJSON)
	require.True(t, ok, "expected *ogen.ListGenresOKApplicationJSON, got %T", result)

	// Should have 4 distinct genres: Action, Comedy, Drama, Sci-Fi & Fantasy
	require.Len(t, *genres, 4)

	// Verify alphabetical sort
	assert.Equal(t, "Action", (*genres)[0].Name)
	assert.Equal(t, "Comedy", (*genres)[1].Name)
	assert.Equal(t, "Drama", (*genres)[2].Name)
	assert.Equal(t, "Sci-Fi & Fantasy", (*genres)[3].Name)

	// Verify Drama is merged: 25 movies + 15 TV shows
	drama := (*genres)[2]
	assert.Equal(t, 28, (*genres)[0].TmdbGenreID) // Action TMDb ID
	assert.Equal(t, int64(25), drama.MovieCount)
	assert.Equal(t, int64(15), drama.TvshowCount)

	// Verify Action: movies only
	assert.Equal(t, int64(10), (*genres)[0].MovieCount)
	assert.Equal(t, int64(0), (*genres)[0].TvshowCount)

	// Verify Sci-Fi & Fantasy: TV only
	assert.Equal(t, int64(0), (*genres)[3].MovieCount)
	assert.Equal(t, int64(8), (*genres)[3].TvshowCount)
}

func TestHandler_ListGenres_EmptyDB(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:        logging.NewTestLogger(),
		movieHandler:  movie.NewHandler(&mockMovieService{}, nil),
		tvshowService: &mockTVService{},
	}

	result, err := handler.ListGenres(context.Background())
	require.NoError(t, err)

	genres, ok := result.(*ogen.ListGenresOKApplicationJSON)
	require.True(t, ok)
	assert.Empty(t, *genres)
}

func TestHandler_ListGenres_MovieError(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
		movieHandler: movie.NewHandler(&mockMovieService{
			err: errors.New("movie db error"),
		}, nil),
		tvshowService: &mockTVService{},
	}

	result, err := handler.ListGenres(context.Background())
	require.NoError(t, err)

	errResp, ok := result.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", result)
	assert.Equal(t, 500, errResp.Code)
}

func TestHandler_ListGenres_TVShowError(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
		movieHandler: movie.NewHandler(&mockMovieService{
			genres: []content.GenreSummary{{TMDbGenreID: 28, Name: "Action", ItemCount: 5}},
		}, nil),
		tvshowService: &mockTVService{err: errors.New("tv db error")},
	}

	result, err := handler.ListGenres(context.Background())
	require.NoError(t, err)

	errResp, ok := result.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", result)
	assert.Equal(t, 500, errResp.Code)
}
