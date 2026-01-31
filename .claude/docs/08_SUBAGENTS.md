# Subagents Reference

> Source: https://code.claude.com/docs/en/sub-agents
> Fetched: 2026-01-31
> Type: html

---

## Overview

Subagents are specialized AI assistants that handle specific tasks. Each runs in its own context window with custom system prompt, specific tool access, and independent permissions.

**Benefits:**
- Preserve context by keeping exploration out of main conversation
- Enforce constraints with limited tool access
- Reuse configurations across projects
- Specialize behavior with focused prompts
- Control costs by routing to cheaper models

---

## Built-in Subagents

### Explore
- **Model:** Haiku (fast, low-latency)
- **Tools:** Read-only (no Write/Edit)
- **Purpose:** File discovery, code search, codebase exploration
- Thoroughness levels: quick, medium, very thorough

### Plan
- **Model:** Inherits from main
- **Tools:** Read-only
- **Purpose:** Codebase research during plan mode

### General-purpose
- **Model:** Inherits from main
- **Tools:** All tools
- **Purpose:** Complex multi-step tasks requiring both exploration and action

---

## Create Subagents

### Using /agents Command

```
/agents
```

- View all available subagents
- Create new subagents with guided setup
- Edit existing configuration
- Delete custom subagents

### Subagent Locations

| Location | Scope | Priority |
|----------|-------|----------|
| `--agents` CLI flag | Current session | 1 (highest) |
| `.claude/agents/` | Current project | 2 |
| `~/.claude/agents/` | All your projects | 3 |
| Plugin's `agents/` | Where plugin enabled | 4 (lowest) |

### CLI-defined Subagents

```bash
claude --agents '{
  "code-reviewer": {
    "description": "Expert code reviewer. Use proactively after code changes.",
    "prompt": "You are a senior code reviewer. Focus on code quality, security, and best practices.",
    "tools": ["Read", "Grep", "Glob", "Bash"],
    "model": "sonnet"
  }
}'
```

---

## Subagent File Format

```markdown
---
name: code-reviewer
description: Reviews code for quality and best practices
tools: Read, Glob, Grep
model: sonnet
---

You are a code reviewer. When invoked, analyze the code and provide
specific, actionable feedback on quality, security, and best practices.
```

---

## Frontmatter Reference

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Unique identifier (lowercase, hyphens) |
| `description` | Yes | When Claude should delegate |
| `tools` | No | Tools subagent can use |
| `disallowedTools` | No | Tools to deny |
| `model` | No | `sonnet`, `opus`, `haiku`, or `inherit` |
| `permissionMode` | No | `default`, `acceptEdits`, `dontAsk`, `bypassPermissions`, `plan` |
| `skills` | No | Skills to preload into context |
| `hooks` | No | Lifecycle hooks scoped to subagent |

---

## Permission Modes

| Mode | Behavior |
|------|----------|
| `default` | Standard permission checking |
| `acceptEdits` | Auto-accept file edits |
| `dontAsk` | Auto-deny permission prompts |
| `bypassPermissions` | Skip all permission checks |
| `plan` | Read-only exploration |

---

## Preload Skills

```yaml
---
name: api-developer
description: Implement API endpoints following team conventions
skills:
  - api-conventions
  - error-handling-patterns
---

Implement API endpoints. Follow the conventions from the preloaded skills.
```

---

## Subagent Hooks

### In Frontmatter

```yaml
---
name: code-reviewer
description: Review code changes with automatic linting
hooks:
  PreToolUse:
    - matcher: "Bash"
      hooks:
        - type: command
          command: "./scripts/validate-command.sh $TOOL_INPUT"
  PostToolUse:
    - matcher: "Edit|Write"
      hooks:
        - type: command
          command: "./scripts/run-linter.sh"
---
```

### In settings.json

```json
{
  "hooks": {
    "SubagentStart": [
      {
        "matcher": "db-agent",
        "hooks": [
          { "type": "command", "command": "./scripts/setup-db-connection.sh" }
        ]
      }
    ]
  }
}
```

---

## Foreground vs Background

- **Foreground:** Blocks main conversation, passes permission prompts through
- **Background:** Runs concurrently, permissions pre-approved, auto-denies unapproved

Toggle: Ask Claude to "run this in the background" or press `Ctrl+B`

Disable: Set `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1`

---

## Example Subagents

### Code Reviewer

```markdown
---
name: code-reviewer
description: Expert code review specialist. Use proactively after code changes.
tools: Read, Grep, Glob, Bash
model: inherit
---

You are a senior code reviewer ensuring high standards of code quality and security.

When invoked:
1. Run git diff to see recent changes
2. Focus on modified files
3. Begin review immediately

Review checklist:
- Code is clear and readable
- Proper error handling
- No exposed secrets
- Good test coverage
```

### Debugger

```markdown
---
name: debugger
description: Debugging specialist for errors and test failures.
tools: Read, Edit, Bash, Grep, Glob
---

You are an expert debugger specializing in root cause analysis.

When invoked:
1. Capture error message and stack trace
2. Identify reproduction steps
3. Isolate the failure location
4. Implement minimal fix
5. Verify solution works
```

### Database Query Validator

```markdown
---
name: db-reader
description: Execute read-only database queries.
tools: Bash
hooks:
  PreToolUse:
    - matcher: "Bash"
      hooks:
        - type: command
          command: "./scripts/validate-readonly-query.sh"
---

Execute SELECT queries only. Cannot modify data.
```

With validation script:

```bash
#!/bin/bash
INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty')

if echo "$COMMAND" | grep -iE '\b(INSERT|UPDATE|DELETE|DROP)\b' > /dev/null; then
  echo "Blocked: Write operations not allowed" >&2
  exit 2
fi
exit 0
```

---

## Disable Subagents

In settings:
```json
{
  "permissions": {
    "deny": ["Task(Explore)", "Task(my-custom-agent)"]
  }
}
```

Via CLI:
```bash
claude --disallowedTools "Task(Explore)"
```
