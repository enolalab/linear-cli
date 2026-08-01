# Authentication commands

## `auth login`

Save an API key as `api_key` in the local configuration file.

```bash
linear-cli auth login --token lin_api_your_key
```

| Flag | Required | Description |
| --- | --- | --- |
| `--token <key>` | Yes | Linear API key to save |

## `auth whoami`

Fetch the authenticated viewer. The returned user includes `id`, `name`, `email`, `displayName`, `active`, and `admin`.

```bash
linear-cli auth whoami
```

See [authentication](../authentication.md) for resolution precedence.
