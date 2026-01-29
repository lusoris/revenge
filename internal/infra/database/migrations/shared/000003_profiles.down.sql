DROP TRIGGER IF EXISTS users_create_default_profile ON users;
DROP FUNCTION IF EXISTS create_default_profile();
DROP TRIGGER IF EXISTS profiles_updated_at ON profiles;
DROP TABLE IF EXISTS profiles;
