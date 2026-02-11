package jobs

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockAuthCleanupRepo implements AuthCleanupRepository for testing.
type mockAuthCleanupRepo struct {
	mock.Mock
}

func (m *mockAuthCleanupRepo) DeleteExpiredAuthTokens(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockAuthCleanupRepo) DeleteRevokedAuthTokens(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockAuthCleanupRepo) DeleteExpiredPasswordResetTokens(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockAuthCleanupRepo) DeleteUsedPasswordResetTokens(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockAuthCleanupRepo) DeleteExpiredEmailVerificationTokens(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockAuthCleanupRepo) DeleteVerifiedEmailTokens(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockAuthCleanupRepo) DeleteOldFailedLoginAttempts(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// TestCleanupArgs_Kind tests the Kind method.
func TestCleanupArgs_Kind(t *testing.T) {
	args := CleanupArgs{}
	assert.Equal(t, CleanupJobKind, args.Kind())
	assert.Equal(t, "cleanup", args.Kind())
}

// TestCleanupArgs_InsertOpts tests the InsertOpts method.
func TestCleanupArgs_InsertOpts(t *testing.T) {
	args := CleanupArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, QueueLow, opts.Queue)
}

// TestNewCleanupWorker tests worker creation.
func TestNewCleanupWorker(t *testing.T) {
	t.Run("with logger and repo", func(t *testing.T) {
		logger := slog.Default()
		repo := &mockAuthCleanupRepo{}
		worker := NewCleanupWorker(repo, logger)

		assert.NotNil(t, worker)
		assert.NotNil(t, worker.logger)
		assert.Equal(t, logger, worker.logger)
		assert.Equal(t, repo, worker.authRepo)
	})

	t.Run("with nil logger", func(t *testing.T) {
		worker := NewCleanupWorker(nil, nil)

		assert.NotNil(t, worker)
		assert.NotNil(t, worker.logger)
	})
}

// TestCleanupWorker_Timeout tests the timeout value.
func TestCleanupWorker_Timeout(t *testing.T) {
	worker := NewCleanupWorker(nil, slog.Default())
	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 1},
	}
	assert.Equal(t, 2*time.Minute, worker.Timeout(job))
}

// TestCleanupWorker_ValidateArgs tests argument validation.
func TestCleanupWorker_ValidateArgs(t *testing.T) {
	worker := NewCleanupWorker(nil, slog.Default())

	tests := []struct {
		name    string
		args    CleanupArgs
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid args",
			args: CleanupArgs{
				TargetType: CleanupTargetExpiredTokens,
				OlderThan:  24 * time.Hour,
				BatchSize:  100,
			},
			wantErr: false,
		},
		{
			name: "empty target type",
			args: CleanupArgs{
				TargetType: "",
				OlderThan:  24 * time.Hour,
			},
			wantErr: true,
			errMsg:  "target_type is required",
		},
		{
			name: "zero older than",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  0,
			},
			wantErr: true,
			errMsg:  "older_than must be positive",
		},
		{
			name: "negative older than",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  -1 * time.Hour,
			},
			wantErr: true,
			errMsg:  "older_than must be positive",
		},
		{
			name: "negative batch size",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  24 * time.Hour,
				BatchSize:  -100,
			},
			wantErr: true,
			errMsg:  "batch_size cannot be negative",
		},
		{
			name: "zero batch size is valid",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  24 * time.Hour,
				BatchSize:  0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := worker.validateArgs(tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCleanupWorker_Work_ExpiredTokens tests cleanup of expired tokens.
func TestCleanupWorker_Work_ExpiredTokens(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredAuthTokens", mock.Anything).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 1},
		Args: CleanupArgs{
			TargetType: CleanupTargetExpiredTokens,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// TestCleanupWorker_Work_RevokedTokens tests cleanup of revoked tokens.
func TestCleanupWorker_Work_RevokedTokens(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteRevokedAuthTokens", mock.Anything).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 2},
		Args: CleanupArgs{
			TargetType: CleanupTargetRevokedTokens,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// TestCleanupWorker_Work_PasswordResets tests cleanup of password reset tokens.
func TestCleanupWorker_Work_PasswordResets(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredPasswordResetTokens", mock.Anything).Return(nil)
	repo.On("DeleteUsedPasswordResetTokens", mock.Anything).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 3},
		Args: CleanupArgs{
			TargetType: CleanupTargetPasswordResets,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// TestCleanupWorker_Work_EmailVerifications tests cleanup of email verification tokens.
func TestCleanupWorker_Work_EmailVerifications(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredEmailVerificationTokens", mock.Anything).Return(nil)
	repo.On("DeleteVerifiedEmailTokens", mock.Anything).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 4},
		Args: CleanupArgs{
			TargetType: CleanupTargetEmailVerifications,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// TestCleanupWorker_Work_FailedLogins tests cleanup of failed login attempts.
func TestCleanupWorker_Work_FailedLogins(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteOldFailedLoginAttempts", mock.Anything).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 5},
		Args: CleanupArgs{
			TargetType: CleanupTargetFailedLogins,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// TestCleanupWorker_Work_All tests cleanup of all targets.
func TestCleanupWorker_Work_All(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	// All 7 repo methods should be called
	repo.On("DeleteExpiredAuthTokens", mock.Anything).Return(nil)
	repo.On("DeleteRevokedAuthTokens", mock.Anything).Return(nil)
	repo.On("DeleteExpiredPasswordResetTokens", mock.Anything).Return(nil)
	repo.On("DeleteUsedPasswordResetTokens", mock.Anything).Return(nil)
	repo.On("DeleteExpiredEmailVerificationTokens", mock.Anything).Return(nil)
	repo.On("DeleteVerifiedEmailTokens", mock.Anything).Return(nil)
	repo.On("DeleteOldFailedLoginAttempts", mock.Anything).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 6},
		Args: CleanupArgs{
			TargetType: CleanupTargetAll,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// TestCleanupWorker_Work_RepoError tests error propagation from repository.
func TestCleanupWorker_Work_RepoError(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	dbErr := errors.New("database connection lost")
	repo.On("DeleteExpiredAuthTokens", mock.Anything).Return(dbErr)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 7},
		Args: CleanupArgs{
			TargetType: CleanupTargetExpiredTokens,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 errors")
	repo.AssertExpectations(t)
}

// TestCleanupWorker_Work_AllPartialError tests that "all" continues on partial errors.
func TestCleanupWorker_Work_AllPartialError(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	// Some succeed, some fail
	repo.On("DeleteExpiredAuthTokens", mock.Anything).Return(nil)
	repo.On("DeleteRevokedAuthTokens", mock.Anything).Return(errors.New("revoked cleanup failed"))
	repo.On("DeleteExpiredPasswordResetTokens", mock.Anything).Return(nil)
	repo.On("DeleteUsedPasswordResetTokens", mock.Anything).Return(nil)
	repo.On("DeleteExpiredEmailVerificationTokens", mock.Anything).Return(nil)
	repo.On("DeleteVerifiedEmailTokens", mock.Anything).Return(nil)
	repo.On("DeleteOldFailedLoginAttempts", mock.Anything).Return(errors.New("failed logins cleanup failed"))

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 8},
		Args: CleanupArgs{
			TargetType: CleanupTargetAll,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "2 errors")
	// All 7 methods should still be called despite errors
	repo.AssertExpectations(t)
}

// TestCleanupWorker_Work_DryRun tests dry run mode skips actual cleanup.
func TestCleanupWorker_Work_DryRun(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())
	// No repo methods should be called in dry run

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 9},
		Args: CleanupArgs{
			TargetType: CleanupTargetAll,
			OlderThan:  24 * time.Hour,
			DryRun:     true,
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
	repo.AssertNotCalled(t, "DeleteExpiredAuthTokens", mock.Anything)
}

// TestCleanupWorker_Work_InvalidArgs tests that invalid args return an error.
func TestCleanupWorker_Work_InvalidArgs(t *testing.T) {
	worker := NewCleanupWorker(nil, slog.Default())

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 10},
		Args: CleanupArgs{
			TargetType: "",
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid arguments")
	assert.Contains(t, err.Error(), "target_type is required")
}

// TestCleanupWorker_Work_UnknownTarget tests that unknown target types return an error.
func TestCleanupWorker_Work_UnknownTarget(t *testing.T) {
	worker := NewCleanupWorker(nil, slog.Default())

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 11},
		Args: CleanupArgs{
			TargetType: "nonexistent",
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown target_type")
}

// TestCleanupTargetConstants tests that target type constants have expected values.
func TestCleanupTargetConstants(t *testing.T) {
	assert.Equal(t, "expired_tokens", CleanupTargetExpiredTokens)
	assert.Equal(t, "revoked_tokens", CleanupTargetRevokedTokens)
	assert.Equal(t, "password_resets", CleanupTargetPasswordResets)
	assert.Equal(t, "email_verifications", CleanupTargetEmailVerifications)
	assert.Equal(t, "failed_logins", CleanupTargetFailedLogins)
	assert.Equal(t, "all", CleanupTargetAll)
}
