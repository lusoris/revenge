# Design Documentation Completion - TODO List

<!-- DESIGN: .complete, README, SCAFFOLD_TEMPLATE, test_output_claude -->


**Purpose**: Detailed action items for completing all design documentation
**Source**: [ANALYSIS.md](ANALYSIS.md) gap analysis
**Automation**: Aligned with YAML-to-Markdown generation system
**Status**: Ready for execution

---

## Overview

**Total Tasks**: 46
- 🔴 **Critical** (15 tasks) - Blocks implementation or adoption
- 🟡 **High** (18 tasks) - Important for completeness
- 🟢 **Medium** (10 tasks) - Nice to have
- ⚪ **Low** (3 tasks) - Can be deferred

**Estimated Effort**: 4 weeks (1 person full-time)

---

## Priority 1: CRITICAL Operations Docs (Week 1)

### Task 1.1: Create DEVELOPMENT.md 🔴

**File**: `data/operations/DEVELOPMENT.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 8 hours
**Dependencies**: None
**Blocks**: Contributor onboarding

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 1.1** (DEVELOPMENT.md questions)
2. [ ] **Update YAML** (`data/operations/DEVELOPMENT.yaml`):
   ```yaml
   technical_summary: |-
     > Developer environment setup for local development

   # Replace PLACEHOLDER content
   wiki_overview: |
     User-friendly overview of setting up development environment
   ```

3. [ ] **Add Content Sections** (edit YAML):
   - `prerequisites` (list): Required tools with versions
   - `environment_setup` (text): Step-by-step environment setup
   - `first_run` (text): Commands to run on first clone
   - `hot_reload` (text): Using air for Go, Vite for frontend
   - `database_setup` (text): PostgreSQL setup, seed data
   - `common_issues` (list): Troubleshooting common problems
   - `sources`: Link to air docs, Vite docs

4. [ ] **Run Generation**:
   ```bash
   python scripts/automation/batch_regenerate.py
   ```

5. [ ] **Verify Output**:
   - `docs/dev/design/operations/DEVELOPMENT.md` (Claude version)
   - `docs/wiki/operations/development.md` (User version)

6. [ ] **Run Pipeline**:
   ```bash
   ./scripts/doc-pipeline.sh --apply
   ```

7. [ ] **Test**: Can a new contributor clone and run using only this doc?

**Acceptance Criteria**:
- ✅ Complete step-by-step setup for Linux/macOS
- ✅ No PLACEHOLDER content
- ✅ All commands tested and working
- ✅ Troubleshooting section covers common issues

---

### Task 1.2: Create SETUP.md 🔴

**File**: `data/operations/SETUP.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 12 hours
**Dependencies**: None
**Blocks**: User adoption

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 1.2** (SETUP.md questions)
2. [ ] **Update YAML** (`data/operations/SETUP.yaml`):
   ```yaml
   technical_summary: >
     > Production deployment guide for self-hosting Revenge

   # Add deployment content
   deployment_methods:
     - name: Docker Compose
       difficulty: Easy
       recommended: true
       use_case: Homelab, NAS deployments
     - name: Kubernetes (K3s)
       difficulty: Advanced
       recommended: false
       use_case: High availability
   ```

3. [ ] **Add Content Sections**:
   - `hardware_requirements`: Minimum/recommended specs
   - `docker_compose_setup`: Full docker-compose.yml example
   - `kubernetes_setup`: Helm chart installation
   - `initial_configuration`: First-time setup wizard
   - `reverse_proxy`: Traefik/Nginx/Caddy examples
   - `upgrade_guide`: How to upgrade versions
   - `backup_restore`: Backup strategies

4. [ ] **Create Assets**:
   - Sample `docker-compose.yml` in docs or repo root
   - Sample `.env.example` file
   - Sample Helm `values.yaml`

5. [ ] **Run Generation + Pipeline**
6. [ ] **Test**: Deploy using only this doc (fresh VM or container)

**Acceptance Criteria**:
- ✅ Docker Compose setup works end-to-end
- ✅ Environment variables documented
- ✅ Reverse proxy configs tested
- ✅ Upgrade path documented

---

### Task 1.3: Create GITFLOW.md 🔴

**File**: `data/operations/GITFLOW.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 6 hours
**Dependencies**: None
**Blocks**: Contributor workflow

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 1.3** (GITFLOW.md questions)
2. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > Git branching strategy and contribution workflow

   # Add workflow content
   branch_strategy:
     main_branch: develop
     feature_prefix: feat/
     bugfix_prefix: fix/
     release_prefix: release/

   commit_convention:
     style: Conventional Commits
     types:
       - feat
       - fix
       - docs
       - chore
       - refactor
       - test
   ```

3. [ ] **Add Content Sections**:
   - `branch_naming`: Branch naming conventions
   - `commit_format`: Commit message format with examples
   - `pr_process`: Pull request creation and review
   - `release_process`: How releases are created
   - `hotfix_process`: Emergency bugfix workflow

4. [ ] **Create Diagrams** (optional):
   - Mermaid diagram of branching strategy
   - PR workflow flowchart

5. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Clear branch naming rules
- ✅ Commit message examples
- ✅ PR checklist provided
- ✅ Release process documented

---

### Task 1.4: Create BEST_PRACTICES.md 🔴

**File**: `data/operations/BEST_PRACTICES.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 10 hours
**Dependencies**: None
**Blocks**: Code quality standards

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 1.4** (BEST_PRACTICES.md questions)
2. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > Go coding patterns, conventions, and best practices

   # Add patterns
   code_organization:
     - pattern: Module structure
       example: |
         internal/content/movie/
         ├── entity.go
         ├── repository.go
         ├── service.go
         └── module.go
   ```

3. [ ] **Add Content Sections**:
   - `module_structure`: Directory layout patterns
   - `error_handling`: Sentinel errors, custom types
   - `context_usage`: Context patterns
   - `testing_patterns`: Table-driven tests, mocking
   - `logging_patterns`: Structured logging with slog
   - `performance_patterns`: Caching, pooling
   - `security_patterns`: Input validation, SQL injection prevention

4. [ ] **Add Code Examples**: Real code snippets from codebase
5. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Module structure clear
- ✅ Error handling patterns with examples
- ✅ Testing patterns with examples
- ✅ Security best practices documented

---

### Task 1.5: Create VERSIONING.md 🔴

**File**: `data/operations/VERSIONING.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 4 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 1.5** (VERSIONING.md questions)
2. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > Semantic versioning strategy and release process

   current_version: v0.1.0
   versioning_scheme: SemVer 2.0.0
   ```

3. [ ] **Add Content Sections**:
   - `semver_rules`: When to bump major/minor/patch
   - `pre_releases`: Alpha/beta versioning
   - `release_cadence`: Frequency of releases
   - `compatibility`: API/DB compatibility guarantees

4. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ SemVer rules clear
- ✅ Release cadence defined
- ✅ Compatibility guarantees documented

---

### Task 1.6: Create BRANCH_PROTECTION.md 🔴

**File**: `data/operations/BRANCH_PROTECTION.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 3 hours
**Dependencies**: Task 1.3 (GITFLOW.md)

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 1.6** (BRANCH_PROTECTION.md questions)
2. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > GitHub branch protection rules for main and develop

   protected_branches:
     - name: develop
       required_approvals: 1
       required_checks:
         - ci/tests
         - ci/lint
         - ci/build
   ```

3. [ ] **Add Content Sections**:
   - `protected_branches`: Which branches are protected
   - `required_checks`: CI checks that must pass
   - `review_requirements`: Number of approvals
   - `force_push`: Force push policy

4. [ ] **Configure GitHub**: Apply these rules in GitHub repo settings
5. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Branch protection rules documented
- ✅ GitHub settings match documentation
- ✅ Cannot merge without approvals

---

### Task 1.7: Create REVERSE_PROXY.md 🟡

**File**: `data/operations/REVERSE_PROXY.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 6 hours
**Dependencies**: Task 1.2 (SETUP.md)

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 1.7** (REVERSE_PROXY.md questions)
2. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > Reverse proxy configuration for Traefik, Nginx, and Caddy

   supported_proxies:
     - Traefik
     - Nginx
     - Caddy
   ```

3. [ ] **Add Content Sections**:
   - `traefik_config`: Traefik v2/v3 configuration
   - `nginx_config`: Nginx configuration
   - `caddy_config`: Caddy configuration
   - `websocket_proxying`: WebSocket upgrade handling
   - `tls_termination`: HTTPS configuration

4. [ ] **Add Config Examples**: Full working configs for each proxy
5. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Working config for Traefik
- ✅ Working config for Nginx
- ✅ Working config for Caddy
- ✅ WebSocket support verified

---

### Task 1.8: Create DATABASE_AUTO_HEALING.md 🟢

**File**: `data/operations/DATABASE_AUTO_HEALING.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 5 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 1.8** (DATABASE_AUTO_HEALING.md questions)
2. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > Database connection recovery and consistency checking

   healing_strategies:
     - Connection pool auto-reconnect (pgxpool)
     - Health check monitoring
     - Migration retry logic
   ```

3. [ ] **Add Content Sections**:
   - `connection_recovery`: pgxpool reconnect behavior
   - `health_checks`: Database health monitoring
   - `consistency_checks`: Background validation jobs
   - `migration_recovery`: Handling failed migrations

4. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Connection recovery documented
- ✅ Health check implementation described
- ✅ Migration failure handling documented

---

## Priority 2: CRITICAL Technical Docs (Week 2)

### Task 2.1: Create API.md 🔴

**File**: `data/technical/API.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 10 hours
**Dependencies**: None
**Blocks**: API implementation

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 2.1** (API.md questions)
2. [ ] **Verify OpenAPI Spec Exists**: Check `api/openapi/spec.yaml`
3. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > REST API reference and OpenAPI specification

   api_version: v1
   base_path: /api/v1
   authentication: Bearer token
   ```

4. [ ] **Add Content Sections**:
   - `openapi_spec_location`: Where to find spec
   - `api_versioning`: Versioning strategy
   - `authentication`: Auth methods (Bearer, API key)
   - `rate_limiting`: Rate limit policies
   - `error_responses`: Standard error format
   - `pagination`: Cursor vs offset
   - `filtering`: Query parameter patterns

5. [ ] **Add Code Examples**: Request/response examples
6. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ OpenAPI spec linked
- ✅ Authentication methods documented
- ✅ Error response format standardized
- ✅ Example requests/responses provided

---

### Task 2.2: Create CONFIGURATION.md 🔴

**File**: `data/technical/CONFIGURATION.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 8 hours
**Dependencies**: None
**Blocks**: Configuration implementation

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 2.2** (CONFIGURATION.md questions)
2. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > Configuration system using koanf (YAML + env vars)

   config_file: config.yaml
   env_prefix: REVENGE_
   ```

3. [ ] **Add Content Sections**:
   - `config_file_location`: Default paths
   - `environment_variables`: Env var mapping
   - `hot_reload`: Which settings can reload
   - `validation`: Config validation rules
   - `secrets_management`: File-based secrets, env expansion

4. [ ] **Create Sample Config**: `config.example.yaml`
5. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Config file format documented
- ✅ Env var naming clear
- ✅ Hot reload behavior described
- ✅ Example config provided

---

### Task 2.3: Create FRONTEND.md 🔴

**File**: `data/technical/FRONTEND.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 10 hours
**Dependencies**: None
**Blocks**: Frontend development

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 2.3** (FRONTEND.md questions)
2. [ ] **Update YAML**:
   ```yaml
   technical_summary: >
     > Frontend architecture with SvelteKit, Svelte 5, Tailwind CSS 4

   framework: SvelteKit 2
   ui_library: Svelte 5 (runes)
   styling: Tailwind CSS 4
   components: shadcn-svelte
   ```

3. [ ] **Add Content Sections**:
   - `component_organization`: Directory structure
   - `state_management`: TanStack Query + Svelte stores
   - `styling_system`: Tailwind config, theming
   - `api_client`: Fetch wrapper, generated client
   - `forms`: Superforms + Zod validation
   - `authentication`: JWT storage, refresh logic

4. [ ] **Add Code Examples**: Component examples, API client usage
5. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Component structure clear
- ✅ State management pattern documented
- ✅ API client usage shown
- ✅ Form handling examples provided

---

## Priority 3: Content Module Completion (Week 2-3)

### Task 3.1: Complete MUSIC_MODULE.md 🟡

**File**: `data/features/music/MUSIC_MODULE.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 6 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 3.1** (MUSIC_MODULE.md questions)
2. [ ] **Update YAML** (`data/features/music/MUSIC_MODULE.yaml`):
   ```yaml
   status_design: ✅ Complete  # Change from 🟡
   overall_status: ✅ Complete

   # Add missing technical content
   lidarr_integration:
     sync_direction: two-way
     metadata_priority:
       - Lidarr (cached)
       - MusicBrainz
       - Last.fm

   supported_formats:
     - MP3
     - FLAC
     - AAC
     - OGG
     - ALAC

   metadata_sources:
     - name: MusicBrainz
       purpose: Album/artist metadata
       priority: 1
     - name: Last.fm
       purpose: User scrobbles, tags
       priority: 2
   ```

3. [ ] **Fill PLACEHOLDER Content**:
   - Replace `wiki_overview: "PLACEHOLDER: ..."`
   - Add specific Lidarr webhook handling
   - Add music file scanning logic
   - Add playlist support details

4. [ ] **Run Generation + Pipeline**
5. [ ] **Verify**: No PLACEHOLDER content remains

**Acceptance Criteria**:
- ✅ Lidarr integration fully documented
- ✅ Supported audio formats listed
- ✅ Metadata priority clear
- ✅ No PLACEHOLDER content

---

### Task 3.2: Complete AUDIOBOOK_MODULE.md 🟡

**File**: `data/features/audiobook/AUDIOBOOK_MODULE.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 5 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 3.2** (AUDIOBOOK_MODULE.md questions)
2. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete

   chaptarr_integration:
     sync_direction: two-way
     metadata_sync: Audnexus via Chaptarr

   supported_formats:
     - M4B (with chapters)
     - MP3 (multi-file)

   progress_tracking:
     resume: per-user, per-audiobook
     sync_to_chaptarr: true
   ```

3. [ ] **Fill PLACEHOLDER Content**
4. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Chaptarr integration documented
- ✅ Chapter marker handling described
- ✅ Progress tracking explained
- ✅ No PLACEHOLDER content

---

### Task 3.3: Complete BOOK_MODULE.md 🟡

**File**: `data/features/book/BOOK_MODULE.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 5 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 3.3** (BOOK_MODULE.md questions)
2. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete

   supported_formats:
     - EPUB
     - PDF
     - MOBI
     - AZW3

   reading_progress:
     tracking: page/chapter level
     sync: across devices
   ```

3. [ ] **Fill PLACEHOLDER Content**
4. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ eBook formats documented
- ✅ Reading progress tracking described
- ✅ Metadata sources clear
- ✅ No PLACEHOLDER content

---

### Task 3.4: Complete PODCASTS.md 🟡

**File**: `data/features/podcasts/PODCASTS.yaml`
**Status**: ✅ Complete (but verify no PLACEHOLDERs)
**Estimated Time**: 2 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Review Existing Content**: Check for PLACEHOLDER markers
2. [ ] **If PLACEHOLDERs exist**:
   - Answer questions
   - Fill content
   - Run generation
3. [ ] **Verify Complete**

---

## Priority 4: Pattern Documentation (Week 3)

### Task 4.1: Complete ARR_INTEGRATION.md 🟡

**File**: `data/patterns/ARR_INTEGRATION.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 4 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 4.1** (ARR_INTEGRATION.md questions)
2. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete

   webhook_events:
     - Download
     - Upgrade
     - Rename
     - Delete

   conflict_resolution:
     strategy: Arr metadata wins by default
     user_override: configurable priority
   ```

3. [ ] **Add Code Examples**:
   - Webhook handler implementation
   - Metadata sync logic
   - Conflict resolution code

4. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Webhook handling pattern clear
- ✅ Metadata sync logic documented
- ✅ Conflict resolution explained
- ✅ Code examples provided

---

### Task 4.2: Complete METADATA_ENRICHMENT.md 🟡

**File**: `data/patterns/METADATA_ENRICHMENT.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 4 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 4.2** (METADATA_ENRICHMENT.md questions)
2. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete

   priority_chain:
     1: Local cache (otter/rueidis)
     2: Arr services (Radarr, Sonarr, etc.)
     3: Internal services (Stash-App)
     4: External APIs (TMDb, StashDB, etc.)
     5: Background enrichment (River jobs)

   cache_strategy:
     metadata_ttl: 24h
     images_ttl: 7d
     search_ttl: 1h
   ```

3. [ ] **Add Code Examples**: Priority chain implementation
4. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Priority chain documented
- ✅ Cache strategy clear
- ✅ Background jobs explained
- ✅ Code examples provided

---

### Task 4.3: Complete WEBHOOK_PATTERNS.md 🟡

**File**: `data/patterns/WEBHOOK_PATTERNS.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 5 hours
**Dependencies**: None

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 4.3** (WEBHOOK_PATTERNS.md questions)
2. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete

   events:
     playback:
       - playback.started
       - playback.stopped
       - playback.paused
     library:
       - library.scan.started
       - library.scan.completed
       - library.item.added
       - library.item.updated
       - library.item.deleted

   security:
     hmac_signatures: true
     algorithm: HMAC-SHA256
   ```

3. [ ] **Add Event Schema**: JSON schema for each event type
4. [ ] **Run Generation + Pipeline**

**Acceptance Criteria**:
- ✅ Event catalog complete
- ✅ Event payloads documented
- ✅ Security (HMAC) explained
- ✅ Retry logic documented

---

## Priority 5: Remaining Technical Docs (Week 3-4)

### Task 5.1: Complete AUDIO_STREAMING.md 🔴

**File**: `data/technical/AUDIO_STREAMING.yaml`
**Status**: 🔴 Not Started → ✅ Complete
**Estimated Time**: 6 hours

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 2.4** (AUDIO_STREAMING.md questions)
2. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete
   overall_status: ✅ Complete

   streaming_protocol: HLS
   supported_codecs:
     - AAC
     - MP3
     - FLAC
     - Opus

   progress_tracking:
     interval: 10s
     storage: per-user, per-track
   ```

3. [ ] **Fill PLACEHOLDER Content**
4. [ ] **Run Generation + Pipeline**

---

### Task 5.2: Complete EMAIL.md 🟡

**File**: `data/technical/EMAIL.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 4 hours

**Action Items**:

1. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete

   smtp_configuration:
     library: go-mail
     tls_required: true
     auth_methods:
       - PLAIN
       - LOGIN
   ```

2. [ ] **Fill PLACEHOLDER Content**
3. [ ] **Run Generation + Pipeline**

---

### Task 5.3: Complete NOTIFICATIONS.md 🟡

**File**: `data/technical/NOTIFICATIONS.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 5 hours

**Action Items**:

1. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete

   notification_channels:
     - Email (SMTP)
     - Push (Firebase FCM)
     - Webhooks
   ```

2. [ ] **Fill PLACEHOLDER Content**
3. [ ] **Run Generation + Pipeline**

---

### Task 5.4: Complete WEBHOOKS.md 🟡

**File**: `data/technical/WEBHOOKS.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 3 hours

**Action Items**:

1. [ ] **Update YAML** (add event schema)
2. [ ] **Fill PLACEHOLDER Content**
3. [ ] **Run Generation + Pipeline**

---

### Task 5.5: Complete WEBSOCKETS.md 🟡

**File**: `data/technical/WEBSOCKETS.yaml`
**Status**: 🟡 In Progress → ✅ Complete
**Estimated Time**: 4 hours

**Action Items**:

1. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete

   websocket_library: gobwas/ws
   endpoints:
     - /api/v1/ws/playback
     - /api/v1/ws/library
   ```

2. [ ] **Fill PLACEHOLDER Content**
3. [ ] **Run Generation + Pipeline**

---

### Task 5.6: Complete OFFLOADING.md 🔴

**File**: `data/technical/OFFLOADING.yaml`
**Status**: 🔴 Not Started → ✅ Complete
**Estimated Time**: 6 hours

**Action Items**:

1. [ ] **Update YAML**:
   ```yaml
   status_design: ✅ Complete
   overall_status: ✅ Complete

   offloading_targets:
     - Blackbeard (transcoding)
     - Local FFmpeg (fallback)
   ```

2. [ ] **Fill PLACEHOLDER Content**
3. [ ] **Run Generation + Pipeline**

---

## Priority 6: Create Missing Documentation (Week 4)

### Task 6.1: Create COLLECTIONS.md 🟡

**File**: `data/features/shared/COLLECTIONS.yaml` (NEW)
**Status**: 🔴 Does Not Exist → ✅ Complete
**Estimated Time**: 5 hours

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 5.1** (COLLECTIONS.md questions)
2. [ ] **Create YAML File**:
   ```yaml
   doc_title: Collections & Playlists
   doc_category: feature
   created_date: '2026-01-31'
   overall_status: ✅ Complete
   status_design: ✅

   technical_summary: >
     > User-created collections of content across types

   feature_name: Collections
   content_types:
     - Movies
     - TV Shows
     - Music
     - All media types

   collection_types:
     - Manual (user-curated)
     - Smart (filter-based)
   ```

3. [ ] **Run Generation + Pipeline**

---

### Task 6.2: Create TRANSCODING.md Service Doc 🟢

**File**: `data/services/TRANSCODING.yaml` (NEW)
**Status**: 🔴 Does Not Exist → ✅ Complete
**Estimated Time**: 4 hours

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 5.2** (TRANSCODING.md questions)
2. [ ] **Create YAML File**:
   ```yaml
   doc_title: Transcoding Service
   doc_category: service
   created_date: '2026-01-31'
   overall_status: ✅ Complete
   status_design: ✅

   service_name: Transcoding Service
   package_path: internal/service/transcoding
   fx_module: transcoding.Module

   technical_summary: >
     > On-demand video/audio transcoding service
   ```

3. [ ] **Run Generation + Pipeline**

---

### Task 6.3: Create EPG.md Service Doc 🟢

**File**: `data/services/EPG.yaml` (NEW)
**Status**: 🔴 Does Not Exist → ✅ Complete
**Estimated Time**: 3 hours

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 5.3** (EPG.md questions)
2. [ ] **Create YAML File**:
   ```yaml
   doc_title: EPG Service
   doc_category: service
   created_date: '2026-01-31'
   overall_status: ✅ Complete
   status_design: ✅

   service_name: EPG Service
   package_path: internal/service/epg
   fx_module: epg.Module

   technical_summary: >
     > Electronic Program Guide for Live TV
   ```

3. [ ] **Run Generation + Pipeline**

---

### Task 6.4: Create OBSERVABILITY.md Pattern Doc 🟢

**File**: `data/patterns/OBSERVABILITY.yaml` (NEW)
**Status**: 🔴 Does Not Exist → ✅ Complete
**Estimated Time**: 5 hours

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 5.4** (OBSERVABILITY.md questions)
2. [ ] **Create YAML File**:
   ```yaml
   doc_title: Observability Pattern
   doc_category: pattern
   created_date: '2026-01-31'
   overall_status: ✅ Complete
   status_design: ✅

   technical_summary: >
     > Metrics, tracing, and logging patterns with Prometheus, OpenTelemetry, Loki
   ```

3. [ ] **Run Generation + Pipeline**

---

### Task 6.5: Create TESTING.md Pattern Doc 🟢

**File**: `data/patterns/TESTING.yaml` (NEW)
**Status**: 🔴 Does Not Exist → ✅ Complete
**Estimated Time**: 4 hours

**Action Items**:

1. [ ] **Answer QUESTIONS.md Section 5.5** (TESTING.md questions)
2. [ ] **Create YAML File**:
   ```yaml
   doc_title: Testing Patterns
   doc_category: pattern
   created_date: '2026-01-31'
   overall_status: ✅ Complete
   status_design: ✅

   technical_summary: >
     > Table-driven tests, testcontainers, mocking patterns
   ```

3. [ ] **Run Generation + Pipeline**

---

## Priority 7: Research Documentation (Ongoing / Low Priority)

### Task 7.1: Create USER_PAIN_POINTS_RESEARCH.md ⚪

**File**: `data/research/USER_PAIN_POINTS_RESEARCH.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 8 hours (research-intensive)

**Action Items**:

1. [ ] **Research Pain Points**:
   - Reddit: r/jellyfin, r/Plex, r/selfhosted
   - GitHub Issues: Jellyfin, Plex complaints, Emby
   - Forums: Jellyfin forum, Emby forum

2. [ ] **Document Findings** in YAML
3. [ ] **Run Generation + Pipeline**

---

### Task 7.2: Create UX_UI_RESOURCES.md ⚪

**File**: `data/research/UX_UI_RESOURCES.yaml`
**Status**: 🔴 Not Started
**Estimated Time**: 4 hours

**Action Items**:

1. [ ] **Catalog Resources**:
   - Design systems (Radix UI, Material Design)
   - UI kits (Figma community)
   - Color schemes
   - Icon libraries

2. [ ] **Document in YAML**
3. [ ] **Run Generation + Pipeline**

---

## Priority 8: PLACEHOLDER Cleanup (Ongoing)

### Task 8.1: Replace All PLACEHOLDER Content

**Files**: All YAML files with `PLACEHOLDER:` markers
**Status**: Ongoing
**Estimated Time**: 2 hours per file × 36 files = 72 hours

**Action Items**:

1. [ ] **Search for PLACEHOLDERs**:
   ```bash
   grep -r "PLACEHOLDER" data --include="*.yaml"
   ```

2. [ ] **For each PLACEHOLDER**:
   - Research actual content
   - Replace with real content
   - Run generation for that file

3. [ ] **Verify Wiki Output**: Ensure user-friendly language

**Target**: Zero PLACEHOLDER markers before M1 release

---

## Workflow Summary

### For Each Documentation Task:

1. **Preparation**:
   - [ ] Review ANALYSIS.md for context
   - [ ] Answer relevant QUESTIONS.md questions
   - [ ] Research if needed (external APIs, tools, etc.)

2. **YAML Update**:
   - [ ] Edit `data/{category}/{FILE}.yaml`
   - [ ] Add/update technical content
   - [ ] Replace PLACEHOLDER values
   - [ ] Update status: `status_design: ✅`

3. **Generation**:
   - [ ] Run: `python scripts/automation/batch_regenerate.py`
   - [ ] Verify: `docs/dev/design/{category}/{FILE}.md`
   - [ ] Verify: `docs/wiki/{category}/{file}.md`

4. **Pipeline**:
   - [ ] Run: `./scripts/doc-pipeline.sh --apply`
   - [ ] Verify: INDEX.md files updated
   - [ ] Verify: Breadcrumbs added
   - [ ] Verify: Links validated

5. **Quality Check**:
   - [ ] No PLACEHOLDER content
   - [ ] No TODO comments (or tracked separately)
   - [ ] All examples tested
   - [ ] Links valid
   - [ ] Wiki version user-friendly

6. **Commit**:
   ```bash
   git add data/{category}/{FILE}.yaml docs/dev/design/{category}/{FILE}.md docs/wiki/{category}/{file}.md
   git commit -m "docs: complete {FILE} documentation"
   ```

---

## Progress Tracking

### Week 1 Checklist (Operations - CRITICAL)
- [ ] Task 1.1: DEVELOPMENT.md
- [ ] Task 1.2: SETUP.md
- [ ] Task 1.3: GITFLOW.md
- [ ] Task 1.4: BEST_PRACTICES.md
- [ ] Task 1.5: VERSIONING.md
- [ ] Task 1.6: BRANCH_PROTECTION.md
- [ ] Task 1.7: REVERSE_PROXY.md
- [ ] Task 1.8: DATABASE_AUTO_HEALING.md

**Success Metric**: 8/8 Operations docs complete (0% → 100%)

### Week 2 Checklist (Technical + Content Modules)
- [ ] Task 2.1: API.md
- [ ] Task 2.2: CONFIGURATION.md
- [ ] Task 2.3: FRONTEND.md
- [ ] Task 3.1: MUSIC_MODULE.md
- [ ] Task 3.2: AUDIOBOOK_MODULE.md
- [ ] Task 3.3: BOOK_MODULE.md
- [ ] Task 3.4: PODCASTS.md (verify)

**Success Metric**: 7/7 docs complete

### Week 3 Checklist (Patterns + Technical)
- [ ] Task 4.1: ARR_INTEGRATION.md
- [ ] Task 4.2: METADATA_ENRICHMENT.md
- [ ] Task 4.3: WEBHOOK_PATTERNS.md
- [ ] Task 5.1: AUDIO_STREAMING.md
- [ ] Task 5.2: EMAIL.md
- [ ] Task 5.3: NOTIFICATIONS.md
- [ ] Task 5.4: WEBHOOKS.md
- [ ] Task 5.5: WEBSOCKETS.md
- [ ] Task 5.6: OFFLOADING.md

**Success Metric**: 9/9 docs complete

### Week 4 Checklist (Missing Docs + Cleanup)
- [ ] Task 6.1: COLLECTIONS.md
- [ ] Task 6.2: TRANSCODING.md
- [ ] Task 6.3: EPG.md
- [ ] Task 6.4: OBSERVABILITY.md
- [ ] Task 6.5: TESTING.md
- [ ] Task 8.1: PLACEHOLDER cleanup (ongoing)

**Success Metric**: 5/5 new docs created, 36 PLACEHOLDERs resolved

### Ongoing (Research - Low Priority)
- [ ] Task 7.1: USER_PAIN_POINTS_RESEARCH.md
- [ ] Task 7.2: UX_UI_RESOURCES.md

---

## Completion Criteria

Documentation is **COMPLETE** when:

- ✅ All 46 tasks checked off
- ✅ Zero PLACEHOLDER markers in YAML files
- ✅ All `status_design: ✅` in YAML
- ✅ All docs generate without errors
- ✅ All docs pass validation (doc-pipeline.sh)
- ✅ Wiki versions user-friendly
- ✅ All links valid

**Target Completion**: 4 weeks from start

---

**Related Files**:
- [ANALYSIS.md](ANALYSIS.md) - Gap analysis
- [QUESTIONS.md](QUESTIONS.md) - Clarification questions
