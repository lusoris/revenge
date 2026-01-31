---
name: coder-template
description: Manage and deploy Coder templates
argument-hint: <push|pull|test|versions> [template-name]
disable-model-invocation: true
allowed-tools: Bash(coder *), Read, Glob
---

# Coder Template Management

Manage Coder templates for the Revenge development environment.

## Usage

```
/coder-template push                  # Push template to Coder
/coder-template pull revenge          # Pull existing template
/coder-template test                  # Validate template locally
/coder-template versions              # List template versions
/coder-template promote               # Promote template version
```

## Arguments

- `$0`: Action (push, pull, test, versions, promote)
- `$1`: Template name (default: revenge)

## Template Location

Template files are in `.coder/`:
- `template.tf` - Main Terraform configuration
- `docs/` - Template documentation

## Task

### push: Push template to Coder
```bash
cd .coder && coder templates push revenge --yes
```

### pull: Pull existing template
```bash
coder templates pull $1 --dest .coder/pulled
```

### test: Validate template
```bash
cd .coder && terraform init && terraform validate
```

### versions: List template versions
```bash
coder templates versions revenge
```

### promote: Promote to active
```bash
coder templates versions promote revenge --version latest
```

## Template Variables

The template supports these variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `deployment_backend` | docker, kubernetes, or swarm | docker |
| `namespace` | K8s namespace | coder-workspaces |
| `storage_class` | K8s storage class | standard |
| `cpu_limit` | CPU limit | 4 |
| `memory_limit` | Memory limit | 8Gi |
| `disk_size` | Disk size | 20Gi |
| `dotfiles_uri` | Dotfiles repo | (empty) |
| `git_clone_url` | Repo to clone | revenge |

## Deployment Backends

### Docker (default)
Simple single-node deployment using Docker containers.

### Kubernetes
Production-like environment with:
- Namespace isolation per user
- Persistent volumes
- Resource limits and requests
- StatefulSet for PostgreSQL

### Docker Swarm
Multi-node cluster deployment (requires Swarm mode).

## Creating Workspaces

After pushing template:
```bash
# Docker backend (default)
coder create my-workspace --template revenge

# Kubernetes backend
coder create my-workspace --template revenge \
  --parameter deployment_backend=kubernetes \
  --parameter namespace=dev-workspaces

# With dotfiles
coder create my-workspace --template revenge \
  --parameter dotfiles_uri=https://github.com/user/dotfiles
```

## Coder Host

Target: https://coder.ancilla.lol
