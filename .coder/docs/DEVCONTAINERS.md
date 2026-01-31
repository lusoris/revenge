# Dev Containers in Coder

> Source: https://coder.com/docs/user-guides/devcontainers

## Overview

Dev containers define development environments as code using `devcontainer.json`. Coder integrates with `@devcontainers/cli` and Docker to build containerized environments.

## Prerequisites

- Coder version 2.24.0+
- Docker in workspace (Docker-in-Docker or socket mounting)
- `@devcontainers/cli` installed

## Key Features

- **Automatic Detection**: Discovers dev container configs from repositories
- **Seamless Startup**: Containers initialize automatically
- **Change Monitoring**: Dashboard shows outdated configs
- **On-Demand Rebuilds**: Rebuild via dashboard button
- **IDE Integration**: VS Code compatible
- **SSH Access**: Direct container connectivity
- **Port Discovery**: Automatic port forwarding

## Configuration

### File Locations

| Path | Description |
|------|-------------|
| `.devcontainer/devcontainer.json` | Recommended |
| `.devcontainer.json` | Repository root |
| `.devcontainer/<folder>/devcontainer.json` | Monorepo support |

### Minimal Configuration

```json
{
  "name": "My Dev Container",
  "image": "mcr.microsoft.com/devcontainers/base:ubuntu"
}
```

### With Features

```json
{
  "name": "Go Development",
  "image": "mcr.microsoft.com/devcontainers/go:1.22",
  "features": {
    "ghcr.io/devcontainers/features/docker-in-docker:2": {},
    "ghcr.io/devcontainers/features/node:1": {
      "version": "20"
    }
  },
  "customizations": {
    "vscode": {
      "extensions": [
        "golang.go",
        "esbenp.prettier-vscode"
      ]
    }
  },
  "forwardPorts": [8080, 3000]
}
```

## How It Works

1. Docker environment initializes
2. System detects repositories with dev container configs
3. Containers appear in dashboard
4. Auto-start if configured
5. Sub-agent created for each container

## Connecting

- **Web terminal**: Via Coder dashboard
- **SSH**: `coder ssh <workspace>.<agent>`
- **VS Code**: "Open in VS Code Desktop" button

## Agent Naming

Agent name derived from workspace folder:
- `/home/coder/my-app` → `my-app`
- Names sanitized to lowercase alphanumeric + hyphens

Custom names via `devcontainer.json`:

```json
{
  "customizations": {
    "coder": {
      "name": "custom-agent-name"
    }
  }
}
```

## Template Configuration

Enable dev containers in template:

```hcl
resource "coder_devcontainer" "main" {
  agent_id         = coder_agent.main.id
  workspace_folder = "/home/coder/project"
}
```

## Limitations

- **Linux only** (not Windows/macOS workspaces)
- Manual rebuild required for config changes
- `forwardPorts` with `host:port` not supported for sidecars
- Use `coder port-forward` for single-container workaround

## Best Practices

1. Use `.devcontainer/devcontainer.json` path
2. Test configs locally before deploying
3. Monitor dashboard for outdated status
4. Use Docker or Envbuilder based on infrastructure
5. Leverage dev container features for common tools
