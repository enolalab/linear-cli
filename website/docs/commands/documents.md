# Documents

## `doc list`

List workspace documents, including document content, project, creator, and timestamps.

```bash
linear-cli doc list --limit 20 --cursor "$CURSOR"
```

| Flag | Default | Description |
| --- | ---: | --- |
| `--limit <n>` | 20 | Maximum result count |
| `--cursor <cursor>` | | Cursor from the previous response |

## `doc get <doc-id>`

Get one document and full content.

```bash
linear-cli doc get abc-123
```

## `doc search <query>`

Search documents by text. It accepts `--limit <n>` (default 20). Although its response includes pagination metadata, it has no `--cursor` flag, so this CLI cannot retrieve later search pages.

```bash
linear-cli doc search "onboarding" --limit 5
```
