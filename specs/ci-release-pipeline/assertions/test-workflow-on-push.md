---
id: test-workflow-on-push
parent: ci-release-pipeline
created: 2026-06-30T00:00:00Z
priority: 1
status: not_started
---

# Test Workflow Runs on Every Push

A GitHub Actions workflow builds and tests on every push to any branch.

## Success Criteria

- `.github/workflows/test.yml` exists
- Triggers on push to all branches
- Uses `actions/setup-go` with `go-version-file: go.mod`
- Runs `go build .` (root package, not cmd/)
- Runs `go test ./... -count=1 -timeout 120s`
