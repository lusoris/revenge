-- name: GetNetwork :one
SELECT * FROM tvshow.networks WHERE id = $1;

-- name: GetNetworkByTMDbID :one
SELECT * FROM tvshow.networks WHERE tmdb_id = $1;

-- name: ListNetworksBySeries :many
SELECT n.* FROM tvshow.networks n
JOIN tvshow.series_networks sn ON n.id = sn.network_id
WHERE sn.series_id = $1
ORDER BY n.name ASC;

-- name: CreateNetwork :one
INSERT INTO tvshow.networks (tmdb_id, name, logo_path, origin_country)
VALUES ($1, $2, $3, $4)
ON CONFLICT (tmdb_id) DO UPDATE SET
    name = EXCLUDED.name,
    logo_path = EXCLUDED.logo_path,
    origin_country = EXCLUDED.origin_country
RETURNING *;

-- name: AddSeriesNetwork :exec
INSERT INTO tvshow.series_networks (series_id, network_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteSeriesNetworks :exec
DELETE FROM tvshow.series_networks WHERE series_id = $1;

-- name: ListSeriesByNetwork :many
SELECT s.* FROM tvshow.series s
JOIN tvshow.series_networks sn ON s.id = sn.series_id
WHERE sn.network_id = $1
ORDER BY s.first_air_date DESC NULLS LAST
LIMIT $2 OFFSET $3;
