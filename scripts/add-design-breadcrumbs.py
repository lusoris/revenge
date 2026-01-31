#!/usr/bin/env python3
"""
Add internal design doc breadcrumbs to all design documents.

Adds/updates DESIGN-BREADCRUMBS section in each doc showing:
- Parent category with link to index
- Related documents in same category
- Topic-related documents from other categories
- Links to cross-reference indexes

Does NOT modify: 00_SOURCE_OF_TRUTH.md

Usage:
    python scripts/add-design-breadcrumbs.py              # Dry run
    python scripts/add-design-breadcrumbs.py --update     # Write files
"""

import argparse
import re
from collections import defaultdict
from pathlib import Path

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOT_FILE = "00_SOURCE_OF_TRUTH.md"

# Files to never modify
SKIP_FILES = {
    SOT_FILE,  # Never modify SOT
    "01_DESIGN_DOC_TEMPLATE.md",
    "02_QUESTIONS_TO_DISCUSS.md",
    "03_DESIGN_DOCS_STATUS.md",
    "DESIGN_INDEX.md",
}

# Topic keywords for finding related docs
TOPIC_KEYWORDS = {
    "authentication": ["auth", "oidc", "oauth", "jwt", "sso", "login", "session", "rbac", "casbin"],
    "metadata": ["metadata", "provider", "tmdb", "tvdb", "musicbrainz", "stashdb", "omdb"],
    "playback": ["player", "stream", "hls", "transcode", "playback", "media", "trickplay", "skip"],
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

BREADCRUMB_START = "<!-- DESIGN-BREADCRUMBS-START -->"
BREADCRUMB_END = "<!-- DESIGN-BREADCRUMBS-END -->"


def find_design_docs() -> list[Path]:
    """Find all design documents."""
    docs = []
    for md_file in sorted(DESIGN_DIR.rglob("*.md")):
        if ".archive" in str(md_file):
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
    topic_docs = defaultdict(list)
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


def get_related_docs(doc_path: Path, all_docs: list[Path], topic_index: dict) -> list[Path]:
    """Get related documents from other categories based on shared topics."""
    content = doc_path.read_text(encoding="utf-8")
    doc_topics = get_doc_topics(content)
    doc_dir = doc_path.parent

    related = set()
    for topic in doc_topics:
        for other_doc in topic_index.get(topic, []):
            # Only include docs from OTHER directories
            if other_doc.parent != doc_dir and other_doc != doc_path:
                related.add(other_doc)

    return sorted(related, key=lambda x: str(x))[:5]  # Limit to 5


def get_relative_path(from_path: Path, to_path: Path) -> str:
    """Calculate relative path between two files."""
    from_dir = from_path.parent
    try:
        return str(to_path.relative_to(from_dir))
    except ValueError:
        # Need to go up
        from_parts = from_dir.parts
        to_parts = to_path.parts
        common = 0
        for a, b in zip(from_parts, to_parts):
            if a == b:
                common += 1
            else:
                break
        ups = len(from_parts) - common
        rel_parts = [".."] * ups + list(to_parts[common:])
        return "/".join(rel_parts)


def get_doc_title(doc_path: Path) -> str:
    """Extract title from document."""
    content = doc_path.read_text(encoding="utf-8")
    match = re.search(r'^#\s+(.+)$', content, re.MULTILINE)
    return match.group(1) if match else doc_path.stem.replace("_", " ")


def generate_breadcrumb_section(
    doc_path: Path,
    category_docs: list[Path],
    related_docs: list[Path],
) -> str:
    """Generate the breadcrumb section content."""
    lines = [
        BREADCRUMB_START,
        "",
        "## Related Design Docs",
        "",
        "> Auto-generated cross-references to related design documentation",
        "",
    ]

    # Parent category link
    parent_index = doc_path.parent / "INDEX.md"
    if parent_index.exists():
        cat_name = doc_path.parent.name.replace("_", " ").title()
        lines.extend([
            f"**Category**: [{cat_name}](INDEX.md)",
            "",
        ])

    # Same category docs
    if category_docs:
        lines.extend(["### In This Section", ""])
        for other_doc in sorted(category_docs, key=lambda x: x.name)[:8]:
            title = get_doc_title(other_doc)
            rel_path = other_doc.name
            lines.append(f"- [{title}]({rel_path})")
        lines.append("")

    # Related docs from other categories
    if related_docs:
        lines.extend(["### Related Topics", ""])
        for other_doc in related_docs:
            title = get_doc_title(other_doc)
            rel_path = get_relative_path(doc_path, other_doc)
            # Get category name for context
            cat = other_doc.parent.name.replace("_", " ").title()
            lines.append(f"- [{title}]({rel_path}) _{cat}_")
        lines.append("")

    # Cross-reference indexes
    design_index_path = get_relative_path(doc_path, DESIGN_DIR / "DESIGN_INDEX.md")
    sot_path = get_relative_path(doc_path, DESIGN_DIR / SOT_FILE)

    lines.extend([
        "### Indexes",
        "",
        f"- [Design Index]({design_index_path}) - All design docs by category/topic",
        f"- [Source of Truth]({sot_path}) - Package versions and status",
        "",
        BREADCRUMB_END,
    ])

    return "\n".join(lines)


def update_doc_breadcrumbs(
    doc_path: Path,
    category_docs: list[Path],
    related_docs: list[Path],
    dry_run: bool = True,
) -> bool:
    """Update breadcrumbs in a document. Returns True if changed."""
    content = doc_path.read_text(encoding="utf-8")

    # Generate new breadcrumb section
    new_section = generate_breadcrumb_section(doc_path, category_docs, related_docs)

    # Check if section already exists
    pattern = re.compile(
        rf'{re.escape(BREADCRUMB_START)}.*?{re.escape(BREADCRUMB_END)}',
        re.DOTALL
    )

    if pattern.search(content):
        # Replace existing section
        new_content = pattern.sub(new_section, content)
    else:
        # Add section before first --- or at end
        # Try to add before "## Related" or at end
        if "\n## Related" in content:
            new_content = content.replace("\n## Related", f"\n{new_section}\n\n## Related")
        elif "\n---\n" in content:
            # Add before last ---
            parts = content.rsplit("\n---\n", 1)
            if len(parts) == 2:
                new_content = parts[0] + f"\n{new_section}\n\n---\n" + parts[1]
            else:
                new_content = content + f"\n\n{new_section}\n"
        else:
            new_content = content + f"\n\n{new_section}\n"

    if new_content == content:
        return False

    if not dry_run:
        doc_path.write_text(new_content, encoding="utf-8")

    return True


def main():
    parser = argparse.ArgumentParser(description="Add design doc breadcrumbs")
    parser.add_argument(
        "--update", "-u", action="store_true", help="Write files (default: dry run)"
    )
    args = parser.parse_args()

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

        changed = update_doc_breadcrumbs(
            doc_path,
            category_docs,
            related_docs,
            dry_run=not args.update
        )

        rel_path = doc_path.relative_to(DESIGN_DIR)
        if changed:
            action = "Updated" if args.update else "Would update"
            print(f"  {action}: {rel_path}")
            updated += 1
        else:
            unchanged += 1

    print(f"\n=== SUMMARY ===")
    print(f"{'Updated' if args.update else 'Would update'}: {updated}")
    print(f"Unchanged: {unchanged}")
    print(f"Skipped: {SOT_FILE} (protected)")

    if not args.update and updated > 0:
        print("\nRun with --update to write changes")


if __name__ == "__main__":
    main()
