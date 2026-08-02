# Labels

## `label list`

List issue labels, optionally scoped to a team. Responses include pagination metadata, but this command has no `--cursor` flag; see [pagination gaps](../pagination.md).

```bash
linear-cli label list --team ENG --limit 50
```

| Flag | Default | Description |
| --- | ---: | --- |
| `--team <key>` | | Team key to filter labels |
| `--limit <n>` | 50 | Maximum result count |

## `label create`

Create a workspace or team-scoped label.

```bash
linear-cli label create --name "Bug" --color "#ff0000" --team ENG
```

| Flag | Required | Description |
| --- | --- | --- |
| `--name <name>` | Yes | Label name |
| `--color <hex>` | No | Hex color such as `#ff0000` |
| `--team <key>` | No | Team key to scope the label |
