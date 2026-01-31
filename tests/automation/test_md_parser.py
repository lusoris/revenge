#!/usr/bin/env python3
"""Tests for markdown parser.

Test coverage for md_parser.py:
- Initialization and sources index loading
- parse_file() main workflow
- Status table parsing (7-dimension)
- Category determination (path-based)
- Section extraction
- Design path resolution
- YAML output generation
- Edge cases
"""

import pytest
import yaml

from scripts.automation.md_parser import MarkdownParser


@pytest.fixture
def parser_setup(tmp_path):
    """Set up markdown parser test environment."""
    # Create directory structure
    (tmp_path / "docs" / "dev" / "design" / "features").mkdir(parents=True)
    (tmp_path / "docs" / "dev" / "sources").mkdir(parents=True)

    # Create SOURCES.yaml
    sources = {
        "sources": {
            "apis": [
                {"id": "tmdb", "name": "TMDb API", "url": "https://tmdb.org/api"},
                {"id": "tvdb", "name": "TVDB API", "url": "https://tvdb.com/api"},
            ],
            "tools": [
                {"id": "fx", "name": "Uber FX", "url": "https://uber-go.github.io/fx/"},
            ],
        }
    }
    (tmp_path / "docs" / "dev" / "sources" / "SOURCES.yaml").write_text(
        yaml.dump(sources)
    )

    return tmp_path


class TestInitialization:
    """Test MarkdownParser initialization."""

    def test_init_sets_paths(self, parser_setup):
        """Test initialization sets correct paths."""
        parser = MarkdownParser(parser_setup)

        assert parser.repo_root == parser_setup
        assert parser.docs_dir == parser_setup / "docs" / "dev" / "design"
        assert parser.sources_dir == parser_setup / "docs" / "dev" / "sources"

    def test_init_loads_sources_index(self, parser_setup):
        """Test initialization loads sources index."""
        parser = MarkdownParser(parser_setup)

        assert "tmdb" in parser.sources_index
        assert "tvdb" in parser.sources_index
        assert "fx" in parser.sources_index

        assert parser.sources_index["tmdb"]["name"] == "TMDb API"
        assert parser.sources_index["tmdb"]["url"] == "https://tmdb.org/api"

    def test_init_with_missing_sources(self, tmp_path):
        """Test initialization when SOURCES.yaml is missing."""
        (tmp_path / "docs" / "dev" / "design").mkdir(parents=True)
        (tmp_path / "docs" / "dev" / "sources").mkdir(parents=True)

        parser = MarkdownParser(tmp_path)

        assert parser.sources_index == {}


class TestStatusTableParsing:
    """Test status table parsing."""

    def test_parse_status_table_complete(self, parser_setup):
        """Test parsing complete 7-dimension status table."""
        content = """# Test Document

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete |
| Sources | ✅ | Complete |
| Instructions | 🟡 | Partial |
| Code | 🔴 | Not Started |
| Linting | 🔴 | Not Started |
| Unit Testing | 🔴 | Not Started |
| Integration Testing | 🔴 | Not Started |

## Other Section
Content here.
"""
        parser = MarkdownParser(parser_setup)
        status = parser._parse_status_table(content)

        assert status is not None
        assert status["status_design"] == "✅"
        assert status["status_sources"] == "✅"
        assert status["status_instructions"] == "🟡"
        assert status["status_code"] == "🔴"
        assert status["status_linting"] == "🔴"
        assert status["status_unit_testing"] == "🔴"
        assert status["status_integration_testing"] == "🔴"
        assert status["overall_status"] == "✅"  # Matches design status

    def test_parse_status_table_missing(self, parser_setup):
        """Test handling missing status table."""
        content = """# Test Document

## Overview
No status table here.
"""
        parser = MarkdownParser(parser_setup)
        status = parser._parse_status_table(content)

        assert status is None


class TestCategoryDetermination:
    """Test document category determination."""

    def test_determine_category_feature(self, parser_setup):
        """Test category detection for feature docs."""
        parser = MarkdownParser(parser_setup)
        md_file = parser_setup / "docs" / "dev" / "design" / "features" / "MOVIE.md"

        category = parser._determine_category(md_file, "")

        assert category == "feature"

    def test_determine_category_service(self, parser_setup):
        """Test category detection for service docs."""
        parser = MarkdownParser(parser_setup)
        md_file = parser_setup / "docs" / "dev" / "design" / "services" / "METADATA.md"

        category = parser._determine_category(md_file, "")

        assert category == "service"

    def test_determine_category_integration(self, parser_setup):
        """Test category detection for integration docs."""
        parser = MarkdownParser(parser_setup)
        md_file = (
            parser_setup / "docs" / "dev" / "design" / "integrations" / "TMDB.md"
        )

        category = parser._determine_category(md_file, "")

        assert category == "integration"

    def test_determine_category_architecture(self, parser_setup):
        """Test category detection for architecture docs."""
        parser = MarkdownParser(parser_setup)
        md_file = parser_setup / "docs" / "dev" / "design" / "architecture" / "01_ARCHITECTURE.md"

        category = parser._determine_category(md_file, "")

        assert category == "architecture"

    def test_determine_category_other(self, parser_setup):
        """Test category for unrecognized paths."""
        parser = MarkdownParser(parser_setup)
        md_file = parser_setup / "docs" / "random" / "FILE.md"

        category = parser._determine_category(md_file, "")

        assert category == "other"


class TestSectionExtraction:
    """Test markdown section extraction."""

    def test_extract_sections_basic(self, parser_setup):
        """Test extracting basic sections."""
        content = """# Document Title

## Overview
Overview content here.

## Architecture
Architecture details.

## Implementation
Implementation notes.
"""
        parser = MarkdownParser(parser_setup)
        sections = parser._extract_sections(content)

        assert len(sections) == 3
        assert sections[0]["name"] == "Overview"
        assert "Overview content here" in sections[0]["content"]
        assert sections[1]["name"] == "Architecture"
        assert sections[2]["name"] == "Implementation"

    def test_extract_sections_skips_generated(self, parser_setup):
        """Test that auto-generated sections are skipped."""
        content = """# Document Title

## Table of Contents
- [Overview](#overview)

## Overview
Real content.

## Related Design Docs
- Link 1
- Link 2

## Sources & Cross-References
- Source 1
"""
        parser = MarkdownParser(parser_setup)
        sections = parser._extract_sections(content)

        # Should only have Overview, skip auto-generated sections
        assert len(sections) == 1
        assert sections[0]["name"] == "Overview"

    def test_extract_sections_empty(self, parser_setup):
        """Test extracting from document with no ## headers."""
        content = """# Document Title

Just some content without sections.
"""
        parser = MarkdownParser(parser_setup)
        sections = parser._extract_sections(content)

        assert len(sections) == 0


class TestDesignPathResolution:
    """Test design path resolution."""

    def test_resolve_design_path_index(self, parser_setup):
        """Test resolving INDEX paths."""
        parser = MarkdownParser(parser_setup)

        assert parser._resolve_design_path("OPERATIONS") == "operations/INDEX.md"
        assert parser._resolve_design_path("ARCHITECTURE") == "architecture/INDEX.md"
        assert parser._resolve_design_path("features") == "features/INDEX.md"

    def test_resolve_design_path_numbered(self, parser_setup):
        """Test resolving numbered design docs."""
        parser = MarkdownParser(parser_setup)

        path = parser._resolve_design_path("01_ARCHITECTURE")
        assert path == "architecture/01_ARCHITECTURE.md"

    def test_resolve_design_path_default(self, parser_setup):
        """Test default path resolution."""
        parser = MarkdownParser(parser_setup)

        path = parser._resolve_design_path("custom_doc")
        assert path == "custom_doc.md"


class TestParseFile:
    """Test complete file parsing."""

    def test_parse_file_complete(self, parser_setup):
        """Test parsing a complete markdown file."""
        md_content = """# Movie Module

<!-- SOURCES: tmdb, tvdb -->
<!-- DESIGN: OPERATIONS, services -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete |
| Sources | ✅ | Complete |
| Instructions | 🟡 | Partial |
| Code | 🔴 | Not Started |
| Linting | 🔴 | Not Started |
| Unit Testing | 🔴 | Not Started |
| Integration Testing | 🔴 | Not Started |

## Overview
Movie management module.

## Architecture
Technical details.
"""
        md_file = parser_setup / "docs" / "dev" / "design" / "features" / "MOVIE.md"
        md_file.write_text(md_content)

        parser = MarkdownParser(parser_setup)
        data = parser.parse_file(md_file)

        # Check extracted data
        assert data["doc_title"] == "Movie Module"
        assert data["doc_category"] == "feature"
        assert data["source_ids"] == ["tmdb", "tvdb"]
        assert data["design_ids"] == ["OPERATIONS", "services"]
        assert data["status_design"] == "✅"
        assert data["overall_status"] == "✅"
        assert len(data["sections"]) == 3  # Status + Overview + Architecture

    def test_parse_file_minimal(self, parser_setup):
        """Test parsing minimal markdown file."""
        md_content = """# Simple Doc

Just content.
"""
        md_file = parser_setup / "docs" / "dev" / "design" / "TEST.md"
        md_file.write_text(md_content)

        parser = MarkdownParser(parser_setup)
        data = parser.parse_file(md_file)

        assert data["doc_title"] == "Simple Doc"
        assert data["doc_category"] == "other"
        assert "source_file" in data
        assert "created_date" in data


class TestYAMLOutput:
    """Test YAML output generation."""

    def test_to_yaml_basic(self, parser_setup):
        """Test basic YAML output generation."""
        parser = MarkdownParser(parser_setup)

        data = {
            "doc_title": "Test Feature",
            "doc_category": "feature",
            "status_design": "✅",
        }

        yaml_output = parser.to_yaml(data, template_type="basic")

        # Should be valid YAML
        parsed = yaml.safe_load(yaml_output)
        assert parsed["doc_title"] == "Test Feature"
        assert parsed["doc_category"] == "feature"


class TestEdgeCases:
    """Test edge cases and error conditions."""

    def test_parse_file_missing_title(self, parser_setup):
        """Test parsing file without # title."""
        md_content = """## Section
No title header.
"""
        md_file = parser_setup / "docs" / "dev" / "design" / "NO_TITLE.md"
        md_file.write_text(md_content)

        parser = MarkdownParser(parser_setup)
        data = parser.parse_file(md_file)

        # Should handle gracefully
        assert "doc_title" not in data or data.get("doc_title") is None

    def test_parse_file_empty(self, parser_setup):
        """Test parsing empty file."""
        md_file = parser_setup / "docs" / "dev" / "design" / "EMPTY.md"
        md_file.write_text("")

        parser = MarkdownParser(parser_setup)
        data = parser.parse_file(md_file)

        # Should return basic structure
        assert isinstance(data, dict)
        assert data["doc_category"] == "other"

    def test_sources_index_malformed_sources_yaml(self, tmp_path):
        """Test handling malformed SOURCES.yaml."""
        (tmp_path / "docs" / "dev" / "design").mkdir(parents=True)
        (tmp_path / "docs" / "dev" / "sources").mkdir(parents=True)

        # Create malformed YAML
        (tmp_path / "docs" / "dev" / "sources" / "SOURCES.yaml").write_text(
            "invalid: yaml: syntax:"
        )

        with pytest.raises(yaml.YAMLError):
            MarkdownParser(tmp_path)


class TestIntegration:
    """Test full integration scenarios."""

    def test_full_parse_workflow(self, parser_setup):
        """Test complete parse and YAML generation workflow."""
        # Create realistic markdown file
        md_content = """# TV Show Module

<!-- SOURCES: tvdb, tmdb -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete |
| Sources | ✅ | Complete |
| Instructions | 🟡 | Partial |
| Code | 🔴 | Not Started |
| Linting | 🔴 | Not Started |
| Unit Testing | 🔴 | Not Started |
| Integration Testing | 🔴 | Not Started |

## Overview
Manages TV shows and episodes.

## Database Schema
Schema details.
"""
        md_file = parser_setup / "docs" / "dev" / "design" / "features" / "TVSHOW.md"
        md_file.write_text(md_content)

        # Parse
        parser = MarkdownParser(parser_setup)
        data = parser.parse_file(md_file)

        # Generate YAML
        yaml_output = parser.to_yaml(data)

        # Verify complete workflow
        assert data["doc_title"] == "TV Show Module"
        assert data["source_ids"] == ["tvdb", "tmdb"]
        assert len(data["sections"]) == 3  # Status + Overview + Database Schema

        # YAML should be valid
        parsed = yaml.safe_load(yaml_output)
        assert parsed["doc_title"] == "TV Show Module"


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
