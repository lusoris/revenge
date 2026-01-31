# JetBrains IDE Troubleshooting

**Purpose**: Solutions for common issues with GoLand and IntelliJ IDEA

**Last Updated**: 2026-01-31

---

## Table of Contents

- [Go SDK Issues](#go-sdk-issues)
- [Gateway Connection Issues](#gateway-connection-issues)
- [Performance Issues](#performance-issues)
- [IDE Backend Issues](#ide-backend-issues)
- [gopls Issues](#gopls-issues)
- [Plugin Issues](#plugin-issues)
- [Database Tools Issues](#database-tools-issues)
- [Frontend Development Issues](#frontend-development-issues)
- [Indexing Issues](#indexing-issues)
- [License Issues](#license-issues)

---

## Go SDK Issues

### Go SDK Not Detected

**Problem**: "Go SDK not configured" or "GOROOT not found"

**Symptoms**:
- Red underlines on all Go imports
- No code completion
- "Cannot resolve symbol" errors

**Solutions**:

1. **Verify Go Installation**:
   ```bash
   go version
   # Should show required Go version (see SOURCE_OF_TRUTH) or higher

   which go
   # Should show /usr/local/go/bin/go (macOS/Linux)
   ```

2. **Configure GOROOT in IDE**:
   - **Settings** → **Go** → **GOROOT**
   - Click **+** → **Local**
   - Select Go installation:
     - macOS/Linux: `/usr/local/go`
     - Windows: `C:\Go`
   - Click **OK**

3. **Restart IDE**:
   - **File** → **Invalidate Caches** → **Invalidate and Restart**

4. **Re-index Project**:
   - **File** → **Invalidate Caches** → **Just Restart**

### Wrong Go Version Detected

**Problem**: IDE uses wrong Go version

**Symptoms**:
- Code runs fine in terminal but IDE shows errors
- "Go version mismatch" warnings

**Solutions**:

1. **Check PATH**:
   ```bash
   echo $PATH
   # Ensure /usr/local/go/bin comes before other Go installations
   ```

2. **Update GOROOT**:
   - **Settings** → **Go** → **GOROOT**
   - Remove old SDK
   - Add correct SDK (`/usr/local/go`)

3. **Verify in IDE Terminal**:
   - **Alt+F12** to open terminal
   - `go version` should match expected version

### GOPATH Issues

**Problem**: "GOPATH not set" or wrong GOPATH

**Solutions**:

1. **Auto-detect GOPATH**:
   - **Settings** → **Go** → **GOPATH**
   - Click **Auto-detect**
   - Should find `~/go` or your custom GOPATH

2. **Manually set GOPATH**:
   - **Settings** → **Go** → **GOPATH**
   - **Project GOPATH**: Leave empty (use global)
   - **Global GOPATH**: `~/go` (or your preference)

3. **For Go Modules (Revenge uses modules)**:
   - GOPATH is less important
   - Ensure **Go Modules** is enabled:
     - **Settings** → **Go** → **Go Modules**
     - ✅ **Enable Go modules integration**

---

## Gateway Connection Issues

### Gateway Can't Connect to Coder

**Problem**: "Connection failed" when trying to connect via Gateway

**Symptoms**:
- "Cannot connect to Coder" error
- "Workspace not found" error
- Gateway hangs on "Connecting..."

**Solutions**:

1. **Verify Coder CLI Login**:
   ```bash
   coder login https://coder.ancilla.lol
   coder list
   # Should show your workspaces
   ```

2. **Verify Workspace is Running**:
   ```bash
   coder list
   # Status should be "Running", not "Stopped"

   # If stopped, start it
   coder start revenge-dev
   ```

3. **Re-login in Gateway**:
   - Open Gateway
   - **Settings** → **Coder**
   - Click **Logout**
   - Click **Login** → Enter `https://coder.ancilla.lol`
   - Authenticate in browser

4. **Check Network**:
   ```bash
   ping coder.ancilla.lol
   # Should respond with low latency (<200ms)
   ```

5. **Restart Gateway**:
   - Quit Gateway completely
   - Reopen and try connecting again

### Gateway Connection Drops

**Problem**: Connection to workspace drops unexpectedly

**Solutions**:

1. **Check Network Stability**:
   ```bash
   ping -c 100 coder.ancilla.lol
   # Look for packet loss
   ```

2. **Increase Timeout** (if on slow network):
   - **Settings** → **Advanced** → **Connection timeout**
   - Increase to 60 seconds

3. **Restart Workspace**:
   ```bash
   coder restart revenge-dev
   ```

4. **Use Wired Connection**:
   - WiFi can be unstable for remote development
   - Use Ethernet if possible

---

## Performance Issues

### IDE is Slow/Laggy

**Problem**: Slow typing, laggy completions, high CPU usage

**Symptoms**:
- Delayed response when typing
- Completions take seconds to appear
- Fans running loud
- IDE feels unresponsive

**Solutions**:

1. **Increase Heap Size**:
   - **Help** → **Change Memory Settings**
   - Increase to **4096 MB** (4 GB)
   - Click **Save and Restart**

2. **Enable Power Save Mode** (temporarily):
   - **File** → **Power Save Mode**
   - Disables background inspections
   - Use when not actively coding

3. **Disable Unused Plugins**:
   - **Settings** → **Plugins**
   - Disable plugins you don't use
   - Restart IDE

4. **Exclude Unnecessary Directories**:
   - Right-click `node_modules/` → **Mark Directory as** → **Excluded**
   - Right-click `bin/` → **Mark Directory as** → **Excluded**
   - Right-click `vendor/` → **Mark Directory as** → **Excluded**

5. **Clear Caches**:
   - **File** → **Invalidate Caches** → **Invalidate and Restart**

6. **Reduce Code Analysis Scope**:
   - **Settings** → **Editor** → **Inspections**
   - Uncheck less important inspections
   - Keep critical ones (errors, warnings)

### High Memory Usage

**Problem**: IDE uses too much RAM, system slows down

**Solutions**:

1. **Check Current Heap Usage**:
   - **Help** → **Diagnostic Tools** → **Activity Monitor**
   - See heap usage

2. **Adjust Heap Size**:
   - **Help** → **Change Memory Settings**
   - If usage consistently near max: Increase to 4GB or 6GB
   - If usage low: Can reduce to save RAM

3. **Close Unused Projects**:
   - **File** → **Close Project**
   - Only keep one project open at a time

4. **Restart IDE Regularly**:
   - Memory leaks can accumulate
   - Restart IDE daily if doing heavy work

### Slow Indexing

**Problem**: Indexing takes forever or never completes

**Symptoms**:
- "Indexing..." status bar message persists
- Code completion doesn't work during indexing
- Project unresponsive

**Solutions**:

1. **Wait for Initial Index** (first time):
   - First index of Revenge project takes 2-5 minutes
   - Check progress: Bottom right corner

2. **Exclude Large Directories**:
   - `node_modules/` should be auto-excluded
   - If not: Right-click → **Mark Directory as** → **Excluded**

3. **Pause Indexing** (if needed urgently):
   - Click **Pause** in indexing status
   - Resume when ready

4. **Invalidate and Restart**:
   - **File** → **Invalidate Caches** → **Invalidate and Restart**

5. **Check Disk Space**:
   - IDE needs disk space for caches
   - Ensure at least 5GB free

---

## IDE Backend Issues

### IDE Backend Won't Start (Gateway)

**Problem**: "Failed to start IDE backend" when using Gateway

**Solutions**:

1. **Delete Backend Cache**:
   ```bash
   coder ssh revenge-dev
   rm -rf ~/.cache/JetBrains/
   exit
   ```

2. **Remove Connection and Re-add**:
   - Open Gateway
   - Remove connection from recent list
   - Add new connection
   - Gateway re-downloads backend

3. **Check Workspace Resources**:
   ```bash
   coder ssh revenge-dev
   df -h  # Check disk space (need 2GB+ free)
   free -h  # Check RAM
   ```

4. **Restart Workspace**:
   ```bash
   coder restart revenge-dev
   ```

### IDE Backend Crashes

**Problem**: Backend crashes during use (Gateway)

**Solutions**:

1. **Check Logs**:
   ```bash
   coder ssh revenge-dev
   tail -f ~/.cache/JetBrains/*/log/idea.log
   ```

2. **Increase Backend Memory**:
   - In workspace, edit backend memory settings
   - Files in `~/.cache/JetBrains/*/` have memory configs

3. **Restart Backend**:
   - Close IDE on local machine
   - Reopen connection in Gateway

---

## gopls Issues

### gopls Not Starting

**Problem**: "gopls is not running" or no Go code intelligence

**Symptoms**:
- No code completion for Go
- No go to definition
- No error highlighting

**Solutions**:

1. **Verify gopls Installed**:
   ```bash
   gopls version
   # Should show version

   # If not found, install
   go install golang.org/x/tools/gopls@latest
   ```

2. **Restart gopls in IDE**:
   - **Settings** → **Go** → **gopls**
   - Click **Restart gopls**

3. **Check gopls Logs**:
   - **Help** → **Show Log in Finder/Explorer**
   - Look for `gopls.log` or errors mentioning gopls

4. **Reinstall gopls**:
   ```bash
   go install golang.org/x/tools/gopls@latest
   ```

5. **Restart IDE**:
   - **File** → **Exit** (fully quit)
   - Reopen IDE

### gopls High CPU Usage

**Problem**: gopls process uses 100% CPU

**Solutions**:

1. **Wait for Initial Analysis**:
   - First time gopls runs, it analyzes entire project
   - Can take 2-3 minutes for Revenge
   - CPU usage normalizes after

2. **Restart gopls**:
   - **Settings** → **Go** → **gopls**
   - Click **Restart gopls**

3. **Exclude Directories from Analysis**:
   - Mark `vendor/`, `bin/`, `node_modules/` as Excluded

4. **Disable Unused gopls Features**:
   - **Settings** → **Go** → **gopls**
   - Disable features you don't need

---

## Plugin Issues

### Plugin Won't Install

**Problem**: Plugin installation fails or hangs

**Solutions**:

1. **Check Internet Connection**:
   - Plugins download from JetBrains servers
   - Verify connection is stable

2. **Clear Plugin Cache**:
   ```bash
   # Close IDE first
   rm -rf ~/.cache/JetBrains/*/plugins/
   # Reopen IDE
   ```

3. **Install Manually**:
   - Download plugin JAR from marketplace
   - **Settings** → **Plugins** → **⚙️** → **Install Plugin from Disk**
   - Select downloaded JAR

4. **Check Compatibility**:
   - Ensure plugin supports your IDE version
   - Some plugins only work with specific versions

### Plugin Causes Errors

**Problem**: IDE errors or crashes after installing plugin

**Solutions**:

1. **Disable Plugin**:
   - Start IDE in Safe Mode:
     ```bash
     # macOS/Linux
     /Applications/GoLand.app/Contents/MacOS/goland --safe-mode
     ```
   - **Settings** → **Plugins**
   - Disable problematic plugin
   - Restart normally

2. **Uninstall Plugin**:
   - **Settings** → **Plugins**
   - Find plugin → **Uninstall**
   - Restart IDE

3. **Check Plugin Updates**:
   - **Settings** → **Plugins** → **Installed**
   - Click **Update All**

---

## Database Tools Issues

### Can't Connect to PostgreSQL

**Problem**: Database connection fails

**Symptoms**:
- "Connection refused" error
- "Authentication failed" error
- "Driver not found" error

**Solutions**:

1. **Download JDBC Driver**:
   - When clicking **Test Connection**
   - IDE prompts to download PostgreSQL driver
   - Click **Download**

2. **Verify PostgreSQL is Running**:
   ```bash
   # Local
   docker-compose -f docker-compose.dev.yml ps
   # Should show PostgreSQL running

   # In Coder workspace
   coder ssh revenge-dev
   pg_isready -h localhost -p 5432
   ```

3. **Check Connection Details**:
   - **Host**: `localhost` (in local/workspace context)
   - **Port**: `5432`
   - **Database**: `revenge`
   - **User**: `revenge`
   - **Password**: `revenge` (dev environment)

4. **Test Connection from Terminal**:
   ```bash
   psql -h localhost -p 5432 -U revenge -d revenge
   # Should connect successfully
   ```

5. **Allow IDE to Access Localhost**:
   - On macOS: System Preferences → Security & Privacy → Firewall
   - Allow JetBrains IDE

### Query Console Not Working

**Problem**: Can't execute queries in query console

**Solutions**:

1. **Verify Connection is Active**:
   - Database should show green dot (connected)
   - If red, reconnect: Right-click → **Refresh**

2. **Check SQL Dialect**:
   - In query console, bottom right
   - Should be **PostgreSQL**

3. **Use Correct Shortcut**:
   - **Cmd+Enter** (macOS) or **Ctrl+Enter** (Windows/Linux)
   - Executes query at cursor

4. **Check for Syntax Errors**:
   - IDE highlights SQL errors
   - Fix syntax before executing

---

## Frontend Development Issues

### TypeScript LSP Not Working

**Problem**: No TypeScript completion or errors in `.ts` or `.svelte` files

**Solutions**:

1. **Verify TypeScript Service**:
   - **Settings** → **Languages & Frameworks** → **TypeScript**
   - Ensure service is running
   - Click **Restart service** if needed

2. **Install Dependencies**:
   ```bash
   cd web
   npm install
   ```

3. **Check tsconfig.json**:
   - Ensure `web/tsconfig.json` exists and is valid
   - IDE reads this for TypeScript configuration

4. **Restart IDE**:
   - **File** → **Invalidate Caches** → **Invalidate and Restart**

### Svelte Syntax Not Highlighted

**Problem**: `.svelte` files show as plain text

**Solutions**:

1. **Install Svelte Plugin**:
   - **Settings** → **Plugins** → **Marketplace**
   - Search **"Svelte"**
   - Install official Svelte plugin
   - Restart IDE

2. **Check File Association**:
   - **Settings** → **Editor** → **File Types**
   - Ensure `.svelte` is associated with **Svelte**

3. **Reinstall Plugin**:
   - **Settings** → **Plugins**
   - Uninstall Svelte plugin
   - Restart IDE
   - Reinstall Svelte plugin

### Prettier Not Formatting

**Problem**: Prettier doesn't format files on save or manually

**Solutions**:

1. **Verify Prettier is Installed**:
   ```bash
   cd web
   npm list prettier
   # Should show prettier in dependencies
   ```

2. **Configure Prettier in IDE**:
   - **Settings** → **Languages & Frameworks** → **JavaScript** → **Prettier**
   - **Prettier package**: `{project}/web/node_modules/prettier`
   - ✅ **On save** (if you want auto-format)

3. **Check File Pattern**:
   - **Settings** → **Prettier**
   - **Run for files**: `{**/*,*}.{js,ts,svelte,json,css}`

4. **Manual Format**:
   - **Cmd/Ctrl+Alt+Shift+P** to format with Prettier
   - Or right-click → **Reformat with Prettier**

---

## Indexing Issues

### Indexing Stuck at 0%

**Problem**: Indexing shows 0% and never progresses

**Solutions**:

1. **Wait 5 Minutes**:
   - Sometimes indexing appears stuck but is working
   - Give it time, especially on large projects

2. **Check Logs**:
   - **Help** → **Show Log in Finder/Explorer**
   - Look for errors

3. **Invalidate Caches**:
   - **File** → **Invalidate Caches** → **Invalidate and Restart**

4. **Delete System Caches Manually**:
   ```bash
   # Close IDE first
   rm -rf ~/.cache/JetBrains/*/system/
   # Reopen IDE
   ```

### Indexing Repeats on Every Startup

**Problem**: IDE re-indexes entire project every time it starts

**Solutions**:

1. **Check .idea/ Folder Permissions**:
   ```bash
   ls -la .idea/
   # Should be writable

   # If not, fix permissions
   chmod -R u+w .idea/
   ```

2. **Exclude .idea/ from Antivirus**:
   - Antivirus can interfere with IDE caches
   - Add `.idea/` and `~/.cache/JetBrains/` to exclusions

3. **Check Disk Space**:
   ```bash
   df -h
   # Ensure 5GB+ free
   ```

4. **Disable "Always update snapshots"**:
   - **Settings** → **Build, Execution, Deployment** → **Build Tools**
   - Uncheck options that force re-indexing

---

## License Issues

### License Not Recognized

**Problem**: "License not found" or "License expired"

**Solutions**:

1. **Activate License**:
   - **Help** → **Register**
   - Choose activation method:
     - **JetBrains Account** (recommended)
     - **Activation code**
     - **License server**

2. **Verify License Server** (if using server):
   - **Help** → **Register** → **License server**
   - Enter server URL
   - Click **Activate**

3. **Check Subscription Status**:
   - Visit https://account.jetbrains.com/
   - Verify subscription is active
   - Check which products are included

4. **Re-activate**:
   - **Help** → **Register**
   - Log out and log back in with JetBrains Account

### Trial Expired

**Problem**: Trial period ended

**Solutions**:

1. **Purchase License**:
   - GoLand: $99/year (individual)
   - IntelliJ IDEA Ultimate: $169/year (individual)
   - https://www.jetbrains.com/store/

2. **Apply for Free License** (if eligible):
   - **Students**: https://www.jetbrains.com/community/education/
   - **Open Source**: https://www.jetbrains.com/community/opensource/
   - **Teachers**: https://www.jetbrains.com/community/education/

3. **Switch to VS Code or Zed** (free alternatives):
   - See [../../.vscode/docs/INDEX.md](../../.vscode/docs/INDEX.md)
   - See [../../.zed/docs/INDEX.md](../../.zed/docs/INDEX.md)

---

## General Troubleshooting Steps

### When All Else Fails

1. **Restart IDE**:
   - **File** → **Exit**
   - Wait 10 seconds
   - Reopen IDE

2. **Invalidate Caches**:
   - **File** → **Invalidate Caches** → **Invalidate and Restart**

3. **Delete IDE Cache Manually**:
   ```bash
   # Close IDE first
   rm -rf ~/.cache/JetBrains/GoLand*/
   # Or for IntelliJ IDEA
   rm -rf ~/.cache/JetBrains/IntelliJIdea*/
   ```

4. **Reset IDE Settings**:
   ```bash
   # Backup first!
   mv ~/.config/JetBrains/ ~/.config/JetBrains.backup/
   # Reopen IDE with fresh settings
   ```

5. **Reinstall IDE**:
   - Uninstall completely
   - Delete all caches and configs
   - Download fresh copy
   - Install

### Collecting Logs for Support

If reporting an issue:

1. **Enable Debug Logging**:
   - **Help** → **Diagnostic Tools** → **Debug Log Settings**
   - Add relevant categories

2. **Collect Logs**:
   - **Help** → **Collect Logs and Diagnostic Data**
   - Saves ZIP with all logs

3. **Attach to Issue**:
   - https://youtrack.jetbrains.com/issues
   - Create new issue
   - Attach logs ZIP

---

## Getting Help

### Official Resources

- **JetBrains Support**: https://www.jetbrains.com/support/
- **GoLand Docs**: https://www.jetbrains.com/help/go/
- **IntelliJ IDEA Docs**: https://www.jetbrains.com/help/idea/
- **Issue Tracker**: https://youtrack.jetbrains.com/issues

### Community

- **Reddit**: r/Jetbrains, r/golang (for Go-specific)
- **Stack Overflow**: Tag `jetbrains` or `goland` or `intellij-idea`

### Project Resources

- **Revenge Issues**: https://github.com/kilianso/revenge/issues
- **Shared Troubleshooting**: [../../.shared/docs/TROUBLESHOOTING.md](../../.shared/docs/TROUBLESHOOTING.md)

---

## Related Documentation

- [SETUP.md](SETUP.md) - Initial setup and configuration
- [REMOTE_DEVELOPMENT.md](REMOTE_DEVELOPMENT.md) - Gateway + Coder
- [INDEX.md](INDEX.md) - Documentation hub
- [../../.coder/docs/TROUBLESHOOTING.md](../../.coder/docs/TROUBLESHOOTING.md) - Coder-specific issues

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
