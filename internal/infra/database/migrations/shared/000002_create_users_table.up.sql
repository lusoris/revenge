-- Create users table in shared schema
-- Basic user authentication and profile information

CREATE TABLE IF NOT EXISTS shared.users (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Authentication
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL, -- Argon2id hash

    -- Profile
    display_name VARCHAR(255),
    avatar_url TEXT,

    -- Settings
    locale VARCHAR(10) DEFAULT 'en-US',
    timezone VARCHAR(50) DEFAULT 'UTC',

    -- QAR Access (legacy:read scope)
    qar_enabled BOOLEAN DEFAULT FALSE,

    -- Account status
    is_active BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,

    -- Soft delete
    deleted_at TIMESTAMPTZ
);

-- Indexes for common queries
CREATE INDEX idx_users_username ON shared.users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email ON shared.users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_active ON shared.users(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON shared.users(created_at DESC);

-- Comments
COMMENT ON TABLE shared.users IS 'User accounts with authentication and profile information';
COMMENT ON COLUMN shared.users.qar_enabled IS 'Grants legacy:read scope for QAR (adult content) access';
COMMENT ON COLUMN shared.users.password_hash IS 'Argon2id password hash';
