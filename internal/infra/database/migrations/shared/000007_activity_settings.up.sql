-- Activity Log: Audit trail for important actions
CREATE TABLE activity_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    profile_id      UUID REFERENCES profiles(id) ON DELETE SET NULL,

    -- Action details
    action          VARCHAR(100) NOT NULL,           -- e.g., 'login', 'playback.start', 'library.scan'
    module          VARCHAR(50),                     -- Content module (movie, tvshow, etc.)
    item_id         UUID,                            -- Related item ID (if applicable)
    item_type       VARCHAR(50),                     -- Item type for display

    -- Context
    details         JSONB,                           -- Additional action-specific data
    ip_address      INET,
    user_agent      TEXT,

    -- Severity
    severity        VARCHAR(20) NOT NULL DEFAULT 'info',  -- debug, info, warn, error

    -- Timestamp
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_activity_log_user ON activity_log(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_activity_log_created ON activity_log(created_at DESC);
CREATE INDEX idx_activity_log_action ON activity_log(action);
CREATE INDEX idx_activity_log_module ON activity_log(module) WHERE module IS NOT NULL;

-- Partitioning hint: Consider partitioning by month for large deployments
-- CREATE TABLE activity_log_YYYY_MM PARTITION OF activity_log FOR VALUES FROM (...) TO (...);

-- Server Settings: Persisted configuration
CREATE TABLE server_settings (
    key             VARCHAR(255) PRIMARY KEY,
    value           JSONB NOT NULL,
    description     TEXT,
    updated_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default settings
INSERT INTO server_settings (key, value, description) VALUES
    ('server.name', '"Revenge Media Server"', 'Server display name'),
    ('server.public', 'false', 'Allow public access without login'),
    ('server.registration_enabled', 'false', 'Allow new user registration'),
    ('server.default_language', '"en"', 'Default language for new users'),
    ('transcoding.enabled', 'true', 'Enable transcoding via Blackbeard'),
    ('transcoding.blackbeard_url', '""', 'Blackbeard service URL'),
    ('library.scan_on_startup', 'true', 'Scan libraries on server start'),
    ('library.watch_filesystem', 'true', 'Watch for file changes'),
    ('cache.enabled', 'true', 'Enable caching layer'),
    ('search.enabled', 'true', 'Enable Typesense search'),
    ('adult.globally_enabled', 'false', 'Allow adult content access server-wide');
