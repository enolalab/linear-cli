# Comments

## `comment list`

List an issue's comments. The issue must be given as a `TEAM-NUMBER` identifier.

```bash
linear-cli comment list --issue ENG-123
```

| Flag | Required | Description |
| --- | --- | --- |
| `--issue <identifier>` | Yes | Issue identifier, such as `ENG-123` |

## `comment create`

Create a markdown comment. A body read from `--body-file` replaces `--body`.

```bash
linear-cli comment create --issue ENG-123 --body "Fixed in PR #456"
```

| Flag | Required | Description |
| --- | --- | --- |
| `--issue <identifier>` | Yes | Issue identifier |
| `--body <markdown>` | One body source | Comment body |
| `--body-file <path>` | One body source | Read the comment body from a file |
