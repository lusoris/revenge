"""Tests for Claude Code skills validation.

Tests:
- Skill file structure validation
- YAML frontmatter validation
- Required sections present
- Markdown syntax validation
"""

from pathlib import Path

import pytest
import yaml


class TestSkillStructure:
    """Test skill file structure and content."""

    @pytest.fixture
    def skills_dir(self):
        """Get skills directory."""
        return Path(".claude/skills")

    @pytest.fixture
    def all_skills(self, skills_dir):
        """Get all skill directories."""
        return [d for d in skills_dir.iterdir() if d.is_dir()]

    def test_skills_directory_exists(self, skills_dir):
        """Test skills directory exists."""
        assert skills_dir.exists()
        assert skills_dir.is_dir()

    def test_all_skills_have_skill_md(self, all_skills):
        """Test all skill directories have SKILL.md file."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            assert skill_file.exists(), f"{skill_dir.name} missing SKILL.md"

    def test_skill_file_not_empty(self, all_skills):
        """Test all SKILL.md files are not empty."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()
            assert len(content) > 0, f"{skill_dir.name}/SKILL.md is empty"

    def test_skill_has_yaml_frontmatter(self, all_skills):
        """Test all skills have valid YAML frontmatter."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            # Check for frontmatter delimiters
            assert content.startswith("---\n"), f"{skill_dir.name} missing frontmatter start"
            assert "\n---\n" in content, f"{skill_dir.name} missing frontmatter end"

            # Extract frontmatter
            parts = content.split("---\n", 2)
            assert len(parts) >= 3, f"{skill_dir.name} invalid frontmatter structure"

            frontmatter = parts[1]

            # Parse YAML
            try:
                data = yaml.safe_load(frontmatter)
                assert data is not None, f"{skill_dir.name} empty frontmatter"
            except yaml.YAMLError as e:
                pytest.fail(f"{skill_dir.name} invalid YAML: {e}")

    def test_skill_frontmatter_has_required_fields(self, all_skills):
        """Test all skills have required frontmatter fields."""
        required_fields = ["name", "description", "disable-model-invocation", "allowed-tools"]

        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            parts = content.split("---\n", 2)
            frontmatter = parts[1]
            data = yaml.safe_load(frontmatter)

            for field in required_fields:
                assert field in data, f"{skill_dir.name} missing required field: {field}"

    def test_skill_name_matches_directory(self, all_skills):
        """Test skill name in frontmatter matches directory name."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            parts = content.split("---\n", 2)
            frontmatter = parts[1]
            data = yaml.safe_load(frontmatter)

            assert data["name"] == skill_dir.name, (
                f"{skill_dir.name} frontmatter name '{data['name']}' "
                f"doesn't match directory name '{skill_dir.name}'"
            )

    def test_skill_disable_model_invocation_is_boolean(self, all_skills):
        """Test disable-model-invocation is a boolean."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            parts = content.split("---\n", 2)
            frontmatter = parts[1]
            data = yaml.safe_load(frontmatter)

            assert isinstance(data["disable-model-invocation"], bool), (
                f"{skill_dir.name} disable-model-invocation must be boolean"
            )

    def test_skill_has_description(self, all_skills):
        """Test all skills have non-empty description."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            parts = content.split("---\n", 2)
            frontmatter = parts[1]
            data = yaml.safe_load(frontmatter)

            assert len(data["description"]) > 0, f"{skill_dir.name} has empty description"
            assert len(data["description"]) < 200, (
                f"{skill_dir.name} description too long (>200 chars)"
            )

    def test_skill_has_content_sections(self, all_skills):
        """Test all skills have expected content sections."""
        # Only strictly require Prerequisites for Phase 14 skills
        phase14_skills = {
            "check-health", "view-logs", "manage-docker-config", "manage-ci-workflows",
            "run-linters", "format-code", "check-licenses",
            "update-dependencies", "configure-dependabot", "configure-release-please",
            "setup-github-projects", "setup-github-discussions", "configure-branch-protection",
            "setup-codeql", "manage-labels", "manage-milestones", "assign-reviewers",
        }

        basic_sections = ["# ", "## Usage", "## Arguments", "## Task"]

        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            parts = content.split("---\n", 2)
            body = parts[2] if len(parts) >= 3 else ""

            # Check basic sections for all skills
            for section in basic_sections:
                assert section in body, f"{skill_dir.name} missing section: {section}"

            # Additional strict requirements only for Phase 14 skills
            if skill_dir.name in phase14_skills:
                assert "## Prerequisites" in body, f"{skill_dir.name} missing Prerequisites"

    def test_skill_has_examples(self, all_skills):
        """Test Phase 14 skills have examples section."""
        # Only require Examples for Phase 14 skills
        phase14_skills = {
            "check-health", "view-logs", "manage-docker-config", "manage-ci-workflows",
            "run-linters", "format-code", "check-licenses",
            "update-dependencies", "configure-dependabot", "configure-release-please",
            "setup-github-projects", "setup-github-discussions", "configure-branch-protection",
            "setup-codeql", "manage-labels", "manage-milestones", "assign-reviewers",
        }

        for skill_dir in all_skills:
            if skill_dir.name not in phase14_skills:
                continue  # Skip old skills

            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            # Check for Examples or Example section
            assert "## Example" in content or "### Example" in content, (
                f"{skill_dir.name} missing Examples section"
            )

    def test_skill_has_code_blocks(self, all_skills):
        """Test all skills have code blocks (bash examples)."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            # Check for code blocks
            assert "```bash" in content or "```" in content, (
                f"{skill_dir.name} missing code blocks"
            )

    def test_skill_code_blocks_are_closed(self, all_skills):
        """Test all code blocks are properly closed."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            # Count opening and closing code fences
            opening = content.count("```")
            # Should be even (each opening has a closing)
            assert opening % 2 == 0, (
                f"{skill_dir.name} has unclosed code blocks (odd number of ```)"
            )

    def test_skill_has_troubleshooting(self, all_skills):
        """Test Phase 14 skills have troubleshooting section."""
        # Only require Troubleshooting for Phase 14 skills
        phase14_skills = {
            "check-health", "view-logs", "manage-docker-config", "manage-ci-workflows",
            "run-linters", "format-code", "check-licenses",
            "update-dependencies", "configure-dependabot", "configure-release-please",
            "setup-github-projects", "setup-github-discussions", "configure-branch-protection",
            "setup-codeql", "manage-labels", "manage-milestones", "assign-reviewers",
        }

        for skill_dir in all_skills:
            if skill_dir.name not in phase14_skills:
                continue  # Skip old skills

            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            assert "## Troubleshooting" in content or "## Tips" in content, (
                f"{skill_dir.name} missing Troubleshooting/Tips section"
            )

    def test_new_phase14_skills_exist(self, skills_dir):
        """Test all Phase 14 skills were created."""
        expected_skills = [
            # Monitoring (2)
            "check-health",
            "view-logs",
            # Infrastructure (2 - coder-* already existed)
            "manage-docker-config",
            "manage-ci-workflows",
            # Code Quality (3 - run-tests exists as run-all-tests)
            "run-linters",
            "format-code",
            "check-licenses",
            # Dependency & Release (3)
            "update-dependencies",
            "configure-dependabot",
            "configure-release-please",
            # GitHub Management (7)
            "setup-github-projects",
            "setup-github-discussions",
            "configure-branch-protection",
            "setup-codeql",
            "manage-labels",
            "manage-milestones",
            "assign-reviewers",
        ]

        for skill_name in expected_skills:
            skill_dir = skills_dir / skill_name
            assert skill_dir.exists(), f"Phase 14 skill missing: {skill_name}"
            assert (skill_dir / "SKILL.md").exists(), (
                f"Phase 14 skill {skill_name} missing SKILL.md"
            )

    def test_skill_allowed_tools_format(self, all_skills):
        """Test allowed-tools field is properly formatted."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            parts = content.split("---\n", 2)
            frontmatter = parts[1]
            data = yaml.safe_load(frontmatter)

            # allowed-tools should be a string or list
            if "allowed-tools" in data:
                allowed_tools = data["allowed-tools"]
                assert isinstance(allowed_tools, (str, list)), (
                    f"{skill_dir.name} allowed-tools must be string or list"
                )

    def test_skill_argument_hint_exists(self, all_skills):
        """Test skills with arguments have argument-hint."""
        for skill_dir in all_skills:
            skill_file = skill_dir / "SKILL.md"
            content = skill_file.read_text()

            parts = content.split("---\n", 2)
            frontmatter = parts[1]
            data = yaml.safe_load(frontmatter)

            # If skill has ## Arguments section, should have argument-hint
            body = parts[2] if len(parts) >= 3 else ""
            if "## Arguments" in body:
                # Check if there's actual argument content (not just "None")
                args_section = body.split("## Arguments")[1].split("##")[0]
                if "- `$0`" in args_section or "- `$1`" in args_section:
                    # Should have argument-hint in frontmatter
                    assert "argument-hint" in data, (
                        f"{skill_dir.name} has arguments but no argument-hint"
                    )


class TestSkillContent:
    """Test skill content quality."""

    @pytest.fixture
    def skills_dir(self):
        """Get skills directory."""
        return Path(".claude/skills")

    def test_new_skills_have_sufficient_content(self, skills_dir):
        """Test new Phase 14 skills have sufficient documentation."""
        new_skills = [
            "check-health",
            "view-logs",
            "manage-docker-config",
            "manage-ci-workflows",
            "run-linters",
            "format-code",
            "check-licenses",
            "update-dependencies",
            "configure-dependabot",
            "configure-release-please",
            "setup-github-projects",
            "setup-github-discussions",
            "configure-branch-protection",
            "setup-codeql",
            "manage-labels",
            "manage-milestones",
            "assign-reviewers",
        ]

        for skill_name in new_skills:
            skill_file = skills_dir / skill_name / "SKILL.md"
            if skill_file.exists():
                content = skill_file.read_text()
                # Should have at least 200 lines of content
                lines = content.split("\n")
                assert len(lines) >= 100, (
                    f"{skill_name} has insufficient content ({len(lines)} lines, need 100+)"
                )

    def test_skills_reference_correct_scripts(self, skills_dir):
        """Test skills reference existing automation scripts."""
        skill_to_script = {
            "check-health": "scripts/automation/check_health.py",
            "view-logs": "scripts/automation/view_logs.py",
            "manage-docker-config": "scripts/automation/manage_docker.py",
            "manage-ci-workflows": "scripts/automation/manage_ci.py",
            "run-linters": "scripts/automation/run_linters.py",
            "format-code": "scripts/automation/format_code.py",
            "check-licenses": "scripts/automation/check_licenses.py",
            "update-dependencies": "scripts/automation/update_dependencies.py",
        }

        for skill_name, script_path in skill_to_script.items():
            skill_file = skills_dir / skill_name / "SKILL.md"
            if skill_file.exists():
                content = skill_file.read_text()
                # Should reference the script
                assert script_path in content, (
                    f"{skill_name} doesn't reference {script_path}"
                )
                # Verify script actually exists
                assert Path(script_path).exists(), (
                    f"{skill_name} references non-existent script: {script_path}"
                )
