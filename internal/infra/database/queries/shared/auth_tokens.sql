-- name: CreateAuthToken :one
INSERT INTO shared.auth_tokens (
    user_id,
    token_hash,
    token_type,
    device_name,
    device_fingerprint,
    ip_address,
    user_agent,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAuthTokenByHash :one
SELECT * FROM shared.auth_tokens
WHERE token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > NOW()
LIMIT 1;

-- name: GetAuthTokensByUserID :many
SELECT * FROM shared.auth_tokens
WHERE user_id = $1
  AND revoked_at IS NULL
  AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: GetAuthTokensByDeviceFingerprint :many
SELECT * FROM shared.auth_tokens
WHERE user_id = $1
  AND device_fingerprint = $2
  AND revoked_at IS NULL
  AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: UpdateAuthTokenLastUsed :exec
UPDATE shared.auth_tokens
SET last_used_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: RevokeAuthToken :exec
UPDATE shared.auth_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: RevokeAuthTokenByHash :exec
UPDATE shared.auth_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token_hash = $1;

-- name: RevokeAllUserAuthTokens :exec
UPDATE shared.auth_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE user_id = $1
  AND revoked_at IS NULL;

-- name: RevokeAllUserAuthTokensExcept :exec
UPDATE shared.auth_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE user_id = $1
  AND id != $2
  AND revoked_at IS NULL;

-- name: DeleteExpiredAuthTokens :exec
DELETE FROM shared.auth_tokens
WHERE expires_at < NOW();

-- name: DeleteRevokedAuthTokens :exec
DELETE FROM shared.auth_tokens
WHERE revoked_at IS NOT NULL
  AND revoked_at < NOW() - INTERVAL '30 days';

-- name: CountActiveAuthTokensByUser :one
SELECT COUNT(*) FROM shared.auth_tokens
WHERE user_id = $1
  AND revoked_at IS NULL
  AND expires_at > NOW();

-- Password Reset Tokens

-- name: CreatePasswordResetToken :one
INSERT INTO shared.password_reset_tokens (
    user_id,
    token_hash,
    ip_address,
    user_agent,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPasswordResetToken :one
SELECT * FROM shared.password_reset_tokens
WHERE token_hash = $1
  AND used_at IS NULL
  AND expires_at > NOW()
LIMIT 1;

-- name: MarkPasswordResetTokenUsed :exec
UPDATE shared.password_reset_tokens
SET used_at = NOW()
WHERE id = $1;

-- name: InvalidateUserPasswordResetTokens :exec
UPDATE shared.password_reset_tokens
SET used_at = NOW()
WHERE user_id = $1
  AND used_at IS NULL;

-- name: DeleteExpiredPasswordResetTokens :exec
DELETE FROM shared.password_reset_tokens
WHERE expires_at < NOW();

-- name: DeleteUsedPasswordResetTokens :exec
DELETE FROM shared.password_reset_tokens
WHERE used_at IS NOT NULL
  AND used_at < NOW() - INTERVAL '7 days';

-- Email Verification Tokens

-- name: CreateEmailVerificationToken :one
INSERT INTO shared.email_verification_tokens (
    user_id,
    token_hash,
    email,
    ip_address,
    user_agent,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetEmailVerificationToken :one
SELECT * FROM shared.email_verification_tokens
WHERE token_hash = $1
  AND verified_at IS NULL
  AND expires_at > NOW()
LIMIT 1;

-- name: MarkEmailVerificationTokenUsed :exec
UPDATE shared.email_verification_tokens
SET verified_at = NOW()
WHERE id = $1;

-- name: InvalidateUserEmailVerificationTokens :exec
UPDATE shared.email_verification_tokens
SET verified_at = NOW()
WHERE user_id = $1
  AND verified_at IS NULL;

-- name: InvalidateEmailVerificationTokensByEmail :exec
UPDATE shared.email_verification_tokens
SET verified_at = NOW()
WHERE email = $1
  AND verified_at IS NULL;

-- name: DeleteExpiredEmailVerificationTokens :exec
DELETE FROM shared.email_verification_tokens
WHERE expires_at < NOW();

-- name: DeleteVerifiedEmailTokens :exec
DELETE FROM shared.email_verification_tokens
WHERE verified_at IS NOT NULL
  AND verified_at < NOW() - INTERVAL '7 days';

-- Failed Login Attempts (Account Lockout / Rate Limiting)

-- name: RecordFailedLoginAttempt :exec
INSERT INTO shared.failed_login_attempts (
    username,
    ip_address
) VALUES (
    $1, $2
);

-- name: CountFailedLoginAttemptsByUsername :one
SELECT COUNT(*) FROM shared.failed_login_attempts
WHERE username = $1
  AND attempted_at > $2;

-- name: CountFailedLoginAttemptsByIP :one
SELECT COUNT(*) FROM shared.failed_login_attempts
WHERE ip_address = $1
  AND attempted_at > $2;

-- name: ClearFailedLoginAttemptsByUsername :exec
DELETE FROM shared.failed_login_attempts
WHERE username = $1;

-- name: DeleteOldFailedLoginAttempts :exec
DELETE FROM shared.failed_login_attempts
WHERE attempted_at < NOW() - INTERVAL '24 hours';
