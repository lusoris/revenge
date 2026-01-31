#!/usr/bin/env python3
"""
Fetch external documentation sources defined in SOURCES.yaml.

Features:
    - Content hash tracking: Only updates files when content actually changes
    - Timestamping: Records fetch time in header and INDEX.yaml
    - Version tracking: Stores content hash for change detection

Usage:
    python scripts/fetch-sources.py                    # Fetch all sources (skip unchanged)
    python scripts/fetch-sources.py --category apis   # Fetch specific category
    python scripts/fetch-sources.py --id tmdb         # Fetch specific source
    python scripts/fetch-sources.py --dry-run         # Preview without fetching
    python scripts/fetch-sources.py --force           # Force update even if unchanged
"""

import argparse
import hashlib
import json
import sys
import time
from datetime import UTC, datetime
from pathlib import Path
from typing import Any

import requests
import yaml
from bs4 import BeautifulSoup


# Optional html2text - fall back to simple text extraction
try:
    import html2text

    HAS_HTML2TEXT = True
except ImportError:
    HAS_HTML2TEXT = False

# Project paths
SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
SOURCES_YAML = SOURCES_DIR / "SOURCES.yaml"
INDEX_YAML = SOURCES_DIR / "INDEX.yaml"


class SourceFetcher:
    """Fetches and processes external documentation sources."""

    def __init__(self, config: dict[str, Any]):
        self.config = config
        self.fetch_config = config.get("fetch_config", {})
        self.delay = self.fetch_config.get("delay_between_requests", 2)
        self.timeout = self.fetch_config.get("timeout", 30)
        self.retry_count = self.fetch_config.get("retry_count", 3)
        self.user_agent = self.fetch_config.get("user_agent", "Revenge-DocFetcher/1.0")
        self.session = requests.Session()
        self.session.headers.update({"User-Agent": self.user_agent})
        if HAS_HTML2TEXT:
            self.html_converter = html2text.HTML2Text()
            self.html_converter.ignore_links = False
            self.html_converter.ignore_images = True
            self.html_converter.body_width = 0  # No wrapping
        else:
            self.html_converter = None
        self.results: list[dict] = []
        self.force_update = False
        # Load existing index for change detection
        self.existing_index = self._load_existing_index()

    def _load_existing_index(self) -> dict:
        """Load existing INDEX.yaml for change detection."""
        if INDEX_YAML.exists():
            with open(INDEX_YAML, encoding="utf-8") as f:
                return yaml.safe_load(f) or {}
        return {}

    def _content_hash(self, content: str) -> str:
        """Calculate SHA256 hash of content."""
        return hashlib.sha256(content.encode("utf-8")).hexdigest()[:16]

    def _content_changed(self, source_id: str, new_hash: str) -> bool:
        """Check if content has changed from previous fetch."""
        existing = self.existing_index.get("sources", {}).get(source_id, {})
        old_hash = existing.get("content_hash")
        return old_hash != new_hash

    def fetch_url(self, url: str) -> requests.Response | None:
        """Fetch URL with retries and error handling."""
        for attempt in range(self.retry_count):
            try:
                response = self.session.get(url, timeout=self.timeout)
                response.raise_for_status()
                return response
            except requests.RequestException as e:
                if attempt < self.retry_count - 1:
                    print(f"  Retry {attempt + 1}/{self.retry_count}: {e}")
                    time.sleep(self.delay)
                else:
                    print(f"  Failed after {self.retry_count} attempts: {e}")
                    return None
        return None

    def process_html(
        self, response: requests.Response, selectors: list[str] | None
    ) -> str:
        """Extract and convert HTML content to markdown."""
        soup = BeautifulSoup(response.text, "lxml")

        # Remove script and style elements
        for element in soup(["script", "style", "nav", "footer", "header"]):
            element.decompose()

        if selectors:
            # Extract specific elements
            content_parts = []
            for selector in selectors:
                elements = soup.select(selector)
                for el in elements:
                    content_parts.append(str(el))
            html_content = "\n".join(content_parts)
        else:
            # Use main content or body
            main = soup.find("main") or soup.find("article") or soup.find("body")
            html_content = str(main) if main else response.text

        if self.html_converter:
            return self.html_converter.handle(html_content)
        else:
            # Fallback: simple text extraction without html2text
            soup = BeautifulSoup(html_content, "lxml")
            return soup.get_text(separator="\n\n", strip=True)

    def process_github_readme(self, url: str) -> str | None:
        """Fetch GitHub README in raw markdown format."""
        # Convert github.com URL to raw.githubusercontent.com
        if "github.com" in url:
            parts = url.replace("https://github.com/", "").split("/")
            if len(parts) >= 2:
                owner, repo = parts[0], parts[1]
                raw_url = (
                    f"https://raw.githubusercontent.com/{owner}/{repo}/HEAD/README.md"
                )
                response = self.fetch_url(raw_url)
                if response:
                    return response.text
        return None

    def process_json(self, response: requests.Response) -> str:
        """Pretty-print JSON content."""
        try:
            data = response.json()
            return json.dumps(data, indent=2)
        except json.JSONDecodeError:
            return response.text

    def process_graphql_schema(self, url: str) -> str | None:
        """Fetch GraphQL schema using introspection."""
        # For GraphQL endpoints, we'd need introspection
        # For now, just note it requires manual handling
        return f"# GraphQL Schema\n\nEndpoint: {url}\n\nNote: Requires introspection query to fetch schema."

    def fetch_source(self, source: dict[str, Any]) -> dict[str, Any]:
        """Fetch and process a single source."""
        source_id = source.get("id", "unknown")
        name = source.get("name", source_id)
        url = source.get("url", "")
        source_type = source.get("type", "html")
        output_path = source.get("output", "")
        selectors = source.get("selectors")
        note = source.get("note", "")

        result = {
            "id": source_id,
            "name": name,
            "url": url,
            "output": output_path,
            "status": "pending",
            "fetched_at": None,
            "error": None,
        }

        print(f"  [{source_id}] {name}")

        # Skip manual sources
        if source_type == "manual":
            result["status"] = "skipped"
            result["error"] = f"Manual source: {note}"
            print(f"    Skipped (manual): {note}")
            return result

        # Fetch content based on type
        content = None

        if source_type == "github_readme":
            content = self.process_github_readme(url)
        elif source_type == "graphql_schema":
            content = self.process_graphql_schema(url)
        elif source_type == "json":
            response = self.fetch_url(url)
            if response:
                content = self.process_json(response)
        else:  # html
            response = self.fetch_url(url)
            if response:
                content = self.process_html(response, selectors)

        if content:
            # Calculate content hash (without header to detect actual content changes)
            content_hash = self._content_hash(content)
            result["content_hash"] = content_hash

            # Check if content changed
            if not self.force_update and not self._content_changed(
                source_id, content_hash
            ):
                # Get previous fetch time from existing index
                existing = self.existing_index.get("sources", {}).get(source_id, {})
                result["status"] = "unchanged"
                result["fetched_at"] = existing.get("fetched_at")
                print(f"    Unchanged (hash: {content_hash})")
                return result

            # Add header with source info
            now = datetime.now(UTC).isoformat()
            header = f"""# {name}

> Source: {url}
> Fetched: {now}
> Content-Hash: {content_hash}
> Type: {source_type}

---

"""
            full_content = header + content

            # Save to output file
            output_file = SOURCES_DIR / output_path
            output_file.parent.mkdir(parents=True, exist_ok=True)
            output_file.write_text(full_content, encoding="utf-8")

            result["status"] = "success"
            result["fetched_at"] = now
            print(f"    Saved to {output_path} (hash: {content_hash})")
        else:
            result["status"] = "failed"
            result["error"] = "Failed to fetch content"
            print("    Failed to fetch")

        return result

    def fetch_category(self, category: str, sources: list[dict]) -> list[dict]:
        """Fetch all sources in a category."""
        print(f"\n=== {category.upper()} ===")
        results = []
        for source in sources:
            result = self.fetch_source(source)
            results.append(result)
            time.sleep(self.delay)
        return results

    def fetch_all(
        self,
        category_filter: str | None = None,
        source_id_filter: str | None = None,
    ) -> list[dict]:
        """Fetch all sources, optionally filtered."""
        all_results = []
        sources = self.config.get("sources", {})

        for category, category_sources in sources.items():
            if category_filter and category != category_filter:
                continue

            if source_id_filter:
                # Filter to specific source
                category_sources = [
                    s for s in category_sources if s.get("id") == source_id_filter
                ]
                if not category_sources:
                    continue

            results = self.fetch_category(category, category_sources)
            all_results.extend(results)

        self.results = all_results
        return all_results

    def update_index(self):
        """Update INDEX.yaml with fetch results."""
        # Load existing index or create new
        if INDEX_YAML.exists():
            with open(INDEX_YAML, encoding="utf-8") as f:
                index = yaml.safe_load(f) or {}
        else:
            index = {}

        # Update with new results
        index["last_updated"] = datetime.now(UTC).isoformat()
        index["total_sources"] = len(self.results)
        index["successful"] = sum(1 for r in self.results if r["status"] == "success")
        index["unchanged"] = sum(1 for r in self.results if r["status"] == "unchanged")
        index["failed"] = sum(1 for r in self.results if r["status"] == "failed")
        index["skipped"] = sum(1 for r in self.results if r["status"] == "skipped")

        # Update individual source statuses
        if "sources" not in index:
            index["sources"] = {}

        for result in self.results:
            entry = {
                "name": result["name"],
                "url": result["url"],
                "output": result["output"],
                "status": result["status"],
                "fetched_at": result["fetched_at"],
                "error": result["error"],
            }
            # Add content_hash if present
            if "content_hash" in result:
                entry["content_hash"] = result["content_hash"]
            index["sources"][result["id"]] = entry

        # Save index
        with open(INDEX_YAML, "w", encoding="utf-8") as f:
            yaml.dump(
                index, f, default_flow_style=False, allow_unicode=True, sort_keys=False
            )

        print("\n=== SUMMARY ===")
        print(f"Total: {index['total_sources']}")
        print(f"Updated: {index['successful']}")
        print(f"Unchanged: {index['unchanged']}")
        print(f"Failed: {index['failed']}")
        print(f"Skipped: {index['skipped']}")


def main():
    parser = argparse.ArgumentParser(description="Fetch external documentation sources")
    parser.add_argument(
        "--category", "-c", help="Fetch only specific category", default=None
    )
    parser.add_argument(
        "--id", "-i", help="Fetch only specific source by ID", default=None
    )
    parser.add_argument(
        "--dry-run", "-n", action="store_true", help="Preview without fetching"
    )
    parser.add_argument(
        "--force", "-f", action="store_true", help="Force update even if unchanged"
    )
    args = parser.parse_args()

    # Load SOURCES.yaml
    if not SOURCES_YAML.exists():
        print(f"Error: {SOURCES_YAML} not found")
        sys.exit(1)

    with open(SOURCES_YAML, encoding="utf-8") as f:
        config = yaml.safe_load(f)

    if args.dry_run:
        print("=== DRY RUN ===")
        sources = config.get("sources", {})
        for category, category_sources in sources.items():
            if args.category and category != args.category:
                continue
            print(f"\n{category}:")
            for source in category_sources:
                if args.id and source.get("id") != args.id:
                    continue
                print(f"  - {source.get('id')}: {source.get('name')}")
                print(f"    URL: {source.get('url')}")
                print(f"    Output: {source.get('output')}")
        return

    # Run fetcher
    fetcher = SourceFetcher(config)
    fetcher.force_update = args.force
    fetcher.fetch_all(category_filter=args.category, source_id_filter=args.id)
    fetcher.update_index()


if __name__ == "__main__":
    main()
