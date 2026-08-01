# Pagination

Paginated responses contain `pagination.hasNextPage` and, when supplied by Linear, `pagination.endCursor`. Pass that cursor unchanged to the same command's `--cursor` flag.

```bash
first=$(linear-cli issue list --team ENG --limit 10)
cursor=$(printf '%s' "$first" | jq -r '.pagination.endCursor')
linear-cli issue list --team ENG --limit 10 --cursor "$cursor"
```

## Commands with `--cursor`

- `team list`
- `user list`
- `issue list`
- `project list`
- `cycle list`
- `doc list`

## Known cursor gaps

`issue search`, `doc search`, and `label list` return a pagination object from Linear but do **not** expose a `--cursor` flag. Their result sets cannot currently be advanced through this CLI. Treat `hasNextPage: true` on those commands as a known implementation limitation, not an invitation to pass an unsupported flag.
