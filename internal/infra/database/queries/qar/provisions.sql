-- Provisions (Content Requests) - QAR obfuscated queries
-- Request system for adult content

-- ============================================================================
-- PROVISIONS (Requests)
-- ============================================================================

-- name: GetProvisionByID :one
SELECT * FROM qar.provisions WHERE id = $1;

-- name: ListProvisions :many
SELECT * FROM qar.provisions
ORDER BY
    CASE WHEN @sort_by::text = 'priority' THEN priority END DESC,
    CASE WHEN @sort_by::text = 'ayes_count' THEN ayes_count END DESC,
    CASE WHEN @sort_by::text = 'created_at' OR @sort_by::text = '' THEN created_at END DESC
LIMIT $1 OFFSET $2;

-- name: ListProvisionsByUser :many
SELECT * FROM qar.provisions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListProvisionsByStatus :many
SELECT * FROM qar.provisions
WHERE status = $1
ORDER BY priority DESC, ayes_count DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListProvisionsByContentType :many
SELECT * FROM qar.provisions
WHERE content_type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPendingProvisions :many
SELECT * FROM qar.provisions
WHERE status = 'pending'
ORDER BY priority DESC, ayes_count DESC, created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchProvisions :many
SELECT * FROM qar.provisions
WHERE title ILIKE '%' || $1 || '%'
ORDER BY priority DESC, ayes_count DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateProvision :one
INSERT INTO qar.provisions (
    user_id,
    content_type,
    request_subtype,
    external_id,
    external_source,
    title,
    release_year,
    manifest,
    estimated_cargo_gb
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateProvisionStatus :one
UPDATE qar.provisions
SET
    status = $2,
    approved_by_user_id = $3,
    approved_at = $4,
    declined_reason = $5
WHERE id = $1
RETURNING *;

-- name: UpdateProvisionPriority :one
UPDATE qar.provisions
SET priority = $2
WHERE id = $1
RETURNING *;

-- name: UpdateProvisionIntegration :one
UPDATE qar.provisions
SET
    integration_id = $2,
    integration_status = $3
WHERE id = $1
RETURNING *;

-- name: UpdateProvisionAvailable :one
UPDATE qar.provisions
SET
    status = 'available',
    actual_cargo_gb = $2,
    available_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetProvisionAutoApproved :one
UPDATE qar.provisions
SET
    auto_approved = true,
    auto_article_id = $2,
    status = 'approved',
    approved_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProvision :exec
DELETE FROM qar.provisions WHERE id = $1;

-- name: CountProvisionsByUser :one
SELECT COUNT(*) FROM qar.provisions WHERE user_id = $1;

-- name: CountProvisionsByStatus :one
SELECT COUNT(*) FROM qar.provisions WHERE status = $1;

-- name: GetProvisionByExternalID :one
SELECT * FROM qar.provisions
WHERE external_source = $1 AND external_id = $2;

-- ============================================================================
-- PROVISION AYES (Votes)
-- ============================================================================

-- name: GetProvisionAye :one
SELECT * FROM qar.provision_ayes
WHERE provision_id = $1 AND user_id = $2;

-- name: ListProvisionAyes :many
SELECT * FROM qar.provision_ayes
WHERE provision_id = $1
ORDER BY voted_at DESC;

-- name: CreateProvisionAye :one
INSERT INTO qar.provision_ayes (provision_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteProvisionAye :exec
DELETE FROM qar.provision_ayes
WHERE provision_id = $1 AND user_id = $2;

-- name: HasUserVoted :one
SELECT EXISTS(
    SELECT 1 FROM qar.provision_ayes
    WHERE provision_id = $1 AND user_id = $2
);

-- name: CountProvisionAyes :one
SELECT COUNT(*) FROM qar.provision_ayes WHERE provision_id = $1;

-- ============================================================================
-- PROVISION MISSIVES (Comments)
-- ============================================================================

-- name: GetProvisionMissive :one
SELECT * FROM qar.provision_missives WHERE id = $1;

-- name: ListProvisionMissives :many
SELECT * FROM qar.provision_missives
WHERE provision_id = $1
ORDER BY created_at ASC;

-- name: CreateProvisionMissive :one
INSERT INTO qar.provision_missives (
    provision_id,
    user_id,
    message,
    is_captain_order
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteProvisionMissive :exec
DELETE FROM qar.provision_missives WHERE id = $1;

-- name: CountProvisionMissives :one
SELECT COUNT(*) FROM qar.provision_missives WHERE provision_id = $1;

-- ============================================================================
-- RATIONS (User Quotas)
-- ============================================================================

-- name: GetRation :one
SELECT * FROM qar.rations WHERE user_id = $1;

-- name: CreateRation :one
INSERT INTO qar.rations (user_id)
VALUES ($1)
RETURNING *;

-- name: UpsertRation :one
INSERT INTO qar.rations (user_id)
VALUES ($1)
ON CONFLICT (user_id) DO UPDATE SET user_id = EXCLUDED.user_id
RETURNING *;

-- name: UpdateRationLimits :one
UPDATE qar.rations
SET
    daily_limit = COALESCE($2, daily_limit),
    weekly_limit = COALESCE($3, weekly_limit),
    monthly_limit = COALESCE($4, monthly_limit),
    cargo_quota_gb = COALESCE($5, cargo_quota_gb)
WHERE user_id = $1
RETURNING *;

-- name: IncrementRationUsage :one
UPDATE qar.rations
SET
    daily_used = daily_used + 1,
    weekly_used = weekly_used + 1,
    monthly_used = monthly_used + 1
WHERE user_id = $1
RETURNING *;

-- name: AddRationCargoUsage :one
UPDATE qar.rations
SET cargo_used_gb = cargo_used_gb + $2
WHERE user_id = $1
RETURNING *;

-- name: ResetDailyRations :exec
UPDATE qar.rations
SET
    daily_used = 0,
    last_reset_daily = CURRENT_DATE
WHERE last_reset_daily < CURRENT_DATE;

-- name: ResetWeeklyRations :exec
UPDATE qar.rations
SET
    weekly_used = 0,
    last_reset_weekly = CURRENT_DATE
WHERE last_reset_weekly < (CURRENT_DATE - INTERVAL '7 days');

-- name: ResetMonthlyRations :exec
UPDATE qar.rations
SET
    monthly_used = 0,
    last_reset_monthly = CURRENT_DATE
WHERE last_reset_monthly < (CURRENT_DATE - INTERVAL '30 days');

-- ============================================================================
-- ARTICLES (Rules)
-- ============================================================================

-- name: GetArticle :one
SELECT * FROM qar.articles WHERE id = $1;

-- name: ListArticles :many
SELECT * FROM qar.articles
ORDER BY priority DESC, created_at ASC;

-- name: ListEnabledArticles :many
SELECT * FROM qar.articles
WHERE enabled = true
ORDER BY priority DESC;

-- name: ListArticlesByContentType :many
SELECT * FROM qar.articles
WHERE content_type = $1 OR content_type IS NULL
ORDER BY priority DESC;

-- name: CreateArticle :one
INSERT INTO qar.articles (
    name,
    description,
    content_type,
    condition_type,
    condition_value,
    action,
    automation_trigger,
    enabled,
    priority
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateArticle :one
UPDATE qar.articles
SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    content_type = $4,
    condition_type = COALESCE($5, condition_type),
    condition_value = COALESCE($6, condition_value),
    action = COALESCE($7, action),
    automation_trigger = $8,
    enabled = COALESCE($9, enabled),
    priority = COALESCE($10, priority)
WHERE id = $1
RETURNING *;

-- name: DeleteArticle :exec
DELETE FROM qar.articles WHERE id = $1;

-- name: SetArticleEnabled :one
UPDATE qar.articles
SET enabled = $2
WHERE id = $1
RETURNING *;

-- ============================================================================
-- CARGO HOLD (Global Storage Quotas)
-- ============================================================================

-- name: GetCargoHold :one
SELECT * FROM qar.cargo_hold WHERE id = 1;

-- name: UpdateCargoHoldUsage :one
UPDATE qar.cargo_hold
SET
    total_used_gb = $1,
    expedition_used_gb = $2,
    voyage_used_gb = $3
WHERE id = 1
RETURNING *;

-- name: UpdateCargoHoldQuotas :one
UPDATE qar.cargo_hold
SET
    total_quota_gb = $1,
    expedition_quota_gb = $2,
    voyage_quota_gb = $3
WHERE id = 1
RETURNING *;

-- name: AddCargoHoldUsage :one
UPDATE qar.cargo_hold
SET
    total_used_gb = total_used_gb + $1,
    expedition_used_gb = CASE WHEN $2 THEN expedition_used_gb + $1 ELSE expedition_used_gb END,
    voyage_used_gb = CASE WHEN NOT $2 THEN voyage_used_gb + $1 ELSE voyage_used_gb END
WHERE id = 1
RETURNING *;
