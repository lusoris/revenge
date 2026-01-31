#!/usr/bin/env python3
"""
Find external URLs in design docs that are missing from SOURCES.yaml.
"""

import re
from collections import defaultdict
from pathlib import Path

import yaml

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
SOURCES_YAML = SOURCES_DIR / "SOURCES.yaml"


def load_sources() -> set[str]:
    """Load all URLs from SOURCES.yaml."""
    with open(SOURCES_YAML, encoding="utf-8") as f:
        config = yaml.safe_load(f)

    urls = set()
    for category, sources in config.get("sources", {}).items():
        for source in sources:
            url = source.get("url", "")
            if url:
                # Normalize
                normalized = url.replace("https://", "").replace("http://", "").rstrip("/")
                urls.add(normalized)
                # Also add domain
                domain = normalized.split("/")[0]
                urls.add(domain)

    return urls


def extract_urls(content: str) -> set[str]:
    """Extract external URLs from content."""
    pattern = r'https?://[^\s\)\]>"\']+'
    urls = set()
    for match in re.findall(pattern, content):
        # Clean up
        url = match.rstrip(".,;:)")
        # Skip internal/local
        if "localhost" in url or "127.0.0.1" in url:
            continue
        if "example.com" in url or "example.org" in url:
            continue
        urls.add(url)
    return urls


def main():
    print("Loading SOURCES.yaml...")
    known_urls = load_sources()
    print(f"  Found {len(known_urls)} known URL patterns")

    print("\nScanning design docs...")
    missing = defaultdict(set)  # url -> set of docs that reference it
    all_urls = set()

    for md_file in sorted(DESIGN_DIR.rglob("*.md")):
        if ".archive" in str(md_file):
            continue

        content = md_file.read_text(encoding="utf-8")
        urls = extract_urls(content)
        all_urls.update(urls)

        for url in urls:
            normalized = url.replace("https://", "").replace("http://", "").rstrip("/")
            domain = normalized.split("/")[0]

            # Check if any known URL matches
            is_known = False
            for known in known_urls:
                if known in normalized or normalized in known:
                    is_known = True
                    break

            if not is_known:
                rel_path = md_file.relative_to(DESIGN_DIR)
                missing[url].add(str(rel_path))

    print(f"  Found {len(all_urls)} total external URLs")
    print(f"  Found {len(missing)} potentially missing from SOURCES.yaml")

    # Group by domain
    by_domain = defaultdict(list)
    for url, docs in sorted(missing.items()):
        domain = url.replace("https://", "").replace("http://", "").split("/")[0]
        by_domain[domain].append((url, docs))

    print("\n" + "=" * 60)
    print("MISSING URLs BY DOMAIN")
    print("=" * 60)

    for domain, items in sorted(by_domain.items()):
        print(f"\n## {domain}")
        for url, docs in items[:5]:  # Limit to 5 per domain
            print(f"  - {url}")
            for doc in sorted(docs)[:3]:
                print(f"      in: {doc}")
        if len(items) > 5:
            print(f"  ... and {len(items) - 5} more URLs from this domain")


if __name__ == "__main__":
    main()
