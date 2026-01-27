-- 000004_oidc.up.sql
-- OIDC provider configuration and user links

-- OIDC Providers (Keycloak, Authentik, Auth0, etc.)
CREATE TABLE oidc_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,       -- Internal name (keycloak, authentik, etc.)
    display_name VARCHAR(255) NOT NULL,       -- UI display name
    issuer_url VARCHAR(512) NOT NULL,         -- OIDC issuer URL
    client_id VARCHAR(255) NOT NULL,
    client_secret_encrypted BYTEA NOT NULL,   -- Encrypted with server key
    scopes TEXT[] NOT NULL DEFAULT ARRAY['openid', 'profile', 'email'],
    enabled BOOLEAN NOT NULL DEFAULT true,
    auto_create_users BOOLEAN NOT NULL DEFAULT true,
    default_admin BOOLEAN NOT NULL DEFAULT false,  -- New users are admins
    claim_mappings JSONB NOT NULL DEFAULT '{}',    -- Custom claim mappings
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_oidc_providers_enabled ON oidc_providers(enabled) WHERE enabled = true;

CREATE TRIGGER update_oidc_providers_updated_at
    BEFORE UPDATE ON oidc_providers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Links between users and OIDC identities
CREATE TABLE oidc_user_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES oidc_providers(id) ON DELETE CASCADE,
    subject VARCHAR(255) NOT NULL,            -- OIDC 'sub' claim
    email VARCHAR(255),                       -- Email from OIDC (may differ from user email)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,
    UNIQUE(provider_id, subject)
);

CREATE INDEX idx_oidc_user_links_user_id ON oidc_user_links(user_id);
CREATE INDEX idx_oidc_user_links_provider_subject ON oidc_user_links(provider_id, subject);
