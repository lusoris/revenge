-- Audit Log Redesign: Partitioned activity log for metadata auditing
-- Async writes via River, monthly partitions for performance
BEGIN;

-- Drop old activity_log if exists and recreate as partitioned
DROP TABLE IF EXISTS activity_log CASCADE;

-- Create partitioned activity_log table
CREATE TABLE activity_log (
    id              UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,  -- No FK to avoid write contention
    action          VARCHAR(50) NOT NULL,
    module          VARCHAR(50) NOT NULL,  -- 'movie', 'tvshow', 'qar', 'user', etc.
    entity_id       UUID NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,  -- 'movie', 'episode', 'expedition', etc.
    changes         JSONB NOT NULL DEFAULT '{}',  -- Field changes: {"field": {"old": x, "new": y}}
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id, created_at)  -- Partition key must be in PK
) PARTITION BY RANGE (created_at);

-- Create partitions for current and next few months
CREATE TABLE activity_log_2026_01 PARTITION OF activity_log
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE activity_log_2026_02 PARTITION OF activity_log
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE activity_log_2026_03 PARTITION OF activity_log
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE activity_log_2026_04 PARTITION OF activity_log
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE activity_log_2026_05 PARTITION OF activity_log
    FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE activity_log_2026_06 PARTITION OF activity_log
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');

-- Default partition for future dates (will be split by maintenance job)
CREATE TABLE activity_log_default PARTITION OF activity_log DEFAULT;

-- Indexes (created on parent, auto-applied to partitions)
CREATE INDEX idx_activity_log_entity ON activity_log(module, entity_type, entity_id);
CREATE INDEX idx_activity_log_user ON activity_log(user_id, created_at DESC);
CREATE INDEX idx_activity_log_action ON activity_log(action, created_at DESC);

-- Action types reference
COMMENT ON TABLE activity_log IS 'Partitioned audit log for metadata changes, written async via River';
COMMENT ON COLUMN activity_log.action IS 'Action type: metadata.edit, metadata.lock, metadata.unlock, metadata.refresh, image.upload, image.select, image.delete, content.delete, user.login, user.logout';
COMMENT ON COLUMN activity_log.module IS 'Module: movie, tvshow, qar, user, library, system';
COMMENT ON COLUMN activity_log.changes IS 'JSONB of field changes: {"field": {"old": value, "new": value}}';

COMMIT;
