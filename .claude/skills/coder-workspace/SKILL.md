---
name: coder-workspace
description: Manage Coder workspaces on coder.ancilla.lol
argument-hint: "<list|start|stop|logs> [workspace-name]"
disable-model-invocation: true
allowed-tools: Bash(coder *)
---

# Coder Workspace Management

Manage Coder workspaces on https://coder.ancilla.lol

## Usage

```
/coder-workspace list                 # List all workspaces
/coder-workspace start revenge        # Start a workspace
/coder-workspace stop revenge         # Stop a workspace
/coder-workspace logs revenge         # View workspace logs
/coder-workspace ssh revenge          # SSH into workspace
/coder-workspace delete revenge       # Delete a workspace
```

## Arguments

- `$0`: Action (list, start, stop, logs, ssh, delete, create)
- `$1`: Workspace name (required for all except list)

## Prerequisites

Ensure Coder CLI is configured:
```bash
coder login https://coder.ancilla.lol
```

## Task

### list: Show all workspaces
```bash
coder list
```

### start: Start a workspace
```bash
coder start $1
```

### stop: Stop a workspace
```bash
coder stop $1
```

### logs: View workspace logs
```bash
coder agent logs $1
```

### ssh: Connect to workspace
```bash
coder ssh $1
```

### delete: Delete a workspace
```bash
coder delete $1 --yes
```

### create: Create new workspace from template
```bash
coder create $1 --template revenge
```

## Common Operations

### Check workspace status
```bash
coder list --output json | jq '.[] | select(.name == "revenge")'
```

### View resource usage
```bash
coder stat cpu
coder stat mem
coder stat disk
```

### Port forward
```bash
coder port-forward $1 --tcp 8080:8080 --tcp 5173:5173
```

## Coder Host

All operations target: https://coder.ancilla.lol
