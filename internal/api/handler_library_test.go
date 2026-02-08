package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/lusoris/revenge/internal/infra/logging"

	"github.com/lusoris/revenge/internal/api/ogen"
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

// TestHandler_CreateLibrary_NotAdmin verifies that CreateLibrary returns 403
// when the caller is not an admin (no user in context means isAdmin returns false).
func TestHandler_CreateLibrary_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.CreateLibraryRequest{}

	result, err := handler.CreateLibrary(ctx, req)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.CreateLibraryForbidden)
	require.True(t, ok, "expected *ogen.CreateLibraryForbidden, got %T", result)
	assert.Equal(t, 403, forbidden.Code)
	assert.Equal(t, "Admin access required", forbidden.Message)
}

// TestHandler_GetLibrary_NoAuth verifies that GetLibrary returns 401
// when no user is present in the context.
func TestHandler_GetLibrary_NoAuth(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.GetLibraryParams{LibraryId: uuid.New()}

	result, err := handler.GetLibrary(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.GetLibraryUnauthorized)
	require.True(t, ok, "expected *ogen.GetLibraryUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_UpdateLibrary_NotAdmin verifies that UpdateLibrary returns 403
// when the caller is not an admin.
func TestHandler_UpdateLibrary_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.UpdateLibraryRequest{}
	params := ogen.UpdateLibraryParams{LibraryId: uuid.New()}

	result, err := handler.UpdateLibrary(ctx, req, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.UpdateLibraryForbidden)
	require.True(t, ok, "expected *ogen.UpdateLibraryForbidden, got %T", result)
	assert.Equal(t, 403, forbidden.Code)
	assert.Equal(t, "Admin access required", forbidden.Message)
}

// TestHandler_DeleteLibrary_NotAdmin verifies that DeleteLibrary returns 403
// when the caller is not an admin.
func TestHandler_DeleteLibrary_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.DeleteLibraryParams{LibraryId: uuid.New()}

	result, err := handler.DeleteLibrary(ctx, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.DeleteLibraryForbidden)
	require.True(t, ok, "expected *ogen.DeleteLibraryForbidden, got %T", result)
	assert.Equal(t, 403, forbidden.Code)
	assert.Equal(t, "Admin access required", forbidden.Message)
}

// TestHandler_TriggerLibraryScan_NotAdmin verifies that TriggerLibraryScan returns 403
// when the caller is not an admin.
func TestHandler_TriggerLibraryScan_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.TriggerLibraryScanReq{ScanType: ogen.TriggerLibraryScanReqScanType("full")}
	params := ogen.TriggerLibraryScanParams{LibraryId: uuid.New()}

	result, err := handler.TriggerLibraryScan(ctx, req, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.TriggerLibraryScanForbidden)
	require.True(t, ok, "expected *ogen.TriggerLibraryScanForbidden, got %T", result)
	assert.Equal(t, 403, forbidden.Code)
	assert.Equal(t, "Admin access required", forbidden.Message)
}

// TestHandler_ListLibraryScans_NoAuth verifies that ListLibraryScans returns 401
// when no user is present in the context.
func TestHandler_ListLibraryScans_NoAuth(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.ListLibraryScansParams{LibraryId: uuid.New()}

	result, err := handler.ListLibraryScans(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.ListLibraryScansUnauthorized)
	require.True(t, ok, "expected *ogen.ListLibraryScansUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Equal(t, "Authentication required", unauthorized.Message)
}

// TestHandler_ListLibraryPermissions_NotAdmin verifies that ListLibraryPermissions returns 403
// when the caller is not an admin.
func TestHandler_ListLibraryPermissions_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.ListLibraryPermissionsParams{LibraryId: uuid.New()}

	result, err := handler.ListLibraryPermissions(ctx, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.ListLibraryPermissionsForbidden)
	require.True(t, ok, "expected *ogen.ListLibraryPermissionsForbidden, got %T", result)
	assert.Equal(t, 403, forbidden.Code)
	assert.Equal(t, "Admin access required", forbidden.Message)
}

// TestHandler_GrantLibraryPermission_NotAdmin verifies that GrantLibraryPermission returns 403
// when the caller is not an admin.
func TestHandler_GrantLibraryPermission_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.GrantLibraryPermissionReq{
		UserId:     uuid.New(),
		Permission: ogen.GrantLibraryPermissionReqPermission("read"),
	}
	params := ogen.GrantLibraryPermissionParams{LibraryId: uuid.New()}

	result, err := handler.GrantLibraryPermission(ctx, req, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.GrantLibraryPermissionForbidden)
	require.True(t, ok, "expected *ogen.GrantLibraryPermissionForbidden, got %T", result)
	assert.Equal(t, 403, forbidden.Code)
	assert.Equal(t, "Admin access required", forbidden.Message)
}

// TestHandler_RevokeLibraryPermission_NotAdmin verifies that RevokeLibraryPermission returns 403
// when the caller is not an admin.
func TestHandler_RevokeLibraryPermission_NotAdmin(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.RevokeLibraryPermissionParams{
		LibraryId:  uuid.New(),
		UserId:     uuid.New(),
		Permission: ogen.RevokeLibraryPermissionPermission("read"),
	}

	result, err := handler.RevokeLibraryPermission(ctx, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.RevokeLibraryPermissionForbidden)
	require.True(t, ok, "expected *ogen.RevokeLibraryPermissionForbidden, got %T", result)
	assert.Equal(t, 403, forbidden.Code)
	assert.Equal(t, "Admin access required", forbidden.Message)
}
