#!/usr/bin/env python3
"""
Generate NAVIGATION.md - A comprehensive navigation map for the design documentation.

This script creates a single-page navigation hub that provides:
- Category-based navigation
- Topic-based quick links
- Depth-aware document listing
- Status overview per category

Usage:
    python3 scripts/generate-navigation-map.py [--output PATH]

The generated NAVIGATION.md is linked from SOT and provides
an alternative entry point to the documentation.
"""

import argparse
import re
from datetime import datetime
from pathlib import Path


# Project paths
SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
DEFAULT_OUTPUT = DESIGN_DIR / "NAVIGATION.md"

# Category descriptions
CATEGORY_INFO = {
    "architecture": {
        "icon": "🏗️",
        "description": "System design, principles, and core architecture",
        "priority": 1,
    },
    "features": {
        "icon": "✨",
        "description": "Content modules and feature specifications",
        "priority": 2,
    },
    "services": {
        "icon": "⚙️",
        "description": "Backend services and business logic",
        "priority": 3,
    },
    "integrations": {
        "icon": "🔌",
        "description": "External APIs, metadata providers, and Arr stack",
        "priority": 4,
    },
    "technical": {
        "icon": "🔧",
        "description": "API design, frontend, and configuration",
        "priority": 5,
    },
    "operations": {
        "icon": "🚀",
        "description": "Deployment, setup, and best practices",
        "priority": 6,
    },
    "research": {
        "icon": "🔬",
        "description": "User research and UX/UI resources",
        "priority": 7,
    },
    "planning": {
        "icon": "📋",
        "description": "Project planning and versioning",
        "priority": 8,
    },
}

# Status emoji mapping
STATUS_EMOJI = {
    "complete": "✅",
    "partial": "🟡",
    "scaffold": "🟡",
    "planned": "🔴",
    "not started": "🔴",
}


def extract_status(content: str) -> str:
    """Extract status from document content."""
    # Look for status table
    status_match = re.search(r"\|\s*Design\s*\|\s*([✅🟡🔴⚪])", content)
    if status_match:
        emoji = status_match.group(1)
        return emoji

    # Look for status in header
    for keyword, emoji in STATUS_EMOJI.items():
        if keyword.lower() in content.lower()[:500]:
            return emoji

    return "⚪"


def extract_title(content: str, filepath: Path) -> str:
    """Extract title from markdown file."""
    # Look for H1
    match = re.search(r"^#\s+(.+)$", content, re.MULTILINE)
    if match:
        return match.group(1).strip()

    # Fallback to filename
    return filepath.stem.replace("_", " ").replace("-", " ").title()


def extract_description(content: str) -> str:
    """Extract description (blockquote after title)."""
    match = re.search(r"^#\s+.+\n+>\s*(.+)$", content, re.MULTILINE)
    if match:
        return match.group(1).strip()
    return ""


def get_category_docs(category_dir: Path) -> list:
    """Get all documents in a category with metadata."""
    docs = []

    for md_file in sorted(category_dir.rglob("*.md")):
        # Skip INDEX files for main listing
        if md_file.name == "INDEX.md":
            continue

        # Calculate relative path from design dir
        rel_path = md_file.relative_to(DESIGN_DIR)

        # Calculate depth (subdirectories within category)
        parts = rel_path.parts
        depth = len(parts) - 2  # Subtract category and filename

        try:
            content = md_file.read_text(encoding="utf-8")
            title = extract_title(content, md_file)
            status = extract_status(content)
            description = extract_description(content)
        except Exception:
            title = md_file.stem.replace("_", " ").title()
            status = "⚪"
            description = ""

        # Determine subcategory
        subcategory = parts[1] if len(parts) > 2 else None

        docs.append(
            {
                "path": rel_path,
                "title": title,
                "status": status,
                "description": description,
                "depth": depth,
                "subcategory": subcategory,
                "filename": md_file.name,
            }
        )

    return docs


def get_subcategories(category_dir: Path) -> list:
    """Get subcategories (subdirectories) in a category."""
    subcats = []
    for subdir in sorted(category_dir.iterdir()):
        if subdir.is_dir() and not subdir.name.startswith("."):
            index_file = subdir / "INDEX.md"
            if index_file.exists():
                try:
                    content = index_file.read_text(encoding="utf-8")
                    title = extract_title(content, index_file)
                except Exception:
                    title = subdir.name.replace("_", " ").title()
            else:
                title = subdir.name.replace("_", " ").title()

            doc_count = len(list(subdir.glob("*.md")))
            subcats.append(
                {
                    "name": subdir.name,
                    "title": title,
                    "doc_count": doc_count,
                    "path": subdir.relative_to(DESIGN_DIR),
                }
            )

    return subcats


def generate_navigation_map() -> str:
    """Generate the complete navigation map."""
    lines = [
        "# Navigation Map",
        "",
        "> Comprehensive navigation hub for design documentation",
        ">",
        f"> Auto-generated: {datetime.now().strftime('%Y-%m-%d %H:%M')}",
        ">",
        "> **Master Reference**: [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md)",
        "",
        "---",
        "",
        "## Quick Start",
        "",
        "| If you want to... | Start here |",
        "|-------------------|------------|",
        "| Understand the system | [Architecture](architecture/INDEX.md) → [01_ARCHITECTURE.md](architecture/01_ARCHITECTURE.md) |",
        "| Learn about a content type | [Features](features/INDEX.md) → Pick a module |",
        "| Implement a service | [Services](services/INDEX.md) → Pick a service |",
        "| Add an integration | [Integrations](integrations/INDEX.md) → Pick a provider |",
        "| Deploy the system | [Operations](operations/INDEX.md) → [SETUP.md](operations/SETUP.md) |",
        "| Check package versions | [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md) |",
        "",
        "---",
        "",
    ]

    # Sort categories by priority
    sorted_categories = sorted(CATEGORY_INFO.items(), key=lambda x: x[1]["priority"])

    # Category overview table
    lines.extend(
        [
            "## Categories Overview",
            "",
            "| Category | Description | Docs |",
            "|----------|-------------|------|",
        ]
    )

    for cat_name, cat_info in sorted_categories:
        cat_dir = DESIGN_DIR / cat_name
        if not cat_dir.exists():
            continue

        docs = get_category_docs(cat_dir)
        doc_count = len(docs)
        icon = cat_info["icon"]
        desc = cat_info["description"]

        lines.append(
            f"| {icon} [{cat_name.title()}](#{cat_name}) | {desc} | {doc_count} |"
        )

    lines.extend(["", "---", ""])

    # Detailed sections for each category
    for cat_name, cat_info in sorted_categories:
        cat_dir = DESIGN_DIR / cat_name
        if not cat_dir.exists():
            continue

        icon = cat_info["icon"]
        desc = cat_info["description"]

        lines.extend(
            [
                f"## {icon} {cat_name.title()}",
                "",
                f"> {desc}",
                "",
                f"**Index**: [{cat_name}/INDEX.md]({cat_name}/INDEX.md)",
                "",
            ]
        )

        # Get subcategories
        subcats = get_subcategories(cat_dir)

        if subcats:
            lines.extend(
                [
                    "### Subcategories",
                    "",
                    "| Subcategory | Documents |",
                    "|-------------|-----------|",
                ]
            )

            for subcat in subcats:
                lines.append(
                    f"| [{subcat['title']}]({subcat['path']}/INDEX.md) | {subcat['doc_count']} |"
                )

            lines.append("")

        # Get direct documents (not in subcategories)
        docs = get_category_docs(cat_dir)
        direct_docs = [d for d in docs if d["depth"] == 0]

        if direct_docs:
            lines.extend(
                [
                    "### Documents",
                    "",
                    "| Document | Status | Description |",
                    "|----------|--------|-------------|",
                ]
            )

            for doc in sorted(direct_docs, key=lambda x: x["title"]):
                title = doc["title"]
                status = doc["status"]
                desc = (
                    doc["description"][:60] + "..."
                    if len(doc["description"]) > 60
                    else doc["description"]
                )
                path = doc["path"]

                lines.append(f"| [{title}]({path}) | {status} | {desc} |")

            lines.append("")

        lines.extend(["---", ""])

    # Deep links section for Q16 - shortcuts to deep directories
    lines.extend(
        [
            "## Deep Directory Shortcuts",
            "",
            "> Direct links to deeply nested documentation (depth 3+)",
            "",
            "| Path | Description |",
            "|------|-------------|",
        ]
    )

    deep_dirs = [
        ("integrations/metadata/adult/", "Adult metadata providers (StashDB, TPDB)"),
        ("integrations/metadata/books/", "Book metadata (OpenLibrary, Audible)"),
        ("integrations/metadata/comics/", "Comics metadata (ComicVine)"),
        ("integrations/metadata/music/", "Music metadata (MusicBrainz, Last.fm)"),
        ("integrations/metadata/video/", "Video metadata (TMDb, TheTVDB)"),
        ("integrations/wiki/adult/", "Adult wiki sources (IAFD, Babepedia)"),
        ("integrations/external/adult/", "Adult external sources"),
    ]

    for path, desc in deep_dirs:
        full_path = DESIGN_DIR / path
        if full_path.exists():
            lines.append(f"| [{path}]({path}INDEX.md) | {desc} |")

    lines.extend(
        [
            "",
            "---",
            "",
            "## Cross-References",
            "",
            "| Resource | Description |",
            "|----------|-------------|",
            "| [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md) | Package versions, module status, config keys |",
            "| [DESIGN_INDEX.md](../sources/DESIGN_INDEX.md) | Auto-generated index of all design docs |",
            "| [SOURCES_INDEX.md](../sources/SOURCES_INDEX.md) | Index of external documentation sources |",
            "| [DESIGN_CROSSREF.md](../sources/DESIGN_CROSSREF.md) | Design ↔ Sources cross-reference map |",
            "",
            "---",
            "",
            "*This file is auto-generated by `scripts/generate-navigation-map.py`*",
            "",
        ]
    )

    return "\n".join(lines)


def main():
    parser = argparse.ArgumentParser(
        description="Generate NAVIGATION.md for design documentation"
    )
    parser.add_argument(
        "--output", "-o", type=Path, default=DEFAULT_OUTPUT, help="Output file path"
    )
    parser.add_argument(
        "--dry-run", action="store_true", help="Print to stdout instead of writing file"
    )

    args = parser.parse_args()

    content = generate_navigation_map()

    if args.dry_run:
        print(content)
    else:
        args.output.write_text(content, encoding="utf-8")
        print(f"Generated: {args.output}")

        # Count stats
        cat_count = len(
            [
                d
                for d in DESIGN_DIR.iterdir()
                if d.is_dir() and not d.name.startswith(".")
            ]
        )
        doc_count = len(list(DESIGN_DIR.rglob("*.md")))
        print(f"  Categories: {cat_count}")
        print(f"  Documents: {doc_count}")


if __name__ == "__main__":
    main()
