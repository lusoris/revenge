#!/usr/bin/env python3
"""SOURCE_OF_TRUTH.md parser - Extract structured data from master document.

This parser reads docs/dev/design/00_SOURCE_OF_TRUTH.md and extracts:
- Content modules (module inventory with status)
- Backend services (service inventory with fx modules)
- Infrastructure components (versions and purpose)
- Go dependencies (all packages with versions and notes)
- Development tools (tools with versions for config sync)
- Core design principles
- API namespaces
- Configuration keys
- QAR terminology mappings

Output: data/shared-sot.yaml (used by all doc generation)

Author: Automation System
Created: 2026-01-31
"""

import re
import sys
from pathlib import Path
from typing import Any

import yaml


class SOTParser:
    """Parse SOURCE_OF_TRUTH.md into structured YAML data."""

    def __init__(self, sot_path: Path):
        """Initialize parser with path to SOURCE_OF_TRUTH.md."""
        self.sot_path = sot_path
        self.content = self.sot_path.read_text()
        self.data: dict[str, Any] = {}

    def parse(self) -> dict[str, Any]:
        """Parse all sections from SOURCE_OF_TRUTH.md."""
        print(f"📖 Parsing {self.sot_path}")

        # Extract metadata from frontmatter
        self.data["metadata"] = self._parse_metadata()

        # Extract tables
        self.data["content_modules"] = self._parse_content_modules()
        self.data["backend_services"] = self._parse_backend_services()
        self.data["infrastructure"] = self._parse_infrastructure()
        self.data["go_dependencies"] = self._parse_go_dependencies()

        # Extract design principles
        self.data["design_principles"] = self._parse_design_principles()

        # TODO: Add more sections as needed
        # self.data["api_namespaces"] = self._parse_api_namespaces()
        # self.data["config_keys"] = self._parse_config_keys()
        # self.data["qar_terminology"] = self._parse_qar_terminology()

        print(f"✅ Parsed {len(self.data)} sections")
        return self.data

    def _parse_metadata(self) -> dict[str, str]:
        """Extract metadata from document frontmatter."""
        metadata = {}

        # Extract version info from top of document
        patterns = {
            "last_updated": r"\*\*Last Updated\*\*:\s*(\S+)",
            "go_version": r"\*\*Go Version\*\*:\s*(\S+)",
            "nodejs_version": r"\*\*Node\.js\*\*:\s*(.+?)(?=\n|$)",
            "python_version": r"\*\*Python\*\*:\s*(\S+)",
            "postgresql_version": r"\*\*PostgreSQL\*\*:\s*(.+?)(?=\n|$)",
            "build_command": r"\*\*Build Command\*\*:\s*`([^`]+)`",
        }

        for key, pattern in patterns.items():
            match = re.search(pattern, self.content)
            if match:
                metadata[key] = match.group(1).strip()

        return metadata

    def _parse_content_modules(self) -> list[dict[str, str]]:
        """Parse Content Modules table."""
        modules = []

        # Find the table section
        table_match = re.search(
            r"## Content Modules\n\n(.*?)(?=\n---|\n##|\Z)",
            self.content,
            re.DOTALL,
        )
        if not table_match:
            print("⚠️  Warning: Content Modules table not found")
            return modules

        table_text = table_match.group(1)

        # Parse table rows (skip header and separator)
        lines = [line.strip() for line in table_text.split("\n") if "|" in line]
        if len(lines) < 3:  # Need header + separator + at least one row
            return modules

        for line in lines[2:]:  # Skip header and separator
            if not line.strip():
                continue

            # Split by | and clean up
            cells = [cell.strip() for cell in line.split("|")[1:-1]]
            if len(cells) < 6:
                continue

            module = {
                "name": cells[0],
                "schema": cells[1],
                "status": cells[2],
                "primary_metadata": cells[3],
                "arr_integration": cells[4],
                "design_doc": self._extract_link(cells[5]),
            }
            modules.append(module)

        print(f"   ✓ Parsed {len(modules)} content modules")
        return modules

    def _parse_backend_services(self) -> list[dict[str, str]]:
        """Parse Backend Services table."""
        services = []

        table_match = re.search(
            r"## Backend Services\n\n(.*?)(?=\n---|\n##|\Z)",
            self.content,
            re.DOTALL,
        )
        if not table_match:
            print("⚠️  Warning: Backend Services table not found")
            return services

        table_text = table_match.group(1)
        lines = [line.strip() for line in table_text.split("\n") if "|" in line]
        if len(lines) < 3:
            return services

        for line in lines[2:]:
            if not line.strip():
                continue

            cells = [cell.strip() for cell in line.split("|")[1:-1]]
            if len(cells) < 5:
                continue

            service = {
                "name": cells[0],
                "package": cells[1],
                "fx_module": cells[2],
                "status": cells[3],
                "design_doc": self._extract_link(cells[4]),
            }
            services.append(service)

        print(f"   ✓ Parsed {len(services)} backend services")
        return services

    def _parse_infrastructure(self) -> list[dict[str, str]]:
        """Parse Infrastructure Components table."""
        components = []

        table_match = re.search(
            r"## Infrastructure Components\n\n(.*?)(?=\n---|\n##|\Z)",
            self.content,
            re.DOTALL,
        )
        if not table_match:
            print("⚠️  Warning: Infrastructure Components table not found")
            return components

        table_text = table_match.group(1)
        lines = [line.strip() for line in table_text.split("\n") if "|" in line]
        if len(lines) < 3:
            return components

        for line in lines[2:]:
            if not line.strip():
                continue

            cells = [cell.strip() for cell in line.split("|")[1:-1]]
            if len(cells) < 5:
                continue

            component = {
                "name": cells[0],
                "package": cells[1],
                "version": cells[2],
                "purpose": cells[3],
                "design_doc": self._extract_link(cells[4]),
            }
            components.append(component)

        print(f"   ✓ Parsed {len(components)} infrastructure components")
        return components

    def _parse_go_dependencies(self) -> dict[str, list[dict[str, str]]]:
        """Parse all Go Dependencies tables."""
        dependencies = {
            "core": [],
            "security": [],
            "observability": [],
            "resilience": [],
            "distributed": [],
        }

        # Map section titles to keys
        section_map = {
            "Go Dependencies (Core)": "core",
            "Go Dependencies (Security & RBAC)": "security",
            "Go Dependencies (Observability)": "observability",
            "Go Dependencies (Resilience)": "resilience",
            "Go Dependencies (Distributed/Clustering)": "distributed",
        }

        for section_title, key in section_map.items():
            # Escape parentheses in regex
            escaped_title = re.escape(section_title)
            table_match = re.search(
                rf"## {escaped_title}\n\n(.*?)(?=\n##|\n---|\Z)",
                self.content,
                re.DOTALL,
            )
            if not table_match:
                continue

            table_text = table_match.group(1)
            lines = [line.strip() for line in table_text.split("\n") if "|" in line]
            if len(lines) < 3:
                continue

            for line in lines[2:]:
                if not line.strip():
                    continue

                cells = [cell.strip() for cell in line.split("|")[1:-1]]
                if len(cells) < 4:
                    continue

                dep = {
                    "package": cells[0].replace("`", ""),
                    "version": cells[1],
                    "purpose": cells[2],
                    "notes": cells[3],
                }
                dependencies[key].append(dep)

        total = sum(len(deps) for deps in dependencies.values())
        print(f"   ✓ Parsed {total} Go dependencies across {len(dependencies)} categories")
        return dependencies

    def _parse_design_principles(self) -> dict[str, Any]:
        """Extract core design principles."""
        principles = {}

        # Database strategy
        db_match = re.search(
            r"\*\*PostgreSQL ONLY\*\* - (.+?)(?=\n\n|$)",
            self.content,
        )
        if db_match:
            principles["database_strategy"] = db_match.group(1).strip()

        # Package update policy
        policy_match = re.search(
            r"\*\*1 Minor Behind\*\* - (.+?)(?=\n\n|$)",
            self.content,
        )
        if policy_match:
            principles["package_update_policy"] = policy_match.group(1).strip()

        # Test coverage
        coverage_match = re.search(
            r"\*\*80% minimum\*\* - (.+?)(?=\n\n|$)",
            self.content,
        )
        if coverage_match:
            principles["test_coverage"] = coverage_match.group(1).strip()

        # Design patterns table
        patterns_match = re.search(
            r"### Design Patterns\n\n(.*?)(?=\n---|\n##|\Z)",
            self.content,
            re.DOTALL,
        )
        if patterns_match:
            patterns_table = patterns_match.group(1)
            lines = [
                line.strip() for line in patterns_table.split("\n") if "|" in line
            ]
            if len(lines) >= 3:
                patterns_list = []
                for line in lines[2:]:
                    if not line.strip():
                        continue
                    cells = [cell.strip() for cell in line.split("|")[1:-1]]
                    if len(cells) >= 3:
                        patterns_list.append({
                            "pattern": cells[0],
                            "decision": cells[1],
                            "notes": cells[2],
                        })
                principles["design_patterns"] = patterns_list

        print("   ✓ Parsed design principles")
        return principles

    def _extract_link(self, markdown_link: str) -> str:
        """Extract the actual link from markdown [text](link) format."""
        match = re.search(r"\[.+?\]\(([^)]+)\)", markdown_link)
        if match:
            return match.group(1)
        return markdown_link.strip()

    def save_yaml(self, output_path: Path) -> None:
        """Save parsed data as YAML."""
        output_path.parent.mkdir(parents=True, exist_ok=True)

        with open(output_path, "w") as f:
            yaml.dump(
                self.data,
                f,
                default_flow_style=False,
                sort_keys=False,
                allow_unicode=True,
            )

        print(f"💾 Saved to {output_path}")


def main():
    """Main entry point."""
    # Paths
    repo_root = Path(__file__).parent.parent.parent
    sot_path = repo_root / "docs" / "dev" / "design" / "00_SOURCE_OF_TRUTH.md"
    output_path = repo_root / "data" / "shared-sot.yaml"

    if not sot_path.exists():
        print(f"❌ Error: {sot_path} not found")
        sys.exit(1)

    # Parse and save
    parser = SOTParser(sot_path)
    parser.parse()
    parser.save_yaml(output_path)

    print("\n✅ SOT parsing complete!")
    print(f"   Output: {output_path}")
    print(f"   Sections: {', '.join(parser.data.keys())}")


if __name__ == "__main__":
    main()
