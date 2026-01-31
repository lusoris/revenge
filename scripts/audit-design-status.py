#!/usr/bin/env python3
"""
Audit design documentation status.

Checks each design document for completeness across dimensions:
- Design: Has architecture/schema/implementation sections
- Sources: Has Developer Resources or external links
- Instructions: Has Implementation Checklist

Usage:
    python scripts/audit-design-status.py              # Audit and report
    python scripts/audit-design-status.py --update    # Update status file
"""

import argparse
import re
from collections import defaultdict
from pathlib import Path


# Project paths
SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
STATUS_FILE = DESIGN_DIR / "03_DESIGN_DOCS_STATUS.md"
QUESTIONS_FILE = DESIGN_DIR / "02_QUESTIONS_TO_DISCUSS.md"


def check_design_status(content: str) -> str:
    """Check if document has design content (schemas, tables, code examples)."""
    indicators = [
        r"```sql",  # SQL schemas
        r"```go",  # Go code
        r"## (Database|Schema|Architecture|Implementation|API|Features)",
        r"\| .+ \| .+ \|",  # Tables (at least 2 columns)
        r"CREATE TABLE",
        r"type \w+ struct",
    ]

    score = sum(
        1 for pattern in indicators if re.search(pattern, content, re.IGNORECASE)
    )

    if score >= 3:
        return "✅"
    elif score >= 1:
        return "🟡"
    return "🔴"


def check_sources_status(content: str) -> str:
    """Check if document has external sources/references."""

    # Count external links
    external_links = len(re.findall(r"\[.+\]\(https?://", content))

    has_resources_section = bool(
        re.search(r"## (Developer Resources|Sources)", content)
    )

    if has_resources_section and external_links >= 3:
        return "✅"
    elif external_links >= 2 or external_links >= 1:
        return "🟡"
    return "🔴"


def check_instructions_status(content: str) -> str:
    """Check if document has implementation instructions."""

    has_checklist = bool(re.search(r"## Implementation", content))
    checkbox_count = len(re.findall(r"- \[[ x]\]", content))

    if has_checklist and checkbox_count >= 5:
        return "✅"
    elif has_checklist or checkbox_count >= 3:
        return "🟡"
    return "🔴"


def audit_document(doc_path: Path) -> dict:
    """Audit a single document."""
    content = doc_path.read_text(encoding="utf-8")

    return {
        "path": doc_path,
        "design": check_design_status(content),
        "sources": check_sources_status(content),
        "instructions": check_instructions_status(content),
        "code": "🔴",  # Always red until implementation
        "linting": "🔴",
        "unit": "🔴",
        "integration": "🔴",
    }


def find_design_docs() -> dict[str, list[Path]]:
    """Find and categorize design documents."""
    categories = defaultdict(list)

    for md_file in sorted(DESIGN_DIR.rglob("*.md")):
        # Skip archives, index files, top-level files
        if ".archive" in str(md_file):
            continue
        if md_file.name in ["INDEX.md", "00_INDEX.md"]:
            continue
        if md_file.parent == DESIGN_DIR:
            continue  # Skip top-level docs like 00_SOURCE_OF_TRUTH.md

        # Categorize by path
        rel_path = md_file.relative_to(DESIGN_DIR)
        parts = rel_path.parts

        if len(parts) >= 2:
            category = parts[0]
            if len(parts) >= 3:
                category = f"{parts[0]}/{parts[1]}"
            categories[category].append(md_file)
        else:
            categories["other"].append(md_file)

    return dict(categories)


def generate_status_section(category: str, docs: list[dict]) -> list[str]:
    """Generate markdown section for a category."""
    lines = []

    # Format category name
    display_name = category.replace("/", " - ").replace("_", " ").title()
    lines.append(f"## {display_name}")
    lines.append("")
    lines.append(
        "| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |"
    )
    lines.append(
        "|----------|--------|---------|--------------|------|---------|------|-------------|"
    )

    design_ok = 0
    sources_ok = 0
    instructions_ok = 0

    for doc in sorted(docs, key=lambda d: d["path"].name):
        rel_path = doc["path"].relative_to(DESIGN_DIR)
        name = doc["path"].stem

        lines.append(
            f"| [{name}]({rel_path}) | {doc['design']} | {doc['sources']} | "
            f"{doc['instructions']} | {doc['code']} | {doc['linting']} | "
            f"{doc['unit']} | {doc['integration']} |"
        )

        if doc["design"] == "✅":
            design_ok += 1
        if doc["sources"] == "✅":
            sources_ok += 1
        if doc["instructions"] == "✅":
            instructions_ok += 1

    total = len(docs)
    lines.append("")
    lines.append(
        f"**Summary**: {design_ok}/{total} Design ✅ | {sources_ok}/{total} Sources ✅ | {instructions_ok}/{total} Instructions ✅"
    )
    lines.append("")
    lines.append("---")
    lines.append("")

    return lines


def generate_status_file(all_audits: dict[str, list[dict]]) -> str:
    """Generate the complete status file."""
    lines = [
        "# Design Documentation Status",
        "",
        "> Auto-generated overview of design document completeness",
        "",
        "**Last Updated**: Auto-generated",
        "",
        "---",
        "",
        "## Status Legend",
        "",
        "| Emoji | Meaning |",
        "|-------|---------|",
        "| ✅ | Complete |",
        "| 🟡 | Partial |",
        "| 🔴 | Not Started |",
        "",
        "---",
        "",
    ]

    # Overall summary first
    total_docs = 0
    total_design = 0
    total_sources = 0
    total_instructions = 0

    summary_rows = []

    for category, docs in sorted(all_audits.items()):
        total_docs += len(docs)
        design_ok = sum(1 for d in docs if d["design"] == "✅")
        sources_ok = sum(1 for d in docs if d["sources"] == "✅")
        instructions_ok = sum(1 for d in docs if d["instructions"] == "✅")

        total_design += design_ok
        total_sources += sources_ok
        total_instructions += instructions_ok

        display_name = category.replace("/", " - ").replace("_", " ").title()
        summary_rows.append(
            f"| {display_name} | {len(docs)} | {design_ok} ({100 * design_ok // len(docs)}%) | "
            f"{sources_ok} ({100 * sources_ok // len(docs)}%) | {instructions_ok} ({100 * instructions_ok // len(docs)}%) |"
        )

    lines.extend(
        [
            "## Overall Summary",
            "",
            f"**Total Documents**: {total_docs}",
            "",
            "| Category | Total | Design ✅ | Sources ✅ | Instructions ✅ |",
            "|----------|-------|-----------|------------|-----------------|",
        ]
    )
    lines.extend(summary_rows)
    lines.extend(
        [
            f"| **TOTAL** | **{total_docs}** | **{total_design} ({100 * total_design // total_docs}%)** | "
            f"**{total_sources} ({100 * total_sources // total_docs}%)** | **{total_instructions} ({100 * total_instructions // total_docs}%)** |",
            "",
            "---",
            "",
        ]
    )

    # Category sections
    for category, docs in sorted(all_audits.items()):
        lines.extend(generate_status_section(category, docs))

    # Notes
    lines.extend(
        [
            "## Notes",
            "",
            "- **Code/Linting/Unit/Integration**: All 🔴 as codebase is at template stage",
            "- **Design**: Schemas, tables, architecture diagrams, Go code examples",
            "- **Sources**: Developer Resources section with external documentation links",
            "- **Instructions**: Implementation Checklist with actionable items",
            "",
            "---",
            "",
            "## Regenerate",
            "",
            "```bash",
            "python scripts/audit-design-status.py --update",
            "```",
        ]
    )

    return "\n".join(lines)


def generate_questions_file() -> str:
    """Generate clean questions file."""
    return """# Questions & Gaps to Discuss

> Active tracking of open items requiring resolution

**Last Updated**: Auto-generated

---

## Status

All major design questions have been resolved. This file tracks any new questions that arise during implementation.

---

## Open Items

*No open items at this time.*

---

## Resolved Items (Reference)

All decisions documented in [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md):

- ✅ Package versions aligned to SOT
- ✅ Namespace migration complete ('c' → 'qar')
- ✅ All 23 Go package decisions finalized
- ✅ Design patterns documented
- ✅ External sources indexed in SOURCES.yaml
- ✅ Cross-reference indexes generated

---

## How to Add New Questions

1. Add under "## Open Items" with clear description
2. Include context and options if known
3. Tag priority: 🔴 Critical | 🟡 Important | 🟢 Nice to have
4. Link related design docs
5. Once resolved, move to "Resolved Items" with decision

---

<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->
"""


def main():
    parser = argparse.ArgumentParser(description="Audit design documentation status")
    parser.add_argument(
        "--update", "-u", action="store_true", help="Update status and questions files"
    )
    args = parser.parse_args()

    print("Finding design documents...")
    categories = find_design_docs()

    print("Auditing documents...")
    all_audits = {}
    for category, docs in categories.items():
        audits = [audit_document(doc) for doc in docs]
        all_audits[category] = audits

        # Print progress
        design_ok = sum(1 for d in audits if d["design"] == "✅")
        sources_ok = sum(1 for d in audits if d["sources"] == "✅")
        instructions_ok = sum(1 for d in audits if d["instructions"] == "✅")
        print(
            f"  {category}: {len(docs)} docs | D:{design_ok} S:{sources_ok} I:{instructions_ok}"
        )

    total_docs = sum(len(docs) for docs in all_audits.values())
    print(f"\nTotal: {total_docs} documents audited")

    if args.update:
        print("\nUpdating status file...")
        status_content = generate_status_file(all_audits)
        STATUS_FILE.write_text(status_content, encoding="utf-8")
        print(f"  Wrote {STATUS_FILE}")

        print("Updating questions file...")
        questions_content = generate_questions_file()
        QUESTIONS_FILE.write_text(questions_content, encoding="utf-8")
        print(f"  Wrote {QUESTIONS_FILE}")
    else:
        print("\nRun with --update to write changes")


if __name__ == "__main__":
    main()
