#!/usr/bin/env python3
"""Validate design document structure against template requirements.

Checks each design doc for required elements:
- Title (# heading)
- Description (> blockquote)
- Status table
- Overview section
- SOT reference
- Breadcrumb markers

Usage:
    python scripts/doc-pipeline/04-validate.py               # Check all
    python scripts/doc-pipeline/04-validate.py --strict      # Fail on warnings
    python scripts/doc-pipeline/04-validate.py --category services
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"

SKIP_FILES = {
    "00_SOURCE_OF_TRUTH.md",
    "01_DESIGN_DOC_TEMPLATE.md",
    "02_QUESTIONS_TO_DISCUSS.md",
    "03_DESIGN_DOCS_STATUS.md",
    "DESIGN_INDEX.md",
    "NAVIGATION.md",
}

REQUIRED_ELEMENTS = {
    "title": {
        "pattern": r"^#\s+.+$",
        "description": "Document title (# heading)",
        "severity": "error",
    },
    "description": {
        "pattern": r"^>\s+.+$",
        "description": "Description blockquote (> text)",
        "severity": "warning",
    },
    "status_table": {
        "pattern": r"\|\s*Dimension\s*\|\s*Status\s*\|",
        "description": "Status table with Dimension/Status columns",
        "severity": "warning",
    },
    "sot_reference": {
        "pattern": r"00_SOURCE_OF_TRUTH\.md",
        "description": "Reference to Source of Truth",
        "severity": "info",
    },
}

RECOMMENDED_ELEMENTS = {
    "overview": {
        "pattern": r"^##\s+Overview",
        "description": "Overview section",
    },
    "source_breadcrumbs": {
        "pattern": r"<!-- SOURCES:",
        "description": "Source breadcrumbs (minimal format)",
    },
    "design_breadcrumbs": {
        "pattern": r"<!-- DESIGN:",
        "description": "Design breadcrumbs (minimal format)",
    },
    "horizontal_rules": {
        "pattern": r"^---$",
        "description": "Section dividers (---)",
    },
}

CATEGORY_REQUIREMENTS = {
    "services": {
        "module_path": {
            "pattern": r"\*\*Module\*\*:\s*`internal/service/",
            "description": "Module path declaration",
            "severity": "info",
        },
        "dependencies": {
            "pattern": r"^##\s+Dependencies",
            "description": "Dependencies section",
            "severity": "info",
        },
    },
    "integrations": {
        "api_reference": {
            "pattern": r"API|api|endpoint|URL",
            "description": "API/endpoint reference",
            "severity": "info",
        },
    },
    "features": {
        "implementation": {
            "pattern": r"^##\s+Implementation|^##\s+Design",
            "description": "Implementation or Design section",
            "severity": "info",
        },
    },
}


def find_design_docs(category: str | None = None) -> list[Path]:
    """Find design documents, optionally filtered by category."""
    docs = []
    for md_file in sorted(DESIGN_DIR.rglob("*.md")):
        if ".archive" in str(md_file):
            continue
        if ".analysis" in str(md_file):
            continue
        if md_file.name in SKIP_FILES:
            continue
        if md_file.name.startswith("INDEX") or md_file.name.startswith("00_INDEX"):
            continue

        if category:
            rel_path = md_file.relative_to(DESIGN_DIR)
            if not str(rel_path).startswith(category):
                continue

        docs.append(md_file)
    return docs


def validate_document(doc_path: Path) -> dict:
    """Validate a single document. Returns validation results."""
    content = doc_path.read_text(encoding="utf-8")
    rel_path = doc_path.relative_to(DESIGN_DIR)

    results = {
        "path": str(rel_path),
        "errors": [],
        "warnings": [],
        "info": [],
        "passed": [],
    }

    # Check required elements
    for _name, spec in REQUIRED_ELEMENTS.items():
        pattern = re.compile(spec["pattern"], re.MULTILINE)
        if pattern.search(content):
            results["passed"].append(spec["description"])
        elif spec["severity"] == "error":
            results["errors"].append(f"Missing: {spec['description']}")
        elif spec["severity"] == "warning":
            results["warnings"].append(f"Missing: {spec['description']}")
        else:
            results["info"].append(f"Missing: {spec['description']}")

    # Check recommended elements
    for _name, spec in RECOMMENDED_ELEMENTS.items():
        pattern = re.compile(spec["pattern"], re.MULTILINE)
        if pattern.search(content):
            results["passed"].append(spec["description"])
        else:
            results["info"].append(f"Recommended: {spec['description']}")

    # Check category-specific requirements
    parts = rel_path.parts
    if parts:
        category = parts[0]
        if category in CATEGORY_REQUIREMENTS:
            for _name, spec in CATEGORY_REQUIREMENTS[category].items():
                pattern = re.compile(spec["pattern"], re.MULTILINE)
                if pattern.search(content):
                    results["passed"].append(spec["description"])
                else:
                    severity = spec.get("severity", "info")
                    msg = f"Category '{category}': {spec['description']}"
                    if severity == "error":
                        results["errors"].append(msg)
                    elif severity == "warning":
                        results["warnings"].append(msg)
                    else:
                        results["info"].append(msg)

    # Additional checks
    lines = content.split("\n")

    if len(lines) < 20:
        results["warnings"].append("Document is very short (< 20 lines)")

    if re.search(r"\bTODO\b|\bFIXME\b|\bXXX\b", content):
        results["info"].append("Contains TODO/FIXME markers")

    if re.search(
        r"\[.*TBD.*\]|\[.*TODO.*\]|placeholder|lorem ipsum", content, re.IGNORECASE
    ):
        results["warnings"].append("Contains placeholder content")

    long_lines = sum(
        1 for line in lines if len(line) > 200 and not line.startswith("|")
    )
    if long_lines > 5:
        results["info"].append(f"{long_lines} lines exceed 200 characters")

    return results


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate document structure")
    parser.add_argument(
        "--category",
        "-c",
        help="Only check specific category (e.g., 'services')",
    )
    parser.add_argument(
        "--strict",
        "-s",
        action="store_true",
        help="Treat warnings as errors",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Show all details",
    )
    parser.add_argument(
        "--quiet",
        "-q",
        action="store_true",
        help="Only show summary",
    )
    args = parser.parse_args()

    print("Finding design documents...")
    docs = find_design_docs(args.category)
    print(f"  Found {len(docs)} documents")

    total_errors = 0
    total_warnings = 0
    total_info = 0
    docs_with_issues = 0

    all_results = []

    for doc_path in docs:
        results = validate_document(doc_path)
        all_results.append(results)

        has_issues = results["errors"] or results["warnings"]
        if has_issues:
            docs_with_issues += 1

        total_errors += len(results["errors"])
        total_warnings += len(results["warnings"])
        total_info += len(results["info"])

        if not args.quiet and (has_issues or args.verbose):
            print(f"\n{results['path']}:")

            for err in results["errors"]:
                print(f"  ❌ {err}")

            for warn in results["warnings"]:
                print(f"  ⚠️  {warn}")

            if args.verbose:
                for info in results["info"]:
                    print(f"  ℹ️  {info}")
                for passed in results["passed"]:
                    print(f"  ✓ {passed}")

    print(f"\n{'=' * 50}")
    print("VALIDATION SUMMARY")
    print(f"{'=' * 50}")
    print(f"Documents checked: {len(docs)}")
    print(f"Documents with issues: {docs_with_issues}")
    print(f"Errors: {total_errors}")
    print(f"Warnings: {total_warnings}")
    print(f"Info: {total_info}")

    if not args.category:
        print("\nBy Category:")
        category_stats: dict[str, dict] = {}
        for result in all_results:
            parts = Path(result["path"]).parts
            cat = parts[0] if parts else "root"
            if cat not in category_stats:
                category_stats[cat] = {"docs": 0, "errors": 0, "warnings": 0}
            category_stats[cat]["docs"] += 1
            category_stats[cat]["errors"] += len(result["errors"])
            category_stats[cat]["warnings"] += len(result["warnings"])

        for cat, stats in sorted(category_stats.items()):
            status = "✓" if stats["errors"] == 0 and stats["warnings"] == 0 else "!"
            print(
                f"  {status} {cat}: {stats['docs']} docs, "
                f"{stats['errors']}E/{stats['warnings']}W",
            )

    if args.strict:
        exit_code = 1 if (total_errors + total_warnings) > 0 else 0
    else:
        exit_code = 1 if total_errors > 0 else 0

    if exit_code == 0:
        print("\n✓ Validation passed!")
    else:
        print("\n✗ Validation failed")

    return exit_code


if __name__ == "__main__":
    sys.exit(main())
