package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/infra/database/db"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatsAggregationArgs_Kind(t *testing.T) {
	args := StatsAggregationArgs{}
	assert.Equal(t, StatsAggregationJobKind, args.Kind())
	assert.Equal(t, "stats_aggregation", args.Kind())
}

func TestStatsAggregationArgs_InsertOpts(t *testing.T) {
	args := StatsAggregationArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueLow, opts.Queue)
}

func TestStatKeyConstants_AreUnique(t *testing.T) {
	keys := []string{
		StatTotalUsers, StatActiveUsers24h, StatTotalLibraries, StatTotalMovies,
		StatTotalSeries, StatTotalEpisodes, StatTotalMoviePlays, StatTotalEpisodePlays,
		StatMovieWatchSeconds, StatEpisodeWatchSeconds,
	}
	seen := make(map[string]bool)
	for _, k := range keys {
		assert.False(t, seen[k], "duplicate stat key: %s", k)
		seen[k] = true
	}
}

func TestStatKeyConstants_Values(t *testing.T) {
	assert.Equal(t, "total_users", StatTotalUsers)
	assert.Equal(t, "active_users_24h", StatActiveUsers24h)
	assert.Equal(t, "total_libraries", StatTotalLibraries)
	assert.Equal(t, "total_movies", StatTotalMovies)
	assert.Equal(t, "total_series", StatTotalSeries)
	assert.Equal(t, "total_episodes", StatTotalEpisodes)
	assert.Equal(t, "total_movie_plays", StatTotalMoviePlays)
	assert.Equal(t, "total_episode_plays", StatTotalEpisodePlays)
	assert.Equal(t, "movie_watch_seconds", StatMovieWatchSeconds)
	assert.Equal(t, "episode_watch_seconds", StatEpisodeWatchSeconds)
}

func TestNewStatsAggregationWorker(t *testing.T) {
	logger := logging.NewTestLogger()
	worker := NewStatsAggregationWorker(&db.Queries{}, logger)
	require.NotNil(t, worker)
	assert.NotNil(t, worker.queries)
	assert.NotNil(t, worker.logger)
}

func TestStatsAggregationWorker_Timeout(t *testing.T) {
	logger := logging.NewTestLogger()
	worker := NewStatsAggregationWorker(&db.Queries{}, logger)
	timeout := worker.Timeout(&river.Job[StatsAggregationArgs]{})
	assert.Equal(t, 2*time.Minute, timeout)
}

func TestStatsAggregationWorker_Work_NilQueriesPanics(t *testing.T) {
	logger := logging.NewTestLogger()
	worker := NewStatsAggregationWorker(nil, logger)
	job := &river.Job[StatsAggregationArgs]{
		JobRow: &rivertype.JobRow{ID: 1},
	}
	assert.Panics(t, func() {
		_ = worker.Work(context.Background(), job)
	})
}

func TestModule_NotNil(t *testing.T) {
	assert.NotNil(t, Module)
}

func TestRegisterStatsAggregationWorker(t *testing.T) {
	logger := logging.NewTestLogger()
	workers := river.NewWorkers()
	worker := NewStatsAggregationWorker(&db.Queries{}, logger)
	registerStatsAggregationWorker(workers, worker)
}

func TestStatsAggregationArgs_ImplementsJobArgs(t *testing.T) {
	var _ river.JobArgs = StatsAggregationArgs{}
}

func TestStatsAggregationWorker_ImplementsWorker(t *testing.T) {
	logger := logging.NewTestLogger()
	worker := NewStatsAggregationWorker(&db.Queries{}, logger)
	var _ river.Worker[StatsAggregationArgs] = worker
}

func TestStatsAggregationJobKind_IsStable(t *testing.T) {
	assert.Equal(t, "stats_aggregation", StatsAggregationJobKind)
}

func TestStatEntry_Fields(t *testing.T) {
	entry := statEntry{key: "total_movies", value: 42}
	assert.Equal(t, "total_movies", entry.key)
	assert.Equal(t, int64(42), entry.value)
}
