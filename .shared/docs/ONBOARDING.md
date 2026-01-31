# New Developer Onboarding Guide

**Purpose**: Get you productive on the Revenge project quickly

**Last Updated**: 2026-01-31

---

## Welcome to Revenge! 🎬

This guide will help you set up your development environment and start contributing. Choose your path based on your role and preferences.

---

## Table of Contents

- [Quick Start (5 minutes)](#quick-start-5-minutes)
- [Role-Based Onboarding Paths](#role-based-onboarding-paths)
- [Local vs Remote Development](#local-vs-remote-development)
- [First Week Checklist](#first-week-checklist)
- [Tool Selection Guide](#tool-selection-guide)
- [Common Workflows](#common-workflows)
- [Getting Help](#getting-help)

---

## Quick Start (5 minutes)

### 1. Clone the Repository

```bash
git clone https://github.com/lusoris/revenge.git
cd revenge
```

###

 2. Read the Source of Truth

**REQUIRED READING** before writing any code:

📖 **[docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)**

This document contains:
- Technology stack and versions
- Architecture decisions
- All dependencies with exact versions
- API structure
- Database schemas
- QAR (adult content) terminology

### 3. Choose Your Development Environment

**Option A: Remote Development (Recommended for this week)**
- Consistent environment
- No local setup needed
- Access from anywhere
- **→ [Jump to Remote Setup](#remote-development-coder)**

**Option B: Local Development**
- Full control
- Faster iteration (frontend)
- Works offline
- **→ [Jump to Local Setup](#local-development)**

---

## Role-Based Onboarding Paths

Choose your primary role to get customized setup instructions:

### 🔧 Backend Developer (Go)

**Tech Stack**:
- Go 1.25.6
- PostgreSQL 18+
- Dragonfly (Redis-compatible cache)
- Typesense (search)
- River (job queue)

**Setup Path**:
1. [Remote Development on Coder](#remote-development-coder) (Recommended)
2. [Install Go Development Tools](#backend-go-tools)
3. [Setup Local Databases](#local-database-setup) (if local)
4. [Configure IDE for Go](#ide-setup-for-go)
5. [Run Your First Build](#first-go-build)

**First Task**: Implement a simple API endpoint

---

### 🎨 Frontend Developer (Svelte)

**Tech Stack**:
- SvelteKit 2
- Svelte 5 (Runes)
- TypeScript
- Tailwind CSS 4
- shadcn-svelte

**Setup Path**:
1. [Local Development](#local-development) (Recommended for frontend)
2. [Install Node.js and Dependencies](#frontend-tools)
3. [Configure IDE for Svelte](#ide-setup-for-svelte)
4. [Run Development Server](#frontend-dev-server)

**First Task**: Create a new component

---

### 🔄 Full-Stack Developer

**Setup Path**:
1. [Remote Development for Backend](#remote-development-coder)
2. [Local Development for Frontend](#local-development)
3. [Install Both Tool Sets](#full-stack-tools)
4. [Configure IDE for Both](#ide-setup-full-stack)

**First Task**: Build a feature end-to-end (API + UI)

---

### 🚀 DevOps Engineer

**Setup Path**:
1. [Install Docker and Coder CLI](#devops-tools)
2. [Understand Coder Template](../../.coder/template.tf)
3. [Review GitHub Actions](../../.github/workflows/)
4. [Setup Kubernetes/Swarm](#orchestration-setup) (if applicable)

**First Task**: Deploy a test workspace

---

### 📝 Documentation Writer

**Setup Path**:
1. [Install Zed](#install-zed) (fast, lightweight)
2. [Install Python for Scripts](#documentation-tools)
3. [Understand Documentation Structure](../../docs/dev/design/INDEX.md)
4. [Setup Claude Code Skills](#claude-code-setup)

**First Task**: Improve a design document

---

## Local vs Remote Development

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

### Recommendation

- **Backend work** → Remote (Coder)
- **Frontend work** → Local (fast HMR)
- **Full-stack** → Hybrid (remote backend, local frontend)
- **Documentation** → Either (Zed is fast locally)

---

## Remote Development (Coder)

### Prerequisites

- Access to Coder server: https://coder.ancilla.lol
- SSH client
- One of: VS Code, Zed, or JetBrains Gateway

### Step 1: Install Coder CLI

```bash
# Linux/macOS
curl -fsSL https://coder.com/install.sh | sh

# Windows (PowerShell)
winget install Coder.Coder

# Verify
coder version
```

### Step 2: Login to Coder

```bash
# Login
coder login https://coder.ancilla.lol

# Follow browser prompts to authenticate
```

### Step 3: Create Your Workspace

```bash
# List available templates
coder templates list

# Create workspace from revenge template
coder create --template revenge my-workspace

# Start workspace
coder start my-workspace
```

### Step 4: Connect Your IDE

#### Option A: VS Code (Desktop)

```bash
# Install Coder extension
code --install-extension coder.coder-remote

# Connect to workspace
coder code my-workspace
```

VS Code will open and connect via Remote-SSH automatically.

#### Option B: VS Code (Browser)

```bash
# Open browser-based VS Code
coder open my-workspace
```

Access code-server in your browser.

#### Option C: Zed (SSH)

```bash
# Get SSH command
coder ssh my-workspace

# In another terminal, connect Zed
zed ssh://coder-workspace
```

See [.coder/docs/ZED_INTEGRATION.md](../../.coder/docs/ZED_INTEGRATION.md) for details.

#### Option D: JetBrains Gateway

```bash
# Install JetBrains Gateway
# Add Coder plugin
# Connect to workspace via GUI
```

See [.coder/docs/JETBRAINS_INTEGRATION.md](../../.coder/docs/JETBRAINS_INTEGRATION.md) for details.

### Step 5: Verify Setup

Inside the workspace:

```bash
# Check Go version
go version
# Should show: go version go1.25.6

# Check Node version
node --version

# Check Python version
python --version

# Check Git
git --version

# Clone the repository (if not already)
git clone https://github.com/lusoris/revenge.git
cd revenge
```

### Step 6: Start Development

```bash
# Install Go dependencies
go mod download

# Run tests
go test ./...

# Start dev server with hot reload
air

# Access at workspace URL (forwarded port)
```

**✅ You're ready to code remotely!**

See [.coder/docs/REMOTE_WORKFLOW.md](../../.coder/docs/REMOTE_WORKFLOW.md) for complete remote development workflow.

---

## Local Development

### Prerequisites

- Git
- Go 1.25.6+
- Node.js 20+
- Python 3.12+ (for scripts)
- Docker & Docker Compose (for databases)

### Step 1: Install Development Tools

#### Go

```bash
# Download Go 1.25.6
# From: https://go.dev/dl/

# Verify
go version
# Should show: go1.25.6

# Set GOEXPERIMENT
export GOEXPERIMENT=greenteagc,jsonv2

# Add to ~/.bashrc or ~/.zshrc
echo 'export GOEXPERIMENT=greenteagc,jsonv2' >> ~/.bashrc
```

#### Node.js

```bash
# Install Node 20 LTS
# Using nvm (recommended)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 20
nvm use 20

# Verify
node --version
# Should show: v20.x.x
```

#### Python

```bash
# Python 3.12+
python --version
# Should show: Python 3.12.x

# Install script dependencies
pip install -r scripts/requirements.txt
```

### Step 2: Install Development Databases

```bash
# Start PostgreSQL, Dragonfly, Typesense
docker-compose -f docker-compose.dev.yml up -d

# Verify services running
docker-compose -f docker-compose.dev.yml ps

# Expected:
# - postgres (port 5432)
# - dragonfly (port 6379)
# - typesense (port 8108)
```

### Step 3: Install Go Tools

```bash
# golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# air (hot reload)
go install github.com/air-verse/air@latest

# gopls (LSP)
go install golang.org/x/tools/gopls@latest

# mockery (mocking)
go install github.com/vektra/mockery/v3@latest

# Verify
golangci-lint version
air -v
gopls version
mockery --version
```

### Step 4: Install Python Tools

```bash
# ruff (linter/formatter)
pip install ruff

# pytest (testing)
pip install pytest pytest-cov

# Verify
ruff --version
pytest --version
```

### Step 5: Install Git Hooks

```bash
# Install pre-commit framework
pip install pre-commit

# Install hooks
pre-commit install
pre-commit install --hook-type commit-msg
pre-commit install --hook-type pre-push

# Test hooks
pre-commit run --all-files
```

### Step 6: Configure IDE

See:
- [.vscode/docs/INDEX.md](../../.vscode/docs/INDEX.md) for VS Code
- [.zed/docs/INDEX.md](../../.zed/docs/INDEX.md) for Zed

### Step 7: Build and Run

```bash
# Download Go dependencies
go mod download

# Run tests
go test ./...

# Build
GOEXPERIMENT=greenteagc,jsonv2 go build ./cmd/revenge

# Run with hot reload
air

# Access at http://localhost:8096
```

**✅ You're ready to code locally!**

---

## First Week Checklist

Use this checklist to track your onboarding progress:

### Day 1: Environment Setup

- [ ] Clone repository
- [ ] Read [SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)
- [ ] Choose local or remote development
- [ ] Install required tools (Go/Node/Python)
- [ ] Setup IDE (VS Code or Zed)
- [ ] Install Git hooks
- [ ] Run first build successfully
- [ ] Run tests successfully

### Day 2: Understand the Codebase

- [ ] Read [ARCHITECTURE.md](../../docs/dev/design/architecture/01_ARCHITECTURE.md)
- [ ] Read [DESIGN_PRINCIPLES.md](../../docs/dev/design/architecture/02_DESIGN_PRINCIPLES.md)
- [ ] Explore project structure
- [ ] Understand module pattern
- [ ] Review existing code for your area (backend/frontend)
- [ ] Run the application locally or remotely

### Day 3: Make Your First Change

- [ ] Pick a "good first issue" or small task
- [ ] Create feature branch
- [ ] Make code changes
- [ ] Write/update tests
- [ ] Run linters and formatters
- [ ] Commit with conventional commit message
- [ ] Create pull request

### Day 4: Code Review Process

- [ ] Review someone else's PR
- [ ] Respond to feedback on your PR
- [ ] Understand CI/CD pipeline
- [ ] Learn how releases work

### Day 5: Deep Dive

- [ ] Explore one design doc deeply
- [ ] Understand one integration (e.g., TMDb, Radarr)
- [ ] Read test patterns
- [ ] Experiment with a larger feature

---

## Tool Selection Guide

### IDE: VS Code vs Zed

| Feature | VS Code | Zed |
|---------|---------|-----|
| **Extensions** | 50,000+ | Limited (built-in mostly) |
| **Debugging** | Full debug UI | Terminal-based |
| **Performance** | Good | Excellent (Rust-based) |
| **Go Support** | Excellent | Excellent |
| **Svelte Support** | Excellent | Good |
| **Remote Dev** | Full support | SSH only |
| **Learning Curve** | Moderate | Low |
| **Customization** | Extensive | Minimal |

**Recommendation**:
- **Backend Go work** → Either (Zed is faster, VS Code has more features)
- **Frontend Svelte work** → VS Code (better extension)
- **Debugging complex issues** → VS Code (full debugger UI)
- **Quick edits / documentation** → Zed (fast startup)
- **Remote development** → Either works well

### Local vs Remote: Decision Tree

```
Start Here
    │
    ├─→ Working on Frontend?
    │   └─→ YES → Local Development (fast HMR)
    │   └─→ NO → Continue
    │
    ├─→ Need consistent environment across team?
    │   └─→ YES → Remote (Coder)
    │   └─→ NO → Continue
    │
    ├─→ Have powerful local machine?
    │   └─→ YES → Either works
    │   └─→ NO → Remote (Coder)
    │
    ├─→ Need to work from multiple devices?
    │   └─→ YES → Remote (Coder)
    │   └─→ NO → Either works
    │
    └─→ Default → Remote (Coder) for backend, Local for frontend
```

---

## Common Workflows

### Workflow 1: Implementing a Feature

1. **Plan**
   - Read relevant design docs
   - Understand requirements
   - Check SOURCE_OF_TRUTH for tech stack

2. **Code**
   - Create feature branch: `git checkout -b feature/my-feature`
   - Write code following existing patterns
   - Add tests (80% coverage minimum)

3. **Test**
   - Run tests: `go test ./...`
   - Run linters: `golangci-lint run`
   - Test manually in dev environment

4. **Commit**
   - Stage changes: `git add .`
   - Commit: `git commit -m "feat: add my feature"`
   - Pre-commit hooks run automatically

5. **Push & PR**
   - Push: `git push origin feature/my-feature`
   - Create PR on GitHub
   - Wait for CI checks and review

### Workflow 2: Fixing a Bug

1. **Reproduce**
   - Understand the bug
   - Create minimal reproduction

2. **Fix**
   - Create bugfix branch: `git checkout -b fix/bug-description`
   - Write failing test first (TDD)
   - Fix the bug
   - Verify test passes

3. **Commit**
   - `git commit -m "fix: resolve bug description"`

4. **Push & PR**
   - Push and create PR
   - Reference issue number in PR description

### Workflow 3: Updating Documentation

1. **Locate**
   - Find doc to update in `docs/dev/design/`

2. **Edit**
   - Use Zed or VS Code
   - Follow existing doc structure
   - Update cross-references if needed

3. **Validate**
   - Run: `python scripts/validate-doc-structure.py`
   - Run: `python scripts/validate-links.py`
   - Generate indexes: `python scripts/generate-design-indexes.py`

4. **Commit**
   - `git commit -m "docs: update design doc"`

---

## Getting Help

### Resources

1. **Documentation**
   - [Design Docs](../../docs/dev/design/INDEX.md)
   - [Source of Truth](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)
   - [Tool Docs](./)

2. **Ask Questions**
   - GitHub Discussions
   - Discord (if available)
   - Comments on issues/PRs

3. **Claude Code**
   - Use Claude Code for code questions
   - Custom skills for automation
   - See [.claude/docs/](../../.claude/docs/)

### Common Issues

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for:
- IDE setup issues
- Build failures
- Test failures
- Git hook problems
- Remote development issues

---

## Next Steps

After completing onboarding:

1. **Pick Your First Task**
   - Look for "good first issue" labels
   - Ask maintainers for recommendations

2. **Join the Community**
   - Introduce yourself
   - Review others' PRs
   - Share knowledge

3. **Explore Advanced Topics**
   - [Metadata System](../../docs/dev/design/architecture/03_METADATA_SYSTEM.md)
   - [RBAC with Casbin](../../docs/dev/design/features/shared/RBAC_CASBIN.md)
   - [River Job Queue](../../docs/dev/design/integrations/infrastructure/RIVER.md)

---

## Summary

You should now have:
- ✅ Development environment setup (local or remote)
- ✅ IDE configured
- ✅ Tools installed
- ✅ First build successful
- ✅ Understanding of project structure
- ✅ Knowledge of common workflows

**Ready to code? Start with a small task and don't hesitate to ask for help!**

---

**Maintained By**: Development Team
**Last Updated**: 2026-01-31
