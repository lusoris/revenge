-- Rollback schema creation
-- Drop schemas in reverse order (qar, shared, keep public)

-- Drop shared utility function (must be before dropping the schema)
DROP FUNCTION IF EXISTS shared.update_updated_at_column() CASCADE;

-- Drop qar schema and all its objects
DROP SCHEMA IF EXISTS qar CASCADE;

-- Drop shared schema and all its objects
DROP SCHEMA IF EXISTS shared CASCADE;

-- Drop extensions
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS pgcrypto;

-- Reset search path to default
ALTER DATABASE revenge RESET search_path;

-- Note: We don't drop the public schema as it's a PostgreSQL default
