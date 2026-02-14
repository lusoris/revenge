# Frontend Container — Fixes & Monitoring

## Session: 2026-02-14 (continued)

### Issues Found & Fixed

#### 1. Brotli Module Missing
- **Problem**: `nginx:alpine` doesn't include the brotli module → `unknown directive "brotli_static"`
- **Attempted fix**: Removed brotli directives (rejected — don't degrade quality)
- **Final fix**: Switched to `fholzer/nginx-brotli:latest` (Alpine-based, brotli compiled in)
- **Gotcha**: entrypoint is `["nginx"]`, so CMD must be `["-g", "daemon off;"]` not `["nginx", ...]`

#### 2. Port 3000 Conflict
- **Problem**: Grafana already on port 3000 in tools profile
- **Fix**: Frontend mapped to `"4000:3000"` in docker-compose.dev.yml

#### 3. Healthcheck IPv6
- **Problem**: Alpine `wget` resolves `localhost` to `[::1]` but nginx only listens on IPv4
- **Fix**: Changed all healthchecks to `http://127.0.0.1:3000/health` (Dockerfile, dev compose, prod compose)

#### 4. API Proxy 404
- **Problem**: nginx `proxy_pass http://revenge:8096/api/` forwards `/api/healthz` to `/api/healthz` on backend, but backend routes are at `/` (ogen registers at root)
- **Fix**: Changed to `proxy_pass http://revenge:8096/` — the trailing slash + `/api/` location means nginx strips the prefix: `/api/healthz` → backend `/healthz`

### Monitoring Setup

#### nginx-prometheus-exporter
- Added as sidecar in docker-compose.dev.yml (tools profile)
- Image: `nginx/nginx-prometheus-exporter:latest`
- Scrapes `http://revenge-frontend-dev:3000/stub_status`
- Exposes Prometheus metrics on port 9113

#### Prometheus
- Added `nginx-frontend` scrape job targeting `nginx-exporter:9113`
- Both targets verified UP

#### Grafana Dashboard: Revenge / Frontend
- File: `deploy/grafana/provisioning/dashboards/revenge-frontend.json`
- UID: `revenge-frontend`, tags: `revenge`, `frontend`, `nginx`
- Panel IDs: 5000-5303 (no conflicts with existing dashboards)
- Sections:
  - 🌐 Frontend Health — up/down, active conns, req/s, reading/writing/waiting
  - 📈 Request Traffic — HTTP request rate + connection states (timeseries)
  - 🔗 Connection Throughput — accepted vs handled, dropped connections rate
  - 📊 Cumulative Totals — total requests, accepted, handled since start
- **Existing dashboards untouched**: overview, features, infrastructure, playback

### Commits Pushed (3)
1. `5598f28e` — fix(web): brotli support, API proxy routing, healthcheck IPv4
2. `31d1f83d` — fix(docker): frontend port 4000, nginx-exporter sidecar, healthcheck IPv4
3. `f4062560` — feat(monitoring): frontend Grafana dashboard + Prometheus nginx scraping

### Verified Working
- ✅ `curl http://localhost:4000/` — serves SPA with `Content-Encoding: br`
- ✅ `curl http://localhost:4000/api/healthz` — proxied to backend, returns healthy
- ✅ `curl http://localhost:4000/health` — nginx health endpoint
- ✅ Prometheus targets: both `revenge` and `nginx-frontend` are UP
- ✅ `nginx_connections_active` metric flowing in Prometheus
- ✅ Grafana: 5 dashboards provisioned (4 existing + 1 new frontend)
- ✅ Frontend container: healthy
