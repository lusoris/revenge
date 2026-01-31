# Coder Web IDEs

> Source: https://coder.com/docs/admin/templates/extending-templates/web-ides

## Overview

Coder supports web-based IDEs through the `coder_app` resource, enabling browser access to development environments.

## code-server

VS Code in the browser:

```hcl
resource "coder_agent" "main" {
  arch = "amd64"
  os   = "linux"

  startup_script = <<EOF
#!/bin/sh
curl -fsSL https://code-server.dev/install.sh | sh -s -- --method=standalone --prefix=/tmp/code-server
/tmp/code-server/bin/code-server --auth none --port 13337 &
EOF
}

resource "coder_app" "code-server" {
  agent_id     = coder_agent.main.id
  slug         = "code-server"
  display_name = "code-server"
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

## VS Code Web

Native VS Code in browser (v1.82.0+):

```hcl
module "vscode-web" {
  source         = "registry.coder.com/modules/vscode-web/coder"
  version        = "1.0.14"
  agent_id       = coder_agent.main.id
  accept_license = true
}
```

Or manual setup:

```hcl
startup_script = <<EOF
#!/bin/sh
curl -Lk 'https://code.visualstudio.com/sha/download?build=stable&os=cli-alpine-x64' --output vscode_cli.tar.gz
mkdir -p /tmp/vscode-cli
tar -xf vscode_cli.tar.gz -C /tmp/vscode-cli
/tmp/vscode-cli/code serve-web --port 13338 --without-connection-token --accept-server-license-terms &
EOF

resource "coder_app" "vscode-web" {
  agent_id     = coder_agent.main.id
  slug         = "vscode-web"
  display_name = "VS Code Web"
  url          = "http://localhost:13338?folder=/home/coder"
  subdomain    = true
  share        = "owner"
}
```

## JupyterLab

```hcl
module "jupyter" {
  source   = "registry.coder.com/modules/jupyter-lab/coder"
  version  = "1.0.0"
  agent_id = coder_agent.main.id
}
```

Or manual:

```hcl
resource "coder_app" "jupyter" {
  agent_id     = coder_agent.main.id
  slug         = "jupyter"
  display_name = "JupyterLab"
  url          = "http://localhost:8888"
  icon         = "/icon/jupyter.svg"
  subdomain    = true
  share        = "owner"

  healthcheck {
    url       = "http://localhost:8888/healthz"
    interval  = 5
    threshold = 10
  }
}
```

## RStudio

```hcl
resource "coder_app" "rstudio" {
  agent_id     = coder_agent.main.id
  slug         = "rstudio"
  display_name = "RStudio"
  icon         = "https://upload.wikimedia.org/wikipedia/commons/d/d0/RStudio_logo_flat.svg"
  url          = "http://localhost:8787"
  subdomain    = true
  share        = "owner"

  healthcheck {
    url       = "http://localhost:8787/healthz"
    interval  = 3
    threshold = 10
  }
}
```

## File Browser

```hcl
module "filebrowser" {
  source   = "registry.coder.com/modules/filebrowser/coder"
  version  = "1.0.8"
  agent_id = coder_agent.main.id
}
```

## Custom Web Apps

Any web service can be exposed:

```hcl
resource "coder_app" "custom" {
  agent_id     = coder_agent.main.id
  slug         = "my-app"
  display_name = "My Application"
  url          = "http://localhost:3000"
  icon         = "/icon/custom.svg"
  subdomain    = true
  share        = "owner"  # owner, authenticated, public

  healthcheck {
    url       = "http://localhost:3000/health"
    interval  = 5
    threshold = 10
  }
}
```

## Share Options

| Value | Access |
|-------|--------|
| `owner` | Only workspace owner |
| `authenticated` | Any authenticated user |
| `public` | Anyone with URL |
