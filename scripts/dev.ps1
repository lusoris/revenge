# Development helper script for Windows

param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('check', 'install-tools', 'setup', 'test', 'lint', 'dev', 'build')]
    [string]$Command
)

$ErrorActionPreference = "Stop"

$ProjectRoot = Split-Path -Parent $PSScriptRoot

Push-Location $ProjectRoot

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

function Test-Requirements {
    Write-Info "Checking requirements..."

    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Error "Go is not installed. Please install Go 1.24 or later."
        exit 1
    }

    $goVersion = (go version) -replace 'go version go', '' -replace ' .*', ''
    Write-Info "Go version: $goVersion"

    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-Warn "Docker is not installed. Docker is optional but recommended."
    }
    else {
        Write-Info "Docker is installed"
    }
}

function Install-Tools {
    Write-Info "Installing development tools..."

    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/cosmtrek/air@latest
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

    Write-Info "Development tools installed"
}

function Initialize-Database {
    Write-Info "Setting up database..."

    Write-Info "Checking PostgreSQL connection..."
    $env:PGPASSWORD = "password"

    # Check if PostgreSQL is reachable
    try {
        $null = & psql -h localhost -U jellyfin -d jellyfin -c "SELECT 1" 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Info "PostgreSQL connection successful"
        }
        else {
            Write-Warning "PostgreSQL not reachable. Please ensure PostgreSQL is running."
            Write-Warning "Start with: docker-compose -f docker-compose.dev.yml up -d postgres"
        }
    }
    catch {
        Write-Warning "PostgreSQL client not found. Install psql or use Docker."
    }

    Write-Info "Database ready"
}

function Invoke-Tests {
    Write-Info "Running tests..."
    go test -v -race -coverprofile=coverage.out ./...
}

function Invoke-Lint {
    Write-Info "Running linter..."
    golangci-lint run --timeout=5m
}

function Start-DevServer {
    Write-Info "Starting development server with hot reload..."
    air
}

function Build-Binary {
    Write-Info "Building binary..."
    New-Item -ItemType Directory -Force -Path "bin" | Out-Null
    go build -o bin\jellyfin-go.exe .\cmd\jellyfin
    Write-Info "Binary built: bin\jellyfin-go.exe"
}

# Execute command
switch ($Command) {
    'check' {
        Test-Requirements
    }
    'install-tools' {
        Install-Tools
    }
    'setup' {
        Test-Requirements
        Install-Tools
        Initialize-Database
        Write-Info "Setup complete! Run '.\scripts\dev.ps1 dev' to start development server"
    }
    'test' {
        Invoke-Tests
    }
    'lint' {
        Invoke-Lint
    }
    'dev' {
        Start-DevServer
    }
    'build' {
        Build-Binary
    }
}

Pop-Location
