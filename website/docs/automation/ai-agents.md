# AI agents

`linear-cli` is designed for non-interactive shell execution: it does not request missing input, returns JSON, and uses categorized exit statuses. Agents should treat every invocation as an external side effect.

## Secure invocation pattern

Provide the key through a narrowly scoped environment variable supplied by the agent runtime or secret manager. Do not place a token in a prompt, argument list, repository file, or generated transcript.

```bash
env LINEAR_API_KEY="$LINEAR_API_KEY" linear-cli issue list --team ENG --limit 10
```

Avoid `set -x` around authenticated commands, redact captured command output before external logging, and do not use `--api-key` unless a controlled environment cannot supply `LINEAR_API_KEY`. Rotate a token exposed in a log or process listing.

## Read before write

Resolve identifiers and allowed values from Linear before making a mutation. For example, obtain workflow state IDs with `status list --team ENG`, then create or update an issue with a deliberate state value.

```bash
linear-cli status list --team ENG
linear-cli issue create --title "Triage incident" --team ENG --state "$STATE_ID"
```

Validate a mutation response before making a dependent mutation. Use `issue get` after a write when the workflow needs confirmation.

## Parse outcomes

Check the process status, then parse the JSON envelope. The exit code identifies the error category; `success` is the JSON equivalent. Retry only errors that make sense for the operation, such as a transient network failure or rate limit, with bounded backoff. Do not blindly retry writes after an ambiguous failure because the first request may have succeeded.

```bash
response=$(linear-cli issue search "login" --limit 5)
status=$?
if [ "$status" -ne 0 ]; then
  printf '%s\n' "$response" >&2
  exit "$status"
fi
printf '%s' "$response" | jq -r '.data[] | .identifier'
```

This CLI does not claim documented feature parity or a compatibility contract with Linear MCP. Agent integrations should use the command reference as the supported surface.
