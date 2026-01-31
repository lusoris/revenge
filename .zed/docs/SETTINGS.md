# Zed IDE Settings Reference

> Source: https://zed.dev/docs/reference/all-settings
> Fetched: 2026-01-31
> Type: html

---

## Overview

Zed provides comprehensive configuration options for customizing the editor's appearance, behavior, and functionality. Settings can be configured globally or on a per-language basis.

---

## Configuration Files

| Location | Scope |
|----------|-------|
| `~/.config/zed/settings.json` | Global (user) settings |
| `.zed/settings.json` | Project-specific settings |

---

## Core Editor Settings

### Display & Appearance

| Setting | Type | Default | Purpose |
|---------|------|---------|---------|
| `buffer_font_family` | string | `.ZedMono` | Font for editor text |
| `buffer_font_size` | integer | 15 | Font size in pixels (6-100) |
| `buffer_font_weight` | integer | 400 | Font weight (100-900) |
| `buffer_line_height` | string | "comfortable" | Line spacing |
| `theme` | string | system | Color scheme |
| `icon_theme` | string | "Zed (Default)" | File/folder icons |
| `cursor_shape` | string | "bar" | bar, block, underline, hollow |
| `cursor_blink` | boolean | true | Cursor animation |
| `current_line_highlight` | string | "all" | none, gutter, line, all |

### Editor Behavior

| Setting | Type | Default | Purpose |
|---------|------|---------|---------|
| `auto_indent` | boolean | true | Auto-adjust indentation |
| `auto_indent_on_paste` | boolean | true | Indent pasted content |
| `hard_tabs` | boolean | false | Use tabs vs spaces |
| `tab_size` | integer | 4 | Spaces per tab |
| `soft_wrap` | string | "none" | Text wrapping |
| `preferred_line_length` | integer | 80 | Wrap column |
| `use_autoclose` | boolean | true | Auto-close brackets |
| `use_auto_surround` | boolean | true | Surround selection |

### File Handling

| Setting | Type | Default | Purpose |
|---------|------|---------|---------|
| `autosave` | string | "off" | Auto-save behavior |
| `format_on_save` | string | "on" | Format before save |
| `formatter` | string | "auto" | Formatting tool |
| `ensure_final_newline_on_save` | boolean | true | Add trailing newline |
| `remove_trailing_whitespace_on_save` | boolean | true | Strip whitespace |

---

## Code Intelligence

| Setting | Type | Default | Purpose |
|---------|------|---------|---------|
| `enable_language_server` | boolean | true | Enable LSP |
| `show_completions_on_input` | boolean | true | Auto completions |
| `show_completion_documentation` | boolean | true | Completion docs |
| `show_edit_predictions` | boolean | true | AI suggestions |
| `auto_signature_help` | boolean | false | Method signatures |
| `hover_popover_enabled` | boolean | true | Hover info |
| `hover_popover_delay` | integer | 300 | Delay in ms |

---

## Git Integration

| Setting | Type | Default |
|---------|------|---------|
| `git.git_gutter` | string | "tracked_files" |
| `git.inline_blame` | object | enabled |
| `git.hunk_style` | string | "staged_hollow" |

---

## Terminal

| Setting | Type | Default |
|---------|------|---------|
| `terminal.dock` | string | "bottom" |
| `terminal.font_size` | integer | null |
| `terminal.shell` | string | "system" |
| `terminal.working_directory` | string | "current_project_directory" |

---

## Language-Specific Configuration

```json
{
  "languages": {
    "Python": {
      "tab_size": 4,
      "format_on_save": "on",
      "formatter": {
        "external": {
          "command": "ruff",
          "arguments": ["format", "-"]
        }
      }
    },
    "Markdown": {
      "soft_wrap": "editor_width",
      "preferred_line_length": 100
    },
    "Go": {
      "tab_size": 4,
      "hard_tabs": true,
      "format_on_save": "on"
    }
  }
}
```

---

## Project Settings Example

`.zed/settings.json` for a Go/TypeScript project:

```json
{
  "tab_size": 2,
  "format_on_save": "on",
  "ensure_final_newline_on_save": true,
  "remove_trailing_whitespace_on_save": true,

  "languages": {
    "Go": {
      "tab_size": 4,
      "hard_tabs": true
    },
    "TypeScript": {
      "tab_size": 2,
      "formatter": {
        "external": {
          "command": "prettier",
          "arguments": ["--stdin-filepath", "{buffer_path}"]
        }
      }
    },
    "Python": {
      "formatter": {
        "external": {
          "command": "ruff",
          "arguments": ["format", "-"]
        }
      }
    }
  },

  "lsp": {
    "gopls": {
      "initialization_options": {
        "hints": {
          "assignVariableTypes": true,
          "compositeLiteralFields": true
        }
      }
    }
  },

  "file_scan_exclusions": [
    "**/.git",
    "**/node_modules",
    "**/dist",
    "**/__pycache__"
  ]
}
```

---

## Keyboard Shortcuts

| Task | macOS | Linux/Windows |
|------|-------|---------------|
| Command Palette | `Cmd+Shift+P` | `Ctrl+Shift+P` |
| Go to File | `Cmd+P` | `Ctrl+P` |
| Go to Symbol | `Cmd+Shift+O` | `Ctrl+Shift+O` |
| Find in Project | `Cmd+Shift+F` | `Ctrl+Shift+F` |
| Toggle Terminal | `Ctrl+\`` | `Ctrl+\`` |
| Settings | `Cmd+,` | `Ctrl+,` |
| Theme Selector | `Cmd+K Cmd+T` | `Ctrl+K Ctrl+T` |
| Agent Panel | `Cmd+Shift+A` | `Ctrl+Shift+A` |
| Inline Assist | `Cmd+Enter` | `Ctrl+Enter` |

---

## Editor Modes

```json
{
  "vim_mode": true
}
```

or

```json
{
  "helix_mode": true
}
```

---

## Privacy & Telemetry

| Setting | Type | Default |
|---------|------|---------|
| `telemetry.diagnostics` | boolean | true |
| `telemetry.metrics` | boolean | true |
| `disable_ai` | boolean | false |
