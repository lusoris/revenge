-- 000001_extensions.up.sql
-- PostgreSQL extensions required by Jellyfin Go
-- PostgreSQL 18+

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
