-- Create auth_tokens table in shared schema
-- Stores JWT refresh tokens for persistent sessions

CREATE TABLE IF NOT EXISTS shared.auth_tokens (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Owner
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Token data (hashed for security)
    token_hash TEXT NOT NULL UNIQUE, -- SHA-256 hash of the actual token
    token_type VARCHAR(20) NOT NULL DEFAULT 'refresh', -- 'refresh', 'access'

    -- Device tracking
    device_name VARCHAR(255), -- e.g., "iPhone 13", "Chrome on Windows"
    device_fingerprint TEXT, -- Browser/device fingerprint
    ip_address INET, -- IP address when token was created
    user_agent TEXT, -- User agent string

    -- Expiry and status
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ, -- When token was manually revoked
    last_used_at TIMESTAMPTZ, -- Track token activity

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_auth_tokens_user_id ON shared.auth_tokens(user_id) WHERE revoked_at IS NULL;
CREATE INDEX idx_auth_tokens_token_hash ON shared.auth_tokens(token_hash) WHERE revoked_at IS NULL;
CREATE INDEX idx_auth_tokens_expires_at ON shared.auth_tokens(expires_at) WHERE revoked_at IS NULL;
CREATE INDEX idx_auth_tokens_device ON shared.auth_tokens(user_id, device_fingerprint) WHERE revoked_at IS NULL;

-- Comments
COMMENT ON TABLE shared.auth_tokens IS 'JWT refresh tokens for persistent user sessions';
COMMENT ON COLUMN shared.auth_tokens.token_hash IS 'SHA-256 hash of the refresh token (never store plaintext)';
COMMENT ON COLUMN shared.auth_tokens.device_fingerprint IS 'Unique identifier for the device/browser';
COMMENT ON COLUMN shared.auth_tokens.last_used_at IS 'Timestamp when token was last used for refresh';
