# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
# FFmpeg development libraries are required for go-astiav (CGO bindings)
RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    pkgconfig \
    ffmpeg-dev \
    libavcodec-dev \
    libavformat-dev \
    libavutil-dev \
    libswscale-dev \
    libswresample-dev

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Enable Go 1.25 experimental features
ENV GOEXPERIMENT=greenteagc,jsonv2

# Build binary with stripped symbols
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -w -s" \
    -a -installsuffix cgo \
    -o revenge \
    ./cmd/revenge

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
# FFmpeg is required for:
# - HLS audio streaming (gohlslib + FFmpeg)
# - Audio transcoding (AAC, MP3, FLAC, Opus)
# - Metadata extraction & thumbnails (go-astiav)
# - Light video processing
# See: docs/dev/design/technical/AUDIO_STREAMING.md
# See: docs/dev/design/00_SOURCE_OF_TRUTH.md (go-astiav)
# FFmpeg libraries (libav*) needed at runtime for go-astiav CGO bindings
# postgresql-client for pg_isready in entrypoint script
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    ffmpeg-libs \
    postgresql-client \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 revenge && \
    adduser -D -u 1000 -G revenge revenge

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/revenge .

# Copy config files
COPY config/casbin_model.conf ./config/casbin_model.conf

# Copy entrypoint script
COPY scripts/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Create directories
RUN mkdir -p /data /config /cache /media && \
    chown -R revenge:revenge /app /data /config /cache

USER revenge

EXPOSE 8096

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8096/health/live || exit 1

VOLUME ["/data", "/config", "/media"]

# Environment variables
ENV REVENGE_DATA_DIR=/data \
    REVENGE_CONFIG_DIR=/config \
    REVENGE_CACHE_DIR=/cache

ENTRYPOINT ["/docker-entrypoint.sh"]
