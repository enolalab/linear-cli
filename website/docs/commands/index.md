# Command reference

All commands inherit these root flags:

| Flag | Description |
| --- | --- |
| `--api-key <key>` | Linear API key; highest authentication precedence |
| `--pretty` | Pretty-print the JSON envelope |
| `--help` | Show Cobra help |

`linear-cli version` returns build metadata: `version`, `commit`, and `date`.

| Command group | Operations |
| --- | --- |
| `auth` | `login`, `whoami` |
| `config` | `set`, `get`, `list` |
| `team` | `list`, `get` |
| `user` | `list`, `get`, `me` |
| `issue` | `list`, `get`, `create`, `update`, `search` |
| `comment` | `list`, `create` |
| `label` | `list`, `create` |
| `status` | `list`, `get` |
| `project` | `list`, `get`, `create`, `update` |
| `cycle` | `list`, `get` |
| `doc` | `list`, `get`, `search` |
| `attachment` | `upload` |

Cobra also provides `completion` and `help`; they are framework commands rather than Linear resource operations.
