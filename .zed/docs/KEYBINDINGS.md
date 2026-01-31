# Zed Keybindings Reference

> Comprehensive keyboard shortcuts for Revenge development

**Last Updated**: 2026-01-31
**Platform Notes**: macOS uses `Cmd`, Linux/Windows use `Ctrl`

---

## Quick Reference Table

### File & Navigation (Most Used)

| Action | macOS | Linux/Windows | Category |
|--------|-------|---------------|----------|
| Go to File | `Cmd+P` | `Ctrl+P` | **File** |
| Go to Symbol | `Cmd+Shift+O` | `Ctrl+Shift+O` | **Symbol** |
| Go to Definition | `F12` | `F12` | **Code** |
| Find References | `Shift+F12` | `Shift+F12` | **Code** |
| Go to Line | `Cmd+G` | `Ctrl+G` | **Navigation** |
| Command Palette | `Cmd+Shift+P` | `Ctrl+Shift+P` | **General** |
| Find in Files | `Cmd+Shift+F` | `Ctrl+Shift+F` | **Search** |
| Find in File | `Cmd+F` | `Ctrl+F` | **Search** |
| Replace in File | `Cmd+Alt+F` | `Ctrl+Alt+F` | **Search** |

### Format & Edit

| Action | macOS | Linux/Windows | Category |
|--------|-------|---------------|----------|
| Format Document | `Cmd+Shift+I` | `Ctrl+Shift+I` | **Format** |
| Rename Symbol | `F2` | `F2` | **Edit** |
| Comment Line | `Cmd+/` | `Ctrl+/` | **Edit** |
| Duplicate Line | `Cmd+D` | `Ctrl+D` | **Edit** |
| Delete Line | `Cmd+Shift+K` | `Ctrl+Shift+K` | **Edit** |

---

## Complete Navigation Shortcuts

### File Operations

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| New File | `Cmd+N` | `Ctrl+N` | Creates untitled file |
| Open File | `Cmd+O` | `Ctrl+O` | Browse filesystem |
| Open Folder | `Cmd+K Cmd+O` | `Ctrl+K Ctrl+O` | Open project |
| Go to File | `Cmd+P` | `Ctrl+P` | Fuzzy file search |
| Recent Files | `Cmd+Shift+E` | `Ctrl+Shift+E` | Recently opened |
| Save | `Cmd+S` | `Ctrl+S` | Save current file |
| Save All | `Cmd+Alt+S` | `Ctrl+Alt+S` | Save all files |
| Close File | `Cmd+W` | `Ctrl+W` | Close current tab |
| Close All | `Cmd+K Cmd+W` | `Ctrl+K Ctrl+W` | Close all tabs |

### Symbol Navigation

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Go to Symbol | `Cmd+Shift+O` | `Ctrl+Shift+O` | Functions, types, etc. |
| Go to Definition | `F12` or `Cmd+Click` | `F12` or `Ctrl+Click` | Jump to definition |
| Go to Declaration | `Cmd+Shift+D` | `Ctrl+Shift+D` | Jump to declaration |
| Go to Implementation | `Cmd+Alt+B` | `Ctrl+Alt+B` | Jump to implementation |
| Peek Definition | `Alt+F12` | `Alt+F12` | View without jumping |
| Find References | `Shift+F12` | `Shift+F12` | All usages of symbol |
| Breadcrumbs | `Cmd+Shift+;` | `Ctrl+Shift+;` | Navigate via breadcrumb |

### Movement

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Go to Line | `Cmd+G` | `Ctrl+G` | Jump to line number |
| Next Error | `F8` | `F8` | Next problem |
| Previous Error | `Shift+F8` | `Shift+F8` | Previous problem |
| Next Hunk | `Cmd+G` | `Ctrl+G` | Next git change |
| Previous Hunk | `Shift+Cmd+G` | `Shift+Ctrl+G` | Previous git change |
| Top of File | `Cmd+Home` | `Ctrl+Home` | Jump to start |
| End of File | `Cmd+End` | `Ctrl+End` | Jump to end |

### Editor Tabs

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Next Tab | `Cmd+Tab` | `Ctrl+Tab` | Cycle forward |
| Previous Tab | `Cmd+Shift+Tab` | `Ctrl+Shift+Tab` | Cycle backward |
| New Tab | `Cmd+T` | `Ctrl+T` | Create new tab |
| Close Tab | `Cmd+W` | `Ctrl+W` | Close current |
| Close Other Tabs | `Cmd+Alt+W` | `Ctrl+Alt+W` | Close all but current |
| Reopen Closed Tab | `Cmd+Shift+T` | `Ctrl+Shift+T` | Restore last closed |

---

## Search & Replace

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Find in File | `Cmd+F` | `Ctrl+F` | Find in current |
| Replace in File | `Cmd+Alt+F` | `Ctrl+Alt+F` | Find + replace |
| Find in Files | `Cmd+Shift+F` | `Ctrl+Shift+F` | Project-wide search |
| Replace in Files | `Cmd+Shift+H` | `Ctrl+Shift+H` | Project-wide replace |
| Next Match | `Cmd+G` | `Ctrl+G` | Jump to next result |
| Previous Match | `Shift+Cmd+G` | `Shift+Ctrl+G` | Jump to previous |
| Replace | `Cmd+Shift+1` | `Ctrl+Shift+1` | Replace current match |
| Replace All | `Cmd+Alt+Enter` | `Ctrl+Alt+Enter` | Replace all matches |

**Pro Tip**: In search/replace, use regex by clicking the `.*` button or pressing `Alt+R`.

---

## Editing Shortcuts

### Text Selection

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Select Word | `Cmd+D` | `Ctrl+D` | Selects word at cursor |
| Select All | `Cmd+A` | `Ctrl+A` | Select entire file |
| Select Line | `Cmd+L` | `Ctrl+L` | Select current line |
| Expand Selection | `Shift+Cmd+→` | `Shift+Ctrl+→` | Grow selection right |
| Shrink Selection | `Shift+Cmd+←` | `Shift+Ctrl+←` | Shrink selection |
| Expand to Bracket | `Cmd+Alt+]` | `Ctrl+Alt+]` | Select to bracket |
| Expand by Line | `Shift+Down` | `Shift+Down` | Add line to selection |

### Text Manipulation

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Delete Line | `Cmd+Shift+K` | `Ctrl+Shift+K` | Delete entire line |
| Duplicate Line | `Cmd+D` or Copy+Paste | `Ctrl+D` | Duplicate line |
| Move Line Up | `Alt+Up` | `Alt+Up` | Move line up |
| Move Line Down | `Alt+Down` | `Alt+Down` | Move line down |
| Insert Line Before | `Cmd+Shift+Enter` | `Ctrl+Shift+Enter` | New line above |
| Insert Line After | `Cmd+Enter` | `Ctrl+Enter` | New line below |
| Comment Line | `Cmd+/` | `Ctrl+/` | Toggle comment |
| Block Comment | `Cmd+Alt+/` | `Ctrl+Alt+/` | Multi-line comment |

### Code Intelligence

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Format Document | `Cmd+Shift+I` | `Ctrl+Shift+I` | Format entire file |
| Format Selection | `Cmd+Alt+I` | `Ctrl+Alt+I` | Format selected code |
| Rename Symbol | `F2` | `F2` | Rename with refactor |
| Extract Function | `Cmd+Alt+E` | `Ctrl+Alt+E` | Extract to new function |
| Find Issues | `Cmd+Shift+M` | `Ctrl+Shift+M` | Open diagnostics |
| Quick Fix | `Cmd+.` | `Ctrl+.` | Suggest fixes |
| Show Signature Help | `Cmd+Shift+Space` | `Ctrl+Shift+Space` | Parameter hints |

---

## Multi-Cursor & Selection

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Add Cursor Up | `Cmd+Alt+Up` | `Ctrl+Alt+Up` | Multi-cursor above |
| Add Cursor Down | `Cmd+Alt+Down` | `Ctrl+Alt+Down` | Multi-cursor below |
| Add Cursor at End of Selection | `Cmd+D` | `Ctrl+D` | Select next occurrence |
| Select All Occurrences | `Cmd+Alt+L` | `Ctrl+Alt+L` | Multi-cursor all matches |
| Toggle Block Selection | `Alt+Shift+Cmd+M` | `Alt+Shift+Ctrl+M` | Column selection |

---

## Git Integration

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Show Git Blame | - | - | View → Inline Blame |
| Next Hunk | `Cmd+G` | `Ctrl+G` | Next changed section |
| Previous Hunk | `Shift+Cmd+G` | `Shift+Ctrl+G` | Previous hunk |
| Stage Hunk | - | - | Click gutter indicator |
| Discard Hunk | - | - | Right-click gutter |
| Open in Git | `Cmd+Alt+G` | `Ctrl+Alt+G` | Git operations menu |
| Undo Last Git Change | `Cmd+Z` | `Ctrl+Z` | Works for staged changes |

**Note**: Git operations in Zed are mainly viewing. Use terminal for commits/push:
```bash
git add .
git commit -m "message"
git push
```

---

## Terminal & Panels

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Toggle Terminal | `Ctrl+\`` | `Ctrl+\`` | Show/hide terminal |
| New Terminal Tab | `Ctrl+Shift+\`` | `Ctrl+Shift+\`` | New terminal |
| Next Terminal | `Ctrl+Tab` | `Ctrl+Tab` | Cycle terminals |
| Previous Terminal | `Ctrl+Shift+Tab` | `Ctrl+Shift+Tab` | Previous terminal |
| Focus Editor | `Cmd+1` | `Ctrl+1` | Switch to editor |
| Focus Terminal | `Cmd+2` | `Ctrl+2` | Switch to terminal |
| Focus Explorer | `Cmd+Shift+E` | `Ctrl+Shift+E` | Toggle file explorer |
| Focus Outline | `Cmd+Shift+O` | `Ctrl+Shift+O` | Show outline panel |

---

## Settings & UI

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Settings | `Cmd+,` | `Ctrl+,` | Open settings |
| Command Palette | `Cmd+Shift+P` | `Ctrl+Shift+P` | Run commands |
| Theme Selector | `Cmd+K Cmd+T` | `Ctrl+K Ctrl+T` | Change theme |
| Font Size Up | `Cmd+Plus` | `Ctrl+Plus` | Increase font |
| Font Size Down | `Cmd+Minus` | `Ctrl+Minus` | Decrease font |
| Reset Font Size | `Cmd+0` | `Ctrl+0` | Default font size |
| Toggle Zoom | - | - | View → Zoom |
| Toggle Full Screen | `Cmd+Ctrl+F` | `F11` | Full screen mode |
| Toggle Sidebar | `Cmd+B` | `Ctrl+B` | Show/hide sidebar |

---

## Claude Code Integration

| Action | macOS | Linux/Windows | Notes |
|--------|-------|---------------|-------|
| Open Agent | `Cmd+Shift+A` | `Ctrl+Shift+A` | Claude Code panel |
| Inline Assist | `Cmd+Enter` | `Ctrl+Enter` | Quick AI help |
| Generate Code | - | - | Use Agent panel |
| Explain Code | - | - | Use Agent panel |
| Fix Issues | - | - | Use Agent panel |

---

## VS Code vs Zed Keybinding Comparison

### Same Bindings (Compatible)

| Action | Binding | Compatible |
|--------|---------|-----------|
| Go to File | `Cmd/Ctrl+P` | ✅ |
| Command Palette | `Cmd/Ctrl+Shift+P` | ✅ |
| Find in File | `Cmd/Ctrl+F` | ✅ |
| Find & Replace | `Cmd/Ctrl+H` | ✅ (mostly) |
| Settings | `Cmd/Ctrl+,` | ✅ |
| Format Document | `Cmd/Ctrl+Shift+I` | ✅ |
| Comment Line | `Cmd/Ctrl+/` | ✅ |
| Multi-cursor | `Cmd/Ctrl+D` | ✅ |
| Rename Symbol | `F2` | ✅ |

### Different Bindings (VS Code vs Zed)

| Action | VS Code | Zed | Notes |
|--------|---------|-----|-------|
| Go to Definition | `F12` / `Cmd+Click` | `F12` / `Cmd+Click` | Same |
| Show Problems | `Cmd/Ctrl+Shift+M` | `Cmd/Ctrl+Shift+M` | Same |
| Find References | `Shift+F12` | `Shift+F12` | Same |
| Terminal | `Ctrl+\`` | `Ctrl+\`` | Same |
| Toggle Sidebar | `Cmd/Ctrl+B` | `Cmd/Ctrl+B` | Same |

**Tip**: Most shortcuts are compatible. Use VS Code cheat sheet as a reference.

---

## Platform Differences (macOS vs Linux)

### Key Modifier Differences

| Feature | macOS | Linux/Windows |
|---------|-------|---------------|
| Primary modifier | `Cmd` | `Ctrl` |
| Secondary modifier | `Alt` | `Alt` |
| Extended key combo | `Cmd+Alt` | `Ctrl+Alt` |

### Common Differences

**Window Management**:
- macOS: Mission Control, Spaces
- Linux/Windows: Window managers (varies by system)

**System-level Shortcuts**:
- macOS: `Cmd+Space` (Spotlight) — may conflict with Zed
- Linux: `Super` key handling varies
- Windows: Windows key handling varies

### Resolving Conflicts

If a system shortcut conflicts with Zed:

**macOS**:
1. System Settings → Keyboard → Shortcuts
2. Disable conflicting shortcut
3. Or remap in Zed's `keybindings.json`

**Linux** (GNOME):
1. Settings → Keyboard → Shortcuts
2. Disable or change the shortcut

**Windows**:
1. Settings → Keyboard → Shortcuts
2. Usually no conflicts with Zed

---

## Custom Keybindings

### Create Custom Bindings

1. Open Command Palette: `Cmd/Ctrl+Shift+P`
2. Run: "Open User Keybindings"
3. Add custom binding:

```json
[
  {
    "bindings": ["ctrl-alt-d"],
    "command": "editor::DuplicateLine"
  },
  {
    "bindings": ["cmd-shift-r"],
    "command": "editor::ToggleComments",
    "context": "Editor"
  }
]
```

### Common Custom Bindings for Revenge

Add to your keybindings file:

```json
[
  {
    "bindings": ["cmd-alt-t"],
    "command": "editor::OpenTerminal",
    "context": "Editor"
  },
  {
    "bindings": ["cmd-alt-g"],
    "command": "vcs::ToggleBlameLine",
    "context": "Editor"
  },
  {
    "bindings": ["cmd-shift-l"],
    "command": "editor::ToggleSoftWrap",
    "context": "Editor"
  },
  {
    "bindings": ["cmd-alt-f"],
    "command": "editor::RenameSymbol",
    "context": "Editor"
  }
]
```

---

## Command Palette Commands

Press `Cmd/Ctrl+Shift+P` to access these:

### File Commands

| Command | Description |
|---------|-------------|
| `File: New` | Create new file |
| `File: Open` | Open file dialog |
| `File: Open Folder` | Open folder/project |
| `File: Save` | Save current file |
| `File: Save As` | Save with new name |
| `File: Close` | Close current file |
| `File: Close Folder` | Close project |

### Editor Commands

| Command | Description |
|---------|-------------|
| `Editor: Format Document` | Format entire file |
| `Editor: Toggle Comments` | Comment/uncomment |
| `Editor: Duplicate Selection` | Duplicate line/selection |
| `Editor: Delete Line` | Delete entire line |
| `Editor: Go to Line` | Jump to line |
| `Editor: Go to Definition` | Jump to definition |
| `Editor: Rename Symbol` | Refactor rename |

### Git Commands

| Command | Description |
|---------|-------------|
| `Git: Clone Repository` | Clone from remote |
| `Git: Pull` | Fetch + merge |
| `Git: Push` | Push to remote |
| `Git: Open in Git Graph` | Visual git history |

### Search Commands

| Command | Description |
|---------|-------------|
| `Search: Find in Files` | Project-wide search |
| `Search: Replace in Files` | Project-wide replace |
| `Search: Find Next Match` | Navigate results |

### View Commands

| Command | Description |
|---------|-------------|
| `View: Toggle Sidebar` | Show/hide sidebar |
| `View: Toggle Terminal` | Show/hide terminal |
| `View: Toggle Zoom` | Zoom in/out |
| `View: Select Theme` | Change editor theme |
| `View: Select Font Family` | Change font |

### LSP Commands

| Command | Description |
|---------|-------------|
| `Editor: Restart Language Server` | Reload LSP |
| `Editor: Format Selection` | Format selected code |
| `Editor: Go to Declaration` | Jump to declaration |
| `Editor: Show Hover` | Show info at cursor |

---

## Pro Tips

### Efficiency Tips

1. **Use Cmd+P frequently** for file navigation
2. **Learn Cmd+Shift+P** - it's your command center
3. **Multi-cursor (Cmd+D)** for quick bulk edits
4. **Go to Symbol (Cmd+Shift+O)** for code structure
5. **Rename (F2)** instead of manual find-replace

### Speed Up Workflows

```bash
# Terminal shortcuts to improve productivity
# Add to ~/.zed/settings.json for quick access

# Hot reload backend
alias ar='air'  # Auto-rebuild on Go changes

# Format and lint
alias ff='go fmt ./... && go vet ./...'  # Format + vet

# Run tests
alias tt='go test ./...'  # Run all tests

# Frontend dev
alias ndev='npm run dev'  # SvelteKit dev server
```

### Debugging Navigation

When debugging Go code:
1. Use `F12` to go to definition
2. Use `Shift+F12` to see all uses
3. Use `Cmd+Shift+O` to jump between functions
4. Hover over variables for type hints

### Refactoring

1. Use `F2` to rename symbols project-wide
2. Use `Cmd+.` to see refactoring suggestions
3. Use `Cmd+X` to cut (better than delete for undo)

---

## Accessibility Features

| Feature | Shortcut | Notes |
|---------|----------|-------|
| Zoom In | `Cmd/Ctrl+Plus` | Larger text |
| Zoom Out | `Cmd/Ctrl+Minus` | Smaller text |
| Reset Zoom | `Cmd/Ctrl+0` | Default size |
| High Contrast | Via Settings | View → Appearance |
| Increase Font Size | Settings → Buffer Font Size | For visually impaired |
| Screen Reader Support | Check Accessibility docs | Partial support |

---

## References

- **Official Zed Keybindings**: https://zed.dev/docs/key-bindings
- **VS Code Keybindings**: https://code.visualstudio.com/docs/getstarted/keybindings
- **Project Settings**: [../settings.json](../settings.json)
- **SETUP.md**: [SETUP.md](SETUP.md) for installation

