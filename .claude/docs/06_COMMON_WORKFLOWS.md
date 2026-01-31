# Common Workflows

> Source: https://code.claude.com/docs/en/common-workflows
> Fetched: 2026-01-31
> Type: html

---

## Understand New Codebases

### Get a Quick Overview

```
> give me an overview of this codebase
> explain the main architecture patterns used here
> what are the key data models?
> how is authentication handled?
```

### Find Relevant Code

```
> find the files that handle user authentication
> how do these authentication files work together?
> trace the login process from front-end to database
```

---

## Fix Bugs Efficiently

```
> I'm seeing an error when I run npm test
> suggest a few ways to fix the @ts-ignore in user.ts
> update user.ts to add the null check you suggested
```

---

## Refactor Code

```
> find deprecated API usage in our codebase
> suggest how to refactor utils.js to use modern JavaScript features
> refactor utils.js to use ES2024 features while maintaining the same behavior
> run tests for the refactored code
```

---

## Use Specialized Subagents

```
> /agents                                    # View available subagents
> review my recent code changes for security issues
> use the code-reviewer subagent to check the auth module
> have the debugger subagent investigate why users can't log in
```

---

## Plan Mode for Safe Code Analysis

**When to Use:**
- Multi-step implementation
- Code exploration before changes
- Interactive development with direction iteration

**How to Use:**
- `Shift+Tab` to cycle through permission modes
- `claude --permission-mode plan` to start in plan mode
- `Ctrl+G` to open plan in text editor

```bash
claude --permission-mode plan -p "Analyze the authentication system and suggest improvements"
```

---

## Work with Tests

```
> find functions in NotificationsService.swift that are not covered by tests
> add tests for the notification service
> add test cases for edge conditions in the notification service
> run the new tests and fix any failures
```

---

## Create Pull Requests

```
> /commit-push-pr                           # Commits, pushes, and opens PR
> summarize the changes I've made to the authentication module
> create a pr
> enhance the PR description with more context about the security improvements
```

---

## Handle Documentation

```
> find functions without proper JSDoc comments in the auth module
> add JSDoc comments to the undocumented functions in auth.js
> improve the generated documentation with more context and examples
> check if the documentation follows our project standards
```

---

## Work with Images

Methods to add images:
1. Drag and drop into Claude Code window
2. Copy and paste with `Ctrl+V`
3. Provide image path: `"Analyze this image: /path/to/your/image.png"`

```
> What does this image show?
> Describe the UI elements in this screenshot
> Here's a screenshot of the error. What's causing it?
> Generate CSS to match this design mockup
```

---

## Reference Files and Directories

```
> Explain the logic in @src/utils/auth.js         # Reference single file
> What's the structure of @src/components?         # Reference directory
> Show me the data from @github:repos/owner/repo/issues  # MCP resources
```

---

## Extended Thinking

Enabled by default, reserves up to 31,999 tokens for step-by-step reasoning.

**Configure:**
- `Option+T` / `Alt+T` - Toggle for current session
- `/config` - Set global default
- `MAX_THINKING_TOKENS` env var - Limit budget

**View thinking:** `Ctrl+O` to toggle verbose mode

---

## Resume Previous Conversations

```bash
claude --continue    # Resume most recent
claude --resume      # Select from recent conversations
/rename auth-refactor  # Name current session
```

**Session Picker Shortcuts:**
| Key | Action |
|-----|--------|
| `↑/↓` | Navigate sessions |
| `→/←` | Expand/collapse grouped |
| `Enter` | Resume highlighted |
| `P` | Preview session |
| `R` | Rename session |
| `/` | Search filter |
| `A` | Toggle all projects |
| `B` | Filter by current branch |

---

## Parallel Sessions with Git Worktrees

```bash
# Create worktree with new branch
git worktree add ../project-feature-a -b feature-a

# Create worktree with existing branch
git worktree add ../project-bugfix bugfix-123

# Run Claude in each worktree
cd ../project-feature-a && claude
cd ../project-bugfix && claude

# Manage worktrees
git worktree list
git worktree remove ../project-feature-a
```

---

## Claude as Unix Utility

### Build Script Integration

```json
{
  "scripts": {
    "lint:claude": "claude -p 'you are a linter. look at changes vs main and report issues related to typos...'"
  }
}
```

### Pipe Data

```bash
cat build-error.txt | claude -p 'concisely explain the root cause' > output.txt
```

### Output Formats

```bash
claude -p 'summarize' --output-format text     # Plain text (default)
claude -p 'analyze' --output-format json       # JSON array with metadata
claude -p 'parse' --output-format stream-json  # Real-time streaming JSON
```

---

## Ask About Capabilities

```
> can Claude Code create pull requests?
> how does Claude Code handle permissions?
> what skills are available?
> how do I use MCP with Claude Code?
> how do I configure Claude Code for Amazon Bedrock?
```
