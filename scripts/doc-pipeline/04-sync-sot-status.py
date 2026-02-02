#!/usr/bin/env python3
"""Sync status from YAML files to SOURCE_OF_TRUTH.md tables.

Reads overall_status from YAML files and updates ONLY status cells in SOT tables:
- Content Modules table
- Backend Services table
- Metadata Providers table
- Arr Ecosystem table

Usage:
    python scripts/doc-pipeline/04-sync-sot-status.py                # Dry-run
    python scripts/doc-pipeline/04-sync-sot-status.py --apply        # Apply changes
    python scripts/doc-pipeline/04-sync-sot-status.py --strict       # Exit 1 on drift
"""

from __future__ import annotations

import argparse
import re
import sys
from datetime import datetime
from pathlib import Path
from typing import Any

import yaml

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
DATA_DIR = PROJECT_ROOT / "data"
SOT_FILE = DESIGN_DIR / "00_SOURCE_OF_TRUTH.md"

VALID_STATUSES = {"✅", "🟡", "🔴", "⚪", "🔵"}

# Mapping of YAML doc_title to SOT table entry name
# This allows us to match YAML files to their SOT table rows
YAML_TO_SOT_MAP = {
    # Content Modules (by doc_title)
    "Movie Module": "Movie",
    "TV Show Module": "TV Show",
    "Music Module": "Music",
    "Audiobook Module": "Audiobook",
    "Book Module": "Book",
    "Podcast Module": "Podcast",
    "Photo Library": "Photo",
    "Comics Module": "Comics",
    "Live TV & DVR": "LiveTV",
    "Adult Content System": "QAR Voyages",  # Maps to first row
    "Gallery Module": "QAR Treasures",

    # Backend Services (by service name)
    "Authentication Service": "Auth",
    "User Service": "User",
    "Session Service": "Session",
    "RBAC Service": "RBAC",
    "Activity Service": "Activity",
    "Settings Service": "Settings",
    "API Keys Service": "API Keys",
    "OIDC Service": "OIDC",
    "Grants Service": "Grants",
    "Fingerprint Service": "Fingerprint",
    "Library Service": "Library",
    "Playback Service": "Playback",
    "Metadata Service": "Metadata",
    "Search Service": "Search",
    "Health Service": "Health",
    "Scrobbling Service": "Scrobbling",

    # Metadata Providers
    "TMDb Provider": "TMDb",
    "TheTVDB Provider": "TheTVDB",
    "MusicBrainz Provider": "MusicBrainz",
    "Last.fm Provider": "Last.fm",
    "Audnexus Provider": "Audnexus",
    "OpenLibrary Provider": "OpenLibrary",
    "ComicVine Provider": "ComicVine",
    "StashDB Provider": "StashDB",
    "ThePornDB Provider": "ThePornDB",

    # Arr Ecosystem
    "Radarr Integration": "Radarr",
    "Sonarr Integration": "Sonarr",
    "Lidarr Integration": "Lidarr",
    "Whisparr Integration": "Whisparr",
    "Chaptarr Integration": "Chaptarr",
    "Prowlarr Integration": "Prowlarr",
}


def parse_yaml_status(yaml_path: Path) -> tuple[str, str | None]:
    """Parse overall_status from YAML file.

    Returns:
        (doc_title, overall_status) or (doc_title, None) if not found
    """
    try:
        with open(yaml_path, encoding="utf-8") as f:
            data = yaml.safe_load(f)

        doc_title = data.get("doc_title", "")
        overall_status = data.get("overall_status", None)

        return doc_title, overall_status
    except Exception as e:
        print(f"⚠️  Error reading {yaml_path}: {e}", file=sys.stderr)
        return "", None


def find_all_yaml_statuses() -> dict[str, str]:
    """Find all YAML files and extract their statuses.

    Returns:
        Dict mapping SOT table names to status emoji
    """
    sot_statuses = {}

    for yaml_path in DATA_DIR.rglob("*.yaml"):
        if ".archive" in str(yaml_path):
            continue

        doc_title, status = parse_yaml_status(yaml_path)

        if not doc_title or not status:
            continue

        # Map doc_title to SOT name
        sot_name = YAML_TO_SOT_MAP.get(doc_title)
        if sot_name:
            sot_statuses[sot_name] = status

    return sot_statuses


def parse_markdown_table(content: str, table_start: int) -> tuple[list[dict], int, int]:
    """Parse a markdown table starting at given position.

    Returns:
        (rows, table_start_pos, table_end_pos)
        Each row is dict with: {line_num, content, columns, status_col_idx}
    """
    lines = content.split("\n")
    start_line = content[:table_start].count("\n")

    rows = []
    current_line = start_line
    header_line = None
    separator_line = None

    # Find header
    while current_line < len(lines):
        line = lines[current_line]
        if line.startswith("|") and "Status" in line:
            header_line = current_line
            header_cols = [col.strip() for col in line.split("|")[1:-1]]
            try:
                status_col_idx = header_cols.index("Status")
            except ValueError:
                return [], table_start, table_start
            break
        current_line += 1

    if header_line is None:
        return [], table_start, table_start

    # Find separator
    current_line = header_line + 1
    if current_line < len(lines) and "---" in lines[current_line]:
        separator_line = current_line
    else:
        return [], table_start, table_start

    # Parse data rows
    current_line = separator_line + 1
    while current_line < len(lines):
        line = lines[current_line]
        if not line.startswith("|"):
            break

        cols = [col.strip() for col in line.split("|")[1:-1]]
        if len(cols) > status_col_idx:
            rows.append({
                "line_num": current_line,
                "content": line,
                "columns": cols,
                "status_col_idx": status_col_idx,
            })

        current_line += 1

    table_end_pos = content.find("\n", table_start)
    for _ in range(current_line - start_line):
        table_end_pos = content.find("\n", table_end_pos + 1)

    return rows, table_start, table_end_pos


def find_table_section(content: str, section_name: str) -> tuple[int, int]:
    """Find start and end of a table section by heading.

    Returns:
        (start_pos, end_pos) or (-1, -1) if not found
    """
    # Look for heading - try both ## and ###
    pattern = rf"^###?\s+{re.escape(section_name)}\s*$"
    match = re.search(pattern, content, re.MULTILINE)
    if not match:
        return -1, -1

    start_pos = match.end()

    # Find next heading of same or higher level or end of file
    next_heading = re.search(r"^###?\s+", content[start_pos:], re.MULTILINE)
    if next_heading:
        end_pos = start_pos + next_heading.start()
    else:
        end_pos = len(content)

    return start_pos, end_pos


def update_table_status(
    content: str,
    table_name: str,
    yaml_statuses: dict[str, str],
    *,
    verbose: bool = False,
) -> tuple[str, int]:
    """Update status column in a specific SOT table.

    Returns:
        (updated_content, num_changes)
    """
    start_pos, end_pos = find_table_section(content, table_name)
    if start_pos == -1:
        if verbose:
            print(f"⚠️  Table section '{table_name}' not found")
        return content, 0

    table_content = content[start_pos:end_pos]
    rows, table_rel_start, table_rel_end = parse_markdown_table(
        content, start_pos + table_content.find("|")
    )

    if not rows:
        if verbose:
            print(f"⚠️  Could not parse table in '{table_name}' section")
        return content, 0

    changes = 0
    lines = content.split("\n")

    for row in rows:
        cols = row["columns"]
        name_col = cols[0]  # First column is usually the name

        # Match to YAML status
        if name_col in yaml_statuses:
            new_status = yaml_statuses[name_col]
            old_status = cols[row["status_col_idx"]]

            if new_status != old_status:
                # Update only the status column
                cols[row["status_col_idx"]] = new_status
                new_line = "| " + " | ".join(cols) + " |"
                lines[row["line_num"]] = new_line
                changes += 1

                if verbose:
                    print(f"  {name_col}: {old_status} → {new_status}")

    if changes > 0:
        return "\n".join(lines), changes

    return content, 0


def create_backup(filepath: Path) -> Path:
    """Create a timestamped backup of a file."""
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    backup_path = filepath.with_suffix(f".{timestamp}.bak")
    backup_path.write_text(filepath.read_text(encoding="utf-8"), encoding="utf-8")
    return backup_path


def main() -> int:
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Sync YAML overall_status to SOT tables",
    )
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Apply changes (default is dry-run)",
    )
    parser.add_argument(
        "--strict",
        action="store_true",
        help="Exit with code 1 if any drift detected",
    )
    parser.add_argument(
        "-v",
        "--verbose",
        action="store_true",
        help="Verbose output",
    )

    args = parser.parse_args()

    # Read SOT
    if not SOT_FILE.exists():
        print(f"❌ SOT file not found: {SOT_FILE}", file=sys.stderr)
        return 1

    sot_content = SOT_FILE.read_text(encoding="utf-8")

    # Collect YAML statuses
    yaml_statuses = find_all_yaml_statuses()

    if args.verbose:
        print(f"Found {len(yaml_statuses)} YAML status mappings")

    # Update tables
    tables_to_sync = [
        "Content Modules",
        "Backend Services",
        "Metadata Providers",
        "Arr Ecosystem",
    ]

    total_changes = 0
    updated_content = sot_content

    for table_name in tables_to_sync:
        if args.verbose:
            print(f"\nProcessing: {table_name}")

        updated_content, changes = update_table_status(
            updated_content,
            table_name,
            yaml_statuses,
            verbose=args.verbose,
        )
        total_changes += changes

    # Report
    print(f"\n{'=' * 50}")
    print("STATUS SYNC SUMMARY")
    print(f"{'=' * 50}")
    print(f"YAML files scanned: {len(yaml_statuses)}")
    print(f"Tables processed: {len(tables_to_sync)}")
    print(f"Status changes detected: {total_changes}")

    if total_changes > 0:
        if args.apply:
            # Create backup
            backup_path = create_backup(SOT_FILE)
            print(f"\nBackup created: {backup_path.name}")

            # Write changes
            SOT_FILE.write_text(updated_content, encoding="utf-8")
            print(f"✅ Applied {total_changes} changes to {SOT_FILE.name}")
        else:
            print("\n=== DRY RUN complete. Use --apply to write changes. ===")
    else:
        print("\n✅ All statuses in sync!")

    if args.strict and total_changes > 0:
        print("\n❌ --strict mode: Drift detected, exiting with code 1")
        return 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
