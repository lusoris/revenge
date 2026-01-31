# Coder Workspaces

> Source: https://coder.com/docs/user-guides

## Overview

Coder provides a comprehensive platform for managing cloud-based development environments. These guides contain information on workspace management, workspace access via IDEs, environment personalization, and workspace scheduling.

## Core Workspace Features

### Creation and Lifecycle

Users can create workspaces from templates and manage their complete lifecycle:

- **Creation**: Initialize workspaces using customizable templates
- **Status Management**: Monitor workspace states during provisioning and operation
- **Stopping/Deletion**: Control resource consumption through workspace lifecycle management

### Workspace Access Methods

Coder supports multiple IDE integrations for remote workspace access:

**Code Editors:**

- Visual Studio Code (desktop and browser)
- JetBrains IDEs (Fleet, Gateway, Toolbox)
- Cursor, Windsurf, Zed
- code-server for browser-based editing

**Terminal and File Access:**

- Web-based terminal interface
- Filebrowser for remote file management
- SSH and remote desktop (RDP) capabilities
- Emacs TRAMP support

**Port Management:**

- Port forwarding for accessing workspace services
- Wildcard access URLs for application hosting

### Workspace Management

Users can perform essential operations:

- Create and delete workspaces
- Rename workspace instances
- Restart and update workspaces
- Schedule automated start/stop times for cost optimization
- Set workspace favorites for quick access

## Personalization Features

### Environment Customization

- Dotfiles integration for applying personal configuration
- Environment variable management
- Custom application installations

### Dev Containers

Coder supports containerized development environments, allowing developers to run containerized development environments using the dev containers specification.

## Workspace Sharing (Beta)

The platform offers collaborative features, enabling users to share workspaces with team members for pair programming and knowledge transfer.

## Scheduling and Cost Control

Workspace scheduling enables cost control with workspace schedules, allowing administrators and users to define automatic start and stop times, reducing infrastructure costs.

## Administration vs. User Scope

Documentation explicitly distinguishes roles:

- **End users**: Workspace creation, management, and personalization
- **Administrators**: Template configuration and control plane management

## Integration Ecosystem

Coder provides integrations with development tools through:

- Native IDE plugins and extensions
- Backstage integration for platform engineering
- Container registry support
- Git operations integration
