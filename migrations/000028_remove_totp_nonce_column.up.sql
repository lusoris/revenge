-- Migration 000028: Remove redundant nonce column from user_totp_secrets
-- The nonce is already prepended to the encrypted_secret by the AES-256-GCM encryption

-- Drop the constraint first
ALTER TABLE public.user_totp_secrets DROP CONSTRAINT IF EXISTS nonce_correct_size;

-- Drop the column
ALTER TABLE public.user_totp_secrets DROP COLUMN IF EXISTS nonce;

-- Update comment
COMMENT ON COLUMN public.user_totp_secrets.encrypted_secret IS 'AES-256-GCM encrypted base32-encoded TOTP secret (nonce prepended)';
