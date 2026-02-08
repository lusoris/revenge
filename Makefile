.PHONY: help build run test lint clean docker-build docker-run migrate

# Variables
BINARY_NAME=revenge
DOCKER_IMAGE=ghcr.io/lusoris/revenge
MIGRATIONS_DIR=internal/infra/database/migrations/shared
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Go 1.25 experimental features
export GOEXPERIMENT=greenteagc,jsonv2

# Database configuration (override with environment variables)
DB_HOST?=localhost
DB_PORT?=5432
DB_USER?=revenge
DB_PASSWORD?=revenge_dev_pass
DB_NAME?=revenge
DB_SSLMODE?=disable
DATABASE_URL?=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# =============================================================================
# Help
# =============================================================================

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# Build
# =============================================================================

build: ## Build the binary
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/revenge

build-linux: ## Build for Linux (Docker targets: amd64 + arm64)
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd/revenge
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 ./cmd/revenge

run: ## Run the application
	go run ./cmd/revenge

dev: ## Run with hot reload (requires air)
	air

# =============================================================================
# Testing - Local and CI use the same targets
# =============================================================================

test: ## Run unit tests (fast, no Docker needed)
	@echo "Running unit tests..."
	go test -race -coverprofile=coverage.out -covermode=atomic -count=1 ./...

test-short: ## Run unit tests in short mode (skip slow tests)
	@echo "Running short tests..."
	go test -short -count=1 ./...

test-integration: ## Run integration tests (requires Docker)
	@echo "Running integration tests..."
	go test -v -race -tags=integration -count=1 ./tests/integration/...

test-all: test test-integration ## Run all tests (unit + integration)

test-live: ## Run live smoke tests against running stack (requires make docker-local)
	@echo "Running live smoke tests against $(or $(REVENGE_TEST_URL),http://localhost:8096)..."
	go test -tags=live -v -count=1 ./tests/live/...

test-coverage: test ## Run tests and open coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# =============================================================================
# Docker - Build, scan, test the real image
# =============================================================================

docker-build: ## Build Docker image locally
	@echo "Building Docker image..."
	docker build -t ${DOCKER_IMAGE}:${VERSION} -t ${DOCKER_IMAGE}:dev -t revenge/revenge:dev .

docker-scan: docker-build ## Build and scan Docker image with Trivy
	@echo "Scanning Docker image with Trivy..."
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
		aquasec/trivy:latest image --severity CRITICAL,HIGH ${DOCKER_IMAGE}:dev

docker-test: docker-build ## Build image and run smoke test with full stack
	@echo "Starting full stack for smoke test..."
	docker compose -f docker-compose.dev.yml up -d --wait
	@echo "Running smoke tests against real image..."
	@sleep 5
	@curl -sf http://localhost:8096/healthz && echo "Health check: OK" || echo "Health check: FAILED"
	docker compose -f docker-compose.dev.yml down

docker-local: docker-build ## Build and run full local stack
	docker compose -f docker-compose.dev.yml up -d --wait
	@echo "Waiting for services to initialize..."
	@sleep 5
	@curl -sf http://localhost:8096/healthz && echo "Revenge is healthy!" || echo "Startup failed - check logs with: docker compose -f docker-compose.dev.yml logs revenge"

docker-up: ## Start dev services (postgres, dragonfly, typesense)
	docker compose -f docker-compose.dev.yml up -d --wait

docker-down: ## Stop dev services
	docker compose -f docker-compose.dev.yml down

# =============================================================================
# CI Pipeline - Runs the same as local but in order
# =============================================================================

ci: lint test docker-build docker-scan ## Full CI pipeline (lint + test + build + scan)

# =============================================================================
# Code Quality
# =============================================================================

lint: ## Run linters
	golangci-lint run --timeout=5m

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...

vuln: ## Run govulncheck
	govulncheck ./...

# =============================================================================
# Database Migrations
# =============================================================================

migrate-up: ## Run database migrations up
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down (one step)
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-down-all: ## Run all database migrations down
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down -all

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force ${VERSION}

migrate-version: ## Show current migration version
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

migrate-create: ## Create a new migration (usage: make migrate-create NAME=create_users_table)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq ${NAME}

# =============================================================================
# Code Generation
# =============================================================================

generate: ogen sqlc ## Run all code generation (ogen, sqlc, go generate)
	go generate ./...

ogen: ## Generate ogen code from OpenAPI spec
	go run github.com/ogen-go/ogen/cmd/ogen@v1.18.0 --target internal/api/ogen --package ogen --clean api/openapi/openapi.yaml

sqlc: ## Generate sqlc code
	sqlc generate

# =============================================================================
# Tools
# =============================================================================

install-tools: ## Install development tools
	@echo "Installing development tools..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.8.0
	go install github.com/air-verse/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/vektra/mockery/v3@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

# =============================================================================
# Cleanup
# =============================================================================

clean: ## Clean build artifacts
	rm -rf bin/ dist/
	rm -f coverage.out coverage.html
	go clean

deps: ## Download and verify dependencies
	go mod download
	go mod verify

tidy: ## Tidy go.mod
	go mod tidy

# =============================================================================
# Full Pipeline
# =============================================================================

all: clean deps lint test build ## Run all checks and build

.DEFAULT_GOAL := help
