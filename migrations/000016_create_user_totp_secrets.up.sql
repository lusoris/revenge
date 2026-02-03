-- TOTP secrets table (encrypted)
-- Migration 000016: Create user_totp_secrets table

CREATE TABLE IF NOT EXISTS public.user_totp_secrets (
    user_id UUID PRIMARY KEY REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Encrypted secret (AES-256-GCM encrypted base32 string)
    encrypted_secret BYTEA NOT NULL,

    -- Encryption metadata
    nonce BYTEA NOT NULL,  -- GCM nonce (12 bytes)

    -- Status
    verified_at TIMESTAMPTZ,  -- NULL until first successful verification
    enabled BOOLEAN NOT NULL DEFAULT false,

    -- Usage tracking
    last_used_at TIMESTAMPTZ,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure secret exists
    CONSTRAINT secret_not_empty CHECK (length(encrypted_secret) > 0),
    CONSTRAINT nonce_correct_size CHECK (length(nonce) = 12)
);

-- Index for enabled TOTP users
CREATE INDEX idx_totp_secrets_user_enabled ON public.user_totp_secrets(user_id) WHERE enabled = true;

-- Index for last used tracking
CREATE INDEX idx_totp_secrets_last_used ON public.user_totp_secrets(last_used_at DESC) WHERE last_used_at IS NOT NULL;

COMMENT ON TABLE public.user_totp_secrets IS 'TOTP (Time-based One-Time Password) secrets for multi-factor authentication';
COMMENT ON COLUMN public.user_totp_secrets.encrypted_secret IS 'AES-256-GCM encrypted base32-encoded TOTP secret';
COMMENT ON COLUMN public.user_totp_secrets.nonce IS 'GCM nonce used for encryption (12 bytes)';
COMMENT ON COLUMN public.user_totp_secrets.verified_at IS 'When the TOTP was first successfully verified (enrollment completion)';
