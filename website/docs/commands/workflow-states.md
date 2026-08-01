# Workflow states

The CLI command is named `status`; it operates on Linear workflow states.

## `status list`

List states for a required team. State types include `backlog`, `unstarted`, `started`, `completed`, and `cancelled`.

```bash
linear-cli status list --team ENG
```

| Flag | Required | Description |
| --- | --- | --- |
| `--team <key>` | Yes | Team key |

## `status get <state-id>`

Get one workflow state by UUID.

```bash
linear-cli status get abc-123-def
```

Use the returned state `id` with `issue create --state` or `issue update --state`. `issue update --status` accepts a state name and resolves it for the issue's team.
