# Coder Template Parameters

> Source: https://coder.com/docs/admin/templates/extending-templates/parameters

## Overview

Parameters allow template authors to customize workspaces at build time by prompting users for configuration details.

## Parameter Types

| Type | Description |
|------|-------------|
| `string` | Text values |
| `bool` | Boolean true/false |
| `number` | Numeric values |
| `list(string)` | Arrays (use `jsonencode()` for defaults) |

## Required vs Optional

```hcl
# Required - no default
data "coder_parameter" "region" {
  name = "region"
  type = "string"
}

# Optional - has default
data "coder_parameter" "cpu" {
  name    = "cpu"
  type    = "number"
  default = 2
}
```

## Parameter Options

Restrict choices with options:

```hcl
data "coder_parameter" "region" {
  name    = "Region"
  type    = "string"
  default = "us-east-1"

  option {
    name  = "US East"
    value = "us-east-1"
    icon  = "/emojis/1f1fa-1f1f8.png"
  }

  option {
    name  = "EU West"
    value = "eu-west-1"
    icon  = "/emojis/1f1ea-1f1fa.png"
  }
}
```

## Mutability

Immutable parameters can only be set at creation or version update:

```hcl
data "coder_parameter" "disk_size" {
  name    = "disk_size"
  type    = "number"
  default = 20
  mutable = false  # Can't change after creation
}
```

## Validation

### Number Validation

```hcl
data "coder_parameter" "cpu" {
  name = "cpu"
  type = "number"

  validation {
    min       = 1
    max       = 8
    monotonic = "increasing"  # Can only increase
    error     = "CPU must be between {min} and {max}"
  }
}
```

### String Validation

```hcl
data "coder_parameter" "project_id" {
  name = "project_id"
  type = "string"

  validation {
    regex = "^[a-z0-9-]+$"
    error = "Project ID must be lowercase alphanumeric"
  }
}
```

## Ephemeral Parameters

Temporary parameters for workspace operations:

```hcl
data "coder_parameter" "force_rebuild" {
  name      = "force_rebuild"
  type      = "bool"
  default   = false
  ephemeral = true  # Only applies during start/update
}
```

## Workspace Presets

Bundle common parameter combinations:

```hcl
data "coder_workspace_preset" "gpu_dev" {
  name        = "GPU Development"
  description = "High-performance workspace with GPU"
  icon        = "/emojis/1f680.png"

  parameters = {
    "machine_type" = "n1-standard-4"
    "attach_gpu"   = "true"
    "disk_size"    = "100"
  }
}
```

## Dynamic Parameters

Coder v2.24.0+ supports conditional, identity-aware parameter forms.

## Parameter Autofill

- URL query: `?param.region=us-east-1`
- Recent values (requires `--experiments=auto-fill-parameters`)
