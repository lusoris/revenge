#!/bin/sh
set -e

# Auto-set GOMEMLIMIT to 90% of container memory limit if not already set.
# This tells Go's GC about the real memory boundary so it can reclaim heap
# before the OOM killer fires. Works with cgroup v2 and v1.
if [ -z "$GOMEMLIMIT" ]; then
  mem_bytes=""
  if [ -f /sys/fs/cgroup/memory.max ]; then
    # cgroup v2
    val=$(cat /sys/fs/cgroup/memory.max)
    if [ "$val" != "max" ]; then
      mem_bytes=$val
    fi
  elif [ -f /sys/fs/cgroup/memory/memory.limit_in_bytes ]; then
    # cgroup v1
    val=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)
    # Huge value means no limit (typically 2^63-1 page-aligned)
    if [ "$val" -lt 1099511627776 ] 2>/dev/null; then
      mem_bytes=$val
    fi
  fi

  if [ -n "$mem_bytes" ]; then
    # 90% of container limit, rounded to MiB
    limit_mib=$(( mem_bytes * 9 / 10 / 1048576 ))
    export GOMEMLIMIT="${limit_mib}MiB"
    echo "==> Auto-set GOMEMLIMIT=${GOMEMLIMIT} (90% of $(( mem_bytes / 1048576 ))MiB container limit)"
  fi
fi

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
