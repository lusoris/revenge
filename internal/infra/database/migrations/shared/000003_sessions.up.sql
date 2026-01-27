-- 000003_sessions.up.sql
-- Sessions table - user authentication sessions

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA-256 of access token
    refresh_token_hash VARCHAR(64) UNIQUE,   -- SHA-256 of refresh token
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    client_name VARCHAR(255),
    client_version VARCHAR(50),
    ip_address INET,
    expires_at TIMESTAMPTZ NOT NULL,
    refresh_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for session lookups
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_refresh_token_hash ON sessions(refresh_token_hash) 
    WHERE refresh_token_hash IS NOT NULL;
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Index for cleanup of expired sessions
CREATE INDEX idx_sessions_expired ON sessions(expires_at) 
    WHERE expires_at < NOW();
