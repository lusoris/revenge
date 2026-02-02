-- name: CreateOIDCProvider :one
-- Creates a new OIDC provider configuration
INSERT INTO shared.oidc_providers (
    name,
    display_name,
    provider_type,
    issuer_url,
    client_id,
    client_secret_encrypted,
    authorization_endpoint,
    token_endpoint,
    userinfo_endpoint,
    jwks_uri,
    end_session_endpoint,
    scopes,
    claim_mappings,
    role_mappings,
    auto_create_users,
    update_user_info,
    allow_linking,
    is_enabled,
    is_default
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
) RETURNING *;

-- name: GetOIDCProvider :one
-- Gets an OIDC provider by ID
SELECT * FROM shared.oidc_providers WHERE id = $1;

-- name: GetOIDCProviderByName :one
-- Gets an OIDC provider by name
SELECT * FROM shared.oidc_providers WHERE name = $1;

-- name: GetDefaultOIDCProvider :one
-- Gets the default OIDC provider
SELECT * FROM shared.oidc_providers WHERE is_default = true AND is_enabled = true;

-- name: ListOIDCProviders :many
-- Lists all OIDC providers
SELECT * FROM shared.oidc_providers ORDER BY display_name;

-- name: ListEnabledOIDCProviders :many
-- Lists all enabled OIDC providers
SELECT * FROM shared.oidc_providers WHERE is_enabled = true ORDER BY display_name;

-- name: UpdateOIDCProvider :one
-- Updates an OIDC provider
UPDATE shared.oidc_providers SET
    display_name = COALESCE(sqlc.narg('display_name'), display_name),
    provider_type = COALESCE(sqlc.narg('provider_type'), provider_type),
    issuer_url = COALESCE(sqlc.narg('issuer_url'), issuer_url),
    client_id = COALESCE(sqlc.narg('client_id'), client_id),
    client_secret_encrypted = COALESCE(sqlc.narg('client_secret_encrypted'), client_secret_encrypted),
    authorization_endpoint = COALESCE(sqlc.narg('authorization_endpoint'), authorization_endpoint),
    token_endpoint = COALESCE(sqlc.narg('token_endpoint'), token_endpoint),
    userinfo_endpoint = COALESCE(sqlc.narg('userinfo_endpoint'), userinfo_endpoint),
    jwks_uri = COALESCE(sqlc.narg('jwks_uri'), jwks_uri),
    end_session_endpoint = COALESCE(sqlc.narg('end_session_endpoint'), end_session_endpoint),
    scopes = COALESCE(sqlc.narg('scopes'), scopes),
    claim_mappings = COALESCE(sqlc.narg('claim_mappings'), claim_mappings),
    role_mappings = COALESCE(sqlc.narg('role_mappings'), role_mappings),
    auto_create_users = COALESCE(sqlc.narg('auto_create_users'), auto_create_users),
    update_user_info = COALESCE(sqlc.narg('update_user_info'), update_user_info),
    allow_linking = COALESCE(sqlc.narg('allow_linking'), allow_linking),
    is_enabled = COALESCE(sqlc.narg('is_enabled'), is_enabled),
    is_default = COALESCE(sqlc.narg('is_default'), is_default)
WHERE id = $1
RETURNING *;

-- name: SetDefaultOIDCProvider :exec
-- Sets a provider as default (clears other defaults first)
UPDATE shared.oidc_providers SET is_default = false WHERE is_default = true;

-- name: DeleteOIDCProvider :exec
-- Deletes an OIDC provider
DELETE FROM shared.oidc_providers WHERE id = $1;

-- name: EnableOIDCProvider :exec
-- Enables an OIDC provider
UPDATE shared.oidc_providers SET is_enabled = true WHERE id = $1;

-- name: DisableOIDCProvider :exec
-- Disables an OIDC provider
UPDATE shared.oidc_providers SET is_enabled = false WHERE id = $1;

-- ============================================================================
-- OIDC User Links
-- ============================================================================

-- name: CreateOIDCUserLink :one
-- Links a user to an OIDC provider
INSERT INTO shared.oidc_user_links (
    user_id,
    provider_id,
    subject,
    email,
    name,
    picture_url,
    access_token_encrypted,
    refresh_token_encrypted,
    token_expires_at,
    last_login_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
) RETURNING *;

-- name: GetOIDCUserLink :one
-- Gets a user link by ID
SELECT * FROM shared.oidc_user_links WHERE id = $1;

-- name: GetOIDCUserLinkBySubject :one
-- Gets a user link by provider and subject
SELECT * FROM shared.oidc_user_links
WHERE provider_id = $1 AND subject = $2;

-- name: GetOIDCUserLinkByUserAndProvider :one
-- Gets a user link by user and provider
SELECT * FROM shared.oidc_user_links
WHERE user_id = $1 AND provider_id = $2;

-- name: ListUserOIDCLinks :many
-- Lists all OIDC links for a user
SELECT l.*, p.name as provider_name, p.display_name as provider_display_name
FROM shared.oidc_user_links l
JOIN shared.oidc_providers p ON l.provider_id = p.id
WHERE l.user_id = $1
ORDER BY l.created_at;

-- name: UpdateOIDCUserLink :one
-- Updates a user link (tokens, user info)
UPDATE shared.oidc_user_links SET
    email = COALESCE(sqlc.narg('email'), email),
    name = COALESCE(sqlc.narg('name'), name),
    picture_url = COALESCE(sqlc.narg('picture_url'), picture_url),
    access_token_encrypted = COALESCE(sqlc.narg('access_token_encrypted'), access_token_encrypted),
    refresh_token_encrypted = COALESCE(sqlc.narg('refresh_token_encrypted'), refresh_token_encrypted),
    token_expires_at = COALESCE(sqlc.narg('token_expires_at'), token_expires_at),
    last_login_at = COALESCE(sqlc.narg('last_login_at'), last_login_at)
WHERE id = $1
RETURNING *;

-- name: UpdateOIDCUserLinkLastLogin :exec
-- Updates the last login timestamp
UPDATE shared.oidc_user_links SET last_login_at = NOW() WHERE id = $1;

-- name: DeleteOIDCUserLink :exec
-- Unlinks a user from an OIDC provider
DELETE FROM shared.oidc_user_links WHERE id = $1;

-- name: DeleteOIDCUserLinkByUserAndProvider :exec
-- Unlinks a user from an OIDC provider by user and provider ID
DELETE FROM shared.oidc_user_links WHERE user_id = $1 AND provider_id = $2;

-- name: CountUserOIDCLinks :one
-- Counts how many OIDC providers a user is linked to
SELECT COUNT(*) FROM shared.oidc_user_links WHERE user_id = $1;

-- ============================================================================
-- OIDC States (OAuth2 flow)
-- ============================================================================

-- name: CreateOIDCState :one
-- Creates a new OAuth2 state for the auth flow
INSERT INTO shared.oidc_states (
    state,
    code_verifier,
    provider_id,
    user_id,
    redirect_url,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOIDCState :one
-- Gets an OAuth2 state by state token
SELECT * FROM shared.oidc_states WHERE state = $1;

-- name: DeleteOIDCState :exec
-- Deletes an OAuth2 state (after use or expiration)
DELETE FROM shared.oidc_states WHERE state = $1;

-- name: DeleteExpiredOIDCStates :execrows
-- Cleans up expired OAuth2 states
DELETE FROM shared.oidc_states WHERE expires_at < NOW();

-- name: DeleteOIDCStatesByProvider :exec
-- Deletes all states for a provider (when disabling)
DELETE FROM shared.oidc_states WHERE provider_id = $1;
