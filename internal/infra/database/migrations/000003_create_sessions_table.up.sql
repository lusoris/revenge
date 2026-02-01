-- Create sessions table in shared schema
-- Tracks user sessions and refresh tokens

CREATE TABLE shared.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- User reference
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Session tokens
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    access_token_hash VARCHAR(255), -- Optional: track active access tokens

    -- Session metadata
    ip_address INET,
    user_agent TEXT,
    device_name VARCHAR(255),

    -- Expiry
    expires_at TIMESTAMPTZ NOT NULL,

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT true,
    revoked_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_sessions_user_id ON shared.sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON shared.sessions(refresh_token) WHERE is_active = true;
CREATE INDEX idx_sessions_expires_at ON shared.sessions(expires_at);
CREATE INDEX idx_sessions_is_active ON shared.sessions(is_active);

-- Comments
COMMENT ON TABLE shared.sessions IS 'User sessions and refresh tokens';
COMMENT ON COLUMN shared.sessions.id IS 'Unique session identifier';
COMMENT ON COLUMN shared.sessions.user_id IS 'User who owns this session';
COMMENT ON COLUMN shared.sessions.refresh_token IS 'Refresh token (secure random string)';
COMMENT ON COLUMN shared.sessions.access_token_hash IS 'Optional hash of current access token';
COMMENT ON COLUMN shared.sessions.ip_address IS 'IP address of the session';
COMMENT ON COLUMN shared.sessions.user_agent IS 'User agent string';
COMMENT ON COLUMN shared.sessions.device_name IS 'Friendly device name (e.g., "iPhone 12")';
COMMENT ON COLUMN shared.sessions.expires_at IS 'When the refresh token expires';
COMMENT ON COLUMN shared.sessions.is_active IS 'Whether the session is active';
COMMENT ON COLUMN shared.sessions.revoked_at IS 'When the session was revoked (if applicable)';
COMMENT ON COLUMN shared.sessions.last_used_at IS 'When the session was last used';
