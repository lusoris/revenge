"""Tests for code quality automation.

Tests:
- Linter runner functionality
- Test runner functionality
- Code formatter functionality
- License checker functionality
"""

import json
from pathlib import Path
from unittest.mock import Mock, patch

import pytest

from scripts.automation.check_licenses import (
    License,
    LicenseChecker,
    LicenseCheckResult,
)
from scripts.automation.format_code import CodeFormatter, FormatResult
from scripts.automation.run_linters import LinterRunner, LintResult
from scripts.automation.run_tests import TestResult, TestRunner


class TestLinterRunner:
    """Test LinterRunner class."""

    @pytest.fixture
    def runner(self):
        """Create runner instance."""
        return LinterRunner(fix=False, verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.run_linters.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_runner_initialization(self, runner):
        """Test runner initialization."""
        assert runner.fix is False
        assert runner.verbose is False
        assert runner.root == Path.cwd()

    def test_run_golangci_lint(self, runner, mock_subprocess):
        """Test running golangci-lint."""
        result = runner.run_golangci_lint()

        assert result.linter == "golangci-lint"
        assert result.success is True
        mock_subprocess.assert_called_once()

        # Check command includes correct arguments
        args = mock_subprocess.call_args[0][0]
        assert "golangci-lint" in args
        assert "run" in args

    def test_run_golangci_lint_with_fix(self, mock_subprocess):
        """Test running golangci-lint with fix."""
        runner = LinterRunner(fix=True, verbose=False)
        runner.run_golangci_lint()

        args = mock_subprocess.call_args[0][0]
        assert "--fix" in args

    def test_run_ruff(self, runner, mock_subprocess):
        """Test running ruff."""
        result = runner.run_ruff()

        assert result.linter == "ruff"
        assert result.success is True

    def test_run_markdownlint(self, runner, mock_subprocess):
        """Test running markdownlint."""
        result = runner.run_markdownlint()

        assert result.linter == "markdownlint"
        assert result.success is True

    def test_run_prettier(self, runner, mock_subprocess):
        """Test running prettier."""
        result = runner.run_prettier()

        assert result.linter == "prettier"
        assert result.success is True

    def test_run_all_linters(self, runner, mock_subprocess):
        """Test running all linters."""
        results = runner.run_all_linters()

        assert "go" in results
        assert "python" in results
        assert "markdown" in results
        assert "prettier" in results

    def test_run_specific_linters(self, runner, mock_subprocess):
        """Test running specific linters only."""
        results = runner.run_all_linters(linters=["go", "python"])

        assert "go" in results
        assert "python" in results
        assert "markdown" not in results
        assert "prettier" not in results

    def test_linter_not_found(self):
        """Test handling of missing linter."""
        runner = LinterRunner(fix=False, verbose=False)

        with patch("scripts.automation.run_linters.subprocess.run") as mock:
            mock.side_effect = FileNotFoundError("golangci-lint not found")

            result = runner.run_golangci_lint()

        assert result.success is False
        assert result.exit_code == 127
        assert "not found" in result.error

    def test_print_summary(self, runner, capsys):
        """Test printing summary."""
        results = {
            "go": LintResult("go", True, "", "", 0),
            "python": LintResult("python", False, "", "error", 1),
        }

        runner.print_summary(results)

        captured = capsys.readouterr()
        assert "Linter Summary" in captured.out
        assert "PASS" in captured.out
        assert "FAIL" in captured.out


class TestTestRunner:
    """Test TestRunner class."""

    @pytest.fixture
    def runner(self):
        """Create runner instance."""
        return TestRunner(coverage=False, watch=False, verbose=False, threshold=80)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.run_tests.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_runner_initialization(self, runner):
        """Test runner initialization."""
        assert runner.coverage is False
        assert runner.watch is False
        assert runner.verbose is False
        assert runner.threshold == 80
        assert runner.root == Path.cwd()

    def test_run_go_tests(self, runner, mock_subprocess):
        """Test running Go tests."""
        mock_subprocess.return_value.stdout = "ok  \tpkg\t0.001s"

        result = runner.run_go_tests()

        assert result.suite == "go"
        assert result.success is True
        assert result.passed >= 0

    def test_run_go_tests_with_coverage(self, mock_subprocess):
        """Test running Go tests with coverage."""
        runner = TestRunner(coverage=True, watch=False, verbose=False, threshold=80)
        mock_subprocess.return_value.stdout = "ok  \tpkg\t0.001s\tcoverage: 85.0% of statements"

        result = runner.run_go_tests()

        assert result.coverage == 85.0

    def test_run_python_tests(self, runner, mock_subprocess):
        """Test running Python tests."""
        mock_subprocess.return_value.stdout = "272 passed in 1.38s"

        result = runner.run_python_tests()

        assert result.suite == "python"
        assert result.success is True
        assert result.passed == 272

    def test_run_python_tests_with_failures(self, runner, mock_subprocess):
        """Test running Python tests with failures."""
        mock_subprocess.return_value.stdout = "100 passed, 5 failed, 2 skipped in 1.38s"
        mock_subprocess.return_value.returncode = 1

        result = runner.run_python_tests()

        assert result.passed == 100
        assert result.failed == 5
        assert result.skipped == 2
        assert result.success is False

    def test_run_frontend_tests_no_dir(self, runner):
        """Test running frontend tests with no frontend directory."""
        with patch.object(Path, "exists", return_value=False):
            result = runner.run_frontend_tests()

        assert result.suite == "frontend"
        assert result.success is True
        assert result.passed == 0

    def test_run_all_tests(self, runner, mock_subprocess):
        """Test running all test suites."""
        mock_subprocess.return_value.stdout = "272 passed in 1.38s"

        results = runner.run_all_tests()

        assert "go" in results
        assert "python" in results
        assert "frontend" in results

    def test_run_specific_tests(self, runner, mock_subprocess):
        """Test running specific test suites only."""
        mock_subprocess.return_value.stdout = "272 passed in 1.38s"

        results = runner.run_all_tests(suites=["python"])

        assert "python" in results
        assert "go" not in results

    def test_check_coverage_threshold_pass(self):
        """Test coverage threshold check passing."""
        # Need coverage=True for threshold checking
        runner = TestRunner(coverage=True, watch=False, verbose=False, threshold=80)

        results = {
            "python": TestResult("python", True, 100, 0, 0, 85.0, "", "", 0),
        }

        result = runner.check_coverage_threshold(results)
        assert result is True

    def test_check_coverage_threshold_fail(self):
        """Test coverage threshold check failing."""
        # Need coverage=True for threshold checking
        runner = TestRunner(coverage=True, watch=False, verbose=False, threshold=80)

        results = {
            "python": TestResult("python", True, 100, 0, 0, 75.0, "", "", 0),
        }

        result = runner.check_coverage_threshold(results)
        assert result is False

    def test_print_summary(self, runner, capsys):
        """Test printing summary."""
        results = {
            "python": TestResult("python", True, 100, 0, 0, 85.0, "", "", 0),
        }

        runner.print_summary(results)

        captured = capsys.readouterr()
        assert "Test Summary" in captured.out
        assert "PASS" in captured.out


class TestCodeFormatter:
    """Test CodeFormatter class."""

    @pytest.fixture
    def formatter(self):
        """Create formatter instance."""
        return CodeFormatter(check=False, verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.format_code.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_formatter_initialization(self, formatter):
        """Test formatter initialization."""
        assert formatter.check is False
        assert formatter.verbose is False
        assert formatter.root == Path.cwd()

    def test_format_go_code(self, formatter, mock_subprocess):
        """Test formatting Go code."""
        result = formatter.format_go_code()

        assert result.formatter == "go"
        assert result.success is True

        # Should call both gofmt and goimports
        assert mock_subprocess.call_count == 2

    def test_format_go_code_check_mode(self, mock_subprocess):
        """Test formatting Go code in check mode."""
        formatter = CodeFormatter(check=True, verbose=False)
        formatter.format_go_code()

        # Check that -l is used instead of -w
        calls = mock_subprocess.call_args_list
        for call in calls:
            args = call[0][0]
            if "gofmt" in args or "goimports" in args:
                assert "-l" in args
                assert "-w" not in args

    def test_format_python_code(self, formatter, mock_subprocess):
        """Test formatting Python code."""
        result = formatter.format_python_code()

        assert result.formatter == "ruff"
        assert result.success is True

    def test_format_frontend_code(self, formatter, mock_subprocess):
        """Test formatting frontend code."""
        result = formatter.format_frontend_code()

        assert result.formatter == "prettier"
        assert result.success is True

    def test_format_all_code(self, formatter, mock_subprocess):
        """Test formatting all code."""
        results = formatter.format_all_code()

        assert "go" in results
        assert "python" in results
        assert "frontend" in results

    def test_format_specific_code(self, formatter, mock_subprocess):
        """Test formatting specific code only."""
        results = formatter.format_all_code(formatters=["python"])

        assert "python" in results
        assert "go" not in results

    def test_formatter_not_found(self):
        """Test handling of missing formatter."""
        formatter = CodeFormatter(check=False, verbose=False)

        with patch("scripts.automation.format_code.subprocess.run") as mock:
            mock.side_effect = FileNotFoundError("ruff not found")

            result = formatter.format_python_code()

        assert result.success is False
        assert result.exit_code == 127
        assert "not found" in result.error

    def test_print_summary(self, formatter, capsys):
        """Test printing summary."""
        results = {
            "go": FormatResult("go", True, 10, "", "", 0),
            "python": FormatResult("python", False, 0, "", "error", 1),
        }

        formatter.print_summary(results)

        captured = capsys.readouterr()
        assert "Formatter Summary" in captured.out
        assert "PASS" in captured.out
        assert "FAIL" in captured.out


class TestLicenseChecker:
    """Test LicenseChecker class."""

    @pytest.fixture
    def checker(self):
        """Create checker instance."""
        return LicenseChecker(strict=False, verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.check_licenses.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_checker_initialization(self, checker):
        """Test checker initialization."""
        assert checker.strict is False
        assert checker.verbose is False
        assert checker.root == Path.cwd()

    def test_normalize_license(self, checker):
        """Test license normalization."""
        assert checker.normalize_license("MIT License") == "MIT"
        assert checker.normalize_license("Apache License 2.0") == "Apache-2.0"
        assert checker.normalize_license("BSD 3-Clause") == "BSD-3-Clause"

    def test_check_license_allowed(self, checker):
        """Test checking allowed license."""
        assert checker.check_license("MIT") == "allowed"
        assert checker.check_license("Apache-2.0") == "allowed"
        assert checker.check_license("BSD-3-Clause") == "allowed"

    def test_check_license_denied(self, checker):
        """Test checking denied license."""
        assert checker.check_license("GPL-3.0") == "denied"
        assert checker.check_license("AGPL-3.0") == "denied"

    def test_check_license_unknown(self, checker):
        """Test checking unknown license."""
        assert checker.check_license("Custom-License") == "unknown"

    def test_check_go_licenses_not_found(self, checker, mock_subprocess):
        """Test checking Go licenses when tool not found."""
        mock_subprocess.return_value.returncode = 127
        mock_subprocess.return_value.stderr = "go-licenses not found"

        result = checker.check_go_licenses()

        assert result.success is True
        assert result.total == 0

    def test_check_go_licenses(self, mock_subprocess):
        """Test checking Go licenses."""
        checker = LicenseChecker(strict=False, verbose=False)

        go_output = "github.com/pkg/errors,https://example.com,MIT\n"
        go_output += "github.com/other/pkg,https://example.com,Apache-2.0\n"

        mock_subprocess.return_value.stdout = go_output
        mock_subprocess.return_value.returncode = 0

        result = checker.check_go_licenses()

        assert result.total == 2
        assert result.allowed == 2
        assert result.denied == 0

    def test_check_npm_licenses_no_frontend(self, checker):
        """Test checking npm licenses with no frontend directory."""
        with patch.object(Path, "exists", return_value=False):
            result = checker.check_npm_licenses()

        assert result.success is True
        assert result.total == 0

    def test_check_npm_licenses(self, mock_subprocess):
        """Test checking npm licenses."""
        checker = LicenseChecker(strict=False, verbose=False)

        npm_output = json.dumps({
            "package@1.0.0": {"licenses": "MIT"},
            "another@2.0.0": {"licenses": "Apache-2.0"},
        })

        mock_subprocess.return_value.stdout = npm_output
        mock_subprocess.return_value.returncode = 0

        with patch.object(Path, "exists", return_value=True):
            result = checker.check_npm_licenses()

        assert result.total == 2
        assert result.allowed == 2
        assert result.denied == 0

    def test_check_python_licenses(self, mock_subprocess):
        """Test checking Python licenses."""
        checker = LicenseChecker(strict=False, verbose=False)

        pip_output = json.dumps([
            {"Name": "pytest", "Version": "7.0.0", "License": "MIT"},
            {"Name": "ruff", "Version": "0.1.0", "License": "MIT"},
        ])

        mock_subprocess.return_value.stdout = pip_output
        mock_subprocess.return_value.returncode = 0

        result = checker.check_python_licenses()

        assert result.total == 2
        assert result.allowed == 2
        assert result.denied == 0

    def test_check_all_licenses(self, checker, mock_subprocess):
        """Test checking all licenses."""
        mock_subprocess.return_value.stdout = "[]"

        results = checker.check_all_licenses()

        assert "go" in results
        assert "npm" in results
        assert "python" in results

    def test_check_specific_licenses(self, checker, mock_subprocess):
        """Test checking specific licenses only."""
        mock_subprocess.return_value.stdout = "[]"

        results = checker.check_all_licenses(ecosystems=["python"])

        assert "python" in results
        assert "go" not in results

    def test_strict_mode_unknown_license(self):
        """Test strict mode failing on unknown license."""
        LicenseChecker(strict=True, verbose=False)

        unknown_license = License("pkg", "1.0.0", "Custom", "go")
        result = LicenseCheckResult(
            ecosystem="go",
            success=False,
            total=1,
            allowed=0,
            denied=0,
            unknown=1,
            unknown_licenses=[unknown_license],
        )

        assert not result.success

    def test_print_summary(self, checker, capsys):
        """Test printing summary."""
        results = {
            "python": LicenseCheckResult(
                ecosystem="python",
                success=True,
                total=10,
                allowed=10,
                denied=0,
                unknown=0,
            ),
        }

        checker.print_summary(results)

        captured = capsys.readouterr()
        assert "License Compliance Summary" in captured.out
        assert "PASS" in captured.out

    def test_generate_report(self, checker, tmp_path):
        """Test generating license report."""
        denied_license = License("bad-pkg", "1.0.0", "GPL-3.0", "npm")

        results = {
            "npm": LicenseCheckResult(
                ecosystem="npm",
                success=False,
                total=2,
                allowed=1,
                denied=1,
                unknown=0,
                denied_licenses=[denied_license],
            ),
        }

        output_file = tmp_path / "report.md"
        checker.generate_report(results, output_file)

        assert output_file.exists()

        content = output_file.read_text()
        assert "License Compliance Report" in content
        assert "bad-pkg" in content
        assert "GPL-3.0" in content


class TestDataClasses:
    """Test data classes."""

    def test_lint_result(self):
        """Test LintResult dataclass."""
        result = LintResult("ruff", True, "output", "", 0)

        assert result.linter == "ruff"
        assert result.success is True
        assert result.output == "output"

    def test_test_result(self):
        """Test TestResult dataclass."""
        result = TestResult("python", True, 100, 0, 0, 85.0, "", "", 0)

        assert result.suite == "python"
        assert result.passed == 100
        assert result.coverage == 85.0

    def test_format_result(self):
        """Test FormatResult dataclass."""
        result = FormatResult("prettier", True, 10, "", "", 0)

        assert result.formatter == "prettier"
        assert result.files_formatted == 10

    def test_license(self):
        """Test License dataclass."""
        lic = License("pytest", "7.0.0", "MIT", "python")

        assert lic.name == "pytest"
        assert lic.version == "7.0.0"
        assert lic.license == "MIT"
        assert lic.ecosystem == "python"

    def test_license_check_result(self):
        """Test LicenseCheckResult dataclass."""
        result = LicenseCheckResult(
            ecosystem="python",
            success=True,
            total=10,
            allowed=10,
            denied=0,
            unknown=0,
        )

        assert result.ecosystem == "python"
        assert result.total == 10
        assert result.success is True
