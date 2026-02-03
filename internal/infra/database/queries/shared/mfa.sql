-- ============================================================================
-- TOTP (Time-based One-Time Password) Queries
-- ============================================================================

-- name: GetUserTOTPSecret :one
-- Get TOTP secret for a user
SELECT * FROM public.user_totp_secrets
WHERE user_id = $1;

-- name: CreateTOTPSecret :one
-- Create a new TOTP secret for a user
INSERT INTO public.user_totp_secrets (
    user_id,
    encrypted_secret,
    nonce,
    enabled
) VALUES (
    $1, $2, $3, false
) RETURNING *;

-- name: UpdateTOTPSecret :exec
-- Update TOTP secret (for re-enrollment)
UPDATE public.user_totp_secrets
SET encrypted_secret = $2,
    nonce = $3,
    verified_at = NULL,
    enabled = false,
    updated_at = NOW()
WHERE user_id = $1;

-- name: VerifyTOTPSecret :exec
-- Mark TOTP as verified and enabled
UPDATE public.user_totp_secrets
SET verified_at = NOW(),
    enabled = true,
    updated_at = NOW()
WHERE user_id = $1 AND verified_at IS NULL;

-- name: UpdateTOTPLastUsed :exec
-- Update last used timestamp for TOTP
UPDATE public.user_totp_secrets
SET last_used_at = NOW(),
    updated_at = NOW()
WHERE user_id = $1;

-- name: EnableTOTP :exec
-- Enable TOTP for a user (after verification)
UPDATE public.user_totp_secrets
SET enabled = true,
    updated_at = NOW()
WHERE user_id = $1;

-- name: DisableTOTP :exec
-- Disable TOTP for a user
UPDATE public.user_totp_secrets
SET enabled = false,
    updated_at = NOW()
WHERE user_id = $1;

-- name: DeleteTOTPSecret :exec
-- Delete TOTP secret for a user
DELETE FROM public.user_totp_secrets
WHERE user_id = $1;

-- ============================================================================
-- WebAuthn Credential Queries
-- ============================================================================

-- name: ListWebAuthnCredentials :many
-- List all WebAuthn credentials for a user
SELECT * FROM public.webauthn_credentials
WHERE user_id = $1
ORDER BY last_used_at DESC NULLS LAST, created_at DESC;

-- name: GetWebAuthnCredential :one
-- Get a specific WebAuthn credential by ID
SELECT * FROM public.webauthn_credentials
WHERE id = $1;

-- name: GetWebAuthnCredentialByCredentialID :one
-- Get a WebAuthn credential by its credential ID
SELECT * FROM public.webauthn_credentials
WHERE credential_id = $1;

-- name: CreateWebAuthnCredential :one
-- Create a new WebAuthn credential
INSERT INTO public.webauthn_credentials (
    user_id,
    credential_id,
    public_key,
    attestation_type,
    transports,
    backup_eligible,
    backup_state,
    user_present,
    user_verified,
    aaguid,
    name
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: UpdateWebAuthnCounter :exec
-- Update sign count and last used timestamp
UPDATE public.webauthn_credentials
SET sign_count = $2,
    last_used_at = NOW()
WHERE credential_id = $1;

-- name: MarkWebAuthnCloneDetected :exec
-- Mark a credential as potentially cloned
UPDATE public.webauthn_credentials
SET clone_detected = true
WHERE credential_id = $1;

-- name: UpdateWebAuthnCredentialName :exec
-- Update the user-facing name of a credential
UPDATE public.webauthn_credentials
SET name = $2
WHERE id = $1 AND user_id = $2;

-- name: DeleteWebAuthnCredential :exec
-- Delete a WebAuthn credential
DELETE FROM public.webauthn_credentials
WHERE id = $1 AND user_id = $2;

-- name: CountWebAuthnCredentials :one
-- Count WebAuthn credentials for a user
SELECT COUNT(*) FROM public.webauthn_credentials
WHERE user_id = $1;

-- ============================================================================
-- Backup Code Queries
-- ============================================================================

-- name: CreateBackupCodes :copyfrom
-- Bulk insert backup codes
INSERT INTO public.mfa_backup_codes (
    user_id,
    code_hash
) VALUES (
    $1, $2
);

-- name: GetUnusedBackupCodes :many
-- Get all unused backup codes for a user
SELECT * FROM public.mfa_backup_codes
WHERE user_id = $1 AND used_at IS NULL
ORDER BY created_at DESC;

-- name: GetBackupCodeByHash :one
-- Get a backup code by its hash
SELECT * FROM public.mfa_backup_codes
WHERE user_id = $1 AND code_hash = $2 AND used_at IS NULL;

-- name: UseBackupCode :exec
-- Mark a backup code as used
UPDATE public.mfa_backup_codes
SET used_at = NOW(),
    used_from_ip = $3
WHERE id = $1 AND user_id = $2 AND used_at IS NULL;

-- name: CountUnusedBackupCodes :one
-- Count unused backup codes for a user
SELECT COUNT(*) FROM public.mfa_backup_codes
WHERE user_id = $1 AND used_at IS NULL;

-- name: DeleteAllBackupCodes :exec
-- Delete all backup codes for a user (when regenerating)
DELETE FROM public.mfa_backup_codes
WHERE user_id = $1;

-- ============================================================================
-- MFA Settings Queries
-- ============================================================================

-- name: GetUserMFASettings :one
-- Get MFA settings for a user
SELECT * FROM public.user_mfa_settings
WHERE user_id = $1;

-- name: CreateUserMFASettings :one
-- Create MFA settings for a user
INSERT INTO public.user_mfa_settings (
    user_id,
    totp_enabled,
    webauthn_enabled,
    backup_codes_generated,
    require_mfa,
    remember_device_enabled,
    remember_device_duration_days,
    trusted_devices
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateMFASettingsTOTPEnabled :exec
-- Update TOTP enabled flag
UPDATE public.user_mfa_settings
SET totp_enabled = $2,
    updated_at = NOW()
WHERE user_id = $1;

-- name: UpdateMFASettingsWebAuthnEnabled :exec
-- Update WebAuthn enabled flag
UPDATE public.user_mfa_settings
SET webauthn_enabled = $2,
    updated_at = NOW()
WHERE user_id = $1;

-- name: UpdateMFASettingsBackupCodesGenerated :exec
-- Update backup codes generated flag
UPDATE public.user_mfa_settings
SET backup_codes_generated = $2,
    updated_at = NOW()
WHERE user_id = $1;

-- name: UpdateMFASettingsRequireMFA :exec
-- Update require MFA flag (admin setting)
UPDATE public.user_mfa_settings
SET require_mfa = $2,
    updated_at = NOW()
WHERE user_id = $1;

-- name: UpdateMFASettingsRememberDevice :exec
-- Update remember device settings
UPDATE public.user_mfa_settings
SET remember_device_enabled = $2,
    remember_device_duration_days = $3,
    updated_at = NOW()
WHERE user_id = $1;

-- name: AddTrustedDevice :exec
-- Add a device fingerprint to trusted devices
UPDATE public.user_mfa_settings
SET trusted_devices = trusted_devices || $2::jsonb,
    updated_at = NOW()
WHERE user_id = $1;

-- NOTE: RemoveTrustedDevice removed - implement in application code
-- Use ClearTrustedDevices to reset all devices

-- name: ClearTrustedDevices :exec
-- Clear all trusted devices for a user
UPDATE public.user_mfa_settings
SET trusted_devices = '[]'::jsonb,
    updated_at = NOW()
WHERE user_id = $1;

-- name: DeleteUserMFASettings :exec
-- Delete MFA settings for a user
DELETE FROM public.user_mfa_settings
WHERE user_id = $1;

-- ============================================================================
-- Combined Status Queries
-- ============================================================================

-- name: GetUserMFAStatus :one
-- Get comprehensive MFA status for a user
SELECT
    (SELECT COUNT(*) > 0 FROM public.user_totp_secrets WHERE user_totp_secrets.user_id = $1 AND enabled = true) as has_totp,
    (SELECT COUNT(*) FROM public.webauthn_credentials WHERE webauthn_credentials.user_id = $1) as webauthn_count,
    (SELECT COUNT(*) FROM public.mfa_backup_codes WHERE mfa_backup_codes.user_id = $1 AND used_at IS NULL) as unused_backup_codes,
    (SELECT require_mfa FROM public.user_mfa_settings WHERE user_mfa_settings.user_id = $1) as require_mfa;

-- name: HasAnyMFAMethod :one
-- Check if user has any MFA method enabled
SELECT (
    EXISTS(SELECT 1 FROM public.user_totp_secrets WHERE user_totp_secrets.user_id = $1 AND enabled = true)
    OR
    EXISTS(SELECT 1 FROM public.webauthn_credentials WHERE webauthn_credentials.user_id = $1)
) as has_mfa;
