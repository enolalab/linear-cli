# Configuration

Configuration is stored at `~/.config/linear-cli/config.yaml`. The CLI creates the directory with owner-only permissions when saving config.

| Key | Purpose |
| --- | --- |
| `api_key` | Fallback Linear API key, after `--api-key` and `LINEAR_API_KEY` |
| `default_team` | Team key used only by `issue list`, `issue create`, and `project create` when those commands omit `--team` |

```bash
linear-cli config set default_team ENG
linear-cli config get default_team
linear-cli config list
```

`config list` redacts a saved API key. `config set` and `config get` accept any key string, but the CLI's documented configuration keys are `api_key` and `default_team`.

`default_team` does not apply to labels, cycles, workflow states, project listing, or any other command. Those commands require their own `--team` where documented.
