# Contributing

## Toolchain requirement

`go.mod` declares Go `1.26.2`, which is the source of truth for local development and documentation. The current CI and release workflows still configure Go `1.23`, and CONTRIBUTING historically stated `1.23+`; those are repository configuration/documentation discrepancies, not an alternate supported toolchain declaration. Contributors should use Go `1.26.2` or a compatible later toolchain and should coordinate a workflow update before relying on CI for the newer module requirement.

## Local setup

```bash
git clone https://github.com/enolalab/linear-cli.git
cd linear-cli
go mod tidy
go test -v -race ./...
go build ./...
```

For docs changes:

```bash
npm install --prefix website
npm run build --prefix website
```

Command definitions in `cmd/*.go` are the source of truth for documentation. Keep `website/docs` in sync with their flags and behavior. Internal configuration, output, and error packages define credential resolution and the machine-readable contract.

## Design constraints

- Keep runtime commands non-interactive.
- Preserve the standard JSON envelope for command output.
- Return typed error categories rather than asking callers to parse prose.
- Add or update command documentation with an implementation change.
