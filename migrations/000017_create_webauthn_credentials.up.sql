-- WebAuthn/FIDO2 credentials table
-- Migration 000017: Create webauthn_credentials table

CREATE TABLE IF NOT EXISTS public.webauthn_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Credential data (from WebAuthn spec)
    credential_id BYTEA NOT NULL UNIQUE,  -- Raw credential ID
    public_key BYTEA NOT NULL,             -- COSE encoded public key

    -- Counter for clone detection (must always increment)
    sign_count INTEGER NOT NULL DEFAULT 0,
    clone_detected BOOLEAN NOT NULL DEFAULT false,  -- Flag if counter decreased

    -- Authenticator data
    aaguid BYTEA,                          -- Authenticator AAGUID (16 bytes, nullable for platform authenticators)
    attestation_type TEXT NOT NULL,        -- "none", "packed", "fido-u2f", "tpm", "android-key", "android-safetynet"

    -- Transports (for UX hints)
    transports TEXT[],                     -- ["usb", "nfc", "ble", "internal", "hybrid"]

    -- Flags from authenticator data
    backup_eligible BOOLEAN NOT NULL DEFAULT false,
    backup_state BOOLEAN NOT NULL DEFAULT false,
    user_present BOOLEAN NOT NULL DEFAULT true,
    user_verified BOOLEAN NOT NULL DEFAULT false,

    -- User-facing metadata
    name TEXT,                             -- User-given name ("My YubiKey", "Touch ID", "Pixel 8 Fingerprint")

    -- Status
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT credential_id_not_empty CHECK (length(credential_id) > 0),
    CONSTRAINT public_key_not_empty CHECK (length(public_key) > 0),
    CONSTRAINT sign_count_non_negative CHECK (sign_count >= 0),
    CONSTRAINT aaguid_correct_size CHECK (aaguid IS NULL OR length(aaguid) = 16)
);

-- Indexes
CREATE INDEX idx_webauthn_credentials_user ON public.webauthn_credentials(user_id);
CREATE INDEX idx_webauthn_credentials_credential_id ON public.webauthn_credentials(credential_id);
CREATE INDEX idx_webauthn_credentials_last_used ON public.webauthn_credentials(last_used_at DESC) WHERE last_used_at IS NOT NULL;
CREATE INDEX idx_webauthn_credentials_clone_detected ON public.webauthn_credentials(user_id, clone_detected) WHERE clone_detected = true;

COMMENT ON TABLE public.webauthn_credentials IS 'WebAuthn/FIDO2 credentials for passwordless and multi-factor authentication';
COMMENT ON COLUMN public.webauthn_credentials.credential_id IS 'Unique credential identifier from the authenticator';
COMMENT ON COLUMN public.webauthn_credentials.public_key IS 'COSE-encoded public key for verifying assertions';
COMMENT ON COLUMN public.webauthn_credentials.sign_count IS 'Signature counter for detecting cloned authenticators (must increment)';
COMMENT ON COLUMN public.webauthn_credentials.clone_detected IS 'Flag indicating potential authenticator cloning (counter decreased)';
COMMENT ON COLUMN public.webauthn_credentials.aaguid IS 'Authenticator Attestation GUID (identifies authenticator model)';
COMMENT ON COLUMN public.webauthn_credentials.transports IS 'Communication methods supported by authenticator';
