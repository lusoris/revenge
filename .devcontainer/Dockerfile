FROM golang:1.25-alpine

# System deps (runtime + CGO build dependencies)
RUN apk add --no-cache \
	# Build tools
	git \
	make \
	gcc \
	musl-dev \
	pkgconfig \
	bash \
	curl \
	wget \
	jq \
	# Runtime
	ffmpeg \
	postgresql16-client \
	# CGO build deps: govips image processing
	vips-dev \
	# CGO build deps: go-astiav FFmpeg bindings
	ffmpeg-dev

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
