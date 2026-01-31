# Tool Comparison Guide

**Purpose**: Help you choose the right tool for each task in the Revenge project

**Last Updated**: 2026-01-31

---

## Overview

This guide provides comprehensive comparisons of all development tools used in the Revenge project, helping you make informed decisions about which tool to use for specific tasks.

---

## Table of Contents

- [VS Code vs Zed - Detailed Comparison](#vs-code-vs-zed---detailed-comparison)
- [Local vs Remote Development](#local-vs-remote-development)
- [IDE Selection Matrix](#ide-selection-matrix)
- [Git Tools Comparison](#git-tools-comparison)
- [Real-World Scenarios](#real-world-scenarios)
- [Decision Trees](#decision-trees)

---

## VS Code vs Zed - Detailed Comparison

### Feature Matrix

| Feature | VS Code | Zed | Winner |
|---------|---------|-----|--------|
| **Performance** | | | |
| Startup time | 2-5 seconds | <1 second | 🟢 Zed |
| Memory usage | 300-800 MB | 100-300 MB | 🟢 Zed |
| File search speed | Fast | Very fast | 🟢 Zed |
| LSP responsiveness | Good | Excellent | 🟢 Zed |
| **Language Support** | | | |
| Go | Excellent | Excellent | 🟡 Tie |
| Python | Excellent | Good | 🟢 VS Code |
| TypeScript | Excellent | Good | 🟢 VS Code |
| Svelte | Excellent | Basic | 🟢 VS Code |
| **Debugging** | | | |
| Graphical debugger | Yes | No (terminal) | 🟢 VS Code |
| Breakpoints UI | Excellent | N/A | 🟢 VS Code |
| Variable inspection | Visual | Terminal | 🟢 VS Code |
| **Extensions** | | | |
| Extension marketplace | 50,000+ | Limited | 🟢 VS Code |
| Custom themes | Thousands | ~20 | 🟢 VS Code |
| Language support | Any | Built-in only | 🟢 VS Code |
| **Collaboration** | | | |
| Live Share | Mature | Beta | 🟢 VS Code |
| Remote SSH | Excellent | Native | 🟡 Tie |
| Remote Containers | Yes | No | 🟢 VS Code |
| **Git Integration** | | | |
| Basic git UI | Yes | Yes | 🟡 Tie |
| GitLens | Yes | No | 🟢 VS Code |
| Inline blame | Via extension | Built-in | 🟡 Tie |
| Visual diff | Excellent | Good | 🟢 VS Code |
| **Configuration** | | | |
| Settings UI | Graphical | JSON only | 🟢 VS Code |
| Keybindings | Graphical | JSON only | 🟢 VS Code |
| Workspace settings | Yes | Yes | 🟡 Tie |
| **Resource Usage** | | | |
| Electron-based | Yes | No (native) | 🟢 Zed |
| Battery impact | Higher | Lower | 🟢 Zed |
| CPU usage (idle) | 2-5% | <1% | 🟢 Zed |

### Detailed Comparison

#### Performance

**VS Code**:
- Electron-based (Chromium + Node.js)
- Startup: 2-5 seconds for typical project
- Memory: 300-800 MB depending on extensions
- File search: Fast with ripgrep backend
- LSP: Good performance, occasional lag with large codebases

**Zed**:
- Native Rust application (GPU-accelerated)
- Startup: <1 second consistently
- Memory: 100-300 MB typical usage
- File search: Very fast (native implementation)
- LSP: Excellent performance, minimal lag

**Verdict**: Zed wins on raw performance, especially for quick edits and large files.

#### Language Support

**Go**:
- **VS Code**: Excellent via golang.go extension + gopls
  - Full debugging support
  - Test explorer
  - Code lens
  - Rich IntelliSense
- **Zed**: Excellent via built-in gopls support
  - All LSP features work
  - Inlay hints
  - No graphical debugger (use delve CLI)
- **Winner**: Tie (both excellent, VS Code has better debugging UI)

**Python**:
- **VS Code**: Excellent via Pylance + Python extension
  - Type checking
  - Jupyter notebook support
  - Debugging
  - Virtual environment detection
- **Zed**: Good via ruff-lsp
  - Linting and formatting
  - Basic IntelliSense
  - No Jupyter support
  - Terminal debugging only
- **Winner**: VS Code (Pylance is superior LSP, Jupyter support)

**TypeScript/JavaScript**:
- **VS Code**: Excellent (TypeScript is a Microsoft project)
  - Best-in-class IntelliSense
  - Refactoring tools
  - Built-in support
- **Zed**: Good via typescript-language-server
  - Standard LSP features work
  - Formatting via prettier
  - Less mature refactoring
- **Winner**: VS Code (superior integration)

**Svelte**:
- **VS Code**: Excellent via svelte.svelte-vscode
  - Official Svelte team extension
  - Svelte 5 runes support
  - Component IntelliSense
- **Zed**: Basic support
  - Syntax highlighting
  - Must use external prettier for formatting
  - No component-specific features
- **Winner**: VS Code (only real option for Svelte)

#### Debugging

**VS Code**:
- Full graphical debugger
- Launch configurations in `.vscode/launch.json`
- Breakpoints UI with conditional breakpoints
- Variable inspection in sidebar
- Debug console with REPL
- Call stack visualization
- Watch expressions
- **Supported**: Go (delve), Python, Node.js, many more

**Zed**:
- No built-in graphical debugger
- Must use terminal-based debuggers:
  - Go: `dlv debug` or `dlv attach`
  - Python: `python -m pdb`
  - Node: `node --inspect`
- Can attach external debuggers
- Community working on debug adapter protocol support

**Verdict**: VS Code is clear winner for debugging. If you need to debug complex issues, use VS Code.

#### Extensions & Customization

**VS Code**:
- 50,000+ extensions in marketplace
- Extensions for every language, framework, tool
- Themes, icon packs, snippets
- Can modify nearly every aspect
- Extension API is well-documented
- Community is massive

**Zed**:
- Limited extension system (still evolving)
- ~20 official themes
- Language support is mostly built-in
- Fewer customization options
- Simpler = less to configure
- Community is growing fast

**Verdict**: VS Code wins if you need specific extensions. Zed wins if you want simplicity.

#### Collaboration

**VS Code Live Share**:
- Mature, production-ready
- Share terminal, servers, debugging sessions
- Voice/audio support
- Follow participants
- Shared editing with cursors

**Zed Collaboration** (Beta):
- Built-in from ground up (not extension)
- Very fast (CRDT-based)
- Shared cursors and selections
- Voice channels
- Still in beta, occasional bugs

**Verdict**: VS Code Live Share is more mature, but Zed's is promising and faster when it works.

#### Remote Development

**VS Code**:
- Remote-SSH: Excellent
- Remote-Containers: Full Docker dev container support
- Code Server: Browser-based VS Code
- GitHub Codespaces: First-class support

**Zed**:
- Remote-SSH: Native, works well
- No container support
- Fast over SSH (native app)
- Lower bandwidth usage than VS Code

**Verdict**: VS Code has more remote options, but Zed is faster for pure SSH.

### When to Use Each

#### Use VS Code When:

1. **Working with Svelte/TypeScript** - Only real option for Svelte 5
2. **Debugging complex issues** - Graphical debugger is invaluable
3. **Need specific extensions** - Database tools, API clients, etc.
4. **Jupyter notebooks** - Python data science work
5. **Docker development** - Dev containers, Docker extension
6. **Pair programming** - Live Share is mature
7. **Learning** - More docs, tutorials, community support

#### Use Zed When:

1. **Quick edits** - Startup time is instant
2. **Documentation writing** - Fast, distraction-free Markdown editing
3. **Go backend work** - gopls works great, debugging with delve CLI is fine
4. **Large files** - Better performance with big files
5. **Battery life matters** - Laptop on the go
6. **Prefer simplicity** - Fewer extensions = less config
7. **Remote SSH** - Faster than VS Code Remote-SSH

### Migration Guide: VS Code → Zed

If you're coming from VS Code, here's what to expect:

#### What Works Identically

- Most keybindings (Cmd/Ctrl+P, Cmd/Ctrl+Shift+P, etc.)
- LSP features (Go to Definition, Find References, Rename)
- Multi-cursor editing
- Terminal integration
- Git integration basics
- File explorer
- Format on save

#### What Works Differently

| Feature | VS Code | Zed | Workaround |
|---------|---------|-----|------------|
| Debugging | Graphical UI | Terminal only | Learn `dlv` or `pdb` CLI |
| Extensions | Marketplace | Built-in | Use native features or external tools |
| Settings | JSON + UI | JSON only | Edit `settings.json` directly |
| Themes | Thousands | ~20 | Pick from available or make custom |
| Tasks | `tasks.json` | Terminal | Run commands directly |
| Snippets | Extension | Limited | Use external snippet manager |

#### What's Missing

- Jupyter notebook support
- Docker/Container development
- Full extension ecosystem
- Settings UI
- Many language-specific features

#### Recommended Approach

**Hybrid workflow**:
- Use Zed for: Go backend, docs, quick edits, config files
- Use VS Code for: Svelte frontend, debugging, Docker, specific extensions

---

## Local vs Remote Development

### Performance Comparison

| Metric | Local | Remote (Coder) |
|--------|-------|----------------|
| **File operations** | Instant | Network latency (~10-50ms) |
| **LSP responsiveness** | Instant | Slight lag (~10-50ms) |
| **Build speed** | Depends on laptop | Depends on server (usually faster) |
| **Test execution** | Depends on laptop | Depends on server |
| **HMR/Hot Reload** | Very fast | Slower (network overhead) |
| **Git operations** | Fast | Fast (server has better bandwidth) |
| **Docker builds** | Depends on laptop | Usually faster (more CPU/RAM) |
| **Database queries** | localhost | localhost (on server) |

### Cost Comparison

#### Local Development

**Hardware Costs**:
- Development laptop: $1,500 - $3,000
- RAM upgrade (32GB+): $200 - $500
- SSD upgrade (1TB NVMe): $100 - $200
- **Total upfront**: $1,800 - $3,700

**Ongoing Costs**:
- Electricity: ~$10/month (8 hours/day usage)
- Hardware refresh: ~$2,000 every 3-4 years ($50/month amortized)
- **Total ongoing**: ~$60/month

**Pros**:
- No network dependency
- Full control
- Works offline
- No ongoing costs after hardware purchase
- Data stays local

**Cons**:
- High upfront cost
- Limited by laptop specs
- Battery drain on laptops
- Can't access from other devices
- Must maintain local environment

#### Remote Development (Coder)

**Infrastructure Costs** (self-hosted):
- Server hardware: $500 - $2,000 (or cloud VM)
- Cloud VM (8 CPU, 16GB RAM): $50 - $150/month
- Storage: $10 - $50/month
- Bandwidth: Usually included
- **Total**: $60 - $200/month

**Pros**:
- Access from any device (tablet, laptop, desktop)
- Powerful server specs
- No battery drain on client device
- Easy to scale resources
- Team can share infrastructure
- Consistent environment across team
- Work from anywhere

**Cons**:
- Network dependency
- Monthly costs
- Latency for LSP/file operations
- Requires infrastructure management
- Can't work offline

### Decision Matrix

| Factor | Local | Remote (Coder) | Winner |
|--------|-------|----------------|--------|
| **Setup time** | 30-60 min | 5 min | 🟢 Remote |
| **Performance (Backend)** | Depends on machine | Consistent | 🟢 Remote |
| **Performance (Frontend HMR)** | Fast | Slower over SSH | 🟢 Local |
| **Resource usage** | Local CPU/RAM | Server CPU/RAM | 🟢 Remote |
| **Access from anywhere** | No | Yes | 🟢 Remote |
| **Works offline** | Yes | No | 🟢 Local |
| **Consistent environment** | No (varies) | Yes (template) | 🟢 Remote |
| **Team collaboration** | Manual sync | Built-in | 🟢 Remote |
| **Initial cost** | High | Low | 🟢 Remote |
| **Ongoing cost** | Low | Monthly | 🟢 Local |

### When to Use Each

#### Use Local Development When:

1. **You work offline frequently** - Flights, commutes, remote areas
2. **You have a powerful laptop** - No need for remote resources
3. **You value instant responsiveness** - No tolerance for network latency
4. **You work with large binary files** - Video editing, game assets
5. **You prefer simplicity** - No infrastructure to manage
6. **Your budget is limited** - No ongoing costs after hardware purchase
7. **Frontend development** - HMR works best locally

#### Use Remote Development When:

1. **You work from multiple devices** - Desktop at home, laptop on the go
2. **Your laptop is low-spec** - Chromebook, older machine, budget laptop
3. **You need consistent environments** - Entire team uses same setup
4. **You run resource-intensive tasks** - Large builds, parallel tests
5. **You value team collaboration** - Shared workspaces, pair programming
6. **You have good internet** - Low latency, high bandwidth
7. **Backend development** - Full stack available remotely

### Hybrid Approach (Recommended)

**Best of both worlds**:

```
Local Machine:
├── VS Code (primary for frontend)
├── Zed (quick edits)
├── Frontend repo clone
└── Light testing

Remote Workspace (Coder):
├── VS Code (browser or SSH)
├── Full development stack
├── PostgreSQL, Dragonfly, Typesense
├── Backend builds and tests
└── Integration testing
```

**Workflow**:
1. **Frontend dev**: Local (fast HMR with `npm run dev`)
2. **Backend dev**: Remote (full stack, powerful resources)
3. **Integration tests**: Remote (requires database, cache, search)
4. **Documentation**: Local (Zed for speed, works offline)
5. **Quick fixes**: Local (instant startup)
6. **Large refactors**: Remote (more RAM, parallel tools)

---

## IDE Selection Matrix

### By Task Type

| Task | Primary Tool | Alternative | Rationale |
|------|-------------|-------------|-----------|
| **Writing Go code** | VS Code | Zed | VS Code has better debugging, gopls works great in both |
| **Writing Python scripts** | VS Code | Zed | VS Code has Pylance and Jupyter support |
| **Writing Svelte/TypeScript** | VS Code | - | Svelte extension is VS Code-only |
| **Debugging Go** | VS Code | Terminal (delve) | Graphical debugger is superior for complex issues |
| **Debugging Python** | VS Code | Terminal (pdb) | Graphical debugger shows variables inline |
| **Writing documentation** | Zed | VS Code | Fast, distraction-free Markdown editing |
| **Editing config files** | Zed | VS Code | Quick edits don't need full IDE |
| **Reviewing PRs** | GitHub Web | VS Code | GitHub UI is designed for reviews |
| **Pair programming** | VS Code (LiveShare) | Zed (beta) | More mature collaboration features |
| **Database queries** | VS Code (extension) | Terminal (psql) | SQL extensions provide autocomplete |
| **Docker/K8s work** | VS Code | Terminal | Docker extension visualizes containers |
| **Git operations** | VS Code (GitLens) | Terminal | Graphical git history is useful |
| **Running tests** | VS Code | Terminal | VS Code can run tests in debug mode |
| **Refactoring** | VS Code | Zed | Better LSP integration and refactoring tools |
| **Quick fixes** | Zed | VS Code | Faster startup for single-file edits |

### By Language

| Language | Primary Tool | Extensions/LSP | Notes |
|----------|-------------|----------------|-------|
| **Go** | VS Code | golang.go (gopls) | Excellent debugging, testing, and LSP support |
| **Python** | VS Code | Pylance, Ruff | Pylance is best Python LSP, Ruff for linting |
| **TypeScript** | VS Code | Built-in | TypeScript originated from Microsoft |
| **Svelte** | VS Code | svelte.svelte-vscode | Official extension, no Zed support |
| **JavaScript** | VS Code | Built-in | Same as TypeScript |
| **HTML/CSS** | VS Code | Built-in | Emmet, autocomplete |
| **SQL** | VS Code | PostgreSQL extension | Query execution, autocomplete |
| **Markdown** | Zed | - | Fast, lightweight, good preview |
| **YAML/JSON** | Zed | - | Schema validation, fast editing |
| **Shell scripts** | Zed | - | Lightweight, fast syntax highlighting |
| **Rust** | Zed or VS Code | rust-analyzer | Zed has first-class Rust support (written in Rust) |
| **Dockerfile** | VS Code | Docker extension | Build, run, debug containers |
| **Terraform** | VS Code | HashiCorp extension | Syntax, validation, completion |

---

## Git Tools Comparison

### Command Line vs GitHub CLI vs GUI Tools

| Feature | Git CLI | GitHub CLI (`gh`) | VS Code (GitLens) | Winner |
|---------|---------|-------------------|-------------------|--------|
| **Speed** | Fastest | Fast | Medium | 🟢 Git CLI |
| **Visualize history** | No (text) | No | Yes | 🟢 VS Code |
| **Branching** | Manual | Manual | Point-and-click | 🟢 VS Code |
| **Merging** | Manual | Manual | UI-guided | 🟢 VS Code |
| **PR creation** | No | Yes | Yes | 🟢 gh CLI |
| **Conflict resolution** | Manual | - | Visual merge tool | 🟢 VS Code |
| **Learning curve** | Steep | Medium | Easy | 🟢 VS Code |
| **Scripting** | Excellent | Good | No | 🟢 Git CLI |
| **Offline** | Yes | No | Yes (local ops) | 🟢 Git CLI |

### When to Use Each

#### Git CLI (Terminal)

**Use When**:
- Scripting or automation
- Complex rebase/cherry-pick operations
- You know exactly what you want to do
- Working with large repos (faster than GUI)
- SSH into remote servers

**Essential Commands**:
```bash
# Daily workflow
git status
git add <files>
git commit -m "message"
git push
git pull --rebase

# Branching
git checkout -b feature/my-feature
git branch -d feature/old-feature

# Stashing
git stash
git stash pop

# Undoing
git reset --soft HEAD~1
git revert <commit>
```

#### GitHub CLI (`gh`)

**Use When**:
- Creating pull requests from terminal
- Viewing PR status and checks
- Reviewing PRs from command line
- Managing issues
- CI/CD workflows

**Essential Commands**:
```bash
# PRs
gh pr create --title "My Feature" --body "Description"
gh pr list
gh pr view 123
gh pr checkout 123
gh pr merge 123

# Issues
gh issue create
gh issue list
```

#### VS Code (GitLens)

**Use When**:
- Exploring file/line history (Git Blame)
- Visualizing repository structure
- Comparing branches visually
- Resolving merge conflicts
- Learning Git (visual feedback)

---

## Real-World Scenarios

### Scenario 1: Implementing a New API Endpoint

**Task**: Add new REST endpoint for movie search

**Recommended Tools**:
- **IDE**: VS Code
- **Environment**: Local or Remote (choose based on preference)

**Workflow**:
```
1. VS Code: Design OpenAPI spec (docs/api/openapi.yaml)
2. Terminal: Run ogen codegen
3. VS Code: Implement handler in internal/api/
4. VS Code: Write unit tests
5. Terminal: Run tests locally
6. Remote (Coder): Run integration tests with full stack
7. Git CLI: Commit changes
8. GitHub CLI: Create PR
```

**Estimated Time**: 2-4 hours

### Scenario 2: Fixing a Complex Bug

**Task**: Database connection pool leak

**Recommended Tools**:
- **IDE**: VS Code (need debugger)
- **Environment**: Remote (Coder) - bug reproduces with full stack

**Workflow**:
```
1. Remote: Reproduce bug with integration tests
2. VS Code: Attach debugger to running process
3. VS Code: Set breakpoints in pgxpool code
4. VS Code: Inspect connection pool state
5. VS Code: Implement fix
6. Remote: Run integration tests to verify
7. Git CLI: Commit fix
8. GitHub CLI: Create PR
```

**Estimated Time**: 4-8 hours

### Scenario 3: Writing Documentation

**Task**: Update design docs for new feature

**Recommended Tools**:
- **IDE**: Zed (fast, distraction-free)
- **Environment**: Local (works offline)

**Workflow**:
```
1. Zed: Edit docs/dev/design/services/METADATA.md
2. Terminal: Run validation (python scripts/validate-doc-structure.py)
3. Terminal: Generate indexes (python scripts/generate-design-indexes.py)
4. Git CLI: Commit changes
5. GitHub CLI: Create PR
```

**Estimated Time**: 1-3 hours

### Scenario 4: Frontend Development

**Task**: New Svelte component for movie details

**Recommended Tools**:
- **IDE**: VS Code (Svelte extension required)
- **Environment**: Local (HMR is fastest)

**Workflow**:
```
1. VS Code: Create component (src/routes/movies/[id]/+page.svelte)
2. VS Code: Implement with Svelte 5 runes
3. Browser: Test with hot reload
4. VS Code: Add Tailwind CSS styling
5. VS Code: Write component tests (Vitest)
6. Terminal: Run tests
7. Git CLI: Commit
```

**Estimated Time**: 2-6 hours

---

## Decision Trees

### Choosing an IDE

```
Need to choose an IDE?
│
├─ Writing code?
│  ├─ Go backend? → VS Code (debugging) or Zed (fast editing)
│  ├─ Frontend (Svelte/TS)? → VS Code (only option)
│  └─ Python scripts? → VS Code or Zed
│
├─ Debugging?
│  ├─ Complex? → VS Code (graphical debugger)
│  └─ Simple? → Terminal (delve/pdb)
│
├─ Documentation?
│  └─ → Zed (fast, distraction-free)
│
└─ Quick edits?
   └─ → Zed (instant startup)
```

### Choosing Environment

```
Local or Remote?
│
├─ Working on Frontend?
│  └─ YES → Local (fast HMR)
│
├─ Need full stack?
│  └─ YES → Remote (databases included)
│
├─ Work from multiple devices?
│  └─ YES → Remote (access anywhere)
│
└─ Work offline frequently?
   └─ YES → Local
```

---

**Last Updated**: 2026-01-31
**Maintained By**: Development Team
