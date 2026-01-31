# Claude Code Memory Management

> Source: https://code.claude.com/docs/en/memory
> Fetched: 2026-01-31
> Type: html

---

## Memory Types (Hierarchy)

| Memory Type | Location | Purpose | Shared With |
|-------------|----------|---------|-------------|
| **Managed policy** | `/etc/claude-code/CLAUDE.md` (Linux) | Organization-wide | All users |
| **Project memory** | `./CLAUDE.md` or `./.claude/CLAUDE.md` | Team-shared | Team via git |
| **Project rules** | `./.claude/rules/*.md` | Modular topic-specific | Team via git |
| **User memory** | `~/.claude/CLAUDE.md` | Personal preferences | Just you (all projects) |
| **Project local** | `./CLAUDE.local.md` | Personal project-specific | Just you (auto-ignored) |

Higher priority overrides lower. All loaded automatically at session start.

---

## CLAUDE.md Imports

Import files using `@path/to/import` syntax:

```markdown
See @README for project overview and @package.json for npm commands.

# Additional Instructions
- git workflow @docs/git-instructions.md
```

- Relative and absolute paths allowed
- Import from home dir: `@~/.claude/my-project-instructions.md`
- Not evaluated inside code spans/blocks
- Max import depth: 5 hops
- Check loaded files with `/memory`

---

## Memory Discovery

Claude Code reads memories recursively from cwd up to root.

**Load from additional directories:**
```bash
CLAUDE_CODE_ADDITIONAL_DIRECTORIES_CLAUDE_MD=1 claude --add-dir ../shared-config
```

---

## Modular Rules with `.claude/rules/`

```
.claude/rules/
├── code-style.md     # Code style guidelines
├── testing.md        # Testing conventions
├── security.md       # Security requirements
└── frontend/
    ├── react.md
    └── styles.md
```

### Path-Specific Rules

```markdown
---
paths:
  - "src/api/**/*.ts"
---

# API Development Rules
- All endpoints must include input validation
```

### Glob Patterns

| Pattern | Matches |
|---------|---------|
| `**/*.ts` | All TypeScript files |
| `src/**/*` | All files under src/ |
| `*.md` | Markdown in root |
| `src/**/*.{ts,tsx}` | Both .ts and .tsx |

---

## Best Practices

- **Be specific**: "Use 2-space indentation" > "Format code properly"
- **Use structure**: Bullet points, markdown headings
- **Review periodically**: Update as project evolves
