# linear-cli

> CLI for [Linear](https://linear.app) — designed for AI Agents.

A command-line interface for the Linear project management tool that serves as a drop-in replacement for Linear's MCP (Model Context Protocol). Built primarily for AI Agents, it provides structured JSON output, meaningful exit codes, and zero interactive prompts.

## Why linear-cli?

| | Linear MCP | linear-cli |
|---|---|---|
| Connection | Persistent (stdio/SSE) | Stateless per-command |
| Protocol overhead | Handshake + negotiation | None |
| Framework support | MCP-compatible only | Any (shell exec) |
| Debugging | Complex | Run command in terminal |
| Composability | Limited | Pipes, xargs, scripts |
| Error handling | Protocol-level | Exit codes + JSON |

## Installation

### Homebrew (macOS / Linux)

```bash
brew install k0walski/tap/linear-cli
```

### Go install

```bash
go install github.com/enolalab/linear-cli@latest
```

### Binary download

Download pre-built binaries from [GitHub Releases](https://github.com/enolalab/linear-cli/releases).

## Quick Start

```bash
# 1. Set your API key
export LINEAR_API_KEY=lin_api_xxxxx

# 2. Verify authentication
linear-cli auth whoami

# 3. List your teams
linear-cli team list

# 4. List issues
linear-cli issue list --team ENG --limit 5

# 5. Get issue details
linear-cli issue get ENG-123
```

## Authentication

Set the `LINEAR_API_KEY` environment variable (recommended for AI Agents):

```bash
export LINEAR_API_KEY=lin_api_xxxxx
```

Or save it to the config file:

```bash
linear-cli auth login --token lin_api_xxxxx
```

You can get a personal API key from [Linear Settings → API](https://linear.app/settings/api).

**Priority:** `--api-key` flag > `LINEAR_API_KEY` env var > config file (`~/.config/linear-cli/config.yaml`).

## Output Format

All commands output a consistent JSON envelope:

**Success:**
```json
{
  "success": true,
  "data": { ... },
  "pagination": {
    "hasNextPage": true,
    "endCursor": "cursor-string"
  }
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "issue not found: ENG-999"
  }
}
```

Use `--pretty` for indented output when debugging.

## Exit Codes

| Code | Meaning | When |
|------|---------|------|
| 0 | Success | Operation completed |
| 1 | General error | Unspecified error |
| 2 | Usage error | Invalid arguments or flags |
| 3 | Auth error | Missing or invalid API key |
| 4 | Not found | Resource doesn't exist |
| 5 | Rate limited | API quota exceeded |
| 6 | Network error | Connection failure |
| 7 | Permission denied | Insufficient access |

## Command Reference

### Authentication & Config

| Command | Description |
|---------|-------------|
| `auth login --token <key>` | Save API key to config |
| `auth whoami` | Show current user |
| `config set <key> <value>` | Set config value |
| `config get <key>` | Get config value |
| `config list` | List all config (API key redacted) |

### Issues

| Command | Key Flags |
|---------|-----------|
| `issue list` | `--team`, `--assignee`, `--status`, `--priority`, `--label`, `--limit`, `--cursor` |
| `issue get <ID>` | Identifier like `ENG-123` |
| `issue create` | `--title` (req), `--team`, `--description`, `--description-file`, `--priority`, `--assignee`, `--label`, `--estimate`, `--cycle`, `--project` |
| `issue update <ID>` | `--title`, `--status`, `--priority`, `--assignee`, `--label`, `--description`, `--description-file` |
| `issue search <query>` | `--limit` |

### Comments

| Command | Key Flags |
|---------|-----------|
| `comment list` | `--issue` (required) |
| `comment create` | `--issue` (req), `--body` or `--body-file` |

### Teams & Users

| Command | Key Flags |
|---------|-----------|
| `team list` | `--limit`, `--cursor` |
| `team get <key>` | Team key like `ENG` |
| `user list` | `--limit`, `--cursor` |
| `user get <id>` | User UUID |
| `user me` | Alias for `auth whoami` |

### Workflow & Labels

| Command | Key Flags |
|---------|-----------|
| `status list` | `--team` (required) |
| `status get <id>` | State UUID |
| `label list` | `--team` |
| `label create` | `--name` (req), `--team`, `--color` |

### Projects & Cycles

| Command | Key Flags |
|---------|-----------|
| `project list` | `--team`, `--limit`, `--cursor` |
| `project get <id>` | Project UUID |
| `project create` | `--name` (req), `--team`, `--description`, `--description-file` |
| `project update <id>` | `--name`, `--description`, `--state` |
| `cycle list` | `--team` (required), `--limit`, `--cursor` |
| `cycle get <id>` | Includes issues |

### Documents

| Command | Key Flags |
|---------|-----------|
| `doc list` | `--limit`, `--cursor` |
| `doc get <id>` | Document UUID |
| `doc search <query>` | `--limit` |

### Attachments

| Command | Key Flags |
|---------|-----------|
| `attachment upload` | `--issue` (req), `--file` (req), `--title` |

## AI Agent Integration

### As Shell Tools

AI Agents can invoke `linear-cli` commands as shell tools. Each command is:
- **Stateless** — no persistent connection needed
- **Atomic** — one command = one operation
- **Predictable** — JSON output with typed exit codes

### Example: Agent Workflow

```
User: "Create 3 tasks for the current sprint in Engineering"

Agent executes:
1. linear-cli team list                          → find team ID
2. linear-cli cycle list --team ENG              → find current cycle
3. linear-cli status list --team ENG             → list workflow states
4. linear-cli issue create --title "Task 1" --team ENG
5. linear-cli issue create --title "Task 2" --team ENG
6. linear-cli issue create --title "Task 3" --team ENG
```

### Example: Filtering & Search

```bash
# Issues assigned to me that are in progress
linear-cli issue list --assignee me --status "In Progress"

# Search for bugs
linear-cli issue search "login error" --limit 10

# Get all workflow states to know valid status values
linear-cli status list --team ENG
```

### Pagination

For large result sets, use cursor-based pagination:

```bash
# First page
linear-cli issue list --team ENG --limit 10
# Response includes: "pagination": { "hasNextPage": true, "endCursor": "abc123" }

# Next page
linear-cli issue list --team ENG --limit 10 --cursor "abc123"
```

### Default Team

Set a default team to avoid passing `--team` every time:

```bash
linear-cli config set default_team ENG
linear-cli issue list  # automatically uses team ENG
```

## Configuration

Config file location: `~/.config/linear-cli/config.yaml`

| Key | Description | Example |
|-----|-------------|---------|
| `api_key` | Linear API key | `lin_api_xxxxx` |
| `default_team` | Default team key | `ENG` |

## Development

```bash
# Build
make build

# Run tests
make test

# Install locally
make install

# Clean
make clean
```

## License

[MIT](LICENSE)
