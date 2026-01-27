# revenge - Setup Guide

## Quick Start Options

### Option 1: Docker Compose (Recommended)

The recommended way to run revenge with all required services:

```yaml
# docker-compose.yml
version: '3.8'

services:
  revenge:
    image: revenge/revenge:latest
    container_name: revenge
    ports:
      - "8096:8096"
    environment:
      - REVENGE_DB_HOST=postgres
      - REVENGE_DB_PORT=5432
      - REVENGE_DB_USER=revenge
      - REVENGE_DB_PASSWORD=changeme
      - REVENGE_DB_NAME=revenge
      - REVENGE_CACHE_URL=redis://dragonfly:6379
      - REVENGE_TYPESENSE_URL=http://typesense:8108
      - REVENGE_TYPESENSE_API_KEY=xyz
    volumes:
      - /path/to/media:/media
      - revenge-config:/config
    depends_on:
      - postgres
      - dragonfly
      - typesense
    restart: unless-stopped

  postgres:
    image: postgres:18-alpine
    container_name: revenge-postgres
    environment:
      - POSTGRES_DB=revenge
      - POSTGRES_USER=revenge
      - POSTGRES_PASSWORD=changeme
    volumes:
      - postgres-data:/var/lib/postgresql/data
    restart: unless-stopped

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly
    container_name: revenge-dragonfly
    restart: unless-stopped

  typesense:
    image: typesense/typesense:27.1
    container_name: revenge-typesense
    environment:
      - TYPESENSE_DATA_DIR=/data
      - TYPESENSE_API_KEY=xyz
    volumes:
      - typesense-data:/data
    restart: unless-stopped

volumes:
  revenge-config:
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
wget https://github.com/lusoris/revenge/releases/latest/revenge
chmod +x revenge

# Set environment variables
export REVENGE_DB_HOST=localhost
export REVENGE_DB_PORT=5432
export REVENGE_DB_USER=revenge
export REVENGE_DB_PASSWORD=changeme
export REVENGE_DB_NAME=revenge
export REVENGE_CACHE_URL=localhost:6379
export REVENGE_TYPESENSE_URL=http://localhost:8108
export REVENGE_TYPESENSE_API_KEY=xyz

# Run
./revenge

# Opens at http://localhost:8096
```

## Environment Variables

```bash
# Database (PostgreSQL - Required)
REVENGE_DB_HOST=localhost
REVENGE_DB_PORT=5432
REVENGE_DB_USER=revenge
REVENGE_DB_PASSWORD=
REVENGE_DB_NAME=revenge
REVENGE_DB_SSLMODE=disable      # disable, require, verify-ca, verify-full

# Cache (Dragonfly/Redis - Required)
REVENGE_CACHE_URL=localhost:6379
REVENGE_CACHE_PASSWORD=

# Search (Typesense - Required)
REVENGE_TYPESENSE_URL=http://localhost:8108
REVENGE_TYPESENSE_API_KEY=

# Server
REVENGE_HOST=0.0.0.0
REVENGE_PORT=8096
REVENGE_BASE_URL=/              # For reverse proxy

# FFmpeg
REVENGE_FFMPEG_PATH=ffmpeg      # Auto-detected if not set
REVENGE_FFMPEG_HWACCEL=auto     # auto, vaapi, nvenc, qsv, amf, videotoolbox, none

# Logging
REVENGE_LOG_LEVEL=info          # debug, info, warn, error
REVENGE_LOG_FORMAT=json         # json or text
```

## Building from Source

### Prerequisites

- Go 1.24 or higher
- FFmpeg (preferably revenge-ffmpeg)
- Git

### Clone and Build

```bash
# Clone repository
git clone https://github.com/your-org/revenge.git
cd revenge

# Download dependencies
go mod download

# Build
go build -o revenge ./cmd/revenge

# Run
./revenge
```

### Development Mode

```bash
# Run with hot reload (using air)
go install github.com/cosmtrek/air@latest
air

# Or directly
go run ./cmd/revenge
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
  user: revenge
  password: changeme
  name: revenge
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
    server_name revenge.example.com;

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
    server_name revenge.example.com;

    ssl_certificate /etc/letsencrypt/live/revenge.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/revenge.example.com/privkey.pem;

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
revenge.example.com {
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

### Import from Existing Revenge

If you have an existing Revenge installation:

1. First-time setup wizard will detect existing Revenge data
2. Select "Import from Revenge"
3. Choose your Revenge data directory
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
  --name revenge \
  --device /dev/dri/renderD128 \
  --device /dev/dri/card0 \
  -p 8096:8096 \
  -v /path/to/media:/media \
  revenge/revenge:latest
```

### NVIDIA (NVENC)

```bash
# Install nvidia-docker2
sudo apt install nvidia-docker2

# Run with NVIDIA runtime
docker run -d \
  --name revenge \
  --gpus all \
  -p 8096:8096 \
  -v /path/to/media:/media \
  revenge/revenge:latest
```

### macOS (VideoToolbox)

Works automatically on macOS, no special configuration needed.

## Systemd Service (Linux)

```ini
# /etc/systemd/system/revenge.service
[Unit]
Description=revenge Media Server
After=network.target

[Service]
Type=simple
User=revenge
Group=revenge
ExecStart=/usr/local/bin/revenge --config /etc/revenge/config.yaml
Restart=on-failure
RestartSec=5s

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/revenge

[Install]
WantedBy=multi-user.target
```

```bash
# Create user
sudo useradd -r -s /bin/false revenge

# Create directories
sudo mkdir -p /var/lib/revenge
sudo mkdir -p /etc/revenge
sudo chown revenge:revenge /var/lib/revenge

# Copy binary
sudo cp revenge /usr/local/bin/
sudo chmod +x /usr/local/bin/revenge

# Copy config
sudo cp config.yaml /etc/revenge/

# Enable and start
sudo systemctl enable revenge
sudo systemctl start revenge

# Check status
sudo systemctl status revenge
```

## Troubleshooting

### Check Logs

```bash
# Docker
docker logs revenge

# Docker Compose
docker-compose logs -f revenge

# Native/Systemd
journalctl -u revenge -f

# Or check log file
tail -f ~/.revenge/logs/revenge.log
```

### Common Issues

**Port already in use:**
```bash
# Check what's using port 8096
sudo lsof -i :8096

# Change port
REVENGE_PORT=8097 ./revenge
```

**Permission denied (media files):**
```bash
# Fix permissions
sudo chown -R revenge:revenge /path/to/media
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
docker-compose exec postgres psql -U revenge -d revenge -c "SELECT 1"
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
docker pull revenge/revenge:latest

# Stop and remove old container
docker stop revenge
docker rm revenge

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
cp -r ~/.revenge ~/.revenge.backup

# Download new version
wget https://github.com/your-org/revenge/releases/latest/revenge
chmod +x revenge

# Replace old binary
sudo cp revenge /usr/local/bin/

# Restart service
sudo systemctl restart revenge
```

## Backup & Restore

### PostgreSQL

```bash
# Backup
docker-compose exec postgres pg_dump -U revenge revenge > backup.sql

# Restore
docker-compose exec -T postgres psql -U revenge revenge < backup.sql
```

### Full Backup (Docker Volumes)

```bash
# Backup all volumes
docker run --rm \
  -v revenge-config:/config \
  -v postgres-data:/postgres \
  -v typesense-data:/typesense \
  -v $(pwd):/backup \
  alpine tar czf /backup/revenge-backup.tar.gz /config /postgres /typesense

# Restore
docker run --rm \
  -v revenge-config:/config \
  -v postgres-data:/postgres \
  -v typesense-data:/typesense \
  -v $(pwd):/backup \
  alpine tar xzf /backup/revenge-backup.tar.gz -C /
```

## Support

- Documentation: [docs/](../docs/)
- GitHub Issues: <https://github.com/lusoris/revenge/issues>
- Discussions: <https://github.com/lusoris/revenge/discussions>
