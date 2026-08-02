# Issues

Issue identifiers use `TEAM-NUMBER`, for example `ENG-123`. The CLI resolves them through Linear before reading or updating.

## `issue list`

```bash
linear-cli issue list --team ENG --assignee me --status "In Progress" --limit 20
```

| Flag | Default | Description |
| --- | ---: | --- |
| `--team <key>` | config fallback | Team key; falls back to `default_team` |
| `--assignee <id\|me>` | | User ID or `me` |
| `--status <name>` | | Workflow state name |
| `--priority <0-4>` | | 0 none, 1 urgent, 2 high, 3 medium, 4 low; applied only when supplied |
| `--label <name>` | | Label name |
| `--limit <n>` | 50 | Maximum result count |
| `--cursor <cursor>` | | Cursor from the previous response |

## `issue get <identifier>`

Return an issue with full markdown description and comments.

```bash
linear-cli issue get ENG-123
```

## `issue create`

`--title` is required. `--team` falls back to `default_team`; if neither is available, the command fails with a usage error. `--description-file` replaces `--description` when provided.

```bash
linear-cli issue create --title "Fix login" --team ENG --priority 1 --label bug
```

| Flag | Description |
| --- | --- |
| `--title <text>` | Required issue title |
| `--team <key>` | Team key or `default_team` fallback |
| `--description <markdown>` | Issue description |
| `--description-file <path>` | Read description from a file |
| `--priority <0-4>` | Priority value |
| `--assignee <user-id>` | Assignee UUID |
| `--label <name>` | Label name, resolved within the team |
| `--estimate <number>` | Point estimate |
| `--cycle <cycle-id>` | Cycle UUID |
| `--project <project-id>` | Project UUID |
| `--state <state-id>` | Workflow state UUID |

## `issue update <identifier>`

Supply at least one flag. `--description-file` replaces `--description` if both are set. `--status` resolves a state name within the issue team; `--state` directly assigns a state UUID and takes precedence if both are supplied.

```bash
linear-cli issue update ENG-123 --status "In Progress" --assignee user-id
```

| Flag | Description |
| --- | --- |
| `--title <text>` | New title |
| `--description <markdown>` | New description |
| `--description-file <path>` | Read new description from a file |
| `--priority <0-4>` | New priority |
| `--assignee <user-id>` | New assignee UUID |
| `--label <name>` | Label name to set; resolved within the team |
| `--estimate <number>` | New estimate |
| `--cycle <cycle-id>` | New cycle UUID |
| `--project <project-id>` | New project UUID |
| `--status <name>` | State name to resolve |
| `--state <state-id>` | State UUID, applied after `--status` |

## `issue search <query>`

Full-text search across issue title, description, and comments. It accepts `--limit <n>` (default 20) but no `--cursor`, even though its response reports pagination metadata.

```bash
linear-cli issue search "login failure" --limit 10
```
