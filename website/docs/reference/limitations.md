# Limitations

## Cursor support

Cursor input is implemented only for `team list`, `user list`, `issue list`, `project list`, `cycle list`, and `doc list`. `issue search`, `doc search`, and `label list` expose Linear's page information in output but do not accept `--cursor`, so later pages are unavailable through those commands.

## Command scope

The documented surface is the command set in this repository. `linear-cli` is a stateless shell CLI and does not publish Linear MCP feature parity or a formal compatibility guarantee. Consumers needing a capability not listed in the command reference should not infer that it is available.

## Configuration scope

`default_team` is intentionally narrow: it applies to `issue list`, `issue create`, and `project create` only. It is not a global default for every command taking a team argument.

## Release channels

No Homebrew formula is configured in the current release setup. GitHub Release binaries and `go install` are the documented installation options.
