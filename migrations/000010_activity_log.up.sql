-- 000010_activity_log.up.sql
-- Activity log for auditing and history

CREATE TYPE activity_type AS ENUM (
    'user_login', 'user_logout', 'user_created', 'user_deleted',
    'playback_start', 'playback_stop', 'playback_progress',
    'library_scan_start', 'library_scan_complete',
    'item_added', 'item_removed', 'item_updated',
    'system_start', 'system_stop', 'system_update'
);

CREATE TABLE activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    type activity_type NOT NULL,
    item_id UUID REFERENCES media_items(id) ON DELETE SET NULL,
    severity VARCHAR(20) NOT NULL DEFAULT 'info',  -- info, warning, error
    overview TEXT,
    short_overview VARCHAR(500),
    data JSONB NOT NULL DEFAULT '{}',         -- Additional structured data
    ip_address INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_activity_log_user_id ON activity_log(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_activity_log_type ON activity_log(type);
CREATE INDEX idx_activity_log_item_id ON activity_log(item_id) WHERE item_id IS NOT NULL;
CREATE INDEX idx_activity_log_severity ON activity_log(severity) WHERE severity != 'info';
CREATE INDEX idx_activity_log_created_at ON activity_log(created_at DESC);

-- Composite index for user activity timeline
CREATE INDEX idx_activity_log_user_timeline ON activity_log(user_id, created_at DESC) 
    WHERE user_id IS NOT NULL;
