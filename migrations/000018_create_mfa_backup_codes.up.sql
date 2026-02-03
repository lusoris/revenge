-- MFA backup/recovery codes table
-- Migration 000018: Create mfa_backup_codes table

CREATE TABLE IF NOT EXISTS public.mfa_backup_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Argon2id hashed code (originally 16 characters hex-encoded)
    code_hash TEXT NOT NULL,

    -- Usage tracking
    used_at TIMESTAMPTZ,
    used_from_ip INET,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT code_hash_not_empty CHECK (length(code_hash) > 0)
);

-- Indexes
CREATE INDEX idx_backup_codes_user ON public.mfa_backup_codes(user_id);
CREATE INDEX idx_backup_codes_unused ON public.mfa_backup_codes(user_id) WHERE used_at IS NULL;
CREATE INDEX idx_backup_codes_created ON public.mfa_backup_codes(created_at DESC);

COMMENT ON TABLE public.mfa_backup_codes IS 'One-time backup codes for MFA account recovery';
COMMENT ON COLUMN public.mfa_backup_codes.code_hash IS 'Argon2id hash of the backup code (codes are 16 chars hex-encoded)';
COMMENT ON COLUMN public.mfa_backup_codes.used_at IS 'When this code was used (NULL if unused)';
COMMENT ON COLUMN public.mfa_backup_codes.used_from_ip IS 'IP address where code was used (for audit trail)';
