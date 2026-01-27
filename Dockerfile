# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -w -s" \
    -a -installsuffix cgo \
    -o jellyfin-go \
    ./cmd/jellyfin

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 jellyfin && \
    adduser -D -u 1000 -G jellyfin jellyfin

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/jellyfin-go .

# Create directories
RUN mkdir -p /data /config /cache /media && \
    chown -R jellyfin:jellyfin /app /data /config /cache

# Switch to non-root user
USER jellyfin

# Expose port
EXPOSE 8096

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8096/health/live || exit 1

# Volume for data
VOLUME ["/data", "/config", "/media"]

# Set environment variables
ENV JELLYFIN_DATA_DIR=/data \
    JELLYFIN_CONFIG_DIR=/config \
    JELLYFIN_CACHE_DIR=/cache

# Run the application
ENTRYPOINT ["/app/jellyfin-go"]
