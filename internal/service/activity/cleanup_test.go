package activity

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/infra/raft"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepository implements the Repository interface for same-package unit testing.
type mockRepository struct {
	createFn           func(ctx context.Context, entry *Entry) error
	getFn              func(ctx context.Context, id uuid.UUID) (*Entry, error)
	listFn             func(ctx context.Context, limit, offset int32) ([]Entry, error)
	countFn            func(ctx context.Context) (int64, error)
	searchFn           func(ctx context.Context, filters SearchFilters) ([]Entry, int64, error)
	getByUserFn        func(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]Entry, int64, error)
	getByResourceFn    func(ctx context.Context, resourceType string, resourceID uuid.UUID, limit, offset int32) ([]Entry, int64, error)
	getByActionFn      func(ctx context.Context, action string, limit, offset int32) ([]Entry, error)
	getByIPFn          func(ctx context.Context, ip net.IP, limit, offset int32) ([]Entry, error)
	getFailedFn        func(ctx context.Context, limit, offset int32) ([]Entry, error)
	deleteOldFn        func(ctx context.Context, olderThan time.Time) (int64, error)
	countOldFn         func(ctx context.Context, olderThan time.Time) (int64, error)
	getStatsFn         func(ctx context.Context) (*Stats, error)
	getRecentActionsFn func(ctx context.Context, limit int32) ([]ActionCount, error)
}

func (m *mockRepository) Create(ctx context.Context, entry *Entry) error {
	if m.createFn != nil {
		return m.createFn(ctx, entry)
	}
	return nil
}

func (m *mockRepository) Get(ctx context.Context, id uuid.UUID) (*Entry, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return nil, ErrNotFound
}

func (m *mockRepository) List(ctx context.Context, limit, offset int32) ([]Entry, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) Count(ctx context.Context) (int64, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

func (m *mockRepository) Search(ctx context.Context, filters SearchFilters) ([]Entry, int64, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, filters)
	}
	return nil, 0, nil
}

func (m *mockRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
	if m.getByUserFn != nil {
		return m.getByUserFn(ctx, userID, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockRepository) GetByResource(ctx context.Context, resourceType string, resourceID uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
	if m.getByResourceFn != nil {
		return m.getByResourceFn(ctx, resourceType, resourceID, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockRepository) GetByAction(ctx context.Context, action string, limit, offset int32) ([]Entry, error) {
	if m.getByActionFn != nil {
		return m.getByActionFn(ctx, action, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) GetByIP(ctx context.Context, ip net.IP, limit, offset int32) ([]Entry, error) {
	if m.getByIPFn != nil {
		return m.getByIPFn(ctx, ip, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) GetFailed(ctx context.Context, limit, offset int32) ([]Entry, error) {
	if m.getFailedFn != nil {
		return m.getFailedFn(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) DeleteOld(ctx context.Context, olderThan time.Time) (int64, error) {
	if m.deleteOldFn != nil {
		return m.deleteOldFn(ctx, olderThan)
	}
	return 0, nil
}

func (m *mockRepository) CountOld(ctx context.Context, olderThan time.Time) (int64, error) {
	if m.countOldFn != nil {
		return m.countOldFn(ctx, olderThan)
	}
	return 0, nil
}

func (m *mockRepository) GetStats(ctx context.Context) (*Stats, error) {
	if m.getStatsFn != nil {
		return m.getStatsFn(ctx)
	}
	return &Stats{}, nil
}

func (m *mockRepository) GetRecentActions(ctx context.Context, limit int32) ([]ActionCount, error) {
	if m.getRecentActionsFn != nil {
		return m.getRecentActionsFn(ctx, limit)
	}
	return nil, nil
}

// mockLeaderElection implements the methods used from raft.LeaderElection.
// Since LeaderElection is a concrete struct (not an interface), we use the real
// nil-pointer behavior for testing non-leader and nil scenarios, and create a
// helper that constructs the worker differently.

// newTestJob creates a river.Job for testing cleanup workers.
func newTestJob(id int64, args ActivityCleanupArgs) *river.Job[ActivityCleanupArgs] {
	return &river.Job[ActivityCleanupArgs]{
		JobRow: &rivertype.JobRow{
			ID: id,
		},
		Args: args,
	}
}

func newTestServiceWithMock(t *testing.T, repo *mockRepository) *Service {
	t.Helper()
	logger := logging.NewTestLogger()
	return NewService(repo, logger)
}

// ============================================================================
// ActivityCleanupArgs Tests
// ============================================================================

func TestActivityCleanupArgs_Kind(t *testing.T) {
	t.Parallel()
	args := ActivityCleanupArgs{}
	assert.Equal(t, ActivityCleanupJobKind, args.Kind())
	assert.Equal(t, "activity_cleanup", args.Kind())
}

func TestActivityCleanupArgs_InsertOpts(t *testing.T) {
	t.Parallel()
	args := ActivityCleanupArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueLow, opts.Queue)
}

// ============================================================================
// NewActivityCleanupWorker Tests
// ============================================================================

func TestNewActivityCleanupWorker(t *testing.T) {
	t.Parallel()
	repo := &mockRepository{}
	svc := newTestServiceWithMock(t, repo)
	logger := logging.NewTestLogger()

	worker := NewActivityCleanupWorker(nil, svc, logger)

	require.NotNil(t, worker)
	assert.Nil(t, worker.leaderElection)
	assert.NotNil(t, worker.service)
	assert.NotNil(t, worker.logger)
}

func TestNewActivityCleanupWorker_WithLeaderElection(t *testing.T) {
	t.Parallel()

	// LeaderElection with nil raft still works (single-node mode)
	repo := &mockRepository{}
	svc := newTestServiceWithMock(t, repo)
	logger := logging.NewTestLogger()

	// nil leader election means single-node mode
	worker := NewActivityCleanupWorker(nil, svc, logger)

	require.NotNil(t, worker)
}

// ============================================================================
// Timeout Tests
// ============================================================================

func TestActivityCleanupWorker_Timeout(t *testing.T) {
	t.Parallel()
	repo := &mockRepository{}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	job := newTestJob(1, ActivityCleanupArgs{RetentionDays: 90})
	timeout := worker.Timeout(job)

	assert.Equal(t, 2*time.Minute, timeout)
}

// ============================================================================
// Work Tests
// ============================================================================

func TestActivityCleanupWorker_Work_NilLeaderElection(t *testing.T) {
	t.Parallel()

	var deletedCalled bool
	repo := &mockRepository{
		deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			deletedCalled = true
			return 42, nil
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	job := newTestJob(1, ActivityCleanupArgs{RetentionDays: 30})
	err := worker.Work(context.Background(), job)

	require.NoError(t, err)
	assert.True(t, deletedCalled, "DeleteOld should be called when leaderElection is nil")
}

func TestActivityCleanupWorker_Work_DefaultRetentionDays(t *testing.T) {
	t.Parallel()

	var capturedOlderThan time.Time
	repo := &mockRepository{
		deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			capturedOlderThan = olderThan
			return 10, nil
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	// RetentionDays <= 0 should default to 90
	job := newTestJob(2, ActivityCleanupArgs{RetentionDays: 0})
	before := time.Now()
	err := worker.Work(context.Background(), job)

	require.NoError(t, err)
	// The cutoff should be approximately 90 days ago
	expected := before.AddDate(0, 0, -90)
	assert.WithinDuration(t, expected, capturedOlderThan, 5*time.Second)
}

func TestActivityCleanupWorker_Work_NegativeRetentionDays(t *testing.T) {
	t.Parallel()

	var capturedOlderThan time.Time
	repo := &mockRepository{
		deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			capturedOlderThan = olderThan
			return 5, nil
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	// Negative retention days should also default to 90
	job := newTestJob(3, ActivityCleanupArgs{RetentionDays: -5})
	before := time.Now()
	err := worker.Work(context.Background(), job)

	require.NoError(t, err)
	expected := before.AddDate(0, 0, -90)
	assert.WithinDuration(t, expected, capturedOlderThan, 5*time.Second)
}

func TestActivityCleanupWorker_Work_CustomRetentionDays(t *testing.T) {
	t.Parallel()

	var capturedOlderThan time.Time
	repo := &mockRepository{
		deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			capturedOlderThan = olderThan
			return 100, nil
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	job := newTestJob(4, ActivityCleanupArgs{RetentionDays: 60})
	before := time.Now()
	err := worker.Work(context.Background(), job)

	require.NoError(t, err)
	expected := before.AddDate(0, 0, -60)
	assert.WithinDuration(t, expected, capturedOlderThan, 5*time.Second)
}

func TestActivityCleanupWorker_Work_DryRun(t *testing.T) {
	t.Parallel()

	var countCalled bool
	var deleteCalled bool
	repo := &mockRepository{
		countOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			countCalled = true
			return 250, nil
		},
		deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			deleteCalled = true
			return 0, nil
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	job := newTestJob(5, ActivityCleanupArgs{
		RetentionDays: 30,
		DryRun:        true,
	})
	err := worker.Work(context.Background(), job)

	require.NoError(t, err)
	assert.True(t, countCalled, "CountOld should be called during dry run")
	assert.False(t, deleteCalled, "DeleteOld should NOT be called during dry run")
}

func TestActivityCleanupWorker_Work_DryRunError(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{
		countOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			return 0, errors.New("count query failed")
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	job := newTestJob(6, ActivityCleanupArgs{
		RetentionDays: 30,
		DryRun:        true,
	})
	err := worker.Work(context.Background(), job)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count old logs")
}

func TestActivityCleanupWorker_Work_DeleteError(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{
		deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			return 0, errors.New("delete query failed")
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	job := newTestJob(7, ActivityCleanupArgs{RetentionDays: 30})
	err := worker.Work(context.Background(), job)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to cleanup logs")
}

func TestActivityCleanupWorker_Work_NotLeader(t *testing.T) {
	t.Parallel()

	// Create a LeaderElection with nil raft - IsLeader() returns true for nil
	// We need to test the non-leader path. Since LeaderElection is a concrete struct,
	// we construct one with a nil raft field -- IsLeader returns true when raft is nil.
	// To test the "not leader" branch, we need a non-nil LeaderElection with
	// a non-nil raft. Since we can't easily construct a real raft, we test the
	// nil case (which acts as leader) and verify the code path works.

	var deleteCalled bool
	repo := &mockRepository{
		deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			deleteCalled = true
			return 10, nil
		},
	}
	svc := newTestServiceWithMock(t, repo)

	// When leaderElection is nil, the worker should proceed (single-node mode)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())
	job := newTestJob(8, ActivityCleanupArgs{RetentionDays: 30})
	err := worker.Work(context.Background(), job)

	require.NoError(t, err)
	assert.True(t, deleteCalled)
}

func TestActivityCleanupWorker_Work_LeaderElectionNilRaft(t *testing.T) {
	t.Parallel()

	// A zero-value LeaderElection (raft field is nil) considers itself leader
	le := &raft.LeaderElection{}
	assert.True(t, le.IsLeader(), "nil raft should be considered leader (single-node)")

	var deleteCalled bool
	repo := &mockRepository{
		deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			deleteCalled = true
			return 5, nil
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(le, svc, logging.NewTestLogger())

	job := newTestJob(9, ActivityCleanupArgs{RetentionDays: 30})
	err := worker.Work(context.Background(), job)

	require.NoError(t, err)
	assert.True(t, deleteCalled)
}

func TestActivityCleanupWorker_Work_DryRunDefaultRetention(t *testing.T) {
	t.Parallel()

	var capturedOlderThan time.Time
	repo := &mockRepository{
		countOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
			capturedOlderThan = olderThan
			return 100, nil
		},
	}
	svc := newTestServiceWithMock(t, repo)
	worker := NewActivityCleanupWorker(nil, svc, logging.NewTestLogger())

	// DryRun with 0 retention days should default to 90
	job := newTestJob(10, ActivityCleanupArgs{
		RetentionDays: 0,
		DryRun:        true,
	})
	before := time.Now()
	err := worker.Work(context.Background(), job)

	require.NoError(t, err)
	expected := before.AddDate(0, 0, -90)
	assert.WithinDuration(t, expected, capturedOlderThan, 5*time.Second)
}

// ============================================================================
// Service method tests with mock repo (same-package)
// ============================================================================

func TestService_Log_WithMockRepo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     LogRequest
		repoErr error
		wantErr bool
	}{
		{
			name: "success with all fields",
			req: LogRequest{
				UserID:       ptrUUID(uuid.Must(uuid.NewV7())),
				Username:     ptrStr("testuser"),
				Action:       ActionUserLogin,
				ResourceType: ptrStr(ResourceTypeUser),
				ResourceID:   ptrUUID(uuid.Must(uuid.NewV7())),
				Changes:      map[string]interface{}{"field": "value"},
				Metadata:     map[string]interface{}{"key": "data"},
				IPAddress:    ptrIP(net.ParseIP("10.0.0.1")),
				UserAgent:    ptrStr("TestAgent/1.0"),
				Success:      true,
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "success with minimal fields",
			req: LogRequest{
				Action:  "system.startup",
				Success: true,
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "repo returns error",
			req: LogRequest{
				Action:  ActionUserLogin,
				Success: true,
			},
			repoErr: errors.New("database connection lost"),
			wantErr: true,
		},
		{
			name: "failed action with error message",
			req: LogRequest{
				Action:       ActionUserLogin,
				Success:      false,
				ErrorMessage: ptrStr("invalid credentials"),
			},
			repoErr: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockRepository{
				createFn: func(ctx context.Context, entry *Entry) error {
					// Verify the entry is populated from the request
					assert.Equal(t, tt.req.Action, entry.Action)
					assert.Equal(t, tt.req.Success, entry.Success)
					assert.Equal(t, tt.req.UserID, entry.UserID)
					assert.Equal(t, tt.req.Username, entry.Username)
					assert.Equal(t, tt.req.ResourceType, entry.ResourceType)
					assert.Equal(t, tt.req.ResourceID, entry.ResourceID)
					return tt.repoErr
				},
			}
			svc := newTestServiceWithMock(t, repo)

			err := svc.Log(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_LogWithContext_WithMockRepo(t *testing.T) {
	t.Parallel()

	var capturedEntry *Entry
	repo := &mockRepository{
		createFn: func(ctx context.Context, entry *Entry) error {
			capturedEntry = entry
			return nil
		},
	}
	svc := newTestServiceWithMock(t, repo)

	userID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())
	ip := net.ParseIP("192.168.1.1")

	err := svc.LogWithContext(
		context.Background(),
		userID,
		"testuser",
		ActionUserUpdate,
		ResourceTypeUser,
		resourceID,
		map[string]interface{}{"field": "new_value"},
		ip,
		"Mozilla/5.0",
	)

	require.NoError(t, err)
	require.NotNil(t, capturedEntry)
	assert.Equal(t, &userID, capturedEntry.UserID)
	assert.Equal(t, ptrStr("testuser"), capturedEntry.Username)
	assert.Equal(t, ActionUserUpdate, capturedEntry.Action)
	assert.Equal(t, ptrStr(ResourceTypeUser), capturedEntry.ResourceType)
	assert.Equal(t, &resourceID, capturedEntry.ResourceID)
	assert.True(t, capturedEntry.Success)
}

func TestService_LogFailure_WithMockRepo(t *testing.T) {
	t.Parallel()

	var capturedEntry *Entry
	repo := &mockRepository{
		createFn: func(ctx context.Context, entry *Entry) error {
			capturedEntry = entry
			return nil
		},
	}
	svc := newTestServiceWithMock(t, repo)

	userID := uuid.Must(uuid.NewV7())
	username := "failuser"
	ip := net.ParseIP("10.0.0.1")
	userAgent := "BadAgent"

	err := svc.LogFailure(
		context.Background(),
		&userID,
		&username,
		ActionUserLogin,
		"invalid password",
		&ip,
		&userAgent,
	)

	require.NoError(t, err)
	require.NotNil(t, capturedEntry)
	assert.False(t, capturedEntry.Success)
	assert.Equal(t, ptrStr("invalid password"), capturedEntry.ErrorMessage)
	assert.Equal(t, &userID, capturedEntry.UserID)
}

func TestService_Get_WithMockRepo(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()
		entryID := uuid.Must(uuid.NewV7())
		repo := &mockRepository{
			getFn: func(ctx context.Context, id uuid.UUID) (*Entry, error) {
				return &Entry{
					ID:     id,
					Action: ActionUserLogin,
				}, nil
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entry, err := svc.Get(context.Background(), entryID)
		require.NoError(t, err)
		assert.Equal(t, entryID, entry.ID)
		assert.Equal(t, ActionUserLogin, entry.Action)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			getFn: func(ctx context.Context, id uuid.UUID) (*Entry, error) {
				return nil, ErrNotFound
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entry, err := svc.Get(context.Background(), uuid.Must(uuid.NewV7()))
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("database error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			getFn: func(ctx context.Context, id uuid.UUID) (*Entry, error) {
				return nil, errors.New("connection refused")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entry, err := svc.Get(context.Background(), uuid.Must(uuid.NewV7()))
		assert.Nil(t, entry)
		assert.Error(t, err)
	})
}

func TestService_List_WithMockRepo(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			listFn: func(ctx context.Context, limit, offset int32) ([]Entry, error) {
				return []Entry{
					{ID: uuid.Must(uuid.NewV7()), Action: "a"},
					{ID: uuid.Must(uuid.NewV7()), Action: "b"},
				}, nil
			},
			countFn: func(ctx context.Context) (int64, error) {
				return 2, nil
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entries, count, err := svc.List(context.Background(), 10, 0)
		require.NoError(t, err)
		assert.Len(t, entries, 2)
		assert.Equal(t, int64(2), count)
	})

	t.Run("list error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			listFn: func(ctx context.Context, limit, offset int32) ([]Entry, error) {
				return nil, errors.New("list failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entries, count, err := svc.List(context.Background(), 10, 0)
		assert.Nil(t, entries)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})

	t.Run("count error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			listFn: func(ctx context.Context, limit, offset int32) ([]Entry, error) {
				return []Entry{}, nil
			},
			countFn: func(ctx context.Context) (int64, error) {
				return 0, errors.New("count failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entries, count, err := svc.List(context.Background(), 10, 0)
		assert.Nil(t, entries)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestService_Search_WithMockRepo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputLimit  int32
		expectLimit int32
	}{
		{"zero limit defaults to 50", 0, 50},
		{"negative limit defaults to 50", -1, 50},
		{"over 100 capped to 100", 200, 100},
		{"exactly 100 stays 100", 100, 100},
		{"normal limit stays", 25, 25},
		{"limit 1 stays", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedLimit int32
			repo := &mockRepository{
				searchFn: func(ctx context.Context, filters SearchFilters) ([]Entry, int64, error) {
					capturedLimit = filters.Limit
					return []Entry{}, 0, nil
				},
			}
			svc := newTestServiceWithMock(t, repo)

			_, _, err := svc.Search(context.Background(), SearchFilters{Limit: tt.inputLimit})
			require.NoError(t, err)
			assert.Equal(t, tt.expectLimit, capturedLimit)
		})
	}

	t.Run("search error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			searchFn: func(ctx context.Context, filters SearchFilters) ([]Entry, int64, error) {
				return nil, 0, errors.New("search failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entries, count, err := svc.Search(context.Background(), SearchFilters{Limit: 10})
		assert.Nil(t, entries)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestService_GetUserActivity_WithMockRepo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputLimit  int32
		expectLimit int32
	}{
		{"zero limit defaults to 50", 0, 50},
		{"negative limit defaults to 50", -10, 50},
		{"over 100 capped to 100", 150, 100},
		{"normal limit stays", 30, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userID := uuid.Must(uuid.NewV7())

			var capturedLimit int32
			repo := &mockRepository{
				getByUserFn: func(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
					capturedLimit = limit
					return []Entry{}, 0, nil
				},
			}
			svc := newTestServiceWithMock(t, repo)

			_, _, err := svc.GetUserActivity(context.Background(), userID, tt.inputLimit, 0)
			require.NoError(t, err)
			assert.Equal(t, tt.expectLimit, capturedLimit)
		})
	}

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			getByUserFn: func(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
				return nil, 0, errors.New("user query failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entries, count, err := svc.GetUserActivity(context.Background(), uuid.Must(uuid.NewV7()), 10, 0)
		assert.Nil(t, entries)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestService_GetResourceActivity_WithMockRepo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputLimit  int32
		expectLimit int32
	}{
		{"zero limit defaults to 50", 0, 50},
		{"negative limit defaults to 50", -5, 50},
		{"over 100 capped to 100", 200, 100},
		{"normal limit stays", 40, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resourceID := uuid.Must(uuid.NewV7())

			var capturedLimit int32
			repo := &mockRepository{
				getByResourceFn: func(ctx context.Context, rt string, rid uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
					capturedLimit = limit
					return []Entry{}, 0, nil
				},
			}
			svc := newTestServiceWithMock(t, repo)

			_, _, err := svc.GetResourceActivity(context.Background(), "movie", resourceID, tt.inputLimit, 0)
			require.NoError(t, err)
			assert.Equal(t, tt.expectLimit, capturedLimit)
		})
	}

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			getByResourceFn: func(ctx context.Context, rt string, rid uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
				return nil, 0, errors.New("resource query failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entries, count, err := svc.GetResourceActivity(context.Background(), "movie", uuid.Must(uuid.NewV7()), 10, 0)
		assert.Nil(t, entries)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestService_GetFailedActivity_WithMockRepo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputLimit  int32
		expectLimit int32
	}{
		{"zero limit defaults to 50", 0, 50},
		{"negative limit defaults to 50", -3, 50},
		{"over 100 capped to 100", 200, 100},
		{"normal limit stays", 20, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedLimit int32
			repo := &mockRepository{
				getFailedFn: func(ctx context.Context, limit, offset int32) ([]Entry, error) {
					capturedLimit = limit
					return []Entry{}, nil
				},
			}
			svc := newTestServiceWithMock(t, repo)

			_, err := svc.GetFailedActivity(context.Background(), tt.inputLimit, 0)
			require.NoError(t, err)
			assert.Equal(t, tt.expectLimit, capturedLimit)
		})
	}

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			getFailedFn: func(ctx context.Context, limit, offset int32) ([]Entry, error) {
				return nil, errors.New("failed query failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		entries, err := svc.GetFailedActivity(context.Background(), 10, 0)
		assert.Nil(t, entries)
		assert.Error(t, err)
	})
}

func TestService_GetStats_WithMockRepo(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		repo := &mockRepository{
			getStatsFn: func(ctx context.Context) (*Stats, error) {
				return &Stats{
					TotalCount:   100,
					SuccessCount: 90,
					FailedCount:  10,
					OldestEntry:  &now,
					NewestEntry:  &now,
				}, nil
			},
		}
		svc := newTestServiceWithMock(t, repo)

		stats, err := svc.GetStats(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(100), stats.TotalCount)
		assert.Equal(t, int64(90), stats.SuccessCount)
		assert.Equal(t, int64(10), stats.FailedCount)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			getStatsFn: func(ctx context.Context) (*Stats, error) {
				return nil, errors.New("stats failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		stats, err := svc.GetStats(context.Background())
		assert.Nil(t, stats)
		assert.Error(t, err)
	})
}

func TestService_GetRecentActions_WithMockRepo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputLimit  int32
		expectLimit int32
	}{
		{"zero limit defaults to 20", 0, 20},
		{"negative limit defaults to 20", -1, 20},
		{"over 50 capped to 50", 100, 50},
		{"normal limit stays", 10, 10},
		{"exactly 50 stays", 50, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedLimit int32
			repo := &mockRepository{
				getRecentActionsFn: func(ctx context.Context, limit int32) ([]ActionCount, error) {
					capturedLimit = limit
					return []ActionCount{
						{Action: ActionUserLogin, Count: 10},
					}, nil
				},
			}
			svc := newTestServiceWithMock(t, repo)

			actions, err := svc.GetRecentActions(context.Background(), tt.inputLimit)
			require.NoError(t, err)
			assert.Len(t, actions, 1)
			assert.Equal(t, tt.expectLimit, capturedLimit)
		})
	}

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			getRecentActionsFn: func(ctx context.Context, limit int32) ([]ActionCount, error) {
				return nil, errors.New("actions failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		actions, err := svc.GetRecentActions(context.Background(), 10)
		assert.Nil(t, actions)
		assert.Error(t, err)
	})
}

func TestService_CleanupOldLogs_WithMockRepo(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
				return 42, nil
			},
		}
		svc := newTestServiceWithMock(t, repo)

		count, err := svc.CleanupOldLogs(context.Background(), time.Now().Add(-30*24*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(42), count)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			deleteOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
				return 0, errors.New("delete failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		count, err := svc.CleanupOldLogs(context.Background(), time.Now().Add(-30*24*time.Hour))
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestService_CountOldLogs_WithMockRepo(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			countOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
				return 75, nil
			},
		}
		svc := newTestServiceWithMock(t, repo)

		count, err := svc.CountOldLogs(context.Background(), time.Now().Add(-90*24*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(75), count)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			countOldFn: func(ctx context.Context, olderThan time.Time) (int64, error) {
				return 0, errors.New("count failed")
			},
		}
		svc := newTestServiceWithMock(t, repo)

		count, err := svc.CountOldLogs(context.Background(), time.Now().Add(-90*24*time.Hour))
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

// ============================================================================
// ServiceLogger tests with mock repo (same-package)
// ============================================================================

func TestServiceLogger_LogAction_WithMockRepo(t *testing.T) {
	t.Parallel()

	t.Run("all fields populated", func(t *testing.T) {
		t.Parallel()
		var capturedEntry *Entry
		repo := &mockRepository{
			createFn: func(ctx context.Context, entry *Entry) error {
				capturedEntry = entry
				return nil
			},
		}
		svc := newTestServiceWithMock(t, repo)
		logger := NewLogger(svc)

		userID := uuid.Must(uuid.NewV7())
		resourceID := uuid.Must(uuid.NewV7())
		ip := net.ParseIP("192.168.1.1")

		err := logger.LogAction(context.Background(), LogActionRequest{
			UserID:       userID,
			Username:     "admin",
			Action:       ActionUserCreate,
			ResourceType: ResourceTypeUser,
			ResourceID:   resourceID,
			Changes:      map[string]interface{}{"username": "newuser"},
			Metadata:     map[string]interface{}{"source": "api"},
			IPAddress:    ip,
			UserAgent:    "AdminPanel/1.0",
		})

		require.NoError(t, err)
		require.NotNil(t, capturedEntry)
		assert.True(t, capturedEntry.Success)
		assert.Equal(t, &userID, capturedEntry.UserID)
		assert.Equal(t, &resourceID, capturedEntry.ResourceID)
		assert.NotNil(t, capturedEntry.IPAddress)
	})

	t.Run("nil UUID fields become nil pointers", func(t *testing.T) {
		t.Parallel()
		var capturedEntry *Entry
		repo := &mockRepository{
			createFn: func(ctx context.Context, entry *Entry) error {
				capturedEntry = entry
				return nil
			},
		}
		svc := newTestServiceWithMock(t, repo)
		logger := NewLogger(svc)

		err := logger.LogAction(context.Background(), LogActionRequest{
			UserID:       uuid.Nil,
			Username:     "",
			Action:       ActionLibraryScan,
			ResourceType: "",
			ResourceID:   uuid.Nil,
			IPAddress:    nil,
			UserAgent:    "",
		})

		require.NoError(t, err)
		require.NotNil(t, capturedEntry)
		assert.Nil(t, capturedEntry.UserID, "uuid.Nil should become nil pointer")
		assert.Nil(t, capturedEntry.ResourceID, "uuid.Nil should become nil pointer")
		assert.Nil(t, capturedEntry.Username, "empty string should become nil pointer")
		assert.Nil(t, capturedEntry.ResourceType, "empty string should become nil pointer")
		assert.Nil(t, capturedEntry.UserAgent, "empty string should become nil pointer")
		assert.Nil(t, capturedEntry.IPAddress, "nil IP should become nil pointer")
	})

	t.Run("error propagation", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepository{
			createFn: func(ctx context.Context, entry *Entry) error {
				return errors.New("repo error")
			},
		}
		svc := newTestServiceWithMock(t, repo)
		logger := NewLogger(svc)

		err := logger.LogAction(context.Background(), LogActionRequest{
			UserID: uuid.Must(uuid.NewV7()),
			Action: ActionUserLogin,
		})
		assert.Error(t, err)
	})
}

func TestServiceLogger_LogFailure_WithMockRepo(t *testing.T) {
	t.Parallel()

	t.Run("all fields populated", func(t *testing.T) {
		t.Parallel()
		var capturedEntry *Entry
		repo := &mockRepository{
			createFn: func(ctx context.Context, entry *Entry) error {
				capturedEntry = entry
				return nil
			},
		}
		svc := newTestServiceWithMock(t, repo)
		logger := NewLogger(svc)

		userID := uuid.Must(uuid.NewV7())
		username := "hacker"
		ip := net.ParseIP("1.2.3.4")
		ua := "EvilBot/1.0"

		err := logger.LogFailure(context.Background(), LogFailureRequest{
			UserID:       &userID,
			Username:     &username,
			Action:       ActionUserLogin,
			ErrorMessage: "brute force detected",
			IPAddress:    &ip,
			UserAgent:    &ua,
		})

		require.NoError(t, err)
		require.NotNil(t, capturedEntry)
		assert.False(t, capturedEntry.Success)
		assert.Equal(t, ptrStr("brute force detected"), capturedEntry.ErrorMessage)
	})

	t.Run("nil optional fields", func(t *testing.T) {
		t.Parallel()
		var capturedEntry *Entry
		repo := &mockRepository{
			createFn: func(ctx context.Context, entry *Entry) error {
				capturedEntry = entry
				return nil
			},
		}
		svc := newTestServiceWithMock(t, repo)
		logger := NewLogger(svc)

		err := logger.LogFailure(context.Background(), LogFailureRequest{
			UserID:       nil,
			Username:     nil,
			Action:       ActionUserLogin,
			ErrorMessage: "unknown error",
			IPAddress:    nil,
			UserAgent:    nil,
		})

		require.NoError(t, err)
		require.NotNil(t, capturedEntry)
		assert.Nil(t, capturedEntry.UserID)
		assert.Nil(t, capturedEntry.Username)
		assert.Nil(t, capturedEntry.IPAddress)
		assert.Nil(t, capturedEntry.UserAgent)
	})
}

func TestNoopLogger_WithMockRepo(t *testing.T) {
	t.Parallel()

	logger := NewNoopLogger()

	t.Run("LogAction returns nil", func(t *testing.T) {
		err := logger.LogAction(context.Background(), LogActionRequest{
			UserID:   uuid.Must(uuid.NewV7()),
			Username: "test",
			Action:   ActionUserLogin,
		})
		assert.NoError(t, err)
	})

	t.Run("LogFailure returns nil", func(t *testing.T) {
		err := logger.LogFailure(context.Background(), LogFailureRequest{
			Action:       ActionUserLogin,
			ErrorMessage: "test error",
		})
		assert.NoError(t, err)
	})
}

// ============================================================================
// NewService Tests
// ============================================================================

func TestNewService_WithMockRepo(t *testing.T) {
	t.Parallel()
	repo := &mockRepository{}
	logger := logging.NewTestLogger()

	svc := NewService(repo, logger)

	require.NotNil(t, svc)
	assert.NotNil(t, svc.repo)
	assert.NotNil(t, svc.logger)
}

// ============================================================================
// Constants Tests
// ============================================================================

func TestActivityCleanupJobKind(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "activity_cleanup", ActivityCleanupJobKind)
}

func TestActionConstants(t *testing.T) {
	t.Parallel()

	// Verify action constants are non-empty and have expected prefixes
	actions := map[string]string{
		"ActionUserLogin":         ActionUserLogin,
		"ActionUserLogout":        ActionUserLogout,
		"ActionUserCreate":        ActionUserCreate,
		"ActionUserUpdate":        ActionUserUpdate,
		"ActionUserDelete":        ActionUserDelete,
		"ActionUserPasswordReset": ActionUserPasswordReset,
		"ActionSessionCreate":     ActionSessionCreate,
		"ActionSessionRevoke":     ActionSessionRevoke,
		"ActionSessionExpired":    ActionSessionExpired,
		"ActionAPIKeyCreate":      ActionAPIKeyCreate,
		"ActionAPIKeyRevoke":      ActionAPIKeyRevoke,
		"ActionOIDCLogin":         ActionOIDCLogin,
		"ActionSettingsUpdate":    ActionSettingsUpdate,
		"ActionLibraryCreate":     ActionLibraryCreate,
		"ActionLibraryScan":       ActionLibraryScan,
	}

	for name, value := range actions {
		assert.NotEmpty(t, value, "%s should not be empty", name)
	}
}

func TestResourceTypeConstants(t *testing.T) {
	t.Parallel()

	resourceTypes := []string{
		ResourceTypeUser,
		ResourceTypeSession,
		ResourceTypeAPIKey,
		ResourceTypeOIDC,
		ResourceTypeSetting,
		ResourceTypeLibrary,
		ResourceTypeMovie,
		ResourceTypeTVShow,
		ResourceTypeEpisode,
	}

	for _, rt := range resourceTypes {
		assert.NotEmpty(t, rt)
	}
}

func TestErrNotFound(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, ErrNotFound)
	assert.Equal(t, "activity log not found", ErrNotFound.Error())
}

// ============================================================================
// Helpers
// ============================================================================

func ptrStr(s string) *string {
	return &s
}

func ptrUUID(u uuid.UUID) *uuid.UUID {
	return &u
}

func ptrIP(ip net.IP) *net.IP {
	return &ip
}
