-- Drop triggers
DROP TRIGGER IF EXISTS trigger_oidc_user_links_updated_at ON shared.oidc_user_links;
DROP TRIGGER IF EXISTS trigger_oidc_providers_updated_at ON shared.oidc_providers;
DROP FUNCTION IF EXISTS shared.update_oidc_updated_at();

-- Drop tables in correct order (respecting foreign keys)
DROP TABLE IF EXISTS shared.oidc_states;
DROP TABLE IF EXISTS shared.oidc_user_links;
DROP TABLE IF EXISTS shared.oidc_providers;
