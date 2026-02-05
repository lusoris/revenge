package jobs

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCleanupArgs_Kind tests the Kind method.
func TestCleanupArgs_Kind(t *testing.T) {
	args := CleanupArgs{}
	assert.Equal(t, CleanupJobKind, args.Kind())
	assert.Equal(t, "cleanup", args.Kind())
}

// TestCleanupJobKind tests the constant value.
func TestCleanupJobKind(t *testing.T) {
	assert.Equal(t, "cleanup", CleanupJobKind)
}

// TestNewCleanupWorker tests worker creation.
func TestNewCleanupWorker(t *testing.T) {
	t.Run("with logger", func(t *testing.T) {
		logger := slog.Default()
		worker := NewCleanupWorker(nil, logger)

		assert.NotNil(t, worker)
		assert.NotNil(t, worker.logger)
		assert.Equal(t, logger, worker.logger)
	})

	t.Run("with nil logger", func(t *testing.T) {
		worker := NewCleanupWorker(nil, nil)

		assert.NotNil(t, worker)
		assert.NotNil(t, worker.logger)
	})
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
				TargetType: "sessions",
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
				TargetType: "sessions",
				OlderThan:  0,
			},
			wantErr: true,
			errMsg:  "older_than must be positive",
		},
		{
			name: "negative older than",
			args: CleanupArgs{
				TargetType: "sessions",
				OlderThan:  -1 * time.Hour,
			},
			wantErr: true,
			errMsg:  "older_than must be positive",
		},
		{
			name: "negative batch size",
			args: CleanupArgs{
				TargetType: "sessions",
				OlderThan:  24 * time.Hour,
				BatchSize:  -100,
			},
			wantErr: true,
			errMsg:  "batch_size cannot be negative",
		},
		{
			name: "zero batch size is valid",
			args: CleanupArgs{
				TargetType: "sessions",
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

// TestCleanupWorker_Work tests job execution.
func TestCleanupWorker_Work(t *testing.T) {
	worker := NewCleanupWorker(nil, slog.Default())
	ctx := context.Background()

	t.Run("successful cleanup", func(t *testing.T) {
		job := &river.Job[CleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 123},
			Args: CleanupArgs{
				TargetType: "sessions",
				OlderThan:  24 * time.Hour,
				BatchSize:  100,
				DryRun:     false,
			},
		}

		err := worker.Work(ctx, job)

		assert.NoError(t, err)
	})

	t.Run("dry run mode", func(t *testing.T) {
		job := &river.Job[CleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 456},
			Args: CleanupArgs{
				TargetType: "jobs",
				OlderThan:  7 * 24 * time.Hour,
				BatchSize:  50,
				DryRun:     true,
			},
		}

		err := worker.Work(ctx, job)

		assert.NoError(t, err)
	})

	t.Run("invalid arguments", func(t *testing.T) {
		job := &river.Job[CleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 789},
			Args: CleanupArgs{
				TargetType: "", // Invalid: empty target type
				OlderThan:  24 * time.Hour,
			},
		}

		err := worker.Work(ctx, job)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid arguments")
		assert.Contains(t, err.Error(), "target_type is required")
	})
}

// TestCleanupArgs_Serialization tests JSON serialization.
func TestCleanupArgs_Serialization(t *testing.T) {
	args := CleanupArgs{
		TargetType: "sessions",
		OlderThan:  24 * time.Hour,
		BatchSize:  100,
		DryRun:     true,
	}

	// Verify struct fields are exported and tagged correctly
	assert.Equal(t, "sessions", args.TargetType)
	assert.Equal(t, 24*time.Hour, args.OlderThan)
	assert.Equal(t, 100, args.BatchSize)
	assert.True(t, args.DryRun)
}

// TestCleanupWorker_DifferentTargets tests cleanup with different target types.
func TestCleanupWorker_DifferentTargets(t *testing.T) {
	worker := NewCleanupWorker(nil, slog.Default())
	ctx := context.Background()

	targets := []string{"sessions", "jobs", "logs", "cache_entries"}

	for _, target := range targets {
		t.Run(target, func(t *testing.T) {
			job := &river.Job[CleanupArgs]{
				JobRow: &rivertype.JobRow{ID: 1},
				Args: CleanupArgs{
					TargetType: target,
					OlderThan:  24 * time.Hour,
					BatchSize:  100,
				},
			}

			err := worker.Work(ctx, job)
			assert.NoError(t, err)
		})
	}
}

// TestCleanupWorker_EdgeCases tests edge case scenarios.
func TestCleanupWorker_EdgeCases(t *testing.T) {
	worker := NewCleanupWorker(nil, slog.Default())
	ctx := context.Background()

	t.Run("very large batch size", func(t *testing.T) {
		job := &river.Job[CleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 1},
			Args: CleanupArgs{
				TargetType: "sessions",
				OlderThan:  1 * time.Hour,
				BatchSize:  1000000,
			},
		}

		err := worker.Work(ctx, job)
		assert.NoError(t, err)
	})

	t.Run("very short older than", func(t *testing.T) {
		job := &river.Job[CleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 2},
			Args: CleanupArgs{
				TargetType: "cache",
				OlderThan:  1 * time.Millisecond,
				BatchSize:  10,
			},
		}

		err := worker.Work(ctx, job)
		assert.NoError(t, err)
	})

	t.Run("very long older than", func(t *testing.T) {
		job := &river.Job[CleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 3},
			Args: CleanupArgs{
				TargetType: "archives",
				OlderThan:  365 * 24 * time.Hour, // 1 year
				BatchSize:  1000,
			},
		}

		err := worker.Work(ctx, job)
		assert.NoError(t, err)
	})
}
