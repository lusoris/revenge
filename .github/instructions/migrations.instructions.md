---
applyTo: "**/migrations/**/*.sql"
---

# SQL Migration Guidelines

> Standards for PostgreSQL migrations in revenge

## Naming Convention

```
NNNNNN_description.up.sql
NNNNNN_description.down.sql
```

Example: `000005_add_user_preferences.up.sql`

## Migration Structure

```sql
-- Always start with transaction
BEGIN;

-- Your changes here
CREATE TABLE ...
ALTER TABLE ...

-- End with commit
COMMIT;
```

## Down Migration Requirements

Every `up.sql` MUST have a corresponding `down.sql` that cleanly reverts changes.

```sql
-- 000005_add_user_preferences.down.sql
BEGIN;
DROP TABLE IF EXISTS user_preferences;
COMMIT;
```

## Type Conventions

- UUIDs: `UUID DEFAULT gen_random_uuid()`
- Timestamps: `TIMESTAMPTZ DEFAULT NOW()`
- Enums: Create separate type first
- Arrays: Use `TEXT[]` or specific type arrays

## Foreign Keys

```sql
-- Always include ON DELETE behavior
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
```

## Indexes

```sql
-- Create indexes for frequently queried columns
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_media_items_library_id ON media_items(library_id);
```
