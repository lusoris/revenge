# JetBrains Remote Development with Coder

**Purpose**: Quick start guide for using JetBrains Gateway with Coder workspaces

**Last Updated**: 2026-01-31

---

## Overview

This guide covers the essentials for remote development with JetBrains Gateway + Coder. For comprehensive details, see [../../.coder/docs/JETBRAINS_INTEGRATION.md](../../.coder/docs/JETBRAINS_INTEGRATION.md).

**Week 1 Priority**: Get remote development working on Coder with JetBrains Gateway.

---

## Quick Start

### Prerequisites

- ✅ JetBrains Gateway installed
- ✅ Coder CLI installed and logged in
- ✅ JetBrains IDE license (GoLand or IntelliJ IDEA Ultimate)
- ✅ Access to https://coder.ancilla.lol

### 5-Minute Setup

```bash
# 1. Install JetBrains Gateway (macOS)
brew install --cask jetbrains-gateway

# 2. Login to Coder
coder login https://coder.ancilla.lol

# 3. Verify workspace
coder list

# 4. Open Gateway
open /Applications/JetBrains\ Gateway.app

# 5. In Gateway:
#    - Select "Coder" provider
#    - Login to https://coder.ancilla.lol
#    - Select your workspace
#    - Choose IDE (GoLand or IntelliJ IDEA)
#    - Click "Connect"
```

Gateway downloads the IDE backend to your workspace and opens the IDE locally.

---

## Installation

### Step 1: Install JetBrains Gateway

**macOS**:
```bash
brew install --cask jetbrains-gateway
```

**Linux**:
```bash
wget https://download.jetbrains.com/idea/gateway/JetBrainsGateway-*.tar.gz
tar -xzf JetBrainsGateway-*.tar.gz
./gateway.sh
```

**Windows**:
- Download from https://www.jetbrains.com/remote-development/gateway/
- Run installer

### Step 2: Install Coder Plugin

1. Open JetBrains Gateway
2. **Settings** → **Plugins**
3. Search **"Coder"**
4. Click **Install**
5. Restart Gateway

---

## Connecting to Coder Workspace

### First-Time Connection

1. **Open JetBrains Gateway**

2. **Select Coder Provider**:
   - Click **Coder** in the provider list
   - If not visible, ensure plugin is installed

3. **Configure Coder URL**:
   - Enter: `https://coder.ancilla.lol`
   - Click **Login**
   - Browser opens for authentication
   - Authorize and return to Gateway

4. **Select Workspace**:
   - Choose your workspace from the list
   - Example: `revenge-dev`

5. **Choose IDE**:
   - **GoLand** - Recommended for Go development
   - **IntelliJ IDEA Ultimate** - For full-stack (Go + Svelte + Python)

6. **Click Connect**:
   - Gateway downloads IDE backend to workspace (~1-2 minutes first time)
   - IDE opens locally, connected to remote workspace
   - Project opens at `/workspace/revenge`

### Subsequent Connections

1. Open Gateway
2. Recent connections show automatically
3. Click **Connect** on your workspace
4. IDE opens immediately (backend already installed)

---

## Workflow

### Daily Development

**Morning**:
```bash
# Start workspace if stopped
coder start revenge-dev

# Open Gateway
open /Applications/JetBrains\ Gateway.app

# Click Connect on revenge-dev
# IDE opens in seconds
```

**During Day**:
- Code in IDE (all JetBrains features available)
- Run/debug configurations work normally
- Database tools connect to workspace services
- Terminal opens in workspace context

**Evening**:
```bash
# Close IDE (File → Close Project)
# IDE backend stops automatically

# Optionally stop workspace to save resources
coder stop revenge-dev
```

### Port Forwarding

Gateway automatically forwards common ports:
- `8096` - Revenge API
- `5173` - Frontend dev server
- `5432` - PostgreSQL
- `6379` - Dragonfly

Access from local browser: `http://localhost:8096`

### Running Services

All services run in the workspace:

```bash
# In IDE terminal (Alt+F12 or View → Tool Windows → Terminal)

# Start Revenge server
air

# Run tests
go test ./...

# Start frontend
cd web
npm run dev
```

Access via forwarded ports from your local machine.

---

## IDE Configuration

### Go SDK

**Auto-Configured**:
- GOROOT: `/usr/local/go` (in workspace)
- GOPATH: Auto-detected
- Go Modules: Enabled
- gopls: Running

**Verify**:
1. **Settings** → **Go** → **GOROOT**
2. Should show `/usr/local/go`
3. **Go Modules**: ✅ Enabled

### Run Configuration

Gateway should auto-create run configurations, but if not:

1. **Run** → **Edit Configurations** → **+** → **Go Build**
2. Configure:
   - **Name**: `Revenge Server`
   - **Package path**: `github.com/lusoris/revenge/cmd/revenge`
   - **Working directory**: `/workspace/revenge`
   - **Environment**: `GOEXPERIMENT=greenteagc,jsonv2`
3. Click **OK**

### Database Connection

1. **Database** tool window (right sidebar)
2. **+** → **Data Source** → **PostgreSQL**
3. Configure:
   - **Host**: `localhost` (in workspace)
   - **Port**: `5432`
   - **Database**: `revenge`
   - **User**: `revenge`
   - **Password**: `revenge`
4. **Test Connection** → Download driver if prompted
5. **OK**

---

## Debugging

### Set Breakpoints

1. Open Go file (e.g., `internal/api/handlers/movie.go`)
2. Click gutter next to line number
3. Red dot appears

### Start Debugger

1. **Run** → **Debug 'Revenge Server'**
2. Or click 🐞 icon in toolbar
3. Server starts in debug mode
4. Execution stops at breakpoint

### Debug Actions

- **F8** - Step Over
- **F7** - Step Into
- **Shift+F8** - Step Out
- **F9** - Resume Program
- **Alt+F8** - Evaluate Expression

### Inspect Variables

- Hover over variable to see value
- **Variables** pane shows all locals
- **Watches** pane for custom expressions

---

## Performance Tips

### Optimize Gateway Connection

**Settings in IDE** → **Appearance & Behavior** → **System Settings**:
- ☑ **Power Save Mode** when not coding (saves CPU)
- ☐ **Reopen projects on startup** (faster Gateway startup)
- **Heap size**: 4GB (for backend in workspace)

### Reduce Bandwidth Usage

**Settings** → **Editor** → **General**:
- ☐ **Show parameter name hints** (reduces bandwidth)
- ☐ **Show chain call type hints**

### Network Latency

If experiencing lag:

1. **Check latency**:
   ```bash
   ping coder.ancilla.lol
   ```
   - <50ms: Excellent
   - 50-100ms: Good
   - 100-200ms: Usable
   - >200ms: May experience lag

2. **Reduce memory usage**:
   - **Help** → **Change Memory Settings** → Reduce to 2GB

3. **Disable unused plugins**:
   - **Settings** → **Plugins** → Disable unused

4. **Enable Power Save Mode**:
   - **File** → **Power Save Mode**

---

## Troubleshooting

### Gateway Can't Connect

**Problem**: "Connection failed" error

**Solutions**:
```bash
# 1. Verify workspace is running
coder list

# 2. Start workspace if stopped
coder start revenge-dev

# 3. Check Coder CLI login
coder login https://coder.ancilla.lol

# 4. Re-login in Gateway
# Settings → Coder → Logout → Login again
```

### IDE Backend Won't Start

**Problem**: "Failed to start IDE backend"

**Solutions**:

1. **Delete backend cache**:
   ```bash
   coder ssh revenge-dev
   rm -rf ~/.cache/JetBrains/
   exit
   ```

2. **Reconnect in Gateway**:
   - Remove connection from recent list
   - Add new connection
   - Gateway re-downloads backend

### Go SDK Not Detected

**Problem**: "Go SDK not configured"

**Solutions**:

```bash
# Verify Go in workspace
coder ssh revenge-dev -- go version

# If missing, install Go in workspace
coder ssh revenge-dev

# Download and install Go
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz (check SOURCE_OF_TRUTH for version)
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz (check SOURCE_OF_TRUTH for version)
exit
```

Then in IDE:
1. **Settings** → **Go** → **GOROOT**
2. Set to `/usr/local/go`
3. **OK**

### Slow Performance

**Problem**: Laggy typing, slow completions

**Solutions**:

1. **Enable Power Save Mode**: **File** → **Power Save Mode**
2. **Reduce heap**: **Help** → **Change Memory Settings** → 2GB
3. **Disable plugins**: **Settings** → **Plugins** → Disable unused
4. **Check network**: `ping coder.ancilla.lol` (should be <100ms)

### Port Forwarding Not Working

**Problem**: Can't access `http://localhost:8096`

**Solutions**:

```bash
# Verify service is running in workspace
coder ssh revenge-dev
curl http://localhost:8096/health
exit

# Manually forward port if needed
coder port-forward revenge-dev --tcp 8096:8096
```

Then access `http://localhost:8096` in local browser.

---

## Comparison: Gateway vs Local

| Aspect | Gateway (Remote) | Local IDE |
|--------|------------------|-----------|
| **Setup** | ✅ One-time Gateway install | ⚠️ Full stack on local machine |
| **Performance** | ⚠️ Depends on network | ✅ Maximum (local CPU/disk) |
| **Resources** | ✅ Use powerful remote machines | ⚠️ Limited by local hardware |
| **Access** | ✅ From any machine | ⚠️ Only from configured machine |
| **Cost** | 💰 Coder workspace cost | ✅ Free (own hardware) |
| **Startup** | ⚠️ 10-20s (backend + connection) | ✅ 5-10s (local launch) |

**When to Use Gateway**:
- Need powerful remote machine
- Want same environment everywhere
- Working from multiple locations
- Team wants consistent environments

**When to Use Local**:
- Maximum performance needed
- Poor/unreliable network
- Prefer working offline
- Have powerful local machine

---

## Advanced Features

### JetBrains Gateway Features

**All features work remotely**:
- ✅ Advanced refactoring (Extract Interface, Move Package, etc.)
- ✅ Debugger with conditional breakpoints
- ✅ Database tools and query console
- ✅ Profiler and performance analysis
- ✅ HTTP client for API testing
- ✅ Diagram generation (class diagrams, Go structure)

### Database Tools

Access workspace PostgreSQL:

1. Connect (see [Database Connection](#database-connection))
2. **Query Console**: Right-click database → **New** → **Query Console**
3. Run queries:
   ```sql
   SELECT * FROM users LIMIT 10;
   ```
4. **Ctrl+Enter** to execute
5. Results in bottom panel

### HTTP Client

Test APIs from IDE:

1. Create `.http` file:
   ```http
   ### Get Movies
   GET http://localhost:8096/api/v1/movies

   ### Search Movies
   GET http://localhost:8096/api/v1/movies/search?q=matrix
   ```

2. Click **▶** next to request
3. Results show in panel

---

## Related Documentation

- **Complete Gateway Guide**: [../../.coder/docs/JETBRAINS_INTEGRATION.md](../../.coder/docs/JETBRAINS_INTEGRATION.md)
- **Remote Workflow**: [../../.coder/docs/REMOTE_WORKFLOW.md](../../.coder/docs/REMOTE_WORKFLOW.md)
- **Local Setup**: [SETUP.md](SETUP.md)
- **Troubleshooting**: [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
- **Tool Comparison**: [../../.shared/docs/TOOL_COMPARISON.md](../../.shared/docs/TOOL_COMPARISON.md)

---

## Best Practices

### 1. Close Project When Done

**File** → **Close Project** stops the IDE backend, saving workspace resources.

### 2. Use Power Save Mode When Not Coding

**File** → **Power Save Mode** disables background inspections, reducing CPU usage.

### 3. Create Run Configurations

Save frequently used run/debug configurations for quick access.

### 4. Use Database Tools

Built-in database tools are powerful - no need for separate DB client.

### 5. Learn Keyboard Shortcuts

- `Cmd/Ctrl+Shift+A` - Find Action (command palette)
- `Shift Shift` - Search Everywhere
- `Cmd/Ctrl+B` - Go to Declaration
- `Alt+Enter` - Quick Fix

Full list: **Help** → **Keyboard Shortcuts PDF**

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
