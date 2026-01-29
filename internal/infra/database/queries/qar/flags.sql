-- Flags (Tags) - QAR obfuscated queries

-- name: GetFlagByID :one
SELECT * FROM qar.flags WHERE id = $1;

-- name: GetFlagByName :one
SELECT * FROM qar.flags WHERE name = $1;

-- name: ListFlags :many
SELECT * FROM qar.flags
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CreateFlag :one
INSERT INTO qar.flags (name, description, parent_id, stashdb_id, waters)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateFlag :one
UPDATE qar.flags
SET
    name = $2,
    description = $3,
    parent_id = $4,
    stashdb_id = $5,
    waters = $6
WHERE id = $1
RETURNING *;

-- name: DeleteFlag :exec
DELETE FROM qar.flags WHERE id = $1;

-- name: GetFlagByStashDBID :one
SELECT * FROM qar.flags WHERE stashdb_id = $1;

-- name: SearchFlags :many
SELECT * FROM qar.flags
WHERE name ILIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: ListFlagsByWaters :many
SELECT * FROM qar.flags
WHERE waters = $1
ORDER BY name;

-- name: ListFlagChildren :many
SELECT * FROM qar.flags
WHERE parent_id = $1
ORDER BY name;

-- name: ListRootFlags :many
SELECT * FROM qar.flags
WHERE parent_id IS NULL
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListExpeditionFlags :many
SELECT f.* FROM qar.flags f
JOIN qar.expedition_flags ef ON ef.flag_id = f.id
WHERE ef.expedition_id = $1
ORDER BY f.name;

-- name: ListVoyageFlags :many
SELECT f.* FROM qar.flags f
JOIN qar.voyage_flags vf ON vf.flag_id = f.id
WHERE vf.voyage_id = $1
ORDER BY f.name;

-- name: AddExpeditionFlag :exec
INSERT INTO qar.expedition_flags (expedition_id, flag_id)
VALUES ($1, $2)
ON CONFLICT (expedition_id, flag_id) DO NOTHING;

-- name: AddVoyageFlag :exec
INSERT INTO qar.voyage_flags (voyage_id, flag_id)
VALUES ($1, $2)
ON CONFLICT (voyage_id, flag_id) DO NOTHING;

-- name: RemoveExpeditionFlag :exec
DELETE FROM qar.expedition_flags WHERE expedition_id = $1 AND flag_id = $2;

-- name: RemoveVoyageFlag :exec
DELETE FROM qar.voyage_flags WHERE voyage_id = $1 AND flag_id = $2;

-- name: ClearExpeditionFlags :exec
DELETE FROM qar.expedition_flags WHERE expedition_id = $1;

-- name: ClearVoyageFlags :exec
DELETE FROM qar.voyage_flags WHERE voyage_id = $1;
