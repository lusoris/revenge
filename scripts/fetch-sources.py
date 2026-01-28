#!/usr/bin/env python3
"""
External Documentation Source Fetcher

Fetches external documentation defined in docs/dev/sources/SOURCES.yaml
and stores them in docs/dev/sources/{category}/.

Usage:
    python scripts/fetch-sources.py              # Fetch all sources
    python scripts/fetch-sources.py --category go  # Fetch only 'go' category
    python scripts/fetch-sources.py --id tmdb      # Fetch single source by ID
    python scripts/fetch-sources.py --dry-run      # Show what would be fetched
"""

import argparse
import hashlib
import os
import re
import sys
import time
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

import requests
import yaml
from bs4 import BeautifulSoup

# =============================================================================
# Configuration
# =============================================================================

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
SOURCES_FILE = PROJECT_ROOT / "docs" / "dev" / "sources" / "SOURCES.yaml"
INDEX_FILE = PROJECT_ROOT / "docs" / "dev" / "sources" / "INDEX.yaml"
OUTPUT_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"

# Safety: Only allow writes within this directory
ALLOWED_OUTPUT_PREFIX = str(OUTPUT_DIR.resolve())

# =============================================================================
# Fetcher Classes
# =============================================================================


class BaseFetcher:
    """Base class for all fetchers."""

    def __init__(self, config: dict[str, Any]):
        self.config = config
        self.delay = config.get("delay_between_requests", 2)
        self.timeout = config.get("timeout", 30)
        self.retry_count = config.get("retry_count", 3)
        self.user_agent = config.get("user_agent", "Revenge-DocFetcher/1.0")

    def fetch(self, source: dict[str, Any]) -> tuple[str, str | None]:
        """
        Fetch a single source.

        Returns:
            Tuple of (content, error_message)
            If successful, error_message is None
        """
        raise NotImplementedError


class GitHubReadmeFetcher(BaseFetcher):
    """Fetches GitHub README files via raw.githubusercontent.com."""

    # Pattern: https://github.com/{owner}/{repo}
    GITHUB_PATTERN = re.compile(r"^https?://github\.com/([^/]+)/([^/]+)/?$")

    def fetch(self, source: dict[str, Any]) -> tuple[str, str | None]:
        url = source["url"]
        match = self.GITHUB_PATTERN.match(url)

        if not match:
            return "", f"Invalid GitHub URL format: {url}"

        owner, repo = match.groups()
        # Try common README filenames
        readme_files = ["README.md", "readme.md", "Readme.md", "README.rst", "README"]

        for readme in readme_files:
            raw_url = f"https://raw.githubusercontent.com/{owner}/{repo}/HEAD/{readme}"
            try:
                response = requests.get(
                    raw_url,
                    headers={"User-Agent": self.user_agent},
                    timeout=self.timeout,
                )
                if response.status_code == 200:
                    content = response.text
                    markdown = self._format_readme(source, content, url, raw_url)
                    return markdown, None
            except requests.RequestException:
                continue

        return "", f"Could not find README for {url}"

    def _format_readme(
        self, source: dict[str, Any], content: str, orig_url: str, raw_url: str
    ) -> str:
        """Format README with metadata header."""
        name = source.get("name", source["id"])
        now = datetime.now(timezone.utc).isoformat()

        header = f"""# {name}

> Auto-fetched from [{orig_url}]({orig_url})
> Raw source: [{raw_url}]({raw_url})
> Last Updated: {now}

---

"""
        return header + content


class HTMLFetcher(BaseFetcher):
    """Fetches HTML pages and extracts documentation content."""

    def fetch(self, source: dict[str, Any]) -> tuple[str, str | None]:
        url = source["url"]
        selectors = source.get("selectors", [])

        for attempt in range(self.retry_count):
            try:
                response = requests.get(
                    url,
                    headers={"User-Agent": self.user_agent},
                    timeout=self.timeout,
                )
                response.raise_for_status()

                soup = BeautifulSoup(response.text, "html.parser")

                # Extract content using selectors if provided
                if selectors:
                    content_parts = []
                    for selector in selectors:
                        elements = soup.select(selector)
                        for el in elements:
                            content_parts.append(
                                el.get_text(separator="\n", strip=True)
                            )
                    content = "\n\n".join(content_parts)
                else:
                    # Try common documentation selectors
                    main = (
                        soup.select_one("main")
                        or soup.select_one("article")
                        or soup.select_one(".content")
                        or soup.select_one("#content")
                        or soup.body
                    )
                    content = main.get_text(separator="\n", strip=True) if main else ""

                if not content:
                    return "", f"No content extracted from {url}"

                # Convert to markdown format
                markdown = self._html_to_markdown(source, content, url)
                return markdown, None

            except requests.RequestException as e:
                if attempt < self.retry_count - 1:
                    time.sleep(self.delay * (attempt + 1))
                    continue
                return "", f"Request failed after {self.retry_count} attempts: {e}"

        return "", "Unknown error"

    def _html_to_markdown(self, source: dict[str, Any], content: str, url: str) -> str:
        """Convert extracted content to markdown with metadata header."""
        name = source.get("name", source["id"])
        now = datetime.now(timezone.utc).isoformat()

        header = f"""# {name}

> Auto-fetched from [{url}]({url})
> Last Updated: {now}

---

"""
        # Clean up content
        lines = content.split("\n")
        cleaned_lines = []
        prev_empty = False

        for line in lines:
            line = line.strip()
            if not line:
                if not prev_empty:
                    cleaned_lines.append("")
                    prev_empty = True
            else:
                cleaned_lines.append(line)
                prev_empty = False

        return header + "\n".join(cleaned_lines)


class JSONFetcher(BaseFetcher):
    """Fetches JSON files (e.g., OpenAPI specs) and saves them directly."""

    def fetch(self, source: dict[str, Any]) -> tuple[str, str | None]:
        url = source["url"]

        for attempt in range(self.retry_count):
            try:
                response = requests.get(
                    url,
                    headers={"User-Agent": self.user_agent},
                    timeout=self.timeout,
                )
                response.raise_for_status()

                # Return raw JSON content
                return response.text, None

            except requests.RequestException as e:
                if attempt < self.retry_count - 1:
                    time.sleep(self.delay * (attempt + 1))
                    continue
                return "", f"Request failed after {self.retry_count} attempts: {e}"

        return "", "Unknown error"


class GraphQLSchemaFetcher(BaseFetcher):
    """Fetches GraphQL schemas via introspection."""

    INTROSPECTION_QUERY = """
    query IntrospectionQuery {
      __schema {
        queryType { name }
        mutationType { name }
        subscriptionType { name }
        types {
          ...FullType
        }
        directives {
          name
          description
          locations
          args {
            ...InputValue
          }
        }
      }
    }

    fragment FullType on __Type {
      kind
      name
      description
      fields(includeDeprecated: true) {
        name
        description
        args {
          ...InputValue
        }
        type {
          ...TypeRef
        }
        isDeprecated
        deprecationReason
      }
      inputFields {
        ...InputValue
      }
      interfaces {
        ...TypeRef
      }
      enumValues(includeDeprecated: true) {
        name
        description
        isDeprecated
        deprecationReason
      }
      possibleTypes {
        ...TypeRef
      }
    }

    fragment InputValue on __InputValue {
      name
      description
      type {
        ...TypeRef
      }
      defaultValue
    }

    fragment TypeRef on __Type {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                  ofType {
                    kind
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
    """

    def fetch(self, source: dict[str, Any]) -> tuple[str, str | None]:
        url = source["url"]

        for attempt in range(self.retry_count):
            try:
                response = requests.post(
                    url,
                    headers={
                        "User-Agent": self.user_agent,
                        "Content-Type": "application/json",
                    },
                    json={"query": self.INTROSPECTION_QUERY},
                    timeout=self.timeout,
                )
                response.raise_for_status()

                data = response.json()
                if "errors" in data:
                    return "", f"GraphQL errors: {data['errors']}"

                # Convert schema to SDL format
                schema = data.get("data", {}).get("__schema", {})
                sdl = self._schema_to_sdl(schema)

                name = source.get("name", source["id"])
                now = datetime.now(timezone.utc).isoformat()

                header = f"""# {name} GraphQL Schema

# Auto-fetched from {url}
# Last Updated: {now}

"""
                return header + sdl, None

            except requests.RequestException as e:
                if attempt < self.retry_count - 1:
                    time.sleep(self.delay * (attempt + 1))
                    continue
                return "", f"Request failed after {self.retry_count} attempts: {e}"

        return "", "Unknown error"

    def _schema_to_sdl(self, schema: dict[str, Any]) -> str:
        """Convert introspection result to SDL format."""
        lines = []
        types = schema.get("types", [])

        # Filter out built-in types
        custom_types = [t for t in types if not t["name"].startswith("__")]

        for t in sorted(custom_types, key=lambda x: x["name"]):
            kind = t["kind"]
            name = t["name"]
            description = t.get("description", "")

            if description:
                lines.append(f'"""{description}"""')

            if kind == "OBJECT":
                interfaces = t.get("interfaces", [])
                implements = ""
                if interfaces:
                    impl_names = [i["name"] for i in interfaces]
                    implements = f" implements {' & '.join(impl_names)}"

                lines.append(f"type {name}{implements} {{")
                for field in t.get("fields", []):
                    field_desc = field.get("description", "")
                    if field_desc:
                        lines.append(f'  """{field_desc}"""')
                    field_type = self._type_ref_to_string(field["type"])
                    args = self._format_args(field.get("args", []))
                    lines.append(f"  {field['name']}{args}: {field_type}")
                lines.append("}")
                lines.append("")

            elif kind == "INPUT_OBJECT":
                lines.append(f"input {name} {{")
                for field in t.get("inputFields", []):
                    field_type = self._type_ref_to_string(field["type"])
                    lines.append(f"  {field['name']}: {field_type}")
                lines.append("}")
                lines.append("")

            elif kind == "ENUM":
                lines.append(f"enum {name} {{")
                for value in t.get("enumValues", []):
                    lines.append(f"  {value['name']}")
                lines.append("}")
                lines.append("")

            elif kind == "INTERFACE":
                lines.append(f"interface {name} {{")
                for field in t.get("fields", []):
                    field_type = self._type_ref_to_string(field["type"])
                    lines.append(f"  {field['name']}: {field_type}")
                lines.append("}")
                lines.append("")

            elif kind == "UNION":
                possible = t.get("possibleTypes", [])
                possible_names = [p["name"] for p in possible]
                lines.append(f"union {name} = {' | '.join(possible_names)}")
                lines.append("")

            elif kind == "SCALAR":
                if name not in ("String", "Int", "Float", "Boolean", "ID"):
                    lines.append(f"scalar {name}")
                    lines.append("")

        return "\n".join(lines)

    def _type_ref_to_string(self, type_ref: dict[str, Any]) -> str:
        """Convert type reference to string representation."""
        kind = type_ref.get("kind")
        name = type_ref.get("name")
        of_type = type_ref.get("ofType")

        if kind == "NON_NULL":
            return f"{self._type_ref_to_string(of_type)}!"
        elif kind == "LIST":
            return f"[{self._type_ref_to_string(of_type)}]"
        else:
            return name or "Unknown"

    def _format_args(self, args: list[dict[str, Any]]) -> str:
        """Format field arguments."""
        if not args:
            return ""

        arg_strs = []
        for arg in args:
            arg_type = self._type_ref_to_string(arg["type"])
            default = arg.get("defaultValue")
            if default:
                arg_strs.append(f"{arg['name']}: {arg_type} = {default}")
            else:
                arg_strs.append(f"{arg['name']}: {arg_type}")

        return f"({', '.join(arg_strs)})"


# =============================================================================
# Main Fetcher Logic
# =============================================================================


def load_sources() -> dict[str, Any]:
    """Load sources configuration."""
    with open(SOURCES_FILE) as f:
        return yaml.safe_load(f)


def load_index() -> dict[str, Any]:
    """Load fetch index."""
    if not INDEX_FILE.exists():
        return {"status": "pending", "last_run": None, "sources": {}}
    with open(INDEX_FILE) as f:
        return yaml.safe_load(f) or {
            "status": "pending",
            "last_run": None,
            "sources": {},
        }


def save_index(index: dict[str, Any]) -> None:
    """Save fetch index."""
    with open(INDEX_FILE, "w") as f:
        yaml.dump(index, f, default_flow_style=False, sort_keys=False)


def get_fetcher(source_type: str, config: dict[str, Any]) -> BaseFetcher | None:
    """Get appropriate fetcher for source type. Returns None for manual types."""
    fetchers = {
        "html": HTMLFetcher,
        "json": JSONFetcher,
        "graphql_schema": GraphQLSchemaFetcher,
        "github_readme": GitHubReadmeFetcher,
    }
    if source_type == "manual":
        return None  # Manual sources are skipped
    fetcher_class = fetchers.get(source_type, HTMLFetcher)
    return fetcher_class(config)


def validate_output_path(output_path: Path) -> bool:
    """Ensure output path is within allowed directory."""
    resolved = output_path.resolve()
    return str(resolved).startswith(ALLOWED_OUTPUT_PREFIX)


def fetch_source(
    source: dict[str, Any],
    config: dict[str, Any],
    index: dict[str, Any],
    dry_run: bool = False,
) -> bool:
    """Fetch a single source and update index."""
    source_id = source["id"]
    source_type = source.get("type", "html")
    output_path = OUTPUT_DIR / source["output"]

    # Skip manual sources
    if source_type == "manual":
        print(f"  ⏭️  Skipping {source['name']} (manual type)")
        index["sources"][source_id] = {
            "status": "skipped",
            "reason": "manual type - requires manual download",
            "note": source.get("note", ""),
        }
        return True

    # Validate output path
    if not validate_output_path(output_path):
        print(f"  ❌ SECURITY: Output path outside allowed directory: {output_path}")
        return False

    if dry_run:
        print(f"  [DRY-RUN] Would fetch {source['name']} -> {source['output']}")
        return True

    print(f"  Fetching {source['name']}...")

    fetcher = get_fetcher(source_type, config)
    if fetcher is None:
        print(f"  ⏭️  No fetcher for type: {source_type}")
        return True

    content, error = fetcher.fetch(source)

    now = datetime.now(timezone.utc).isoformat()

    if error:
        print(f"  ❌ Error: {error}")
        # Keep old file if exists, just update status
        index["sources"][source_id] = {
            "status": "failed",
            "last_fetched": index.get("sources", {})
            .get(source_id, {})
            .get("last_fetched"),
            "error": error,
            "attempted": now,
        }
        return False

    # Create output directory if needed
    output_path.parent.mkdir(parents=True, exist_ok=True)

    # Calculate checksum
    checksum = hashlib.sha256(content.encode()).hexdigest()

    # Write content
    with open(output_path, "w") as f:
        f.write(content)

    # Update index
    index["sources"][source_id] = {
        "status": "success",
        "last_fetched": now,
        "file_size": len(content),
        "checksum": checksum,
        "error": None,
    }

    print(f"  ✅ Saved to {source['output']} ({len(content)} bytes)")
    return True


def main():
    parser = argparse.ArgumentParser(description="Fetch external documentation sources")
    parser.add_argument("--category", help="Fetch only sources in this category")
    parser.add_argument("--id", help="Fetch only source with this ID")
    parser.add_argument(
        "--dry-run", action="store_true", help="Show what would be fetched"
    )
    parser.add_argument(
        "--force", action="store_true", help="Force re-fetch even if recent"
    )
    args = parser.parse_args()

    print("=" * 60)
    print("Revenge External Documentation Fetcher")
    print("=" * 60)

    # Load configuration
    sources_config = load_sources()
    fetch_config = sources_config.get("fetch_config", {})
    index = load_index()

    # Collect sources to fetch
    sources_to_fetch = []
    for category, sources in sources_config.get("sources", {}).items():
        if args.category and category != args.category:
            continue

        for source in sources:
            if args.id and source["id"] != args.id:
                continue
            source["_category"] = category
            sources_to_fetch.append(source)

    if not sources_to_fetch:
        print("No sources matched the criteria.")
        return 1

    print(f"Found {len(sources_to_fetch)} sources to fetch")
    print()

    # Fetch sources
    success_count = 0
    fail_count = 0
    delay = fetch_config.get("delay_between_requests", 2)

    for i, source in enumerate(sources_to_fetch):
        category = source["_category"]
        print(f"[{i + 1}/{len(sources_to_fetch)}] {category}/{source['id']}")

        if fetch_source(source, fetch_config, index, args.dry_run):
            success_count += 1
        else:
            fail_count += 1

        # Delay between requests (not on last one, not on dry-run)
        if i < len(sources_to_fetch) - 1 and not args.dry_run:
            time.sleep(delay)

    # Update index
    if not args.dry_run:
        index["status"] = "completed" if fail_count == 0 else "partial"
        index["last_run"] = datetime.now(timezone.utc).isoformat()
        save_index(index)

    print()
    print("=" * 60)
    print(f"Complete: {success_count} success, {fail_count} failed")
    print("=" * 60)

    return 0 if fail_count == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
