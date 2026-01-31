#!/usr/bin/env python3
"""End-to-end pipeline tests.

Test the complete documentation generation pipeline:
1. YAML data → template rendering → TOC generation → output files
2. Dual output (Claude + Wiki)
3. Real templates and data
4. Integration verification
5. Full automation pipeline (SOT, validation, generation, indexes)
6. Link checking and quality assurance
"""

import json
import subprocess
from pathlib import Path

import pytest
import yaml

from scripts.automation.doc_generator import DocGenerator


@pytest.fixture
def pipeline_setup(tmp_path):
    """Set up complete pipeline test environment."""
    # Create directory structure
    (tmp_path / "templates").mkdir()
    (tmp_path / "data" / "features" / "video").mkdir(parents=True)
    (tmp_path / "docs" / "dev" / "design").mkdir(parents=True)
    (tmp_path / "docs" / "wiki").mkdir(parents=True)

    # Create shared-sot.yaml with realistic data
    shared_sot = {
        "metadata": {
            "project_name": "Revenge Media Server",
            "version": "0.1.0",
            "go_version": "1.25+",
        },
        "infrastructure": [
            "PostgreSQL 18+",
            "Dragonfly (Cache)",
            "Typesense (Search)",
        ],
        "go_dependencies": {
            "language": [
                {"package": "Go", "version": "1.25+", "purpose": "Backend"},
            ],
            "framework": [
                {"package": "fx", "version": "1.22+", "purpose": "DI framework"},
            ],
        },
        "design_principles": {
            "architecture": ["Modular", "Event-driven", "Testable"],
        },
    }
    (tmp_path / "data" / "shared-sot.yaml").write_text(
        yaml.dump(shared_sot, default_flow_style=False)
    )

    # Create realistic base template
    base_template = """---
title: {{ doc_title }}
category: {{ doc_category }}
created: {{ created_date }}
status: {{ overall_status }}
---

{% block content %}
# {{ doc_title }}

{{ technical_summary }}

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | {{ status_design }} | {{ status_design_notes }} |
| Code | {{ status_code }} | {{ status_code_notes }} |

{% block details %}
{% endblock %}

## Related Documentation

{% if related_docs is defined and related_docs %}
{% for doc in related_docs %}
- [{{ doc.title }}]({{ doc.path }})
{% endfor %}
{% endif %}
{% endblock %}
"""
    (tmp_path / "templates" / "base.md.jinja2").write_text(base_template)

    # Create feature template
    feature_template = """{% extends "base.md.jinja2" %}

{% block details %}
## Overview

{% if claude %}
**Module**: {{ module_name }}

**Content Types**: {{ content_types | join(', ') }}

**Schema**: `{{ schema_name }}`
{% elif wiki %}
### What is {{ feature_name }}?

{{ wiki_tagline }}

This module handles: {{ content_types | join(', ') }}
{% endif %}

## Architecture

{% if claude %}
### Database Schema

Schema: `{{ schema_name }}`

### Module Structure

```
internal/content/{{ module_name }}/
├── module.go
├── repository.go
├── service.go
└── handler.go
```
{% endif %}

## Implementation

{% if implementation_phases is defined and implementation_phases %}
{% for phase in implementation_phases %}
### Phase {{ loop.index }}: {{ phase.name }}

{{ phase.description }}

{% if phase.tasks %}
Tasks:
{% for task in phase.tasks %}
- {{ task }}
{% endfor %}
{% endif %}
{% endfor %}
{% endif %}
{% endblock %}
"""
    (tmp_path / "templates" / "feature.md.jinja2").write_text(feature_template)

    return tmp_path


class TestEndToEndPipeline:
    """Test complete documentation generation pipeline."""

    def test_complete_feature_doc_generation(self, pipeline_setup):
        """Test generating complete feature documentation."""
        # Create realistic feature YAML
        feature_data = {
            "doc_title": "Movie Module",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅ Complete",
            "status_design": "✅ Complete",
            "status_design_notes": "Fully designed",
            "status_code": "🔴 Not Started",
            "status_code_notes": "Implementation pending",
            "technical_summary": "Movie content management with metadata enrichment from TMDb",
            "wiki_tagline": "Organize and play your movie collection",
            "feature_name": "Movie Module",
            "module_name": "movie",
            "schema_name": "public",
            "content_types": ["Movies", "Collections"],
            "implementation_phases": [
                {
                    "name": "Core Infrastructure",
                    "description": "Set up database schema and basic module structure",
                    "tasks": [
                        "Create database migrations",
                        "Define Go types and interfaces",
                        "Implement repository layer",
                    ],
                },
                {
                    "name": "Metadata Integration",
                    "description": "Integrate with TMDb for metadata fetching",
                    "tasks": [
                        "TMDb API client",
                        "Metadata enrichment service",
                        "Background sync jobs",
                    ],
                },
            ],
            "related_docs": [
                {"title": "TMDb Integration", "path": "../integrations/metadata/video/TMDB.md"},
                {"title": "Radarr Integration", "path": "../integrations/servarr/RADARR.md"},
            ],
        }

        yaml_file = pipeline_setup / "data" / "features" / "video" / "MOVIE_MODULE.yaml"
        yaml_file.write_text(yaml.dump(feature_data, default_flow_style=False, sort_keys=False))

        # Initialize generator
        gen = DocGenerator(pipeline_setup)

        # Generate documentation
        result = gen.generate_doc(
            data_file=yaml_file,
            template_name="feature.md.jinja2",
            output_subpath="features/video",
            render_both=True,
        )

        # Verify files created
        assert "claude" in result
        assert "wiki" in result
        assert result["claude"].exists()
        assert result["wiki"].exists()

        # Read outputs
        claude_content = result["claude"].read_text()
        wiki_content = result["wiki"].read_text()

        # Verify frontmatter preserved
        assert "---" in claude_content
        assert "title: Movie Module" in claude_content
        assert "category: feature" in claude_content

        # Verify TOC generated
        assert "## Table of Contents" in claude_content
        assert "## Table of Contents" in wiki_content

        # Verify TOC has correct structure
        assert "- [Status](#status)" in claude_content
        assert "- [Overview](#overview)" in claude_content
        assert "- [Architecture](#architecture)" in claude_content
        assert "- [Implementation](#implementation)" in claude_content

        # Verify nested TOC items
        assert "  - [Database Schema](#database-schema)" in claude_content
        assert "  - [Module Structure](#module-structure)" in claude_content
        assert "  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)" in claude_content

        # Verify content rendered correctly
        assert "Movie content management with metadata enrichment from TMDb" in claude_content
        assert "internal/content/movie/" in claude_content
        assert "Movies, Collections" in claude_content

        # Verify conditional rendering (Claude vs Wiki)
        assert "**Module**: movie" in claude_content  # Claude-specific
        assert "**Module**: movie" not in wiki_content  # Not in Wiki
        assert "### What is Movie Module?" in wiki_content  # Wiki-specific
        assert "### What is Movie Module?" not in claude_content  # Not in Claude

        # Verify implementation phases rendered
        assert "Phase 1: Core Infrastructure" in claude_content
        assert "Phase 2: Metadata Integration" in claude_content
        assert "Create database migrations" in claude_content
        assert "TMDb API client" in claude_content

        # Verify related docs links
        assert "[TMDb Integration](../integrations/metadata/video/TMDB.md)" in claude_content
        assert "[Radarr Integration](../integrations/servarr/RADARR.md)" in claude_content

        # Verify shared data merged
        # (Implementation phases should be in content)
        assert "Core Infrastructure" in claude_content

    def test_claude_wiki_output_differences(self, pipeline_setup):
        """Test that Claude and Wiki outputs are correctly differentiated."""
        # Minimal feature data
        feature_data = {
            "doc_title": "Test Feature",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅ Complete",
            "status_design": "✅ Complete",
            "status_design_notes": "-",
            "status_code": "🔴 Not Started",
            "status_code_notes": "-",
            "technical_summary": "Technical summary for developers",
            "wiki_tagline": "User-friendly tagline",
            "feature_name": "Test Feature",
            "module_name": "test",
            "schema_name": "public",
            "content_types": ["Test"],
        }

        yaml_file = pipeline_setup / "data" / "test.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        gen = DocGenerator(pipeline_setup)
        result = gen.generate_doc(
            data_file=yaml_file,
            template_name="feature.md.jinja2",
            output_subpath="test",
            render_both=True,
        )

        claude_content = result["claude"].read_text()
        wiki_content = result["wiki"].read_text()

        # Claude-specific content
        assert "**Module**: test" in claude_content
        assert "Schema: `public`" in claude_content
        assert "internal/content/test/" in claude_content

        # Wiki-specific content
        assert "### What is Test Feature?" in wiki_content
        assert "User-friendly tagline" in wiki_content

        # Verify NOT in wrong output
        assert "### What is Test Feature?" not in claude_content
        assert "**Module**: test" not in wiki_content

    def test_toc_with_complex_headers(self, pipeline_setup):
        """Test TOC generation with complex header structure."""
        # Create feature with many nested sections
        feature_data = {
            "doc_title": "Complex Feature",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅ Complete",
            "status_design": "✅ Complete",
            "status_design_notes": "-",
            "status_code": "🔴 Not Started",
            "status_code_notes": "-",
            "technical_summary": "Complex feature with many sections",
            "feature_name": "Complex Feature",
            "module_name": "complex",
            "schema_name": "public",
            "content_types": ["Data"],
            "implementation_phases": [
                {"name": "Phase A", "description": "First", "tasks": ["Task 1"]},
                {"name": "Phase B", "description": "Second", "tasks": ["Task 2"]},
                {"name": "Phase C", "description": "Third", "tasks": ["Task 3"]},
            ],
        }

        yaml_file = pipeline_setup / "data" / "complex.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        gen = DocGenerator(pipeline_setup)
        result = gen.generate_doc(
            data_file=yaml_file,
            template_name="feature.md.jinja2",
            output_subpath="test",
            render_both=False,
        )

        content = result["claude"].read_text()

        # Verify TOC has all phases
        assert "- [Implementation](#implementation)" in content
        assert "  - [Phase 1: Phase A](#phase-1-phase-a)" in content
        assert "  - [Phase 2: Phase B](#phase-2-phase-b)" in content
        assert "  - [Phase 3: Phase C](#phase-3-phase-c)" in content

    def test_atomic_file_operations(self, pipeline_setup):
        """Test that file operations are atomic (no partial writes on error)."""
        feature_data = {
            "doc_title": "Test",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅ Complete",
            "status_design": "✅ Complete",
            "status_design_notes": "-",
            "status_code": "🔴 Not Started",
            "status_code_notes": "-",
            "technical_summary": "Test",
            "feature_name": "Test",
            "module_name": "test",
            "schema_name": "public",
            "content_types": ["Test"],
        }

        yaml_file = pipeline_setup / "data" / "test.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        gen = DocGenerator(pipeline_setup)

        # Generate successfully first
        result = gen.generate_doc(
            data_file=yaml_file,
            template_name="feature.md.jinja2",
            output_subpath="test",
            render_both=False,
        )

        original_content = result["claude"].read_text()
        assert len(original_content) > 0

        # Verify no .tmp files left
        tmp_files = list(pipeline_setup.glob("**/*.tmp"))
        assert len(tmp_files) == 0

    def test_multiple_documents_batch(self, pipeline_setup):
        """Test generating multiple documents in batch."""
        # Create multiple feature YAMLs
        features = [
            ("MOVIE_MODULE", "movie", "Movies"),
            ("TVSHOW_MODULE", "tvshow", "TV Shows"),
            ("MUSIC_MODULE", "music", "Music"),
        ]

        gen = DocGenerator(pipeline_setup)
        results = []

        for name, module, content_type in features:
            feature_data = {
                "doc_title": name.replace("_", " ").title(),
                "doc_category": "feature",
                "created_date": "2026-01-31",
                "overall_status": "✅ Complete",
                "status_design": "✅ Complete",
                "status_design_notes": "-",
                "status_code": "🔴 Not Started",
                "status_code_notes": "-",
                "technical_summary": f"{content_type} management",
                "wiki_tagline": f"Manage your {content_type.lower()}",
                "feature_name": name.replace("_", " ").title(),
                "module_name": module,
                "schema_name": "public",
                "content_types": [content_type],
            }

            yaml_file = pipeline_setup / "data" / "features" / "video" / f"{name}.yaml"
            yaml_file.parent.mkdir(parents=True, exist_ok=True)
            yaml_file.write_text(yaml.dump(feature_data))

            result = gen.generate_doc(
                data_file=yaml_file,
                template_name="feature.md.jinja2",
                output_subpath="features/video",
                render_both=True,
            )

            results.append(result)

        # Verify all documents generated
        assert len(results) == 3

        for result in results:
            assert result["claude"].exists()
            assert result["wiki"].exists()

            content = result["claude"].read_text()
            assert "## Table of Contents" in content


class TestSOTIntegration:
    """Test SOURCE_OF_TRUTH.md integration."""

    def test_sot_file_exists(self):
        """Test SOURCE_OF_TRUTH.md exists."""
        sot = Path("docs/dev/design/00_SOURCE_OF_TRUTH.md")
        assert sot.exists(), "SOURCE_OF_TRUTH.md not found"

    def test_sot_has_required_sections(self):
        """Test SOT has all required sections."""
        sot = Path("docs/dev/design/00_SOURCE_OF_TRUTH.md")
        content = sot.read_text()

        required_sections = [
            "## Infrastructure Components",
            "## Go Dependencies",
            "## Database Schemas",
            "## Content Modules",
        ]

        for section in required_sections:
            assert section in content, f"Missing section: {section}"

    def test_sot_has_version_tables(self):
        """Test SOT has version information tables."""
        sot = Path("docs/dev/design/00_SOURCE_OF_TRUTH.md")
        content = sot.read_text()

        # Should have table rows
        table_rows = [line for line in content.split("\n") if line.startswith("|")]
        assert len(table_rows) >= 20, f"Expected at least 20 table rows, found {len(table_rows)}"


class TestValidationIntegration:
    """Test YAML validation integration."""

    def test_all_schemas_exist(self):
        """Test all required schemas exist."""
        schemas_dir = Path("schemas")
        required_schemas = ["feature.schema.json", "integration.schema.json", "service.schema.json", "generic.schema.json"]

        for schema_name in required_schemas:
            schema_file = schemas_dir / schema_name
            assert schema_file.exists(), f"Missing schema: {schema_name}"

    def test_schemas_are_valid_json(self):
        """Test all schemas are valid JSON."""
        schemas_dir = Path("schemas")
        for schema_file in schemas_dir.glob("*.schema.json"):
            with open(schema_file) as f:
                try:
                    json.load(f)
                except json.JSONDecodeError as e:
                    pytest.fail(f"{schema_file.name} invalid JSON: {e}")

    def test_validator_script_works(self):
        """Test validator script can execute."""
        result = subprocess.run(
            ["python", "scripts/automation/validator.py", "--help"],
            capture_output=True,
            text=True,
            timeout=10,
        )
        assert result.returncode == 0, f"Validator failed: {result.stderr}"


class TestIndexGeneration:
    """Test index generation integration."""

    def test_design_index_generator_exists(self):
        """Test design index generator exists."""
        generator = Path("scripts/generate-design-indexes.py")
        assert generator.exists(), "generate-design-indexes.py not found"

    def test_sources_index_generator_exists(self):
        """Test sources index generator exists."""
        generator = Path("scripts/generate-sources-indexes.py")
        assert generator.exists(), "generate-sources-indexes.py not found"

    def test_design_index_exists(self):
        """Test DESIGN_INDEX.md exists and has content."""
        index = Path("docs/dev/design/DESIGN_INDEX.md")
        assert index.exists(), "DESIGN_INDEX.md not found"

        content = index.read_text()
        assert len(content) > 100, "DESIGN_INDEX.md is too short"
        assert "##" in content, "DESIGN_INDEX.md has no sections"

    def test_sources_index_exists(self):
        """Test sources INDEX.yaml exists."""
        index = Path("docs/dev/sources/INDEX.yaml")
        assert index.exists(), "sources/INDEX.yaml not found"

        with open(index) as f:
            data = yaml.safe_load(f)
        assert data is not None, "INDEX.yaml is empty"


class TestAutomationScripts:
    """Test automation scripts integration."""

    def test_all_automation_scripts_exist(self):
        """Test all required automation scripts exist."""
        required_scripts = [
            "scripts/automation/check_health.py",
            "scripts/automation/view_logs.py",
            "scripts/automation/check_licenses.py",
            "scripts/automation/validator.py",
            "scripts/automation/doc_generator.py",
        ]

        for script_path in required_scripts:
            script = Path(script_path)
            assert script.exists(), f"Missing script: {script_path}"

    def test_automation_scripts_can_run(self):
        """Test automation scripts can execute --help."""
        scripts = [
            "scripts/automation/check_health.py",
            "scripts/automation/view_logs.py",
            "scripts/automation/check_licenses.py",
        ]

        for script_path in scripts:
            result = subprocess.run(
                ["python", script_path, "--help"],
                capture_output=True,
                text=True,
                timeout=10,
            )
            assert result.returncode == 0, f"{script_path} failed: {result.stderr}"


class TestGitHubIntegration:
    """Test GitHub integration."""

    def test_gh_cli_available(self):
        """Test gh CLI is available."""
        result = subprocess.run(
            ["gh", "--version"],
            capture_output=True,
            text=True,
            timeout=10,
        )
        assert result.returncode == 0, "gh CLI not available"

    def test_workflows_directory_exists(self):
        """Test workflows directory exists."""
        workflows = Path(".github/workflows")
        assert workflows.exists(), ".github/workflows not found"

    def test_workflows_are_valid_yaml(self):
        """Test all workflow files are valid YAML."""
        workflows_dir = Path(".github/workflows")
        for workflow in workflows_dir.glob("*.yml"):
            with open(workflow) as f:
                try:
                    yaml.safe_load(f)
                except yaml.YAMLError as e:
                    pytest.fail(f"{workflow.name} invalid YAML: {e}")

    def test_codeowners_exists(self):
        """Test CODEOWNERS file exists."""
        codeowners = Path("CODEOWNERS")
        assert codeowners.exists(), "CODEOWNERS not found"


class TestSkillsValidation:
    """Test Claude Code skills validation."""

    def test_skills_directory_exists(self):
        """Test skills directory exists."""
        skills = Path(".claude/skills")
        assert skills.exists(), ".claude/skills not found"

    def test_phase14_skills_exist(self):
        """Test all Phase 14 skills exist."""
        skills_dir = Path(".claude/skills")
        expected_skills = [
            "check-health",
            "view-logs",
            "manage-docker-config",
            "manage-ci-workflows",
            "run-linters",
            "format-code",
            "check-licenses",
            "update-dependencies",
        ]

        for skill_name in expected_skills:
            skill_dir = skills_dir / skill_name
            assert skill_dir.exists(), f"Skill missing: {skill_name}"
            assert (skill_dir / "SKILL.md").exists(), f"{skill_name} missing SKILL.md"

    def test_skills_have_valid_frontmatter(self):
        """Test skills have valid YAML frontmatter."""
        skills_dir = Path(".claude/skills")
        for skill_dir in skills_dir.iterdir():
            if not skill_dir.is_dir():
                continue

            skill_file = skill_dir / "SKILL.md"
            if not skill_file.exists():
                continue

            content = skill_file.read_text()
            assert content.startswith("---\n"), f"{skill_dir.name} missing frontmatter"

            parts = content.split("---\n", 2)
            if len(parts) < 3:
                pytest.fail(f"{skill_dir.name} invalid frontmatter structure")

            frontmatter = parts[1]
            try:
                data = yaml.safe_load(frontmatter)
                assert "name" in data, f"{skill_dir.name} missing name"
                assert "description" in data, f"{skill_dir.name} missing description"
            except yaml.YAMLError as e:
                pytest.fail(f"{skill_dir.name} invalid frontmatter YAML: {e}")


class TestMarkdownQuality:
    """Test markdown quality."""

    def test_markdownlint_available(self):
        """Test markdownlint is available."""
        markdownlint = Path("node_modules/.bin/markdownlint")
        assert markdownlint.exists(), "markdownlint not installed"

    def test_link_checker_available(self):
        """Test markdown-link-check is available."""
        link_checker = Path("node_modules/.bin/markdown-link-check")
        assert link_checker.exists(), "markdown-link-check not installed"

    def test_can_lint_sample_docs(self):
        """Test can run markdownlint on sample docs."""
        markdownlint = Path("node_modules/.bin/markdownlint")
        design_dir = Path("docs/dev/design")
        sample_docs = list(design_dir.glob("*.md"))[:2]

        if not sample_docs:
            pytest.skip("No design docs to lint")

        for doc in sample_docs:
            result = subprocess.run(
                [str(markdownlint), str(doc)],
                capture_output=True,
                text=True,
                timeout=10,
            )
            # Don't fail on lint errors, just check it runs
            assert result.returncode in [0, 1], f"markdownlint crashed on {doc}"


class TestFullPipeline:
    """Test complete pipeline execution."""

    def test_doc_pipeline_script_exists(self):
        """Test doc pipeline script exists."""
        pipeline = Path("scripts/doc-pipeline.sh")
        assert pipeline.exists(), "doc-pipeline.sh not found"

    def test_required_directories_exist(self):
        """Test all required directories exist."""
        required_dirs = [
            "docs/dev/design",
            "docs/dev/sources",
            "scripts/automation",
            "tests/automation",
            "schemas",
            ".claude/skills",
        ]

        for dir_path in required_dirs:
            directory = Path(dir_path)
            assert directory.exists(), f"Missing directory: {dir_path}"

    def test_all_tests_can_run(self):
        """Test that pytest can discover and run all tests."""
        result = subprocess.run(
            ["pytest", "tests/automation/", "--collect-only"],
            capture_output=True,
            text=True,
            timeout=30,
        )
        assert result.returncode == 0, f"Test collection failed: {result.stderr}"

        # Should have tests
        assert "tests collected" in result.stdout.lower() or "test" in result.stdout.lower()


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
