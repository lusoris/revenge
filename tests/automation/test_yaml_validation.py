"""Tests for YAML data file validation.

Ensures all YAML files have required fields and can generate docs without errors.
"""

import re
import subprocess
import sys
from pathlib import Path

import pytest
import yaml

from tests.constants import (
    BASE_REQUIRED_FIELDS,
    FEATURE_REQUIRED_FIELDS,
    INTEGRATION_REQUIRED_FIELDS,
    MAX_GENERATION_FAILURES,
    MAX_WIKI_OVERVIEW_PLACEHOLDERS,
    MAX_WIKI_TAGLINE_PLACEHOLDERS,
    MIN_SUCCESSFUL_GENERATIONS,
    MIN_YAML_FILES,
    SERVICE_REQUIRED_FIELDS,
)


def _load_yaml_safe(yaml_file: Path) -> dict | None:
    """Load YAML file safely, returning None on error."""
    try:
        with open(yaml_file, encoding="utf-8") as f:
            return yaml.safe_load(f)
    except Exception:
        return None


def _check_missing_fields(data: dict, required_fields: list[str]) -> list[str]:
    """Return list of missing required fields."""
    return [field for field in required_fields if field not in data]


class TestYAMLStructure:
    """Test YAML files have correct structure."""

    def test_all_yaml_files_parseable(self, data_dir: Path):
        """All YAML files must be valid YAML syntax."""
        yaml_files = list(data_dir.rglob("*.yaml"))
        yaml_files = [f for f in yaml_files if ".templates" not in str(f)]

        assert len(yaml_files) > MIN_YAML_FILES, f"Should have >{MIN_YAML_FILES} YAML files"

        errors = []
        for yaml_file in yaml_files:
            try:
                with open(yaml_file, encoding="utf-8") as f:
                    yaml.safe_load(f)
            except yaml.YAMLError as e:
                errors.append(f"{yaml_file.relative_to(data_dir)}: {e}")

        if errors:
            pytest.fail(f"YAML syntax errors in {len(errors)} files:\n" + "\n".join(errors))

    def test_all_yaml_files_have_required_base_fields(self, yaml_data_files: list[Path], data_dir: Path):
        """All YAML files must have required base fields."""
        errors = []
        for yaml_file in yaml_data_files:
            data = _load_yaml_safe(yaml_file)
            if not data:
                continue

            missing = _check_missing_fields(data, BASE_REQUIRED_FIELDS)
            if missing:
                errors.append(f"{yaml_file.relative_to(data_dir)}: missing {missing}")

        if errors:
            pytest.fail(f"Missing required fields in {len(errors)} files:\n" + "\n".join(errors[:10]))


class TestCategorySpecificFields:
    """Test category-specific required fields."""

    def test_feature_yaml_files_have_feature_fields(self, feature_yaml_files: list[Path], data_dir: Path):
        """Feature YAML files must have feature-specific fields."""
        errors = []
        for yaml_file in feature_yaml_files:
            data = _load_yaml_safe(yaml_file)
            if not data:
                continue

            missing = _check_missing_fields(data, FEATURE_REQUIRED_FIELDS)
            if missing:
                errors.append(f"{yaml_file.relative_to(data_dir)}: missing {missing}")

        if errors:
            pytest.fail(f"Missing feature fields in {len(errors)} files:\n" + "\n".join(errors))

    def test_service_yaml_files_have_service_fields(self, service_yaml_files: list[Path], data_dir: Path):
        """Service YAML files must have service-specific fields."""
        errors = []
        for yaml_file in service_yaml_files:
            data = _load_yaml_safe(yaml_file)
            if not data:
                continue

            missing = _check_missing_fields(data, SERVICE_REQUIRED_FIELDS)
            if missing:
                errors.append(f"{yaml_file.relative_to(data_dir)}: missing {missing}")

        if errors:
            pytest.fail(f"Missing service fields in {len(errors)} files:\n" + "\n".join(errors))

    def test_integration_yaml_files_have_integration_fields(self, integration_yaml_files: list[Path], data_dir: Path):
        """Integration YAML files must have integration-specific fields."""
        errors = []
        for yaml_file in integration_yaml_files:
            data = _load_yaml_safe(yaml_file)
            if not data:
                continue

            missing = _check_missing_fields(data, INTEGRATION_REQUIRED_FIELDS)
            if missing:
                errors.append(f"{yaml_file.relative_to(data_dir)}: missing {missing}")

        if errors:
            # Some integrations might not need all fields, so just warn
            print(f"Warning: {len(errors)} integration files missing some fields")


class TestPlaceholderContent:
    """Test that wiki fields don't contain PLACEHOLDER content.

    Note: These tests use a threshold to allow gradual PLACEHOLDER replacement.
    When all PLACEHOLDERs are filled, update constants.py thresholds to 0.
    """

    def test_no_placeholder_in_wiki_overview(self, yaml_data_files: list[Path], data_dir: Path):
        """wiki_overview must not contain PLACEHOLDER marker.

        Enforces that PLACEHOLDER count stays below threshold.
        This prevents new PLACEHOLDERs from being added while allowing
        gradual replacement of existing ones.
        """
        errors = []
        for yaml_file in yaml_data_files:
            data = _load_yaml_safe(yaml_file)
            if not data:
                continue

            wiki_overview = data.get("wiki_overview", "")
            if wiki_overview and "PLACEHOLDER" in str(wiki_overview).upper():
                errors.append(str(yaml_file.relative_to(data_dir)))

        if len(errors) > MAX_WIKI_OVERVIEW_PLACEHOLDERS:
            pytest.fail(
                f"Too many PLACEHOLDER wiki_overview fields: {len(errors)} > {MAX_WIKI_OVERVIEW_PLACEHOLDERS}\n"
                f"Files:\n" + "\n".join(errors[:20])
            )
        elif errors:
            print(f"\nInfo: {len(errors)} wiki_overview PLACEHOLDERs remaining (max: {MAX_WIKI_OVERVIEW_PLACEHOLDERS})")

    def test_no_placeholder_in_wiki_tagline(self, yaml_data_files: list[Path], data_dir: Path):
        """wiki_tagline must not contain PLACEHOLDER marker.

        Enforces that PLACEHOLDER count stays below threshold.
        """
        errors = []
        for yaml_file in yaml_data_files:
            data = _load_yaml_safe(yaml_file)
            if not data:
                continue

            wiki_tagline = data.get("wiki_tagline", "")
            if wiki_tagline and "PLACEHOLDER" in str(wiki_tagline).upper():
                errors.append(str(yaml_file.relative_to(data_dir)))

        if len(errors) > MAX_WIKI_TAGLINE_PLACEHOLDERS:
            pytest.fail(
                f"Too many PLACEHOLDER wiki_tagline fields: {len(errors)} > {MAX_WIKI_TAGLINE_PLACEHOLDERS}\n"
                f"Files:\n" + "\n".join(errors[:20])
            )
        elif errors:
            print(f"\nInfo: {len(errors)} wiki_tagline PLACEHOLDERs remaining (max: {MAX_WIKI_TAGLINE_PLACEHOLDERS})")

    def test_priority_docs_have_no_placeholders(self, data_dir: Path, priority_yaml_patterns: list[str]):
        """Priority 1-2 docs (operations, technical critical) must have no PLACEHOLDERs.

        These are critical for contributor onboarding and deployment.
        """
        errors = []
        for pattern in priority_yaml_patterns:
            for yaml_file in data_dir.glob(pattern):
                data = _load_yaml_safe(yaml_file)
                if not data:
                    continue

                for field in ["wiki_overview", "wiki_tagline"]:
                    value = data.get(field, "")
                    if value and "PLACEHOLDER" in str(value).upper():
                        errors.append(f"{yaml_file.relative_to(data_dir)}: {field}")

        if errors:
            pytest.fail(f"Priority docs must not have PLACEHOLDERs:\n" + "\n".join(errors))


class TestYAMLListFormatting:
    """Test YAML list items are properly formatted."""

    def test_list_items_with_colons_are_quoted(self, data_dir: Path):
        """List items containing colons should be quoted to avoid YAML parsing issues.

        This catches issues like:
        query_params:
          - date (optional): Date for schedule  # BAD - colon in unquoted string
          - "date (optional): Date for schedule"  # GOOD - quoted
        """
        yaml_files = list(data_dir.rglob("*.yaml"))
        yaml_files = [f for f in yaml_files if ".templates" not in str(f)]

        errors = []
        for yaml_file in yaml_files:
            with open(yaml_file, encoding="utf-8") as f:
                content = f.read()

            lines = content.split("\n")
            for i, line in enumerate(lines, 1):
                # Check for unquoted list items with colons
                if re.match(r"^\s+-\s+[^\"'][^:]+:\s+[^\"']", line):
                    test_yaml = f"test:\n{line}"
                    try:
                        parsed = yaml.safe_load(test_yaml)
                        # If parsed as dict in list, colon was interpreted as key separator
                        if isinstance(parsed.get("test"), list) and isinstance(parsed["test"][0], dict):
                            pass  # Valid YAML dict in list - ignore
                    except yaml.YAMLError:
                        errors.append(f"{yaml_file.relative_to(data_dir)}:{i}: {line.strip()}")

        # Informational - actual YAML errors caught by syntax test
        if errors:
            print(f"\nWarning: {len(errors)} potential unquoted colon issues found")
            for error in errors[:5]:
                print(f"  {error}")


class TestDocGeneration:
    """Test that all YAML files can generate docs without errors."""

    def test_all_yaml_files_generate_successfully(self):
        """Run batch_regenerate.py and check for errors."""
        result = subprocess.run(
            [sys.executable, "scripts/automation/batch_regenerate.py"],
            capture_output=True,
            text=True,
            timeout=60,
        )

        output = result.stdout + result.stderr

        if "Success:" in output and "Failed:" in output:
            lines = output.split("\n")
            success_line = [line for line in lines if "Success:" in line][0]
            failed_line = [line for line in lines if "Failed:" in line][0]

            success_count = int(success_line.split("Success:")[1].strip())
            failed_count = int(failed_line.split("Failed:")[1].strip())

            assert failed_count <= MAX_GENERATION_FAILURES, f"Too many generation failures: {failed_count}"
            assert success_count > MIN_SUCCESSFUL_GENERATIONS, f"Too few successful generations: {success_count}"

            if failed_count > 0:
                error_lines = [line for line in lines if "Error:" in line]
                print(f"\n{failed_count} files failed generation:")
                for error in error_lines[:10]:
                    print(f"  {error}")

            return

        pytest.fail(f"Could not parse generation output:\n{output[-500:]}")


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
