# Implementation Plan - Documentation Automation System

**Created**: 2026-01-31
**Based On**: 15_FINAL_ANSWERS_SUMMARY.md (32 decisions)
**Timeline**: 16-25 days (aggressive) or 20-30 days (conservative)
**Team Size**: 1-2 developers

---

## Overview

**Objective**: Build comprehensive documentation automation system that:
- Auto-generates design docs, user docs, API docs, project files from SOURCE_OF_TRUTH.md
- Syncs all config files (IDE, Coder, CI/CD) from SOT
- Integrates with GitHub (Projects, Discussions, branch protection, CodeQL)
- Prevents loops, validates everything, alerts on failure
- Migrates 136+ existing docs to template-based system

**Approach**: Pilot → Build → Test → Migrate → Integrate

---

## Phase Breakdown

### Phase 1: Foundation & Tooling (Days 1-3)
- Set up project structure
- Install dependencies
- Create bot user account
- Fetch missing sources
- Build SOT parser

### Phase 2: Template System (Days 4-7)
- Create Jinja2 templates (base + doc types)
- Create JSON schemas for validation
- Build template testing framework
- Pilot with 3 docs

### Phase 3: Data Extraction & Migration (Days 8-11)
- Build markdown parser for existing docs
- Extract data from 136+ docs → YAML
- Validate extracted data
- Multi-stage migration (10% → 50% → 100%)

### Phase 4: Validation Pipeline (Days 12-14)
- Build YAML schema validation
- Integrate markdown-link-check
- Build SOT reference validator
- Integrate gitleaks secret scanning
- Build full validation pipeline

### Phase 5: Generation Pipeline (Days 15-17)
- Build generation scripts
- Implement atomic operations
- Implement loop prevention (cooldown, bot author check)
- Build PR creation automation
- Test end-to-end generation

### Phase 6: GitHub Integration (Days 18-21)
- Configure GitHub Projects
- Configure GitHub Discussions
- Set up branch protection rules
- Enable CodeQL scanning
- Build repository settings sync
- Create automation workflows

### Phase 7: Config Synchronization (Days 22-23)
- Build config sync scripts
- Sync IDE settings from SOT
- Sync CI/CD configs from SOT
- Sync language version files from SOT

### Phase 8: Testing & Refinement (Days 24-25)
- Full end-to-end testing
- Load testing (100+ docs)
- Failure scenario testing
- Documentation of automation system
- Create troubleshooting guide

---

## Phase 1: Foundation & Tooling (Days 1-3)

### Day 1: Project Setup

#### 1.1 Create Project Structure
```bash
.github/
  automation-config.yml          # NEW: Automation settings
  workflows/
    doc-generation.yml           # NEW: Doc generation workflow
    source-fetching.yml          # Existing, may need updates
scripts/
  automation/                    # NEW: Automation scripts
    __init__.py
    sot_parser.py               # Parse SOURCE_OF_TRUTH.md
    doc_generator.py            # Generate docs from templates
    validator.py                # Validation pipeline
    pr_creator.py               # Create batched PRs
    config_sync.py              # Sync configs from SOT
  requirements.txt              # Update with new deps
templates/                      # NEW: Jinja2 templates
  base.md.jinja2               # Base template
  feature.md.jinja2            # Feature doc template
  service.md.jinja2            # Service doc template
  integration.md.jinja2        # Integration doc template
  partials/                    # Reusable partials
    status_table.jinja2
    implementation_checklist.jinja2
schemas/                        # NEW: JSON schemas
  feature.schema.json
  service.schema.json
  integration.schema.json
data/                           # NEW: Data files (generated from migration)
  features/
    video/
      MOVIE_MODULE.yaml
    music/
      MUSIC_MODULE.yaml
  services/
    AUTH.yaml
  integrations/
    metadata/
      video/
        TMDB.yaml
.automation-lock                # Generated during runs
```

#### 1.2 Install Dependencies
Update `scripts/requirements.txt`:
```txt
# Existing
pyyaml>=6.0
requests>=2.31
beautifulsoup4>=4.12
lxml>=5.0
html2text>=2024.2
jinja2>=3.1.5
ruff>=0.4
pytest>=8.0

# NEW
yamale>=4.0               # YAML schema validation
jsonschema>=4.17          # JSON schema validation
PyGithub>=2.1             # GitHub API integration
python-frontmatter>=1.0   # YAML frontmatter parsing
markdown>=3.5             # Markdown parsing
mistune>=3.0              # Faster markdown parsing alternative
```

Install npm tools:
```bash
npm install -g markdown-link-check markdown-toc
```

Install gitleaks:
```bash
# Linux
wget https://github.com/gitleaks/gitleaks/releases/download/v8.18.1/gitleaks_8.18.1_linux_x64.tar.gz
tar -xzf gitleaks_8.18.1_linux_x64.tar.gz
mv gitleaks /usr/local/bin/

# Or via brew
brew install gitleaks
```

#### 1.3 Create Bot User Account
1. Create GitHub account: `revenge-bot`
2. Add as collaborator to repository
3. Generate personal access token (classic):
   - Scope: `repo`, `workflow`, `write:discussion`
4. Add to repository secrets: `BOT_GITHUB_TOKEN`
5. Configure git:
   ```bash
   git config user.name "Revenge Bot"
   git config user.email "bot@revenge.dev"
   ```

---

### Day 2: Source Fetching

#### 2.1 Update SOURCES.yaml
Add 17 new sources to `docs/dev/sources/SOURCES.yaml`:

```yaml
# DEVOPS - Development Operations & GitHub
devops:
  # GitHub Documentation
  - id: github-readme-guide
    name: "GitHub README Best Practices"
    url: "https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes"
    type: web_page
    output: "devops/github-readme.md"

  - id: github-contributing-guide
    name: "GitHub CONTRIBUTING Guide"
    url: "https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/setting-guidelines-for-repository-contributors"
    type: web_page
    output: "devops/github-contributing.md"

  - id: github-templates-guide
    name: "GitHub Issue/PR Templates"
    url: "https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests"
    type: web_page
    output: "devops/github-templates.md"

  - id: github-metadata-guide
    name: "GitHub Repository Metadata"
    url: "https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-repository-languages"
    type: web_page
    output: "devops/github-metadata.md"

  - id: github-issues-docs
    name: "GitHub Issues Documentation"
    url: "https://docs.github.com/en/issues"
    type: web_page
    output: "devops/github-issues.md"

  - id: github-projects-docs
    name: "GitHub Projects Documentation"
    url: "https://docs.github.com/en/issues/planning-and-tracking-with-projects"
    type: web_page
    output: "devops/github-projects.md"

  - id: github-discussions-docs
    name: "GitHub Discussions Documentation"
    url: "https://docs.github.com/en/discussions"
    type: web_page
    output: "devops/github-discussions.md"

  - id: github-security-docs
    name: "GitHub Advanced Security Documentation"
    url: "https://docs.github.com/en/code-security"
    type: web_page
    output: "devops/github-security.md"

  # Documentation Style Guides
  - id: google-style-guide
    name: "Google Developer Documentation Style Guide"
    url: "https://developers.google.com/style"
    type: web_page
    output: "devops/google-style-guide.md"

  - id: writethedocs-guide
    name: "Write the Docs Best Practices"
    url: "https://www.writethedocs.org/guide/"
    type: web_page
    output: "devops/writethedocs.md"

  - id: markdown-guide
    name: "Markdown Guide"
    url: "https://www.markdownguide.org/basic-syntax/"
    type: web_page
    output: "devops/markdown-guide.md"

# APIS - API Documentation Standards
apis:
  - id: openapi-spec
    name: "OpenAPI 3.1 Specification"
    url: "https://spec.openapis.org/oas/v3.1.0"
    type: web_page
    output: "apis/openapi-spec.md"

  - id: stoplight-api-guide
    name: "Stoplight API Design Guide"
    url: "https://stoplight.io/api-design-guide"
    type: web_page
    output: "apis/stoplight-guide.md"

  - id: readme-api-docs
    name: "API Documentation Best Practices"
    url: "https://readme.com/blog/api-documentation-best-practices"
    type: web_page
    output: "apis/readme-api-docs.md"
```

#### 2.2 Fetch All Sources
```bash
cd /home/kilian/dev/revenge
python scripts/fetch-sources.py
```

Verify all 17 new sources fetched successfully.

---

### Day 3: SOT Parser

#### 3.1 Build `scripts/automation/sot_parser.py`
Parse SOURCE_OF_TRUTH.md to extract:
- Go dependencies (package → version)
- Development tools (tool → version → config sync paths)
- Infrastructure components (component → version)
- Module inventory (module → status → metadata source)
- API namespaces
- Configuration keys

**File**: `scripts/automation/sot_parser.py`
```python
#!/usr/bin/env python3
"""
Parse SOURCE_OF_TRUTH.md to extract structured data.
Used for generating shared data files and validating doc references.
"""

import re
from pathlib import Path
from typing import Dict, List, Any
import yaml


class SOTParser:
    """Parser for SOURCE_OF_TRUTH.md markdown document."""

    def __init__(self, sot_path: Path):
        self.sot_path = sot_path
        with open(sot_path) as f:
            self.content = f.read()

    def parse_table(self, section_heading: str, table_start_marker: str) -> List[Dict[str, str]]:
        """
        Parse markdown table from SOT.

        Args:
            section_heading: e.g., "## Go Dependencies"
            table_start_marker: First column header, e.g., "| Package |"

        Returns:
            List of dicts, one per table row
        """
        # Find section
        section_pattern = rf"{re.escape(section_heading)}.*?(?=\n## |\Z)"
        section_match = re.search(section_pattern, self.content, re.DOTALL)
        if not section_match:
            return []

        section_text = section_match.group(0)

        # Find table
        lines = section_text.split('\n')
        table_lines = []
        in_table = False

        for line in lines:
            if table_start_marker in line:
                in_table = True
                header_line = line
                continue

            if in_table:
                if line.startswith('|') and '---' not in line:
                    table_lines.append(line)
                elif not line.startswith('|'):
                    break

        if not table_lines:
            return []

        # Parse header
        headers = [h.strip() for h in header_line.split('|') if h.strip()]

        # Parse rows
        rows = []
        for line in table_lines:
            cells = [c.strip() for c in line.split('|') if c.strip()]
            if len(cells) == len(headers):
                row = dict(zip(headers, cells))
                rows.append(row)

        return rows

    def extract_go_dependencies(self) -> Dict[str, str]:
        """Extract Go dependencies: package → version."""
        rows = self.parse_table("## Go Dependencies", "| Package |")
        deps = {}
        for row in rows:
            package = row.get('Package', '').strip('`')
            version = row.get('Version', '')
            if package and version:
                deps[package] = version
        return deps

    def extract_dev_tools(self) -> List[Dict[str, Any]]:
        """Extract development tools with versions and config sync paths."""
        rows = self.parse_table("## Development Tools", "| Tool |")
        tools = []
        for row in rows:
            tool = {
                'name': row.get('Tool', ''),
                'version': row.get('Version', ''),
                'purpose': row.get('Purpose', ''),
                'status': row.get('Status', ''),
                'config_sync': row.get('Config Sync', '').split(', ') if row.get('Config Sync') else []
            }
            tools.append(tool)
        return tools

    def extract_infrastructure(self) -> Dict[str, str]:
        """Extract infrastructure components: component → version."""
        rows = self.parse_table("## Infrastructure Components", "| Component |")
        infra = {}
        for row in rows:
            component = row.get('Component', '')
            version = row.get('Version', '')
            if component and version:
                infra[component] = version
        return infra

    def extract_modules(self) -> List[Dict[str, str]]:
        """Extract content modules with status and metadata sources."""
        rows = self.parse_table("## Content Modules", "| Module |")
        return rows

    def extract_all(self) -> Dict[str, Any]:
        """Extract all structured data from SOT."""
        return {
            'go_dependencies': self.extract_go_dependencies(),
            'dev_tools': self.extract_dev_tools(),
            'infrastructure': self.extract_infrastructure(),
            'modules': self.extract_modules(),
        }


def main():
    """Test SOT parser."""
    sot_path = Path('docs/dev/design/00_SOURCE_OF_TRUTH.md')
    parser = SOTParser(sot_path)
    data = parser.extract_all()

    print(f"Go dependencies: {len(data['go_dependencies'])}")
    print(f"Dev tools: {len(data['dev_tools'])}")
    print(f"Infrastructure: {len(data['infrastructure'])}")
    print(f"Modules: {len(data['modules'])}")

    # Save to shared data file
    output_path = Path('data/shared-sot.yaml')
    output_path.parent.mkdir(parents=True, exist_ok=True)
    with open(output_path, 'w') as f:
        yaml.dump(data, f, default_flow_style=False, sort_keys=False)

    print(f"\n✅ Extracted SOT data to {output_path}")


if __name__ == '__main__':
    main()
```

#### 3.2 Test SOT Parser
```bash
python scripts/automation/sot_parser.py
```

Verify `data/shared-sot.yaml` created with extracted data.

---

**End of Phase 1 (Days 1-3)**

---

## Phase 2: Template System (Days 4-7)

### Day 4: Base Template & JSON Schemas

#### 4.1 Create Base Template
**File**: `templates/base.md.jinja2`
```jinja2
{#- Base template for all design docs -#}
---
{%- block frontmatter %}
sources: {{ sources | default([]) | tojson }}
design_refs: {{ design_refs | default([]) | tojson }}
category: {{ category }}
last_updated: {{ last_updated }}
{%- endblock %}
---

{% block title %}# {{ doc_title }}{% endblock %}

{% block status_table %}
## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | {{ status_design }} | {{ status_design_notes }} |
| Sources | {{ status_sources }} | {{ status_sources_notes }} |
| Instructions | {{ status_instructions }} | {{ status_instructions_notes }} |
| Code | {{ status_code }} | {{ status_code_notes }} |
| Linting | {{ status_linting }} | {{ status_linting_notes }} |
| Unit Testing | {{ status_unit_testing }} | {{ status_unit_testing_notes }} |
| Integration Testing | {{ status_integration_testing }} | {{ status_integration_testing_notes }} |

**Module**: `{{ code_location }}`
{% endblock %}

---

{% block toc %}
<!-- toc -->
<!-- tocstop -->
{% endblock %}

---

{% block overview %}
## Overview

{{ overview_content }}
{% endblock %}

---

{% block architecture %}
{%- if architecture_diagram %}
## Architecture

{{ architecture_diagram }}
{%- endif %}
{% endblock %}

---

{% block database_schema %}
{%- if database_tables %}
## Database Schema

{% for table in database_tables %}
### `{{ table.name }}`

{{ table.description }}

```sql
{{ table.schema }}
```

**Indexes**: {{ table.indexes | join(', ') }}

{%- if table.relationships %}
**Relationships**:
{% for rel in table.relationships %}
- {{ rel }}
{%- endfor %}
{%- endif %}

{% endfor %}
{%- endif %}
{% endblock %}

---

{% block implementation_checklist %}
## Implementation Checklist

{% for phase in implementation_phases %}
### Phase {{ phase.phase }}: {{ phase.name }}
{% for task in phase.tasks %}
- [ ] {{ task }}
{%- endfor %}

{% endfor %}
{% endblock %}

---

{% block related_documents %}
## Related Documents

{% for doc in related_docs %}
- [{{ doc.name }}]({{ doc.path }})
{%- endfor %}
{% endblock %}
```

#### 4.2 Create Feature Template
**File**: `templates/feature.md.jinja2`
```jinja2
{% extends "base.md.jinja2" %}

{% block title %}# {{ feature_name }}{% endblock %}

{% block overview %}
## Overview

{{ claude_overview }}

**Technical Purpose**: {{ technical_purpose }}

**Architecture Pattern**: {{ architecture_pattern }}
{% endblock %}

{% block api_endpoints %}
{%- if api_endpoints %}

---

## API Endpoints

**Namespace**: `{{ api_namespace }}`

{% for endpoint in api_endpoints %}
### {{ endpoint.method }} `{{ endpoint.path }}`

{{ endpoint.description }}

**Authentication**: {{ endpoint.auth_required }}
**RBAC Scope**: `{{ endpoint.rbac_scope }}`
**Rate Limit**: {{ endpoint.rate_limit }}

**Request Example**:
```
{{ endpoint.request_example }}
```

**Response Example**:
```json
{{ endpoint.response_example }}
```

{% endfor %}
{%- endif %}
{% endblock %}
```

#### 4.3 Create JSON Schemas
**File**: `schemas/feature.schema.json`
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "feature_name",
    "category",
    "status_design",
    "status_sources",
    "status_instructions",
    "status_code",
    "status_linting",
    "status_unit_testing",
    "status_integration_testing",
    "claude_overview",
    "implementation_phases"
  ],
  "properties": {
    "feature_name": {"type": "string"},
    "category": {"type": "string"},
    "sources": {
      "type": "array",
      "items": {"type": "string"}
    },
    "design_refs": {
      "type": "array",
      "items": {"type": "string"}
    },
    "status_design": {
      "type": "string",
      "enum": ["✅", "🟡", "🔴", "⚪"]
    },
    "status_design_notes": {"type": "string"},
    "status_sources": {
      "type": "string",
      "enum": ["✅", "🟡", "🔴", "⚪"]
    },
    "status_sources_notes": {"type": "string"},
    "status_instructions": {
      "type": "string",
      "enum": ["✅", "🟡", "🔴", "⚪"]
    },
    "status_instructions_notes": {"type": "string"},
    "status_code": {
      "type": "string",
      "enum": ["✅", "🟡", "🔴", "⚪"]
    },
    "status_code_notes": {"type": "string"},
    "status_linting": {
      "type": "string",
      "enum": ["✅", "🟡", "🔴", "⚪"]
    },
    "status_linting_notes": {"type": "string"},
    "status_unit_testing": {
      "type": "string",
      "enum": ["✅", "🟡", "🔴", "⚪"]
    },
    "status_unit_testing_notes": {"type": "string"},
    "status_integration_testing": {
      "type": "string",
      "enum": ["✅", "🟡", "🔴", "⚪"]
    },
    "status_integration_testing_notes": {"type": "string"},
    "code_location": {"type": "string"},
    "last_updated": {"type": "string", "format": "date"},
    "claude_overview": {"type": "string"},
    "technical_purpose": {"type": "string"},
    "architecture_pattern": {"type": "string"},
    "architecture_diagram": {"type": "string"},
    "database_tables": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["name", "description", "schema", "indexes"],
        "properties": {
          "name": {"type": "string"},
          "description": {"type": "string"},
          "schema": {"type": "string"},
          "indexes": {
            "type": "array",
            "items": {"type": "string"}
          },
          "relationships": {
            "type": "array",
            "items": {"type": "string"}
          }
        }
      }
    },
    "api_namespace": {"type": "string"},
    "api_endpoints": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["method", "path", "description"],
        "properties": {
          "method": {"type": "string", "enum": ["GET", "POST", "PUT", "PATCH", "DELETE"]},
          "path": {"type": "string"},
          "description": {"type": "string"},
          "request_example": {"type": "string"},
          "response_example": {"type": "string"},
          "auth_required": {"type": "string"},
          "rbac_scope": {"type": "string"},
          "rate_limit": {"type": "string"}
        }
      }
    },
    "implementation_phases": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["phase", "name", "tasks"],
        "properties": {
          "phase": {"type": "integer", "minimum": 1},
          "name": {"type": "string"},
          "tasks": {
            "type": "array",
            "items": {"type": "string"}
          }
        }
      }
    },
    "related_docs": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["name", "path"],
        "properties": {
          "name": {"type": "string"},
          "path": {"type": "string"}
        }
      }
    }
  }
}
```

---

**Days 5-7**: Continue building:
- Service template + schema
- Integration template + schema
- Template testing framework
- **Pilot migration** (MOVIE_MODULE, MUSIC_MODULE, TMDB)

---

## Phase 3: Data Extraction & Migration (Days 8-11)

**Build markdown parser, extract 136+ docs, multi-stage migration**

---

## Phase 4: Validation Pipeline (Days 12-14)

**Build comprehensive validation: YAML schema, links, SOT refs, secrets**

---

## Phase 5: Generation Pipeline (Days 15-17)

**Build generation scripts with atomic operations, loop prevention, PR automation**

---

## Phase 6: GitHub Integration (Days 18-21)

**Configure Projects, Discussions, branch protection, CodeQL, settings sync**

---

## Phase 7: Config Synchronization (Days 22-23)

**Sync IDE, CI/CD, language files from SOT**

---

## Phase 8: Testing & Refinement (Days 24-25)

**E2E testing, load testing, failure scenarios, documentation**

---

## Success Criteria

✅ All 136+ docs migrated to template system
✅ SOT → docs generation working
✅ Validation pipeline passing (YAML, lint, links, SOT refs, secrets)
✅ Loop prevention working (no infinite loops)
✅ PR automation working (batched by trigger type)
✅ Auto-merge for docs-only PRs
✅ GitHub integration complete (Projects, Discussions, branch protection, CodeQL)
✅ Config sync working (IDE, CI/CD, language files)
✅ Failure alerting working (GitHub issues created)
✅ 17 new sources fetched and documented
✅ Zero regressions (all existing docs still render correctly)
✅ 80%+ test coverage for automation scripts

---

**Status**: Ready for implementation
**Next Step**: Begin Phase 1, Day 1 - Project Setup

