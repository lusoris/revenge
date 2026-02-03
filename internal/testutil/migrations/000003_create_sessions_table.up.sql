-- Create sessions table in shared schema
-- User session management with JWT tokens

CREATE TABLE IF NOT EXISTS shared.sessions (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- User reference
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Session data
    token_hash TEXT NOT NULL UNIQUE, -- SHA256 hash of JWT token
    refresh_token_hash TEXT UNIQUE, -- SHA256 hash of refresh token (optional)

    -- Session metadata
    ip_address INET,
    user_agent TEXT,
    device_name VARCHAR(255),

    -- Scopes and permissions
    scopes TEXT[] DEFAULT ARRAY[]::TEXT[], -- e.g., ['legacy:read'] for QAR access

    -- Session lifecycle
    expires_at TIMESTAMPTZ NOT NULL,
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Revocation
    revoked_at TIMESTAMPTZ,
    revoke_reason TEXT
);

-- Indexes for session lookups and cleanup
CREATE INDEX idx_sessions_user_id ON shared.sessions(user_id) WHERE revoked_at IS NULL;
CREATE INDEX idx_sessions_token_hash ON shared.sessions(token_hash) WHERE revoked_at IS NULL;
CREATE INDEX idx_sessions_expires_at ON shared.sessions(expires_at) WHERE revoked_at IS NULL;
CREATE INDEX idx_sessions_last_activity ON shared.sessions(last_activity_at DESC);

-- Comments
COMMENT ON TABLE shared.sessions IS 'User sessions with JWT token management';
COMMENT ON COLUMN shared.sessions.token_hash IS 'SHA256 hash of JWT access token for validation';
COMMENT ON COLUMN shared.sessions.refresh_token_hash IS 'SHA256 hash of refresh token for token rotation';
COMMENT ON COLUMN shared.sessions.scopes IS 'OAuth2-style scopes (e.g., legacy:read for QAR access)';
