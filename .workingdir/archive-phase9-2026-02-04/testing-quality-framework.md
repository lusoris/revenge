# Testing & Quality Assurance Framework

**Created:** 2026-02-03
**Author:** AI Assistant
**Purpose:** Systematic approach to test live system, find bugs, and improve code quality

---

## Table of Contents

1. [Overview](#overview)
2. [Framework Architecture](#framework-architecture)
3. [Benchmark Testing](#benchmark-testing)
4. [Security Testing](#security-testing)
5. [Live System Testing](#live-system-testing)
6. [Integration with CI/CD](#integration-with-cicd)
7. [Workflow Steps](#workflow-steps)
8. [Templates](#templates)

---

## Overview

### Philosophy

> **"Test the running image → Find bugs → Create tests → Fix → Lint → Repeat"**

This framework enables **backwards development**:
1. Start with live, running Docker container
2. Test actual behavior (not assumptions)
3. Discover bugs through real usage patterns
4. Create regression tests for discovered issues
5. Fix bugs with confidence
6. Ensure code quality through linting
7. Document everything for spec updates

### Benefits

- ✅ **Real-world validation**: Test actual running system, not mocks
- ✅ **Automated bug discovery**: Find issues before users do
- ✅ **Performance awareness**: Know when code gets slower
- ✅ **Security-first**: Catch vulnerabilities early
- ✅ **Living documentation**: Tests document expected behavior
- ✅ **Spec alignment**: Track changes needed in specifications

---

## Framework Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Live Docker Container                     │
│                     (revenge-dev:8096)                       │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    Testing Layers                            │
├─────────────────────────────────────────────────────────────┤
│  1. HTTP API Tests        → Real endpoints                   │
│  2. Integration Tests     → Service + DB                     │
│  3. Benchmark Tests       → Performance                      │
│  4. Security Scans        → Vulnerabilities                  │
│  5. Unit Tests           → Component logic                   │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    Bug Discovery                             │
│  → Document in .workingdir/bugs/                            │
│  → Create regression test                                    │
│  → Fix implementation                                        │
│  → Verify fix with test                                      │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    Quality Checks                            │
│  → Linting (golangci-lint)                                  │
│  → Security (gosec)                                         │
│  → Formatting (gofmt)                                       │
│  → Code coverage                                            │
└─────────────────────────────────────────────────────────────┘
```

---

## Benchmark Testing

### Purpose

Benchmarks measure **performance characteristics**:
- Operations per second
- Memory allocations
- CPU usage patterns
- Latency distributions

### Go Benchmark Basics

```go
// File: internal/crypto/password_bench_test.go
package crypto_test

import (
	"testing"
	"github.com/lusoris/revenge/internal/crypto"
)

// BenchmarkHashPassword measures password hashing performance
func BenchmarkHashPassword(b *testing.B) {
	hasher, _ := crypto.NewPasswordHasher(nil)

	// Reset timer to exclude setup
	b.ResetTimer()

	// Run b.N iterations (automatically determined)
	for b.Loop() {
		_, _ = hasher.HashPassword("TestPassword123!")
	}
}

// BenchmarkHashPasswordParallel tests concurrent performance
func BenchmarkHashPasswordParallel(b *testing.B) {
	hasher, _ := crypto.NewPasswordHasher(nil)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = hasher.HashPassword("TestPassword123!")
		}
	})
}

// BenchmarkVerifyPassword measures verification speed
func BenchmarkVerifyPassword(b *testing.B) {
	hasher, _ := crypto.NewPasswordHasher(nil)
	hash, _ := hasher.HashPassword("TestPassword123!")

	b.ResetTimer()
	for b.Loop() {
		_ = hasher.VerifyPassword("TestPassword123!", hash)
	}
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkHashPassword ./internal/crypto

# With memory statistics
go test -bench=. -benchmem ./internal/crypto

# Compare before/after performance
go test -bench=. -benchmem ./internal/crypto > old.txt
# Make changes...
go test -bench=. -benchmem ./internal/crypto > new.txt
benchstat old.txt new.txt
```

### Benchmark Output Interpretation

```
BenchmarkHashPassword-8          100  10234567 ns/op  4096 B/op  32 allocs/op
                     │           │    │               │          │
                     │           │    │               │          └─ Allocations per op
                     │           │    │               └─ Bytes allocated per op
                     │           │    └─ Nanoseconds per operation
                     │           └─ Number of iterations
                     └─ GOMAXPROCS
```

### What to Benchmark

| Component | What to Measure | Why |
|-----------|-----------------|-----|
| **Crypto** | Hash/Verify speed | Security vs performance trade-off |
| **Database** | Query execution | Identify slow queries |
| **Cache** | Get/Set operations | Ensure cache is faster than DB |
| **JWT** | Token generation | Auth bottleneck prevention |
| **JSON** | Marshal/Unmarshal | API response time impact |
| **Validation** | Input sanitization | Security overhead |

---

## Security Testing

### gosec - Go Security Checker

gosec performs **static analysis** to find security vulnerabilities.

#### Installation

```bash
# Install gosec
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin latest

# Verify installation
gosec --version
```

#### Basic Usage

```bash
# Scan entire project
gosec ./...

# Scan with detailed output
gosec -fmt=json -out=security-report.json ./...

# Focus on specific rules
gosec -include=G101,G102,G103 ./...

# Exclude certain rules
gosec -exclude=G304 ./...

# Track suppressions
gosec -track-suppressions -fmt=sarif -out=results.sarif ./...
```

#### Security Rules (Most Critical)

| Rule | Description | Severity |
|------|-------------|----------|
| **G101** | Hard-coded credentials | HIGH |
| **G102** | Bind to all interfaces | MEDIUM |
| **G103** | Unsafe block usage | HIGH |
| **G104** | Unchecked errors | MEDIUM |
| **G107** | URL as taint input | MEDIUM |
| **G201** | SQL injection (format) | HIGH |
| **G202** | SQL injection (concat) | HIGH |
| **G304** | File path traversal | HIGH |
| **G401** | MD5/SHA1 usage | MEDIUM |
| **G402** | Bad TLS settings | HIGH |
| **G404** | Weak RNG | MEDIUM |

#### Suppressing False Positives

```go
// Suppress specific rule with justification
func main() {
	// #nosec G304 -- File path is validated by Casbin policy
	data, err := os.ReadFile(userProvidedPath)

	// Multiple rules suppression
	// #nosec G201 G202 -- Query parameters are sanitized via SQLC
	db.Query("SELECT * FROM users WHERE id = " + sanitizedID)
}
```

#### Integration with CI/CD

```yaml
# .github/workflows/security.yml
name: Security Scan

on:
  push:
    branches: [develop, main]
  pull_request:
  schedule:
    - cron: '0 0 * * 0'  # Weekly

jobs:
  gosec:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: '-fmt=sarif -out=results.sarif -track-suppressions ./...'

      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: results.sarif
```

### Security Testing Checklist

- [ ] **Passwords**: Never hard-coded, always hashed (Argon2id)
- [ ] **Secrets**: Not in code, use environment variables
- [ ] **SQL**: Use parameterized queries (SQLC generates safe code)
- [ ] **File Paths**: Validate against directory traversal
- [ ] **Input Validation**: Sanitize all user inputs
- [ ] **TLS**: Strong cipher suites, no InsecureSkipVerify
- [ ] **Randomness**: Use crypto/rand, not math/rand
- [ ] **Error Messages**: Don't leak sensitive information

---

## Live System Testing

### HTTP API Testing (Against Running Container)

```go
// File: tests/live/api_health_test.go
package live_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

const BaseURL = "http://localhost:8096"

// TestHealthEndpoints validates all health checks
func TestHealthEndpoints(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping live API tests in short mode")
	}

	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{"Liveness", "/health/live", "healthy"},
		{"Readiness", "/health/ready", "healthy"},
		{"Startup", "/health/startup", "healthy"},
	}

	client := &http.Client{Timeout: 5 * time.Second}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(BaseURL + tt.endpoint)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected 200, got %d", resp.StatusCode)
			}

			var result map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			if result["status"] != tt.want {
				t.Errorf("Expected status=%s, got %s", tt.want, result["status"])
			}
		})
	}
}
```

### Authentication Flow Testing

```go
// File: tests/live/auth_flow_test.go
package live_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestCompleteAuthFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping live API tests")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	timestamp := time.Now().Unix()

	// 1. Register user
	registerData := map[string]string{
		"username": fmt.Sprintf("testuser_%d", timestamp),
		"email":    fmt.Sprintf("test%d@example.com", timestamp),
		"password": "SecureP@ssw0rd!",
	}

	registerJSON, _ := json.Marshal(registerData)
	resp, err := client.Post(
		BaseURL+"/api/v1/auth/register",
		"application/json",
		bytes.NewBuffer(registerJSON),
	)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 on register, got %d", resp.StatusCode)
	}

	var registerResp struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	json.NewDecoder(resp.Body).Decode(&registerResp)

	// 2. Login
	loginData := map[string]string{
		"username": registerData["username"],
		"password": registerData["password"],
	}

	loginJSON, _ := json.Marshal(loginData)
	resp, err = client.Post(
		BaseURL+"/api/v1/auth/login",
		"application/json",
		bytes.NewBuffer(loginJSON),
	)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	defer resp.Body.Close()

	var loginResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResp)

	if loginResp.AccessToken == "" {
		t.Error("Access token not returned")
	}

	// 3. Access protected endpoint
	req, _ := http.NewRequest("GET", BaseURL+"/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Protected endpoint failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 on /users/me, got %d", resp.StatusCode)
	}

	// 4. Refresh token
	refreshData := map[string]string{
		"refresh_token": loginResp.RefreshToken,
	}

	refreshJSON, _ := json.Marshal(refreshData)
	resp, err = client.Post(
		BaseURL+"/api/v1/auth/refresh",
		"application/json",
		bytes.NewBuffer(refreshJSON),
	)
	if err != nil {
		t.Fatalf("Token refresh failed: %v", err)
	}
	defer resp.Body.Close()

	var refreshResp struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&refreshResp)

	if refreshResp.AccessToken == loginResp.AccessToken {
		t.Error("Refresh token should return NEW access token")
	}
}
```

### Running Live Tests

```bash
# Only run live tests (requires running container)
go test -v -tags=live ./tests/live/...

# Skip live tests in CI
go test -short ./tests/live/...

# Run with race detection
go test -race -tags=live ./tests/live/...
```

---

## Integration with CI/CD

### Makefile Targets

```makefile
# File: Makefile

# Run all tests
.PHONY: test
test:
	go test -v -race -cover ./...

# Run benchmarks
.PHONY: bench
bench:
	go test -bench=. -benchmem -run=^$$ ./...

# Security scan
.PHONY: security
security:
	gosec -fmt=json -out=.workingdir/security-report.json ./...

# Live system tests (requires Docker)
.PHONY: test-live
test-live:
	docker-compose up -d
	sleep 5  # Wait for services
	go test -v -tags=live ./tests/live/...
	docker-compose down

# Full quality check
.PHONY: quality
quality: test bench security
	golangci-lint run ./...
	gofmt -s -w .
	@echo "✅ All quality checks passed!"

# Generate coverage report
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
```

### GitHub Actions Workflow

```yaml
# File: .github/workflows/quality.yml
name: Quality Assurance

on:
  push:
    branches: [develop, main]
  pull_request:

jobs:
  tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

  benchmarks:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Run benchmarks
        run: go test -bench=. -benchmem ./... > benchmarks.txt

      - name: Store benchmark results
        uses: actions/upload-artifact@v3
        with:
          name: benchmarks
          path: benchmarks.txt

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: '-fmt=sarif -out=results.sarif ./...'

      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: results.sarif

  live-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:18
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s

      dragonfly:
        image: docker.dragonflydb.io/dragonflydb/dragonfly:latest

    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Build Docker image
        run: docker build -t revenge:test .

      - name: Run container
        run: docker-compose up -d

      - name: Wait for services
        run: sleep 10

      - name: Run live tests
        run: go test -v -tags=live ./tests/live/...

      - name: Cleanup
        if: always()
        run: docker-compose down
```

---

## Workflow Steps

### 1. Initial Setup

```bash
# Install tools
make install-tools

# Or manually:
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/perf/cmd/benchstat@latest
```

### 2. Start Live System

```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up -d

# Verify containers are healthy
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Check logs
docker logs revenge-dev --tail=50 -f
```

### 3. Test Live Endpoints

```bash
# Health checks
curl http://localhost:8096/health/live
curl http://localhost:8096/health/ready

# Register user
curl -X POST http://localhost:8096/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"Test123!"}'

# Login
curl -X POST http://localhost:8096/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"Test123!"}'
```

### 4. Run Test Suite

```bash
# Unit tests
go test -v ./internal/...

# Integration tests
go test -v ./tests/integration/...

# Live tests
go test -v -tags=live ./tests/live/...

# All tests with coverage
go test -v -cover -coverprofile=coverage.out ./...
```

### 5. Run Benchmarks

```bash
# Baseline benchmarks
go test -bench=. -benchmem ./... > baseline.txt

# After optimization
go test -bench=. -benchmem ./... > optimized.txt

# Compare
benchstat baseline.txt optimized.txt
```

### 6. Security Scan

```bash
# Full security scan
gosec -fmt=json -out=.workingdir/security-report.json ./...

# Check for high-severity issues
gosec -severity=high ./...

# Specific checks
gosec -include=G101,G201,G401 ./...
```

### 7. Document Findings

```bash
# Create bug report
cat > .workingdir/bugs/bug-XX-description.md <<EOF
# Bug #XX: Description

## Discovery
- **Date**: $(date +%Y-%m-%d)
- **Found by**: Live system testing
- **Severity**: High/Medium/Low

## Reproduction Steps
1. Step one
2. Step two

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Root Cause
Why this happens

## Fix
How to fix it

## Tests Added
- [ ] Unit test
- [ ] Integration test
- [ ] Live test

## Related Files
- File 1
- File 2
EOF
```

### 8. Fix and Verify

```bash
# Make fixes
git checkout -b fix/bug-XX-description

# Run tests to verify fix
go test -v ./...

# Ensure no new issues
gosec ./...
golangci-lint run ./...

# Commit with reference
git commit -m "fix: Bug #XX - Description

- Fixed root cause
- Added regression test
- Verified with live system

Fixes #XX"
```

---

## Templates

### Benchmark Test Template

```go
// File: <package>/<file>_bench_test.go
package <package>_test

import (
	"testing"
	"<module>/<package>"
)

// Benchmark<Function> measures <what it measures>
func Benchmark<Function>(b *testing.B) {
	// Setup (not timed)
	setup := prepareTestData()

	// Reset timer before actual benchmark
	b.ResetTimer()

	// Run benchmark
	for b.Loop() {
		<package>.<Function>(setup)
	}
}

// Benchmark<Function>Parallel tests concurrent performance
func Benchmark<Function>Parallel(b *testing.B) {
	setup := prepareTestData()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			<package>.<Function>(setup)
		}
	})
}
```

### Live API Test Template

```go
// File: tests/live/<feature>_test.go
package live_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

const BaseURL = "http://localhost:8096"

func Test<Feature>(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping live API test")
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// 1. Setup

	// 2. Execute

	// 3. Verify

	// 4. Cleanup
}
```

### Bug Report Template

```markdown
# Bug #XX: [Short Description]

## Discovery
- **Date**: YYYY-MM-DD
- **Found by**: [Live testing / Integration test / Unit test]
- **Severity**: [Critical / High / Medium / Low]
- **Component**: [Service / Package affected]

## Reproduction Steps

1. Step 1
2. Step 2
3. Step 3

## Expected Behavior

What should happen when following the steps above.

## Actual Behavior

What actually happens (include error messages, logs).

## Root Cause Analysis

### Why it happens
Explanation of the underlying issue.

### Affected Code
```go
// File: path/to/file.go
// Lines: XX-YY
// Problematic code snippet
```

## Fix Implementation

### Changes Made
1. Change 1
2. Change 2

### Modified Files
- `path/to/file1.go`
- `path/to/file2.go`

## Testing

### Tests Added
- [ ] Unit test: `Test<Function>`
- [ ] Integration test: `Test<Feature>Integration`
- [ ] Live test: `Test<Feature>Live`
- [ ] Benchmark: `Benchmark<Function>`

### Verification
- [ ] All existing tests pass
- [ ] New tests pass
- [ ] No new security issues (gosec)
- [ ] No new linting issues
- [ ] Benchmarks show acceptable performance

## Documentation Updates
- [ ] Code comments
- [ ] API documentation
- [ ] OpenAPI spec
- [ ] README/CHANGELOG

## Related Issues
- Fixes #XX
- Related to #YY
```

---

## Next Steps

1. **Create benchmarks** for crypto service (`internal/crypto/password_bench_test.go`)
2. **Run security scan** with gosec
3. **Create live test suite** (`tests/live/`)
4. **Document any bugs found** in `.workingdir/bugs/`
5. **Set up CI/CD pipeline** with quality checks
6. **Establish baseline metrics** for performance tracking

---

## Success Criteria

- ✅ All tests passing (unit, integration, live)
- ✅ No high-severity security issues
- ✅ Benchmark baselines established
- ✅ Code coverage > 80%
- ✅ All bugs documented and tracked
- ✅ CI/CD pipeline running quality checks
- ✅ Live system validated functional

---

**End of Framework Documentation**
