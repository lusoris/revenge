-- Add session_id column to auth_tokens to link refresh tokens to sessions.
-- This enables embedding session_id in JWT access tokens, which is required
-- for the /sessions/current endpoint to identify the active session.

ALTER TABLE shared.auth_tokens
ADD COLUMN session_id UUID REFERENCES shared.sessions (id) ON DELETE SET NULL;

-- Index for looking up auth tokens by session
CREATE INDEX idx_auth_tokens_session_id ON shared.auth_tokens (session_id)
WHERE
    session_id IS NOT NULL
    AND revoked_at IS NULL;

COMMENT ON COLUMN shared.auth_tokens.session_id IS 'Links this refresh token to the session record it was created with';
