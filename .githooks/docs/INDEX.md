# Git Hooks Documentation

> Reference documentation for Git hooks

## Documents

| Document | Description |
|----------|-------------|
| [HOOKS.md](HOOKS.md) | Project Git hooks reference and usage |

## External Documentation

Comprehensive Git hooks documentation (auto-fetched from git-scm.com):

| Document | Description |
|----------|-------------|
| [git-hooks/hooks.md](git-hooks/hooks.md) | Complete Git hooks reference (githooks manual) |
| [git-hooks/hooks-guide.md](git-hooks/hooks-guide.md) | Pro Git book chapter on Git hooks |
| [git-hooks/book.md](git-hooks/book.md) | Pro Git book (complete) |
| [git-hooks/main.md](git-hooks/main.md) | Git documentation overview |

## Quick Links

- [Git Hooks Manual](https://git-scm.com/docs/githooks)
- [Pro Git: Hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks)
- [Git Documentation](https://git-scm.com/doc)

## Project Hooks

| Hook | Purpose |
|------|---------|
| [pre-commit](../pre-commit) | Lint, format check, tests, and security scanning |
| [commit-msg](../commit-msg) | Enforce Conventional Commits format |
| [pre-push](../pre-push) | Run full test suite and linter before push |

## Setup

Configure Git to use project hooks:

```bash
git config core.hooksPath .githooks
```

Or run the setup script:

```bash
./scripts/setup-hooks.sh
```

## Hook Behavior

### pre-commit

Runs before commit message is created:

- ✅ Code formatting check (`gofmt`)
- ✅ Go vet
- ✅ Security check (no hardcoded secrets)
- ✅ Tests on affected packages
- ⚠️ Warning for TODO/FIXME without issue numbers

**Bypass**: `git commit --no-verify`

### commit-msg

Validates commit message format:

- ✅ Conventional Commits format
- ✅ Subject length check (≤ 72 chars)
- ❌ No WIP commits on main/develop

**Bypass**: `git commit --no-verify`

### pre-push

Runs before push:

- ❌ Direct push to main blocked
- ✅ Full test suite on develop branch
- ✅ Linter checks (golangci-lint)
- ✅ Build verification

**Bypass**: `git push --no-verify`
