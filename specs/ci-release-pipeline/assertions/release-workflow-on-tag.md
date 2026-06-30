---
id: release-workflow-on-tag
parent: ci-release-pipeline
created: 2026-06-30T00:00:00Z
priority: 1
status: done
depends-on: test-workflow-on-push
---

# Release Workflow Builds Cross-Platform Binaries on Tag

Pushing a `v*` tag triggers a release workflow that cross-compiles binaries for
all target platforms and creates a GitHub Release with the artifacts.

## Success Criteria

- `.github/workflows/publish.yml` exists
- Triggers on push of `v*` tags and on `workflow_dispatch`
- Has `permissions: contents: write`
- Runs tests before building
- Builds 6 binaries in `dist/`:
  - `research-assistant-linux-amd64`
  - `research-assistant-linux-arm64`
  - `research-assistant-darwin-amd64`
  - `research-assistant-darwin-arm64`
  - `research-assistant-windows-amd64.exe`
  - `research-assistant-windows-arm64.exe`
- Version is injected via `-ldflags "-X main.version=$VERSION"` from the tag
- Uses `softprops/action-gh-release@v3` to create the release with `dist/*`
  and `generate_release_notes: true`
- Build target is `.` (root package), not `./cmd/...`
