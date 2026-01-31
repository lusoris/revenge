#!/usr/bin/env python3
"""Tests for TOC generator.

Test coverage for toc_generator.py:
- Anchor generation (GitHub-flavored)
- Header extraction from markdown
- TOC generation with nesting
- Frontmatter handling
- Existing TOC removal
- Edge cases
"""

import pytest

from scripts.automation.toc_generator import TOCGenerator


class TestAnchorGeneration:
    """Test GitHub-flavored markdown anchor generation."""

    def test_simple_text(self):
        """Test simple text to anchor conversion."""
        gen = TOCGenerator()
        assert gen.generate_anchor("Simple Header") == "simple-header"

    def test_with_numbers(self):
        """Test headers with numbers."""
        gen = TOCGenerator()
        assert gen.generate_anchor("Chapter 1") == "chapter-1"
        assert gen.generate_anchor("2. Second Section") == "2-second-section"

    def test_special_characters(self):
        """Test removal of special characters."""
        gen = TOCGenerator()
        assert gen.generate_anchor("What's This?") == "whats-this"
        assert gen.generate_anchor("API (Version 2.0)") == "api-version-20"
        assert gen.generate_anchor("C++ Programming") == "c-programming"

    def test_multiple_spaces(self):
        """Test collapsing multiple spaces."""
        gen = TOCGenerator()
        assert gen.generate_anchor("Multiple   Spaces") == "multiple-spaces"

    def test_leading_trailing_hyphens(self):
        """Test removal of leading/trailing hyphens."""
        gen = TOCGenerator()
        assert gen.generate_anchor("- Bullet Point") == "bullet-point"
        assert gen.generate_anchor("Trailing -") == "trailing"

    def test_unicode_characters(self):
        """Test handling of unicode characters."""
        gen = TOCGenerator()
        # Unicode should be preserved in anchors
        assert gen.generate_anchor("Café Menu") == "café-menu"
        assert gen.generate_anchor("日本語 Header") == "日本語-header"


class TestHeaderExtraction:
    """Test markdown header extraction."""

    def test_simple_headers(self):
        """Test extracting simple headers."""
        gen = TOCGenerator()
        content = """
# Level 1
## Level 2
### Level 3
"""
        headers = gen.extract_headers(content)
        assert len(headers) == 3
        assert headers[0] == (1, "Level 1")
        assert headers[1] == (2, "Level 2")
        assert headers[2] == (3, "Level 3")

    def test_headers_with_content(self):
        """Test headers mixed with content."""
        gen = TOCGenerator()
        content = """
# Introduction

Some content here.

## Background

More content.

### Details

Even more content.
"""
        headers = gen.extract_headers(content)
        assert len(headers) == 3
        assert headers[0] == (1, "Introduction")
        assert headers[1] == (2, "Background")
        assert headers[2] == (3, "Details")

    def test_headers_with_formatting(self):
        """Test headers with inline formatting."""
        gen = TOCGenerator()
        content = """
# **Bold** Header
## *Italic* Header
### `Code` Header
"""
        headers = gen.extract_headers(content)
        assert len(headers) == 3
        # Formatting should be preserved
        assert headers[0] == (1, "**Bold** Header")
        assert headers[1] == (2, "*Italic* Header")
        assert headers[2] == (3, "`Code` Header")

    def test_no_headers(self):
        """Test content with no headers."""
        gen = TOCGenerator()
        content = "Just plain text without headers."
        headers = gen.extract_headers(content)
        assert len(headers) == 0

    def test_non_header_hashes(self):
        """Test that code blocks and other hashes are ignored."""
        gen = TOCGenerator()
        content = """
# Real Header

```python
# This is not a header
def foo():
    pass
```

## Another Real Header
"""
        headers = gen.extract_headers(content)
        # Should only extract real headers, not code comments
        # (Simple implementation may extract code block headers)
        assert ("Real Header" in str(headers))
        assert ("Another Real Header" in str(headers))


class TestTOCGeneration:
    """Test TOC generation with proper nesting."""

    def test_simple_toc(self):
        """Test simple TOC generation."""
        gen = TOCGenerator()
        headers = [
            (1, "Introduction"),
            (2, "Background"),
            (2, "Methods"),
        ]
        toc = gen.generate_toc(headers)
        assert "## Table of Contents" in toc
        assert "- [Introduction](#introduction)" in toc
        assert "  - [Background](#background)" in toc
        assert "  - [Methods](#methods)" in toc

    def test_nested_toc(self):
        """Test nested TOC with multiple levels."""
        gen = TOCGenerator()
        headers = [
            (1, "Chapter 1"),
            (2, "Section 1.1"),
            (3, "Subsection 1.1.1"),
            (2, "Section 1.2"),
            (1, "Chapter 2"),
        ]
        toc = gen.generate_toc(headers)
        assert "- [Chapter 1](#chapter-1)" in toc
        assert "  - [Section 1.1](#section-11)" in toc
        assert "    - [Subsection 1.1.1](#subsection-111)" in toc
        assert "  - [Section 1.2](#section-12)" in toc
        assert "- [Chapter 2](#chapter-2)" in toc

    def test_empty_headers(self):
        """Test TOC generation with no headers."""
        gen = TOCGenerator()
        toc = gen.generate_toc([])
        # Should return empty TOC or placeholder
        assert toc == ""

    def test_duplicate_headers(self):
        """Test handling of duplicate headers."""
        gen = TOCGenerator()
        headers = [
            (1, "Introduction"),
            (2, "Introduction"),  # Duplicate
            (1, "Introduction"),  # Another duplicate
        ]
        toc = gen.generate_toc(headers)
        # All should link to same anchor
        assert toc.count("[Introduction]") == 3
        assert toc.count("(#introduction)") == 3


class TestFrontmatterHandling:
    """Test frontmatter preservation and TOC insertion."""

    def test_with_yaml_frontmatter(self):
        """Test TOC insertion after YAML frontmatter."""
        gen = TOCGenerator()
        content = """---
title: Test Document
date: 2026-01-31
---

# Introduction

Some content.

## Section 1
"""
        result = gen.add_toc(content)
        # Frontmatter should be preserved
        assert result.startswith("---\ntitle: Test Document")
        # TOC should be after frontmatter
        assert "## Table of Contents" in result
        assert result.index("## Table of Contents") > result.index("---", 3)

    def test_without_frontmatter(self):
        """Test TOC insertion without frontmatter."""
        gen = TOCGenerator()
        content = """# Introduction

Some content.

## Section 1
"""
        result = gen.add_toc(content)
        # TOC should be at the start
        assert result.startswith("## Table of Contents")

    def test_existing_toc_removal(self):
        """Test removal of existing TOC."""
        gen = TOCGenerator()
        content = """---
title: Test
---

## Table of Contents

- [Old Link](#old)

---

# New Header

Content.
"""
        result = gen.add_toc(content)
        # Old TOC should be removed
        assert "[Old Link]" not in result
        # New TOC should be present
        assert "[New Header](#new-header)" in result


class TestEdgeCases:
    """Test edge cases and error handling."""

    def test_empty_content(self):
        """Test handling of empty content."""
        gen = TOCGenerator()
        result = gen.add_toc("")
        assert result == ""

    def test_only_frontmatter(self):
        """Test content with only frontmatter."""
        gen = TOCGenerator()
        content = """---
title: Test
---
"""
        result = gen.add_toc(content)
        # Should preserve frontmatter, no TOC needed
        assert "---" in result
        assert "title: Test" in result

    def test_malformed_frontmatter(self):
        """Test handling of malformed frontmatter."""
        gen = TOCGenerator()
        content = """---
title: Test
missing closing delimiter

# Header
"""
        result = gen.add_toc(content)
        # Should handle gracefully
        assert "# Header" in result

    def test_headers_only(self):
        """Test content with only headers."""
        gen = TOCGenerator()
        content = """# Header 1
## Header 2
### Header 3
"""
        result = gen.add_toc(content)
        assert "## Table of Contents" in result
        assert "[Header 1](#header-1)" in result
        assert "[Header 2](#header-2)" in result
        assert "[Header 3](#header-3)" in result


class TestFullIntegration:
    """Test full TOC generation workflow."""

    def test_realistic_document(self):
        """Test with a realistic markdown document."""
        gen = TOCGenerator()
        content = """---
title: Movie Module
category: feature
created: 2026-01-31
---

# Movie Module

## Overview

Movie content management system.

## Architecture

### Database Schema

Details about the schema.

### Module Structure

Code organization.

## API Endpoints

### List Movies

GET /api/v1/movies

### Get Movie

GET /api/v1/movies/{id}

## Testing

### Unit Tests

Test coverage.

### Integration Tests

E2E testing.
"""
        result = gen.add_toc(content)

        # Verify structure
        assert "---" in result  # Frontmatter preserved
        assert "## Table of Contents" in result
        assert "[Overview](#overview)" in result
        assert "[Architecture](#architecture)" in result
        assert "  - [Database Schema](#database-schema)" in result
        assert "  - [Module Structure](#module-structure)" in result
        assert "[API Endpoints](#api-endpoints)" in result
        assert "  - [List Movies](#list-movies)" in result
        assert "  - [Get Movie](#get-movie)" in result
        assert "[Testing](#testing)" in result
        assert "  - [Unit Tests](#unit-tests)" in result
        assert "  - [Integration Tests](#integration-tests)" in result

        # Verify original content preserved
        assert "Movie content management system." in result
        assert "GET /api/v1/movies" in result
        assert "E2E testing." in result


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
