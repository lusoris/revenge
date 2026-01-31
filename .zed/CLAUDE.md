# Claude Code Instructions - Zed Editor

**Tool**: Zed Editor
**Purpose**: Fast, modern alternative IDE for Revenge development
**Documentation**: [docs/INDEX.md](docs/INDEX.md)

---

## Entry Point for Claude Code

When working with Zed configuration for the Revenge project, always start by reading:

1. **Source of Truth**: [/docs/dev/design/00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md)
   - Technology stack
   - Go dependencies and versions
   - Frontend stack (Svelte, TypeScript)
   - Development tools

2. **Tech Stack**: [/docs/dev/design/technical/TECH_STACK.md](../docs/dev/design/technical/TECH_STACK.md)
   - Complete technology overview
   - Development tools list
   - Deployment platforms

3. **Zed Documentation**: [docs/INDEX.md](docs/INDEX.md)
   - Zed-specific setup
   - Settings reference

---

## Project Context

**Revenge** is a modern media server built with:
- **Backend**: Go (see SOURCE_OF_TRUTH) (greenteagc, jsonv2 experiments)
- **Frontend**: SvelteKit 2 + Svelte 5 + Tailwind CSS 4
- **Database**: PostgreSQL 18+ (ONLY)
- **Cache**: Dragonfly (Redis-compatible)
- **Search**: Typesense
- **Job Queue**: River (PostgreSQL-backed)

---

## Zed Configuration Overview

### Current Setup

**Configuration File**: `settings.json`

**Features Configured**:
- **LSP**: gopls (Go), ruff (Python)
- **Formatters**: ruff (Python), prettier (TypeScript/Svelte)
- **Format on save**: Enabled for all languages
- **Git**: Gutter with inline blame enabled
- **File scan exclusions**: .git, node_modules, dist, __pycache__, .archive, bin, tmp

**Language-Specific Settings**:
- **Go**: Tab size 4, hard tabs
- **TypeScript/Svelte/JSON/YAML**: Tab size 2, soft tabs
- **Python**: Tab size 4, soft tabs

---

## Common Tasks

### Modifying Settings

1. Edit `.zed/settings.json`
2. Consider EditorConfig precedence (if `.editorconfig` exists)
3. Test the setting works
4. Document in [docs/SETTINGS.md](docs/SETTINGS.md)

### Adding LSP Configuration

1. Check Zed's supported LSPs
2. Add to `settings.json`:
   ```json
   {
     "lsp": {
       "language-name": {
         "binary": "lsp-binary-name",
         "args": []
       }
     }
   }
   ```
3. Test LSP starts correctly
4. Document in [docs/SETTINGS.md](docs/SETTINGS.md)

### Configuring Formatter

1. Add formatter to `settings.json`:
   ```json
   {
     "languages": {
       "Language Name": {
         "format_on_save": true,
         "formatter": "external",
         "external_formatter": {
           "command": "formatter-command",
           "arguments": ["--option"]
         }
       }
     }
   }
   ```
2. Test formatting works
3. Document in [docs/SETTINGS.md](docs/SETTINGS.md)

---

## Technology-Specific Guidance

### Go Development

**LSP**: gopls

**Configuration**:
```json
{
  "lsp": {
    "gopls": {
      "binary": "gopls",
      "settings": {
        "gopls": {
          "buildFlags": ["-tags=integration"],
          "experimentalUseInvalidMetadata": true
        }
      }
    }
  }
}
```

**Formatter**: goimports (via gopls)

**Dependencies**: See [SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md#go-dependencies-core)

**Build Command**: `GOEXPERIMENT=greenteagc,jsonv2 go build ./...`

### Svelte/TypeScript Development

**LSP**: typescript-language-server (Zed default)

**Formatter**: prettier

**Configuration**:
```json
{
  "languages": {
    "TypeScript": {
      "format_on_save": true,
      "formatter": "prettier"
    },
    "Svelte": {
      "format_on_save": true,
      "formatter": "prettier"
    }
  }
}
```

### Python Development (Scripts)

**LSP**: ruff-lsp

**Formatter**: ruff

**Configuration**:
```json
{
  "lsp": {
    "ruff": {
      "binary": "ruff-lsp"
    }
  },
  "languages": {
    "Python": {
      "format_on_save": true,
      "formatter": "ruff"
    }
  }
}
```

**Python Version**: 3.12+ required

**Configuration File**: See `ruff.toml` in project root

---

## Best Practices

1. **Always reference SOURCE_OF_TRUTH** for package versions
2. **Use project settings** (`.zed/settings.json`) for project-specific configuration
3. **Follow EditorConfig** where applicable (see `.editorconfig`)
4. **Test LSP and formatters** after configuration changes
5. **Keep settings synchronized** with VS Code where possible

---

## Troubleshooting

### gopls not starting

1. Check gopls is installed: `gopls version`
2. Verify Go (check SOURCE_OF_TRUTH for version) is installed: `go version`
3. Check Zed log: View → Debug → Log File
4. Restart Zed

### Formatter not working

1. Check formatter is installed (e.g., `prettier --version`)
2. Verify `format_on_save` is enabled in settings
3. Check Zed supports the language
4. Restart Zed

### LSP errors

1. Check LSP binary path is correct
2. Verify LSP is installed and in PATH
3. Check Zed log for errors
4. Try restarting LSP: Cmd/Ctrl+Shift+P → "Restart LSP"

---

## Migration from VS Code

If migrating from VS Code, note these differences:

| Feature | VS Code | Zed |
|---------|---------|-----|
| Extensions | Marketplace | Built-in only (limited) |
| Debug | Full debug UI | Terminal-based |
| Tasks | tasks.json | Terminal commands |
| Settings Sync | Built-in | Manual |
| Remote Dev | Remote-SSH, Coder | SSH only |

**Recommendation**: Use Zed for editing, VS Code for debugging complex issues.

---

## Remote Development (Coder)

Zed supports remote development via SSH:

1. Connect to Coder workspace:
   ```bash
   coder ssh <workspace-name>
   ```

2. Open Zed:
   ```bash
   zed /path/to/revenge
   ```

3. Zed will use remote LSPs and formatters

**See Also**: [.coder/docs/ZED_INTEGRATION.md](../.coder/docs/ZED_INTEGRATION.md) (to be created)

---

## Related Documentation

- **Main Documentation**: [../docs/dev/design/INDEX.md](../docs/dev/design/INDEX.md)
- **Development Guide**: [../docs/dev/design/operations/DEVELOPMENT.md](../docs/dev/design/operations/DEVELOPMENT.md)
- **Tech Stack**: [../docs/dev/design/technical/TECH_STACK.md](../docs/dev/design/technical/TECH_STACK.md)
- **Zed Docs**: [docs/INDEX.md](docs/INDEX.md)
- **VS Code Comparison**: [.vscode/CLAUDE.md](../.vscode/CLAUDE.md)

---

## Quick Commands

```bash
# Open Zed in project
zed .

# Format current file
Cmd/Ctrl+Shift+I

# Open command palette
Cmd/Ctrl+Shift+P

# Go to definition
F12

# Find references
Shift+F12
```

---

**Last Updated**: 2026-01-31
**Maintained By**: Development Team
