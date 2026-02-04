-- Rollback Migration 000028: Re-add nonce column to user_totp_secrets
-- Note: This will require re-extracting nonces from existing encrypted_secrets

-- Add back the nonce column
ALTER TABLE public.user_totp_secrets ADD COLUMN IF NOT EXISTS nonce BYTEA;

-- Extract nonce from encrypted_secret (first 12 bytes)
UPDATE public.user_totp_secrets SET nonce = substring(encrypted_secret from 1 for 12);

-- Make nonce NOT NULL after populating
ALTER TABLE public.user_totp_secrets ALTER COLUMN nonce SET NOT NULL;

-- Re-add the constraint
ALTER TABLE public.user_totp_secrets ADD CONSTRAINT nonce_correct_size CHECK (length(nonce) = 12);

-- Update comment
COMMENT ON COLUMN public.user_totp_secrets.nonce IS 'GCM nonce used for encryption (12 bytes)';
