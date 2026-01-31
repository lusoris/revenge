#!/usr/bin/env bash
#
# Doc Pipeline Runner
#
# Runs the documentation pipeline in the correct order:
# 0. Regenerate docs from YAML data (Claude + Wiki)
# 1. Generate INDEX.md files
# 2. Add design breadcrumbs
# 3. Sync status tables
# 4. Validate document structure
# 5. Fix broken links
# 6. Generate meta files (DESIGN_INDEX.md)
#
# Usage:
#   ./scripts/doc-pipeline.sh              # Dry run (preview mode)
#   ./scripts/doc-pipeline.sh --apply      # Actually run (writes files)
#   ./scripts/doc-pipeline.sh --step 0     # Run only step 0 (regenerate)
#   ./scripts/doc-pipeline.sh --step 1     # Run only step 1
#   ./scripts/doc-pipeline.sh --validate   # Only run validation (no changes)
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
VALIDATE_ONLY=""
ADD_MISSING=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --apply)
            APPLY="--apply"
            shift
            ;;
        --step)
            STEP="$2"
            shift 2
            ;;
        --validate)
            VALIDATE_ONLY="1"
            shift
            ;;
        --add-missing)
            ADD_MISSING="--add-missing"
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --apply        Actually run (default: dry-run)"
            echo "  --step N       Run only step N (1-6)"
            echo "  --validate     Only run validation step"
            echo "  --add-missing  Add status tables to docs without one"
            echo "  -h, --help     Show this help"
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
echo -e "${BLUE}║        Doc Pipeline Runner           ║${NC}"
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
        # Don't exit on validation failures
        if [[ "$step_num" != "4" ]]; then
            exit 1
        fi
    fi
    echo ""
}

# Validation only mode
if [[ -n "$VALIDATE_ONLY" ]]; then
    echo -e "${BLUE}Running validation only...${NC}"
    echo ""
    python3 "$SCRIPT_DIR/doc-pipeline/04-validate.py"
    exit $?
fi

# Step 0: Regenerate docs from YAML data
run_step 0 "Regenerate Docs from YAML" "automation/batch_regenerate.py"

# Step 1: Generate INDEX.md files
INDEX_ARGS=()
[[ -n "$APPLY" ]] && INDEX_ARGS+=("--apply")
run_step 1 "Generate INDEX.md Files" "doc-pipeline/01-indexes.py" "${INDEX_ARGS[@]}"

# Step 2: Add design breadcrumbs
BREADCRUMB_ARGS=()
[[ -n "$APPLY" ]] && BREADCRUMB_ARGS+=("--apply")
run_step 2 "Add Design Breadcrumbs" "doc-pipeline/02-breadcrumbs.py" "${BREADCRUMB_ARGS[@]}"

# Step 3: Sync status tables
STATUS_ARGS=()
[[ -n "$APPLY" ]] && STATUS_ARGS+=("--apply")
[[ -n "$ADD_MISSING" ]] && STATUS_ARGS+=("--add-missing")
run_step 3 "Sync Status Tables" "doc-pipeline/03-status.py" "${STATUS_ARGS[@]}"

# Step 4: Validate document structure
run_step 4 "Validate Document Structure" "doc-pipeline/04-validate.py"

# Step 5: Fix broken links
FIX_ARGS=()
[[ -n "$APPLY" ]] && FIX_ARGS+=("--apply")
run_step 5 "Fix Broken Links" "doc-pipeline/05-fix.py" "${FIX_ARGS[@]}"

# Step 6: Generate meta files
META_ARGS=()
[[ -n "$APPLY" ]] && META_ARGS+=("--apply")
run_step 6 "Generate Meta Files" "doc-pipeline/06-meta.py" "${META_ARGS[@]}"

echo -e "${BLUE}╔══════════════════════════════════════╗${NC}"
echo -e "${BLUE}║       Doc Pipeline Complete          ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════╝${NC}"

if [[ -z "$APPLY" ]]; then
    echo ""
    echo -e "${YELLOW}This was a dry run. Use --apply to actually run.${NC}"
fi
