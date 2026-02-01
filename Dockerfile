# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files (go.sum may not exist if no dependencies)
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Enable Go 1.25 experimental features:
# - greenteagc: New garbage collector with 10-40% memory reduction
# - jsonv2: Faster JSON encoding/decoding
ENV GOEXPERIMENT=greenteagc,jsonv2

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -w -s" \
    -a -installsuffix cgo \
    -o revenge \
    ./cmd/revenge

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 revenge && \
    adduser -D -u 1000 -G revenge revenge

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/revenge .

# Create directories
RUN mkdir -p /data /config /cache /media && \
    chown -R revenge:revenge /app /data /config /cache

# Switch to non-root user
USER revenge

# Expose port
EXPOSE 8096

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8096/health/live || exit 1

# Volume for data
VOLUME ["/data", "/config", "/media"]

# Set environment variables
ENV REVENGE_DATA_DIR=/data \
    REVENGE_CONFIG_DIR=/config \
    REVENGE_CACHE_DIR=/cache

# Run the application
ENTRYPOINT ["/app/revenge"]
