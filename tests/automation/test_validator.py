#!/usr/bin/env python3
"""Tests for YAML validator.

Test coverage for validator.py:
- Initialization and schema loading
- Single file validation (valid, invalid, errors)
- Directory validation (recursive scanning)
- doc_category detection and routing
- YAML parsing error handling
- JSON schema validation errors
- Results printing and formatting
- Edge cases (missing fields, unknown categories)
"""

import json

import pytest
import yaml

from scripts.automation.validator import YAMLValidator


@pytest.fixture
def validator_setup(tmp_path):
    """Set up validator test environment."""
    # Create directory structure
    (tmp_path / "schemas").mkdir()
    (tmp_path / "data").mkdir()

    # Create minimal feature schema
    feature_schema = {
        "$schema": "http://json-schema.org/draft-07/schema#",
        "type": "object",
        "properties": {
            "doc_category": {"type": "string", "enum": ["feature"]},
            "doc_title": {"type": "string"},
            "status": {"type": "string"},
        },
        "required": ["doc_category", "doc_title", "status"],
    }
    (tmp_path / "schemas" / "feature.schema.json").write_text(
        json.dumps(feature_schema, indent=2)
    )

    # Create minimal service schema
    service_schema = {
        "$schema": "http://json-schema.org/draft-07/schema#",
        "type": "object",
        "properties": {
            "doc_category": {"type": "string", "enum": ["service"]},
            "doc_title": {"type": "string"},
            "package": {"type": "string"},
        },
        "required": ["doc_category", "doc_title", "package"],
    }
    (tmp_path / "schemas" / "service.schema.json").write_text(
        json.dumps(service_schema, indent=2)
    )

    return tmp_path


class TestInitialization:
    """Test YAMLValidator initialization."""

    def test_init_sets_paths(self, validator_setup):
        """Test initialization sets correct paths."""
        validator = YAMLValidator(validator_setup)

        assert validator.repo_root == validator_setup
        assert validator.schemas_dir == validator_setup / "schemas"
        assert validator.data_dir == validator_setup / "data"

    def test_init_loads_schemas(self, validator_setup):
        """Test initialization loads schemas."""
        validator = YAMLValidator(validator_setup)

        assert "feature" in validator.schemas
        assert "service" in validator.schemas
        assert len(validator.schemas) == 2

    def test_init_with_missing_schema(self, tmp_path):
        """Test initialization with missing schema files."""
        (tmp_path / "schemas").mkdir()
        (tmp_path / "data").mkdir()

        validator = YAMLValidator(tmp_path)

        # Should initialize with empty schemas dict
        assert len(validator.schemas) == 0

    def test_init_creates_validators(self, validator_setup):
        """Test that schema validators are created."""
        from jsonschema import Draft7Validator

        validator = YAMLValidator(validator_setup)

        # Validators should be Draft7Validator instances
        assert isinstance(validator.schemas["feature"], Draft7Validator)
        assert isinstance(validator.schemas["service"], Draft7Validator)


class TestSingleFileValidation:
    """Test single file validation."""

    def test_validate_valid_feature(self, validator_setup):
        """Test validating a valid feature file."""
        # Create valid feature YAML
        feature_data = {
            "doc_category": "feature",
            "doc_title": "Test Feature",
            "status": "Complete",
        }
        yaml_file = validator_setup / "data" / "test_feature.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        is_valid, errors = validator.validate_file(yaml_file)

        assert is_valid is True
        assert errors == []

    def test_validate_invalid_feature_missing_field(self, validator_setup):
        """Test validating feature with missing required field."""
        # Create feature missing required field
        feature_data = {
            "doc_category": "feature",
            "doc_title": "Test Feature",
            # Missing 'status'
        }
        yaml_file = validator_setup / "data" / "invalid.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        is_valid, errors = validator.validate_file(yaml_file)

        assert is_valid is False
        assert len(errors) > 0
        assert any("status" in error.lower() for error in errors)

    def test_validate_invalid_feature_wrong_type(self, validator_setup):
        """Test validating feature with wrong field type."""
        # Create feature with wrong type
        feature_data = {
            "doc_category": "feature",
            "doc_title": 123,  # Should be string
            "status": "Complete",
        }
        yaml_file = validator_setup / "data" / "wrong_type.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        is_valid, errors = validator.validate_file(yaml_file)

        assert is_valid is False
        assert len(errors) > 0

    def test_validate_missing_doc_category(self, validator_setup):
        """Test validating file with missing doc_category."""
        # Create YAML without doc_category
        data = {
            "doc_title": "Test",
            "status": "Complete",
        }
        yaml_file = validator_setup / "data" / "no_category.yaml"
        yaml_file.write_text(yaml.dump(data))

        validator = YAMLValidator(validator_setup)
        is_valid, errors = validator.validate_file(yaml_file)

        assert is_valid is False
        assert len(errors) == 1
        assert "Missing 'doc_category' field" in errors[0]

    def test_validate_unknown_doc_category(self, validator_setup):
        """Test validating file with unknown doc_category."""
        # Create YAML with unknown category
        data = {
            "doc_category": "unknown_type",
            "doc_title": "Test",
        }
        yaml_file = validator_setup / "data" / "unknown.yaml"
        yaml_file.write_text(yaml.dump(data))

        validator = YAMLValidator(validator_setup)
        is_valid, errors = validator.validate_file(yaml_file)

        assert is_valid is False
        assert len(errors) == 1
        assert "Unknown doc_category" in errors[0]

    def test_validate_invalid_yaml_syntax(self, validator_setup):
        """Test validating file with invalid YAML syntax."""
        # Create file with invalid YAML
        yaml_file = validator_setup / "data" / "invalid.yaml"
        yaml_file.write_text("invalid: yaml: syntax:")

        validator = YAMLValidator(validator_setup)
        is_valid, errors = validator.validate_file(yaml_file)

        assert is_valid is False
        assert len(errors) == 1
        assert "YAML parsing error" in errors[0]

    def test_validate_service_file(self, validator_setup):
        """Test validating a service file."""
        # Create valid service YAML
        service_data = {
            "doc_category": "service",
            "doc_title": "Test Service",
            "package": "internal/service/test",
        }
        yaml_file = validator_setup / "data" / "test_service.yaml"
        yaml_file.write_text(yaml.dump(service_data))

        validator = YAMLValidator(validator_setup)
        is_valid, errors = validator.validate_file(yaml_file)

        assert is_valid is True
        assert errors == []

    def test_validate_reports_multiple_errors(self, validator_setup):
        """Test that validation reports all errors."""
        # Create feature with multiple issues
        feature_data = {
            "doc_category": "feature",
            # Missing doc_title and status
        }
        yaml_file = validator_setup / "data" / "multi_error.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        is_valid, errors = validator.validate_file(yaml_file)

        assert is_valid is False
        # Should have errors for missing doc_title and status
        assert len(errors) >= 2


class TestDirectoryValidation:
    """Test directory validation."""

    def test_validate_directory_single_file(self, validator_setup):
        """Test validating directory with single file."""
        # Create one valid file
        feature_data = {
            "doc_category": "feature",
            "doc_title": "Test",
            "status": "Complete",
        }
        yaml_file = validator_setup / "data" / "test.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        results = validator.validate_directory(validator_setup / "data")

        assert len(results) == 1
        file_path = next(iter(results.keys()))
        is_valid, _errors = results[file_path]
        assert is_valid is True

    def test_validate_directory_multiple_files(self, validator_setup):
        """Test validating directory with multiple files."""
        # Create two valid files
        for i in range(2):
            feature_data = {
                "doc_category": "feature",
                "doc_title": f"Feature {i}",
                "status": "Complete",
            }
            yaml_file = validator_setup / "data" / f"feature_{i}.yaml"
            yaml_file.write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        results = validator.validate_directory(validator_setup / "data")

        assert len(results) == 2
        assert all(is_valid for is_valid, _ in results.values())

    def test_validate_directory_recursive(self, validator_setup):
        """Test validating directory recursively."""
        # Create nested structure
        (validator_setup / "data" / "features").mkdir()
        (validator_setup / "data" / "services").mkdir()

        # Create file in subdirectory
        feature_data = {
            "doc_category": "feature",
            "doc_title": "Nested Feature",
            "status": "Complete",
        }
        (validator_setup / "data" / "features" / "nested.yaml").write_text(
            yaml.dump(feature_data)
        )

        validator = YAMLValidator(validator_setup)
        results = validator.validate_directory(validator_setup / "data")

        assert len(results) == 1
        assert "data/features/nested.yaml" in next(iter(results.keys()))

    def test_validate_directory_skips_shared_sot(self, validator_setup):
        """Test that validation skips shared-sot.yaml."""
        # Create shared-sot.yaml (special file, no doc_category)
        sot_data = {"metadata": {"version": "1.0"}}
        (validator_setup / "data" / "shared-sot.yaml").write_text(yaml.dump(sot_data))

        # Create regular file
        feature_data = {
            "doc_category": "feature",
            "doc_title": "Test",
            "status": "Complete",
        }
        (validator_setup / "data" / "test.yaml").write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        results = validator.validate_directory(validator_setup / "data")

        # Should only validate test.yaml, skip shared-sot.yaml
        assert len(results) == 1
        assert "shared-sot.yaml" not in str(results.keys())

    def test_validate_directory_mixed_valid_invalid(self, validator_setup):
        """Test validating directory with mix of valid and invalid files."""
        # Create valid file
        valid_data = {
            "doc_category": "feature",
            "doc_title": "Valid",
            "status": "Complete",
        }
        (validator_setup / "data" / "valid.yaml").write_text(yaml.dump(valid_data))

        # Create invalid file
        invalid_data = {
            "doc_category": "feature",
            "doc_title": "Invalid",
            # Missing status
        }
        (validator_setup / "data" / "invalid.yaml").write_text(yaml.dump(invalid_data))

        validator = YAMLValidator(validator_setup)
        results = validator.validate_directory(validator_setup / "data")

        assert len(results) == 2

        # Check valid file (exact path match to avoid matching "invalid.yaml")
        valid_result = results["data/valid.yaml"]
        assert valid_result[0] is True

        # Check invalid file
        invalid_result = results["data/invalid.yaml"]
        assert invalid_result[0] is False

    def test_validate_directory_empty(self, validator_setup):
        """Test validating empty directory."""
        validator = YAMLValidator(validator_setup)
        results = validator.validate_directory(validator_setup / "data")

        assert len(results) == 0

    def test_validate_directory_supports_yml_extension(self, validator_setup):
        """Test that validation supports both .yaml and .yml extensions."""
        # Create .yml file
        feature_data = {
            "doc_category": "feature",
            "doc_title": "Test",
            "status": "Complete",
        }
        (validator_setup / "data" / "test.yml").write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        results = validator.validate_directory(validator_setup / "data")

        assert len(results) == 1
        assert any(".yml" in path for path in results)


class TestResultsPrinting:
    """Test results printing."""

    def test_print_results_all_valid(self, validator_setup, capsys):
        """Test printing results when all files are valid."""
        results = {
            "data/test1.yaml": (True, []),
            "data/test2.yaml": (True, []),
        }

        validator = YAMLValidator(validator_setup)
        all_valid = validator.print_results(results)

        assert all_valid is True

        # Check output
        captured = capsys.readouterr()
        assert "✅ All 2 files passed validation!" in captured.out
        assert "VALID FILES (2)" in captured.out

    def test_print_results_some_invalid(self, validator_setup, capsys):
        """Test printing results with some invalid files."""
        results = {
            "data/valid.yaml": (True, []),
            "data/invalid.yaml": (False, ["  • root: Missing required field 'status'"]),
        }

        validator = YAMLValidator(validator_setup)
        all_valid = validator.print_results(results)

        assert all_valid is False

        captured = capsys.readouterr()
        assert "❌ INVALID FILES (1)" in captured.out
        assert "data/invalid.yaml" in captured.out
        assert "✅ VALID FILES (1)" in captured.out
        assert "❌ 1 file(s) failed validation" in captured.out

    def test_print_results_all_invalid(self, validator_setup, capsys):
        """Test printing results when all files are invalid."""
        results = {
            "data/test1.yaml": (False, ["Error 1"]),
            "data/test2.yaml": (False, ["Error 2"]),
        }

        validator = YAMLValidator(validator_setup)
        all_valid = validator.print_results(results)

        assert all_valid is False

        captured = capsys.readouterr()
        assert "❌ INVALID FILES (2)" in captured.out
        assert "❌ 2 file(s) failed validation" in captured.out


class TestEdgeCases:
    """Test edge cases and error conditions."""

    def test_empty_yaml_file(self, validator_setup):
        """Test validating empty YAML file."""
        yaml_file = validator_setup / "data" / "empty.yaml"
        yaml_file.write_text("")

        validator = YAMLValidator(validator_setup)
        is_valid, _errors = validator.validate_file(yaml_file)

        # Empty YAML loads as None, which doesn't have .get()
        assert is_valid is False

    def test_yaml_with_null_doc_category(self, validator_setup):
        """Test validating YAML with null doc_category."""
        data = {
            "doc_category": None,
            "doc_title": "Test",
        }
        yaml_file = validator_setup / "data" / "null_category.yaml"
        yaml_file.write_text(yaml.dump(data))

        validator = YAMLValidator(validator_setup)
        is_valid, _errors = validator.validate_file(yaml_file)

        # None should be treated as missing
        assert is_valid is False

    def test_malformed_json_schema(self, tmp_path):
        """Test handling malformed JSON schema."""
        (tmp_path / "schemas").mkdir()
        (tmp_path / "data").mkdir()

        # Create invalid JSON schema
        (tmp_path / "schemas" / "feature.schema.json").write_text("invalid json")

        # Should handle gracefully
        with pytest.raises(json.JSONDecodeError):
            YAMLValidator(tmp_path)

    def test_validation_error_formatting(self, validator_setup):
        """Test that validation errors are formatted correctly."""
        # Create feature with nested error
        feature_data = {
            "doc_category": "feature",
            "doc_title": "Test",
            "status": "Complete",
            "extra_field": {"nested": {"deep": 123}},  # Extra field, not in schema
        }
        yaml_file = validator_setup / "data" / "test.yaml"
        yaml_file.write_text(yaml.dump(feature_data))

        validator = YAMLValidator(validator_setup)
        _is_valid, errors = validator.validate_file(yaml_file)

        # Extra field is allowed (no additionalProperties: false)
        # This test just verifies error formatting works
        assert isinstance(errors, list)


class TestIntegration:
    """Test full integration scenarios."""

    def test_full_validation_workflow(self, validator_setup):
        """Test complete validation workflow."""
        # Create mix of valid and invalid files
        (validator_setup / "data" / "features").mkdir()

        # Valid feature
        valid_data = {
            "doc_category": "feature",
            "doc_title": "Movie Module",
            "status": "Complete",
        }
        (validator_setup / "data" / "features" / "MOVIE_MODULE.yaml").write_text(
            yaml.dump(valid_data)
        )

        # Invalid feature
        invalid_data = {
            "doc_category": "feature",
            "doc_title": "Incomplete Feature",
            # Missing status
        }
        (validator_setup / "data" / "features" / "INCOMPLETE.yaml").write_text(
            yaml.dump(invalid_data)
        )

        # Valid service
        service_data = {
            "doc_category": "service",
            "doc_title": "Metadata Service",
            "package": "internal/service/metadata",
        }
        (validator_setup / "data" / "METADATA.yaml").write_text(yaml.dump(service_data))

        # Initialize and validate
        validator = YAMLValidator(validator_setup)
        results = validator.validate_directory(validator_setup / "data")

        # Should have 3 files
        assert len(results) == 3

        # Count valid/invalid
        valid_count = sum(1 for is_valid, _ in results.values() if is_valid)
        assert valid_count == 2  # Movie + Metadata
        assert len(results) - valid_count == 1  # Incomplete


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
