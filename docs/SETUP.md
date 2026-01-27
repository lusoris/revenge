# Jellyfin Go - Setup Guide

## Quick Start Options

### Option 1: Docker Compose (Recommended)

The recommended way to run Jellyfin Go with all required services:

```yaml
# docker-compose.yml
version: '3.8'

services:
  jellyfin:
    image: jellyfin/jellyfin-go:latest
    container_name: jellyfin-go
    ports:
      - "8096:8096"
    environment:
      - JELLYFIN_DB_HOST=postgres
      - JELLYFIN_DB_PORT=5432
      - JELLYFIN_DB_USER=jellyfin
      - JELLYFIN_DB_PASSWORD=changeme
      - JELLYFIN_DB_NAME=jellyfin
      - JELLYFIN_CACHE_URL=redis://dragonfly:6379
      - JELLYFIN_TYPESENSE_URL=http://typesense:8108
      - JELLYFIN_TYPESENSE_API_KEY=xyz
    volumes:
      - /path/to/media:/media
      - jellyfin-config:/config
    depends_on:
      - postgres
      - dragonfly
      - typesense
    restart: unless-stopped

  postgres:
    image: postgres:18-alpine
    container_name: jellyfin-postgres
    environment:
      - POSTGRES_DB=jellyfin
      - POSTGRES_USER=jellyfin
      - POSTGRES_PASSWORD=changeme
    volumes:
      - postgres-data:/var/lib/postgresql/data
    restart: unless-stopped

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly
    container_name: jellyfin-dragonfly
    restart: unless-stopped

  typesense:
    image: typesense/typesense:27.1
    container_name: jellyfin-typesense
    environment:
      - TYPESENSE_DATA_DIR=/data
      - TYPESENSE_API_KEY=xyz
    volumes:
      - typesense-data:/data
    restart: unless-stopped

volumes:
  jellyfin-config:
  postgres-data:
  typesense-data:
```

```bash
docker-compose up -d
```

Access at `http://localhost:8096`

### Option 2: Native Binary

Requires PostgreSQL 18+, Dragonfly, and Typesense running locally or remotely.

```bash
# Download binary
wget https://github.com/lusoris/jellyfin-go/releases/latest/jellyfin-go
chmod +x jellyfin-go

# Set environment variables
export JELLYFIN_DB_HOST=localhost
export JELLYFIN_DB_PORT=5432
export JELLYFIN_DB_USER=jellyfin
export JELLYFIN_DB_PASSWORD=changeme
export JELLYFIN_DB_NAME=jellyfin
export JELLYFIN_CACHE_URL=localhost:6379
export JELLYFIN_TYPESENSE_URL=http://localhost:8108
export JELLYFIN_TYPESENSE_API_KEY=xyz

# Run
./jellyfin-go

# Opens at http://localhost:8096
```

## Environment Variables

```bash
# Database (PostgreSQL - Required)
JELLYFIN_DB_HOST=localhost
JELLYFIN_DB_PORT=5432
JELLYFIN_DB_USER=jellyfin
JELLYFIN_DB_PASSWORD=
JELLYFIN_DB_NAME=jellyfin
JELLYFIN_DB_SSLMODE=disable      # disable, require, verify-ca, verify-full

# Cache (Dragonfly/Redis - Required)
JELLYFIN_CACHE_URL=localhost:6379
JELLYFIN_CACHE_PASSWORD=

# Search (Typesense - Required)
JELLYFIN_TYPESENSE_URL=http://localhost:8108
JELLYFIN_TYPESENSE_API_KEY=

# Server
JELLYFIN_HOST=0.0.0.0
JELLYFIN_PORT=8096
JELLYFIN_BASE_URL=/              # For reverse proxy

# FFmpeg
JELLYFIN_FFMPEG_PATH=ffmpeg      # Auto-detected if not set
JELLYFIN_FFMPEG_HWACCEL=auto     # auto, vaapi, nvenc, qsv, amf, videotoolbox, none

# Logging
JELLYFIN_LOG_LEVEL=info          # debug, info, warn, error
JELLYFIN_LOG_FORMAT=json         # json or text
```

## Building from Source

### Prerequisites

- Go 1.24 or higher
- FFmpeg (preferably jellyfin-ffmpeg)
- Git

### Clone and Build

```bash
# Clone repository
git clone https://github.com/your-org/jellyfin-go.git
cd jellyfin-go

# Download dependencies
go mod download

# Build
go build -o jellyfin-go ./cmd/jellyfin

# Run
./jellyfin-go
```

### Development Mode

```bash
# Run with hot reload (using air)
go install github.com/cosmtrek/air@latest
air

# Or directly
go run ./cmd/jellyfin
```

## Configuration File

Create `config.yaml` in the same directory as the binary:

```yaml
server:
  host: 0.0.0.0
  port: 8096
  base_url: /
  read_timeout: 30s
  write_timeout: 30s

database:
  host: localhost
  port: 5432
  user: jellyfin
  password: changeme
  name: jellyfin
  sslmode: disable
  max_connections: 25
  min_connections: 5

cache:
  addr: localhost:6379
  password: ""
  db: 0

search:
  url: http://localhost:8108
  api_key: xyz

ffmpeg:
  path: ffmpeg  # Auto-detected if empty
  hwaccel: auto # auto, vaapi, nvenc, qsv, amf, videotoolbox, none
  max_concurrent_jobs: 5

logging:
  level: info
  format: json
```

## Reverse Proxy Setup

### NGINX

```nginx
server {
    listen 80;
    server_name jellyfin.example.com;

    # For Let's Encrypt
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$server_name$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name jellyfin.example.com;

    ssl_certificate /etc/letsencrypt/live/jellyfin.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/jellyfin.example.com/privkey.pem;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    location / {
        proxy_pass http://localhost:8096;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $server_name;

        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # Timeouts for streaming
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }

    # Streaming optimization
    location /api/stream/ {
        proxy_pass http://localhost:8096;
        proxy_buffering off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Caddy (Easier)

```caddyfile
jellyfin.example.com {
    reverse_proxy localhost:8096
}
```

## First-Time Setup Wizard

1. Open `http://localhost:8096` in your browser
2. Select language
3. Create admin account
4. Add media libraries (movies, TV shows, music, etc.)
5. Configure metadata providers (optional)
6. Done!

### Import from Existing Jellyfin

If you have an existing Jellyfin installation:

1. First-time setup wizard will detect existing Jellyfin data
2. Select "Import from Jellyfin"
3. Choose your Jellyfin data directory
4. Migration wizard will backup and convert your data
5. Progress shown in real-time

**Migration includes:**
- Users and passwords
- Libraries and media items
- Watch history and progress
- Playlists and collections
- Settings and preferences

## Hardware Acceleration

### Linux (VAAPI)

```bash
# Install drivers
sudo apt install intel-media-va-driver  # Intel
# or
sudo apt install mesa-va-drivers         # AMD

# Run Docker with GPU access
docker run -d \
  --name jellyfin-go \
  --device /dev/dri/renderD128 \
  --device /dev/dri/card0 \
  -p 8096:8096 \
  -v /path/to/media:/media \
  jellyfin/jellyfin-go:latest
```

### NVIDIA (NVENC)

```bash
# Install nvidia-docker2
sudo apt install nvidia-docker2

# Run with NVIDIA runtime
docker run -d \
  --name jellyfin-go \
  --gpus all \
  -p 8096:8096 \
  -v /path/to/media:/media \
  jellyfin/jellyfin-go:latest
```

### macOS (VideoToolbox)

Works automatically on macOS, no special configuration needed.

## Systemd Service (Linux)

```ini
# /etc/systemd/system/jellyfin-go.service
[Unit]
Description=Jellyfin Go Media Server
After=network.target

[Service]
Type=simple
User=jellyfin
Group=jellyfin
ExecStart=/usr/local/bin/jellyfin-go --config /etc/jellyfin-go/config.yaml
Restart=on-failure
RestartSec=5s

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/jellyfin-go

[Install]
WantedBy=multi-user.target
```

```bash
# Create user
sudo useradd -r -s /bin/false jellyfin

# Create directories
sudo mkdir -p /var/lib/jellyfin-go
sudo mkdir -p /etc/jellyfin-go
sudo chown jellyfin:jellyfin /var/lib/jellyfin-go

# Copy binary
sudo cp jellyfin-go /usr/local/bin/
sudo chmod +x /usr/local/bin/jellyfin-go

# Copy config
sudo cp config.yaml /etc/jellyfin-go/

# Enable and start
sudo systemctl enable jellyfin-go
sudo systemctl start jellyfin-go

# Check status
sudo systemctl status jellyfin-go
```

## Troubleshooting

### Check Logs

```bash
# Docker
docker logs jellyfin-go

# Docker Compose
docker-compose logs -f jellyfin

# Native/Systemd
journalctl -u jellyfin-go -f

# Or check log file
tail -f ~/.jellyfin-go/logs/jellyfin.log
```

### Common Issues

**Port already in use:**
```bash
# Check what's using port 8096
sudo lsof -i :8096

# Change port
JELLYFIN_PORT=8097 ./jellyfin-go
```

**Permission denied (media files):**
```bash
# Fix permissions
sudo chown -R jellyfin:jellyfin /path/to/media
sudo chmod -R 755 /path/to/media
```

**Hardware acceleration not working:**
```bash
# Check FFmpeg capabilities
ffmpeg -hwaccels

# Test VAAPI
ffmpeg -hwaccel vaapi -hwaccel_device /dev/dri/renderD128 -i input.mp4 output.mp4
```

**Database connection failed (PostgreSQL):**
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check credentials
docker-compose exec postgres psql -U jellyfin -d jellyfin -c "SELECT 1"
```

## Performance Tuning

### PostgreSQL Optimization

```bash
# Increase shared_buffers
echo "shared_buffers = 256MB" >> /var/lib/postgresql/data/postgresql.conf

# Increase work_mem
echo "work_mem = 16MB" >> /var/lib/postgresql/data/postgresql.conf

# Restart PostgreSQL
docker-compose restart postgres
```

### Dragonfly Memory

```bash
# Set max memory in docker-compose.yml
dragonfly:
  image: docker.dragonflydb.io/dragonflydb/dragonfly
  command: --maxmemory 1gb
```

## Upgrade

### Docker

```bash
# Pull latest image
docker pull jellyfin/jellyfin-go:latest

# Stop and remove old container
docker stop jellyfin-go
docker rm jellyfin-go

# Start new container (data preserved in volume)
docker run -d ...
```

### Docker Compose

```bash
docker-compose pull
docker-compose up -d
```

### Native Binary

```bash
# Backup data first
cp -r ~/.jellyfin-go ~/.jellyfin-go.backup

# Download new version
wget https://github.com/your-org/jellyfin-go/releases/latest/jellyfin-go
chmod +x jellyfin-go

# Replace old binary
sudo cp jellyfin-go /usr/local/bin/

# Restart service
sudo systemctl restart jellyfin-go
```

## Backup & Restore

### PostgreSQL

```bash
# Backup
docker-compose exec postgres pg_dump -U jellyfin jellyfin > backup.sql

# Restore
docker-compose exec -T postgres psql -U jellyfin jellyfin < backup.sql
```

### Full Backup (Docker Volumes)

```bash
# Backup all volumes
docker run --rm \
  -v jellyfin-config:/config \
  -v postgres-data:/postgres \
  -v typesense-data:/typesense \
  -v $(pwd):/backup \
  alpine tar czf /backup/jellyfin-backup.tar.gz /config /postgres /typesense

# Restore
docker run --rm \
  -v jellyfin-config:/config \
  -v postgres-data:/postgres \
  -v typesense-data:/typesense \
  -v $(pwd):/backup \
  alpine tar xzf /backup/jellyfin-backup.tar.gz -C /
```

## Support

- Documentation: [docs/](../docs/)
- GitHub Issues: <https://github.com/lusoris/jellyfin-go/issues>
- Discussions: <https://github.com/lusoris/jellyfin-go/discussions>
