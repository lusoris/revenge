# Claude Code Quickstart

> Source: https://code.claude.com/docs/en/quickstart
> Fetched: 2026-01-31
> Type: html

---

## Prerequisites

- Terminal or command prompt
- A code project to work with
- [Claude subscription](https://claude.com/pricing) (Pro, Max, Teams, Enterprise) or [Console account](https://console.anthropic.com/)

---

## Step 1: Install

**macOS, Linux, WSL:**
```bash
curl -fsSL https://claude.ai/install.sh | bash
```

**Windows PowerShell:**
```powershell
irm https://claude.ai/install.ps1 | iex
```

**Homebrew:**
```sh
brew install --cask claude-code
```

**WinGet:**
```powershell
winget install Anthropic.ClaudeCode
```

---

## Step 2: Log In

```bash
claude
# Follow prompts to log in
```

Or use `/login` command within a session.

Account types:
- Claude Pro, Max, Teams, Enterprise (recommended)
- Claude Console (API with pre-paid credits)
- Amazon Bedrock, Google Vertex AI, Microsoft Foundry

---

## Step 3: Start Session

```bash
cd /path/to/your/project
claude
```

Commands:
- `/help` - Available commands
- `/resume` - Continue previous conversation

---

## Step 4: Ask Questions

```
what does this project do?
what technologies does this project use?
where is the main entry point?
explain the folder structure
what can Claude Code do?
```

---

## Step 5: Make Code Changes

```
add a hello world function to the main file
```

Claude will:
1. Find the appropriate file
2. Show proposed changes
3. Ask for approval
4. Make the edit

---

## Step 6: Use Git

```
what files have I changed?
commit my changes with a descriptive message
create a new branch called feature/quickstart
show me the last 5 commits
help me resolve merge conflicts
```

---

## Step 7: Fix Bugs / Add Features

```
add input validation to the user registration form
there's a bug where users can submit empty forms - fix it
```

Claude will:
- Locate relevant code
- Understand context
- Implement solution
- Run tests if available

---

## Step 8: Common Workflows

**Refactor:**
```
refactor the authentication module to use async/await instead of callbacks
```

**Write tests:**
```
write unit tests for the calculator functions
```

**Update docs:**
```
update the README with installation instructions
```

**Code review:**
```
review my changes and suggest improvements
```

---

## Essential Commands

| Command | Description | Example |
|---------|-------------|---------|
| `claude` | Start interactive mode | `claude` |
| `claude "task"` | Run one-time task | `claude "fix the build error"` |
| `claude -p "query"` | Query and exit | `claude -p "explain this function"` |
| `claude -c` | Continue recent conversation | `claude -c` |
| `claude -r` | Resume previous conversation | `claude -r` |
| `claude commit` | Create Git commit | `claude commit` |
| `/clear` | Clear history | `/clear` |
| `/help` | Show commands | `/help` |
| `exit` or `Ctrl+C` | Exit | `exit` |

---

## Pro Tips

### Be Specific
Instead of: `fix the bug`
Try: `fix the login bug where users see a blank screen after entering wrong credentials`

### Use Step-by-Step
```
1. create a new database table for user profiles
2. create an API endpoint to get and update user profiles
3. build a webpage that allows users to see and edit their information
```

### Let Claude Explore First
```
analyze the database schema
```

### Shortcuts
- `?` - See keyboard shortcuts
- `Tab` - Command completion
- `↑` - Command history
- `/` - See all commands and skills
