#!/usr/bin/env python3
"""YAML Data Validator - Validate YAML data files against JSON schemas.

This validator:
1. Loads YAML data files
2. Determines schema based on doc_category
3. Validates against JSON schema
4. Reports validation errors
5. Can validate individual files or entire directories

Author: Automation System
Created: 2026-01-31
"""

import json
import sys
from pathlib import Path

import yaml
from jsonschema import Draft7Validator


class YAMLValidator:
    """Validate YAML data files against JSON schemas."""

    def __init__(self, repo_root: Path):
        """Initialize validator with repository root."""
        self.repo_root = repo_root
        self.schemas_dir = repo_root / "schemas"
        self.data_dir = repo_root / "data"

        # Load schemas
        self.schemas = self._load_schemas()

    def _load_schemas(self) -> dict[str, Draft7Validator]:
        """Load all JSON schemas from schemas directory."""
        schemas = {}

        schema_files = {
            "feature": self.schemas_dir / "feature.schema.json",
            "service": self.schemas_dir / "service.schema.json",
            "integration": self.schemas_dir / "integration.schema.json",
            "generic": self.schemas_dir / "generic.schema.json",
        }

        for category, schema_path in schema_files.items():
            if not schema_path.exists():
                print(f"⚠️  Warning: Schema not found: {schema_path}")
                continue

            with open(schema_path) as f:
                schema = json.load(f)
                schemas[category] = Draft7Validator(schema)
                print(f"✓ Loaded schema: {category}")

        # Map additional categories to generic schema
        if "generic" in schemas:
            for cat in [
                "architecture",
                "operations",
                "technical",
                "pattern",
                "research",
                "other",
            ]:
                schemas[cat] = schemas["generic"]

        return schemas

    def validate_file(self, yaml_file: Path) -> tuple[bool, list[str]]:
        """Validate a single YAML file.

        Args:
            yaml_file: Path to YAML data file

        Returns:
            Tuple of (is_valid, errors)
        """
        # Load YAML data
        with open(yaml_file) as f:
            try:
                data = yaml.safe_load(f)
            except yaml.YAMLError as e:
                return False, [f"YAML parsing error: {e}"]

        # Handle empty files (yaml.safe_load returns None)
        if data is None:
            return False, ["Empty YAML file"]

        # Determine schema based on doc_category
        doc_category = data.get("doc_category")
        if not doc_category:
            return False, ["Missing 'doc_category' field"]

        if doc_category not in self.schemas:
            return False, [f"Unknown doc_category: {doc_category}"]

        # Validate against schema
        validator = self.schemas[doc_category]
        errors = []

        for error in validator.iter_errors(data):
            # Format error message
            path = " → ".join(str(p) for p in error.path) if error.path else "root"
            msg = f"  • {path}: {error.message}"
            errors.append(msg)

        is_valid = len(errors) == 0
        return is_valid, errors

    def validate_directory(self, directory: Path) -> dict[str, tuple[bool, list[str]]]:
        """Validate all YAML files in a directory recursively.

        Args:
            directory: Directory to scan

        Returns:
            Dict mapping file paths to (is_valid, errors) tuples
        """
        results = {}

        # Find all YAML files
        yaml_files = list(directory.glob("**/*.yaml")) + list(
            directory.glob("**/*.yml"),
        )

        print(f"\n📁 Scanning {directory.relative_to(self.repo_root)}")
        print(f"   Found {len(yaml_files)} YAML files\n")

        for yaml_file in sorted(yaml_files):
            # Skip shared-sot.yaml (special file)
            if yaml_file.name == "shared-sot.yaml":
                continue

            is_valid, errors = self.validate_file(yaml_file)
            results[str(yaml_file.relative_to(self.repo_root))] = (is_valid, errors)

        return results

    def print_results(self, results: dict[str, tuple[bool, list[str]]]):
        """Print validation results in a readable format."""
        valid_count = sum(1 for is_valid, _ in results.values() if is_valid)
        invalid_count = len(results) - valid_count

        print("\n" + "=" * 70)
        print("VALIDATION RESULTS")
        print("=" * 70)

        # Print invalid files first
        if invalid_count > 0:
            print(f"\n❌ INVALID FILES ({invalid_count}):\n")
            for file_path, (is_valid, errors) in results.items():
                if not is_valid:
                    print(f"  {file_path}")
                    for error in errors:
                        print(error)
                    print()

        # Print valid files summary
        if valid_count > 0:
            print(f"✅ VALID FILES ({valid_count}):")
            for file_path, (is_valid, _) in results.items():
                if is_valid:
                    print(f"  • {file_path}")

        # Summary
        print("\n" + "=" * 70)
        if invalid_count == 0:
            print(f"✅ All {valid_count} files passed validation!")
        else:
            print(f"❌ {invalid_count} file(s) failed validation")
            print(f"✅ {valid_count} file(s) passed validation")
        print("=" * 70 + "\n")

        return invalid_count == 0


def main():
    """Main entry point - Validate all YAML data files."""
    repo_root = Path(__file__).parent.parent.parent

    # Initialize validator
    validator = YAMLValidator(repo_root)

    if len(validator.schemas) == 0:
        print("❌ Error: No schemas loaded!")
        sys.exit(1)

    # Validate all data files
    results = validator.validate_directory(repo_root / "data")

    # Print results
    all_valid = validator.print_results(results)

    # Exit with appropriate code
    sys.exit(0 if all_valid else 1)


if __name__ == "__main__":
    main()
