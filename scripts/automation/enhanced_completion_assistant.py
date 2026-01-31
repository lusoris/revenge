#!/usr/bin/env python3
"""Enhanced YAML Completion Assistant - Extract category-specific fields from markdown.

This tool performs deeper content extraction beyond basic summaries:
- For integrations: Extracts integration_name, external_service, api_base_url, auth_method
- For services: Extracts service_name, package_path, fx_module
- For features: Extracts content_types from content

Author: Automation System
Created: 2026-01-31
"""

import re
import sys
from pathlib import Path

import yaml


class EnhancedCompletionAssistant:
    """Enhanced completion assistant for category-specific fields."""

    def __init__(self, repo_root: Path):
        """Initialize assistant.

        Args:
            repo_root: Repository root path
        """
        self.repo_root = repo_root
        self.data_dir = repo_root / "data"
        self.docs_dir = repo_root / "docs" / "dev" / "design"

    def get_original_markdown(self, yaml_file: Path) -> Path | None:
        """Get path to original markdown file.

        Args:
            yaml_file: Path to YAML data file

        Returns:
            Path to original markdown file, or None if not found
        """
        # Convert data path to docs path
        rel_path = yaml_file.relative_to(self.data_dir)
        md_file = self.docs_dir / rel_path.with_suffix(".md")

        return md_file if md_file.exists() else None

    def extract_integration_fields(self, md_content: str, doc_title: str) -> dict:
        """Extract integration-specific fields from markdown.

        Args:
            md_content: Markdown content
            doc_title: Document title

        Returns:
            Dict with integration fields
        """
        fields = {}

        # Extract integration_name from title
        # "TMDb Integration" -> "TMDb"
        # "Radarr" -> "Radarr"
        integration_name = doc_title.replace(" Integration", "").strip()
        if integration_name:
            fields["integration_name"] = integration_name

        # Generate integration_id (lowercase, no spaces)
        # "TMDb" -> "tmdb"
        # "The Movie Database" -> "the_movie_database"
        if integration_name:
            integration_id = integration_name.lower().replace(" ", "_")
            integration_id = re.sub(r"[^a-z0-9_]", "", integration_id)
            fields["integration_id"] = integration_id

        # Extract external_service (usually same as integration_name or in first paragraph)
        if integration_name:
            fields["external_service"] = integration_name

        # Try to extract API base URL from markdown
        # Look for patterns like:
        # - https://api.themoviedb.org/3
        # - Base URL: https://...
        # - API endpoint: https://...
        api_url_patterns = [
            r"https://api\.[a-z0-9\-\.]+/[a-z0-9/]*",
            r"https://[a-z0-9\-\.]+/api/[a-z0-9/]*",
            r"(?:Base URL|API endpoint|API URL):\s*(https://[^\s\)]+)",
        ]

        for pattern in api_url_patterns:
            match = re.search(pattern, md_content, re.IGNORECASE)
            if match:
                url = match.group(1) if match.lastindex else match.group(0)
                # Clean up URL (remove trailing slashes, punctuation)
                url = re.sub(r"[,\.\)]+$", "", url).rstrip("/")
                if url.startswith("http"):
                    fields["api_base_url"] = url
                    break

        # Try to detect auth method from content
        auth_keywords = {
            "api_key": ["API key", "API token", "access key"],
            "oauth": ["OAuth", "OAuth2", "authorization code"],
            "bearer": ["bearer token", "JWT"],
            "basic": ["basic auth", "username and password"],
            "none": ["no authentication", "public API"],
        }

        content_lower = md_content.lower()
        for auth_method, keywords in auth_keywords.items():
            if any(keyword.lower() in content_lower for keyword in keywords):
                fields["auth_method"] = auth_method
                break

        return fields

    def extract_service_fields(self, md_content: str, doc_title: str) -> dict:
        """Extract service-specific fields from markdown.

        Args:
            md_content: Markdown content
            doc_title: Document title

        Returns:
            Dict with service fields
        """
        fields = {}

        # Extract service_name from title
        # "Authentication Service" -> "Authentication Service"
        # "User Service" -> "User Service"
        service_name = doc_title.replace(" Service", "").strip() + " Service"
        if not service_name.endswith(" Service"):
            service_name += " Service"
        fields["service_name"] = service_name

        # Try to extract package_path from markdown
        # Look for patterns like:
        # - `internal/service/auth`
        # - Package: internal/service/user
        # - Located in: internal/service/metadata
        package_patterns = [
            r"`(internal/service/[a-z0-9_/]+)`",
            r"(?:Package|Located in|Path):\s*`?(internal/service/[a-z0-9_/]+)`?",
        ]

        for pattern in package_patterns:
            match = re.search(pattern, md_content, re.IGNORECASE)
            if match:
                fields["package_path"] = match.group(1)
                break

        # If no package_path found, generate from service name
        if "package_path" not in fields:
            # "Authentication Service" -> "auth"
            # "User Service" -> "user"
            service_slug = (
                service_name.replace(" Service", "").lower().replace(" ", "_")
            )
            fields["package_path"] = f"internal/service/{service_slug}"

        # Generate fx_module name
        # "internal/service/auth" -> "AuthModule"
        # "internal/service/user" -> "UserModule"
        if "package_path" in fields:
            parts = fields["package_path"].split("/")
            if parts:
                module_name = parts[-1].capitalize() + "Module"
                fields["fx_module"] = module_name

        return fields

    def extract_feature_fields(self, md_content: str, doc_title: str) -> dict:
        """Extract feature-specific fields from markdown.

        Args:
            md_content: Markdown content
            doc_title: Document title

        Returns:
            Dict with feature fields
        """
        fields = {}

        # Try to extract content_types from markdown
        # Look for patterns like:
        # - Content Types: Movies, Collections
        # - Handles: TV Shows, Seasons, Episodes
        # - Manages: Books, Authors, Series
        content_type_patterns = [
            r"Content Types?:\s*([A-Z][a-zA-Z, ]+)",
            r"(?:Handles|Manages|Supports):\s*([A-Z][a-zA-Z, ]+)",
        ]

        for pattern in content_type_patterns:
            match = re.search(pattern, md_content)
            if match:
                # Split by comma and clean up
                types_str = match.group(1)
                types = [t.strip() for t in types_str.split(",")]
                # Remove trailing punctuation
                types = [re.sub(r"[,\.\)]+$", "", t) for t in types]
                if types:
                    fields["content_types"] = types
                    break

        # If no content_types found, try to infer from title
        if "content_types" not in fields:
            title_lower = doc_title.lower()
            if "movie" in title_lower:
                fields["content_types"] = ["Movies", "Collections"]
            elif "tv" in title_lower or "show" in title_lower:
                fields["content_types"] = ["TV Shows", "Seasons", "Episodes"]
            elif "music" in title_lower:
                fields["content_types"] = ["Artists", "Albums", "Tracks"]
            elif "book" in title_lower:
                fields["content_types"] = ["Books", "Authors", "Series"]
            elif "comic" in title_lower:
                fields["content_types"] = ["Comics", "Issues", "Series"]
            elif "audiobook" in title_lower:
                fields["content_types"] = ["Audiobooks", "Authors", "Series"]
            elif "podcast" in title_lower:
                fields["content_types"] = ["Podcasts", "Episodes"]
            elif "photo" in title_lower:
                fields["content_types"] = ["Albums", "Photos"]
            elif "adult" in title_lower or "qar" in title_lower:
                fields["content_types"] = ["Scenes", "Performers", "Studios"]

        return fields

    def suggest_completions(
        self, yaml_file: Path, category: str
    ) -> dict | None:
        """Suggest completions for a YAML file.

        Args:
            yaml_file: Path to YAML file
            category: Document category (feature/service/integration)

        Returns:
            Dict of suggested fields, or None if no suggestions
        """
        # Load YAML
        with open(yaml_file) as f:
            data = yaml.safe_load(f)

        # Get original markdown
        md_file = self.get_original_markdown(yaml_file)
        if not md_file or not md_file.exists():
            return None

        # Read markdown content
        with open(md_file) as f:
            md_content = f.read()

        doc_title = data.get("doc_title", "")

        # Extract category-specific fields
        if category == "integration":
            return self.extract_integration_fields(md_content, doc_title)
        elif category == "service":
            return self.extract_service_fields(md_content, doc_title)
        elif category == "feature":
            return self.extract_feature_fields(md_content, doc_title)

        return None

    def apply_suggestions(
        self, yaml_file: Path, suggestions: dict, dry_run: bool = False
    ) -> bool:
        """Apply suggestions to YAML file.

        Args:
            yaml_file: Path to YAML file
            suggestions: Dict of field suggestions
            dry_run: If True, don't write changes

        Returns:
            True if changes applied, False otherwise
        """
        if not suggestions:
            return False

        # Load existing YAML
        with open(yaml_file) as f:
            data = yaml.safe_load(f)

        # Apply suggestions (only if field doesn't exist or is placeholder)
        updated = False
        for field, value in suggestions.items():
            if field not in data or str(data[field]).startswith("PLACEHOLDER"):
                data[field] = value
                updated = True

        if not updated:
            return False

        if dry_run:
            return True

        # Write updated YAML
        with open(yaml_file, "w") as f:
            yaml.safe_dump(
                data, f, default_flow_style=False, allow_unicode=True, sort_keys=False
            )

        return True

    def process_category(
        self, category: str, auto_apply: bool = False, dry_run: bool = False
    ) -> dict:
        """Process all files in a category.

        Args:
            category: Category to process (feature/service/integration)
            auto_apply: If True, auto-apply suggestions
            dry_run: If True, don't write changes

        Returns:
            Statistics dict
        """
        # Find YAML files for category
        category_dirs = {
            "feature": ["features"],
            "service": ["services"],
            "integration": ["integrations"],
        }

        yaml_files = []
        for dir_name in category_dirs.get(category, []):
            category_path = self.data_dir / dir_name
            if category_path.exists():
                yaml_files.extend(category_path.glob("**/*.yaml"))

        if not yaml_files:
            print(f"\n⚠️  No YAML files found for category: {category}")
            return {"total": 0, "suggested": 0, "applied": 0}

        print(f"\n{'='*70}")
        print(f"ENHANCED COMPLETION - {category.upper()}")
        print(f"{'='*70}\n")
        print(f"Found {len(yaml_files)} files")
        print(f"Mode: {'AUTO-APPLY' if auto_apply else 'DRY-RUN'}\n")

        stats = {"total": len(yaml_files), "suggested": 0, "applied": 0}

        for i, yaml_file in enumerate(sorted(yaml_files), 1):
            rel_path = yaml_file.relative_to(self.data_dir)
            print(f"[{i}/{len(yaml_files)}]")
            print(f"\n📝 {rel_path}")

            suggestions = self.suggest_completions(yaml_file, category)

            if not suggestions:
                print("   ⊘ No suggestions")
                continue

            stats["suggested"] += 1

            # Show suggestions
            print("   Suggestions:")
            for field, value in suggestions.items():
                # Truncate long values
                value_str = str(value)
                if len(value_str) > 60:
                    value_str = value_str[:57] + "..."
                print(f"   ✓ {field}: {value_str}")

            # Apply if auto mode
            if auto_apply:
                applied = self.apply_suggestions(yaml_file, suggestions, dry_run)
                if applied:
                    stats["applied"] += 1
                    if not dry_run:
                        print(f"   ✅ Updated {yaml_file.name}")
                    else:
                        print(f"   ℹ️  Would update {yaml_file.name}")
                else:
                    print("   ⊘ No changes needed")

            print()

        # Summary
        print(f"{'='*70}")
        print("SUMMARY")
        print(f"{'='*70}")
        print(f"Total files: {stats['total']}")
        print(f"Files with suggestions: {stats['suggested']}")
        print(f"Files updated: {stats['applied']}")
        print(f"{'='*70}\n")

        return stats


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent.parent

    # Parse arguments
    args = sys.argv[1:]

    category = None
    auto_apply = False
    dry_run = True

    for arg in args:
        if arg == "--feature":
            category = "feature"
        elif arg == "--service":
            category = "service"
        elif arg == "--integration":
            category = "integration"
        elif arg == "--auto":
            auto_apply = True
            dry_run = False
        elif arg == "--dry-run":
            dry_run = True

    if not category:
        print("Usage: python enhanced_completion_assistant.py --<category> [--auto]")
        print()
        print("Categories:")
        print("  --feature      Process feature modules")
        print("  --service      Process backend services")
        print("  --integration  Process integrations")
        print()
        print("Options:")
        print("  --auto         Auto-apply suggestions (default: dry-run)")
        print("  --dry-run      Show what would be done (default)")
        print()
        print("Examples:")
        print("  python enhanced_completion_assistant.py --integration --dry-run")
        print("  python enhanced_completion_assistant.py --service --auto")
        sys.exit(1)

    # Initialize assistant
    assistant = EnhancedCompletionAssistant(repo_root)

    # Process category
    stats = assistant.process_category(category, auto_apply, dry_run)

    # Exit code based on results
    sys.exit(0 if stats["applied"] > 0 or stats["suggested"] > 0 else 1)


if __name__ == "__main__":
    main()
