# JetBrains IDE Documentation

> Complete reference documentation for JetBrains IDEs (GoLand, IntelliJ IDEA) setup for Revenge development

**For**: Revenge Media Server Project
**Last Updated**: 2026-01-31

---

## Core Documentation

| Document | Description | Priority |
|----------|-------------|----------|
| [SETUP.md](SETUP.md) | Local IDE setup, plugins, and project configuration | **START HERE** |
| [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md) | JetBrains Gateway + Coder integration | Essential |
| [TROUBLESHOOTING.md](TROUBLESHOOTING.md) | Common issues and solutions | Reference |
| [CLAUDE.md](CLAUDE.md) | IDE-specific instructions for Claude Code | Reference |

---

## Quick Start

### For New Users

1. **Choose Your IDE**:
   - **GoLand** - Recommended for Go-focused development
   - **IntelliJ IDEA Ultimate** - For multi-language (Go + Svelte + Python)

2. **Local Development**: Follow [SETUP.md](SETUP.md)
   - Install IDE and required plugins
   - Configure Go SDK (1.25.6+)
   - Set up project structure
   - Configure run configurations

3. **Remote Development**: Follow [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md)
   - Install JetBrains Gateway
   - Connect to Coder workspace
   - Use IDE remotely via Gateway

4. **Troubleshoot**: Use [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
   - SDK not detected
   - Performance issues
   - Plugin problems

### Quick Links to Coder Integration

Since remote development on Coder is a Week 1 priority:
- **Complete Gateway Guide**: [../../.coder/docs/JETBRAINS_INTEGRATION.md](../../.coder/docs/JETBRAINS_INTEGRATION.md)
- **Remote Workflow**: [../../.coder/docs/REMOTE_WORKFLOW.md](../../.coder/docs/REMOTE_WORKFLOW.md)

---

## Key Shortcuts (macOS vs Windows/Linux)

| Action | macOS | Windows/Linux | Description |
|--------|-------|---------------|-------------|
| Search Everywhere | `Shift Shift` | `Shift Shift` | Find anything |
| Find Action | `Cmd+Shift+A` | `Ctrl+Shift+A` | Command palette |
| Go to Declaration | `Cmd+B` | `Ctrl+B` | Jump to definition |
| Find Usages | `Cmd+F7` | `Alt+F7` | Find where used |
| Refactor | `Ctrl+T` | `Ctrl+Alt+Shift+T` | Refactoring menu |
| Rename | `Shift+F6` | `Shift+F6` | Rename symbol |
| Extract Method | `Cmd+Alt+M` | `Ctrl+Alt+M` | Extract to method |
| Run | `Ctrl+R` | `Shift+F10` | Run current config |
| Debug | `Ctrl+D` | `Shift+F9` | Debug current config |
| Quick Fix | `Alt+Enter` | `Alt+Enter` | Show intentions |

Full shortcuts: **Help** → **Keyboard Shortcuts PDF**

---

## Language Support

### Go Development

- **IDE**: GoLand (recommended) or IntelliJ IDEA + Go plugin
- **SDK**: Go (see SOURCE_OF_TRUTH) with GOEXPERIMENT=greenteagc,jsonv2
- **Features**: Advanced refactoring, excellent debugger, gopls integration
- **Setup**: [SETUP.md - Go](SETUP.md#go-development)

### Python (Scripts)

- **IDE**: IntelliJ IDEA Ultimate + Python plugin
- **Version**: Python (see SOURCE_OF_TRUTH)
- **Features**: Code intelligence, debugging, package management
- **Setup**: [SETUP.md - Python](SETUP.md#python-development)

### TypeScript/Svelte (Frontend)

- **IDE**: IntelliJ IDEA Ultimate (WebStorm features included)
- **Version**: Node.js (see SOURCE_OF_TRUTH)
- **Features**: TypeScript support, Svelte plugin available
- **Setup**: [SETUP.md - Frontend](SETUP.md#frontend-development)

---

## IDE Comparison

### GoLand vs IntelliJ IDEA Ultimate

| Feature | GoLand | IntelliJ IDEA Ultimate |
|---------|--------|----------------------|
| **Go Support** | ✅ Built-in, optimized | ✅ Via plugin |
| **Go Refactoring** | ✅ Excellent | ✅ Excellent |
| **Go Debugging** | ✅ Excellent | ✅ Excellent |
| **Database Tools** | ✅ Built-in | ✅ Built-in |
| **Python Support** | ⚠️ Via plugin | ✅ Built-in |
| **TypeScript/Svelte** | ⚠️ Via plugin | ✅ Built-in |
| **Price** | 💰 $99/year | 💰 $169/year |
| **Best For** | Go-only development | Full-stack development |

**Recommendation**:
- **Backend developers** → GoLand (cheaper, optimized for Go)
- **Full-stack developers** → IntelliJ IDEA Ultimate (all languages)

---

## Project Structure

```
.jetbrains/
└── docs/
    ├── INDEX.md                    # This file
    ├── SETUP.md                    # Local IDE setup
    ├── REMOTE_DEVELOPMENT.md       # Gateway + Coder
    ├── TROUBLESHOOTING.md          # Common issues
    └── CLAUDE.md                   # IDE-specific Claude Code instructions
```

**Note**: JetBrains IDEs use `.idea/` folder for project settings (git-ignored). There are no shared project settings files like VS Code's `settings.json`.

---

## Common Issues & Solutions

| Issue | Solution | Docs |
|-------|----------|------|
| Go SDK not detected | Configure GOROOT in settings | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#go-sdk-not-detected) |
| Gateway won't connect | Verify Coder login, restart workspace | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#gateway-connection-issues) |
| Slow performance | Reduce memory usage, disable plugins | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#performance-issues) |
| IDE backend crashes | Clear cache, reinstall backend | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#ide-backend-issues) |
| gopls not working | Check Go version, restart IDE | [TROUBLESHOOTING.md](TROUBLESHOOTING.md#gopls-issues) |

---

## Remote Development (Week 1 Priority)

For remote development on Coder:

1. **Install JetBrains Gateway**: [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md#installation)
2. **Install Coder Plugin**: [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md#coder-plugin)
3. **Connect to Workspace**: [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md#connecting)
4. **Start Coding**: [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md#workflow)

Or see the complete guide: [../../.coder/docs/JETBRAINS_INTEGRATION.md](../../.coder/docs/JETBRAINS_INTEGRATION.md)

---

## External Resources

### Official Documentation
- **GoLand**: https://www.jetbrains.com/go/
- **IntelliJ IDEA**: https://www.jetbrains.com/idea/
- **JetBrains Gateway**: https://www.jetbrains.com/remote-development/gateway/
- **Go Plugin**: https://plugins.jetbrains.com/plugin/9568-go

### Project Documentation
- **Revenge Design Index**: [../../docs/dev/design/DESIGN_INDEX.md](../../docs/dev/design/DESIGN_INDEX.md)
- **Tech Stack**: [../../docs/dev/design/technical/TECH_STACK.md](../../docs/dev/design/technical/TECH_STACK.md)
- **Development Guide**: [../../docs/dev/design/operations/DEVELOPMENT.md](../../docs/dev/design/operations/DEVELOPMENT.md)
- **Source of Truth**: [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)

### Related IDE Docs
- **VS Code**: [../../.vscode/docs/INDEX.md](../../.vscode/docs/INDEX.md)
- **Zed**: [../../.zed/docs/INDEX.md](../../.zed/docs/INDEX.md)
- **Coder**: [../../.coder/docs/INDEX.md](../../.coder/docs/INDEX.md)

---

## Getting Help

1. **Check SETUP.md** for installation and configuration
2. **Check TROUBLESHOOTING.md** for common problems
3. **Check IDE logs**: **Help** → **Show Log in Finder/Explorer**
4. **Search JetBrains issues**: https://youtrack.jetbrains.com/issues
5. **Search Revenge issues**: https://github.com/kilianso/revenge/issues

---

## Development Workflow

### Typical Day with GoLand/IntelliJ

```bash
# 1. Open project
# File → Open → /path/to/revenge

# 2. Run development server
# Click ▶ in toolbar or Ctrl+R (macOS) / Shift+F10 (Windows/Linux)

# 3. Edit code
# - Format on save (auto-configured)
# - Use Quick Fix (Alt+Enter) for suggestions
# - Refactor with Ctrl+T / Ctrl+Alt+Shift+T

# 4. Debug
# Set breakpoints (click gutter)
# Click 🐞 or Ctrl+D (macOS) / Shift+F9 (Windows/Linux)

# 5. Run tests
# Right-click test → Run 'TestName'
# Or use Ctrl+Shift+R for all tests

# 6. Use terminal for git
# View → Tool Windows → Terminal
# Or use built-in Git tools (Cmd/Ctrl+K for commit)
```

### Performance Tips

- **Use Power Save Mode** when not actively coding (**File** → **Power Save Mode**)
- **Close unused projects** to save memory
- **Disable unused plugins** (**Settings** → **Plugins**)
- **Increase heap size** if needed (**Help** → **Change Memory Settings**)

---

## Why JetBrains IDEs?

### Advantages over VS Code/Zed

| Feature | JetBrains | VS Code | Zed |
|---------|-----------|---------|-----|
| **Go Refactoring** | ✅ Advanced | ⚠️ Basic | ⚠️ Basic |
| **Debugging** | ✅ Excellent | ✅ Good | ❌ None |
| **Database Tools** | ✅ Built-in | ⚠️ Extension | ❌ None |
| **Code Analysis** | ✅ Deep | ⚠️ Basic | ⚠️ Basic |
| **Refactoring Safety** | ✅ Excellent | ⚠️ Good | ⚠️ Basic |

### When to Use JetBrains

**Use GoLand/IntelliJ when**:
- Complex refactoring needed (renaming across files, extracting interfaces, etc.)
- Advanced debugging (conditional breakpoints, evaluate expression, etc.)
- Database development (built-in query console, schema tools)
- You value deep code intelligence over speed

**Use VS Code when**:
- General development with good extension ecosystem
- Need balance of features and performance
- Prefer free, open-source tools

**Use Zed when**:
- Quick edits and navigation
- Maximum performance is critical
- Lightweight workflow preferred

---

**Maintained by**: Revenge Development Team
**Last Updated**: 2026-01-31
