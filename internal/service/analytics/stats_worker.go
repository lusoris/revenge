// Package analytics provides periodic stats aggregation for server-wide metrics.
package analytics

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/lusoris/revenge/internal/infra/database/db"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/riverqueue/river"
	"golang.org/x/sync/errgroup"
)

// Stat key constants used as server_stats.stat_key values.
const (
	StatTotalUsers          = "total_users"
	StatActiveUsers24h      = "active_users_24h"
	StatTotalLibraries      = "total_libraries"
	StatTotalMovies         = "total_movies"
	StatTotalSeries         = "total_series"
	StatTotalEpisodes       = "total_episodes"
	StatTotalMoviePlays     = "total_movie_plays"
	StatTotalEpisodePlays   = "total_episode_plays"
	StatMovieWatchSeconds   = "movie_watch_seconds"
	StatEpisodeWatchSeconds = "episode_watch_seconds"
)

// StatsAggregationJobKind is the unique identifier for stats aggregation jobs.
const StatsAggregationJobKind = "stats_aggregation"

// StatsAggregationArgs defines the arguments for the periodic stats aggregation job.
type StatsAggregationArgs struct{}

// Kind returns the job kind identifier.
func (StatsAggregationArgs) Kind() string {
	return StatsAggregationJobKind
}

// InsertOpts returns the default insert options for stats aggregation jobs.
func (StatsAggregationArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueLow,
		MaxAttempts: 3,
	}
}

// StatsAggregationWorker computes server-wide aggregate stats and persists them.
// Leader election is handled by River's built-in leader election.
type StatsAggregationWorker struct {
	river.WorkerDefaults[StatsAggregationArgs]
	queries *db.Queries
	logger  *slog.Logger
}

// NewStatsAggregationWorker creates a new stats aggregation worker.
func NewStatsAggregationWorker(
	queries *db.Queries,
	logger *slog.Logger,
) *StatsAggregationWorker {
	return &StatsAggregationWorker{
		queries: queries,
		logger:  logger.With("component", "stats-aggregation"),
	}
}

// Timeout returns the maximum execution time for stats aggregation jobs.
func (w *StatsAggregationWorker) Timeout(_ *river.Job[StatsAggregationArgs]) time.Duration {
	return 2 * time.Minute
}

// Work executes the stats aggregation job.
// Leader election is handled by River's periodic job scheduler.
func (w *StatsAggregationWorker) Work(ctx context.Context, job *river.Job[StatsAggregationArgs]) error {
	w.logger.Info("starting stats aggregation", slog.Int64("job_id", job.ID))
	start := time.Now()

	stats, err := w.collectStats(ctx)
	if err != nil {
		w.logger.Error("failed to collect stats",
			slog.Int64("job_id", job.ID),
			slog.Any("error", err),
		)
		return fmt.Errorf("collect stats: %w", err)
	}

	if err := w.persistStats(ctx, stats); err != nil {
		w.logger.Error("failed to persist stats",
			slog.Int64("job_id", job.ID),
			slog.Any("error", err),
		)
		return fmt.Errorf("persist stats: %w", err)
	}

	w.logger.Info("stats aggregation completed",
		slog.Int64("job_id", job.ID),
		slog.Duration("elapsed", time.Since(start)),
		slog.Int("stats_count", len(stats)),
	)

	return nil
}

// statEntry pairs a stat key with its computed value.
type statEntry struct {
	key   string
	value int64
}

// collectStats runs all aggregate queries concurrently and returns the results.
func (w *StatsAggregationWorker) collectStats(ctx context.Context) ([]statEntry, error) {
	if w.queries == nil {
		return nil, fmt.Errorf("stats_worker: queries is nil (dependency not injected)")
	}

	type queryFunc func(context.Context) (int64, error)

	// Map each stat key to its aggregate query.
	queries := []struct {
		key string
		fn  queryFunc
	}{
		{StatTotalUsers, w.queries.CountActiveUsers},
		{StatActiveUsers24h, w.queries.CountActiveUsersLast24h},
		{StatTotalLibraries, w.queries.CountTotalLibraries},
		{StatTotalMovies, w.queries.CountTotalMovies},
		{StatTotalSeries, w.queries.CountTotalSeries},
		{StatTotalEpisodes, w.queries.CountTotalEpisodes},
		{StatTotalMoviePlays, w.queries.CountTotalMovieWatches},
		{StatTotalEpisodePlays, w.queries.CountTotalEpisodeWatches},
		{StatMovieWatchSeconds, w.queries.SumMovieWatchDurationSeconds},
		{StatEpisodeWatchSeconds, w.queries.SumEpisodeWatchDurationSeconds},
	}

	var mu sync.Mutex
	results := make([]statEntry, 0, len(queries))

	g, gctx := errgroup.WithContext(ctx)
	for _, q := range queries {
		g.Go(func() error {
			val, err := q.fn(gctx)
			if err != nil {
				return fmt.Errorf("query %s: %w", q.key, err)
			}
			mu.Lock()
			results = append(results, statEntry{key: q.key, value: val})
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

// persistStats upserts all collected stats into shared.server_stats.
func (w *StatsAggregationWorker) persistStats(ctx context.Context, stats []statEntry) error {
	for _, s := range stats {
		if err := w.queries.UpsertServerStat(ctx, db.UpsertServerStatParams{
			StatKey:   s.key,
			StatValue: s.value,
		}); err != nil {
			return fmt.Errorf("upsert %s: %w", s.key, err)
		}
	}
	return nil
}
