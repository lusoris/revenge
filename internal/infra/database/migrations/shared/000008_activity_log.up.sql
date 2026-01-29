-- 000008_activity_log.up.sql
-- Activity and audit logging for security and user actions

BEGIN;

-- Activity types enum
CREATE TYPE activity_type AS ENUM (
    'user_login',
    'user_logout',
    'user_created',
    'user_updated',
    'user_deleted',
    'password_changed',
    'session_created',
    'session_expired',
    'library_created',
    'library_updated',
    'library_deleted',
    'library_scanned',
    'content_played',
    'content_rated',
    'settings_changed',
    'api_error',
    'security_event'
);

-- Activity severity levels
CREATE TYPE activity_severity AS ENUM ('info', 'warning', 'error', 'critical');

-- Activity log table
CREATE TABLE activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,  -- NULL for system events
    type activity_type NOT NULL,
    severity activity_severity NOT NULL DEFAULT 'info',
    message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',  -- Additional context (IP, user agent, etc.)
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_activity_log_user_id ON activity_log(user_id);
CREATE INDEX idx_activity_log_type ON activity_log(type);
CREATE INDEX idx_activity_log_severity ON activity_log(severity);
CREATE INDEX idx_activity_log_created_at ON activity_log(created_at DESC);
CREATE INDEX idx_activity_log_user_created ON activity_log(user_id, created_at DESC) WHERE user_id IS NOT NULL;

-- Composite index for filtering by type and date
CREATE INDEX idx_activity_log_type_date ON activity_log(type, created_at DESC);

COMMIT;
