#!/usr/bin/env bash
# Universal setup script for Revenge
# Supports: Linux (apt, yum, pacman), macOS (brew), Windows (via WSL)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }
step() { echo -e "${BLUE}[STEP]${NC} $1"; }

# Detect OS and package manager
detect_system() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
        PKG_MANAGER="brew"
    elif [[ -f /etc/arch-release ]]; then
        OS="linux"
        PKG_MANAGER="pacman"
    elif [[ -f /etc/fedora-release ]] || [[ -f /etc/redhat-release ]]; then
        OS="linux"
        PKG_MANAGER="yum"
    elif [[ -f /etc/debian_version ]]; then
        OS="linux"
        PKG_MANAGER="apt"
    else
        OS="linux"
        PKG_MANAGER="unknown"
    fi

    info "Detected: $OS with $PKG_MANAGER"
}

# Check if command exists
has_command() {
    command -v "$1" &> /dev/null
}

# Install Go
install_go() {
    if has_command go; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        info "Go $GO_VERSION is already installed"
        return 0
    fi

    step "Installing Go 1.25..."

    case "$PKG_MANAGER" in
        brew)
            brew install go@1.25 || brew install go
            ;;
        apt)
            sudo apt update
            sudo apt install -y wget
            wget https://go.dev/dl/go1.25.linux-amd64.tar.gz
            sudo rm -rf /usr/local/go
            sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
            rm go1.25.linux-amd64.tar.gz
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
            export PATH=$PATH:/usr/local/go/bin
            ;;
        yum)
            sudo yum install -y wget
            wget https://go.dev/dl/go1.25.linux-amd64.tar.gz
            sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
            rm go1.25.linux-amd64.tar.gz
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
            export PATH=$PATH:/usr/local/go/bin
            ;;
        pacman)
            sudo pacman -Sy --noconfirm go
            ;;
        *)
            error "Cannot auto-install Go. Please install Go 1.25 manually from https://go.dev/dl/"
            exit 1
            ;;
    esac

    info "Go installed successfully"
}

# Install FFmpeg
install_ffmpeg() {
    if has_command ffmpeg; then
        info "FFmpeg is already installed"
        return 0
    fi

    step "Installing FFmpeg..."

    case "$PKG_MANAGER" in
        brew)
            brew install ffmpeg
            ;;
        apt)
            sudo apt update
            sudo apt install -y ffmpeg
            ;;
        yum)
            sudo yum install -y epel-release
            sudo yum install -y ffmpeg
            ;;
        pacman)
            sudo pacman -Sy --noconfirm ffmpeg
            ;;
        *)
            warn "Cannot auto-install FFmpeg. Please install manually."
            return 1
            ;;
    esac

    info "FFmpeg installed successfully"
}

# Install Docker (optional)
install_docker() {
    if has_command docker; then
        info "Docker is already installed"
        return 0
    fi

    read -p "Install Docker? (recommended but optional) [y/N]: " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        warn "Skipping Docker installation"
        return 0
    fi

    step "Installing Docker..."

    case "$PKG_MANAGER" in
        brew)
            warn "Please install Docker Desktop for Mac manually from https://docker.com"
            ;;
        apt)
            curl -fsSL https://get.docker.com -o get-docker.sh
            sudo sh get-docker.sh
            sudo usermod -aG docker $USER
            rm get-docker.sh
            info "Docker installed. Please log out and back in for group changes to take effect."
            ;;
        yum)
            sudo yum install -y yum-utils
            sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
            sudo yum install -y docker-ce docker-ce-cli containerd.io
            sudo systemctl start docker
            sudo systemctl enable docker
            sudo usermod -aG docker $USER
            ;;
        pacman)
            sudo pacman -Sy --noconfirm docker docker-compose
            sudo systemctl start docker
            sudo systemctl enable docker
            sudo usermod -aG docker $USER
            ;;
        *)
            warn "Cannot auto-install Docker. Please install manually."
            ;;
    esac
}

# Install Go development tools
install_go_tools() {
    step "Installing Go development tools..."

    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/cosmtrek/air@latest
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

    info "Go tools installed successfully"
}

# Setup project
setup_project() {
    cd "$PROJECT_ROOT"

    step "Downloading Go dependencies..."
    go mod download
    go mod verify

    step "Installing Git hooks..."
    if [[ -f "$SCRIPT_DIR/install-hooks.sh" ]]; then
        bash "$SCRIPT_DIR/install-hooks.sh"
    fi

    step "Building project..."
    go build -o bin/revenge ./cmd/revenge

    info "Project setup complete!"
}

# Main installation flow
main() {
    echo ""
    echo "╔════════════════════════════════════════╗"
    echo "║   Revenge - Setup Assistant       ║"
    echo "╚════════════════════════════════════════╝"
    echo ""

    detect_system
    echo ""

    install_go
    install_ffmpeg
    install_docker
    install_go_tools
    setup_project

    echo ""
    echo "╔════════════════════════════════════════╗"
    echo "║   ✅ Setup Complete!                   ║"
    echo "╚════════════════════════════════════════╝"
    echo ""
    info "Next steps:"
    echo "  1. Run the application: ./bin/revenge"
    echo "  2. Or use development mode: make dev"
    echo "  3. Read the docs: cat README.md"
    echo ""
}

main "$@"
