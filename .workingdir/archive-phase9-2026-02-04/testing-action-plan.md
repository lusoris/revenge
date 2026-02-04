# Testing Action Plan - Systematisches Bug Finding

**Created**: 2026-02-03
**Strategy**: Docker → Live Tests → Integration Tests → Fixes → Quality Gates

---

## 🎯 Sofort-Aktionen (Jetzt)

### 1. Linting Pass ⏳

**Warum jetzt?**
- Findet Code Smells, Design Flaws, potentielle Bugs
- Schnell durchzuführen
- Keine Infrastruktur nötig

**Aktion**:
```bash
# Install golangci-lint if needed
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run on entire codebase
golangci-lint run --config .golangci.yml ./... > .workingdir/lint-results.txt 2>&1

# Analyze results
cat .workingdir/lint-results.txt
```

**Erwartete Findings**:
- Unused variables/imports
- Error handling issues
- Code duplication
- Security issues (weak crypto, SQL injection risks)
- Performance problems

---

### 2. Security Scan mit Snyk ⏳

**Warum wichtig?**
- `.github/instructions/snyk_rules.instructions.md` schreibt das vor!
- "Always run snyk_code_scan for new first party code"
- Wir haben 2 neue Services + Fixes

**Aktion**:
```bash
# Scan crypto service (NEW)
snyk code test internal/crypto/

# Scan modified services
snyk code test internal/service/user/
snyk code test internal/service/auth/

# Full scan
snyk test --severity-threshold=medium
```

**Follow Snyk Rules**:
1. Run scan on NEW code ✅
2. Fix issues found ⏳
3. Rescan after fixes ⏳
4. Repeat until clean ⏳

---

### 3. Docker Build Testing Setup ⏳

**Ziel**: Build schlägt fehl wenn Tests/Linting/Security failen

**Dateien erstellen**:

#### `Dockerfile.test` (Multi-Stage Testing)
```dockerfile
# Base Stage
FROM golang:1.23-alpine AS base
WORKDIR /app
RUN apk add --no-cache git make

# Dependencies
FROM base AS deps
COPY go.mod go.sum ./
RUN go mod download

# Source Code
FROM deps AS source
COPY . .

# Testing Stage
FROM source AS testing
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go test -v -race -coverprofile=/tmp/coverage.out ./...
RUN golangci-lint run ./...

# Integration Tests
FROM testing AS integration
# Start dependencies in background (needs docker-compose)
RUN go test -v -tags=integration ./tests/integration/...

# Security Scanning
FROM testing AS security
RUN apk add --no-cache npm
RUN npm install -g snyk
RUN snyk test --severity-threshold=high || true

# Coverage Report
FROM testing AS coverage
RUN go tool cover -func=/tmp/coverage.out | tee /tmp/coverage-summary.txt
RUN go tool cover -html=/tmp/coverage.out -o /tmp/coverage.html

# Production Build
FROM source AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /revenge ./cmd/revenge

# Final Production Image
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates
COPY --from=builder /revenge /usr/local/bin/revenge
CMD ["revenge"]
```

#### `docker-compose.test.yml`
```yaml
version: '3.8'

services:
  # Main test runner
  test:
    build:
      context: .
      dockerfile: Dockerfile.test
      target: testing
    depends_on:
      postgres:
        condition: service_healthy
      dragonfly:
        condition: service_started
      typesense:
        condition: service_started
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=test
      - DB_PASSWORD=test
      - DB_NAME=revenge_test
      - CACHE_URL=redis://dragonfly:6379
      - SEARCH_URL=http://typesense:8108
      - SEARCH_API_KEY=test-key
    volumes:
      - ./coverage:/coverage
    command: |
      sh -c "
        echo '=== Running Unit Tests ===' &&
        go test -v -race -coverprofile=/coverage/coverage.out ./... &&
        echo '=== Running Integration Tests ===' &&
        go test -v -tags=integration ./tests/integration/... &&
        echo '=== Generating Coverage Report ===' &&
        go tool cover -func=/coverage/coverage.out &&
        echo '=== Done ==='
      "

  # Integration test runner
  integration:
    build:
      context: .
      dockerfile: Dockerfile.test
      target: integration
    depends_on:
      postgres:
        condition: service_healthy
      dragonfly:
        condition: service_started
      typesense:
        condition: service_started
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=test
      - DB_PASSWORD=test
      - DB_NAME=revenge_test
      - CACHE_URL=redis://dragonfly:6379
      - SEARCH_URL=http://typesense:8108
      - SEARCH_API_KEY=test-key

  # Database for testing
  postgres:
    image: postgres:18-alpine
    environment:
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
      POSTGRES_DB: revenge_test
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U test"]
      interval: 5s
      timeout: 5s
      retries: 5
    tmpfs:
      - /var/lib/postgresql/data

  # Cache for testing
  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly:latest
    command: --logtostderr

  # Search for testing
  typesense:
    image: typesense/typesense:0.25.2
    command: '--data-dir=/data --api-key=test-key --enable-cors'
    tmpfs:
      - /data
```

**Usage**:
```bash
# Run all tests in Docker
docker-compose -f docker-compose.test.yml up --abort-on-container-exit test

# Run only integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit integration

# Cleanup
docker-compose -f docker-compose.test.yml down -v

# Build testing image only
docker build --target testing -t revenge:test -f Dockerfile.test .

# Build security stage
docker build --target security -t revenge:secure -f Dockerfile.test .
```

---

## 📋 Service Testing Roadmap

### Phase 1: Core Services (Done ✅)
- [x] User Service (11 tests)
- [x] Auth Service (12 tests)

### Phase 2: Session & Settings (Next)
- [ ] Session Service
  - [ ] Create session
  - [ ] Get active sessions
  - [ ] Revoke session
  - [ ] Revoke all sessions
  - [ ] Session expiry
  - [ ] Concurrent sessions limit
  - [ ] Device tracking

- [ ] Settings Service
  - [ ] Get user settings
  - [ ] Update user settings
  - [ ] Get server settings
  - [ ] Update server settings (admin only)
  - [ ] Settings validation
  - [ ] Default settings

### Phase 3: Authorization (Week 2)
- [ ] RBAC Service
  - [ ] Create role
  - [ ] Assign role to user
  - [ ] Check permission
  - [ ] Policy enforcement
  - [ ] Role hierarchy
  - [ ] Custom permissions
  - [ ] Audit logging

### Phase 4: Security Features (Week 2)
- [ ] API Keys Service
  - [ ] Generate API key
  - [ ] List API keys
  - [ ] Revoke API key
  - [ ] API key authentication
  - [ ] Scope/permissions
  - [ ] Rate limiting per key

### Phase 5: Content & Activity (Week 3)
- [ ] Library Service
  - [ ] Add item to library
  - [ ] Get library items
  - [ ] Update item metadata
  - [ ] Delete item
  - [ ] Search library
  - [ ] Pagination

- [ ] Activity Service
  - [ ] Log activity
  - [ ] Get user activity
  - [ ] Activity feed
  - [ ] Activity filtering
  - [ ] Retention policy

### Phase 6: External Integration (Week 3)
- [ ] OIDC Service
  - [ ] Google OAuth flow
  - [ ] GitHub OAuth flow
  - [ ] Token exchange
  - [ ] User linking
  - [ ] Profile sync

---

## 🔍 Design Flaws Investigation Plan

### Investigation 1: Email Verification Flow

**Hypothesis**: Email verification not fully implemented

**Tests to Write**:
```go
func TestEmailVerification_Complete(t *testing.T) {
    // 1. Register user
    user := registerUser(t, "test@example.com")
    assert.False(t, user.EmailVerified)

    // 2. Get verification token (how?)
    // TODO: Check if email service exists

    // 3. Verify email
    // TODO: Find verify endpoint

    // 4. Confirm verified
    user = getUser(t, user.ID)
    assert.True(t, user.EmailVerified)
}
```

**Questions**:
- [ ] Is there an email service?
- [ ] Where are verification tokens stored?
- [ ] What endpoint verifies emails?
- [ ] Does registration send emails?

---

### Investigation 2: User Activation

**Hypothesis**: Users need admin approval or email verification

**Tests**:
```go
func TestUserActivation_AutoActivate(t *testing.T) {
    // Should users auto-activate after email verification?
    user := registerAndVerifyEmail(t)
    assert.True(t, user.IsActive)
}

func TestUserActivation_AdminApproval(t *testing.T) {
    // Or do admins need to approve?
    user := registerUser(t)
    assert.False(t, user.IsActive)

    adminService.ActivateUser(ctx, user.ID)

    user = getUser(t, user.ID)
    assert.True(t, user.IsActive)
}
```

**Questions**:
- [ ] What activates a user?
- [ ] Can inactive users login?
- [ ] Is there an admin activation endpoint?

---

### Investigation 3: Rate Limiting

**Hypothesis**: No rate limiting implemented (security risk!)

**Live Test**:
```bash
# Attempt 100 rapid registrations
for i in {1..100}; do
  curl -X POST http://localhost:8096/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"spam$i\",\"email\":\"spam$i@test.com\",\"password\":\"Test123!\"}" &
done
wait

# Expected: Should be rate limited after N attempts
# Actual: ?
```

**Questions**:
- [ ] Is rate limiting implemented?
- [ ] What endpoints need rate limiting?
- [ ] Should we add middleware?

---

## 🚀 Execution Order

### TODAY (2026-02-03)

1. **Linting** (15 min)
   ```bash
   golangci-lint run ./... > .workingdir/lint-results.txt 2>&1
   ```

2. **Security Scan** (10 min)
   ```bash
   snyk code test internal/crypto/ > .workingdir/snyk-crypto.txt
   snyk code test internal/service/ > .workingdir/snyk-services.txt
   ```

3. **Analyze Results** (20 min)
   - Review lint findings
   - Review security issues
   - Categorize by severity
   - Create fix plan

4. **Docker Testing Setup** (30 min)
   - Create `Dockerfile.test`
   - Create `docker-compose.test.yml`
   - Test build pipeline
   - Document usage

5. **Fix High-Priority Issues** (1 hour)
   - Fix critical security issues
   - Fix major linting issues
   - Re-run scans
   - Verify fixes

6. **Session Service Tests** (1 hour)
   - Write 7 integration tests
   - Run against live DB
   - Document findings
   - Fix bugs if found

### TONIGHT (Commits & Documentation)

7. **Documentation Update** (15 min)
   - Update `.workingdir/` with findings
   - Document all bugs found
   - Update testing matrix

8. **Git Commits** (30 min)
   ```bash
   # Crypto service
   git add internal/crypto/
   git commit -m "feat(crypto): add shared password hashing service"

   # Service refactoring
   git add internal/service/user/ internal/service/auth/
   git commit -m "refactor(services): migrate to shared crypto service"

   # Tests
   git add tests/integration/service/
   git commit -m "test(services): add user and auth integration tests"

   # Bug fixes
   git add internal/service/auth/jwt.go
   git commit -m "fix(auth): use millisecond precision for JWT timestamps

Fixes token duplication when generated within same second.
Refs: Bug #29"

   # Quality improvements
   git add Dockerfile.test docker-compose.test.yml
   git commit -m "ci: add Docker-based testing pipeline"

   # Documentation
   git add .workingdir/
   git commit -m "docs: add testing workflow and bug reports"
   ```

---

## 📊 Success Criteria

**Phase 1 Complete When**:
- ✅ Linting: 0 critical issues
- ✅ Security: 0 high/critical vulnerabilities
- ✅ Docker tests: All stages build successfully
- ✅ Session service: All tests passing
- ✅ Documentation: All findings documented
- ✅ Commits: Clean git history

**Overall Goal**:
- 🎯 All 9 services tested
- 🎯 >80% code coverage
- 🎯 Zero security issues
- 🎯 Automated quality pipeline
- 🎯 Complete documentation

---

**Status**: Ready to execute
**Next Action**: Run linting pass
