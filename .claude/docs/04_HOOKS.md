# Claude Code Hooks Reference

> Source: https://code.claude.com/docs/en/hooks
> Fetched: 2026-01-31
> Type: html

---

## Hook Lifecycle

| Hook | When it fires |
|------|---------------|
| `SessionStart` | Session begins or resumes |
| `UserPromptSubmit` | User submits a prompt |
| `PreToolUse` | Before tool execution |
| `PermissionRequest` | When permission dialog appears |
| `PostToolUse` | After tool succeeds |
| `PostToolUseFailure` | After tool fails |
| `SubagentStart` | When spawning a subagent |
| `SubagentStop` | When subagent finishes |
| `Stop` | Claude finishes responding |
| `PreCompact` | Before context compaction |
| `SessionEnd` | Session terminates |
| `Notification` | Claude Code sends notifications |

---

## Configuration

Location: `.claude/settings.json`, `~/.claude/settings.json`

```json
{
  "hooks": {
    "EventName": [
      {
        "matcher": "ToolPattern",
        "hooks": [
          {
            "type": "command",
            "command": "your-command-here"
          }
        ]
      }
    ]
  }
}
```

**Matcher options:**
- Simple strings: `Write` matches only Write tool
- Regex: `Edit|Write` or `Notebook.*`
- `*` or empty: matches all tools

---

## Common Matchers

- `Task` - Subagent tasks
- `Bash` - Shell commands
- `Glob` - File pattern matching
- `Grep` - Content search
- `Read` - File reading
- `Edit` - File editing
- `Write` - File writing
- `WebFetch`, `WebSearch` - Web operations

---

## Hook Types

### Command Hook
```json
{
  "type": "command",
  "command": "npx prettier --write \"$FILE_PATH\"",
  "timeout": 30
}
```

### Prompt-Based Hook (for Stop, SubagentStop)
```json
{
  "type": "prompt",
  "prompt": "Evaluate if Claude should stop: $ARGUMENTS"
}
```

---

## Exit Codes

| Exit Code | Behavior |
|-----------|----------|
| 0 | Success. stdout shown in verbose mode |
| 2 | Blocking error. stderr fed to Claude |
| Other | Non-blocking error. stderr shown |

---

## JSON Output Control

### PreToolUse Decision
```json
{
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "allow|deny|ask",
    "permissionDecisionReason": "My reason",
    "updatedInput": { "field": "new value" },
    "additionalContext": "Extra info for Claude"
  }
}
```

### Stop/SubagentStop Decision
```json
{
  "decision": "block",
  "reason": "Must continue because..."
}
```

### UserPromptSubmit
```json
{
  "decision": "block",
  "reason": "Security policy violation",
  "hookSpecificOutput": {
    "additionalContext": "Current time: ..."
  }
}
```

---

## Environment Variables

- `CLAUDE_PROJECT_DIR` - Absolute path to project root
- `CLAUDE_ENV_FILE` - File for persisting env vars (SessionStart only)
- `CLAUDE_CODE_REMOTE` - "true" for web, empty for local

---

## Best Practices

1. **Validate and sanitize inputs**
2. **Always quote shell variables** - Use `"$VAR"`
3. **Block path traversal** - Check for `..`
4. **Use absolute paths** - Use `$CLAUDE_PROJECT_DIR`
5. **Skip sensitive files** - `.env`, `.git/`, keys

---

## Debugging

```bash
claude --debug  # See hook execution details
```

Check: `/hooks` menu, JSON syntax, permissions, full paths
