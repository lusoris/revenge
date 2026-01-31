---
name: validate-tools
description: Validate all required development tools are installed and at correct versions
argument-hint: [local|remote]
disable-model-invocation: false
allowed-tools: Bash(*), Read(*)
---

# Validate Development Tools

Validates that all required development tools are installed and at the correct versions for Revenge development.

## Usage

```
/validate-tools              # Validate tools (auto-detect environment)
/validate-tools local        # Validate for local development
/validate-tools remote       # Validate for remote/Coder development
```

## Arguments

- `$0`: Environment (optional: local, remote) - Auto-detects if not provided

## Prerequisites

None - this skill checks prerequisites!

## Source of Truth

Required versions are defined in:
- [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)

## Task

### Step 1: Read Required Versions from SOURCE_OF_TRUTH

First, read the required versions from the authoritative source:

```bash
# Read SOURCE_OF_TRUTH to get exact required versions
cat docs/dev/design/00_SOURCE_OF_TRUTH.md
```

You should extract the following required versions:
- Go version
- Node.js version
- Python version
- PostgreSQL version
- Svelte version
- SvelteKit version
- Any other critical tool versions listed

**IMPORTANT**: All version requirements must come from SOURCE_OF_TRUTH, not hardcoded values.

### Step 2: Check Installed Tool Versions

Check what's currently installed:

**Go**:
```bash
go version
go env GOEXPERIMENT
```

**Git**:
```bash
git --version
```

**Node.js**:
```bash
node --version
npm --version
```

**Python**:
```bash
python3 --version
```

**Go Tools**:
```bash
gopls version
golangci-lint --version
air -v 2>&1 || echo "not installed"
```

**Frontend dependencies** (in web/ directory):
```bash
cd web && npm list svelte --depth=0 2>/dev/null || echo "not installed"
cd web && npm list @sveltejs/kit --depth=0 2>/dev/null || echo "not installed"
cd ..
```

**Python tools**:
```bash
ruff --version 2>&1 || echo "not installed"
pytest --version 2>&1 || echo "not installed"
```

**Database & Services (Local Only)**:
```bash
docker --version 2>&1 || echo "not installed"
docker-compose --version 2>&1 || echo "not installed"
```

### Remote Development (Remote/Coder Only)

**Coder CLI**:
```bash
coder version
# Required: latest
```

**SSH**:
```bash
ssh -V
# Required: Any recent version
```

### IDE-Specific (Optional)

**VS Code**:
```bash
code --version
# Optional but recommended
```

**Zed**:
```bash
zed --version
# Optional
```

### Environment Variables

Check important environment variables:
```bash
echo "GOEXPERIMENT: $GOEXPERIMENT"
# Should include: greenteagc,jsonv2

echo "PATH: $PATH"
# Should include: $HOME/go/bin, /usr/local/go/bin
```

## Validation Logic

After checking all tools:

1. **Read SOURCE_OF_TRUTH** to get exact required versions:
   ```bash
   # Parse docs/dev/design/00_SOURCE_OF_TRUTH.md
   # Extract version requirements from the table
   ```

2. **Compare installed vs required versions**
   - Use the versions from SOURCE_OF_TRUTH (NOT hardcoded values)
   - Parse semantic versions (e.g., "1.25.6" vs "1.24.0")
   - Handle version ranges (e.g., "20+" means 20.0.0 or higher)

3. **Report results in table format**:
   ```
   ## Validation Results

   | Tool | Required* | Installed | Status |
   |------|-----------|-----------|--------|
   | Go | <from SOT> | <detected> | ✅/⚠️/❌ |
   | Node.js | <from SOT> | <detected> | ✅/⚠️/❌ |
   | Python | <from SOT> | <detected> | ✅/⚠️/❌ |
   | ...more tools... |

   *Required versions from: docs/dev/design/00_SOURCE_OF_TRUTH.md
   ```

4. **Provide installation instructions** for missing/wrong versions:
   - Link to setup documentation with correct versions
   - Platform-specific install commands (fetch version from SOURCE_OF_TRUTH)
   - References:
     - **SOURCE_OF_TRUTH**: docs/dev/design/00_SOURCE_OF_TRUTH.md
     - Local: .shared/docs/ONBOARDING.md
     - Tools: .zed/docs/SETUP.md, .jetbrains/docs/SETUP.md

## Success Criteria

- All required tools installed
- All versions meet minimum requirements
- GOEXPERIMENT configured correctly
- PATH includes Go binaries

## Example Output

```
# Validating Development Tools

Reading required versions from: docs/dev/design/00_SOURCE_OF_TRUTH.md

## Environment: Local Development

## Core Tools
✅ Go <installed-version> (required: <from SOT>)
✅ Git <installed-version> (required: 2.0+)
✅ GOEXPERIMENT=greenteagc,jsonv2

## Backend Tools
✅ gopls <installed-version>
✅ golangci-lint <installed-version>
⚠️ air not found (optional, recommended for hot reload)

## Frontend Tools
✅ Node.js <installed-version> (required: <from SOT>)
✅ npm <installed-version>
✅ Svelte <installed-version> (required: <from SOT>)
✅ SvelteKit <installed-version> (required: <from SOT>)

## Python Tools
✅ Python <installed-version> (required: <from SOT>)
✅ ruff <installed-version>
❌ pytest not found (optional)

## Services
✅ Docker <installed-version>
✅ Docker Compose <installed-version>

## Summary
✅ 14 tools validated
⚠️ 1 optional tool missing (air)
❌ 1 optional tool missing (pytest)

All required tools are installed and at correct versions!
(Versions verified against: docs/dev/design/00_SOURCE_OF_TRUTH.md)

### Optional Improvements
Install air for hot reload:
  go install github.com/cosmtrek/air@latest

Install pytest for Python testing:
  pip install pytest
```
