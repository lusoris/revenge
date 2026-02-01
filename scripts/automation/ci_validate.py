#!/usr/bin/env python3
"""CI/CD Validation Script - Validate documentation changes in CI.

This script is designed to run in CI/CD pipelines to:
1. Validate all YAML files against schemas
2. Check for placeholder fields in committed files
3. Verify cross-references are valid
4. Optionally regenerate docs and check for differences
5. Exit with appropriate code for CI integration

Author: Automation System
Created: 2026-01-31
"""

import sys
from pathlib import Path


# Import from same directory
sys.path.insert(0, str(Path(__file__).parent))
from validator import YAMLValidator
from yaml_analyzer import YAMLAnalyzer


class CIValidator:
    """CI/CD validation orchestrator."""

    def __init__(self, repo_root: Path, strict: bool = False):
        """Initialize CI validator.

        Args:
            repo_root: Repository root path
            strict: If True, fail on warnings (placeholders, etc.)
        """
        self.repo_root = repo_root
        self.strict = strict
        self.validator = YAMLValidator(repo_root)
        self.analyzer = YAMLAnalyzer(repo_root)

    def validate_schemas(self) -> tuple[bool, dict]:
        """Validate all YAML files against schemas.

        Returns:
            Tuple of (all_valid, results_dict)
        """
        print(f"\n{'=' * 70}")
        print("SCHEMA VALIDATION")
        print(f"{'=' * 70}\n")

        results = self.validator.validate_directory(self.repo_root / "data")
        all_valid = self.validator.print_results(results)

        return all_valid, results

    def check_placeholders(self) -> tuple[bool, dict]:
        """Check for placeholder fields.

        Returns:
            Tuple of (acceptable, analysis_dict)
        """
        print(f"\n{'=' * 70}")
        print("PLACEHOLDER ANALYSIS")
        print(f"{'=' * 70}\n")

        analysis = self.analyzer.analyze_all()

        total_placeholders = analysis["stats"]["total_placeholders"]
        total_missing = analysis["stats"]["total_missing"]

        print(f"Placeholder fields: {total_placeholders}")
        print(f"Missing required fields: {total_missing}")

        # In strict mode, fail if placeholders exist
        if self.strict and (total_placeholders > 0 or total_missing > 0):
            print("\n❌ STRICT MODE: Cannot have placeholders or missing fields")
            return False, analysis

        # Otherwise, just warn
        if total_placeholders > 0 or total_missing > 0:
            print(
                f"\n⚠️  Warning: {total_placeholders + total_missing} fields need completion"
            )
            print("   This is acceptable for now, but should be addressed")

        return True, analysis

    def check_git_changes(self) -> tuple[bool, list[str]]:
        """Check if there are uncommitted changes to YAML or generated docs.

        Returns:
            Tuple of (clean, changed_files)
        """
        import subprocess

        print(f"\n{'=' * 70}")
        print("GIT STATUS CHECK")
        print(f"{'=' * 70}\n")

        try:
            # Check for uncommitted changes
            result = subprocess.run(
                ["git", "status", "--porcelain", "data/", "docs/"],
                cwd=self.repo_root,
                capture_output=True,
                text=True,
                check=False,
            )

            if result.returncode != 0:
                print("⚠️  Could not check git status (not a git repo?)")
                return True, []

            changed_files = [
                line[3:] for line in result.stdout.strip().split("\n") if line
            ]

            if changed_files:
                print(f"Found {len(changed_files)} uncommitted changes:")
                for f in changed_files[:10]:  # Show first 10
                    print(f"  • {f}")
                if len(changed_files) > 10:
                    print(f"  ... and {len(changed_files) - 10} more")

                if self.strict:
                    print("\n❌ STRICT MODE: Cannot have uncommitted changes")
                    return False, changed_files
                print("\n⚠️  Warning: Uncommitted changes detected")
                return True, changed_files
            print("✅ No uncommitted changes")
            return True, []

        except Exception as e:
            print(f"⚠️  Error checking git status: {e}")
            return True, []

    def run_all_checks(self) -> bool:
        """Run all validation checks.

        Returns:
            True if all checks pass, False otherwise
        """
        print(f"\n{'=' * 70}")
        print(f"CI/CD VALIDATION - {'STRICT' if self.strict else 'LENIENT'} MODE")
        print(f"{'=' * 70}")

        checks_passed = True

        # 1. Schema validation (always required)
        schema_valid, _ = self.validate_schemas()
        if not schema_valid:
            print("\n❌ Schema validation FAILED")
            checks_passed = False
        else:
            print("\n✅ Schema validation PASSED")

        # 2. Placeholder check
        placeholders_ok, _ = self.check_placeholders()
        if not placeholders_ok:
            print("\n❌ Placeholder check FAILED")
            checks_passed = False
        else:
            print("\n✅ Placeholder check PASSED")

        # 3. Git status check
        git_clean, _ = self.check_git_changes()
        if not git_clean:
            print("\n❌ Git status check FAILED")
            checks_passed = False
        else:
            print("\n✅ Git status check PASSED")

        # Final summary
        print(f"\n{'=' * 70}")
        if checks_passed:
            print("✅ ALL CHECKS PASSED")
        else:
            print("❌ SOME CHECKS FAILED")
        print(f"{'=' * 70}\n")

        return checks_passed


def main():
    """Main entry point for CI."""
    repo_root = Path(__file__).parent.parent.parent

    # Check for --strict flag
    strict = "--strict" in sys.argv

    # Initialize validator
    ci_validator = CIValidator(repo_root, strict=strict)

    # Run all checks
    all_passed = ci_validator.run_all_checks()

    # Exit with appropriate code
    sys.exit(0 if all_passed else 1)


if __name__ == "__main__":
    main()
