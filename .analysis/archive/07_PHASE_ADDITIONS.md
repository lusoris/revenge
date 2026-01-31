# Additional Phase Tasks

**Date**: 2026-01-31
**Purpose**: Address user feedback about detailed TODOs and missing design docs

---

## Phase 1 Additions

### Task 1.5: Analyze Missing Design Docs

**Objective**: Identify gaps in design documentation coverage

**Process**:
1. Read through all existing design docs
2. Compare with planned features (from VERSIONING.md)
3. Identify missing docs
4. Categorize by priority (MVP-critical, v0.x, v1.0)

**Expected Gaps** (examples):
- Playback service (exists but may need detail)
- Scrobbling integrations (Trakt, Last.fm, etc.)
- Analytics service (may need expansion)
- Notification service (may need expansion)
- Plugin architecture (deferred but needs doc)

---

### Task 1.6: Scaffold Missing Design Docs

**Objective**: Create placeholder docs for all missing designs

**Template for Scaffolded Docs**:
```markdown
# [Feature Name]

> **Status**: 🔴 PLANNED - Not yet designed

## Overview

[Brief description of what this will be]

## Scope

### MVP (v0.3.x)
- Not in MVP scope

### Post-MVP
- Planned for: [version]

## Design TODO

- [ ] Define architecture
- [ ] Define database schema
- [ ] Define API endpoints
- [ ] Define integrations
- [ ] Define testing strategy

## Related

- **Depends on**: [other features]
- **Blocks**: [other features]
- **Related to**: [similar features]

---

**Note**: This is a placeholder. Design work pending.
```

**Process**:
1. Create doc with scaffold
2. Add to appropriate category in `docs/dev/design/`
3. Link from DESIGN_INDEX.md
4. Add to status tracking

---

### Task 1.7: Create Detailed Phase TODOs

**Objective**: Break down each roadmap phase into actionable tasks

**For Each Phase** (v0.1.x, v0.2.x, v0.3.x, etc.):

1. **Infrastructure Tasks**
   - Database setup
   - Cache setup
   - Job queue setup
   - etc.

2. **Service Implementation**
   - One TODO per service
   - Break down complex services into subtasks

3. **Integration Tasks**
   - One TODO per external integration
   - Include testing

4. **Frontend Tasks**
   - UI components
   - Pages/routes
   - State management

5. **Testing Tasks**
   - Unit tests
   - Integration tests
   - E2E tests

6. **Documentation Tasks**
   - Update design docs with impl details
   - Generate wiki docs
   - API documentation

**Format** (detailed example for v0.1.x):
```markdown
# v0.1.x - Core Foundation TODO

## Infrastructure Setup

### Database (PostgreSQL 18+)
- [ ] Install PostgreSQL 18 locally
- [ ] Create development database
- [ ] Set up connection pooling (pgxpool)
- [ ] Configure pgx with prepared statements
- [ ] Test connection and basic queries

### Migrations
- [ ] Set up golang-migrate
- [ ] Create initial schema migration
- [ ] Create rollback migrations
- [ ] Test migration up/down
- [ ] Document migration workflow

### Cache (Dragonfly)
- [ ] Install Dragonfly locally (or Docker)
- [ ] Set up rueidis client
- [ ] Configure connection pooling
- [ ] Test basic set/get operations
- [ ] Implement cache invalidation patterns

### Job Queue (River)
- [ ] Set up River with PostgreSQL
- [ ] Define job types
- [ ] Implement job handlers
- [ ] Configure workers
- [ ] Test job scheduling and execution

## Core Services

### Auth Service
- [ ] Define Auth interface (`internal/service/auth/service.go`)
- [ ] Implement password hashing (bcrypt)
- [ ] Implement JWT generation/validation
- [ ] Implement session management
- [ ] Unit tests (80%+ coverage)
- [ ] Integration tests with database

### User Service
- [ ] Define User interface
- [ ] Implement CRUD operations
- [ ] Implement user repository (sqlc)
- [ ] Add validation (go-playground/validator)
- [ ] Unit tests
- [ ] Integration tests

### Session Service
- [ ] Define Session interface
- [ ] Implement session storage (cache + DB)
- [ ] Implement session expiration
- [ ] Implement session refresh
- [ ] Unit tests
- [ ] Integration tests

### RBAC Service
- [ ] Set up Casbin
- [ ] Define RBAC model
- [ ] Implement policy enforcement
- [ ] Integrate with Auth middleware
- [ ] Unit tests
- [ ] Integration tests

## Library Management

### File Scanner
- [ ] Implement fsnotify watcher
- [ ] Implement directory traversal
- [ ] Implement file type detection
- [ ] Implement media file parsing
- [ ] Unit tests
- [ ] Integration tests

### Library Service
- [ ] Define Library interface
- [ ] Implement library CRUD
- [ ] Implement scan triggering
- [ ] Implement scan status tracking
- [ ] Unit tests
- [ ] Integration tests

## API Layer

### OpenAPI Specification
- [ ] Define API v1 spec (YAML)
- [ ] Define authentication endpoints
- [ ] Define user management endpoints
- [ ] Define library endpoints
- [ ] Validate spec with tools

### ogen Generation
- [ ] Generate server code with ogen
- [ ] Implement handlers
- [ ] Add authentication middleware
- [ ] Add RBAC middleware
- [ ] Add error handling
- [ ] Add request validation

## Observability

### Logging (slog)
- [ ] Set up structured logging
- [ ] Configure log levels
- [ ] Add contextual logging
- [ ] Implement log rotation
- [ ] Test log output

### Metrics (Prometheus)
- [ ] Set up Prometheus metrics
- [ ] Add service metrics
- [ ] Add HTTP metrics
- [ ] Add database metrics
- [ ] Create Grafana dashboards

### Tracing (OpenTelemetry)
- [ ] Set up OTLP exporter
- [ ] Add trace spans to services
- [ ] Add trace propagation
- [ ] Configure sampling
- [ ] Test with Jaeger

### Health Checks
- [ ] Implement /health endpoint
- [ ] Check database connectivity
- [ ] Check cache connectivity
- [ ] Check external dependencies
- [ ] Add readiness checks

## Testing

### Unit Tests
- [ ] Auth service (80%+ coverage)
- [ ] User service (80%+ coverage)
- [ ] Session service (80%+ coverage)
- [ ] RBAC service (80%+ coverage)
- [ ] Library service (80%+ coverage)

### Integration Tests
- [ ] Database integration tests (testcontainers)
- [ ] Cache integration tests
- [ ] End-to-end API tests
- [ ] Performance tests

### Test Infrastructure
- [ ] Set up testcontainers
- [ ] Set up test fixtures
- [ ] Set up test data generators
- [ ] Set up CI test runner

## Documentation

### Design Docs
- [ ] Update Auth design with impl details
- [ ] Update User design with impl details
- [ ] Update Library design with impl details
- [ ] Add code examples

### API Documentation
- [ ] Generate API docs from OpenAPI
- [ ] Add usage examples
- [ ] Document authentication flow
- [ ] Document error responses

## DevOps

### Docker
- [ ] Create Dockerfile (multi-stage)
- [ ] Create docker-compose.yml (dev)
- [ ] Test local Docker build
- [ ] Optimize image size

### CI/CD
- [ ] Set up GitHub Actions for tests
- [ ] Set up code coverage reporting
- [ ] Set up linting (golangci-lint)
- [ ] Set up security scanning (gosec)

## Exit Criteria

- [ ] All services implemented and tested
- [ ] 80%+ test coverage achieved
- [ ] API responds to requests
- [ ] Health checks pass
- [ ] No critical linting errors
- [ ] Documentation complete
```

---

## Phase 1 Updated Deliverables

With these additions, Phase 1 now includes:

1. ✅ MVP_DEFINITION.md
2. ✅ IMPLEMENTATION_ROADMAP.md
3. ✅ Milestone TODO files (5+)
4. ✅ **NEW**: Gap analysis of missing design docs
5. ✅ **NEW**: Scaffolded docs for all gaps
6. ✅ **NEW**: Detailed phase TODOs (breaking down each milestone)
7. ✅ SOURCE_OF_TRUTH.md updates

---

## Updated Phase 1 Timeline

- **Original estimate**: 1-2 days
- **With additions**: 2-3 days

**Breakdown**:
- Day 1: MVP definition + roadmap
- Day 2: Gap analysis + detailed TODOs
- Day 3: Scaffold missing docs + SOT updates

---

## Questions for User

1. **Should scaffolded docs be part of Phase 1 or separate phase?**
   - Option A: Include in Phase 1 (create stubs now)
   - Option B: Create separate "Phase 1.5: Design Completion" (fill in stubs)

2. **How detailed should phase TODOs be?**
   - Option A: Very detailed (like example above, 100+ tasks per phase)
   - Option B: Moderate (20-30 main tasks, sub-tasks as needed)
   - Option C: High-level (5-10 major tasks, detailed during implementation)

3. **Should we prioritize filling in missing design docs?**
   - Some designs are placeholders (Analytics, Notification, etc.)
   - Should we complete these before implementation?
   - Or create detailed designs as-needed during development?

---

**Recommendation**:
- Create scaffolds in Phase 1 (low effort)
- Write detailed TODOs (helps planning)
- Fill in design docs as-needed (avoid over-planning)

**This allows**:
- Clear visibility of what's missing
- Structured TODOs for implementation
- Flexibility to refine during development
