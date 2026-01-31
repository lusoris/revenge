# Claude Code Instructions - Git Hooks

**Tool**: Git Hooks + Pre-commit Framework
**Purpose**: Automated code quality checks and commit validation
**Documentation**: [docs/INDEX.md](docs/INDEX.md)

---

## Entry Point for Claude Code

When working with Git hooks for the Revenge project, always start by reading:

1. **Source of Truth**: [/docs/dev/design/00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md)
   - Commit message format
   - Code quality standards

2. **Development Guide**: [/docs/dev/design/operations/DEVELOPMENT.md](../docs/dev/design/operations/DEVELOPMENT.md)
   - Git workflow
   - Commit conventions

3. **Git Hooks Documentation**: [docs/INDEX.md](docs/INDEX.md)
   - Hook reference
   - Customization guide

---

## Git Hooks Overview

### Configured Hooks

1. **pre-commit** - Lint and validate before commit
   - Go: fmt, imports, vet, build, mod-tidy
   - Python: ruff formatting
   - General: trailing whitespace, YAML/JSON validation
   - Security: detect-secrets
   - Docker: hadolint

2. **commit-msg** - Enforce conventional commit format
   - Validates: `type(scope): description`
   - Types: feat, fix, docs, test, chore, refactor, perf, ci, build, revert

3. **pre-push** - Run tests before push
   - Unit tests
   - Integration tests (if enabled)

### Pre-commit Framework

**Config**: `.pre-commit-config.yaml`

**Hooks**:
- General hooks (trailing-whitespace, end-of-file-fixer, check-yaml, check-json)
- Go hooks (go-fmt, go-imports, go-vet, go-build, go-mod-tidy)
- Security hooks (detect-secrets)
- Conventional commits (commitizen)
- Docker hooks (hadolint)

---

## Common Tasks

### Adding a New Hook

1. Edit `.pre-commit-config.yaml`
2. Add hook configuration:
   ```yaml
   - repo: https://github.com/example/hooks
     rev: v1.0.0
     hooks:
       - id: hook-name
   ```
3. Test: `pre-commit run hook-name --all-files`
4. Document in [docs/HOOKS.md](docs/HOOKS.md)

### Modifying Existing Hook

1. Edit `.pre-commit-config.yaml`
2. Update hook args or configuration
3. Test: `pre-commit run --all-files`
4. Document changes in [docs/HOOKS.md](docs/HOOKS.md)

### Disabling a Hook Temporarily

```bash
# Skip all hooks
git commit --no-verify -m "commit message"

# Skip specific hook
SKIP=hook-name git commit -m "commit message"
```

**⚠️ Use sparingly** - Only when absolutely necessary

---

## Commit Message Format

**Format**: `type(scope): description`

**Types**:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `test` - Test changes
- `chore` - Maintenance tasks
- `refactor` - Code refactoring
- `perf` - Performance improvements
- `ci` - CI/CD changes
- `build` - Build system changes
- `revert` - Revert previous commit

**Examples**:
```
feat(auth): add OIDC authentication
fix(database): resolve connection pool leak
docs(api): update OpenAPI spec
test(movie): add integration tests
```

---

## Hook Behavior

### pre-commit

**Runs on**: `git commit`

**Checks**:
1. Go formatting (go fmt)
2. Go imports (goimports)
3. Go linting (go vet)
4. Go build (ensure it compiles)
5. Go mod tidy (dependencies)
6. Python formatting (ruff)
7. Trailing whitespace
8. YAML/JSON validity
9. Secret detection
10. Dockerfile linting

**Failure**: Commit is aborted

### commit-msg

**Runs on**: `git commit`

**Checks**:
- Commit message follows conventional commits format

**Failure**: Commit is aborted with message format error

### pre-push

**Runs on**: `git push`

**Checks**:
- Unit tests pass
- (Optional) Integration tests pass

**Failure**: Push is aborted

---

## Best Practices

1. **Run hooks locally** - Don't rely on CI to catch issues
2. **Fix hook failures** - Don't bypass hooks without good reason
3. **Keep hooks fast** - Hooks should complete in < 30 seconds
4. **Document custom hooks** - Update HOOKS.md when adding hooks
5. **Test hooks** - Run `pre-commit run --all-files` before pushing

---

## Troubleshooting

### Hook fails but code is correct

1. Check hook configuration in `.pre-commit-config.yaml`
2. Verify tool versions match project requirements
3. Run tool manually to see detailed error
4. Check for race conditions (e.g., parallel go builds)

### Hook is too slow

1. Identify slow hook: `pre-commit run --verbose --all-files`
2. Consider:
   - Running on changed files only
   - Splitting into multiple hooks
   - Moving to CI instead

### Can't bypass hook when needed

```bash
# Bypass all hooks
git commit --no-verify -m "emergency fix"

# Bypass specific hook
SKIP=go-build git commit -m "docs only change"
```

**⚠️ Only use when**:
- Emergency hotfix
- Hook is broken and needs fixing
- Documentation-only change
- Reverting a bad commit

---

## CI Integration

Git hooks run locally, but CI also runs the same checks:

**GitHub Actions** (`.github/workflows/ci.yml`):
- Runs all pre-commit hooks
- Runs full test suite
- Checks commit message format

**Difference**:
- **Hooks**: Fast feedback, run on developer machine
- **CI**: Authoritative, runs on clean environment

---

## Related Documentation

- **Git Hooks Docs**: [docs/INDEX.md](docs/INDEX.md)
- **Development Guide**: [../docs/dev/design/operations/DEVELOPMENT.md](../docs/dev/design/operations/DEVELOPMENT.md)
- **Best Practices**: [../docs/dev/design/operations/BEST_PRACTICES.md](../docs/dev/design/operations/BEST_PRACTICES.md)
- **GitHub Actions**: [../.github/docs/INDEX.md](../.github/docs/INDEX.md)

---

## Quick Commands

```bash
# Install hooks
pre-commit install

# Run all hooks on all files
pre-commit run --all-files

# Run specific hook
pre-commit run hook-name --all-files

# Update hook versions
pre-commit autoupdate

# Bypass hooks (use sparingly)
git commit --no-verify -m "message"

# Skip specific hook
SKIP=hook-name git commit -m "message"
```

---

**Last Updated**: 2026-01-31
**Maintained By**: Development Team
