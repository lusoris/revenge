# Claude Code Instructions - GitHub Configuration

**Tool**: GitHub (Actions, Dependabot, Templates)
**Purpose**: CI/CD automation, dependency management, issue/PR templates
**Documentation**: [docs/INDEX.md](docs/INDEX.md)

---

## Entry Point for Claude Code

When working with GitHub configuration for the Revenge project, always start by reading:

1. **Source of Truth**: [/docs/dev/design/00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md)
   - Release versioning
   - Dependency policy
   - CI/CD requirements

2. **Operations Guide**: [/docs/dev/design/operations/INDEX.md](../docs/dev/design/operations/INDEX.md)
   - GitFlow workflow
   - Branch protection
   - Versioning strategy

3. **GitHub Docs**: [docs/INDEX.md](docs/INDEX.md)
   - Actions reference
   - Dependabot configuration

---

## GitHub Configuration Overview

### Workflows (15 total)

Located in `.github/workflows/`:

1. **ci.yml** - Run tests and linting
2. **coverage.yml** - Generate coverage reports
3. **dev.yml** - Development builds
4. **doc-validation.yml** - Validate documentation
5. **source-refresh.yml** - Weekly external doc refresh
6. **security.yml** - Security scanning
7. **release.yml** - Release creation
8. **release-please.yml** - Automated releases
9. **auto-label.yml** - Auto-label PRs
10. **pr-checks.yml** - PR validation
11. **stale.yml** - Mark stale issues
12. **build-status.yml** - Build status checks
13. **fetch-sources.yml** - Fetch external sources
14. **dependency-update.yml** - Update dependencies
15. **labels.yml** - Label sync

### Other Configuration

- **dependabot.yml** - Dependency updates
- **labeler.yml** - Auto-labeling rules
- **CODEOWNERS** - Code ownership rules

---

## Common Tasks

### Adding a New Workflow

1. Create `.github/workflows/my-workflow.yml`
2. Follow existing workflow patterns
3. Test with `act` (local runner) or in a PR
4. Document in [docs/WORKFLOWS.md](docs/WORKFLOWS.md) (to be created)

### Modifying Existing Workflow

1. Edit workflow file in `.github/workflows/`
2. Understand trigger conditions
3. Test changes in a PR or branch
4. Document changes in [docs/WORKFLOWS.md](docs/WORKFLOWS.md)

### Managing Secrets

1. Add secret: Settings → Secrets and variables → Actions
2. Reference in workflow: `${{ secrets.SECRET_NAME }}`
3. Document in [docs/SECRETS.md](docs/SECRETS.md) (to be created)
4. Never commit secrets to repository

### Configuring Dependabot

1. Edit `dependabot.yml`
2. Configure package ecosystem, schedule, automerge rules
3. Test with manual trigger
4. Document in [docs/DEPENDABOT.md](docs/DEPENDABOT.md)

---

## Workflow Details

### ci.yml - Continuous Integration

**Triggers**: Push to main/develop, Pull requests

**Jobs**:
1. Lint Go code (golangci-lint)
2. Run Go tests
3. Build application
4. Lint Python scripts (ruff)
5. Validate YAML/JSON

**Matrix**: Go versions (latest stable)

### coverage.yml - Code Coverage

**Triggers**: Push to main/develop, Pull requests

**Jobs**:
1. Run Go tests with coverage
2. Generate coverage report
3. Upload to Codecov (if configured)
4. Comment coverage on PR

**Requirement**: 80% minimum coverage

### doc-validation.yml - Documentation Validation

**Triggers**: Pull requests to docs/, Push to main

**Jobs**:
1. Validate document structure
2. Check for broken links
3. Verify cross-references
4. Validate status tables

### source-refresh.yml - External Source Refresh

**Triggers**: Schedule (weekly), Manual

**Jobs**:
1. Fetch external documentation sources
2. Update SOURCES.yaml
3. Create PR if changes detected

### release-please.yml - Automated Releases

**Triggers**: Push to main

**Jobs**:
1. Analyze commits (conventional commits)
2. Generate changelog
3. Bump version
4. Create release PR or tag

---

## Dependabot Configuration

**File**: `dependabot.yml`

**Settings**:
- **Schedule**: Weekly on Mondays before 6am (Europe/Berlin)
- **Managers**: gomod, dockerfile, docker-compose, github-actions
- **Auto-merge**: Minor and patch Go dependencies
- **Manual review**: Major version updates
- **Concurrency**: 5 PRs max, 2 per hour

**Customization**:
```yaml
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "06:00"
      timezone: "Europe/Berlin"
```

---

## Best Practices

1. **Test workflows locally** - Use `act` before pushing
2. **Use matrix builds** - Test across Go versions, OS, etc.
3. **Cache dependencies** - Speed up workflows with caching
4. **Minimal permissions** - Use least-privilege principle
5. **Document workflows** - Update docs when adding/changing workflows

---

## Troubleshooting

### Workflow fails with permission error

1. Check workflow permissions in YAML
2. Verify GITHUB_TOKEN has required scopes
3. Check branch protection rules
4. Review CODEOWNERS

### Dependabot PR fails

1. Check for breaking changes in dependency
2. Verify tests pass locally
3. Check for conflicts with other dependencies
4. Review Dependabot logs

### Workflow timeout

1. Identify slow steps in workflow logs
2. Add caching for dependencies
3. Parallelize jobs where possible
4. Consider splitting into multiple workflows

### Secret not available in workflow

1. Verify secret is added in repository settings
2. Check secret name matches workflow reference
3. Ensure workflow has permission to access secrets
4. For organization secrets, verify repository access

---

## Security Considerations

1. **Never log secrets** - Be careful with debug output
2. **Use GITHUB_TOKEN** - Don't create PATs unless necessary
3. **Minimal permissions** - Only grant required permissions
4. **Pin action versions** - Use specific commit SHAs for security
5. **Review third-party actions** - Audit before using

---

## Related Documentation

- **GitHub Docs**: [docs/INDEX.md](docs/INDEX.md)
- **GitFlow**: [../docs/dev/design/operations/GITFLOW.md](../docs/dev/design/operations/GITFLOW.md)
- **Branch Protection**: [../docs/dev/design/operations/BRANCH_PROTECTION.md](../docs/dev/design/operations/BRANCH_PROTECTION.md)
- **Versioning**: [../docs/dev/design/operations/VERSIONING.md](../docs/dev/design/operations/VERSIONING.md)

---

## Quick Commands

```bash
# Test workflow locally (requires act)
act -j job-name

# List available workflows
gh workflow list

# Run workflow manually
gh workflow run workflow-name.yml

# View workflow runs
gh run list

# View specific run
gh run view RUN_ID

# Download workflow artifacts
gh run download RUN_ID

# Enable Dependabot
gh api repos/:owner/:repo/automated-security-fixes -X PUT
```

---

**Last Updated**: 2026-01-31
**Maintained By**: Development Team
