#!/usr/bin/env python3
"""Add minimal design breadcrumbs to design documents.

Uses simplified format: <!-- DESIGN: category, related1, related2 -->
Links related documents by topic detection.

Usage:
    python scripts/doc-pipeline/02-breadcrumbs.py           # Dry run
    python scripts/doc-pipeline/02-breadcrumbs.py --apply   # Write files
"""

from __future__ import annotations

import argparse
import re
import sys
from collections import defaultdict
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOT_FILE = "00_SOURCE_OF_TRUTH.md"

SKIP_FILES = {
    SOT_FILE,
    "01_DESIGN_DOC_TEMPLATE.md",
    "02_QUESTIONS_TO_DISCUSS.md",
    "03_DESIGN_DOCS_STATUS.md",
    "DESIGN_INDEX.md",
    "NAVIGATION.md",
}

# Old verbose breadcrumb markers
OLD_BREADCRUMB_START = "<!-- DESIGN-BREADCRUMBS-START -->"
OLD_BREADCRUMB_END = "<!-- DESIGN-BREADCRUMBS-END -->"

# New minimal format
NEW_BREADCRUMB_PATTERN = re.compile(r"<!-- DESIGN: ([^>]+) -->")

TOPIC_KEYWORDS = {
    "authentication": [
        "auth",
        "oidc",
        "oauth",
        "jwt",
        "sso",
        "login",
        "session",
        "rbac",
        "casbin",
    ],
    "metadata": [
        "metadata",
        "provider",
        "tmdb",
        "tvdb",
        "musicbrainz",
        "stashdb",
        "omdb",
    ],
    "playback": [
        "player",
        "stream",
        "hls",
        "transcode",
        "playback",
        "media",
        "trickplay",
        "skip",
    ],
    "search": ["search", "typesense", "index", "query", "filter"],
    "database": ["postgresql", "database", "schema", "migration", "sqlc"],
    "caching": ["cache", "dragonfly", "redis", "otter", "sturdyc"],
    "jobs": ["job", "queue", "river", "worker", "background"],
    "api": ["api", "endpoint", "rest", "graphql", "openapi", "ogen"],
    "frontend": ["svelte", "frontend", "ui", "component", "client"],
    "scrobbling": ["scrobble", "trakt", "lastfm", "listenbrainz", "sync", "simkl"],
    "adult": ["adult", "qar", "stash", "whisparr", "performer", "scene", "tpdb"],
    "video": ["movie", "tvshow", "episode", "series", "radarr", "sonarr"],
    "music": ["music", "artist", "album", "track", "lidarr", "spotify", "discogs"],
    "books": ["book", "audiobook", "readarr", "chapter", "audible", "goodreads"],
    "photos": ["photo", "image", "gallery", "exif", "immich"],
    "livetv": ["livetv", "dvr", "epg", "channel", "tvheadend", "ersatztv"],
    "comics": ["comic", "comicvine", "marvel", "manga"],
    "anime": ["anime", "anilist", "myanimelist", "kitsu"],
    "wiki": ["wiki", "fandom", "wikipedia", "tvtropes"],
    "notifications": ["notification", "alert", "email", "push"],
    "users": ["user", "profile", "settings", "preferences", "avatar"],
}


def find_design_docs() -> list[Path]:
    """Find all design documents."""
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
        docs.append(md_file)
    return docs


def get_doc_topics(content: str) -> set[str]:
    """Detect topics in document content."""
    content_lower = content.lower()
    topics = set()
    for topic, keywords in TOPIC_KEYWORDS.items():
        if any(kw in content_lower for kw in keywords):
            topics.add(topic)
    return topics


def build_topic_index(docs: list[Path]) -> dict[str, list[Path]]:
    """Build topic to documents index."""
    topic_docs: dict[str, list[Path]] = defaultdict(list)
    for doc_path in docs:
        content = doc_path.read_text(encoding="utf-8")
        topics = get_doc_topics(content)
        for topic in topics:
            topic_docs[topic].append(doc_path)
    return dict(topic_docs)


def get_category_docs(doc_path: Path, all_docs: list[Path]) -> list[Path]:
    """Get other documents in the same category/directory."""
    doc_dir = doc_path.parent
    return [d for d in all_docs if d.parent == doc_dir and d != doc_path]


def get_related_docs(
    doc_path: Path,
    all_docs: list[Path],
    topic_index: dict[str, list[Path]],
) -> list[Path]:
    """Get related documents from other categories based on shared topics."""
    content = doc_path.read_text(encoding="utf-8")
    doc_topics = get_doc_topics(content)
    doc_dir = doc_path.parent

    related = set()
    for topic in doc_topics:
        for other_doc in topic_index.get(topic, []):
            if other_doc.parent != doc_dir and other_doc != doc_path:
                related.add(other_doc)

    return sorted(related, key=lambda x: str(x))[:5]


def get_category_name(doc_path: Path) -> str:
    """Get the category name for a document."""
    try:
        rel = doc_path.relative_to(DESIGN_DIR)
        parts = rel.parts[:-1]  # Exclude filename
        if parts:
            return "/".join(parts)
    except ValueError:
        pass
    return "root"


def remove_old_breadcrumbs(content: str) -> str:
    """Remove old verbose breadcrumb sections."""
    pattern = re.compile(
        rf"{re.escape(OLD_BREADCRUMB_START)}.*?{re.escape(OLD_BREADCRUMB_END)}\n*",
        re.DOTALL,
    )
    content = pattern.sub("", content)

    while "\n\n\n\n" in content:
        content = content.replace("\n\n\n\n", "\n\n\n")

    return content


def generate_minimal_breadcrumb(
    category: str,
    related_docs: list[Path],
) -> str:
    """Generate minimal breadcrumb comment."""
    parts = [category]
    for doc in related_docs[:3]:
        parts.append(doc.stem)
    return f"<!-- DESIGN: {', '.join(parts)} -->"


def update_document(
    doc_path: Path,
    category_docs: list[Path],
    related_docs: list[Path],
    *,
    dry_run: bool = True,
) -> bool:
    """Update breadcrumbs in a document. Returns True if changed."""
    content = doc_path.read_text(encoding="utf-8")

    # Get category
    category = get_category_name(doc_path)

    # Remove old verbose breadcrumbs
    new_content = remove_old_breadcrumbs(content)

    # Remove any existing minimal breadcrumb
    new_content = NEW_BREADCRUMB_PATTERN.sub("", new_content)

    # Generate new minimal breadcrumb
    breadcrumb = generate_minimal_breadcrumb(category, related_docs)

    # Find title and insert after it
    lines = new_content.split("\n")
    title_idx = None
    for i, line in enumerate(lines):
        if line.startswith("# "):
            title_idx = i
            break

    if title_idx is not None:
        # Check if there's already a source breadcrumb after title
        insert_idx = title_idx + 1
        for i in range(title_idx + 1, min(title_idx + 4, len(lines))):
            if lines[i].startswith("<!-- SOURCES:"):
                insert_idx = i + 1
                break

        lines.insert(insert_idx, breadcrumb)
        if insert_idx == title_idx + 1:
            lines.insert(insert_idx, "")

    new_content = "\n".join(lines)

    while "\n\n\n\n" in new_content:
        new_content = new_content.replace("\n\n\n\n", "\n\n\n")

    changed = new_content != content

    if changed and not dry_run:
        doc_path.write_text(new_content, encoding="utf-8")

    return changed


def main() -> int:
    parser = argparse.ArgumentParser(description="Add design doc breadcrumbs")
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Write files (default: dry-run)",
    )
    args = parser.parse_args()

    dry_run = not args.apply

    if dry_run:
        print("=== DRY RUN (use --apply to write) ===\n")

    print("Finding design documents...")
    docs = find_design_docs()
    print(f"  Found {len(docs)} documents (excluding SOT)")

    print("Building topic index...")
    topic_index = build_topic_index(docs)
    print(f"  {len(topic_index)} topics detected")

    print("Processing documents...")
    updated = 0
    unchanged = 0

    for doc_path in docs:
        category_docs = get_category_docs(doc_path, docs)
        related_docs = get_related_docs(doc_path, docs, topic_index)

        changed = update_document(
            doc_path, category_docs, related_docs, dry_run=dry_run,
        )

        rel_path = doc_path.relative_to(DESIGN_DIR)
        if changed:
            action = "Updated" if not dry_run else "Would update"
            print(f"  {action}: {rel_path}")
            updated += 1
        else:
            unchanged += 1

    print("\n=== SUMMARY ===")
    action = "Updated" if not dry_run else "Would update"
    print(f"{action}: {updated}")
    print(f"Unchanged: {unchanged}")
    print(f"Skipped: {SOT_FILE} (protected)")

    if dry_run and updated > 0:
        print("\n=== DRY RUN complete. Use --apply to write. ===")

    return 0


if __name__ == "__main__":
    sys.exit(main())
