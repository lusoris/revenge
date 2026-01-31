#!/usr/bin/env python3
"""
Generate cross-reference indexes for documentation sources.

Outputs:
1. SOURCES_INDEX.md - Index of all external sources
2. DESIGN_CROSSREF.md - Cross-reference between design docs and sources

Usage:
    python scripts/generate-crossref.py
"""

import re
from collections import defaultdict
from pathlib import Path

import yaml


# Project paths
SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_YAML = SOURCES_DIR / "SOURCES.yaml"
INDEX_YAML = SOURCES_DIR / "INDEX.yaml"

# Output files
SOURCES_INDEX_MD = SOURCES_DIR / "SOURCES_INDEX.md"
DESIGN_CROSSREF_MD = SOURCES_DIR / "DESIGN_CROSSREF.md"


def load_sources() -> dict:
    """Load SOURCES.yaml."""
    with open(SOURCES_YAML, encoding="utf-8") as f:
        return yaml.safe_load(f)


def load_index() -> dict:
    """Load INDEX.yaml with fetch status."""
    if INDEX_YAML.exists():
        with open(INDEX_YAML, encoding="utf-8") as f:
            return yaml.safe_load(f) or {}
    return {}


def find_design_docs() -> list[Path]:
    """Find all markdown files in design directory."""
    return sorted(DESIGN_DIR.rglob("*.md"))


def extract_urls_from_doc(doc_path: Path) -> set[str]:
    """Extract all URLs from a markdown document."""
    content = doc_path.read_text(encoding="utf-8")
    # Match URLs in markdown links [text](url) and bare URLs
    url_pattern = r'https?://[^\s\)\]>"\']+'
    urls = set(re.findall(url_pattern, content))
    # Clean trailing punctuation
    return {url.rstrip(".,;:)") for url in urls}


def extract_packages_from_doc(doc_path: Path) -> set[str]:
    """Extract Go package references from a document."""
    content = doc_path.read_text(encoding="utf-8")
    packages = set()

    # Match github.com/... package paths
    github_pattern = r"github\.com/[\w\-]+/[\w\-]+(?:/[\w\-]+)*"
    packages.update(re.findall(github_pattern, content))

    # Match golang.org/x/... packages
    golang_pattern = r"golang\.org/x/\w+"
    packages.update(re.findall(golang_pattern, content))

    # Match go.uber.org/... packages
    uber_pattern = r"go\.uber\.org/\w+"
    packages.update(re.findall(uber_pattern, content))

    return packages


def generate_sources_index(config: dict, index: dict) -> str:
    """Generate markdown index of all sources."""
    lines = [
        "# External Sources Index",
        "",
        "> Auto-generated cross-reference of all external documentation sources",
        "> ",
        "> Run `python scripts/generate-crossref.py` to regenerate",
        "",
        "---",
        "",
    ]

    sources = config.get("sources", {})
    index_sources = index.get("sources", {})

    # Summary statistics
    total = sum(len(s) for s in sources.values())
    fetched = sum(
        1 for s in index_sources.values() if s.get("status") in ("success", "unchanged")
    )
    failed = sum(1 for s in index_sources.values() if s.get("status") == "failed")
    skipped = sum(1 for s in index_sources.values() if s.get("status") == "skipped")

    lines.extend(
        [
            "## Summary",
            "",
            "| Metric | Count |",
            "|--------|-------|",
            f"| Total Sources | {total} |",
            f"| Fetched | {fetched} |",
            f"| Failed | {failed} |",
            f"| Skipped | {skipped} |",
            "",
            "---",
            "",
        ]
    )

    # Sources by category
    for category, category_sources in sources.items():
        lines.append(f"## {category.replace('_', ' ').title()}")
        lines.append("")
        lines.append("| ID | Name | Type | Status | Output |")
        lines.append("|----|------|------|--------|--------|")

        for source in category_sources:
            source_id = source.get("id", "unknown")
            name = source.get("name", source_id)
            source_type = source.get("type", "html")
            url = source.get("url", "")
            output = source.get("output", "")

            # Get status from index
            idx_entry = index_sources.get(source_id, {})
            status = idx_entry.get("status", "pending")

            # Status emoji
            status_emoji = {
                "success": "✅",
                "unchanged": "✅",
                "failed": "❌",
                "skipped": "⏭️",
                "pending": "⏳",
            }.get(status, "❓")

            # Link to source and output
            name_link = f"[{name}]({url})"
            output_link = f"[{output}]({output})" if output else "-"

            lines.append(
                f"| `{source_id}` | {name_link} | {source_type} | {status_emoji} | {output_link} |"
            )

        lines.append("")

    return "\n".join(lines)


def generate_design_crossref(config: dict, design_docs: list[Path]) -> str:
    """Generate cross-reference between design docs and sources."""
    lines = [
        "# Design Documents ↔ Sources Cross-Reference",
        "",
        "> Auto-generated mapping between design documents and external sources",
        "> ",
        "> Run `python scripts/generate-crossref.py` to regenerate",
        "",
        "---",
        "",
    ]

    # Build source URL -> source mapping
    sources = config.get("sources", {})
    url_to_source = {}
    package_to_source = {}

    for category, category_sources in sources.items():
        for source in category_sources:
            source_id = source.get("id")
            url = source.get("url", "")
            name = source.get("name", source_id)

            if url:
                # Normalize URL for matching
                normalized = (
                    url.replace("https://", "").replace("http://", "").rstrip("/")
                )
                url_to_source[normalized] = {
                    "id": source_id,
                    "name": name,
                    "url": url,
                    "category": category,
                }

                # Extract package path from pkg.go.dev or github URLs
                if "pkg.go.dev/" in url:
                    pkg = url.split("pkg.go.dev/")[-1]
                    package_to_source[pkg] = {
                        "id": source_id,
                        "name": name,
                        "url": url,
                        "category": category,
                    }
                elif "github.com/" in url:
                    # github.com/owner/repo -> owner/repo
                    parts = url.replace("https://github.com/", "").split("/")
                    if len(parts) >= 2:
                        pkg = f"github.com/{parts[0]}/{parts[1]}"
                        package_to_source[pkg] = {
                            "id": source_id,
                            "name": name,
                            "url": url,
                            "category": category,
                        }

    # Analyze each design doc
    doc_references = {}  # doc_path -> set of source IDs
    source_referenced_by = defaultdict(set)  # source_id -> set of doc paths

    for doc_path in design_docs:
        rel_path = doc_path.relative_to(PROJECT_ROOT)
        urls = extract_urls_from_doc(doc_path)
        packages = extract_packages_from_doc(doc_path)

        found_sources = set()

        # Match URLs
        for url in urls:
            normalized = url.replace("https://", "").replace("http://", "").rstrip("/")
            # Check for exact or partial match
            for source_url, source_info in url_to_source.items():
                if source_url in normalized or normalized in source_url:
                    found_sources.add(source_info["id"])
                    break

        # Match packages
        for pkg in packages:
            if pkg in package_to_source:
                found_sources.add(package_to_source[pkg]["id"])
            else:
                # Try partial match
                for source_pkg, source_info in package_to_source.items():
                    if pkg in source_pkg or source_pkg in pkg:
                        found_sources.add(source_info["id"])
                        break

        if found_sources:
            doc_references[str(rel_path)] = found_sources
            for source_id in found_sources:
                source_referenced_by[source_id].add(str(rel_path))

    # Output: Design docs by referenced sources
    lines.extend(
        [
            "## Design Documents → Sources",
            "",
            "Which external sources are referenced by each design document.",
            "",
        ]
    )

    for doc_path, source_ids in sorted(doc_references.items()):
        doc_name = Path(doc_path).stem
        lines.append(f"### [{doc_name}](../../{doc_path})")
        lines.append("")
        for source_id in sorted(source_ids):
            # Find source info
            for category, category_sources in sources.items():
                for source in category_sources:
                    if source.get("id") == source_id:
                        name = source.get("name", source_id)
                        url = source.get("url", "")
                        lines.append(f"- `{source_id}`: [{name}]({url})")
                        break
        lines.append("")

    # Output: Sources by referencing docs
    lines.extend(
        [
            "---",
            "",
            "## Sources → Design Documents",
            "",
            "Which design documents reference each external source.",
            "",
        ]
    )

    for source_id, doc_paths in sorted(source_referenced_by.items()):
        # Find source info
        source_name = source_id
        source_url = ""
        for category, category_sources in sources.items():
            for source in category_sources:
                if source.get("id") == source_id:
                    source_name = source.get("name", source_id)
                    source_url = source.get("url", "")
                    break

        lines.append(f"### `{source_id}`: [{source_name}]({source_url})")
        lines.append("")
        lines.append("Referenced by:")
        for doc_path in sorted(doc_paths):
            doc_name = Path(doc_path).stem
            lines.append(f"- [{doc_name}](../../{doc_path})")
        lines.append("")

    # Output: Unreferenced sources
    all_source_ids = set()
    for category_sources in sources.values():
        for source in category_sources:
            all_source_ids.add(source.get("id"))

    unreferenced = all_source_ids - set(source_referenced_by.keys())
    if unreferenced:
        lines.extend(
            [
                "---",
                "",
                "## Unreferenced Sources",
                "",
                "Sources defined in SOURCES.yaml but not referenced in any design document.",
                "",
            ]
        )
        for source_id in sorted(unreferenced):
            for category, category_sources in sources.items():
                for source in category_sources:
                    if source.get("id") == source_id:
                        name = source.get("name", source_id)
                        url = source.get("url", "")
                        lines.append(f"- `{source_id}`: [{name}]({url})")
                        break
        lines.append("")

    return "\n".join(lines)


def main():
    print("Loading sources configuration...")
    config = load_sources()
    index = load_index()

    print("Finding design documents...")
    design_docs = find_design_docs()
    print(f"  Found {len(design_docs)} design documents")

    print("\nGenerating sources index...")
    sources_index = generate_sources_index(config, index)
    SOURCES_INDEX_MD.write_text(sources_index, encoding="utf-8")
    print(f"  Wrote {SOURCES_INDEX_MD}")

    print("\nGenerating design cross-reference...")
    design_crossref = generate_design_crossref(config, design_docs)
    DESIGN_CROSSREF_MD.write_text(design_crossref, encoding="utf-8")
    print(f"  Wrote {DESIGN_CROSSREF_MD}")

    print("\nDone!")


if __name__ == "__main__":
    main()
