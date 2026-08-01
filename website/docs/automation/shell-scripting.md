# Shell scripting

Use compact JSON as the default transport between commands. `--pretty` is useful for inspection but not required for `jq`.

## Fail safely

```bash
#!/usr/bin/env bash
set -euo pipefail

: "${LINEAR_API_KEY:?Set LINEAR_API_KEY through your secret manager}"
issues=$(linear-cli issue list --team ENG --limit 20)
printf '%s' "$issues" | jq -r '.data[] | select(.priority == 1) | .identifier'
```

The CLI sends structured errors to standard output, so capture the response before deciding how to report it. Preserve the exit status; it has the categories documented in [output and errors](../output-and-errors.md).

## Iterate supported cursors

Only commands with a documented `--cursor` can advance a result set. This loop applies to `issue list`:

```bash
cursor=""
while :; do
  args=(issue list --team ENG --limit 50)
  if [ -n "$cursor" ]; then args+=(--cursor "$cursor"); fi
  page=$(linear-cli "${args[@]}")
  printf '%s\n' "$page" | jq -c '.data[]'
  [ "$(printf '%s' "$page" | jq -r '.pagination.hasNextPage')" = true ] || break
  cursor=$(printf '%s' "$page" | jq -r '.pagination.endCursor')
done
```

Do not use this pattern for `issue search`, `doc search`, or `label list`: they report pagination information but have no `--cursor` flag.

## File inputs

Use `--description-file` and `--body-file` for multiline markdown so shell quoting does not change the content. Use `attachment upload --file` for binary files; the command owns the presigned upload sequence.
