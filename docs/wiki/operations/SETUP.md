## Table of Contents

- [Production Deployment Setup](#production-deployment-setup)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Troubleshooting](#troubleshooting)
    - [Container won't start](#container-wont-start)
    - [Cannot connect to database](#cannot-connect-to-database)
    - [Reverse proxy returns 502 Bad Gateway](#reverse-proxy-returns-502-bad-gateway)
    - [Slow performance / high CPU](#slow-performance-high-cpu)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Dragonfly Documentation
    url: ../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
  - name: FFmpeg Documentation
    url: ../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
  - name: FFmpeg Codecs
    url: ../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
  - name: FFmpeg Formats
    url: ../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
  - name: go-astiav (FFmpeg bindings)
    url: ../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
  - name: go-astiav GitHub README
    url: ../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
  - name: Go io
    url: ../sources/go/stdlib/io.md
    note: Auto-resolved from go-io
  - name: pgx PostgreSQL Driver
    url: ../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
  - name: rueidis
    url: ../sources/tooling/rueidis.md
    note: Auto-resolved from rueidis
  - name: rueidis GitHub README
    url: ../sources/tooling/rueidis-guide.md
    note: Auto-resolved from rueidis-docs
  - name: Typesense API
    url: ../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: ../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
  - title: operations
    path: operations/INDEX.md
  - title: TECH_STACK
    path: technical/TECH_STACK.md
  - title: REVERSE_PROXY
    path: operations/REVERSE_PROXY.md
  - title: 00_SOURCE_OF_TRUTH
    path: 00_SOURCE_OF_TRUTH.md
---

# Production Deployment Setup




> Deploy Revenge to your homelab or server with Docker, Kubernetes, or bare metal


This guide covers everything you need to deploy Revenge in production. The recommended approach uses Docker Compose for simple homelab setups - just create a docker-compose.yml, set a few environment variables, and run docker compose up. For larger deployments, Kubernetes (K3s) provides high availability and auto-scaling. Includes reverse proxy configuration (Traefik, Caddy, nginx), SSL setup with LetsEncrypt, backup strategies, and security hardening tips.


---




## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---


## How It Works

<!-- User-friendly explanation -->




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->








## Troubleshooting
### Container won't start

**Cause**: Environment misconfiguration or dependency issues

**Solution**: 1. Check logs: `docker compose logs revenge`
2. Verify environment variables in .env
3. Ensure PostgreSQL is ready: `docker compose ps postgres`

### Cannot connect to database

**Cause**: PostgreSQL not ready or network issue

**Solution**: 1. Check PostgreSQL health: `docker exec revenge-postgres pg_isready`
2. Verify DATABASE_URL in environment
3. Check network: `docker compose exec revenge ping postgres`

### Reverse proxy returns 502 Bad Gateway

**Cause**: Backend not running or unhealthy

**Solution**: 1. Verify Revenge is running: `docker compose ps`
2. Check backend is healthy: `curl http://localhost:8080/health/live`
3. Review proxy logs (Traefik/nginx/Caddy)

### Slow performance / high CPU

**Cause**: Resource constraints or heavy transcoding load

**Solution**: 1. Check Docker resource limits
2. Review active transcoding jobs: `docker compose logs revenge | grep transcode`
3. Increase dragonfly cache size
4. Consider offloading transcoding to Blackbeard


## Related Documentation
### See Also
<!-- Related wiki pages -->



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)