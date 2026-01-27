-- 000005_libraries.down.sql
DROP TRIGGER IF EXISTS update_libraries_updated_at ON libraries;
DROP TABLE IF EXISTS libraries;
DROP TYPE IF EXISTS library_type;
