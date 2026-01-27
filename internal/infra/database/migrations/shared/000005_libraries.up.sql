-- 000005_libraries.up.sql
-- Libraries - collections of media content

CREATE TYPE library_type AS ENUM ('movies', 'tvshows', 'music', 'photos', 'mixed');

CREATE TABLE libraries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type library_type NOT NULL,
    paths TEXT[] NOT NULL DEFAULT '{}',       -- Array of filesystem paths
    settings JSONB NOT NULL DEFAULT '{}',     -- Library-specific settings
    is_visible BOOLEAN NOT NULL DEFAULT true,
    scan_interval_hours INT DEFAULT 24,
    last_scan_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_libraries_type ON libraries(type);
CREATE INDEX idx_libraries_visible ON libraries(is_visible) WHERE is_visible = true;

CREATE TRIGGER update_libraries_updated_at
    BEFORE UPDATE ON libraries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
