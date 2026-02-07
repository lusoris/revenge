# Revenge Development Environment - Coder Template
#
# Backends: Docker (default), Kubernetes, K3s, Docker Swarm
# IDEs:     Zed (SSH, default), VS Code (browser/desktop), JetBrains Gateway, Terminal
# Host:     https://coder.ancilla.lol
#
# Versions sourced from go.mod + docker-compose.yml (NOT design docs)

terraform {
  required_providers {
    coder = {
      source  = "coder/coder"
      version = ">= 2.4.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 3.0.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
  }
}

# =============================================================================
# Image versions - Single source of truth for container tags
# Keep in sync with docker-compose.yml and docker-compose.dev.yml
# =============================================================================

locals {
  # Actual versions from docker-compose files and go.mod
  go_version      = "1.25.6"
  postgres_image  = "postgres:18-alpine"
  dragonfly_image = "docker.dragonflydb.io/dragonflydb/dragonfly:latest"
  typesense_image = "typesense/typesense:0.25.2"
  workspace_image = "golang:1.25-alpine"

  username     = data.coder_workspace_owner.me.name
  workspace_id = data.coder_workspace.me.id

  # Shared env vars matching docker-compose.dev.yml
  app_env = {
    GOEXPERIMENT           = "greenteagc,jsonv2"
    REVENGE_LOG_LEVEL      = "debug"
    REVENGE_LOG_FORMAT     = "console"
    REVENGE_DATABASE_URL   = "postgres://revenge:revenge_dev_pass@postgres:5432/revenge?sslmode=disable"
    REVENGE_CACHE_ENABLED  = "true"
    REVENGE_CACHE_URL      = "redis://dragonfly:6379"
    REVENGE_SEARCH_ENABLED = "true"
    REVENGE_SEARCH_URL     = "http://typesense:8108"
    REVENGE_SEARCH_API_KEY = "dev_api_key"
  }

  # K8s/K3s share the same resources
  use_k8s = contains(["kubernetes", "k3s"], var.deployment_backend)
}

# =============================================================================
# Data Sources
# =============================================================================

data "coder_workspace" "me" {}
data "coder_workspace_owner" "me" {}
data "coder_provisioner" "me" {}

# =============================================================================
# Parameters
# =============================================================================

variable "deployment_backend" {
  description = "Backend type"
  type        = string
  default     = "docker"
  validation {
    condition     = contains(["docker", "kubernetes", "k3s", "swarm"], var.deployment_backend)
    error_message = "Must be: docker, kubernetes, k3s, or swarm."
  }
}

variable "namespace" {
  description = "Kubernetes/K3s namespace"
  type        = string
  default     = "coder-workspaces"
}

variable "storage_class" {
  description = "Kubernetes storage class for persistent volumes"
  type        = string
  default     = "local-path"
}

variable "disk_size" {
  description = "Disk size for workspace volume"
  type        = string
  default     = "30Gi"
}

variable "dotfiles_uri" {
  description = "Dotfiles repo URI (optional)"
  type        = string
  default     = ""
}

variable "git_clone_url" {
  description = "Git repository URL to clone"
  type        = string
  default     = "https://github.com/lusoris/revenge.git"
}

data "coder_parameter" "cpu" {
  name         = "cpu"
  display_name = "CPU Cores"
  description  = "Number of CPU cores for the workspace"
  default      = "4"
  mutable      = true
  option {
    name  = "2 Cores"
    value = "2"
  }
  option {
    name  = "4 Cores (recommended)"
    value = "4"
  }
  option {
    name  = "8 Cores"
    value = "8"
  }
  option {
    name  = "16 Cores"
    value = "16"
  }
}

data "coder_parameter" "memory" {
  name         = "memory"
  display_name = "Memory (GB)"
  description  = "RAM allocation"
  default      = "8"
  mutable      = true
  option {
    name  = "4 GB"
    value = "4"
  }
  option {
    name  = "8 GB (recommended)"
    value = "8"
  }
  option {
    name  = "16 GB"
    value = "16"
  }
  option {
    name  = "32 GB"
    value = "32"
  }
}

data "coder_parameter" "ide" {
  name         = "ide"
  display_name = "IDE"
  description  = "Primary development environment"
  default      = "zed"
  mutable      = false
  option {
    name  = "Zed (SSH) - Recommended"
    value = "zed"
    icon  = "/icon/terminal.svg"
  }
  option {
    name  = "VS Code (Browser)"
    value = "vscode-browser"
    icon  = "/icon/code.svg"
  }
  option {
    name  = "VS Code (Desktop)"
    value = "vscode-desktop"
    icon  = "/icon/code.svg"
  }
  option {
    name  = "JetBrains Gateway"
    value = "jetbrains"
    icon  = "/icon/gateway.svg"
  }
  option {
    name  = "Terminal Only"
    value = "terminal"
    icon  = "/icon/terminal.svg"
  }
}

# =============================================================================
# Workspace presets
# =============================================================================

data "coder_workspace_preset" "lightweight" {
  name        = "Lightweight"
  description = "Quick edits, docs, reviews (Zed + 2 CPU, 4 GB)"
  parameters = {
    "cpu"    = "2"
    "memory" = "4"
    "ide"    = "zed"
  }
}

data "coder_workspace_preset" "standard" {
  name        = "Standard Development"
  description = "Daily Go/Svelte development (Zed + 4 CPU, 8 GB)"
  parameters = {
    "cpu"    = "4"
    "memory" = "8"
    "ide"    = "zed"
  }
}

data "coder_workspace_preset" "heavy" {
  name        = "Heavy Build"
  description = "Full builds, integration tests (8 CPU, 16 GB)"
  parameters = {
    "cpu"    = "8"
    "memory" = "16"
    "ide"    = "zed"
  }
}

# =============================================================================
# Coder Agent
# =============================================================================

resource "coder_agent" "main" {
  arch                   = data.coder_provisioner.me.arch
  os                     = "linux"
  startup_script_timeout = 600

  startup_script = <<-EOT
    #!/bin/sh
    set -eu

    export GOEXPERIMENT=greenteagc,jsonv2

    # ── System dependencies (CGO: vips + ffmpeg) ───────────────────────
    apk add --no-cache \
      bash curl git gcc musl-dev pkgconfig \
      ffmpeg-dev vips-dev \
      postgresql-client \
      make tar zstd openssh-client

    # ── Dotfiles ───────────────────────────────────────────────────────
    if [ -n "${var.dotfiles_uri}" ]; then
      coder dotfiles -y "${var.dotfiles_uri}"
    fi

    # ── Clone repo ─────────────────────────────────────────────────────
    if [ ! -d "/workspace/revenge" ]; then
      git clone --depth 1 ${var.git_clone_url} /workspace/revenge
    fi
    cd /workspace/revenge

    # ── Go dependencies ────────────────────────────────────────────────
    go mod download &

    # ── Dev tools (install in background) ──────────────────────────────
    (
      go install golang.org/x/tools/gopls@latest
      go install github.com/air-verse/air@latest
      go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
      go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
      go install github.com/vektra/mockery/v3@latest
      go install github.com/go-delve/delve/cmd/dlv@latest
      go install golang.org/x/vuln/cmd/govulncheck@latest
      # golangci-lint v2 requires curl install (go install only works for v1)
      curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.8.0
    ) &

    # ── Node deps (for frontend, when web/ exists) ─────────────────────
    if [ -d "web" ] && command -v npm >/dev/null 2>&1; then
      (cd web && npm install --prefer-offline) &
    fi

    # ── Wait for background installs ──────────────────────────────────
    wait

    # ── Code generation ────────────────────────────────────────────────
    if command -v sqlc >/dev/null 2>&1; then
      sqlc generate 2>/dev/null || true
    fi

    # ── Database migrations (wait for postgres) ────────────────────────
    for i in $(seq 1 30); do
      if pg_isready -h postgres -U revenge >/dev/null 2>&1; then
        migrate -path internal/infra/database/migrations/shared -database "$${REVENGE_DATABASE_URL}" up 2>/dev/null || true
        break
      fi
      sleep 2
    done

    # ── Start code-server if VS Code browser selected ──────────────────
    if [ "${data.coder_parameter.ide.value}" = "vscode-browser" ]; then
      code-server --auth none --port 13337 /workspace/revenge &
    fi

    echo "Revenge dev environment ready!"
  EOT

  env = {
    GOEXPERIMENT        = "greenteagc,jsonv2"
    GIT_AUTHOR_NAME     = coalesce(data.coder_workspace_owner.me.full_name, data.coder_workspace_owner.me.name)
    GIT_AUTHOR_EMAIL    = data.coder_workspace_owner.me.email
    GIT_COMMITTER_NAME  = coalesce(data.coder_workspace_owner.me.full_name, data.coder_workspace_owner.me.name)
    GIT_COMMITTER_EMAIL = data.coder_workspace_owner.me.email
    CODER_HOST          = "https://coder.ancilla.lol"
    REVENGE_LOG_LEVEL      = local.app_env.REVENGE_LOG_LEVEL
    REVENGE_LOG_FORMAT     = local.app_env.REVENGE_LOG_FORMAT
    REVENGE_DATABASE_URL   = local.app_env.REVENGE_DATABASE_URL
    REVENGE_CACHE_ENABLED  = local.app_env.REVENGE_CACHE_ENABLED
    REVENGE_CACHE_URL      = local.app_env.REVENGE_CACHE_URL
    REVENGE_SEARCH_ENABLED = local.app_env.REVENGE_SEARCH_ENABLED
    REVENGE_SEARCH_URL     = local.app_env.REVENGE_SEARCH_URL
    REVENGE_SEARCH_API_KEY = local.app_env.REVENGE_SEARCH_API_KEY
  }

  # ── Metadata probes ───────────────────────────────────────────────
  metadata {
    display_name = "CPU Usage"
    key          = "cpu_usage"
    script       = "coder stat cpu"
    interval     = 10
    timeout      = 1
  }
  metadata {
    display_name = "RAM Usage"
    key          = "ram_usage"
    script       = "coder stat mem"
    interval     = 10
    timeout      = 1
  }
  metadata {
    display_name = "Disk Usage"
    key          = "disk_usage"
    script       = "coder stat disk --path /workspace"
    interval     = 60
    timeout      = 1
  }
  metadata {
    display_name = "Go Version"
    key          = "go_version"
    script       = "go version 2>/dev/null | awk '{print $3}' || echo 'N/A'"
    interval     = 86400
    timeout      = 3
  }
  metadata {
    display_name = "PostgreSQL"
    key          = "pg_status"
    script       = "pg_isready -h postgres -U revenge &>/dev/null && echo 'healthy' || echo 'down'"
    interval     = 30
    timeout      = 3
  }
}

# =============================================================================
# Coder Apps (IDE launchers & service UIs)
# =============================================================================

resource "coder_app" "code_server" {
  count        = data.coder_parameter.ide.value == "vscode-browser" ? 1 : 0
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

resource "coder_app" "revenge_api" {
  agent_id     = coder_agent.main.id
  slug         = "revenge-api"
  display_name = "Revenge API"
  url          = "http://localhost:8096"
  icon         = "/icon/globe.svg"
  subdomain    = true
  share        = "owner"

  healthcheck {
    url       = "http://localhost:8096/health/live"
    interval  = 10
    threshold = 3
  }
}

resource "coder_app" "frontend" {
  agent_id     = coder_agent.main.id
  slug         = "frontend"
  display_name = "Frontend (SvelteKit)"
  url          = "http://localhost:5173"
  icon         = "/icon/widgets.svg"
  subdomain    = true
  share        = "owner"
}

resource "coder_app" "terminal" {
  agent_id     = coder_agent.main.id
  slug         = "terminal"
  display_name = "Terminal"
  icon         = "/icon/terminal.svg"
  url          = "http://localhost"
  subdomain    = false
  share        = "owner"
}

# =============================================================================
# Docker Backend (default)
# =============================================================================

resource "docker_image" "revenge" {
  count        = var.deployment_backend == "docker" ? 1 : 0
  name         = local.workspace_image
  keep_locally = true
}

resource "docker_image" "postgres" {
  count        = var.deployment_backend == "docker" ? 1 : 0
  name         = local.postgres_image
  keep_locally = true
}

resource "docker_image" "dragonfly" {
  count        = var.deployment_backend == "docker" ? 1 : 0
  name         = local.dragonfly_image
  keep_locally = true
}

resource "docker_image" "typesense" {
  count        = var.deployment_backend == "docker" ? 1 : 0
  name         = local.typesense_image
  keep_locally = true
}

# ── Docker Volumes ─────────────────────────────────────────────────

resource "docker_volume" "workspace" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-workspace-${local.workspace_id}"
  lifecycle { prevent_destroy = true }
}

resource "docker_volume" "postgres_data" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-postgres-${local.workspace_id}"
  lifecycle { prevent_destroy = true }
}

resource "docker_volume" "dragonfly_data" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-dragonfly-${local.workspace_id}"
}

resource "docker_volume" "typesense_data" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-typesense-${local.workspace_id}"
}

resource "docker_volume" "go_cache" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-go-cache-${local.workspace_id}"
}

resource "docker_volume" "go_pkg" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-go-pkg-${local.workspace_id}"
}

# ── Docker Network ─────────────────────────────────────────────────

resource "docker_network" "revenge" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-net-${local.workspace_id}"
}

# ── PostgreSQL Container ───────────────────────────────────────────

resource "docker_container" "postgres" {
  count = var.deployment_backend == "docker" ? 1 : 0
  image = docker_image.postgres[0].image_id
  name  = "revenge-postgres-${local.workspace_id}"

  networks_advanced {
    name    = docker_network.revenge[0].name
    aliases = ["postgres"]
  }

  env = [
    "POSTGRES_USER=revenge",
    "POSTGRES_PASSWORD=revenge_dev_pass",
    "POSTGRES_DB=revenge",
  ]

  volumes {
    container_path = "/var/lib/postgresql/data"
    volume_name    = docker_volume.postgres_data[0].name
  }

  healthcheck {
    test     = ["CMD-SHELL", "pg_isready -U revenge"]
    interval = "5s"
    timeout  = "5s"
    retries  = 5
  }
}

# ── Dragonfly Container ───────────────────────────────────────────

resource "docker_container" "dragonfly" {
  count = var.deployment_backend == "docker" ? 1 : 0
  image = docker_image.dragonfly[0].image_id
  name  = "revenge-dragonfly-${local.workspace_id}"

  networks_advanced {
    name    = docker_network.revenge[0].name
    aliases = ["dragonfly"]
  }

  command = ["--cache_mode"]

  volumes {
    container_path = "/data"
    volume_name    = docker_volume.dragonfly_data[0].name
  }

  healthcheck {
    test     = ["CMD", "redis-cli", "ping"]
    interval = "5s"
    timeout  = "3s"
    retries  = 5
  }
}

# ── Typesense Container ───────────────────────────────────────────

resource "docker_container" "typesense" {
  count = var.deployment_backend == "docker" ? 1 : 0
  image = docker_image.typesense[0].image_id
  name  = "revenge-typesense-${local.workspace_id}"

  networks_advanced {
    name    = docker_network.revenge[0].name
    aliases = ["typesense"]
  }

  env = [
    "TYPESENSE_DATA_DIR=/data",
    "TYPESENSE_API_KEY=dev_api_key",
    "TYPESENSE_ENABLE_CORS=true",
  ]

  volumes {
    container_path = "/data"
    volume_name    = docker_volume.typesense_data[0].name
  }

  healthcheck {
    test     = ["CMD", "curl", "-f", "http://localhost:8108/health"]
    interval = "5s"
    timeout  = "3s"
    retries  = 5
  }
}

# ── Main Workspace Container ──────────────────────────────────────

resource "docker_container" "workspace" {
  count = var.deployment_backend == "docker" ? data.coder_workspace.me.start_count : 0
  image = docker_image.revenge[0].image_id
  name  = "revenge-workspace-${local.workspace_id}"

  hostname   = data.coder_workspace.me.name
  domainname = "coder.ancilla.lol"

  networks_advanced {
    name = docker_network.revenge[0].name
  }

  env = [
    "CODER_AGENT_TOKEN=${coder_agent.main.token}",
    "GOEXPERIMENT=greenteagc,jsonv2",
    "REVENGE_DATABASE_URL=${local.app_env.REVENGE_DATABASE_URL}",
    "REVENGE_CACHE_URL=${local.app_env.REVENGE_CACHE_URL}",
    "REVENGE_SEARCH_URL=${local.app_env.REVENGE_SEARCH_URL}",
    "REVENGE_SEARCH_API_KEY=${local.app_env.REVENGE_SEARCH_API_KEY}",
    "REVENGE_LOG_LEVEL=debug",
    "REVENGE_LOG_FORMAT=console",
  ]

  command = ["sh", "-c", coder_agent.main.init_script]

  volumes {
    container_path = "/workspace"
    volume_name    = docker_volume.workspace[0].name
  }
  volumes {
    container_path = "/home/coder/.cache/go-build"
    volume_name    = docker_volume.go_cache[0].name
  }
  volumes {
    container_path = "/go/pkg"
    volume_name    = docker_volume.go_pkg[0].name
  }

  cpu_shares = data.coder_parameter.cpu.value * 1024
  memory     = data.coder_parameter.memory.value * 1024

  depends_on = [
    docker_container.postgres[0],
    docker_container.dragonfly[0],
    docker_container.typesense[0],
  ]
}

# =============================================================================
# Kubernetes / K3s Backend
# =============================================================================

resource "kubernetes_namespace" "workspace" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name = "${var.namespace}-${local.username}"
    labels = {
      "coder.workspace" = local.workspace_id
      "coder.owner"     = local.username
    }
  }
}

# ── PVCs ──────────────────────────────────────────────────────────

resource "kubernetes_persistent_volume_claim" "workspace" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name      = "workspace-pvc"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    access_modes       = ["ReadWriteOnce"]
    storage_class_name = var.storage_class
    resources {
      requests = { storage = var.disk_size }
    }
  }
}

resource "kubernetes_persistent_volume_claim" "postgres" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name      = "postgres-pvc"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    access_modes       = ["ReadWriteOnce"]
    storage_class_name = var.storage_class
    resources {
      requests = { storage = "10Gi" }
    }
  }
}

# ── PostgreSQL (K8s/K3s) ──────────────────────────────────────────

resource "kubernetes_stateful_set" "postgres" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name      = "postgres"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    service_name = "postgres"
    replicas     = 1
    selector {
      match_labels = { app = "postgres" }
    }
    template {
      metadata {
        labels = { app = "postgres" }
      }
      spec {
        container {
          name  = "postgres"
          image = local.postgres_image
          env {
            name  = "POSTGRES_USER"
            value = "revenge"
          }
          env {
            name  = "POSTGRES_PASSWORD"
            value = "revenge_dev_pass"
          }
          env {
            name  = "POSTGRES_DB"
            value = "revenge"
          }
          port {
            container_port = 5432
          }
          volume_mount {
            name       = "postgres-data"
            mount_path = "/var/lib/postgresql/data"
          }
          resources {
            limits   = { cpu = "1", memory = "1Gi" }
            requests = { cpu = "250m", memory = "256Mi" }
          }
        }
        volume {
          name = "postgres-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.postgres[0].metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "postgres" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name      = "postgres"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    selector = { app = "postgres" }
    port {
      port        = 5432
      target_port = 5432
    }
    type = "ClusterIP"
  }
}

# ── Dragonfly (K8s/K3s) ──────────────────────────────────────────

resource "kubernetes_deployment" "dragonfly" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name      = "dragonfly"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    replicas = 1
    selector {
      match_labels = { app = "dragonfly" }
    }
    template {
      metadata {
        labels = { app = "dragonfly" }
      }
      spec {
        container {
          name  = "dragonfly"
          image = local.dragonfly_image
          args  = ["--cache_mode", "--maxmemory=512mb"]
          port {
            container_port = 6379
          }
          resources {
            limits   = { cpu = "500m", memory = "512Mi" }
            requests = { cpu = "100m", memory = "128Mi" }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "dragonfly" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name      = "dragonfly"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    selector = { app = "dragonfly" }
    port {
      port        = 6379
      target_port = 6379
    }
    type = "ClusterIP"
  }
}

# ── Typesense (K8s/K3s) ──────────────────────────────────────────

resource "kubernetes_deployment" "typesense" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name      = "typesense"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    replicas = 1
    selector {
      match_labels = { app = "typesense" }
    }
    template {
      metadata {
        labels = { app = "typesense" }
      }
      spec {
        container {
          name  = "typesense"
          image = local.typesense_image
          env {
            name  = "TYPESENSE_DATA_DIR"
            value = "/data"
          }
          env {
            name  = "TYPESENSE_API_KEY"
            value = "dev_api_key"
          }
          env {
            name  = "TYPESENSE_ENABLE_CORS"
            value = "true"
          }
          port {
            container_port = 8108
          }
          resources {
            limits   = { cpu = "500m", memory = "512Mi" }
            requests = { cpu = "100m", memory = "128Mi" }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "typesense" {
  count = local.use_k8s ? 1 : 0
  metadata {
    name      = "typesense"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    selector = { app = "typesense" }
    port {
      port        = 8108
      target_port = 8108
    }
    type = "ClusterIP"
  }
}

# ── Main Workspace Pod (K8s/K3s) ─────────────────────────────────

resource "kubernetes_pod" "workspace" {
  count = local.use_k8s ? data.coder_workspace.me.start_count : 0
  metadata {
    name      = "workspace"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
    labels = {
      "coder.workspace" = local.workspace_id
    }
  }
  spec {
    container {
      name    = "dev"
      image   = local.workspace_image
      command = ["sh", "-c", coder_agent.main.init_script]

      env {
        name  = "CODER_AGENT_TOKEN"
        value = coder_agent.main.token
      }
      env {
        name  = "GOEXPERIMENT"
        value = "greenteagc,jsonv2"
      }
      env {
        name  = "REVENGE_DATABASE_URL"
        value = local.app_env.REVENGE_DATABASE_URL
      }
      env {
        name  = "REVENGE_CACHE_URL"
        value = local.app_env.REVENGE_CACHE_URL
      }
      env {
        name  = "REVENGE_SEARCH_URL"
        value = local.app_env.REVENGE_SEARCH_URL
      }
      env {
        name  = "REVENGE_SEARCH_API_KEY"
        value = local.app_env.REVENGE_SEARCH_API_KEY
      }
      env {
        name  = "REVENGE_LOG_LEVEL"
        value = "debug"
      }

      volume_mount {
        name       = "workspace"
        mount_path = "/workspace"
      }

      resources {
        limits = {
          cpu    = data.coder_parameter.cpu.value
          memory = "${data.coder_parameter.memory.value}Gi"
        }
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
      }
    }
    volume {
      name = "workspace"
      persistent_volume_claim {
        claim_name = kubernetes_persistent_volume_claim.workspace[0].metadata[0].name
      }
    }
  }
  depends_on = [
    kubernetes_stateful_set.postgres[0],
    kubernetes_deployment.dragonfly[0],
    kubernetes_deployment.typesense[0],
  ]
}

# =============================================================================
# Docker Swarm Backend
# =============================================================================
#
# For Swarm, infra services (postgres, dragonfly, typesense) run as part of
# the stack defined in deploy/docker-swarm-stack.yml.
#
# Deploy the stack first:
#   docker stack deploy -c deploy/docker-swarm-stack.yml revenge
#
# Then create the Coder workspace:
#   coder create my-workspace --template revenge --variable deployment_backend=swarm

# =============================================================================
# Outputs
# =============================================================================

output "workspace_url" {
  value       = "https://coder.ancilla.lol/@${local.username}/${data.coder_workspace.me.name}"
  description = "URL to access the workspace"
}

output "deployment_backend" {
  value       = var.deployment_backend
  description = "Active deployment backend"
}

output "ide" {
  value       = data.coder_parameter.ide.value
  description = "Selected IDE"
}

output "zed_ssh" {
  value       = "zed ssh://${data.coder_workspace.me.name}/workspace/revenge"
  description = "Zed SSH connection command"
}
