# Claude Code Instructions for JetBrains IDEs

**Purpose**: JetBrains-specific guidance for Claude Code when assisting with the Revenge project

**Last Updated**: 2026-01-31

---

## IDE Context

When the user mentions JetBrains, GoLand, IntelliJ IDEA, or Gateway, they are using JetBrains IDEs.

### Project IDEs

**Revenge supports**:
- **GoLand** - Primary IDE for Go development (recommended for backend developers)
- **IntelliJ IDEA Ultimate** - For full-stack development (Go + Svelte + Python)
- **JetBrains Gateway** - For remote development on Coder workspaces

---

## File References

### When referencing code in JetBrains IDEs

Use standard file:line format that JetBrains recognizes:

**Good**:
```
internal/api/handlers/movie.go:42
cmd/revenge/main.go:15-25
```

**Explanation**: JetBrains IDEs support `Cmd/Ctrl+Click` on `filename:line` to navigate.

---

## Run Configurations

### When suggesting running code

JetBrains uses **Run Configurations** instead of launch.json (VS Code) or tasks.

**Suggest creating run configuration**:
```
To run the server in GoLand:
1. Run → Edit Configurations
2. + → Go Build
3. Name: "Revenge Server"
4. Package path: github.com/lusoris/revenge/cmd/revenge
5. Environment: GOEXPERIMENT=greenteagc,jsonv2
6. Click OK
7. Run with Ctrl+R (macOS) or Shift+F10 (Windows/Linux)
```

**For tests**:
```
To run tests:
1. Right-click internal/ folder
2. Run 'go test internal/...'
```

---

## Code Actions

### When suggesting code changes

JetBrains has powerful refactoring tools - suggest using them:

**Refactoring**:
- **Rename**: `Shift+F6` - Suggest for renaming symbols across project
- **Extract Method**: `Cmd/Ctrl+Alt+M` - Suggest for extracting code blocks
- **Inline**: `Cmd/Ctrl+Alt+N` - Suggest for inlining variables/functions
- **Move**: `F6` - Suggest for moving files or symbols

**Quick Fixes**:
- **Alt+Enter** - Always suggest this for showing available fixes

**Example**:
```
To rename the function across all files:
1. Place cursor on function name
2. Press Shift+F6
3. Enter new name
4. Press Enter - JetBrains updates all references automatically
```

---

## Code Style and Formatting

### JetBrains Formatting

JetBrains has built-in formatters:

**For Go**:
- Uses `gofmt` via gopls automatically
- No need to install separate formatters
- Respects `.editorconfig` in project

**Format commands**:
- **Format file**: `Cmd/Ctrl+Alt+L`
- **Optimize imports**: `Cmd/Ctrl+Alt+O`
- **Reformat and optimize**: Both together

**When suggesting formatting**:
```
Format the file with Cmd/Ctrl+Alt+L and optimize imports with Cmd/Ctrl+Alt+O.
```

---

## Debugging

### When user wants to debug

JetBrains has an excellent debugger - suggest using it:

**Setting breakpoints**:
```
To debug:
1. Click gutter next to line number to set breakpoint (red dot appears)
2. Run → Debug 'Revenge Server' (or Ctrl+D / Shift+F9)
3. Execution stops at breakpoint
```

**Debug actions**:
- **F8** - Step Over
- **F7** - Step Into
- **Shift+F8** - Step Out
- **F9** - Resume
- **Alt+F8** - Evaluate Expression

**Example**:
```
To debug this function:
1. Set breakpoint on line 42 (click gutter)
2. Press Ctrl+D (macOS) or Shift+F9 (Windows/Linux) to debug
3. When stopped, press Alt+F8 to evaluate expressions
4. Use F8 to step through code
```

---

## Database Tools

### When user works with database

JetBrains has built-in database tools - suggest using them instead of external clients:

**Connecting to PostgreSQL**:
```
To query the database:
1. Database tool window (right sidebar)
2. + → Data Source → PostgreSQL
3. Host: localhost, Port: 5432, Database: revenge, User: revenge, Password: revenge
4. Test Connection → Download driver if prompted
5. OK
```

**Running queries**:
```
To run SQL:
1. Right-click database → New → Query Console
2. Write SQL:
   SELECT * FROM users LIMIT 10;
3. Ctrl+Enter to execute
```

---

## Terminal

### When suggesting terminal commands

JetBrains has integrated terminal:

**Opening terminal**:
- **Alt+F12** - Opens terminal at project root
- Or **View** → **Tool Windows** → **Terminal**

**Example**:
```
Run in the IDE terminal (Alt+F12):
go test ./...
```

---

## Plugins

### When suggesting plugins

JetBrains has a rich plugin ecosystem:

**Installing plugins**:
```
To install [Plugin Name]:
1. Settings → Plugins → Marketplace
2. Search "[Plugin Name]"
3. Click Install
4. Restart IDE if prompted
```

**Useful plugins for Revenge**:
- **.env files support** - Environment file syntax
- **GitToolBox** - Enhanced Git integration
- **Rainbow Brackets** - Colorful bracket matching
- **Key Promoter X** - Learn shortcuts

---

## Version Control (Git)

### When suggesting Git operations

JetBrains has powerful Git integration:

**Commit**:
```
To commit changes:
1. Cmd/Ctrl+K to open commit dialog
2. Select files
3. Write message
4. Ensure "Reformat code" and "Optimize imports" are checked
5. Click Commit (or Commit and Push)
```

**Viewing changes**:
- **Cmd/Ctrl+Alt+Shift+D** - Show diff
- **Alt+9** - Show Version Control panel

**Branches**:
- **Git** → **Branches** - Manage branches
- **Cmd/Ctrl+Shift+`** - VCS operations popup

---

## Common IDE Operations

### When guiding users through IDE tasks

**Search**:
- **Shift Shift** - Search Everywhere (files, symbols, actions, settings)
- **Cmd/Ctrl+Shift+A** - Find Action (command palette)
- **Cmd/Ctrl+Shift+F** - Find in Files
- **Cmd/Ctrl+Shift+R** - Replace in Files

**Navigation**:
- **Cmd/Ctrl+B** - Go to Declaration
- **Cmd/Ctrl+Alt+B** - Go to Implementation
- **Cmd/Ctrl+F7** - Find Usages
- **Cmd/Ctrl+E** - Recent Files
- **Cmd/Ctrl+Shift+E** - Recent Locations

**Always suggest keyboard shortcuts** - JetBrains users value them.

---

## Remote Development (Gateway)

### When user mentions Gateway or Coder

User is doing remote development via JetBrains Gateway:

**Connection issues**:
```
If Gateway won't connect:
1. Verify workspace is running: coder list
2. Start if stopped: coder start revenge-dev
3. Re-login in Gateway: Settings → Coder → Logout → Login
```

**Port forwarding**:
- Gateway auto-forwards ports (8096, 5173, 5432, 6379)
- Access at `http://localhost:PORT` from local browser

**Performance tips**:
```
If Gateway is slow:
1. Enable Power Save Mode: File → Power Save Mode
2. Reduce heap: Help → Change Memory Settings → 2GB
3. Disable unused plugins: Settings → Plugins
```

---

## Project Structure

### JetBrains-specific project files

**.idea/** folder:
- Contains project settings (git-ignored)
- Auto-generated by IDE
- Don't manually edit unless necessary

**No shared settings**:
- Unlike VS Code's `settings.json`, JetBrains settings are in `.idea/`
- Team shares settings via `.editorconfig` (already configured)

---

## Language-Specific Tips

### Go

**gopls integration**:
- Built-in, no setup needed
- **Settings** → **Go** → **gopls** to configure

**Go tools**:
- IDE auto-installs Go tools (goimports, dlv, etc.)
- Prompts if missing

### Python

**For scripts in Revenge**:
- Install Python plugin (IntelliJ IDEA only)
- Configure interpreter: **Settings** → **Project** → **Python Interpreter**

### TypeScript/Svelte

**Frontend development**:
- TypeScript built-in (IntelliJ IDEA Ultimate)
- Install Svelte plugin for `.svelte` files
- Prettier integration: **Settings** → **Prettier**

---

## Performance Considerations

### When user reports slowness

**First suggestions**:
1. **Increase heap**: Help → Change Memory Settings → 4GB
2. **Power Save Mode**: File → Power Save Mode (when not coding)
3. **Invalidate caches**: File → Invalidate Caches → Invalidate and Restart

**For remote (Gateway)**:
4. **Check network latency**: `ping coder.ancilla.lol` (should be <100ms)
5. **Reduce bandwidth usage**: Settings → Editor → Disable inline hints

---

## Troubleshooting

### Common issues and solutions

**Go SDK not detected**:
```
Settings → Go → GOROOT → + → Local → /usr/local/go
```

**gopls not working**:
```
Settings → Go → gopls → Restart gopls
```

**Gateway connection fails**:
```
coder login https://coder.ancilla.lol
coder start revenge-dev
Gateway → Settings → Coder → Re-login
```

**Slow performance**:
```
Help → Change Memory Settings → 4096 MB
File → Power Save Mode (temporarily)
```

Full troubleshooting: [.jetbrains/docs/TROUBLESHOOTING.md](.jetbrains/docs/TROUBLESHOOTING.md)

---

## Documentation References

### When user needs help

**JetBrains docs** (in `.jetbrains/docs/`):
- **INDEX.md** - Documentation hub
- **SETUP.md** - Installation and configuration
- **REMOTE_DEVELOPMENT.md** - Gateway + Coder
- **TROUBLESHOOTING.md** - Common issues

**Coder docs** (in `.coder/docs/`):
- **JETBRAINS_INTEGRATION.md** - Complete Gateway guide
- **REMOTE_WORKFLOW.md** - Remote development workflow

**Shared docs** (in `.shared/docs/`):
- **TOOL_COMPARISON.md** - When to use GoLand vs VS Code vs Zed
- **WORKFLOWS.md** - Development workflows

---

## Code Examples

### When providing code examples for JetBrains

**Include shortcuts and IDE features**:

```go
// To extract this to a method:
// 1. Select lines 42-45
// 2. Press Cmd/Ctrl+Alt+M
// 3. Name: validateUser
// 4. Press Enter

if user == nil {
    return nil, errors.New("user is nil")
}
if user.Email == "" {
    return nil, errors.New("email is required")
}
```

**Suggest using IDE features**:
```
Instead of manually renaming everywhere, use Shift+F6 on the function name.
JetBrains will safely rename all references across the entire project.
```

---

## Best Practices

### When assisting JetBrains users

1. **Suggest keyboard shortcuts** - JetBrains users value efficiency
2. **Recommend refactoring tools** - JetBrains excels at refactoring
3. **Mention IDE features** - Database tools, HTTP client, etc.
4. **Use correct terminology**:
   - "Run Configuration" (not launch.json)
   - "Inspection" (not linter/diagnostic)
   - "Quick Fix" (not code action)
5. **Reference JetBrains docs** when appropriate

---

## Comparison to Other IDEs

### When user asks about alternatives

**JetBrains advantages**:
- Advanced refactoring (Extract Interface, Move Package, etc.)
- Excellent debugger with conditional breakpoints
- Built-in database tools
- Deep code analysis

**When to suggest alternatives**:
- **VS Code**: Free, lighter, good extensions
- **Zed**: Ultra-fast, minimal, Rust-based

See: [.shared/docs/TOOL_COMPARISON.md](.shared/docs/TOOL_COMPARISON.md)

---

## Example Interactions

### Good responses for JetBrains users

**User**: "How do I rename this function everywhere?"

**Response**:
```
To rename the function across the entire project:
1. Place cursor on the function name
2. Press Shift+F6
3. Enter the new name
4. Press Enter

JetBrains will safely rename all references, including in comments and strings
if you select those options. Preview the changes before applying if needed.
```

**User**: "The debugger isn't stopping at my breakpoint"

**Response**:
```
Check these common issues:
1. Ensure breakpoint is enabled (red filled circle, not hollow)
2. Verify you're running Debug (Ctrl+D / Shift+F9), not Run (Ctrl+R / Shift+F10)
3. Check if breakpoint condition is met (right-click breakpoint → Edit)
4. Ensure code is not optimized out

To debug:
- Set breakpoint: Click gutter next to line number
- Start debugger: Ctrl+D (macOS) or Shift+F9 (Windows/Linux)
- Wait for execution to stop
```

**User**: "Gateway won't connect to my Coder workspace"

**Response**:
```
Try these steps:
1. Verify workspace is running:
   coder list

2. Start workspace if stopped:
   coder start revenge-dev

3. Re-login in Gateway:
   - Open Gateway
   - Settings → Coder → Logout
   - Click Login → https://coder.ancilla.lol
   - Authenticate in browser

4. Check network connectivity:
   ping coder.ancilla.lol

If still failing, see: .coder/docs/TROUBLESHOOTING.md
```

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
