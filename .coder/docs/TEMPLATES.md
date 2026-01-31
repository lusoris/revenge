# Coder Templates Reference

> Source: https://coder.com/docs/admin/templates
> Fetched: 2026-01-31
> Type: html

---

## Overview

Templates form the foundation of Coder workspaces, written in Terraform to define underlying infrastructure. They enable organizations to standardize development environments across teams.

---

## Core Concepts

### What Are Templates?

Templates are written in Terraform and define the underlying infrastructure that all Coder workspaces run on. Templates provide a standardized way to configure workspace resources, ensuring consistency across your deployment.

### Learning Path

1. **Foundational Learning**: Start by creating a basic template from scratch
2. **Terraform Familiarity**: Consult HashiCorp's tutorials if you're new to Terraform
3. **Starter Templates**: Leverage pre-built templates for popular platforms (AWS, Kubernetes, Docker)

---

## Template Creation

### Creating Templates

- **From Scratch**: Build custom templates tailored to specific needs
- **Starter Templates**: Import templates with sensible defaults
- **Cloning Existing**: Modify proven templates for new use cases

### Common Customizations

| Area | Description |
|------|-------------|
| Container Images | Pre-installed languages, tools, dependencies |
| Parameters | Build-time variables (disk size, instance type, region) |
| IDE & Features | JetBrains IDEs, code editors, RDP, dotfiles |

---

## Terraform Provider Resources

### Core Resources

```hcl
# Workspace data
data "coder_workspace" "me" {}
data "coder_workspace_owner" "me" {}

# Application resource (IDE, web UI)
resource "coder_app" "code-server" {
  agent_id     = coder_agent.main.id
  slug         = "code-server"
  display_name = "VS Code"
  url          = "http://localhost:8080/?folder=/home/coder"
  icon         = "/icon/code.svg"
  subdomain    = true
}

# Build parameters
data "coder_parameter" "region" {
  name         = "region"
  display_name = "Region"
  type         = "string"
  default      = "us-east-1"
  icon         = "/icon/aws.svg"

  option {
    name  = "US East"
    value = "us-east-1"
    icon  = "/emojis/flag-us.png"
  }
}

# Metadata display
resource "coder_metadata" "workspace_info" {
  resource_id = docker_container.workspace.id

  item {
    key   = "CPU"
    value = "4 cores"
  }
}

# Scripts
resource "coder_script" "startup" {
  agent_id     = coder_agent.main.id
  display_name = "Startup"
  script       = file("${path.module}/startup.sh")
  run_on_start = true
}
```

---

## Best Practices

### Deployment Strategy

1. Begin with a universal template handling basic development tasks
2. Create specialized templates as needs evolve

### Key Implementation Areas

| Area | Recommendation |
|------|----------------|
| Image Management | Develop optimized container images |
| Dev Containers | Enable native support via `@devcontainers/cli` |
| Template Hardening | Protect critical resources (user disks) |
| CI/CD Integration | Version control templates, use GitOps |
| Access Control | Implement permissions and policies |

### Resource Persistence

Configure templates to prevent destruction of specific resources:

```hcl
resource "docker_volume" "user_data" {
  name = "user-data-${data.coder_workspace.me.id}"

  lifecycle {
    prevent_destroy = true
  }
}
```

---

## Icon Configuration

### Built-in Icons

Icons can be specified via URL:
- Bundled: `/icon/coder.svg`, `/icon/aws.svg`, `/icon/docker.svg`
- Emojis: `/emojis/1f3f3-fe0f.png`
- External: Any HTTPS URL

### Usage

```hcl
data "coder_parameter" "my_parameter" {
  icon = "/icon/coder.svg"

  option {
    icon = "https://example.com/icon.png"
  }
}

resource "coder_app" "my_app" {
  icon = "/icon/code.svg"
}
```

---

## Workspace Configuration

### Environment Variables

```hcl
resource "coder_agent" "main" {
  os   = "linux"
  arch = "amd64"

  env = {
    GIT_AUTHOR_NAME  = data.coder_workspace_owner.me.name
    GIT_AUTHOR_EMAIL = data.coder_workspace_owner.me.email
  }
}
```

### Startup Scripts

```hcl
resource "coder_agent" "main" {
  startup_script = <<-EOF
    #!/bin/bash
    set -e

    # Install dependencies
    pip install -r requirements.txt

    # Start services
    docker-compose up -d
  EOF
}
```

---

## Advanced Features

### Dev Container Support

```hcl
resource "coder_devcontainer" "main" {
  agent_id       = coder_agent.main.id
  workspace_folder = "/workspaces/project"
}
```

### Port Forwarding

```hcl
resource "coder_app" "api" {
  agent_id  = coder_agent.main.id
  slug      = "api"
  url       = "http://localhost:8000"
  subdomain = true
  share     = "authenticated"
}
```

### Metadata Display

```hcl
resource "coder_metadata" "info" {
  resource_id = docker_container.main.id
  hide        = false

  item {
    key   = "Container"
    value = docker_container.main.name
  }

  item {
    key       = "CPU Usage"
    value     = "Monitoring..."
    sensitive = false
  }
}
```
