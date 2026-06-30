---
id: template-embedded-in-binary
parent: scaffold-init
created: 2026-06-30T00:00:00Z
priority: 1
status: not_started
---

# Doc Template Is Embedded in the Binary

The `doc-template.yaml` content is embedded in the Go binary using `go:embed` so the
tool has zero external file dependencies. Users distribute and run a single binary.

## Success Criteria

- `doc-template.yaml` is embedded via Go's `embed` package
- `init` writes the embedded template into the scaffolded directory
- The backoffice mode's `read_file` tool can still read the template at runtime
  from the project directory (no behavior change for existing usage)
- No external files are required to run the binary after `init`
