.PHONY: help build run test lint clean docker-build docker-run migrate

# Variables
BINARY_NAME=revenge
DOCKER_IMAGE=revenge/revenge
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Database configuration (override with environment variables)
DB_HOST?=localhost
DB_PORT?=5432
DB_USER?=revenge
DB_PASSWORD?=revenge_dev_pass
DB_NAME?=revenge
DB_SSLMODE?=disable
DATABASE_URL?=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	go mod tidy

build: ## Build the binary
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/revenge

build-all: ## Build for all platforms
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd/revenge
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 ./cmd/revenge
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 ./cmd/revenge
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 ./cmd/revenge
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe ./cmd/revenge

run: ## Run the application
	@echo "Running ${BINARY_NAME}..."
	go run ./cmd/revenge

dev: ## Run with hot reload (requires air)
	@echo "Starting development server with hot reload..."
	air

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/integration/...

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run --timeout=5m

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ${DOCKER_IMAGE}:${VERSION} -t ${DOCKER_IMAGE}:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8096:8096 -v revenge-data:/data ${DOCKER_IMAGE}:latest

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services..."
	docker-compose down

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down (one step)
	@echo "Running migrations down..."
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-down-all: ## Run all database migrations down
	@echo "Running all migrations down..."
	migrate -path migrations -database "$(DATABASE_URL)" down -all

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@echo "Forcing migration version: ${VERSION}..."
	migrate -path migrations -database "$(DATABASE_URL)" force ${VERSION}

migrate-version: ## Show current migration version
	@echo "Current migration version:"
	migrate -path migrations -database "$(DATABASE_URL)" version

migrate-create: ## Create a new migration (usage: make migrate-create NAME=create_users_table)
	@echo "Creating migration: ${NAME}..."
	migrate create -ext sql -dir migrations -seq ${NAME}

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

generate: ## Run go generate
	@echo "Running go generate..."
	go generate ./...

sqlc: ## Generate sqlc code
	@echo "Generating sqlc code..."
	sqlc generate

all: clean deps lint test build ## Run all checks and build

.DEFAULT_GOAL := help
