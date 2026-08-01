# Authentication

Commands use a Linear API key. Credential resolution is exact and ordered:

1. `--api-key <key>`
2. `LINEAR_API_KEY`
3. `api_key` in `~/.config/linear-cli/config.yaml`

The command flag wins over the environment; the environment wins over the config file.

```bash
LINEAR_API_KEY=lin_api_environment linear-cli team list
linear-cli --api-key lin_api_override team list
linear-cli auth login --token lin_api_saved
```

`auth login` writes `api_key` to the local config file. It requires `--token`; it does not prompt. `auth whoami` confirms the authenticated user. `user me` is an alias for `auth whoami`.

## Keep tokens out of command history

Prefer `LINEAR_API_KEY` from a process environment, secret manager, or CI secret. Avoid putting `--api-key` and `auth login --token` values directly in checked-in scripts, logs, or process listings. See [AI agents](./automation/ai-agents.md) for a safe invocation pattern.
