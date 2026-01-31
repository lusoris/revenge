# Coder CLI Reference

> Source: https://coder.com/docs/reference/cli

## Global Flags

| Flag | Environment | Description |
|------|-------------|-------------|
| `--url` | `$CODER_URL` | Deployment URL |
| `--token` | `$CODER_SESSION_TOKEN` | Auth token |
| `-v, --verbose` | `$CODER_VERBOSE` | Verbose output |
| `--global-config` | `$CODER_CONFIG_DIR` | Config dir (~/.config/coderv2) |

## Authentication

```bash
# Login to Coder
coder login https://coder.example.com

# View current token
coder login token

# Logout
coder logout

# Check current user
coder whoami
```

## Workspace Management

```bash
# Create workspace
coder create my-workspace --template=python

# List workspaces
coder list

# Show workspace details
coder show my-workspace

# Start/Stop/Restart
coder start my-workspace
coder stop my-workspace
coder restart my-workspace

# Update to latest template
coder update my-workspace

# Delete workspace
coder delete my-workspace

# Rename workspace
coder rename old-name new-name

# Mark as favorite
coder favorite my-workspace
coder unfavorite my-workspace
```

## Workspace Access

```bash
# SSH into workspace
coder ssh my-workspace

# Run command via SSH
coder ssh my-workspace -- ls -la

# Open in browser
coder open my-workspace

# Open in VS Code
coder open vscode my-workspace

# Port forwarding
coder port-forward my-workspace 8080:8080 3000:3000

# Configure SSH config
coder config-ssh
```

## Scheduling

```bash
# Show schedule
coder schedule show my-workspace

# Set auto-start
coder schedule start my-workspace --schedule "0 9 * * MON-FRI"

# Set auto-stop
coder schedule stop my-workspace --stop-at "18:00"

# Extend running workspace
coder schedule extend my-workspace --duration 2h
```

## Templates

```bash
# Initialize from examples
coder templates init

# Create template
coder templates create my-template

# List templates
coder templates list

# Pull template code
coder templates pull my-template

# Push changes
coder templates push my-template

# Edit template
coder templates edit my-template

# List versions
coder templates versions list my-template

# Promote version
coder templates versions promote my-template@v1.2.3
```

## Monitoring

```bash
# View logs
coder logs my-workspace

# Test connectivity
coder ping my-workspace

# Speed test
coder speedtest my-workspace

# Resource stats
coder stat cpu my-workspace
coder stat mem my-workspace
coder stat disk my-workspace
```

## User Management

```bash
# List users
coder users list

# Create user
coder users create --email user@example.com

# Show user
coder users show username

# Suspend/Activate
coder users suspend username
coder users activate username

# Edit roles
coder users edit-roles username --roles owner
```

## Tokens

```bash
# Create token
coder tokens create --name my-token

# List tokens
coder tokens list

# Remove token
coder tokens remove my-token
```

## Dotfiles

```bash
# Apply dotfiles from repo
coder dotfiles https://github.com/user/dotfiles
```

## Server Admin

```bash
# Start server
coder server

# Create admin user
coder server create-admin-user

# Database encryption
coder server dbcrypt rotate
```

## Diagnostics

```bash
# Network check
coder netcheck

# Support bundle
coder support bundle

# Version info
coder version
```

## Shell Completion

```bash
# Bash
coder completion bash > /etc/bash_completion.d/coder

# Zsh
coder completion zsh > "${fpath[1]}/_coder"

# Fish
coder completion fish > ~/.config/fish/completions/coder.fish
```
