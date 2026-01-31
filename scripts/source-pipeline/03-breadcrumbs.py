#!/usr/bin/env python3
"""Add minimal source breadcrumbs to design documents.

Uses the simplified format: <!-- SOURCES: id1, id2, id3 -->
This replaces the verbose SOURCE-BREADCRUMBS section with a compact one-liner.

Usage:
    python scripts/source-pipeline/03-breadcrumbs.py           # Dry run
    python scripts/source-pipeline/03-breadcrumbs.py --apply   # Write files
    python scripts/source-pipeline/03-breadcrumbs.py --file X.md --apply
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

import yaml


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_YAML = SOURCES_DIR / "SOURCES.yaml"

# Old verbose breadcrumb markers (to be replaced)
OLD_BREADCRUMB_START = "<!-- SOURCE-BREADCRUMBS-START -->"
OLD_BREADCRUMB_END = "<!-- SOURCE-BREADCRUMBS-END -->"

# New minimal format
NEW_BREADCRUMB_PATTERN = re.compile(r"<!-- SOURCES: ([^>]+) -->")

# Keyword to source mappings
KEYWORD_SOURCE_MAP = {
    "tmdb": ["tmdb-api"],
    "themoviedb": ["tmdb-api"],
    "thetvdb": ["thetvdb-api"],
    "musicbrainz": ["musicbrainz-api"],
    "last.fm": ["lastfm-api"],
    "lastfm": ["lastfm-api"],
    "listenbrainz": ["listenbrainz-api"],
    "trakt": ["trakt-api"],
    "simkl": ["simkl-api"],
    "anilist": ["anilist-schema"],
    "myanimelist": ["myanimelist-api"],
    "comicvine": ["comicvine-api"],
    "openlibrary": ["openlibrary-api"],
    "spotify": ["spotify-api"],
    "discogs": ["discogs-api"],
    "omdb": ["omdb-api"],
    "keycloak": ["keycloak"],
    "authelia": ["authelia"],
    "authentik": ["authentik"],
    "oidc": ["oidc-core"],
    "oauth": ["oauth2-rfc"],
    "jwt": ["jwt-rfc"],
    "radarr": ["radarr-openapi"],
    "sonarr": ["sonarr-openapi"],
    "lidarr": ["lidarr-openapi"],
    "whisparr": ["whisparr-openapi"],
    "readarr": ["readarr-openapi"],
    "postgresql": ["pgx", "postgresql-json", "postgresql-arrays"],
    "postgres": ["pgx", "postgresql-json"],
    "sqlc": ["sqlc", "sqlc-config"],
    "dragonfly": ["dragonfly"],
    "redis": ["rueidis", "rueidis-guide"],
    "typesense": ["typesense", "typesense-go"],
    "svelte": ["svelte5", "svelte-runes", "sveltekit"],
    "sveltekit": ["sveltekit"],
    "shadcn": ["shadcn-svelte"],
    "vidstack": ["vidstack"],
    "tanstack": ["tanstack-query"],
    "hls": ["hls-rfc", "gohlslib"],
    "webrtc": ["webrtc", "pion-webrtc"],
    "xmltv": ["xmltv", "xmltv-format"],
    "m3u": ["m3u8"],
    "ffmpeg": ["ffmpeg", "ffmpeg-codecs", "ffmpeg-formats", "go-astiav"],
    "blurhash": ["go-blurhash"],
    "vips": ["govips-guide"],
    "prometheus": ["prometheus", "prometheus-metrics"],
    "opentelemetry": ["opentelemetry"],
    "grafana": ["grafana", "grafana-dashboards"],
    "jaeger": ["jaeger", "jaeger-go"],
    "testcontainers": ["testcontainers"],
    "testify": ["testify"],
    "mockery": ["mockery", "mockery-guide"],
    "fx": ["fx", "fx-guide"],
    "koanf": ["koanf"],
    "ogen": ["ogen", "ogen-guide"],
    "river": ["river", "river-guide"],
    "casbin": ["casbin", "casbin-guide", "casbin-pgx"],
}


def load_sources() -> dict:
    """Load SOURCES.yaml."""
    with open(SOURCES_YAML, encoding="utf-8") as f:
        return yaml.safe_load(f)


def build_source_index(config: dict) -> dict:
    """Build URL/package -> source info mapping."""
    sources = config.get("sources", {})
    url_index = {}
    package_index = {}
    id_to_info = {}

    for category, category_sources in sources.items():
        for source in category_sources:
            source_id = source.get("id")
            url = source.get("url", "")
            name = source.get("name", source_id)
            output = source.get("output", "")

            info = {
                "id": source_id,
                "name": name,
                "url": url,
                "output": output,
                "category": category,
            }

            id_to_info[source_id] = info

            if url:
                normalized = (
                    url.replace("https://", "").replace("http://", "").rstrip("/")
                )
                url_index[normalized] = info

                if "pkg.go.dev/" in url:
                    pkg = url.split("pkg.go.dev/")[-1]
                    package_index[pkg] = info
                elif "github.com/" in url:
                    parts = url.replace("https://github.com/", "").split("/")
                    if len(parts) >= 2:
                        pkg = f"github.com/{parts[0]}/{parts[1]}"
                        package_index[pkg] = info

    return {"urls": url_index, "packages": package_index, "ids": id_to_info}


def extract_references(content: str, source_index: dict) -> list[str]:
    """Extract source IDs referenced in document content."""
    found = set()
    content_lower = content.lower()

    # Extract URLs
    url_pattern = r'https?://[^\s\)\]>"\']+'
    urls = set(re.findall(url_pattern, content))

    for url in urls:
        normalized = (
            url.replace("https://", "")
            .replace("http://", "")
            .rstrip("/")
            .rstrip(".,;:)")
        )
        for source_url, info in source_index["urls"].items():
            if source_url in normalized or normalized in source_url:
                found.add(info["id"])
                break

    # Extract package references
    github_pattern = r"github\.com/[\w\-]+/[\w\-]+(?:/[\w\-]+)*"
    packages = set(re.findall(github_pattern, content))

    golang_pattern = r"golang\.org/x/\w+"
    packages.update(re.findall(golang_pattern, content))

    uber_pattern = r"go\.uber\.org/\w+"
    packages.update(re.findall(uber_pattern, content))

    for pkg in packages:
        if pkg in source_index["packages"]:
            found.add(source_index["packages"][pkg]["id"])
        else:
            for source_pkg, info in source_index["packages"].items():
                if pkg in source_pkg or source_pkg in pkg:
                    found.add(info["id"])
                    break

    # Keyword-based matching
    for keyword, source_ids in KEYWORD_SOURCE_MAP.items():
        if keyword.lower() in content_lower:
            for source_id in source_ids:
                if source_id in source_index["ids"]:
                    found.add(source_id)

    return sorted(found)


def remove_old_breadcrumbs(content: str) -> str:
    """Remove old verbose breadcrumb sections."""
    # Remove old SOURCE-BREADCRUMBS section
    pattern = re.compile(
        rf"{re.escape(OLD_BREADCRUMB_START)}.*?{re.escape(OLD_BREADCRUMB_END)}\n*",
        re.DOTALL,
    )
    content = pattern.sub("", content)

    # Clean up any resulting multiple blank lines
    while "\n\n\n\n" in content:
        content = content.replace("\n\n\n\n", "\n\n\n")

    return content


def generate_minimal_breadcrumb(source_ids: list[str]) -> str:
    """Generate minimal breadcrumb comment."""
    if not source_ids:
        return ""
    return f"<!-- SOURCES: {', '.join(source_ids)} -->"


def update_document(
    doc_path: Path,
    source_index: dict,
    *,
    dry_run: bool = True,
) -> tuple[bool, list[str]]:
    """Update a document with minimal breadcrumbs. Returns (changed, source_ids)."""
    content = doc_path.read_text(encoding="utf-8")

    # Extract references
    source_ids = extract_references(content, source_index)

    # Remove old verbose breadcrumbs
    new_content = remove_old_breadcrumbs(content)

    # Remove any existing minimal breadcrumb
    new_content = NEW_BREADCRUMB_PATTERN.sub("", new_content)

    # Clean up whitespace at start
    lines = new_content.split("\n")

    # Find the title line (first # heading)
    title_idx = None
    for i, line in enumerate(lines):
        if line.startswith("# "):
            title_idx = i
            break

    if title_idx is not None and source_ids:
        breadcrumb = generate_minimal_breadcrumb(source_ids)
        # Insert after title
        lines.insert(title_idx + 1, "")
        lines.insert(title_idx + 2, breadcrumb)

    new_content = "\n".join(lines)

    # Clean up multiple blank lines
    while "\n\n\n\n" in new_content:
        new_content = new_content.replace("\n\n\n\n", "\n\n\n")

    changed = new_content != content

    if changed and not dry_run:
        doc_path.write_text(new_content, encoding="utf-8")

    return changed, source_ids


def find_design_docs(specific_file: str | None = None) -> list[Path]:
    """Find design documents to process."""
    if specific_file:
        path = Path(specific_file)
        if not path.is_absolute():
            path = DESIGN_DIR / specific_file
        return [path] if path.exists() else []

    docs = []
    for md_file in DESIGN_DIR.rglob("*.md"):
        if ".archive" in str(md_file):
            continue
        if md_file.name == "INDEX.md":
            continue
        docs.append(md_file)

    return sorted(docs)


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Add minimal source breadcrumbs to design docs"
    )
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Write changes (default: dry-run)",
    )
    parser.add_argument(
        "--file",
        "-f",
        help="Update specific file only",
    )
    args = parser.parse_args()

    dry_run = not args.apply

    if dry_run:
        print("=== DRY RUN (use --apply to write) ===\n")

    print("Loading sources configuration...")
    config = load_sources()
    source_index = build_source_index(config)

    print("Finding design documents...")
    docs = find_design_docs(args.file)
    print(f"  Found {len(docs)} documents")

    updated = 0
    unchanged = 0

    for doc in docs:
        rel_path = doc.relative_to(PROJECT_ROOT)
        changed, source_ids = update_document(doc, source_index, dry_run=dry_run)

        if changed:
            action = "Updated" if not dry_run else "Would update"
            refs = f"({len(source_ids)} refs)" if source_ids else "(no refs)"
            print(f"  {action}: {rel_path} {refs}")
            updated += 1
        else:
            unchanged += 1

    print("\n=== SUMMARY ===")
    action = "Updated" if not dry_run else "Would update"
    print(f"{action}: {updated}")
    print(f"Unchanged: {unchanged}")

    if dry_run and updated > 0:
        print("\n=== DRY RUN complete. Use --apply to write. ===")

    return 0


if __name__ == "__main__":
    sys.exit(main())
