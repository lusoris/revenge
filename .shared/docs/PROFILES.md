# Pre-Configured Development Profiles

**Purpose**: Quick-start configurations for different developer roles and environments

**Last Updated**: 2026-01-31

---

## Overview

Development profiles are pre-configured settings bundles that optimize your development environment for your specific role. Choose one, apply it, and you're ready to go.

### Quick Selection

| Your Role | Profile | Setup Time | IDE Choice |
|-----------|---------|-----------|-----------|
| Go backend work | **Backend Developer** | 5 min | VS Code or Zed |
| Svelte/TypeScript frontend | **Frontend Developer** | 5 min | VS Code (better TS support) |
| Multiple specialties | **Full-Stack Developer** | 10 min | VS Code (complete setup) |
| Server/K8s/Docker operations | **DevOps Engineer** | 10 min | VS Code (terminal-heavy) |
| Docs/Content creation | **Documentation Writer** | 3 min | Zed (lightweight) |
| Using Coder for remote dev | **Remote Development** | 15 min | Zed (SSH optimized) |

---

## Profile 1: Backend Developer

**Focus**: Go development with PostgreSQL, testing, and debugging

**Best For**: Building APIs, services, and backend features

### Setup Steps

#### 1. Copy Profile Files

```bash
# VS Code
cp profiles/backend-developer/.vscode/* .vscode/

# Zed
cp profiles/backend-developer/.zed/* .zed/
```

#### 2. Verify Configuration

```bash
# Check Go is installed and latest
go version  # Should be 1.25+

# Check go-related tools
gopls version
golangci-lint version
goimports -h
```

#### 3. Install VS Code Extensions

From the project root:
```bash
# Use VS Code: Cmd/Ctrl+Shift+P → "Extensions: Show Recommended Extensions"
# Select recommended extensions (50+ total)
# Focus on: Go, Test Explorer, Database, etc.
```

#### 4. Configure Debugger

Create `.vscode/launch.json` with Go debugging:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Current Test",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${fileDirname}",
      "args": ["-test.v"],
      "showLog": true
    },
    {
      "name": "Attach to Process",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": "${command:pickGoProcess}"
    }
  ]
}
```

#### 5. Test the Setup

```bash
# Build
make build

# Run tests
go test -v ./...

# Debug a test: Use VS Code's Debug menu → "Debug Current Test"
```

### Profile Configuration

#### VS Code Settings (`.vscode/settings.json`)

```json
{
  // Go-focused settings
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.formatTool": "goimports",
  "go.testFlags": ["-v", "-race"],
  "go.coverOnTestPackage": true,

  "gopls": {
    "ui.semanticTokens": true,
    "ui.completion.usePlaceholders": true,
    "analyses": {
      "unusedparams": true,
      "shadow": true,
      "nilfunc": true,
      "unreachable": true
    },
    "staticcheck": true
  },

  // Editor defaults (for Go files)
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": "explicit"
  },
  "editor.tabSize": 4,
  "editor.insertSpaces": false,
  "editor.rulers": [100, 120],

  // Testing integration
  "go.enableCodeLens": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  },

  // Database (optional)
  "database.connections": [
    {
      "name": "revenge-dev",
      "driver": "PostgreSQL",
      "host": "localhost",
      "port": 5432,
      "user": "postgres",
      "password": "postgres",
      "database": "revenge_dev"
    }
  ]
}
```

#### Zed Settings (`.zed/settings.json`)

```json
{
  "tab_size": 4,
  "hard_tabs": true,
  "format_on_save": "on",
  "ensure_final_newline_on_save": true,
  "remove_trailing_whitespace_on_save": true,

  "languages": {
    "Go": {
      "tab_size": 4,
      "hard_tabs": true,
      "format_on_save": "on"
    }
  },

  "lsp": {
    "gopls": {
      "binary": "gopls",
      "initialization_options": {
        "hints": {
          "assignVariableTypes": true,
          "parameterNames": true,
          "rangeVariableTypes": true
        },
        "analyses": {
          "unusedparams": true,
          "shadow": true,
          "staticcheck": true
        }
      }
    }
  },

  "file_scan_exclusions": [
    "**/vendor",
    "**/bin",
    "**/.archive"
  ]
}
```

### Common Tasks

#### Run and Debug Tests

```bash
# Terminal
go test -v ./...

# VS Code: Click "Debug" above test function
# Zed: Use terminal

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Format Code

```bash
# VS Code: Cmd/Ctrl+Shift+P → "Format Document"
# Or use keyboard shortcut: Shift+Alt+F

# Terminal
go fmt ./...
goimports -w .
```

#### Run Linter

```bash
# Full lint
golangci-lint run ./...

# Fast mode (configured in settings)
golangci-lint run --fast ./...

# Specific check
golangci-lint run --enable=gosec ./...
```

#### Database Migrations

```bash
# Using golang-migrate
migrate -path db/migrations -database "$DB_URL" up

# Create new migration
migrate create -ext sql -dir db/migrations -seq add_users_table

# Rollback
migrate -path db/migrations -database "$DB_URL" down 1
```

### Troubleshooting

**gopls not responding**:
```bash
gopls -inspect=binding
pkill gopls
go get -u github.com/golang/tools/cmd/gopls
```

**Import organization issues**:
```bash
# Force goimports
goimports -w .

# Check VS Code has goimports configured
grep -A5 "formatTool" .vscode/settings.json
```

---

## Profile 2: Frontend Developer

**Focus**: Svelte 5, TypeScript, SvelteKit development

**Best For**: Building UI components, pages, and interactions

### Setup Steps

#### 1. Copy Profile Files

```bash
# VS Code (recommended for frontend)
cp profiles/frontend-developer/.vscode/* .vscode/

# Or Zed
cp profiles/frontend-developer/.zed/* .zed/
```

#### 2. Verify Dependencies

```bash
# Node.js and npm
node --version  # Should be 20+
npm --version

# Frontend tools
npm list svelte
npm list sveltekit
npm list tailwindcss
npm list prettier
```

#### 3. Install Extensions

**VS Code** (critical for frontend):

```json
{
  "recommendations": [
    "svelte.svelte-vscode",
    "esbenp.prettier-vscode",
    "bradlc.vscode-tailwindcss",
    "dbaeumer.vscode-eslint",
    "Vue.volar"
  ]
}
```

Command to install:
```bash
code --install-extension svelte.svelte-vscode
code --install-extension esbenp.prettier-vscode
code --install-extension bradlc.vscode-tailwindcss
```

#### 4. Configure Formatter

Ensure `.prettierrc.json` exists:

```json
{
  "useTabs": false,
  "tabWidth": 2,
  "semi": true,
  "singleQuote": false,
  "trailingComma": "es5",
  "plugins": ["prettier-plugin-svelte"]
}
```

#### 5. Test Setup

```bash
# Start dev server
npm run dev

# Format check
npm run format:check

# Lint check
npm run lint

# Build
npm run build
```

### Profile Configuration

#### VS Code Settings

```json
{
  // Svelte/TypeScript focused
  "[svelte]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 2,
    "editor.insertSpaces": true,
    "editor.defaultFormatter": "svelte.svelte-vscode"
  },

  "[typescript]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 2,
    "editor.insertSpaces": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode",
    "editor.codeActionsOnSave": {
      "source.fixAll.eslint": "explicit"
    }
  },

  "[javascript]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },

  // TypeScript settings
  "typescript.enablePromptUseWorkspaceTsdk": true,
  "typescript.tsdk": "node_modules/typescript/lib",
  "typescript.preferences.useImportType": true,

  // Tailwind
  "tailwindCSS.validate": true,
  "tailwindCSS.lint.cssConflict": "warning",

  // Editor
  "editor.formatOnSave": true,
  "editor.rulers": [100, 120],
  "editor.wordWrap": "on",
  "[html]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[css]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

#### Zed Settings

```json
{
  "tab_size": 2,
  "format_on_save": "on",

  "languages": {
    "TypeScript": {
      "tab_size": 2,
      "formatter": {
        "external": {
          "command": "prettier",
          "arguments": ["--stdin-filepath", "{buffer_path}"]
        }
      }
    },
    "Svelte": {
      "tab_size": 2,
      "formatter": {
        "external": {
          "command": "prettier",
          "arguments": [
            "--stdin-filepath",
            "{buffer_path}",
            "--plugin",
            "prettier-plugin-svelte"
          ]
        }
      }
    },
    "CSS": {
      "tab_size": 2,
      "formatter": {
        "external": {
          "command": "prettier",
          "arguments": ["--stdin-filepath", "{buffer_path}"]
        }
      }
    }
  }
}
```

### Common Tasks

#### Start Development Server

```bash
npm run dev
# Open http://localhost:5173
```

#### Format Code

```bash
# VS Code: Cmd/Ctrl+Shift+P → "Format Document"

# Terminal
npm run format

# Check without fixing
npm run format:check
```

#### Run Linter

```bash
npm run lint

# Fix issues
npm run lint -- --fix
```

#### Build for Production

```bash
npm run build

# Preview build
npm run preview
```

#### Create New Component

```bash
# Create component file in src/lib/components/
# Use existing components as templates
# Format on save will handle formatting

touch src/lib/components/MyComponent.svelte
```

### Troubleshooting

**Svelte intellisense not working**:
```bash
# Restart Svelte plugin
# Cmd/Ctrl+Shift+P → "Reload Window"

# Verify extension
code --list-extensions | grep svelte
```

**Prettier not formatting**:
```bash
# Check prettier is installed
npm list prettier prettier-plugin-svelte

# Check .prettierrc.json exists
cat .prettierrc.json

# Verify VS Code has correct formatter
grep -A2 "svelte" .vscode/settings.json
```

**TypeScript errors in Svelte**:
```bash
# Clear TypeScript cache
rm -rf node_modules/.vite

# Restart VS Code
```

---

## Profile 3: Full-Stack Developer

**Focus**: Both backend and frontend, all tools enabled

**Best For**: Working on complete features from API to UI

### Setup Steps

#### 1. Install Both Backend and Frontend Extensions

```bash
# Use BOTH Backend Developer and Frontend Developer profiles
# Copy .vscode settings from both (merge overlapping keys)
# Copy .zed settings from both
```

#### 2. Install All Tools

```bash
# Go tools
go install github.com/golang/tools/cmd/gopls@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Node tools
npm install -D svelte sveltekit tailwindcss prettier typescript

# Python tools (for scripts)
pip install ruff
```

#### 3. Configure Both IDEs

Use VS Code for primary development (better support for both Go and frontend):

```json
{
  // Merge of Backend Developer + Frontend Developer settings
  // Go settings (backend)
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.formatTool": "goimports",

  // Frontend settings
  "[svelte]": {
    "editor.defaultFormatter": "svelte.svelte-vscode"
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },

  // General settings
  "editor.formatOnSave": true,
  "editor.rulers": [100, 120],
  "files.exclude": {
    "**/node_modules": true,
    "**/vendor": true,
    "**/bin": true
  }
}
```

#### 4. Set Up Multi-Root Workspace

Create `.vscode/revenge.code-workspace`:

```json
{
  "folders": [
    {
      "path": ".",
      "name": "Revenge (Full-Stack)"
    },
    {
      "path": "frontend",
      "name": "Frontend (SvelteKit)"
    }
  ],
  "settings": {
    "editor.formatOnSave": true,
    "editor.rulers": [100, 120]
  },
  "extensions": {
    "recommendations": [
      "golang.go",
      "svelte.svelte-vscode",
      "esbenp.prettier-vscode",
      "bradlc.vscode-tailwindcss",
      "ms-python.python"
    ]
  }
}
```

Open with:
```bash
code revenge.code-workspace
```

### Profile Configuration

Combined settings from Backend and Frontend profiles. See sections above for detailed config.

Key difference: Enable all file type handlers:

```json
{
  "files.exclude": {
    "**/node_modules": true,
    "**/vendor": true,
    "**/.archive": true
  },

  "[go]": { /* Backend settings */ },
  "[svelte]": { /* Frontend settings */ },
  "[typescript]": { /* Frontend settings */ },
  "[python]": { /* Script settings */ }
}
```

### Common Full-Stack Workflow

```bash
# 1. Start backend in one terminal
go run cmd/revenge/main.go

# 2. Start frontend in another terminal
npm run dev

# 3. Access http://localhost:5173 (frontend)
# Frontend proxies API requests to http://localhost:8080

# 4. Make changes:
# - Backend: Format with goimports, test with go test
# - Frontend: Format with prettier, preview in browser

# 5. Debug both:
# - Backend: F5 in VS Code (debug config for Go)
# - Frontend: Browser DevTools (F12)
```

### Troubleshooting

**Both formatters running and conflicting**:
```bash
# Disable one formatter per file type
# Keep goimports for Go only
# Keep prettier for frontend only
```

**LSP conflicts**:
```bash
# gopls and prettier shouldn't conflict
# But if they do, restart both:
pkill gopls
code --reload-extension esbenp.prettier-vscode
```

---

## Profile 4: DevOps Engineer

**Focus**: Kubernetes, Docker, infrastructure, configuration

**Best For**: Building deployment pipelines, managing infrastructure

### Setup Steps

#### 1. Install Tools

```bash
# Docker
docker --version
docker-compose --version

# Kubernetes
kubectl version
helm version
k3s --version  # if using K3s

# Terraform
terraform --version

# Cloud CLI (AWS, GCP, Azure)
aws --version  # or gcloud, az
```

#### 2. Configure VS Code for Infrastructure

```bash
# Install extensions
code --install-extension ms-azuretools.vscode-docker
code --install-extension hashicorp.terraform
code --install-extension ms-kubernetes-tools.vscode-kubernetes-tools
code --install-extension redhat.vscode-yaml
```

#### 3. Set Up Kubernetes Context

```bash
# Check available contexts
kubectl config get-contexts

# Switch context if needed
kubectl config use-context <context-name>

# Verify connection
kubectl cluster-info
```

#### 4. Configure Cloud Credentials

```bash
# AWS (if using)
mkdir -p ~/.aws
cat > ~/.aws/config << 'EOF'
[default]
region = us-east-1
EOF

# GCP
gcloud auth login
gcloud config set project <project-id>

# Or use environment variables
export AWS_ACCESS_KEY_ID=...
export REVENGE_DATABASE_URL=...
```

### Profile Configuration

#### VS Code Settings

```json
{
  // Docker
  "[dockerfile]": {
    "editor.defaultFormatter": "ms-azuretools.vscode-docker"
  },

  // Kubernetes/YAML
  "[yaml]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 2,
    "editor.insertSpaces": true
  },

  // Terraform
  "[hcl]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 2
  },
  "terraform.path": "terraform",
  "terraform.format": {
    "on_save": true
  },

  // JSON (for config)
  "[json]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 2
  },

  // Terminal (important for DevOps)
  "terminal.integrated.defaultProfile.linux": "bash",
  "terminal.integrated.fontSize": 12,
  "terminal.integrated.cwd": "${workspaceFolder}"
}
```

### Common Tasks

#### Deploy with Docker Compose

```bash
# Start all services
docker-compose -f docker-compose.yml up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f revenge

# Stop
docker-compose down
```

#### Deploy to Kubernetes

```bash
# Create namespace
kubectl create namespace revenge

# Apply manifests
kubectl apply -f k8s/configmap.yaml -n revenge
kubectl apply -f k8s/deployment.yaml -n revenge

# Check status
kubectl get pods -n revenge

# View logs
kubectl logs -f <pod-name> -n revenge

# Debug pod
kubectl exec -it <pod-name> -n revenge -- bash
```

#### Terraform Deployment

```bash
# Initialize
terraform init

# Plan changes
terraform plan -out=tfplan

# Apply
terraform apply tfplan

# Destroy
terraform destroy
```

#### SSH to Remote Server

```bash
# Connect via SSH
ssh <user>@<host>

# Or use Coder for managed environment
coder ssh <workspace-name>
```

---

## Profile 5: Documentation Writer

**Focus**: Markdown editing, preview, and validation

**Best For**: Writing documentation, design docs, guides

### Setup Steps

#### 1. Install Minimal Extensions

```bash
# VS Code
code --install-extension yzhang.markdown-all-in-one
code --install-extension DavidAnson.vscode-markdownlint
code --install-extension esbenp.prettier-vscode

# Or use Zed (lighter, no extensions needed)
```

#### 2. Configure Markdown Settings

```json
{
  "[markdown]": {
    "editor.formatOnSave": true,
    "editor.wordWrap": "on",
    "editor.defaultFormatter": "esbenp.prettier-vscode",
    "editor.tabSize": 2,
    "files.trimTrailingWhitespace": false
  },

  "markdown.preview.breaks": true,
  "markdown.validate.enabled": true,
  "markdownlint.config": {
    "MD007": { "indent": 2 }
  }
}
```

#### 3. Use Zed (Recommended)

Zed is lighter and perfect for doc writing:

```json
{
  "languages": {
    "Markdown": {
      "soft_wrap": "editor_width",
      "preferred_line_length": 100
    }
  }
}
```

### Common Tasks

#### Write a Design Document

```bash
# Start from template
cp docs/dev/design/TEMPLATE.md docs/dev/design/my-feature/MY_FEATURE.md

# Edit in VS Code or Zed
code docs/dev/design/my-feature/MY_FEATURE.md

# Preview (VS Code: Cmd+K V, Zed: terminal command)
```

#### Validate Markdown

```bash
# Check syntax
markdownlint docs/

# Check links
python3 scripts/validate-links.py

# Check structure
python3 scripts/validate-doc-structure.py
```

#### Format Documentation

```bash
# Format all markdown
prettier --write '**/*.md'

# Check formatting
prettier --check '**/*.md'
```

---

## Profile 6: Remote Development (Coder-Optimized)

**Focus**: SSH-based remote development via Coder

**Best For**: Development in containerized environments, CI/CD pipelines

### Setup Steps

#### 1. Install Coder CLI

```bash
# macOS
brew install coder

# Linux
curl -L https://coder.com/install.sh | sh

# Verify
coder version
```

#### 2. Connect to Coder Workspace

```bash
# List available workspaces
coder list

# Connect via SSH
coder ssh <workspace-name>

# Or create new workspace
coder create <workspace-name> --template=revenge
```

#### 3. Open IDE in Remote

**Zed (Recommended)**:
```bash
# From Coder SSH session
zed /home/coder/workspace/revenge &

# Or from local machine with Zed remote SSH extension
# Zed → Open Remote → SSH → <workspace-ssh-config>
```

**VS Code (via SSH)** :
```bash
# Install Remote - SSH extension
code --install-extension ms-vscode-remote.remote-ssh

# Connect to Coder workspace
# Cmd+Shift+P → "Remote-SSH: Connect to Host"
# Select coder instance from SSH config
```

#### 4. Configure Remote Settings

`.zed/settings.json` (optimized for remote):

```json
{
  "tab_size": 2,
  "format_on_save": "on",

  "languages": {
    "Go": {
      "tab_size": 4,
      "hard_tabs": true
    }
  },

  "lsp": {
    "gopls": {
      "initialization_options": {
        "analyses": {
          "unusedparams": false
        }
      }
    }
  },

  "file_scan_exclusions": [
    "**/.git",
    "**/node_modules",
    "**/__pycache__",
    "**/vendor"
  ]
}
```

**VS Code Remote SSH** (`.vscode/settings.json`):

```json
{
  "remote.SSH.showLoginTerminal": true,
  "remote.SSH.useLocalServer": true,

  "editor.formatOnSave": true,
  "go.useLanguageServer": true
}
```

#### 5. Configure Coder Workspace

`.coder/template.tf` (Terraform configuration):

```hcl
resource "coder_workspace" "revenge" {
  name = "revenge-${lower(data.coder_workspace.me.name)}"

  agents {
    id = coder_agent.main.id
  }

  tags = {
    os       = "linux"
    arch     = "amd64"
    lifetime = "24h"
  }
}

resource "coder_agent" "main" {
  os       = "linux"
  arch     = "amd64"
  startup_script = file("${path.module}/startup.sh")
}
```

### Common Remote Workflow

```bash
# 1. Connect to workspace
coder ssh revenge-dev

# 2. Open editor
zed /home/coder/workspace/revenge &

# 3. Edit, test, commit as usual
cd /home/coder/workspace/revenge
go build ./...
go test ./...

# 4. Push changes
git push origin feature-branch

# 5. Create PR from GitHub/GitLab

# 6. Disconnect when done
exit
```

### Remote-Specific Tips

#### Run Background Services

```bash
# Start PostgreSQL
docker run -d --name revenge-db -p 5432:5432 postgres:18

# Start Redis/Dragonfly
docker run -d --name revenge-cache -p 6379:6379 dragonflydb/dragonfly

# Connect from Revenge
export REVENGE_DATABASE_URL="postgres://postgres:postgres@localhost:5432/revenge"
```

#### Port Forwarding

```bash
# From local machine, access remote services
# Coder handles this automatically via SSH tunneling

# Or explicit port forward
ssh -L 8080:localhost:8080 coder@workspace

# Now http://localhost:8080 connects to remote
```

#### File Sync

```bash
# One-way sync from local to remote
rsync -avz --delete . coder@workspace:~/workspace/revenge

# Or use VS Code's built-in file sync
# Automatically syncs on save
```

#### Monitor Resource Usage

```bash
# Inside Coder workspace
htop  # Monitor CPU/memory
docker stats  # Monitor containers
```

---

## Installation and Customization

### Installing a Profile

#### Automatic Installation (Script - Planned)

```bash
# Install specific profile
./scripts/setup-profile.sh backend-developer

# Install with Zed
./scripts/setup-profile.sh backend-developer --editor=zed

# Install with VS Code
./scripts/setup-profile.sh backend-developer --editor=vscode
```

#### Manual Installation

1. **Identify your profile** from the list above

2. **Copy settings**:
   ```bash
   # For VS Code
   cp .shared/profiles/[profile-name]/.vscode/* .vscode/

   # For Zed
   cp .shared/profiles/[profile-name]/.zed/* .zed/
   ```

3. **Verify tools are installed**:
   ```bash
   # Backend
   go version
   gopls version
   golangci-lint version

   # Frontend
   node --version
   npm list svelte sveltekit

   # Others per profile
   ```

4. **Test the setup**:
   ```bash
   # Go backend
   go build ./...

   # Frontend
   npm run build

   # Both
   make test
   ```

### Customizing a Profile

Profiles are starting points, not restrictions. Customize as needed:

#### Add More Extensions

```json
{
  "recommendations": [
    // ... existing extensions ...
    "your.new-extension"
  ]
}
```

#### Modify Settings

```json
{
  // Keep what works, change what doesn't
  "editor.fontSize": 14,  // Personal preference
  "editor.fontFamily": "Fira Code",

  // But keep language-specific settings aligned with EditorConfig
  "[go]": {
    "editor.tabSize": 4,  // Must stay 4 for Go
    "editor.insertSpaces": false
  }
}
```

#### Create Custom Profile

```bash
# Copy existing profile
cp -r .shared/profiles/backend-developer .shared/profiles/my-profile

# Customize files
# Then use:
cp .shared/profiles/my-profile/.vscode/* .vscode/
```

---

## Profile Comparison Matrix

| Feature | Backend | Frontend | Full-Stack | DevOps | Docs | Remote |
|---------|---------|----------|-----------|--------|------|--------|
| **Go support** | ✅ Full | ⚠️ Basic | ✅ Full | ⚠️ Basic | ❌ | ✅ Full |
| **Svelte/TS** | ❌ | ✅ Full | ✅ Full | ❌ | ❌ | ⚠️ Basic |
| **Kubernetes** | ❌ | ❌ | ❌ | ✅ Full | ❌ | ✅ Full |
| **Docker** | ⚠️ Basic | ❌ | ⚠️ Basic | ✅ Full | ❌ | ✅ Full |
| **Database tools** | ✅ Full | ❌ | ✅ Full | ⚠️ Basic | ❌ | ⚠️ Basic |
| **Debug support** | ✅ Full | ✅ Full | ✅ Full | ⚠️ Basic | ❌ | ⚠️ Basic |
| **Extensions** | 20+ | 15+ | 30+ | 12+ | 3+ | 0 (SSH) |
| **Setup time** | 5 min | 5 min | 10 min | 10 min | 3 min | 15 min |
| **Recommended IDE** | VS Code | VS Code | VS Code | VS Code | Zed | Zed |

---

## Troubleshooting Profile Setup

### Problem: Extensions Won't Install

**Solution**:
```bash
# Install manually
code --install-extension golang.go
code --install-extension svelte.svelte-vscode

# Or use VS Code UI
# Cmd+Shift+P → "Extensions: Install from VSIX"
```

### Problem: Settings Not Taking Effect

**Solution**:
```bash
# Check settings file exists
ls -la .vscode/settings.json
ls -la .zed/settings.json

# Reload editor
# Cmd+Shift+P → "Reload Window" (VS Code)
# Cmd+Shift+P → "Reload Window" (Zed)
```

### Problem: Different Settings Per Machine

**Solution**:
```bash
# Use User settings for personal preferences
# Use Workspace settings (in .vscode/) for team consistency

# Check precedence
# User settings: ~/.config/Code/User/settings.json
# Workspace settings: .vscode/settings.json (this wins)
```

### Problem: LSP Not Using New Settings

**Solution**:
```bash
# Restart LSP
# Cmd+Shift+P → "Restart Language Server"

# Or kill the process
pkill gopls
pkill ruff-lsp

# Reload
code . -r
```

---

## FAQ

**Q: Can I use multiple profiles?**
A: Yes, pick the one closest to your role, then customize. You can merge settings from multiple profiles.

**Q: How do I switch between profiles?**
A: Copy new profile settings over existing ones:
```bash
cp .shared/profiles/new-profile/.vscode/* .vscode/
```

**Q: Should I commit profile changes to git?**
A: Yes, workspace settings (`.vscode/`, `.zed/`) should be committed. User settings should not.

**Q: What if my organization has different standards?**
A: Modify EditorConfig (`.editorconfig`) first, then all IDE settings will align automatically.

**Q: Can I use both VS Code and Zed?**
A: Yes, switch between them. Both read EditorConfig, so formatting stays consistent.

---

## Related Documentation

- [SETTINGS_GUIDE.md](SETTINGS_GUIDE.md) - Detailed settings reference
- [TOOL_COMPARISON.md](TOOL_COMPARISON.md) - VS Code vs Zed comparison
- [INTEGRATION.md](INTEGRATION.md) - How tools integrate
- [ONBOARDING.md](ONBOARDING.md) - New developer setup
- [SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md) - Technology stack

---

**Maintained By**: Development Team
**Last Updated**: 2026-01-31
