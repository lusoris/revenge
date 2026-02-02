-- Create casbin_rule table for Casbin pgx-adapter
CREATE TABLE IF NOT EXISTS shared.casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100) NOT NULL,
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);

COMMENT ON TABLE shared.casbin_rule IS 'Casbin RBAC policy rules and role mappings';
COMMENT ON COLUMN shared.casbin_rule.ptype IS 'Policy type: p (policy) or g (role)';
COMMENT ON COLUMN shared.casbin_rule.v0 IS 'Subject (user/role)';
COMMENT ON COLUMN shared.casbin_rule.v1 IS 'Object (resource)';
COMMENT ON COLUMN shared.casbin_rule.v2 IS 'Action (read/write/delete)';

-- Indexes for faster policy lookups
CREATE INDEX idx_casbin_rule_ptype ON shared.casbin_rule (ptype);
CREATE INDEX idx_casbin_rule_v0 ON shared.casbin_rule (v0);
CREATE INDEX idx_casbin_rule_v1 ON shared.casbin_rule (v1);
CREATE INDEX idx_casbin_rule_v0_v1 ON shared.casbin_rule (v0, v1);

-- Insert default admin role and policies
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    -- Admin has all permissions on all resources
    ('p', 'admin', '*', '*'),
    -- User role - basic permissions
    ('p', 'user', 'profile', 'read'),
    ('p', 'user', 'profile', 'write'),
    ('p', 'user', 'library', 'read'),
    ('p', 'user', 'playback', 'read'),
    ('p', 'user', 'playback', 'write'),
    -- Guest role - read-only
    ('p', 'guest', 'library', 'read');

-- First user (ID from migration 000002) gets admin role
-- Note: This assumes the first user UUID, should be updated with actual first user ID
COMMENT ON TABLE shared.casbin_rule IS 'Default policies created. First user should be assigned admin role manually.';
