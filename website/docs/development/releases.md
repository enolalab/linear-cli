# Releases

Release tags matching `v*` trigger the GitHub release workflow, which invokes GoReleaser. The GoReleaser configuration builds archives for Linux, macOS, and Windows on `amd64` and `arm64`, and publishes GitHub Release assets.

```bash
git tag v0.1.0
git push --tags
```

## Installation channels

The configured release output offers GitHub Release binaries and `go install github.com/enolalab/linear-cli@<tag>`. There is no Homebrew tap configuration in `.goreleaser.yaml`; the repository removed that configuration. Do not announce or document a Homebrew formula until a tap and release automation are added.

## Go version discrepancy

The module declares Go `1.26.2`, while the current CI and release workflow request Go `1.23`. The module declaration is the documented development requirement; the automation version is stale and needs maintainers to align it before a release pipeline can be treated as matching the declared toolchain. This documentation does not change runtime or release configuration.
