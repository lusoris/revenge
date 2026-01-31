# Settings Management Guide

**Purpose**: Comprehensive guide to understanding and managing settings across all development tools

**Last Updated**: 2026-01-31

---

## Quick Start

For most developers, you can:

1. **Copy a profile** - Use one of the pre-configured profiles in [PROFILES.md](PROFILES.md)
2. **EditorConfig handles most settings** - `.editorconfig` provides universal defaults
3. **IDE-specific tweaks** - Minor adjustments per tool as needed
4. **Environment variables** - Runtime configuration via env vars

---

## Settings Hierarchy

Settings are applied in priority order, with higher priority overriding lower:

```
1. User Settings (Highest Priority)
   └─ ~/.config/Code/User/settings.json (VS Code)
   └─ ~/.config/zed/settings.json (Zed)
   └─ Environment variables (Runtime)

2. Workspace Settings (Project-Level)
   └─ .vscode/settings.json
   └─ .zed/settings.json
   └─ .editorconfig

3. IDE Defaults
   └─ Built-in language settings

4. Hard-coded Tool Defaults (Lowest Priority)
   └─ goimports defaults
   └─ prettier defaults
   └─ ruff defaults
```

**Key Principle**: Workspace settings (in the project) should match `.editorconfig` to ensure consistency across all developers and all tools.

---

## Settings by Category

### Indentation

**EditorConfig** (applies to all files and all tools):

```ini
[*.go]
indent_style = tab
indent_size = 4

[*.{ts,tsx,js,jsx,svelte}]
indent_style = space
indent_size = 2

[*.py]
indent_style = space
indent_size = 4

[*.{json,yaml,toml}]
indent_style = space
indent_size = 2
```

**VS Code Override** (if needed, but avoid):

```json
{
  "[go]": {
    "editor.tabSize": 4,
    "editor.insertSpaces": false
  },
  "[typescript]": {
    "editor.tabSize": 2,
    "editor.insertSpaces": true
  }
}
```

**Zed Equivalent**:

```json
{
  "languages": {
    "Go": {
      "tab_size": 4,
      "hard_tabs": true
    },
    "TypeScript": {
      "tab_size": 2
    }
  }
}
```

### Line Length

**EditorConfig**:

```ini
# Default: no limit
[*.{ts,tsx,js,jsx,svelte}]
max_line_length = 100

[*.md]
max_line_length = 120

[*.py]
max_line_length = 88  # Matches Black/Ruff standard
```

**VS Code**:

```json
{
  "editor.rulers": [100, 120],
  "editor.wordWrap": "on",
  "[markdown]": {
    "editor.wordWrap": "on"
  }
}
```

**Zed**:

```json
{
  "languages": {
    "Markdown": {
      "soft_wrap": "editor_width",
      "preferred_line_length": 100
    }
  }
}
```

### End of Line (EOL) & File Endings

**EditorConfig** (universal):

```ini
[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.md]
trim_trailing_whitespace = false  # Markdown needs trailing spaces for line breaks
```

**VS Code**:

```json
{
  "files.trimTrailingWhitespace": true,
  "files.insertFinalNewline": true,
  "[markdown]": {
    "files.trimTrailingWhitespace": false
  }
}
```

**Zed**:

```json
{
  "ensure_final_newline_on_save": true,
  "remove_trailing_whitespace_on_save": true
}
```

### Format on Save

**EditorConfig**: Cannot enforce (just provides style config)

**VS Code**:

```json
{
  "editor.formatOnSave": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  },
  "[python]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "charliermarsh.ruff"
  },
  "[typescript]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[svelte]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "svelte.svelte-vscode"
  }
}
```

**Zed**:

```json
{
  "format_on_save": "on"
}
```

---

## Tool-Specific Settings

### VS Code Configuration

**File**: `.vscode/settings.json`

**Core Settings**:

```json
{
  // Editor basics
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": "explicit"
  },
  "editor.rulers": [100, 120],
  "editor.tabSize": 4,
  "editor.insertSpaces": false,  // Go uses tabs

  // File handling
  "files.trimTrailingWhitespace": true,
  "files.insertFinalNewline": true,
  "files.exclude": {
    "**/bin": true,
    "**/tmp": true,
    "**/coverage.out": true,
    "**/coverage.html": true
  },

  // Search
  "search.exclude": {
    "**/node_modules": true,
    "**/vendor": true,
    "**/bin": true,
    "**/*.sum": true
  }
}
```

**Go Extension Settings**:

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.formatTool": "goimports",
  "go.testFlags": ["-v", "-race"],
  "go.coverOnTestPackage": true,

  "gopls": {
    "ui.semanticTokens": true,
    "ui.completion.usePlaceholders": true,
    "analyses": {
      "unusedparams": true,
      "shadow": true
    },
    "staticcheck": true
  }
}
```

**Python Settings**:

```json
{
  "[python]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 4,
    "editor.insertSpaces": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    },
    "editor.defaultFormatter": "charliermarsh.ruff"
  },
  "python.analysis.typeCheckingMode": "basic",
  "python.testing.pytestEnabled": true,
  "python.testing.pytestArgs": ["tests"],

  "ruff.format.args": [],
  "ruff.lint.args": [],
  "ruff.organizeImports": true
}
```

**Frontend Settings**:

```json
{
  "[typescript]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 2,
    "editor.insertSpaces": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },

  "[svelte]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 2,
    "editor.insertSpaces": true,
    "editor.defaultFormatter": "svelte.svelte-vscode"
  }
}
```

**Terminal Settings**:

```json
{
  "terminal.integrated.defaultProfile.linux": "bash",
  "terminal.integrated.defaultProfile.windows": "PowerShell",
  "terminal.integrated.cwd": "${workspaceFolder}"
}
```

**Git Settings**:

```json
{
  "git.enableSmartCommit": true,
  "git.confirmSync": false
}
```

### Zed Configuration

**File**: `.zed/settings.json`

**Core Settings**:

```json
{
  "tab_size": 2,
  "format_on_save": "on",
  "ensure_final_newline_on_save": true,
  "remove_trailing_whitespace_on_save": true,

  "file_scan_exclusions": [
    "**/.git",
    "**/node_modules",
    "**/dist",
    "**/__pycache__",
    "**/.archive",
    "**/bin",
    "**/tmp"
  ]
}
```

**Language-Specific Settings**:

```json
{
  "languages": {
    "Go": {
      "tab_size": 4,
      "hard_tabs": true,
      "format_on_save": "on"
    },
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
    "TypeScript": {
      "tab_size": 2,
      "formatter": {
        "external": {
          "command": "prettier",
          "arguments": ["--stdin-filepath", "{buffer_path}"]
        }
      }
    },
    "Svelte": {
      "tab_size": 2,
      "formatter": {
        "external": {
          "command": "prettier",
          "arguments": ["--stdin-filepath", "{buffer_path}", "--plugin", "prettier-plugin-svelte"]
        }
      }
    }
  }
}
```

**LSP Configuration**:

```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "hints": {
          "assignVariableTypes": true,
          "compositeLiteralFields": true,
          "constantValues": true,
          "functionTypeParameters": true,
          "parameterNames": true,
          "rangeVariableTypes": true
        },
        "analyses": {
          "unusedparams": true,
          "shadow": true
        },
        "staticcheck": true
      }
    },
    "ruff": {
      "initialization_options": {
        "settings": {
          "configurationPreference": "filesystemFirst"
        }
      }
    }
  }
}
```

**Git Integration**:

```json
{
  "git": {
    "git_gutter": "tracked_files",
    "inline_blame": {
      "enabled": true
    }
  }
}
```

### EditorConfig

**File**: `.editorconfig`

Complete universal configuration:

```ini
root = true

# Universal
[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

# Go
[*.go]
indent_style = tab
indent_size = 4

[go.{mod,sum}]
indent_style = tab
indent_size = 4

# Frontend
[*.{js,jsx,ts,tsx,mjs,cjs,mts,cts}]
indent_style = space
indent_size = 2
max_line_length = 100

[*.svelte]
indent_style = space
indent_size = 2
max_line_length = 100

# Styling
[*.{css,scss,sass,less,pcss}]
indent_style = space
indent_size = 2

# Python
[*.py]
indent_style = space
indent_size = 4
max_line_length = 88

# Config files
[*.{json,jsonc,yaml,yml,toml}]
indent_style = space
indent_size = 2

[*.sh]
indent_style = space
indent_size = 2

# Markdown (special: preserve trailing spaces)
[*.{md,mdx}]
indent_style = space
indent_size = 2
trim_trailing_whitespace = false
max_line_length = 120

# Terraform
[*.{tf,tfvars,hcl}]
indent_style = space
indent_size = 2

# Docker
[Dockerfile*]
indent_style = space
indent_size = 2

[docker-compose*.{yml,yaml}]
indent_style = space
indent_size = 2

# Makefiles (must use tabs)
[Makefile]
indent_style = tab

[*.mk]
indent_style = tab

# SQL & HTML
[*.sql]
indent_style = space
indent_size = 2

[*.{html,htm,tmpl,tpl}]
indent_style = space
indent_size = 2

# GraphQL
[*.{graphql,gql}]
indent_style = space
indent_size = 2
```

### Ruff Configuration

**File**: `ruff.toml`

```toml
target-version = "py312"
line-length = 88

exclude = [
    ".git",
    ".venv",
    "venv",
    "__pycache__",
    "*.egg-info",
    "build",
    "dist",
    ".archive",
    "node_modules",
]

[lint]
select = [
    "E",      # pycodestyle errors
    "W",      # pycodestyle warnings
    "F",      # Pyflakes
    "I",      # isort
    "B",      # flake8-bugbear
    "C4",     # flake8-comprehensions
    "UP",     # pyupgrade
    "ARG",    # flake8-unused-arguments
    "SIM",    # flake8-simplify
]

ignore = [
    "E501",    # Line too long (handled by formatter)
    "E741",    # Ambiguous variable name
    "PLR0913", # Too many arguments
    "ARG001",  # Unused function argument
]

[lint.per-file-ignores]
"tests/*" = ["ARG", "PLR2004"]

[lint.isort]
known-first-party = ["revenge"]
force-single-line = false
lines-after-imports = 2

[format]
quote-style = "double"
indent-style = "space"
skip-magic-trailing-comma = false
line-ending = "auto"
docstring-code-format = true
```

---

## Settings Synchronization with Validation Scripts

### Manual Validation

Check that settings are synchronized across tools:

```bash
# Check EditorConfig is valid
editorconfig-cli . 2>&1 | grep -i error

# Verify Go formatting matches goimports
go fmt ./...
goimports -w .

# Check Python formatting
ruff format --check .

# Check TypeScript formatting
prettier --check 'src/**/*.{ts,tsx,svelte}'
```

### Continuous Validation

The project includes scripts to validate settings:

```bash
# Validate all settings files exist and are valid
scripts/validate-doc-structure.py

# Check settings consistency
python3 scripts/validate-links.py
```

### IDE vs EditorConfig Conflicts

**If you see differences**:

1. **EditorConfig always wins** for initial file creation
2. **IDE formatter** applies on save
3. **Command-line tools** (goimports, ruff) are source of truth

**Resolution order**:
```bash
# 1. Fix EditorConfig (root cause)
# 2. Update IDE settings to match
# 3. Run command-line formatter
git diff                           # See changes
go fmt ./...                       # Format Go
ruff format .                      # Format Python
prettier --write .                 # Format frontend
```

---

## Environment-Specific Settings

### Local Development

**File**: None (use workspace defaults)

**Configuration**:
- Use all formatters (slower but safer)
- Enable all linters
- Full debug output
- Watch mode enabled

**Example**:
```json
{
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": "explicit"
  },
  "go.lintFlags": ["--fast"]
}
```

### Remote Development (Coder)

**File**: `.coder/template.tf` (configures environment)

**Considerations**:
- Use `--fast` flag for gopls to reduce latency
- Disable resource-heavy extensions
- Use lighter formatters when possible

**Example Zed Settings**:
```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "analyses": {
          "unusedparams": false
        }
      }
    }
  }
}
```

### Production (CI/CD)

**Environment Variables**:
```bash
# Disable interactive formatters
export GOFMT_FLAGS="-w"
export PRETTIER_CONFIG=".prettierrc.json"
export RUFF_FORMAT_ARGS="--check"
```

**No IDE settings apply** - Use command-line tools only:
```bash
golangci-lint run ./...
ruff check .
prettier --check .
```

---

## Best Practices

### 1. Use EditorConfig as Source of Truth

EditorConfig settings should be your **single source of truth**. All other tools should respect these settings.

```bash
# Verify EditorConfig is present and complete
cat .editorconfig | grep -v '^#' | grep -v '^$'
```

### 2. Keep IDE Settings Minimal

Only override EditorConfig in IDE settings when:
- The tool is IDE-specific (e.g., VSCode Go extension settings)
- EditorConfig doesn't support the setting
- The setting is personal preference and won't affect others

**Don't do this**:
```json
{
  "editor.tabSize": 4,           // EditorConfig already handles this
  "editor.insertSpaces": false   // EditorConfig already handles this
}
```

**Do this instead**:
```json
{
  "go.lintTool": "golangci-lint",  // VSCode-specific
  "go.formatTool": "goimports"     // VSCode-specific
}
```

### 3. Formatters Should Match

All formatters should produce the same output:

| Tool | Formatter | Config |
|------|-----------|--------|
| Go | goimports | `.editorconfig` + built-in |
| Python | ruff | `ruff.toml` + `.editorconfig` |
| TypeScript/Svelte | prettier | `.prettierrc.json` (if exists) + `.editorconfig` |

### 4. Format on Save Everywhere

Enable "format on save" in all IDEs:

- **VS Code**: `editor.formatOnSave: true`
- **Zed**: `format_on_save: "on"`
- **Command line**: Use pre-commit hooks

### 5. Validate Before Committing

```bash
# Pre-commit validation script
#!/bin/bash
set -e

echo "Checking formatting..."
go fmt ./...
ruff format .
prettier --write .

echo "Running linters..."
golangci-lint run ./...
ruff check .

echo "Running tests..."
go test ./...

git diff --exit-code || echo "⚠️ Uncommitted changes after format"
```

---

## Troubleshooting Settings Issues

### Problem: Different Formatting Between IDE and CLI

**Cause**: IDE formatter has different config than CLI tool

**Solution**:
```bash
# 1. Check what formatters are being used
go version
ruff --version
prettier --version

# 2. Compare outputs
go fmt ./cmd/main.go          # See what Go formatter does
goimports -d ./cmd/main.go    # See what goimports does
ruff format --diff .          # See what ruff does

# 3. Update IDE settings to match if needed
# Prefer keeping IDE settings minimal and letting EditorConfig + CLI tools take precedence
```

### Problem: EditorConfig Not Working in VS Code

**Cause**: Missing VS Code extension

**Solution**:
```bash
# Install EditorConfig extension
code --install-extension EditorConfig.EditorConfig

# Verify it's working
cat .editorconfig  # Should see effect on file creation
```

### Problem: LSP Not Recognizing Settings

**Cause**: LSP server needs restart after settings change

**Solution**:
```bash
# VS Code: Cmd+Shift+P → "Restart Language Server"
# Zed: Cmd+Shift+P → "Restart LSP"
# Or simply reload the editor
```

### Problem: Format on Save Not Working

**Cause**: Default formatter not configured

**Solution**:

```json
{
  "[language]": {
    "editor.defaultFormatter": "extension.name",
    "editor.formatOnSave": true
  }
}
```

### Problem: Different Settings on Different Machines

**Cause**: User settings in home directory override workspace settings

**Solution**:
```bash
# Check both locations
cat ~/.config/Code/User/settings.json    # User settings (highest priority)
cat .vscode/settings.json                # Workspace settings

# Remove conflicting user settings
# Workspace settings should win
```

---

## Settings Checklist

Before committing code:

- [ ] EditorConfig present and valid (`.editorconfig`)
- [ ] VS Code settings match EditorConfig (`.vscode/settings.json`)
- [ ] Zed settings match EditorConfig (`.zed/settings.json`)
- [ ] Format on save is enabled in your IDE
- [ ] All files end with newline
- [ ] No trailing whitespace (except markdown)
- [ ] Correct indentation per file type
- [ ] Line length respected (100 for code, 120 for markdown)
- [ ] Imports organized (Go: goimports, Python: ruff, TypeScript: prettier)
- [ ] No formatter errors in terminal

---

## Related Documentation

- [PROFILES.md](PROFILES.md) - Pre-configured settings profiles
- [INTEGRATION.md](INTEGRATION.md) - How tools integrate
- [TOOL_COMPARISON.md](TOOL_COMPARISON.md) - IDE comparison
- [SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md) - Technology stack
- [DEVELOPMENT.md](../../docs/dev/design/operations/DEVELOPMENT.md) - Development workflow

---

**Maintained By**: Development Team
**Last Updated**: 2026-01-31
