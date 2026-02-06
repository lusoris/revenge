-- MFA user settings table
-- Migration 000019: Create user_mfa_settings table

CREATE TABLE IF NOT EXISTS public.user_mfa_settings (
    user_id UUID PRIMARY KEY REFERENCES shared.users(id) ON DELETE CASCADE,

    -- MFA method enablement
    totp_enabled BOOLEAN NOT NULL DEFAULT false,
    webauthn_enabled BOOLEAN NOT NULL DEFAULT false,
    backup_codes_generated BOOLEAN NOT NULL DEFAULT false,

    -- Enforcement
    require_mfa BOOLEAN NOT NULL DEFAULT false,  -- Force MFA for this user (admin setting)

    -- Remember device settings
    remember_device_enabled BOOLEAN NOT NULL DEFAULT true,
    remember_device_duration_days INTEGER NOT NULL DEFAULT 30,

    -- Trusted devices (stored as JSON array of device fingerprints)
    trusted_devices JSONB NOT NULL DEFAULT '[]'::jsonb,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT remember_duration_positive CHECK (remember_device_duration_days > 0),
    CONSTRAINT trusted_devices_is_array CHECK (jsonb_typeof(trusted_devices) = 'array')
);

-- Indexes
CREATE INDEX idx_mfa_settings_require_mfa ON public.user_mfa_settings(user_id) WHERE require_mfa = true;
CREATE INDEX idx_mfa_settings_totp_enabled ON public.user_mfa_settings(user_id) WHERE totp_enabled = true;
CREATE INDEX idx_mfa_settings_webauthn_enabled ON public.user_mfa_settings(user_id) WHERE webauthn_enabled = true;

COMMENT ON TABLE public.user_mfa_settings IS 'Per-user MFA configuration and enforcement settings';
COMMENT ON COLUMN public.user_mfa_settings.totp_enabled IS 'Whether TOTP (authenticator app) is enabled for this user';
COMMENT ON COLUMN public.user_mfa_settings.webauthn_enabled IS 'Whether WebAuthn (passkeys/security keys) is enabled for this user';
COMMENT ON COLUMN public.user_mfa_settings.backup_codes_generated IS 'Whether backup recovery codes have been generated';
COMMENT ON COLUMN public.user_mfa_settings.require_mfa IS 'Admin override to force MFA for this user';
COMMENT ON COLUMN public.user_mfa_settings.trusted_devices IS 'JSON array of device fingerprints that skip MFA challenge';
