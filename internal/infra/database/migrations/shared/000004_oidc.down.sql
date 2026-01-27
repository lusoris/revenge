-- 000004_oidc.down.sql
DROP TRIGGER IF EXISTS update_oidc_providers_updated_at ON oidc_providers;
DROP TABLE IF EXISTS oidc_user_links;
DROP TABLE IF EXISTS oidc_providers;
