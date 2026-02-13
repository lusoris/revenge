# Zed Editor Setup Guide

> Complete installation and configuration guide for Revenge development with Zed

**Last Updated**: 2026-01-31
**For**: Revenge Media Server Project

---

## Installation

### macOS

Using Homebrew:
```bash
brew install zed
```

Or download from [zed.dev](https://zed.dev/download)

### Linux

#### Ubuntu/Debian
```bash
# Add Zed repository
curl https://zed.dev/api/releases/download/0.197.2/zed-linux-x86_64.AppImage -L -o zed.AppImage
chmod +x zed.AppImage
sudo mv zed.AppImage /usr/local/bin/zed
```

#### Fedora/RHEL
```bash
sudo dnf copr enable maximbaz/AppImages
sudo dnf install zed
```

#### Arch
```bash
yay -S zed-bin
```

### Windows

Download from [zed.dev](https://zed.dev/download) or use Winget:
```powershell
winget install Zed.Zed
```

---

## Initial Configuration

### 1. Open Settings

**macOS**: `Cmd+,` (or `Zed → Settings`)
**Linux/Windows**: `Ctrl+,`

### 2. Font Configuration

Add to your user settings (`~/.config/zed/settings.json`):

```json
{
  "buffer_font_family": "JetBrains Mono",
  "buffer_font_size": 14,
  "buffer_font_weight": 400,
  "buffer_line_height": "comfortable"
}
```

**Install JetBrains Mono:**

macOS:
```bash
brew install jetbrains-mono
```

Linux:
```bash
# Download and install
mkdir -p ~/.local/share/fonts
cd ~/.local/share/fonts
wget https://github.com/JetBrains/JetBrainsMono/releases/download/v2.304/JetBrainsMono-2.304.zip
unzip JetBrainsMono-2.304.zip
fc-cache -fv
```

Windows:
Download from [JetBrains Mono GitHub](https://github.com/JetBrains/JetBrainsMono) and install fonts.

### 3. Apply Project Settings

Project settings are in `.zed/settings.json` and automatically apply when opening the Revenge project.

---

## Language Server Protocol (LSP) Setup

### Go (gopls)

**Automatic**: Zed auto-discovers gopls if Go is installed.

**Manual installation** (if needed):
```bash
go install github.com/golang/tools/gopls@latest
```

**Version requirement**: Go 1.26.0+

**Check installation**:
```bash
gopls version
```

**Expected output**:
```
golang.org/x/tools/gopls v0.16.0
    golang.org/x/sys v0.XX.X
    golang.org/x/text v0.XX.X
    golang.org/x/tools v0.XX.X (modified)
```

### Python (Ruff LSP)

**Installation**:
```bash
# Install ruff package
pip install ruff

# Or using uv (faster)
uv pip install ruff
```

**Version requirement**: Python 3.12+

**Check installation**:
```bash
ruff --version
```

**Configuration in `.zed/settings.json`**:
```json
{
  "lsp": {
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

### TypeScript

**Automatic**: Zed includes built-in TypeScript language server.

**Manual setup** (optional advanced):
```bash
npm install -g typescript typescript-language-server
```

**No configuration needed** - works out of box for `.ts`, `.tsx`, `.js`, `.jsx` files.

### Svelte

**Installation**:
```bash
npm install -g svelte-language-server
```

**Configuration in `.zed/settings.json`**:
```json
{
  "lsp": {
    "svelte": {
      "binary": {
        "path": "svelte-language-server",
        "arguments": ["--stdio"]
      }
    }
  }
}
```

**Formatter** (prettier with Svelte plugin):
```bash
npm install --save-dev prettier prettier-plugin-svelte
```

Configured in `.zed/settings.json`:
```json
{
  "languages": {
    "Svelte": {
      "formatter": {
        "external": {
          "command": "prettier",
          "arguments": [
            "--stdin-filepath", "{buffer_path}",
            "--plugin", "prettier-plugin-svelte"
          ]
        }
      }
    }
  }
}
```

---

## Formatter Setup

### Go (gofmt via gopls)

**Automatic** - Enabled in project settings.

**Behavior**: Formats on save, managed by gopls.

### Python (Ruff)

**Installation**:
```bash
pip install ruff
# or
uv pip install ruff
```

**Configuration in `.zed/settings.json`**:
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
    }
  }
}
```

**Project configuration** in `ruff.toml`:
```toml
line-length = 88
target-version = "py312"

[lint]
select = ["E", "W", "F", "I", "UP"]
ignore = ["E501"]

[format]
quote-style = "double"
indent-width = 4
```

### TypeScript/JavaScript/Svelte (Prettier)

**Installation**:
```bash
npm install --save-dev prettier prettier-plugin-svelte
```

**Configuration in `.zed/settings.json`**:
```json
{
  "languages": {
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
          "arguments": [
            "--stdin-filepath", "{buffer_path}",
            "--plugin", "prettier-plugin-svelte"
          ]
        }
      }
    }
  }
}
```

**Project configuration** in `.prettierrc.json`:
```json
{
  "printWidth": 100,
  "tabWidth": 2,
  "useTabs": false,
  "semi": true,
  "singleQuote": false,
  "trailingComma": "es5",
  "arrowParens": "always",
  "plugins": ["prettier-plugin-svelte"]
}
```

---

## Git Integration

**Enabled by default** in project settings (`.zed/settings.json`):

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

### Features

- **Git gutter**: Shows modified/added/deleted lines in margin
- **Inline blame**: Hover over lines to see last commit
- **Diff navigation**: `Cmd/Ctrl+F` in project view to find changed files

### Common Git Commands

| Action | Shortcut | Menu |
|--------|----------|------|
| Toggle Git gutter | - | View → Appearance → Git Gutter |
| Show blame | - | View → Inline Blame |
| Navigate hunks | Cmd/Ctrl+G | - |

---

## Revenge Project-Specific Setup

### Clone and Open Project

```bash
# Clone the project
git clone https://github.com/lusoris/revenge.git
cd revenge

# Open in Zed
zed .
```

### First-Time Setup

**1. Install Go dependencies**:
```bash
go mod download
```

**2. Install Frontend dependencies**:
```bash
npm install
```

**3. Install development tools**:
```bash
# gopls (Go)
go install github.com/golang/tools/gopls@latest

# Air (hot reload)
go install github.com/cosmtrek/air@latest

# sqlc (SQL code generation)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Mockery (mock generation)
go install github.com/vektra/mockery/v3@latest

# Python tools
pip install ruff ruff-lsp
```

**4. Create environment file** (if needed):
```bash
cp .env.example .env
# Edit .env with your configuration
```

**5. Start development**:
```bash
# Backend (Go with hot reload)
air

# Frontend (SvelteKit with hot reload)
npm run dev

# Both in parallel (different terminals)
```

### Project Structure in Zed

```
revenge/
├── cmd/revenge/          # Backend entry point
├── internal/             # Backend packages
│   ├── api/              # HTTP handlers
│   ├── content/          # Content modules
│   ├── service/          # Business logic
│   └── infra/            # Infrastructure
├── pkg/                  # Public packages
├── frontend/             # SvelteKit app
│   ├── src/
│   │   ├── lib/          # Components, stores
│   │   └── routes/       # Pages
│   └── static/           # Static assets
├── docs/                 # Documentation
├── scripts/              # Utility scripts
├── .zed/                 # Zed config
├── .vscode/              # VS Code config
├── go.mod               # Go dependencies
├── package.json         # Node dependencies
├── .editorconfig        # Universal editor settings
└── Makefile             # Build commands
```

### File Exclusions

Already configured in `.zed/settings.json`:
```json
{
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

This speeds up file searching and symbol navigation.

---

## Keyboard Shortcuts Quick Reference

### Navigation

| Action | macOS | Linux/Windows |
|--------|-------|---------------|
| Go to File | `Cmd+P` | `Ctrl+P` |
| Go to Symbol | `Cmd+Shift+O` | `Ctrl+Shift+O` |
| Go to Definition | `F12` | `F12` |
| Go to Declaration | `Cmd+Shift+D` | `Ctrl+Shift+D` |
| Find References | `Shift+F12` | `Shift+F12` |
| Find in Files | `Cmd+Shift+F` | `Ctrl+Shift+F` |

### Editing

| Action | macOS | Linux/Windows |
|--------|-------|---------------|
| Format Document | `Cmd+Shift+I` | `Ctrl+Shift+I` |
| Comment Line | `Cmd+/` | `Ctrl+/` |
| Multi-cursor | `Cmd+D` | `Ctrl+D` |
| Select Word | `Cmd+D` (repeat) | `Ctrl+D` (repeat) |
| Rename Symbol | `F2` | `F2` |

### Git/Version Control

| Action | macOS | Linux/Windows |
|--------|-------|---------------|
| Toggle Git Blame | - | View → Inline Blame |
| Next Hunk | `Cmd+G` | `Ctrl+G` |
| Previous Hunk | `Shift+Cmd+G` | `Shift+Ctrl+G` |

### Terminal

| Action | macOS | Linux/Windows |
|--------|-------|---------------|
| Toggle Terminal | `Ctrl+\`` | `Ctrl+\`` |
| Create New Terminal | `Ctrl+Shift+\`` | `Ctrl+Shift+\`` |
| Next Terminal | `Ctrl+Tab` | `Ctrl+Tab` |

### AI/Claude

| Action | macOS | Linux/Windows |
|--------|-------|---------------|
| Inline Assist | `Cmd+Enter` | `Ctrl+Enter` |
| Open Agent | `Cmd+Shift+A` | `Ctrl+Shift+A` |

### General

| Action | macOS | Linux/Windows |
|--------|-------|---------------|
| Command Palette | `Cmd+Shift+P` | `Ctrl+Shift+P` |
| Settings | `Cmd+,` | `Ctrl+,` |
| Theme Selector | `Cmd+K Cmd+T` | `Ctrl+K Ctrl+T` |
| Show Problem Panel | `Cmd+Shift+M` | `Ctrl+Shift+M` |
| Quick Fix | `Cmd+.` | `Ctrl+.` |

---

## Themes and Appearance

### Recommended Themes

**Dark** (default in Zed):
- One Dark Pro
- Dracula
- Nord

**Light**:
- GitHub Light
- Solarized Light

### Set Theme

**Via Palette**: `Cmd/Ctrl+K Cmd/Ctrl+T`

**Via Settings** (add to `~/.config/zed/settings.json`):
```json
{
  "theme": "One Dark Pro"
}
```

### Font Customization

In `.zed/settings.json`:
```json
{
  "buffer_font_family": "JetBrains Mono",
  "buffer_font_size": 14,
  "buffer_font_weight": 400,
  "buffer_line_height": "comfortable",
  "ui_font_size": 12
}
```

---

## Remote Development (Coder)

### Connect to Coder Workspace

```bash
# From local machine
coder ssh <workspace-name>

# Then open Zed
zed /path/to/revenge
```

Zed will use the remote machine's LSPs and formatters automatically.

### SSH Configuration

If using direct SSH (not Coder):
```bash
# Configure ~/.ssh/config
Host revenge-dev
  HostName your-dev-server.com
  User developer
  IdentityFile ~/.ssh/id_ed25519
  ForwardAgent yes

# Connect
coder ssh revenge-dev
zed .
```

---

## Troubleshooting Setup

### Go Development

**Issue**: gopls not starting
```bash
# Check Go version
go version  # Should be 1.25.6 or higher

# Check gopls
gopls version

# Reinstall gopls
go install github.com/golang/tools/gopls@latest
```

**Issue**: "gopls crashed"
1. Check Zed logs: `Cmd/Ctrl+Shift+P` → "Toggle Log"
2. Restart gopls: `Cmd/Ctrl+Shift+P` → "Restart LSP"
3. If persists, check Go workspace setup

### Python Development

**Issue**: ruff-lsp not working
```bash
# Ensure ruff is installed
pip show ruff

# Reinstall if needed
pip install --upgrade ruff
```

**Issue**: "configurationPreference" warning
- Ensure `ruff.toml` exists in project root
- Check `.pylintrc` or `pyproject.toml` don't conflict

### TypeScript/Svelte

**Issue**: TypeScript server not recognizing paths
1. Check `tsconfig.json` exists and has correct `paths`
2. Restart TypeScript server: `Cmd/Ctrl+Shift+P` → "Restart LSP"

**Issue**: Prettier not formatting Svelte files
```bash
# Verify installation
npm list prettier prettier-plugin-svelte

# Reinstall if needed
npm install --save-dev prettier prettier-plugin-svelte
```

### Git Integration

**Issue**: Git gutter not showing
1. Verify file is tracked: `git add filename`
2. Enable in settings: View → Appearance → Git Gutter
3. Ensure `.git` directory exists

---

## Performance Tuning

### Memory and CPU Usage

If Zed is slow, adjust these settings:

```json
{
  "enable_language_server": true,
  "show_completions_on_input": true,
  "hover_popover_delay": 500,
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

### Disable Unused Features

```json
{
  "telemetry": {
    "diagnostics": false,
    "metrics": false
  },
  "show_edit_predictions": false,
  "inline_assist": {
    "enabled": true
  }
}
```

---

## VS Code Comparison

If migrating from VS Code:

| Feature | VS Code | Zed |
|---------|---------|-----|
| Extensions | Rich ecosystem | Limited (built-in only) |
| Debugging | Native debug UI | Use terminal/external |
| Tasks | tasks.json | Terminal/Makefile |
| Settings Sync | Built-in | Manual |
| Remote Dev | Remote-SSH, Dev Containers | SSH only |
| Performance | Good (Electron) | Excellent (Rust) |
| Keybindings | Extensive library | Good defaults |

**Migration tips**:
1. Use Zed for editing
2. Use VS Code for debugging complex issues
3. Keybindings are mostly compatible
4. Settings need to be re-entered in JSON format

---

## Next Steps

1. **Read KEYBINDINGS.md** for comprehensive shortcut reference
2. **Read TROUBLESHOOTING.md** for common issues and solutions
3. **Install recommended extensions** (Vim/Helix mode if preferred)
4. **Customize settings** in `~/.config/zed/settings.json`
5. **Open Revenge project** and start developing

---

## References

- **Official Zed Docs**: https://zed.dev/docs
- **Zed Settings Reference**: https://zed.dev/docs/configuring-zed
- **Project Settings**: [../settings.json](../settings.json)
- **EditorConfig**: [../../.editorconfig](../../.editorconfig)
- **Tech Stack**: [../../docs/dev/design/technical/TECH_STACK.md](../../docs/dev/design/technical/TECH_STACK.md)

