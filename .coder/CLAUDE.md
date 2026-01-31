# Claude Code Instructions - Coder

**Tool**: Coder (Remote Development Platform)
**Purpose**: Remote workspace management for Revenge development
**Server**: https://coder.ancilla.lol
**Documentation**: [docs/INDEX.md](docs/INDEX.md)

---

## Entry Point for Claude Code

When working with Coder configuration for the Revenge project, always start by reading:

1. **Source of Truth**: [/docs/dev/design/00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md)
   - Technology stack
   - Container orchestration patterns
   - Deployment options

2. **Tech Stack**: [/docs/dev/design/technical/TECH_STACK.md](../docs/dev/design/technical/TECH_STACK.md)
   - Deployment platforms
   - Container requirements

3. **Coder Documentation**: [docs/INDEX.md](docs/INDEX.md)
   - Template reference
   - Workspace management
   - IDE integrations

---

## Coder Template Overview

**File**: `template.tf`
**Size**: 22KB
**Deployment Backends**: Docker, Kubernetes, Docker Swarm

### Template Features

- **Multi-backend support**: Choose deployment target
- **IDE integration**: VS Code (browser/desktop), Zed, JetBrains Gateway
- **Resource management**: CPU/memory limits configurable
- **Persistent storage**: Workspace data persists across rebuilds
- **Git integration**: Pre-configured Git credentials

### Parameters

See [docs/PARAMETERS.md](docs/PARAMETERS.md) for all configurable parameters.

---

## Common Tasks

### Modifying the Template

1. Edit `template.tf`
2. Test locally: `coder template plan`
3. Apply: `coder template push revenge`
4. Document changes in [docs/TEMPLATES.md](docs/TEMPLATES.md)

### Adding a New IDE Integration

1. Add IDE configuration to `template.tf`
2. Update [docs/WEB_IDES.md](docs/WEB_IDES.md) or [docs/VSCODE.md](docs/VSCODE.md)
3. Test IDE connection
4. Document setup process

### Changing Resource Limits

1. Modify resource parameters in `template.tf`
2. Consider: development vs production workloads
3. Test workspace creation
4. Document in [docs/PARAMETERS.md](docs/PARAMETERS.md)

---

## IDE Integrations

### VS Code

**Supported**:
- Browser (code-server)
- Desktop (Remote-SSH via Coder)

**See**: [docs/VSCODE.md](docs/VSCODE.md)

### Zed

**Supported**:
- Remote SSH

**See**: [docs/ZED_INTEGRATION.md](docs/ZED_INTEGRATION.md) (to be created)

### JetBrains Gateway

**Supported**:
- GoLand, IntelliJ IDEA

**See**: [docs/WEB_IDES.md](docs/WEB_IDES.md)

---

## Deployment Backends

### Docker (Local Development)

```bash
# Create workspace
coder create --template revenge my-workspace

# Start workspace
coder start my-workspace

# SSH into workspace
coder ssh my-workspace
```

### Kubernetes (Production)

```bash
# Apply template with K8s backend
coder template push revenge --variable backend=kubernetes

# Create workspace
coder create --template revenge prod-workspace
```

### Docker Swarm

```bash
# Apply template with Swarm backend
coder template push revenge --variable backend=swarm

# Create workspace
coder create --template revenge swarm-workspace
```

---

## Best Practices

1. **Test locally first** - Use Docker backend for template testing
2. **Version control** - Keep template.tf in Git
3. **Document parameters** - Update PARAMETERS.md when adding variables
4. **Resource limits** - Set appropriate CPU/memory based on workload
5. **Persistent storage** - Ensure critical data is in persistent volumes

---

## Troubleshooting

### Workspace won't start

1. Check Coder server status: `coder server status`
2. Check template: `coder template plan`
3. Review workspace logs: `coder logs my-workspace`
4. Verify backend availability (Docker/K8s/Swarm)

### IDE can't connect

1. Verify workspace is running: `coder ls`
2. Check SSH connection: `coder ssh my-workspace`
3. Verify IDE integration is configured in template
4. Check firewall/network settings

### Persistent data lost

1. Check volume configuration in template
2. Verify workspace wasn't deleted (stopped != deleted)
3. Check backend storage provider status

---

## Related Documentation

- **Coder Docs**: [docs/INDEX.md](docs/INDEX.md)
- **VS Code Integration**: [docs/VSCODE.md](docs/VSCODE.md)
- **Template Reference**: [docs/TEMPLATES.md](docs/TEMPLATES.md)
- **CLI Reference**: [docs/CLI.md](docs/CLI.md)
- **Main Project Docs**: [../docs/dev/design/INDEX.md](../docs/dev/design/INDEX.md)

---

## Quick Commands

```bash
# List templates
coder templates ls

# Push template
coder template push revenge

# Create workspace
coder create --template revenge my-workspace

# Start workspace
coder start my-workspace

# SSH into workspace
coder ssh my-workspace

# Stop workspace
coder stop my-workspace

# Delete workspace
coder delete my-workspace

# Open in VS Code
coder code my-workspace
```

---

**Last Updated**: 2026-01-31
**Maintained By**: Development Team
