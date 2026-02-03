#!/bin/sh
set -e

echo "==> Waiting for database to be ready..."
timeout=60
counter=0
until pg_isready -h postgres -U revenge > /dev/null 2>&1; do
  counter=$((counter + 1))
  if [ $counter -gt $timeout ]; then
    echo "ERROR: Database failed to become ready in ${timeout} seconds"
    exit 1
  fi
  echo "Waiting for postgres... ($counter/$timeout)"
  sleep 1
done

echo "==> Database is ready, running migrations..."
/app/revenge migrate up

echo "==> Starting revenge server..."
exec /app/revenge "$@"
