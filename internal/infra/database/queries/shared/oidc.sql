-- OIDC Providers
-- name: GetOIDCProviderByID :one
SELECT * FROM oidc_providers WHERE id = $1;

-- name: GetOIDCProviderBySlug :one
SELECT * FROM oidc_providers WHERE slug = $1;

-- name: ListOIDCProviders :many
SELECT * FROM oidc_providers ORDER BY name ASC;

-- name: ListEnabledOIDCProviders :many
SELECT * FROM oidc_providers WHERE enabled = true ORDER BY name ASC;

-- name: CreateOIDCProvider :one
INSERT INTO oidc_providers (
    name, slug, enabled,
    issuer_url, client_id, client_secret_enc, scopes,
    claim_mapping, role_mapping,
    auto_provision, default_role
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: UpdateOIDCProvider :one
UPDATE oidc_providers SET
    name = COALESCE(sqlc.narg('name'), name),
    enabled = COALESCE(sqlc.narg('enabled'), enabled),
    issuer_url = COALESCE(sqlc.narg('issuer_url'), issuer_url),
    client_id = COALESCE(sqlc.narg('client_id'), client_id),
    client_secret_enc = COALESCE(sqlc.narg('client_secret_enc'), client_secret_enc),
    scopes = COALESCE(sqlc.narg('scopes'), scopes),
    claim_mapping = COALESCE(sqlc.narg('claim_mapping'), claim_mapping),
    role_mapping = COALESCE(sqlc.narg('role_mapping'), role_mapping),
    auto_provision = COALESCE(sqlc.narg('auto_provision'), auto_provision),
    default_role = COALESCE(sqlc.narg('default_role'), default_role)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteOIDCProvider :exec
DELETE FROM oidc_providers WHERE id = $1;

-- OIDC User Links
-- name: GetOIDCLinkByID :one
SELECT * FROM oidc_user_links WHERE id = $1;

-- name: GetOIDCLinkBySubject :one
SELECT * FROM oidc_user_links
WHERE provider_id = $1 AND subject = $2;

-- name: ListOIDCLinksByUser :many
SELECT l.*, p.name as provider_name, p.slug as provider_slug
FROM oidc_user_links l
JOIN oidc_providers p ON l.provider_id = p.id
WHERE l.user_id = $1;

-- name: CreateOIDCLink :one
INSERT INTO oidc_user_links (
    user_id, provider_id, subject, email, name, groups
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateOIDCLinkLogin :exec
UPDATE oidc_user_links SET
    last_login_at = NOW(),
    email = COALESCE(sqlc.narg('email'), email),
    name = COALESCE(sqlc.narg('name'), name),
    groups = COALESCE(sqlc.narg('groups'), groups)
WHERE id = sqlc.arg('id');

-- name: DeleteOIDCLink :exec
DELETE FROM oidc_user_links WHERE id = $1;

-- name: DeleteOIDCLinksByUser :exec
DELETE FROM oidc_user_links WHERE user_id = $1;

-- name: DeleteOIDCLinksByProvider :exec
DELETE FROM oidc_user_links WHERE provider_id = $1;
