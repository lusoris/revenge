---
name: run-all-tests
description: Run all test suites (Go, frontend, Python) with comprehensive reporting
argument-hint: [--coverage|--race|--quick|--watch]
disable-model-invocation: false
allowed-tools: Bash(*), Read(*), Write(*)
---

# Run All Tests

Runs comprehensive test suites across all components of the Revenge project: Go backend, frontend (Svelte), and Python scripts.

## Usage

```
/run-all-tests                  # Run all tests with standard settings
/run-all-tests --coverage       # Run with coverage reports
/run-all-tests --race           # Run with race detector (Go only)
/run-all-tests --quick          # Run quick tests only (skip slow integration tests)
/run-all-tests --watch          # Run in watch mode (frontend only)
```

## Arguments

- `$0`: Mode flag (optional: --coverage, --race, --quick, --watch)

## Prerequisites

- All dependencies installed (`go mod download`, `npm install` in web/)
- Services running if integration tests included (PostgreSQL, Dragonfly, etc.)

## Task

Run comprehensive test suites and provide detailed results with clear pass/fail status.

### Step 1: Pre-Test Validation

Verify test environment is ready:

```bash
# Check Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go not found. Install Go (check SOURCE_OF_TRUTH for version)"
    exit 1
fi

# Check if in project root
if [ ! -f "go.mod" ]; then
    echo "❌ Not in project root (go.mod not found)"
    exit 1
fi

# Check services are running (for integration tests)
if ! docker ps | grep -q postgres; then
    echo "⚠️ PostgreSQL not running (integration tests may fail)"
    echo "Start services: docker-compose -f docker-compose.dev.yml up -d"
fi
```

### Step 2: Go Backend Tests

**Standard tests**:
```bash
echo "🧪 Running Go tests..."

# Run all tests with verbose output
go test ./... -v -timeout 5m

# Capture exit code
GO_EXIT=$?
```

**With coverage** (if --coverage flag):
```bash
echo "🧪 Running Go tests with coverage..."

# Run tests with coverage
go test ./... -coverprofile=coverage.out -covermode=atomic -timeout 5m

# Show coverage summary
go tool cover -func=coverage.out | grep total:

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
echo "📊 Coverage report: coverage.html"

# Capture exit code
GO_EXIT=$?
```

**With race detector** (if --race flag):
```bash
echo "🧪 Running Go tests with race detector..."

# Run with race detector (slower but catches race conditions)
go test ./... -race -timeout 10m

# Capture exit code
GO_EXIT=$?
```

**Quick tests** (if --quick flag):
```bash
echo "🧪 Running quick Go tests (skipping integration)..."

# Run only unit tests (skip integration tests with build tag)
go test ./... -short -timeout 2m

# Capture exit code
GO_EXIT=$?
```

### Step 3: Frontend Tests

**Check if web/ directory exists**:
```bash
if [ ! -d "web" ]; then
    echo "⚠️ web/ directory not found, skipping frontend tests"
    FRONTEND_EXIT=0
else
    cd web

    # Check if node_modules exists
    if [ ! -d "node_modules" ]; then
        echo "📦 Installing frontend dependencies..."
        npm install
    fi

    # Run tests based on mode
    if [ "$MODE" == "--watch" ]; then
        echo "🧪 Running frontend tests in watch mode..."
        npm run test:watch
        FRONTEND_EXIT=$?
    elif [ "$MODE" == "--coverage" ]; then
        echo "🧪 Running frontend tests with coverage..."
        npm run test:coverage
        FRONTEND_EXIT=$?
    else
        echo "🧪 Running frontend tests..."
        npm run test
        FRONTEND_EXIT=$?
    fi

    cd ..
fi
```

### Step 4: Python Script Tests (if applicable)

```bash
if [ -d "scripts" ] && [ -f "scripts/requirements.txt" ]; then
    echo "🧪 Running Python tests..."

    # Check if pytest is available
    if command -v pytest &> /dev/null; then
        # Run pytest on scripts
        pytest scripts/ -v
        PYTHON_EXIT=$?
    else
        echo "⚠️ pytest not found, skipping Python tests"
        echo "Install: pip install pytest"
        PYTHON_EXIT=0
    fi
else
    echo "ℹ️ No Python tests found"
    PYTHON_EXIT=0
fi
```

### Step 5: Linting (Optional Quality Check)

Run linters to catch code quality issues:

**Go linting**:
```bash
if command -v golangci-lint &> /dev/null; then
    echo "🔍 Running Go linter..."
    golangci-lint run ./...
    LINT_EXIT=$?
else
    echo "⚠️ golangci-lint not found (optional)"
    LINT_EXIT=0
fi
```

**Frontend linting**:
```bash
if [ -d "web" ]; then
    echo "🔍 Running frontend linter..."
    cd web
    npm run lint
    FRONTEND_LINT_EXIT=$?
    cd ..
else
    FRONTEND_LINT_EXIT=0
fi
```

**Python linting**:
```bash
if command -v ruff &> /dev/null && [ -d "scripts" ]; then
    echo "🔍 Running Python linter..."
    ruff check scripts/
    PYTHON_LINT_EXIT=$?
else
    PYTHON_LINT_EXIT=0
fi
```

### Step 6: Test Results Summary

Generate comprehensive test report:

```bash
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "           TEST RESULTS SUMMARY           "
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Go tests
if [ $GO_EXIT -eq 0 ]; then
    echo "✅ Go Tests: PASSED"
else
    echo "❌ Go Tests: FAILED (exit code: $GO_EXIT)"
fi

# Frontend tests
if [ $FRONTEND_EXIT -eq 0 ]; then
    echo "✅ Frontend Tests: PASSED"
else
    echo "❌ Frontend Tests: FAILED (exit code: $FRONTEND_EXIT)"
fi

# Python tests
if [ $PYTHON_EXIT -eq 0 ]; then
    echo "✅ Python Tests: PASSED"
else
    echo "❌ Python Tests: FAILED (exit code: $PYTHON_EXIT)"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Linting results (if run)
if [ -n "$LINT_EXIT" ]; then
    echo ""
    echo "Linting Results:"
    [ $LINT_EXIT -eq 0 ] && echo "✅ Go Lint: PASSED" || echo "⚠️ Go Lint: WARNINGS"
    [ $FRONTEND_LINT_EXIT -eq 0 ] && echo "✅ Frontend Lint: PASSED" || echo "⚠️ Frontend Lint: WARNINGS"
    [ $PYTHON_LINT_EXIT -eq 0 ] && echo "✅ Python Lint: PASSED" || echo "⚠️ Python Lint: WARNINGS"
fi

# Overall result
TOTAL_EXIT=$((GO_EXIT + FRONTEND_EXIT + PYTHON_EXIT))

echo ""
if [ $TOTAL_EXIT -eq 0 ]; then
    echo "🎉 ALL TESTS PASSED!"
    echo ""
    echo "✨ Great work! Your code is ready for review."
    exit 0
else
    echo "💥 SOME TESTS FAILED"
    echo ""
    echo "Please fix failing tests before committing."
    echo ""
    echo "For help, see:"
    echo "  - .shared/docs/WORKFLOWS.md#testing-workflow"
    echo "  - .shared/docs/TROUBLESHOOTING.md"
    exit 1
fi
```

### Step 7: Coverage Report (if --coverage)

If coverage flag was used, display coverage summary:

```bash
if [ "$MODE" == "--coverage" ]; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "          COVERAGE SUMMARY               "
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    # Go coverage
    echo "Go Backend Coverage:"
    go tool cover -func=coverage.out | grep total:
    echo "HTML Report: coverage.html"

    # Frontend coverage (if generated)
    if [ -f "web/coverage/coverage-summary.json" ]; then
        echo ""
        echo "Frontend Coverage:"
        cat web/coverage/coverage-summary.json | jq '.total'
        echo "HTML Report: web/coverage/index.html"
    fi

    echo ""
    echo "Target: 80%+ coverage for critical paths"
    echo "See: .shared/docs/WORKFLOWS.md#testing"
fi
```

## Test Output Examples

### Successful Run

```
🧪 Running Go tests...
✅ All Go tests passed (127 tests)

🧪 Running frontend tests...
✅ All frontend tests passed (42 tests)

🧪 Running Python tests...
✅ All Python tests passed (8 tests)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           TEST RESULTS SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Go Tests: PASSED
✅ Frontend Tests: PASSED
✅ Python Tests: PASSED

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🎉 ALL TESTS PASSED!

✨ Great work! Your code is ready for review.
```

### With Coverage

```
🧪 Running Go tests with coverage...
✅ All Go tests passed
📊 Coverage: 84.7% of statements

🧪 Running frontend tests with coverage...
✅ All frontend tests passed
📊 Coverage: 78.3% of statements

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           TEST RESULTS SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Go Tests: PASSED
✅ Frontend Tests: PASSED

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
          COVERAGE SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Go Backend Coverage:
total:  (statements)    84.7%
HTML Report: coverage.html

Frontend Coverage:
Lines: 78.3%
Statements: 76.8%
Functions: 81.2%
Branches: 74.5%
HTML Report: web/coverage/index.html

Target: 80%+ coverage for critical paths
```

### With Failures

```
🧪 Running Go tests...
❌ FAIL: TestMovieHandler (0.03s)
❌ FAIL: TestSearchAPI (0.02s)
2 tests failed

🧪 Running frontend tests...
✅ All frontend tests passed (42 tests)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
           TEST RESULTS SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

❌ Go Tests: FAILED (exit code: 1)
✅ Frontend Tests: PASSED

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💥 SOME TESTS FAILED

Please fix failing tests before committing.

For help, see:
  - .shared/docs/WORKFLOWS.md#testing-workflow
  - .shared/docs/TROUBLESHOOTING.md
```

## Performance Considerations

**Quick mode** (--quick):
- Skips integration tests
- Runs in ~30 seconds vs 2-3 minutes
- Use for rapid feedback during development

**Race detector** (--race):
- Slower (2-3x) but catches concurrency bugs
- Recommended before commits
- Required before merges to main

**Coverage** (--coverage):
- Adds ~20% overhead
- Generates detailed HTML reports
- Required for PR reviews

## Integration with CI/CD

This skill mimics what CI/CD runs:
- GitHub Actions runs: `/run-all-tests --coverage --race`
- Pre-commit hooks run: `/run-all-tests --quick`

See:
- .github/workflows/ - GitHub Actions config
- .githooks/ - Git hooks

## Troubleshooting

**Tests fail on fresh clone**:
```bash
# Install dependencies first
go mod download
cd web && npm install

# Start services
docker-compose -f docker-compose.dev.yml up -d

# Wait for services
sleep 5

# Try again
/run-all-tests
```

**Race detector fails but tests pass**:
- Race conditions detected (not necessarily bugs, but worth investigating)
- Review race detector output carefully
- See: https://go.dev/blog/race-detector

**Coverage too low**:
- Aim for 80%+ coverage on critical paths
- 100% coverage not required (diminishing returns)
- Focus on business logic, not boilerplate

**Frontend tests hang**:
- Check if dev server is running (should not be during tests)
- Kill any npm processes: `pkill -f "npm\|node"`

For more help:
- .shared/docs/WORKFLOWS.md#testing-workflow
- .shared/docs/TROUBLESHOOTING.md
