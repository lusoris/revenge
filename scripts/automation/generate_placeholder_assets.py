#!/usr/bin/env python3
"""Placeholder asset generation script.

Generates placeholder images for documentation:
- Branding assets (logos, icons, favicons)
- Screenshot placeholders for features
- Social media images (OpenGraph, Twitter cards)
- Diagram placeholders

Uses Pillow for programmatic image generation.

Requirements:
- Pillow (PIL) installed

Usage:
    python scripts/automation/generate_placeholder_assets.py --all
    python scripts/automation/generate_placeholder_assets.py --screenshots
    python scripts/automation/generate_placeholder_assets.py --branding
    python scripts/automation/generate_placeholder_assets.py --social

Author: Automation System
Created: 2026-01-31
"""

import argparse
import sys
from pathlib import Path
from typing import ClassVar


try:
    from PIL import Image, ImageDraw, ImageFont
except ImportError:
    print("❌ Error: Pillow not installed")
    print("   Install: pip install Pillow")
    sys.exit(1)


class PlaceholderGenerator:
    """Generate placeholder assets for documentation."""

    # Brand colors (temporary palette)
    COLORS: ClassVar[dict[str, str]] = {
        "primary": "#0066CC",
        "secondary": "#6B46C1",
        "accent": "#F59E0B",
        "background": "#FFFFFF",
        "surface": "#F9FAFB",
        "border": "#E5E7EB",
        "text": "#111827",
        "text_light": "#6B7280",
    }

    def __init__(self, assets_dir: Path):
        """Initialize generator.

        Args:
            assets_dir: Root assets directory (docs/assets)
        """
        self.assets_dir = assets_dir
        self.branding_dir = assets_dir / "branding"
        self.placeholders_dir = assets_dir / "placeholders"
        self.screenshots_dir = self.placeholders_dir / "screenshots"
        self.diagrams_dir = self.placeholders_dir / "diagrams"
        self.social_dir = assets_dir / "social"

        # Ensure directories exist
        for directory in [
            self.branding_dir,
            self.screenshots_dir,
            self.diagrams_dir,
            self.social_dir,
        ]:
            directory.mkdir(parents=True, exist_ok=True)

    def _hex_to_rgb(self, hex_color: str) -> tuple[int, int, int]:
        """Convert hex color to RGB tuple.

        Args:
            hex_color: Hex color string (e.g., "#0066CC")

        Returns:
            RGB tuple (r, g, b)
        """
        hex_color = hex_color.lstrip("#")
        return tuple(int(hex_color[i : i + 2], 16) for i in (0, 2, 4))

    def _get_font(self, size: int) -> ImageFont.FreeTypeFont | ImageFont.ImageFont:
        """Get font for text rendering.

        Args:
            size: Font size in pixels

        Returns:
            Font object
        """
        try:
            # Try to load a nice sans-serif font
            return ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", size)
        except OSError:
            try:
                return ImageFont.truetype("/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf", size)
            except OSError:
                # Fall back to default font
                return ImageFont.load_default()

    def create_placeholder_image(
        self,
        filename: str,
        size: tuple[int, int],
        title: str,
        subtitle: str = "Coming Soon",
        note: str = "(Implementation: Post-v1.0)",
    ) -> Path:
        """Create a placeholder image with text.

        Args:
            filename: Output filename
            size: Image size (width, height)
            title: Main title text
            subtitle: Subtitle text
            note: Additional note text

        Returns:
            Path to created image
        """
        # Create image with surface background
        img = Image.new("RGB", size, color=self._hex_to_rgb(self.COLORS["surface"]))
        draw = ImageDraw.Draw(img)

        # Draw border
        border_width = 4
        draw.rectangle(
            [
                (border_width, border_width),
                (size[0] - border_width, size[1] - border_width),
            ],
            outline=self._hex_to_rgb(self.COLORS["border"]),
            width=border_width,
        )

        # Calculate text positions (centered)
        center_x = size[0] // 2
        center_y = size[1] // 2

        # Draw title
        title_font = self._get_font(size=min(size[0] // 15, 60))
        title_bbox = draw.textbbox((0, 0), title, font=title_font)
        title_width = title_bbox[2] - title_bbox[0]
        title_height = title_bbox[3] - title_bbox[1]
        title_pos = (center_x - title_width // 2, center_y - title_height - 40)
        draw.text(
            title_pos,
            title,
            fill=self._hex_to_rgb(self.COLORS["text"]),
            font=title_font,
        )

        # Draw subtitle
        subtitle_font = self._get_font(size=min(size[0] // 25, 40))
        subtitle_bbox = draw.textbbox((0, 0), subtitle, font=subtitle_font)
        subtitle_width = subtitle_bbox[2] - subtitle_bbox[0]
        subtitle_pos = (center_x - subtitle_width // 2, center_y + 10)
        draw.text(
            subtitle_pos,
            subtitle,
            fill=self._hex_to_rgb(self.COLORS["text_light"]),
            font=subtitle_font,
        )

        # Draw note
        note_font = self._get_font(size=min(size[0] // 35, 28))
        note_bbox = draw.textbbox((0, 0), note, font=note_font)
        note_width = note_bbox[2] - note_bbox[0]
        note_pos = (center_x - note_width // 2, center_y + 60)
        draw.text(
            note_pos,
            note,
            fill=self._hex_to_rgb(self.COLORS["text_light"]),
            font=note_font,
        )

        # Save image
        output_path = self.screenshots_dir / filename
        img.save(output_path, "PNG", optimize=True)
        return output_path

    def create_logo_placeholder(self, size: tuple[int, int], variant: str = "default") -> Path:
        """Create placeholder logo.

        Args:
            size: Logo size (width, height)
            variant: Logo variant (default, light, dark, icon, wordmark)

        Returns:
            Path to created logo
        """
        # Create image
        img = Image.new("RGBA", size, color=(0, 0, 0, 0))  # Transparent background
        draw = ImageDraw.Draw(img)

        # Determine text color based on variant
        if variant == "light":
            text_color = (255, 255, 255, 255)  # White for dark backgrounds
            filename = "logo-light.png"
        elif variant == "dark":
            text_color = (*self._hex_to_rgb(self.COLORS["text"]), 255)
            filename = "logo-dark.png"
        elif variant == "icon":
            # Draw simple geometric icon (circle with "RV" initials)
            circle_radius = min(size) // 2 - 10
            circle_center = (size[0] // 2, size[1] // 2)
            draw.ellipse(
                [
                    (circle_center[0] - circle_radius, circle_center[1] - circle_radius),
                    (circle_center[0] + circle_radius, circle_center[1] + circle_radius),
                ],
                fill=(*self._hex_to_rgb(self.COLORS["primary"]), 255),
                outline=(*self._hex_to_rgb(self.COLORS["secondary"]), 255),
                width=4,
            )

            # Draw "RV" text
            font = self._get_font(size=circle_radius)
            text = "RV"
            bbox = draw.textbbox((0, 0), text, font=font)
            text_width = bbox[2] - bbox[0]
            text_height = bbox[3] - bbox[1]
            text_pos = (
                circle_center[0] - text_width // 2,
                circle_center[1] - text_height // 2,
            )
            draw.text(text_pos, text, fill=(255, 255, 255, 255), font=font)

            filename = "logo-icon.png"
        elif variant == "wordmark":
            text_color = (*self._hex_to_rgb(self.COLORS["text"]), 255)
            filename = "logo-wordmark.png"
        else:  # default
            text_color = (*self._hex_to_rgb(self.COLORS["primary"]), 255)
            filename = "logo.png"

        # For non-icon variants, draw text logo
        if variant != "icon":
            font = self._get_font(size=size[1] // 3)
            text = "REVENGE" if variant != "icon" else "RV"
            bbox = draw.textbbox((0, 0), text, font=font)
            text_width = bbox[2] - bbox[0]
            text_height = bbox[3] - bbox[1]
            text_pos = (
                (size[0] - text_width) // 2,
                (size[1] - text_height) // 2,
            )
            draw.text(text_pos, text, fill=text_color, font=font)

        # Save logo
        output_path = self.branding_dir / filename
        img.save(output_path, "PNG", optimize=True)
        return output_path

    def create_favicon(self) -> Path:
        """Create favicon (ICO format with multiple sizes).

        Returns:
            Path to created favicon
        """
        # Create icon in multiple sizes
        sizes = [16, 32, 48, 64]
        images = []

        for size in sizes:
            img = Image.new("RGBA", (size, size), color=(0, 0, 0, 0))
            draw = ImageDraw.Draw(img)

            # Draw circle background
            draw.ellipse(
                [(2, 2), (size - 2, size - 2)],
                fill=(*self._hex_to_rgb(self.COLORS["primary"]), 255),
            )

            # Draw "RV" text
            font_size = max(8, size // 3)
            font = self._get_font(size=font_size)
            text = "RV"
            bbox = draw.textbbox((0, 0), text, font=font)
            text_width = bbox[2] - bbox[0]
            text_height = bbox[3] - bbox[1]
            text_pos = (
                (size - text_width) // 2,
                (size - text_height) // 2,
            )
            draw.text(text_pos, text, fill=(255, 255, 255, 255), font=font)

            images.append(img)

        # Save as ICO
        output_path = self.branding_dir / "favicon.ico"
        images[0].save(
            output_path,
            format="ICO",
            sizes=[(img.width, img.height) for img in images],
            append_images=images[1:],
        )
        return output_path

    def create_social_image(self, image_type: str = "og") -> Path:
        """Create social media image (OpenGraph or Twitter card).

        Args:
            image_type: Type of image (og, twitter, home, docs, wiki)

        Returns:
            Path to created image
        """
        # Determine size based on type
        if image_type == "twitter":
            size = (1200, 600)
            filename = "twitter-card.png"
        else:
            size = (1200, 630)
            if image_type == "home":
                filename = "og-image-home.png"
            elif image_type == "docs":
                filename = "og-image-docs.png"
            elif image_type == "wiki":
                filename = "og-image-wiki.png"
            else:
                filename = "og-image.png"

        # Create image with gradient background
        img = Image.new("RGB", size, color=self._hex_to_rgb(self.COLORS["primary"]))
        draw = ImageDraw.Draw(img)

        # Draw gradient effect (simple two-tone)
        for y in range(size[1]):
            alpha = y / size[1]
            r1, g1, b1 = self._hex_to_rgb(self.COLORS["primary"])
            r2, g2, b2 = self._hex_to_rgb(self.COLORS["secondary"])
            r = int(r1 * (1 - alpha) + r2 * alpha)
            g = int(g1 * (1 - alpha) + g2 * alpha)
            b = int(b1 * (1 - alpha) + b2 * alpha)
            draw.line([(0, y), (size[0], y)], fill=(r, g, b))

        # Draw title
        title_font = self._get_font(size=80)
        title = "REVENGE"
        title_bbox = draw.textbbox((0, 0), title, font=title_font)
        title_width = title_bbox[2] - title_bbox[0]
        title_pos = ((size[0] - title_width) // 2, size[1] // 2 - 80)
        draw.text(title_pos, title, fill=(255, 255, 255), font=title_font)

        # Draw subtitle based on type
        subtitle_font = self._get_font(size=40)
        subtitles = {
            "home": "Modern Self-Hosted Media Server",
            "docs": "Documentation",
            "wiki": "Wiki",
            "og": "Modern Self-Hosted Media Server",
            "twitter": "Modern Self-Hosted Media Server",
        }
        subtitle = subtitles.get(image_type, "Modern Self-Hosted Media Server")
        subtitle_bbox = draw.textbbox((0, 0), subtitle, font=subtitle_font)
        subtitle_width = subtitle_bbox[2] - subtitle_bbox[0]
        subtitle_pos = ((size[0] - subtitle_width) // 2, size[1] // 2 + 20)
        draw.text(subtitle_pos, subtitle, fill=(255, 255, 255, 200), font=subtitle_font)

        # Save image
        output_path = self.social_dir / filename
        img.save(output_path, "PNG", optimize=True)
        return output_path

    def generate_screenshots(self):
        """Generate placeholder screenshots for common features."""
        print("\n📸 Generating screenshot placeholders...")

        screenshots = [
            ("movie-library.png", (1200, 800), "Movie Library", "Browse your movie collection"),
            ("movie-details.png", (1200, 800), "Movie Details", "View detailed movie information"),
            ("tv-library.png", (1200, 800), "TV Shows Library", "Browse your TV shows"),
            ("tv-details.png", (1200, 800), "TV Show Details", "View show and episode information"),
            ("player-controls.png", (1600, 900), "Video Player", "Watch with advanced playback controls"),
            ("settings-general.png", (1200, 800), "Settings", "Configure your server"),
            ("dashboard.png", (1600, 900), "Dashboard", "Monitor your media server"),
            ("search-results.png", (1200, 800), "Search", "Find content across all libraries"),
        ]

        for filename, size, title, subtitle in screenshots:
            path = self.create_placeholder_image(filename, size, title, subtitle)
            print(f"   ✓ Created: {path.name}")

    def generate_branding(self):
        """Generate placeholder branding assets."""
        print("\n🎨 Generating branding placeholders...")

        # Logo variants
        variants = [
            ("default", (1024, 1024)),
            ("light", (1024, 1024)),
            ("dark", (1024, 1024)),
            ("icon", (512, 512)),
            ("wordmark", (800, 200)),
        ]

        for variant, size in variants:
            path = self.create_logo_placeholder(size, variant)
            print(f"   ✓ Created: {path.name}")

        # Favicon
        favicon_path = self.create_favicon()
        print(f"   ✓ Created: {favicon_path.name}")

        # Apple touch icon (simplified - just a PNG version of icon)
        icon_img = Image.open(self.branding_dir / "logo-icon.png")
        icon_img = icon_img.resize((180, 180), Image.Resampling.LANCZOS)
        apple_touch_path = self.branding_dir / "apple-touch-icon.png"
        icon_img.save(apple_touch_path, "PNG", optimize=True)
        print(f"   ✓ Created: {apple_touch_path.name}")

    def generate_social(self):
        """Generate placeholder social media images."""
        print("\n🌐 Generating social media placeholders...")

        social_types = ["og", "home", "docs", "wiki", "twitter"]

        for social_type in social_types:
            path = self.create_social_image(social_type)
            print(f"   ✓ Created: {path.name}")

    def generate_all(self):
        """Generate all placeholder assets."""
        print(f"\n{'='*70}")
        print("Placeholder Asset Generation")
        print(f"{'='*70}\n")
        print(f"Assets directory: {self.assets_dir}")

        self.generate_branding()
        self.generate_screenshots()
        self.generate_social()

        print(f"\n{'='*70}")
        print("✅ All placeholder assets generated!")
        print(f"{'='*70}\n")


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Generate placeholder documentation assets")
    parser.add_argument(
        "--all",
        action="store_true",
        help="Generate all assets",
    )
    parser.add_argument(
        "--screenshots",
        action="store_true",
        help="Generate screenshot placeholders only",
    )
    parser.add_argument(
        "--branding",
        action="store_true",
        help="Generate branding assets only",
    )
    parser.add_argument(
        "--social",
        action="store_true",
        help="Generate social media images only",
    )
    parser.add_argument(
        "--assets-dir",
        type=Path,
        default=Path.cwd() / "docs" / "assets",
        help="Assets directory (default: docs/assets)",
    )

    args = parser.parse_args()

    # Verify assets directory exists
    if not args.assets_dir.exists():
        print(f"❌ Error: Assets directory not found: {args.assets_dir}")
        print("   Create it first or specify --assets-dir")
        sys.exit(1)

    generator = PlaceholderGenerator(args.assets_dir)

    # Generate based on arguments
    if args.all or (not args.screenshots and not args.branding and not args.social):
        generator.generate_all()
    else:
        if args.branding:
            generator.generate_branding()
        if args.screenshots:
            generator.generate_screenshots()
        if args.social:
            generator.generate_social()


if __name__ == "__main__":
    main()
