---
id: ci-release-pipeline
created: 2026-06-30T00:00:00Z
priority: 2
---

# CI Release Pipeline

The tool needs GitHub Actions workflows for testing on every push and building
cross-platform binaries on tagged releases. Follows the same pattern as spekk-cli
(tag-triggered release, cross-compiled binaries, GitHub Release with softprops)
but without dev releases or sandbox binaries.

## Context

- Entry point is the root package (`main.go`), not a `cmd/` subdirectory
- Binary name: `research-assistant`
- Module: `github.com/thinknimble/ai-assisted-research-harness`
- Target platforms: linux (amd64/arm64), darwin (amd64/arm64), windows (amd64/arm64)
- Reference workflow: spekk-cli `.github/workflows/publish.yml` and `test.yml`
