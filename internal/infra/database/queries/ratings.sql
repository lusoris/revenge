-- name: GetRatingSystemByID :one
SELECT id, code, name, country_codes, is_active, sort_order, created_at
FROM rating_systems
WHERE id = $1;

-- name: GetRatingSystemByCode :one
SELECT id, code, name, country_codes, is_active, sort_order, created_at
FROM rating_systems
WHERE code = $1;

-- name: ListRatingSystems :many
SELECT id, code, name, country_codes, is_active, sort_order, created_at
FROM rating_systems
WHERE is_active = true
ORDER BY sort_order, name;

-- name: ListRatingSystemsByCountry :many
SELECT id, code, name, country_codes, is_active, sort_order, created_at
FROM rating_systems
WHERE is_active = true
  AND $1 = ANY(country_codes)
ORDER BY sort_order, name;

-- name: GetRatingByID :one
SELECT r.id, r.system_id, r.code, r.name, r.description, r.min_age,
       r.normalized_level, r.sort_order, r.is_adult, r.icon_url, r.created_at,
       rs.code as system_code, rs.name as system_name
FROM ratings r
JOIN rating_systems rs ON r.system_id = rs.id
WHERE r.id = $1;

-- name: GetRatingBySystemAndCode :one
SELECT r.id, r.system_id, r.code, r.name, r.description, r.min_age,
       r.normalized_level, r.sort_order, r.is_adult, r.icon_url, r.created_at,
       rs.code as system_code, rs.name as system_name
FROM ratings r
JOIN rating_systems rs ON r.system_id = rs.id
WHERE r.system_id = $1 AND r.code = $2;

-- name: ListRatingsBySystem :many
SELECT r.id, r.system_id, r.code, r.name, r.description, r.min_age,
       r.normalized_level, r.sort_order, r.is_adult, r.icon_url, r.created_at,
       rs.code as system_code, rs.name as system_name
FROM ratings r
JOIN rating_systems rs ON r.system_id = rs.id
WHERE r.system_id = $1
ORDER BY r.sort_order, r.normalized_level;

-- name: ListRatingsByNormalizedLevel :many
SELECT r.id, r.system_id, r.code, r.name, r.description, r.min_age,
       r.normalized_level, r.sort_order, r.is_adult, r.icon_url, r.created_at,
       rs.code as system_code, rs.name as system_name
FROM ratings r
JOIN rating_systems rs ON r.system_id = rs.id
WHERE r.normalized_level <= $1
ORDER BY r.normalized_level, r.system_id, r.sort_order;

-- name: GetRatingEquivalents :many
SELECT r.id, r.system_id, r.code, r.name, r.description, r.min_age,
       r.normalized_level, r.sort_order, r.is_adult, r.icon_url, r.created_at,
       rs.code as system_code, rs.name as system_name
FROM ratings r
JOIN rating_systems rs ON r.system_id = rs.id
JOIN rating_equivalents re ON r.id = re.equivalent_rating_id
WHERE re.rating_id = $1
ORDER BY rs.sort_order, r.sort_order;

-- name: GetContentRatings :many
SELECT cr.id, cr.content_id, cr.content_type, cr.rating_id, cr.source, cr.created_at,
       r.code as rating_code, r.name as rating_name, r.normalized_level, r.is_adult,
       rs.code as system_code, rs.name as system_name
FROM content_ratings cr
JOIN ratings r ON cr.rating_id = r.id
JOIN rating_systems rs ON r.system_id = rs.id
WHERE cr.content_id = $1 AND cr.content_type = $2
ORDER BY rs.sort_order, r.sort_order;

-- name: GetContentMinLevel :one
SELECT content_id, content_type, min_level, is_adult
FROM content_min_rating_levels
WHERE content_id = $1 AND content_type = $2;

-- name: GetContentDisplayRating :one
-- Get the rating to display for content, preferring the specified system
SELECT cr.id, cr.content_id, cr.content_type, cr.rating_id, cr.source, cr.created_at,
       r.code as rating_code, r.name as rating_name, r.normalized_level, r.is_adult,
       rs.code as system_code, rs.name as system_name
FROM content_ratings cr
JOIN ratings r ON cr.rating_id = r.id
JOIN rating_systems rs ON r.system_id = rs.id
WHERE cr.content_id = $1 AND cr.content_type = $2
ORDER BY
    CASE WHEN rs.code = $3 THEN 0 ELSE 1 END,  -- Prefer specified system
    rs.sort_order,
    r.sort_order
LIMIT 1;

-- name: CreateContentRating :one
INSERT INTO content_ratings (content_id, content_type, rating_id, source)
VALUES ($1, $2, $3, $4)
ON CONFLICT (content_id, rating_id) DO UPDATE SET source = EXCLUDED.source
RETURNING id, content_id, content_type, rating_id, source, created_at;

-- name: DeleteContentRating :exec
DELETE FROM content_ratings
WHERE content_id = $1 AND rating_id = $2;

-- name: DeleteAllContentRatings :exec
DELETE FROM content_ratings
WHERE content_id = $1 AND content_type = $2;

-- name: IsContentAllowed :one
-- Check if content is allowed for a user's rating level
SELECT EXISTS (
    SELECT 1 FROM content_min_rating_levels
    WHERE content_id = $1
      AND content_type = $2
      AND min_level <= $3
      AND (NOT is_adult OR $4 = true)
) AS allowed;

-- name: FilterAllowedContentIDs :many
-- Filter a list of content IDs to only those allowed for a user
SELECT cmrl.content_id
FROM content_min_rating_levels cmrl
WHERE cmrl.content_id = ANY($1::uuid[])
  AND cmrl.content_type = $2
  AND cmrl.min_level <= $3
  AND (NOT cmrl.is_adult OR $4 = true);

-- name: RefreshContentMinRatingLevels :exec
REFRESH MATERIALIZED VIEW CONCURRENTLY content_min_rating_levels;
