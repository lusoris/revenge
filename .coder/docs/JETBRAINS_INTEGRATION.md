# JetBrains Gateway Integration with Coder

**Purpose**: Complete guide for using JetBrains IDEs (GoLand, IntelliJ IDEA) with Coder

**Last Updated**: 2026-01-31

> **📋 Version Requirements**: For all tool and package versions, see [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)

---

## Overview

JetBrains Gateway provides remote development with GoLand and IntelliJ IDEA through Coder workspaces, offering advanced Go development features.

---

## Why JetBrains + Coder?

**Advantages**:
- ✅ Best-in-class Go support
- ✅ Advanced refactoring tools
- ✅ Powerful debugger
- ✅ Database tools built-in
- ✅ All JetBrains features remotely

**Best for**:
- Complex Go refactoring
- Advanced debugging
- Database development
- Existing JetBrains users

**Tradeoffs**:
- Requires JetBrains license (paid)
- Higher resource usage than Zed
- Longer startup time than VS Code

---

## Prerequisites

1. **JetBrains Gateway** - Free download
2. **JetBrains IDE License** - GoLand or IntelliJ IDEA Ultimate
3. **Coder CLI** - For workspace management
4. **Good internet connection** - Gateway uses more bandwidth

---

## Installation

### Step 1: Install JetBrains Gateway

**macOS**:
```bash
brew install --cask jetbrains-gateway

# Or download from:
# https://www.jetbrains.com/remote-development/gateway/
```

**Linux**:
```bash
# Download and extract
wget https://download.jetbrains.com/idea/gateway/JetBrainsGateway-*.tar.gz
tar -xzf JetBrainsGateway-*.tar.gz

# Run
./gateway.sh
```

**Windows**:
- Download installer from https://www.jetbrains.com/remote-development/gateway/
- Run installer

### Step 2: Install Coder Plugin

1. Open JetBrains Gateway
2. **Settings** → **Plugins**
3. Search for **"Coder"**
4. Click **Install**
5. Restart Gateway

### Step 3: Configure Coder

1. In Gateway, select **Coder** provider
2. Enter Coder URL: `https://coder.ancilla.lol`
3. Click **Login**
4. Browser opens → Authenticate
5. Gateway connects to Coder

---

## Connecting to Workspace

### First-Time Connection

1. **Open Gateway**
2. Select **Coder** provider
3. Choose workspace: `my-workspace`
4. Select IDE:
   - **GoLand** (recommended for Go)
   - **IntelliJ IDEA Ultimate** (for multi-language)
5. Click **Connect**

**What happens**:
- Gateway downloads IDE backend to workspace
- Launches IDE client locally
- Connects via SSH
- Opens `/workspace/revenge`

### Subsequent Connections

1. Open Gateway
2. Recent connections shown
3. Click **Connect** on workspace
4. IDE opens immediately

---

## IDE Setup

### GoLand Configuration

**Automatic**:
- Go SDK detected (`/usr/local/go`)
- GOPATH configured
- Modules enabled
- gopls running

**Verify**:
1. **File** → **Settings** → **Go** → **GOROOT**
   - Should show: `/usr/local/go`
2. **Go Modules** enabled
3. **gopls** running (status bar)

**Manual Configuration** (if needed):
1. **File** → **Settings** → **Go** → **GOROOT**
2. Click **Add SDK** → **Local** → `/usr/local/go`
3. Enable **Go Modules** (vgo)

### IntelliJ IDEA Configuration

**Go Plugin**:
1. **File** → **Settings** → **Plugins**
2. Search "Go"
3. Install official Go plugin
4. Restart IDE
5. Configure GOROOT (same as GoLand above)

### Project Structure

**Open Project**:
- Gateway automatically opens `/workspace/revenge`
- If not, **File** → **Open** → `/workspace/revenge`

**Mark Directories**:
- `cmd/revenge` → Sources Root
- `internal/` → Sources Root
- `tests/` → Test Sources Root

### Run Configurations

**Go Application**:
1. **Run** → **Edit Configurations**
2. **+** → **Go Build**
3. **Name**: Revenge Server
4. **Run kind**: Package
5. **Package path**: `github.com/lusoris/revenge/cmd/revenge`
6. **Working directory**: `/workspace/revenge`
7. **Environment**: `GOEXPERIMENT=greenteagc,jsonv2`
8. Click **OK**

**Run**: Click ▶ or `Shift+F10`

### Debugging

**Set Breakpoints**:
1. Click gutter next to line number
2. Red dot appears

**Start Debugger**:
1. **Run** → **Debug 'Revenge Server'**
2. Or click 🐞 icon
3. Execution stops at breakpoint

**Debug Actions**:
- **F8** - Step Over
- **F7** - Step Into
- **Shift+F8** - Step Out
- **F9** - Resume

**Inspect Variables**:
- Hover over variable
- Variables pane shows all locals
- Evaluate expression: `Alt+F8`

---

## Database Tools

GoLand/IntelliJ IDEA include Database tools:

### Connect to PostgreSQL

1. **Database** tool window (right sidebar)
2. **+** → **Data Source** → **PostgreSQL**
3. **Host**: `localhost`
4. **Port**: `5432`
5. **Database**: `revenge`
6. **User**: `revenge`
7. **Password**: `revenge` (from workspace)
8. Test Connection
9. Click **OK**

### Query Console

1. Right-click database → **New** → **Query Console**
2. Write SQL:
```sql
SELECT * FROM users LIMIT 10;
```
3. **Ctrl+Enter** to execute

---

## Common Tasks

### Running Tests

**All Tests**:
1. Right-click `internal/` folder
2. **Run 'go test internal/...'**

**Specific Package**:
1. Right-click package folder
2. **Run 'go test <package>'**

**Single Test**:
1. Click green ▶ next to `func TestXxx`
2. **Run 'TestXxx'**

**With Coverage**:
1. Right-click folder/file
2. **Run '...' with Coverage**
3. Coverage report shows in editor

### Refactoring

**Rename**:
1. Right-click symbol → **Refactor** → **Rename**
2. Or `Shift+F6`
3. Enter new name
4. **Refactor** → Updates all references

**Extract Method**:
1. Select code block
2. **Refactor** → **Extract** → **Method**
3. Or `Cmd/Ctrl+Alt+M`
4. Name method
5. **OK** → Method extracted

**Move**:
1. Right-click file/symbol → **Refactor** → **Move**
2. Or `F6`
3. Select destination
4. **Refactor**

### Code Generation

**Generate**:
- `Cmd/Ctrl+N` in code
- Choose:
  - Constructor
  - Getter/Setter
  - Test
  - Implement Methods

**Live Templates**:
- Type `for` → `Tab` → for loop
- Type `if err` → `Tab` → error check
- Type `t` → `Tab` → test function

---

## Performance

### Recommended Settings

**Settings** → **Appearance & Behavior** → **System Settings**:
- ☑ **Power Save Mode** - When not actively coding
- ☐ **Reopen projects on startup** - Faster Gateway startup
- **Memory Settings**: 4GB heap (for remote backend)

**Settings** → **Editor** → **General**:
- ☐ **Show parameter name hints** - Reduces bandwidth
- ☐ **Show chain call type hints** - Reduces bandwidth

### Port Forwarding

Gateway automatically forwards common ports:
- `8096` - Revenge API
- `5173` - Frontend dev server
- `5432` - PostgreSQL
- `6379` - Dragonfly

Access at `http://localhost:8096` from your local browser.

---

## Troubleshooting

### Gateway Can't Connect

**Problem**: "Connection failed" error

**Solutions**:
```bash
# 1. Verify workspace is running
coder list

# 2. Check Coder CLI is logged in
coder login https://coder.ancilla.lol

# 3. Restart workspace
coder stop my-workspace
coder start my-workspace

# 4. Re-login in Gateway
# Settings → Coder → Logout → Login again
```

### IDE Backend Won't Start

**Problem**: "Failed to start IDE backend"

**Solutions**:
1. **Delete IDE backend cache**:
   ```bash
   coder ssh my-workspace
   rm -rf ~/.cache/JetBrains/
   ```

2. **Reconnect in Gateway**:
   - Remove connection
   - Add new connection
   - Let Gateway re-download backend

### Slow Performance

**Problem**: Laggy typing, slow IntelliSense

**Solutions**:
1. **Check network latency**:
   ```bash
   ping coder.ancilla.lol
   ```

2. **Reduce memory usage**:
   - **Help** → **Change Memory Settings**
   - Reduce heap to 2GB (if on slow connection)

3. **Disable unused plugins**:
   - **Settings** → **Plugins**
   - Disable unused plugins

4. **Enable Power Save Mode**:
   - **File** → **Power Save Mode**
   - Disables inspections (temporary)

### Go SDK Not Detected

**Problem**: "Go SDK not configured"

**Solutions**:
```bash
# Verify Go is installed in workspace
coder ssh my-workspace -- go version

# If missing, install in workspace
coder ssh my-workspace
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

# In IDE: Settings → Go → GOROOT
# Set to: /usr/local/go
```

---

## Best Practices

### 1. Use Power Save Mode When Not Coding

**File** → **Power Save Mode**

Disables background inspections, saves CPU/bandwidth.

### 2. Close Project When Done

**File** → **Close Project**

Stops IDE backend on server, saves resources.

### 3. Use Database Tools

Instead of separate DB client, use built-in Database tools.

### 4. Create Run Configurations

Save time with pre-configured run/debug configurations.

### 5. Learn Keyboard Shortcuts

- `Cmd/Ctrl+Shift+A` - Action search
- `Shift Shift` - Search Everywhere
- `Cmd/Ctrl+B` - Go to Declaration
- `Alt+Enter` - Quick Fix

---

## Comparison: GoLand vs VS Code vs Zed

| Feature | GoLand | VS Code | Zed |
|---------|--------|---------|-----|
| **Go Refactoring** | ✅ Advanced | ⚠️ Basic | ⚠️ Basic |
| **Debugging** | ✅ Excellent | ✅ Good | ❌ Terminal only |
| **Database Tools** | ✅ Built-in | ⚠️ Extension | ❌ None |
| **Performance** | ⚠️ Medium | ✅ Good | ✅ Excellent |
| **Startup Time** | ⚠️ Slow (10-20s) | ✅ Medium (5-10s) | ✅ Fast (<1s) |
| **Resource Usage** | ⚠️ High | ✅ Medium | ✅ Low |
| **Cost** | 💰 Paid | ✅ Free | ✅ Free |

**Recommendation**:
- **Complex refactoring** → GoLand
- **General development** → VS Code
- **Quick edits** → Zed

---

## Keyboard Shortcuts

| Action | macOS | Windows/Linux |
|--------|-------|---------------|
| Search Everywhere | `Shift Shift` | `Shift Shift` |
| Find Action | `Cmd+Shift+A` | `Ctrl+Shift+A` |
| Go to Declaration | `Cmd+B` | `Ctrl+B` |
| Find Usages | `Cmd+F7` | `Alt+F7` |
| Refactor | `Ctrl+T` | `Ctrl+Alt+Shift+T` |
| Rename | `Shift+F6` | `Shift+F6` |
| Extract Method | `Cmd+Alt+M` | `Ctrl+Alt+M` |
| Run | `Ctrl+R` | `Shift+F10` |
| Debug | `Ctrl+D` | `Shift+F9` |
| Quick Fix | `Alt+Enter` | `Alt+Enter` |

Full list: **Help** → **Keyboard Shortcuts PDF**

---

## Related Documentation

- [REMOTE_WORKFLOW.md](REMOTE_WORKFLOW.md) - Complete remote development guide
- [ZED_INTEGRATION.md](ZED_INTEGRATION.md) - Zed Editor integration
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues
- [JetBrains Gateway Docs](https://www.jetbrains.com/help/idea/remote-development-overview.html)

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
