#!/usr/bin/env python3
"""Generate and update INDEX.md files for all design doc directories.

Creates standardized INDEX.md files for each directory containing:
- Category title and description
- Document listings with status
- Links to subdirectories
- Related category links
- SOT reference

Usage:
    python scripts/doc-pipeline/01-indexes.py           # Dry run
    python scripts/doc-pipeline/01-indexes.py --apply   # Write files
"""

from __future__ import annotations

import argparse
import re
import sys
from collections import defaultdict
from pathlib import Path
from typing import Any


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
SOT_FILE = DESIGN_DIR / "00_SOURCE_OF_TRUTH.md"

CATEGORY_META = {
    "architecture": {
        "title": "Architecture Documentation",
        "desc": "Core system design and architectural decisions",
        "related": ["technical", "features"],
    },
    "features": {
        "title": "Features Documentation",
        "desc": "Content modules and feature specifications",
        "related": ["architecture", "integrations"],
    },
    "features/shared": {
        "title": "Shared Features",
        "desc": "Features that apply across all modules",
        "related": ["features", "architecture"],
    },
    "features/playback": {
        "title": "Playback Features",
        "desc": "Media playback and streaming features",
        "related": ["features/shared", "integrations/transcoding"],
    },
    "features/video": {
        "title": "Video Module",
        "desc": "Movies and TV Shows features",
        "related": ["integrations/metadata/video", "integrations/servarr"],
    },
    "features/photos": {
        "title": "Photos Module",
        "desc": "Photo library features",
        "related": ["features/shared"],
    },
    "features/podcasts": {
        "title": "Podcasts Module",
        "desc": "Podcast management features",
        "related": ["integrations/audiobook"],
    },
    "features/livetv": {
        "title": "Live TV Module",
        "desc": "Live TV and DVR features",
        "related": ["integrations/livetv"],
    },
    "features/comics": {
        "title": "Comics Module",
        "desc": "Comic book features",
        "related": ["integrations/metadata/comics"],
    },
    "features/adult": {
        "title": "Adult Content Module",
        "desc": "Adult content features (isolated in qar schema)",
        "related": ["integrations/metadata/adult", "integrations/external/adult"],
    },
    "integrations": {
        "title": "External Integrations",
        "desc": "Third-party services and APIs",
        "related": ["architecture", "technical"],
    },
    "integrations/metadata": {
        "title": "Metadata Providers",
        "desc": "External sources for media metadata",
        "related": ["features"],
    },
    "integrations/metadata/video": {
        "title": "Video Metadata Providers",
        "desc": "TMDB, TVDB, OMDB, Fanart.tv",
        "related": ["features/video"],
    },
    "integrations/metadata/music": {
        "title": "Music Metadata Providers",
        "desc": "MusicBrainz, Last.fm, Spotify, Discogs",
        "related": [],
    },
    "integrations/metadata/books": {
        "title": "Book Metadata Providers",
        "desc": "Open Library, Google Books, Goodreads",
        "related": [],
    },
    "integrations/metadata/comics": {
        "title": "Comic Metadata Providers",
        "desc": "ComicVine, Marvel API, GCD",
        "related": ["features/comics"],
    },
    "integrations/metadata/adult": {
        "title": "Adult Metadata Providers",
        "desc": "StashDB, TPDB, FreeOnes",
        "related": ["features/adult"],
    },
    "integrations/scrobbling": {
        "title": "Scrobbling Services",
        "desc": "Activity tracking and sync services",
        "related": ["features/shared"],
    },
    "integrations/auth": {
        "title": "Authentication Providers",
        "desc": "OIDC/SSO providers",
        "related": ["architecture", "features/shared"],
    },
    "integrations/servarr": {
        "title": "Servarr Stack",
        "desc": "Media management automation (Radarr, Sonarr, etc.)",
        "related": ["features/video", "operations"],
    },
    "integrations/anime": {
        "title": "Anime Integration",
        "desc": "Anime-specific metadata and tracking",
        "related": ["features/video", "integrations/scrobbling"],
    },
    "integrations/audiobook": {
        "title": "Audiobook & Podcast Integration",
        "desc": "Native audiobook and podcast management",
        "related": ["features/podcasts"],
    },
    "integrations/wiki": {
        "title": "Wiki Integration",
        "desc": "Supplementary information sources",
        "related": ["features/shared"],
    },
    "integrations/wiki/adult": {
        "title": "Adult Wiki Integration",
        "desc": "Adult performer information sources",
        "related": ["features/adult"],
    },
    "integrations/external": {
        "title": "External Services",
        "desc": "Third-party integrations",
        "related": [],
    },
    "integrations/external/adult": {
        "title": "Adult External Services",
        "desc": "Social media integration for adult content",
        "related": ["features/adult"],
    },
    "integrations/infrastructure": {
        "title": "Infrastructure Components",
        "desc": "Core stack: PostgreSQL, Dragonfly, Typesense, River",
        "related": ["architecture", "operations"],
    },
    "integrations/transcoding": {
        "title": "Transcoding Services",
        "desc": "External transcoding (Blackbeard)",
        "related": ["features/playback"],
    },
    "integrations/livetv": {
        "title": "Live TV Integration",
        "desc": "PVR backend integration",
        "related": ["features/livetv"],
    },
    "integrations/casting": {
        "title": "Casting Protocols",
        "desc": "Chromecast, DLNA device casting",
        "related": ["features/playback"],
    },
    "operations": {
        "title": "Operations Documentation",
        "desc": "Deployment, setup, and operational guides",
        "related": ["architecture", "integrations/infrastructure"],
    },
    "technical": {
        "title": "Technical Documentation",
        "desc": "API specs, frontend architecture, tech stack",
        "related": ["architecture", "features"],
    },
    "research": {
        "title": "Research Documentation",
        "desc": "Technology research and evaluations",
        "related": ["architecture"],
    },
    "services": {
        "title": "Services Documentation",
        "desc": "Internal service specifications",
        "related": ["architecture", "technical"],
    },
    "planning": {
        "title": "Planning Documentation",
        "desc": "Project planning and roadmap",
        "related": ["architecture"],
    },
}


def get_relative_path(from_dir: Path, to_path: Path) -> str:
    """Calculate relative path from one directory to another."""
    try:
        return str(to_path.relative_to(from_dir))
    except ValueError:
        parts_from = from_dir.parts
        parts_to = to_path.parts
        common = 0
        for a, b in zip(parts_from, parts_to, strict=False):
            if a == b:
                common += 1
            else:
                break
        ups = len(parts_from) - common
        return "/".join([".."] * ups + list(parts_to[common:]))


def extract_doc_info(doc_path: Path) -> dict[str, Any]:
    """Extract metadata from a design document."""
    content = doc_path.read_text(encoding="utf-8")

    title_match = re.search(r"^#\s+(.+)$", content, re.MULTILINE)
    title = title_match.group(1) if title_match else doc_path.stem.replace("_", " ")

    desc_match = re.search(r"^>\s*(.+)$", content, re.MULTILINE)
    desc = desc_match.group(1) if desc_match else ""

    has_schema = bool(re.search(r"```sql|CREATE TABLE", content, re.IGNORECASE))
    has_code = bool(re.search(r"```go|type \w+ struct", content))
    has_checklist = bool(re.search(r"## Implementation", content))
    checkbox_count = len(re.findall(r"- \[[ x]\]", content))
    external_links = len(re.findall(r"\[.+\]\(https?://", content))

    if has_schema or has_code:
        status = "✅ Designed"
    elif has_checklist and checkbox_count >= 3:
        status = "🟡 Partial"
    elif external_links >= 2 or desc:
        status = "🟡 Planned"
    else:
        status = "🔴 Draft"

    return {
        "path": doc_path,
        "title": title,
        "desc": desc,
        "status": status,
    }


def find_directories_with_docs() -> dict[str, list[Path]]:
    """Find all directories containing design docs."""
    dirs: dict[str, list[Path]] = defaultdict(list)

    for md_file in sorted(DESIGN_DIR.rglob("*.md")):
        if ".archive" in str(md_file):
            continue
        if ".analysis" in str(md_file):
            continue
        if md_file.name in ["INDEX.md", "00_INDEX.md"]:
            continue
        if md_file.parent == DESIGN_DIR:
            continue

        rel_dir = str(md_file.parent.relative_to(DESIGN_DIR))
        dirs[rel_dir].append(md_file)

    return dict(dirs)


def generate_index(dir_path: str, docs: list[Path]) -> str:
    """Generate INDEX.md content for a directory."""
    abs_dir = DESIGN_DIR / dir_path

    meta = CATEGORY_META.get(
        dir_path,
        {
            "title": dir_path.replace("/", " - ").replace("_", " ").title(),
            "desc": "",
            "related": [],
        },
    )

    sot_rel = get_relative_path(abs_dir, SOT_FILE)
    parent_rel = ".." if "/" in dir_path else "../"
    sources_rel = get_relative_path(abs_dir, SOURCES_DIR / "SOURCES.md")

    lines = [
        f"# {meta['title']}",
        "",
        f"← Back to [Design Docs]({parent_rel})",
        "",
    ]

    if meta.get("desc"):
        lines.extend([f"> {meta['desc']}", ""])

    lines.extend(
        [
            f"**Source of Truth**: [{sot_rel.split('/')[-1]}]({sot_rel})",
            "",
            "---",
            "",
        ],
    )

    # Subdirectories
    subdirs = []
    for item in sorted(abs_dir.iterdir()):
        if item.is_dir() and not item.name.startswith("."):
            sub_index = item / "INDEX.md"
            if sub_index.exists() or any(item.glob("*.md")):
                subdirs.append(item)

    if subdirs:
        lines.extend(["## Subdirectories", ""])
        for subdir in subdirs:
            sub_rel = f"{subdir.name}/INDEX.md"
            sub_key = f"{dir_path}/{subdir.name}"
            sub_meta = CATEGORY_META.get(sub_key, {})
            sub_title = sub_meta.get("title", subdir.name.replace("_", " ").title())
            sub_desc = sub_meta.get("desc", "")
            if sub_desc:
                lines.append(f"- [{sub_title}]({sub_rel}) - {sub_desc}")
            else:
                lines.append(f"- [{sub_title}]({sub_rel})")
        lines.extend(["", "---", ""])

    # Documents in this directory
    direct_docs = [d for d in docs if d.parent == abs_dir]
    if direct_docs:
        lines.extend(["## Documents", ""])
        lines.append("| Document | Description | Status |")
        lines.append("|----------|-------------|--------|")

        for doc in sorted(direct_docs, key=lambda x: x.name):
            info = extract_doc_info(doc)
            rel_path = doc.name
            title = info["title"]
            desc = info.get("desc", "")
            if len(desc) > 60:
                desc = desc[:60] + "..."
            status = info["status"]
            lines.append(f"| [{title}]({rel_path}) | {desc} | {status} |")

        lines.extend(["", "---", ""])

    # Sources reference (simplified)
    lines.extend(
        [
            f"**Sources**: [External Sources Index]({sources_rel})",
            "",
        ],
    )

    # Related documentation
    if meta.get("related"):
        lines.extend(["## Related", ""])
        for rel in meta["related"]:
            rel_meta = CATEGORY_META.get(rel, {})
            rel_title = rel_meta.get("title", rel.replace("/", " - ").title())
            rel_path = get_relative_path(abs_dir, DESIGN_DIR / rel)
            lines.append(f"- [{rel_title}]({rel_path}/)")
        lines.append("")

    lines.extend(
        [
            "---",
            "",
            "## Status Legend",
            "",
            f"> See [{sot_rel.split('/')[-1]}]({sot_rel}#status-system) "
            "for full status definitions",
            "",
            "Quick reference: ✅ Complete | 🟡 Partial | 🔴 Not Started | ⚪ N/A",
            "",
        ],
    )

    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(description="Generate design doc indexes")
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Write files (default: dry-run)",
    )
    args = parser.parse_args()

    dry_run = not args.apply

    if dry_run:
        print("=== DRY RUN (use --apply to write) ===\n")

    print("Finding design doc directories...")
    dirs = find_directories_with_docs()
    print(f"  Found {len(dirs)} directories with docs")

    updated = 0
    unchanged = 0

    for dir_path, docs in sorted(dirs.items()):
        abs_dir = DESIGN_DIR / dir_path
        index_path = abs_dir / "INDEX.md"

        new_content = generate_index(dir_path, docs)

        if index_path.exists():
            old_content = index_path.read_text(encoding="utf-8")
            if old_content == new_content:
                unchanged += 1
                continue

        if dry_run:
            print(f"  Would update: {dir_path}/INDEX.md")
        else:
            index_path.write_text(new_content, encoding="utf-8")
            print(f"  Updated: {dir_path}/INDEX.md")
        updated += 1

    print("\n=== SUMMARY ===")
    action = "Updated" if not dry_run else "Would update"
    print(f"{action}: {updated}")
    print(f"Unchanged: {unchanged}")

    if dry_run and updated > 0:
        print("\n=== DRY RUN complete. Use --apply to write. ===")

    return 0


if __name__ == "__main__":
    sys.exit(main())
