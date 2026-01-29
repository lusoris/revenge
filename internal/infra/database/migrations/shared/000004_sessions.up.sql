-- Sessions table: Active login sessions
CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    profile_id      UUID REFERENCES profiles(id) ON DELETE SET NULL,
    token_hash      VARCHAR(255) NOT NULL UNIQUE,    -- SHA-256 of session token

    -- Device info
    device_name     VARCHAR(255),
    device_type     VARCHAR(50),                     -- browser, mobile, tv, desktop
    client_name     VARCHAR(100),                    -- App name (e.g., "Revenge Web", "Revenge iOS")
    client_version  VARCHAR(50),
    ip_address      INET,
    user_agent      TEXT,

    -- Session state
    is_active       BOOLEAN NOT NULL DEFAULT true,
    last_activity   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ NOT NULL,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at) WHERE is_active = true;
CREATE INDEX idx_sessions_last_activity ON sessions(last_activity);

-- API Keys table: For external integrations (Servarr, etc.)
CREATE TABLE api_keys (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,           -- User-defined name
    key_hash        VARCHAR(255) NOT NULL UNIQUE,    -- SHA-256 of API key
    key_prefix      VARCHAR(10) NOT NULL,            -- First 8 chars for identification

    -- Permissions
    scopes          TEXT[] NOT NULL DEFAULT '{}',    -- Array of permission scopes

    -- Usage tracking
    last_used_at    TIMESTAMPTZ,
    use_count       BIGINT NOT NULL DEFAULT 0,

    -- Timestamps
    expires_at      TIMESTAMPTZ,                     -- NULL = never expires
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_key_prefix ON api_keys(key_prefix);
