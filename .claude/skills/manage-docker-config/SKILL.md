---
name: manage-docker-config
description: Manage Docker configurations (sync from SOT, build/push images, compose operations)
argument-hint: "[--sync|--validate|--build|--push|--up|--down] [--tag TAG] [--registry REGISTRY]"
disable-model-invocation: false
allowed-tools: Bash(*), Read(*), Write(*)
---

# Manage Docker Config

Manages Docker configurations including syncing from SOURCE_OF_TRUTH, building/pushing images, and docker-compose operations.

## Usage

```
/manage-docker-config --sync             # Sync Dockerfile and docker-compose.yml from SOT
/manage-docker-config --validate         # Validate Docker configurations
/manage-docker-config --build            # Build Docker images
/manage-docker-config --build --tag v1.0 # Build with specific tag
/manage-docker-config --push             # Push images to registry
/manage-docker-config --up               # Start services (docker-compose up -d)
/manage-docker-config --down             # Stop services (docker-compose down)
```

## Arguments

- `$0`: Action (--sync, --validate, --build, --push, --up, --down)
- `$1`: Tag (optional: --tag TAG for build/push)
- `$2`: Registry (optional: --registry REGISTRY for push)

## Prerequisites

- Python 3.10+ installed
- Docker and docker-compose installed
- Docker daemon running
- Access to container registry (for push operations)
- SOURCE_OF_TRUTH.md with infrastructure versions

## Task

Manage Docker configurations and operations with version synchronization from SOURCE_OF_TRUTH.

### Step 1: Verify Script and Docker

```bash
if [ ! -f "scripts/automation/manage_docker.py" ]; then
    echo "❌ Docker management script not found"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found. Install Docker first"
    exit 1
fi

if ! docker info &> /dev/null; then
    echo "❌ Docker daemon not running"
    exit 1
fi
```

### Step 2: Sync Configurations from SOT

**Sync Dockerfile and docker-compose.yml**:
```bash
# Sync versions from SOURCE_OF_TRUTH.md
python scripts/automation/manage_docker.py --sync-dockerfile
python scripts/automation/manage_docker.py --sync-compose

# Dry-run mode (preview changes)
python scripts/automation/manage_docker.py --sync-dockerfile --dry-run
```

**What gets synced**:
- PostgreSQL version from SOT Infrastructure table
- Go version from SOT tech stack
- Node.js version from SOT tech stack
- Alpine Linux version from SOT
- Service versions in docker-compose.yml

### Step 3: Validate Configurations

**Validate all Docker configs**:
```bash
python scripts/automation/manage_docker.py --validate
```

**Validation checks**:
- Dockerfile exists and is valid
- docker-compose.yml exists and is valid
- Base images are accessible
- Service definitions are correct
- Volume mounts exist
- Network configurations are valid

### Step 4: Build Images

**Build all images**:
```bash
# Build with latest tag
python scripts/automation/manage_docker.py --build

# Build with specific tag
python scripts/automation/manage_docker.py --build --tag v1.0.0

# Build with custom registry
python scripts/automation/manage_docker.py --build --tag v1.0.0 --registry ghcr.io/user/repo
```

**Build targets**:
- Backend (Go application)
- Frontend (Node.js/SvelteKit)
- Migrations (database migrations)

### Step 5: Push Images

**Push to registry**:
```bash
# Push to default registry (from env or config)
python scripts/automation/manage_docker.py --push --tag v1.0.0

# Push to specific registry
python scripts/automation/manage_docker.py --push --tag v1.0.0 --registry ghcr.io/user/repo

# Login to registry first if needed
docker login ghcr.io
```

### Step 6: Docker Compose Operations

**Start services**:
```bash
# Start all services in background
python scripts/automation/manage_docker.py --up

# Start in foreground (see logs)
docker-compose up

# Start specific services
docker-compose up -d postgres dragonfly
```

**Stop services**:
```bash
# Stop all services
python scripts/automation/manage_docker.py --down

# Stop and remove volumes
docker-compose down -v
```

**View running services**:
```bash
docker-compose ps
```

## Configuration Sync Details

### Dockerfile Sync

SOURCE_OF_TRUTH.md:
```markdown
| Go | 1.25+ | Backend | ✅ |
| Node.js | 20.x | Frontend | ✅ |
| Alpine | 3.20+ | Base Image | ✅ |
```

Synced to Dockerfile:
```dockerfile
FROM golang:1.25-alpine3.20 AS backend-builder
FROM node:20-alpine3.20 AS frontend-builder
```

### docker-compose.yml Sync

SOURCE_OF_TRUTH.md:
```markdown
| PostgreSQL | 18+ | Database | ✅ |
| Dragonfly | latest | Cache | ✅ |
| Typesense | 27.1 | Search | ✅ |
```

Synced to docker-compose.yml:
```yaml
services:
  postgres:
    image: postgres:18-alpine

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly:latest

  typesense:
    image: typesense/typesense:27.1
```

## Examples

**Full sync and rebuild**:
```bash
# 1. Sync configs from SOT
/manage-docker-config --sync

# 2. Validate configs
/manage-docker-config --validate

# 3. Rebuild images
/manage-docker-config --build --tag latest

# 4. Restart services
/manage-docker-config --down
/manage-docker-config --up
```

**Release workflow**:
```bash
# 1. Sync to latest versions
/manage-docker-config --sync

# 2. Build with version tag
/manage-docker-config --build --tag v1.2.3

# 3. Push to registry
/manage-docker-config --push --tag v1.2.3 --registry ghcr.io/user/repo

# 4. Update production
# (deployment-specific commands)
```

**Development workflow**:
```bash
# Start services for development
/manage-docker-config --up

# Check status
docker-compose ps

# View logs
/view-logs --docker postgres

# Stop when done
/manage-docker-config --down
```

## Troubleshooting

**"Version not found in SOURCE_OF_TRUTH"**:
```bash
# Check SOT has required versions
cat docs/dev/design/00_SOURCE_OF_TRUTH.md | grep -A5 "Infrastructure"
```

**"Docker daemon not running"**:
```bash
# Start Docker daemon
# macOS: Open Docker Desktop
# Linux: sudo systemctl start docker
```

**"Permission denied accessing /var/run/docker.sock"**:
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Re-login or run:
newgrp docker
```

**"Image push failed - authentication required"**:
```bash
# Login to container registry
docker login ghcr.io

# Or use token
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

**"Service already running"**:
```bash
# Stop existing services first
/manage-docker-config --down

# Then start
/manage-docker-config --up
```

## Tips

1. **Always sync before building**:
   Ensures you're building with latest versions from SOT

2. **Use tags for releases**:
   Never push `latest` to production - use semantic versioning

3. **Dry-run for safety**:
   Use `--dry-run` flag to preview changes before applying

4. **Monitor resource usage**:
   ```bash
   docker stats
   ```

5. **Clean up regularly**:
   ```bash
   # Remove unused images
   docker image prune -a

   # Remove unused volumes
   docker volume prune

   # Full cleanup
   docker system prune -a --volumes
   ```

## Exit Codes

- `0`: Success
- `1`: Failure (error in operation)

## Related Skills

- `/check-health` - Check system health including Docker services
- `/view-logs` - View Docker container logs
- `/manage-ci-workflows` - Manage CI/CD workflows
