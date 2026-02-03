-- Remove MFA tracking from sessions table

DROP INDEX IF EXISTS shared.idx_sessions_mfa_verified;

ALTER TABLE shared.sessions
    DROP COLUMN IF EXISTS mfa_verified_at,
    DROP COLUMN IF EXISTS mfa_verified;
