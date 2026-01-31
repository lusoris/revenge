# Shared Tool Documentation

**Purpose**: Cross-tool integration, workflows, and onboarding guides

**Last Updated**: 2026-01-31

---

## Overview

This folder contains documentation that spans multiple development tools and provides guidance on how they work together in the Revenge project.

---

## Documentation Index

### Getting Started

- **[ONBOARDING.md](ONBOARDING.md)** - New developer onboarding guide
  - Role-based setup paths (Backend, Frontend, Full-Stack, DevOps)
  - First week checklist
  - Tool selection guide
  - Local vs Remote development decision matrix

### Integration & Workflows

- **[INTEGRATION.md](INTEGRATION.md)** - How all tools integrate
  - Tool ecosystem overview
  - Data flow between tools
  - Configuration synchronization
  - Recommended tool combinations
  - Settings precedence hierarchy

- **[WORKFLOWS.md](WORKFLOWS.md)** - Common development workflows
  - Local development workflow
  - Remote development workflow (Coder)
  - Code review workflow
  - Release workflow
  - Testing workflow

### Configuration

- **[SETTINGS_GUIDE.md](SETTINGS_GUIDE.md)** - Unified settings documentation
  - Settings across all tools
  - Settings precedence (EditorConfig → IDE → User)
  - Cross-tool settings matrix
  - Environment-specific settings

- **[PROFILES.md](PROFILES.md)** - Pre-configured settings profiles
  - Backend developer profile (Go-focused)
  - Frontend developer profile (Svelte/TypeScript-focused)
  - Full-stack developer profile
  - Remote development profile (optimized for Coder)

### Tool Comparison & Selection

- **[TOOL_COMPARISON.md](TOOL_COMPARISON.md)** - When to use which tool
  - VS Code vs Zed feature matrix
  - Local vs Remote development comparison
  - IDE selection guide
  - Performance characteristics

### Troubleshooting

- **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** - Cross-tool issues
  - IDE conflicts
  - Environment synchronization issues
  - Remote vs local problems
  - Common error patterns
  - Integration issues

### Development Assistance

- **[CLAUDE.md](CLAUDE.md)** - Claude Code instructions
  - Project overview and tech stack
  - Tool-specific guidance references
  - Code style and conventions
  - Common commands and workflows
  - Best practices for assistance

---

## Quick Links

### Tool-Specific Documentation

- [Claude Code](../../.claude/docs/INDEX.md) - AI assistant
- [VS Code](../../.vscode/docs/INDEX.md) - Primary IDE
- [Zed](../../.zed/docs/INDEX.md) - Modern alternative IDE
- [Coder](../../.coder/docs/INDEX.md) - Remote development platform
- [Git Hooks](../../.githooks/docs/INDEX.md) - Pre-commit automation
- [GitHub Actions](../../.github/docs/INDEX.md) - CI/CD

### Project Documentation

- [Design Documentation](../../docs/dev/design/INDEX.md) - Architecture and design
- [Source of Truth](../../docs/dev/design/00_SOURCE_OF_TRUTH.md) - Technology stack
- [Development Guide](../../docs/dev/design/operations/DEVELOPMENT.md) - Setup and workflow

---

## Contributing

When adding new tools or configurations:

1. Document tool-specific details in the tool's own `.tool/docs/` folder
2. Document cross-tool integration here in `.shared/docs/`
3. Update this INDEX.md with links to new documentation
4. Update INTEGRATION.md if the new tool interacts with existing tools
5. Update ONBOARDING.md if new developers need to know about the tool

---

## File Naming Conventions

- `ALL_CAPS.md` - Major documentation files
- `lowercase-with-dashes.md` - Supporting documentation
- `INDEX.md` - Always the entry point for each folder

---

## Status Legend

- ✅ Complete - Production ready
- 🟡 Partial - Exists but needs expansion
- 🔴 Planned - Not yet created
- ⚠️ Needs update - Outdated information

---

## Current Status

| Document | Status | Notes |
|----------|--------|-------|
| INDEX.md | ✅ | This file |
| ONBOARDING.md | ✅ | Role-based onboarding (5,000+ words) |
| INTEGRATION.md | ✅ | Tool ecosystem integration (4,300+ words) |
| WORKFLOWS.md | ✅ | Development workflows (3,500+ words) |
| SETTINGS_GUIDE.md | ✅ | Settings management (915 lines) |
| PROFILES.md | ✅ | 6 developer profiles (1,376 lines) |
| TOOL_COMPARISON.md | ✅ | Tool selection guide (6,500+ words) |
| TROUBLESHOOTING.md | ✅ | Cross-tool issues (4,500+ words) |
| CLAUDE.md | ✅ | Claude Code instructions (470+ lines) |

---

**Maintained By**: Development Team
