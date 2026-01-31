# Claude Code Best Practices

> Source: https://code.claude.com/docs/en/best-practices
> Fetched: 2026-01-31
> Type: html

---

## Core Constraint

**Claude's context window fills up fast, and performance degrades as it fills.**

The context window holds your entire conversation, including every message, every file Claude reads, and every command output. LLM performance degrades as context fills - Claude may start "forgetting" earlier instructions or making more mistakes.

---

## Give Claude a Way to Verify Its Work

**The single highest-leverage thing you can do.**

Claude performs dramatically better when it can verify its own work: run tests, compare screenshots, validate outputs.

| Strategy | Before | After |
|----------|--------|-------|
| **Provide verification criteria** | "implement a function that validates email addresses" | "write a validateEmail function. example test cases: user@example.com is true, invalid is false, user@.com is false. run the tests after implementing" |
| **Verify UI changes visually** | "make the dashboard look better" | "[paste screenshot] implement this design. take a screenshot of the result and compare it to the original. list differences and fix them" |
| **Address root causes** | "the build is failing" | "the build fails with this error: [paste error]. fix it and verify the build succeeds. address the root cause, don't suppress the error" |

---

## Explore First, Then Plan, Then Code

Separate research and planning from implementation.

### Four-Phase Workflow

1. **Explore**: Enter Plan Mode. Claude reads files and answers questions without making changes.
2. **Plan**: Ask Claude to create a detailed implementation plan. Press `Ctrl+G` to open in editor.
3. **Implement**: Switch to Normal Mode. Let Claude code, verifying against its plan.
4. **Commit**: Ask Claude to commit with a descriptive message and create a PR.

**Skip planning** for tasks where scope is clear and fix is small. Planning is most useful when uncertain about approach or when modifying multiple files.

---

## Provide Specific Context in Your Prompts

| Strategy | Before | After |
|----------|--------|-------|
| **Scope the task** | "add tests for foo.py" | "write a test for foo.py covering the edge case where the user is logged out. avoid mocks." |
| **Point to sources** | "why does ExecutionFactory have such a weird api?" | "look through ExecutionFactory's git history and summarize how its api came to be" |
| **Reference patterns** | "add a calendar widget" | "look at how existing widgets are implemented on the home page. HotDogWidget.php is a good example. follow the pattern..." |
| **Describe symptoms** | "fix the login bug" | "users report that login fails after session timeout. check the auth flow in src/auth/, especially token refresh..." |

### Provide Rich Content

- **Reference files with `@`** instead of describing where code lives
- **Paste images directly** - copy/paste or drag and drop
- **Give URLs** for documentation and API references
- **Pipe in data** - `cat error.log | claude`
- **Let Claude fetch** - tell Claude to pull context itself using Bash, MCP, or file reads

---

## Configure Your Environment

### Write an Effective CLAUDE.md

Run `/init` to generate a starter file, then refine.

```markdown
# Code style
- Use ES modules (import/export) syntax, not CommonJS (require)
- Destructure imports when possible

# Workflow
- Be sure to typecheck when you're done making a series of code changes
- Prefer running single tests, not the whole test suite
```

**Include:**
- Bash commands Claude can't guess
- Code style rules that differ from defaults
- Testing instructions and preferred test runners
- Repository etiquette (branch naming, PR conventions)
- Architectural decisions specific to your project
- Common gotchas or non-obvious behaviors

**Exclude:**
- Anything Claude can figure out by reading code
- Standard language conventions Claude already knows
- Detailed API documentation (link to docs instead)
- Information that changes frequently
- File-by-file descriptions of the codebase
- Self-evident practices like "write clean code"

### CLAUDE.md Locations

| Location | Purpose |
|----------|---------|
| `~/.claude/CLAUDE.md` | Applies to all Claude sessions |
| `./CLAUDE.md` | Check into git to share with team |
| `CLAUDE.local.md` | Personal overrides, .gitignore it |
| Parent directories | Useful for monorepos |
| Child directories | Pulled in on demand |

---

## Configure Permissions

Use `/permissions` to allowlist safe commands or `/sandbox` for OS-level isolation.

---

## Use CLI Tools

Tell Claude Code to use CLI tools like `gh`, `aws`, `gcloud`, `sentry-cli` for external services. CLI tools are the most context-efficient way to interact with external services.

---

## Connect MCP Servers

Run `claude mcp add` to connect external tools like Notion, Figma, or databases.

---

## Set Up Hooks

Use hooks for actions that must happen every time with zero exceptions.

Example: "Write a hook that runs eslint after every file edit"

---

## Create Skills

Create `SKILL.md` files in `.claude/skills/` for domain knowledge and reusable workflows.

```markdown
---
name: api-conventions
description: REST API design conventions for our services
---
# API Conventions
- Use kebab-case for URL paths
- Use camelCase for JSON properties
- Always include pagination for list endpoints
```

---

## Create Custom Subagents

Define specialized assistants in `.claude/agents/` that Claude can delegate to.

```markdown
---
name: security-reviewer
description: Reviews code for security vulnerabilities
tools: Read, Grep, Glob, Bash
model: opus
---
You are a senior security engineer. Review code for:
- Injection vulnerabilities
- Authentication and authorization flaws
- Secrets or credentials in code
```

---

## Communicate Effectively

### Ask Codebase Questions

Ask Claude questions you'd ask a senior engineer:
- How does logging work?
- How do I make a new API endpoint?
- What edge cases does `CustomerOnboardingFlowImpl` handle?

### Let Claude Interview You

For larger features, have Claude interview you first:

```
I want to build [brief description]. Interview me in detail using the AskUserQuestion tool.

Ask about technical implementation, UI/UX, edge cases, concerns, and tradeoffs.
Keep interviewing until we've covered everything, then write a complete spec to SPEC.md.
```

---

## Manage Your Session

### Course-Correct Early and Often

- **`Esc`**: Stop Claude mid-action
- **`Esc + Esc` or `/rewind`**: Open rewind menu and restore previous state
- **`"Undo that"`**: Have Claude revert changes
- **`/clear`**: Reset context between unrelated tasks

### Manage Context Aggressively

- Use `/clear` frequently between tasks
- Auto compaction triggers when context limits approached
- Run `/compact <instructions>` for more control
- Customize compaction behavior in CLAUDE.md

### Use Subagents for Investigation

Subagents run in separate context windows and report back summaries:

```
Use subagents to investigate how our authentication system handles token
refresh, and whether we have any existing OAuth utilities I should reuse.
```

### Rewind with Checkpoints

Double-tap `Escape` or run `/rewind` to open checkpoint menu. Restore conversation only, code only, or both.

### Resume Conversations

```bash
claude --continue    # Resume most recent
claude --resume      # Select from recent
```

Use `/rename` to give sessions descriptive names.

---

## Automate and Scale

### Run Headless Mode

```bash
# One-off queries
claude -p "Explain what this project does"

# Structured output
claude -p "List all API endpoints" --output-format json

# Streaming for real-time processing
claude -p "Analyze this log file" --output-format stream-json
```

### Run Multiple Claude Sessions

- **Claude Desktop**: Manage multiple local sessions visually
- **Claude Code on web**: Run on Anthropic's cloud in isolated VMs

### Writer/Reviewer Pattern

| Session A (Writer) | Session B (Reviewer) |
|--------------------|---------------------|
| `Implement a rate limiter for our API endpoints` | |
| | `Review the rate limiter implementation. Look for edge cases, race conditions...` |
| `Address these issues: [Session B output]` | |

### Fan Out Across Files

```bash
for file in $(cat files.txt); do
  claude -p "Migrate $file from React to Vue. Return OK or FAIL." \
    --allowedTools "Edit,Bash(git commit *)"
done
```

---

## Avoid Common Failure Patterns

| Pattern | Problem | Fix |
|---------|---------|-----|
| **Kitchen sink session** | Unrelated tasks polluting context | `/clear` between unrelated tasks |
| **Correcting over and over** | Context polluted with failed approaches | After two corrections, `/clear` and better prompt |
| **Over-specified CLAUDE.md** | Claude ignores instructions | Ruthlessly prune |
| **Trust-then-verify gap** | Plausible code without edge case handling | Always provide verification |
| **Infinite exploration** | Claude reads hundreds of files | Scope investigations or use subagents |

---

## Develop Your Intuition

Pay attention to what works:
- Prompt structure that produces great output
- Context you provided
- Mode you were in

When Claude struggles, ask why:
- Context too noisy?
- Prompt too vague?
- Task too big for one pass?

Over time, you'll know when to be specific vs. open-ended, when to plan vs. explore, when to clear context vs. let it accumulate.
