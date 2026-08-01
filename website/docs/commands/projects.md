# Projects

## `project list`

List projects, optionally filtered by accessible team.

```bash
linear-cli project list --team ENG --limit 50 --cursor "$CURSOR"
```

| Flag | Default | Description |
| --- | ---: | --- |
| `--team <key>` | | Team key filter |
| `--limit <n>` | 50 | Maximum result count |
| `--cursor <cursor>` | | Cursor from the previous response |

## `project get <project-id>`

Return project details including lead, members, teams, and issue summaries.

```bash
linear-cli project get abc-123
```

## `project create`

`--name` is required. `--team` falls back to `default_team`; it is a usage error if neither is set. `--description-file` replaces `--description`.

```bash
linear-cli project create --name "Q2 Roadmap" --team ENG --description-file ./plan.md
```

| Flag | Required | Description |
| --- | --- | --- |
| `--name <name>` | Yes | Project name |
| `--team <key>` | No | Team key or `default_team` fallback |
| `--description <markdown>` | No | Project description |
| `--description-file <path>` | No | Read description from a file |

## `project update <project-id>`

Supply at least one update field.

```bash
linear-cli project update abc-123 --state completed
```

| Flag | Description |
| --- | --- |
| `--name <name>` | New project name |
| `--description <markdown>` | New description |
| `--state <state>` | `planned`, `started`, `paused`, `completed`, or `cancelled` |
