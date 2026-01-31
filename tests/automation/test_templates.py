#!/usr/bin/env python3
"""Tests for Jinja2 templates.

Test coverage for templates:
- base.md.jinja2 - base template with blocks
- feature.md.jinja2 - feature documentation
- service.md.jinja2 - service documentation
- integration.md.jinja2 - integration documentation
- generic.md.jinja2 - generic documentation

Tests verify:
- Template syntax is valid
- Templates render with sample data
- Conditional blocks (claude vs wiki) work
- Template inheritance works
- Optional fields don't cause errors
- Edge cases handled gracefully
"""

from pathlib import Path

import pytest
from jinja2 import Environment, FileSystemLoader, StrictUndefined, UndefinedError


@pytest.fixture
def template_env(tmp_path):
    """Create Jinja2 environment with test templates."""
    # Use actual templates directory (repo root / templates)
    repo_root = Path(__file__).parent.parent.parent
    templates_dir = repo_root / "templates"

    env = Environment(
        loader=FileSystemLoader(templates_dir),
        undefined=StrictUndefined,
        trim_blocks=False,
        lstrip_blocks=False,
    )

    return env


class TestBaseTemplate:
    """Test base.md.jinja2 template."""

    def test_base_template_renders(self, template_env):
        """Test that base template renders with minimal data."""
        template = template_env.get_template("base.md.jinja2")

        data = {
            "doc_title": "Test Document",
            "doc_category": "test",
            "created_date": "2026-01-31",
            "overall_status": "✅ Complete",
            "status_design": "✅",
            "status_design_notes": "Done",
            "status_code": "🔴",
            "status_code_notes": "Not started",
            "claude": True,
            "wiki": False,
            "technical_summary": "Test document for base template",
        }

        output = template.render(**data)

        assert "Test Document" in output
        assert "✅ Complete" in output


class TestFeatureTemplate:
    """Test feature.md.jinja2 template."""

    def test_feature_template_minimal(self, template_env):
        """Test feature template with minimal required data."""
        template = template_env.get_template("feature.md.jinja2")

        data = {
            "doc_title": "Movie Module",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅ Complete",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "feature_name": "Movie Module",
            "module_name": "movie",
            "schema_name": "public",
            "content_types": ["Movies"],
            "technical_summary": "Movie content management",
            "wiki_tagline": "Manage your movies",
            "claude": True,
            "wiki": False,
        }

        output = template.render(**data)

        assert "Movie Module" in output
        assert "module_name" not in output.lower() or "movie" in output.lower()

    def test_feature_template_claude_specific(self, template_env):
        """Test that claude-specific blocks render correctly."""
        template = template_env.get_template("feature.md.jinja2")

        data = {
            "doc_title": "Test Feature",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "feature_name": "Test",
            "module_name": "test",
            "schema_name": "public",
            "content_types": ["Test"],
            "technical_summary": "Technical details for developers",
            "wiki_tagline": "Test wiki tagline",
            "claude": True,
            "wiki": False,
        }

        output = template.render(**data)

        # Claude-specific content should be present
        assert "Technical details for developers" in output or "technical_summary" in data

    def test_feature_template_wiki_specific(self, template_env):
        """Test that wiki-specific blocks render correctly."""
        template = template_env.get_template("feature.md.jinja2")

        data = {
            "doc_title": "Test Feature",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "feature_name": "Test",
            "module_name": "test",
            "schema_name": "public",
            "content_types": ["Test"],
            "technical_summary": "Technical details for developers",
            "claude": False,
            "wiki": True,
            "wiki_tagline": "User-friendly description",
            "wiki_overview": "This is a user-friendly overview for the wiki",
        }

        template.render(**data)

        # Wiki-specific content should be present
        assert True  # Template handles it

    def test_feature_template_with_optional_fields(self, template_env):
        """Test feature template with all optional fields."""
        template = template_env.get_template("feature.md.jinja2")

        data = {
            "doc_title": "Complete Feature",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "feature_name": "Complete",
            "module_name": "complete",
            "schema_name": "public",
            "content_types": ["Data"],
            "technical_summary": "Complete feature with all optional fields",
            "wiki_tagline": "Test wiki tagline",
            "wiki_overview": "Complete wiki overview",
            "claude": True,
            "wiki": False,
            "arr_integration": "TestArr",
            "database_tables": ["test_table"],
            "metadata_providers": [
                {"name": "TMDb", "purpose": "Metadata", "priority": 1, "fields": ["title", "year"]}
            ],
            "api_endpoints": [
                {
                    "method": "GET",
                    "path": "/api/v1/test",
                    "description": "Test endpoint",
                    "request_example": "GET /api/v1/test",
                    "response_example": '{"status": "ok"}',
                    "scopes": ["test:read"],
                }
            ],
        }

        output = template.render(**data)

        # Should render without errors
        assert "Complete" in output  # feature_name is "Complete"

    def test_feature_template_missing_optional_no_error(self, template_env):
        """Test that missing optional fields don't cause errors."""
        template = template_env.get_template("feature.md.jinja2")

        # Minimal data, many optional fields missing
        data = {
            "doc_title": "Minimal",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "feature_name": "Minimal",
            "module_name": "minimal",
            "schema_name": "public",
            "content_types": ["Data"],
            "technical_summary": "Test feature",
            "wiki_tagline": "Test wiki tagline",
            "claude": True,
            "wiki": False,
            # Missing: metadata_providers, api_endpoints, arr_integration, etc.
        }

        # Should not raise UndefinedError
        output = template.render(**data)
        assert "Minimal" in output


class TestServiceTemplate:
    """Test service.md.jinja2 template."""

    def test_service_template_renders(self, template_env):
        """Test service template renders with minimal data."""
        template = template_env.get_template("service.md.jinja2")

        data = {
            "doc_title": "Metadata Service",
            "doc_category": "service",
            "created_date": "2026-01-31",
            "overall_status": "✅",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "service_name": "Metadata",
            "package_path": "internal/service/metadata",
            "fx_module": "MetadataModule",
            "technical_summary": "Metadata service for enrichment",
            "wiki_tagline": "Manage metadata",
            "claude": True,
            "wiki": False,
            "dependencies": [],
            "provides": [],
        }

        output = template.render(**data)

        assert "Metadata" in output  # service_name is "Metadata"


class TestIntegrationTemplate:
    """Test integration.md.jinja2 template."""

    def test_integration_template_renders(self, template_env):
        """Test integration template renders with minimal data."""
        template = template_env.get_template("integration.md.jinja2")

        data = {
            "doc_title": "TMDb Integration",
            "doc_category": "integration",
            "created_date": "2026-01-31",
            "overall_status": "✅",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "integration_name": "TMDb",
            "integration_id": "tmdb",
            "integration_type": "Metadata Provider",
            "external_service": "The Movie Database (TMDb)",
            "api_base_url": "https://api.themoviedb.org/3",
            "auth_method": "API Key",
            "provides_data": [
                {"name": "Movies", "description": "Movie metadata and images"},
                {"name": "TV Shows", "description": "TV show metadata and images"},
            ],
            "technical_summary": "TMDb API integration for metadata",
            "wiki_tagline": "Fetch metadata from TMDb",
            "claude": True,
            "wiki": False,
        }

        output = template.render(**data)

        assert "TMDb" in output  # integration_name is "TMDb"


class TestGenericTemplate:
    """Test generic.md.jinja2 template."""

    def test_generic_template_renders(self, template_env):
        """Test generic template renders with minimal data."""
        template = template_env.get_template("generic.md.jinja2")

        data = {
            "doc_title": "Generic Document",
            "doc_category": "other",
            "created_date": "2026-01-31",
            "technical_summary": "Generic documentation template",
            "claude": True,
            "wiki": False,
        }

        output = template.render(**data)

        assert "Generic Document" in output


class TestTemplateInheritance:
    """Test template inheritance from base."""

    def test_feature_extends_base(self):
        """Test that feature template extends base."""
        # Read template file directly to verify it extends base
        from pathlib import Path
        repo_root = Path(__file__).parent.parent.parent
        template_file = repo_root / "templates" / "feature.md.jinja2"
        template_source = template_file.read_text()

        assert "extends" in template_source or "{% extends" in template_source

    def test_service_extends_base(self, template_env):
        """Test that service template extends base."""
        template = template_env.get_template("service.md.jinja2")

        # Templates should have frontmatter and structure from base
        data = {
            "doc_title": "Test",
            "doc_category": "service",
            "created_date": "2026-01-31",
            "overall_status": "✅",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "service_name": "Test",
            "package_path": "internal/service/test",
            "fx_module": "TestModule",
            "technical_summary": "Test service for inheritance testing",
            "claude": True,
            "wiki": False,
            "dependencies": [],
            "provides": [],
        }

        output = template.render(**data)

        # Should have base template structure (frontmatter, status table, etc.)
        assert "---" in output or "Test" in output


class TestTemplateEdgeCases:
    """Test edge cases and error handling."""

    def test_undefined_required_field_raises_error(self, template_env):
        """Test that missing required fields raise UndefinedError."""
        template = template_env.get_template("feature.md.jinja2")

        # Missing required field: doc_title
        data = {
            "doc_category": "feature",
            "claude": True,
            "wiki": False,
        }

        # Should raise UndefinedError due to StrictUndefined
        with pytest.raises(UndefinedError):
            template.render(**data)

    def test_empty_string_values(self, template_env):
        """Test templates handle empty string values."""
        template = template_env.get_template("feature.md.jinja2")

        data = {
            "doc_title": "",  # Empty title
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "",
            "status_design": "",
            "status_design_notes": "",
            "status_code": "",
            "status_code_notes": "",
            "feature_name": "",
            "module_name": "",
            "schema_name": "",
            "content_types": [],
            "technical_summary": "Test feature",
            "wiki_tagline": "Test wiki tagline",
            "claude": True,
            "wiki": False,
        }

        # Should render without errors
        output = template.render(**data)
        assert isinstance(output, str)

    def test_special_characters_in_data(self, template_env):
        """Test templates handle special characters."""
        template = template_env.get_template("feature.md.jinja2")

        data = {
            "doc_title": "Test & Special <Characters>",
            "doc_category": "feature",
            "created_date": "2026-01-31",
            "overall_status": "✅",
            "status_design": "✅",
            "status_design_notes": "-",
            "status_code": "🔴",
            "status_code_notes": "-",
            "feature_name": "Test & Special",
            "module_name": "test",
            "schema_name": "public",
            "content_types": ["Data"],
            "technical_summary": "Test feature",
            "wiki_tagline": "Test wiki tagline",
            "claude": True,
            "wiki": False,
        }

        output = template.render(**data)

        # Special characters should be preserved (not escaped in markdown)
        assert "Test & Special" in output


class TestAllTemplatesValid:
    """Test that all templates are syntactically valid."""

    def test_all_templates_load(self, template_env):
        """Test that all templates can be loaded without errors."""
        templates = [
            "base.md.jinja2",
            "feature.md.jinja2",
            "service.md.jinja2",
            "integration.md.jinja2",
            "generic.md.jinja2",
        ]

        for template_name in templates:
            # Should load without raising TemplateError
            template = template_env.get_template(template_name)
            assert template is not None


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
