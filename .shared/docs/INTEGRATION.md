# Tool Integration Guide

**Purpose**: Understanding how all development tools work together in the Revenge project

**Last Updated**: 2026-01-31

---

## Table of Contents

- [Tool Ecosystem Overview](#tool-ecosystem-overview)
- [Integration Architecture](#integration-architecture)
- [Data Flow](#data-flow)
- [Configuration Synchronization](#configuration-synchronization)
- [Recommended Tool Combinations](#recommended-tool-combinations)
- [Settings Precedence](#settings-precedence)

---

## Tool Ecosystem Overview

The Revenge project uses 8 primary development tools that work together:

```
┌─────────────────────────────────────────────────────────────┐
│                    REVENGE PROJECT                          │
│                   (Media Server)                            │
└─────────────────────────────────────────────────────────────┘
                            │
            ┌───────────────┼───────────────┐
            │               │               │
      ┌─────▼─────┐   ┌────▼────┐    ┌────▼────┐
      │   LOCAL   │   │ REMOTE  │    │   CI/CD │
      │    DEV    │   │   DEV   │    │         │
      └───────────┘   └─────────┘    └─────────┘
            │               │               │
    ┌───────┼───────┐       │       ┌───────┼───────┐
    │       │       │       │       │       │       │
┌───▼──┐ ┌─▼──┐ ┌─▼───┐ ┌─▼────┐ ┌▼──────┐ │   ┌──▼────┐
│VSCode│ │Zed │ │Claude│ │Coder │ │GitHub │ │   │Docker │
│      │ │    │ │Code  │ │      │ │Actions│ │   │       │
└──┬───┘ └┬───┘ └──┬───┘ └──┬───┘ └───┬───┘ │   └───┬───┘
   │      │       │        │         │      │       │
   └──────┴───────┴────────┴─────────┴──────┴───────┘
                            │
                ┌───────────┼───────────┐
                │           │           │
          ┌─────▼─────┐ ┌──▼───┐ ┌────▼────┐
          │ Git Hooks │ │ Deps │ │ Secrets │
          │(pre-commit│ │(Dbot)│ │ (GitHub)│
          └───────────┘ └──────┘ └─────────┘
```

### Tool Categories

#### 1. IDEs (Code Editing)
- **VS Code** - Comprehensive IDE with extensive extensions
- **Zed** - Fast, modern alternative IDE
- **JetBrains Gateway** - For GoLand/IntelliJ users (via Coder)

#### 2. AI Assistant
- **Claude Code** - AI-powered coding assistant with custom skills

#### 3. Remote Development
- **Coder** - Remote workspace management (https://coder.ancilla.lol)

#### 4. Version Control
- **Git Hooks** - Pre-commit validation and linting
- **GitHub** - Repository hosting and collaboration

#### 5. CI/CD
- **GitHub Actions** - Automated testing, building, deployment
- **Dependabot** - Automated dependency updates

#### 6. Containerization
- **Docker** - Local testing and production deployment
- **Docker Compose** - Multi-service orchestration

---

## Integration Architecture

### Layer 1: Code Editing (IDEs)

**VS Code** ←→ **Zed** (Choose one or use both)

Both IDEs:
- Share the same `.editorconfig` settings
- Use the same formatters (goimports, ruff, prettier)
- Connect to the same LSPs (gopls, ruff-lsp, typescript-language-server)
- Can connect to Coder for remote development

**Key Files**:
- `.vscode/settings.json` - VS Code configuration
- `.zed/settings.json` - Zed configuration
- `.editorconfig` - Universal editor settings

### Layer 2: AI Assistance

**Claude Code** integrates with:
- **Git** - Via configured permissions (status, diff, commit, push)
- **Build tools** - Can run go build, tests, linting
- **Documentation** - Generates docs, updates indexes
- **Coder** - Via custom skills (workspace management)

**Key Files**:
- `.claude/settings.local.json` - Permissions and hooks
- `.claude/skills/` - Custom automation skills
- `.claude/docs/` - Claude Code documentation

### Layer 3: Remote Development

**Coder** provides:
- Remote workspaces accessible via SSH or web
- Support for all IDEs (VS Code, Zed, JetBrains)
- Persistent storage for workspace data
- Pre-configured development environment

**Key Files**:
- `.coder/template.tf` - Terraform template for workspaces
- `.coder/docs/` - Coder-specific documentation

**Integration with**:
- **VS Code** - Via Remote-SSH or code-server (browser)
- **Zed** - Via SSH
- **JetBrains Gateway** - Direct integration
- **Git** - Credentials configured in workspace
- **Docker** - Available inside workspace

### Layer 4: Quality Assurance

**Git Hooks** run before:
- **Commit** - Linting, formatting, tests
- **Push** - Full test suite

**GitHub Actions** run on:
- **Push to main/develop** - CI pipeline
- **Pull Request** - PR checks
- **Schedule** - Weekly source refresh, dependency updates
- **Release tag** - Release creation

**Key Files**:
- `.githooks/` - Hook scripts
- `.pre-commit-config.yaml` - pre-commit framework config
- `.github/workflows/` - GitHub Actions workflows
- `.github/dependabot.yml` - Dependency update config

---

## Data Flow

### Local Development Flow

```
1. Developer → IDE (VS Code/Zed)
2. IDE → LSP (gopls/ruff-lsp) → Code intelligence
3. IDE → Formatter (goimports/ruff/prettier) → Code formatting
4. Developer → Save file
5. File watcher (Air) → Auto-rebuild → Dev server restart
6. Developer → Git commit
7. Git → Pre-commit hooks → Lint/Format/Test
8. If pass → Commit created
9. Developer → Git push
10. Git → Pre-push hook → Run tests
11. If pass → Push to remote
12. GitHub → Trigger Actions → CI/CD
```

### Remote Development Flow (Coder)

```
1. Developer → Coder CLI → Create/Start workspace
2. Coder → Provision workspace (Docker/K8s/Swarm)
3. Developer → Connect IDE (VS Code/Zed/JetBrains)
4. IDE → SSH tunnel → Coder workspace
5. [Same as local flow inside workspace]
6. Workspace → Git push → GitHub
7. GitHub → Actions → CI/CD
```

### CI/CD Flow

```
1. Push to GitHub → Trigger workflow
2. GitHub Actions → Checkout code
3. Actions → Setup Go/Node/Python
4. Actions → Cache dependencies
5. Actions → Run linters (golangci-lint, ruff)
6. Actions → Run tests
7. Actions → Build binaries
8. Actions → Upload artifacts
9. If main branch → Deploy / Create release
```

---

## Configuration Synchronization

### Settings Precedence Hierarchy

```
┌─────────────────────────────────────┐
│  1. .editorconfig (Highest)         │  ← Universal settings
├─────────────────────────────────────┤
│  2. IDE Workspace Settings          │  ← Project-specific
│     .vscode/settings.json            │
│     .zed/settings.json               │
├─────────────────────────────────────┤
│  3. IDE User Settings               │  ← Personal preferences
│     ~/.config/Code/User/settings.json│
│     ~/.config/zed/settings.json      │
├─────────────────────────────────────┤
│  4. Tool-specific configs           │  ← Tool defaults
│     ruff.toml                        │
│     .golangci.yml                    │
│     .prettierrc                      │
└─────────────────────────────────────┘
```

### Settings Synchronization Matrix

| Setting | .editorconfig | .vscode | .zed | ruff.toml | Notes |
|---------|---------------|---------|------|-----------|-------|
| **Indentation (Go)** | Tab, 4 | Tab, 4 | Tab, 4 | N/A | Hard tabs |
| **Indentation (TS/JS)** | Space, 2 | Space, 2 | Space, 2 | N/A | Soft tabs |
| **Indentation (Python)** | Space, 4 | Space, 4 | Space, 4 | 4 | Soft tabs |
| **Line length (Python)** | N/A | 88 | 88 | 88 | ruff default |
| **Trim trailing whitespace** | Yes | Yes | Yes | N/A | All files |
| **Insert final newline** | Yes | Yes | Yes | N/A | All files |
| **Format on save** | N/A | Yes | Yes | N/A | All languages |

### Formatter Synchronization

**Go**:
- Tool: `goimports` (via gopls)
- Config: `.vscode/settings.json`, `.zed/settings.json`
- Format on save: Enabled
- Tabs: Hard tabs, size 4

**Python**:
- Tool: `ruff`
- Config: `ruff.toml`
- Format on save: Enabled
- Tabs: Soft tabs (spaces), size 4
- Line length: 88 characters

**TypeScript/Svelte**:
- Tool: `prettier`
- Config: `.prettierrc` (if exists) or defaults
- Format on save: Enabled
- Tabs: Soft tabs (spaces), size 2

### LSP Synchronization

All IDEs use the same LSPs:

**gopls** (Go):
```json
{
  "gopls": {
    "buildFlags": ["-tags=integration"],
    "analyses": {
      "unusedparams": true
    }
  }
}
```

**ruff-lsp** (Python):
- Configured via `ruff.toml`
- Same rules enforced in all IDEs

**typescript-language-server** (TypeScript/Svelte):
- Uses project's `tsconfig.json`
- Svelte plugin enabled

---

## Recommended Tool Combinations

### Scenario 1: Local Backend Development (Go)

**Recommended**:
- **IDE**: VS Code (mature Go support) OR Zed (fast performance)
- **AI**: Claude Code with go-specific skills
- **Testing**: Local Docker Compose for PostgreSQL/Dragonfly
- **Git**: Local Git hooks for pre-commit checks

**Why**:
- Full LSP support (gopls)
- Fast hot reload with Air
- Excellent debugging support
- Local database for testing

---

### Scenario 2: Remote Backend Development (Go)

**Recommended**:
- **Platform**: Coder workspace
- **IDE**: Zed (via SSH) for speed, or VS Code (desktop) for debugging
- **AI**: Claude Code (runs locally, connects to remote repo)
- **Testing**: Docker inside Coder workspace
- **Git**: Remote Git hooks

**Why**:
- Consistent environment
- No local resource usage
- Access from any device
- Shared team configuration

---

### Scenario 3: Frontend Development (Svelte/TypeScript)

**Recommended**:
- **IDE**: VS Code (best Svelte support)
- **AI**: Claude Code
- **Dev server**: Local (fast HMR)
- **Git**: Local hooks

**Why**:
- svelte-vscode extension is most mature
- HMR works best locally
- Browser dev tools integration

---

### Scenario 4: Full-Stack Development

**Recommended**:
- **IDE**: VS Code (supports both Go and Svelte well)
- **Backend**: Remote on Coder (resource-intensive builds)
- **Frontend**: Local (fast iteration)
- **AI**: Claude Code
- **Git**: Local hooks, push to remote

**Why**:
- Best of both worlds
- Frontend HMR locally
- Backend resources on server
- Unified IDE experience

---

### Scenario 5: Documentation / Scripts

**Recommended**:
- **IDE**: Zed (lightweight, fast startup)
- **AI**: Claude Code with doc generation skills
- **Git**: Local hooks

**Why**:
- Fast editing
- Markdown support
- Python script support

---

## Settings Precedence

### Precedence Rules

1. **EditorConfig beats IDE** - If `.editorconfig` specifies tab settings, IDE must follow
2. **Workspace beats User** - Workspace settings override user settings
3. **Language-specific beats General** - Language-specific settings override general
4. **Formatter config beats IDE** - `ruff.toml` overrides IDE format settings for Python

### Example: Python Indentation

```
Precedence (highest to lowest):
1. .editorconfig       → indent_size = 4, indent_style = space
2. ruff.toml           → (honors EditorConfig)
3. .vscode/settings.json → "python.formatting.provider": "ruff"
4. User settings       → (ignored if workspace settings exist)
```

Result: All tools indent Python with 4 spaces.

---

## Integration Best Practices

### 1. Use EditorConfig for Universal Settings

Always define core settings in `.editorconfig`:
- Indentation (tabs vs spaces, size)
- End of line (LF)
- Trim trailing whitespace
- Insert final newline

### 2. Let Formatters Handle Formatting

Don't configure formatting in IDE settings. Use:
- **Go**: goimports (via gopls)
- **Python**: ruff
- **TypeScript/Svelte**: prettier

IDE just triggers the formatter on save.

### 3. Share LSP Configurations

Both VS Code and Zed should use identical LSP settings:
- gopls config in both `.vscode/settings.json` and `.zed/settings.json`
- ruff config centralized in `ruff.toml`

### 4. Synchronize Git Hooks

Same hooks run locally and in CI:
- Local: pre-commit framework
- CI: Same checks in GitHub Actions

### 5. Document Tool-Specific Quirks

If a tool behaves differently, document in:
- `.tool/docs/TROUBLESHOOTING.md`
- `.shared/docs/TOOL_COMPARISON.md`

---

## Verification

### Check Settings Synchronization

```bash
# Check Go indentation
grep -r "tab" .editorconfig .vscode/settings.json .zed/settings.json

# Check Python line length
grep -r "88" ruff.toml .vscode/settings.json .zed/settings.json

# Check formatters
grep -r "formatOnSave" .vscode/settings.json .zed/settings.json
```

### Test Integration

```bash
# Test Git hooks work
git add . && git commit -m "test: verify hooks"

# Test formatters match
go fmt ./...
ruff format .
prettier --write "web/**/*.{ts,svelte}"

# Test LSPs work
# Open file in VS Code → Check Go to Definition works
# Open same file in Zed → Check Go to Definition works
```

---

## Troubleshooting Integration Issues

### Problem: Formatters conflict between IDEs

**Solution**:
1. Check `.editorconfig` is respected
2. Verify both IDEs use same formatter binary
3. Check formatter version matches

### Problem: LSP works in VS Code but not Zed

**Solution**:
1. Verify LSP binary is in PATH
2. Check Zed logs: `View → Debug → Log File`
3. Restart LSP: Command palette → "Restart LSP"

### Problem: Git hooks don't run

**Solution**:
1. Verify hooks installed: `pre-commit install`
2. Check Git config: `git config core.hooksPath`
3. Re-install: `pre-commit install --install-hooks`

### Problem: Settings not syncing to Coder

**Solution**:
1. Check Coder template includes settings
2. Verify persistent volume mounted
3. Manually copy settings: `scp settings.json coder-workspace:~/.config/`

---

## Related Documentation

- [ONBOARDING.md](ONBOARDING.md) - Getting started guide
- [WORKFLOWS.md](WORKFLOWS.md) - Development workflows
- [TOOL_COMPARISON.md](TOOL_COMPARISON.md) - When to use which tool
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues

### Tool-Specific Docs

- [Claude Code Integration](../../.claude/docs/INDEX.md)
- [VS Code Setup](../../.vscode/docs/INDEX.md)
- [Zed Setup](../../.zed/docs/INDEX.md)
- [Coder Remote Dev](../../.coder/docs/INDEX.md)
- [Git Hooks](../../.githooks/docs/INDEX.md)
- [GitHub Actions](../../.github/docs/INDEX.md)

---

**Maintained By**: Development Team
**Last Updated**: 2026-01-31
