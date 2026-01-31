# JetBrains IDE Setup for Revenge

**Purpose**: Complete installation and configuration guide for GoLand and IntelliJ IDEA

**Last Updated**: 2026-01-31

---

## Overview

This guide covers local setup for JetBrains IDEs. For remote development with JetBrains Gateway + Coder, see [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md) or [../../.coder/docs/JETBRAINS_INTEGRATION.md](../../.coder/docs/JETBRAINS_INTEGRATION.md).

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [IDE Installation](#ide-installation)
- [Go Development](#go-development)
- [Python Development](#python-development)
- [Frontend Development](#frontend-development)
- [Database Tools](#database-tools)
- [Essential Plugins](#essential-plugins)
- [Project Configuration](#project-configuration)
- [Run Configurations](#run-configurations)
- [Code Style Settings](#code-style-settings)
- [Version Control](#version-control)

---

## Prerequisites

> **📋 Version Requirements**: See [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md) for exact required versions.

### Required Software

- **Go** - https://go.dev/dl/ (version: see SOURCE_OF_TRUTH)
- **Python** - https://www.python.org/downloads/ (version: see SOURCE_OF_TRUTH)
- **Node.js** - https://nodejs.org/ (version: see SOURCE_OF_TRUTH)
- **Git** - https://git-scm.com/downloads
- **Docker** - https://www.docker.com/get-started (for local services)

### Verify Installation

```bash
# Check Go version (compare with SOURCE_OF_TRUTH)
go version

# Check Python version (compare with SOURCE_OF_TRUTH)
python3 --version

# Check Node version (compare with SOURCE_OF_TRUTH)
node --version

# Check Git
git --version

# Check Docker
docker --version
```

---

## IDE Installation

### GoLand

**Download**: https://www.jetbrains.com/go/download/

**macOS**:
```bash
# Using Homebrew
brew install --cask goland

# Or download DMG from website
# Open DMG → Drag GoLand to Applications
```

**Linux**:
```bash
# Download tar.gz
wget https://download.jetbrains.com/go/goland-*.tar.gz

# Extract
tar -xzf goland-*.tar.gz -C ~/Applications/

# Run
~/Applications/GoLand-*/bin/goland.sh

# Optional: Create desktop entry
~/Applications/GoLand-*/bin/goland.sh --install-desktop-entry
```

**Windows**:
- Download installer from website
- Run installer
- Follow wizard

### IntelliJ IDEA Ultimate

**Download**: https://www.jetbrains.com/idea/download/

**macOS**:
```bash
# Using Homebrew
brew install --cask intellij-idea

# Or download DMG from website
```

**Linux**:
```bash
# Download tar.gz
wget https://download.jetbrains.com/idea/ideaIU-*.tar.gz

# Extract
tar -xzf ideaIU-*.tar.gz -C ~/Applications/

# Run
~/Applications/idea-IU-*/bin/idea.sh

# Optional: Create desktop entry
~/Applications/idea-IU-*/bin/idea.sh --install-desktop-entry
```

**Windows**:
- Download installer
- Run installer
- Follow wizard

---

## Go Development

### Configure Go SDK

**GoLand (Automatic)**:
1. Open project: **File** → **Open** → Select `/path/to/revenge`
2. GoLand detects Go SDK automatically
3. Verify: **Settings** → **Go** → **GOROOT**
   - Should show: `/usr/local/go` (Linux/macOS) or `C:\Go` (Windows)

**IntelliJ IDEA (Manual)**:
1. Install Go plugin:
   - **Settings** → **Plugins**
   - Search **"Go"**
   - Install official Go plugin
   - Restart IDE

2. Configure GOROOT:
   - **Settings** → **Languages & Frameworks** → **Go** → **GOROOT**
   - Click **+** → **Local**
   - Select Go installation directory:
     - macOS/Linux: `/usr/local/go`
     - Windows: `C:\Go`
   - Click **OK**

3. Verify GOPATH (should be auto-detected):
   - **Settings** → **Go** → **GOPATH**
   - Should show: `~/go` or auto-detected location

### Configure Go Modules

**Settings** → **Go** → **Go Modules**:
- ✅ **Enable Go modules integration**
- ✅ **Vendoring mode**: Off (unless using vendor/)
- **Environment**: `GOEXPERIMENT=greenteagc,jsonv2`

### Configure gopls

**Settings** → **Languages & Frameworks** → **Go** → **Build Tags & Vendoring**:
- **Build tags**: (leave empty unless specific tags needed)
- **OS**: `linux` (or your OS)
- **Arch**: `amd64` (or your arch)

**gopls settings** (optional, defaults are good):
- **Settings** → **Go** → **Inspections**
- Ensure gopls inspections are enabled

### Install Go Tools

```bash
# Install goimports (if not auto-installed by IDE)
go install golang.org/x/tools/cmd/goimports@latest

# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install delve debugger (if not auto-installed)
go install github.com/go-delve/delve/cmd/dlv@latest
```

IDE will prompt to install missing tools automatically.

---

## Python Development

### IntelliJ IDEA Python Plugin

1. **Settings** → **Plugins**
2. Search **"Python"**
3. Install **Python** plugin (official)
4. Restart IDE

### Configure Python Interpreter

1. **Settings** → **Project: revenge** → **Python Interpreter**
2. Click gear icon → **Add Interpreter** → **Add Local Interpreter**
3. Select **System Interpreter**
4. Choose Python (version from SOURCE_OF_TRUTH):
   - macOS/Linux: `/usr/bin/python3` or `/usr/local/bin/python3`
   - Windows: `C:\Python3xx\python.exe` (xx = version from SOURCE_OF_TRUTH)
5. Click **OK**

### Install Python Tools

```bash
# Install ruff for linting and formatting
pip install ruff

# Install pytest for testing
pip install pytest

# Verify installation
ruff --version
pytest --version
```

### Configure ruff

**Settings** → **Tools** → **External Tools**:

1. Click **+** to add new tool
2. **Name**: `ruff format`
3. **Program**: `ruff`
4. **Arguments**: `format $FilePath$`
5. **Working directory**: `$ProjectFileDir$`
6. Click **OK**

Repeat for `ruff check`:
- **Name**: `ruff check`
- **Arguments**: `check $FilePath$`

---

## Frontend Development

### Node.js and npm

Verify Node.js installation:
```bash
node --version  # Should be v20+
npm --version
```

### Install Frontend Dependencies

```bash
cd web
npm install
```

### TypeScript Plugin

**IntelliJ IDEA Ultimate** includes TypeScript support built-in.

**Verify**:
1. **Settings** → **Languages & Frameworks** → **TypeScript**
2. Ensure TypeScript service is running
3. **TypeScript version**: Use bundled or project version (from `web/node_modules`)

### Svelte Plugin

1. **Settings** → **Plugins**
2. Search **"Svelte"**
3. Install **Svelte** plugin
4. Restart IDE

### Configure Prettier

**Settings** → **Languages & Frameworks** → **JavaScript** → **Prettier**:
- ✅ **Automatic Prettier configuration**
- **Prettier package**: `{project}/web/node_modules/prettier`
- ✅ **On save** (optional, recommended)
- **Run for files**: `{**/*,*}.{js,ts,svelte,json,css,scss}`

---

## Database Tools

### Built-in Database Tools

GoLand and IntelliJ IDEA Ultimate include powerful database tools.

### Connect to PostgreSQL

1. **Database** tool window (right sidebar)
   - If not visible: **View** → **Tool Windows** → **Database**

2. Click **+** → **Data Source** → **PostgreSQL**

3. Configure connection:
   - **Host**: `localhost`
   - **Port**: `5432`
   - **Database**: `revenge`
   - **User**: `revenge`
   - **Password**: `revenge` (dev environment)
   - **URL**: `jdbc:postgresql://localhost:5432/revenge`

4. Click **Test Connection**
   - If driver missing, IDE will prompt to download
   - Click **Download** to install PostgreSQL JDBC driver

5. Click **OK**

### Query Console

1. Right-click database → **New** → **Query Console**
2. Write SQL:
   ```sql
   SELECT * FROM users LIMIT 10;
   ```
3. Execute: **Ctrl+Enter** (macOS: Cmd+Enter)
4. Results appear in bottom panel

### Schema Navigation

- **Tables** → Browse tables
- Right-click table → **Modify Table** → Edit schema
- **Diagrams** → Visualize relationships
- **Export Data** → Export to CSV, JSON, SQL

---

## Essential Plugins

### Recommended Plugins

**Settings** → **Plugins** → **Marketplace**:

1. **.env files support** - Environment variable files syntax highlighting
2. **GitToolBox** - Enhanced Git integration with inline blame
3. **Rainbow Brackets** - Colorize matching brackets
4. **Key Promoter X** - Learn keyboard shortcuts
5. **Atom Material Icons** - Better file icons
6. **Makefile Language** - Makefile support (if using Makefiles)

**For IntelliJ IDEA only** (GoLand has these built-in):
- **Go** - Official Go language support
- **Docker** - Docker integration

**Optional**:
- **String Manipulation** - Advanced text operations
- **Grep Console** - Colorize console output
- **CSV Editor** - Edit CSV files with table view

### Install Plugin

1. **Settings** → **Plugins** → **Marketplace**
2. Search plugin name
3. Click **Install**
4. Restart IDE if prompted

---

## Project Configuration

### Open Revenge Project

1. **File** → **Open**
2. Navigate to `/path/to/revenge`
3. Click **OK**
4. IDE will index project (may take 1-2 minutes first time)

### Mark Directories

**Right-click directories** → **Mark Directory as**:

- `cmd/revenge` → **Sources Root** (or keep as default)
- `internal/` → **Sources Root** (or keep as default)
- `tests/` → **Test Sources Root**
- `web/node_modules` → **Excluded** (if not auto-excluded)
- `bin/` → **Excluded**
- `vendor/` → **Excluded** (if using vendoring)

### Project Structure

**File** → **Project Structure** → **Project**:
- **Project SDK**: Go (see SOURCE_OF_TRUTH)
- **Project language level**: (default)

**Modules**:
- Should auto-detect `revenge` module

**Libraries**:
- Should auto-detect Go libraries from `go.mod`

---

## Run Configurations

### Create "Revenge Server" Configuration

1. **Run** → **Edit Configurations**
2. Click **+** → **Go Build**
3. Configure:
   - **Name**: `Revenge Server`
   - **Run kind**: `Package`
   - **Package path**: `github.com/lusoris/revenge/cmd/revenge`
   - **Output directory**: `bin/`
   - **Working directory**: `$ProjectFileDir$`
   - **Environment**:
     ```
     GOEXPERIMENT=greenteagc,jsonv2
     ```
   - **Go tool arguments** (optional):
     ```
     -v
     ```
   - **Program arguments** (optional):
     ```
     --config config.dev.yaml
     ```
4. Click **OK**

### Run Server

- Click **▶ Run** icon in toolbar
- Or **Ctrl+R** (macOS) / **Shift+F10** (Windows/Linux)
- Server starts at `http://localhost:8096`

### Create "Run Tests" Configuration

1. **Run** → **Edit Configurations**
2. Click **+** → **Go Test**
3. Configure:
   - **Name**: `All Tests`
   - **Test kind**: `Directory`
   - **Directory**: `$ProjectFileDir$/internal`
   - **Pattern**: `.*`
   - **Working directory**: `$ProjectFileDir$`
   - **Environment**: `GOEXPERIMENT=greenteagc,jsonv2`
4. Click **OK**

### Run Tests

- Click **▶ Run 'All Tests'**
- Or right-click `internal/` → **Run 'go test internal/...'**

---

## Code Style Settings

### EditorConfig

The project includes `.editorconfig` which IDE respects automatically.

**Verify**: **Settings** → **Editor** → **Code Style**
- ✅ **Enable EditorConfig support**

### Go Code Style

**Settings** → **Editor** → **Code Style** → **Go**:

**Tabs and Indents**:
- ✅ **Use tab character**
- **Tab size**: `4`
- **Indent**: `4`
- **Continuation indent**: `4`

**Imports**:
- ✅ **Optimize imports on the fly**
- ✅ **Group stdlib imports**
- **Import order**:
  1. Standard library
  2. Third-party packages
  3. Current project

**Other**:
- ✅ **Add parentheses to receiver in method declarations**
- Format with **goimports** (automatic)

### Python Code Style

**Settings** → **Editor** → **Code Style** → **Python**:

**Tabs and Indents**:
- ☐ **Use tab character** (use spaces)
- **Tab size**: `4`
- **Indent**: `4`

**Imports**:
- ✅ **Optimize imports on the fly**
- ✅ **Sort imports**

### TypeScript/JavaScript Code Style

**Settings** → **Editor** → **Code Style** → **TypeScript**:

**Tabs and Indents**:
- ☐ **Use tab character** (use spaces)
- **Tab size**: `2`
- **Indent**: `2`

**Prettier Integration**:
- ✅ Use Prettier for formatting (configured above)

---

## Version Control

### Git Integration

**Settings** → **Version Control** → **Git**:
- **Path to Git executable**: `/usr/bin/git` (auto-detected)
- ✅ **Use credential helper**

### Configure Git in IDE

**Settings** → **Version Control** → **Commit**:
- ✅ **Analyze code** before commit
- ✅ **Check TODO** (Reformat code, optimize imports)
- ✅ **Perform code analysis**

### Commit Changes

1. **Cmd/Ctrl+K** to open Commit dialog
2. Select files to commit
3. Write commit message (follows conventional commit format)
4. ✅ **Reformat code**
5. ✅ **Optimize imports**
6. Click **Commit** or **Commit and Push**

### Useful Git Shortcuts

- **Cmd/Ctrl+K** - Commit
- **Cmd/Ctrl+Shift+K** - Push
- **Cmd/Ctrl+T** - Update Project (pull)
- **Alt+9** - Show Version Control panel
- **Cmd/Ctrl+Alt+Z** - Revert changes

---

## Format on Save

### Enable Auto-Format

**Settings** → **Tools** → **Actions on Save**:
- ✅ **Reformat code**
  - **File types**: Go, Python, JavaScript, TypeScript, Svelte
- ✅ **Optimize imports**
- ✅ **Run code cleanup** (optional)

Now files auto-format when you save (**Cmd/Ctrl+S**).

---

## Keyboard Shortcuts

### Essential Shortcuts

| Action | macOS | Windows/Linux |
|--------|-------|---------------|
| Search Everywhere | `Shift Shift` | `Shift Shift` |
| Find Action | `Cmd+Shift+A` | `Ctrl+Shift+A` |
| Go to File | `Cmd+Shift+O` | `Ctrl+Shift+N` |
| Go to Symbol | `Cmd+Alt+O` | `Ctrl+Alt+Shift+N` |
| Go to Declaration | `Cmd+B` | `Ctrl+B` |
| Go to Implementation | `Cmd+Alt+B` | `Ctrl+Alt+B` |
| Find Usages | `Cmd+F7` | `Alt+F7` |
| Rename | `Shift+F6` | `Shift+F6` |
| Extract Method | `Cmd+Alt+M` | `Ctrl+Alt+M` |
| Quick Fix | `Alt+Enter` | `Alt+Enter` |
| Format Code | `Cmd+Alt+L` | `Ctrl+Alt+L` |
| Optimize Imports | `Ctrl+Alt+O` | `Ctrl+Alt+O` |
| Run | `Ctrl+R` | `Shift+F10` |
| Debug | `Ctrl+D` | `Shift+F9` |
| Commit | `Cmd+K` | `Ctrl+K` |
| Push | `Cmd+Shift+K` | `Ctrl+Shift+K` |

### Learn More Shortcuts

**Help** → **Keyboard Shortcuts PDF**

Or install **Key Promoter X** plugin to learn shortcuts as you work.

---

## Performance Tuning

### Increase Heap Size

**Help** → **Change Memory Settings**:
- Default: 2048 MB
- Recommended for Revenge: **4096 MB** (4 GB)
- Click **Save and Restart**

### Exclude Unnecessary Folders

**Settings** → **Project Structure** → **Modules**:
- Mark as **Excluded**:
  - `node_modules/`
  - `bin/`
  - `vendor/` (if using vendoring)
  - `.git/`

### Power Save Mode

When not actively coding:
- **File** → **Power Save Mode**
- Disables background inspections (saves CPU)
- Re-enable when coding

---

## Verification

### Verify Go Setup

1. Open `cmd/revenge/main.go`
2. Check:
   - ✅ No red underlines (imports resolved)
   - ✅ Code completion works (**Ctrl+Space**)
   - ✅ Go to Definition works (**Cmd/Ctrl+B**)
   - ✅ Quick documentation works (**F1** or **Cmd+Q**)

3. Run tests:
   - Right-click `internal/` → **Run 'go test internal/...'**
   - Tests should pass

4. Run server:
   - **Run** → **Run 'Revenge Server'**
   - Server should start without errors

### Verify Frontend Setup

1. Open `web/src/App.svelte`
2. Check:
   - ✅ Svelte syntax highlighting
   - ✅ TypeScript completion works
   - ✅ Imports resolved

3. Run frontend:
   ```bash
   # In IDE terminal (Alt+F12)
   cd web
   npm run dev
   ```
   - Should start at `http://localhost:5173`

---

## Next Steps

1. **Configure Run Configurations** - See [Run Configurations](#run-configurations)
2. **Learn Keyboard Shortcuts** - See [Keyboard Shortcuts](#keyboard-shortcuts)
3. **Set Up Remote Development** - See [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md)
4. **Explore Database Tools** - See [Database Tools](#database-tools)

---

## Related Documentation

- [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md) - JetBrains Gateway + Coder
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues
- [../../.coder/docs/JETBRAINS_INTEGRATION.md](../../.coder/docs/JETBRAINS_INTEGRATION.md) - Complete Gateway guide
- [../../.shared/docs/TOOL_COMPARISON.md](../../.shared/docs/TOOL_COMPARISON.md) - Compare IDEs

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
