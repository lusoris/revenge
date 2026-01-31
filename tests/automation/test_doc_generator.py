#!/usr/bin/env python3
"""Tests for doc generator.

Test coverage for doc_generator.py:
- Initialization and setup
- Shared data loading
- Data merging (shared + doc-specific)
- Template rendering
- TOC generation integration
- Dual output (Claude + Wiki)
- Atomic file writes
- Error handling
"""

from unittest.mock import patch

import pytest
import yaml

from scripts.automation.doc_generator import DocGenerator


@pytest.fixture
def tmp_repo(tmp_path):
    """Create temporary repository structure."""
    # Create directory structure
    (tmp_path / "templates").mkdir()
    (tmp_path / "data").mkdir()
    (tmp_path / "docs" / "dev" / "design").mkdir(parents=True)
    (tmp_path / "docs" / "wiki").mkdir(parents=True)

    # Create minimal shared-sot.yaml
    shared_sot = {
        "metadata": {
            "project_name": "TestProject",
            "version": "1.0.0",
        },
        "infrastructure": ["PostgreSQL", "Redis"],
        "go_dependencies": {
            "language": [{"package": "Go", "version": "1.25+"}],
        },
    }
    (tmp_path / "data" / "shared-sot.yaml").write_text(
        yaml.dump(shared_sot, default_flow_style=False)
    )

    # Create minimal base template
    base_template = """# {{ doc_title }}

{% block overview %}
Overview content
{% endblock %}

{% block details %}
Details content
{% endblock %}
"""
    (tmp_path / "templates" / "base.md.jinja2").write_text(base_template)

    # Create feature template
    feature_template = """{% extends "base.md.jinja2" %}

{% block overview %}
Feature: {{ feature_name }}
{% endblock %}

{% block details %}
Category: {{ doc_category }}
{% endblock %}
"""
    (tmp_path / "templates" / "feature.md.jinja2").write_text(feature_template)

    return tmp_path


class TestInitialization:
    """Test DocGenerator initialization."""

    def test_init_sets_paths(self, tmp_repo):
        """Test initialization sets correct paths."""
        gen = DocGenerator(tmp_repo)

        assert gen.repo_root == tmp_repo
        assert gen.templates_dir == tmp_repo / "templates"
        assert gen.data_dir == tmp_repo / "data"
        assert gen.output_dir_claude == tmp_repo / "docs" / "dev" / "design"
        assert gen.output_dir_wiki == tmp_repo / "docs" / "wiki"

    def test_init_loads_shared_data(self, tmp_repo):
        """Test initialization loads shared-sot.yaml."""
        gen = DocGenerator(tmp_repo)

        assert gen.shared_data is not None
        assert "metadata" in gen.shared_data
        assert gen.shared_data["metadata"]["project_name"] == "TestProject"

    def test_init_creates_jinja_environment(self, tmp_repo):
        """Test initialization creates Jinja2 environment."""
        gen = DocGenerator(tmp_repo)

        assert gen.env is not None
        # Check template can be loaded
        template = gen.env.get_template("base.md.jinja2")
        assert template is not None

    def test_init_creates_toc_generator(self, tmp_repo):
        """Test initialization creates TOC generator."""
        gen = DocGenerator(tmp_repo)

        assert gen.toc_generator is not None


class TestSharedDataLoading:
    """Test shared data loading."""

    def test_load_existing_shared_data(self, tmp_repo):
        """Test loading existing shared-sot.yaml."""
        gen = DocGenerator(tmp_repo)
        data = gen._load_shared_data()

        assert data is not None
        assert "metadata" in data
        assert data["metadata"]["project_name"] == "TestProject"

    def test_load_missing_shared_data(self, tmp_repo):
        """Test handling missing shared-sot.yaml."""
        # Remove shared-sot.yaml
        (tmp_repo / "data" / "shared-sot.yaml").unlink()

        gen = DocGenerator(tmp_repo)
        data = gen._load_shared_data()

        assert data == {}

    def test_load_invalid_yaml(self, tmp_repo):
        """Test handling invalid YAML in shared-sot.yaml."""
        # Write invalid YAML
        (tmp_repo / "data" / "shared-sot.yaml").write_text("invalid: yaml: content:")

        with pytest.raises(yaml.YAMLError):
            gen = DocGenerator(tmp_repo)


class TestDataMerging:
    """Test data merging logic."""

    def test_merge_adds_shared_metadata(self, tmp_repo):
        """Test merging adds shared metadata."""
        gen = DocGenerator(tmp_repo)

        shared = {"metadata": {"project": "Test", "version": "1.0"}}
        doc = {"doc_title": "Feature"}

        merged = gen._merge_data(shared, doc)

        assert "project" in merged
        assert "version" in merged
        assert merged["project"] == "Test"
        assert merged["version"] == "1.0"

    def test_merge_adds_shared_infrastructure(self, tmp_repo):
        """Test merging adds shared infrastructure."""
        gen = DocGenerator(tmp_repo)

        shared = {"infrastructure": ["PostgreSQL", "Redis"]}
        doc = {"doc_title": "Feature"}

        merged = gen._merge_data(shared, doc)

        assert "shared_infrastructure" in merged
        assert merged["shared_infrastructure"] == ["PostgreSQL", "Redis"]

    def test_merge_doc_overrides_shared(self, tmp_repo):
        """Test doc-specific data overrides shared data."""
        gen = DocGenerator(tmp_repo)

        shared = {"metadata": {"version": "1.0"}}
        doc = {"version": "2.0", "doc_title": "Feature"}

        merged = gen._merge_data(shared, doc)

        # Doc-specific version should override shared
        assert merged["version"] == "2.0"

    def test_merge_preserves_all_doc_fields(self, tmp_repo):
        """Test merging preserves all doc-specific fields."""
        gen = DocGenerator(tmp_repo)

        shared = {}
        doc = {
            "doc_title": "Test Feature",
            "doc_category": "feature",
            "feature_name": "Test",
            "custom_field": "value",
        }

        merged = gen._merge_data(shared, doc)

        for key, value in doc.items():
            assert key in merged
            assert merged[key] == value


class TestDocumentGeneration:
    """Test document generation."""

    def test_generate_claude_only(self, tmp_repo):
        """Test generating Claude output only."""
        gen = DocGenerator(tmp_repo)

        # Create test YAML
        doc_yaml = tmp_repo / "data" / "test_feature.yaml"
        doc_yaml.write_text(yaml.dump({
            "doc_title": "Test Feature",
            "doc_category": "feature",
            "feature_name": "Test",
        }))

        result = gen.generate_doc(
            data_file=doc_yaml,
            template_name="feature.md.jinja2",
            output_subpath="features/test",
            render_both=False,
        )

        assert "claude" in result
        assert result["claude"].exists()

        # Verify content
        content = result["claude"].read_text()
        assert "Feature: Test" in content

    def test_generate_both_outputs(self, tmp_repo):
        """Test generating both Claude and Wiki outputs."""
        gen = DocGenerator(tmp_repo)

        # Create test YAML
        doc_yaml = tmp_repo / "data" / "test_feature.yaml"
        doc_yaml.write_text(yaml.dump({
            "doc_title": "Test Feature",
            "doc_category": "feature",
            "feature_name": "Test",
        }))

        result = gen.generate_doc(
            data_file=doc_yaml,
            template_name="feature.md.jinja2",
            output_subpath="features/test",
            render_both=True,
        )

        assert "claude" in result
        assert "wiki" in result
        assert result["claude"].exists()
        assert result["wiki"].exists()

    def test_generate_creates_output_directories(self, tmp_repo):
        """Test generation creates output directories."""
        gen = DocGenerator(tmp_repo)

        # Create test YAML
        doc_yaml = tmp_repo / "data" / "test_feature.yaml"
        doc_yaml.write_text(yaml.dump({
            "doc_title": "Test",
            "doc_category": "feature",
            "feature_name": "Test",
        }))

        # Use deep subdirectory path
        result = gen.generate_doc(
            data_file=doc_yaml,
            template_name="feature.md.jinja2",
            output_subpath="features/video/movies",
            render_both=False,
        )

        # Verify directory was created
        assert (tmp_repo / "docs" / "dev" / "design" / "features" / "video" / "movies").exists()

    def test_generate_with_toc(self, tmp_repo):
        """Test TOC is added to generated output."""
        gen = DocGenerator(tmp_repo)

        # Create test YAML with headers
        doc_yaml = tmp_repo / "data" / "test_feature.yaml"
        doc_yaml.write_text(yaml.dump({
            "doc_title": "Test Feature",
            "doc_category": "feature",
            "feature_name": "Test",
        }))

        result = gen.generate_doc(
            data_file=doc_yaml,
            template_name="feature.md.jinja2",
            output_subpath="features/test",
            render_both=False,
        )

        content = result["claude"].read_text()

        # Verify TOC was added
        assert "## Table of Contents" in content

    def test_generate_applies_data_merging(self, tmp_repo):
        """Test generation applies data merging."""
        gen = DocGenerator(tmp_repo)

        # Create simple template that uses shared data
        (tmp_repo / "templates" / "test.md.jinja2").write_text("""
# {{ doc_title }}

Project: {{ project_name }}
Infrastructure: {{ shared_infrastructure | join(', ') }}
""")

        doc_yaml = tmp_repo / "data" / "test.yaml"
        doc_yaml.write_text(yaml.dump({
            "doc_title": "Test",
        }))

        result = gen.generate_doc(
            data_file=doc_yaml,
            template_name="test.md.jinja2",
            output_subpath="test",
            render_both=False,
        )

        content = result["claude"].read_text()

        # Verify shared data was merged
        assert "Project: TestProject" in content
        assert "PostgreSQL, Redis" in content


class TestAtomicWrites:
    """Test atomic file write operations."""

    def test_save_output_creates_file(self, tmp_repo):
        """Test _save_output creates output file."""
        gen = DocGenerator(tmp_repo)

        output_path = tmp_repo / "test_output.md"
        content = "# Test Content\n\nSome text."

        result = gen._save_output(content, output_path)

        assert result == output_path
        assert output_path.exists()
        assert output_path.read_text() == content

    def test_save_output_creates_parent_dirs(self, tmp_repo):
        """Test _save_output creates parent directories."""
        gen = DocGenerator(tmp_repo)

        output_path = tmp_repo / "nested" / "deep" / "path" / "output.md"
        content = "Test"

        result = gen._save_output(content, output_path)

        assert result == output_path
        assert output_path.exists()

    def test_save_output_overwrites_existing(self, tmp_repo):
        """Test _save_output overwrites existing file."""
        gen = DocGenerator(tmp_repo)

        output_path = tmp_repo / "test.md"
        output_path.write_text("Old content")

        new_content = "New content"
        result = gen._save_output(new_content, output_path)

        assert result == output_path
        assert output_path.read_text() == new_content

    def test_save_output_atomic_on_error(self, tmp_repo):
        """Test _save_output cleans up temp file on error."""
        gen = DocGenerator(tmp_repo)

        output_path = tmp_repo / "test.md"

        # Mock file write to raise error
        with patch("builtins.open", side_effect=OSError("Write error")):
            with pytest.raises(IOError):
                gen._save_output("content", output_path)

        # Verify no temp files left
        temp_files = list(tmp_repo.glob("*.tmp"))
        assert len(temp_files) == 0


class TestErrorHandling:
    """Test error handling."""

    def test_generate_invalid_yaml(self, tmp_repo):
        """Test handling invalid YAML in data file."""
        gen = DocGenerator(tmp_repo)

        # Create invalid YAML
        doc_yaml = tmp_repo / "data" / "invalid.yaml"
        doc_yaml.write_text("invalid: yaml: content:")

        with pytest.raises(yaml.YAMLError):
            gen.generate_doc(
                data_file=doc_yaml,
                template_name="feature.md.jinja2",
                output_subpath="test",
            )

    def test_generate_missing_template(self, tmp_repo):
        """Test handling missing template."""
        gen = DocGenerator(tmp_repo)

        doc_yaml = tmp_repo / "data" / "test.yaml"
        doc_yaml.write_text(yaml.dump({"doc_title": "Test"}))

        with pytest.raises(Exception):  # Jinja2 will raise TemplateNotFound
            gen.generate_doc(
                data_file=doc_yaml,
                template_name="nonexistent.jinja2",
                output_subpath="test",
            )

    def test_generate_undefined_variable_strict(self, tmp_repo):
        """Test StrictUndefined raises error for missing variables."""
        gen = DocGenerator(tmp_repo)

        # Create template with undefined variable
        (tmp_repo / "templates" / "strict.md.jinja2").write_text("""
# {{ doc_title }}

{{ undefined_variable }}
""")

        doc_yaml = tmp_repo / "data" / "test.yaml"
        doc_yaml.write_text(yaml.dump({"doc_title": "Test"}))

        with pytest.raises(Exception):  # Jinja2 UndefinedError
            gen.generate_doc(
                data_file=doc_yaml,
                template_name="strict.md.jinja2",
                output_subpath="test",
            )


class TestIntegration:
    """Test full integration scenarios."""

    def test_full_workflow_feature_doc(self, tmp_repo):
        """Test complete workflow for feature document."""
        gen = DocGenerator(tmp_repo)

        # Create realistic feature YAML
        doc_yaml = tmp_repo / "data" / "features" / "video" / "MOVIE_MODULE.yaml"
        doc_yaml.parent.mkdir(parents=True, exist_ok=True)
        doc_yaml.write_text(yaml.dump({
            "doc_title": "Movie Module",
            "doc_category": "feature",
            "feature_name": "Movie Module",
            "technical_summary": "Movie content management",
        }))

        # Generate both outputs
        result = gen.generate_doc(
            data_file=doc_yaml,
            template_name="feature.md.jinja2",
            output_subpath="features/video",
            render_both=True,
        )

        # Verify both files created
        assert result["claude"].exists()
        assert result["wiki"].exists()

        # Verify paths are correct
        assert result["claude"] == tmp_repo / "docs" / "dev" / "design" / "features" / "video" / "MOVIE_MODULE.md"
        assert result["wiki"] == tmp_repo / "docs" / "wiki" / "features" / "video" / "MOVIE_MODULE.md"

        # Verify content has TOC
        claude_content = result["claude"].read_text()
        wiki_content = result["wiki"].read_text()

        assert "## Table of Contents" in claude_content
        assert "## Table of Contents" in wiki_content

        # Verify feature name in content
        assert "Movie Module" in claude_content
        assert "Movie Module" in wiki_content


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
