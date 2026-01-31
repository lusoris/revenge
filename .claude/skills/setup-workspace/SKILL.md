---
name: setup-workspace
description: Set up complete development environment (local or Coder workspace)
argument-hint: "[local|remote] [--full|--minimal]"
disable-model-invocation: false
allowed-tools: Bash(*), Read(*), Write(*), Edit(*)
---

# Setup Development Workspace

Sets up a complete development environment for Revenge, either locally or in a Coder workspace.

## Usage

```
/setup-workspace                    # Auto-detect environment, full setup
/setup-workspace local              # Local development setup
/setup-workspace remote             # Coder workspace setup
/setup-workspace local --minimal    # Minimal setup (skip optional tools)
/setup-workspace local --full       # Full setup with all tools
```

## Arguments

- `$0`: Environment (optional: local, remote) - Auto-detects if not provided
- `$1`: Mode (optional: --full, --minimal) - Default: --full

## Prerequisites

Basic system tools (git, curl, bash) should be installed.

## Task

Set up a complete development environment based on the detected or specified environment.

### Step 1: Detect Environment

```bash
# Check if running in Coder workspace
if [ -n "$CODER_WORKSPACE_NAME" ]; then
    echo "Environment: Coder workspace"
    ENV="remote"
elif coder list 2>/dev/null | grep -q "revenge"; then
    echo "Environment: Coder CLI available"
    ENV="remote"
else
    echo "Environment: Local machine"
    ENV="local"
fi
```

### Step 2: Validate Prerequisites

Run validation first:
```bash
# Use validate-tools skill
/validate-tools $ENV
```

### Step 3: Get Required Versions from SOURCE_OF_TRUTH

**CRITICAL**: Read required versions from the authoritative source:

```bash
# Parse docs/dev/design/00_SOURCE_OF_TRUTH.md to get:
# - Go version
# - Node.js version
# - Python version
# - PostgreSQL version
# - Other tool versions

# Store these in variables for use in download URLs
```

### Step 4: Install Missing Tools

Based on validation results and versions from SOURCE_OF_TRUTH:

**Go (if missing)**:
```bash
# Get version from SOURCE_OF_TRUTH (e.g., "1.25.6")
GO_VERSION="<from SOURCE_OF_TRUTH>"

# Linux
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

# macOS
brew install go
# Then verify version matches SOURCE_OF_TRUTH

# Verify
go version
```

**Node.js (if missing)**:
```bash
# Get major version from SOURCE_OF_TRUTH (e.g., "20")
NODE_VERSION="<from SOURCE_OF_TRUTH>"

# Linux (using nvm)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install ${NODE_VERSION}
nvm use ${NODE_VERSION}

# macOS
brew install node@${NODE_VERSION}

# Verify
node --version
```

**Python (if missing)**:
```bash
# Get version from SOURCE_OF_TRUTH (e.g., "3.12")
PYTHON_VERSION="<from SOURCE_OF_TRUTH>"

# Linux
sudo apt install python${PYTHON_VERSION} python3-pip

# macOS
brew install python@${PYTHON_VERSION}

# Verify
python3 --version
```

**Docker (Local only, if missing)**:
```bash
# macOS
brew install --cask docker

# Linux
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Verify
docker --version
```

### Step 4: Install Development Tools

**Go tools**:
```bash
# gopls (LSP)
go install golang.org/x/tools/gopls@latest

# golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# air (hot reload)
go install github.com/cosmtrek/air@latest

# delve (debugger)
go install github.com/go-delve/delve/cmd/dlv@latest

# Verify
gopls version
golangci-lint --version
air -v
```

**Python tools**:
```bash
# ruff (linter/formatter)
pip install ruff

# pytest (testing)
pip install pytest

# Verify
ruff --version
pytest --version
```

### Step 5: Clone Repository (if not already cloned)

```bash
# Check if already in repo
if [ ! -d ".git" ]; then
    echo "Cloning repository..."
    git clone https://github.com/lusoris/revenge.git
    cd revenge
else
    echo "Already in revenge repository"
fi
```

### Step 6: Install Dependencies

**Go modules**:
```bash
go mod download
go mod verify
```

**Frontend dependencies**:
```bash
cd web
npm install
cd ..
```

**Python dependencies** (if requirements.txt exists):
```bash
if [ -f "scripts/requirements.txt" ]; then
    pip install -r scripts/requirements.txt
fi
```

### Step 7: Set Up Services (Local only)

```bash
# Start PostgreSQL, Dragonfly, Typesense via Docker Compose
docker-compose -f docker-compose.dev.yml up -d

# Wait for services to be ready
sleep 5

# Verify services
docker-compose -f docker-compose.dev.yml ps
```

### Step 8: Run Database Migrations (if applicable)

```bash
# Check if migrations exist
if [ -d "migrations" ]; then
    # Install golang-migrate if needed
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

    # Run migrations
    migrate -path migrations -database "postgres://revenge:revenge@localhost:5432/revenge?sslmode=disable" up
fi
```

### Step 9: Configure Environment

**Set GOEXPERIMENT**:
```bash
# Add to shell profile if not already present
if ! grep -q "GOEXPERIMENT=greenteagc,jsonv2" ~/.bashrc; then
    echo 'export GOEXPERIMENT=greenteagc,jsonv2' >> ~/.bashrc
fi

# Apply for current session
export GOEXPERIMENT=greenteagc,jsonv2
```

**Add Go binaries to PATH**:
```bash
if ! grep -q 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' ~/.bashrc; then
    echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
fi
```

### Step 10: Verify Setup

**Run validation again**:
```bash
/validate-tools $ENV
```

**Test build**:
```bash
# Build backend
go build -o bin/revenge ./cmd/revenge

# Check binary
./bin/revenge --version
```

**Test frontend build**:
```bash
cd web
npm run build
cd ..
```

**Run quick test**:
```bash
# Run a simple test
go test ./internal/config -v
```

### Step 11: IDE-Specific Setup (Optional)

**VS Code**:
```bash
# Install recommended extensions
code --install-extension golang.go
code --install-extension svelte.svelte-vscode
code --install-extension esbenp.prettier-vscode
```

**Zed** (create config if doesn't exist):
```bash
# Config already in .zed/settings.json
echo "Zed configuration: .zed/settings.json"
```

### Step 12: Print Setup Summary

After setup completes, print a summary:

```
# Setup Complete!

## Environment: Local Development

## Installed Tools (verified against SOURCE_OF_TRUTH)
✅ Go <version>
✅ Node.js <version>
✅ Python <version>
✅ Docker <version>
✅ gopls, golangci-lint, air, dlv
✅ ruff, pytest
✅ npm dependencies (web/)

## Services Running
✅ PostgreSQL (localhost:5432)
✅ Dragonfly (localhost:6379)
✅ Typesense (localhost:8108)

## Next Steps

1. Start development server:
   air

2. Start frontend:
   cd web && npm run dev

3. Run tests:
   go test ./...

4. Open in IDE:
   - VS Code: code .
   - Zed: zed .
   - JetBrains: Open project in GoLand/IntelliJ IDEA

## Documentation
- Getting Started: .shared/docs/ONBOARDING.md
- Development Guide: docs/dev/design/operations/DEVELOPMENT.md
- Workflows: .shared/docs/WORKFLOWS.md

Happy coding! 🚀
```

## Remote/Coder Specific Steps

If environment is remote (Coder workspace):

### Skip Docker Setup
Services are already running in workspace container.

### Configure SSH
```bash
# On local machine
coder config-ssh

# Verify connection
ssh coder-revenge-dev
```

### IDE Connection
Provide connection instructions for each IDE:

**VS Code**:
```bash
# Browser
coder open revenge-dev

# Desktop
code --remote ssh-remote+coder-revenge-dev /workspace/revenge
```

**Zed**:
```bash
# Connect via SSH
# Zed → File → Open Remote → SSH → coder-revenge-dev
```

**JetBrains Gateway**:
- Open Gateway
- Select Coder provider
- Connect to revenge-dev workspace

## Minimal Setup Mode

If `--minimal` flag is provided:
- Skip optional tools (air, pytest, delve)
- Skip IDE-specific setup
- Skip documentation generation
- Only install required core tools

## Error Handling

If any step fails:
1. Print clear error message
2. Show command that failed
3. Provide troubleshooting link:
   - Local: .shared/docs/TROUBLESHOOTING.md
   - Remote: .coder/docs/TROUBLESHOOTING.md
4. Allow user to continue or abort

## Success Criteria

- All required tools installed and validated
- Repository cloned and dependencies installed
- Services running (local) or accessible (remote)
- Environment variables configured
- Can build backend and frontend successfully
- Tests run successfully
