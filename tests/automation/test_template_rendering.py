"""Tests for Jinja2 template rendering (templates/*.md.jinja2).

Ensures templates handle undefined variables gracefully and render correctly
for both Claude (design docs) and Wiki outputs.
"""

from pathlib import Path

import pytest
from jinja2 import Environment, FileSystemLoader, StrictUndefined


class TestTemplateLoading:
    """Test template files exist and load correctly."""

    @pytest.fixture
    def jinja_env(self):
        """Create Jinja2 environment with templates/ directory."""
        templates_dir = Path("templates")
        return Environment(
            loader=FileSystemLoader(templates_dir),
            undefined=StrictUndefined,  # Catch undefined variables
            trim_blocks=True,
            lstrip_blocks=True,
        )

    def test_all_templates_exist(self):
        """All expected template files must exist."""
        templates_dir = Path("templates")
        expected_templates = [
            "base.md.jinja2",
            "feature.md.jinja2",
            "service.md.jinja2",
            "integration.md.jinja2",
            "generic.md.jinja2",
        ]

        for template in expected_templates:
            template_path = templates_dir / template
            assert template_path.exists(), f"Template not found: {template}"

    def test_templates_load_without_errors(self, jinja_env):
        """Templates must load without Jinja2 syntax errors."""
        templates = [
            "feature.md.jinja2",
            "service.md.jinja2",
            "integration.md.jinja2",
            "generic.md.jinja2",
        ]

        for template_name in templates:
            # Should not raise TemplateSyntaxError
            template = jinja_env.get_template(template_name)
            assert template is not None


class TestServiceTemplate:
    """Test service.md.jinja2 template."""

    @pytest.fixture
    def jinja_env(self):
        """Create Jinja2 environment."""
        from jinja2 import Environment, FileSystemLoader

        return Environment(
            loader=FileSystemLoader("templates"),
            trim_blocks=True,
            lstrip_blocks=True,
        )

    def test_renders_with_minimal_data(self, jinja_env):
        """Service template must render with minimal YAML data."""
        template = jinja_env.get_template("service.md.jinja2")

        minimal_data = {
            "claude": True,
            "wiki": False,
            "doc_title": "Auth Service",
            "service_name": "Auth Service",
            "package_path": "internal/service/auth",
            "fx_module": "auth.Module",
            "technical_summary": "Authentication service",
            # IMPORTANT: No 'dependencies' field (must handle gracefully)
            # IMPORTANT: No 'provides' field (must handle gracefully)
            # IMPORTANT: No 'wiki_how_it_works' field (must handle gracefully)
        }

        # Should not raise UndefinedError
        output = template.render(**minimal_data)
        assert "Auth Service" in output
        assert "internal/service/auth" in output

    def test_dependencies_field_optional(self, jinja_env):
        """Service template must handle missing 'dependencies' field."""
        template = jinja_env.get_template("service.md.jinja2")

        data = {
            "claude": True,
            "service_name": "Test Service",
            "package_path": "internal/service/test",
            "fx_module": "test.Module",
            "technical_summary": "Test",
            # No 'dependencies' field
        }

        output = template.render(**data)
        assert "No external service dependencies" in output

    def test_provides_field_optional(self, jinja_env):
        """Service template must handle missing 'provides' field."""
        template = jinja_env.get_template("service.md.jinja2")

        data = {
            "claude": True,
            "service_name": "Test Service",
            "package_path": "internal/service/test",
            "fx_module": "test.Module",
            "technical_summary": "Test",
            # No 'provides' field
        }

        output = template.render(**data)
        assert "<!-- Service provides -->" in output

    def test_wiki_how_it_works_optional(self, jinja_env):
        """Service template must handle missing 'wiki_how_it_works' field."""
        template = jinja_env.get_template("service.md.jinja2")

        data = {
            "wiki": True,
            "claude": False,
            "service_name": "Test Service",
            "wiki_tagline": "Test tagline",
            "wiki_overview": "Test overview",
            # No 'wiki_how_it_works' field
        }

        output = template.render(**data)
        assert "<!-- How it works -->" in output


class TestIntegrationTemplate:
    """Test integration.md.jinja2 template."""

    @pytest.fixture
    def jinja_env(self):
        """Create Jinja2 environment."""
        from jinja2 import Environment, FileSystemLoader

        return Environment(
            loader=FileSystemLoader("templates"),
            trim_blocks=True,
            lstrip_blocks=True,
        )

    def test_renders_with_minimal_data(self, jinja_env):
        """Integration template must render with minimal YAML data."""
        template = jinja_env.get_template("integration.md.jinja2")

        minimal_data = {
            "claude": True,
            "wiki": False,
            "doc_title": "TMDb Integration",
            "integration_name": "TMDb Integration",
            "external_service": "TMDb",
            "integration_id": "tmdb",
            "technical_summary": "Primary metadata provider",
            # IMPORTANT: No 'api_base_url' field (must handle gracefully)
            # IMPORTANT: No 'auth_method' field (must handle gracefully)
            # IMPORTANT: No 'provides_data' field (must handle gracefully)
            # IMPORTANT: No 'auth_config' field (must handle gracefully)
            # IMPORTANT: No 'rate_limits' field (must handle gracefully)
            # IMPORTANT: No 'api_endpoints' field (must handle gracefully)
            # IMPORTANT: No 'cache_ttl' field (must handle gracefully)
            # IMPORTANT: No 'prerequisites' field (must handle gracefully)
        }

        # Should not raise UndefinedError
        output = template.render(**minimal_data)
        assert "TMDb Integration" in output

    def test_api_base_url_optional(self, jinja_env):
        """Integration template must handle missing 'api_base_url' field."""
        template = jinja_env.get_template("integration.md.jinja2")

        data = {
            "claude": True,
            "integration_name": "Test Integration",
            "external_service": "Test",
            "integration_id": "test",
            "technical_summary": "Test",
            # No 'api_base_url' field
        }

        output = template.render(**data)
        # Should not contain API Base URL line if field is missing
        assert "**API Base URL**" not in output or "**API Base URL**: ``" in output

    def test_auth_method_optional(self, jinja_env):
        """Integration template must handle missing 'auth_method' field."""
        template = jinja_env.get_template("integration.md.jinja2")

        data = {
            "claude": True,
            "integration_name": "Test Integration",
            "external_service": "Test",
            "integration_id": "test",
            "technical_summary": "Test",
            # No 'auth_method' field
        }

        output = template.render(**data)
        # Should not contain Authentication line if field is missing
        lines_with_auth = [line for line in output.split("\n") if "**Authentication**" in line]
        assert len(lines_with_auth) == 0

    def test_provides_data_optional(self, jinja_env):
        """Integration template must handle missing 'provides_data' field."""
        template = jinja_env.get_template("integration.md.jinja2")

        data = {
            "claude": True,
            "integration_name": "Test",
            "external_service": "Test",
            "integration_id": "test",
            "technical_summary": "Test",
            # No 'provides_data' field
        }

        output = template.render(**data)
        assert "<!-- Data provided by integration -->" in output

    def test_prerequisites_optional(self, jinja_env):
        """Integration template must handle missing 'prerequisites' field."""
        template = jinja_env.get_template("integration.md.jinja2")

        data = {
            "wiki": True,
            "claude": False,
            "integration_name": "Test",
            "wiki_tagline": "Test",
            "wiki_overview": "Test",
            # No 'prerequisites' field
        }

        # Should not raise UndefinedError
        output = template.render(**data)
        assert output is not None
        assert "Test" in output


class TestFeatureTemplate:
    """Test feature.md.jinja2 template."""

    @pytest.fixture
    def jinja_env(self):
        """Create Jinja2 environment."""
        from jinja2 import Environment, FileSystemLoader

        return Environment(
            loader=FileSystemLoader("templates"),
            trim_blocks=True,
            lstrip_blocks=True,
        )

    def test_renders_with_minimal_data(self, jinja_env):
        """Feature template must render with minimal YAML data."""
        template = jinja_env.get_template("feature.md.jinja2")

        minimal_data = {
            "claude": True,
            "wiki": False,
            "doc_title": "Movie Module",
            "feature_name": "Movie Module",
            "technical_summary": "Movie content module",
            # IMPORTANT: No 'content_types' field (must handle gracefully)
        }

        # Should not raise UndefinedError
        output = template.render(**minimal_data)
        assert "Movie Module" in output

    def test_content_types_optional(self, jinja_env):
        """Feature template must handle missing 'content_types' field."""
        template = jinja_env.get_template("feature.md.jinja2")

        data = {
            "claude": True,
            "feature_name": "Test Feature",
            "technical_summary": "Test",
            # No 'content_types' field
        }

        output = template.render(**data)
        # Should render empty join when content_types is missing
        assert "Content module for" in output


class TestTemplateDefaults:
    """Test all templates handle undefined variables with defaults."""

    @pytest.fixture
    def jinja_env(self):
        """Create Jinja2 environment."""
        from jinja2 import Environment, FileSystemLoader

        return Environment(
            loader=FileSystemLoader("templates"),
            trim_blocks=True,
            lstrip_blocks=True,
        )

    def test_service_template_has_defaults_for_all_optional_fields(self, jinja_env):
        """Service template must use | default() for all optional fields."""
        template_path = Path("templates/service.md.jinja2")
        content = template_path.read_text()

        # Check that optional fields use | default()
        optional_fields = ["dependencies", "provides", "wiki_how_it_works"]

        for field in optional_fields:
            # Either field uses | default() OR is in {%- else %} block
            assert (
                f"{field} | default" in content or
                f"if {field}" not in content  # Field not checked without default
            ), f"Field '{field}' should use | default() or not be checked"

    def test_integration_template_has_defaults_for_all_optional_fields(self, jinja_env):
        """Integration template must use | default() for all optional fields."""
        template_path = Path("templates/integration.md.jinja2")
        content = template_path.read_text()

        optional_fields = [
            "api_base_url",
            "auth_method",
            "provides_data",
            "auth_config",
            "rate_limits",
            "api_endpoints",
            "cache_ttl",
            "prerequisites",
            "wiki_how_it_works",
        ]

        for field in optional_fields:
            # Either field uses | default() OR is in {%- else %} block
            assert (
                f"{field} | default" in content or
                f"if {field}" not in content
            ), f"Field '{field}' should use | default() or not be checked"

    def test_feature_template_has_defaults_for_all_optional_fields(self, jinja_env):
        """Feature template must use | default() for all optional fields."""
        template_path = Path("templates/feature.md.jinja2")
        content = template_path.read_text()

        # content_types must use | default()
        assert "content_types | default" in content, (
            "Field 'content_types' should use | default()"
        )


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
