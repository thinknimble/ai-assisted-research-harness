---
id: init-creates-directory-structure
parent: scaffold-init
created: 2026-06-30T00:00:00Z
priority: 1
status: done
---

# Init Creates Directory Structure

Running `research-assistant init [path]` creates a ready-to-use research directory
with all required folders and files.

## Success Criteria

- `research-assistant init ~/my-research` creates the directory at the given path
- If no path argument is given, scaffolds in the current working directory
- The created structure contains: `raw/`, `formatted/`, and `doc-template.yaml`
- If the target directory has content that is not a research directory (no `raw/`
  or `formatted/`), the command warns and exits without overwriting
- The command prints a success message with next steps when complete
  (e.g., "Run `research-assistant --mode backoffice` to start")
