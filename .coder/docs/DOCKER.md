# Docker in Coder Workspaces

> Source: https://coder.com/docs/admin/templates/extending-templates/docker-in-workspaces

## Overview

Multiple approaches exist for running Docker within Coder workspaces.

## Methods Comparison

| Method | Security | Performance | Requirements |
|--------|----------|-------------|--------------|
| **Sysbox** | High | Good | Sysbox runtime on nodes |
| **Envbox** | Medium | Good | Privileged outer container |
| **Rootless Podman** | High | Good | FUSE device support |
| **Privileged Sidecar** | Low | Good | None (insecure) |

## Sysbox Runtime

Secure docker-in-docker via custom runtime.

### Docker-Based Templates

```hcl
resource "docker_container" "workspace" {
  runtime = "sysbox-runc"
  # ...
}

resource "coder_agent" "main" {
  startup_script = <<EOF
#!/bin/sh
sudo dockerd &
EOF
}
```

### Kubernetes Templates

```hcl
resource "kubernetes_pod" "dev" {
  spec {
    runtime_class_name = "sysbox-runc"
    security_context {
      run_as_user = 1000
      fs_group    = 1000
    }
  }
}
```

## Envbox

Sysbox bundled in container image - no runtime installation needed.

**Pros:**
- No runtime installation on nodes
- Unlimited pod capacity
- Familiar Docker experience

**Cons:**
- Outer container requires privileged mode
- Slower initial startup
- Requires compatible kernel

Starter template: [kubernetes-envbox](https://github.com/coder/coder/tree/main/examples/templates/kubernetes-envbox)

## Rootless Podman

Docker alternative without privileges.

### Prerequisites

1. Enable smart-device-manager for FUSE
2. Label nodes: `smarter-device-manager=enabled`
3. Disable/set SELinux to permissive

### Bottlerocket Config

```hcl
[settings.kernel.sysctl]
"user.max_user_namespaces" = "65536"
```

## Privileged Sidecar (Insecure)

**Warning:** Workspaces can gain root access to host.

### Docker

```hcl
resource "docker_container" "dind" {
  image      = "docker:dind"
  privileged = true
  entrypoint = ["dockerd", "-H", "tcp://0.0.0.0:2375"]
}

resource "docker_container" "workspace" {
  env = ["DOCKER_HOST=${docker_container.dind.name}:2375"]
}
```

### Kubernetes

```hcl
spec {
  container {
    name  = "docker-sidecar"
    image = "docker:dind"
    security_context {
      privileged  = true
      run_as_user = 0
    }
    command = ["dockerd", "-H", "tcp://127.0.0.1:2375"]
  }

  container {
    name = "dev"
    env {
      name  = "DOCKER_HOST"
      value = "localhost:2375"
    }
  }
}
```

## Systemd Support

Combine Sysbox with systemd:

```hcl
resource "kubernetes_pod" "dev" {
  spec {
    runtime_class_name = "sysbox-runc"
    security_context {
      run_as_user = 0  # Required for systemd
      fs_group    = 0
    }
    container {
      command = ["sh", "-c", <<EOF
sudo -u coder --preserve-env=CODER_AGENT_TOKEN /bin/bash -- <<-'EOT' &
while [[ ! $(systemctl is-system-running) =~ ^(running|degraded) ]]; do
  sleep 2
done
${coder_agent.main.init_script}
EOT
exec /sbin/init
EOF
      ]
    }
  }
}
```

## Best Practices

1. Prefer Sysbox/Envbox over privileged sidecars
2. Verify kernel support before deploying
3. Monitor Docker resource consumption
4. Leverage Envbox image caching
5. Review [Sysbox limitations](https://github.com/nestybox/sysbox/blob/master/docs/user-guide/limitations.md)
