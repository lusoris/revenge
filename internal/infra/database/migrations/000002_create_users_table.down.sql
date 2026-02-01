-- Drop users table and related objects

DROP TRIGGER IF EXISTS update_users_updated_at ON shared.users;
DROP FUNCTION IF EXISTS shared.update_updated_at_column();
DROP TABLE IF EXISTS shared.users CASCADE;
