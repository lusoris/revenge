-- OIDC Providers: External authentication providers
CREATE TABLE oidc_providers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,           -- Display name
    slug            VARCHAR(50) NOT NULL UNIQUE,     -- URL-safe identifier
    enabled         BOOLEAN NOT NULL DEFAULT true,

    -- OIDC Configuration
    issuer_url          TEXT NOT NULL,               -- Discovery URL base
    client_id           TEXT NOT NULL,
    client_secret_enc   BYTEA NOT NULL,              -- Encrypted client secret
    scopes              TEXT[] NOT NULL DEFAULT ARRAY['openid', 'profile', 'email'],

    -- Claim mapping (dot notation for nested claims)
    claim_mapping       JSONB NOT NULL DEFAULT '{
        "sub": "sub",
        "email": "email",
        "name": "name",
        "username": "preferred_username",
        "groups": "groups"
    }'::jsonb,

    -- Role mapping (groups -> roles)
    role_mapping        JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- User provisioning
    auto_provision      BOOLEAN NOT NULL DEFAULT true,
    default_role        VARCHAR(50) NOT NULL DEFAULT 'user',

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trigger for updated_at
CREATE TRIGGER oidc_providers_updated_at
    BEFORE UPDATE ON oidc_providers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- OIDC User Links: Connect Revenge users to OIDC identities
CREATE TABLE oidc_user_links (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_id     UUID NOT NULL REFERENCES oidc_providers(id) ON DELETE CASCADE,
    subject         VARCHAR(255) NOT NULL,           -- OIDC 'sub' claim

    -- Cached claims (for display/sync)
    email           VARCHAR(255),
    name            VARCHAR(255),
    groups          TEXT[],

    -- Timestamps
    linked_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at   TIMESTAMPTZ,

    UNIQUE (provider_id, subject)
);

-- Indexes
CREATE INDEX idx_oidc_user_links_user ON oidc_user_links(user_id);
CREATE INDEX idx_oidc_user_links_provider ON oidc_user_links(provider_id);
