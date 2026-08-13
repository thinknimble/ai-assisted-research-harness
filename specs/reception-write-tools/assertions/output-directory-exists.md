---
id: output-directory-exists
parent: reception-write-tools
created: 2026-08-13T00:00:00Z
priority: 1
status: not_started
---

# Output Directory Exists in Research Repos

The research directory structure includes an `output/` directory for reception-generated files.

## Success Criteria

- `research-assistant init` creates `output/` alongside `raw/` and `formatted/`
- Existing research directories (those with `raw/` and `formatted/` already) get `output/` created on next init if it does not already exist
- The init success message includes `output/` in the directory listing
- `writeFile()` in `docs.go` works with `"output"` as the dir argument — no changes needed to its sandboxing logic since it already resolves `filepath.Join(projectRoot, dir, filename)`
