#!/usr/bin/env python3
"""
Add source breadcrumbs to design documents.

Adds a "Sources & Cross-References" section to each design doc that links to:
1. The specific sources referenced in that document
2. The cross-reference indexes

Usage:
    python scripts/add-source-breadcrumbs.py              # Update all design docs
    python scripts/add-source-breadcrumbs.py --dry-run    # Preview changes
    python scripts/add-source-breadcrumbs.py --file X.md  # Update specific file
"""

import argparse
import re
from pathlib import Path

import yaml


# Project paths
SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_YAML = SOURCES_DIR / "SOURCES.yaml"

# Breadcrumb markers
BREADCRUMB_START = "<!-- SOURCE-BREADCRUMBS-START -->"
BREADCRUMB_END = "<!-- SOURCE-BREADCRUMBS-END -->"

# Keyword to source mappings for better cross-referencing
# Maps keywords/topics to source IDs (from SOURCES.yaml)
KEYWORD_SOURCE_MAP = {
    # API integrations
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
    # Auth providers
    "keycloak": ["keycloak"],
    "authelia": ["authelia"],
    "authentik": ["authentik"],
    "oidc": ["oidc-core"],
    "oauth": ["oauth2-rfc"],
    "jwt": ["jwt-rfc"],
    # Servarr stack
    "radarr": ["radarr-openapi"],
    "sonarr": ["sonarr-openapi"],
    "lidarr": ["lidarr-openapi"],
    "whisparr": ["whisparr-openapi"],
    "readarr": ["readarr-openapi"],
    # Database
    "postgresql": ["pgx", "postgresql-json", "postgresql-arrays"],
    "postgres": ["pgx", "postgresql-json"],
    "sqlc": ["sqlc", "sqlc-config"],
    # Caching
    "dragonfly": ["dragonfly"],
    "redis": ["rueidis", "rueidis-guide"],
    "typesense": ["typesense", "typesense-go"],
    # Frontend
    "svelte": ["svelte5", "svelte-runes", "sveltekit"],
    "sveltekit": ["sveltekit"],
    "shadcn": ["shadcn-svelte"],
    "vidstack": ["vidstack"],
    "tanstack": ["tanstack-query"],
    # Protocols
    "hls": ["hls-rfc", "gohlslib"],
    "webrtc": ["webrtc", "pion-webrtc"],
    "xmltv": ["xmltv", "xmltv-format"],
    "m3u": ["m3u8"],
    # Media processing
    "ffmpeg": ["ffmpeg", "ffmpeg-codecs", "ffmpeg-formats", "go-astiav"],
    "blurhash": ["go-blurhash"],
    "vips": ["govips-guide"],
    # Observability
    "prometheus": ["prometheus", "prometheus-metrics"],
    "opentelemetry": ["opentelemetry"],
    "grafana": ["grafana", "grafana-dashboards"],
    "jaeger": ["jaeger", "jaeger-go"],
    # Testing
    "testcontainers": ["testcontainers"],
    "testify": ["testify"],
    "mockery": ["mockery", "mockery-guide"],
    # Go tooling
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

            if url:
                # Normalize URL
                normalized = (
                    url.replace("https://", "").replace("http://", "").rstrip("/")
                )
                url_index[normalized] = info

                # Extract package path
                if "pkg.go.dev/" in url:
                    pkg = url.split("pkg.go.dev/")[-1]
                    package_index[pkg] = info
                elif "github.com/" in url:
                    parts = url.replace("https://github.com/", "").split("/")
                    if len(parts) >= 2:
                        pkg = f"github.com/{parts[0]}/{parts[1]}"
                        package_index[pkg] = info

    return {"urls": url_index, "packages": package_index}


def extract_references(content: str, source_index: dict) -> list[dict]:
    """Extract source references from document content."""
    found = {}
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
                found[info["id"]] = info
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
            info = source_index["packages"][pkg]
            found[info["id"]] = info
        else:
            # Partial match
            for source_pkg, info in source_index["packages"].items():
                if pkg in source_pkg or source_pkg in pkg:
                    found[info["id"]] = info
                    break

    # Keyword-based matching (new)
    for keyword, source_ids in KEYWORD_SOURCE_MAP.items():
        if keyword.lower() in content_lower:
            for source_id in source_ids:
                # Look up source by ID in the index
                for info in source_index["urls"].values():
                    if info["id"] == source_id:
                        found[info["id"]] = info
                        break
                for info in source_index["packages"].values():
                    if info["id"] == source_id:
                        found[info["id"]] = info
                        break

    return list(found.values())


def calculate_relative_path(from_file: Path, to_file: Path) -> str:
    """Calculate relative path from one file to another."""
    from_dir = from_file.parent
    try:
        return str(to_file.relative_to(from_dir))
    except ValueError:
        # Need to go up directories
        rel = Path()
        current = from_dir
        while True:
            try:
                target_rel = to_file.relative_to(current)
                return str(rel / target_rel)
            except ValueError:
                rel = rel / ".."
                current = current.parent
                if current == current.parent:
                    break
        return str(to_file)


def generate_breadcrumb_section(doc_path: Path, references: list[dict]) -> str:
    """Generate the breadcrumb section content."""
    lines = [
        BREADCRUMB_START,
        "",
        "## Sources & Cross-References",
        "",
        "> Auto-generated section linking to external documentation sources",
        "",
    ]

    # Calculate relative paths to indexes
    sources_index_rel = calculate_relative_path(
        doc_path,
        SOURCES_DIR / "SOURCES_INDEX.md",
    )
    crossref_rel = calculate_relative_path(doc_path, SOURCES_DIR / "DESIGN_CROSSREF.md")

    lines.extend(
        [
            "### Cross-Reference Indexes",
            "",
            f"- [All Sources Index]({sources_index_rel}) - Complete list of external documentation",
            f"- [Design ↔ Sources Map]({crossref_rel}) - Which docs reference which sources",
            "",
        ],
    )

    if references:
        lines.extend(
            [
                "### Referenced Sources",
                "",
                "| Source | Documentation |",
                "|--------|---------------|",
            ],
        )

        for ref in sorted(references, key=lambda x: x["name"]):
            name = ref["name"]
            url = ref["url"]
            output = ref["output"]

            if output:
                output_rel = calculate_relative_path(doc_path, SOURCES_DIR / output)
                lines.append(f"| [{name}]({url}) | [Local]({output_rel}) |")
            else:
                lines.append(f"| [{name}]({url}) | - |")

        lines.append("")

    lines.append(BREADCRUMB_END)

    return "\n".join(lines)


def find_insertion_point(content: str) -> tuple[int, str]:
    """Find where to insert breadcrumbs and what section follows."""
    lines = content.split("\n")

    # Look for existing breadcrumb section
    start_idx = None
    end_idx = None
    for i, line in enumerate(lines):
        if BREADCRUMB_START in line:
            start_idx = i
        elif BREADCRUMB_END in line:
            end_idx = i
            break

    if start_idx is not None and end_idx is not None:
        return start_idx, "replace"

    # Look for "Related Documentation" section (insert before it)
    for i, line in enumerate(lines):
        if line.strip().startswith("## Related") or line.strip().startswith(
            "## See Also",
        ):
            return i, "insert"

    # Look for last "---" separator before end
    last_separator = None
    for i, line in enumerate(lines):
        if line.strip() == "---":
            last_separator = i

    if last_separator:
        return last_separator, "insert"

    # Append at end
    return len(lines), "append"


def update_document(doc_path: Path, source_index: dict, dry_run: bool = False) -> bool:
    """Update a single document with breadcrumbs. Returns True if changed."""
    content = doc_path.read_text(encoding="utf-8")

    # Extract references
    references = extract_references(content, source_index)

    # Generate breadcrumb section
    breadcrumbs = generate_breadcrumb_section(doc_path, references)

    # Find insertion point
    insert_idx, action = find_insertion_point(content)
    lines = content.split("\n")

    if action == "replace":
        # Find end marker
        end_idx = None
        for i, line in enumerate(lines):
            if BREADCRUMB_END in line:
                end_idx = i
                break

        if end_idx is None:
            end_idx = insert_idx + 1

        new_lines = [*lines[:insert_idx], breadcrumbs, *lines[end_idx + 1 :]]
    elif action == "insert":
        new_lines = [*lines[:insert_idx], "", breadcrumbs, "", *lines[insert_idx:]]
    else:  # append
        new_lines = [*lines, "", "---", "", breadcrumbs]

    new_content = "\n".join(new_lines)

    # Clean up multiple blank lines
    while "\n\n\n\n" in new_content:
        new_content = new_content.replace("\n\n\n\n", "\n\n\n")

    if new_content != content:
        if not dry_run:
            doc_path.write_text(new_content, encoding="utf-8")
        return True
    return False


def find_design_docs(specific_file: str | None = None) -> list[Path]:
    """Find design documents to process."""
    if specific_file:
        path = Path(specific_file)
        if not path.is_absolute():
            path = DESIGN_DIR / specific_file
        return [path] if path.exists() else []

    # Find all markdown files, excluding archives and INDEX files
    docs = []
    for md_file in DESIGN_DIR.rglob("*.md"):
        # Skip archives
        if ".archive" in str(md_file):
            continue
        # Skip INDEX files
        if md_file.name == "INDEX.md":
            continue
        docs.append(md_file)

    return sorted(docs)


def main():
    parser = argparse.ArgumentParser(
        description="Add source breadcrumbs to design docs",
    )
    parser.add_argument(
        "--dry-run",
        "-n",
        action="store_true",
        help="Preview changes without writing",
    )
    parser.add_argument("--file", "-f", help="Update specific file only")
    args = parser.parse_args()

    print("Loading sources configuration...")
    config = load_sources()
    source_index = build_source_index(config)

    print("Finding design documents...")
    docs = find_design_docs(args.file)
    print(f"  Found {len(docs)} documents")

    if args.dry_run:
        print("\n=== DRY RUN ===\n")

    updated = 0
    unchanged = 0

    for doc in docs:
        rel_path = doc.relative_to(PROJECT_ROOT)
        changed = update_document(doc, source_index, args.dry_run)

        if changed:
            print(f"  {'Would update' if args.dry_run else 'Updated'}: {rel_path}")
            updated += 1
        else:
            unchanged += 1

    print("\n=== SUMMARY ===")
    print(f"{'Would update' if args.dry_run else 'Updated'}: {updated}")
    print(f"Unchanged: {unchanged}")


if __name__ == "__main__":
    main()
