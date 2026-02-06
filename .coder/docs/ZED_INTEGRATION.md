# Zed Integration with Coder

**Purpose**: Complete guide for using Zed Editor with Coder remote workspaces

**Last Updated**: 2026-01-31

> **📋 Version Requirements**: For all tool and package versions, see [../../docs/dev/design/technical/TECH_STACK.md](../../docs/dev/design/technical/TECH_STACK.md)

---

## Overview

Zed is a fast, lightweight editor built in Rust that works excellently with Coder via SSH. This guide covers setup, configuration, and best practices.

---

## Why Zed + Coder?

**Advantages**:
- ✅ Native performance (faster than VS Code Remote)
- ✅ Low bandwidth usage
- ✅ Excellent SSH support (built-in)
- ✅ Full LSP support (gopls, ruff-lsp)
- ✅ Minimal CPU/RAM on client machine
- ✅ Works great on low-spec laptops

**Best for**:
- Go backend development
- Documentation writing
- Quick file edits
- Low-bandwidth connections
- Battery-constrained laptops

---

## Prerequisites

1. **Zed Installed** - See [.zed/docs/SETUP.md](../../.zed/docs/SETUP.md)
2. **Coder CLI** - `curl -fsSL https://coder.com/install.sh | sh`
3. **SSH Client** - Built into Linux/macOS, or Git Bash on Windows

---

## Setup

### Step 1: Login to Coder

```bash
# Login
coder login https://coder.ancilla.lol

# Verify
coder list
```

### Step 2: Create or Start Workspace

```bash
# Create new workspace
coder create my-workspace --template revenge --parameter ide=zed

# Or start existing workspace
coder start my-workspace
```

### Step 3: Get SSH Connection Info

```bash
# Display SSH config
coder ssh my-workspace

# Output:
# Host my-workspace
#   HostName coder.ancilla.lol
#   ProxyCommand coder ssh --stdio my-workspace
#   User coder
#   StrictHostKeyChecking no
#   UserKnownHostsFile=/dev/null
```

### Step 4: Add to SSH Config (Optional)

```bash
# Edit SSH config
nano ~/.ssh/config

# Add (copy from coder ssh output):
Host coder-revenge
  HostName coder.ancilla.lol
  ProxyCommand coder ssh --stdio my-workspace
  User coder
  StrictHostKeyChecking no
  UserKnownHostsFile=/dev/null
```

---

## Connecting with Zed

### Method 1: Direct SSH URI

```bash
# Connect directly
zed ssh://my-workspace

# Or with custom SSH config name
zed ssh://coder-revenge
```

### Method 2: Zed Remote UI

1. Open Zed
2. **File** → **Open Remote**
3. Select **SSH**
4. Choose `my-workspace` or `coder-revenge`
5. Zed connects and opens remote workspace

### Method 3: Command Line with Path

```bash
# Open specific directory
zed ssh://my-workspace/workspace/revenge

# Open specific file
zed ssh://my-workspace/workspace/revenge/cmd/revenge/main.go
```

---

## Configuration

### Workspace Settings

Inside the remote workspace, Zed automatically loads `.zed/settings.json`:

```json
{
  "tab_size": 2,
  "format_on_save": "on",
  "ensure_final_newline_on_save": true,

  "languages": {
    "Go": {
      "tab_size": 4,
      "hard_tabs": true,
      "format_on_save": "on"
    },
    "Python": {
      "tab_size": 4,
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
        "analyses": {
          "unusedparams": true,
          "shadow": true
        },
        "staticcheck": true
      }
    }
  }
}
```

**All project settings work remotely** - no additional configuration needed!

### Performance Optimization

For remote development, optimize Zed settings:

**~/.config/zed/settings.json** (user settings):

```json
{
  // Reduce network traffic
  "file_scan_exclusions": [
    "**/.git",
    "**/node_modules",
    "**/dist",
    "**/__pycache__",
    "**/bin",
    "**/tmp"
  ],

  // Disable resource-intensive features
  "lsp": {
    "gopls": {
      "initialization_options": {
        "hints": {
          // Disable inlay hints (reduces bandwidth)
          "assignVariableTypes": false,
          "compositeLiteralFields": false
        }
      }
    }
  },

  // Git integration (local)
  "git": {
    "git_gutter": "tracked_files",
    "inline_blame": {
      "enabled": false  // Reduces latency
    }
  }
}
```

---

## Development Workflow

### Opening Project

```bash
# SSH into workspace first to verify
coder ssh my-workspace

# Then open in Zed
zed ssh://my-workspace/workspace/revenge
```

### Terminal Usage

**Option 1: Zed Integrated Terminal**
- Press `Ctrl+\`` to toggle terminal
- Terminal runs on remote workspace
- Full SSH session

**Option 2: Separate SSH Session**
```bash
# In separate terminal
coder ssh my-workspace

# Run commands
cd /workspace/revenge
go test ./...
air
```

### File Operations

**Remote Files**:
- Open: `Cmd/Ctrl+P` → Type filename
- Save: `Cmd/Ctrl+S` (saves to remote)
- Search: `Cmd/Ctrl+Shift+F` (searches remote files)

**Local Files**:
- Use separate Zed instance
- `zed .` (local) vs `zed ssh://...` (remote)

### Git Operations

**In Zed**:
- Git gutter shows changes
- Status in sidebar

**In Terminal** (recommended for commits):
```bash
# SSH into workspace
coder ssh my-workspace

# Git operations
git status
git add .
git commit -m "feat: implement feature"
git push
```

---

## LSP Configuration

### Go (gopls)

**Auto-discovered** - works out of the box.

**Verify**:
1. Open Go file in Zed
2. Hover over function → Should see docs
3. `F12` (Go to Definition) → Should jump

**Troubleshooting**:
```bash
# Inside workspace, check gopls
coder ssh my-workspace -- which gopls
coder ssh my-workspace -- gopls version

# If missing, install
coder ssh my-workspace -- go install golang.org/x/tools/gopls@latest
```

### Python (ruff-lsp)

**Configured in workspace** - uses external ruff formatter.

**Verify**:
1. Open Python file
2. Save → Should auto-format with ruff

**Troubleshooting**:
```bash
# Check ruff is installed
coder ssh my-workspace -- which ruff
coder ssh my-workspace -- ruff --version

# If missing
coder ssh my-workspace -- pip install ruff
```

### TypeScript/Svelte

**Auto-discovered** - Zed uses built-in TypeScript server.

**For Svelte**, ensure prettier plugin is installed:
```bash
coder ssh my-workspace
cd /workspace/revenge
npm install -g prettier-plugin-svelte
```

---

## Common Tasks

### Running Development Server

**Method 1: Zed Terminal**
1. `Ctrl+\`` to open terminal
2. `air` (Go hot reload)
3. Access http://localhost:8096 (forwarded by Coder)

**Method 2: Separate SSH**
```bash
coder ssh my-workspace
air
```

### Running Tests

```bash
# In Zed terminal or separate SSH
go test ./...
go test -v -race ./...
```

### Debugging

Zed doesn't have graphical debugger. Use `delve` CLI:

```bash
# SSH into workspace
coder ssh my-workspace

# Start debugger
dlv debug ./cmd/revenge
(dlv) break main.main
(dlv) continue
(dlv) print variable
```

**Or use VS Code for debugging**, Zed for editing.

---

## Performance Comparison

| Metric | Zed + Coder | VS Code Remote |
|--------|-------------|----------------|
| Startup time | <1 second | 5-10 seconds |
| Memory (client) | 100-200 MB | 500-800 MB |
| Bandwidth | Low | Medium-High |
| LSP responsiveness | Excellent | Good |
| File search | Very fast | Fast |
| Battery impact | Minimal | Moderate |

**Verdict**: Zed is significantly faster and lighter for remote development.

---

## Troubleshooting

### Connection Issues

**Problem**: `zed ssh://my-workspace` doesn't connect

**Solutions**:
```bash
# 1. Verify workspace is running
coder list

# 2. Test SSH manually
coder ssh my-workspace -- echo "OK"

# 3. Check Coder CLI is logged in
coder login https://coder.ancilla.lol

# 4. Try with full SSH config
zed ssh://coder-revenge
```

### LSP Not Working

**Problem**: gopls or ruff-lsp not starting

**Solutions**:
```bash
# Check Zed logs
# macOS: ~/Library/Logs/Zed/Zed.log
# Linux: ~/.local/share/zed/logs/Zed.log

# Restart LSP in Zed
# Cmd/Ctrl+Shift+P → "Restart Language Server"

# Verify LSP binary exists on remote
coder ssh my-workspace -- which gopls
coder ssh my-workspace -- which ruff-lsp

# Reinstall if missing
coder ssh my-workspace
go install golang.org/x/tools/gopls@latest
pip install ruff-lsp
```

### Format on Save Not Working

**Problem**: Files not formatting automatically

**Solutions**:

1. **Check Zed settings** (`.zed/settings.json`):
```json
{
  "format_on_save": "on"
}
```

2. **Check formatter is available**:
```bash
coder ssh my-workspace -- which goimports
coder ssh my-workspace -- which ruff
```

3. **Manual format**:
- `Cmd/Ctrl+Shift+I` in Zed

### Slow Performance

**Problem**: Laggy typing, slow LSP

**Solutions**:

1. **Check network latency**:
```bash
ping coder.ancilla.lol
```

2. **Disable LSP hints**:
```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "hints": {
          "assignVariableTypes": false
        }
      }
    }
  }
}
```

3. **Reduce file scanning**:
```json
{
  "file_scan_exclusions": [
    "**/.git",
    "**/node_modules",
    "**/bin"
  ]
}
```

---

## Best Practices

### 1. Use Zed for Editing, Terminal for Git

```bash
# Edit in Zed
zed ssh://my-workspace

# Git in separate terminal
coder ssh my-workspace
git status
git commit -m "message"
```

### 2. Keep Workspace Settings in Sync

All Zed settings in `.zed/settings.json` are version-controlled and work remotely.

### 3. Use VS Code for Debugging

```bash
# Edit with Zed (fast)
zed ssh://my-workspace

# Debug with VS Code (graphical debugger)
coder code my-workspace
```

### 4. Close Unused Tabs

Zed keeps SSH connection alive for each tab. Close unused tabs to reduce connections.

### 5. Use Cmd+P for Navigation

Faster than file explorer:
- `Cmd/Ctrl+P` → Search files
- `Cmd/Ctrl+Shift+O` → Search symbols
- `Cmd/Ctrl+Shift+F` → Search content

---

## Keyboard Shortcuts (Remote)

All standard Zed shortcuts work over SSH:

| Action | macOS | Linux/Windows |
|--------|-------|---------------|
| Go to File | `Cmd+P` | `Ctrl+P` |
| Command Palette | `Cmd+Shift+P` | `Ctrl+Shift+P` |
| Go to Definition | `F12` | `F12` |
| Format Document | `Cmd+Shift+I` | `Ctrl+Shift+I` |
| Toggle Terminal | `Ctrl+\`` | `Ctrl+\`` |
| Find & Replace | `Cmd+H` | `Ctrl+H` |

See [.zed/docs/KEYBINDINGS.md](../../.zed/docs/KEYBINDINGS.md) for complete reference.

---

## Comparison: Zed vs VS Code Remote

| Feature | Zed | VS Code Remote |
|---------|-----|----------------|
| **Setup** | Simple (SSH) | Medium (extension) |
| **Performance** | Excellent | Good |
| **LSP** | Full support | Full support |
| **Debugging** | Terminal only | Graphical UI ✅ |
| **Extensions** | Limited | Extensive ✅ |
| **Battery usage** | Low | Medium |
| **Bandwidth** | Low | Medium |
| **Startup time** | <1 sec | 5-10 sec |

**Recommendation**: Use Zed for editing, VS Code for debugging.

---

## Related Documentation

- [REMOTE_WORKFLOW.md](REMOTE_WORKFLOW.md) - Complete remote development guide
- [.zed/docs/SETUP.md](../../.zed/docs/SETUP.md) - Zed installation and setup
- [.zed/docs/KEYBINDINGS.md](../../.zed/docs/KEYBINDINGS.md) - Keyboard shortcuts
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
