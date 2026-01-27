-- OIDC provider queries for Jellyfin Go

-- =============================================================================
-- PROVIDERS
-- =============================================================================

-- name: GetOIDCProviderByID :one
SELECT * FROM oidc_providers WHERE id = $1 LIMIT 1;

-- name: GetOIDCProviderByName :one
SELECT * FROM oidc_providers WHERE name = $1 LIMIT 1;

-- name: ListOIDCProviders :many
SELECT * FROM oidc_providers ORDER BY display_name ASC;

-- name: ListEnabledOIDCProviders :many
SELECT * FROM oidc_providers
WHERE enabled = true
ORDER BY display_name ASC;

-- name: CreateOIDCProvider :one
INSERT INTO oidc_providers (
    name, display_name, issuer_url, client_id, client_secret_encrypted,
    scopes, enabled, auto_create_users, default_admin, claim_mappings
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateOIDCProvider :one
UPDATE oidc_providers
SET
    display_name = COALESCE($2, display_name),
    issuer_url = COALESCE($3, issuer_url),
    client_id = COALESCE($4, client_id),
    client_secret_encrypted = COALESCE($5, client_secret_encrypted),
    scopes = COALESCE($6, scopes),
    enabled = COALESCE($7, enabled),
    auto_create_users = COALESCE($8, auto_create_users),
    default_admin = COALESCE($9, default_admin),
    claim_mappings = COALESCE($10, claim_mappings),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOIDCProviderEnabled :exec
UPDATE oidc_providers
SET enabled = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteOIDCProvider :exec
DELETE FROM oidc_providers WHERE id = $1;

-- =============================================================================
-- USER LINKS
-- =============================================================================

-- name: GetOIDCUserLink :one
SELECT * FROM oidc_user_links
WHERE provider_id = $1 AND subject = $2
LIMIT 1;

-- name: GetOIDCUserLinkByUserID :many
SELECT
    l.*,
    p.name AS provider_name,
    p.display_name AS provider_display_name
FROM oidc_user_links l
JOIN oidc_providers p ON l.provider_id = p.id
WHERE l.user_id = $1;

-- name: CreateOIDCUserLink :one
INSERT INTO oidc_user_links (
    user_id, provider_id, subject, email
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateOIDCUserLinkLastLogin :exec
UPDATE oidc_user_links
SET last_login_at = NOW()
WHERE id = $1;

-- name: DeleteOIDCUserLink :exec
DELETE FROM oidc_user_links WHERE id = $1;

-- name: DeleteOIDCUserLinksByUser :exec
DELETE FROM oidc_user_links WHERE user_id = $1;

-- name: DeleteOIDCUserLinksByProvider :exec
DELETE FROM oidc_user_links WHERE provider_id = $1;

-- name: OIDCUserLinkExists :one
SELECT EXISTS(
    SELECT 1 FROM oidc_user_links
    WHERE provider_id = $1 AND subject = $2
);

-- =============================================================================
-- JOIN QUERIES
-- =============================================================================

-- name: GetUserByOIDCLink :one
SELECT u.*
FROM users u
JOIN oidc_user_links l ON u.id = l.user_id
WHERE l.provider_id = $1 AND l.subject = $2
LIMIT 1;
