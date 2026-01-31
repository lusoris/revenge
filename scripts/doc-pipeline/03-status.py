#!/usr/bin/env python3
"""Sync status tables in design documents.

Ensures consistent status table format across all design docs:
- Standardizes Dimension/Status/Notes columns
- Validates status values (✅, 🟡, 🔴, ⚪)
- Can add status tables to docs without one

Usage:
    python scripts/doc-pipeline/03-status.py                  # Check all
    python scripts/doc-pipeline/03-status.py --apply          # Fix formatting
    python scripts/doc-pipeline/03-status.py --add-missing --apply
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOT_FILE = DESIGN_DIR / "00_SOURCE_OF_TRUTH.md"

VALID_STATUSES = {"✅", "🟡", "🔴", "⚪", "-"}

STANDARD_DIMENSIONS = [
    "Design",
    "Sources",
    "Instructions",
    "Code",
    "Linting",
    "Unit Testing",
    "Integration Testing",
]

SKIP_FILES = {
    "00_SOURCE_OF_TRUTH.md",
    "01_DESIGN_DOC_TEMPLATE.md",
    "02_QUESTIONS_TO_DISCUSS.md",
    "03_DESIGN_DOCS_STATUS.md",
    "DESIGN_INDEX.md",
    "NAVIGATION.md",
}

STATUS_TABLE_PATTERN = re.compile(
    r"\|\s*Dimension\s*\|\s*Status\s*\|.*?\n\|[-\s|]+\n((?:\|.*?\n)+)",
    re.MULTILINE,
)

ROW_PATTERN = re.compile(r"\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|(?:\s*([^|]*)\s*\|)?")


def parse_status_table(content: str) -> list[dict] | None:
    """Parse a status table from content."""
    match = STATUS_TABLE_PATTERN.search(content)
    if not match:
        return None

    rows = []
    for line in match.group(1).strip().split("\n"):
        row_match = ROW_PATTERN.match(line)
        if row_match:
            dimension = row_match.group(1).strip()
            status = row_match.group(2).strip()
            notes = row_match.group(3).strip() if row_match.group(3) else ""

            rows.append(
                {
                    "dimension": dimension,
                    "status": status,
                    "notes": notes,
                },
            )

    return rows


def validate_status_table(rows: list[dict]) -> list[str]:
    """Validate a status table. Returns list of issues."""
    issues = []

    for row in rows:
        status = row["status"]
        status_char = status[0] if status else ""
        if status_char not in VALID_STATUSES and status not in VALID_STATUSES:
            issues.append(f"Invalid status '{status}' for {row['dimension']}")

    return issues


def generate_status_table(rows: list[dict], *, include_notes: bool = True) -> str:
    """Generate a properly formatted status table."""
    if include_notes:
        lines = [
            "| Dimension | Status | Notes |",
            "|-----------|--------|-------|",
        ]
        for row in rows:
            lines.append(f"| {row['dimension']} | {row['status']} | {row['notes']} |")
    else:
        lines = [
            "| Dimension | Status |",
            "|-----------|--------|",
        ]
        for row in rows:
            lines.append(f"| {row['dimension']} | {row['status']} |")

    return "\n".join(lines)


def fix_status_table(content: str) -> tuple[str, bool]:
    """Fix status table formatting. Returns (new_content, changed)."""
    match = STATUS_TABLE_PATTERN.search(content)
    if not match:
        return content, False

    rows = parse_status_table(content)
    if not rows:
        return content, False

    has_notes = any(row["notes"] for row in rows)
    new_table = generate_status_table(rows, include_notes=has_notes)

    full_pattern = re.compile(
        r"\|\s*Dimension\s*\|\s*Status\s*\|.*?\n\|[-\s|]+\n(?:\|.*?\n)+",
        re.MULTILINE,
    )

    full_match = full_pattern.search(content)
    if not full_match:
        return content, False

    new_content = (
        content[: full_match.start()] + new_table + content[full_match.end() :]
    )

    return new_content, new_content != content


def find_design_docs(category: str | None = None) -> list[Path]:
    """Find design documents."""
    docs = []
    for md_file in sorted(DESIGN_DIR.rglob("*.md")):
        if ".archive" in str(md_file):
            continue
        if ".analysis" in str(md_file):
            continue
        if md_file.name in SKIP_FILES:
            continue
        if md_file.name.startswith("INDEX") or md_file.name.startswith("00_INDEX"):
            continue

        if category:
            rel_path = md_file.relative_to(DESIGN_DIR)
            if not str(rel_path).startswith(category):
                continue

        docs.append(md_file)
    return docs


def ensure_status_table(content: str) -> tuple[str, bool]:
    """Ensure document has a status table. Returns (new_content, changed)."""
    if parse_status_table(content):
        return content, False

    default_rows = [
        {"dimension": "Design", "status": "🔴", "notes": ""},
        {"dimension": "Sources", "status": "🔴", "notes": ""},
        {"dimension": "Instructions", "status": "🔴", "notes": ""},
        {"dimension": "Code", "status": "🔴", "notes": ""},
        {"dimension": "Linting", "status": "🔴", "notes": ""},
        {"dimension": "Unit Testing", "status": "🔴", "notes": ""},
        {"dimension": "Integration Testing", "status": "🔴", "notes": ""},
    ]

    new_table = (
        "## Status\n\n"
        + generate_status_table(default_rows, include_notes=True)
        + "\n\n---"
    )

    lines = content.split("\n")
    insert_pos = 0

    for i, line in enumerate(lines):
        if line.startswith("# "):
            insert_pos = i + 1
            continue
        if insert_pos > 0 and not line.strip():
            insert_pos = i + 1
            continue
        if insert_pos > 0 and line.startswith(">"):
            insert_pos = i + 1
            continue
        if insert_pos > 0 and line.startswith("<!-- "):
            insert_pos = i + 1
            continue
        if insert_pos > 0 and not line.strip():
            insert_pos = i + 1
            break
        if insert_pos > 0:
            break

    new_lines = [*lines[:insert_pos], "", new_table, "", *lines[insert_pos:]]
    return "\n".join(new_lines), True


def main() -> int:
    parser = argparse.ArgumentParser(description="Sync status tables")
    parser.add_argument(
        "--category",
        "-c",
        help="Only check specific category",
    )
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Apply fixes (default: dry-run)",
    )
    parser.add_argument(
        "--add-missing",
        action="store_true",
        help="Add status tables to docs without one",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Show details",
    )
    args = parser.parse_args()

    dry_run = not args.apply

    if dry_run:
        print("=== DRY RUN (use --apply to write) ===\n")

    docs = find_design_docs(args.category)
    print(f"Checking {len(docs)} documents...")

    stats = {
        "has_table": 0,
        "no_table": 0,
        "issues": 0,
        "fixed": 0,
    }

    for doc_path in docs:
        content = doc_path.read_text(encoding="utf-8")
        rel_path = doc_path.relative_to(DESIGN_DIR)

        rows = parse_status_table(content)

        if not rows:
            stats["no_table"] += 1
            if args.verbose:
                print(f"  No status table: {rel_path}")

            if args.add_missing:
                new_content, changed = ensure_status_table(content)
                if changed:
                    if not dry_run:
                        doc_path.write_text(new_content, encoding="utf-8")
                        print(f"  Added status table: {rel_path}")
                    else:
                        print(f"  Would add status table: {rel_path}")
                    stats["fixed"] += 1
            continue

        stats["has_table"] += 1

        issues = validate_status_table(rows)
        if issues:
            stats["issues"] += len(issues)
            print(f"\n{rel_path}:")
            for issue in issues:
                print(f"  ⚠️  {issue}")

        new_content, changed = fix_status_table(content)
        if changed:
            if not dry_run:
                doc_path.write_text(new_content, encoding="utf-8")
                if args.verbose:
                    print(f"  Fixed formatting: {rel_path}")
            elif args.verbose:
                print(f"  Would fix formatting: {rel_path}")
            stats["fixed"] += 1

    print(f"\n{'=' * 50}")
    print("STATUS TABLE SUMMARY")
    print(f"{'=' * 50}")
    print(f"Documents checked: {len(docs)}")
    print(f"With status table: {stats['has_table']}")
    print(f"Without status table: {stats['no_table']}")
    print(f"Validation issues: {stats['issues']}")
    action = "Fixed" if not dry_run else "Would fix"
    print(f"{action}: {stats['fixed']}")

    if len(docs) > 0:
        coverage = 100 * stats["has_table"] // len(docs)
        print(f"\nStatus table coverage: {coverage}%")

    if stats["no_table"] > 0 and not args.add_missing:
        print("\nUse --add-missing --apply to add status tables to docs without one")

    if dry_run and stats["fixed"] > 0:
        print("\n=== DRY RUN complete. Use --apply to write. ===")

    return 0 if stats["issues"] == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
