# Custom Typesense image with healthcheck support
FROM typesense/typesense:0.25.2

# Install curl for healthchecks using apt (Debian-based)
USER root
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Keep original entrypoint and command (runs as root by default)
