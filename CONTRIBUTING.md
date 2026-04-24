# Contributing to linear-cli

Thank you for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- A [Linear](https://linear.app) account with an API key

### Getting Started

```bash
# Clone the repo
git clone https://github.com/enolalab/linear-cli.git
cd linear-cli

# Install dependencies
go mod tidy

# Build
make build

# Run
export LINEAR_API_KEY=lin_api_your_key_here
./linear-cli --help
```

## Project Structure

```
linear-cli/
├── main.go                 # Entry point
├── cmd/                    # Cobra command definitions
│   ├── root.go             # Root command + global flags
│   ├── auth.go             # auth login, whoami
│   ├── config.go           # config set/get/list
│   ├── team.go             # team list/get
│   ├── user.go             # user list/get/me
│   ├── issue.go            # issue list/get/create/update/search
│   ├── comment.go          # comment list/create
│   ├── label.go            # label list/create
│   ├── status.go           # status list/get (workflow states)
│   ├── project.go          # project list/get/create/update
│   ├── cycle.go            # cycle list/get
│   ├── doc.go              # doc list/get/search
│   └── attachment.go       # attachment upload
└── internal/
    ├── api/client.go       # GraphQL HTTP client
    ├── config/config.go    # Viper-based config management
    ├── errors/errors.go    # Typed errors + exit codes
    └── output/output.go    # JSON envelope formatter
```

## Design Principles

This CLI is designed for **AI Agent consumption**, not human users. Keep these in mind:

1. **JSON output always** — No tables, no colors, no emojis in default output
2. **No interactive prompts** — Never ask for user input; fail with a clear error
3. **Meaningful exit codes** — AI Agents use exit codes, not error messages
4. **Detailed help text** — AI Agents read `--help` to discover capabilities
5. **Predictable behavior** — Same input = same output, always

## Adding a New Command

1. Create `cmd/<resource>.go`
2. Define Cobra commands with comprehensive `Long` descriptions and `Example` fields
3. Use `getAPIClient()` to get the API client
4. Use `output.PrintSuccess()` / return `CLIError` for output
5. Register in `cmd/root.go` via `register<Resource>Commands()`
6. Add to command reference in `README.md`

### Command Template

```go
var myCmd = &cobra.Command{
    Use:   "list",
    Short: "Short description",
    Long: `Detailed description including output format.

Output JSON:
  { "success": true, "data": [...] }`,
    Example: `  linear-cli myresource list --flag value`,
    RunE: func(cmd *cobra.Command, args []string) error {
        client, err := getAPIClient()
        if err != nil {
            return err
        }

        // GraphQL query + execution...

        output.PrintSuccess(result, pagination)
        return nil
    },
}
```

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `ci:` CI/CD changes
- `chore:` Maintenance

## Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Run `make test` and `make build`
5. Commit with conventional commit message
6. Push and open a PR

## Releasing

Releases are automated via GitHub Actions + GoReleaser:

```bash
git tag v0.1.0
git push --tags
```

This triggers the release workflow which builds binaries for all platforms and updates the Homebrew tap.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
