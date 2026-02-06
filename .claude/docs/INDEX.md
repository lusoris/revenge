# Claude Code & API Documentation Index

> Fetched from official Anthropic documentation
> Last updated: 2026-01-31
> Sources: code.claude.com, platform.claude.com

---

## Documentation Map

```
platform.claude.com/docs/         <- API & Developer Platform Entry Point
│
├── /en/home                      <- Documentation Home
│   ├── Claude Developer Platform
│   │   ├── Get started          -> API Quickstart
│   │   ├── Features overview    -> API capabilities
│   │   ├── API reference        -> Full API docs
│   │   └── Release notes        -> API changelog
│   │
│   ├── Claude Code              <- Links to code.claude.com
│   │   ├── Quickstart           -> code.claude.com/docs/en/quickstart
│   │   ├── Reference            -> code.claude.com/docs/en/overview
│   │   └── Changelog            -> GitHub CHANGELOG.md
│   │
│   └── Learning Resources
│       ├── Anthropic Courses    -> anthropic.skilljar.com
│       ├── Claude Cookbook      -> platform.claude.com/cookbooks
│       └── Claude Quickstarts   -> GitHub anthropic-quickstarts
│
code.claude.com/docs/             <- Claude Code Entry Point
│
├── /en/overview                  <- Claude Code Overview
├── /en/quickstart                <- Getting Started
├── /en/best-practices            <- Best Practices
├── /en/common-workflows          <- Common Workflows
├── /en/memory                    <- CLAUDE.md & Memory
├── /en/mcp                       <- MCP Integration
├── /en/hooks                     <- Hooks System
├── /en/skills                    <- Skills System
└── /en/sub-agents                <- Subagents System
```

---

## Local Documentation Files

| File | Topic | Source |
|------|-------|--------|
| [01_OVERVIEW.md](01_OVERVIEW.md) | Claude Code Overview | code.claude.com |
| [02_MEMORY.md](02_MEMORY.md) | CLAUDE.md & Memory Management | code.claude.com |
| [03_MCP.md](03_MCP.md) | Model Context Protocol | code.claude.com |
| [04_HOOKS.md](04_HOOKS.md) | Hooks Lifecycle & Config | code.claude.com |
| [05_CODING_STANDARDS.md](05_CODING_STANDARDS.md) | Best Practices Guide | code.claude.com |
| [06_COMMON_WORKFLOWS.md](06_COMMON_WORKFLOWS.md) | Common Workflows | code.claude.com |
| [07_SKILLS.md](07_SKILLS.md) | Skills System Reference | code.claude.com |
| [08_SUBAGENTS.md](08_SUBAGENTS.md) | Subagents Reference | code.claude.com |
| [09_API_OVERVIEW.md](09_API_OVERVIEW.md) | Claude API Overview | platform.claude.com |
| [10_QUICKSTART.md](10_QUICKSTART.md) | Claude Code Quickstart | code.claude.com |

---

## Quick Reference

### Claude Code Features

| Feature | Description | Doc |
|---------|-------------|-----|
| **CLAUDE.md** | Project-specific instructions and memory | [02_MEMORY](02_MEMORY.md) |
| **Skills** | Reusable prompts and workflows | [07_SKILLS](07_SKILLS.md) |
| **Subagents** | Specialized AI assistants for specific tasks | [08_SUBAGENTS](08_SUBAGENTS.md) |
| **Hooks** | Automated actions on lifecycle events | [04_HOOKS](04_HOOKS.md) |
| **MCP** | Connect to external tools and services | [03_MCP](03_MCP.md) |

### Key Concepts

| Concept | Description |
|---------|-------------|
| **Context Window** | Holds conversation, files read, command output - fills up fast |
| **Plan Mode** | Read-only exploration before making changes |
| **Verification** | Give Claude tests/screenshots to verify its own work |
| **Subagent Isolation** | Run tasks in separate context windows |

### Best Practices Summary

1. **Give Claude verification criteria** - tests, screenshots, expected outputs
2. **Explore first, then plan, then code** - use Plan Mode
3. **Provide specific context** - reference files with `@`, paste images
4. **Write effective CLAUDE.md** - concise, actionable instructions
5. **Manage context aggressively** - use `/clear` between tasks
6. **Use subagents for investigation** - keep main context clean

---

## External Resources

- **Claude Code Changelog**: https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md
- **API Release Notes**: platform.claude.com/docs/en/release-notes/api
- **Claude Cookbook**: platform.claude.com/cookbooks
- **Anthropic Courses**: anthropic.skilljar.com
- **Discord Community**: anthropic.com/discord
