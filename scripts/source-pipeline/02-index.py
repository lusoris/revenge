#!/usr/bin/env python3
"""Generate SOURCES.md index for external documentation sources.

Combines source listing with cross-reference information into a single
consolidated file, replacing the previous SOURCES_INDEX.md and DESIGN_CROSSREF.md.

Usage:
    python scripts/source-pipeline/02-index.py           # Dry run (default)
    python scripts/source-pipeline/02-index.py --apply   # Write files
"""

from __future__ import annotations

import argparse
import re
import sys
from collections import defaultdict
from pathlib import Path

import yaml


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_YAML = SOURCES_DIR / "SOURCES.yaml"
INDEX_YAML = SOURCES_DIR / "INDEX.yaml"

OUTPUT_FILE = SOURCES_DIR / "SOURCES.md"


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
    url_pattern = r'https?://[^\s\)\]>"\']+'
    urls = set(re.findall(url_pattern, content))
    return {url.rstrip(".,;:)") for url in urls}


def extract_packages_from_doc(doc_path: Path) -> set[str]:
    """Extract Go package references from a document."""
    content = doc_path.read_text(encoding="utf-8")
    packages = set()

    github_pattern = r"github\.com/[\w\-]+/[\w\-]+(?:/[\w\-]+)*"
    packages.update(re.findall(github_pattern, content))

    golang_pattern = r"golang\.org/x/\w+"
    packages.update(re.findall(golang_pattern, content))

    uber_pattern = r"go\.uber\.org/\w+"
    packages.update(re.findall(uber_pattern, content))

    return packages


def build_source_mapping(config: dict) -> tuple[dict, dict]:
    """Build URL and package to source info mappings."""
    sources = config.get("sources", {})
    url_to_source = {}
    package_to_source = {}

    for category, category_sources in sources.items():
        for source in category_sources:
            source_id = source.get("id")
            url = source.get("url", "")
            name = source.get("name", source_id)
            info = {"id": source_id, "name": name, "url": url, "category": category}

            if url:
                normalized = (
                    url.replace("https://", "").replace("http://", "").rstrip("/")
                )
                url_to_source[normalized] = info

                if "pkg.go.dev/" in url:
                    pkg = url.split("pkg.go.dev/")[-1]
                    package_to_source[pkg] = info
                elif "github.com/" in url:
                    parts = url.replace("https://github.com/", "").split("/")
                    if len(parts) >= 2:
                        pkg = f"github.com/{parts[0]}/{parts[1]}"
                        package_to_source[pkg] = info

    return url_to_source, package_to_source


def analyze_design_docs(
    design_docs: list[Path],
    url_to_source: dict,
    package_to_source: dict,
) -> tuple[dict, dict]:
    """Analyze design docs for source references."""
    doc_references: dict[str, set[str]] = {}
    source_referenced_by: dict[str, set[str]] = defaultdict(set)

    for doc_path in design_docs:
        rel_path = doc_path.relative_to(PROJECT_ROOT)
        urls = extract_urls_from_doc(doc_path)
        packages = extract_packages_from_doc(doc_path)

        found_sources = set()

        for url in urls:
            normalized = url.replace("https://", "").replace("http://", "").rstrip("/")
            for source_url, source_info in url_to_source.items():
                if source_url in normalized or normalized in source_url:
                    found_sources.add(source_info["id"])
                    break

        for pkg in packages:
            if pkg in package_to_source:
                found_sources.add(package_to_source[pkg]["id"])
            else:
                for source_pkg, source_info in package_to_source.items():
                    if pkg in source_pkg or source_pkg in pkg:
                        found_sources.add(source_info["id"])
                        break

        if found_sources:
            doc_references[str(rel_path)] = found_sources
            for source_id in found_sources:
                source_referenced_by[source_id].add(str(rel_path))

    return doc_references, dict(source_referenced_by)


def generate_sources_md(
    config: dict,
    index: dict,
    source_referenced_by: dict,
) -> str:
    """Generate consolidated SOURCES.md content."""
    lines = [
        "# External Documentation Sources",
        "",
        "> Consolidated index of all external documentation sources",
        "> with cross-references to design documents.",
        ">",
        "> Auto-generated by `scripts/source-pipeline/02-index.py`",
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
    referenced = len(source_referenced_by)
    unreferenced = total - referenced

    lines.extend(
        [
            "## Summary",
            "",
            "| Metric | Count |",
            "|--------|-------|",
            f"| Total Sources | {total} |",
            f"| Fetched | {fetched} |",
            f"| Failed | {failed} |",
            f"| Referenced by Design Docs | {referenced} |",
            f"| Unreferenced | {unreferenced} |",
            "",
            "---",
            "",
        ]
    )

    # Sources by category with reference info
    for category, category_sources in sources.items():
        lines.append(f"## {category.replace('_', ' ').title()}")
        lines.append("")
        lines.append("| ID | Name | Status | Referenced By |")
        lines.append("|----|------|--------|---------------|")

        for source in category_sources:
            source_id = source.get("id", "unknown")
            name = source.get("name", source_id)
            url = source.get("url", "")

            idx_entry = index_sources.get(source_id, {})
            status = idx_entry.get("status", "pending")

            status_emoji = {
                "success": "✅",
                "unchanged": "✅",
                "failed": "❌",
                "skipped": "⏭️",
                "pending": "⏳",
            }.get(status, "❓")

            refs = source_referenced_by.get(source_id, set())
            ref_count = len(refs)
            ref_text = f"{ref_count} docs" if ref_count > 0 else "-"

            name_link = f"[{name}]({url})" if url else name
            lines.append(
                f"| `{source_id}` | {name_link} | {status_emoji} | {ref_text} |"
            )

        lines.append("")

    # Cross-reference section
    lines.extend(
        [
            "---",
            "",
            "## Cross-References",
            "",
            "### Sources → Design Documents",
            "",
            "Which design documents reference each source.",
            "",
        ]
    )

    for source_id in sorted(source_referenced_by.keys()):
        doc_paths = source_referenced_by[source_id]

        source_name = source_id
        source_url = ""
        for category_sources in sources.values():
            for source in category_sources:
                if source.get("id") == source_id:
                    source_name = source.get("name", source_id)
                    source_url = source.get("url", "")
                    break

        lines.append(f"#### `{source_id}`: [{source_name}]({source_url})")
        lines.append("")
        for doc_path in sorted(doc_paths):
            doc_name = Path(doc_path).stem
            lines.append(f"- [{doc_name}](../{doc_path})")
        lines.append("")

    # Unreferenced sources
    all_source_ids = set()
    for category_sources in sources.values():
        for source in category_sources:
            all_source_ids.add(source.get("id"))

    unreferenced_ids = all_source_ids - set(source_referenced_by.keys())
    if unreferenced_ids:
        lines.extend(
            [
                "---",
                "",
                "## Unreferenced Sources",
                "",
                "Sources defined but not yet referenced in any design document.",
                "These are kept for future design docs.",
                "",
            ]
        )
        for source_id in sorted(unreferenced_ids):
            for category_sources in sources.values():
                for source in category_sources:
                    if source.get("id") == source_id:
                        name = source.get("name", source_id)
                        url = source.get("url", "")
                        lines.append(f"- `{source_id}`: [{name}]({url})")
                        break
        lines.append("")

    lines.extend(
        [
            "---",
            "",
            "*Auto-generated by `scripts/source-pipeline/02-index.py`*",
            "",
        ]
    )

    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(description="Generate consolidated sources index")
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Write output file (default: dry-run)",
    )
    args = parser.parse_args()

    dry_run = not args.apply

    if dry_run:
        print("=== DRY RUN (use --apply to write) ===\n")

    print("Loading sources configuration...")
    config = load_sources()
    index = load_index()

    print("Finding design documents...")
    design_docs = find_design_docs()
    print(f"  Found {len(design_docs)} design documents")

    print("Building source mappings...")
    url_to_source, package_to_source = build_source_mapping(config)

    print("Analyzing cross-references...")
    _, source_referenced_by = analyze_design_docs(
        design_docs, url_to_source, package_to_source
    )
    print(f"  Found {len(source_referenced_by)} sources referenced")

    print("\nGenerating SOURCES.md...")
    content = generate_sources_md(config, index, source_referenced_by)

    if dry_run:
        print(f"  Would write {len(content)} bytes to {OUTPUT_FILE}")
        print("\n=== DRY RUN complete. Use --apply to write. ===")
    else:
        OUTPUT_FILE.write_text(content, encoding="utf-8")
        print(f"  Wrote {OUTPUT_FILE}")
        print("\nDone!")

    return 0


if __name__ == "__main__":
    sys.exit(main())
