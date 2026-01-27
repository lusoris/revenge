-- 000002_users.down.sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS users;
-- Note: update_updated_at_column() function is kept for other tables
