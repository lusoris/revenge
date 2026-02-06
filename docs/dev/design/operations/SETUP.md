# Deployment

<!-- DESIGN: operations -->

**Docker Image**: `ghcr.io/lusoris/revenge`
**Helm Chart**: `ghcr.io/lusoris/charts/revenge`

> Production deployment via Docker Compose or Kubernetes Helm chart

---

## Docker Image

Multi-stage build (Dockerfile):

**Builder stage** (`golang:1.25-alpine`):
- Build deps: git, make, gcc, musl-dev, pkgconfig, ffmpeg-dev, vips-dev
- CGO_ENABLED=1, GOEXPERIMENT=greenteagc,jsonv2
- Output: stripped binary (`-w -s`)

**Runtime stage** (`alpine:latest`):
- Runtime deps: ca-certificates, ffmpeg, ffmpeg-libs, postgresql-client, tzdata
- Non-root user: `revenge` (uid/gid 1000)
- Port: 8096
- Volumes: `/data`, `/config`, `/cache`, `/media` (ro)
- Healthcheck: `GET /health/live` (30s interval, 3 retries)
- Entrypoint: `/docker-entrypoint.sh` (waits for PostgreSQL, runs migrations, starts server)

## Docker Compose

### docker-compose.yml — Lightweight Production

Services: revenge (GHCR image), postgres:18, dragonfly, typesense:0.25.2

Configurable via environment: `REVENGE_VERSION`, `DB_PASSWORD`, `TYPESENSE_API_KEY`, `MEDIA_PATH`.

### docker-compose.prod.yml — Full Production Build

Services: revenge (local build), postgres:18, dragonfly, typesense:0.25.2

- Dragonfly: `--maxmemory=2gb`, `--maxmemory-policy=allkeys-lru`
- Stricter health checks (10s interval, 10s start-period)
- Required env vars fail with `:?` if unset

### docker-compose.dev.yml — Development

Services: revenge (local build), postgres:18-alpine, dragonfly, typesense, pgadmin (optional)

- Revenge ports: 8096 (HTTP), 2345 (Delve debugger)
- Debug logging, dev credentials

### Overlay Files

| File | Purpose |
|------|---------|
| docker-compose.nfs.yml | NFS media mount (NFSv4.1, read-only) |
| docker-compose.s3.yml | MinIO/S3 storage (avatar bucket, auto-created) |

## Helm Chart

`charts/revenge/` — Kubernetes deployment with all dependencies included.

### Dependencies

| Component | Default | Persistence |
|-----------|---------|-------------|
| PostgreSQL | postgres:18-alpine | 10Gi |
| Dragonfly | latest (`--cache_mode`) | - |
| Typesense | 0.25.2 | 5Gi |

### Key Values

```yaml
image:
  repository: ghcr.io/lusoris/revenge
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8096

resources:
  limits:
    memory: 1Gi
  requests:
    cpu: 100m
    memory: 256Mi

revenge:
  server:
    port: 8096
  database:
    host: revenge-postgresql
    port: 5432
  cache:
    host: revenge-dragonfly
    port: 6379
  search:
    host: revenge-typesense
    port: 8108
```

Probes: liveness at `/health/live` (15s delay, 30s period), readiness at `/health/ready` (5s delay, 10s period).

Ingress disabled by default. Media persistence optional with NFS support.

### Templates

deployment.yaml, service.yaml, ingress.yaml (optional), media-pvc.yaml (conditional), _helpers.tpl.

## Quick Start

```bash
# Docker Compose (simplest)
docker compose -f docker-compose.yml up -d

# Production build
docker compose -f docker-compose.prod.yml up -d

# Kubernetes
helm install revenge oci://ghcr.io/lusoris/charts/revenge
```

## Related Documentation

- [CI_CD.md](CI_CD.md) - Automated builds and releases
- [DEVELOPMENT.md](DEVELOPMENT.md) - Development environment setup
