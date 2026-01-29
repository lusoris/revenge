-- Crew (Performers) - QAR obfuscated queries

-- name: GetCrewByID :one
SELECT * FROM qar.crew WHERE id = $1;

-- name: ListCrew :many
SELECT * FROM qar.crew
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CreateCrew :one
INSERT INTO qar.crew (
    name,
    disambiguation,
    gender,
    christening,
    death_date,
    birth_city,
    origin,
    nationality,
    rigging,
    compass,
    height_cm,
    weight_kg,
    measurements,
    cup_size,
    breast_type,
    markings,
    anchors,
    maiden_voyage,
    last_port,
    bio,
    stash_id,
    charter,
    registry,
    manifest,
    twitter,
    instagram,
    image_path
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27
)
RETURNING *;

-- name: UpdateCrew :one
UPDATE qar.crew
SET
    name = $2,
    disambiguation = $3,
    gender = $4,
    christening = $5,
    death_date = $6,
    birth_city = $7,
    origin = $8,
    nationality = $9,
    rigging = $10,
    compass = $11,
    height_cm = $12,
    weight_kg = $13,
    measurements = $14,
    cup_size = $15,
    breast_type = $16,
    markings = $17,
    anchors = $18,
    maiden_voyage = $19,
    last_port = $20,
    bio = $21,
    stash_id = $22,
    charter = $23,
    registry = $24,
    manifest = $25,
    twitter = $26,
    instagram = $27,
    image_path = $28
WHERE id = $1
RETURNING *;

-- name: DeleteCrew :exec
DELETE FROM qar.crew WHERE id = $1;

-- name: GetCrewByCharter :one
SELECT * FROM qar.crew WHERE charter = $1;

-- name: GetCrewByRegistry :one
SELECT * FROM qar.crew WHERE registry = $1;

-- name: SearchCrew :many
SELECT * FROM qar.crew
WHERE name ILIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: ListCrewNames :many
SELECT * FROM qar.crew_names WHERE crew_id = $1;

-- name: AddCrewName :exec
INSERT INTO qar.crew_names (crew_id, name)
VALUES ($1, $2)
ON CONFLICT (crew_id, name) DO NOTHING;

-- name: RemoveCrewName :exec
DELETE FROM qar.crew_names WHERE crew_id = $1 AND name = $2;

-- name: ListCrewPortraits :many
SELECT * FROM qar.crew_portraits WHERE crew_id = $1;

-- name: AddCrewPortrait :one
INSERT INTO qar.crew_portraits (crew_id, path, type, source, primary_image)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: SetPrimaryPortrait :exec
UPDATE qar.crew_portraits
SET primary_image = (id = $2)
WHERE crew_id = $1;

-- name: ListExpeditionCrew :many
SELECT c.* FROM qar.crew c
JOIN qar.expedition_crew ec ON ec.crew_id = c.id
WHERE ec.expedition_id = $1
ORDER BY c.name;

-- name: ListVoyageCrew :many
SELECT c.* FROM qar.crew c
JOIN qar.voyage_crew vc ON vc.crew_id = c.id
WHERE vc.voyage_id = $1
ORDER BY c.name;

-- name: AddExpeditionCrew :exec
INSERT INTO qar.expedition_crew (expedition_id, crew_id, character_name)
VALUES ($1, $2, $3)
ON CONFLICT (expedition_id, crew_id) DO NOTHING;

-- name: AddVoyageCrew :exec
INSERT INTO qar.voyage_crew (voyage_id, crew_id, role)
VALUES ($1, $2, $3)
ON CONFLICT (voyage_id, crew_id) DO NOTHING;

-- name: RemoveExpeditionCrew :exec
DELETE FROM qar.expedition_crew WHERE expedition_id = $1 AND crew_id = $2;

-- name: RemoveVoyageCrew :exec
DELETE FROM qar.voyage_crew WHERE voyage_id = $1 AND crew_id = $2;
