#!/usr/bin/env fish
# Universal setup script for Jellyfin Go (Fish shell)

set SCRIPT_DIR (dirname (status --current-filename))
set PROJECT_ROOT (dirname $SCRIPT_DIR)

# Colors
set -g RED '\033[0;31m'
set -g GREEN '\033[0;32m'
set -g YELLOW '\033[1;33m'
set -g BLUE '\033[0;34m'
set -g NC '\033[0m'

function info
    echo -e "$GREEN[INFO]$NC $argv"
end

function warn
    echo -e "$YELLOW[WARN]$NC $argv"
end

function error
    echo -e "$RED[ERROR]$NC $argv"
end

function step
    echo -e "$BLUE[STEP]$NC $argv"
end

# Detect OS and package manager
function detect_system
    if test (uname) = "Darwin"
        set -g OS "macos"
        set -g PKG_MANAGER "brew"
    else if test -f /etc/arch-release
        set -g OS "linux"
        set -g PKG_MANAGER "pacman"
    else if test -f /etc/fedora-release; or test -f /etc/redhat-release
        set -g OS "linux"
        set -g PKG_MANAGER "yum"
    else if test -f /etc/debian_version
        set -g OS "linux"
        set -g PKG_MANAGER "apt"
    else
        set -g OS "linux"
        set -g PKG_MANAGER "unknown"
    end
    
    info "Detected: $OS with $PKG_MANAGER"
end

# Check if command exists
function has_command
    command -v $argv[1] > /dev/null 2>&1
end

# Install Go
function install_go
    if has_command go
        set GO_VERSION (go version | awk '{print $3}' | sed 's/go//')
        info "Go $GO_VERSION is already installed"
        return 0
    end
    
    step "Installing Go 1.24..."
    
    switch $PKG_MANAGER
        case brew
            brew install go@1.24; or brew install go
        case apt
            sudo apt update
            sudo apt install -y wget
            wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
            sudo rm -rf /usr/local/go
            sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
            rm go1.24.linux-amd64.tar.gz
            fish_add_path /usr/local/go/bin
        case yum
            sudo yum install -y wget
            wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
            sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
            rm go1.24.linux-amd64.tar.gz
            fish_add_path /usr/local/go/bin
        case pacman
            sudo pacman -Sy --noconfirm go
        case '*'
            error "Cannot auto-install Go. Please install Go 1.24 manually from https://go.dev/dl/"
            exit 1
    end
    
    info "Go installed successfully"
end

# Install FFmpeg
function install_ffmpeg
    if has_command ffmpeg
        info "FFmpeg is already installed"
        return 0
    end
    
    step "Installing FFmpeg..."
    
    switch $PKG_MANAGER
        case brew
            brew install ffmpeg
        case apt
            sudo apt update
            sudo apt install -y ffmpeg
        case yum
            sudo yum install -y epel-release
            sudo yum install -y ffmpeg
        case pacman
            sudo pacman -Sy --noconfirm ffmpeg
        case '*'
            warn "Cannot auto-install FFmpeg. Please install manually."
            return 1
    end
    
    info "FFmpeg installed successfully"
end

# Install Go tools
function install_go_tools
    step "Installing Go development tools..."
    
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/cosmtrek/air@latest
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    
    info "Go tools installed successfully"
end

# Setup project
function setup_project
    cd $PROJECT_ROOT
    
    step "Downloading Go dependencies..."
    go mod download
    go mod verify
    
    step "Installing Git hooks..."
    if test -f "$SCRIPT_DIR/install-hooks.sh"
        bash "$SCRIPT_DIR/install-hooks.sh"
    end
    
    step "Building project..."
    go build -o bin/jellyfin-go ./cmd/jellyfin
    
    info "Project setup complete!"
end

# Main
echo ""
echo "╔════════════════════════════════════════╗"
echo "║   Jellyfin Go - Setup Assistant       ║"
echo "╚════════════════════════════════════════╝"
echo ""

detect_system
echo ""

install_go
install_ffmpeg
install_go_tools
setup_project

echo ""
echo "╔════════════════════════════════════════╗"
echo "║   ✅ Setup Complete!                   ║"
echo "╚════════════════════════════════════════╝"
echo ""
info "Next steps:"
echo "  1. Run the application: ./bin/jellyfin-go"
echo "  2. Or use development mode: make dev"
echo "  3. Read the docs: cat README.md"
echo ""
