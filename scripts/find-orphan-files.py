#!/usr/bin/env python3
"""
Find orphan files not linked from any INDEX.md or other document.

Identifies:
- Markdown files not referenced from any index
- Files not cross-referenced from other docs
- Potentially abandoned or forgotten documentation

Usage:
    python scripts/find-orphan-files.py                    # Check all
    python scripts/find-orphan-files.py --include-sources  # Include sources/
    python scripts/find-orphan-files.py --fix              # Add to nearest INDEX.md
"""

import argparse
import re
import sys
from collections import defaultdict
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"

# Files that are allowed to be "orphans" (root-level special files)
ALLOWED_ORPHANS = {
    "00_SOURCE_OF_TRUTH.md",
    "01_DESIGN_DOC_TEMPLATE.md",
    "02_QUESTIONS_TO_DISCUSS.md",
    "03_DESIGN_DOCS_STATUS.md",
    "DESIGN_INDEX.md",
    "SOURCES_INDEX.md",
    "DESIGN_CROSSREF.md",
    "SOURCES.yaml",
    "INDEX.yaml",
}

# Link pattern
LINK_PATTERN = re.compile(r"\[([^\]]*)\]\(([^)]+\.md)\)")


def find_all_markdown_files(directory: Path) -> list[Path]:
    """Find all markdown files in directory."""
    files = []
    for md_file in sorted(directory.rglob("*.md")):
        if ".archive" in str(md_file):
            continue
        files.append(md_file)
    return files


def extract_links(file_path: Path) -> set[Path]:
    """Extract all markdown file links from a document."""
    content = file_path.read_text(encoding="utf-8")
    links = set()

    for match in LINK_PATTERN.finditer(content):
        link_url = match.group(2)

        # Skip external links
        if link_url.startswith(("http://", "https://")):
            continue

        # Resolve relative path
        try:
            if "#" in link_url:
                link_url = link_url.split("#")[0]
            if not link_url:
                continue

            target = (file_path.parent / link_url).resolve()
            if target.exists():
                links.add(target)
        except (ValueError, OSError):
            pass

    return links


def build_link_graph(files: list[Path]) -> dict[Path, set[Path]]:
    """Build a graph of which files link to which other files."""
    graph = defaultdict(set)

    for file_path in files:
        links = extract_links(file_path)
        for target in links:
            graph[target].add(file_path)

    return dict(graph)


def find_index_files(directory: Path) -> list[Path]:
    """Find all INDEX.md files."""
    indexes = []
    for md_file in directory.rglob("*.md"):
        if md_file.name in ["INDEX.md", "00_INDEX.md"]:
            indexes.append(md_file)
    return indexes


def main():
    parser = argparse.ArgumentParser(description="Find orphan documentation files")
    parser.add_argument(
        "--include-sources",
        action="store_true",
        help="Include sources/ directory",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Show linking details",
    )
    args = parser.parse_args()

    directories = [DESIGN_DIR]
    if args.include_sources:
        directories.append(SOURCES_DIR)

    all_files = []
    for directory in directories:
        all_files.extend(find_all_markdown_files(directory))

    print(f"Found {len(all_files)} markdown files")

    # Build link graph
    print("Building link graph...")
    link_graph = build_link_graph(all_files)

    # Find files linked from indexes
    index_files = []
    for directory in directories:
        index_files.extend(find_index_files(directory))

    print(f"Found {len(index_files)} index files")

    # Track which files are linked
    linked_from_index = set()
    linked_from_any = set()

    for index_path in index_files:
        links = extract_links(index_path)
        linked_from_index.update(links)

    for file_path in all_files:
        links = extract_links(file_path)
        linked_from_any.update(links)

    # Categorize orphans
    orphans_no_index = []  # Not in any INDEX
    orphans_no_links = []  # Not linked from anywhere

    for file_path in all_files:
        # Skip index files themselves
        if file_path.name in ["INDEX.md", "00_INDEX.md"]:
            continue

        # Skip allowed orphans
        if file_path.name in ALLOWED_ORPHANS:
            continue

        # Check if in design dir root (special handling)
        try:
            rel_path = file_path.relative_to(DESIGN_DIR)
            if len(rel_path.parts) == 1:
                continue  # Root-level files are okay
        except ValueError:
            pass

        # Check orphan status
        in_index = file_path in linked_from_index
        in_any = file_path in linked_from_any or file_path in link_graph

        if not in_index:
            orphans_no_index.append(file_path)

        if not in_any:
            orphans_no_links.append(file_path)

    # Report
    print(f"\n{'=' * 60}")
    print("ORPHAN FILE ANALYSIS")
    print(f"{'=' * 60}")

    if orphans_no_index:
        print(f"\n## Not linked from any INDEX.md ({len(orphans_no_index)} files)")
        print("These files exist but aren't listed in their directory's INDEX.md\n")

        by_dir = defaultdict(list)
        for f in orphans_no_index:
            try:
                rel = f.relative_to(DESIGN_DIR)
                by_dir[str(rel.parent)].append(f)
            except ValueError:
                try:
                    rel = f.relative_to(SOURCES_DIR)
                    by_dir[f"sources/{rel.parent}"].append(f)
                except ValueError:
                    by_dir["other"].append(f)

        for dir_name, files in sorted(by_dir.items()):
            print(f"  {dir_name}/")
            for f in sorted(files):
                linked_count = len(link_graph.get(f, []))
                status = (
                    f"(linked from {linked_count} docs)"
                    if linked_count
                    else "(no links)"
                )
                print(f"    - {f.name} {status}")

    if orphans_no_links:
        print(f"\n## Not linked from ANYWHERE ({len(orphans_no_links)} files)")
        print("These files have NO incoming links from any document\n")

        for f in sorted(orphans_no_links):
            try:
                rel = f.relative_to(PROJECT_ROOT)
            except ValueError:
                rel = f
            print(f"  - {rel}")

    # Summary
    well_linked = len(all_files) - len(orphans_no_index)
    print(f"\n{'=' * 60}")
    print("SUMMARY")
    print(f"{'=' * 60}")
    print(f"Total files: {len(all_files)}")
    print(f"Linked from INDEX: {well_linked} ({100 * well_linked // len(all_files)}%)")
    print(f"Not in any INDEX: {len(orphans_no_index)}")
    print(f"Completely orphaned: {len(orphans_no_links)}")

    if args.verbose:
        print("\n## Link Statistics")
        most_linked = sorted(
            [(f, len(sources)) for f, sources in link_graph.items()],
            key=lambda x: x[1],
            reverse=True,
        )[:10]

        print("\nMost linked files:")
        for f, count in most_linked:
            try:
                rel = f.relative_to(PROJECT_ROOT)
            except ValueError:
                rel = f
            print(f"  {count:3d} links → {rel}")

    if orphans_no_links:
        print("\n⚠️  Consider removing or linking the completely orphaned files")
        return 1

    print("\n✓ No completely orphaned files found")
    return 0


if __name__ == "__main__":
    sys.exit(main())
