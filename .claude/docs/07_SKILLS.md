# Skills Reference

> Source: https://code.claude.com/docs/en/skills
> Fetched: 2026-01-31
> Type: html

---

## Overview

Skills extend Claude's capabilities with project-specific knowledge and workflows. Claude uses them automatically when relevant, or invoke directly with `/skill-name`.

**Locations:**

| Location | Path | Scope |
|----------|------|-------|
| Enterprise | Managed settings | All org users |
| Personal | `~/.claude/skills/<skill-name>/SKILL.md` | All your projects |
| Project | `.claude/skills/<skill-name>/SKILL.md` | This project only |
| Plugin | `<plugin>/skills/<skill-name>/SKILL.md` | Where plugin enabled |

Priority: Enterprise > Personal > Project

---

## Create a Skill

```bash
mkdir -p ~/.claude/skills/explain-code
```

Create `SKILL.md`:

```yaml
---
name: explain-code
description: Explains code with visual diagrams and analogies. Use when explaining how code works.
---

When explaining code, always include:

1. **Start with an analogy**: Compare the code to something from everyday life
2. **Draw a diagram**: Use ASCII art to show flow, structure, or relationships
3. **Walk through the code**: Explain step-by-step what happens
4. **Highlight a gotcha**: What's a common mistake or misconception?
```

---

## Frontmatter Reference

| Field | Required | Description |
|-------|----------|-------------|
| `name` | No | Display name, defaults to directory name |
| `description` | Recommended | When to use this skill |
| `argument-hint` | No | Hint for autocomplete, e.g. `[issue-number]` |
| `disable-model-invocation` | No | `true` = only user can invoke |
| `user-invocable` | No | `false` = hidden from `/` menu |
| `allowed-tools` | No | Tools Claude can use without permission |
| `model` | No | Model to use when skill active |
| `context` | No | `fork` = run in subagent context |
| `agent` | No | Subagent type when `context: fork` |
| `hooks` | No | Lifecycle hooks scoped to skill |

---

## String Substitutions

| Variable | Description |
|----------|-------------|
| `$ARGUMENTS` | All arguments passed when invoking |
| `$ARGUMENTS[N]` | Specific argument by 0-based index |
| `$N` | Shorthand for `$ARGUMENTS[N]` |
| `${CLAUDE_SESSION_ID}` | Current session ID |

---

## Skill Types

### Reference Content (Knowledge)

```yaml
---
name: api-conventions
description: API design patterns for this codebase
---

When writing API endpoints:
- Use RESTful naming conventions
- Return consistent error formats
- Include request validation
```

### Task Content (Actions)

```yaml
---
name: deploy
description: Deploy the application to production
context: fork
disable-model-invocation: true
---

Deploy the application:
1. Run the test suite
2. Build the application
3. Push to the deployment target
```

---

## Control Who Invokes

| Frontmatter | User Can Invoke | Claude Can Invoke |
|-------------|-----------------|-------------------|
| (default) | Yes | Yes |
| `disable-model-invocation: true` | Yes | No |
| `user-invocable: false` | No | Yes |

---

## Supporting Files

```
my-skill/
├── SKILL.md           # Main instructions (required)
├── template.md        # Template for Claude to fill in
├── examples/
│   └── sample.md      # Example output
└── scripts/
    └── validate.sh    # Script Claude can execute
```

Reference from SKILL.md:
```markdown
## Additional resources
- For complete API details, see [reference.md](reference.md)
- For usage examples, see [examples.md](examples.md)
```

---

## Dynamic Context Injection

Use `!`command`` to run shell commands before skill content is sent:

```yaml
---
name: pr-summary
description: Summarize changes in a pull request
context: fork
agent: Explore
---

## Pull request context
- PR diff: !`gh pr diff`
- PR comments: !`gh pr view --comments`
- Changed files: !`gh pr diff --name-only`

## Your task
Summarize this pull request...
```

---

## Run in Subagent

Add `context: fork` for isolated execution:

```yaml
---
name: deep-research
description: Research a topic thoroughly
context: fork
agent: Explore
---

Research $ARGUMENTS thoroughly:
1. Find relevant files using Glob and Grep
2. Read and analyze the code
3. Summarize findings with specific file references
```

---

## Example: Fix GitHub Issue

```yaml
---
name: fix-issue
description: Fix a GitHub issue
disable-model-invocation: true
---

Fix GitHub issue $ARGUMENTS following our coding standards.

1. Read the issue description
2. Understand the requirements
3. Implement the fix
4. Write tests
5. Create a commit
```

Usage: `/fix-issue 123`

---

## Example: Codebase Visualizer

Generate interactive HTML tree view:

```yaml
---
name: codebase-visualizer
description: Generate an interactive tree visualization of your codebase.
allowed-tools: Bash(python *)
---

Run the visualization script from your project root:

```bash
python ~/.claude/skills/codebase-visualizer/scripts/visualize.py .
```
```

---

## Permission Control

Disable all skills:
```
# Add to deny rules:
Skill
```

Allow/deny specific skills:
```
Skill(commit)           # Allow commit skill
Skill(review-pr *)      # Allow review-pr with any args
Skill(deploy *)         # Deny deploy skill
```
