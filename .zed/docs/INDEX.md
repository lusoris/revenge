# Zed Documentation

> Complete reference documentation for Zed editor setup and configuration for Revenge development

**For**: Revenge Media Server Project
**Last Updated**: 2026-01-31

---

## Core Documentation

| Document | Description | Priority |
|----------|-------------|----------|
| [SETUP.md](SETUP.md) | Installation, LSP setup, formatters, and project configuration | **START HERE** |
| [KEYBINDINGS.md](KEYBINDINGS.md) | Complete keyboard shortcuts reference + VS Code comparison | Essential |
| [TROUBLESHOOTING.md](TROUBLESHOOTING.md) | Solutions for common issues and debugging tips | Reference |
| [SETTINGS.md](SETTINGS.md) | Zed settings reference and language-specific configuration | Reference |

---

## Quick Start

### For New Users

1. **Installation**: Follow [SETUP.md](SETUP.md)
   - Install Zed on your platform
   - Install required tools (Go 1.25.6, Python 3.12+, Node 20+)
   - Configure formatters and LSPs

2. **Learn Shortcuts**: Check [KEYBINDINGS.md](KEYBINDINGS.md)
   - Essential navigation shortcuts
   - Editing commands
   - Git integration

3. **Troubleshoot**: Use [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
   - If anything doesn't work
   - Common LSP/formatter issues
   - Performance optimization

### Configuration Files

- **Project Settings**: [../settings.json](../settings.json) - Revenge-specific Zed config
- **Universal Config**: See `~/.config/zed/settings.json` for user settings
- **EditorConfig**: [../../.editorconfig](../../.editorconfig) - Universal editor rules
- **Source of Truth**: [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md) - Tech stack versions

---

## Key Shortcuts (macOS vs Linux/Windows)

| Action | macOS | Linux/Windows | More Info |
|--------|-------|---------------|-----------|
| Go to File | `Cmd+P` | `Ctrl+P` | [KEYBINDINGS.md](KEYBINDINGS.md#file--navigation-most-used) |
| Command Palette | `Cmd+Shift+P` | `Ctrl+Shift+P` | [KEYBINDINGS.md](KEYBINDINGS.md#settings--ui) |
| Format Document | `Cmd+Shift+I` | `Ctrl+Shift+I` | [KEYBINDINGS.md](KEYBINDINGS.md#code-intelligence) |
| Go to Definition | `F12` | `F12` | [KEYBINDINGS.md](KEYBINDINGS.md#symbol-navigation) |
| Find & Replace | `Cmd+H` | `Ctrl+H` | [KEYBINDINGS.md](KEYBINDINGS.md#search--replace) |
| Settings | `Cmd+,` | `Ctrl+,` | [KEYBINDINGS.md](KEYBINDINGS.md#settings--ui) |
| Toggle Terminal | `Ctrl+\`` | `Ctrl+\`` | [KEYBINDINGS.md](KEYBINDINGS.md#terminal--panels) |
| Theme Selector | `Cmd+K Cmd+T` | `Ctrl+K Ctrl+T` | [KEYBINDINGS.md](KEYBINDINGS.md#settings--ui) |

---

## Language Setup

### Go Development
- **LSP**: gopls (auto-discovered)
- **Formatter**: gofmt (via gopls)
- **Setup**: [SETUP.md - Go](SETUP.md#go-gopls)
- **Troubleshooting**: [TROUBLESHOOTING.md - gopls](TROUBLESHOOTING.md#gopls-not-starting)
- **Required Version**: Go 1.25.6+

### Python Development (Scripts)
- **LSP**: ruff-lsp
- **Formatter**: ruff
- **Setup**: [SETUP.md - Python](SETUP.md#python-ruff-lsp)
- **Troubleshooting**: [TROUBLESHOOTING.md - Ruff](TROUBLESHOOTING.md#python-ruff-lsp-issues)
- **Required Version**: Python 3.12+

### TypeScript/JavaScript
- **LSP**: Built-in TypeScript server
- **Formatter**: prettier
- **Setup**: [SETUP.md - TypeScript](SETUP.md#typescript)
- **Troubleshooting**: [TROUBLESHOOTING.md - TypeScript](TROUBLESHOOTING.md#typescriptsvelte-lsp-not-working)
- **Required Version**: Node 20+

### Svelte Frontend
- **LSP**: svelte-language-server
- **Formatter**: prettier with svelte plugin
- **Setup**: [SETUP.md - Svelte](SETUP.md#svelte)
- **Troubleshooting**: [TROUBLESHOOTING.md - Svelte](TROUBLESHOOTING.md#typescriptsvelte-lsp-not-working)

---

## Project Structure

```
.zed/
├── settings.json          # Project Zed configuration
└── docs/
    ├── INDEX.md           # This file
    ├── SETUP.md           # Installation & setup
    ├── KEYBINDINGS.md     # Keyboard shortcuts
    ├── TROUBLESHOOTING.md # Common issues
    └── SETTINGS.md        # Settings reference
```

---

## Common Issues & Solutions

| Issue | Solution | Docs |
|-------|----------|------|
| gopls not starting | Check Go version, reinstall gopls | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#gopls-not-starting) |
| Prettier not formatting | Install plugin, configure args | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#prettier-not-formatting-typescriptsvelte) |
| High memory usage | Disable LSP hints, limit file scans | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#zed-using-too-much-memory) |
| Git gutter not showing | Enable in settings, verify tracking | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#git-gutter-not-showing-changes) |
| Format on save not working | Check format_on_save setting | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#format-on-save-not-working) |
| TypeScript LSP slow | Check tsconfig.json, restart LSP | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#typescriptsvelte-lsp-not-working) |

---

## External Resources

### Official Documentation
- **Zed Documentation**: https://zed.dev/docs
- **Zed Settings Reference**: https://zed.dev/docs/configuring-zed
- **Zed Keybindings**: https://zed.dev/docs/key-bindings

### Project Documentation
- **Revenge Design Index**: [../../docs/dev/design/DESIGN_INDEX.md](../../docs/dev/design/DESIGN_INDEX.md)
- **Tech Stack**: [../../docs/dev/design/technical/TECH_STACK.md](../../docs/dev/design/technical/TECH_STACK.md)
- **Development Guide**: [../../docs/dev/design/operations/DEVELOPMENT.md](../../docs/dev/design/operations/DEVELOPMENT.md)
- **Source of Truth**: [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)

### Related IDE Docs
- **VS Code Settings**: [../../.vscode/settings.json](../../.vscode/settings.json)

---

## Getting Help

1. **Check SETUP.md** first for installation issues
2. **Check TROUBLESHOOTING.md** for common problems
3. **Check Zed logs** (View → Toggle Log Panel)
4. **Search Zed issues**: https://github.com/zed-industries/zed/issues
5. **Search Revenge issues**: https://github.com/lusoris/revenge/issues

---

## Development Workflow

### Typical Day

```bash
# 1. Open project
zed .

# 2. Check keyboard shortcuts
# Cmd/Ctrl+Shift+P for commands

# 3. Edit code (formatters auto-run)
# Go: Tab size 4, hard tabs
# Python: Tab size 4, ruff format
# TypeScript/Svelte: Tab size 2, prettier

# 4. Use terminal for git/build
# Ctrl+` to open terminal in Zed

# 5. Use Cmd+Shift+P for quick actions
# Format, rename, search, etc.
```

### Performance Tips
- Open subdirectory instead of root for large repos
- Close unused tabs
- Use Cmd/Ctrl+P for file search (faster than explorer)
- Disable LSP hints if CPU usage is high

---

**Maintained by**: Revenge Development Team
**Last Updated**: 2026-01-31
