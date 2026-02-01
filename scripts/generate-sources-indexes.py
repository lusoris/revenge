#!/usr/bin/env python3
"""
Generate INDEX.md files for all docs/dev/sources/ subdirectories.

Creates browsable index files for each source category showing:
- All fetched source documents
- Fetch status and content hash
- Links to related design docs

Usage:
    python scripts/generate-sources-indexes.py              # Dry run
    python scripts/generate-sources-indexes.py --update     # Write files
"""

import argparse
from pathlib import Path

import yaml


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
SOURCES_YAML = SOURCES_DIR / "SOURCES.yaml"
INDEX_YAML = SOURCES_DIR / "INDEX.yaml"

# Category metadata
CATEGORY_META = {
    "apis": {
        "title": "API Documentation",
        "desc": "External API references and OpenAPI specs",
    },
    "database": {
        "title": "Database Documentation",
        "desc": "PostgreSQL, sqlc, and database tooling",
    },
    "frontend": {
        "title": "Frontend Documentation",
        "desc": "Svelte, SvelteKit, and UI libraries",
    },
    "go": {
        "title": "Go Documentation",
        "desc": "Go standard library and language references",
    },
    "go/stdlib": {
        "title": "Go Standard Library",
        "desc": "Core Go packages documentation",
    },
    "go/x": {
        "title": "Go Extended Libraries",
        "desc": "golang.org/x/ packages",
    },
    "infrastructure": {
        "title": "Infrastructure Documentation",
        "desc": "Dragonfly, Typesense, and infrastructure components",
    },
    "media": {
        "title": "Media Processing",
        "desc": "FFmpeg, image processing, and media libraries",
    },
    "observability": {
        "title": "Observability",
        "desc": "Prometheus, OpenTelemetry, logging, and monitoring",
    },
    "protocols": {
        "title": "Protocols & Standards",
        "desc": "HLS, WebRTC, HTTP specs, and streaming protocols",
    },
    "security": {
        "title": "Security",
        "desc": "OAuth, OIDC, JWT, and authentication standards",
    },
    "testing": {
        "title": "Testing",
        "desc": "Go testing patterns and frameworks",
    },
    "tooling": {
        "title": "Go Tooling",
        "desc": "Libraries and tools used in the backend",
    },
    "standards": {
        "title": "Standards & Conventions",
        "desc": "Versioning, commit conventions, and workflows",
    },
    "casting": {
        "title": "Casting Protocols",
        "desc": "Chromecast, DLNA, and device casting",
    },
    "distributed": {
        "title": "Distributed Systems",
        "desc": "Clustering, consensus, and distributed patterns",
    },
    "livetv": {
        "title": "Live TV",
        "desc": "EPG, DVR, and live streaming",
    },
    "orchestration": {
        "title": "Orchestration",
        "desc": "Kubernetes, Docker, and deployment",
    },
    "wiki": {
        "title": "Wiki APIs",
        "desc": "Wikipedia, Fandom, and wiki integration",
    },
}


def load_sources_config() -> dict:
    """Load SOURCES.yaml configuration."""
    if SOURCES_YAML.exists():
        with open(SOURCES_YAML) as f:
            return yaml.safe_load(f) or {}
    return {}


def load_index() -> dict:
    """Load existing INDEX.yaml with fetch status."""
    if INDEX_YAML.exists():
        with open(INDEX_YAML) as f:
            return yaml.safe_load(f) or {}
    return {}


def find_source_directories() -> list[Path]:
    """Find all directories in sources/ that contain .md files."""
    dirs = set()
    for md_file in SOURCES_DIR.rglob("*.md"):
        # Skip root-level files
        if md_file.parent == SOURCES_DIR:
            continue
        # Skip INDEX.md files
        if md_file.name == "INDEX.md":
            continue
        dirs.add(md_file.parent)
    return sorted(dirs)


def get_source_files(directory: Path) -> list[Path]:
    """Get all source files in a directory (non-recursive)."""
    files = []
    for f in sorted(directory.iterdir()):
        if f.is_file() and f.suffix in [".md", ".json", ".graphql", ".yaml"]:
            if f.name != "INDEX.md":
                files.append(f)
    return files


def get_file_info(file_path: Path, index_data: dict) -> dict:
    """Get information about a source file."""
    # Try to find in index
    sources = index_data.get("sources", {})

    # Find matching source by output file
    rel_path = str(file_path.relative_to(SOURCES_DIR))
    for source_id, info in sources.items():
        if info.get("output") == rel_path:
            return {
                "id": source_id,
                "name": info.get("name", source_id),
                "hash": info.get("content_hash", ""),
                "fetched": info.get("last_fetched", ""),
                "status": "fetched",
            }

    # Not in index - manual file
    return {
        "id": file_path.stem,
        "name": file_path.stem.replace("-", " ").replace("_", " ").title(),
        "hash": "",
        "fetched": "",
        "status": "manual",
    }


def generate_index(directory: Path, files: list[Path], index_data: dict) -> str:
    """Generate INDEX.md content for a directory."""
    rel_dir = directory.relative_to(SOURCES_DIR)
    dir_key = str(rel_dir).replace("\\", "/")

    meta = CATEGORY_META.get(
        dir_key,
        {
            "title": directory.name.replace("_", " ").replace("-", " ").title(),
            "desc": "",
        },
    )

    # Calculate parent path
    depth = len(rel_dir.parts)
    parent_path = "../" * depth

    lines = [
        f"# {meta['title']}",
        "",
        f"← Back to [Sources]({parent_path}SOURCES_INDEX.md)",
        "",
    ]

    if meta.get("desc"):
        lines.extend([f"> {meta['desc']}", ""])

    lines.extend(
        [
            "---",
            "",
            "## Documents",
            "",
            "| Document | Status | Last Fetched |",
            "|----------|--------|--------------|",
        ],
    )

    for file_path in files:
        info = get_file_info(file_path, index_data)
        name = info["name"]
        rel_path = file_path.name

        if info["status"] == "fetched":
            status = "✅ Auto"
            fetched = info["fetched"][:10] if info["fetched"] else "-"
        else:
            status = "📝 Manual"
            fetched = "-"

        # Determine file type icon
        if file_path.suffix == ".json":
            icon = "📋"
        elif file_path.suffix == ".graphql":
            icon = "🔷"
        else:
            icon = "📄"

        lines.append(f"| {icon} [{name}]({rel_path}) | {status} | {fetched} |")

    lines.extend(
        [
            "",
            "---",
            "",
            "## Legend",
            "",
            "- ✅ Auto = Fetched automatically by `fetch-sources.py`",
            "- 📝 Manual = Manually maintained document",
            "- 📄 Markdown | 📋 JSON/OpenAPI | 🔷 GraphQL",
            "",
            "---",
            "",
            f"*See [SOURCES.yaml]({parent_path}SOURCES.yaml) for fetch configuration*",
            "",
        ],
    )

    return "\n".join(lines)


def main():
    parser = argparse.ArgumentParser(description="Generate sources index files")
    parser.add_argument(
        "--update",
        "-u",
        action="store_true",
        help="Write files (default: dry run)",
    )
    args = parser.parse_args()

    print("Loading configuration...")
    load_sources_config()
    index_data = load_index()
    print(f"  {len(index_data.get('sources', {}))} sources in index")

    print("Finding source directories...")
    directories = find_source_directories()
    print(f"  Found {len(directories)} directories")

    updated = 0
    unchanged = 0

    for directory in directories:
        files = get_source_files(directory)
        if not files:
            continue

        index_path = directory / "INDEX.md"
        new_content = generate_index(directory, files, index_data)

        # Check if update needed
        if index_path.exists():
            old_content = index_path.read_text(encoding="utf-8")
            if old_content == new_content:
                rel_path = index_path.relative_to(SOURCES_DIR)
                print(f"  {rel_path} - unchanged")
                unchanged += 1
                continue

        rel_path = index_path.relative_to(SOURCES_DIR)
        if args.update:
            index_path.write_text(new_content, encoding="utf-8")
            print(f"  {rel_path} - updated")
        else:
            print(f"  {rel_path} - would update")
        updated += 1

    print("\n=== SUMMARY ===")
    print(f"{'Updated' if args.update else 'Would update'}: {updated}")
    print(f"Unchanged: {unchanged}")

    if not args.update and updated > 0:
        print("\nRun with --update to write changes")


if __name__ == "__main__":
    main()
