package playbackjobs

import (
	"context"
	"log/slog"
	"time"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/playback"
	"github.com/lusoris/revenge/internal/playback/transcode"
	"github.com/riverqueue/river"
)

// PlaybackCleanupJobKind is the unique identifier for playback cleanup jobs.
const PlaybackCleanupJobKind = "playback_cleanup"

// CleanupArgs defines the arguments for the playback cleanup job.
type CleanupArgs struct{}

// Kind returns the job kind identifier.
func (CleanupArgs) Kind() string {
	return PlaybackCleanupJobKind
}

// InsertOpts returns the default insert options.
func (CleanupArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueLow,
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByPeriod: 5 * time.Minute,
		},
	}
}

// CleanupWorker periodically logs playback session health.
// Actual session expiration is handled by L1Cache (otter) TTL.
// This worker serves as a monitoring heartbeat and safety net.
type CleanupWorker struct {
	river.WorkerDefaults[CleanupArgs]
	sessions *playback.SessionManager
	pipeline *transcode.PipelineManager
	logger   *slog.Logger
}

// NewCleanupWorker creates a new playback cleanup worker.
func NewCleanupWorker(sessions *playback.SessionManager, pipeline *transcode.PipelineManager, logger *slog.Logger) *CleanupWorker {
	return &CleanupWorker{
		sessions: sessions,
		pipeline: pipeline,
		logger:   logger,
	}
}

// Timeout returns the maximum execution time for cleanup jobs.
func (w *CleanupWorker) Timeout(_ *river.Job[CleanupArgs]) time.Duration {
	return 1 * time.Minute
}

// Work executes the playback cleanup job.
func (w *CleanupWorker) Work(_ context.Context, _ *river.Job[CleanupArgs]) error {
	count := w.sessions.ActiveCount()
	w.logger.Info("playback session health check",
		slog.Int("active_sessions", count),
	)
	return nil
}
