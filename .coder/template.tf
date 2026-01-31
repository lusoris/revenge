# Revenge Development Environment - Coder Template
#
# Supports multiple deployment backends:
# - Docker (default, for local development)
# - Kubernetes (for production-like environments)
# - Docker Swarm (for multi-node clusters)
#
# Host: https://coder.ancilla.lol

terraform {
  required_providers {
    coder = {
      source  = "coder/coder"
      version = ">= 2.0.0"
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
# Variables
# =============================================================================

variable "deployment_backend" {
  description = "Backend type: docker, kubernetes, or swarm"
  type        = string
  default     = "docker"
  validation {
    condition     = contains(["docker", "kubernetes", "swarm"], var.deployment_backend)
    error_message = "deployment_backend must be docker, kubernetes, or swarm."
  }
}

variable "namespace" {
  description = "Kubernetes namespace (only used when deployment_backend=kubernetes)"
  type        = string
  default     = "coder-workspaces"
}

variable "storage_class" {
  description = "Kubernetes storage class for persistent volumes"
  type        = string
  default     = "standard"
}

variable "cpu_limit" {
  description = "CPU limit for workspace container"
  type        = string
  default     = "4"
}

variable "memory_limit" {
  description = "Memory limit for workspace container"
  type        = string
  default     = "8Gi"
}

variable "disk_size" {
  description = "Disk size for workspace volume"
  type        = string
  default     = "20Gi"
}

variable "dotfiles_uri" {
  description = "Dotfiles repo URI (optional)"
  default     = ""
  type        = string
}

variable "git_clone_url" {
  description = "Git repository URL to clone"
  type        = string
  default     = "https://github.com/lusoris/revenge.git"
}

# =============================================================================
# Data Sources
# =============================================================================

data "coder_workspace" "me" {}
data "coder_workspace_owner" "me" {}
data "coder_provisioner" "me" {}

data "coder_parameter" "cpu" {
  name         = "cpu"
  display_name = "CPU Cores"
  description  = "Number of CPU cores"
  default      = "2"
  mutable      = true
  option {
    name  = "1 Core"
    value = "1"
  }
  option {
    name  = "2 Cores"
    value = "2"
  }
  option {
    name  = "4 Cores"
    value = "4"
  }
  option {
    name  = "8 Cores"
    value = "8"
  }
}

data "coder_parameter" "memory" {
  name         = "memory"
  display_name = "Memory"
  description  = "Amount of RAM"
  default      = "4"
  mutable      = true
  option {
    name  = "2 GB"
    value = "2"
  }
  option {
    name  = "4 GB"
    value = "4"
  }
  option {
    name  = "8 GB"
    value = "8"
  }
  option {
    name  = "16 GB"
    value = "16"
  }
}

data "coder_parameter" "ide" {
  name         = "ide"
  display_name = "IDE"
  description  = "Preferred development environment"
  default      = "vscode"
  mutable      = false
  option {
    name  = "VS Code (Browser)"
    value = "vscode"
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

locals {
  username     = data.coder_workspace_owner.me.name
  workspace_id = data.coder_workspace.me.id
  home_volume  = "revenge-home-${local.workspace_id}"

  # Common environment variables
  common_env = {
    DATABASE_URL      = "postgresql://revenge:revenge@postgres:5432/revenge?sslmode=disable"
    CACHE_URL         = "redis://:revenge@dragonfly:6379/0"
    TYPESENSE_URL     = "http://typesense:8108"
    TYPESENSE_API_KEY = "revenge"
  }
}

# =============================================================================
# Coder Agent
# =============================================================================

resource "coder_agent" "main" {
  arch                   = data.coder_provisioner.me.arch
  os                     = "linux"
  startup_script_timeout = 300
  startup_script         = <<-EOT
    #!/bin/bash
    set -e

    # Apply dotfiles if configured
    if [ -n "${var.dotfiles_uri}" ]; then
      coder dotfiles -y "${var.dotfiles_uri}"
    fi

    # Clone repository if not exists
    if [ ! -d "/workspace/revenge" ]; then
      git clone ${var.git_clone_url} /workspace/revenge
    fi

    cd /workspace/revenge

    # Install Go dependencies
    if command -v go &> /dev/null; then
      go mod download
    fi

    # Install Python dependencies for doc pipeline
    if command -v pip3 &> /dev/null; then
      pip3 install --user pyyaml requests beautifulsoup4 lxml html2text ruff pytest 2>/dev/null || true
    elif command -v python3 &> /dev/null; then
      python3 -m pip install --user pyyaml requests beautifulsoup4 lxml html2text ruff pytest 2>/dev/null || true
    fi

    # Install Node dependencies for frontend
    if [ -d "web" ] && command -v npm &> /dev/null; then
      cd web && npm install && cd ..
    fi

    # Generate code
    if command -v sqlc &> /dev/null; then
      sqlc generate || true
    fi
    if command -v go &> /dev/null; then
      go generate ./... || true
    fi

    # Run database migrations (wait for postgres)
    for i in {1..30}; do
      if pg_isready -h postgres -U revenge &>/dev/null; then
        go run ./cmd/revenge migrate up || true
        break
      fi
      sleep 2
    done

    # Start code-server if IDE is vscode
    if [ "${data.coder_parameter.ide.value}" = "vscode" ]; then
      code-server --auth none --port 13337 /workspace/revenge &
    fi

    echo "🚀 Revenge development environment ready!"
  EOT

  env = {
    GIT_AUTHOR_NAME     = coalesce(data.coder_workspace_owner.me.full_name, data.coder_workspace_owner.me.name)
    GIT_AUTHOR_EMAIL    = data.coder_workspace_owner.me.email
    GIT_COMMITTER_NAME  = coalesce(data.coder_workspace_owner.me.full_name, data.coder_workspace_owner.me.name)
    GIT_COMMITTER_EMAIL = data.coder_workspace_owner.me.email
    CODER_HOST          = "https://coder.ancilla.lol"
  }

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
    script       = "go version | awk '{print $3}'"
    interval     = 86400
    timeout      = 3
  }

  metadata {
    display_name = "Node Version"
    key          = "node_version"
    script       = "node --version"
    interval     = 86400
    timeout      = 3
  }
}

# =============================================================================
# Docker Backend (default)
# =============================================================================

# Only create Docker resources when backend is docker
resource "docker_image" "revenge" {
  count        = var.deployment_backend == "docker" ? 1 : 0
  name         = "ghcr.io/lusoris/revenge-dev:latest"
  keep_locally = true
}

resource "docker_image" "postgres" {
  count        = var.deployment_backend == "docker" ? 1 : 0
  name         = "postgres:16-alpine"
  keep_locally = true
}

resource "docker_image" "dragonfly" {
  count        = var.deployment_backend == "docker" ? 1 : 0
  name         = "docker.dragonflydb.io/dragonflydb/dragonfly:latest"
  keep_locally = true
}

resource "docker_image" "typesense" {
  count        = var.deployment_backend == "docker" ? 1 : 0
  name         = "typesense/typesense:27.1"
  keep_locally = true
}

# Docker Volumes with lifecycle protection
resource "docker_volume" "workspace" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-workspace-${local.workspace_id}"

  lifecycle {
    prevent_destroy = true
  }
}

resource "docker_volume" "postgres_data" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-postgres-${local.workspace_id}"

  lifecycle {
    prevent_destroy = true
  }
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

# Docker Network
resource "docker_network" "revenge" {
  count = var.deployment_backend == "docker" ? 1 : 0
  name  = "revenge-network-${local.workspace_id}"
}

# PostgreSQL Container
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
    "POSTGRES_PASSWORD=revenge",
    "POSTGRES_DB=revenge",
  ]

  volumes {
    container_path = "/var/lib/postgresql/data"
    volume_name    = docker_volume.postgres_data[0].name
  }

  healthcheck {
    test     = ["CMD-SHELL", "pg_isready -U revenge"]
    interval = "10s"
    timeout  = "5s"
    retries  = 5
  }
}

# Dragonfly Container
resource "docker_container" "dragonfly" {
  count = var.deployment_backend == "docker" ? 1 : 0
  image = docker_image.dragonfly[0].image_id
  name  = "revenge-dragonfly-${local.workspace_id}"

  networks_advanced {
    name    = docker_network.revenge[0].name
    aliases = ["dragonfly"]
  }

  command = [
    "dragonfly",
    "--requirepass=revenge",
    "--maxmemory=512mb",
  ]

  volumes {
    container_path = "/data"
    volume_name    = docker_volume.dragonfly_data[0].name
  }
}

# Typesense Container
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
    "TYPESENSE_API_KEY=revenge",
    "TYPESENSE_ENABLE_CORS=true",
  ]

  volumes {
    container_path = "/data"
    volume_name    = docker_volume.typesense_data[0].name
  }
}

# Main Workspace Container
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
    "DATABASE_URL=${local.common_env.DATABASE_URL}",
    "CACHE_URL=${local.common_env.CACHE_URL}",
    "TYPESENSE_URL=${local.common_env.TYPESENSE_URL}",
    "TYPESENSE_API_KEY=${local.common_env.TYPESENSE_API_KEY}",
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

  # Resource limits
  cpu_shares = data.coder_parameter.cpu.value * 1024
  memory     = data.coder_parameter.memory.value * 1024

  depends_on = [
    docker_container.postgres[0],
    docker_container.dragonfly[0],
    docker_container.typesense[0],
  ]
}

# =============================================================================
# Kubernetes Backend
# =============================================================================

resource "kubernetes_namespace" "workspace" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name = "${var.namespace}-${local.username}"
    labels = {
      "coder.workspace" = local.workspace_id
      "coder.owner"     = local.username
    }
  }
}

resource "kubernetes_persistent_volume_claim" "workspace" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name      = "workspace-pvc"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    access_modes       = ["ReadWriteOnce"]
    storage_class_name = var.storage_class
    resources {
      requests = {
        storage = var.disk_size
      }
    }
  }
}

resource "kubernetes_persistent_volume_claim" "postgres" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name      = "postgres-pvc"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    access_modes       = ["ReadWriteOnce"]
    storage_class_name = var.storage_class
    resources {
      requests = {
        storage = "10Gi"
      }
    }
  }
}

# PostgreSQL StatefulSet (Kubernetes)
resource "kubernetes_stateful_set" "postgres" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name      = "postgres"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    service_name = "postgres"
    replicas     = 1
    selector {
      match_labels = {
        app = "postgres"
      }
    }
    template {
      metadata {
        labels = {
          app = "postgres"
        }
      }
      spec {
        container {
          name  = "postgres"
          image = "postgres:16-alpine"
          env {
            name  = "POSTGRES_USER"
            value = "revenge"
          }
          env {
            name  = "POSTGRES_PASSWORD"
            value = "revenge"
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
            limits = {
              cpu    = "1"
              memory = "1Gi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
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

# Kubernetes Service for PostgreSQL
resource "kubernetes_service" "postgres" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name      = "postgres"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    selector = {
      app = "postgres"
    }
    port {
      port        = 5432
      target_port = 5432
    }
    type = "ClusterIP"
  }
}

# Dragonfly Deployment (Kubernetes)
resource "kubernetes_deployment" "dragonfly" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name      = "dragonfly"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "dragonfly"
      }
    }
    template {
      metadata {
        labels = {
          app = "dragonfly"
        }
      }
      spec {
        container {
          name  = "dragonfly"
          image = "docker.dragonflydb.io/dragonflydb/dragonfly:latest"
          args  = ["dragonfly", "--requirepass=revenge", "--maxmemory=512mb"]
          port {
            container_port = 6379
          }
          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "dragonfly" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name      = "dragonfly"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    selector = {
      app = "dragonfly"
    }
    port {
      port        = 6379
      target_port = 6379
    }
    type = "ClusterIP"
  }
}

# Typesense Deployment (Kubernetes)
resource "kubernetes_deployment" "typesense" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name      = "typesense"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "typesense"
      }
    }
    template {
      metadata {
        labels = {
          app = "typesense"
        }
      }
      spec {
        container {
          name  = "typesense"
          image = "typesense/typesense:27.1"
          env {
            name  = "TYPESENSE_DATA_DIR"
            value = "/data"
          }
          env {
            name  = "TYPESENSE_API_KEY"
            value = "revenge"
          }
          env {
            name  = "TYPESENSE_ENABLE_CORS"
            value = "true"
          }
          port {
            container_port = 8108
          }
          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "typesense" {
  count = var.deployment_backend == "kubernetes" ? 1 : 0
  metadata {
    name      = "typesense"
    namespace = kubernetes_namespace.workspace[0].metadata[0].name
  }
  spec {
    selector = {
      app = "typesense"
    }
    port {
      port        = 8108
      target_port = 8108
    }
    type = "ClusterIP"
  }
}

# Main Workspace Pod (Kubernetes)
resource "kubernetes_pod" "workspace" {
  count = var.deployment_backend == "kubernetes" ? data.coder_workspace.me.start_count : 0
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
      image   = "ghcr.io/lusoris/revenge-dev:latest"
      command = ["sh", "-c", coder_agent.main.init_script]
      env {
        name  = "CODER_AGENT_TOKEN"
        value = coder_agent.main.token
      }
      env {
        name  = "DATABASE_URL"
        value = local.common_env.DATABASE_URL
      }
      env {
        name  = "CACHE_URL"
        value = local.common_env.CACHE_URL
      }
      env {
        name  = "TYPESENSE_URL"
        value = local.common_env.TYPESENSE_URL
      }
      env {
        name  = "TYPESENSE_API_KEY"
        value = local.common_env.TYPESENSE_API_KEY
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
# Coder Apps
# =============================================================================

resource "coder_app" "code_server" {
  count        = data.coder_parameter.ide.value == "vscode" ? 1 : 0
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

resource "coder_app" "frontend" {
  agent_id     = coder_agent.main.id
  slug         = "frontend"
  display_name = "Frontend (Vite)"
  url          = "http://localhost:5173"
  icon         = "/icon/svelte.svg"
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
# Outputs
# =============================================================================

output "workspace_url" {
  value       = "https://coder.ancilla.lol/@${local.username}/${data.coder_workspace.me.name}"
  description = "URL to access the workspace"
}

output "deployment_backend" {
  value       = var.deployment_backend
  description = "Active deployment backend"
.}
