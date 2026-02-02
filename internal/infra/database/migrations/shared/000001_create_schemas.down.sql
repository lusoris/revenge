-- Drop schemas (reverse migration)
-- Note: This will CASCADE delete all tables in these schemas!

DROP SCHEMA IF EXISTS qar CASCADE;
DROP SCHEMA IF EXISTS shared CASCADE;
-- Don't drop public schema as it's a PostgreSQL default
