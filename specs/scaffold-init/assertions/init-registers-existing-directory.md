---
id: init-registers-existing-directory
parent: scaffold-init
created: 2026-06-30T00:00:00Z
priority: 1
status: in_progress
locked-by: builder-Paris-MacBook-Pro-2.local-22524-1782850759
depends-on: init-creates-directory-structure
---

# Init Registers an Existing Research Directory Without Overwriting

When `research-assistant init` targets a directory that already contains `raw/`
and `formatted/`, it skips scaffolding and registers the directory in the global
config. This lets users adopt existing research directories they set up manually.

## Success Criteria

- If the target directory contains `raw/` and `formatted/`, init skips creating
  files and registers the directory in `~/.research-assistant/config.yaml`
- A message confirms registration: "Found existing research directory, registered
  as {name}"
- Missing pieces are filled in (e.g., if `doc-template.yaml` is absent, it is
  written; existing files are never overwritten)
- If the target directory has content but no `raw/` or `formatted/` (not a
  research directory), init warns and exits without changes
