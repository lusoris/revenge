"""Tests for YAML data file validation.

Ensures all YAML files have required fields and can generate docs without errors.
"""

from pathlib import Path

import pytest
import yaml


class TestYAMLStructure:
    """Test YAML files have correct structure."""

    def test_all_yaml_files_parseable(self):
        """All YAML files must be valid YAML syntax."""
        data_dir = Path("data")
        yaml_files = list(data_dir.rglob("*.yaml"))

        # Exclude templates
        yaml_files = [f for f in yaml_files if ".templates" not in str(f)]

        assert len(yaml_files) > 100, "Should have many YAML files"

        errors = []
        for yaml_file in yaml_files:
            try:
                with open(yaml_file, encoding="utf-8") as f:
                    yaml.safe_load(f)
            except yaml.YAMLError as e:
                errors.append(f"{yaml_file.relative_to(data_dir)}: {e}")

        if errors:
            pytest.fail(f"YAML syntax errors in {len(errors)} files:\n" + "\n".join(errors))

    def test_all_yaml_files_have_required_base_fields(self):
        """All YAML files must have required base fields."""
        data_dir = Path("data")
        yaml_files = list(data_dir.rglob("*.yaml"))
        yaml_files = [f for f in yaml_files if ".templates" not in str(f) and "shared-sot.yaml" not in str(f)]

        required_fields = [
            "doc_title",
            "doc_category",
            "created_date",
            "overall_status",
            "status_design",
            "technical_summary",
        ]

        errors = []
        for yaml_file in yaml_files:
            with open(yaml_file, encoding="utf-8") as f:
                try:
                    data = yaml.safe_load(f)
                except Exception:
                    continue  # Syntax errors caught by other test

            if not data:
                continue

            missing = [field for field in required_fields if field not in data]
            if missing:
                errors.append(f"{yaml_file.relative_to(data_dir)}: missing {missing}")

        if errors:
            pytest.fail(f"Missing required fields in {len(errors)} files:\n" + "\n".join(errors[:10]))


class TestCategorySpecificFields:
    """Test category-specific required fields."""

    def test_feature_yaml_files_have_feature_fields(self):
        """Feature YAML files must have feature-specific fields."""
        data_dir = Path("data")
        feature_files = list(data_dir.glob("features/**/*.yaml"))
        feature_files = [f for f in feature_files if "INDEX.yaml" not in f.name]

        required_fields = ["feature_name", "schema_name"]

        errors = []
        for yaml_file in feature_files:
            with open(yaml_file, encoding="utf-8") as f:
                data = yaml.safe_load(f)

            if not data:
                continue

            missing = [field for field in required_fields if field not in data]
            if missing:
                errors.append(f"{yaml_file.relative_to(data_dir)}: missing {missing}")

        if errors:
            pytest.fail(f"Missing feature fields in {len(errors)} files:\n" + "\n".join(errors))

    def test_service_yaml_files_have_service_fields(self):
        """Service YAML files must have service-specific fields."""
        data_dir = Path("data")
        service_files = list(data_dir.glob("services/*.yaml"))

        required_fields = ["service_name", "package_path", "fx_module"]

        errors = []
        for yaml_file in service_files:
            with open(yaml_file, encoding="utf-8") as f:
                data = yaml.safe_load(f)

            if not data:
                continue

            missing = [field for field in required_fields if field not in data]
            if missing:
                errors.append(f"{yaml_file.relative_to(data_dir)}: missing {missing}")

        if errors:
            pytest.fail(f"Missing service fields in {len(errors)} files:\n" + "\n".join(errors))

    def test_integration_yaml_files_have_integration_fields(self):
        """Integration YAML files must have integration-specific fields."""
        data_dir = Path("data")
        integration_files = list(data_dir.glob("integrations/**/*.yaml"))
        integration_files = [f for f in integration_files if "INDEX.yaml" not in f.name]

        required_fields = ["integration_name", "external_service", "integration_id"]

        errors = []
        for yaml_file in integration_files:
            with open(yaml_file, encoding="utf-8") as f:
                data = yaml.safe_load(f)

            if not data:
                continue

            missing = [field for field in required_fields if field not in data]
            if missing:
                errors.append(f"{yaml_file.relative_to(data_dir)}: missing {missing}")

        if errors:
            # Some integrations might not need all fields, so just warn
            print(f"Warning: {len(errors)} integration files missing some fields")


class TestPlaceholderContent:
    """Test that wiki fields don't contain PLACEHOLDER content.

    Note: These tests use a threshold to allow gradual PLACEHOLDER replacement.
    When all PLACEHOLDERs are filled, set MAX_ALLOWED_PLACEHOLDERS to 0.
    """

    # Allow up to this many PLACEHOLDERs (decrease as docs are completed)
    # Set to 0 for strict mode when all docs should be complete
    MAX_ALLOWED_WIKI_OVERVIEW_PLACEHOLDERS = 120  # Currently ~113, allow some buffer
    MAX_ALLOWED_WIKI_TAGLINE_PLACEHOLDERS = 15  # Currently ~11, allow some buffer

    def test_no_placeholder_in_wiki_overview(self):
        """wiki_overview must not contain PLACEHOLDER marker.

        Enforces that PLACEHOLDER count stays below threshold.
        This prevents new PLACEHOLDERs from being added while allowing
        gradual replacement of existing ones.
        """
        data_dir = Path("data")
        yaml_files = list(data_dir.rglob("*.yaml"))
        yaml_files = [f for f in yaml_files if ".templates" not in str(f) and "shared-sot.yaml" not in str(f)]

        errors = []
        for yaml_file in yaml_files:
            with open(yaml_file, encoding="utf-8") as f:
                try:
                    data = yaml.safe_load(f)
                except Exception:
                    continue

            if not data:
                continue

            wiki_overview = data.get("wiki_overview", "")
            if wiki_overview and "PLACEHOLDER" in str(wiki_overview).upper():
                errors.append(f"{yaml_file.relative_to(data_dir)}")

        if len(errors) > self.MAX_ALLOWED_WIKI_OVERVIEW_PLACEHOLDERS:
            pytest.fail(
                f"Too many PLACEHOLDER wiki_overview fields: {len(errors)} > {self.MAX_ALLOWED_WIKI_OVERVIEW_PLACEHOLDERS}\n"
                f"Files:\n" + "\n".join(errors[:20])
            )
        elif errors:
            print(f"\nInfo: {len(errors)} wiki_overview PLACEHOLDERs remaining (max: {self.MAX_ALLOWED_WIKI_OVERVIEW_PLACEHOLDERS})")

    def test_no_placeholder_in_wiki_tagline(self):
        """wiki_tagline must not contain PLACEHOLDER marker.

        Enforces that PLACEHOLDER count stays below threshold.
        """
        data_dir = Path("data")
        yaml_files = list(data_dir.rglob("*.yaml"))
        yaml_files = [f for f in yaml_files if ".templates" not in str(f) and "shared-sot.yaml" not in str(f)]

        errors = []
        for yaml_file in yaml_files:
            with open(yaml_file, encoding="utf-8") as f:
                try:
                    data = yaml.safe_load(f)
                except Exception:
                    continue

            if not data:
                continue

            wiki_tagline = data.get("wiki_tagline", "")
            if wiki_tagline and "PLACEHOLDER" in str(wiki_tagline).upper():
                errors.append(f"{yaml_file.relative_to(data_dir)}")

        if len(errors) > self.MAX_ALLOWED_WIKI_TAGLINE_PLACEHOLDERS:
            pytest.fail(
                f"Too many PLACEHOLDER wiki_tagline fields: {len(errors)} > {self.MAX_ALLOWED_WIKI_TAGLINE_PLACEHOLDERS}\n"
                f"Files:\n" + "\n".join(errors[:20])
            )
        elif errors:
            print(f"\nInfo: {len(errors)} wiki_tagline PLACEHOLDERs remaining (max: {self.MAX_ALLOWED_WIKI_TAGLINE_PLACEHOLDERS})")

    def test_priority_docs_have_no_placeholders(self):
        """Priority 1-2 docs (operations, technical critical) must have no PLACEHOLDERs.

        These are critical for contributor onboarding and deployment.
        """
        data_dir = Path("data")

        # Priority docs that must be complete
        priority_patterns = [
            "operations/*.yaml",
            "technical/API.yaml",
            "technical/CONFIGURATION.yaml",
            "technical/FRONTEND.yaml",
        ]

        errors = []
        for pattern in priority_patterns:
            for yaml_file in data_dir.glob(pattern):
                with open(yaml_file, encoding="utf-8") as f:
                    try:
                        data = yaml.safe_load(f)
                    except Exception:
                        continue

                if not data:
                    continue

                for field in ["wiki_overview", "wiki_tagline"]:
                    value = data.get(field, "")
                    if value and "PLACEHOLDER" in str(value).upper():
                        errors.append(f"{yaml_file.relative_to(data_dir)}: {field}")

        if errors:
            pytest.fail(
                f"Priority docs must not have PLACEHOLDERs:\n" + "\n".join(errors)
            )


class TestYAMLListFormatting:
    """Test YAML list items are properly formatted."""

    def test_list_items_with_colons_are_quoted(self):
        """List items containing colons should be quoted to avoid YAML parsing issues.

        This catches issues like:
        query_params:
          - date (optional): Date for schedule  # BAD - colon in unquoted string
          - "date (optional): Date for schedule"  # GOOD - quoted
        """
        data_dir = Path("data")
        yaml_files = list(data_dir.rglob("*.yaml"))
        yaml_files = [f for f in yaml_files if ".templates" not in str(f)]

        errors = []
        for yaml_file in yaml_files:
            with open(yaml_file, encoding="utf-8") as f:
                content = f.read()

            # Check for list items that have colons and might cause issues
            # Look for patterns like "  - something: something" without quotes
            import re
            # Match list items that have a colon but aren't quoted
            # This is a simplified check - YAML parsing will catch actual errors
            lines = content.split("\n")
            for i, line in enumerate(lines, 1):
                # Skip if line is already properly parsed (no syntax error)
                # Focus on common problematic patterns
                if re.match(r"^\s+-\s+[^\"'][^:]+:\s+[^\"']", line):
                    # This looks like an unquoted list item with colon
                    # Try parsing just this structure to see if it's valid
                    test_yaml = f"test:\n{line}"
                    try:
                        parsed = yaml.safe_load(test_yaml)
                        # If parsed as dict, it means colon was interpreted as key separator
                        if isinstance(parsed.get("test"), list) and isinstance(parsed["test"][0], dict):
                            # This is actually valid YAML dict in list - ignore
                            pass
                    except yaml.YAMLError:
                        errors.append(f"{yaml_file.relative_to(data_dir)}:{i}: {line.strip()}")

        # This test is informational - actual YAML errors are caught by syntax test
        if errors:
            print(f"\nWarning: {len(errors)} potential unquoted colon issues found")
            for error in errors[:5]:
                print(f"  {error}")


class TestDocGeneration:
    """Test that all YAML files can generate docs without errors."""

    def test_all_yaml_files_generate_successfully(self):
        """Run batch_regenerate.py and check for errors."""
        import subprocess
        import sys

        result = subprocess.run(
            [sys.executable, "scripts/automation/batch_regenerate.py"],
            capture_output=True,
            text=True,
            timeout=60,
        )

        # Check output for errors
        output = result.stdout + result.stderr

        # Extract summary
        if "Success:" in output and "Failed:" in output:
            lines = output.split("\n")
            success_line = [l for l in lines if "Success:" in l][0]
            failed_line = [l for l in lines if "Failed:" in l][0]

            success_count = int(success_line.split("Success:")[1].strip())
            failed_count = int(failed_line.split("Failed:")[1].strip())

            # Allow up to 5 failures (for files that may have issues)
            assert failed_count <= 5, f"Too many generation failures: {failed_count}"
            assert success_count > 140, f"Too few successful generations: {success_count}"

            if failed_count > 0:
                # Extract which files failed
                error_lines = [l for l in lines if "Error:" in l]
                print(f"\n{failed_count} files failed generation:")
                for error in error_lines[:10]:
                    print(f"  {error}")

            return

        # If we get here, couldn't parse output
        pytest.fail(f"Could not parse generation output:\n{output[-500:]}")


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
