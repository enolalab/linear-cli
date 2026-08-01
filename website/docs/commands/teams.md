# Teams

## `team list`

List accessible workspace teams. Results contain `id`, `name`, `key`, `description`, `private`, `timezone`, and `issueCount`.

```bash
linear-cli team list --limit 50 --cursor "$CURSOR"
```

| Flag | Default | Description |
| --- | ---: | --- |
| `--limit <n>` | 50 | Maximum result count |
| `--cursor <cursor>` | | Cursor from the previous response |

## `team get <team-key>`

Get a team by key such as `ENG`. The response includes members, workflow states, and labels in addition to team details.

```bash
linear-cli team get ENG
```
