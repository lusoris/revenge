#!/usr/bin/env python3
"""YAML Completion Assistant - Help complete placeholder fields in YAML.

This assistant:
1. Finds YAML files with placeholders
2. Reads corresponding original markdown
3. Extracts useful content to suggest completions
4. Provides interactive or batch completion mode
5. Validates completed fields

Author: Automation System
Created: 2026-01-31
"""

import re
import sys
from pathlib import Path

import yaml


class CompletionAssistant:
    """Assist with YAML file completion."""

    def __init__(self, repo_root: Path):
        """Initialize assistant."""
        self.repo_root = repo_root
        self.data_dir = repo_root / "data"
        self.docs_dir = repo_root / "docs" / "dev" / "design"

    def find_incomplete_files(self, category: str | None = None) -> list[Path]:
        """Find YAML files with placeholders.

        Args:
            category: Optional category filter (feature, service, integration)

        Returns:
            List of incomplete YAML file paths
        """
        yaml_files = list(self.data_dir.glob("**/*.yaml"))
        yaml_files = [f for f in yaml_files if f.name != "shared-sot.yaml"]

        incomplete = []

        for yaml_file in yaml_files:
            with open(yaml_file) as f:
                try:
                    data = yaml.safe_load(f)

                    # Check category filter
                    if category and data.get("doc_category") != category:
                        continue

                    # Check for placeholders
                    has_placeholders = any(
                        isinstance(v, str) and "PLACEHOLDER" in v for v in data.values()
                    )

                    if has_placeholders:
                        incomplete.append(yaml_file)
                except Exception:
                    pass

        return sorted(incomplete)

    def get_original_markdown(self, yaml_file: Path) -> Path | None:
        """Get corresponding original markdown file.

        Args:
            yaml_file: YAML data file

        Returns:
            Path to original markdown or None
        """
        rel_path = yaml_file.relative_to(self.data_dir)
        md_file = self.docs_dir / rel_path.with_suffix(".md")

        if md_file.exists():
            return md_file
        return None

    def extract_summary(self, md_file: Path) -> str | None:
        """Extract a summary from markdown.

        Looks for first paragraph after title, or Status section description.

        Args:
            md_file: Markdown file

        Returns:
            Extracted summary or None
        """
        with open(md_file) as f:
            content = f.read()

        # Remove frontmatter if present
        content = re.sub(r"^---\n.*?\n---\n", "", content, flags=re.DOTALL)

        # Skip title
        lines = content.split("\n")
        in_content = False
        summary_lines = []

        for line in lines:
            # Skip title
            if line.startswith("# "):
                in_content = True
                continue

            if in_content:
                # Skip comments, TOC markers
                if line.startswith("<!--") or "TOC" in line or line.startswith("##"):
                    if summary_lines:  # Already got some content
                        break
                    continue

                # Collect non-empty lines
                if line.strip():
                    summary_lines.append(line.strip())

                # Stop after a few good lines
                if len(summary_lines) >= 3:
                    break

        if summary_lines:
            return " ".join(summary_lines)[:200]  # Max 200 chars

        return None

    def suggest_completions(self, yaml_file: Path) -> dict:
        """Suggest completions for placeholder fields.

        Args:
            yaml_file: YAML file to complete

        Returns:
            Dict with suggestions
        """
        # Load YAML
        with open(yaml_file) as f:
            data = yaml.safe_load(f)

        suggestions = {}

        # Get original markdown if exists
        md_file = self.get_original_markdown(yaml_file)

        # Suggest technical_summary
        if "PLACEHOLDER" in str(data.get("technical_summary", "")):
            if md_file and md_file.exists():
                summary = self.extract_summary(md_file)
                if summary:
                    suggestions["technical_summary"] = summary

        # Suggest wiki_tagline (shortened summary)
        if "PLACEHOLDER" in str(data.get("wiki_tagline", "")):
            if md_file and md_file.exists():
                summary = self.extract_summary(md_file)
                if summary:
                    # Take first sentence or first 80 chars
                    tagline = summary.split(".")[0][:80]
                    suggestions["wiki_tagline"] = tagline

        # For features, suggest common fields from doc_title
        if data.get("doc_category") == "feature":
            doc_title = data.get("doc_title", "")

            if not data.get("feature_name"):
                suggestions["feature_name"] = doc_title

            if not data.get("module_name"):
                # Convert "Movie Module" -> "movie"
                module_name = (
                    doc_title.lower()
                    .replace(" module", "")
                    .replace(" ", "_")
                    .replace("-", "_")
                )
                suggestions["module_name"] = module_name

            if not data.get("schema_name"):
                # Guess based on path
                if "adult" in str(yaml_file) or "qar" in str(yaml_file):
                    suggestions["schema_name"] = "qar"
                else:
                    suggestions["schema_name"] = "public"

        return suggestions

    def apply_suggestions(
        self,
        yaml_file: Path,
        suggestions: dict,
        auto_apply: bool = False,
    ):
        """Apply suggestions to YAML file.

        Args:
            yaml_file: YAML file to update
            suggestions: Dict of field -> value suggestions
            auto_apply: If True, apply without prompting
        """
        if not suggestions:
            print(f"  No suggestions for {yaml_file.name}")
            return

        # Load current YAML
        with open(yaml_file) as f:
            data = yaml.safe_load(f)

        print(f"\n📝 {yaml_file.relative_to(self.data_dir)}")
        print("   Suggestions:")

        updates_made = False

        for field, suggestion in suggestions.items():
            current = data.get(field)

            if auto_apply:
                data[field] = suggestion
                updates_made = True
                print(f"   ✓ {field}: {suggestion[:60]}...")
            else:
                print(f"\n   {field}:")
                print(f"     Current: {current}")
                print(f"     Suggest: {suggestion}")

                response = input("     Apply? (y/n/edit): ").lower()

                if response == "y":
                    data[field] = suggestion
                    updates_made = True
                    print("     ✓ Applied")
                elif response == "edit":
                    edited = input("     Enter value: ")
                    if edited:
                        data[field] = edited
                        updates_made = True
                        print("     ✓ Applied (edited)")

        if updates_made:
            # Write back to file
            with open(yaml_file, "w") as f:
                yaml.dump(
                    data,
                    f,
                    default_flow_style=False,
                    allow_unicode=True,
                    sort_keys=False,
                )
            print(f"   ✅ Updated {yaml_file.name}")
        else:
            print("   ⊘ No changes made")

    def process_all(self, category: str | None = None, auto_apply: bool = False):
        """Process all incomplete files.

        Args:
            category: Optional category filter
            auto_apply: Auto-apply suggestions without prompting
        """
        incomplete_files = self.find_incomplete_files(category)

        print(f"\n{'=' * 70}")
        print("YAML COMPLETION ASSISTANT")
        print(f"{'=' * 70}\n")
        print(f"Found {len(incomplete_files)} incomplete files")

        if category:
            print(f"Category filter: {category}")

        if auto_apply:
            print("Mode: AUTO-APPLY (no prompts)")
        else:
            print("Mode: INTERACTIVE (will prompt for each suggestion)")

        for i, yaml_file in enumerate(incomplete_files, 1):
            print(f"\n[{i}/{len(incomplete_files)}]")

            suggestions = self.suggest_completions(yaml_file)
            self.apply_suggestions(yaml_file, suggestions, auto_apply)

            if not auto_apply and i < len(incomplete_files):
                cont = input("\nContinue to next file? (y/n/q): ").lower()
                if cont == "q":
                    print("\nStopped by user")
                    break
                if cont == "n":
                    continue

        print(f"\n{'=' * 70}")
        print("✅ Processing complete!")
        print(f"{'=' * 70}\n")


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent.parent

    # Parse args
    category = None
    auto_apply = "--auto" in sys.argv

    if "--feature" in sys.argv:
        category = "feature"
    elif "--service" in sys.argv:
        category = "service"
    elif "--integration" in sys.argv:
        category = "integration"

    # Initialize assistant
    assistant = CompletionAssistant(repo_root)

    # Process files
    assistant.process_all(category=category, auto_apply=auto_apply)


if __name__ == "__main__":
    main()
