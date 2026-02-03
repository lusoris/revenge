-- Add MFA tracking to sessions table
-- Tracks whether a session has passed MFA verification

ALTER TABLE shared.sessions
    ADD COLUMN mfa_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN mfa_verified_at TIMESTAMPTZ;

-- Index for querying MFA-verified sessions
CREATE INDEX idx_sessions_mfa_verified ON shared.sessions(user_id, mfa_verified) 
    WHERE revoked_at IS NULL;

-- Comments
COMMENT ON COLUMN shared.sessions.mfa_verified IS 'Whether this session has passed MFA verification';
COMMENT ON COLUMN shared.sessions.mfa_verified_at IS 'Timestamp when MFA was verified for this session';
