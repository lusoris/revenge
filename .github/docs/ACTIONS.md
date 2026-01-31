# GitHub Actions Workflow Syntax

> Source: https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions

## Overview

A workflow is a configurable automated process composed of one or more jobs. Workflows are defined using YAML files stored in `.github/workflows/`.

## Core Workflow Structure

### Basic Template

```yaml
name: Workflow Name
run-name: Custom Run Name

on:
  push:
    branches: [main]
  pull_request:

jobs:
  job-id:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "Hello World"
```

## Triggers (`on`)

### Event Types

```yaml
# Single event
on: push

# Multiple events
on: [push, pull_request]

# With activity types
on:
  pull_request:
    types: [opened, synchronize, reopened]
```

### Branch/Path Filtering

```yaml
on:
  push:
    branches:
      - main
      - 'releases/**'
    paths:
      - 'src/**'
      - '!src/docs/**'
```

### Schedule (Cron)

```yaml
on:
  schedule:
    - cron: '0 0 * * 0'  # Every Sunday at midnight UTC
```

Format: `minute hour day month day-of-week`

### Manual Triggers

```yaml
on:
  workflow_dispatch:
    inputs:
      environment:
        description: 'Target environment'
        required: true
        default: 'staging'
        type: choice
        options: [staging, production]
```

## Permissions

```yaml
permissions:
  contents: read
  pull-requests: write
  issues: write
```

Available: `read`, `write`, `none`, `read-all`, `write-all`

## Environment Variables

```yaml
env:
  NODE_ENV: production

jobs:
  build:
    env:
      CI: true
    steps:
      - run: echo ${{ env.NODE_ENV }}
        env:
          STEP_VAR: value
```

## Jobs

### Basic Job

```yaml
jobs:
  build:
    name: Build Application
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
```

### Job Dependencies

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps: [...]

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps: [...]
```

### Conditional Execution

```yaml
jobs:
  deploy:
    if: github.ref == 'refs/heads/main'
```

### Matrix Strategy

```yaml
jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest]
        node: [18, 20, 22]
        exclude:
          - os: windows-latest
            node: 18
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ matrix.node }}
```

### Concurrency

```yaml
jobs:
  deploy:
    concurrency:
      group: production
      cancel-in-progress: false
```

## Steps

### Using Actions

```yaml
steps:
  - uses: actions/checkout@v4
  - uses: actions/setup-node@v4
    with:
      node-version: '20'
```

### Running Commands

```yaml
steps:
  - run: npm ci
  - run: npm test
  - run: |
      echo "Multi-line"
      echo "commands"
```

### Step Options

```yaml
steps:
  - id: step-id
    name: Step Name
    if: success()
    continue-on-error: true
    timeout-minutes: 10
    working-directory: ./app
    run: npm test
```

## Expressions & Contexts

### Common Contexts

| Context | Description |
|---------|-------------|
| `github` | Workflow execution details |
| `env` | Environment variables |
| `vars` | Repository variables |
| `secrets` | Repository secrets |
| `inputs` | Workflow inputs |
| `jobs` | Job outputs |
| `steps` | Step outputs |

### Status Functions

| Function | Description |
|----------|-------------|
| `success()` | All previous steps succeeded |
| `failure()` | Any previous step failed |
| `always()` | Run regardless of status |
| `cancelled()` | Workflow was cancelled |

### Examples

```yaml
steps:
  - if: ${{ github.event_name == 'push' }}
    run: echo "Push event"

  - if: ${{ always() }}
    run: echo "Always runs"

  - if: ${{ contains(github.event.head_commit.message, '[skip ci]') }}
    run: echo "Skip CI"
```

## Outputs

### Step Outputs

```yaml
steps:
  - id: version
    run: echo "version=1.0.0" >> $GITHUB_OUTPUT

  - run: echo "Version is ${{ steps.version.outputs.version }}"
```

### Job Outputs

```yaml
jobs:
  build:
    outputs:
      version: ${{ steps.version.outputs.version }}
    steps:
      - id: version
        run: echo "version=1.0.0" >> $GITHUB_OUTPUT

  deploy:
    needs: build
    steps:
      - run: echo "${{ needs.build.outputs.version }}"
```

## Services

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: postgres
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
```

## Container Jobs

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    container:
      image: node:20
      env:
        NODE_ENV: test
      volumes:
        - my_docker_volume:/volume_mount
```

## Reusable Workflows

### Calling

```yaml
jobs:
  call-workflow:
    uses: ./.github/workflows/reusable.yml
    with:
      config-path: .github/config.json
    secrets: inherit
```

### Defining

```yaml
on:
  workflow_call:
    inputs:
      config-path:
        required: true
        type: string
    secrets:
      token:
        required: true
```

## Defaults

```yaml
defaults:
  run:
    shell: bash
    working-directory: ./app
```
