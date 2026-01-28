terraform {
  required_providers {
    coder = {
      source = "coder/coder"
    }
    docker = {
      source = "kreuzwerker/docker"
    }
  }
}

locals {
  username = data.coder_workspace_owner.me.name
}

data "coder_workspace" "me" {}
data "coder_workspace_owner" "me" {}

# Docker image for development environment
resource "docker_image" "revenge" {
  name = "ghcr.io/lusoris/revenge-dev:latest"
  keep_locally = true
}

# PostgreSQL database container
resource "docker_container" "postgres" {
  image = docker_image.postgres.image_id
  name  = "revenge-postgres-${data.coder_workspace.me.id}"

  env = [
    "POSTGRES_USER=revenge",
    "POSTGRES_PASSWORD=revenge",
    "POSTGRES_DB=revenge",
  ]

  ports {
    internal = 5432
    external = 5432
  }

  volumes {
    container_path = "/var/lib/postgresql/data"
    volume_name    = docker_volume.postgres_data.name
  }

  healthcheck {
    test     = ["CMD-SHELL", "pg_isready -U revenge"]
    interval = "10s"
    timeout  = "5s"
    retries  = 5
  }
}

resource "docker_image" "postgres" {
  name = "postgres:18-alpine"
}

resource "docker_volume" "postgres_data" {
  name = "revenge-postgres-${data.coder_workspace.me.id}"
}

# Dragonfly cache container (Redis-compatible)
resource "docker_container" "dragonfly" {
  image = docker_image.dragonfly.image_id
  name  = "revenge-dragonfly-${data.coder_workspace.me.id}"

  command = [
    "dragonfly",
    "--requirepass=revenge",
    "--maxmemory=512mb",
  ]

  ports {
    internal = 6379
    external = 6379
  }

  volumes {
    container_path = "/data"
    volume_name    = docker_volume.dragonfly_data.name
  }
}

resource "docker_image" "dragonfly" {
  name = "docker.dragonflydb.io/dragonflydb/dragonfly:latest"
}

resource "docker_volume" "dragonfly_data" {
  name = "revenge-dragonfly-${data.coder_workspace.me.id}"
}

# Typesense search container
resource "docker_container" "typesense" {
  image = docker_image.typesense.image_id
  name  = "revenge-typesense-${data.coder_workspace.me.id}"

  env = [
    "TYPESENSE_DATA_DIR=/data",
    "TYPESENSE_API_KEY=revenge",
    "TYPESENSE_ENABLE_CORS=true",
  ]

  ports {
    internal = 8108
    external = 8108
  }

  volumes {
    container_path = "/data"
    volume_name    = docker_volume.typesense_data.name
  }
}

resource "docker_image" "typesense" {
  name = "typesense/typesense:27.1"
}

resource "docker_volume" "typesense_data" {
  name = "revenge-typesense-${data.coder_workspace.me.id}"
}

# Main workspace container
resource "docker_container" "workspace" {
  count = data.coder_workspace.me.start_count
  image = docker_image.revenge.image_id
  name  = "revenge-workspace-${data.coder_workspace.me.id}"

  hostname = data.coder_workspace.me.name

  env = [
    "CODER_AGENT_TOKEN=${coder_agent.main.token}",
    "DATABASE_URL=postgresql://revenge:revenge@postgres:5432/revenge?sslmode=disable",
    "CACHE_URL=redis://:revenge@dragonfly:6379/0",
    "TYPESENSE_URL=http://typesense:8108",
    "TYPESENSE_API_KEY=revenge",
  ]

  command = ["sh", "-c", coder_agent.main.init_script]

  volumes {
    container_path = "/workspace"
    volume_name    = docker_volume.workspace.name
  }

  volumes {
    container_path = "/home/coder/.cache/go-build"
    volume_name    = docker_volume.go_cache.name
  }

  volumes {
    container_path = "/go/pkg"
    volume_name    = docker_volume.go_pkg.name
  }

  # Link to service containers
  links = [
    docker_container.postgres.name,
    docker_container.dragonfly.name,
    docker_container.typesense.name,
  ]

  # Development ports
  ports {
    internal = 8080
    external = 8080
  }

  ports {
    internal = 5173
    external = 5173
  }
}

resource "docker_volume" "workspace" {
  name = "revenge-workspace-${data.coder_workspace.me.id}"
}

resource "docker_volume" "go_cache" {
  name = "revenge-go-cache-${data.coder_workspace.me.id}"
}

resource "docker_volume" "go_pkg" {
  name = "revenge-go-pkg-${data.coder_workspace.me.id}"
}

resource "coder_agent" "main" {
  arch                   = data.coder_provisioner.me.arch
  os                     = "linux"
  startup_script_timeout = 180
  startup_script         = <<-EOT
    #!/bin/bash
    set -e

    # Clone repository if not exists
    if [ ! -d "/workspace/revenge" ]; then
      git clone https://github.com/lusoris/revenge.git /workspace/revenge
    fi

    cd /workspace/revenge

    # Install Go dependencies
    go mod download

    # Generate code
    sqlc generate || true
    go generate ./... || true

    # Run database migrations
    go run ./cmd/revenge migrate up || true

    echo "🚀 Revenge development environment ready!"
  EOT

  env = {
    GIT_AUTHOR_NAME     = coalesce(data.coder_workspace_owner.me.full_name, data.coder_workspace_owner.me.name)
    GIT_AUTHOR_EMAIL    = "${data.coder_workspace_owner.me.email}"
    GIT_COMMITTER_NAME  = coalesce(data.coder_workspace_owner.me.full_name, data.coder_workspace_owner.me.name)
    GIT_COMMITTER_EMAIL = "${data.coder_workspace_owner.me.email}"
  }

  metadata {
    display_name = "CPU Usage"
    key          = "0_cpu_usage"
    script       = "coder stat cpu"
    interval     = 10
    timeout      = 1
  }

  metadata {
    display_name = "RAM Usage"
    key          = "1_ram_usage"
    script       = "coder stat mem"
    interval     = 10
    timeout      = 1
  }

  metadata {
    display_name = "Disk Usage"
    key          = "3_disk_usage"
    script       = "coder stat disk --path /workspace"
    interval     = 60
    timeout      = 1
  }
}

data "coder_provisioner" "me" {}

# Code Server (VS Code in browser)
resource "coder_app" "code-server" {
  agent_id     = coder_agent.main.id
  slug         = "code-server"
  display_name = "VS Code"
  url          = "http://localhost:13337/?folder=/workspace/revenge"
  icon         = "/icon/code.svg"
  subdomain    = false
  share        = "owner"

  healthcheck {
    url       = "http://localhost:13337/healthz"
    interval  = 5
    threshold = 6
  }
}

# Revenge API
resource "coder_app" "revenge-api" {
  agent_id     = coder_agent.main.id
  slug         = "revenge"
  display_name = "Revenge API"
  url          = "http://localhost:8080"
  icon         = "/icon/globe.svg"
  subdomain    = true
  share        = "owner"

  healthcheck {
    url       = "http://localhost:8080/health"
    interval  = 10
    threshold = 3
  }
}

# Frontend dev server
resource "coder_app" "frontend" {
  agent_id     = coder_agent.main.id
  slug         = "frontend"
  display_name = "Frontend (Vite)"
  url          = "http://localhost:5173"
  icon         = "/icon/svelte.svg"
  subdomain    = true
  share        = "owner"
}

# PostgreSQL UI
resource "coder_app" "postgres-ui" {
  agent_id     = coder_agent.main.id
  slug         = "pgadmin"
  display_name = "PostgreSQL Admin"
  url          = "http://localhost:8081"
  icon         = "/icon/database.svg"
  subdomain    = false
  share        = "owner"
}
