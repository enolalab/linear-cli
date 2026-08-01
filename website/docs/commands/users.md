# Users

## `user list`

List workspace users with `id`, identity fields, activity and role flags, avatar URL, and creation time.

```bash
linear-cli user list --limit 50 --cursor "$CURSOR"
```

| Flag | Default | Description |
| --- | ---: | --- |
| `--limit <n>` | 50 | Maximum result count |
| `--cursor <cursor>` | | Cursor from the previous response |

## `user get <user-id>`

Get one user by Linear UUID.

```bash
linear-cli user get abc123-def456
```

## `user me`

Alias for `auth whoami`.

```bash
linear-cli user me
```
