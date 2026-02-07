-- Create schemas for Revenge
-- public: Main content (movies, TV shows, music, etc.)
-- shared: Shared services (users, sessions, settings, etc.)
-- qar: Adult content (isolated with access control)

-- public schema already exists by default in PostgreSQL
-- Ensure it's configured correctly
COMMENT ON SCHEMA public IS 'Main content: movies, TV shows, music, audiobooks, books, podcasts';

-- Create shared schema for shared services
CREATE SCHEMA IF NOT EXISTS shared;
COMMENT ON SCHEMA shared IS 'Shared services: users, sessions, settings, RBAC, activity';

-- Create qar schema for adult content (requires legacy:read scope)
CREATE SCHEMA IF NOT EXISTS qar;
COMMENT ON SCHEMA qar IS 'QAR (Adult content): voyages, expeditions, treasures - requires legacy:read scope';

-- Enable required PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm;   -- Trigram text search (used for fuzzy title matching)
CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- Cryptographic functions (used for gen_random_uuid fallback)

-- Shared utility function: auto-update updated_at timestamp on row modification
-- Used by movie tables (000021-000026) and any future shared-schema tables
CREATE OR REPLACE FUNCTION shared.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Set search path to include all schemas
-- Application will set this per-connection based on user permissions
-- Default search path: public, shared (qar requires explicit scope)
-- Note: We don't use ALTER DATABASE here because it requires knowing the DB name
-- Instead, the application will SET search_path per connection
-- or via pgxpool config: "search_path=public,shared"
