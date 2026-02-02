-- OIDC Provider Configuration
-- Stores configuration for each OIDC provider (Authentik, Keycloak, generic, etc.)

CREATE TABLE IF NOT EXISTS shared.oidc_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Provider identification
    name VARCHAR(100) NOT NULL UNIQUE,           -- Internal name (e.g., "authentik", "keycloak")
    display_name VARCHAR(255) NOT NULL,          -- User-facing name
    provider_type VARCHAR(50) NOT NULL DEFAULT 'generic', -- generic, authentik, keycloak, etc.

    -- OIDC Configuration
    issuer_url VARCHAR(500) NOT NULL,            -- OIDC issuer URL
    client_id VARCHAR(255) NOT NULL,
    client_secret_encrypted BYTEA NOT NULL,      -- Encrypted client secret

    -- Endpoints (optional - can be discovered)
    authorization_endpoint VARCHAR(500),
    token_endpoint VARCHAR(500),
    userinfo_endpoint VARCHAR(500),
    jwks_uri VARCHAR(500),
    end_session_endpoint VARCHAR(500),

    -- Scopes
    scopes TEXT[] NOT NULL DEFAULT ARRAY['openid', 'profile', 'email'],

    -- Claim mappings (JSON)
    claim_mappings JSONB NOT NULL DEFAULT '{
        "username": "preferred_username",
        "email": "email",
        "name": "name",
        "picture": "picture",
        "roles": "groups"
    }'::jsonb,

    -- Role mappings (provider role -> revenge role)
    role_mappings JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- Behavior settings
    auto_create_users BOOLEAN NOT NULL DEFAULT true,
    update_user_info BOOLEAN NOT NULL DEFAULT true,
    allow_linking BOOLEAN NOT NULL DEFAULT true,     -- Allow existing users to link

    -- Status
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    is_default BOOLEAN NOT NULL DEFAULT false,       -- Default provider for SSO

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Ensure only one default provider
CREATE UNIQUE INDEX idx_oidc_providers_single_default
    ON shared.oidc_providers (is_default)
    WHERE is_default = true;

-- Index for lookups
CREATE INDEX idx_oidc_providers_enabled ON shared.oidc_providers (is_enabled) WHERE is_enabled = true;
CREATE INDEX idx_oidc_providers_name ON shared.oidc_providers (name);

-- OIDC User Links
-- Links Revenge users to OIDC provider subjects

CREATE TABLE IF NOT EXISTS shared.oidc_user_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Link to user
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Link to provider
    provider_id UUID NOT NULL REFERENCES shared.oidc_providers(id) ON DELETE CASCADE,

    -- Provider subject (unique ID from provider)
    subject VARCHAR(500) NOT NULL,               -- sub claim from ID token

    -- Cached user info from provider
    email VARCHAR(255),
    name VARCHAR(255),
    picture_url VARCHAR(500),

    -- Provider tokens (encrypted)
    access_token_encrypted BYTEA,
    refresh_token_encrypted BYTEA,
    token_expires_at TIMESTAMPTZ,

    -- Metadata
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Each user can only link to a provider once
    CONSTRAINT unique_user_provider UNIQUE (user_id, provider_id),
    -- Each subject can only be linked once per provider
    CONSTRAINT unique_provider_subject UNIQUE (provider_id, subject)
);

-- Indexes
CREATE INDEX idx_oidc_user_links_user_id ON shared.oidc_user_links (user_id);
CREATE INDEX idx_oidc_user_links_provider_id ON shared.oidc_user_links (provider_id);
CREATE INDEX idx_oidc_user_links_subject ON shared.oidc_user_links (provider_id, subject);

-- OIDC State Storage
-- Temporary storage for OAuth2 state during auth flow

CREATE TABLE IF NOT EXISTS shared.oidc_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- State token (random string)
    state VARCHAR(255) NOT NULL UNIQUE,

    -- PKCE verifier (for code_challenge)
    code_verifier VARCHAR(128),

    -- Provider being used
    provider_id UUID NOT NULL REFERENCES shared.oidc_providers(id) ON DELETE CASCADE,

    -- Optional: user linking flow
    user_id UUID REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Redirect URL after auth
    redirect_url VARCHAR(500),

    -- Expiration (short-lived)
    expires_at TIMESTAMPTZ NOT NULL,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for state lookup
CREATE INDEX idx_oidc_states_state ON shared.oidc_states (state);
-- Index for cleanup
CREATE INDEX idx_oidc_states_expires_at ON shared.oidc_states (expires_at);

-- Trigger to update updated_at
CREATE OR REPLACE FUNCTION shared.update_oidc_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_oidc_providers_updated_at
    BEFORE UPDATE ON shared.oidc_providers
    FOR EACH ROW
    EXECUTE FUNCTION shared.update_oidc_updated_at();

CREATE TRIGGER trigger_oidc_user_links_updated_at
    BEFORE UPDATE ON shared.oidc_user_links
    FOR EACH ROW
    EXECUTE FUNCTION shared.update_oidc_updated_at();
