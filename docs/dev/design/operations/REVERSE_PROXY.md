# Revenge - Reverse Proxy & Deployment Best Practices

<!-- SOURCES: dragonfly, gohlslib, m3u8, pgx, postgresql-arrays, postgresql-json, prometheus, prometheus-metrics, river, rueidis, rueidis-docs -->

<!-- DESIGN: operations, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Production deployment with nginx, Caddy, Traefik, and Docker optimization.


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Reverse Proxy Configurations](#reverse-proxy-configurations)
  - [Nginx (Recommended)](#nginx-recommended)
  - [Caddy (Simpler Alternative)](#caddy-simpler-alternative)
  - [Traefik (Docker-native)](#traefik-docker-native)
- [Server Configuration](#server-configuration)
  - [Revenge Backend Settings](#revenge-backend-settings)
  - [Headers Middleware](#headers-middleware)
- [Streaming Optimization](#streaming-optimization)
  - [Byte Range Support](#byte-range-support)
  - [Nginx Slice Module (for large files)](#nginx-slice-module-for-large-files)
- [Caching Strategy](#caching-strategy)
  - [Nginx Proxy Cache](#nginx-proxy-cache)
  - [CDN Headers](#cdn-headers)
- [Load Balancing](#load-balancing)
  - [Multiple Backend Instances](#multiple-backend-instances)
  - [Docker Swarm / Kubernetes](#docker-swarm-kubernetes)
- [Security Best Practices](#security-best-practices)
  - [Fail2ban Configuration](#fail2ban-configuration)
  - [WAF Rules (ModSecurity)](#waf-rules-modsecurity)
  - [IP Allowlist for Admin](#ip-allowlist-for-admin)
- [Monitoring & Observability](#monitoring-observability)
  - [Prometheus Metrics Endpoint](#prometheus-metrics-endpoint)
  - [Nginx Status](#nginx-status)
  - [Log Format for Analysis](#log-format-for-analysis)
- [Docker Production Setup](#docker-production-setup)
  - [Optimized Dockerfile](#optimized-dockerfile)
  - [Docker Compose Production](#docker-compose-production)
- [Configuration Summary](#configuration-summary)
  - [Environment Variables](#environment-variables)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Quick Reference](#quick-reference)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | 🔴 |
| Sources | 🔴 |
| Instructions | 🔴 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
---

## Overview

Revenge is designed to run behind a reverse proxy for:
- TLS termination
- Load balancing
- Caching static assets
- Rate limiting
- WebSocket support
- Large file uploads (media)

---

## Reverse Proxy Configurations

### Nginx (Recommended)

```nginx
# /etc/nginx/sites-available/revenge
upstream revenge_backend {
    server 127.0.0.1:8096;
    keepalive 64;
}

# Rate limiting zones
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=5r/m;
limit_conn_zone $binary_remote_addr zone=conn_limit:10m;

server {
    listen 80;
    server_name revenge.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name revenge.example.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/revenge.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/revenge.example.com/privkey.pem;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_session_tickets off;

    # Modern SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # HSTS
    add_header Strict-Transport-Security "max-age=63072000" always;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Logging
    access_log /var/log/nginx/revenge_access.log;
    error_log /var/log/nginx/revenge_error.log;

    # Client body size (for uploads)
    client_max_body_size 0;  # Unlimited for large media uploads

    # Timeouts for streaming
    proxy_connect_timeout 60s;
    proxy_send_timeout 3600s;
    proxy_read_timeout 3600s;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml application/json application/javascript application/xml+rss application/atom+xml image/svg+xml;

    # Static assets with long cache
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        proxy_pass http://revenge_backend;
        proxy_http_version 1.1;
        proxy_set_header Connection "";

        # Cache for 1 year (immutable assets)
        add_header Cache-Control "public, max-age=31536000, immutable";

        # Nginx caching
        proxy_cache_valid 200 1y;
        proxy_cache_use_stale error timeout updating;
    }

    # API endpoints
    location /api/ {
        limit_req zone=api_limit burst=20 nodelay;
        limit_conn conn_limit 50;

        proxy_pass http://revenge_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Connection "";

        # No caching for API
        add_header Cache-Control "no-store, no-cache, must-revalidate";
    }

    # Auth endpoints (stricter rate limiting)
    location /api/v1/auth/ {
        limit_req zone=auth_limit burst=5 nodelay;

        proxy_pass http://revenge_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket support
    location /api/v1/ws {
        proxy_pass http://revenge_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket timeout
        proxy_read_timeout 86400s;
        proxy_send_timeout 86400s;
    }

    # Video/Audio streaming
    location ~* ^/api/v1/(stream|video|audio)/ {
        proxy_pass http://revenge_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Range $http_range;
        proxy_set_header If-Range $http_if_range;
        proxy_set_header Connection "";

        # Buffering for streaming
        proxy_buffering on;
        proxy_buffer_size 128k;
        proxy_buffers 4 256k;
        proxy_busy_buffers_size 256k;

        # Allow seeking
        proxy_force_ranges on;

        # Long timeout for streaming
        proxy_read_timeout 3600s;
    }

    # HLS segments (cacheable)
    location ~* \.(m3u8|ts|m4s)$ {
        proxy_pass http://revenge_backend;
        proxy_http_version 1.1;
        proxy_set_header Connection "";

        # Cache segments
        proxy_cache_valid 200 1m;
        add_header Cache-Control "public, max-age=60";
    }

    # Default location
    location / {
        proxy_pass http://revenge_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Connection "";
    }
}
```

### Caddy (Simpler Alternative)

```caddyfile
# /etc/caddy/Caddyfile
revenge.example.com {
    # Enable compression
    encode gzip zstd

    # Rate limiting
    rate_limit {
        zone api {
            key {remote_host}
            events 10
            window 1s
        }
        zone auth {
            key {remote_host}
            events 5
            window 1m
        }
    }

    # Auth endpoints (strict rate limit)
    @auth path /api/v1/auth/*
    rate_limit @auth {
        zone auth
    }

    # API endpoints
    @api path /api/*
    rate_limit @api {
        zone api
    }

    # WebSocket
    @websocket {
        path /api/v1/ws
        header Connection *Upgrade*
        header Upgrade websocket
    }

    # Static assets
    @static {
        path *.js *.css *.png *.jpg *.jpeg *.gif *.ico *.svg *.woff *.woff2
    }
    header @static Cache-Control "public, max-age=31536000, immutable"

    # HLS segments
    @hls {
        path *.m3u8 *.ts *.m4s
    }
    header @hls Cache-Control "public, max-age=60"

    # Reverse proxy
    reverse_proxy localhost:8096 {
        # Health check
        health_uri /api/v1/health
        health_interval 30s

        # Streaming timeout
        transport http {
            read_timeout 3600s
            write_timeout 3600s
        }

        # Headers
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-Proto {scheme}
    }
}
```

### Traefik (Docker-native)

```yaml
# docker-compose.yml (with Traefik)
version: "3.9"

services:
  traefik:
    image: traefik:v3.0
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.letsencrypt.acme.httpchallenge=true"
      - "--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web"
      - "--certificatesresolvers.letsencrypt.acme.email=admin@example.com"
      - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - traefik-certs:/letsencrypt

  revenge:
    image: ghcr.io/lusoris/revenge:latest
    labels:
      - "traefik.enable=true"
      # HTTP to HTTPS redirect
      - "traefik.http.routers.revenge-http.rule=Host(`revenge.example.com`)"
      - "traefik.http.routers.revenge-http.entrypoints=web"
      - "traefik.http.routers.revenge-http.middlewares=redirect-to-https"
      - "traefik.http.middlewares.redirect-to-https.redirectscheme.scheme=https"
      # HTTPS
      - "traefik.http.routers.revenge.rule=Host(`revenge.example.com`)"
      - "traefik.http.routers.revenge.entrypoints=websecure"
      - "traefik.http.routers.revenge.tls=true"
      - "traefik.http.routers.revenge.tls.certresolver=letsencrypt"
      # Middlewares
      - "traefik.http.routers.revenge.middlewares=revenge-chain"
      - "traefik.http.middlewares.revenge-chain.chain.middlewares=security-headers,rate-limit,compress"
      # Security headers
      - "traefik.http.middlewares.security-headers.headers.stsSeconds=63072000"
      - "traefik.http.middlewares.security-headers.headers.frameDeny=true"
      - "traefik.http.middlewares.security-headers.headers.contentTypeNosniff=true"
      # Rate limiting
      - "traefik.http.middlewares.rate-limit.ratelimit.average=100"
      - "traefik.http.middlewares.rate-limit.ratelimit.burst=50"
      # Compression
      - "traefik.http.middlewares.compress.compress=true"
      # Service
      - "traefik.http.services.revenge.loadbalancer.server.port=8096"
      - "traefik.http.services.revenge.loadbalancer.server.scheme=http"

volumes:
  traefik-certs:
```

---

## Server Configuration

### Revenge Backend Settings

```yaml
# configs/config.yaml
server:
  # Listen address (internal only behind proxy)
  host: "127.0.0.1"
  port: 8096

  # Trust proxy headers
  trusted_proxies:
    - "127.0.0.1"
    - "10.0.0.0/8"
    - "172.16.0.0/12"
    - "192.168.0.0/16"

  # Timeouts
  read_timeout: 30s
  write_timeout: 3600s    # Long for streaming
  idle_timeout: 120s

  # Request limits
  max_header_bytes: 1048576   # 1MB

  # Graceful shutdown
  shutdown_timeout: 30s

# For proper IP logging
logging:
  real_ip_header: "X-Real-IP"
  forwarded_for_header: "X-Forwarded-For"
```

### Headers Middleware

```go
// internal/api/middleware/proxy.go
package middleware

import (
    "net"
    "net/http"
    "strings"
)

type ProxyConfig struct {
    TrustedProxies    []string
    RealIPHeader      string
    ForwardedForHeader string
}

func TrustedProxy(cfg ProxyConfig) func(http.Handler) http.Handler {
    // Parse trusted CIDRs
    trustedNets := make([]*net.IPNet, 0, len(cfg.TrustedProxies))
    for _, cidr := range cfg.TrustedProxies {
        if !strings.Contains(cidr, "/") {
            cidr += "/32"
        }
        _, ipnet, err := net.ParseCIDR(cidr)
        if err == nil {
            trustedNets = append(trustedNets, ipnet)
        }
    }

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get remote IP
            remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
            ip := net.ParseIP(remoteIP)

            // Check if from trusted proxy
            trusted := false
            for _, ipnet := range trustedNets {
                if ipnet.Contains(ip) {
                    trusted = true
                    break
                }
            }

            if trusted {
                // Use real IP header
                if realIP := r.Header.Get(cfg.RealIPHeader); realIP != "" {
                    r.RemoteAddr = realIP + ":0"
                } else if xff := r.Header.Get(cfg.ForwardedForHeader); xff != "" {
                    // X-Forwarded-For: client, proxy1, proxy2
                    ips := strings.Split(xff, ",")
                    if len(ips) > 0 {
                        r.RemoteAddr = strings.TrimSpace(ips[0]) + ":0"
                    }
                }
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Streaming Optimization

### Byte Range Support

```go
// Ensure proper Range header handling
func StreamHandler(w http.ResponseWriter, r *http.Request) {
    // File info
    fileSize := getFileSize()

    // Parse Range header
    rangeHeader := r.Header.Get("Range")
    if rangeHeader == "" {
        // Full file
        w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
        w.Header().Set("Accept-Ranges", "bytes")
        w.WriteHeader(http.StatusOK)
        // Stream full file
        return
    }

    // Parse range: "bytes=0-1023"
    ranges, err := parseRange(rangeHeader, fileSize)
    if err != nil {
        http.Error(w, "Invalid Range", http.StatusRequestedRangeNotSatisfiable)
        return
    }

    if len(ranges) == 1 {
        // Single range
        ra := ranges[0]
        w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", ra.Start, ra.End, fileSize))
        w.Header().Set("Content-Length", fmt.Sprintf("%d", ra.End-ra.Start+1))
        w.Header().Set("Accept-Ranges", "bytes")
        w.WriteHeader(http.StatusPartialContent)
        // Stream range
    }
}
```

### Nginx Slice Module (for large files)

```nginx
# For very large files, use nginx slice module
location ~* ^/api/v1/stream/video/ {
    slice 1m;
    proxy_cache video_cache;
    proxy_cache_key $uri$is_args$args$slice_range;
    proxy_set_header Range $slice_range;
    proxy_cache_valid 200 206 1h;
    proxy_pass http://revenge_backend;
}
```

---

## Caching Strategy

### Nginx Proxy Cache

```nginx
# Define cache zone
proxy_cache_path /var/cache/nginx/revenge levels=1:2 keys_zone=revenge_cache:100m max_size=10g inactive=7d use_temp_path=off;

# Use in location
location /api/v1/images/ {
    proxy_cache revenge_cache;
    proxy_cache_valid 200 7d;
    proxy_cache_valid 404 1m;
    proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
    proxy_cache_lock on;

    add_header X-Cache-Status $upstream_cache_status;

    proxy_pass http://revenge_backend;
}
```

### CDN Headers

```go
// Set proper cache headers for CDN compatibility
func ImageHandler(w http.ResponseWriter, r *http.Request) {
    // ETag for cache validation
    etag := calculateETag(imageData)
    w.Header().Set("ETag", etag)

    // Check If-None-Match
    if match := r.Header.Get("If-None-Match"); match == etag {
        w.WriteHeader(http.StatusNotModified)
        return
    }

    // Cache control
    w.Header().Set("Cache-Control", "public, max-age=604800, stale-while-revalidate=86400")
    w.Header().Set("Vary", "Accept-Encoding")

    // Content type
    w.Header().Set("Content-Type", contentType)
    w.Write(imageData)
}
```

---

## Load Balancing

### Multiple Backend Instances

```nginx
upstream revenge_backend {
    least_conn;  # Least connections algorithm

    server 127.0.0.1:8096 weight=5;
    server 127.0.0.1:8097 weight=5;
    server 127.0.0.1:8098 weight=5 backup;

    keepalive 64;
}
```

### Docker Swarm / Kubernetes

```yaml
# docker-compose.yml for swarm
version: "3.9"

services:
  revenge:
    image: ghcr.io/lusoris/revenge:latest
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
        failure_action: rollback
      rollback_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      resources:
        limits:
          cpus: '2'
          memory: 4G
        reservations:
          cpus: '0.5'
          memory: 512M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8096/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

---

## Security Best Practices

### Fail2ban Configuration

```ini
# /etc/fail2ban/jail.d/revenge.conf
[revenge-auth]
enabled = true
port = http,https
filter = revenge-auth
logpath = /var/log/nginx/revenge_access.log
maxretry = 5
findtime = 600
bantime = 3600

# /etc/fail2ban/filter.d/revenge-auth.conf
[Definition]
failregex = ^<HOST> .* "POST /api/v1/auth/login.*" (401|403)
            ^<HOST> .* "POST /api/v1/auth/register.*" (400|429)
ignoreregex =
```

### WAF Rules (ModSecurity)

```nginx
# /etc/nginx/modsec/revenge.conf
SecRule REQUEST_URI "@beginsWith /api/v1/auth" \
    "id:1001,phase:1,pass,nolog,ctl:ruleEngine=On"

SecRule REQUEST_METHOD "!@pm GET HEAD OPTIONS" \
    "id:1002,phase:1,deny,status:405,msg:'Method not allowed'"

# SQL injection protection
SecRule ARGS "@detectSQLi" \
    "id:1003,phase:2,deny,status:403,msg:'SQL Injection detected'"

# XSS protection
SecRule ARGS "@detectXSS" \
    "id:1004,phase:2,deny,status:403,msg:'XSS detected'"
```

### IP Allowlist for Admin

```nginx
# Restrict admin endpoints
location /api/v1/admin/ {
    allow 10.0.0.0/8;
    allow 192.168.1.0/24;
    deny all;

    proxy_pass http://revenge_backend;
}
```

---

## Monitoring & Observability

### Prometheus Metrics Endpoint

```nginx
# Expose metrics securely
location /metrics {
    allow 127.0.0.1;
    allow 10.0.0.0/8;
    deny all;

    proxy_pass http://revenge_backend;
}
```

### Nginx Status

```nginx
# Nginx status for monitoring
location /nginx_status {
    stub_status on;
    allow 127.0.0.1;
    deny all;
}
```

### Log Format for Analysis

```nginx
# JSON log format
log_format json_combined escape=json '{'
    '"time_local":"$time_local",'
    '"remote_addr":"$remote_addr",'
    '"remote_user":"$remote_user",'
    '"request":"$request",'
    '"status": "$status",'
    '"body_bytes_sent":"$body_bytes_sent",'
    '"request_time":"$request_time",'
    '"http_referrer":"$http_referer",'
    '"http_user_agent":"$http_user_agent",'
    '"http_x_forwarded_for":"$http_x_forwarded_for",'
    '"upstream_response_time":"$upstream_response_time",'
    '"upstream_cache_status":"$upstream_cache_status"'
'}';

access_log /var/log/nginx/revenge_access.json json_combined;
```

---

## Docker Production Setup

### Optimized Dockerfile

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /revenge ./cmd/revenge

# Runtime stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -g '' revenge

WORKDIR /app
COPY --from=builder /revenge .
COPY configs/defaults.yaml /app/configs/

USER revenge
EXPOSE 8096

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8096/api/v1/health || exit 1

ENTRYPOINT ["./revenge"]
```

### Docker Compose Production

```yaml
# docker-compose.prod.yml
version: "3.9"

services:
  revenge:
    image: ghcr.io/lusoris/revenge:latest
    restart: unless-stopped
    environment:
      - REVENGE_DATABASE_URL=postgres://revenge:${DB_PASSWORD}@postgres:5432/revenge?sslmode=require
      - REVENGE_CACHE_URL=redis://dragonfly:6379
      - REVENGE_LOG_LEVEL=info
      - REVENGE_LOG_FORMAT=json
    volumes:
      - ./config:/app/config:ro
      - media:/media:ro
      - cache:/app/cache
    networks:
      - internal
      - proxy
    depends_on:
      postgres:
        condition: service_healthy
      dragonfly:
        condition: service_started
    deploy:
      resources:
        limits:
          cpus: '4'
          memory: 8G

  postgres:
    image: postgres:18-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: revenge
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: revenge
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - internal
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U revenge"]
      interval: 10s
      timeout: 5s
      retries: 5
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly
    restart: unless-stopped
    command: --maxmemory 2gb --proactor_threads 2
    volumes:
      - dragonfly_data:/data
    networks:
      - internal
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G

networks:
  internal:
    driver: bridge
  proxy:
    external: true

volumes:
  postgres_data:
  dragonfly_data:
  media:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /mnt/media
  cache:
```

---

## Configuration Summary

### Environment Variables

```bash
# Server
REVENGE_SERVER_HOST=127.0.0.1
REVENGE_SERVER_PORT=8096

# Database
REVENGE_DATABASE_URL=postgres://user:pass@localhost/revenge

# Cache
REVENGE_CACHE_URL=redis://localhost:6379

# Logging
REVENGE_LOG_LEVEL=info
REVENGE_LOG_FORMAT=json

# Security
REVENGE_TRUSTED_PROXIES=127.0.0.1,10.0.0.0/8
REVENGE_REAL_IP_HEADER=X-Real-IP

# Secrets (use Docker secrets in production)
REVENGE_JWT_SECRET=file:/run/secrets/jwt_secret
```


---

## Quick Reference

| Aspect | Recommendation |
|--------|----------------|
| TLS | TLS 1.2+ only, strong ciphers |
| HTTP | HTTP/2 enabled |
| Compression | gzip for text, skip for media |
| Timeouts | 30s connect, 3600s streaming |
| Rate Limiting | 10 req/s general, 5 req/min auth |
| Caching | Images 7d, API none, HLS 1m |
| Headers | HSTS, X-Frame-Options, CSP |
| Monitoring | Health checks, metrics endpoint |
| Logs | JSON format for analysis |
