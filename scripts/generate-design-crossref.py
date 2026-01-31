#!/usr/bin/env python3
"""
Generate cross-reference index for design documents.

Analyzes internal links and topic relationships between design docs
to create a comprehensive cross-reference index.

Output: docs/dev/design/DESIGN_INDEX.md

Usage:
    python scripts/generate-design-crossref.py
"""

import re
from collections import defaultdict
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
OUTPUT_FILE = DESIGN_DIR / "DESIGN_INDEX.md"

# Files to skip
SKIP_FILES = {
    "00_SOURCE_OF_TRUTH.md",
    "01_DESIGN_DOC_TEMPLATE.md",
    "02_QUESTIONS_TO_DISCUSS.md",
    "03_DESIGN_DOCS_STATUS.md",
    "DESIGN_INDEX.md",
}

# Topic keywords for categorization
TOPIC_KEYWORDS = {
    "authentication": ["auth", "oidc", "oauth", "jwt", "sso", "login", "session"],
    "metadata": ["metadata", "provider", "tmdb", "tvdb", "musicbrainz", "stashdb"],
    "playback": ["player", "stream", "hls", "transcode", "playback", "media"],
    "search": ["search", "typesense", "index", "query", "filter"],
    "database": ["postgresql", "database", "schema", "migration", "sqlc"],
    "caching": ["cache", "dragonfly", "redis", "otter", "sturdyc"],
    "jobs": ["job", "queue", "river", "worker", "background"],
    "api": ["api", "endpoint", "rest", "graphql", "openapi"],
    "frontend": ["svelte", "frontend", "ui", "component", "client"],
    "scrobbling": ["scrobble", "trakt", "lastfm", "listenbrainz", "sync"],
    "adult": ["adult", "qar", "stash", "whisparr", "performer", "scene"],
    "video": ["movie", "tvshow", "episode", "series", "radarr", "sonarr"],
    "music": ["music", "artist", "album", "track", "lidarr"],
    "books": ["book", "audiobook", "readarr", "chapter"],
    "photos": ["photo", "image", "gallery", "exif"],
    "livetv": ["livetv", "dvr", "epg", "channel", "tvheadend"],
}


def find_design_docs() -> list[Path]:
    """Find all design documents."""
    docs = []
    for md_file in sorted(DESIGN_DIR.rglob("*.md")):
        if ".archive" in str(md_file):
            continue
        if md_file.name in SKIP_FILES:
            continue
        if md_file.name.startswith("INDEX"):
            continue
        if md_file.name.startswith("00_INDEX"):
            continue
        docs.append(md_file)
    return docs


def get_doc_info(doc_path: Path) -> dict:
    """Extract information from a design document."""
    content = doc_path.read_text(encoding="utf-8")
    rel_path = doc_path.relative_to(DESIGN_DIR)

    # Get title
    title_match = re.search(r"^#\s+(.+)$", content, re.MULTILINE)
    title = title_match.group(1) if title_match else doc_path.stem.replace("_", " ")

    # Get description
    desc_match = re.search(r"^>\s*(.+)$", content, re.MULTILINE)
    desc = desc_match.group(1) if desc_match else ""

    # Find internal links to other design docs
    internal_links = set()
    # Match markdown links like [text](path.md) or [text](../path.md)
    link_pattern = r"\[([^\]]+)\]\(([^)]+\.md)\)"
    for match in re.finditer(link_pattern, content):
        link_path = match.group(2)
        # Skip external links and source links
        if link_path.startswith("http") or "sources/" in link_path:
            continue
        # Resolve relative path
        try:
            if link_path.startswith("../"):
                resolved = (doc_path.parent / link_path).resolve()
            else:
                resolved = (doc_path.parent / link_path).resolve()
            if (
                resolved.exists() and DESIGN_DIR in resolved.parents
            ) or resolved.parent == DESIGN_DIR:
                internal_links.add(str(resolved.relative_to(DESIGN_DIR)))
        except (ValueError, OSError):
            pass

    # Detect topics
    content_lower = content.lower()
    topics = []
    for topic, keywords in TOPIC_KEYWORDS.items():
        if any(kw in content_lower for kw in keywords):
            topics.append(topic)

    # Get category from path
    parts = rel_path.parts
    category = "/".join(parts[:-1]) if len(parts) > 1 else "root"

    return {
        "path": str(rel_path),
        "title": title,
        "desc": desc[:100] + "..." if len(desc) > 100 else desc,
        "category": category,
        "topics": topics,
        "links_to": sorted(internal_links),
        "linked_from": [],  # Will be populated later
    }


def build_crossref(docs: dict[str, dict]) -> dict[str, dict]:
    """Build bidirectional cross-references."""
    # Add linked_from references
    for doc_path, info in docs.items():
        for linked_path in info["links_to"]:
            if linked_path in docs:
                docs[linked_path]["linked_from"].append(doc_path)

    # Sort linked_from
    for info in docs.values():
        info["linked_from"] = sorted(set(info["linked_from"]))

    return docs


def build_topic_index(docs: dict[str, dict]) -> dict[str, list[str]]:
    """Build topic to documents index."""
    topic_docs = defaultdict(list)
    for doc_path, info in docs.items():
        for topic in info["topics"]:
            topic_docs[topic].append(doc_path)
    return dict(topic_docs)


def build_category_index(docs: dict[str, dict]) -> dict[str, list[str]]:
    """Build category to documents index."""
    cat_docs = defaultdict(list)
    for doc_path, info in docs.items():
        cat_docs[info["category"]].append(doc_path)
    return dict(cat_docs)


def generate_index(docs: dict[str, dict], topics: dict, categories: dict) -> str:
    """Generate the DESIGN_INDEX.md content."""
    lines = [
        "# Design Documentation Index",
        "",
        "> Auto-generated cross-reference index for all design documents",
        "",
        "**Source of Truth**: [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md)",
        "",
        "---",
        "",
        "## Quick Stats",
        "",
        f"- **Total Documents**: {len(docs)}",
        f"- **Categories**: {len(categories)}",
        f"- **Topics**: {len(topics)}",
        "",
        "---",
        "",
        "## By Category",
        "",
    ]

    # Category index
    for category in sorted(categories.keys()):
        doc_paths = categories[category]
        cat_title = category.replace("/", " → ").replace("_", " ").title()
        lines.append(f"### {cat_title}")
        lines.append("")
        lines.append("| Document | Topics | Links |")
        lines.append("|----------|--------|-------|")

        for doc_path in sorted(doc_paths):
            info = docs[doc_path]
            title = info["title"]
            topic_str = ", ".join(info["topics"][:3]) if info["topics"] else "-"
            link_count = len(info["links_to"]) + len(info["linked_from"])
            lines.append(f"| [{title}]({doc_path}) | {topic_str} | {link_count} |")

        lines.append("")

    # Topic index
    lines.extend(
        [
            "---",
            "",
            "## By Topic",
            "",
        ]
    )

    for topic in sorted(topics.keys()):
        doc_paths = topics[topic]
        lines.append(f"### {topic.title()}")
        lines.append("")
        for doc_path in sorted(doc_paths):
            info = docs[doc_path]
            lines.append(f"- [{info['title']}]({doc_path})")
        lines.append("")

    # Cross-reference graph (documents with most connections)
    lines.extend(
        [
            "---",
            "",
            "## Most Connected Documents",
            "",
            "> Documents with the most internal cross-references",
            "",
            "| Document | Links To | Linked From | Total |",
            "|----------|----------|-------------|-------|",
        ]
    )

    # Sort by total connections
    sorted_docs = sorted(
        docs.items(),
        key=lambda x: len(x[1]["links_to"]) + len(x[1]["linked_from"]),
        reverse=True,
    )[:20]

    for doc_path, info in sorted_docs:
        links_to = len(info["links_to"])
        linked_from = len(info["linked_from"])
        total = links_to + linked_from
        if total > 0:
            lines.append(
                f"| [{info['title']}]({doc_path}) | {links_to} | {linked_from} | {total} |"
            )

    lines.append("")

    # Orphan documents (no connections)
    orphans = [
        (path, info)
        for path, info in docs.items()
        if not info["links_to"] and not info["linked_from"]
    ]

    if orphans:
        lines.extend(
            [
                "---",
                "",
                "## Orphan Documents",
                "",
                "> Documents with no internal cross-references (may need linking)",
                "",
            ]
        )
        for doc_path, info in sorted(orphans, key=lambda x: x[0]):
            lines.append(f"- [{info['title']}]({doc_path})")
        lines.append("")

    # Footer
    lines.extend(
        [
            "---",
            "",
            "*Generated by `scripts/generate-design-crossref.py`*",
            "",
        ]
    )

    return "\n".join(lines)


def main():
    print("Finding design documents...")
    doc_files = find_design_docs()
    print(f"  Found {len(doc_files)} documents")

    print("Analyzing documents...")
    docs = {}
    for doc_path in doc_files:
        info = get_doc_info(doc_path)
        docs[info["path"]] = info

    print("Building cross-references...")
    docs = build_crossref(docs)
    topics = build_topic_index(docs)
    categories = build_category_index(docs)

    print(f"  {len(topics)} topics detected")
    print(f"  {len(categories)} categories")

    # Count connections
    total_links = sum(len(info["links_to"]) for info in docs.values())
    print(f"  {total_links} internal cross-references")

    print("Generating index...")
    content = generate_index(docs, topics, categories)

    OUTPUT_FILE.write_text(content, encoding="utf-8")
    print(f"  Wrote {OUTPUT_FILE}")

    print("\nDone!")


if __name__ == "__main__":
    main()
