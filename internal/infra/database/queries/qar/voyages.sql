-- Voyages (Adult Scenes) - QAR obfuscated queries

-- name: GetVoyageByID :one
SELECT * FROM qar.voyages WHERE id = $1;

-- name: ListVoyages :many
SELECT * FROM qar.voyages
ORDER BY title
LIMIT $1 OFFSET $2;

-- name: ListVoyagesByFleet :many
SELECT * FROM qar.voyages
WHERE fleet_id = $1
ORDER BY title
LIMIT $2 OFFSET $3;

-- name: CreateVoyage :one
INSERT INTO qar.voyages (
    fleet_id,
    title,
    sort_title,
    overview,
    launch_date,
    distance,
    port_id,
    whisparr_id,
    stash_id,
    charter,
    registry,
    path,
    size_bytes,
    video_codec,
    audio_codec,
    resolution,
    oshash,
    coordinates,
    md5,
    cover_path
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING *;

-- name: UpdateVoyage :one
UPDATE qar.voyages
SET
    fleet_id = $2,
    title = $3,
    sort_title = $4,
    overview = $5,
    launch_date = $6,
    distance = $7,
    port_id = $8,
    whisparr_id = $9,
    stash_id = $10,
    charter = $11,
    registry = $12,
    path = $13,
    size_bytes = $14,
    video_codec = $15,
    audio_codec = $16,
    resolution = $17,
    oshash = $18,
    coordinates = $19,
    md5 = $20,
    cover_path = $21
WHERE id = $1
RETURNING *;

-- name: DeleteVoyage :exec
DELETE FROM qar.voyages WHERE id = $1;

-- name: GetVoyageByPath :one
SELECT * FROM qar.voyages WHERE path = $1;

-- name: GetVoyageByOshash :one
SELECT * FROM qar.voyages WHERE oshash = $1;

-- name: GetVoyageByCoordinates :one
SELECT * FROM qar.voyages WHERE coordinates = $1;

-- name: GetVoyageByCharter :one
SELECT * FROM qar.voyages WHERE charter = $1;

-- name: CountVoyagesByFleet :one
SELECT COUNT(*) FROM qar.voyages WHERE fleet_id = $1;

-- name: SearchVoyages :many
SELECT * FROM qar.voyages
WHERE title ILIKE '%' || $1 || '%'
ORDER BY title
LIMIT $2 OFFSET $3;

-- name: ListVoyagesByPort :many
SELECT * FROM qar.voyages
WHERE port_id = $1
ORDER BY launch_date DESC
LIMIT $2 OFFSET $3;
