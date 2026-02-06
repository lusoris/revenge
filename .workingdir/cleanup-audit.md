# Cleanup Audit - 2026-02-06

## GitHub Remote State

### Releases (3)
- v0.1.2 - Errors Package 100% Coverage (2026-02-02, Latest)
- v0.1.1 - Test Coverage Sprint (2026-02-02)
- v0.1.0 - Skeleton Release (2026-02-02)

### Tags (5)
- v0.1.3 (no release), v0.1.2, v0.1.1, v0.1.0, v0.0.0 (no release)

### Packages
- No container images published (or token lacks read:packages scope)
- No Helm charts published

### Workflow Runs
- 30 total workflow runs

---

## GitHub Actions Workflows (22 files)

### Keep (Core CI/CD - Docker/Helm only)
1. ci.yml - Main CI (lint, test, build)
2. dev.yml - Development builds (8-job pipeline)
3. release-please.yml - Automated releases
4. release.yml - GoReleaser + Docker build on tags
5. deploy.yml - Validate Helm/Docker configs
6. security.yml - CodeQL, Trivy, gosec, govulncheck
7. codeql.yml - Dedicated CodeQL analysis

### Consider Keeping (Useful automation)
8. pr-checks.yml - PR validation
9. stale.yml - Auto-close stale issues
10. labels.yml - Label sync
11. dependabot.yml (config, not workflow)
12. coverage.yml - Coverage reports

### Remove (Redundant/overkill for current stage)
13. auto-label.yml - Complex auto-labeling
14. build-status.yml - Just logs info, doesn't do anything useful
15. dependency-update.yml - Conflicts with Dependabot
16. deploy-k3s.yml - K3s deployment test (overkill, no k3s yet)
17. deploy-k8s.yml - K8s deployment test (overkill, no k8s yet)
18. deploy-swarm.yml - Swarm deployment test (overkill)
19. doc-validation.yml - Only useful if keeping design pipeline
20. fetch-sources.yml - Only useful if keeping source pipeline
21. source-refresh.yml - Only useful if keeping source pipeline
22. validate-sot.yml - Validates SOURCE_OF_TRUTH format
23. _versions.yml - Version extraction from SOT (fragile)

---

## Claude Skills (29)

### Keep (Actually useful)
1. coder-template - Manage Coder templates
2. coder-workspace - Manage Coder workspaces
3. setup-workspace - Set up dev environment

### Maybe Keep (If simplified)
4. configure-dependabot - Simple GitHub config
5. configure-release-please - Simple GitHub config

### Remove (Reference non-existent or broken scripts)
All 24 others reference Python scripts in scripts/automation/ that are either placeholders or part of the design pipeline system.

---

## Scripts (72+ files)

### Keep (Core development)
- Makefile (25 targets - essential)
- docker-entrypoint.sh

### Pipeline Scripts (doc-pipeline/, source-pipeline/, deploy-pipeline/)
- 7 doc-pipeline scripts
- 3 source-pipeline scripts
- 1 deploy-pipeline script
- 2 orchestrator bash scripts

### Automation Scripts (scripts/automation/ - 36+ files)
- Most are either placeholders or design-pipeline adjacent
- Some are useful (run_tests, format_code, manage_docker)

### Root Scripts (scripts/ - 25+ files)
- Mostly doc generation, validation, fixing scripts

---

## Documentation (711 markdown, 160 YAML, 138 directories)

### docs/dev/design/ - ~90% hand-written design documents
### docs/dev/sources/ - External documentation fetched by pipeline
### Generated files: INDEX.md, breadcrumbs, cross-references, status tables
