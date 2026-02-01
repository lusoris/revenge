#!/usr/bin/env python3
"""YAML Analyzer - Analyze migrated YAML files for completion status.

This analyzer:
1. Scans all YAML files in data/ directory
2. Identifies placeholder fields
3. Checks for missing required fields
4. Generates completion priority report
5. Suggests next steps for manual completion

Author: Automation System
Created: 2026-01-31
"""

from collections import defaultdict
from pathlib import Path

import yaml


class YAMLAnalyzer:
    """Analyze YAML data files for completion status."""

    def __init__(self, repo_root: Path):
        """Initialize analyzer."""
        self.repo_root = repo_root
        self.data_dir = repo_root / "data"

        # Required fields by category
        self.required_fields = {
            "feature": [
                "feature_name",
                "module_name",
                "schema_name",
                "content_types",
                "metadata_providers",
            ],
            "service": ["service_name", "package_path", "fx_module"],
            "integration": [
                "integration_name",
                "integration_id",
                "external_service",
                "api_base_url",
                "auth_method",
            ],
        }

    def analyze_file(self, yaml_file: Path) -> dict:
        """Analyze a single YAML file.

        Returns:
            Dict with analysis results
        """
        with open(yaml_file) as f:
            try:
                data = yaml.safe_load(f)
            except yaml.YAMLError as e:
                return {"error": f"YAML parse error: {e}"}

        result = {
            "file": str(yaml_file.relative_to(self.data_dir)),
            "category": data.get("doc_category", "unknown"),
            "title": data.get("doc_title", "UNKNOWN"),
            "placeholders": [],
            "missing_required": [],
            "has_sources": "sources" in data,
            "has_design_refs": "design_refs" in data,
        }

        # Find placeholder fields
        for key, value in data.items():
            if isinstance(value, str) and "PLACEHOLDER" in value:
                result["placeholders"].append(key)

        # Check for missing required fields
        category = result["category"]
        if category in self.required_fields:
            for field in self.required_fields[category]:
                if field not in data or not data[field]:
                    result["missing_required"].append(field)

        # Calculate completion score (0-100)
        total_fields = len(data)
        incomplete_fields = len(result["placeholders"]) + len(
            result["missing_required"],
        )
        result["completion_score"] = (
            int((total_fields - incomplete_fields) / total_fields * 100)
            if total_fields > 0
            else 0
        )

        return result

    def analyze_all(self) -> dict:
        """Analyze all YAML files in data directory.

        Returns:
            Dict with statistics and file-level results
        """
        yaml_files = list(self.data_dir.glob("**/*.yaml"))

        # Skip shared-sot.yaml
        yaml_files = [f for f in yaml_files if f.name != "shared-sot.yaml"]

        print(f"\n{'=' * 70}")
        print("YAML COMPLETION ANALYSIS")
        print(f"{'=' * 70}\n")
        print(f"Analyzing {len(yaml_files)} YAML files...\n")

        stats = {
            "total": len(yaml_files),
            "by_category": defaultdict(int),
            "by_completion": {"0-25%": 0, "26-50%": 0, "51-75%": 0, "76-100%": 0},
            "total_placeholders": 0,
            "total_missing": 0,
        }

        results = []

        for yaml_file in sorted(yaml_files):
            analysis = self.analyze_file(yaml_file)
            results.append(analysis)

            # Update stats
            category = analysis.get("category", "unknown")
            stats["by_category"][category] += 1

            # Completion bucket
            score = analysis.get("completion_score", 0)
            if score <= 25:
                stats["by_completion"]["0-25%"] += 1
            elif score <= 50:
                stats["by_completion"]["26-50%"] += 1
            elif score <= 75:
                stats["by_completion"]["51-75%"] += 1
            else:
                stats["by_completion"]["76-100%"] += 1

            stats["total_placeholders"] += len(analysis.get("placeholders", []))
            stats["total_missing"] += len(analysis.get("missing_required", []))

        return {"stats": stats, "results": results}

    def print_report(self, analysis: dict):
        """Print analysis report."""
        stats = analysis["stats"]
        results = analysis["results"]

        # Summary
        print(f"{'=' * 70}")
        print("SUMMARY")
        print(f"{'=' * 70}")
        print(f"Total files: {stats['total']}")
        print("\nBy category:")
        for category, count in sorted(stats["by_category"].items()):
            print(f"  {category}: {count} files")

        print("\nBy completion:")
        for bucket, count in sorted(stats["by_completion"].items()):
            print(f"  {bucket}: {count} files")

        print(f"\nPlaceholder fields found: {stats['total_placeholders']}")
        print(f"Missing required fields: {stats['total_missing']}")

        # Top priority files (features, services, integrations with low completion)
        print(f"\n{'=' * 70}")
        print("TOP PRIORITY FOR COMPLETION")
        print(f"{'=' * 70}\n")

        priority_categories = ["feature", "service", "integration"]
        priority_files = [
            r
            for r in results
            if r.get("category") in priority_categories
            and r.get("completion_score", 100) < 75
        ]

        # Sort by category (feature first) then by completion score
        priority_files.sort(
            key=lambda x: (
                priority_categories.index(x.get("category", "other")),
                x.get("completion_score", 0),
            ),
        )

        if priority_files:
            for i, result in enumerate(priority_files[:20], 1):  # Top 20
                print(f"{i}. {result['file']} ({result['category']})")
                print(f"   Completion: {result['completion_score']}%")
                if result["placeholders"]:
                    print(f"   Placeholders: {', '.join(result['placeholders'][:5])}")
                if result["missing_required"]:
                    print(
                        f"   Missing: {', '.join(result['missing_required'][:5])}",
                    )
                print()
        else:
            print("✅ No high-priority files need completion!\n")

        # Files with missing sources/design_refs
        print(f"{'=' * 70}")
        print("FILES MISSING CROSS-REFERENCES")
        print(f"{'=' * 70}\n")

        missing_refs = [
            r
            for r in results
            if not r.get("has_sources") and not r.get("has_design_refs")
        ]

        if missing_refs:
            print(f"Found {len(missing_refs)} files without sources or design refs:\n")
            for result in missing_refs[:10]:  # Show first 10
                print(f"  • {result['file']}")
            if len(missing_refs) > 10:
                print(f"  ... and {len(missing_refs) - 10} more")
            print()
        else:
            print("✅ All files have sources or design refs!\n")

        print(f"{'=' * 70}\n")


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent.parent

    analyzer = YAMLAnalyzer(repo_root)

    # Run analysis
    analysis = analyzer.analyze_all()

    # Print report
    analyzer.print_report(analysis)


if __name__ == "__main__":
    main()
