-- Create API keys table for programmatic access
-- Keys are SHA-256 hashed (like auth tokens)
-- Scopes define permissions (read, write, admin)

CREATE TABLE IF NOT EXISTS shared.api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    key_hash VARCHAR(64) NOT NULL UNIQUE, -- SHA-256 hash (32 bytes = 64 hex chars)
    key_prefix VARCHAR(16) NOT NULL, -- First 8 chars for identification (rv_xxxxxxxx)
    scopes TEXT[] NOT NULL DEFAULT '{}', -- ['read', 'write', 'admin']
    is_active BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_api_keys_user_id ON shared.api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON shared.api_keys(key_hash) WHERE is_active = true;
CREATE INDEX idx_api_keys_key_prefix ON shared.api_keys(key_prefix);
CREATE INDEX idx_api_keys_expires_at ON shared.api_keys(expires_at) WHERE expires_at IS NOT NULL AND is_active = true;

-- Comments
COMMENT ON TABLE shared.api_keys IS 'API keys for programmatic access with scope-based permissions';
COMMENT ON COLUMN shared.api_keys.key_hash IS 'SHA-256 hash of the API key (never store plaintext)';
COMMENT ON COLUMN shared.api_keys.key_prefix IS 'First 8 chars for key identification (rv_xxxxxxxx)';
COMMENT ON COLUMN shared.api_keys.scopes IS 'Permission scopes: read, write, admin';
COMMENT ON COLUMN shared.api_keys.is_active IS 'Key can be revoked by setting to false';
COMMENT ON COLUMN shared.api_keys.last_used_at IS 'Timestamp of last successful authentication';
