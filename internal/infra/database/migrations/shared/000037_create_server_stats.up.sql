-- Server-wide pre-computed statistics.
-- A simple key-value store for aggregate stats, refreshed periodically
-- by the stats_aggregation River job.
CREATE TABLE IF NOT EXISTS shared.server_stats (
    stat_key    TEXT        PRIMARY KEY,
    stat_value  BIGINT      NOT NULL DEFAULT 0,
    computed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  shared.server_stats IS 'Pre-computed server-wide aggregate statistics';
COMMENT ON COLUMN shared.server_stats.stat_key    IS 'Unique metric identifier (e.g. total_users, total_movies)';
COMMENT ON COLUMN shared.server_stats.stat_value  IS 'Integer metric value';
COMMENT ON COLUMN shared.server_stats.computed_at IS 'When this value was last recomputed';
