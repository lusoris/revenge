# Remote Development Workflow (Coder)

**Purpose**: Complete guide to remote development on Coder for the Revenge project

**Coder Host**: https://coder.ancilla.lol

**Last Updated**: 2026-01-31

> **📋 Version Requirements**: For all tool and package versions, see [../../docs/dev/design/technical/TECH_STACK.md](../../docs/dev/design/technical/TECH_STACK.md)

---

## Overview

This guide covers the complete remote development workflow using Coder workspaces, from creation to daily development tasks.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Workspace Lifecycle](#workspace-lifecycle)
- [Connecting with IDEs](#connecting-with-ides)
- [Daily Development Workflow](#daily-development-workflow)
- [Common Tasks](#common-tasks)
- [Performance Tips](#performance-tips)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Tools

1. **Coder CLI** - Workspace management
2. **SSH Client** - Remote connection
3. **One IDE**:
   - VS Code (browser or desktop)
   - Zed
   - JetBrains Gateway

### Installation

```bash
# Install Coder CLI
curl -fsSL https://coder.com/install.sh | sh

# Verify installation
coder version

# Login to Coder
coder login https://coder.ancilla.lol
# Follow browser prompts to authenticate
```

---

## Workspace Lifecycle

### 1. Create Workspace

```bash
# List available templates
coder templates list

# Create workspace from revenge template
coder create my-workspace --template revenge

# Or with custom parameters
coder create my-workspace \
  --template revenge \
  --parameter cpu=8 \
  --parameter memory=16 \
  --parameter ide=vscode-browser
```

**Parameters**:
- `cpu`: 2, 4 (recommended), 8, or 16 cores
- `memory`: 4, 8 (recommended), 16, or 32 GB
- `ide`: vscode-browser, vscode-desktop, zed, jetbrains, terminal

### 2. Start/Stop Workspace

```bash
# Start workspace
coder start my-workspace

# Stop workspace (saves resources)
coder stop my-workspace

# Check status
coder list
```

### 3. Delete Workspace

```bash
# Delete when no longer needed
coder delete my-workspace --yes

# Warning: This deletes all data in the workspace
# Make sure to push code to git first!
```

---

## Connecting with IDEs

### Option A: VS Code (Browser)

**Best for**: Quick access, no local setup

```bash
# Open browser-based VS Code
coder open my-workspace

# Or get URL manually
coder list
# Access at: https://coder.ancilla.lol/workspaces/my-workspace/code
```

**What you get**:
- Full VS Code in browser
- All extensions pre-installed
- Terminal access
- Git integration
- No local installation needed

### Option B: VS Code (Desktop)

**Best for**: Full IDE features, better performance

**Step 1: Install Coder Extension**
```bash
code --install-extension coder.coder-remote
```

**Step 2: Connect**
```bash
# Open VS Code and connect
coder code my-workspace

# VS Code will open and connect via SSH automatically
```

**What you get**:
- Full VS Code desktop features
- Local debugger UI
- Better performance than browser
- Local extensions + remote extensions

### Option C: Zed (SSH)

**Best for**: Lightweight, fast performance

**Step 1: Get SSH Config**
```bash
# Display SSH connection info
coder ssh my-workspace

# Output shows:
# Host my-workspace
#   HostName <coder-host>
#   ProxyCommand coder ssh --stdio my-workspace
#   User coder
```

**Step 2: Connect in Zed**
```bash
# Method 1: Direct SSH command
zed ssh://my-workspace

# Method 2: Add to ~/.ssh/config and connect via Zed UI
# File → Open Remote → SSH → my-workspace
```

**What you get**:
- Fast, native Rust performance
- Low bandwidth usage
- Full LSP support (gopls, ruff-lsp)
- Git integration

### Option D: JetBrains Gateway

**Best for**: GoLand/IntelliJ users

**Step 1: Install JetBrains Gateway**
- Download from https://www.jetbrains.com/remote-development/gateway/

**Step 2: Install Coder Plugin**
- Open Gateway
- Plugins → Search "Coder"
- Install and restart

**Step 3: Connect**
- Select Coder provider
- Choose my-workspace
- Select IDE (GoLand or IntelliJ IDEA)
- Connect

**What you get**:
- Full JetBrains IDE features
- Advanced refactoring tools
- Excellent Go support
- Remote debugging

---

## Daily Development Workflow

### Morning Routine

```bash
# 1. Start workspace
coder start my-workspace

# 2. Connect with your IDE
# - VS Code Browser: coder open my-workspace
# - VS Code Desktop: coder code my-workspace
# - Zed: zed ssh://my-workspace
# - JetBrains: Open Gateway → Connect

# 3. Inside workspace, check services
docker ps
# PostgreSQL, Dragonfly, Typesense should be running

# 4. Pull latest code
cd /workspace/revenge
git pull origin develop
```

### Development Cycle

**Backend (Go)**:
```bash
# Start development server with hot reload
air

# In another terminal, run tests
go test ./...

# Access API at http://localhost:8096
# (Automatically forwarded by Coder)
```

**Frontend (Svelte)**:
```bash
cd web
npm run dev

# Access at http://localhost:5173
# (Automatically forwarded by Coder)
```

### End of Day

```bash
# 1. Commit and push changes
git add .
git commit -m "feat: implement feature X"
git push

# 2. Stop workspace (optional, saves resources)
coder stop my-workspace

# Your data persists in the workspace volume
```

---

## Common Tasks

### Port Forwarding

Coder automatically forwards common ports. To manually forward:

```bash
# Forward additional ports
coder port-forward my-workspace --tcp 8080:8080 --tcp 3000:3000

# List forwarded ports
coder ports my-workspace
```

### File Transfer

**Upload to workspace**:
```bash
# Using SCP
scp -r local-folder coder.my-workspace:/workspace/

# Using rsync (faster for large files)
rsync -avz local-folder coder.my-workspace:/workspace/
```

**Download from workspace**:
```bash
# Using SCP
scp -r coder.my-workspace:/workspace/logs ./local-logs

# Using rsync
rsync -avz coder.my-workspace:/workspace/logs ./
```

### Running Tests

```bash
# SSH into workspace
coder ssh my-workspace

# Go tests
go test ./...
go test -v -race ./...
go test -coverprofile=coverage.out ./...

# Frontend tests
cd web
npm run test
npm run test:e2e

# Python tests (for scripts)
pytest scripts/
```

### Viewing Logs

```bash
# Workspace logs
coder agent logs my-workspace

# Service logs (inside workspace)
coder ssh my-workspace
docker logs revenge-postgres
docker logs revenge-dragonfly
docker logs revenge-typesense
```

### Workspace Shell Access

```bash
# SSH into workspace
coder ssh my-workspace

# Or one-off command
coder ssh my-workspace -- git status
coder ssh my-workspace -- go test ./...
```

---

## Performance Tips

### 1. Use Local IDE for Frontend

Frontend development benefits from local HMR (Hot Module Reload):

```bash
# Run backend remotely on Coder
coder ssh my-workspace -- air

# Run frontend locally
cd web
npm run dev
```

### 2. Minimize VS Code Extensions (Browser)

Browser VS Code is slower with many extensions:

```json
{
  "remote.SSH.enableExtensions": [
    "golang.go",
    "charliermarsh.ruff",
    "svelte.svelte-vscode"
  ]
}
```

### 3. Use Zed for Quick Edits

Zed is faster than VS Code for quick file edits:

```bash
# Edit single file with Zed
zed ssh://my-workspace /workspace/revenge/config.yaml

# VS Code for complex tasks (debugging, refactoring)
```

### 4. Close Workspace When Not Using

```bash
# Stop workspace to save resources
coder stop my-workspace

# Workspace volume persists
# Start again when needed: coder start my-workspace
```

### 5. Use Git Efficiently

```bash
# Commit frequently, push to remote
# Don't rely on workspace as backup

# Use shallow clones for faster setup
git clone --depth 1 https://github.com/lusoris/revenge.git
```

---

## Troubleshooting

### Workspace Won't Start

**Symptoms**:
- `coder start my-workspace` hangs
- Workspace status stuck at "Starting"

**Solutions**:
```bash
# Check workspace status
coder list

# View logs
coder agent logs my-workspace

# Restart workspace
coder stop my-workspace
coder start my-workspace

# If still failing, rebuild
coder delete my-workspace --yes
coder create my-workspace --template revenge
```

### Can't Connect to Workspace

**Symptoms**:
- SSH connection refused
- IDE won't connect

**Solutions**:
```bash
# Verify workspace is running
coder list

# Test SSH connection
coder ssh my-workspace -- echo "Connection OK"

# Check SSH keys
ssh-add -l

# Re-login to Coder
coder logout
coder login https://coder.ancilla.lol
```

### Port Forwarding Not Working

**Symptoms**:
- Can't access http://localhost:8096
- Services not accessible

**Solutions**:
```bash
# Check which ports are forwarded
coder ports my-workspace

# Manually forward port
coder port-forward my-workspace --tcp 8096:8096

# Inside workspace, verify service is listening
coder ssh my-workspace -- netstat -tulpn | grep 8096
```

### Slow Performance

**Symptoms**:
- LSP is slow
- File operations lag
- Terminal has latency

**Solutions**:
```bash
# Check network latency
ping coder.ancilla.lol

# Reduce LSP features in IDE
# VS Code: Disable inline hints, reduce analysis
# Zed: Disable LSP hints in settings

# Use terminal for git operations instead of IDE
coder ssh my-workspace
git status
git commit -m "message"
```

### Workspace Disk Full

**Symptoms**:
- "No space left on device"
- Build failures
- Can't commit changes

**Solutions**:
```bash
# SSH into workspace
coder ssh my-workspace

# Check disk usage
df -h
du -sh /workspace/*

# Clean up
docker system prune -a
rm -rf node_modules
rm -rf bin/ tmp/
go clean -cache -testcache -modcache

# If still full, increase disk_size parameter
# (requires recreating workspace)
```

---

## Workspace Management

### Resource Usage

```bash
# Check CPU/memory usage
coder stat my-workspace

# View detailed stats
coder list --output json | jq '.[] | select(.name == "my-workspace")'
```

### Multiple Workspaces

```bash
# Create workspaces for different tasks
coder create revenge-main --template revenge
coder create revenge-hotfix --template revenge

# List all workspaces
coder list

# Switch between workspaces
coder stop revenge-main
coder start revenge-hotfix
```

### Workspace Templates

```bash
# List available templates
coder templates list

# View template parameters
coder templates show revenge

# Update workspace with new template version
coder update my-workspace
```

---

## Best Practices

### 1. Commit and Push Frequently

```bash
# Workspaces can be deleted/rebuilt
# Git is your backup, not the workspace

git commit -m "WIP: feature in progress"
git push
```

### 2. Use .gitignore for Large Files

```gitignore
# Don't commit large binaries or generated files
bin/
tmp/
*.log
coverage.out
node_modules/
```

### 3. Stop Workspace When Not Using

```bash
# Save resources for the team
coder stop my-workspace

# Workspace persists, just stopped
# Start again: coder start my-workspace
```

### 4. Keep Workspace Clean

```bash
# Regular cleanup
docker system prune
go clean -cache
npm cache clean --force
```

### 5. Document Workspace-Specific Setup

```bash
# If your workspace needs special setup, document it
echo "export MY_CUSTOM_VAR=value" >> ~/.bashrc
```

---

## Related Documentation

- [ZED_INTEGRATION.md](ZED_INTEGRATION.md) - Zed-specific Coder setup
- [JETBRAINS_INTEGRATION.md](JETBRAINS_INTEGRATION.md) - JetBrains Gateway setup
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Detailed troubleshooting
- [../template.tf](../template.tf) - Workspace template configuration

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
