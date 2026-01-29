-- Fleets (Libraries) - QAR obfuscated queries

-- name: GetFleetByID :one
SELECT * FROM qar.fleets WHERE id = $1;

-- name: ListFleets :many
SELECT * FROM qar.fleets
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListFleetsByOwner :many
SELECT * FROM qar.fleets
WHERE owner_user_id = $1
ORDER BY name;

-- name: ListFleetsByType :many
SELECT * FROM qar.fleets
WHERE fleet_type = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: CreateFleet :one
INSERT INTO qar.fleets (
    name,
    fleet_type,
    paths,
    stashdb_endpoint,
    tpdb_enabled,
    whisparr_sync,
    auto_tag_crew,
    fingerprint_on_scan,
    owner_user_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateFleet :one
UPDATE qar.fleets
SET
    name = $2,
    fleet_type = $3,
    paths = $4,
    stashdb_endpoint = $5,
    tpdb_enabled = $6,
    whisparr_sync = $7,
    auto_tag_crew = $8,
    fingerprint_on_scan = $9,
    owner_user_id = $10
WHERE id = $1
RETURNING *;

-- name: DeleteFleet :exec
DELETE FROM qar.fleets WHERE id = $1;

-- name: CountFleetExpeditions :one
SELECT COUNT(*) FROM qar.expeditions WHERE fleet_id = $1;

-- name: CountFleetVoyages :one
SELECT COUNT(*) FROM qar.voyages WHERE fleet_id = $1;

-- name: GetFleetStats :one
SELECT
    COUNT(DISTINCT e.id) AS expedition_count,
    COUNT(DISTINCT v.id) AS voyage_count,
    COALESCE(SUM(e.size_bytes), 0) + COALESCE(SUM(v.size_bytes), 0) AS total_size_bytes
FROM qar.fleets f
LEFT JOIN qar.expeditions e ON e.fleet_id = f.id
LEFT JOIN qar.voyages v ON v.fleet_id = f.id
WHERE f.id = $1;
