-- Create failed_login_attempts table for account lockout / rate limiting
-- Tracks failed login attempts to prevent brute-force attacks
-- Part of A7.5: Service-Level Rate Limiting

CREATE TABLE IF NOT EXISTS shared.failed_login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    attempted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for counting attempts by username in time window
CREATE INDEX idx_failed_login_username ON shared.failed_login_attempts(username, attempted_at DESC);

-- Index for counting attempts by IP address in time window (future use for IP-based rate limiting)
CREATE INDEX idx_failed_login_ip ON shared.failed_login_attempts(ip_address, attempted_at DESC);

-- Add comment
COMMENT ON TABLE shared.failed_login_attempts IS 'Tracks failed login attempts for account lockout and rate limiting';
COMMENT ON COLUMN shared.failed_login_attempts.username IS 'Username or email used in the failed login attempt';
COMMENT ON COLUMN shared.failed_login_attempts.ip_address IS 'IP address from which the failed login attempt originated';
COMMENT ON COLUMN shared.failed_login_attempts.attempted_at IS 'Timestamp when the failed login attempt occurred';
COMMENT ON COLUMN shared.failed_login_attempts.created_at IS 'Timestamp when this record was created (for auditing)';
