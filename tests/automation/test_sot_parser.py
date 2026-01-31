#!/usr/bin/env python3
"""Tests for SOT parser.

Test coverage for sot_parser.py:
- Initialization and file loading
- Metadata extraction from frontmatter
- Content modules table parsing
- Backend services table parsing
- Infrastructure components table parsing
- Go dependencies parsing (all 5 categories)
- Design principles extraction
- Link extraction from markdown
- YAML output generation
- Edge cases and error handling
"""


import pytest
import yaml

from scripts.automation.sot_parser import SOTParser


@pytest.fixture
def minimal_sot(tmp_path):
    """Create minimal SOURCE_OF_TRUTH.md for testing."""
    sot_content = """# Source of Truth

**Last Updated**: 2026-01-31
**Go Version**: 1.25+
**Node.js**: 20.x (LTS)
**Python**: 3.12+
**PostgreSQL**: 18+ (required)
**Build Command**: `GOEXPERIMENT=greenteagc,jsonv2 go build ./...`

## Content Modules

| Module | Schema | Status | Primary Metadata | Arr Integration | Design Doc |
|--------|--------|--------|------------------|-----------------|------------|
| Movie | `public` | ✅ Complete | [TMDb](../integrations/metadata/video/TMDB.md) | [Radarr](../integrations/servarr/RADARR.md) | [MOVIE_MODULE.md](../features/video/MOVIE_MODULE.md) |
| TV Show | `public` | 🟡 Partial | [TVDB](../integrations/metadata/video/THETVDB.md) | [Sonarr](../integrations/servarr/SONARR.md) | [TVSHOW_MODULE.md](../features/video/TVSHOW_MODULE.md) |

## Backend Services

| Service | Package | fx Module | Status | Design Doc |
|---------|---------|-----------|--------|------------|
| Metadata | `internal/service/metadata` | `MetadataModule` | ✅ Complete | [METADATA.md](../services/METADATA.md) |
| Library | `internal/service/library` | `LibraryModule` | 🟡 Partial | [LIBRARY.md](../services/LIBRARY.md) |

## Infrastructure Components

| Component | Package | Version | Purpose | Design Doc |
|-----------|---------|---------|---------|------------|
| PostgreSQL | `pgx` | 18+ | Primary database | [DATABASE.md](../technical/DATABASE.md) |
| Dragonfly | `rueidis` | 1.0+ | Distributed cache | [DRAGONFLY.md](../integrations/infrastructure/DRAGONFLY.md) |

## Go Dependencies (Core)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `fx` | 1.22+ | Dependency injection | Core framework |
| `ogen` | 1.6+ | OpenAPI code gen | API framework |

## Go Dependencies (Security & RBAC)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `casbin` | 2.100+ | RBAC engine | Policy-based access control |

## Go Dependencies (Observability)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `slog` | stdlib | Structured logging | Go 1.21+ |

## Go Dependencies (Resilience)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `gobreaker` | 0.6+ | Circuit breaker | Failure isolation |

## Go Dependencies (Distributed/Clustering)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `raft` | 1.7+ | Consensus | Leader election |

## Design Principles

**PostgreSQL ONLY** - No SQLite, MySQL, or other DB engines

**1 Minor Behind** - Stay one minor version behind latest stable

**80% minimum** - All packages must maintain 80%+ test coverage

### Design Patterns

| Pattern | Decision | Notes |
|---------|----------|-------|
| Database | **PostgreSQL ONLY** - No SQLite, MySQL, or other DB engines | Single DB for simplicity |
| Packages | **1 Minor Behind** - Stay one minor version behind latest stable | Stability over bleeding edge |
| Testing | **80% minimum** - All packages must maintain 80%+ test coverage | Quality gate |
"""
    sot_file = tmp_path / "00_SOURCE_OF_TRUTH.md"
    sot_file.write_text(sot_content)
    return sot_file


class TestInitialization:
    """Test SOTParser initialization."""

    def test_init_reads_file(self, minimal_sot):
        """Test initialization reads SOT file."""
        parser = SOTParser(minimal_sot)

        assert parser.sot_path == minimal_sot
        assert len(parser.content) > 0
        assert "Source of Truth" in parser.content
        assert parser.data == {}

    def test_init_with_missing_file(self, tmp_path):
        """Test initialization with non-existent file."""
        missing_file = tmp_path / "missing.md"

        with pytest.raises(FileNotFoundError):
            SOTParser(missing_file)


class TestMetadataParsing:
    """Test metadata extraction."""

    def test_parse_metadata_all_fields(self, minimal_sot):
        """Test parsing all metadata fields."""
        parser = SOTParser(minimal_sot)
        metadata = parser._parse_metadata()

        assert metadata["last_updated"] == "2026-01-31"
        assert metadata["go_version"] == "1.25+"
        assert metadata["nodejs_version"] == "20.x (LTS)"
        assert metadata["python_version"] == "3.12+"
        assert metadata["postgresql_version"] == "18+ (required)"
        assert metadata["build_command"] == "GOEXPERIMENT=greenteagc,jsonv2 go build ./..."

    def test_parse_metadata_partial(self, tmp_path):
        """Test parsing with some metadata missing."""
        sot_content = """# Source of Truth

**Last Updated**: 2026-01-31
**Go Version**: 1.25+
"""
        sot_file = tmp_path / "partial.md"
        sot_file.write_text(sot_content)

        parser = SOTParser(sot_file)
        metadata = parser._parse_metadata()

        assert metadata["last_updated"] == "2026-01-31"
        assert metadata["go_version"] == "1.25+"
        assert "nodejs_version" not in metadata
        assert "python_version" not in metadata


class TestContentModulesParsing:
    """Test content modules table parsing."""

    def test_parse_content_modules(self, minimal_sot):
        """Test parsing content modules table."""
        parser = SOTParser(minimal_sot)
        modules = parser._parse_content_modules()

        assert len(modules) == 2

        # First module
        assert modules[0]["name"] == "Movie"
        assert modules[0]["schema"] == "`public`"
        assert modules[0]["status"] == "✅ Complete"
        assert modules[0]["primary_metadata"] == "[TMDb](../integrations/metadata/video/TMDB.md)"
        assert modules[0]["arr_integration"] == "[Radarr](../integrations/servarr/RADARR.md)"
        assert modules[0]["design_doc"] == "../features/video/MOVIE_MODULE.md"

        # Second module
        assert modules[1]["name"] == "TV Show"
        assert modules[1]["status"] == "🟡 Partial"

    def test_parse_content_modules_missing(self, tmp_path):
        """Test handling missing content modules section."""
        sot_file = tmp_path / "no_modules.md"
        sot_file.write_text("# Source of Truth\n\nNo modules here.")

        parser = SOTParser(sot_file)
        modules = parser._parse_content_modules()

        assert modules == []

    def test_parse_content_modules_malformed(self, tmp_path):
        """Test handling malformed table."""
        sot_content = """# Source of Truth

## Content Modules

| Module | Schema |
|--------|--------|
| Movie | public |
"""
        sot_file = tmp_path / "malformed.md"
        sot_file.write_text(sot_content)

        parser = SOTParser(sot_file)
        modules = parser._parse_content_modules()

        # Should skip rows with insufficient columns
        assert len(modules) == 0


class TestBackendServicesParsing:
    """Test backend services table parsing."""

    def test_parse_backend_services(self, minimal_sot):
        """Test parsing backend services table."""
        parser = SOTParser(minimal_sot)
        services = parser._parse_backend_services()

        assert len(services) == 2

        # First service
        assert services[0]["name"] == "Metadata"
        assert services[0]["package"] == "`internal/service/metadata`"
        assert services[0]["fx_module"] == "`MetadataModule`"
        assert services[0]["status"] == "✅ Complete"
        assert services[0]["design_doc"] == "../services/METADATA.md"

        # Second service
        assert services[1]["name"] == "Library"
        assert services[1]["status"] == "🟡 Partial"

    def test_parse_backend_services_missing(self, tmp_path):
        """Test handling missing backend services section."""
        sot_file = tmp_path / "no_services.md"
        sot_file.write_text("# Source of Truth\n")

        parser = SOTParser(sot_file)
        services = parser._parse_backend_services()

        assert services == []


class TestInfrastructureParsing:
    """Test infrastructure components table parsing."""

    def test_parse_infrastructure(self, minimal_sot):
        """Test parsing infrastructure table."""
        parser = SOTParser(minimal_sot)
        components = parser._parse_infrastructure()

        assert len(components) == 2

        # PostgreSQL
        assert components[0]["name"] == "PostgreSQL"
        assert components[0]["package"] == "`pgx`"
        assert components[0]["version"] == "18+"
        assert components[0]["purpose"] == "Primary database"
        assert components[0]["design_doc"] == "../technical/DATABASE.md"

        # Dragonfly
        assert components[1]["name"] == "Dragonfly"
        assert components[1]["package"] == "`rueidis`"

    def test_parse_infrastructure_missing(self, tmp_path):
        """Test handling missing infrastructure section."""
        sot_file = tmp_path / "no_infra.md"
        sot_file.write_text("# Source of Truth\n")

        parser = SOTParser(sot_file)
        components = parser._parse_infrastructure()

        assert components == []


class TestGoDependenciesParsing:
    """Test Go dependencies parsing."""

    def test_parse_go_dependencies_all_categories(self, minimal_sot):
        """Test parsing all Go dependency categories."""
        parser = SOTParser(minimal_sot)
        deps = parser._parse_go_dependencies()

        # Check all categories exist
        assert "core" in deps
        assert "security" in deps
        assert "observability" in deps
        assert "resilience" in deps
        assert "distributed" in deps

        # Core deps
        assert len(deps["core"]) == 2
        assert deps["core"][0]["package"] == "fx"
        assert deps["core"][0]["version"] == "1.22+"
        assert deps["core"][0]["purpose"] == "Dependency injection"
        assert deps["core"][0]["notes"] == "Core framework"

        assert deps["core"][1]["package"] == "ogen"

        # Security deps
        assert len(deps["security"]) == 1
        assert deps["security"][0]["package"] == "casbin"

        # Observability deps
        assert len(deps["observability"]) == 1
        assert deps["observability"][0]["package"] == "slog"

        # Resilience deps
        assert len(deps["resilience"]) == 1
        assert deps["resilience"][0]["package"] == "gobreaker"

        # Distributed deps
        assert len(deps["distributed"]) == 1
        assert deps["distributed"][0]["package"] == "raft"

    def test_parse_go_dependencies_backticks_removed(self, minimal_sot):
        """Test that backticks are removed from package names."""
        parser = SOTParser(minimal_sot)
        deps = parser._parse_go_dependencies()

        # Package names should not have backticks
        for category_deps in deps.values():
            for dep in category_deps:
                assert "`" not in dep["package"]

    def test_parse_go_dependencies_missing_category(self, tmp_path):
        """Test handling missing dependency category."""
        sot_content = """# Source of Truth

## Go Dependencies (Core)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `fx` | 1.22+ | DI | Framework |
"""
        sot_file = tmp_path / "partial_deps.md"
        sot_file.write_text(sot_content)

        parser = SOTParser(sot_file)
        deps = parser._parse_go_dependencies()

        # Core should have 1 dep
        assert len(deps["core"]) == 1

        # Other categories should be empty
        assert len(deps["security"]) == 0
        assert len(deps["observability"]) == 0
        assert len(deps["resilience"]) == 0
        assert len(deps["distributed"]) == 0


class TestDesignPrinciplesParsing:
    """Test design principles extraction."""

    def test_parse_design_principles(self, minimal_sot):
        """Test parsing design principles."""
        parser = SOTParser(minimal_sot)
        principles = parser._parse_design_principles()

        # Check text-based principles (parser extracts description after " - ")
        assert "database_strategy" in principles
        assert "No SQLite, MySQL, or other DB engines" in principles["database_strategy"]

        assert "package_update_policy" in principles
        assert "Stay one minor version behind latest stable" in principles["package_update_policy"]

        assert "test_coverage" in principles
        assert "All packages must maintain 80%+ test coverage" in principles["test_coverage"]

        # Check design patterns table
        assert "design_patterns" in principles
        patterns = principles["design_patterns"]
        assert len(patterns) == 3

        assert patterns[0]["pattern"] == "Database"
        assert patterns[0]["decision"] == "**PostgreSQL ONLY** - No SQLite, MySQL, or other DB engines"
        assert patterns[0]["notes"] == "Single DB for simplicity"

    def test_parse_design_principles_partial(self, tmp_path):
        """Test parsing with some principles missing."""
        sot_content = """# Source of Truth

**PostgreSQL ONLY** - No other databases allowed.
"""
        sot_file = tmp_path / "partial_principles.md"
        sot_file.write_text(sot_content)

        parser = SOTParser(sot_file)
        principles = parser._parse_design_principles()

        assert "database_strategy" in principles
        assert "package_update_policy" not in principles
        assert "test_coverage" not in principles


class TestLinkExtraction:
    """Test markdown link extraction."""

    def test_extract_link_with_markdown(self, minimal_sot):
        """Test extracting link from markdown format."""
        parser = SOTParser(minimal_sot)

        link = parser._extract_link("[TMDb](../integrations/metadata/video/TMDB.md)")
        assert link == "../integrations/metadata/video/TMDB.md"

    def test_extract_link_plain_text(self, minimal_sot):
        """Test extracting from plain text (no markdown)."""
        parser = SOTParser(minimal_sot)

        link = parser._extract_link("plain/path/to/file.md")
        assert link == "plain/path/to/file.md"

    def test_extract_link_with_whitespace(self, minimal_sot):
        """Test extracting link with extra whitespace."""
        parser = SOTParser(minimal_sot)

        link = parser._extract_link("  [Link](path.md)  ")
        assert link == "path.md"


class TestYAMLSaving:
    """Test YAML output generation."""

    def test_save_yaml_creates_file(self, minimal_sot, tmp_path):
        """Test saving parsed data to YAML file."""
        parser = SOTParser(minimal_sot)
        parser.data = {"test": "data", "nested": {"key": "value"}}

        output_file = tmp_path / "output.yaml"
        parser.save_yaml(output_file)

        assert output_file.exists()

        # Verify content
        with open(output_file) as f:
            loaded = yaml.safe_load(f)

        assert loaded == parser.data

    def test_save_yaml_creates_parent_dirs(self, minimal_sot, tmp_path):
        """Test saving YAML creates parent directories."""
        parser = SOTParser(minimal_sot)
        parser.data = {"test": "data"}

        output_file = tmp_path / "nested" / "deep" / "output.yaml"
        parser.save_yaml(output_file)

        assert output_file.exists()
        assert output_file.parent.exists()

    def test_save_yaml_unicode(self, minimal_sot, tmp_path):
        """Test saving YAML with unicode characters."""
        parser = SOTParser(minimal_sot)
        parser.data = {
            "status": "✅ Complete",
            "partial": "🟡 Partial",
            "not_started": "🔴 Not Started",
        }

        output_file = tmp_path / "unicode.yaml"
        parser.save_yaml(output_file)

        with open(output_file) as f:
            loaded = yaml.safe_load(f)

        assert loaded["status"] == "✅ Complete"
        assert loaded["partial"] == "🟡 Partial"


class TestFullParsing:
    """Test complete parsing workflow."""

    def test_parse_all_sections(self, minimal_sot):
        """Test parsing all sections."""
        parser = SOTParser(minimal_sot)
        data = parser.parse()

        # Check all expected sections exist
        assert "metadata" in data
        assert "content_modules" in data
        assert "backend_services" in data
        assert "infrastructure" in data
        assert "go_dependencies" in data
        assert "design_principles" in data

        # Verify counts
        assert len(data["content_modules"]) == 2
        assert len(data["backend_services"]) == 2
        assert len(data["infrastructure"]) == 2

        # Verify Go dependencies
        go_deps = data["go_dependencies"]
        total_deps = sum(len(deps) for deps in go_deps.values())
        assert total_deps == 6  # 2 core + 1 each from 4 other categories

    def test_parse_returns_data(self, minimal_sot):
        """Test that parse() returns the parsed data."""
        parser = SOTParser(minimal_sot)
        data = parser.parse()

        assert data == parser.data
        assert isinstance(data, dict)


class TestEdgeCases:
    """Test edge cases and error conditions."""

    def test_empty_file(self, tmp_path):
        """Test parsing empty SOT file."""
        sot_file = tmp_path / "empty.md"
        sot_file.write_text("")

        parser = SOTParser(sot_file)
        data = parser.parse()

        # Should return empty or minimal data
        assert isinstance(data, dict)
        assert len(data["content_modules"]) == 0
        assert len(data["backend_services"]) == 0

    def test_table_with_no_rows(self, tmp_path):
        """Test table with header but no data rows."""
        sot_content = """# Source of Truth

## Content Modules

| Module | Schema | Status | Primary Metadata | Arr Integration | Design Doc |
|--------|--------|--------|------------------|-----------------|------------|
"""
        sot_file = tmp_path / "empty_table.md"
        sot_file.write_text(sot_content)

        parser = SOTParser(sot_file)
        modules = parser._parse_content_modules()

        assert len(modules) == 0

    def test_malformed_markdown_links(self, minimal_sot):
        """Test handling malformed markdown links."""
        parser = SOTParser(minimal_sot)

        # Missing closing parenthesis
        link = parser._extract_link("[Text](incomplete")
        assert link == "[Text](incomplete"

        # No brackets
        link = parser._extract_link("just/a/path")
        assert link == "just/a/path"

    def test_special_characters_in_cells(self, tmp_path):
        """Test handling special characters in table cells."""
        sot_content = """# Source of Truth

## Content Modules

| Module | Schema | Status | Primary Metadata | Arr Integration | Design Doc |
|--------|--------|--------|------------------|-----------------|------------|
| Movie & TV | `public` | ✅ Complete | [TMDb](link.md) | [Arr](arr.md) | [Doc](doc.md) |
"""
        sot_file = tmp_path / "special.md"
        sot_file.write_text(sot_content)

        parser = SOTParser(sot_file)
        modules = parser._parse_content_modules()

        assert len(modules) == 1
        assert modules[0]["name"] == "Movie & TV"


class TestIntegration:
    """Test full integration scenarios."""

    def test_full_workflow(self, minimal_sot, tmp_path):
        """Test complete parse → save workflow."""
        parser = SOTParser(minimal_sot)
        data = parser.parse()

        output_file = tmp_path / "shared-sot.yaml"
        parser.save_yaml(output_file)

        # Load saved file and verify
        with open(output_file) as f:
            loaded = yaml.safe_load(f)

        assert loaded == data
        assert len(loaded["content_modules"]) == 2
        assert len(loaded["backend_services"]) == 2


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
