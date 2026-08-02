---
slug: /
---

# Linear CLI

`linear-cli` is a non-interactive command-line interface for Linear. It writes a consistent JSON response to standard output, making it suitable for scripts and AI agents that can execute shell commands.

```bash
export LINEAR_API_KEY=lin_api_your_key
linear-cli auth whoami
linear-cli issue list --team ENG --limit 5
```

Every command accepts `--pretty` for indented JSON. Commands that call Linear also accept `--api-key`; see [authentication](./authentication.md) for the exact credential order.

## Quick workflow

1. Install a release binary or use `go install`.
2. Set `LINEAR_API_KEY` in the process environment.
3. Discover workspace data with `team list`, `user list`, and `status list`.
4. Use IDs returned in JSON for resource-specific operations.

The CLI is a shell interface to the implemented command set. It does **not** publish a feature-parity or compatibility contract with Linear MCP.

## Next steps

- [Install the CLI](./installation.md)
- [Configure credentials](./authentication.md)
- [Understand JSON and exit codes](./output-and-errors.md)
- [Browse every command](./commands/index.md)
