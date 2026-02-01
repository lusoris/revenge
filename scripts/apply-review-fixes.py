#!/usr/bin/env python3
"""Apply 'Review' status link fixes from FIXES_REPORT.md.

This script applies link fixes that were marked as 'Review' status,
meaning they have suggestions but need manual verification.

Author: Automation System
Created: 2026-02-01
"""

import re
import sys
from pathlib import Path


def parse_fixes_report(report_path: Path) -> dict[str, list[tuple[int, str, str]]]:
    """Parse FIXES_REPORT.md and extract Review fixes.

    Returns:
        Dict mapping file paths to list of (line_num, old_link, new_link) tuples
    """
    with open(report_path) as f:
        content = f.read()

    fixes = {}
    current_file = None

    # Parse the report
    for line in content.split("\n"):
        # Match file headers like "### docs/dev/design/services/EPG.md"
        file_match = re.match(r"^### (.+\.md)$", line)
        if file_match:
            current_file = file_match.group(1)
            continue

        # Match table rows with Review status
        # Format: | 265 | `features/livetv/LIVE_TV_DVR.md` | Review | `../features/livetv/LIVE_TV_DVR.md` |
        row_match = re.match(
            r"^\|\s*(\d+)\s*\|\s*`([^`]+)`\s*\|\s*Review\s*\|\s*`([^`]+)`\s*\|", line
        )
        if row_match and current_file:
            line_num = int(row_match.group(1))
            old_link = row_match.group(2)
            new_link = row_match.group(3)

            if current_file not in fixes:
                fixes[current_file] = []
            fixes[current_file].append((line_num, old_link, new_link))

    return fixes


def apply_fix(file_path: Path, line_num: int, old_link: str, new_link: str) -> bool:
    """Apply a single fix to a file.

    Returns:
        True if fix was applied, False otherwise
    """
    try:
        with open(file_path) as f:
            lines = f.readlines()

        # Line numbers in report are 1-indexed
        if line_num < 1 or line_num > len(lines):
            print(f"    ⚠️  Line {line_num} out of range")
            return False

        line_content = lines[line_num - 1]

        # Replace the old link with the new one
        # Handle both markdown link format [text](link) and bare links
        new_content = line_content.replace(f"]({old_link})", f"]({new_link})")
        new_content = new_content.replace(f"`{old_link}`", f"`{new_link}`")

        if new_content == line_content:
            print(f"    ⚠️  No change on line {line_num}")
            return False

        lines[line_num - 1] = new_content

        with open(file_path, "w") as f:
            f.writelines(lines)

        return True

    except Exception as e:
        print(f"    ❌ Error applying fix: {e}")
        return False


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent
    report_path = (
        repo_root / "docs" / "dev" / "design" / ".analysis" / "FIXES_REPORT.md"
    )

    if not report_path.exists():
        print(f"❌ Error: {report_path} not found")
        print("Run the link checker first to generate the report.")
        sys.exit(1)

    print("🔍 Parsing FIXES_REPORT.md...")
    fixes = parse_fixes_report(report_path)

    if not fixes:
        print("✅ No Review fixes found")
        return

    print(f"📁 Found {len(fixes)} files with Review fixes")

    total_applied = 0

    for file_path_str, file_fixes in sorted(fixes.items()):
        file_path = repo_root / file_path_str

        if not file_path.exists():
            print(f"\n⚠️  {file_path_str} - File not found")
            continue

        print(f"\n📄 {file_path_str}")
        print(f"   {len(file_fixes)} fixes to apply")

        applied = 0
        for line_num, old_link, new_link in file_fixes:
            if apply_fix(file_path, line_num, old_link, new_link):
                print(f"    ✓ Line {line_num}: {old_link} → {new_link}")
                applied += 1
            else:
                print(f"    ✗ Line {line_num}: Failed to apply fix")

        total_applied += applied

    print(f"\n✅ Applied {total_applied} fixes across {len(fixes)} files")


if __name__ == "__main__":
    main()
