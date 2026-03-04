# incident.io MCP Server (Read-Only)

A **read-only** MCP server for incident.io that lets AI assistants (Claude, Windsurf, etc.) query incident data safely.

> ⚠️ **Read-Only**: This server can only read data - no create, update, or delete operations.

## � Installation

### Step 1: Get API Key
1. Go to incident.io → **Settings** → **API Keys**
2. Create a new API key (starts with `incio_`)

### Step 2: Install
```bash
# Clone the repo
git clone https://github.com/incident-io/incidentio-mcp-golang.git
cd incidentio-mcp-golang

# Set up environment
cp .env.example .env
# Edit .env and add: INCIDENT_IO_API_KEY=your-key-here

# Build
go build -o incidentio-mcp ./cmd/mcp-server

# Run
./incidentio-mcp
```

### Step 3: Connect to Windsurf/Claude

**Windsurf**: Add to `~/.codeium/windsurf/mcp_config.json`:
```json
{
  "incidentio-local": {
    "command": "/path/to/incidentio-mcp",
    "env": {
      "INCIDENT_IO_API_KEY": "your-key"
    }
  }
}
```

## 🛠️ Available Tools

### Incident Management (Read-Only)

- `list_incidents` - List incidents with optional filters (severity, status, date range, custom fields)
- `get_incident` - Get details of a specific incident by ID or reference (e.g., INC-123)
- `list_incident_statuses` - List all available incident statuses
- `list_incident_types` - List available incident types

### Alerts & Follow-ups
- `list_alerts` - List alerts with filters
- `get_alert` - Get alert details
- `list_follow_ups` - List follow-up tasks
- `get_follow_up` - Get follow-up details

### On-Call & Schedules
- `list_schedules` - List on-call schedules
- `get_schedule` - Get schedule details
- `get_current_on_call` - Get current on-call person

### Custom Fields & Catalogs
- `list_custom_fields` - List custom fields
- `search_custom_fields` - Search custom fields
- `list_catalog_types` - List catalog types
- `list_catalog_entries` - List catalog entries

### Severities & Statuses
- `list_severities` - List severity levels
- `list_incident_statuses` - List incident statuses

### Other Tools
- `list_users` - List organization users
- `list_workflows` - List workflows
- `list_actions` - List actions

## � Example Queries

Ask your AI assistant (Claude/Windsurf):

- "Show me all active incidents"
- "List P1 incidents from the last 30 days"
- "Who is on-call for the indexing team?"
- "Get details for incident INC-123"
- "Show me all outstanding follow-ups"
- "List alerts created in the last 24 hours"

## 🔧 Troubleshooting

| Issue | Solution |
|-------|----------|
| **API key error** | Verify key starts with `incio_` |
| **Go not found** | Install Go 1.21+ from [go.dev](https://go.dev/dl/) |
| **Build fails** | Run `go mod tidy` first |
| **404 errors** | Check incident ID is valid |

## 📚 More Info

- **30+ read-only tools** for incident.io
- **Complete API coverage** for incidents, alerts, schedules, custom fields
- **Safe for AI** - no create/update/delete operations
- **MIT License**

For detailed documentation, see the `docs/` folder.
