-- Rollback schema creation
-- Drop schemas in reverse order (qar, shared, keep public)

-- Drop qar schema and all its objects
DROP SCHEMA IF EXISTS qar CASCADE;

-- Drop shared schema and all its objects
DROP SCHEMA IF EXISTS shared CASCADE;

-- Reset search path to default
ALTER DATABASE revenge RESET search_path;

-- Note: We don't drop the public schema as it's a PostgreSQL default
