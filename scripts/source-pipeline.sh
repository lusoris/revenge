#!/usr/bin/env bash
#
# Source Pipeline Runner
#
# Runs the source documentation pipeline in the correct order:
# 1. Fetch external sources
# 2. Generate source index
# 3. Add source breadcrumbs to design docs
#
# Usage:
#   ./scripts/source-pipeline.sh              # Dry run
#   ./scripts/source-pipeline.sh --apply      # Actually run
#   ./scripts/source-pipeline.sh --step 1     # Run only step 1
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
VENV_DIR="$PROJECT_ROOT/.venv"

# Activate venv if it exists
if [[ -d "$VENV_DIR" ]]; then
    # shellcheck source=/dev/null
    source "$VENV_DIR/bin/activate"
elif [[ ! -f "$SCRIPT_DIR/requirements.txt" ]]; then
    echo "Error: No venv found. Create one with:"
    echo "  python3 -m venv .venv && source .venv/bin/activate && pip install -r scripts/requirements.txt"
    exit 1
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default options
APPLY=""
STEP=""
FORCE=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --apply)
            APPLY="--apply"
            shift
            ;;
        --force)
            FORCE="--force"
            shift
            ;;
        --step)
            STEP="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --apply    Actually run (default: dry-run)"
            echo "  --force    Force update unchanged sources"
            echo "  --step N   Run only step N (1-3)"
            echo "  -h, --help Show this help"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Check for uncommitted changes
if [[ -n "$APPLY" ]]; then
    if [[ -n "$(git -C "$PROJECT_ROOT" status --porcelain)" ]]; then
        echo -e "${YELLOW}Warning: You have uncommitted changes.${NC}"
        echo "Consider committing or stashing before running with --apply."
        read -p "Continue anyway? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
fi

echo -e "${BLUE}╔══════════════════════════════════════╗${NC}"
echo -e "${BLUE}║      Source Pipeline Runner          ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════╝${NC}"
echo ""

if [[ -z "$APPLY" ]]; then
    echo -e "${YELLOW}=== DRY RUN MODE ===${NC}"
    echo ""
fi

run_step() {
    local step_num=$1
    local step_name=$2
    local script=$3
    shift 3
    local args=("$@")

    if [[ -n "$STEP" && "$STEP" != "$step_num" ]]; then
        echo -e "${YELLOW}Skipping step $step_num: $step_name${NC}"
        return 0
    fi

    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}Step $step_num: $step_name${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    if python3 "$SCRIPT_DIR/$script" "${args[@]}"; then
        echo ""
        echo -e "${GREEN}✓ Step $step_num complete${NC}"
    else
        echo ""
        echo -e "${RED}✗ Step $step_num failed${NC}"
        exit 1
    fi
    echo ""
}

# Step 1: Fetch external sources
FETCH_ARGS=()
[[ -n "$APPLY" ]] && FETCH_ARGS+=("--apply")
[[ -n "$FORCE" ]] && FETCH_ARGS+=("--force")
run_step 1 "Fetch External Sources" "source-pipeline/01-fetch.py" "${FETCH_ARGS[@]}"

# Step 2: Generate source index
INDEX_ARGS=()
[[ -n "$APPLY" ]] && INDEX_ARGS+=("--apply")
run_step 2 "Generate Source Index" "source-pipeline/02-index.py" "${INDEX_ARGS[@]}"

# Step 3: Add source breadcrumbs
BREADCRUMB_ARGS=()
[[ -n "$APPLY" ]] && BREADCRUMB_ARGS+=("--apply")
run_step 3 "Add Source Breadcrumbs" "source-pipeline/03-breadcrumbs.py" "${BREADCRUMB_ARGS[@]}"

echo -e "${BLUE}╔══════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Source Pipeline Complete         ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════╝${NC}"

if [[ -z "$APPLY" ]]; then
    echo ""
    echo -e "${YELLOW}This was a dry run. Use --apply to actually run.${NC}"
fi
