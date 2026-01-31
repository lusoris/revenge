# Cross-Tool Troubleshooting Guide

**Purpose**: Common issues that span multiple development tools

**Last Updated**: 2026-01-31

---

## Quick Diagnosis

Having an issue? Find your symptom below:

| Symptom | Likely Cause | Quick Fix |
|---------|--------------|-----------|
| LSP not working in any IDE | LSP not installed or not in PATH | Install gopls/ruff-lsp, check PATH |
| Formatter works in VS Code but not Zed | Different formatter binary or config | Check both use same tool |
| Git hooks don't run | Hooks not installed | `pre-commit install` |
| Can't connect to Coder workspace | Workspace stopped or network issue | `coder start workspace-name` |
| Build fails locally but not in CI | Environment mismatch | Check Go version, GOEXPERIMENT |
| Tests pass locally but fail in CI | Missing dependencies in CI | Check GitHub Actions workflow |
| Hot reload not working | Air config issue or file watcher limit | Check `.air.toml`, increase inotify |
| Slow IDE performance | Too many extensions or large workspace | Disable unused extensions, exclude folders |

---

## Table of Contents

- [IDE Issues](#ide-issues)
- [Environment Synchronization](#environment-synchronization)
- [Remote vs Local Problems](#remote-vs-local-problems)
- [Git and Version Control](#git-and-version-control)
- [Build and Test Issues](#build-and-test-issues)
- [Integration Problems](#integration-problems)
- [Performance Issues](#performance-issues)
- [Network and Connectivity](#network-and-connectivity)

---

## IDE Issues

### LSP Not Starting

**Symptoms**:
- No code completion
- No go-to-definition
- No error highlighting

**Diagnosis**:
```bash
# Check if LSP is installed
gopls version      # For Go
ruff-lsp --version # For Python
which typescript-language-server  # For TypeScript

# Check if it's in PATH
echo $PATH

# Check IDE logs
# VS Code: View → Output → Select "Go" or "Ruff"
# Zed: View → Debug → Log File
```

**Solutions**:

1. **Install missing LSP**:
   ```bash
   # gopls
   go install golang.org/x/tools/gopls@latest

   # ruff-lsp
   pip install ruff-lsp

   # TypeScript
   npm install -g typescript-language-server
   ```

2. **Add to PATH**:
   ```bash
   # Add Go bin to PATH
   export PATH="$PATH:$(go env GOPATH)/bin"

   # Add Python bin to PATH
   export PATH="$PATH:$HOME/.local/bin"

   # Add to ~/.bashrc or ~/.zshrc
   ```

3. **Restart IDE**:
   - VS Code: Reload Window (Cmd/Ctrl+Shift+P → "Reload Window")
   - Zed: Restart Zed

4. **Check LSP configuration**:
   - VS Code: Check `.vscode/settings.json`
   - Zed: Check `.zed/settings.json`
   - Ensure LSP binary path is correct

---

### Formatter Not Working

**Symptoms**:
- Format on save doesn't work
- Manual format command does nothing
- Code remains unformatted

**Diagnosis**:
```bash
# Check if formatter is installed
goimports -h       # For Go
ruff --version     # For Python
prettier --version # For TypeScript/Svelte

# Try formatting manually
goimports -w file.go
ruff format file.py
prettier --write file.ts
```

**Solutions**:

1. **Install formatter**:
   ```bash
   go install golang.org/x/tools/cmd/goimports@latest
   pip install ruff
   npm install -g prettier
   ```

2. **Check format on save is enabled**:
   - VS Code: `.vscode/settings.json` → `"editor.formatOnSave": true`
   - Zed: `.zed/settings.json` → `"format_on_save": true`

3. **Check formatter is selected**:
   - VS Code: Right-click in file → "Format Document With..." → Select formatter
   - Zed: Check `.zed/settings.json` → `"formatter"` for language

4. **Check for conflicting formatters**:
   - Only one formatter should be enabled per language
   - Disable competing extensions (e.g., Black vs Ruff for Python)

---

### Debugger Not Working (VS Code)

**Symptoms**:
- Breakpoints not hit
- Debugger won't start
- "Could not attach to process" error

**Solutions**:

1. **Check launch.json configuration**:
   ```json
   {
     "version": "0.2.0",
     "configurations": [
       {
         "name": "Launch Revenge",
         "type": "go",
         "request": "launch",
         "mode": "debug",
         "program": "${workspaceFolder}/cmd/revenge"
       }
     ]
   }
   ```

2. **Install Delve**:
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   ```

3. **Check Go extension is installed**:
   - Install: `golang.go` extension

4. **Disable optimizations**:
   ```bash
   # Build without optimizations for debugging
   go build -gcflags="all=-N -l" ./cmd/revenge
   ```

---

### Extensions Conflict

**Symptoms**:
- Unexpected behavior
- Multiple formatters fighting
- IDE becomes slow

**Solutions**:

1. **Identify conflicting extensions**:
   - VS Code: Extensions → Filter @enabled
   - Look for duplicate functionality

2. **Disable unused extensions**:
   - Keep only one formatter per language
   - Keep only one LSP provider per language

3. **Workspace-specific extensions**:
   - VS Code: Use `.vscode/extensions.json` to recommend specific extensions
   - Disable others in workspace scope

---

## Environment Synchronization

### Different Behavior: Local vs CI

**Symptoms**:
- Tests pass locally but fail in CI
- Build succeeds locally but fails in CI
- Linter passes locally but fails in CI

**Diagnosis**:
```bash
# Check Go version
go version
# Compare with .github/workflows/ci.yml

# Check GOEXPERIMENT
echo $GOEXPERIMENT
# Should be: greenteagc,jsonv2

# Check dependencies
go mod verify
go mod tidy
git diff go.mod go.sum
```

**Solutions**:

1. **Match Go version**:
   - Local should match CI (1.25.6)
   - Update local Go if needed

2. **Set GOEXPERIMENT**:
   ```bash
   export GOEXPERIMENT=greenteagc,jsonv2
   # Add to ~/.bashrc or ~/.zshrc
   ```

3. **Clean and rebuild**:
   ```bash
   go clean -cache
   go mod tidy
   go build ./...
   ```

4. **Run CI checks locally**:
   ```bash
   # Run linters
   golangci-lint run

   # Run tests with same flags as CI
   go test -race -coverprofile=coverage.out ./...
   ```

---

### Settings Drift Between IDEs

**Symptoms**:
- Code formatted differently in VS Code vs Zed
- Different indentation
- Different line length

**Solutions**:

1. **Check .editorconfig is honored**:
   - Both IDEs should respect `.editorconfig`
   - Install EditorConfig extension if needed

2. **Synchronize formatter configs**:
   ```bash
   # Check all formatting configs match
   grep -r "line_length\|line-length" ruff.toml .vscode/ .zed/
   grep -r "tab\|indent" .editorconfig .vscode/ .zed/
   ```

3. **Use same formatter binary**:
   - Verify both IDEs use the same `ruff`, `goimports`, `prettier` binary
   - Check versions match: `ruff --version`, `goimports -h`, `prettier --version`

4. **Standardize on .editorconfig**:
   - Move all universal settings to `.editorconfig`
   - Remove duplicates from IDE configs

---

## Remote vs Local Problems

### Can't Connect to Coder Workspace

**Symptoms**:
- "Connection refused"
- "Workspace not found"
- SSH timeout

**Diagnosis**:
```bash
# Check workspace status
coder list

# Check workspace is running
coder status workspace-name

# Test SSH connection
coder ssh workspace-name -- echo "Connection OK"
```

**Solutions**:

1. **Start workspace if stopped**:
   ```bash
   coder start workspace-name
   ```

2. **Recreate workspace if broken**:
   ```bash
   coder stop workspace-name
   coder delete workspace-name
   coder create --template revenge workspace-name
   ```

3. **Check network connectivity**:
   ```bash
   # Ping Coder server
   ping coder.ancilla.lol

   # Check firewall isn't blocking
   ```

4. **Re-login to Coder**:
   ```bash
   coder logout
   coder login https://coder.ancilla.lol
   ```

---

### Slow Performance on Remote

**Symptoms**:
- Laggy typing
- Slow file operations
- Delayed code completion

**Solutions**:

1. **Check network latency**:
   ```bash
   # Measure latency to Coder server
   ping -c 10 coder.ancilla.lol
   # Should be < 50ms ideally
   ```

2. **Use Zed instead of VS Code**:
   - Zed is faster over SSH
   - Less network traffic

3. **Disable resource-intensive extensions**:
   - Disable unused extensions
   - Disable Git decorations if not needed
   - Reduce refresh frequency

4. **Increase workspace resources**:
   - Edit Coder template to allocate more CPU/RAM
   - Restart workspace

5. **Use code-server (VS Code browser)**:
   - All computation happens on server
   - No network latency for typing

---

### File Sync Issues

**Symptoms**:
- Changes not appearing in IDE
- File watcher not detecting changes
- Hot reload not working

**Solutions**:

1. **Increase inotify limit (Linux)**:
   ```bash
   # Temporary
   sudo sysctl fs.inotify.max_user_watches=524288

   # Permanent
   echo "fs.inotify.max_user_watches=524288" | sudo tee -a /etc/sysctl.conf
   sudo sysctl -p
   ```

2. **Restart file watcher**:
   ```bash
   # Stop Air
   pkill air

   # Restart
   air
   ```

3. **Check .air.toml excludes**:
   - Ensure excluded folders don't include your code
   - Ensure included extensions match your files

---

## Git and Version Control

### Pre-commit Hooks Don't Run

**Symptoms**:
- Commit succeeds without running hooks
- Linters not executed before commit
- Can commit unformatted code

**Diagnosis**:
```bash
# Check if hooks installed
ls -la .git/hooks/
# Should see pre-commit, commit-msg, pre-push

# Check pre-commit is installed
pre-commit --version

# Check Git hooks path
git config core.hooksPath
# Should be empty or .git/hooks
```

**Solutions**:

1. **Install pre-commit framework**:
   ```bash
   pip install pre-commit
   ```

2. **Install hooks**:
   ```bash
   pre-commit install
   pre-commit install --hook-type commit-msg
   pre-commit install --hook-type pre-push
   ```

3. **Test hooks**:
   ```bash
   pre-commit run --all-files
   ```

4. **If hooks still don't run**:
   ```bash
   # Remove and reinstall
   pre-commit uninstall
   rm -rf .git/hooks/pre-commit .git/hooks/commit-msg .git/hooks/pre-push
   pre-commit install --install-hooks
   ```

---

### Commit Message Rejected

**Symptoms**:
- "commit-msg hook failed"
- "Invalid commit message format"

**Solutions**:

1. **Use conventional commit format**:
   ```
   feat: add new feature
   fix: resolve bug
   docs: update documentation
   test: add tests
   chore: update dependencies
   ```

2. **Check format exactly**:
   - Type: `feat`, `fix`, `docs`, `test`, `chore`, `refactor`, `perf`, `ci`, `build`, `revert`
   - Optional scope: `feat(api): ...`
   - Colon and space: `: `
   - Description: lowercase, no period at end

3. **Bypass if necessary (emergency only)**:
   ```bash
   git commit --no-verify -m "message"
   ```

---

### Merge Conflicts in Generated Files

**Symptoms**:
- Conflicts in `go.sum`
- Conflicts in lock files

**Solutions**:

1. **For go.sum conflicts**:
   ```bash
   # Accept both changes
   git checkout --ours go.sum
   git checkout --theirs go.sum

   # Regenerate
   go mod tidy

   # Stage and continue
   git add go.sum
   git merge --continue
   ```

2. **For lock files (package-lock.json, etc.)**:
   ```bash
   # Delete lock file
   rm package-lock.json

   # Regenerate
   npm install

   # Stage and continue
   git add package-lock.json
   git merge --continue
   ```

---

## Build and Test Issues

### Build Fails with "Cannot Find Package"

**Symptoms**:
- `package ... is not in GOROOT`
- `no required module provides package`

**Solutions**:

1. **Download dependencies**:
   ```bash
   go mod download
   go mod tidy
   ```

2. **Verify go.mod**:
   ```bash
   go mod verify
   ```

3. **Clean cache and rebuild**:
   ```bash
   go clean -modcache
   go mod download
   ```

4. **Check GOPATH and GOROOT**:
   ```bash
   go env GOPATH
   go env GOROOT
   # Should be valid paths
   ```

---

### Tests Fail with "No Such File or Directory"

**Symptoms**:
- Tests looking for test fixtures fail
- File paths incorrect in tests

**Solutions**:

1. **Check working directory**:
   ```bash
   # Tests run from package directory
   # Use relative paths from test file location
   ```

2. **Use testdata folder**:
   ```
   pkg/
     mypackage/
       mypackage.go
       mypackage_test.go
       testdata/
         fixture.json
   ```

3. **Embed test data** (Go 1.16+):
   ```go
   import _ "embed"

   //go:embed testdata/fixture.json
   var fixtureData string
   ```

---

### Race Detector Fails

**Symptoms**:
- `go test -race` fails
- `WARNING: DATA RACE`

**Solutions**:

1. **Fix the race condition**:
   - Use mutexes: `sync.Mutex`, `sync.RWMutex`
   - Use channels for synchronization
   - Use atomic operations: `sync/atomic`

2. **Common race patterns**:
   ```go
   // BAD: Race condition
   var count int
   for i := 0; i < 10; i++ {
       go func() { count++ }()
   }

   // GOOD: Use mutex
   var count int
   var mu sync.Mutex
   for i := 0; i < 10; i++ {
       go func() {
           mu.Lock()
           count++
           mu.Unlock()
       }()
   }

   // BETTER: Use atomic
   var count int64
   for i := 0; i < 10; i++ {
       go func() {
           atomic.AddInt64(&count, 1)
       }()
   }
   ```

---

## Integration Problems

### Docker Compose Services Won't Start

**Symptoms**:
- PostgreSQL fails to start
- "port already in use"
- Container exits immediately

**Solutions**:

1. **Check port conflicts**:
   ```bash
   # Check if port in use
   lsof -i :5432  # PostgreSQL
   lsof -i :6379  # Dragonfly
   lsof -i :8108  # Typesense

   # Kill conflicting process or change port
   ```

2. **Check Docker status**:
   ```bash
   docker ps -a
   docker logs postgres
   docker logs dragonfly
   docker logs typesense
   ```

3. **Clean restart**:
   ```bash
   docker-compose -f docker-compose.dev.yml down -v
   docker-compose -f docker-compose.dev.yml up -d
   ```

4. **Check disk space**:
   ```bash
   df -h
   # Docker needs space for volumes
   ```

---

### Can't Connect to Database

**Symptoms**:
- "connection refused"
- "authentication failed"
- "database does not exist"

**Solutions**:

1. **Check PostgreSQL is running**:
   ```bash
   docker-compose -f docker-compose.dev.yml ps postgres
   # Should show "Up"
   ```

2. **Check connection string**:
   ```bash
   # Default for docker-compose.dev.yml
   postgres://revenge:password@localhost:5432/revenge
   ```

3. **Test connection**:
   ```bash
   PGPASSWORD=password psql -h localhost -U revenge -d revenge -c "SELECT 1"
   ```

4. **Check network**:
   ```bash
   # Ensure port is exposed
   docker port <postgres-container> 5432
   ```

---

## Performance Issues

### IDE is Slow

**Symptoms**:
- Laggy typing
- Slow code completion
- High CPU usage

**Solutions**:

1. **Disable unused extensions**:
   - VS Code: Disable workspace-level
   - Reduce active extensions to essentials

2. **Exclude large folders**:
   - Add to `.vscode/settings.json`:
     ```json
     {
       "files.watcherExclude": {
         "**/.git/objects/**": true,
         "**/node_modules/**": true,
         "**/tmp/**": true,
         "**/.archive/**": true
       }
     }
     ```

3. **Reduce file watchers**:
   - Disable auto-save
   - Increase debounce delay

4. **Use Zed for quick edits**:
   - Zed is significantly faster for simple editing tasks

---

### Builds are Slow

**Symptoms**:
- `go build` takes > 30 seconds
- Incremental builds not working

**Solutions**:

1. **Use build cache**:
   ```bash
   # Should be enabled by default
   go env GOCACHE
   ```

2. **Parallel builds**:
   ```bash
   # Use multiple cores
   go build -p 8 ./...
   ```

3. **Disable CGo if not needed**:
   ```bash
   CGO_ENABLED=0 go build ./...
   ```

4. **Use Air for incremental builds**:
   ```bash
   # Air only rebuilds changed packages
   air
   ```

---

## Network and Connectivity

### Can't Fetch Dependencies

**Symptoms**:
- `go get` fails
- `npm install` fails
- Timeout errors

**Solutions**:

1. **Check network connection**:
   ```bash
   ping github.com
   ping npmjs.org
   ```

2. **Configure proxy if behind corporate firewall**:
   ```bash
   # Go proxy
   export GOPROXY=https://proxy.golang.org,direct

   # NPM proxy
   npm config set proxy http://proxy:port
   npm config set https-proxy http://proxy:port
   ```

3. **Use mirror/alternative registry**:
   ```bash
   # Go: use Athens or Artifactory
   export GOPROXY=https://athens.mycompany.com

   # NPM: use alternative registry
   npm config set registry https://registry.npmmirror.com
   ```

---

## Getting More Help

If your issue isn't covered here:

1. **Check tool-specific docs**:
   - [VS Code Troubleshooting](../../.vscode/docs/SETTINGS.md)
   - [Zed Troubleshooting](../../.zed/docs/TROUBLESHOOTING.md)
   - [Coder Troubleshooting](../../.coder/docs/TROUBLESHOOTING.md)
   - [Git Hooks Troubleshooting](../../.githooks/docs/TROUBLESHOOTING.md)

2. **Search GitHub Issues**:
   - Someone may have had the same problem
   - Open a new issue if not

3. **Ask Claude Code**:
   - Use Claude Code to diagnose and fix issues
   - Claude has access to all project documentation

4. **Community Support**:
   - GitHub Discussions
   - Discord (if available)

---

**Maintained By**: Development Team
**Last Updated**: 2026-01-31
