---
id: ignore-file-filters-listings
parent: research-ignore
created: 2026-08-18T17:00:00Z
priority: 1
status: not_started
---

# .researchignore Filters Directory Listings

`listDir()` in `docs.go` reads a `.researchignore` file from the project root and excludes matching files from its results.

## Success Criteria

- `.researchignore` lives at the research project root (same level as `raw/`, `formatted/`, `output/`)
- Patterns use gitignore-compatible syntax: globs (`*.tmp`), directory patterns (`raw/drafts/`), negation (`!keep-this.md`), comments (`# comment`), and `**` recursive matching
- Matching is applied against the file's path relative to the project root (e.g. `raw/drafts/notes.md`)
- Files matching any non-negated pattern are excluded from the slice `listDir()` returns
- If `.researchignore` does not exist or is empty, `listDir()` returns all files (backwards compatible)
- A malformed `.researchignore` (e.g. invalid pattern) does not crash — it is skipped with a warning to stderr
- Both `list_raw_files` (backoffice) and `list_stubs` (reception) respect the filtering because they call `listDir()`
- `unprocessedRawFiles()` inherits the filtering — ignored raw files are not queued for processing

**Note:** Use `github.com/sabhiram/go-gitignore` (or equivalent) for pattern matching rather than hand-rolling gitignore semantics. The library handles edge cases like trailing slashes for directory-only patterns and `**` globbing.
