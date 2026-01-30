-- Resource Grants: Polymorphic per-resource access control
-- Grants access to specific resources (libraries, playlists, collections) to users
BEGIN;

-- Polymorphic resource grants table
-- Grant knows what resource type it's for - no central registry needed
CREATE TABLE resource_grants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Polymorphic reference (grant owns the reference)
    resource_type   VARCHAR(50) NOT NULL,   -- 'movie_library', 'tv_library', 'playlist', 'collection', etc.
    resource_id     UUID NOT NULL,          -- UUID of the actual resource

    -- Grant level: view, edit, manage, owner
    grant_type      VARCHAR(20) NOT NULL DEFAULT 'view',

    -- Audit trail
    granted_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    granted_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ,            -- Optional expiration

    -- Each user can only have one grant per resource
    UNIQUE (user_id, resource_type, resource_id),

    -- Validate grant_type
    CONSTRAINT valid_grant_type CHECK (grant_type IN ('view', 'edit', 'manage', 'owner'))
);

-- Index for looking up user's grants
CREATE INDEX idx_resource_grants_user ON resource_grants(user_id);

-- Index for looking up grants by resource
CREATE INDEX idx_resource_grants_resource ON resource_grants(resource_type, resource_id);

-- Index for expiration cleanup (only where expiration is set)
CREATE INDEX idx_resource_grants_expires ON resource_grants(expires_at)
    WHERE expires_at IS NOT NULL;

-- Grant types reference:
-- view:   Can view/browse the resource
-- edit:   Can view + edit the resource
-- manage: Can view + edit + delete/add items
-- owner:  Full control including sharing and deletion

COMMENT ON TABLE resource_grants IS 'Polymorphic per-resource access grants for sharing libraries, playlists, etc.';
COMMENT ON COLUMN resource_grants.resource_type IS 'Type of resource: movie_library, tv_library, music_library, adult_library, playlist, collection';
COMMENT ON COLUMN resource_grants.grant_type IS 'Access level: view, edit, manage, owner';

COMMIT;
