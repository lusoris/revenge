#!/bin/bash

# Development helper script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

function warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

function error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

function check_requirements() {
    info "Checking requirements..."

    if ! command -v go &> /dev/null; then
        error "Go is not installed. Please install Go 1.24 or later."
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    info "Go version: $GO_VERSION"

    if ! command -v docker &> /dev/null; then
        warn "Docker is not installed. Docker is optional but recommended."
    else
        info "Docker is installed"
    fi
}

function install_tools() {
    info "Installing development tools..."

    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/cosmtrek/air@latest
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

    info "Development tools installed"
}

function setup_db() {
    info "Setting up database..."

    info "Checking PostgreSQL connection..."

    # Check if PostgreSQL is reachable
    if PGPASSWORD=password psql -h localhost -U revenge -d revenge -c "SELECT 1" &>/dev/null; then
        info "PostgreSQL connection successful"
    else
        warn "PostgreSQL not reachable. Please ensure PostgreSQL is running."
        warn "Start with: docker-compose -f docker-compose.dev.yml up -d postgres"
    fi

    info "Database ready"
}

function run_tests() {
    info "Running tests..."
    go test -v -race -coverprofile=coverage.out ./...
}

function run_lint() {
    info "Running linter..."
    golangci-lint run --timeout=5m
}

function dev_server() {
    info "Starting development server with hot reload..."
    air
}

function build_binary() {
    info "Building binary..."
    go build -o bin/revenge ./cmd/revenge
    info "Binary built: bin/revenge"
}

# Main menu
case "${1:-}" in
    check)
        check_requirements
        ;;
    install-tools)
        install_tools
        ;;
    setup)
        check_requirements
        install_tools
        setup_db
        info "Setup complete! Run './scripts/dev.sh dev' to start development server"
        ;;
    test)
        run_tests
        ;;
    lint)
        run_lint
        ;;
    dev)
        dev_server
        ;;
    build)
        build_binary
        ;;
    *)
        echo "Usage: $0 {check|install-tools|setup|test|lint|dev|build}"
        echo ""
        echo "Commands:"
        echo "  check         - Check requirements"
        echo "  install-tools - Install development tools"
        echo "  setup         - Full development setup"
        echo "  test          - Run tests"
        echo "  lint          - Run linter"
        echo "  dev           - Start development server with hot reload"
        echo "  build         - Build binary"
        exit 1
        ;;
esac
