# VS Code Integration with Coder

> Source: https://coder.com/docs/user-guides/workspace-access/vscode

## Overview

Coder supports VS Code through desktop and browser-based approaches.

## VS Code Desktop

### One-Click Access

Click "VS Code Desktop" in dashboard to:
1. Install Coder Remote extension
2. Authenticate with Coder
3. Connect to workspace

### Manual Extension Install

```
ext install coder.coder-remote
```

Or download VSIX from [GitHub releases](https://github.com/coder/vscode-coder/releases).

## Extension Management

### Installation Methods

| Method | Use Case |
|--------|----------|
| Public marketplaces | code-server's interface |
| Custom Docker images | Bundle in workspace image |
| VSIX files | Command-line installation |
| Marketplace CLI | Remote installation |

### Adding to Custom Images

1. Download VSIX files to image
2. Install via startup script:

```hcl
resource "coder_agent" "main" {
  startup_script = "code-server --install-extension /vsix/extension.vsix"
}
```

### Marketplace Installation

```bash
SERVICE_URL=https://extensions.coder.com/api \
ITEM_URL=https://extensions.coder.com/item \
/path/to/code-server --install-extension <extension-name>
```

## Template Configuration

### code-server App

```hcl
resource "coder_app" "code-server" {
  agent_id     = coder_agent.main.id
  slug         = "code-server"
  display_name = "VS Code Browser"
  url          = "http://localhost:13337/?folder=/home/coder"
  icon         = "/icon/code.svg"
  subdomain    = false

  healthcheck {
    url       = "http://localhost:13337/healthz"
    interval  = 2
    threshold = 10
  }
}
```

### VS Code Desktop App

```hcl
resource "coder_app" "vscode-desktop" {
  agent_id     = coder_agent.main.id
  slug         = "vscode-desktop"
  display_name = "VS Code Desktop"
  url          = "vscode://coder.coder-remote/open?workspace=${data.coder_workspace.me.id}"
  icon         = "/icon/code.svg"
  external     = true
}
```

## Settings Sync

VS Code settings can sync across workspaces:

1. Sign in to VS Code with GitHub/Microsoft
2. Enable Settings Sync
3. Settings persist across workspace rebuilds

## Recommended Extensions

For Go development:

```json
{
  "customizations": {
    "vscode": {
      "extensions": [
        "golang.go",
        "ms-vscode.vscode-typescript-next",
        "esbenp.prettier-vscode",
        "dbaeumer.vscode-eslint",
        "bradlc.vscode-tailwindcss",
        "svelte.svelte-vscode"
      ],
      "settings": {
        "go.useLanguageServer": true,
        "editor.formatOnSave": true
      }
    }
  }
}
```

## Troubleshooting

### Connection Issues

1. Verify workspace is running: `coder list`
2. Check agent status: `coder show <workspace>`
3. Test SSH: `coder ssh <workspace>`
4. Check logs: `coder logs <workspace>`

### Extension Issues

1. Verify extension compatibility
2. Check extension logs in VS Code
3. Try reinstalling extension
4. Check workspace has required dependencies
