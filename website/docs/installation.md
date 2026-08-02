# Installation

## Release binary

Download the archive for your platform from [GitHub Releases](https://github.com/enolalab/linear-cli/releases), extract it, and place `linear-cli` on your `PATH`.

```bash
linear-cli version
```

## Go install

The module's declared Go version is `1.26.2`; use that version or a compatible newer toolchain when installing from source.

```bash
go install github.com/enolalab/linear-cli@latest
linear-cli version
```

## Homebrew status

The repository does not currently configure a Homebrew tap or generate a formula. The historical `brew install k0walski/tap/linear-cli` instruction is therefore not a supported installation path. Use a release binary or `go install` until a tap is configured and published.

## First command

Create a Linear personal API key in [Linear Settings > API](https://linear.app/settings/api), then keep it in the environment for the command invocation:

```bash
export LINEAR_API_KEY=lin_api_your_key
linear-cli auth whoami
```
