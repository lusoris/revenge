FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm

# Install additional tools
RUN apt-get update && apt-get install -y \
	ffmpeg \
	sqlite3 \
	postgresql-client \
	curl \
	wget \
	jq \
	&& rm -rf /var/lib/apt/lists/*

# Install Go tools
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
	&& go install github.com/cosmtrek/air@latest \
	&& go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
	&& go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
	&& go install golang.org/x/vuln/cmd/govulncheck@latest \
	&& go install github.com/go-delve/delve/cmd/dlv@latest

# Set workspace
WORKDIR /workspace
