# Zed Troubleshooting Guide

> Solutions for common issues with Zed editor in Revenge development

**Last Updated**: 2026-01-31

---

## Table of Contents

1. [LSP Issues](#lsp-issues)
2. [Formatter Issues](#formatter-issues)
3. [Performance Problems](#performance-problems)
4. [Git Integration](#git-integration)
5. [Debugging](#debugging)
6. [Log Files and Diagnostics](#log-files-and-diagnostics)

---

## LSP Issues

### gopls Not Starting

**Symptoms**:
- No code completion for Go files
- No go-to-definition support
- Error message: "gopls crashed" or similar

**Diagnosis**:

```bash
# Check Go version
go version
# Should be 1.25.6 or higher for Revenge

# Check gopls installation
gopls version
# Should show version info

# Check if gopls is in PATH
which gopls
# Should show: /path/to/gopls
```

**Solutions**:

**1. Install/Reinstall gopls**:
```bash
# Remove old version
go clean -i github.com/golang/tools/gopls

# Install latest
go install github.com/golang/tools/gopls@latest

# Verify
gopls version
```

**2. Restart Zed's Language Server**:
- Press `Cmd/Ctrl+Shift+P`
- Search: "Restart Language Server"
- Select "Restart LSP"
- Wait 5-10 seconds for gopls to start

**3. Check Zed Logs**:
```bash
# macOS
tail -f ~/Library/Logs/Zed/zed.log

# Linux
tail -f ~/.local/share/zed/zed.log

# Look for gopls startup messages
```

**4. Verify Go Workspace**:
```bash
# Check go.mod exists
ls -la go.mod

# Download dependencies
go mod download

# Tidy up
go mod tidy
```

**5. Update gopls Configuration**:

Edit `.zed/settings.json`:
```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "hints": {
          "assignVariableTypes": true,
          "compositeLiteralFields": true,
          "constantValues": true
        },
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

**6. Check for Go version conflicts**:
```bash
# If you have multiple Go versions
go version

# Ensure correct version in PATH
which go
echo $PATH

# May need to update PATH in ~/.zprofile or ~/.bashrc
export PATH="/usr/local/go/bin:$PATH"
```

**Still not working?**:
- Check Zed's debug log for exact error
- Try creating simple test file: `package main\n\nfunc main() {}`
- Verify gopls can start independently: `gopls` (should start a server)

---

### gopls Slow or Crashes on Large Projects

**Symptoms**:
- High CPU/memory usage
- Frequent crashes
- Slow completions

**Solutions**:

**1. Increase gopls timeout**:
```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "experimentalMemoryModel": true
      }
    }
  }
}
```

**2. Disable expensive features**:
```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "analyses": {
          "shadow": false
        },
        "staticcheck": false
      }
    }
  }
}
```

**3. Check for errant build processes**:
```bash
# See what's using CPU
ps aux | grep go

# Kill any old gopls processes
pkill -f gopls

# Restart Zed
```

**4. Limit project scope**:
- If working in subdirectory, open that in Zed instead
- `.zed/settings.json`:
  ```json
  {
    "file_scan_exclusions": [
      "**/.git",
      "**/node_modules",
      "**/dist",
      "**/__pycache__",
      "**/.archive",
      "**/bin",
      "**/tmp",
      "**/vendor"
    ]
  }
  ```

---

### Python (Ruff) LSP Issues

**Symptoms**:
- No Python linting
- Ruff-lsp not found
- Import errors not shown

**Diagnosis**:

```bash
# Check ruff installation
ruff --version

# Check ruff-lsp (if using)
pip show ruff
```

**Solutions**:

**1. Install Ruff**:
```bash
# Using pip
pip install ruff --upgrade

# Using uv (faster)
uv pip install ruff

# Verify
ruff --version
```

**2. Update Zed configuration**:
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
  },
  "languages": {
    "Python": {
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

**3. Check ruff.toml**:
```bash
# Should exist at project root
cat ruff.toml

# Key settings for Revenge
# line-length = 88
# target-version = "py312"
```

**4. Restart Python LSP**:
- `Cmd/Ctrl+Shift+P` → "Restart Language Server"
- Select Python/Ruff
- Wait for restart

**Still not working?**:
```bash
# Check Python version
python3 --version  # Should be 3.12+

# Verify installation
python3 -m pip show ruff

# Try reinstalling
pip uninstall ruff
pip install ruff
```

---

### TypeScript/Svelte LSP Not Working

**Symptoms**:
- No TypeScript errors
- No Svelte syntax highlighting
- Completions not working for TypeScript

**Diagnosis**:

```bash
# Check TypeScript installation
npm list -g typescript

# Check if in project
npm list typescript

# Check Svelte LSP (if installed globally)
npm list -g svelte-language-server
```

**Solutions**:

**1. Install TypeScript dependencies**:
```bash
# In project directory
npm install --save-dev typescript
npm install --save-dev svelte

# Optionally install Svelte language server
npm install --save-dev svelte-language-server
```

**2. Configure TypeScript in Zed**:
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
    }
  }
}
```

**3. Check tsconfig.json**:
```bash
# Should exist in project root
cat tsconfig.json

# Verify paths are correct
# Should have:
# "compilerOptions": {
#   "paths": { ... }
# }
```

**4. Restart TypeScript server**:
- `Cmd/Ctrl+Shift+P` → "Restart Language Server"
- Select TypeScript
- Wait 5-10 seconds

**Still having issues?**:
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Restart Zed
```

---

## Formatter Issues

### Format on Save Not Working

**Symptoms**:
- File not formatted when saved
- Manual format (Cmd/Ctrl+Shift+I) works
- No error messages

**Diagnosis**:

```bash
# Check format_on_save setting
# In .zed/settings.json should have:
# "format_on_save": "on"

# Check formatter is installed and working
go fmt -h  # For Go
prettier --version  # For TypeScript/Svelte
ruff format --help  # For Python
```

**Solutions**:

**1. Verify setting is enabled**:
```json
{
  "format_on_save": "on"
}
```

**2. Check language-specific formatter**:
```json
{
  "languages": {
    "Go": {
      "format_on_save": "on"
    },
    "Python": {
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

**3. Ensure formatters are installed**:
```bash
# Go (built-in via gopls)
gopls version

# Python
ruff --version

# TypeScript/Svelte
npm list prettier prettier-plugin-svelte
```

**4. Check file is saved in correct format**:
```bash
# Verify file extension
ls -la your_file.go
ls -la your_file.py
ls -la your_file.ts
```

**5. Manual format test**:
- `Cmd/Ctrl+Shift+I` (format document)
- Check if it works manually
- If yes, issue is with on-save trigger

**Still not working?**:
- Check Zed logs: `View → Toggle Log Panel`
- Look for formatter errors
- Try closing and reopening file

---

### Prettier Not Formatting TypeScript/Svelte

**Symptoms**:
- Prettier command fails
- Plugin not found error
- Svelte files not formatting

**Solutions**:

**1. Install Prettier with Svelte plugin**:
```bash
npm install --save-dev prettier prettier-plugin-svelte
```

**2. Configure in .zed/settings.json**:
```json
{
  "languages": {
    "TypeScript": {
      "formatter": {
        "external": {
          "command": "prettier",
          "arguments": ["--stdin-filepath", "{buffer_path}"]
        }
      }
    },
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

**3. Create .prettierrc.json**:
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

**4. Test Prettier manually**:
```bash
# Format a file
prettier --write src/routes/+page.svelte

# Check for errors
prettier --check src/routes/+page.svelte
```

**Still failing?**:
- Ensure prettier can find the plugin: `npm list prettier-plugin-svelte`
- Try without plugin first to isolate issue
- Check Zed logs for exact error

---

### Ruff Format Errors

**Symptoms**:
- "ruff format" command not found
- Format fails silently
- Python files not formatting

**Solutions**:

**1. Verify ruff format is available**:
```bash
ruff format --help
# If command not found, reinstall:
pip install --upgrade ruff
```

**2. Configure correctly**:
```json
{
  "languages": {
    "Python": {
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

**3. Test ruff format**:
```bash
# Create test file
echo 'x=1' > test.py

# Format it
ruff format test.py

# Should become
echo 'x = 1' > test.py
```

**4. Check ruff version**:
```bash
ruff --version
# Should be recent (2024+)

# If old, upgrade
pip install --upgrade ruff
```

---

## Performance Problems

### Zed Using Too Much Memory

**Symptoms**:
- High memory usage
- Zed becoming sluggish
- System running out of memory

**Solutions**:

**1. Check memory usage**:
```bash
# macOS
top -l1 | grep Zed

# Linux
ps aux | grep -i zed

# Note the RSS (memory) column
```

**2. Disable heavy features**:
```json
{
  "show_completions_on_input": false,
  "show_inline_diagnostics": false,
  "hover_popover_enabled": false,
  "show_edit_predictions": false
}
```

**3. Reduce file scan scope**:
```json
{
  "file_scan_exclusions": [
    "**/.git",
    "**/node_modules",
    "**/dist",
    "**/__pycache__",
    "**/.archive",
    "**/bin",
    "**/tmp",
    "**/vendor",
    "**/.next",
    "**/build"
  ]
}
```

**4. Limit LSP features**:
```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "analyses": {
          "shadow": false,
          "unused": false
        },
        "staticcheck": false
      }
    }
  }
}
```

**5. Restart Zed**:
- Close Zed completely
- Wait 10 seconds
- Reopen

**Still high?**:
- Close unused tabs
- Close unused projects
- Restart LSP servers individually

---

### Zed Responding Slowly

**Symptoms**:
- Typing lag
- Menu delays
- Slow file opening

**Solutions**:

**1. Check CPU usage**:
```bash
# macOS/Linux
top
# Look for high CPU processes

# Kill problematic processes
pkill -f gopls  # If gopls stuck
pkill -f ruff   # If ruff stuck
```

**2. Reduce visual effects**:
```json
{
  "cursor_blink": false,
  "current_line_highlight": "gutter"
}
```

**3. Disable extensions/features**:
- Disable "show completions on input"
- Disable LSP inline hints
- Disable telemetry: `"telemetry.diagnostics": false`

**4. Check open files**:
- Close unnecessary tabs
- Large files can be slow
- Consider splitting very large files

**5. Update Zed**:
```bash
# Check version
zed --version

# Update via package manager
brew upgrade zed  # macOS
sudo dnf upgrade zed  # Fedora
```

---

### High CPU Usage from gopls

**Symptoms**:
- gopls using 100% CPU
- Fan spinning loudly
- Battery drain on laptop

**Solutions**:

**1. Restart gopls**:
```bash
pkill -f gopls
# Zed will restart it automatically
```

**2. Disable expensive analyses**:
```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "analyses": {
          "shadow": false,
          "unused": false,
          "unreachable": false
        },
        "staticcheck": false,
        "gofumpt": false
      }
    }
  }
}
```

**3. Exclude vendor directories**:
```json
{
  "file_scan_exclusions": [
    "**/vendor",
    "**/.git"
  ]
}
```

**4. Check for errant processes**:
```bash
# See all gopls processes
ps aux | grep gopls

# Should be just 1-2
# Kill extras if any
pkill -f gopls
```

**5. Update gopls**:
```bash
go install github.com/golang/tools/gopls@latest
```

---

## Git Integration

### Git Gutter Not Showing Changes

**Symptoms**:
- No modified line indicators
- Blame line not showing
- Git changes not highlighted

**Solutions**:

**1. Enable git gutter in settings**:
```json
{
  "git": {
    "git_gutter": "tracked_files"
  }
}
```

**2. Verify file is tracked**:
```bash
# Check git status
git status

# File should not be in "Untracked files"
# If untracked, add it
git add your_file.go
```

**3. Ensure git is working**:
```bash
# Check git installation
which git

# Check repository
git log --oneline -1

# Should show recent commit
```

**4. Check file permissions**:
```bash
# Verify .git directory exists
ls -la .git

# Check if readable
git rev-parse --git-dir
```

**5. Restart Zed**:
- Close and reopen Zed
- Sometimes git integration needs restart

---

### Git Blame (Inline) Not Working

**Symptoms**:
- No author/date info on hover
- Blame disabled option grayed out

**Solutions**:

**1. Enable inline blame**:
```json
{
  "git": {
    "inline_blame": {
      "enabled": true
    }
  }
}
```

**2. Hover properly**:
- Click on a line (not in comment)
- Hover briefly to see blame info
- May have delay of 300-500ms

**3. Check git history**:
```bash
# Ensure file has commits
git log -p -- src/file.go

# Should show history
```

---

### Git Merge Conflicts Showing Incorrectly

**Symptoms**:
- Conflict markers not highlighted
- Resolution UI not appearing

**Solutions**:

**1. Verify conflict markers exist**:
```bash
# Should see <<<<<<, ======, >>>>>>
cat file_with_conflict.go | grep -E "^<<<<<<|^======|^>>>>>>"
```

**2. Resolve manually**:
- Zed may not have native conflict UI
- Edit conflict markers manually
- Stage the resolved file: `git add file_with_conflict.go`

**3. Use terminal for complex merges**:
```bash
# Use VS Code or terminal for complex conflicts
code .

# Or use git mergetool
git mergetool
```

---

## Debugging

### Issues with Code Debugging

**Note**: Zed doesn't have native debugging UI like VS Code.

**Solutions**:

**1. Use terminal debugging**:
```bash
# Add debug breakpoint in code
import "runtime/debug"
debug.SetTraceback("all")

# Run with dlv debugger
dlv debug ./cmd/revenge

# Or use VS Code for debugging
code .
```

**2. Use VS Code for complex debugging**:
- Zed is excellent for editing
- VS Code is better for interactive debugging
- Use both side-by-side if needed

**3. Log-based debugging**:
```go
// Add logging instead
log.Printf("value: %v", myVar)

// Or use structured logging
slog.Debug("debug info", "key", value)
```

---

## Log Files and Diagnostics

### Finding Zed Logs

**Location**:

macOS:
```bash
# Zed log
~/Library/Logs/Zed/zed.log

# Real-time
tail -f ~/Library/Logs/Zed/zed.log
```

Linux:
```bash
# Zed log
~/.local/share/zed/zed.log

# Real-time
tail -f ~/.local/share/zed/zed.log
```

Windows:
```powershell
# Log location
$env:APPDATA\Zed\logs\zed.log

# Real-time (PowerShell)
Get-Content $env:APPDATA\Zed\logs\zed.log -Wait
```

### Reading Zed Logs

Look for these keywords:
- `ERROR` - Actual errors
- `gopls` - Go language server output
- `ruff` - Python language server
- `typescript` - TypeScript server
- `crashed` - LSP crashes
- `timeout` - Slow operations

Example error:
```
[ERROR] gopls (3456) crashed: signal: segmentation fault
[ERROR] Failed to spawn LSP: gopls exited with code 1
```

### Enable Debug Logging

In `~/.config/zed/settings.json`:
```json
{
  "log_level": "debug"
}
```

Then restart Zed and check logs for more details.

### Collecting Diagnostic Info

**For bug reports**:

```bash
# Get Zed version
zed --version

# Get Go version (if using Go)
go version

# Get Python version
python3 --version

# Check installed formatters
prettier --version
ruff --version
gopls version

# Get recent log excerpt
tail -100 ~/.local/share/zed/zed.log > zed_logs.txt

# Describe reproduction steps
# Include exact error message
# Include platform (macOS/Linux/Windows)
```

---

## Getting Help

### Check Zed Documentation

- **Official Docs**: https://zed.dev/docs
- **Settings Reference**: https://zed.dev/docs/configuring-zed
- **Troubleshooting**: https://zed.dev/docs/troubleshooting

### Check Project Documentation

- **SETUP.md**: [SETUP.md](SETUP.md) - Installation guide
- **KEYBINDINGS.md**: [KEYBINDINGS.md](KEYBINDINGS.md) - Keyboard reference
- **Settings.json**: [../settings.json](../settings.json) - Project config
- **SOURCE_OF_TRUTH**: [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)

### Zed Community

- **GitHub Issues**: https://github.com/zed-industries/zed/issues
- **Zed Discord**: https://zed.dev/chat
- **Discussions**: https://github.com/zed-industries/zed/discussions

### Revenge Project

- **GitHub Issues**: https://github.com/lusoris/revenge/issues
- **Project Docs**: [../../docs/dev/design/INDEX.md](../../docs/dev/design/INDEX.md)

---

## Checklist for Troubleshooting

Use this checklist when diagnosing issues:

- [ ] Restart Zed completely
- [ ] Check Zed logs for errors
- [ ] Verify tool installation (go, ruff, prettier, etc.)
- [ ] Restart language servers
- [ ] Check file is saved and in correct format
- [ ] Verify project configuration in `.zed/settings.json`
- [ ] Check `.editorconfig` for conflicts
- [ ] Try disabling LSP hints and completions temporarily
- [ ] Check for file permission issues
- [ ] Ensure you're on latest Zed version
- [ ] Search GitHub issues for similar problems
- [ ] Try reproducing in simple test file

---

**Last Updated**: 2026-01-31
**For**: Revenge Media Server Project

