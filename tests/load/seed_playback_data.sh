#!/bin/bash
# tests/load/seed_playback_data.sh — Insert minimal movie + movie_file records
# so that playback_load.js has content IDs to create sessions against.
#
# Requires:
#   - PostgreSQL container running (revenge-postgres-dev or $DB_CONTAINER)
#   - The BBB test video mounted at /movies inside the app container
#
# Idempotent: uses ON CONFLICT DO NOTHING, safe to run repeatedly.

set -euo pipefail

DB_CONTAINER="${DB_CONTAINER:-revenge-postgres-dev}"
DB_USER="${DB_USER:-revenge}"
DB_NAME="${DB_NAME:-revenge}"
NUM_MOVIES="${NUM_MOVIES:-5}"

psql_exec() {
    docker exec "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -tAc "$1"
}

echo "Seeding playback test data (${NUM_MOVIES} movies)..."

# Verify DB connection
if ! docker exec "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" >/dev/null 2>&1; then
    echo "ERROR: Cannot connect to database container '$DB_CONTAINER'"
    exit 1
fi

# File metadata for the test_movie_XX.mp4 files (pre-generated in the NFS mount)
# test_movie_01: 30s 1080p h264+aac 40594638 bytes
# test_movie_02: 30s 1080p h264+aac 31704781 bytes
# test_movie_03: 30s 2160p h264+aac 86084619 bytes
# test_movie_04: 30s 2160p h264+aac 86096558 bytes
# test_movie_05: 30s 720p  h264+aac  9070777 bytes
FILE_SIZES=(0 40594638 31704781 86084619 86096558 9070777)
RESOLUTIONS=('' '1080p' '1080p' '2160p' '2160p' '720p')
BITRATES=(0 10800 8450 22950 22950 2420)

# Insert movies with deterministic UUIDs so we can reference them for movie_files
for i in $(seq 1 "$NUM_MOVIES"); do
    IDX=$(printf "%04d" "$i")
    MOVIE_UUID="00000000-0000-4000-a000-10ad0000${IDX}"
    FILE_UUID="00000000-0000-4000-a000-f11e0000${IDX}"
    TITLE="Load Test Movie ${i}"
    YEAR=$((2000 + i))

    # Movie 1 uses the real filename, others use numbered symlinks
    if [ "$i" -eq 1 ]; then
        FILE_SUFFIX=""
    else
        FILE_SUFFIX="_${i}"
    fi

    IDX2=$(printf "%02d" "$i")

    psql_exec "
        INSERT INTO movie.movies (id, title, year, runtime, overview, original_language)
        VALUES ('${MOVIE_UUID}', '${TITLE}', ${YEAR}, 1,
                'Synthetic movie for load testing playback sessions.', 'en')
        ON CONFLICT (id) DO NOTHING;
    "

    psql_exec "
        INSERT INTO movie.movie_files
            (id, movie_id, file_path, file_size, file_name, resolution, video_codec, audio_codec, container, duration_seconds, bitrate_kbps)
        VALUES ('${FILE_UUID}', '${MOVIE_UUID}',
                '/movies/test_movie_${IDX2}.mp4', ${FILE_SIZES[$i]},
                'test_movie_${IDX2}.mp4', '${RESOLUTIONS[$i]}', 'h264', 'aac', 'mp4', 30, ${BITRATES[$i]})
        ON CONFLICT (id) DO NOTHING;
    "
done

# Verify
COUNT=$(psql_exec "SELECT count(*) FROM movie.movies WHERE id::text LIKE '00000000-0000-4000-a000-10ad%'")
echo "Seeded ${COUNT} movies with movie_files for playback_load.js"
