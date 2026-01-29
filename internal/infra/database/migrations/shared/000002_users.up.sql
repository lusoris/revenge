-- Users table: Account-level authentication
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username        VARCHAR(100) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE,
    password_hash   VARCHAR(255),                    -- NULL if OIDC-only user
    is_admin        BOOLEAN NOT NULL DEFAULT false,
    is_disabled     BOOLEAN NOT NULL DEFAULT false,

    -- Content access
    max_rating_level    INT NOT NULL DEFAULT 100,    -- 0-100 normalized rating level
    adult_enabled       BOOLEAN NOT NULL DEFAULT false,  -- Access to schema c

    -- Preferences
    preferred_language  VARCHAR(10) DEFAULT 'en',    -- ISO 639-1
    preferred_rating_system VARCHAR(20) DEFAULT 'mpaa',

    -- Timestamps
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
