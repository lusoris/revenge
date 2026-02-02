-- Create password_reset_tokens table in shared schema
-- Stores one-time tokens for password reset requests

CREATE TABLE IF NOT EXISTS shared.password_reset_tokens (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Owner
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Token data (hashed for security)
    token_hash TEXT NOT NULL UNIQUE, -- SHA-256 hash of the actual token

    -- Request metadata
    ip_address INET, -- IP address when reset was requested
    user_agent TEXT, -- User agent string

    -- Token lifecycle
    expires_at TIMESTAMPTZ NOT NULL, -- Typically 1 hour from creation
    used_at TIMESTAMPTZ, -- When token was consumed (prevents reuse)

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_password_reset_tokens_user_id ON shared.password_reset_tokens(user_id) WHERE used_at IS NULL;
CREATE INDEX idx_password_reset_tokens_token_hash ON shared.password_reset_tokens(token_hash) WHERE used_at IS NULL;
CREATE INDEX idx_password_reset_tokens_expires_at ON shared.password_reset_tokens(expires_at) WHERE used_at IS NULL;

-- Comments
COMMENT ON TABLE shared.password_reset_tokens IS 'One-time tokens for password reset flow';
COMMENT ON COLUMN shared.password_reset_tokens.token_hash IS 'SHA-256 hash of the reset token (never store plaintext)';
COMMENT ON COLUMN shared.password_reset_tokens.used_at IS 'Timestamp when token was used (prevents reuse)';
