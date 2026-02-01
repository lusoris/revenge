-- Create users table in shared schema
-- This is the foundation for authentication and authorization

CREATE TABLE shared.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Authentication
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- bcrypt or argon2

    -- Profile
    username VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    avatar_url TEXT,

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_admin BOOLEAN NOT NULL DEFAULT false,
    email_verified BOOLEAN NOT NULL DEFAULT false,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,

    -- Soft delete
    deleted_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_users_email ON shared.users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username ON shared.users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_is_active ON shared.users(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON shared.users(created_at);

-- Updated timestamp trigger
CREATE OR REPLACE FUNCTION shared.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON shared.users
    FOR EACH ROW
    EXECUTE FUNCTION shared.update_updated_at_column();

-- Comments
COMMENT ON TABLE shared.users IS 'User accounts for authentication and authorization';
COMMENT ON COLUMN shared.users.id IS 'Unique user identifier (UUID v4)';
COMMENT ON COLUMN shared.users.email IS 'User email address (unique, used for login)';
COMMENT ON COLUMN shared.users.password_hash IS 'Hashed password (bcrypt or argon2)';
COMMENT ON COLUMN shared.users.username IS 'Unique username (alphanumeric, used for @mentions)';
COMMENT ON COLUMN shared.users.display_name IS 'Display name shown in UI';
COMMENT ON COLUMN shared.users.is_active IS 'Whether the user account is active (can log in)';
COMMENT ON COLUMN shared.users.is_admin IS 'Whether the user has admin privileges';
COMMENT ON COLUMN shared.users.email_verified IS 'Whether the email address has been verified';
COMMENT ON COLUMN shared.users.deleted_at IS 'Soft delete timestamp (NULL = not deleted)';
