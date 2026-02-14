#!/bin/bash
# tests/load/run.sh - Convenience script to run k6 load tests
# Usage: ./tests/load/run.sh [test] [profile]
#   test: all_endpoints | realistic_usage | auth_stress (default: all_endpoints)
#   profile: smoke | gentle | spike | soak | stress (default: smoke)

set -e

TEST="${1:-all_endpoints}"
PROFILE="${2:-smoke}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Validate test
case "$TEST" in
    all_endpoints|realistic_usage|auth_stress|playback_load|api_key_load|write_operations)
        TEST_FILE="${SCRIPT_DIR}/${TEST}.js"
        ;;
    *)
        echo "Unknown test: $TEST"
        echo "Available: all_endpoints, realistic_usage, auth_stress, playback_load, api_key_load, write_operations"
        exit 1
        ;;
esac

# Validate profile
case "$PROFILE" in
    smoke|gentle|spike|soak|stress)
        ;;
    *)
        echo "Unknown profile: $PROFILE"
        echo "Available: smoke, gentle, spike, soak, stress"
        exit 1
        ;;
esac

echo "========================================"
echo "k6 Load Test"
echo "========================================"
echo "Test:     $TEST"
echo "Profile:  $PROFILE"
echo "Target:   ${BASE_URL:-http://localhost:8096}"
echo "========================================"
echo ""

# Auto-seed playback test data if running playback_load
if [[ "$TEST" == "playback_load" ]]; then
    echo "Seeding playback test data..."
    "${SCRIPT_DIR}/seed_playback_data.sh"
    echo ""
fi

# Run k6
k6 run \
    --env PROFILE="$PROFILE" \
    --env BASE_URL="${BASE_URL:-http://localhost:8096}" \
    --env TEST_USER="${TEST_USER:-dbg_ext1}" \
    --env TEST_PASS="${TEST_PASS:-TestPass123!}" \
    --out json="results/${TEST}_${PROFILE}_$(date +%Y%m%d_%H%M%S).json" \
    "$TEST_FILE"
