#!/usr/bin/env pwsh
# Universal setup script for Jellyfin Go on Windows
# Supports: winget, choco, scoop

param(
    [switch]$SkipDocker
)

$ErrorActionPreference = "Stop"

$ProjectRoot = Split-Path -Parent $PSScriptRoot

function Write-Step {
    param([string]$Message)
    Write-Host "[STEP] $Message" -ForegroundColor Blue
}

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Detect package manager
function Get-PackageManager {
    if (Get-Command winget -ErrorAction SilentlyContinue) {
        return "winget"
    }
    elseif (Get-Command choco -ErrorAction SilentlyContinue) {
        return "choco"
    }
    elseif (Get-Command scoop -ErrorAction SilentlyContinue) {
        return "scoop"
    }
    else {
        return "none"
    }
}

# Install Go
function Install-Go {
    if (Get-Command go -ErrorAction SilentlyContinue) {
        $goVersion = (go version) -replace 'go version go', '' -replace ' .*', ''
        Write-Info "Go $goVersion is already installed"
        return
    }

    Write-Step "Installing Go 1.24..."

    $pkgManager = Get-PackageManager

    switch ($pkgManager) {
        "winget" {
            winget install GoLang.Go.1.24
        }
        "choco" {
            choco install golang -y
        }
        "scoop" {
            scoop install go
        }
        default {
            Write-Error "No package manager found. Please install winget, choco, or scoop first."
            Write-Info "Or download Go manually from https://go.dev/dl/"
            exit 1
        }
    }

    Write-Info "Go installed successfully. You may need to restart your terminal."
}

# Install FFmpeg
function Install-FFmpeg {
    if (Get-Command ffmpeg -ErrorAction SilentlyContinue) {
        Write-Info "FFmpeg is already installed"
        return
    }

    Write-Step "Installing FFmpeg..."

    $pkgManager = Get-PackageManager

    switch ($pkgManager) {
        "winget" {
            winget install Gyan.FFmpeg
        }
        "choco" {
            choco install ffmpeg -y
        }
        "scoop" {
            scoop install ffmpeg
        }
        default {
            Write-Warn "Cannot auto-install FFmpeg. Please install manually."
        }
    }
}

# Install Docker
function Install-Docker {
    if (Get-Command docker -ErrorAction SilentlyContinue) {
        Write-Info "Docker is already installed"
        return
    }

    if ($SkipDocker) {
        Write-Info "Skipping Docker installation (--SkipDocker flag)"
        return
    }

    $install = Read-Host "Install Docker Desktop? (recommended but optional) [y/N]"
    if ($install -notmatch '^[Yy]') {
        Write-Warn "Skipping Docker installation"
        return
    }

    Write-Step "Installing Docker Desktop..."

    $pkgManager = Get-PackageManager

    switch ($pkgManager) {
        "winget" {
            winget install Docker.DockerDesktop
        }
        "choco" {
            choco install docker-desktop -y
        }
        default {
            Write-Warn "Cannot auto-install Docker. Please install manually from https://docker.com"
        }
    }

    Write-Info "Docker installed. You may need to restart your computer."
}

# Install Go tools
function Install-GoTools {
    Write-Step "Installing Go development tools..."

    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/cosmtrek/air@latest
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

    Write-Info "Go tools installed successfully"
}

# Setup project
function Initialize-Project {
    Push-Location $ProjectRoot

    Write-Step "Downloading Go dependencies..."
    go mod download
    go mod verify

    Write-Step "Installing Git hooks..."
    $hooksScript = Join-Path $PSScriptRoot "install-hooks.ps1"
    if (Test-Path $hooksScript) {
        & $hooksScript
    }

    Write-Step "Building project..."
    go build -o bin/jellyfin-go.exe ./cmd/jellyfin

    Pop-Location

    Write-Info "Project setup complete!"
}

# Main
Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   Jellyfin Go - Setup Assistant       ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

$pkgManager = Get-PackageManager
Write-Info "Detected package manager: $pkgManager"
Write-Host ""

Install-Go
Install-FFmpeg
Install-Docker
Install-GoTools
Initialize-Project

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║   ✅ Setup Complete!                   ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Info "Next steps:"
Write-Host "  1. Run the application: .\bin\jellyfin-go.exe"
Write-Host "  2. Or use development mode: make dev"
Write-Host "  3. Read the docs: cat README.md"
Write-Host ""
