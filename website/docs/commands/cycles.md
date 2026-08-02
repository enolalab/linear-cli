# Cycles

## `cycle list`

List cycles (sprints) for a required team.

```bash
linear-cli cycle list --team ENG --limit 20 --cursor "$CURSOR"
```

| Flag | Default | Description |
| --- | ---: | --- |
| `--team <key>` | | Required team key |
| `--limit <n>` | 20 | Maximum result count |
| `--cursor <cursor>` | | Cursor from the previous response |

## `cycle get <cycle-id>`

Get a cycle and its issues.

```bash
linear-cli cycle get abc-123
```
