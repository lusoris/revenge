#!/usr/bin/env python3
"""
Replace placeholder text with proper 🔴 Not implemented markers.

Replaces:
- TBD → 🔴 Not implemented
- TODO → 🔴 Not implemented
- placeholder → 🔴 Not implemented
- [placeholder] → 🔴 Not implemented

Usage:
    python3 scripts/fix-placeholders.py           # Dry-run
    python3 scripts/fix-placeholders.py --fix    # Apply changes
"""

import argparse
import re
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"

# Patterns to replace
REPLACEMENTS = [
    # Standalone TBD
    (r"\bTBD\b(?!\s*[-:])", "🔴 Not implemented"),
    # TBD with description
    (r"\bTBD\s*[-:]\s*", "🔴 Not implemented - "),
    # TODO items (but not code comments)
    (r"^(?!\s*//|\s*#|\s*\*)\s*TODO\s*[-:]?\s*", "🔴 TODO: ", re.MULTILINE),
    # Placeholder text
    (r"\[placeholder\]", "🔴 Not implemented"),
    (r"\bplaceholder content\b", "🔴 Not implemented"),
    (r"\bplaceholder\b(?!\s+for)", "🔴 Not implemented"),
]


def fix_file(filepath: Path, dry_run: bool = True) -> list:
    """Fix placeholders in a file."""
    try:
        content = filepath.read_text(encoding="utf-8")
        original = content
    except Exception:
        return []

    changes = []

    for pattern, replacement, *flags in REPLACEMENTS:
        flag = flags[0] if flags else 0
        matches = list(re.finditer(pattern, content, flag))

        for match in matches:
            old_text = match.group(0)
            changes.append(
                {
                    "line": content[: match.start()].count("\n") + 1,
                    "old": old_text.strip(),
                    "new": replacement.strip()
                    if not replacement.endswith(" ")
                    else replacement.strip() + " ...",
                }
            )

        content = re.sub(pattern, replacement, content, flags=flag)

    if content != original and not dry_run:
        filepath.write_text(content, encoding="utf-8")

    return changes if content != original else []


def main():
    parser = argparse.ArgumentParser(description="Fix placeholder text")
    parser.add_argument("--fix", action="store_true", help="Apply changes")
    parser.add_argument("--verbose", "-v", action="store_true", help="Show details")
    args = parser.parse_args()

    total_changes = 0
    files_changed = 0

    for filepath in sorted(DESIGN_DIR.rglob("*.md")):
        if ".archive" in str(filepath) or ".analysis" in str(filepath):
            continue

        changes = fix_file(filepath, dry_run=not args.fix)

        if changes:
            files_changed += 1
            total_changes += len(changes)
            rel_path = filepath.relative_to(PROJECT_ROOT)

            if args.verbose or not args.fix:
                print(f"\n{rel_path}:")
                for change in changes:
                    print(
                        f"  Line {change['line']}: '{change['old']}' → '{change['new']}'"
                    )

    print(
        f"\n{'Applied' if args.fix else 'Would apply'} {total_changes} changes in {files_changed} files"
    )

    if not args.fix and total_changes > 0:
        print("\nRun with --fix to apply changes")


if __name__ == "__main__":
    main()
