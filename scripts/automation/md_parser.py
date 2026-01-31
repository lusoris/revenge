#!/usr/bin/env python3
"""Markdown Parser - Extract structured data from design docs.

This parser:
1. Extracts metadata from existing markdown files
2. Parses status tables
3. Identifies document type
4. Creates YAML data templates
5. Preserves existing content sections

Author: Automation System
Created: 2026-01-31
"""

import re
import sys
from pathlib import Path

import yaml


class MarkdownParser:
    """Parse design documentation markdown files."""

    def __init__(self, repo_root: Path):
        """Initialize parser with repository root."""
        self.repo_root = repo_root
        self.docs_dir = repo_root / "docs" / "dev" / "design"
        self.sources_dir = repo_root / "docs" / "dev" / "sources"

        # Load SOURCES.yaml for source resolution
        self.sources_index = self._load_sources_index()

    def _load_sources_index(self) -> dict:
        """Load SOURCES.yaml to resolve source IDs to URLs."""
        sources_file = self.sources_dir / "SOURCES.yaml"
        if not sources_file.exists():
            print(f"⚠️  Warning: {sources_file} not found")
            return {}

        with open(sources_file) as f:
            sources_yaml = yaml.safe_load(f)

        # Build index: source_id -> source details
        # Structure: sources -> category -> list of sources
        index = {}
        for category_sources in sources_yaml.get("sources", {}).values():
            if isinstance(category_sources, list):
                for source in category_sources:
                    source_id = source.get("id")
                    if source_id:
                        index[source_id] = {
                            "name": source.get("name", source_id),
                            "url": source.get("url", ""),
                        }
        return index

    def parse_file(self, md_file: Path) -> dict:
        """Parse a markdown file and extract structured data.

        Args:
            md_file: Path to markdown file

        Returns:
            Dict with extracted data
        """
        with open(md_file) as f:
            content = f.read()

        data = {}

        # Extract title
        title_match = re.search(r"^#\s+(.+)$", content, re.MULTILINE)
        if title_match:
            data["doc_title"] = title_match.group(1).strip()

        # Extract sources comment
        sources_match = re.search(
            r"<!--\s*SOURCES:\s*(.+?)\s*-->", content, re.IGNORECASE
        )
        if sources_match:
            source_ids = [s.strip() for s in sources_match.group(1).split(",")]
            data["source_ids"] = source_ids

        # Extract design refs comment
        design_match = re.search(
            r"<!--\s*DESIGN:\s*(.+?)\s*-->", content, re.IGNORECASE
        )
        if design_match:
            design_ids = [d.strip() for d in design_match.group(1).split(",")]
            data["design_ids"] = design_ids

        # Extract status table
        status = self._parse_status_table(content)
        if status:
            data.update(status)

        # Determine document category
        doc_category = self._determine_category(md_file, content)
        data["doc_category"] = doc_category

        # Extract content sections
        sections = self._extract_sections(content)
        data["sections"] = sections

        # Add file metadata
        data["source_file"] = str(md_file.relative_to(self.repo_root))
        data["created_date"] = "2026-01-31"  # Default, can be overridden

        return data

    def _parse_status_table(self, content: str) -> dict | None:
        """Parse 7-dimension status table.

        Returns:
            Dict with status_* fields or None if table not found
        """
        # Find status table (7 rows after header)
        pattern = r"\|\s*Dimension\s*\|\s*Status\s*\|.*?\n" r"\|[-\s|]+\n" r"((?:\|.+\n){7})"

        match = re.search(pattern, content, re.MULTILINE | re.DOTALL)
        if not match:
            return None

        rows = match.group(1).strip().split("\n")
        status_data = {}

        dimensions = [
            "design",
            "sources",
            "instructions",
            "code",
            "linting",
            "unit_testing",
            "integration_testing",
        ]

        for i, row in enumerate(rows):
            if i >= len(dimensions):
                break

            # Parse row: | Dimension | 🔴 |
            parts = [p.strip() for p in row.split("|")]
            if len(parts) >= 3:
                emoji_status = parts[2].strip()
                status_data[f"status_{dimensions[i]}"] = emoji_status
                status_data[f"status_{dimensions[i]}_notes"] = "-"

        # Determine overall status (use design status as proxy)
        if "status_design" in status_data:
            status_data["overall_status"] = status_data["status_design"]

        return status_data

    def _determine_category(self, md_file: Path, content: str) -> str:
        """Determine document category based on file path and content.

        Returns:
            One of: feature, service, integration, architecture, operations, technical
        """
        path_str = str(md_file)

        # Check path-based categories
        if "/features/" in path_str:
            return "feature"
        elif "/services/" in path_str:
            return "service"
        elif "/integrations/" in path_str:
            return "integration"
        elif "/architecture/" in path_str:
            return "architecture"
        elif "/operations/" in path_str:
            return "operations"
        elif "/technical/" in path_str:
            return "technical"
        elif "/patterns/" in path_str:
            return "pattern"
        elif "/research/" in path_str:
            return "research"

        return "other"

    def _resolve_design_path(self, design_id: str) -> str:
        """Resolve design ID to relative path.

        Args:
            design_id: Design doc identifier (e.g., "01_ARCHITECTURE", "operations")

        Returns:
            Relative path to design doc
        """
        # Common patterns for design doc IDs
        if design_id.upper() in ["OPERATIONS", "ARCHITECTURE", "SERVICES", "FEATURES"]:
            return f"{design_id.lower()}/INDEX.md"
        elif re.match(r"^\d{2}_[A-Z_]+$", design_id):
            # Numbered docs like 01_ARCHITECTURE
            return f"architecture/{design_id}.md"
        else:
            # Default: assume it's a relative path
            return f"{design_id}.md"

    def _extract_sections(self, content: str) -> list[dict]:
        """Extract major content sections from markdown.

        Returns:
            List of dicts with section name and content
        """
        sections = []

        # Split by ## headers (skip # title)
        section_pattern = r"^##\s+(.+?)$"
        matches = list(re.finditer(section_pattern, content, re.MULTILINE))

        for i, match in enumerate(matches):
            section_name = match.group(1).strip()
            start_pos = match.end()

            # Find end of section (next ## or end of file)
            end_pos = matches[i + 1].start() if i + 1 < len(matches) else len(content)

            section_content = content[start_pos:end_pos].strip()

            # Skip auto-generated sections
            if section_name in [
                "Table of Contents",
                "Related Design Docs",
                "Sources & Cross-References",
            ]:
                continue

            sections.append({"name": section_name, "content": section_content})

        return sections

    def to_yaml(self, data: dict, template_type: str = "basic") -> str:
        """Convert extracted data to YAML format.

        Args:
            data: Extracted data dict
            template_type: Type of YAML template to generate

        Returns:
            YAML string with data and placeholders
        """
        # Build clean YAML data (no comment keys)
        yaml_data = {
            "doc_title": data.get("doc_title", "PLACEHOLDER"),
            "doc_category": data.get("doc_category", "other"),
            "created_date": data.get("created_date", "2026-01-31"),
        }

        # Add status fields
        for key in [
            "overall_status",
            "status_design",
            "status_design_notes",
            "status_sources",
            "status_sources_notes",
            "status_instructions",
            "status_instructions_notes",
            "status_code",
            "status_code_notes",
            "status_linting",
            "status_linting_notes",
            "status_unit_testing",
            "status_unit_testing_notes",
            "status_integration_testing",
            "status_integration_testing_notes",
        ]:
            if key in data:
                yaml_data[key] = data[key]

        # Add placeholder fields based on category
        yaml_data["technical_summary"] = "PLACEHOLDER: Brief technical summary"
        yaml_data["wiki_tagline"] = "PLACEHOLDER: User-friendly tagline"
        yaml_data["wiki_overview"] = "PLACEHOLDER: User-friendly overview"

        # Resolve sources
        if "source_ids" in data:
            sources = []
            for source_id in data["source_ids"]:
                # Try exact match first
                if source_id in self.sources_index:
                    sources.append(
                        {
                            "name": self.sources_index[source_id]["name"],
                            "url": self.sources_index[source_id]["url"],
                            "note": f"Auto-resolved from {source_id}",
                        }
                    )
                else:
                    # Try partial match (e.g., "pgx" matches "go-pgx" or contains "pgx")
                    matched = False
                    for full_id, source_info in self.sources_index.items():
                        if (
                            source_id in full_id
                            or source_id.lower() in source_info["name"].lower()
                        ):
                            sources.append(
                                {
                                    "name": source_info["name"],
                                    "url": source_info["url"],
                                    "note": f"Auto-resolved from {source_id} → {full_id}",
                                }
                            )
                            matched = True
                            break

                    if not matched:
                        sources.append(
                            {
                                "name": source_id,
                                "url": "PLACEHOLDER_URL",
                                "note": f"⚠️ Source '{source_id}' not found - needs manual resolution",
                            }
                        )
            if sources:
                yaml_data["sources"] = sources

        # Resolve design refs
        if "design_ids" in data:
            design_refs = []
            for design_id in data["design_ids"]:
                # Try to find the design doc file
                design_path = self._resolve_design_path(design_id)
                design_refs.append({"title": design_id, "path": design_path})
            if design_refs:
                yaml_data["design_refs"] = design_refs

        # Convert to YAML
        yaml_str = yaml.dump(
            yaml_data, default_flow_style=False, allow_unicode=True, sort_keys=False
        )

        # Add header comment
        header = f"""# Auto-generated YAML from {data.get('source_file', 'markdown')}
# Migration Date: {data.get('created_date', '2026-01-31')}
#
# TODO: Complete the PLACEHOLDER fields below
# TODO: Review auto-resolved sources and design_refs
# TODO: Add feature/service/integration-specific fields as needed

"""
        return header + yaml_str


def main():
    """Main entry point - Test parser with VERSIONING.md."""
    repo_root = Path(__file__).parent.parent.parent

    parser = MarkdownParser(repo_root)

    # Test with VERSIONING.md
    test_file = repo_root / "docs" / "dev" / "design" / "operations" / "VERSIONING.md"

    if not test_file.exists():
        print(f"❌ Test file not found: {test_file}")
        sys.exit(1)

    print(f"📄 Parsing: {test_file.relative_to(repo_root)}\n")

    # Parse file
    data = parser.parse_file(test_file)

    # Print extracted data
    print("=" * 70)
    print("EXTRACTED DATA")
    print("=" * 70)
    print(f"Title: {data.get('doc_title')}")
    print(f"Category: {data.get('doc_category')}")
    print(f"Source IDs: {data.get('source_ids', [])}")
    print(f"Design IDs: {data.get('design_ids', [])}")
    print(f"Status Design: {data.get('status_design', 'N/A')}")
    print(f"Sections: {len(data.get('sections', []))}")
    print()

    # Generate YAML
    yaml_output = parser.to_yaml(data)
    print("=" * 70)
    print("GENERATED YAML")
    print("=" * 70)
    print(yaml_output)


if __name__ == "__main__":
    main()
