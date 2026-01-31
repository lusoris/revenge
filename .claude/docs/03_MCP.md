# MCP (Model Context Protocol)

> Source: https://code.claude.com/docs/en/mcp
> Fetched: 2026-01-31
> Type: html

---

## What MCP Enables

Connect Claude to external tools and services:
- Implement features from issue trackers
- Analyze monitoring data (Sentry, Statsig)
- Query databases (PostgreSQL)
- Integrate designs (Figma, Slack)
- Automate workflows (Gmail drafts)

---

## Installing MCP Servers

### Option 1: Remote HTTP Server (Recommended)

```bash
claude mcp add --transport http <name> <url>

# Example: Connect to Notion
claude mcp add --transport http notion https://mcp.notion.com/mcp

# With Bearer token
claude mcp add --transport http secure-api https://api.example.com/mcp \
  --header "Authorization: Bearer your-token"
```

### Option 2: Remote SSE Server (Deprecated)

```bash
claude mcp add --transport sse <name> <url>
```

### Option 3: Local stdio Server

```bash
claude mcp add --transport stdio --env AIRTABLE_API_KEY=YOUR_KEY airtable \
  -- npx -y airtable-mcp-server
```

---

## Managing Servers

```bash
claude mcp list              # List all servers
claude mcp get github        # Get details
claude mcp remove github     # Remove server
/mcp                         # Check status in Claude Code
```

---

## MCP Scopes

| Scope | Storage | Purpose |
|-------|---------|---------|
| **Local** | `~/.claude.json` (per project) | Personal dev servers |
| **Project** | `.mcp.json` (checked in) | Team-shared tools |
| **User** | `~/.claude.json` | Cross-project utilities |

**Precedence**: Local > Project > User

---

## Practical Examples

### PostgreSQL Database
```bash
claude mcp add --transport stdio db -- npx -y @bytebase/dbhub \
  --dsn "postgresql://readonly:pass@prod.db.com:5432/analytics"
```

### Sentry Error Monitoring
```bash
claude mcp add --transport http sentry https://mcp.sentry.dev/mcp
# Then: /mcp to authenticate
```

### GitHub Integration
```bash
claude mcp add --transport http github https://api.githubcopilot.com/mcp/
```

---

## MCP Resources

Reference resources via @ mentions:
```
Analyze @github:issue://123
Review @postgres:schema://users
Fetch @docs:file://api/authentication
```

---

## MCP Tool Search

Automatically enables when tools exceed 10% of context:
- `ENABLE_TOOL_SEARCH=auto` (default, 10% threshold)
- `ENABLE_TOOL_SEARCH=auto:5` (custom 5% threshold)
- `ENABLE_TOOL_SEARCH=true` (always enabled)
- `ENABLE_TOOL_SEARCH=false` (disabled)

---

## Managed MCP Configuration

### Option 1: Exclusive Control (managed-mcp.json)

Deploy to system directories:
- macOS: `/Library/Application Support/ClaudeCode/managed-mcp.json`
- Linux: `/etc/claude-code/managed-mcp.json`
- Windows: `C:\Program Files\ClaudeCode\managed-mcp.json`

### Option 2: Policy-Based Control

```json
{
  "allowedMcpServers": [
    { "serverName": "github" },
    { "serverUrl": "https://mcp.company.com/*" }
  ],
  "deniedMcpServers": [
    { "serverName": "dangerous-server" }
  ]
}
```
