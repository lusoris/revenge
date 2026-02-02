-- Create email_verification_tokens table in shared schema
-- Stores one-time tokens for email verification

CREATE TABLE IF NOT EXISTS shared.email_verification_tokens (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Owner
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Token data (hashed for security)
    token_hash TEXT NOT NULL UNIQUE, -- SHA-256 hash of the actual token

    -- Email to verify (supports email change flow)
    email VARCHAR(255) NOT NULL, -- Email being verified

    -- Request metadata
    ip_address INET, -- IP address when verification was requested
    user_agent TEXT, -- User agent string

    -- Token lifecycle
    expires_at TIMESTAMPTZ NOT NULL, -- Typically 24 hours from creation
    verified_at TIMESTAMPTZ, -- When token was used (prevents reuse)

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_email_verification_tokens_user_id ON shared.email_verification_tokens(user_id) WHERE verified_at IS NULL;
CREATE INDEX idx_email_verification_tokens_token_hash ON shared.email_verification_tokens(token_hash) WHERE verified_at IS NULL;
CREATE INDEX idx_email_verification_tokens_expires_at ON shared.email_verification_tokens(expires_at) WHERE verified_at IS NULL;
CREATE INDEX idx_email_verification_tokens_email ON shared.email_verification_tokens(email) WHERE verified_at IS NULL;

-- Comments
COMMENT ON TABLE shared.email_verification_tokens IS 'One-time tokens for email verification and email change flow';
COMMENT ON COLUMN shared.email_verification_tokens.token_hash IS 'SHA-256 hash of the verification token (never store plaintext)';
COMMENT ON COLUMN shared.email_verification_tokens.email IS 'Email address being verified (may differ from user.email during change)';
COMMENT ON COLUMN shared.email_verification_tokens.verified_at IS 'Timestamp when token was used (prevents reuse)';
