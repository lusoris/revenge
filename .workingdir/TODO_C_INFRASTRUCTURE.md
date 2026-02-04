# Phase C: Infrastructure & Release

**Effort**: 8-16h
**Dependencies**: Phase B complete

---

## C1: Docker Compose Stack

- [ ] `docker-compose.yml` with all services:
  - revenge (backend)
  - revenge-frontend
  - postgresql
  - dragonfly
  - typesense
  - traefik (reverse proxy)
- [ ] `.env.example` with all config
- [ ] Health check integration
- [ ] Volume mounts for persistence

---

## C2: Docker Images

- [ ] Backend multi-stage Dockerfile (verified)
- [ ] Frontend multi-stage Dockerfile
- [ ] Combined nginx config
- [ ] GitHub Actions for image builds

---

## C3: Documentation

- [ ] Getting started guide
- [ ] Installation guide (Docker)
- [ ] Configuration reference
- [ ] Radarr setup guide
- [ ] API authentication guide

---

## C4: MVP Verification

- [ ] Movies display in frontend
- [ ] Search works end-to-end
- [ ] Radarr sync imports movies
- [ ] Watch progress saves and restores
- [ ] Authentication works (login/logout)
- [ ] MFA works (TOTP + backup codes)
- [ ] RBAC enforced on admin pages
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes
- [ ] Docker Compose stack works
