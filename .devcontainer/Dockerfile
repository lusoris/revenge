FROM mcr.microsoft.com/devcontainers/go:1-1.25-bookworm

# System deps
RUN apt-get update && apt-get install -y --no-install-recommends \
	ffmpeg \
	postgresql-client \
	curl \
	wget \
	jq \
	&& rm -rf /var/lib/apt/lists/*

# Go dev tools (versions match go.mod / Makefile)
ENV GOEXPERIMENT=greenteagc,jsonv2
RUN go install golang.org/x/tools/gopls@latest \
	&& go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
	&& go install github.com/air-verse/air@latest \
	&& go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
	&& go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
	&& go install github.com/vektra/mockery/v3@latest \
	&& go install github.com/go-delve/delve/cmd/dlv@latest \
	&& go install golang.org/x/vuln/cmd/govulncheck@latest

WORKDIR /workspace
