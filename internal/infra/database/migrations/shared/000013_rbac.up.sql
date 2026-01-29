-- RBAC: Add role column to users table
-- Roles: admin, moderator, user, guest (as defined in ARCHITECTURE_V2.md)

-- Create role enum type
CREATE TYPE user_role AS ENUM ('admin', 'moderator', 'user', 'guest');

-- Add role column to users table
ALTER TABLE users ADD COLUMN role user_role NOT NULL DEFAULT 'user';

-- Migrate existing is_admin values to role
UPDATE users SET role = 'admin' WHERE is_admin = true;

-- Create index for role queries
CREATE INDEX idx_users_role ON users(role);

-- Note: is_admin column is kept for backwards compatibility during migration period
-- It can be removed in a future migration once all code uses the role column
