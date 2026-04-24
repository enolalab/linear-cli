# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Initial release with full Linear MCP feature parity
- **Authentication**: `auth login`, `auth whoami` commands
- **Configuration**: `config set/get/list` with default team support
- **Issues**: `issue list/get/create/update/search` with rich filtering
- **Comments**: `comment list/create` with file-based body input
- **Teams**: `team list/get` with member, state, and label details
- **Users**: `user list/get/me`
- **Workflow States**: `status list/get`
- **Labels**: `label list/create`
- **Projects**: `project list/get/create/update`
- **Cycles**: `cycle list/get` with issue details
- **Documents**: `doc list/get/search`
- **Attachments**: `attachment upload` with 3-step presigned URL flow
- JSON output by default with `--pretty` flag
- Structured error responses with 7 meaningful exit codes
- Human-readable identifier resolution (ENG-123 → UUID)
- Cursor-based pagination support
- Environment variable authentication (`LINEAR_API_KEY`)
- Cross-platform builds via GoReleaser (linux/darwin/windows × amd64/arm64)
- Homebrew tap and `go install` distribution
- CI/CD with GitHub Actions (test, lint, build, release)
