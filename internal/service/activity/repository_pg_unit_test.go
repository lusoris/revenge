package activity

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// errDBTX is a mock DBTX that always returns errors.
type errDBTX struct {
	err error
}

func (e *errDBTX) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, e.err
}

func (e *errDBTX) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return nil, e.err
}

func (e *errDBTX) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return &errRow{err: e.err}
}

func (e *errDBTX) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, e.err
}

// errRow is a mock pgx.Row that always returns an error on Scan.
type errRow struct {
	err error
}

func (r *errRow) Scan(dest ...any) error {
	return r.err
}

// noRowsDBTX is a mock DBTX that returns pgx.ErrNoRows from QueryRow.
type noRowsDBTX struct{}

func (n *noRowsDBTX) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (n *noRowsDBTX) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return &emptyRows{}, nil
}

func (n *noRowsDBTX) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return &errRow{err: pgx.ErrNoRows}
}

func (n *noRowsDBTX) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

// emptyRows is a mock pgx.Rows that has no data.
type emptyRows struct{}

func (r *emptyRows) Close()                                        {}
func (r *emptyRows) Err() error                                    { return nil }
func (r *emptyRows) CommandTag() pgconn.CommandTag                 { return pgconn.CommandTag{} }
func (r *emptyRows) FieldDescriptions() []pgconn.FieldDescription  { return nil }
func (r *emptyRows) Next() bool                                    { return false }
func (r *emptyRows) Scan(dest ...any) error                        { return errors.New("no rows") }
func (r *emptyRows) Values() ([]any, error)                        { return nil, nil }
func (r *emptyRows) RawValues() [][]byte                           { return nil }
func (r *emptyRows) Conn() *pgx.Conn                               { return nil }

// successExecDBTX returns successful results for Exec operations.
type successExecDBTX struct {
	rowsAffected int64
}

func (s *successExecDBTX) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag("DELETE " + intToStr(s.rowsAffected)), nil
}

func (s *successExecDBTX) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return &emptyRows{}, nil
}

func (s *successExecDBTX) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return &errRow{err: pgx.ErrNoRows}
}

func (s *successExecDBTX) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func intToStr(n int64) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

// ============================================================================
// RepositoryPg Error Path Tests
// ============================================================================

func TestRepositoryPg_Create_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("connection refused")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	entry := &Entry{
		Action:  ActionUserLogin,
		Success: true,
	}

	err := repo.Create(context.Background(), entry)
	assert.Error(t, err)
	assert.Equal(t, dbErr, err)
}

func TestRepositoryPg_Create_WithChangesAndMetadata(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("insert failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	userID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())
	ip := net.ParseIP("10.0.0.1")

	entry := &Entry{
		UserID:       &userID,
		Username:     ptrStr("testuser"),
		Action:       ActionUserLogin,
		ResourceType: ptrStr(ResourceTypeUser),
		ResourceID:   &resourceID,
		Changes:      map[string]interface{}{"key": "value"},
		Metadata:     map[string]interface{}{"meta": "data"},
		IPAddress:    &ip,
		UserAgent:    ptrStr("TestAgent"),
		Success:      true,
	}

	err := repo.Create(context.Background(), entry)
	assert.Error(t, err)
}

func TestRepositoryPg_Get_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	entry, err := repo.Get(context.Background(), uuid.Must(uuid.NewV7()))
	assert.Nil(t, entry)
	assert.Error(t, err)
	assert.Equal(t, dbErr, err)
}

func TestRepositoryPg_Get_NotFound_Unit(t *testing.T) {
	t.Parallel()
	queries := db.New(&noRowsDBTX{})
	repo := NewRepositoryPg(queries)

	entry, err := repo.Get(context.Background(), uuid.Must(uuid.NewV7()))
	assert.Nil(t, entry)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRepositoryPg_List_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("list query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	entries, err := repo.List(context.Background(), 10, 0)
	assert.Nil(t, entries)
	assert.Error(t, err)
}

func TestRepositoryPg_List_Empty(t *testing.T) {
	t.Parallel()
	queries := db.New(&noRowsDBTX{})
	repo := NewRepositoryPg(queries)

	entries, err := repo.List(context.Background(), 10, 0)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestRepositoryPg_Count_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("count query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	count, err := repo.Count(context.Background())
	assert.Equal(t, int64(0), count)
	assert.Error(t, err)
}

func TestRepositoryPg_Search_CountError(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("search count failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	filters := SearchFilters{Limit: 10}
	entries, count, err := repo.Search(context.Background(), filters)
	assert.Nil(t, entries)
	assert.Equal(t, int64(0), count)
	assert.Error(t, err)
}

func TestRepositoryPg_Search_WithFilters(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("search failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	userID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())
	action := ActionUserLogin
	resourceType := ResourceTypeUser
	success := true
	now := time.Now()

	filters := SearchFilters{
		UserID:       &userID,
		Action:       &action,
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
		Success:      &success,
		StartTime:    &now,
		EndTime:      &now,
		Limit:        10,
		Offset:       0,
	}
	entries, count, err := repo.Search(context.Background(), filters)
	assert.Nil(t, entries)
	assert.Equal(t, int64(0), count)
	assert.Error(t, err)
}

func TestRepositoryPg_GetByUser_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("user query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	entries, count, err := repo.GetByUser(context.Background(), uuid.Must(uuid.NewV7()), 10, 0)
	assert.Nil(t, entries)
	assert.Equal(t, int64(0), count)
	assert.Error(t, err)
}

func TestRepositoryPg_GetByResource_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("resource query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	entries, count, err := repo.GetByResource(context.Background(), ResourceTypeUser, uuid.Must(uuid.NewV7()), 10, 0)
	assert.Nil(t, entries)
	assert.Equal(t, int64(0), count)
	assert.Error(t, err)
}

func TestRepositoryPg_GetByAction_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("action query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	entries, err := repo.GetByAction(context.Background(), ActionUserLogin, 10, 0)
	assert.Nil(t, entries)
	assert.Error(t, err)
}

func TestRepositoryPg_GetByIP_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("ip query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	ip := net.ParseIP("192.168.1.1")
	entries, err := repo.GetByIP(context.Background(), ip, 10, 0)
	assert.Nil(t, entries)
	assert.Error(t, err)
}

func TestRepositoryPg_GetByIP_InvalidIP(t *testing.T) {
	t.Parallel()
	// An invalid/nil IP should return nil, nil
	queries := db.New(&noRowsDBTX{})
	repo := NewRepositoryPg(queries)

	entries, err := repo.GetByIP(context.Background(), nil, 10, 0)
	assert.Nil(t, entries)
	assert.NoError(t, err)
}

func TestRepositoryPg_GetByIP_ValidIP(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("ip query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	ip := net.ParseIP("2001:db8::1")
	entries, err := repo.GetByIP(context.Background(), ip, 10, 0)
	assert.Nil(t, entries)
	assert.Error(t, err)
}

func TestRepositoryPg_GetFailed_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("failed query failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	entries, err := repo.GetFailed(context.Background(), 10, 0)
	assert.Nil(t, entries)
	assert.Error(t, err)
}

func TestRepositoryPg_DeleteOld_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("delete failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	count, err := repo.DeleteOld(context.Background(), time.Now())
	assert.Equal(t, int64(0), count)
	assert.Error(t, err)
}

func TestRepositoryPg_DeleteOld_Success(t *testing.T) {
	t.Parallel()
	queries := db.New(&successExecDBTX{rowsAffected: 42})
	repo := NewRepositoryPg(queries)

	count, err := repo.DeleteOld(context.Background(), time.Now())
	require.NoError(t, err)
	assert.Equal(t, int64(42), count)
}

func TestRepositoryPg_CountOld_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("count old failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	count, err := repo.CountOld(context.Background(), time.Now())
	assert.Equal(t, int64(0), count)
	assert.Error(t, err)
}

func TestRepositoryPg_GetStats_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("stats failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	stats, err := repo.GetStats(context.Background())
	assert.Nil(t, stats)
	assert.Error(t, err)
}

func TestRepositoryPg_GetRecentActions_Error(t *testing.T) {
	t.Parallel()
	dbErr := errors.New("recent actions failed")
	queries := db.New(&errDBTX{err: dbErr})
	repo := NewRepositoryPg(queries)

	actions, err := repo.GetRecentActions(context.Background(), 10)
	assert.Nil(t, actions)
	assert.Error(t, err)
}

func TestRepositoryPg_GetRecentActions_Empty(t *testing.T) {
	t.Parallel()
	queries := db.New(&noRowsDBTX{})
	repo := NewRepositoryPg(queries)

	actions, err := repo.GetRecentActions(context.Background(), 10)
	require.NoError(t, err)
	assert.Empty(t, actions)
}

func TestRepositoryPg_GetFailed_Empty(t *testing.T) {
	t.Parallel()
	queries := db.New(&noRowsDBTX{})
	repo := NewRepositoryPg(queries)

	entries, err := repo.GetFailed(context.Background(), 10, 0)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestRepositoryPg_GetByAction_Empty(t *testing.T) {
	t.Parallel()
	queries := db.New(&noRowsDBTX{})
	repo := NewRepositoryPg(queries)

	entries, err := repo.GetByAction(context.Background(), ActionUserLogin, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestRepositoryPg_GetByIP_Empty(t *testing.T) {
	t.Parallel()
	queries := db.New(&noRowsDBTX{})
	repo := NewRepositoryPg(queries)

	ip := net.ParseIP("192.168.1.1")
	entries, err := repo.GetByIP(context.Background(), ip, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

// ============================================================================
// newRepository (module.go) Tests
// ============================================================================

func TestNewRepository(t *testing.T) {
	t.Parallel()
	queries := db.New(&noRowsDBTX{})
	repo := newRepository(queries)

	require.NotNil(t, repo)
	// Verify it's a RepositoryPg
	_, ok := repo.(*RepositoryPg)
	assert.True(t, ok, "newRepository should return *RepositoryPg")
}

func TestNewRepository_NilQueries(t *testing.T) {
	t.Parallel()
	repo := newRepository(nil)
	require.NotNil(t, repo)
}
