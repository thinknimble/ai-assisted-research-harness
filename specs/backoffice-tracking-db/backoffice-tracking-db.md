---
id: backoffice-tracking-db
created: 2026-08-26T00:00:00Z
priority: 1
---

# Backoffice Tracking Database

Backoffice mode uses fragile basename matching between `raw/` and `formatted/` to determine which files need processing. When the LLM names a formatted stub differently from the raw file's basename (common with subdirectory-organized raw files), the matching breaks — files reappear as "unprocessed" every run, and the LLM hallucinates user intent to justify overwriting existing stubs.

After this spec is complete, each research repo has a SQLite database (`.research-assistant.db`) that tracks which raw files have been processed and what formatted stubs they produced. The basename matching in `unprocessedRawFiles()` is replaced entirely by DB lookups.

## Context

- `backoffice.go:unprocessedRawFiles()` — current basename matching logic (lines 180–209)
- `backoffice.go:handleBackofficeTool("write_formatted_stub")` — where stubs are written (needs to record the mapping)
- `init.go` — creates the directory scaffold (needs to create/open the DB)
- `config.go` — global config lives at `~/.research-assistant/config.yaml`; the tracking DB is per-repo, not global
- Formatted stubs already contain a `path` frontmatter field pointing back to their raw file — useful for backfill migration
