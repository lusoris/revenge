"""Tests for asset generation scripts.

Test coverage:
- generate_placeholder_assets.py: Placeholder generation, branding, social images
- generate_badges.py: Badge generation from various sources

Tests use mocking to avoid network calls and file system operations where appropriate.
"""

import subprocess
from unittest.mock import Mock, patch

import pytest

from scripts.automation.generate_badges import BadgeGenerator
from scripts.automation.generate_placeholder_assets import PlaceholderGenerator


@pytest.fixture
def temp_assets_dir(tmp_path):
    """Create temporary assets directory."""
    assets_dir = tmp_path / "docs" / "assets"
    assets_dir.mkdir(parents=True)

    # Create subdirectories
    (assets_dir / "branding").mkdir()
    (assets_dir / "placeholders" / "screenshots").mkdir(parents=True)
    (assets_dir / "placeholders" / "diagrams").mkdir(parents=True)
    (assets_dir / "social").mkdir()
    (assets_dir / "badges").mkdir()

    return assets_dir


@pytest.fixture
def temp_repo(tmp_path):
    """Create temporary repository directory."""
    repo = tmp_path / "repo"
    repo.mkdir()

    # Create .git directory
    (repo / ".git").mkdir()

    # Create go.mod
    (repo / "go.mod").write_text("module github.com/owner/repo\n\ngo 1.25\n")

    # Create LICENSE
    (repo / "LICENSE").write_text("MIT License\n\nCopyright (c) 2026")

    return repo


@pytest.fixture
def mock_pil():
    """Mock PIL Image and ImageDraw."""
    with patch("scripts.automation.generate_placeholder_assets.Image") as mock_image, \
         patch("scripts.automation.generate_placeholder_assets.ImageDraw") as mock_draw, \
         patch("scripts.automation.generate_placeholder_assets.ImageFont") as mock_font:

        # Mock image instance
        mock_img = Mock()
        mock_image.new.return_value = mock_img
        mock_image.open.return_value = mock_img

        # Mock draw instance
        mock_draw_inst = Mock()
        mock_draw.Draw.return_value = mock_draw_inst
        mock_draw_inst.textbbox.return_value = (0, 0, 100, 50)

        # Mock font
        mock_font.truetype.return_value = Mock()
        mock_font.load_default.return_value = Mock()

        yield {
            "Image": mock_image,
            "ImageDraw": mock_draw,
            "ImageFont": mock_font,
            "img": mock_img,
            "draw": mock_draw_inst,
        }


class TestPlaceholderGenerator:
    """Test PlaceholderGenerator class."""

    def test_init(self, temp_assets_dir):
        """Test generator initialization."""
        generator = PlaceholderGenerator(temp_assets_dir)

        assert generator.assets_dir == temp_assets_dir
        assert generator.branding_dir == temp_assets_dir / "branding"
        assert generator.screenshots_dir == temp_assets_dir / "placeholders" / "screenshots"
        assert generator.social_dir == temp_assets_dir / "social"

    def test_hex_to_rgb(self, temp_assets_dir):
        """Test hex to RGB conversion."""
        generator = PlaceholderGenerator(temp_assets_dir)

        assert generator._hex_to_rgb("#FFFFFF") == (255, 255, 255)
        assert generator._hex_to_rgb("#000000") == (0, 0, 0)
        assert generator._hex_to_rgb("#0066CC") == (0, 102, 204)
        assert generator._hex_to_rgb("0066CC") == (0, 102, 204)  # Without #

    def test_get_font_fallback(self, temp_assets_dir):
        """Test font loading with fallback."""
        generator = PlaceholderGenerator(temp_assets_dir)

        # Should not raise even if fonts missing
        font = generator._get_font(20)
        assert font is not None

    def test_create_placeholder_image(self, temp_assets_dir, mock_pil):
        """Test placeholder image creation."""
        generator = PlaceholderGenerator(temp_assets_dir)

        path = generator.create_placeholder_image(
            "test.png",
            (1200, 800),
            "Test Feature",
            "Coming Soon"
        )

        assert path == temp_assets_dir / "placeholders" / "screenshots" / "test.png"
        mock_pil["img"].save.assert_called_once()

    def test_create_logo_placeholder_default(self, temp_assets_dir, mock_pil):
        """Test default logo creation."""
        generator = PlaceholderGenerator(temp_assets_dir)

        path = generator.create_logo_placeholder((1024, 1024), variant="default")

        assert path == temp_assets_dir / "branding" / "logo.png"
        mock_pil["img"].save.assert_called_once()

    def test_create_logo_placeholder_variants(self, temp_assets_dir, mock_pil):
        """Test all logo variants."""
        generator = PlaceholderGenerator(temp_assets_dir)

        variants = {
            "light": "logo-light.png",
            "dark": "logo-dark.png",
            "icon": "logo-icon.png",
            "wordmark": "logo-wordmark.png",
        }

        for variant, expected_name in variants.items():
            mock_pil["img"].save.reset_mock()
            path = generator.create_logo_placeholder((512, 512), variant=variant)
            assert path.name == expected_name
            mock_pil["img"].save.assert_called_once()

    def test_create_favicon(self, temp_assets_dir, mock_pil):
        """Test favicon creation."""
        generator = PlaceholderGenerator(temp_assets_dir)

        path = generator.create_favicon()

        assert path == temp_assets_dir / "branding" / "favicon.ico"
        # Should create multiple sizes and save as ICO
        assert mock_pil["img"].save.call_count >= 1

    def test_create_social_image_types(self, temp_assets_dir, mock_pil):
        """Test all social image types."""
        generator = PlaceholderGenerator(temp_assets_dir)

        social_types = {
            "og": "og-image.png",
            "home": "og-image-home.png",
            "docs": "og-image-docs.png",
            "wiki": "og-image-wiki.png",
            "twitter": "twitter-card.png",
        }

        for social_type, expected_name in social_types.items():
            mock_pil["Image"].new.reset_mock()
            path = generator.create_social_image(social_type)
            assert path.name == expected_name
            # Should create gradient
            assert mock_pil["Image"].new.called

    def test_generate_screenshots(self, temp_assets_dir, mock_pil, capsys):
        """Test screenshot generation workflow."""
        generator = PlaceholderGenerator(temp_assets_dir)

        generator.generate_screenshots()

        captured = capsys.readouterr()
        assert "Generating screenshot placeholders" in captured.out
        assert "Created:" in captured.out
        # Should create multiple screenshots
        assert mock_pil["img"].save.call_count >= 5

    def test_generate_branding(self, temp_assets_dir, mock_pil, capsys):
        """Test branding generation workflow."""
        generator = PlaceholderGenerator(temp_assets_dir)

        generator.generate_branding()

        captured = capsys.readouterr()
        assert "Generating branding placeholders" in captured.out
        assert "Created:" in captured.out
        # Should create logos + favicon + apple touch icon
        assert mock_pil["img"].save.call_count >= 5

    def test_generate_social(self, temp_assets_dir, mock_pil, capsys):
        """Test social media generation workflow."""
        generator = PlaceholderGenerator(temp_assets_dir)

        generator.generate_social()

        captured = capsys.readouterr()
        assert "Generating social media placeholders" in captured.out
        # Should create 5 social images
        assert mock_pil["Image"].new.call_count >= 5

    def test_generate_all(self, temp_assets_dir, mock_pil, capsys):
        """Test complete generation workflow."""
        generator = PlaceholderGenerator(temp_assets_dir)

        generator.generate_all()

        captured = capsys.readouterr()
        assert "Placeholder Asset Generation" in captured.out
        assert "branding placeholders" in captured.out
        assert "screenshot placeholders" in captured.out
        assert "social media placeholders" in captured.out
        assert "All placeholder assets generated" in captured.out


class TestBadgeGenerator:
    """Test BadgeGenerator class."""

    def test_init(self, temp_repo):
        """Test generator initialization."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        generator = BadgeGenerator(badges_dir, temp_repo)

        assert generator.badges_dir == badges_dir
        assert generator.repo_root == temp_repo
        assert badges_dir.exists()

    def test_run_command_success(self, temp_repo):
        """Test successful command execution."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch("subprocess.run") as mock_run:
            mock_run.return_value.stdout = "output\n"
            mock_run.return_value.returncode = 0

            result = generator._run_command(["echo", "test"])

            assert result == "output"

    def test_run_command_failure(self, temp_repo):
        """Test command execution failure."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = subprocess.CalledProcessError(1, "cmd")

            result = generator._run_command(["false"])

            assert result is None

    def test_download_badge_success(self, temp_repo):
        """Test successful badge download."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch("requests.get") as mock_get:
            mock_response = Mock()
            mock_response.content = b"<svg>badge</svg>"
            mock_get.return_value = mock_response

            path = generator._download_badge("https://example.com/badge.svg", "test.svg")

            assert path == badges_dir / "test.svg"
            assert path.exists()
            assert path.read_bytes() == b"<svg>badge</svg>"

    def test_download_badge_failure(self, temp_repo, capsys):
        """Test badge download failure."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch("scripts.automation.generate_badges.requests.get") as mock_get:
            import requests
            mock_get.side_effect = requests.RequestException("Network error")

            path = generator._download_badge("https://example.com/badge.svg", "test.svg")

            assert path is None
            captured = capsys.readouterr()
            assert "Failed to download badge" in captured.out

    def test_parse_go_coverage(self, temp_repo):
        """Test parsing Go coverage report."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        coverage_file = temp_repo / "coverage.out"
        coverage_file.write_text("mode: set\nfile.go:1.1,2.2 1 1\n")

        with patch.object(generator, "_run_command") as mock_cmd:
            mock_cmd.return_value = "total: (statement coverage) 85.4%"

            coverage = generator._parse_go_coverage(coverage_file)

            assert coverage == 85.4

    def test_parse_go_coverage_failure(self, temp_repo):
        """Test Go coverage parsing failure."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        coverage_file = temp_repo / "coverage.out"
        coverage_file.write_text("invalid")

        with patch.object(generator, "_run_command") as mock_cmd:
            mock_cmd.return_value = None

            coverage = generator._parse_go_coverage(coverage_file)

            assert coverage == 0.0

    def test_generate_coverage_badge_with_file(self, temp_repo):
        """Test coverage badge generation with coverage file."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        coverage_file = temp_repo / "coverage.out"
        coverage_file.write_text("mode: set\n")

        with patch.object(generator, "_parse_go_coverage", return_value=85.0), \
             patch.object(generator, "_download_badge", return_value=badges_dir / "coverage.svg") as mock_dl:

            path = generator.generate_coverage_badge(coverage_file)

            assert path == badges_dir / "coverage.svg"
            # Should use green color for >= 80%
            mock_dl.assert_called_once()
            call_args = mock_dl.call_args[0][0]
            assert "brightgreen" in call_args

    def test_generate_coverage_badge_no_file(self, temp_repo, capsys):
        """Test coverage badge generation without coverage file."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch.object(generator, "_download_badge", return_value=badges_dir / "coverage.svg") as mock_dl:

            generator.generate_coverage_badge()

            captured = capsys.readouterr()
            assert "Coverage file not found" in captured.out
            # Should use red color for 0%
            call_args = mock_dl.call_args[0][0]
            assert "red" in call_args

    def test_generate_build_badge_github_repo(self, temp_repo):
        """Test build badge generation for GitHub repo."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch.object(generator, "_run_command", return_value="https://github.com/owner/repo.git"), \
             patch.object(generator, "_download_badge", return_value=badges_dir / "build.svg") as mock_dl:

            path = generator.generate_build_badge()

            assert path == badges_dir / "build.svg"
            mock_dl.assert_called_once()
            call_args = mock_dl.call_args[0][0]
            assert "owner/repo" in call_args
            assert "workflow/status" in call_args

    def test_generate_build_badge_non_github(self, temp_repo, capsys):
        """Test build badge generation for non-GitHub repo."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch.object(generator, "_run_command", return_value="https://gitlab.com/owner/repo.git"), \
             patch.object(generator, "_download_badge", return_value=badges_dir / "build.svg"):

            generator.generate_build_badge()

            captured = capsys.readouterr()
            assert "Not a GitHub repository" in captured.out

    def test_generate_version_badge_from_git(self, temp_repo):
        """Test version badge generation from git tags."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch.object(generator, "_run_command", return_value="v1.2.3"), \
             patch.object(generator, "_download_badge", return_value=badges_dir / "version.svg") as mock_dl:

            path = generator.generate_version_badge()

            assert path == badges_dir / "version.svg"
            call_args = mock_dl.call_args[0][0]
            assert "version-v1.2.3" in call_args

    def test_generate_version_badge_fallback(self, temp_repo):
        """Test version badge generation with fallback."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch.object(generator, "_run_command", return_value=None), \
             patch.object(generator, "_download_badge", return_value=badges_dir / "version.svg") as mock_dl:

            generator.generate_version_badge()

            call_args = mock_dl.call_args[0][0]
            assert "0.0.0-dev" in call_args

    def test_generate_license_badge_mit(self, temp_repo):
        """Test license badge generation for MIT license."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch.object(generator, "_download_badge", return_value=badges_dir / "license.svg") as mock_dl:

            path = generator.generate_license_badge()

            assert path == badges_dir / "license.svg"
            call_args = mock_dl.call_args[0][0]
            assert "license-MIT" in call_args

    def test_generate_license_badge_apache(self, temp_repo):
        """Test license badge generation for Apache license."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        (temp_repo / "LICENSE").write_text("Apache License\nVersion 2.0")
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch.object(generator, "_download_badge", return_value=badges_dir / "license.svg") as mock_dl:

            generator.generate_license_badge()

            call_args = mock_dl.call_args[0][0]
            assert "license-Apache-2.0" in call_args

    def test_generate_all_badges(self, temp_repo, capsys):
        """Test generating all badges."""
        badges_dir = temp_repo / "docs" / "assets" / "badges"
        badges_dir.mkdir(parents=True)
        generator = BadgeGenerator(badges_dir, temp_repo)

        with patch.object(generator, "_download_badge", return_value=badges_dir / "test.svg"), \
             patch.object(generator, "_run_command", return_value=None):

            generator.generate_all()

            captured = capsys.readouterr()
            assert "Badge Generation" in captured.out
            assert "Coverage" in captured.out
            assert "Build" in captured.out
            assert "Version" in captured.out
            assert "License" in captured.out


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
