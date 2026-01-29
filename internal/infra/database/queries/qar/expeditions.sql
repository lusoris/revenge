-- Expeditions (Adult Movies) - QAR obfuscated queries

-- name: GetExpeditionByID :one
SELECT * FROM qar.expeditions WHERE id = $1;

-- name: ListExpeditions :many
SELECT * FROM qar.expeditions
ORDER BY title
LIMIT $1 OFFSET $2;

-- name: ListExpeditionsByFleet :many
SELECT * FROM qar.expeditions
WHERE fleet_id = $1
ORDER BY title
LIMIT $2 OFFSET $3;

-- name: CreateExpedition :one
INSERT INTO qar.expeditions (
    fleet_id,
    title,
    sort_title,
    original_title,
    overview,
    launch_date,
    runtime_ticks,
    port_id,
    director,
    series,
    path,
    size_bytes,
    container,
    video_codec,
    audio_codec,
    resolution,
    coordinates,
    oshash,
    whisparr_id,
    charter,
    registry,
    has_file,
    is_hdr,
    is_3d
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15, $16,
    $17, $18, $19, $20, $21, $22, $23, $24
)
RETURNING *;

-- name: UpdateExpedition :one
UPDATE qar.expeditions
SET
    fleet_id = $2,
    title = $3,
    sort_title = $4,
    original_title = $5,
    overview = $6,
    launch_date = $7,
    runtime_ticks = $8,
    port_id = $9,
    director = $10,
    series = $11,
    path = $12,
    size_bytes = $13,
    container = $14,
    video_codec = $15,
    audio_codec = $16,
    resolution = $17,
    coordinates = $18,
    oshash = $19,
    whisparr_id = $20,
    charter = $21,
    registry = $22,
    has_file = $23,
    is_hdr = $24,
    is_3d = $25
WHERE id = $1
RETURNING *;

-- name: DeleteExpedition :exec
DELETE FROM qar.expeditions WHERE id = $1;

-- name: GetExpeditionByPath :one
SELECT * FROM qar.expeditions WHERE path = $1;

-- name: GetExpeditionByCoordinates :one
SELECT * FROM qar.expeditions WHERE coordinates = $1;

-- name: GetExpeditionByCharter :one
SELECT * FROM qar.expeditions WHERE charter = $1;

-- name: CountExpeditionsByFleet :one
SELECT COUNT(*) FROM qar.expeditions WHERE fleet_id = $1;

-- name: SearchExpeditions :many
SELECT * FROM qar.expeditions
WHERE title ILIKE '%' || $1 || '%'
   OR original_title ILIKE '%' || $1 || '%'
ORDER BY title
LIMIT $2 OFFSET $3;
